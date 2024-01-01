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
	"reflect"
	"testing"
)

func TestPrioritySort(t *testing.T) {
	type args struct {
		sorter PrioritySorter[Runner]
		sign   int
	}
	tests := []struct {
		name string
		args args
		want PrioritySorter[Runner]
	}{
		{
			name: "ordered#PrioritySort#ASC",
			args: args{
				sorter: PrioritySorter[Runner]{
					NewDog(HighPriority),
					NewTurtle(LowPriority),
					NewCat(HighPriority + DefaultStep),
					NewRabbit(DefaultPriority),
				},
				sign: 1,
			},
			want: PrioritySorter[Runner]{
				NewDog(HighPriority),
				NewCat(HighPriority + DefaultStep),
				NewRabbit(DefaultPriority),
				NewTurtle(LowPriority),
			},
		},
		{
			name: "ordered#PrioritySort#DESC",
			args: args{
				sorter: PrioritySorter[Runner]{
					NewDog(HighPriority),
					NewTurtle(LowPriority),
					NewCat(HighPriority + DefaultStep),
					NewRabbit(DefaultPriority),
				},
				sign: -1,
			},
			want: PrioritySorter[Runner]{
				NewTurtle(LowPriority),
				NewRabbit(DefaultPriority),
				NewCat(HighPriority + DefaultStep),
				NewDog(HighPriority),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrioritySort(tt.args.sorter, tt.args.sign)
			if !reflect.DeepEqual(tt.args.sorter, tt.want) {
				t.Errorf("Expected %+v, but got %+v", tt.want, tt.args.sorter)
			}
		})
	}
}
