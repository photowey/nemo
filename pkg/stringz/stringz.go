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
)

func HasSuffix(source, suffix string) bool {
	return strings.HasSuffix(source, suffix)
}

func IsNotSuffix(source, suffix string) bool {
	return !HasSuffix(source, suffix)
}

func IsBlankString(str string) bool {
	return EmptyString == str
}

func IsNotBlankString(str string) bool {
	return !IsBlankString(str)
}

func Format(sources ...string) string {
	return fmt.Sprintf("%v", sources)
}

func Concat(sources ...string) string {
	return Join("", sources...)
}

func Join(separator string, sources ...string) string {
	return strings.Join(sources, separator)
}
