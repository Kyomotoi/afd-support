package api

import (
	"afd-support/afdian"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func auth(token string, params string, ts int, _md5 string) bool {
	content := fmt.Sprintf("token%sparams%sts%s", token, params, strconv.Itoa(ts))
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

func (a *APIService) prepare(r *http.Request) (string, *APIRequest, error) {
	ar := &APIRequest{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, ar)
	if err != nil {
		logrus.Error("API 解析返回数据失败")
		return "请求结构有误", &APIRequest{}, err
	}

	if !auth(ar.Token, ar.Data, ar.Ts, ar.Auth) {
		logrus.Info(fmt.Sprintf("API %s(x): %s", r.Method, ar.Data))
		return "验证未通过", &APIRequest{}, errors.New("校验未通过")
	}

	logrus.Info(fmt.Sprintf("API %s(v): %s", r.Method, ar.Data))

	return "OK", ar, nil
}

func (a *APIService) GetOrders(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("API GetOrders 收到请求")
	c := &Context{
		rw:  rw,
		req: req,
	}

	msg, ar, err := a.prepare(req)
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

	msg, ar, err := a.prepare(req)
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

	msg, ar, err := a.prepare(req)
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
