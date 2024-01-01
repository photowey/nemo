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

const (
	MaxInt64        int64 = 1<<63 - 1 // 9223372036854775807
	MinInt64        int64 = -1 << 63  // -9223372036854775808
	MaxInt32        int32 = 1<<31 - 1 // 2147483647
	MinInt32        int32 = -1 << 31  // -2147483648
	MaxInt          int   = 1<<31 - 1 // 2147483647
	MinInt          int   = -1 << 31  // 2147483648
	LowPriority           = MaxInt64  // 9223372036854775807
	HighPriority          = MinInt64  // -9223372036854775808
	DefaultPriority int64 = 0         // 0
	DefaultStep     int64 = 100       // 100
)
