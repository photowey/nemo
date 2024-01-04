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
