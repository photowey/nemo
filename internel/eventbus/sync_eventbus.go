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
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
)

// ----------------------------------------------------------------

type eventBus struct {
	listenerMap map[string][]EventListener[Event]
}

// ----------------------------------------------------------------

func (bus *eventBus) Register(listener EventListener[Event]) error {
	if listener != nil {
		if len(listener.Topic()) == 0 {
			return listenerTopicEmptyError
		}

		for _, topic := range listener.Topic() {
			if _, ok := bus.listenerMap[topic]; !ok {
				bus.listenerMap[topic] = make(EventListenerGroup, 0)
			}

			bus.listenerMap[topic] = append(bus.listenerMap[topic], listener)
		}
		return nil
	}

	return listenerNilError
}

func (bus *eventBus) Post(event Event) error {
	return bus.onEvent(event)
}

// ----------------------------------------------------------------

func (bus *eventBus) onEvent(event Event) error {
	topic := event.Topic()
	if stringz.IsBlankString(topic) {
		topic = event.Name()
	}

	if stringz.IsBlankString(topic) {
		return eventTopicOrNameEmptyError
	}

	listeners := bus.listenerMap[topic]
	sorter := ordered.NewSorter(listeners...)

	ordered.Sort(sorter, 1)
	for _, h := range listeners {
		if h.Supports(topic) {
			if err := h.OnEvent(event); err != nil {
				return nil
			}
		}
	}

	return nil
}
