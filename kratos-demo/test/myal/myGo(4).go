package main

import (
    "fmt"

    "github.com/garyburd/redigo/redis"
)

// 12:00:00
func main2() {
    getRobotAndAddMoney()
}

func getUserIdByName() {
    conn, err := redis.Dial("tcp", "172.13.3.51:30071") // 172.13.3.51:30071 南美 172.13.3.51:30171 土耳其
    if err != nil {
        fmt.Println("connect redis error :", err)
        return
    }
    defer conn.Close()
    for i := 1000000; i < 1070000; i++ {
        key := fmt.Sprintf("HPlayer:%d", i)
        name, err := redis.String(conn.Do("HGET", key, "Nick"))
        if err != nil {
            continue
        }
        if name == "Belle Anu" {
            fmt.Println(i)
            break
        }
    }
    fmt.Println("end")
}

func getRobotAndAddMoney() {
    conn, err := redis.Dial("tcp", "172.13.3.51:30071")
    if err != nil {
        fmt.Println("connect redis error :", err)
        return
    }
    defer conn.Close()
    pid := []int{}
    rid := []int{}
    for i := 1000000; i < 1070000; i++ {
        key := fmt.Sprintf("HPlayer:%d", i)
        isRobot, err := redis.Bool(conn.Do("HGET", key, "IsRobot"))
        if err != nil {
            continue
        }
        if isRobot == true {
            rid = append(rid, i)
            key1 := fmt.Sprintf("HProp:%d", i)
            conn.Do("HSET", key1, 2, 10000000000)
        } else {
            pid = append(pid, i)
        }
    }

    show(rid)
}

func show(rid []int) {
    fmt.Printf("[")

    for _, v := range rid {
        fmt.Printf("%d,", v)
    }

    fmt.Printf("]\n")
}
