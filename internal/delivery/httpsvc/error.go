package httpsvc

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"net/http"
)

// http errors
var (
	ErrInvalidArgument            = echo.NewHTTPError(http.StatusBadRequest, "invalid argument")
	ErrNotEnoughBalance           = echo.NewHTTPError(http.StatusBadRequest, "balance not enough")
	ErrNotFound                   = echo.NewHTTPError(http.StatusNotFound, "record not found")
	ErrInternal                   = echo.NewHTTPError(http.StatusInternalServerError, "internal system error")
	ErrEntityTooLarge             = echo.NewHTTPError(http.StatusRequestEntityTooLarge, "entity too large")
	ErrUnauthenticated            = echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	ErrUnauthorized               = echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrEmailTokenNotMatch         = echo.NewHTTPError(http.StatusUnauthorized, "email or token not match")
	ErrEmailPasswordNotMatch      = echo.NewHTTPError(http.StatusUnauthorized, "email or password not match")
	ErrPermissionDenied           = echo.NewHTTPError(http.StatusForbidden, "permission denied")
	ErrSpaceNotSelected           = echo.NewHTTPError(http.StatusBadRequest, "space not selected")
	ErrLoginByEmailPasswordLocked = echo.NewHTTPError(http.StatusLocked, "user is locked from logging in using email and password")
	ErrInvitationExpired          = echo.NewHTTPError(http.StatusBadRequest, "invitation expired")
	ErrFailedPrecondition         = echo.NewHTTPError(http.StatusPreconditionFailed, "precondition failed")
)

//// httpValidationOrInternalErr return valdiation or internal error
//func httpValidationOrInternalErr(err error) error {
//	switch t := err.(type) {
//	case validator.ValidationErrors:
//		_ = t
//		errVal := err.(validator.ValidationErrors)
//
//		fields := map[string]interface{}{}
//		for _, ve := range errVal {
//			fields[ve.Field()] = fmt.Sprintf("Failed on the '%s' tag", ve.Tag())
//		}
//
//		return echo.NewHTTPError(http.StatusBadRequest, utils.Dump(fields))
//	default:
//		return ErrInternal
//	}
//}

// httpValidationOrInternalErr return valdiation or internal error
func httpValidationOrInternalErr(err error) error {
	// Memeriksa apakah err merupakan validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Jika tidak ada kesalahan validasi, mengembalikan kesalahan internal
		return ErrInternal
	}

	// Mengubah validator.ValidationErrors menjadi map dengan kunci field dan nilai pesan kesalahan
	fields := make(map[string]string)
	for _, validationError := range validationErrors {
		// Mengambil tag yang digunakan untuk menyebabkan kesalahan validasi
		tag := validationError.Tag()
		// Menambahkan kesalahan validasi ke map
		fields[validationError.Field()] = fmt.Sprintf("Failed on the '%s' tag", tag)
	}

	// Mengembalikan kesalahan HTTP dengan status bad request dan daftar field yang gagal validasi
	return echo.NewHTTPError(http.StatusBadRequest, fields)
}
