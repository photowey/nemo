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

package environment

import (
	"os"
	"reflect"
	"testing"

	"github.com/photowey/nemo/internel/eventbus"
	"github.com/photowey/nemo/pkg/collection"
)

func TestNew(t *testing.T) {
	testdataDir := "testdata"

	if err := os.Mkdir(testdataDir, os.ModePerm); err != nil {
		t.Errorf("nemo: Mkdir testdata failed:%v", err)
	}
	defer os.RemoveAll(testdataDir)

	type args struct {
		sources []PropertySource
	}
	tests := []struct {
		name string
		args args
		want Environment
	}{
		{
			name: "environment#New",
			args: args{
				sources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"},
				},
			},
			want: &StandardEnvironment{
				configMap: make(collection.AnyMap),
				propertySources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"},
				},
				profiles: make(collection.StringSlice, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.sources...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStandardEnvironmentEvent(t *testing.T) {
	environment := New()

	type args struct {
		name string
		data Environment
	}
	tests := []struct {
		name string
		args args
		want eventbus.Event
	}{
		{
			name: "environment#NewStandardEnvironmentEvent",
			args: args{
				name: PrepareEnvironmentEventName,
				data: environment,
			},
			want: &StandardEnvironmentEvent{
				event: PrepareEnvironmentEventName,
				data:  environment,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStandardEnvironmentEvent(tt.args.name, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStandardEnvironmentEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardEnvironment_Start(t *testing.T) {

	properties := make(collection.AnyMap)
	properties["hello"] = "world"

	type fields struct {
		configMap       collection.AnyMap
		propertySources []PropertySource
		profiles        collection.StringSlice
	}
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "environment#Start",
			fields: fields{
				configMap: make(collection.AnyMap),
				propertySources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"},
				},
				profiles: collection.StringSlice{"dev"},
			},
			args: args{
				[]Option{
					WithAbsolutePaths("/opt/data", "/opt/configs"),
					WithConfigNames("application", "config", "configs"),
					WithConfigTypes("yaml", "yml", "toml"),
					WithSearchPaths("resources", "configs"),
					WithProfiles("dev", "test"),
					WithSources(PropertySource{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"}),
					WithProperties(properties),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &StandardEnvironment{
				configMap:       tt.fields.configMap,
				propertySources: tt.fields.propertySources,
				profiles:        tt.fields.profiles,
			}
			if err := e.Start(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
