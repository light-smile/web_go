package common

import (
	"testing"
)

type info struct {
	Name string
	Age  int
}

func TestFindCom(t *testing.T) {
	com := FindComArgs{
		ColumnName: []string{
			"name",
			"age",
		},
		Values: []interface{}{
			"zhangsan",
			18,
		},
	}
	var res []info
	com.FindCom().Find(&res)
}
