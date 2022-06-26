package verilogrunner

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CompileVerilogReq struct {
	Ping string
}

func (vr *VerilogRunner) processCompileVerilog(ws *websocket.Conn) error {
	for {
		var req CompileVerilogReq
		err := ws.ReadJSON(&req)
		if err != nil {
			return err
		}
		fmt.Println(req)

		err = ws.WriteJSON(struct {
			Reply string `json:"reply"`
		}{
			Reply: "Echo...",
		})
		if err != nil {
			return err
		}
	}
}

func (vr *VerilogRunner) CompileVerilog(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		closeSocketErr := ws.Close()
		if closeSocketErr != nil {
			panic(err)
		}
	}()

	err = vr.processCompileVerilog(ws)
	if err != nil {
		panic(err)
	}
}
