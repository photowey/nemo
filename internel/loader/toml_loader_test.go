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

package loader

import (
	"path/filepath"
	"testing"
)

func TestTomlConfigLoader_Load(t *testing.T) {
	testFile := determineTestSourceFilePath()
	testdataDir := filepath.Dir(testFile)

	absPath := filepath.Clean(filepath.Join(testdataDir, "../../tests/testdata/application.toml"))
	badAbsPath := filepath.Clean(filepath.Join(testdataDir, "../../tests/testdata/config.toml")) // not found

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
			name: "loader#toml_ok",
			args: args{
				path:      absPath,
				targetPtr: &ctx,
			},
			wantErr: false,
		},
		{
			name: "loader#toml_failed",
			args: args{
				path:      badAbsPath,
				targetPtr: &ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcl := &TomlConfigLoader{}
			if err := tcl.Load(tt.args.path, tt.args.targetPtr); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
