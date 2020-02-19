package test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os/exec"
)

type Helper struct {
	cmd *exec.Cmd
	output *bytes.Buffer
}

func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) Start(command string) string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("gen uuid failed: %v", err)
	}

	uid := base64.URLEncoding.EncodeToString(b)

	socket := "/tmp/uds-"+uid+".sock"

	output := &bytes.Buffer{}

	cmd := exec.Command(command)
	cmd.Env = append(cmd.Env, "DUNDERGITCALL_UDS="+socket)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("start app failed: %v", err)
	}
	cmd.Stdout = output
	cmd.Stderr = output

	h.cmd = cmd
	h.output = output

	return socket
}

func (h *Helper) Logs() string {
	if h.cmd == nil {
		return ""
	}

	out, err := ioutil.ReadAll(h.output)
	if err != nil {
		log.Fatalf("read logs failed: %v", err)
	}

	return string(out)
}

func (h *Helper) Stop() {
	if h.cmd == nil {
		return
	}

	err := h.cmd.Process.Kill()
	if err != nil {
		log.Printf("failed to stop app: %v", err)
	}
}
