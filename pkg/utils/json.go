package utils

//
//import (
//	"encoding/json"
//	"errors"
//
//	"github.com/go-playground/validator/v10"
//
//	"micro/pkg/rest"
//)
//
//type JsonError struct {
//	FieldName    string `json:"fieldName,omitempty"`
//	Message      string `json:"message,omitempty"`
//	ExpectedType string `json:"expectedType,omitempty"`
//	ReceivedType string `json:"ReceivedType,omitempty"`
//	err          error
//}
//
//const (
//	InvalidType  = "تایپ نامعتبر است"
//	InvalidInput = "ورودی نامعتبر است"
//)
//
//func JsonErrorHandler(err error) rest.ApiResponse {
//	if err == nil {
//		return nil
//	}
//
//	var unmarshalTypeError *json.UnmarshalTypeError
//	if errors.As(err, &unmarshalTypeError) {
//		return rest.BadRequest(
//			100400,
//			InvalidType,
//			rest.FA).AddDetails(unmarshalTypeError.Field, unmarshalTypeError.Value).
//			AddCause(err.Error())
//	}
//
//	var validationErrors *validator.ValidationErrors
//
//	if errors.As(err, &validationErrors) {
//		res := rest.BadRequest(
//			100400,
//			InvalidInput,
//			rest.FA).AddDetails(unmarshalTypeError.Field, unmarshalTypeError.Value)
//		for _, fieldError := range *validationErrors {
//			res = res.AddDetails(fieldError.Field(), fieldError.Value())
//		}
//
//		res = res.AddCause(err.Error())
//		return res
//	}
//
//	return rest.BadRequest(
//		100400,
//		InvalidInput,
//		rest.FA,
//	)
//}
//
//func (e JsonError) Error() string {
//	return e.err.Error()
//}
