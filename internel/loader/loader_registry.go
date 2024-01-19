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

package loader

import (
	"github.com/photowey/nemo/pkg/mapz"
)

var (
	_ Registry = (*registry)(nil)

	_registry = &registry{
		loaderMap: make(map[string]ConfigLoader, 1<<2),
	}
)

type Registry interface {
	Register(loader ConfigLoader)
	Contains(name string) bool
	Loaders() []ConfigLoader
}

type registry struct {
	loaderMap map[string]ConfigLoader
}

func (r registry) Register(loader ConfigLoader) {
	r.loaderMap[loader.Name()] = loader
}

func (r registry) Contains(name string) bool {
	_, ok := r.loaderMap[name]
	return ok
}

func (r registry) Loaders() []ConfigLoader {
	return mapz.Values[string, ConfigLoader](r.loaderMap)
}

func Register(loader ConfigLoader) {
	_registry.Register(loader)
}

func Contains(name string) bool {
	return _registry.Contains(name)
}

func Loaders() []ConfigLoader {
	return _registry.Loaders()
}
