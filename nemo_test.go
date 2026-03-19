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

package nemo

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRootPackageNewAndBind(t *testing.T) {
	type DemoConfig struct {
		Timeout time.Duration `binder:"timeout" default:"3s"`
		Enabled bool          `binder:"enabled" required:"true"`
	}

	env := New().(*StandardEnvironment)
	if err := env.Start(
		WithProperties(MixedMap{
			"demo": MixedMap{
				"enabled": "true",
			},
		}),
	); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	cfg, err := Bind[DemoConfig](env, "demo")
	if err != nil {
		t.Fatalf("Bind[T]() error = %v", err)
	}
	if !cfg.Enabled || cfg.Timeout != 3*time.Second {
		t.Fatalf("expected bound config, got %+v", cfg)
	}
}

func TestRootPackageBindErrorHelpers(t *testing.T) {
	type DemoConfig struct {
		Enabled bool `binder:"enabled" required:"true"`
	}

	env := New().(*StandardEnvironment)
	if err := env.Start(WithProperties(MixedMap{"demo": MixedMap{}})); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	_, err := Bind[DemoConfig](env, "demo")
	if err == nil {
		t.Fatalf("expected bind error")
	}

	bindErr, ok := AsBindError(err)
	if !ok {
		t.Fatalf("expected AsBindError to unwrap bind error")
	}
	if bindErr.Kind != MissingRequiredErrorKind {
		t.Fatalf("expected missing required error kind, got %s", bindErr.Kind)
	}
	if !IsBindErrorKind(err, MissingRequiredErrorKind) {
		t.Fatalf("expected IsBindErrorKind to detect missing required error kind")
	}

	var typed *BindError
	if !errors.As(err, &typed) {
		t.Fatalf("expected errors.As to unwrap bind error")
	}

	diagnostic := FormatBindDiagnostic(err)
	if !strings.Contains(diagnostic, "MissingRequired") && !strings.Contains(diagnostic, string(MissingRequiredErrorKind)) {
		t.Fatalf("expected diagnostic to contain bind error kind, got %q", diagnostic)
	}
}
