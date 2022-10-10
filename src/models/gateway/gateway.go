/*
 * @Author: jlz
 * @Date: 2022-08-31 13:21:45
 * @LastEditTime: 2022-08-31 17:13:25
 * @LastEditors: jlz
 * @Description:网关档案结结构体
 */

package gatewayModel

import (
	"dnds_go/common"
	"dnds_go/global"
	"time"
)

type Gateway struct {
	Id             *int   `json:"id" gorm:"type:int;column:id;primaryKey;autoIncrement;not null"`
	GatewayName    string `json:"gatewayName" binding:"required,max=50" gorm:"type:varchar(20);column:gateway_name;not null"`       // 网关名称
	GatewayAddress string `json:"gatewayAddress" binding:"required,max=50" gorm:"type:varchar(50);column:gateway_address;not null"` // 网关地址
	SiteName       string `json:"siteName" binding:"max=30" gorm:"type:varchar(30);column:site_name"`                               // 站房名称
	// SiteArea       string `json:"siteArea" binding:"max=30" gorm:"type:varchar(30);column:site_area"`                                 // 站房所属地市
	// SiteZone       string `json:"siteZone" binding:"max=20" gorm:"type:varchar(20);column:site_zone"`                                 // 站房所属区/县
	AreaIds      string `json:"areaIds" binding:"max=50" gorm:"type:varchar(50);column:area_ids"`                   // areaIds
	SiteType     string `json:"siteType" binding:"max=20" gorm:"type:varchar(20);column:site_type"`                 // 站点类型
	SiteVolt     string `json:"siteVolt" binding:"max=20" gorm:"type:varchar(20);column:site_volt"`                 // 电压等级
	SiteLevel    string `json:"siteLevel" binding:"max=20" gorm:"type:varchar(20);column:site_level"`               // 站点等级
	SiteBuild    string `json:"siteBuild" binding:"max=20" gorm:"type:varchar(20);column:site_build"`               // 建设厂商名称
	SiteBDZ      string `json:"siteBDZ" binding:"max=20" gorm:"type:varchar(20);column:siteBDZ"`                    // 所属变电站
	StationLine  string `json:"stationLine" binding:"max=225" gorm:"type:varchar(225);column:station_line"`         // 站点进线线路
	Longitude    string `json:"longitude" binding:"max=30" gorm:"type:varchar(30);column:longitude"`                // 经度
	Latitude     string `json:"latitude" binding:"max=30" gorm:"type:varchar(30);column:latitude"`                  // 纬度
	Inlines      *int   `json:"inlines"  gorm:"type:int;column:inlines"`                                            // 进线数
	Outlines     *int   `json:"outlines"  gorm:"type:int;column:outlines"`                                          // 出线数
	DevOps       string `json:"devOps" binding:"max=225" gorm:"type:varchar(225);column:dev_ops"`                   //运维班组
	OnlineStatus *int   `json:"onlineStatus" gorm:"type:int;column:online_status;default:0;comment:在线状态：0-离线，1-在线"` //在线状态：0-离线，1-在线
	CreateTime   string `json:"createTime" gorm:"type:varchar(50);column:create_time;"`
	ReSgAtt1     string `json:"reSgAtt1" gorm:"type:varchar(50);column:ReSgAtt1;"` // 遥信 1 接入属性
	ReSgAtt2     string `json:"reSgAtt2" gorm:"type:varchar(50);column:ReSgAtt2;"` // 遥信 2 接入属性
	ReSgAtt3     string `json:"reSgAtt3" gorm:"type:varchar(50);column:ReSgAtt3;"` // 遥信 3 接入属性
	ReSgAtt4     string `json:"reSgAtt4" gorm:"type:varchar(50);column:ReSgAtt4;"` // 遥信 4 接入属性
	ReCoAtt1     string `json:"reCoAtt1" gorm:"type:varchar(50);column:ReCoAtt1;"` // 遥信 1 接入属性
	ReCoAtt2     string `json:"reCoAtt2" gorm:"type:varchar(50);column:ReCoAtt2;"` // 遥信 2 接入属性
	ReCoAtt3     string `json:"reCoAtt3" gorm:"type:varchar(50);column:ReCoAtt3;"` // 遥信 3 接入属性
	ReCoAtt4     string `json:"reCoAtt4" gorm:"type:varchar(50);column:ReCoAtt4;"` // 遥信 4 接入属性
}

func (g Gateway) TableName() string {
	return "gateway"
}

func (g *Gateway) Create() (err error) {
	exit, err := g.FindByAddress()
	if err != nil {
		return
	}
	if len(exit) > 0 {
		return common.NewMyErr("网关地址已存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	g.CreateTime = time.Now().Format(global.TimeFormatYmdhis)
	res := global.DB.Create(g)
	err = res.Error
	return
}

func (g *Gateway) Update() (err error) {
	exit, err := g.FindById()
	if err != nil {
		return
	}
	if len(exit) < 1 {
		return common.NewMyErr("网关档案不存在", nil)
	}
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Model(&Gateway{}).Where("id = ?", g.Id).Updates(g)
	err = res.Error
	return
}
func (g *Gateway) GetList() (gateways []GatewayGetListRes, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	// var gateway Gateway
	res := global.DB.Table("gateway").Find(&gateways)
	err = res.Error
	return
}
func (g *Gateway) FindByAddress() (gateways []Gateway, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("gateway_address = ?", g.GatewayAddress).Find(&gateways)
	err = res.Error
	return
}
func (g *Gateway) FindById() (gateways []Gateway, err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("id = ?", g.Id).Find(&gateways)
	err = res.Error
	return
}

type GatewayUpdateReq struct {
	Id             *int `binding:"required"`
	GatewayName    string
	GatewayAddress string
	SiteName       string
	AreaIds        string
	SiteType       string
	SiteLevel      string
	SiteVolt       string
	SiteBuild      string
	SiteBDZ        string
	StationLine    string
	Longitude      string
	Latitude       string
	Inlines        string
	Outlines       string
	DevOps         string
	ReSgAtt1       string
	ReSgAtt2       string
	ReSgAtt3       string
	ReSgAtt4       string
	ReCoAtt1       string
	ReCoAtt2       string
	ReCoAtt3       string
	ReCoAtt4       string
}

type GatewayGetListRes struct {
	Id             *int `binding:"required"`
	GatewayName    string
	GatewayAddress string
	SiteName       string
	AreaIds        string
	SiteType       string
	SiteLevel      string
	SiteVolt       string
	SiteBuild      string
	SiteBDZ        string
	StationLine    string
	Longitude      string
	Latitude       string
	Inlines        string
	Outlines       string
	DevOps         string
	ReSgAtt1       string
	ReSgAtt2       string
	ReSgAtt3       string
	ReSgAtt4       string
	ReCoAtt1       string
	ReCoAtt2       string
	ReCoAtt3       string
	ReCoAtt4       string
}
