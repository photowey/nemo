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

	"github.com/photowey/nemo/pkg/collection"
)

type PropertySource struct {
	Property string            // the name of the PropertySource.
	FilePath string            // the path of config file.
	Name     string            // the name of config file.
	Type     reflect.Type      // the type of PropertySource, only support map now.
	Map      collection.AnyMap // the map context, when the Type is map.
}

type Environment struct {
	configMap       collection.AnyMap
	propertySources []PropertySource
}
