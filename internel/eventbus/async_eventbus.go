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

type asyncEventBus struct {
	listenerMap map[string][]EventListener[Event]
}

// ----------------------------------------------------------------

func (bus *asyncEventBus) Register(listener EventListener[Event]) error {
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

func (bus *asyncEventBus) Post(event Event) error {
	return bus.onEvent(event)
}

// ----------------------------------------------------------------

func (bus *asyncEventBus) onEvent(_ Event) error {
	return nil
}
