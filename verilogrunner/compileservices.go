package verilogrunner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (vr *VerilogRunner) CompileVerilog(c *gin.Context) {
	id := uuid.Must(uuid.NewRandom()).String()

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"id":     id,
	})
}
