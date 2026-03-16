/*
 * Copyright 漏 2023 the original author or authors.
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
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/mapz"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
	"github.com/photowey/nemo/pkg/valuez"
	"gopkg.in/ini.v1"
)

const (
	Ini = "ini"
)

const (
	iniStep     = 500
	iniPriority = ordered.HighPriority + iniStep*ordered.DefaultStep
)

var (
	iniSupportedConfigTypes = stringz.InitStringSlice(Ini)
)

var (
	_ ConfigLoader = (*IniConfigLoader)(nil)
)

func init() {
	Register(NewIniConfigLoader())
}

type IniConfigLoader struct{}

func NewIniConfigLoader() ConfigLoader {
	return &IniConfigLoader{}
}

func (icl *IniConfigLoader) Supports(strategy string) bool {
	return collection.ArrayContains(iniSupportedConfigTypes, strategy)
}

func (icl *IniConfigLoader) Order() int64 {
	return iniPriority
}

func (icl *IniConfigLoader) Name() string {
	return Ini
}

func (icl *IniConfigLoader) Load(path string, targetPtr any) error {
	if valuez.IsNil(targetPtr) {
		return fmt.Errorf("nemo: load ini config file, targetPtr can't be nil")
	}

	ctx := make(map[string]any)
	if err := icl.LoadMap(path, ctx); err != nil {
		return err
	}

	return mapstructure.Decode(ctx, targetPtr)
}

func (icl *IniConfigLoader) LoadMap(path string, ctx map[string]any) error {
	if valuez.IsNil(ctx) {
		return fmt.Errorf("nemo: load ini config file, ctx can't be nil")
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		for _, key := range section.Keys() {
			if sectionName == ini.DefaultSection {
				mapz.NestedSet(ctx, key.Name(), key.Value())
				continue
			}

			mapz.NestedSet(ctx, stringz.Concat(sectionName, stringz.Dot, key.Name()), key.Value())
		}
	}

	return nil
}
