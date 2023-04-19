package server

import (
    "net/http"

    "mytest/internal/conf"
    "mytest/internal/service"

    "github.com/gin-gonic/gin"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, svc *service.Service) *gin.Engine {
    srv := gin.Default()
    srv.POST("/hello", OnHello)
    return srv
}

func OnHello(c *gin.Context) {
    c.JSON(http.StatusOK, "hello")
}
