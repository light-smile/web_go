/*
 * @Author: jlz
 * @Date: 2022-08-31 14:13:54
 * @LastEditTime: 2022-08-31 18:04:16
 * @LastEditors: jlz
 * @Description:网关档案
 */

package gatewayService

import (
	"dnds_go/common"
	gatewayModel "dnds_go/src/models/gateway"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var req gatewayModel.Gateway
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	if err := req.Create(); err != nil {
		if ok := common.MyErrHandle(c, err); ok {
			return
		}
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "网关档案创建失败")
	}
	common.Succ(c, nil, "网关档案创建成功")
}

func Update(c *gin.Context) {
	var req gatewayModel.GatewayUpdateReq
	if ok := common.BindJson(c, &req); !ok {
		return
	}
	var inlines, outlines *int
	if req.Inlines != "" {
		i, err := strconv.Atoi(req.Inlines)
		if err != nil {
			return
		}
		inlines = &i
	}
	if req.Outlines != "" {
		i, err := strconv.Atoi(req.Outlines)
		if err != nil {
			return
		}
		outlines = &i
	}
	gateway := gatewayModel.Gateway{
		Id:             req.Id,
		GatewayName:    req.GatewayName,
		GatewayAddress: req.GatewayAddress,
		SiteName:       req.SiteName,
		AreaIds:        req.AreaIds,
		SiteType:       req.SiteType,
		SiteVolt:       req.SiteVolt,
		SiteLevel:      req.SiteLevel,
		SiteBuild:      req.SiteBuild,
		SiteBDZ:        req.SiteBDZ,
		StationLine:    req.StationLine,
		Longitude:      req.Longitude,
		Latitude:       req.Latitude,
		Inlines:        inlines,
		Outlines:       outlines,
		DevOps:         req.DevOps,
		ReSgAtt1:       req.ReSgAtt1,
		ReSgAtt2:       req.ReSgAtt2,
		ReSgAtt3:       req.ReSgAtt3,
		ReSgAtt4:       req.ReSgAtt4,
		ReCoAtt1:       req.ReCoAtt1,
		ReCoAtt2:       req.ReCoAtt2,
		ReCoAtt3:       req.ReCoAtt3,
		ReCoAtt4:       req.ReCoAtt4,
	}
	if err := gateway.Update(); err != nil {
		common.LogErr(req, err)
		common.Fail(c, common.ErrorCommonCode, nil, "网关档案更新失败")
		return
	}
	common.Succ(c, nil, "网关档案修改成功")
}

func GetList(c *gin.Context) {
	var gateway gatewayModel.Gateway
	res, err := gateway.GetList()
	if err != nil {
		common.LogErr(nil, err)
		common.Fail(c, common.ErrorCommonCode, nil, "网关档案列表获取失败")
		return
	}
	common.Succ(c, res, "网关档案列表获取成功")
}
