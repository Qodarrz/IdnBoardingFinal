package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Qodarrz/fiber-app/repository"
	"github.com/go-playground/validator/v10"
)

// CustomValidator struct untuk validasi kustom
type CustomValidator struct {
	Validator *validator.Validate
	valRepo   repository.ValidationRepositoryInterface
}

// InitCustomValidation menginisialisasi validasi kustom
func InitCustomValidation(valRepo repository.ValidationRepositoryInterface) *CustomValidator {
	v := validator.New()
	customValidator := &CustomValidator{
		Validator: v,
		valRepo:   valRepo,
	}
	v.RegisterValidationCtx("unique", customValidator.UniqueValidation)
	v.RegisterValidationCtx("in", customValidator.InValidation)
	return customValidator
}

// Implementasi metode Validate untuk memenuhi interface echo.Validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// UniqueValidation untuk memeriksa apakah nilai sudah ada di database
func (v *CustomValidator) UniqueValidation(ctx context.Context, fl validator.FieldLevel) bool {
	param := fl.Param()          // Contoh: "username.users"
	value := fl.Field().String() // Value dari field (misal: "john_doe")

	parent := fl.Parent()
	idField := parent.FieldByName("ID")

	var id uint32 = 0
	if idField.IsValid() {
		switch idField.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			id = uint32(idField.Uint())
		case reflect.String:
			// kalau id berupa string, bisa parsing
			idStr := idField.String()
			if idParsed, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				id = uint32(idParsed)
			}
		}
	}

	// Pisahkan parameter menjadi field dan table
	parts := strings.Split(param, ".")
	if len(parts) != 2 {
		return false // Format parameter salah
	}
	field := parts[0]
	table := parts[1]

	// Panggil repository untuk cek apakah data unique
	isUnique, err := v.valRepo.IsUnique(ctx, table, field, value, id)
	if err != nil {
		fmt.Println("Error checking uniqueness:", err)
		return false
	}

	return isUnique
}

func (v *CustomValidator) InValidation(ctx context.Context, fl validator.FieldLevel) bool {

	value := fl.Field().String()

	param := fl.Param()
	allowed := strings.Split(param, ".")

	for _, item := range allowed {
		if value == item {
			return true
		}
	}

	return false
}

func MapJSONFields(structType reflect.Type) map[string]string {
	jsonMap := make(map[string]string)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonName := field.Name

		if tag, ok := field.Tag.Lookup("json"); ok {
			jsonName = strings.Split(tag, ",")[0]
		}

		jsonMap[field.Name] = jsonName
	}

	return jsonMap
}

func HandleValidationErrors(err error, fieldMap map[string]string) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		jsonField := fieldMap[err.Field()]

		switch err.Tag() {
		case "required":
			errors[jsonField] = fmt.Sprintf("%s wajib diisi", jsonField)
		case "min":
			errors[jsonField] = fmt.Sprintf("%s minimal %s karakter", jsonField, err.Param())
		case "max":
			errors[jsonField] = fmt.Sprintf("%s maksimal %s karakter", jsonField, err.Param())
		case "unique":
			errors[jsonField] = fmt.Sprintf("%s sudah digunakan", jsonField)
		case "in":
			errors[jsonField] = fmt.Sprintf("%s format tidak valid", jsonField)
		case "email":
			errors[jsonField] = "format email tidak valid"
		default:
			errors[jsonField] = fmt.Sprintf("%s tidak valid", jsonField)
		}
	}

	return errors
}
