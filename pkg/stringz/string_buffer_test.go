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

package stringz

import (
	"reflect"
	"testing"
)

func TestNewStringBuffer(t *testing.T) {
	type args struct {
		size []int
	}
	tests := []struct {
		name string
		args args
		want *StringBuffer
	}{
		{
			name: "stringz#NewStringBuffer",
			args: args{
				size: []int{10},
			},
			want: &StringBuffer{
				buffer: make([]string, 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringBuffer(tt.args.size...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStringBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_Append(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		needle string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringBuffer
	}{
		{
			name: "stringz#StringBuffer#Append",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				needle: "d",
			},
			want: &StringBuffer{
				buffer: []string{"a", "b", "c", "d"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.Append(tt.args.needle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_Join(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		key    string
		value  string
		joiner string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "stringz#StringBuffer#Join",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				key:    "a",
				value:  "b",
				joiner: ",",
			},
			want: "a,b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.Join(tt.args.key, tt.args.value, tt.args.joiner); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_Length(t *testing.T) {
	type fields struct {
		buffer []string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "stringz#StringBuffer#Length",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_String(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		separators []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "stringz#StringBuffer#String",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				separators: []string{","},
			},
			want: "a,b,c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.String(tt.args.separators...); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_ToSortString(t *testing.T) {
	type fields struct {
		buffer []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "stringz#StringBuffer#ToSortString",
			fields: fields{
				buffer: []string{"z", "a", "c", "b"},
			},
			want: "a&b&c&z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.ToSortString(); got != tt.want {
				t.Errorf("ToSortString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_ToSortStrings(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		separator string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "stringz#StringBuffer#ToSortStrings_&",
			fields: fields{
				buffer: []string{"z", "a", "c", "b"},
			},
			args: args{
				separator: "&",
			},
			want: "a&b&c&z",
		},
		{
			name: "stringz#StringBuffer#ToSortStrings_space",
			fields: fields{
				buffer: []string{"z", "a", "c", "b"},
			},
			args: args{
				separator: "",
			},
			want: "abcz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.ToSortStrings(tt.args.separator); got != tt.want {
				t.Errorf("ToSortStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_ToString(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		separators []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "stringz#StringBuffer#ToString_comma",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				separators: []string{","},
			},
			want: "a,b,c",
		},
		{
			name: "stringz#StringBuffer#ToString_space",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				separators: []string{""},
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.ToString(tt.args.separators...); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_ToStrings(t *testing.T) {
	type fields struct {
		buffer []string
	}
	type args struct {
		separator string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "stringz#StringBuffer#ToStrings_&",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				separator: "&",
			},
			want: "a&b&c",
		},
		{
			name: "stringz#StringBuffer#ToStrings_space",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			args: args{
				separator: "",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.ToStrings(tt.args.separator); got != tt.want {
				t.Errorf("ToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuffer_cloneSlice(t *testing.T) {
	type fields struct {
		buffer []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "stringz#StringBuffer#cloneSlice",
			fields: fields{
				buffer: []string{"a", "b", "c"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &StringBuffer{
				buffer: tt.fields.buffer,
			}
			if got := sb.cloneSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_implode(t *testing.T) {
	type args struct {
		haystack  []string
		separator string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringz#implode",
			args: args{
				haystack:  []string{"a", "b", "c"},
				separator: ",",
			},
			want: "a,b,c",
		},
		{
			name: "stringz#implode",
			args: args{
				haystack:  []string{"a", "b", "c"},
				separator: "",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := implode(tt.args.haystack, tt.args.separator); got != tt.want {
				t.Errorf("implode() = %v, want %v", got, tt.want)
			}
		})
	}
}
