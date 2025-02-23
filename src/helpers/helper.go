package helpers

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Validator instance
var validate = validator.New()

// ValidateStruct memvalidasi struct dan mengembalikan map error jika ada
func ValidateStruct(s interface{}) map[string]string {
	errs := make(map[string]string)
	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = formatErrorMessage(e)
		}
	}
	return errs
}

// formatErrorMessage otomatis membuat pesan error berdasarkan tag validator
func formatErrorMessage(e validator.FieldError) string {
	return fmt.Sprintf("Field '%s' failed on validation '%s' (parameter: %s)", e.Field(), e.Tag(), e.Param())
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	AppName string      `json:"app_name"`
}

// AppName global constant
const AppName = "DockerControlGo"

// SuccessResponse untuk response sukses
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
		AppName: AppName,
	})
}

// ErrorResponse untuk response error
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	errorMsg := message
	// Jika err bertipe error, gunakan pesan error sebagai message
	var responseData interface{}
	if e, ok := err.(error); ok {
		errorMsg = e.Error()
		responseData = nil // Data di-null kan jika error berbentuk string
	} else {
		responseData = err // Jika bukan error, tetap gunakan err sebagai data
	}
	return c.Status(statusCode).JSON(Response{
		Code:    statusCode,
		Message: errorMsg,
		Data:    responseData,
		AppName: AppName,
	})
}

// MaskPrivateFields - Masking field yang bertanda `private:"true"`
func MaskPrivateFields(data interface{}, mask bool) interface{} {
	if !mask {
		return data // Jika mask false, kembalikan data apa adanya
	}

	val := reflect.ValueOf(data)

	// Handle jika data berupa slice (array of struct)
	if val.Kind() == reflect.Slice {
		maskedSlice := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			maskedSlice[i] = MaskPrivateFields(val.Index(i).Interface(), mask)
		}
		return maskedSlice
	}

	// Handle jika data berupa pointer ke struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer
	}

	// Jika bukan struct, langsung return
	if val.Kind() != reflect.Struct {
		return data
	}

	typ := val.Type()
	maskedData := make(map[string]interface{})

	// Loop semua field dalam struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Cek apakah field memiliki tag private
		privateTag := field.Tag.Get("private")
		if privateTag == "true" {
			maskedData[field.Name] = "******" // Masking field private
		} else {
			// Jika field adalah struct atau slice, proses secara rekursif
			if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Slice {
				maskedData[field.Name] = MaskPrivateFields(fieldValue.Interface(), mask)
			} else {
				maskedData[field.Name] = fieldValue.Interface()
			}
		}
	}

	return maskedData
}
