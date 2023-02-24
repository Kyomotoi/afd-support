package api

import (
	"afd-support/afdian"
	"afd-support/configs"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func auth(token string, data string, ts int, _md5 string) bool {
	content := fmt.Sprintf("token%sdata%sts%s", token, data, strconv.Itoa(ts))
	has := fmt.Sprintf("%s", md5.Sum([]byte(content)))
	return has == _md5
}

func New(token string, AfdianToken string, AfdianUserID string) *APIService {
	return &APIService{
		Token: token,
		AfdianItem: &afdian.AfdianAPIService{
			APIToken: AfdianToken,
			UserID:   AfdianUserID,
		},
	}
}

func (a *APIService) prepare(c *Context) (string, *APIRequest, error) {
	if configs.Conf.API.IsLimitHost {
		if configs.Conf.API.Only != strings.Split(c.req.RemoteAddr, ":")[0] {
			logrus.Warn(fmt.Sprintf("外界请求 %s 已拦截", c.req.Host))
			return "拦截外界请求", &APIRequest{}, errors.New("out of limit of config")
		}
	}

	ar := &APIRequest{}
	body, _ := io.ReadAll(c.req.Body)
	err := json.Unmarshal(body, ar)
	if err != nil {
		logrus.Error("API 解析返回数据失败")
		return "请求结构有误", &APIRequest{}, err
	}

	if !auth(ar.Token, ar.Data, ar.Ts, ar.Auth) {
		logrus.Info(fmt.Sprintf("API %s(x): %s", c.req.Method, ar.Data))
		return "验证未通过", &APIRequest{}, errors.New("校验未通过")
	}

	logrus.Info(fmt.Sprintf("API %s(v): %s", c.req.Method, ar.Data))

	return "OK", ar, nil
}

func (a *APIService) GetOrders(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("API GetOrders 收到请求")
	c := &Context{
		rw:  rw,
		req: req,
	}

	msg, ar, err := a.prepare(c)
	if err != nil {
		c.Send(403, msg)
		return
	}

	resp, err := a.AfdianItem.QueryOrder(ar.Data)
	if err != nil {
		c.Send(500, "Afdian 返回信息处理失败")
		return
	}

	data, _ := json.Marshal(resp)
	c.Send(200, string(data))
}

func (a *APIService) GetSponsorsWithPage(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("API GetSponsorsWithPage 收到请求")
	c := &Context{
		rw:  rw,
		req: req,
	}

	msg, ar, err := a.prepare(c)
	if err != nil {
		c.Send(403, msg)
		return
	}

	resp, err := a.AfdianItem.QuerySponsorWithPage(ar.Data)
	if err != nil {
		c.Send(500, "Afdian 返回信息处理失败")
		return
	}

	data, _ := json.Marshal(resp)
	c.Send(200, string(data))
}

func (a *APIService) GetUserIDByURL(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("API GetUserIDByProfileURL 收到请求")
	c := &Context{
		rw:  rw,
		req: req,
	}

	msg, ar, err := a.prepare(c)
	if err != nil {
		c.Send(403, msg)
		return
	}

	resp, err := afdian.GetUserIDbyProfileURL(ar.Data)
	if err != nil {
		c.Send(500, "Afdian 返回信息处理失败")
		return
	}

	data, _ := json.Marshal(resp)
	c.Send(200, string(data))
}
