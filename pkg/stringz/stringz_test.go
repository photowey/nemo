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
	"reflect"
	"testing"
)

func TestArrayContains(t *testing.T) {
	type args struct {
		haystack []string
		needle   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#ArrayContains_true",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "a",
			},
			want: true,
		},
		{
			name: "stringz#ArrayNotContains_false",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayContains(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("ArrayContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayNotContains(t *testing.T) {
	type args struct {
		haystack []string
		needle   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#ArrayNotContains_true",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "a",
			},
			want: false,
		},
		{
			name: "stringz#ArrayNotContains_false",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "d",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayNotContains(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("ArrayNotContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloneSlice(t *testing.T) {
	type args struct {
		src []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "stringz#CloneSlice",
			args: args{
				src: []string{"a", "b", "c"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneSlice(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	type args struct {
		sources []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringz#Concat",
			args: args{
				sources: []string{"a", "b", "c"},
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Concat(tt.args.sources...); got != tt.want {
				t.Errorf("Concat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExplode(t *testing.T) {
	type args struct {
		haystack  string
		separator string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "stringz#Explode",
			args: args{
				haystack:  "a,b,c",
				separator: ",",
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Explode(tt.args.haystack, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Explode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		sources []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringz#Format",
			args: args{
				sources: []string{"a", "b", "c"},
			},
			want: "[a b c]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.sources...); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	type args struct {
		source string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#HasSuffix_true",
			args: args{
				source: "abc.x",
				suffix: "x",
			},
			want: true,
		},
		{
			name: "stringz#HasSuffix_false",
			args: args{
				source: "abc.y",
				suffix: "x",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.source, tt.args.suffix); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplode(t *testing.T) {
	type args struct {
		haystack  []string
		separator string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringz#Implode_,",
			args: args{
				haystack:  []string{"a", "b", "c"},
				separator: ",",
			},
			want: "a,b,c",
		},
		{
			name: "stringz#Implode_@",
			args: args{
				haystack:  []string{"a", "b", "c"},
				separator: "@",
			},
			want: "a@b@c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Implode(tt.args.haystack, tt.args.separator); got != tt.want {
				t.Errorf("Implode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitStringSlice(t *testing.T) {
	type args struct {
		haystack []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "stringz#InitStringSlice",
			args: args{
				haystack: []string{"a", "b", "c"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitStringSlice(tt.args.haystack...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBlankString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#IsBlankString_true",
			args: args{
				str: " ",
			},
			want: true,
		},
		{
			name: "stringz#IsBlankString_false",
			args: args{
				str: "a",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlankString(tt.args.str); got != tt.want {
				t.Errorf("IsBlankString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotBlankString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#IsNotBlankString_true",
			args: args{
				str: "a",
			},
			want: true,
		},
		{
			name: "stringz#IsNotBlankString_false",
			args: args{
				str: " ",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotBlankString(tt.args.str); got != tt.want {
				t.Errorf("IsNotBlankString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotSuffix(t *testing.T) {
	type args struct {
		source string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "stringz#IsNotSuffix_true",
			args: args{
				source: "abc.y",
				suffix: "x",
			},
			want: true,
		},
		{
			name: "stringz#IsNotSuffix_false",
			args: args{
				source: "abc.x",
				suffix: "x",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotSuffix(tt.args.source, tt.args.suffix); got != tt.want {
				t.Errorf("IsNotSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		separator string
		sources   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringz#Join_,",
			args: args{
				separator: ",",
				sources:   []string{"a", "b", "c"},
			},
			want: "a,b,c",
		},
		{
			name: "stringz#Join_@",
			args: args{
				separator: "@",
				sources:   []string{"a", "b", "c"},
			},
			want: "a@b@c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.args.separator, tt.args.sources...); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeStringSlice(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "stringz#MakeStringSlice",
			args: args{
				length: 3,
			},
			want: []string{"", "", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeStringSlice(tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
