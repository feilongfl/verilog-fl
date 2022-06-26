package verilogrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CompileVerilogReq struct {
	Command int
	Data    string

	exec bool
}

type CompileVerilogResp struct {
	Command int
	Data    string
	ID      string

	exec bool
}

const (
	CompileVerilogReqCommand_Ping     = iota // client -> srv
	CompileVerilogReqCommand_Pong            // srv -> client
	CompileVerilogReqCommand_Exec            // client -> srv
	CompileVerilogReqCommand_Buildlog        // srv -> client
	CompileVerilogReqCommand_Timing          // srv -> client
	CompileVerilogReqCommand_Bye             // srv -> client

	CompileVerilogReqCommand_Error
)

const buildpath = "runner"

func (cvr CompileVerilogReq) Json() []byte {
	j, err := json.Marshal(cvr)
	if err != nil {
		return nil
	}

	return j
}

func (cvr CompileVerilogReq) String() string {
	return string(cvr.Json())
}

func (cvr CompileVerilogResp) Json() []byte {
	j, err := json.Marshal(cvr)
	if err != nil {
		return nil
	}

	return j
}

func (cvr CompileVerilogResp) String() string {
	return string(cvr.Json())
}

func (cvr CompileVerilogReq) processPing(id uuid.UUID) *CompileVerilogResp {
	ret := &CompileVerilogResp{
		Command: CompileVerilogReqCommand_Pong,
		Data:    "",
		ID:      id.String(),
		exec:    cvr.exec,
	}

	return ret
}

func (cvr CompileVerilogReq) processExecCompile(ws *websocket.Conn, id uuid.UUID, path string) *CompileVerilogResp {
	var err error
	cmd := exec.Command("make", "-C", path)
	log.Println(cmd.String())

	cmd.Stdout, err = ws.NextWriter(websocket.TextMessage)
	if err != nil {
		return &CompileVerilogResp{
			Command: CompileVerilogReqCommand_Error,
			Data:    "websocket stdout writter error",
			ID:      id.String(),
			exec:    true,
		}
	}

	var stderrmsg bytes.Buffer
	cmd.Stderr = &stderrmsg

	err = cmd.Run()
	if err != nil {
		ws.WriteJSON(&CompileVerilogResp{
			Command: CompileVerilogReqCommand_Buildlog,
			Data:    stderrmsg.String(),
			ID:      id.String(),
			exec:    true,
		})

		return &CompileVerilogResp{
			Command: CompileVerilogReqCommand_Error,
			Data:    "command run error",
			ID:      id.String(),
			exec:    true,
		}
	}

	return &CompileVerilogResp{
		Command: CompileVerilogReqCommand_Bye,
		Data:    "Bye",
		ID:      id.String(),
		exec:    true,
	}
}

func (cvr CompileVerilogReq) processExec(ws *websocket.Conn, id uuid.UUID) *CompileVerilogResp {
	if cvr.exec {
		return cvr.processError(id, "already run")
	}

	// parse data to sv
	sv := cvr.Data

	// mkdir
	path := filepath.Join(".", buildpath, id.String())
	log.Printf("[%s] mkdir: %s\n", id.String(), path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Printf("[%s] error: %s\n", id.String(), err.Error())
		return &CompileVerilogResp{
			Command: CompileVerilogReqCommand_Error,
			Data:    "workspace not alloced",
			ID:      id.String(),
			exec:    true,
		}
	}

	// drop folder
	// defer func() {
	// 	log.Printf("[%s] drop dir: %s\n", id.String(), path)
	// 	err = os.RemoveAll(path)
	// 	if err != nil {
	// 		log.Printf("[%s] error: %s\n", id.String(), err.Error())
	// 	}
	// }()

	// storage file
	err = ioutil.WriteFile("output.txt", []byte(sv), 0644)
	if err != nil {
		log.Printf("[%s] error: %s\n", id.String(), err.Error())
		return &CompileVerilogResp{
			Command: CompileVerilogReqCommand_Error,
			Data:    "sv can't storage",
			ID:      id.String(),
			exec:    true,
		}
	}

	// run build command
	return cvr.processExecCompile(ws, id, path)
}

func (cvr CompileVerilogReq) processError(id uuid.UUID, msg string) *CompileVerilogResp {
	ret := &CompileVerilogResp{
		Command: CompileVerilogReqCommand_Error,
		Data:    msg,
		ID:      id.String(),
		exec:    cvr.exec,
	}

	return ret
}

func (cvr CompileVerilogReq) process(ws *websocket.Conn, id uuid.UUID) *CompileVerilogResp {
	switch cvr.Command {
	case CompileVerilogReqCommand_Ping:
		return cvr.processPing(id)

	case CompileVerilogReqCommand_Exec:
		if cvr.exec {
			return cvr.processError(id, "already run")
		} else {
			return cvr.processExec(ws, id)
		}

	default:
		return cvr.processError(id, "unknow command")
	}
}

func (vr *VerilogRunner) processCompileVerilog(ws *websocket.Conn) error {
	id := uuid.Must(uuid.NewRandom())
	exec := false // only trig once

	for {
		var req CompileVerilogReq
		err := ws.ReadJSON(&req)
		if err != nil {
			return err
		}

		req.exec = req.exec || exec

		fmt.Println(req)
		resp := req.process(ws, id)
		fmt.Println(resp)
		exec = exec || resp.exec

		err = ws.WriteJSON(resp)
		if err != nil {
			return err
		}

		if exec {
			return nil
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
