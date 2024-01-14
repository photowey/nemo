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

package filez

import (
	"testing"
)

func TestIsDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "filez#IsDir_test",
			args: args{
				path: "./testdata",
			},
			want: true,
		},
		{
			name: "filez#IsDir_false_1",
			args: args{
				path: "./testdata_not_exist",
			},
			want: false,
		},

		{
			name: "filez#IsDir_false_2",
			args: args{
				path: "./testdata/application.yml",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDir(tt.args.path); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "filez#IsFile_true",
			args: args{
				path: "./testdata/application.yml",
			},
			want: true,
		},
		{
			name: "filez#IsFile_false_1",
			args: args{
				path: "./testdata/application-not-exist.yml",
			},
			want: false,
		},
		{
			name: "filez#IsFile_false_2",
			args: args{
				path: "./testdata",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.args.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
