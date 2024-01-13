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

	"github.com/photowey/nemo/pkg/stringz"
	"github.com/photowey/nemo/pkg/valuez"
)

func NestedGet(ctx map[string]any, key string) (any, bool) {
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

		next, ok := value.(map[string]any)
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

func NestedSet(ctx map[string]any, key string, value any) {
	keys := strings.Split(key, stringz.Dot)
	lastKey := keys[len(keys)-1]
	keys = keys[:len(keys)-1]

	current := ctx
	for _, k := range keys {
		next, ok := current[k].(map[string]any)
		if !ok {
			next = make(map[string]any)
			current[k] = next
		}
		current = next
	}

	current[lastKey] = value
}

func MergeMixedMaps(target map[string]any, source map[string]any) {
	for key, sourceValue := range source {
		targetValue, exists := target[key]

		if exists {
			if IsMixedMap(targetValue) && IsMixedMap(sourceValue) {
				MergeMixedMaps(targetValue.(map[string]any), sourceValue.(map[string]any))
			} else {
				target[key] = sourceValue
			}
		} else {
			target[key] = sourceValue
		}
	}
}

// ----------------------------------------------------------------

func Contains[K comparable, V any](key K, ctx map[K]V) bool {
	if valuez.IsNil(ctx) {
		return false
	}

	_, ok := ctx[key]

	return ok
}

func NestedContains(key string, ctx map[string]any) bool {
	if valuez.IsNil(ctx) {
		return false
	}

	keys := strings.Split(key, stringz.Dot)

	for i, k := range keys {
		value, ok := ctx[k]
		if !ok {
			return false
		}

		if i < len(keys)-1 {
			if subMap, isMap := value.(map[string]any); isMap {
				ctx = subMap
			} else {
				return false
			}
		}
	}

	return true
}

// ----------------------------------------------------------------

func IsMap[K comparable, V any](src any) bool {
	_, ok := src.(map[K]V)
	return ok
}

func IsMixedMap(src any) bool {
	_, ok := src.(map[string]any)
	return ok
}
