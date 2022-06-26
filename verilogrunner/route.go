package verilogrunner

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (vr *VerilogRunner) addRouteCompiler() {
	vr.server.POST("/compile", vr.CompileVerilog)
}

func (vr *VerilogRunner) addRouteStatic() {
	vr.server.StaticFile("/favicon.ico", "static/images/favicon.png")
}

func (vr *VerilogRunner) addRouteStatus() {
	status := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"server":  "verilog-fl",
			"version": "v0.0.1",
		})
	}

	vr.server.GET("/", status)
	vr.server.GET("/status", status)
}

func (vr *VerilogRunner) addRoute() {
	vr.addRouteStatic()
	vr.addRouteStatus()

	vr.addRouteCompiler()
}
