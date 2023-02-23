package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Em string `json:"em"`
}

type Context struct {
	rw  http.ResponseWriter
	req *http.Request
}

func (c *Context) write(status int, b []byte) {
	c.rw.WriteHeader(status)
	_, err := c.rw.Write(b)
	if err != nil {
		logrus.Error(err)
	}
}

func (c *Context) Send(status int, s any) {
	var con []byte
	if reflect.TypeOf(s).Name() == "string" {
		con, _ = json.Marshal(&Response{
			Em: fmt.Sprint(s),
		})
	} else {
		con = []byte(fmt.Sprint(s))
	}
	c.write(status, con)
}
