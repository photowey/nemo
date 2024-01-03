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
	"os"
	"reflect"
	"strings"

	"github.com/photowey/nemo/internel/eventbus"
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
)

const (
	PrepareEnvironmentEventName  = "nemo.environment.prepare.event"
	PreLoadEnvironmentEventName  = "nemo.environment.load.pre.event"
	PostLoadEnvironmentEventName = "nemo.environment.load.post.event"
	ConfusedEnvironmentEventName = "nemo.environment.value.confused.event"
)

const (
	DefaultSystemPropertySourceName = "os.env"
	EvnSeparator                    = "="
	EvnValidLength                  = 2
)

var (
	_ eventbus.Event = (*StandardEnvironmentEvent)(nil)
)

var (
	confusedEnvironments = collection.StringSlice{
		"::",
	}
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
	Properties    collection.AnyMap      // Properties -> map data-structure -> can also be replaced by PropertySource
	ThrowLevel    int                    // 0: all 1:anyone
}

func (opts *Options) validate() {
	// validate options?
}

// ----------------------------------------------------------------

type PropertySources = []PropertySource

type PropertySource struct {
	ordered.Ordered                   // declare ordered
	Priority        int64             // the of priority of the PropertySource.
	Property        string            // the name of the PropertySource.
	FilePath        string            // the path of config file. -> e.g.: /opt/data | /opt/configs
	Name            string            // the name of config file. -> e.g.: config.yaml | config.yml config.toml
	Suffix          string            // the suffix of config file -> Name == config Suffix == yml -> Full name == config.yml
	Type            reflect.Type      // the type of PropertySource, only support map now. // or string ?
	Map             collection.AnyMap // the map context, when the Type is map.
}

// Order | the priority value of PropertySource sort.
//
// The smaller the value, the higher the priority.
func (ps PropertySource) Order() int64 {
	return ps.Priority
}

// ----------------------------------------------------------------

type Environment interface {
	Start(opts ...Option) error
	Destroy() error
	Refresh(opts ...Option) error
	LoadMap(sourceMap collection.AnyMap) error
	LoadPropertySource(sources ...PropertySource) error
	Get(key string) (any, bool)
	NestedGet(key string) (any, bool)
	Set(key string, value any) bool
	NestedSet(key string, value any) bool
	Contains(key string) bool
}

// ----------------------------------------------------------------

// errors 抽象设计

// ----------------------------------------------------------------

type StandardEnvironment struct {
	configMap       collection.AnyMap      // core config container
	propertySources []PropertySource       // config sources
	profiles        collection.StringSlice // Profiles active e.g.: dev test prod ...
}

// ----------------------------------------------------------------

func New(sources ...PropertySource) Environment {
	return &StandardEnvironment{
		configMap:       make(collection.AnyMap),
		propertySources: sources,
		profiles:        make(collection.StringSlice, 0),
	}
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) Start(opts ...Option) error {
	optz := initOptions(opts...)
	e.translateToPropertySources(optz)

	// prepare
	eventPrepare := NewStandardEnvironmentEvent(PrepareEnvironmentEventName, e)
	if err := eventbus.Post(eventPrepare); err != nil {
		return nil
	}

	// pre load
	preLoadEvent := NewStandardEnvironmentEvent(PreLoadEnvironmentEventName, e)
	if err := eventbus.Post(preLoadEvent); err != nil {
		return nil
	}

	// on load
	if err := e.onLoad(); err != nil {
		return nil
	}

	// post load
	postLoadEvent := NewStandardEnvironmentEvent(PostLoadEnvironmentEventName, e)
	if err := eventbus.Post(postLoadEvent); err != nil {
		return nil
	}

	return nil
}

func (e *StandardEnvironment) translateToPropertySources(opts *Options) {

}

func (e *StandardEnvironment) Destroy() error {
	return nil
}

func (e *StandardEnvironment) Refresh(opts ...Option) error {
	// Destroy ...

	return e.Start(opts...)
}

func (e *StandardEnvironment) LoadMap(sourceMap collection.AnyMap) error {
	return nil
}

func (e *StandardEnvironment) LoadPropertySource(sources ...PropertySource) error {
	for _, source := range sources {
		if stringz.IsNotBlankString(source.FilePath) {
			if err := e.loadConfig(source.FilePath, source.Name, source.Type); err != nil {
				return err
			}
		}

		if source.Type != nil && source.Type.Kind() == reflect.Map {
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
	return nil, true
}

func (e *StandardEnvironment) Set(key string, value any) bool {
	return e.setProperty(key, value)
}

func (e *StandardEnvironment) NestedSet(key string, value any) bool {
	return true
}

func (e *StandardEnvironment) Contains(key string) bool {
	return true
}

func (e *StandardEnvironment) setProperty(key string, value any) bool {
	return true
}

func (e *StandardEnvironment) getProperty(key string) (any, bool) {
	return nil, true
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) onLoad() error {
	e.loadSystemEnvVars()

	sorter := ordered.NewSorter(e.propertySources...)
	ordered.Sort(sorter, -1)

	for _, source := range e.propertySources {
		if err := e.LoadPropertySource(source); err != nil {
			return err
		}
	}

	return nil
}

func (e *StandardEnvironment) loadSystemEnvVars() {
	envVars := make(collection.AnyMap)

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

	if stringz.ArrayContains(confusedEnvironments, env) {
		confusedEvent := eventbus.NewStandardAnyEvent(ConfusedEnvironmentEventName, env)
		_ = eventbus.Post(confusedEvent)
	}
}

func initSystemEnvPropertySource(envVars collection.AnyMap) PropertySource {
	return PropertySource{
		Priority: ordered.HighPriority,
		Property: DefaultSystemPropertySourceName,
		Type:     reflect.TypeOf(collection.AnyMap{}),
		Map:      envVars,
	}
}

func (e *StandardEnvironment) loadSystemEnvMap(envVars collection.AnyMap) {
	e.loadSystemEnvMapDelayed(envVars)
}

func (e *StandardEnvironment) loadSystemEnvMapDelayed(envVars collection.AnyMap) {
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

func WithProperties(properties collection.AnyMap) Option {
	return func(opts *Options) {
		opts.Properties = properties
	}
}

// ----------------------------------------------------------------

func initOptions(opts ...Option) *Options {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	populateIfNecessary(options)

	return options
}

func populateIfNecessary(opts *Options) {
	opts.validate()
}

func newOptions() *Options {
	return &Options{
		AbsolutePaths: make(collection.StringSlice, 0),
		ConfigNames:   make(collection.StringSlice, 0),
		ConfigTypes:   make(collection.StringSlice, 0),
		SearchPaths:   make(collection.StringSlice, 0),
		Profiles:      make(collection.StringSlice, 0),
		Sources:       make(PropertySources, 0),
		Properties:    make(collection.AnyMap),
	}
}
