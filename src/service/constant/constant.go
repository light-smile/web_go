/*
 * @Author: jlz
 * @Date: 2022-08-27 15:12:25
 * @LastEditTime: 2022-08-31 13:15:57
 * @LastEditors: jlz
 * @Description:
 */

package constantService

import (
	"dnds_go/common"
	"dnds_go/global"
	constantDao "dnds_go/src/dao/constant"
	constantModel "dnds_go/src/models/constant"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateConstant(c *gin.Context) {
	var cst constantModel.Constant
	if ok := common.BindJson(c, &cst); !ok {
		return
	}

	if err := cst.Create(); err != nil {
		if ok := common.IsMyErr(err); ok {
			common.Fail(c, common.ErrorCommonCode, nil, err.Error())
			return
		}
		common.LogErr(cst, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量添加失败")
		return
	}
	common.Succ(c, nil, "常量添加成功")
}

func UpdateConstant(c *gin.Context) {
	var req constantModel.ConstantUpdateReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	cst := constantModel.Constant{
		Coding: &req.Coding,
		Name:   req.Name,
		EnName: req.EnName,
	}
	if err := cst.Update(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(cst, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量修改失败")
		return
	}
	common.Succ(c, nil, "常量修改成功")
}

// GetConstant 获取全部常量
func GetConstant(c *gin.Context) {
	var cst constantModel.Constant
	csts, err := cst.GetAll()
	if err != nil {
		common.LogErr(nil, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量获取失败")
	}
	common.Succ(c, csts, "常量获取成功")

}

func DeleteConstant(c *gin.Context) {
	var req constantModel.ConstantDelReq

	if ok := common.BindJson(c, &req); !ok {
		return
	}
	cst := constantModel.Constant{
		Coding: &req.Coding,
	}
	if err := cst.Delete(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量删除失败")
	}
	// 如何常量有子项，子项也要删除
	cstDetail := constantModel.CstDetail{
		Coding: &req.Coding,
	}
	if err := cstDetail.DeleteByCoding(); err != nil {
		if ok := common.IsMyErr(err); ok {
			common.Succ(c, nil, "常量删除成功")
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量子项删除失败")
	}
	common.Succ(c, nil, "常量删除成功")
}

// 创建常量子项
func CreateCstDetail(c *gin.Context) {
	var req constantModel.CstDetail
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	if err := req.Create(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量子项创建失败")
		return
	}
	common.Succ(c, nil, "常量子项创建成功")
}

func UpdateCstDetail(c *gin.Context) {
	var req constantModel.CstDetailUpdateReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	cstDetail := constantModel.CstDetail{
		Coding:   &req.Coding,
		DetailID: req.DetailID,
	}
	if err := cstDetail.Update(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量子项更新失败")
	}
	common.Succ(c, nil, "常量子项修改成功")
}
func DeleteCstDetail(c *gin.Context) {
	var req constantModel.CstDetailDeleteReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}

	cstDetail := constantModel.CstDetail{
		DetailID: req.Id,
	}
	if err := cstDetail.DeleteById(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "常量子项删除失败")
	}
	common.Succ(c, nil, "常量子项删除成功")
}

// 获取常量子项通过coding或者模糊搜索name
func GetDetailAllByCoding(c *gin.Context) {
	var req constantModel.CstDetailAllReq
	coding := c.Query("coding")
	i, err := strconv.Atoi(coding)
	if err != nil {
		common.Fail(c, common.ErrorBadParamCode, nil, "coding不是数字类型")
		return
	}
	req.Coding = i
	req.LikeName = c.Query("likeName")

	CstDetail := constantModel.CstDetail{
		Coding:     &req.Coding,
		DetailName: req.LikeName,
	}
	if CstDetail.DetailName != "" {
		cstDetails, err := CstDetail.FindByNameLike()
		if err != nil {
			if ok := common.MyErrHandle(c, err); ok {
				return
			}
			common.LogErr(req, err)
			common.Fail(c, common.ErrorCommonCode, nil, "获取常量子项列表失败")
			return
		}
		common.Succ(c, cstDetails, "获取常量子项列表成功")
		return
	}
	cstDetails, err := CstDetail.GetAll()
	if err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "获取常量子项列表失败")
	}
	common.Succ(c, cstDetails, "获取常量子项列表成功")
}

func CreateEventCst(c *gin.Context) {
	var req constantModel.EventCst
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	if err := req.Create(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "事件常量创建失败")
	}
	common.Succ(c, nil, "事件常量创建成功")
}

func UpdateEventCst(c *gin.Context) {
	var req constantModel.EventCstUpdateReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	event := constantModel.EventCst{
		Id:            req.Id,
		Name:          req.Name,
		Code:          req.EventCode,
		PointIdentify: req.PointIdentify,
		SensorType:    &req.SensorType,
	}
	err := event.Update()
	if err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "事件常量更新失败")
		return
	}
	common.Succ(c, nil, "事件常量更新成功")
}

// 获取事件列表  名称/设备类型(sensortype)
func GetEventCstList(c *gin.Context) {
	name := c.Query("keyword")
	sensorType := c.Query("typeCode")
	// 没有传设备类型编码为-1,标志查询条件不包含sensorType
	var st *int
	if sensorType != "" {
		s, err := strconv.Atoi(sensorType)
		if err != nil {
			common.Fail(c, common.ErrorCommonCode, nil, "设备类型编码格式不正确")
			return
		}
		st = &s
	}
	eventCst := constantModel.EventCst{
		Name:       name,
		SensorType: st,
	}
	eventCsts, err := constantDao.GetEventCstList(eventCst)
	if err != nil {
		common.Fail(c, common.ErrorCommonCode, nil, "事件列表获取失败")
		return
	}
	common.Succ(c, eventCsts, "事件列表获取成功")
}

func DeleteEventCst(c *gin.Context) {
	var req constantModel.EventCstDeleteReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	eventCst := constantModel.EventCst{
		Id: req.Id,
	}
	if err := eventCst.Delete(); err != nil {
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "事件常量删除失败")
		return
	}
	common.Succ(c, nil, "事件常量删除成功")
}

/**
 * @description: 将constant_detail表中数据导入到event表中
 * @param {*gin.Context} c
 * @return {*}
 * @use:
 */
func TestData(c *gin.Context) {
	// global.DB.Table()
	// 仅测试时使用
	return
	var constant constantModel.CstDetail
	res, _ := constant.GetAll()
	num := 0
	var events []constantModel.EventCst
	for _, v := range res {
		num++
		event := constantModel.EventCst{
			Id:            v.DetailID,
			Code:          v.DetailValue,
			Name:          v.DetailName,
			PointIdentify: v.PointIdentifier,
		}
		events = append(events, event)
		// event.Create()
		// event := constantModel.EventCst{
		// 	Id: v.DetailID,
		// }
		// s, _ := event.FindById()
		// if len(s) == 0 {
		// 	fmt.Println(v.DetailID, "id")
		// }
	}
	global.DB.Create(&events)
	common.Succ(c, num, "ok")
}

func GetSensorTypeByPoint(c *gin.Context) {
	// 仅测试时使用
	return
	constantDao.GetSensorTypeByPoint()
	common.Succ(c, nil, "ok")
}
