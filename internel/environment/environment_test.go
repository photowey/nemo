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
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/photowey/nemo/internel/binder"
	"github.com/photowey/nemo/internel/eventbus"
	"github.com/photowey/nemo/pkg/collection"
)

func TestNew(t *testing.T) {
	testdataDir := "testdata"

	if err := os.Mkdir(testdataDir, os.ModePerm); err != nil {
		t.Errorf("nemo: Mkdir testdata failed:%v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(testdataDir); err != nil {
			t.Logf("failed to remove test data dir %s: %v", testdataDir, err)
		}
	})

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
				configMap: make(collection.MixedMap),
				propertySources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"},
				},
				initialPropertySources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "testdata", Name: "application-dev", Suffix: "yaml"},
				},
				profiles:  make(collection.StringSlice, 0),
				threshold: NoneSuccessThreshold,
				binder:    binder.New(),
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
	absTestDataDir := filepath.Clean(filepath.Join(testSourceDir(), "../../tests/testdata"))
	absMissingDir := filepath.Clean(filepath.Join(testSourceDir(), "../../tests/missing"))
	properties := make(collection.MixedMap)
	properties["hello"] = "world"

	type fields struct {
		configMap       collection.MixedMap
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
				configMap: make(collection.MixedMap),
				propertySources: []PropertySource{
					{Priority: 1, Property: "dev", FilePath: "../../tests/testdata", Name: "application-dev", Suffix: "yml"},
				},
				profiles: collection.StringSlice{"dev"},
			},
			args: args{
				[]Option{
					WithAbsolutePaths(absMissingDir, absTestDataDir),
					WithConfigNames("application", "config", "configs"),
					WithConfigTypes("yaml", "yml", "toml"),
					WithSearchPaths("resources", "configs"),
					WithProfiles("dev", "test"),
					WithSources(PropertySource{Priority: 1, Property: "dev", FilePath: "../../tests/testdata", Name: "application-dev", Suffix: "yml"}),
					WithProperties(properties),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &StandardEnvironment{
				configMap:              tt.fields.configMap,
				propertySources:        tt.fields.propertySources,
				initialPropertySources: append(make([]PropertySource, 0), tt.fields.propertySources...),
				profiles:               tt.fields.profiles,
				threshold:              NoneSuccessThreshold,
				binder:                 binder.New(),
			}
			if err := e.Start(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStandardEnvironment_StartWithThresholdError(t *testing.T) {
	e := New().(*StandardEnvironment)

	err := e.Start(
		WithSources(PropertySource{
			Priority: 1,
			Property: "missing",
			FilePath: "../../tests/testdata",
			Name:     "missing",
			Suffix:   "yml",
		}),
		WithThreshold(AllSuccessThreshold),
	)
	if err == nil {
		t.Fatalf("expected threshold load error")
	}
}

func TestStandardEnvironment_StartLoadsAbsoluteFilePath(t *testing.T) {
	e := New().(*StandardEnvironment)
	absPath := filepath.Clean(filepath.Join(testSourceDir(), "../../tests/testdata/application.yml"))

	err := e.Start(
		WithAbsolutePaths(absPath),
		WithThreshold(AnyoneSuccessThreshold),
	)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	value, ok := e.Get("nemo.application.name")
	if !ok || value != "nemoapp" {
		t.Fatalf("expected nemo.application.name to be loaded, got value=%v ok=%v", value, ok)
	}
}

func TestStandardEnvironment_StartLoadsProfileSpecificConfig(t *testing.T) {
	e := New().(*StandardEnvironment)
	searchPath := filepath.Clean(filepath.Join(testSourceDir(), "../../tests/testdata"))

	err := e.Start(
		WithSearchPaths(searchPath),
		WithConfigNames("application"),
		WithConfigTypes("yml"),
		WithProfiles("dev"),
		WithThreshold(AnyoneSuccessThreshold),
	)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	value, ok := e.Get("nemo.profiles.active")
	if !ok || value != "dev" {
		t.Fatalf("expected profile-specific config override, got value=%v ok=%v", value, ok)
	}
}

func TestStandardEnvironment_RefreshResetsContext(t *testing.T) {
	e := New().(*StandardEnvironment)

	if err := e.Start(WithProperties(collection.MixedMap{
		"nemo": collection.MixedMap{
			"test": collection.MixedMap{
				"first": "world",
			},
		},
	})); err != nil {
		t.Fatalf("initial Start() error = %v", err)
	}
	if value, ok := e.Get("nemo.test.first"); !ok || value != "world" {
		t.Fatalf("expected initial property to be present")
	}

	if err := e.Refresh(WithProperties(collection.MixedMap{
		"nemo": collection.MixedMap{
			"test": collection.MixedMap{
				"second": "moon",
			},
		},
	})); err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}

	if _, ok := e.Get("nemo.test.first"); ok {
		t.Fatalf("expected refresh to clear old properties")
	}
	if value, ok := e.Get("nemo.test.second"); !ok || value != "moon" {
		t.Fatalf("expected refreshed property to be present")
	}
}

func testSourceDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("failed to determine source directory")
	}
	return filepath.Dir(filename)
}
