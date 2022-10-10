package common

import (
	"dnds_go/global"

	"gorm.io/gorm"
)

// 翻页封装，可用返回值基础操作
func PageInate(page int, pageSize int) func() *gorm.DB {
	return func() *gorm.DB {
		if page == 0 {
			page = 1
		}
		// 页数限制 和 页大小限制
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return global.DB.Offset(offset).Limit(pageSize)
	}
}
