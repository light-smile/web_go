package testService

import (
	// "gin-skeleton/helper/response"
	"dnds_go/common"
	"dnds_go/global"
	"dnds_go/provider"
	testDao "dnds_go/src/dao/test"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 测试 数据库添加 返回数据库所有数据
func Create(c *gin.Context) {
	// var user testModle.User
	// if err := c.ShouldBindJSON(&user); err != nil {
	// 	common.BadParam(c, nil, "参数错误")
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// result, err := testDao.FindOrCreate(user)
	// if err != nil {
	// 	global.Logger.Sugar().Errorf("Create user err: %s", err.Error())
	// }

	// common.Succ(c, result, "success")
	res, err := testDao.Create()
	if err != nil {
		global.Logger.Error(err.Error())
	}
	common.Succ(c, res, "success")
}

// 测试配置功能是否正常
func TestConf(c *gin.Context) {
	data := viper.GetString("Server.port")

	common.Succ(c, data, "success")
}

// @Summary
// @Schemes
// @Description 返回hello 测试web框架
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /test/helle [get]
// 返回 hello

func Hello(c *gin.Context) {
	s := c.Query("text")
	fmt.Println(s)
	common.Succ(c, s, "success")
}

type TestMq struct {
}

// @Summary
// @Schemes
// @Description 测试mq 发布消息
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string} testMqtt
// @Router /test/mqSend [get]
// test Mq发布消息
func MqSend(c *gin.Context) {
	data := provider.Entity{
		Code: "0",
		Data: "testMqtt",
		Msg:  "success",
	}
	pubData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	provider.MqClient.Publish("test/hello", false, pubData)
	common.Succ(c, "testMqtt", "success")
}

type Userss struct {
	UserName string
	Sex      string
	Age      string
}

// 测试 导出excel功能
func ExportEexcel(c *gin.Context) {
	excel := common.ExcelInfo{
		Titles:   []string{"用户名", "性别", "年龄"},
		DataKeys: []string{"username", "sex", "age"},
		FileName: "dnds_go",
		ShowTime: true,
		SpecificColWidth: map[string]int{
			"用户名": 20,
		},
	}
	// data := common.ExcelData{
	// 	map[string]interface{}{"username": "海带", "sex": "19.90", "age": "123.png"},
	// 	map[string]interface{}{"username": "白菜", "sex": "9.90", "age": "logo-mate.png"},
	// 	map[string]interface{}{"username": "萝卜", "sex": "4.90", "age": "logo-mate.png"},
	// }
	dataSturct := []Userss{
		{
			UserName: "还带",
			Sex:      "123",
			Age:      "213",
		}, {
			UserName: "还带",
			Sex:      "123",
			Age:      "213",
		}, {
			UserName: "还带",
			Sex:      "123",
			Age:      "213",
		},
	}
	f, err := excel.CreateExcel(dataSturct)
	if err != nil {
		zap.S().Error(err)
	}
	_ = excel.DownloadExcel(c, f)
}

// 测试参数校验
type Users struct {
	FirstName string `json:"firstName" binding:"required"`

	LastName  string     `json:"lastName" binding:"required"`
	Age       uint8      `json:"age" binding:"gte=0,lte=130"`
	Email     string     `json:"email" binding:"required,email"`
	Addresses []*Address `json:"addresses" binding:"required,dive,required"`
}
type Address struct {
	Street string `json:"street" binding:"required"`
	City   string `json:"city" binding:"required" name:"城市"`
	Planet string `json:"planet" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
}

func (u *Users) valid(c *gin.Context) error {
	// if u.FirstName != "name" {
	// 	common.Fail(c, 2, "", "firstName != name")
	// 	return errors.New("firstName != name")
	// }
	return nil
}

type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func (s *SignUpParam) valid(c *gin.Context) error {
	if s.Age == 0 {
		common.Fail(c, common.ErrorBadParamCode, "", "年龄不能为空")
		return errors.New("年龄不能为空")
	}
	return nil
}

// @Summary
// @Schemes
// @Description 测试参数校验
// @Tags Test
// @Accept json
// @Produce json
// @Params Users body Users true "Users"
// @Success 200 {object} Users
// @Router /test/testValidator [post]
// gin 默认参数校验
func TestValidator(c *gin.Context) {
	var users Users
	// var u SignUpParam
	if err := c.ShouldBindJSON(&users); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			common.Fail(c, common.ErrorBadParamCode, nil, "参数格式有误")
			return
		}
		common.DefaultBadParam(c, errs)
		return
	}

	if ok := users.valid(c); ok != nil {
		return
	}
	common.Succ(c, users, "success")
}

type Bee struct {
	Id     int
	Name   string `chn:"手机号"  valid:"Required;Match(/^Bee.*/)" ` // Name 不能为空并且以 Bee 开头
	Age    int    `valid:"Range(1, 140)"`                        // 1 <= Age <= 140，超出此范围即为不合法
	Email  string `valid:"Email; MaxSize(100)"`                  // Email 字段需要符合邮箱格式，并且最大长度不能大于 100 个字符
	Mobile string `json:"mobile" valid:"Mobile" chn:"手机号" `      // Mobile 必须为正确的手机号
	// IP     string `valid:"IP" `                               // IP 必须为一个正确的 IPv4 地址
}

func (u *Bee) Valid(v *validation.Validation) {
	if strings.Index(u.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("Name", "名称里不能含有 admin")
	}
}

func TestBeeValidator(c *gin.Context) {
	var data Bee
	// if err := provider.BeeValidator(c, &data, Bee{}); err != nil {
	// 	return
	// }
	if err := c.ShouldBindJSON(&data); err != nil {
		common.Fail(c, common.ErrorBadParamCode, err, "参数绑定失败")
		return
	}

	if err := provider.BeeValidator(c, data); err != nil {
		return
	}

	common.Succ(c, "ok", "succ")
}

func GoroutineTest(c *gin.Context) {
	nums := c.Query("nums")
	num, _ := strconv.Atoi(nums)
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(10 * time.Second)
		}()
	}
	time.Sleep(10 * time.Second)
	fmt.Println(time.Now().UnixMicro(), num)
	common.Succ(c, "goroutine数量:"+nums, "success")
}
