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

var (
	_ Event = (*StandardAnyEvent)(nil)
)

// ----------------------------------------------------------------

type Event interface {
	Name() string
	Data() any
}

// ----------------------------------------------------------------

type StandardAnyEvent struct {
	event string
	data  any
}

func NewStandardAnyEvent(name string, data any) Event {
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
