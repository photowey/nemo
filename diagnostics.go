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
