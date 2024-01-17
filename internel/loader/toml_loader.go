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
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
)

// xxx.toml

const (
	Toml = "toml"
)

const (
	tomlStep     = 100
	tomlPriority = ordered.HighPriority + tomlStep*ordered.DefaultStep
)

var (
	tomlSupportedConfigTypes = stringz.InitStringSlice(Toml)
)

var (
	_ ConfigLoader = (*TomlConfigLoader)(nil)
)

type TomlConfigLoader struct {
}

func (tcl *TomlConfigLoader) Supports(strategy string) bool {
	return collection.ArrayContains(tomlSupportedConfigTypes, strategy)
}

func (tcl *TomlConfigLoader) Order() int64 {
	return tomlPriority
}

func (tcl *TomlConfigLoader) Load(path string, targetPtr any) {

}
