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
			name: "collection#ArrayContains_true",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "b",
			},
			want: true,
		},
		{
			name: "collection#ArrayContains_false",
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
			name: "collection#ArrayNotContains_true",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "d",
			},
			want: true,
		},
		{
			name: "collection#ArrayNotContains_false",
			args: args{
				haystack: []string{"a", "b", "c"},
				needle:   "b",
			},
			want: false,
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
			name: "collection#CloneSlice_true",
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

func TestIsEmptyMap(t *testing.T) {
	type args[K comparable, V any] struct {
		src map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[string, string]{
		{
			name: "collection#IsEmptyMap_true",
			args: args[string, string]{
				src: map[string]string{},
			},
			want: true,
		},
		{
			name: "collection#IsEmptyMap_false",
			args: args[string, string]{
				src: map[string]string{
					"a": "b",
				},
			},
			want: false,
		},
	}

	int64Tests := []testCase[string, int64]{
		{
			name: "collection#IsEmptyMap_int64_true",
			args: args[string, int64]{
				src: map[string]int64{},
			},
			want: true,
		},
		{
			name: "collection#IsEmptyMap_int64_false",
			args: args[string, int64]{
				src: map[string]int64{
					"a": int64(10086),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyMap(tt.args.src); got != tt.want {
				t.Errorf("IsEmptyMap() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range int64Tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyMap(tt.args.src); got != tt.want {
				t.Errorf("IsEmptyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptySlice(t *testing.T) {
	type args[T any] struct {
		src []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "collection#IsNotEmptySlice_true",
			args: args[string]{
				src: []string{},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptySlice_false",
			args: args[string]{
				src: []string{
					"a",
				},
			},
			want: false,
		},
	}

	int64Tests := []testCase[int64]{
		{
			name: "collection#IsNotEmptySlice_int64_true",
			args: args[int64]{
				src: []int64{},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptySlice_int64_false",
			args: args[int64]{
				src: []int64{
					int64(10086),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptySlice(tt.args.src); got != tt.want {
				t.Errorf("IsEmptySlice() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range int64Tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptySlice(tt.args.src); got != tt.want {
				t.Errorf("IsEmptySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmptyMap(t *testing.T) {
	type args[K comparable, V any] struct {
		src map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[string, string]{
		{
			name: "collection#IsNotEmptyMap_true",
			args: args[string, string]{
				src: map[string]string{
					"a": "b",
				},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptyMap_false",
			args: args[string, string]{
				src: map[string]string{},
			},
			want: false,
		},
	}

	int64tests := []testCase[string, int64]{
		{
			name: "collection#IsNotEmptyMap_int64_true",
			args: args[string, int64]{
				src: map[string]int64{
					"a": int64(10086),
				},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptyMap_int64_false",
			args: args[string, int64]{
				src: map[string]int64{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmptyMap(tt.args.src); got != tt.want {
				t.Errorf("IsNotEmptyMap() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range int64tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmptyMap(tt.args.src); got != tt.want {
				t.Errorf("IsNotEmptyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmptySlice(t *testing.T) {
	type args[T any] struct {
		src []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "collection#IsNotEmptySlice_true",
			args: args[string]{
				src: []string{
					"a",
				},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptySlice_false",
			args: args[string]{
				src: []string{},
			},
			want: false,
		},
	}
	int64Tests := []testCase[int64]{
		{
			name: "collection#IsNotEmptySlice_int64_true",
			args: args[int64]{
				src: []int64{
					int64(10086),
				},
			},
			want: true,
		},
		{
			name: "collection#IsNotEmptySlice_int64_false",
			args: args[int64]{
				src: []int64{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmptySlice(tt.args.src); got != tt.want {
				t.Errorf("IsNotEmptySlice() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range int64Tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmptySlice(tt.args.src); got != tt.want {
				t.Errorf("IsNotEmptySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
