package helpers

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ToPtr[T any](v T) *T {
	return &v
}

type Response struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Token      string      `json:"token,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func BasicResponse(status bool, message string) Response {
	return Response{
		Status:  status,
		Message: message,
	}
}

func AuthResponseToken(status bool, message, token string) Response {
	return Response{
		Status:  status,
		Message: message,
		Token:   token,
	}
}

func SuccessResponseWithData(status bool, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func SuccessResponseWithDataPagination(status bool, message string, data interface{}, pagination interface{}) Response {
	return Response{
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

func ErrorResponseRequest(status bool, message string, err interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Errors:  err,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}

func BindAndValidate(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return fmt.Errorf("invalid request payload")
	}

	v := reflect.ValueOf(request).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		formVal := c.FormValue(field.Tag.Get("json"))
		if formVal == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Uint64:
			if parsed, err := strconv.ParseUint(formVal, 10, 64); err == nil {
				fieldValue.SetUint(parsed)
			}
		case reflect.Bool:
			if parsed, err := strconv.ParseBool(formVal); err == nil {
				fieldValue.SetBool(parsed)
			}
		}
	}

	if err := validate.Struct(request); err != nil {
		requestType := reflect.TypeOf(request).Elem()
		fieldMap := MapJSONFields(requestType)
		validationErrors := HandleValidationErrors(err, fieldMap)
		return &ValidationError{
			Message: "Bad Request",
			Errors:  validationErrors,
		}
	}

	return nil
}


// func BindFormWithFile[T any](ctx *fiber.Ctx, fileField string) (*T, string, error) {
// 	dto := new(T)

// 	// Parse semua field ke DTO
// 	if err := ctx.BodyParser(dto); err != nil {
// 		return nil, "", err
// 	}

// 	// Ambil file
// 	fileHeader, err := ctx.FormFile(fileField)
// 	if err != nil {
// 		// File tidak wajib
// 		return dto, "", nil
// 	}

// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer file.Close()

// 	// Upload ke Blob
// 	url, err := UploadFileToBlob(file, fileHeader.Filename)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return dto, url, nil
// }