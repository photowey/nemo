/*
 * Copyright © 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package environment

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/photowey/nemo/internel/binder"
	"github.com/photowey/nemo/internel/eventbus"
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/filez"
	"github.com/photowey/nemo/pkg/mapz"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
)

type ActiveProfile string

func (p ActiveProfile) String() string {
	return string(p)
}

const (
	DevActiveProfile         = ActiveProfile("dev")
	TestActiveProfile        = ActiveProfile("test")
	ProdActiveProfile        = ActiveProfile("prod")
	StagingActiveProfile     = ActiveProfile("staging")
	IntegrationActiveProfile = ActiveProfile("integration")
	DemoActiveProfile        = ActiveProfile("demo")
	PreActiveProfile         = ActiveProfile("pre")
	TrainingActiveProfile    = ActiveProfile("training")
	BackupActiveProfile      = ActiveProfile("backup")
	DefaultActiveProfile     = ActiveProfile("default")
)

const (
	PrepareEnvironmentEventName  = "nemo.environment.prepare.event"
	PreLoadEnvironmentEventName  = "nemo.environment.load.pre.event"
	PostLoadEnvironmentEventName = "nemo.environment.load.post.event"
	ConfusedEnvironmentEventName = "nemo.environment.value.confused.event"
)

const (
	DefaultSystemPropertySourceName = "os.env"
	DefaultOptionPropertySourceName = "opt.properties"
	EvnSeparator                    = "="
	EvnValidLength                  = 2
)

const (
	AbsoluteFilePriority = ordered.HighPriority + 10*ordered.DefaultStep
	AbsolutePathPriority = ordered.HighPriority + 20*ordered.DefaultStep
	SearchPathPriority   = ordered.HighPriority + 30*ordered.DefaultStep
)

var (
	NoneSuccessThreshold   = SuccessThreshold(0) // default threshold
	AnyoneSuccessThreshold = SuccessThreshold(1)
	AllSuccessThreshold    = SuccessThreshold(2)
)

type SuccessThreshold int

func (st SuccessThreshold) Int() int {
	return int(st)
}

var (
	_ eventbus.Event = (*StandardEnvironmentEvent)(nil)
)

var (
	supportedConfigTypes = stringz.InitStringSlice("yaml", "yml", "toml", "properties")
	defaultConfigNames   = stringz.InitStringSlice("ini", "conf", "config", "configs", "application")
)

var (
	confusedEnvironments = stringz.InitStringSlice("::")
)

// ----------------------------------------------------------------

func RegisterSpecialEnvironment(key, value string) {
	_ = key
	confusedEnvironments = append(confusedEnvironments, value)
}

// ----------------------------------------------------------------

type StandardEnvironmentEvent struct {
	event string
	data  Environment
}

func NewStandardEnvironmentEvent(name string, data Environment) eventbus.Event {
	return &StandardEnvironmentEvent{
		event: name,
		data:  data,
	}
}

func (e *StandardEnvironmentEvent) Name() string {
	return e.event
}

func (e *StandardEnvironmentEvent) Topic() string {
	return e.event
}

func (e *StandardEnvironmentEvent) Data() any {
	return e.data
}

// ----------------------------------------------------------------

type Option func(opts *Options)

type Options struct {
	AbsolutePaths collection.StringSlice // /opt/data | /opt/configs | ...
	ConfigNames   collection.StringSlice // application | config | configs | ...
	ConfigTypes   collection.StringSlice // yaml/yml | toml -> but, json not support now.
	SearchPaths   collection.StringSlice // ./resources | ./configs ...
	Profiles      collection.StringSlice // dev | test | prod | ...
	Sources       PropertySources        // PropertySource
	Properties    collection.MixedMap    // Properties -> map data-structure -> can also be replaced by PropertySource
	Threshold     SuccessThreshold       // the allow count of config files successfully loaded
}

func (opts *Options) validate() (err error) {
	if err = validateAbsolutePaths(opts.AbsolutePaths); err != nil {
		return
	}
	if err = validateConfigNames(opts.ConfigNames); err != nil {
		return
	}
	if err = validateConfigTypes(opts.ConfigTypes); err != nil {
		return
	}
	if err = validateSearchPaths(opts.SearchPaths); err != nil {
		return
	}
	if err = validateProfiles(opts.Profiles); err != nil {
		return

	}
	if err = validateSources(opts.Sources); err != nil {
		return
	}
	if err = validateProperties(opts.Properties); err != nil {
		return
	}
	if err = validateThreshold(opts.Threshold); err != nil {
		return
	}

	return
}

func validateAbsolutePaths(absolutePaths collection.StringSlice) error {
	for _, absolutePath := range absolutePaths {
		if stringz.IsNotBlankString(absolutePath) {
			if ok := path.IsAbs(absolutePath); !ok {
				return fmt.Errorf("nemo: the candidate path:[%s] is not absolute path", absolutePath)
			}
		}
	}

	return nil
}

func validateConfigNames(configNames collection.StringSlice) error {
	return nil
}

func validateConfigTypes(configTypes collection.StringSlice) error {
	for _, configType := range configTypes {
		if collection.ArrayNotContains(supportedConfigTypes, configType) {
			return fmt.Errorf("nemo: config type:[%s] not supported", configType)
		}
	}

	return nil
}

func validateSearchPaths(searchPaths collection.StringSlice) error {
	return nil
}

func validateProfiles(profiles collection.StringSlice) error {
	return nil
}

func validateSources(sources PropertySources) error {
	for _, source := range sources {
		if source.IsEmptySource() {
			return fmt.Errorf("nemo: the property source can't be empty")
		}
		if source.IsMapSource() {
			if len(source.Map) == 0 {
				return fmt.Errorf("nemo: the `Map` of map property source can't be empty")
			}
		}
	}

	return nil
}

func validateProperties(properties collection.MixedMap) error {
	return nil
}

func validateThreshold(threshold SuccessThreshold) error {
	th := threshold.Int()
	if th < NoneSuccessThreshold.Int() || th > AllSuccessThreshold.Int() {
		return fmt.Errorf("nemo: invalid threshold:[%d]", th)
	}

	return nil
}

// ----------------------------------------------------------------

type PropertySources = []PropertySource

type PropertySource struct {
	ordered.Ordered                     // declare ordered
	Priority        int64               // the of priority of the PropertySource.
	Property        string              // the name of the PropertySource.
	FilePath        string              // the path of config file. -> e.g.: /opt/data | /opt/configs
	Name            string              // the name of config file. -> e.g.: config.yaml | config.yml config.toml
	Suffix          string              // the suffix of config file -> Name == config Suffix == yml -> Full name == config.yml
	Type            reflect.Type        // the type of PropertySource, only support map now. // or string ?
	Map             collection.MixedMap // the map context, when the Type is map.
}

// Order | the priority value of PropertySource sort.
//
// The smaller the value, the higher the priority.
func (ps PropertySource) Order() int64 {
	return ps.Priority
}

func (ps PropertySource) IsMapSource() bool {
	if ps.Type != nil && ps.Type.Kind() == reflect.Map {
		return true
	}

	return false
}

func (ps PropertySource) IsFileSource() bool {
	if ps.IsMapSource() {
		return false
	}

	if stringz.IsBlankString(ps.FilePath) &&
		stringz.IsBlankString(ps.Name) &&
		stringz.IsBlankString(ps.Suffix) {
		return false
	}

	return true
}

func (ps PropertySource) IsEmptySource() bool {
	if ps.IsMapSource() || ps.IsFileSource() {
		return false
	}

	return true
}

// ----------------------------------------------------------------

type Environment interface {
	Start(opts ...Option) error
	Destroy() error
	Refresh(opts ...Option) error
	LoadMap(sourceMap collection.MixedMap) error
	LoadPropertySources(sources ...PropertySource) error
	Get(key string) (any, bool)
	NestedGet(key string) (any, bool)
	Set(key string, value any)
	NestedSet(key string, value any)
	Contains(key string) bool
	ActiveProfiles() collection.StringSlice
	ActiveProfilesString() string
	ActiveDefaultProfile() bool
	Bind(prefix string, target any) error
}

// ----------------------------------------------------------------

// errors 抽象设计

// ----------------------------------------------------------------

type StandardEnvironment struct {
	configMap       collection.MixedMap    // core config container
	propertySources []PropertySource       // config sources
	profiles        collection.StringSlice // Profiles active e.g.: dev test prod ...
	threshold       SuccessThreshold       // threshold
	binder          *binder.Binder         // default binder
}

// ----------------------------------------------------------------

func New(sources ...PropertySource) Environment {
	return &StandardEnvironment{
		configMap:       make(collection.MixedMap),
		propertySources: sources,
		profiles:        make(collection.StringSlice, 0),
		threshold:       NoneSuccessThreshold, // default threshold
		binder:          binder.New(),
	}
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) Start(opts ...Option) error {
	optz, err := initOptions(opts...)
	if err != nil {
		return err
	}

	err = e.translateToPropertySources(optz)
	if err != nil {
		return err
	}

	// prepare
	eventPrepare := NewStandardEnvironmentEvent(PrepareEnvironmentEventName, e)
	if err = eventbus.Post(eventPrepare); err != nil {
		return nil
	}

	// pre load
	preLoadEvent := NewStandardEnvironmentEvent(PreLoadEnvironmentEventName, e)
	if err = eventbus.Post(preLoadEvent); err != nil {
		return nil
	}

	// on load
	if err = e.onLoad(); err != nil {
		return nil
	}

	// post load
	postLoadEvent := NewStandardEnvironmentEvent(PostLoadEnvironmentEventName, e)
	if err = eventbus.Post(postLoadEvent); err != nil {
		return nil
	}

	return nil
}

func (e *StandardEnvironment) Destroy() error {
	return nil
}

func (e *StandardEnvironment) Refresh(opts ...Option) error {
	// Destroy ...

	return e.Start(opts...)
}

func (e *StandardEnvironment) LoadMap(sourceMap collection.MixedMap) error {
	e.mergeMap(sourceMap)

	return nil
}

func (e *StandardEnvironment) LoadPropertySources(sources ...PropertySource) error {
	for _, source := range sources {
		if stringz.IsNotBlankString(source.FilePath) {
			if err := e.loadConfig(source.FilePath, source.Name, source.Type); err != nil {
				return err
			}
		}

		if source.IsMapSource() {
			if err := e.LoadMap(source.Map); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *StandardEnvironment) Get(key string) (any, bool) {
	return e.getProperty(key)
}

func (e *StandardEnvironment) NestedGet(key string) (any, bool) {
	return e.getProperty(key)
}

func (e *StandardEnvironment) Set(key string, value any) {
	e.setProperty(key, value)
}

func (e *StandardEnvironment) NestedSet(key string, value any) {
	e.setProperty(key, value)
}

func (e *StandardEnvironment) Contains(key string) bool {
	return mapz.NestedContains(key, e.configMap)
}

func (e *StandardEnvironment) ActiveProfiles() collection.StringSlice {
	return e.profiles
}

func (e *StandardEnvironment) ActiveProfilesString() string {
	return stringz.Implode(e.profiles, stringz.SymbolComma)
}

func (e *StandardEnvironment) ActiveDefaultProfile() bool {
	return collection.ArrayContains(e.profiles, DefaultActiveProfile.String())
}

func (e *StandardEnvironment) Bind(prefix string, target any) error {
	e.binder.Bind(prefix, target, e.configMap)

	return nil
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) setProperty(key string, value any) {
	mapz.NestedSet(e.configMap, key, value)
}

func (e *StandardEnvironment) getProperty(key string) (any, bool) {
	return mapz.NestedGet(e.configMap, key)
}

func (e *StandardEnvironment) mergeMap(ctx collection.MixedMap) {
	if collection.IsEmptyMap(ctx) {
		return
	}

	mapz.MergeMixedMaps(e.configMap, ctx)
}

func (e *StandardEnvironment) translateToPropertySources(opts *Options) error {
	e.translateSources(opts)
	e.translateProfiles(opts)
	e.translateProperties(opts)
	err := e.translatePaths(opts)
	if err != nil {
		return err
	}

	return nil
}

func (e *StandardEnvironment) translateSources(opts *Options) {
	if collection.IsNotEmptySlice(opts.Sources) {
		e.propertySources = append(e.propertySources, opts.Sources...)
	}
}

func (e *StandardEnvironment) translateProfiles(opts *Options) {
	if collection.IsNotEmptySlice(opts.Profiles) {
		e.profiles = append(e.profiles, opts.Profiles...)
	}

	// default profile active.
	if collection.IsEmptySlice(e.profiles) {
		e.profiles = append(e.profiles, DefaultActiveProfile.String())
	}
}

func (e *StandardEnvironment) translateProperties(opts *Options) {
	ps := initPropertySource(opts.Properties, ordered.DefaultPriority, DefaultOptionPropertySourceName)

	e.propertySources = append(e.propertySources, ps)
}

func (e *StandardEnvironment) translatePaths(opts *Options) error {
	absolutePaths := opts.AbsolutePaths
	for _, absolutePath := range absolutePaths {
		e.translatePathToPropertySourceIfNecessary(AbsoluteFilePriority, AbsolutePathPriority, absolutePath, opts)
	}

	searchPaths := opts.SearchPaths
	if collection.IsNotEmptySlice(searchPaths) {
		for _, searchPath := range searchPaths {
			abs, err := filez.ToAbsIfNecessary(searchPath)
			if err != nil {
				return nil
			}
			e.translatePathToPropertySourceIfNecessary(ordered.DefaultPriority, SearchPathPriority, abs, opts)
		}
	}

	return nil
}

func (e *StandardEnvironment) translatePathToPropertySourceIfNecessary(filePriority, priority int64, abs string, opts *Options) {
	if stringz.IsBlankString(abs) {
		return
	}

	abs = filepath.Clean(abs)

	if filez.IsFile(abs) {
		e.translateFileToPropertySource(filePriority, abs)
		return
	}

	configNames := opts.ConfigNames
	if collection.IsEmptySlice(configNames) {
		for _, configName := range defaultConfigNames {
			if collection.ArrayNotContains(configNames, configName) {
				configNames = append(configNames, configName)
			}
		}
	}

	configTypes := opts.ConfigTypes
	if collection.IsEmptySlice(configTypes) {
		// custom ?
		for _, supportedConfigType := range supportedConfigTypes {
			if collection.ArrayNotContains(configTypes, supportedConfigType) {
				configTypes = append(configTypes, supportedConfigType)
			}
		}
	}

	for _, configName := range configNames {
		file := filepath.Join(abs, configName)
		file = filepath.Clean(file)

		if filez.IsFile(file) {
			e.translateFileToPropertySource(filePriority, file)

			continue
		}

		for _, configType := range configTypes {
			ps := PropertySource{
				Priority: priority,
				Property: abs,
				FilePath: abs,
				Name:     configName,
				Suffix:   configType,
			}

			e.propertySources = append(e.propertySources, ps)
		}
	}
}

func (e *StandardEnvironment) translateFileToPropertySource(priority int64, abs string) {
	abs = filepath.Clean(abs)

	dir := filepath.Dir(abs)
	fileName := filepath.Base(abs)
	ext := filepath.Ext(fileName)

	ps := PropertySource{
		Priority: priority,
		Property: abs,
		FilePath: dir,
		Name:     fileName,
		Suffix:   ext,
	}

	e.propertySources = append(e.propertySources, ps)
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) onLoad() error {
	e.loadSystemEnvVars()

	sorter := ordered.NewSorter(e.propertySources...)
	ordered.Sort(sorter, -1)

	th := e.threshold

	okCounter := 0
	errs := make([]error, 0)

	for _, source := range e.propertySources {
		if err := e.LoadPropertySources(source); err != nil {
			if AllSuccessThreshold.Int() == th.Int() {
				return err
			}
			errs = append(errs, err)

			continue
		}

		okCounter++
	}

	if AnyoneSuccessThreshold.Int() == th.Int() && okCounter == 0 {
		sb := stringz.NewStringBuffer(len(errs))
		for _, err := range errs {
			sb.Append(err.Error())
		}

		return fmt.Errorf("nemo: failed to load the all property sources, messages:[%s]", sb.String())
	}

	return nil
}

func (e *StandardEnvironment) loadSystemEnvVars() {
	envVars := make(collection.MixedMap)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, EvnSeparator, EvnValidLength) // k=v -> 2
		// warning:
		// -> ${env} == [ =::=::\ ] ?
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]

			postParse(env)
		}
	}

	e.loadSystemEnvMap(envVars)
}

func postParse(env string) {
	for _, confusedEnv := range confusedEnvironments {
		if strings.Contains(env, confusedEnv) {
			// post: confused event?
			confusedEvent := eventbus.NewStandardAnyEvent(ConfusedEnvironmentEventName, env)
			_ = eventbus.Post(confusedEvent)
		}
	}

	if collection.ArrayContains(confusedEnvironments, env) {
		confusedEvent := eventbus.NewStandardAnyEvent(ConfusedEnvironmentEventName, env)
		_ = eventbus.Post(confusedEvent)
	}
}

func initSystemEnvPropertySource(envVars collection.MixedMap) PropertySource {
	return initPropertySource(envVars, ordered.HighPriority, DefaultSystemPropertySourceName)
}
func initPropertySource(ctx collection.MixedMap, priority int64, property string) PropertySource {
	return PropertySource{
		Priority: priority,
		Property: property,
		Type:     reflect.TypeOf(collection.MixedMap{}),
		Map:      ctx,
	}
}

func (e *StandardEnvironment) loadSystemEnvMap(envVars collection.MixedMap) {
	e.loadSystemEnvMapDelayed(envVars)
}

func (e *StandardEnvironment) loadSystemEnvMapDelayed(envVars collection.MixedMap) {
	envPs := initSystemEnvPropertySource(envVars)
	e.propertySources = append(e.propertySources, envPs)
}

func (e *StandardEnvironment) loadConfig(path, name string, _ reflect.Type) error {
	return nil
}

// ----------------------------------------------------------------

func WithAbsolutePaths(absolutePaths ...string) Option {
	return func(opts *Options) {
		opts.AbsolutePaths = append(opts.AbsolutePaths, absolutePaths...)
	}
}

func WithConfigNames(configNames ...string) Option {
	return func(opts *Options) {
		opts.ConfigNames = append(opts.ConfigNames, configNames...)
	}
}

func WithConfigTypes(configTypes ...string) Option {
	return func(opts *Options) {
		opts.ConfigTypes = append(opts.ConfigTypes, configTypes...)
	}
}

func WithSearchPaths(searchPaths ...string) Option {
	return func(opts *Options) {
		opts.SearchPaths = append(opts.SearchPaths, searchPaths...)
	}
}

func WithProfiles(profiles ...string) Option {
	return func(opts *Options) {
		opts.Profiles = append(opts.Profiles, profiles...)
	}
}

func WithSources(sources ...PropertySource) Option {
	return func(opts *Options) {
		opts.Sources = append(opts.Sources, sources...)
	}
}

func WithProperties(properties collection.MixedMap) Option {
	return func(opts *Options) {
		opts.Properties = properties
	}
}

// ----------------------------------------------------------------

func initOptions(opts ...Option) (*Options, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	if err := populateIfNecessary(options); err != nil {
		return nil, err
	}

	return options, nil
}

func populateIfNecessary(opts *Options) error {
	return opts.validate()
}

func newOptions() *Options {
	return &Options{
		AbsolutePaths: make(collection.StringSlice, 0),
		ConfigNames:   make(collection.StringSlice, 0),
		ConfigTypes:   make(collection.StringSlice, 0),
		SearchPaths:   make(collection.StringSlice, 0),
		Profiles:      make(collection.StringSlice, 0),
		Sources:       make(PropertySources, 0),
		Properties:    make(collection.MixedMap),
	}
}
