package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator is a wrapper around the validator.Validate instance
type CustomValidator struct {
	validator  *govalidator.Validate
	translator *ErrorTranslator
}

// ErrorTranslator handles translating validation errors to custom messages
type ErrorTranslator struct {
	fieldErrors     map[string]map[string]string // map[field]map[tag]message
	leafFieldErrors map[string]map[string]string // map[leafName]map[tag]message for fallback
	defaultErrors   map[string]string            // map[tag]message
	defaultMessage  string
}

// ValidationContext stores context-specific validation settings
type ValidationContext struct {
	validator  *CustomValidator
	translator *ErrorTranslator
}

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   any    `json:"value,omitempty"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// Error implements the error interface
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("validation failed: ")
	for i, err := range v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return sb.String()
}

// AsMap converts ValidationErrors to a map for consistent error formatting
func (v ValidationErrors) AsMap() map[string]interface{} {
	return FormatValidationErrors(v)
}

// normalizePath standardizes field paths for consistent format handling
// - Handles array indices consistently: profile.addresses[0].street
// - Handles map keys: profile.metadata["key"].value
// The path is normalized to ensure consistent lookup regardless of source format
func normalizePath(path string) string {
	// Trim leading/trailing spaces
	path = strings.TrimSpace(path)

	// Handle empty path case
	if path == "" {
		return ""
	}

	// Replace spaces around dots with just dots
	path = strings.ReplaceAll(path, " .", ".")
	path = strings.ReplaceAll(path, ". ", ".")

	// Split by dots to handle each segment
	segments := strings.Split(path, ".")
	normalizedSegments := make([]string, 0, len(segments))

	for _, segment := range segments {
		// Trim spaces from segment
		segment = strings.TrimSpace(segment)

		// Skip empty segments
		if segment == "" {
			continue
		}

		// Normalize array/map notation if present
		if strings.Contains(segment, "[") && strings.Contains(segment, "]") {
			// Extract the field name part (before the bracket)
			fieldName := segment
			if idx := strings.Index(segment, "["); idx > 0 {
				fieldName = strings.TrimSpace(segment[:idx])
			}

			// Extract all index/key parts but normalize them
			normalizedIndices := make([]string, 0)
			remaining := segment
			for strings.Contains(remaining, "[") && strings.Contains(remaining, "]") {
				start := strings.Index(remaining, "[")
				end := strings.Index(remaining, "]")

				if start >= 0 && end > start {
					// Extract the content between brackets
					indexContent := strings.TrimSpace(remaining[start+1 : end])
					// Create normalized index with no spaces
					normalizedIndex := "[" + indexContent + "]"
					normalizedIndices = append(normalizedIndices, normalizedIndex)

					// Move past this index/key
					if end+1 < len(remaining) {
						remaining = remaining[end+1:]
					} else {
						remaining = ""
					}
				} else {
					break
				}
			}

			// Reconstruct the segment with normalized field name and indices
			normalized := fieldName + strings.Join(normalizedIndices, "")
			normalizedSegments = append(normalizedSegments, normalized)
		} else {
			// Regular field segment without array/map notation
			normalizedSegments = append(normalizedSegments, segment)
		}
	}

	// Join segments with dots
	return strings.Join(normalizedSegments, ".")
}

// extractLeafName extracts just the leaf field name from a path
// For example:
// - "profile.firstName" returns "firstName"
// - "items[0].name" returns "name"
// - "data.points[0][1]" returns "points[0][1]"
func extractLeafName(path string) string {
	path = strings.TrimSpace(path)

	if path == "" {
		return ""
	}

	// If there are no dots, it's already a leaf
	if !strings.Contains(path, ".") {
		return path
	}

	// Split by dots and take the last segment
	segments := strings.Split(path, ".")
	return segments[len(segments)-1]
}

// NewErrorTranslator creates a new ErrorTranslator
func NewErrorTranslator() *ErrorTranslator {
	return &ErrorTranslator{
		fieldErrors:     make(map[string]map[string]string),
		leafFieldErrors: make(map[string]map[string]string),
		defaultErrors:   make(map[string]string),
		defaultMessage:  "Invalid value",
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
// The field path is normalized for consistent lookup during validation
func (t *ErrorTranslator) SetFieldError(field, tag, message string) {
	// Normalize the field path for consistent lookup
	normalizedPath := normalizePath(field)

	// Store by normalized path
	if _, ok := t.fieldErrors[normalizedPath]; !ok {
		t.fieldErrors[normalizedPath] = make(map[string]string)
	}
	t.fieldErrors[normalizedPath][tag] = message

	// Also store by leaf name for fallback lookups
	leafName := extractLeafName(field)
	if leafName != "" && leafName != normalizedPath {
		if _, ok := t.leafFieldErrors[leafName]; !ok {
			t.leafFieldErrors[leafName] = make(map[string]string)
		}
		t.leafFieldErrors[leafName][tag] = message
	}
}

// SetDefaultError sets a default error message for a validation tag
func (t *ErrorTranslator) SetDefaultError(tag, message string) {
	t.defaultErrors[tag] = message
}

// SetDefaultMessage sets the default error message for all validations
func (t *ErrorTranslator) SetDefaultMessage(message string) {
	t.defaultMessage = message
}

// Translate translates a validation error to a custom message
func (t *ErrorTranslator) Translate(field string, tag string) string {
	// Normalize the field path for consistent lookup
	normalizedPath := normalizePath(field)

	// First, check if there's a custom message for the normalized full field path and tag
	if fieldMessages, ok := t.fieldErrors[normalizedPath]; ok {
		if message, ok := fieldMessages[tag]; ok {
			return message
		}
	}

	// If no normalized path match, check if there's a leaf field name match
	leafName := extractLeafName(field)
	if leafName != "" && leafName != normalizedPath {
		// Check in the dedicated leaf field errors map first
		if fieldMessages, ok := t.leafFieldErrors[leafName]; ok {
			if message, ok := fieldMessages[tag]; ok {
				return message
			}
		}

		// For backward compatibility, also check in the main fieldErrors map
		if fieldMessages, ok := t.fieldErrors[leafName]; ok {
			if message, ok := fieldMessages[tag]; ok {
				return message
			}
		}
	}

	// Check if there's a default message for this tag
	if message, ok := t.defaultErrors[tag]; ok {
		return message
	}

	// Return the default message
	return t.defaultMessage
}

// NewCustomValidator creates a new CustomValidator instance
func NewCustomValidator() *CustomValidator {
	v := govalidator.New()
	// Register function to get field name from json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	// Register custom validation tag for passwords
	v.RegisterValidation("password", validatePassword)

	translator := NewErrorTranslator()
	setDefaultMessages(translator)

	return &CustomValidator{
		validator:  v,
		translator: translator,
	}
}

// Validate validates the provided struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validationErrors := ValidationErrors{}

		for _, err := range err.(govalidator.ValidationErrors) {
			tag := err.Tag()

			// Get the JSON field name path
			jsonFieldPath := getJSONFieldPath(i, err)

			// Translate the error using the full path - our enhanced translator
			// will handle both full path and leaf name lookups
			message := cv.translator.Translate(jsonFieldPath, tag)

			// Create a validation error
			validationError := ValidationError{
				Field:   jsonFieldPath,
				Message: message,
				Tag:     tag,
				Value:   err.Value(),
			}

			validationErrors = append(validationErrors, validationError)
		}

		return validationErrors
	}

	return nil
}

// getJSONFieldPath returns the JSON field path for a validation error
// For nested fields, it returns a dot-separated path like "profile.firstName"
// For array/slice fields, it returns indexed paths like "items[0].name"
func getJSONFieldPath(obj interface{}, fieldError govalidator.FieldError) string {
	// Build the namespace path based on JSON field names rather than struct field names
	namespace := fieldError.Namespace()
	parts := strings.Split(namespace, ".")

	// The first part is the type name, so we skip it
	parts = parts[1:]

	// Build a new path with JSON names
	var jsonParts []string
	currentObj := obj
	currentType := reflect.TypeOf(currentObj).Elem()
	currentValue := reflect.ValueOf(currentObj).Elem()

	for i, part := range parts {
		// Check if this part refers to an array/slice index
		indexMatch := regexp.MustCompile(`^(\w+)\[(\d+)\]$`).FindStringSubmatch(part)

		if len(indexMatch) == 3 {
			// This is an array/slice index reference like 'Items[0]'
			fieldName := indexMatch[1]
			indexStr := indexMatch[2]
			index, _ := strconv.Atoi(indexStr)

			// Find the struct field for the array/slice
			field, found := currentType.FieldByName(fieldName)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name for the array/slice field
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				jsonName = fieldName
			}

			// Add the field name and index to the path
			jsonParts = append(jsonParts, fmt.Sprintf("%s[%d]", jsonName, index))

			// Update currentType/currentValue for the next iteration
			fieldValue := currentValue.FieldByName(fieldName)
			if !fieldValue.IsValid() || index >= fieldValue.Len() {
				// If the field value is invalid or index is out of bounds, we can't continue
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}

			// Get the element at the specified index
			elemValue := fieldValue.Index(index)

			// Update currentType and currentValue based on the element type
			if elemValue.Kind() == reflect.Struct {
				currentType = elemValue.Type()
				currentValue = elemValue
			} else if elemValue.Kind() == reflect.Ptr && elemValue.Elem().Kind() == reflect.Struct {
				currentType = elemValue.Elem().Type()
				currentValue = elemValue.Elem()
			} else {
				// If the element is not a struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		} else {
			// Regular struct field
			field, found := currentType.FieldByName(part)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				// If there's no JSON tag or it's "-", use the original field name
				jsonParts = append(jsonParts, part)
			} else {
				jsonParts = append(jsonParts, jsonName)
			}

			// Update currentType and currentValue for the next iteration
			fieldValue := currentValue.FieldByName(part)

			if field.Type.Kind() == reflect.Struct {
				currentType = field.Type
				if fieldValue.IsValid() {
					currentValue = fieldValue
				}
			} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				currentType = field.Type.Elem()
				if fieldValue.IsValid() && !fieldValue.IsNil() {
					currentValue = fieldValue.Elem()
				}
			} else {
				// If the field is not a struct or pointer to struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		}
	}

	return strings.Join(jsonParts, ".")
}

// Translator gets the CustomValidator's ErrorTranslator instance
func (cv *CustomValidator) Translator() *ErrorTranslator {
	return cv.translator
}

// Clone creates a copy of the ErrorTranslator
func (t *ErrorTranslator) Clone() *ErrorTranslator {
	clone := NewErrorTranslator()
	clone.defaultMessage = t.defaultMessage

	// Copy default errors
	for tag, msg := range t.defaultErrors {
		clone.defaultErrors[tag] = msg
	}

	// Copy field errors
	for field, tagMsgs := range t.fieldErrors {
		clone.fieldErrors[field] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.fieldErrors[field][tag] = msg
		}
	}

	// Copy leaf field errors
	for field, tagMsgs := range t.leafFieldErrors {
		clone.leafFieldErrors[field] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.leafFieldErrors[field][tag] = msg
		}
	}

	return clone
}

// Clone creates a copy of the CustomValidator with a new translator
func (cv *CustomValidator) Clone() *CustomValidator {
	return &CustomValidator{
		validator:  cv.validator,
		translator: cv.translator.Clone(),
	}
}

// NewValidationContext creates a new validation context with a cloned validator
func NewValidationContext(baseValidator *CustomValidator) *ValidationContext {
	cloned := baseValidator.Clone()
	return &ValidationContext{
		validator:  cloned,
		translator: cloned.translator,
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
func (vc *ValidationContext) SetFieldError(field, tag, message string) *ValidationContext {
	vc.translator.SetFieldError(field, tag, message)
	return vc
}

// SetDefaultError sets a default error message for a validation tag
func (vc *ValidationContext) SetDefaultError(tag, message string) *ValidationContext {
	vc.translator.SetDefaultError(tag, message)
	return vc
}

// SetDefaultMessage sets the default error message for all validations
func (vc *ValidationContext) SetDefaultMessage(message string) *ValidationContext {
	vc.translator.SetDefaultMessage(message)
	return vc
}

// Validate validates the provided struct using this context's validator
func (vc *ValidationContext) Validate(i interface{}) error {
	return vc.validator.Validate(i)
}

// BindAndValidate binds and validates a request body to a struct
func BindAndValidate(c echo.Context, i interface{}) error {
	// Bind the request body to the struct
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate the struct
	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

// BindAndValidateWithContext binds and validates a request body to a struct using a custom validation context
func BindAndValidateWithContext(c echo.Context, i interface{}, vc *ValidationContext) error {
	// Bind the request body to the struct
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate the struct using the context's validator
	if err := vc.Validate(i); err != nil {
		return err
	}

	return nil
}
