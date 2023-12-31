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

const (
	NoopListenerName = "Noop"
)

var (
	_     EventListener[Event[any], any] = (*NoopEventListener[Event[any], any])(nil)
	_noop EventListener[Event[any], any] = &NoopEventListener[Event[any], any]{}
)

func init() {
	_ = Register(_noop)
}

type EventListener[E Event[D], D any] interface {
	Name() string
	Supports(event string) bool
	OnEvent(event E[D]) error
}

// ----------------------------------------------------------------

type NoopEventListener[E Event[D], D any] struct {
}

func (noop NoopEventListener[E, D]) Name() string {
	return NoopListenerName
}

func (noop NoopEventListener[E, D]) Supports(event string) bool {
	return false
}

func (noop NoopEventListener[E, D]) OnEvent(event E[D]) error {
	return nil
}
