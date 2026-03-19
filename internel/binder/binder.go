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

package binder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/photowey/nemo/pkg/collection"
	"github.com/photowey/nemo/pkg/mapz"
	"github.com/photowey/nemo/pkg/stringz"
)

const (
	binderTag   = "binder"
	defaultTag  = "default"
	requiredTag = "required"
)

type ErrorKind string

const (
	InvalidTargetErrorKind    ErrorKind = "invalid_target"
	InvalidTagErrorKind       ErrorKind = "invalid_tag"
	UnsettableFieldErrorKind  ErrorKind = "unsettable_field"
	MissingRequiredErrorKind  ErrorKind = "missing_required"
	UnsupportedTypeErrorKind  ErrorKind = "unsupported_type"
	ConversionFailedErrorKind ErrorKind = "conversion_failed"
)

type BindError struct {
	Kind      ErrorKind
	FieldPath string
	Key       string
	Cause     error
}

func (e *BindError) Error() string {
	if e == nil {
		return "<nil>"
	}
	switch e.Kind {
	case InvalidTargetErrorKind:
		return "nemo: binder target must be a pointer to struct"
	case InvalidTagErrorKind:
		return fmt.Sprintf("nemo: field %s has invalid binding tag configuration: %v", e.FieldPath, e.Cause)
	case UnsettableFieldErrorKind:
		return fmt.Sprintf("nemo: field %s is not settable", e.FieldPath)
	case MissingRequiredErrorKind:
		return fmt.Sprintf("nemo: required bind key %q for field %s is missing", e.Key, e.FieldPath)
	case UnsupportedTypeErrorKind:
		return fmt.Sprintf("nemo: field %s with key %q uses an unsupported bind type: %v", e.FieldPath, e.Key, e.Cause)
	case ConversionFailedErrorKind:
		return fmt.Sprintf("nemo: bind field %s with key %q failed: %v", e.FieldPath, e.Key, e.Cause)
	default:
		if e.Cause != nil {
			return e.Cause.Error()
		}
		return "nemo: bind error"
	}
}

func (e *BindError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type Binder struct {
	Prefix string
}

func New() *Binder {
	return &Binder{}
}

func NewBinder(prefix string) *Binder {
	return &Binder{Prefix: prefix}
}

func (b *Binder) DefaultBind(target any, ctx collection.MixedMap) error {
	return b.Bind(b.Prefix, target, ctx)
}

func (b *Binder) Bind(prefix string, target any, ctx collection.MixedMap) error {
	if stringz.IsNotBlankString(prefix) && stringz.IsNotSuffix(prefix, stringz.Dot) {
		prefix += stringz.Dot
	}

	if target == nil {
		return &BindError{Kind: InvalidTargetErrorKind}
	}

	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return &BindError{Kind: InvalidTargetErrorKind}
	}

	return b.bindStruct(prefix, targetType.Elem(), targetValue.Elem(), ctx, targetType.Elem().Name())
}

func (b *Binder) bindStruct(prefix string, tt reflect.Type, tv reflect.Value, ctx collection.MixedMap, typePath string) error {
	if stringz.IsNotBlankString(prefix) && stringz.IsNotSuffix(prefix, stringz.Dot) {
		prefix += stringz.Dot
	}

	for i := 0; i < tt.NumField(); i++ {
		t := tt.Field(i)
		v := tv.Field(i)
		fieldPath := typePath + "." + t.Name

		tag := t.Tag.Get(binderTag)
		required, err := parseRequiredTag(t.Tag.Get(requiredTag), fieldPath)
		if err != nil {
			return err
		}
		defaultValue := t.Tag.Get(defaultTag)
		if stringz.IsBlankString(tag) {
			if required || defaultValue != "" {
				return &BindError{
					Kind:      InvalidTagErrorKind,
					FieldPath: fieldPath,
					Cause:     fmt.Errorf("constraints require a binder tag"),
				}
			}
			continue
		}
		if !v.CanSet() {
			return &BindError{Kind: UnsettableFieldErrorKind, FieldPath: fieldPath, Key: keyFor(prefix, tag)}
		}
		key := stringz.Concat(prefix, strings.ToLower(tag))

		switch {
		case t.Type.Kind() == reflect.Struct:
			if defaultValue != "" {
				return &BindError{
					Kind:      InvalidTagErrorKind,
					FieldPath: fieldPath,
					Key:       key,
					Cause:     fmt.Errorf("default tag is not supported on struct fields"),
				}
			}
			if required && !mapz.NestedContains(key, ctx) {
				return &BindError{Kind: MissingRequiredErrorKind, FieldPath: fieldPath, Key: key}
			}
			sub := reflect.New(t.Type).Interface()
			if err := b.bindStruct(key, t.Type, reflect.ValueOf(sub).Elem(), ctx, fieldPath); err != nil {
				return err
			}
			v.Set(reflect.ValueOf(sub).Elem())
		case t.Type.Kind() == reflect.Ptr:
			if err := b.bindPointerField(key, t.Type, v, ctx, fieldPath, required, defaultValue); err != nil {
				return err
			}
		default:
			value, ok := mapz.NestedGet(ctx, key)
			if !ok {
				if defaultValue != "" {
					value = defaultValue
					ok = true
				} else if required {
					return &BindError{Kind: MissingRequiredErrorKind, FieldPath: fieldPath, Key: key}
				}
			}
			if ok {
				converted, err := convertValue(value, t.Type)
				if err != nil {
					kind := ConversionFailedErrorKind
					if isUnsupportedTypeError(err) {
						kind = UnsupportedTypeErrorKind
					}
					return &BindError{Kind: kind, FieldPath: fieldPath, Key: key, Cause: err}
				}
				v.Set(converted)
			}
		}
	}

	return nil
}

func (b *Binder) bindPointerField(key string, targetType reflect.Type, fieldValue reflect.Value, ctx collection.MixedMap, fieldPath string, required bool, defaultValue string) error {
	elemType := targetType.Elem()
	if elemType.Kind() == reflect.Struct {
		if defaultValue != "" {
			return &BindError{
				Kind:      InvalidTagErrorKind,
				FieldPath: fieldPath,
				Key:       key,
				Cause:     fmt.Errorf("default tag is not supported on pointer-to-struct fields"),
			}
		}
		if required && !mapz.NestedContains(key, ctx) {
			return &BindError{Kind: MissingRequiredErrorKind, FieldPath: fieldPath, Key: key}
		}
		if !mapz.NestedContains(key, ctx) {
			return nil
		}

		sub := reflect.New(elemType)
		if err := b.bindStruct(key, elemType, sub.Elem(), ctx, fieldPath); err != nil {
			return err
		}
		fieldValue.Set(sub)
		return nil
	}

	value, ok := mapz.NestedGet(ctx, key)
	if !ok {
		if defaultValue != "" {
			value = defaultValue
			ok = true
		} else if required {
			return &BindError{Kind: MissingRequiredErrorKind, FieldPath: fieldPath, Key: key}
		}
	}
	if !ok {
		return nil
	}

	converted, err := convertValue(value, elemType)
	if err != nil {
		kind := ConversionFailedErrorKind
		if isUnsupportedTypeError(err) {
			kind = UnsupportedTypeErrorKind
		}
		return &BindError{Kind: kind, FieldPath: fieldPath, Key: key, Cause: err}
	}

	ptr := reflect.New(elemType)
	ptr.Elem().Set(converted)
	fieldValue.Set(ptr)
	return nil
}

func parseRequiredTag(raw, fieldPath string) (bool, error) {
	if raw == "" {
		return false, nil
	}
	required, err := strconv.ParseBool(raw)
	if err != nil {
		return false, &BindError{
			Kind:      InvalidTagErrorKind,
			FieldPath: fieldPath,
			Cause:     fmt.Errorf("required tag value %q is invalid", raw),
		}
	}
	return required, nil
}

func convertValue(value any, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	if targetType == reflect.TypeOf(time.Duration(0)) {
		switch current := value.(type) {
		case string:
			duration, err := time.ParseDuration(current)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(duration), nil
		case int64:
			return reflect.ValueOf(time.Duration(current)), nil
		case int:
			return reflect.ValueOf(time.Duration(current)), nil
		}
	}

	raw := reflect.ValueOf(value)
	if raw.Type().AssignableTo(targetType) {
		return raw, nil
	}
	if raw.Type().ConvertibleTo(targetType) {
		return raw.Convert(targetType), nil
	}

	if targetType.Kind() == reflect.Slice {
		return convertSliceValue(value, targetType)
	}

	str, ok := value.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("value of type %T can't be converted to %s", value, targetType)
	}

	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(str).Convert(targetType), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(str)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(str, 10, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(str, 10, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(str, targetType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	default:
		return reflect.Value{}, fmt.Errorf("value of type %T can't be converted to %s", value, targetType)
	}
}

func convertSliceValue(value any, targetType reflect.Type) (reflect.Value, error) {
	elemType := targetType.Elem()
	items := make([]any, 0)

	switch current := value.(type) {
	case []any:
		items = append(items, current...)
	case []string:
		for _, item := range current {
			items = append(items, item)
		}
	case string:
		if current == "" {
			return reflect.MakeSlice(targetType, 0, 0), nil
		}
		parts := strings.Split(current, ",")
		for _, part := range parts {
			items = append(items, strings.TrimSpace(part))
		}
	default:
		return reflect.Value{}, fmt.Errorf("value of type %T can't be converted to %s", value, targetType)
	}

	slice := reflect.MakeSlice(targetType, 0, len(items))
	for _, item := range items {
		converted, err := convertValue(item, elemType)
		if err != nil {
			return reflect.Value{}, err
		}
		slice = reflect.Append(slice, converted)
	}
	return slice, nil
}

func keyFor(prefix, tag string) string {
	return stringz.Concat(prefix, strings.ToLower(tag))
}

func isUnsupportedTypeError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "can't be converted")
}
