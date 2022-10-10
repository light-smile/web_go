package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName string     `validate:"required"`
	LastName  string     `validate:"required"`
	Age       uint8      `validate:"gte=0,lte=130"`
	Email     string     `validate:"required,email"`
	Addresses []*Address `validate:"required,dive,required"`
}
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func TestString(t *testing.T) {
	validate := validator.New()

	var boolTest bool

	err := validate.Var(boolTest, "required")
	if err != nil {
		fmt.Println(err)
	}
	var stringTest string = ""
	err = validate.Var(stringTest, "required")
	if err != nil {
		fmt.Println(err)
	}
}

func TestStruct(t *testing.T) {
	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName: "Badger",
		LastName:  "Smith",
		Age:       135,
		Email:     "Badger.Smith@gmail.com",
		Addresses: []*Address{address},
	}

	validate := validator.New()
	jb, _ := json.Marshal(&user)

	fmt.Println(string(jb))
	err := validate.Struct(user)
	if err != nil {
		fmt.Println("=== error msg ====")
		fmt.Println(err)

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		fmt.Println("\r\n=========== error field info ====================")
		for _, err := range err.(validator.ValidationErrors) {
			// 列出效验出错字段的信息
			fmt.Println("Namespace: ", err.Namespace())
			fmt.Println("Fild: ", err.Field())
			fmt.Println("StructNamespace: ", err.StructNamespace())
			fmt.Println("StructField: ", err.StructField())
			fmt.Println("Tag: ", err.Tag())
			fmt.Println("ActualTag: ", err.ActualTag())
			fmt.Println("Kind: ", err.Kind())
			fmt.Println("Type: ", err.Type())
			fmt.Println("Value: ", err.Value())
			fmt.Println("Param: ", err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}
}

type Users struct {
	Name string `form:"name" json:"name" validate:"required,CustomerValidation"` //注意：required和CustomerValidation之间不能有空格，否则panic。CustomerValidation：自定义tag-函数标签
	Age  uint8  ` form:"age" json:"age" validate:"gte=0,lte=80"`                 //注意：gte=0和lte=80之间不能有空格，否则panic
}

var validate *validator.Validate

func TestCustomValidator(t *testing.T) {
	validate = validator.New()
	validate.RegisterValidation("CustomerValidation", CustomerValidationFunc) //注册自定义函数，前一个参数是struct里tag自定义，后一个参数是自定义的函数

	user := &Users{
		Name: "jimmy",
		Age:  86,
	}

	fmt.Println("first value: ", user)
	err := validate.Struct(user)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}

	user.Name = "tom"
	user.Age = 29
	fmt.Println("second value: ", user)
	err = validate.Struct(user)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
}

// 自定义函数
func CustomerValidationFunc(f1 validator.FieldLevel) bool {
	// f1 包含了字段相关信息
	// f1.Field() 获取当前字段信息
	// f1.Param() 获取tag对应的参数
	// f1.FieldName() 获取字段名称
	fmt.Println(f1.Field().String(), ":field")
	fmt.Printf("%v, %s", f1, ":f1")
	fmt.Println(f1.FieldName(), ":FieldName")

	return f1.Field().String() == "jimmy"
}
