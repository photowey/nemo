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

package ordered

import (
	"sort"
)

type PriorityOrdered interface {
	Ordered
}

type PrioritySorter[T PriorityOrdered] []PriorityOrdered

func (sorter PrioritySorter[T]) Len() int {
	return len(sorter)
}

func (sorter PrioritySorter[T]) Less(i, j int) bool {
	return sorter[i].Order() < sorter[j].Order()
}

func (sorter PrioritySorter[T]) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}

func NewPrioritySorter[T PriorityOrdered](actors ...T) PrioritySorter[T] {
	sorter := make(PrioritySorter[T], len(actors))
	for i, actor := range actors {
		sorter[i] = actor
	}

	return sorter
}

func PrioritySort[T PriorityOrdered](sorter PrioritySorter[T], sign int) {
	sort.Slice(sorter, func(i, j int) bool {
		if sign > 0 {
			return sorter[i].Order() < sorter[j].Order()
		}
		return sorter[i].Order() > sorter[j].Order()
	})
}
