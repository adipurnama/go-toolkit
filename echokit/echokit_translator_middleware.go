package echokit

import (
	"context"
	"strings"

	"github.com/adipurnama/go-toolkit/web"
	en_locale "github.com/go-playground/locales/en"
	id_locale "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	echo "github.com/labstack/echo/v4"
)

// Validate validates request body for incoming echo.Context request
// returns web.HTTPError contains field errors (if any).
func Validate(ctx echo.Context, req interface{}) *web.HTTPError {
	if err := ctx.Validate(req); err != nil {
		return web.NewHTTPValidationError(ctx.Request().Context(), err)
	}

	return nil
}

// ValidatorTranslatorMiddleware adds request body validator's translator
// based on 'Accept-Lang' header.
// currently only suuports ID & EN locale.
func ValidatorTranslatorMiddleware(v *validator.Validate) echo.MiddlewareFunc {
	en := en_locale.New()
	id := id_locale.New()
	uni := ut.New(en, en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	_ = en_translations.RegisterDefaultTranslations(v, transEN)
	_ = id_translations.RegisterDefaultTranslations(v, transID)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			trans := transEN

			lang := ctx.Request().Header.Get("Accept-Lang")
			if strings.ToLower(lang) == "id" ||
				strings.ToLower(lang) == "id-id" {
				trans = transID
			}

			rCtx := context.WithValue(ctx.Request().Context(), web.ContextKeyTranslator, trans)
			req := ctx.Request().WithContext(rCtx)
			ctx.SetRequest(req)

			return next(ctx)
		}
	}
}
