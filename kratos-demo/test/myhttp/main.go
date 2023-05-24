package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    v1 "kratos-demo/api/helloworld/v1"
    "kratos-demo/internal/base"

    "github.com/go-kratos/kratos/v2/log"
    "google.golang.org/protobuf/proto"
)

func main() {
    //post()
    //testPostHttp()
    testPostHttp2()
}

func get() {
    resp, err := http.Get("http://127.0.0.1:8000/api/abcd")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("===============%+v", string(body))
}

func post() {
    req := v1.BetReq{
        GameID: 1,
        Uid:    1,
        Data:   "hello",
    }

    reqParam, err := json.Marshal(&req)
    if err != nil {
        log.Errorf("%+v", err)
        return
    }

    url := "http://127.0.0.1:8000/api/OnBetReq"
    reqBody := strings.NewReader(string(reqParam))
    httpReq, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        log.Errorf("%+v", err)
        return
    }

    //httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    //httpReq.Header.Set("Cookie", "name=anny")
    httpReq.Header.Add("Content-Type", "application/json")

    httpRsp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        log.Errorf("%+v", err)
        return
    }

    body, err := ioutil.ReadAll(httpRsp.Body)
    if err != nil {
        log.Errorf("%+v", err)
        return
    }

    //log.Infof("%+v", string(body))

    rsp := v1.BetRsp{}
    if err = json.Unmarshal(body, &rsp); err != nil {
        log.Errorf("%+v", err)
        return
    }

    if err = proto.Unmarshal(body, &rsp); err != nil {
        log.Errorf("%+v", err)
        return
    }
    log.Infof("%+v", rsp)
}

// 修改供应商信息
func testPostHttp() error {

    // json.Marshal
    reqParam, err := json.Marshal(&v1.BetReq{
        GameID: 1001,
        Uid:    10,
        Data:   "hello",
    })
    if err != nil {
        log.Error("Marshal RequestParam fail, err:%v", err)
        return err
    }

    // 准备: HTTP请求
    url := "http://127.0.0.1:8000/api/OnBetReq"
    reqBody := strings.NewReader(string(reqParam))
    httpReq, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        log.Error("NewRequest fail, url: %s, reqBody: %s, err: %v", url, reqBody, err)
        return err
    }
    httpReq.Header.Add("Content-Type", "application/json")
    //httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    //httpReq.Header.Set("Cookie", "name=anny")

    // DO: HTTP请求
    httpRsp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        log.Error("do http fail, url: %s, reqBody: %s, err:%v", url, reqBody, err)
        return err
    }
    defer httpRsp.Body.Close()

    // Read: HTTP结果
    rspBody, err := ioutil.ReadAll(httpRsp.Body)
    if err != nil {
        log.Error("ReadAll failed, url: %s, reqBody: %s, err: %v", url, reqBody, err)
        return err
    }

    // unmarshal: 解析HTTP返回的结果
    // 		body: {"Result":{"RequestId":"12131","HasError":true,"ResponseItems":{"ErrorMsg":"错误信息"}}}
    var result v1.BetRsp
    if err = json.Unmarshal(rspBody, &result); err != nil {
        log.Error("Unmarshal fail, err:%v", err)
        return err
    }

    log.Infof("%+v", result)
    return nil
}

func testPostHttp2() {
    // json.Marshal
    reqParam, err := json.Marshal(&v1.BetReq{
        GameID: 1001,
        Uid:    10,
        Data:   "hello",
    })
    if err != nil {
        log.Error("Marshal RequestParam fail, err:%v", err)
        return
    }

    //http.Post
    url := "http://127.0.0.1:8000/api/OnBetReq"
    reqBody := strings.NewReader(string(reqParam))
    resp, err := http.Post(url, "application/json", reqBody)
    if err != nil {
        log.Error("err: %v", url, reqBody, err)
        return
    }

    //http.Post
    rspBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error("err: %v", url, reqBody, err)
        return
    }

    // json.Unmarshal
    var result v1.BetRsp
    if err = json.Unmarshal(rspBody, &result); err != nil {
        log.Error("Unmarshal fail, err:%v", err)
        return
    }

    log.Infof("%+v", &result)
    log.Infof("%+v", base.ToJSON(&result))
}
