/*
 * @Author: jlz
 * @Date: 2022-09-08 14:40:26
 * @LastEditors: jlz
 * @Description: 实现gorm通用查询方法
 */

package common

import (
	"dnds_go/global"
	"gorm.io/gorm"
)

type FindComArgs struct {
	ColumnName []string      // 列名称数组
	Values     []interface{} // 值数组，与列顺序保持一致
}

// 可以对结果继续进行操作
func (f *FindComArgs) FindCom() *gorm.DB {
	var conditionStr string
	for i := 0; i < len(f.ColumnName); i++ {
		if i == 0 {
			conditionStr += f.ColumnName[i] + "= ?"
			continue
		}
		conditionStr += "And" + f.ColumnName[i] + "= ?"
	}
	return global.DB.Where(conditionStr, f.Values...)
}
func FindComById(v int) *gorm.DB {
	return global.DB.Where("id = ?", v)
}
