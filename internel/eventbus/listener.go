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

package eventbus

import (
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/ordered"
)

const (
	NoopListenerName  = "Noop"
	NoopListenerTopic = "nemo.noop.event"
)

var (
	_noop EventListener[Event] = &NoopEventListener[Event]{}
)

func init() {
	_ = Register(_noop)
}

type EventListener[E Event] interface {
	ordered.Ordered
	Name() string
	Topic() collection.StringSlice // []string
	Supports(event string) bool
	OnEvent(event E) error
}

// ----------------------------------------------------------------

type NoopEventListener[E Event] struct{}

func (noop NoopEventListener[E]) Order() int64 {
	return ordered.DefaultPriority
}

func (noop NoopEventListener[E]) Name() string {
	return NoopListenerName
}

func (noop NoopEventListener[E]) Topic() collection.StringSlice {
	return collection.StringSlice{NoopListenerTopic}
}

func (noop NoopEventListener[E]) Supports(event string) bool {
	return NoopListenerTopic == event
}

func (noop NoopEventListener[E]) OnEvent(event E) error {
	return nil
}
