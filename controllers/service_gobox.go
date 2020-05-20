package controllers

import (
	"github.com/astaxie/beego"
	gversion "github.com/mcuadros/go-version"
	//"github.com/guerillagrow/gobox/models"
)

type ServiceGoBox struct {
	beego.Controller
}

func (c *ServiceGoBox) Get() {

	// TODO(bonedaddy): replace
	currentVersion := "1.0.0"
	localVersion := "1.0.1"
	updateable := gversion.Compare(localVersion, currentVersion, "<")

	res := JSONResp{
		Meta: map[string]interface{}{
			"status": 200,
		},
		Data: map[string]interface{}{
			"current_version":  currentVersion,
			"local_version":    localVersion,
			"update_available": updateable,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()

}
