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
	"strings"

	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/mapz"
	"github.com/photowey/nemo/pkg/stringz"
)

const (
	binderTag = "binder"
)

type Binder struct {
	Prefix string
}

func New() *Binder {
	return &Binder{}
}

func NewBinder(prefix string) *Binder {
	return &Binder{Prefix: prefix}
}

func (b *Binder) DefaultBind(target any, ctx collection.MixedMap) {
	b.Bind(b.Prefix, target, ctx)
}

func (b *Binder) Bind(prefix string, target any, ctx collection.MixedMap) {
	if stringz.IsNotBlankString(prefix) && stringz.IsNotSuffix(prefix, stringz.Dot) {
		prefix += stringz.Dot
	}

	tt := reflect.TypeOf(target).Elem()
	tv := reflect.ValueOf(target).Elem()

	for i := 0; i < tt.NumField(); i++ {
		t := tt.Field(i)
		v := tv.Field(i)

		tag := t.Tag.Get(binderTag)
		key := stringz.Concat(prefix, strings.ToLower(tag))

		if t.Type.Kind() == reflect.Struct {
			sub := reflect.New(t.Type).Interface()
			b.Bind(key, sub, ctx)
			v.Set(reflect.ValueOf(sub).Elem())
		} else {
			if value, ok := mapz.NestedGet(ctx, key); ok {
				v.Set(reflect.ValueOf(value).Convert(t.Type))
			}
		}
	}
}
