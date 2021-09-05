package controllers

// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"myapp/service"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	trans                ut.Translator
	globalTagTranslation = map[string]string{
		"required_with": "{0} is required if {1} available",
	}
)

func init() {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	binding.Validator = new(defaultValidator)
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	var errMsgs []string
	for i, e := range err {
		if e == nil {
			continue
		}
		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}
	return strings.Join(errMsgs, "\n")
}

var _ binding.StructValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	v.lazyinit()
	return bindError(v.validate.Struct(obj))
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		en_translations.RegisterDefaultTranslations(v.validate, trans)
		validateAddCustomTagTranslation(v.validate)
	})
}

func bindError(err error) error {
	if err == nil {
		return err
	}
	validationErrors := err.(validator.ValidationErrors)
	enErr := fmt.Errorf("%+v", validationErrors[0].Translate(trans))
	return enErr
}

func errH(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}

func validateAddCustomTagTranslation(validate *validator.Validate) {
	for key, val := range globalTagTranslation {
		validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
			return ut.Add(key, val, true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(key, fe.Field(), fe.Param())

			return t
		})
	}
}

func userData(c *gin.Context) (map[string]interface{}, error) {
	tokenData := service.CtxVal(c.Request.Context())
	if tokenData == nil {
		return map[string]interface{}{
			"is_login": false,
		}, nil
	}

	getUser, err := service.UserGetByID(c.Request.Context(), tokenData.ID)
	if err != nil {
		return map[string]interface{}{
			"is_login": false,
		}, err
	}

	return map[string]interface{}{
		"is_login": true,
		"data":     getUser,
	}, nil
}
