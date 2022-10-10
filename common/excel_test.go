package common

/*
	测试excel
*/
import (
	"fmt"
	"testing"
)

func TestExportExcel(t *testing.T) {
	excel := ExcelInfo{
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
	type User struct {
		UserName string
		Sex      string
		Age      string
	}
	dataSturct := []User{
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
		fmt.Errorf("create excel err: %v", err.Error())
	}
	if err := f.SaveAs("testExcel.xlsx"); err != nil {
		fmt.Errorf("save excel err: %v", err.Error())
	}
}
