package tests

import (
	"io/ioutil"
	"net"
	"testing"
	"time"

	"git.corezoid.com/gitcall/go-runner/test"
	"github.com/stretchr/testify/assert"
)

func TestShouldHandleTaskSuccessfully(t *testing.T)  {
	request := `{ "method": "Usercode.Run", "params": [{"data": {"key": "value"}}], "id":1}\\n`
	expectedResponse := `{"id":1,"result":{"data":{"foo":123,"key":"value"}},"error":null}
`

	helper := test.NewHelper()
	socket := helper.Start("../../build/gorunner-success")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("unix", socket)
		if err != nil {
			attempts--
			if attempts <= 0 {
				t.Log(helper.Logs())
				t.Fatalf("connection failed: %v", err)
			}

			<-time.NewTimer(time.Second).C

			continue
		}

		break
	}

	_, err = conn.Write([]byte(request))
	assert.NoError(t, err)

	resp, err := ioutil.ReadAll(conn)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, string(resp))
}

func TestShouldHandleTaskError(t *testing.T)  {
	request := `{ "method": "Usercode.Run", "params": [{"data": {"key": "value"}}], "id":1}\\n`
	expectedResponse := `{"id":1,"result":null,"error":"error-happened"}
`

	helper := test.NewHelper()
	socket := helper.Start("../../build/gorunner-error")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("unix", socket)
		if err != nil {
			attempts--
			if attempts <= 0 {
				t.Log(helper.Logs())
				t.Fatalf("connection failed: %v", err)
			}

			<-time.NewTimer(time.Second).C

			continue
		}

		break
	}

	_, err = conn.Write([]byte(request))
	assert.NoError(t, err)

	resp, err := ioutil.ReadAll(conn)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, string(resp))
}
