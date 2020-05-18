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
func (v MyValidate) GetErrorMessages(err error) map[string]string {
	if err == nil {
		return map[string]string{}
	}
	messages := map[string]string{}
	for k, m := range err.(validator.ValidationErrors).Translate(v.trans) {
		messages[k] = m
	}
	return messages
}

// GetErrorMessage エラーメッセージの取得
func (v MyValidate) GetErrorMessage(err error, key string) string {
	if err == nil {
		return ""
	}
	return err.(validator.ValidationErrors).Translate(v.trans)[key]
}

// hasErrorMessage エラーメッセージの存在チェック
func (v MyValidate) HasErrorMessage(err error, key string) bool {
	return v.GetErrorMessage(err, key) != ""
}

// hasErrorMessage エラーメッセージの存在チェック
func (v MyValidate) PushErrorMessage(m map[string]string, key string, message string) map[string]string {
	if m == nil {
		m = map[string]string{}
	}
	m[key] = message
	return m
}
