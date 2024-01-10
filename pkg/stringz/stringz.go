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

package stringz

import (
	"fmt"
	"strings"
)

const (
	Dot         = "."
	EmptyString = ""
	SymbolComma = ","
)

func HasSuffix(source, suffix string) bool {
	return strings.HasSuffix(source, suffix)
}

func IsNotSuffix(source, suffix string) bool {
	return !HasSuffix(source, suffix)
}

func IsBlankString(str string) bool {
	return EmptyString == str || EmptyString == strings.TrimSpace(str)
}

func IsNotBlankString(str string) bool {
	return !IsBlankString(str)
}

// ----------------------------------------------------------------

func Implode(haystack []string, separator string) string {
	if len(haystack) == 0 {
		return EmptyString
	}

	var buf strings.Builder
	for _, str := range haystack {
		if EmptyString == str {
			continue
		}
		buf.WriteString(str)
		buf.WriteString(separator)
	}

	return strings.TrimRight(buf.String(), separator)
}

func Explode(haystack string, separator string) []string {
	if EmptyString == haystack {
		return MakeStringSlice(0)
	}

	return strings.Split(haystack, separator)
}

func MakeStringSlice(length int) []string {
	return make([]string, length)
}

func InitStringSlice(haystack ...string) []string {
	slice := MakeStringSlice(len(haystack))
	copy(slice, haystack)

	return slice
}

// ----------------------------------------------------------------

func Format(sources ...string) string {
	return fmt.Sprintf("%v", sources)
}

func Concat(sources ...string) string {
	return Join(EmptyString, sources...)
}

func Join(separator string, sources ...string) string {
	return strings.Join(sources, separator)
}
