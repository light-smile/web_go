/*
 * @Author: jlz
 * @Date: 2022-08-31 14:56:45
 * @LastEditTime: 2022-08-31 15:46:16
 * @LastEditors: jlz
 * @Description: 将以空格为分隔的字符串单词首字母转换为大写
 */

package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestCapitalize(t *testing.T) {
	Capitalize(
		`id
		gatewayName
	gatewayAddress 
	siteName
	areaIds
	siteType
	siteLevel
	siteVolt
	siteBuild
	siteBDZ
	stationLine
	longitude
	latitude
	inlines
	outlines
	dev0ps
	reSgAtt1
	reSgAtt2
	reSgAtt3
	reSgAtt4
	reCoAtt1
	reCoAtt2
	reCoAtt3
	reCoAtt4`,
		`String
	String
	String
	String
	String
	String
	String
	String
	String
	String
	String
	String
	String
	Int
	Int
	String
	String
	String
	String
	String
	String
	String
	String
	String`,
	)
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

func Capitalize(name string, nameType string) string {
	s := strings.Fields(name)
	t := strings.Fields(nameType)
	// fmt.Println(s, "sss")
	for i := 0; i < len(s); i++ {
		v := FirstUpper(s[i])
		fmt.Println(v, FirstLower(t[i]))
	}
	return ""
}
