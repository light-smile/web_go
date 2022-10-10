/*
 * @Author: jlz
 * @Date: 2022-08-30 14:44:06
 * @LastEditTime: 2022-08-31 13:14:04
 * @LastEditors: jlz
 * @Description:
 */
package constantDao

import (
	"dnds_go/global"
	constantModel "dnds_go/src/models/constant"
	"dnds_go/tool"
	"fmt"
	"strconv"
)

func GetEventCstList(c constantModel.EventCst) (res []constantModel.EventCstGetListRes, err error) {
	fmt.Println(c, "adfaldfjl")
	if c.SensorType == nil {
		r := global.DB.Table("SYS_CONSTANT_EVENT").
			Select(
				"SYS_CONSTANT_EVENT.NAME as EventName,SYS_CONSTANT_EVENT.CODE as EventCode,SYS_CONSTANT_EVENT.POINT_IDENTIFY as PointIdentify,SYS_CONSTANT_EVENT.SENSOR_TYPE as SensorType,END_TYPE.TYPE_NAME as TypeName,END_TYPE_POINT.POINT_NAME as PointName").Joins("left join END_TYPE on END_TYPE.TYPE_CODE = SYS_CONSTANT_EVENT.SENSOR_TYPE").Joins("left join END_TYPE_POINT on END_TYPE_POINT.POINT_IDENTIFICATION = SYS_CONSTANT_EVENT.POINT_IDENTIFY").Where("name LIKE ? ", tool.LikeStr(c.Name)).Find(&res)
		err = r.Error

	}
	r := global.DB.Table("SYS_CONSTANT_EVENT").
		Select(
			"SYS_CONSTANT_EVENT.NAME as EventName,SYS_CONSTANT_EVENT.CODE as EventCode,SYS_CONSTANT_EVENT.POINT_IDENTIFY as PointIdentify,SYS_CONSTANT_EVENT.SENSOR_TYPE as SensorType,END_TYPE.TYPE_NAME as TypeName,END_TYPE_POINT.POINT_NAME as PointName").Joins("left join END_TYPE on END_TYPE.TYPE_CODE = SYS_CONSTANT_EVENT.SENSOR_TYPE").Joins("left join END_TYPE_POINT on END_TYPE_POINT.POINT_IDENTIFICATION = SYS_CONSTANT_EVENT.POINT_IDENTIFY").Where("name LIKE ? AND sensorType = ?", tool.LikeStr(c.Name), c.SensorType).Find(&res)
	err = r.Error

	return
}

type SensorType struct {
	PointIdentify string
	SensorType    string
}

/**
 * @description: 通过测点标识获取设备类型(sensortype)
 * @return {*}
 * @use:
 */
func GetSensorTypeByPoint() {
	return // 仅测试时使用
	var Ss []SensorType
	global.DB.Table("SYS_CONSTANT_EVENT").Select("SYS_CONSTANT_EVENT.POINT_IDENTIFY as PointIdentify,END_TYPE_POINT.TYPE_ID as SensorType").Joins("left join END_TYPE_POINT on END_TYPE_POINT.POINT_IDENTIFICATION = SYS_CONSTANT_EVENT.POINT_IDENTIFY").Scan(&Ss)
	for _, v := range Ss {
		if v.SensorType != "" && v.PointIdentify != "" {
			i, _ := strconv.Atoi(v.SensorType)
			global.DB.Table("SYS_CONSTANT_EVENT").Where("POINT_IDENTIFY = ?", v.PointIdentify).Update("SENSOR_TYPE", i)
		}
	}

}
