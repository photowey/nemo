/*
 * Copyright © 2023-present the nemo authors. All rights reserved.
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

package loader

import (
	"path/filepath"
	"testing"

	"github.com/photowey/nemo/pkg/mapz"
)

func TestPropertiesConfigLoader_Load(t *testing.T) {
	testFile := determineTestSourceFilePath()
	testdataDir := filepath.Dir(testFile)

	absPath := filepath.Clean(filepath.Join(testdataDir, "../../tests/testdata/application.properties"))
	badAbsPath := filepath.Clean(filepath.Join(testdataDir, "../../tests/testdata/config.properties")) // not found

	ctx := make(map[string]any)

	type args struct {
		path      string
		targetPtr any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "loader#prroperties_ok",
			args: args{
				path:      absPath,
				targetPtr: &ctx,
			},
			wantErr: false,
		},
		{
			name: "loader#prroperties_failed",
			args: args{
				path:      badAbsPath,
				targetPtr: &ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcl := &PropertiesConfigLoader{}
			if err := pcl.Load(tt.args.path, tt.args.targetPtr); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPropertiesConfigLoader_LoadMap(t *testing.T) {
	testFile := determineTestSourceFilePath()
	testdataDir := filepath.Dir(testFile)
	absPath := filepath.Clean(filepath.Join(testdataDir, "../../tests/testdata/application.properties"))

	ctx := make(map[string]any)
	pcl := NewPropertiesConfigLoader()
	if err := pcl.LoadMap(absPath, ctx); err != nil {
		t.Fatalf("LoadMap() error = %v", err)
	}

	if got, ok := mapz.NestedGet(ctx, "nemo.application.name"); !ok || got != "\"nemoapp\"" {
		t.Fatalf("expected nemo.application.name to be loaded, got value=%v ok=%v", got, ok)
	}
}
