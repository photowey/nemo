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
	"errors"
	"reflect"
	"testing"
	"time"

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

type StringMain struct {
	A string  `binder:"a"`
	B int     `binder:"b"`
	C bool    `binder:"c"`
	D float64 `binder:"d"`
}

type TaggedConfig struct {
	RequiredName string `binder:"name" required:"true"`
	DefaultHost  string `binder:"host" default:"127.0.0.1"`
	DefaultPort  int    `binder:"port" default:"8080"`
}

type InvalidTaggedConfig struct {
	Flag bool `required:"true"`
}

type InvalidRequiredConfig struct {
	Name string `binder:"name" required:"sometimes"`
}

type RichConfig struct {
	Timeout time.Duration `binder:"timeout" default:"5s"`
	Tags    []string      `binder:"tags" default:"a,b,c"`
	Ports   []int         `binder:"ports"`
	PortPtr *int          `binder:"portPtr" default:"7002"`
	Sub     *Sub          `binder:"sub"`
}

func TestBinder_Bind(t *testing.T) {
	type args struct {
		prefix string
		target Main
		ctx    collection.MixedMap
	}
	tests := []struct {
		name string
		args args
		want Main
		wantErr bool
	}{
		{
			name: "builder#Bind",
			args: args{
				prefix: "a.b.c",
				target: Main{},
				ctx: collection.MixedMap{
					"a": collection.MixedMap{
						"b": collection.MixedMap{
							"c": collection.MixedMap{
								"d": "Hello",
								"e": 42,
								"f": true,
								"g": collection.MixedMap{
									"h": 3.14,
								},
								"sub": collection.MixedMap{
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
			err := b.Bind(tt.args.prefix, &tt.args.target, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		ctx    collection.MixedMap
	}
	tests := []struct {
		name string
		args args
		want Main
		wantErr bool
	}{
		{
			name: "builder#DefaultBind",
			args: args{
				prefix: "a.b.c",
				target: Main{},
				ctx: collection.MixedMap{
					"a": collection.MixedMap{
						"b": collection.MixedMap{
							"c": collection.MixedMap{
								"d": "Hello",
								"e": 42,
								"f": true,
								"g": collection.MixedMap{
									"h": 3.14,
								},
								"sub": collection.MixedMap{
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
			err := b.DefaultBind(&tt.args.target, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DefaultBind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Expected %+v, but got %+v", tt.want, tt.args.target)
			}
		})
	}
}

func TestBinder_BindStringConversions(t *testing.T) {
	target := StringMain{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{
			"a": "hello",
			"b": "42",
			"c": "true",
			"d": "3.14",
		},
	}

	err := New().Bind("cfg", &target, ctx)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	want := StringMain{A: "hello", B: 42, C: true, D: 3.14}
	if !reflect.DeepEqual(target, want) {
		t.Fatalf("Bind() target = %+v, want %+v", target, want)
	}
}

func TestBinder_BindInvalidConversion(t *testing.T) {
	target := StringMain{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{
			"a": "hello",
			"b": "not-a-number",
		},
	}

	err := New().Bind("cfg", &target, ctx)
	if err == nil {
		t.Fatalf("expected conversion error")
	}
}

func TestBinder_BindRequiredField(t *testing.T) {
	target := TaggedConfig{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{},
	}

	err := New().Bind("cfg", &target, ctx)
	if err == nil {
		t.Fatalf("expected required field error")
	}
	var bindErr *BindError
	if !errors.As(err, &bindErr) {
		t.Fatalf("expected BindError, got %T", err)
	}
	if bindErr.Kind != MissingRequiredErrorKind {
		t.Fatalf("expected missing required bind error kind, got %s", bindErr.Kind)
	}
}

func TestBinder_BindDefaultValues(t *testing.T) {
	target := TaggedConfig{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{
			"name": "demo",
		},
	}

	err := New().Bind("cfg", &target, ctx)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	want := TaggedConfig{
		RequiredName: "demo",
		DefaultHost:  "127.0.0.1",
		DefaultPort:  8080,
	}
	if !reflect.DeepEqual(target, want) {
		t.Fatalf("Bind() target = %+v, want %+v", target, want)
	}
}

func TestBinder_RejectsConstraintsWithoutBinderTag(t *testing.T) {
	target := InvalidTaggedConfig{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{},
	}

	err := New().Bind("cfg", &target, ctx)
	if err == nil {
		t.Fatalf("expected binder constraint validation error")
	}
	var bindErr *BindError
	if !errors.As(err, &bindErr) {
		t.Fatalf("expected BindError, got %T", err)
	}
	if bindErr.Kind != InvalidTagErrorKind {
		t.Fatalf("expected invalid tag error kind, got %s", bindErr.Kind)
	}
}

func TestBinder_RejectsInvalidRequiredTagValue(t *testing.T) {
	target := InvalidRequiredConfig{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{
			"name": "demo",
		},
	}

	err := New().Bind("cfg", &target, ctx)
	if err == nil {
		t.Fatalf("expected invalid required tag error")
	}
	var bindErr *BindError
	if !errors.As(err, &bindErr) {
		t.Fatalf("expected BindError, got %T", err)
	}
	if bindErr.Kind != InvalidTagErrorKind {
		t.Fatalf("expected invalid tag error kind, got %s", bindErr.Kind)
	}
}

func TestBinder_BindRichTypes(t *testing.T) {
	target := RichConfig{}
	ctx := collection.MixedMap{
		"cfg": collection.MixedMap{
			"ports":   []any{"8080", "9090"},
			"timeout": "3s",
			"sub": collection.MixedMap{
				"x": "nested",
				"y": 12,
			},
		},
	}

	err := New().Bind("cfg", &target, ctx)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	if target.Timeout != 3*time.Second {
		t.Fatalf("expected timeout to be parsed, got %v", target.Timeout)
	}
	if !reflect.DeepEqual(target.Tags, []string{"a", "b", "c"}) {
		t.Fatalf("expected default tags, got %+v", target.Tags)
	}
	if !reflect.DeepEqual(target.Ports, []int{8080, 9090}) {
		t.Fatalf("expected parsed ports, got %+v", target.Ports)
	}
	if target.PortPtr == nil || *target.PortPtr != 7002 {
		t.Fatalf("expected default pointer port, got %+v", target.PortPtr)
	}
	if target.Sub == nil || target.Sub.X != "nested" || target.Sub.Y != 12 {
		t.Fatalf("expected bound nested pointer struct, got %+v", target.Sub)
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
