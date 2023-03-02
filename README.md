# afd-support

本项目基于爱发电集合了一套订单记录/查询系统。

接受订单：通过 Webhook、用户主动请求。

查询信息：再封装 Afdian API，对外暴露一套接口，并使用签名保证信息安全，可选的限制请求源。

数据完整：当远程数据库出现问题时，自动将新数据转存至本地数据库，定时上传

## Webhook
前往 https://afdian.net/dashboard/dev Webhook 处绑定链接<br>
链接组成为：`{self.host}:{self.port}/{webhook.point}`

## 接口相关

以下所有接口一律使用以下结构，通过 `POST` 进行请求：
```json
{
    "token": ...,
    "data": ...,
    "ts": ...,
    "auth": ...
}
```
其中：
- token (string) 验证密钥
- data (string) 传入的值
- ts (int) 进行操作时的秒级时间戳
- auth (string) 签名验证

auth 的计算规则为：md5(token{token}data{data}ts{ts})

所返回的结构体可在 afdian/models.go 查看。

## 接口详情
### /order

- 为 afdian.net/api/open/query-order 可选参数 `page`、`out_trade_no` 的整合。
- 当所传入的值字符小于20时，传入参数指向 `page` 否则 `out_trade_no`。

### /sponsors

- 为 afdian.net/api/open/query-sponsor 参数 `page` 的整合。

### /getuserid

- 通过 Afdian 个人空间链接获取该用户的 user_id。
- 示例传入：`"data": "https://afdian.net/a/Kyomotoi"`
