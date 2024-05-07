package utils

//
//import (
//	"github.com/labstack/echo/v4"
//	"giftcard/pkg/responser"
//	"giftcard/pkg/validator"
//)
//
//func ReadRequest[T any](c echo.Context) (req T, err error) {
//	var body T
//	if err := c.Bind(&body); err != nil {
//		return body, err
//		//responser.NewErrorBuilder().SetMessage("پیام قابل پردازش نیست.").SetStatusCode(422).Build().Respond(c)
//	}
//
//	if errRess := validator.ValidateRequestDto(c.Request().Context(), body); errRess != nil {
//		return body, errRess
//	}
//
//	return body, nil
//}
//
//func ReadRequestAndMapToEntity[T, K any](c echo.Context, fn func(T) K) (entity K, err error) {
//	var body T
//	if err := c.Bind(&body); err != nil {
//		return entity, responser.NewErrorBuilder().SetMessage("پیام قابل پردازش نیست.").SetStatusCode(422).Build()
//	}
//
//	if errRess := validator.ValidateRequestDto(c.Request().Context(), body); errRess != nil {
//		return entity, errRess
//	}
//
//	return fn(body), nil
//}
