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
	"fmt"

	"github.com/magiconair/properties"
	"github.com/mitchellh/mapstructure"
	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/ordered"
	"github.com/photowey/nemo/pkg/stringz"
	"github.com/photowey/nemo/pkg/valuez"
)

// xxx.properties

const (
	Properties = "properties"
)

const (
	propertiesStep     = 1000
	propertiesPriority = ordered.HighPriority + propertiesStep*ordered.DefaultStep
)

var (
	propertiesSupportedConfigTypes = stringz.InitStringSlice(Properties)
)

var (
	_ ConfigLoader = (*PropertiesConfigLoader)(nil)
)

func init() {
	Register(NewPropertiesConfigLoader())
}

type PropertiesConfigLoader struct {
}

func NewPropertiesConfigLoader() ConfigLoader {
	return &PropertiesConfigLoader{}
}

func (pcl *PropertiesConfigLoader) Supports(strategy string) bool {
	return collection.ArrayContains(propertiesSupportedConfigTypes, strategy)
}

func (pcl *PropertiesConfigLoader) Order() int64 {
	return propertiesPriority
}

func (pcl *PropertiesConfigLoader) Name() string {
	return Properties
}

func (pcl *PropertiesConfigLoader) Load(path string, targetPtr any) error {
	if valuez.IsNil(targetPtr) {
		return fmt.Errorf("nemo: load properties config file, targetPtr can't be nil")
	}

	ctx := make(map[string]any)
	err := pcl.LoadMap(path, ctx)
	if err != nil {
		return err
	}

	return mapstructure.Decode(ctx, targetPtr)
}

func (pcl *PropertiesConfigLoader) LoadMap(path string, ctx map[string]any) error {
	if valuez.IsNil(ctx) {
		return fmt.Errorf("nemo: load properties config file, ctx can't be nil")
	}

	ppt, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return err
	}

	for _, key := range ppt.Keys() {
		value, _ := ppt.Get(key)
		ctx[key] = value
	}

	return nil
}
