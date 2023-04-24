package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/xjieinfo/xjgo/xjcore/xjexcel"
)

func main() {
    router := gin.Default()
    router.GET("/dump", dump)
    router.Run(":8000")
}

// DumpReq url绑定参数（也可以绑定到body里）
type DumpReq struct {
    GameID int32 `form:"gameID" binding:"required"`
    Uid    int32 `form:"uid" binding:"required"`
}

type User struct {
    Name    string `excel:"column:B;desc:姓名;width:30"`
    Age     int    `excel:"column:C;desc:年龄;width:10"`
    Address string `excel:"column:D;desc:地址;width:50"`
}

// curl -X GET "http://127.0.0.1:8000/dump"
func dump(c *gin.Context) {

    req := DumpReq{
        GameID: 0,
        Uid:    0,
    }

    //if err := c.ShouldBindJSON(req); err != nil {
    //    log.Printf("err:%+v", err)
    //    return
    //}

    // curl -X GET "http://127.0.0.1:8000/dump?gameID=1004&uid=1"
    if err := c.ShouldBindQuery(&req); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }

    path := "./temp/"
    fileName := fmt.Sprintf("%+v-%+v-%+v.xls", req.GameID, req.Uid, time.Now().Format("20060102-150405"))
    filePathName := fmt.Sprintf("%+v%+v", path, fileName)

    log.Printf("%+v %+v", fileName, filePathName)

    fileInfo, err := os.Stat(path)
    if err != nil {
        log.Printf("File does exist. File information: %+v", fileInfo)
        if err = os.MkdirAll(path, 0777); err != nil {
            fmt.Printf("创建文件夹失败")
            return
        }
    }

    list := []User{
        {
            Name:    "张三",
            Age:     18,
            Address: "北京东三环",
        },
        {
            Name:    "李四",
            Age:     21,
            Address: "上海人民路",
        },
        {
            Name:    "王五",
            Age:     22,
            Address: "长沙开福区",
        },
    }
    f := xjexcel.ListToExcel(list, "员工信息表", "员工表")
    _ = f.SaveAs(filePathName)

    c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
    c.Writer.Header().Add("Content-Type", "application/octet-stream")
    c.File(filePathName)
}
