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
	"errors"
)

var (
	listenerNilError = errors.New("nemo: listener can't be nil on `Register` action")
)

var (
	_eventbus = &eventBus[Event[any], EventListener[Event[any], any], any]{
		listeners: make([]EventListener[Event[any], any], 0),
	}
)

// ----------------------------------------------------------------

func init() {

}

// ----------------------------------------------------------------

type EventBus[E Event[D], T EventListener[E, D], D any] interface {
	Register(listener T[E[D], D]) error
	Post(event E[D]) error
}

type eventBus[E Event[D], T EventListener[E, D], D any] struct {
	listeners []T[E[D]]
}

// ----------------------------------------------------------------

func (bus *eventBus[E, T, D]) Register(listener T[E[D], D]) error {
	if listener != nil {
		bus.listeners = append(bus.listeners, listener)
		return nil
	}

	return listenerNilError
}

func (bus *eventBus[E, T, D]) Post(event E[D]) error {
	return bus.onEvent(event)
}

// ----------------------------------------------------------------

func (bus *eventBus[E, T, D]) onEvent(event E[D]) error {
	// sort ?
	for _, h := range bus.listeners {
		if h.Supports(event.Name()) {
			return h.OnEvent(event)
		}
	}

	return nil
}

// ----------------------------------------------------------------

func Register(listener EventListener[Event[any], any]) error {
	return _eventbus.Register(listener)
}

func Post(event Event[any]) error {
	return _eventbus.Post(event)
}
