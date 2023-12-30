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

package binder

import (
	"reflect"
	"testing"

	"github.com/photowey/nemo/pkg/collection"
)

type Sub struct {
	X string `binder:"x"`
	Y int    `binder:"y"`
}

type Main struct {
	A   string  `binder:"d"`
	B   int     `binder:"e"`
	C   bool    `binder:"f"`
	Z   float64 `binder:"g.h"`
	Sub Sub     `binder:"sub"`
}

func TestBinder_Bind(t *testing.T) {
	type args struct {
		prefix string
		target Main
		ctx    collection.AnyMap
	}
	tests := []struct {
		name string
		args args
		want Main
	}{
		{
			name: "builder#Bind",
			args: args{
				prefix: "a.b.c",
				target: Main{},
				ctx: collection.AnyMap{
					"a": collection.AnyMap{
						"b": collection.AnyMap{
							"c": collection.AnyMap{
								"d": "Hello",
								"e": 42,
								"f": true,
								"g": collection.AnyMap{
									"h": 3.14,
								},
								"sub": collection.AnyMap{
									"x": "Nested",
									"y": 123,
								},
							},
						},
					},
				},
			},
			want: Main{A: "Hello", B: 42, C: true, Z: 3.14, Sub: Sub{X: "Nested", Y: 123}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New()
			b.Bind(tt.args.prefix, &tt.args.target, tt.args.ctx)
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Expected %+v, but got %+v", tt.want, tt.args.target)
			}
		})
	}
}

func TestBinder_DefaultBind(t *testing.T) {
	type args struct {
		prefix string
		target Main
		ctx    collection.AnyMap
	}
	tests := []struct {
		name string
		args args
		want Main
	}{
		{
			name: "builder#DefaultBind",
			args: args{
				prefix: "a.b.c",
				target: Main{},
				ctx: collection.AnyMap{
					"a": collection.AnyMap{
						"b": collection.AnyMap{
							"c": collection.AnyMap{
								"d": "Hello",
								"e": 42,
								"f": true,
								"g": collection.AnyMap{
									"h": 3.14,
								},
								"sub": collection.AnyMap{
									"x": "Nested",
									"y": 123,
								},
							},
						},
					},
				},
			},
			want: Main{A: "Hello", B: 42, C: true, Z: 3.14, Sub: Sub{X: "Nested", Y: 123}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBinder(tt.args.prefix)
			b.DefaultBind(&tt.args.target, tt.args.ctx)
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Expected %+v, but got %+v", tt.want, tt.args.target)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Binder
	}{
		{
			name: "builder#New",
			want: &Binder{
				Prefix: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBinder(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want *Binder
	}{
		{
			name: "builder#NewBinder",
			args: args{
				prefix: "a.b.c",
			},
			want: &Binder{
				Prefix: "a.b.c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBinder(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBinder() = %v, want %v", got, tt.want)
			}
		})
	}
}
