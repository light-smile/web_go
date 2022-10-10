/*
 * @Author: jlz
 * @Date: 2022-08-24 17:13:52
 * @LastEditTime: 2022-08-31 10:57:57
 * @LastEditors: jlz
 * @Description: 测试 插入数据性能
 */
package test

import (
	"bytes"
	testModle "dnds_go/src/models/test"
	"encoding/json"
	"net/http"
	"testing"
)

const url = "127.0.0.1"

func TestHttpReq(t *testing.T) {
	data := testModle.User{
		Name: "test",
		Age:  111,
	}
	dataJson, _ := json.Marshal(data)
	reader := bytes.NewBuffer(dataJson)
	resp, _ := http.Post("http:"+url+":3000/test/add", "application/json;charset=utf-8", reader)
	defer resp.Body.Close()
}
