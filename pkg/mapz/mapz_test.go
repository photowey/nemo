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
	"reflect"
	"testing"

	"github.com/photowey/nemo/pkg/collection"
)

func TestNestedGet(t *testing.T) {
	type args struct {
		ctx collection.MixedMap
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   any
		wantOk bool
	}{
		{
			name: "mapz#NestedGet_true",
			args: args{
				ctx: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				key: "b.c",
			},
			want:   2,
			wantOk: true,
		},
		{
			name: "mapz#NestedGet_false_1",
			args: args{
				ctx: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				key: "a.c",
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "mapz#NestedGet_false_2",
			args: args{
				ctx: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				key: "x.y",
			},
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NestedGet(tt.args.ctx, tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NestedGet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("NestedGet() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func TestNestedSet(t *testing.T) {
	type args struct {
		ctx   collection.MixedMap
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want collection.MixedMap
	}{
		{
			name: "mapz#NestedSet_1",
			args: args{
				ctx: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				key:   "b.d",
				value: 3,
			},
			want: collection.MixedMap{
				"a": 1,
				"b": collection.MixedMap{
					"c": 2,
					"d": 3,
				},
			},
		},
		{
			name: "mapz#NestedSet_2",
			args: args{
				ctx: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				key: "b.e",
				value: collection.MixedMap{
					"c": 2,
					"d": 3,
				},
			},
			want: collection.MixedMap{
				"a": 1,
				"b": collection.MixedMap{
					"c": 2,
					"e": collection.MixedMap{
						"c": 2,
						"d": 3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NestedSet(tt.args.ctx, tt.args.key, tt.args.value)

			if !reflect.DeepEqual(tt.args.ctx, tt.want) {
				t.Errorf("Expected %+v, but got %+v", tt.want, tt.args.ctx)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
	type args[K comparable, V any] struct {
		src any
	}

	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[string, string]{
		{
			name: "mapz#IsMap_true",
			args: args[string, string]{
				src: map[string]string{
					"a": "1",
				},
			},
			want: true,
		},
		{
			name: "mapz#IsMap_false",
			args: args[string, string]{
				src: "1",
			},
			want: false,
		},
	}
	int64Tests := []testCase[string, int64]{
		{
			name: "mapz#IsMap_int64_true",
			args: args[string, int64]{
				src: map[string]int64{
					"a": int64(10086),
				},
			},
			want: true,
		},
		{
			name: "mapz#IsMap_int64_false",
			args: args[string, int64]{
				src: 10086,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMap[string, string](tt.args.src); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range int64Tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMap[string, int64](tt.args.src); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMixedMap(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mapz#IsMixedMap_true",
			args: args{
				src: collection.MixedMap{
					"a": 1,
				},
			},
			want: true,
		},
		{
			name: "mapz#IsMixedMap_false",
			args: args{
				src: map[string]string{
					"a": "1",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMixedMap(tt.args.src); got != tt.want {
				t.Errorf("IsMixedMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMixedMaps(t *testing.T) {
	type args struct {
		target collection.MixedMap
		source collection.MixedMap
	}
	tests := []struct {
		name string
		args args
		want collection.MixedMap
	}{
		{
			name: "mapz#MergeMixedMaps",
			args: args{
				target: collection.MixedMap{
					"a": 1,
					"b": collection.MixedMap{
						"c": 2,
					},
				},
				source: collection.MixedMap{
					"b": collection.MixedMap{
						"d": 3,
					},
				},
			},
			want: collection.MixedMap{
				"a": 1,
				"b": collection.MixedMap{
					"c": 2,
					"d": 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeMixedMaps(tt.args.target, tt.args.source)
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("MergeMixedMaps() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args[K comparable, V any] struct {
		key K
		ctx map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[string, string]{
		{
			name: "mapz#Contains_string_true",
			args: args[string, string]{
				key: "a",
				ctx: map[string]string{
					"a": "1",
				},
			},
			want: true,
		},
		{
			name: "mapz#Contains_string_false",
			args: args[string, string]{
				key: "a.b",
				ctx: map[string]string{
					"b": "1",
				},
			},
			want: false,
		},
		{
			name: "mapz#Contains_string_false",
			args: args[string, string]{
				key: "a",
				ctx: map[string]string{
					"b": "1",
				},
			},
			want: false,
		},
	}

	intTests := []testCase[string, int64]{
		{
			name: "mapz#Contains_int64_true",
			args: args[string, int64]{
				key: "a",
				ctx: map[string]int64{
					"a": int64(10086),
				},
			},
			want: true,
		},
		{
			name: "mapz#Contains_int64_false",
			args: args[string, int64]{
				key: "a",
				ctx: map[string]int64{
					"b": int64(10086),
				},
			},
			want: false,
		},
		{
			name: "mapz#Contains_int64_false",
			args: args[string, int64]{
				key: "a.b",
				ctx: map[string]int64{
					"b": int64(10086),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.key, tt.args.ctx); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.key, tt.args.ctx); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNestedContains(t *testing.T) {
	type args struct {
		key string
		ctx map[string]any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mapz#NestedContains_true",
			args: args{
				key: "a.b",
				ctx: map[string]any{
					"a": map[string]any{
						"b": 1,
					},
				},
			},
			want: true,
		},
		{
			name: "mapz#NestedContains_false",
			args: args{
				key: "a.b",
				ctx: map[string]any{
					"a": map[string]any{},
				},
			},
			want: false,
		},
		{
			name: "mapz#NestedContains_single_key_true",
			args: args{
				key: "a",
				ctx: map[string]any{
					"a": map[string]any{
						"b": 1,
					},
				},
			},
			want: true,
		},
		{
			name: "mapz#NestedContains_single_key_false",
			args: args{
				key: "b",
				ctx: map[string]any{
					"a": map[string]any{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NestedContains(tt.args.key, tt.args.ctx); got != tt.want {
				t.Errorf("NestedContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
