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

package mapz

import (
	"strings"

	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/stringz"
)

func NestedGet(ctx collection.MixedMap, key string) (any, bool) {
	keys := strings.Split(key, stringz.Dot)
	current := ctx

	for i, k := range keys {
		value, ok := current[k]
		if !ok {
			if i == len(keys)-1 {
				return current, true
			}
			return nil, false
		}

		next, ok := value.(collection.MixedMap)
		if !ok {
			if i == len(keys)-1 {
				return value, true
			}
			return nil, false
		}

		current = next
	}

	return current, true
}

func NestedSet(ctx collection.MixedMap, key string, value any) {
	keys := strings.Split(key, stringz.Dot)
	lastKey := keys[len(keys)-1]
	keys = keys[:len(keys)-1]

	current := ctx
	for _, k := range keys {
		next, ok := current[k].(collection.MixedMap)
		if !ok {
			next = make(collection.MixedMap)
			current[k] = next
		}
		current = next
	}

	current[lastKey] = value
}
