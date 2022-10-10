package common

/*
	封装导出excel表格
*/
import (
	"dnds_go/tool"
	"errors"
	"fmt"
	"image"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// excel实体结构体，可以设置excel生成属性
type ExcelInfo struct {
	FileName         string          // 导出文件名称
	Titles           []string        // 标题数组
	DataKeys         []string        // 数据字段名称
	ShowTime         bool            // 文件名后是否添加时间 默认不添加
	Suffix           string          // 导出格式  默认 .xlsx
	SavePath         string          //  保存路径
	SheetName        string          // 表名  默认Sheet
	RowHeight        int             // 行高 默认30
	ColWidth         int             // 列宽 默认50
	TitleSytle       *excelize.Style // 默认加粗， 水平、垂直居中
	DataSytle        *excelize.Style // 默认， 水平、垂直居中
	SpecificColWidth map[string]int
	Scale            float32 // 图片宽高压缩比 默认 0.2
	AutoScale        bool    // 是否开启自动判断图片压缩比，开启后上边Scale无效，会读取图片占用一定资源
}

// 数据源是结构体，解析结构体，生成表格
func (e *ExcelInfo) CreateExcel(d interface{}) (*excelize.File, error) {
	data, ok := tool.CreateAnyTypeSlice(d)
	if !ok {
		return nil, errors.New("excel数据源必须为切片类型")
	}
	f := excelize.NewFile()

	sheetName := "Sheet"
	if e.SheetName == "" {
		e.SheetName = sheetName
	}

	index := f.NewSheet(sheetName)
	if e.RowHeight == 0 {
		e.RowHeight = 20
	}
	if e.ColWidth == 0 {
		e.ColWidth = 20
	}

	if len(e.Titles) != len(data) {
		return nil, errors.New("生成excel标题列数和数据列数不匹配")
	}
	titleStyle, err := e.GetTitleStytle(f)
	if err != nil {
		return nil, err
	}
	// 设置表头
	for i, title := range e.Titles {
		cellNum := fmt.Sprintf("%c1", 65+i)
		fmt.Println(cellNum)
		err := f.SetCellValue(e.SheetName, cellNum, title)
		if err != nil {
			return nil, err
		}
		// 设置 行宽
		curCol := string('A' + byte(i))
		if v, ok := e.SpecificColWidth[title]; ok {
			// 使用指定宽度
			f.SetColWidth(e.SheetName, curCol, curCol, float64(v))
		} else {
			// 使用默认宽度
			f.SetColWidth(e.SheetName, curCol, curCol, float64(e.ColWidth))
		}

		// 设置 表头样式
		_ = f.SetCellStyle(e.SheetName, cellNum, cellNum, titleStyle)
		err = f.SetRowHeight(e.SheetName, 1+i, float64(e.RowHeight))
		if err != nil {
			return nil, err
		}
	}
	// 生成默认data部分样式
	dataStyle, err := e.GetDataStytle(f)
	if err != nil {
		return nil, err
	}
	rowNum := 1
	for i, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface {
		}, 0)
		for l := 0; l < t.NumField(); l++ {
			line := fmt.Sprintf("%c%d", 65+l, i+2)
			val := value.Field(l).Interface()
			row = append(row, val)

			// 设置 数据区样式
			err = f.SetCellStyle(e.SheetName, line, line, dataStyle)
			if err != nil {
				return nil, err
			}
		}
		rowNum++
		err := f.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		// _ = f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID)
		if err != nil {

			return nil, err
		}
		// 设置工作簿的默认工作表
		f.SetActiveSheet(index)
	}
	return f, nil
}

//ExcelDatA 导出数据 行格式
type ExcelData []map[string]interface{}

// 数据源是map 解析map，生成表格
func (e *ExcelInfo) CreateEexcelByMap(data ExcelData) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sheet"
	if e.SheetName == "" {
		e.SheetName = sheetName
	}

	index := f.NewSheet(sheetName)
	if e.RowHeight == 0 {
		e.RowHeight = 20
	}
	if e.ColWidth == 0 {
		e.ColWidth = 20
	}

	if len(e.Titles) != len(data) {
		return nil, errors.New("生成excel标题列数和数据列数不匹配")
	}
	titleStyle, err := e.GetTitleStytle(f)
	if err != nil {
		return nil, err
	}
	// 设置表头
	for i, title := range e.Titles {
		cellNum := fmt.Sprintf("%c1", 65+i)
		fmt.Println(cellNum)
		err := f.SetCellValue(e.SheetName, cellNum, title)
		if err != nil {
			return nil, err
		}
		// 设置 行宽
		curCol := string('A' + byte(i))
		if v, ok := e.SpecificColWidth[title]; ok {
			// 使用指定宽度
			f.SetColWidth(e.SheetName, curCol, curCol, float64(v))
		} else {
			// 使用默认宽度
			f.SetColWidth(e.SheetName, curCol, curCol, float64(e.ColWidth))
		}

		// 设置 表头样式
		err = f.SetCellStyle(e.SheetName, cellNum, cellNum, titleStyle)
		if err != nil {
			return nil, err
		}
		err = f.SetRowHeight(e.SheetName, 1+i, float64(e.RowHeight))
		if err != nil {
			return nil, err
		}
	}

	// 生成默认data部分样式
	dataStyle, err := e.GetDataStytle(f)
	if err != nil {
		return nil, err
	}
	// 设置单元格的值
	for i, rows := range data { //一行

		for num, key := range e.DataKeys { //每列
			line := fmt.Sprintf("%c%d", 65+num, i+2)

			if strings.Index(key, "IMG_") == -1 { //正常值
				err = f.SetCellValue(e.SheetName, line, rows[key])
				if err != nil {
					return nil, err
				}
			} else { //图片
				link := fmt.Sprintf("%s", rows[key])
				format := ``
				if e.AutoScale {
					//获取图片宽高
					file, err := os.Open(link)
					if err != nil {
						return nil, err
					}
					img, _, err := image.Decode(file)
					if err != nil {
						return nil, err
					}
					err = file.Close()
					if err != nil {
						return nil, err
					}
					b := img.Bounds()
					yScale := fmt.Sprintf("%f", float64(e.RowHeight)/float64(b.Max.Y))
					format = `{  "x_scale": ` + yScale + `, "y_scale": ` + yScale + ` }`
				} else {
					format = `{  "x_scale": ` + fmt.Sprintf("%f", e.Scale) + `, "y_scale": ` + fmt.Sprintf("%f", e.Scale) + ` }`
				}

				err = f.AddPicture("Sheet1", line, link, format)
				if err != nil {
					return nil, err
				}
			}
			// 设置 数据区样式
			if num == 0 {
				err = f.SetCellStyle(e.SheetName, line, line, dataStyle)
				if err != nil {
					return nil, err
				}
			}
		}
		//设置行高
		err = f.SetRowHeight("Sheet1", 2+i, float64(e.RowHeight))
		if err != nil {
			return nil, err
		}
	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f, nil
}

// 生成title 样式
func (e *ExcelInfo) GetTitleStytle(f *excelize.File) (int, error) {
	if e.TitleSytle != nil {
		return f.NewStyle(e.TitleSytle)
	}
	return f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center", // 垂直居中
			Horizontal: "center", // 水平居中
		},
	})

}

// 生成data 样式
func (e *ExcelInfo) GetDataStytle(f *excelize.File) (int, error) {
	if e.TitleSytle != nil {
		return f.NewStyle(e.DataSytle)
	}
	return f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
	})
}

// 向前端返回excel
func (e *ExcelInfo) DownloadExcel(c *gin.Context, f *excelize.File) error {
	curTime := time.Now().Format("20060102150405")
	if e.ShowTime && e.FileName != "" {
		e.FileName = e.FileName + curTime
	}
	if e.FileName == "" {
		e.FileName = curTime
	}
	if e.Suffix == "" {
		e.Suffix = ".xlsx"
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+e.FileName+e.Suffix)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	return f.Write(c.Writer)
}

// MaxCharCount 最多26个字符A-Z
const MaxCharCount = 26

// headers 列名切片， 表头
// rows 数据源是一个二维数组
func CreateExcel(sheetName string, headers []string, rows [][]interface{}) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetIndex := f.NewSheet(sheetName)
	_ = f.SetColWidth(sheetName, "A", "H", 30) // 设置列宽度

	maxColumnRowNameLen := 1 + len(strconv.Itoa(len(rows)))
	columnCount := len(headers)

	if columnCount > MaxCharCount {
		maxColumnRowNameLen++
	} else if columnCount > MaxCharCount*MaxCharCount {
		maxColumnRowNameLen += 2
	}
	columnNames := make([][]byte, 0, columnCount)

	for i, header := range headers {
		columnName := getColumnName(i, maxColumnRowNameLen)
		columnNames = append(columnNames, columnName)
		// 初始化excel表头，这里的index从1开始要注意
		curColumnName := getColumnRowName(columnName, 1)
		err := f.SetCellValue(sheetName, curColumnName, header)
		if err != nil {
			return nil, err
		}
	}

	for rowIndex, row := range rows {

		for columnIndex, columnName := range columnNames {
			// 从第二行开始
			err := f.SetCellValue(sheetName, getColumnRowName(columnName, rowIndex+2), row[columnIndex])
			if err != nil {
				return nil, err
			}
		}
	}

	f.SetActiveSheet(sheetIndex)

	return f, nil
}

// getColumnName 生成列名
// Excel的列名规则是从A-Z往后排;超过Z以后用两个字母表示，比如AA,AB,AC;两个字母不够以后用三个字母表示，比如AAA,AAB,AAC
// 这里做数字到列名的映射：0 -> A, 1 -> B, 2 -> C
// maxColumnRowNameLen 表示名称框的最大长度，假设数据是10行，1000列，则最后一个名称框是J1000(如果有表头，则是J1001),是4位
// 这里根据 maxColumnRowNameLen 生成切片，后面生成名称框的时候可以复用这个切片，而无需扩容
func getColumnName(column, maxColumnRowNameLen int) []byte {
	const A = 'A'
	if column < MaxCharCount {
		// 第一次就分配好切片的容量
		slice := make([]byte, 0, maxColumnRowNameLen)
		return append(slice, byte(A+column))
	} else {
		// 递归生成类似AA,AB,AAA,AAB这种形式的列名
		return append(getColumnName(column/MaxCharCount-1, maxColumnRowNameLen), byte(A+column%MaxCharCount))
	}
}

// getColumnRowName 生成名称框
// Excel的名称框是用A1,A2,B1,B2来表示的，这里需要传入前一步生成的列名切片，然后直接加上行索引来生成名称框，就无需每次分配内存
func getColumnRowName(columnName []byte, rowIndex int) (columnRowName string) {
	l := len(columnName)
	columnName = strconv.AppendInt(columnName, int64(rowIndex), 10)
	columnRowName = string(columnName)
	// 将列名恢复回去
	columnName = columnName[:l]
	return
}

// 数据源是二维数组向前端返回excel
func DownloadExcel(c *gin.Context, f *excelize.File, fileName string, showTime bool) error {
	curTime := time.Now().Format("20060102150405")
	if showTime && fileName != "" {
		fileName = fileName + curTime
	}
	if fileName == "" {
		fileName = curTime
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName+".xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	return f.Write(c.Writer)
}
