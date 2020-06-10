package validate

import (
	"context"
	"reflect"
	"strings"

	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"

	// val "gopkg.in/go-playground/validator.v10"
	en_translations "gopkg.in/go-playground/validator.v10/translations/en"
)

var ()

//Validate is struct for Basic needs
type Validate struct {
	*ut.UniversalTranslator
	*validator.Validate
}

// var validate = validator.New()
var val Validate

//RegisterCustomTags needs to be called atleast once
func RegisterCustomTags() {
	en := en.New()
	val.UniversalTranslator = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := val.GetTranslator("en")

	val.Validate = validator.New()
	en_translations.RegisterDefaultTranslations(val.Validate, trans)

}

//GeneralFieldValidation is overriding default validation like : required, one of, numeric,email etc
//We overriding the message
func GeneralFieldValidation(trans ut.Translator) {

	val.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} can't be null!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	val.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} value should be email!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
	val.RegisterTranslation("numeric", trans, func(ut ut.Translator) error {
		return ut.Add("numeric", "{0} value should be number!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("numeric", fe.Field())
		return t
	})
	val.RegisterTranslation("one of", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} is failed to parse! Please use one of enumeration value!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("one of", fe.Field())
		return t
	})
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

}

//Struct Override from standard struct validation from go validate
func Struct(v interface{}) error {
	var trans, _ = val.GetTranslator("en")
	GeneralFieldValidation(trans)
	if err := val.StructCtx(context.Background(), v); err != nil {
		errs := err.(validator.ValidationErrors)
		var temp []string
		for _, e := range errs {
			temp = append(temp, fmt.Sprint(e.Translate(trans)))
		}
		err = fmt.Errorf(strings.Join(temp, "\n"))

		return err
	}
	return nil
}
