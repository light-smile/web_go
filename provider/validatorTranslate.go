package provider

import (
	"dnds_go/common"
	"dnds_go/global"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// InitTrans 初始化翻译器 对validate(请求参数)的校验结果进行翻译
func InitTrans() (err error) {
	locale := "zh"
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json 字段名 tag的自定义方法
		// 返回字段名优先级 name > json > 结构体字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
			if name != "" {
				return name
			}
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			return name
		})

		zhT := zh.New() // 中文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(zhT, zhT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// // 注册翻译器
		// switch locale {
		// case "en":
		// 	err = enTranslations.RegisterDefaultTranslations(v, trans)
		// case "zh":
		// 	err = zhTranslations.RegisterDefaultTranslations(v, trans)
		// default:
		// 	err = enTranslations.RegisterDefaultTranslations(v, trans)
		// }
		err = zhTranslations.RegisterDefaultTranslations(v, global.Trans)
		return
	}
	return
}

// 修改 beego 错误提示
var MessageTmpls = map[string]string{

	"Required": "不能为空",

	"Min": "最小值 为 %d",

	"Max": "最大值 为 %d",

	"Range": "范围 为 %d 到 %d",

	"MinSize": "最短长度 为 %d",

	"MaxSize": "最大长度 为 %d",

	"Length": "长度必须 为 %d",

	"Alpha": "必须是有效的字母",

	"Numeric": "必须是有效的数字",

	"AlphaNumeric": "必须是有效的字母或数字",

	"Match": "必须匹配 %s",

	"NoMatch": "必须不匹配 %s",

	"AlphaDash": "必须是有效的字母、数字或连接符号(-_)",

	"Email": "必须是有效的电子邮件地址",

	"IP": "必须是有效的IP地址",

	"Base64": "必须是有效的base64字符",

	"Mobile": "必须是有效的手机号码",

	"Tel": "必须是有效的电话号码",

	"Phone": "必须是有效的电话或移动电话号码",

	"ZipCode": "必须是有效的邮政编码",
}

// beego 的校验组件，可以不使用
type Bee struct {
	Id     int
	Name   string `  valid:"Required;Match(/^Bee.*/)"  chn:"手机号"` // Name 不能为空并且以 Bee 开头
	Age    int    `valid:"Range(1, 140)" chn:"年龄"`                // 1 <= Age <= 140，超出此范围即为不合法
	Email  string `valid:"Email; MaxSize(100)"`                   // Email 字段需要符合邮箱格式，并且最大长度不能大于 100 个字符
	Mobile string `json:"mobile" valid:"Mobile" chn:"手机号" `       // Mobile 必须为正确的手机号
	// IP     string `valid:"IP" `                               // IP 必须为一个正确的 IPv4 地址
}

// 校验规则为Match(正则)时 chn 无效
func BeeValidator(c *gin.Context, data interface{}) error {

	valid := validation.Validation{}
	validation.SetDefaultMessage(MessageTmpls)
	b, err := valid.Valid(data)
	if err != nil {
		// handle error
		common.BBadParam(c, err.Error())
		return err
	}
	// fmt.Println(b, "bbbb")
	// 参数校验失败的处理
	if !b {
		// fmt.Println(&valid.Errors)
		for _, err := range valid.Errors {
			st := reflect.TypeOf(Bee{})
			filed, _ := st.FieldByName(err.Field)
			var chn = filed.Tag.Get("chn")
			if chn == "" {
				common.BBadParam(c, err.Message)
				return errors.New(err.Message)
			}
			common.BBadParam(c, chn+err.Message)
			return errors.New(chn + err.Message)

		}
	}
	return nil
}
