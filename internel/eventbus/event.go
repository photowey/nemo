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
	PrepareEnvironmentEventName = "nemo.prepare.environment.event"
)

// ----------------------------------------------------------------

var (
	_ Event[environment.Environment] = (*PrepareEnvironmentEvent)(nil)
)

// ----------------------------------------------------------------

type Event[T any] interface {
	Name() string
	Data() T
}

type PrepareEnvironmentEvent struct {
	event string
	env   environment.Environment
}

func NewPrepareEnvironmentEvent(env environment.Environment) Event[environment.Environment] {
	return &PrepareEnvironmentEvent{
		event: PrepareEnvironmentEventName,
		env:   env,
	}
}

func (e *PrepareEnvironmentEvent) Name() string {
	return e.event
}

func (e *PrepareEnvironmentEvent) Data() environment.Environment {
	return e.env
}
