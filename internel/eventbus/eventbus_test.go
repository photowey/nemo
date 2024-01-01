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
	"testing"
)

func TestPost(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "eventbus#Post",
			args: args{
				event: NewStandardAnyEvent("hello", "Hello world"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Post(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostAsync(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "eventbus#PostAsync",
			args: args{
				event: NewStandardAnyEvent("hello", "Hello world"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostAsync(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
