package main

import (
    "net/http"

    v1 "kratos-demo/test/myproto/api"

    "github.com/gin-gonic/gin"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/swagger-api/openapiv2"
)

func main() {
    //r.HandlePrefix("/q/", openAPIhandler) //http://127.0.0.1:8000/q/swagger-ui/
    log.Infof("%+v", "http://127.0.0.1:8000/q/swagger-ui/")
    http.Handle("/q", openapiv2.NewHandler())
    //http.HandleFunc("/", homeHandler)
    http.ListenAndServe(":8000", nil)

}

//func homeHandler(w http.ResponseWriter, r *http.Request) {
//    fmt.Fprintf(w, "Hello, Golang!")
//}

func main3() {
    r := gin.Default()
    r.POST("/dump", dump)

    //
    //h.Handler = openAPIhandler

    log.Fatal(r.Run(":8000"))

}
func dump(c *gin.Context) {
    req := v1.HelloRequest{}
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Errorf("dump fail. %+v", err)
        c.JSON(http.StatusOK, &v1.HelloReply{
            Code:    1,
            Msg:     err.Error(),
            Message: "",
        })
    }

    c.JSON(http.StatusOK, &v1.HelloReply{
        Code:    0,
        Msg:     "",
        Message: req.Name,
    })
}

////https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code
//
////protoc --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative   greeter.proto

////protoc --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative   greeter.proto

//protoc -I=. -I=$GOPATH/pkg/mod/http://github.com/gogo/protobuf@latest --gofast_out=plugins=grpc:. *.proto
