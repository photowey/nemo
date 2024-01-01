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
	"testing"
)

func TestPriority(t *testing.T) {
	t.Logf("max int64:%d\n", MaxInt64)
	t.Logf("min int64:%d\n", MinInt64)

	t.Logf("max int32:%d\n", MaxInt32)
	t.Logf("min int32:%d\n", MinInt32)

	t.Logf("max int:%d\n", MaxInt)
	t.Logf("min int:%d\n", MinInt)

	t.Logf("high priority :%d\n", HighPriority)
	t.Logf("low priority :%d\n", LowPriority)
	t.Logf("default priority :%d\n", DefaultPriority)
	t.Logf("default strp :%d\n", DefaultStep)
}
