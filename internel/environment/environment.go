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
	"reflect"

	"github.com/photowey/nemo/internel/eventbus"
	"github.com/photowey/nemo/pkg/collection"
)

type Option func(opts *Options)

type Options struct {
	AbsolutePaths []string          // /opt/data | /opt/configs | ...
	ConfigNames   []string          // application | config | configs | ...
	ConfigTypes   []string          // yaml/yml | toml -> but, json not support now.
	SearchPaths   []string          // ./resources | ./configs ...
	Profiles      []string          // dev | test | prod | ...
	Sources       []PropertySource  // PropertySource
	Properties    collection.AnyMap // Properties -> map data-structure -> can also be replaced by PropertySource
	ThrowLevel    int               // 0: all 1:anyone
}

// ----------------------------------------------------------------

type PropertySource struct {
	Property string            // the name of the PropertySource.
	FilePath string            // the path of config file. -> e.g.: /opt/data | /opt/configs
	Name     string            // the name of config file. -> e.g.: config.yaml | config.yml config.toml
	Suffix   string            // the suffix of config file -> Name == config Suffix == yml -> Full name == config.yml
	Type     reflect.Type      // the type of PropertySource, only support map now.
	Map      collection.AnyMap // the map context, when the Type is map.
}

type Environment interface {
	Start(opts ...Option) error
	LoadMap(sourceMap collection.AnyMap) error
	LoadPropertySource(sources ...PropertySource) error
	Get(key string) (any, bool)
	Set(key string, value any) bool
	Contains(key string) bool
}

// ----------------------------------------------------------------

type StandardEnvironment struct {
	configMap       collection.AnyMap // core config container
	propertySources []PropertySource  // config sources
	profiles        []string          // Profiles active e.g.: dev test prod ...
}

// ----------------------------------------------------------------

func New(sources ...PropertySource) Environment {
	return &StandardEnvironment{
		configMap:       make(collection.AnyMap),
		propertySources: sources,
		profiles:        make([]string, 0),
	}
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) Start(opts ...Option) error {
	initOptions(opts...)

	// prepare
	eventPrepare := eventbus.NewStandardEnvironmentEvent(eventbus.PrepareEnvironmentEventName, e)
	if err := eventbus.Post(eventPrepare); err != nil {
		return nil
	}

	// pre load
	preLoadEvent := eventbus.NewStandardEnvironmentEvent(eventbus.PreLoadEnvironmentEventName, e)
	if err := eventbus.Post(preLoadEvent); err != nil {
		return nil
	}

	// on load
	if err := e.onLoad(); err != nil {
		return nil
	}

	// post load
	postLoadEvent := eventbus.NewStandardEnvironmentEvent(eventbus.PostLoadEnvironmentEventName, e)
	if err := eventbus.Post(postLoadEvent); err != nil {
		return nil
	}

	return nil
}

func (e *StandardEnvironment) LoadMap(sourceMap collection.AnyMap) error {
	return nil
}

func (e *StandardEnvironment) LoadPropertySource(sources ...PropertySource) error {
	return nil
}

func (e *StandardEnvironment) Get(key string) (any, bool) {
	return nil, true
}

func (e *StandardEnvironment) Set(key string, value any) bool {
	return true
}

func (e *StandardEnvironment) Contains(key string) bool {
	return true
}

// ----------------------------------------------------------------

func (e *StandardEnvironment) onLoad() error {
	return nil
}

func (e *StandardEnvironment) loadSystemEnvVars() {

}

func (e *StandardEnvironment) loadConfig(path, name string, _ reflect.Type) {

}

// ----------------------------------------------------------------

func WithAbsolutePaths(absolutePaths ...string) Option {
	return func(opts *Options) {
		opts.AbsolutePaths = absolutePaths
	}
}

func WithConfigNames(configNames ...string) Option {
	return func(opts *Options) {
		opts.ConfigNames = configNames
	}
}

func WithConfigTypes(configTypes ...string) Option {
	return func(opts *Options) {
		opts.ConfigTypes = configTypes
	}
}

func WithSearchPaths(searchPaths ...string) Option {
	return func(opts *Options) {
		opts.SearchPaths = searchPaths
	}
}

func WithProfiles(profiles ...string) Option {
	return func(opts *Options) {
		opts.Profiles = profiles
	}
}

func WithSources(sources ...PropertySource) Option {
	return func(opts *Options) {
		opts.Sources = sources
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

	// do something

	return options
}

func newOptions() *Options {
	return &Options{
		AbsolutePaths: make([]string, 0),
		ConfigNames:   make([]string, 0),
		ConfigTypes:   make([]string, 0),
		SearchPaths:   make([]string, 0),
		Profiles:      make([]string, 0),
		Sources:       make([]PropertySource, 0),
		Properties:    make(collection.AnyMap),
	}
}
