package gin

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tzh "github.com/go-playground/validator/v10/translations/zh"
)

var (
	// 中文翻译
	zhTrans ut.Translator
	// 不支持
	errNotSupport = errors.New("not support")
)

func TranslateZH(err error) string {
	switch _err := err.(type) {
	case validator.ValidationErrors:
		var str []string
		for i := 0; i < len(_err); i++ {
			str = append(str, _err[i].Translate(zhTrans))
		}
		return strings.Join(str, ",")
	case binding.SliceValidationError:
		var str []string
		for i := 0; i < len(_err); i++ {
			if errs, ok := _err[i].(validator.ValidationErrors); ok {
				for i := 0; i < len(errs); i++ {
					str = append(str, errs[i].Translate(zhTrans))
				}
			}
		}
		return strings.Join(str, ",")
	case *json.UnmarshalTypeError:
		return fmt.Sprintf("%s字段类型不正确，应为%s", _err.Field, _err.Type.String())
	case *json.UnsupportedTypeError:
		return fmt.Sprintf("字段类型%s无法解析", _err.Type.String())
	default:
		return "未知错误"
	}
}

func ZH_CN(fieldLabel string, customRegister func(v *validator.Validate, t ut.Translator) error) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_zh := zh.New()
		_ut := ut.New(_zh, _zh)
		_t, _ := _ut.GetTranslator("zh")
		zhTrans = _t
		// 设置
		if fieldLabel != "" {
			v.RegisterTagNameFunc(func(field reflect.StructField) string {
				label := field.Tag.Get(fieldLabel)
				if label == "" {
					return field.Name
				}
				return label
			})
		}
		// 默认翻译
		if err := tzh.RegisterDefaultTranslations(v, _t); err != nil {
			return err
		}
		// 自定义
		if customRegister != nil {
			return customRegister(v, _t)
		}
		//
		return nil
	}
	return errNotSupport
}
