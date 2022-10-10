/*
 * @Author: jlz
 * @Date: 2022-09-08 13:24:16
 * @LastEditors: jlz
 * @Description:
 */

package constantModel

import (
	"dnds_go/common"
	"dnds_go/global"
	"dnds_go/tool"

	"gorm.io/gorm"
)

// Constant
// @Description:
type Constant struct {
	//
	//  Coding
	//  @Description:
	//
	Coding *int    `gorm:"type:int;column:CODING;size:10;primaryKey;autoIncrement;not null"` //
	Name   string  `json:"name" binding:"required" gorm:"type:varchar(50);column:NAME;unique;not null;"`
	EnName string  `json:"enName" gorm:"type:varchar(50);column:EN_NAME;"`
	Status *uint16 `json:"status" gorm:"type:tinyint;size:3;column:STATUS;not null;default:1"`
}

// @aram k1查询字段名称
// @aram k2查询字段名称
// @aram v 更新的值
// @return err
// @Description: 将符合字段名为k1、k2，值为v1、v2 的条件的值更改为 v
func (c Constant) UpdateCom(k1 string, k2 string, v1 interface{}, v2 interface{}, update interface{}) (err error) {
	global.L.Lock()
	defer global.L.Unlock()
	var condition string
	var value []interface{}
	value = append(value, v1)
	if k2 != "" {
		condition = k1 + "= ?"
		value = append(value, v2)
	}
	condition = k1 + "= ? AND" + k2 + "= ?"
	res := global.DB.Where(condition, value...).Updates(update)
	err = res.Error
	return
}
func (c Constant) TableName() string {
	return "SYS_CONSTANT"
}

// Create
// @receiver c
// @return err
// @Description:
func (c *Constant) Create() (err error) {
	cs, err := c.FindByName()
	if err != nil {
		return err
	}
	if len(cs) > 0 {
		return common.NewMyErr("常量名称已存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Create(c)
	if err = res.Error; err != nil {
		// global.Logger.Error(err.Error())
		return
	}
	return
}

func (c *Constant) Update() (err error) {
	constant, err := c.FindByCoding()
	if err != nil {
		return
	}
	if len(constant) == 0 {
		err = common.NewMyErr("常量不存在", nil)
		return
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Model(&Constant{}).Where("coding = ?", c.Coding).Updates(c)
	if err = res.Error; err != nil {
		return
	}
	return
}

func (c *Constant) GetAll() (csts []Constant, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Find(&csts)
	if err = res.Error; err != nil {
		return nil, err
	}
	return
}
func (c *Constant) Delete() (err error) {
	constant, err := c.FindByCoding()
	if err != nil {
		return
	}
	if len(constant) == 0 {
		err = common.NewMyErr("常量不存在", nil)
		return
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Delete(&Constant{}, c.Coding)
	if err = res.Error; err != nil {
		return
	}
	return
}

func (c *Constant) FindByName() (res []Constant, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("name = ?", c.Name).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *Constant) FindByEnName() (res []Constant, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("en_name = ?", c.EnName).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *Constant) FindByCoding() (res []Constant, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("coding = ?", c.Coding).Find(&res).Error; err != nil {
		return
	}
	return
}

type ConstantUpdateReq struct {
	Coding int    `json:"coding" binding:"required"`
	Name   string `json:"name"`
	EnName string `json:"enName"`
}
type ConstantDelReq struct {
	Coding int `json:"coding" binding:"required"`
}

type CstDetail struct {
	DetailID        int     `gorm:"type:int;column:DETAIL_ID;size:10;primaryKey;autoIncrement;not null"`
	Coding          *int    `json:"coding" binding:"required" gorm:"type:int;column:CODING;size:10;not null"`
	DetailName      string  `json:"detailName" binding:"required" gorm:"type:varchar;size:50;column:DETAIL_NAME;not null;"`
	EnName          string  `json:"enName" gorm:"type:varchar;size:50;column:EN_NAME;"`
	DetailValue     string  `json:"detailValue" binding:"required" gorm:"type:varchar;size:50;column:DETAIL_VALUE;not null"`
	Status          *uint16 `json:"status" gorm:"type:tineint;size:3;column:STATUS;not null;default:1"`
	PointIdentifier string  `gorm:"type:varchar;size:30;column:POINT_IDENTIFIER;"`
}

func (c CstDetail) TableName() string {
	return "SYS_CONSTANT_DETAIL"
}

func (c *CstDetail) Create() (err error) {
	exit, err := c.FindByCodingAndName()
	if err != nil {
		return err
	}
	if len(exit) > 0 {
		return common.NewMyErr("常量子项名称已存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Create(c)
	if err = res.Error; err != nil {
		return
	}
	return
}

func (c *CstDetail) Update() (err error) {
	exit, err := c.FindById()
	if err != nil {
		return err
	}
	if len(exit) < 1 {
		return common.NewMyErr("常量子项不存在", nil)
	}
	constant := Constant{
		Coding: c.Coding,
	}
	cExit, err := constant.FindByCoding()
	if err != nil {
		return err
	}
	if len(cExit) < 1 {
		return common.NewMyErr("常量不存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Model(&CstDetail{}).Where("detail_id = ?", c.DetailID).Updates(c)
	if err = res.Error; err != nil {
		return
	}
	return
}

func (c *CstDetail) DeleteById() (err error) {
	exit, err := c.FindById()
	if err != nil {
		return err
	}
	if len(exit) < 1 {
		return common.NewMyErr("常量子项不存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Delete(&CstDetail{}, c.DetailID)
	err = res.Error
	return
}
func (c *CstDetail) DeleteByCoding() (err error) {
	exit, err := c.FindByCoding()
	if err != nil {
		return err
	}
	if len(exit) < 1 {
		return common.NewMyErr("常量子项不存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("coding = ?", c.Coding).Delete(&CstDetail{})
	err = res.Error
	return
}
func (c *CstDetail) GetAll() (cstDetails []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if c.Coding == nil {
		res := global.DB.Find(&cstDetails)
		err = res.Error
		return
	}
	res := global.DB.Where("coding = ?", c.Coding).Find(&cstDetails)
	err = res.Error
	return
}
func (c *CstDetail) FindByName() (res []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("detail_name = ?", c.DetailName).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *CstDetail) FindByNameLike() (res []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	// pageInate := common.PageInate(1, 10)
	if err = global.DB.Where("detail_name LIKE ?", tool.LikeStr(c.DetailName)).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *CstDetail) FindById() (res []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("detail_id = ?", c.DetailID).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *CstDetail) FindByCoding() (res []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("coding = ?", c.Coding).Find(&res).Error; err != nil {
		return
	}
	return
}
func (c *CstDetail) FindByCodingAndName() (res []CstDetail, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	if err = global.DB.Where("coding = ? AND detail_name = ?", c.Coding, c.DetailName).Find(&res).Error; err != nil {
		return
	}
	return
}

type CstDetailUpdateReq struct {
	Coding          int    `json:"coding" binding:"required"`
	DetailID        int    `json:"detailId" binding:"required"`
	DetailName      string `json:"detailName"`
	EnName          string `json:"enName"`
	DetailValue     string `json:"detailValue"`
	PointIdentifier string `json:"pointIdentifier"`
}
type CstDetailAllReq struct {
	Coding   int    `json:"coding" binding:"required"`
	LikeName string `json:"likeName"`
}
type CstDetailDeleteReq struct {
	Id int `json:"id" binding:"required"`
}

type EventCst struct {
	Id            int    `gorm:"type:int;column:ID;size:10;primaryKey;autoIncrement;not null"`
	Name          string `json:"name" binding:"required,max=50"  gorm:"type:varchar;size:50;column:NAME;not null;"`
	Code          string `json:"code" binding:"max=50"  gorm:"type:varchar;size:50;column:CODE;"`
	PointIdentify string `json:"pointIdentify" binding:"max=50"  gorm:"type:varchar;size:50;column:POINT_IDENTIFY;"`
	SensorType    *int   `json:"sensorType" binding:"omitempty,number" gorm:"type:int;size:10;column:SENSOR_TYPE;"`
}

func (e EventCst) TableName() string {
	return "SYS_CONSTANT_EVENT"
}
func (e *EventCst) Create() (err error) {
	// exit, err := e.FindByName()
	// if err != nil {
	// 	return
	// }
	// if len(exit) > 0 {
	// 	return common.NewMyErr("事件常量名称已存在", nil)
	// }
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Create(e)
	err = res.Error
	return
}

func (e *EventCst) Update() (err error) {
	exit, err := e.FindById()
	if err != nil {
		return
	}
	if len(exit) == 0 {
		return common.NewMyErr("事件常量不存在，无法编辑", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()

	res := global.DB.Model(&EventCst{}).Where("id = ?", e.Id).Updates(e)
	err = res.Error
	return
}
func (e *EventCst) GetList() (events []EventCst, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	var res *gorm.DB
	if e.SensorType == nil {
		res = global.DB.Where("name LIKE ?", tool.LikeStr(e.Name)).Find(&events)
		err = res.Error
		return
	}
	res = global.DB.Where("name LIKE ? AND sensor_type = ?", tool.LikeStr(e.Name), e.SensorType).Find(&events)
	err = res.Error
	return
}
func (e *EventCst) Delete() (err error) {
	global.L.Lock()
	defer global.L.Unlock()
	// global.DB.Delete(&CstDetail{}, c.DetailID)
	res := global.DB.Delete(&EventCst{}, e)
	err = res.Error
	return
}

func (e *EventCst) FindByName() (events []EventCst, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("name = ?", e.Name).Find(&events)
	err = res.Error
	return
}
func (e *EventCst) FindById() (events []EventCst, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("id = ?", e.Id).Find(&events)
	err = res.Error
	return
}

type EventCstUpdateReq struct {
	Id            int    `json:"id" binding:"required,number"`
	EventCode     string `json:"eventCode"`
	PointIdentify string `json:"pointIdentify"`
	SensorType    int    `json:"sensorType"`
	Name          string `json:"name"`
}
type EventCstDeleteReq struct {
	Id int `json:"id" binding:"required,number"`
}
type EventCstGetListRes struct {
	Id            int    `json:"id"`
	EventCode     string `json:"eventCode"`
	PointIdentify string `json:"pointIdentify"`
	SensorType    int    `json:"sensorType"`
	EventName     string `json:"eventName"`
	TypeName      string `json:"typeName"`
	PointName     string `json:"pointName"`
}
