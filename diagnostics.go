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
	"fmt"
	"strings"
)

func FormatBindDiagnostic(err error) string {
	if err == nil {
		return "Nemo Bind Diagnostic\n  Status: no error"
	}

	bindErr, ok := AsBindError(err)
	if !ok {
		return fmt.Sprintf("Nemo Bind Diagnostic\n  Error: %v", err)
	}

	lines := []string{
		"Nemo Bind Diagnostic",
		fmt.Sprintf("  Kind: %s", bindErr.Kind),
	}
	if bindErr.FieldPath != "" {
		lines = append(lines, fmt.Sprintf("  Field: %s", bindErr.FieldPath))
	}
	if bindErr.Key != "" {
		lines = append(lines, fmt.Sprintf("  Key: %s", bindErr.Key))
	}
	if bindErr.Cause != nil {
		lines = append(lines, fmt.Sprintf("  Cause: %v", bindErr.Cause))
	}
	if suggestion := bindSuggestion(bindErr.Kind); suggestion != "" {
		lines = append(lines, fmt.Sprintf("  Suggestion: %s", suggestion))
	}

	return strings.Join(lines, "\n")
}

func bindSuggestion(kind ErrorKind) string {
	switch kind {
	case InvalidTargetErrorKind:
		return "pass a pointer to a struct target"
	case InvalidTagErrorKind:
		return "check binder/required/default tags on the target field"
	case UnsettableFieldErrorKind:
		return "export the field or bind into a settable target"
	case MissingRequiredErrorKind:
		return "provide the missing property or remove required:\"true\""
	case UnsupportedTypeErrorKind:
		return "change the field to a supported bind type or add a custom conversion path"
	case ConversionFailedErrorKind:
		return "check the source value format and target field type"
	default:
		return ""
	}
}
