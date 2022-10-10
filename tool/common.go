/*
 * @Author: jlz
 * @Date: 2022-08-16 17:32:40
 * @LastEditTime: 2022-08-31 15:20:46
 * @LastEditors: jlz
 * @Description:
 */

package tool

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
)

// 将 任意类型的切片转换为 []interface 类型的切片
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := IsSlice(slice)
	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// 判断interface 是否是切片
func IsSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == reflect.Slice {
		ok = true
	}
	return
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func NewUuid() (uu string) {
	u := uuid.NewV4()
	return u.String()
}

// 字符串左边添加%
func LeftLike(str string) string {
	return "%" + str
}

// 字符串右边添加%
func RightLike(str string) string {
	return str + "%"
}

// 字符串左右添加%
func LikeStr(str string) string {
	return "%" + str + "%"
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
