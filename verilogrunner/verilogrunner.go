package verilogrunner

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type VerilogRunner struct {
	server *gin.Engine
	addr   string
}

func NewVerilogRunner(cmd *cobra.Command) *VerilogRunner {
	ret := &VerilogRunner{}

	ret.server = gin.Default()
	ret.addr = cmd.Flags().Lookup("listen").Value.String()

	return ret
}

func (vr *VerilogRunner) Run() {
	vr.addRoute()

	vr.server.Run(vr.addr)
}

func RunVerilogRunner(cmd *cobra.Command, args []string) {
	vr := *NewVerilogRunner(cmd)

	vr.Run()
}
