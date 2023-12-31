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
	"github.com/photowey/nemo/internel/environment"
)

const (
	PrepareEnvironmentEventName  = "nemo.prepare.environment.event"
	PreLoadEnvironmentEventName  = "nemo.pre.environment.event"
	PostLoadEnvironmentEventName = "nemo.post.environment.event"
)

// ----------------------------------------------------------------

var (
	_ Event[any] = (*StandardEnvironmentEvent)(nil)
	_ Event[any] = (*StandardAnyEvent)(nil)
)

// ----------------------------------------------------------------

type Event[T any] interface {
	Name() string
	Data() T
}

// ----------------------------------------------------------------

type StandardEnvironmentEvent struct {
	event string
	data  environment.Environment
}

func NewStandardEnvironmentEvent(name string, data environment.Environment) Event[any] {
	return &StandardEnvironmentEvent{
		event: name,
		data:  data,
	}
}

func (e *StandardEnvironmentEvent) Name() string {
	return e.event
}

func (e *StandardEnvironmentEvent) Data() any {
	return e.data
}

// ----------------------------------------------------------------

type StandardAnyEvent struct {
	event string
	data  any
}

func NewStandardAnyEvent(name string, data any) Event[any] {
	return &StandardAnyEvent{
		event: name,
		data:  data,
	}
}

func (e *StandardAnyEvent) Name() string {
	return e.event
}

func (e *StandardAnyEvent) Data() any {
	return e.data
}
