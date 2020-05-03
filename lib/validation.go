package lib

import (
	"reflect"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

type MyValidate struct {
	uni      *ut.UniversalTranslator
	Validate *validator.Validate
	trans    ut.Translator
}

func (v MyValidate) InitValidate() MyValidate {
	ja := ja.New()
	v.uni = ut.New(ja, ja)
	t, _ := v.uni.GetTranslator("ja")
	v.trans = t
	v.Validate = validator.New()
	v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		fieldName := fld.Tag.Get("jaFieldName")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})
	ja_translations.RegisterDefaultTranslations(v.Validate, v.trans)

	return v
}

// GetErrorMessages エラーメッセージ群の取得
func (v MyValidate) GetErrorMessages(err error) []string {
	if err == nil {
		return []string{}
	}
	var messages []string
	for _, m := range err.(validator.ValidationErrors).Translate(v.trans) {
		messages = append(messages, m)
	}
	return messages
}
