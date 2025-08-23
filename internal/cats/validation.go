package cats

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)
	for _, e := range errs {
		errorMap[e.Field()] = e.Tag()
	}
	return errorMap
}
