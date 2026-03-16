/*
 * Copyright © 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nemo

import (
	"errors"

	"github.com/photowey/nemo/internel/binder"
	"github.com/photowey/nemo/internel/environment"
	"github.com/photowey/nemo/pkg/collection"
)

type Environment = environment.Environment
type StandardEnvironment = environment.StandardEnvironment
type PropertySource = environment.PropertySource
type PropertySources = environment.PropertySources
type Option = environment.Option
type SuccessThreshold = environment.SuccessThreshold
type ActiveProfile = environment.ActiveProfile
type Binder = binder.Binder
type BindError = binder.BindError
type ErrorKind = binder.ErrorKind
type MixedMap = collection.MixedMap

var (
	NoneSuccessThreshold   = environment.NoneSuccessThreshold
	AnyoneSuccessThreshold = environment.AnyoneSuccessThreshold
	AllSuccessThreshold    = environment.AllSuccessThreshold
	InvalidTargetErrorKind   = binder.InvalidTargetErrorKind
	InvalidTagErrorKind      = binder.InvalidTagErrorKind
	UnsettableFieldErrorKind = binder.UnsettableFieldErrorKind
	MissingRequiredErrorKind = binder.MissingRequiredErrorKind
	UnsupportedTypeErrorKind = binder.UnsupportedTypeErrorKind
	ConversionFailedErrorKind = binder.ConversionFailedErrorKind
)

func New(sources ...PropertySource) Environment {
	return environment.New(sources...)
}

func RegisterSpecialEnvironment(key, value string) {
	environment.RegisterSpecialEnvironment(key, value)
}

func WithAbsolutePaths(absolutePaths ...string) Option {
	return environment.WithAbsolutePaths(absolutePaths...)
}

func WithConfigNames(configNames ...string) Option {
	return environment.WithConfigNames(configNames...)
}

func WithConfigTypes(configTypes ...string) Option {
	return environment.WithConfigTypes(configTypes...)
}

func WithSearchPaths(searchPaths ...string) Option {
	return environment.WithSearchPaths(searchPaths...)
}

func WithProfiles(profiles ...string) Option {
	return environment.WithProfiles(profiles...)
}

func WithSources(sources ...PropertySource) Option {
	return environment.WithSources(sources...)
}

func WithProperties(properties MixedMap) Option {
	return environment.WithProperties(properties)
}

func WithThreshold(threshold SuccessThreshold) Option {
	return environment.WithThreshold(threshold)
}

func Bind[T any](env Environment, prefix string) (T, error) {
	var target T
	err := env.Bind(prefix, &target)
	return target, err
}

func MustBind[T any](env Environment, prefix string) T {
	target, err := Bind[T](env, prefix)
	if err != nil {
		panic(err)
	}
	return target
}

func AsBindError(err error) (*BindError, bool) {
	var bindErr *BindError
	if errors.As(err, &bindErr) {
		return bindErr, true
	}
	return nil, false
}

func IsBindErrorKind(err error, kind ErrorKind) bool {
	bindErr, ok := AsBindError(err)
	return ok && bindErr.Kind == kind
}
