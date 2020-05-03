package main

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		fmt.Println("ValidateStruct")
		if err := v.validate.Struct(obj); err != nil {
			fmt.Println("v.validate.Struct")
			// fmt.Println(err)

			errs := err.(validator.ValidationErrors)
			fmt.Println(errs)

			japanese := ja_JP.New()
			uni := ut.New(japanese, japanese)

			trans, _ := uni.GetTranslator("ja_JP")
			_ = trans.Add("ConfirmForm.Url", "フォームユーアルエル", false)
			_ = trans.Add("Url", "ユーアルエル", false)
			_ = trans.Add("Query", "クエリ", false)

			fmt.Println(errs[0].Translate(trans))

			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		trans := ut.New(ja_JP.New(), ja_JP.New())

		ja, _ := trans.GetTranslator("ja")
		_ = ja.Add("Url", "こんにちは", false)

		validate := validator.New()
		ja_translations.RegisterDefaultTranslations(validate, ja)

		// _ = trans.Add("Url", "ユーアルエル", false)
		// _ = ja_translations.RegisterDefaultTranslations(v.validate, trans)
		//

		// add any custom validations etc. here
	})
}

func TransFunc(ut ut.Translator, fe validator.FieldError) string {
	fld, _ := ut.T(fe.Field())
	t, err := ut.T(fe.Tag(), fld)
	if err != nil {
		return fe.(error).Error()
	}
	return t
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
