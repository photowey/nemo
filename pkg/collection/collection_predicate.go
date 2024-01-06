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

package collection

import (
	"strings"
)

func IsEmptySlice[T any](src []T) bool {
	return len(src) == 0
}

func IsNotEmptySlice[T any](src []T) bool {
	return !IsEmptySlice(src)
}

func IsEmptyMap[K comparable, V any](src map[K]V) bool {
	return len(src) == 0
}

func IsNotEmptyMap[K comparable, V any](src map[K]V) bool {
	return !IsEmptyMap(src)
}

func ArrayContains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.EqualFold(a, needle) {
			return true
		}
	}

	return false
}

func ArrayNotContains(haystack []string, needle string) bool {
	return !ArrayContains(haystack, needle)
}

func CloneSlice(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)

	return dst
}
