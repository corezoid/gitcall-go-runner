package tests

import (
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/corezoid/gitcall-go-runner/test"
	"github.com/stretchr/testify/assert"
)

func TestShouldHandleTaskSuccessfully(t *testing.T) {
	request := `{ "method": "Usercode.Run", "params": [{"data": {"case":"success","key": "value"}}], "id":1}\\n`
	expectedResponse := `{"id":1,"result":{"data":{"case":"success","foo":123,"key":"value"}},"error":null}
`

	helper := test.NewHelper()
	addr := helper.Start("../../build/testrunner")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("tcp", addr)
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

func TestShouldHandleInternalError(t *testing.T) {
	request := `{ "method": "Usercode.Run", "params": [], "id":1}\\n`
	expectedResponse := `{"id":1,"result":null,"error":"'data' was not found"}
`

	helper := test.NewHelper()
	addr := helper.Start("../../build/testrunner")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("tcp", addr)
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

func TestShouldHandleTaskError(t *testing.T) {
	request := `{ "method": "Usercode.Run", "params": [{"data": {"case":"error","key": "value"}}], "id":1}\\n`
	expectedResponse := `{"id":1,"result":null,"error":"usercode:error-happened"}
`

	helper := test.NewHelper()
	addr := helper.Start("../../build/testrunner")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("tcp", addr)
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

func TestShouldHandleTaskPanic(t *testing.T) {
	request := `{ "method": "Usercode.Run", "params": [{"data": {"case":"panic","key": "value"}}], "id":1}\\n`
	expectedResponse := `{"id":1,"result":null,"error":"usercode:something went wrong"}
`

	helper := test.NewHelper()
	addr := helper.Start("../../build/testrunner")
	defer helper.Stop()

	var conn net.Conn
	var err error

	attempts := 10
	for {
		conn, err = net.Dial("tcp", addr)
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
