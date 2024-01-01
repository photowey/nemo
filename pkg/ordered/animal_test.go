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

type Runner interface {
	Ordered
	Name() string
}

// ----------------------------------------------------------------

type Dog struct {
	order int64
}

func NewDog(order int64) *Dog {
	return &Dog{order: order}
}

func (dog *Dog) Order() int64 {
	return dog.order
}

func (dog *Dog) Name() string {
	return "dog"
}

// ----------------------------------------------------------------

type Cat struct {
	order int64
}

func NewCat(order int64) *Cat {
	return &Cat{order: order}
}

func (cat *Cat) Order() int64 {
	return cat.order
}

func (cat *Cat) Name() string {
	return "cat"
}

// ----------------------------------------------------------------

type Rabbit struct {
	order int64
}

func NewRabbit(order int64) *Rabbit {
	return &Rabbit{order: order}
}

func (rabbit *Rabbit) Order() int64 {
	return rabbit.order
}

func (rabbit *Rabbit) Name() string {
	return "rabbit"
}

// ----------------------------------------------------------------

type Turtle struct {
	order int64
}

func NewTurtle(order int64) *Turtle {
	return &Turtle{order: order}
}

func (turtle *Turtle) Order() int64 {
	return turtle.order
}

func (turtle *Turtle) Name() string {
	return "turtle"
}

// ----------------------------------------------------------------
