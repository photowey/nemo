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
	"sort"
	"strings"
)

const (
	single = 1
)

var (
	DefaultEmptySeparator = ""
	DefaultSeparator      = ","
	AndSeparator          = "&"
	EqualJoiner           = "="
)

type StringBuffer struct {
	buffer []string
}

func NewStringBuffer(size ...int) *StringBuffer {
	length := 0
	switch len(size) {
	case single:
		length = size[0]
	}

	return &StringBuffer{
		buffer: make([]string, length),
	}
}

func (sb *StringBuffer) Append(needle string) *StringBuffer {
	buf := append(sb.buffer, needle)
	sb.buffer = buf

	return sb
}

func (sb *StringBuffer) Join(key, value, joiner string) string {
	return key + joiner + value
}

func (sb *StringBuffer) String(separators ...string) string {
	return sb.ToString(separators...)
}

func (sb *StringBuffer) ToString(separators ...string) string {
	separator := AndSeparator
	switch len(separators) {
	case 1:
		separator = separators[0]
	}

	return implode(sb.buffer, separator)
}

func (sb *StringBuffer) ToStrings(separator string) string {
	return sb.ToString(separator)
}

func (sb *StringBuffer) ToSortString() string {
	cloneSlice := sb.cloneSlice()
	sort.Strings(cloneSlice)

	return implode(cloneSlice, AndSeparator)
}

func (sb *StringBuffer) ToSortStrings(separator string) string {
	cloneSlice := sb.cloneSlice()
	sort.Strings(cloneSlice)

	return implode(cloneSlice, separator)
}

func (sb *StringBuffer) Length() int {
	return len(sb.buffer)
}

func (sb *StringBuffer) cloneSlice() []string {
	cloneSlice := make([]string, len(sb.buffer))
	copy(cloneSlice, sb.buffer)

	return cloneSlice
}

func implode(haystack []string, separator string) string {
	var buf strings.Builder
	for _, str := range haystack {
		if "" == str {
			continue
		}
		buf.WriteString(str)
		buf.WriteString(separator)
	}

	return strings.TrimRight(buf.String(), separator)
}
