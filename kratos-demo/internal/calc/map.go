package calc

import (
    "math/rand"
    "sync/atomic"

    "kratos-demo/internal/base"
)

type Map struct {
    points [][]int //5*5

    ////计算次数
    //betCnt int

    //挖矿次数
    digCnt int

    ////函数处理池
    //pool *Pool

    _winCnt   int64 //= int64(0)
    _loseCnt  int64 //= int64(0)
    _totalCnt int64 //= int64(0)
}

type stPoint struct {
    x   int
    y   int
    val int
}

type stInput struct{}

func New(digCnt int) (m *Map) {
    m = &Map{}

    m.points = [][]int{
        {0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0},
    }
    for i := 0; i < len(m.points); i++ {
        for j := 0; j < len(m.points[i]); j++ {
            m.points[i][j] = base.RandRange(0, 3)
        }
    }

    m.digCnt = digCnt

    //m.pool = NewPool()

    return m
}

func (m *Map) getRandomPoints(n int) (points []stPoint) {

    //打乱
    for i := 0; i < len(m.points); i++ {
        for j := 0; j < len(m.points[i]); j++ {
            points = append(points, stPoint{
                x:   i,
                y:   j,
                val: m.points[i][j],
            })
        }
    }
    for i := len(points) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        points[i], points[j] = points[j], points[i]
    }

    //取n个随机的点
    if totalCnt := len(m.points) * len(m.points[0]); n > totalCnt {
        n = totalCnt - 1
    }
    points = points[:n]
    return
}

func (m *Map) calc() (value int) {
    nPoints := m.getRandomPoints(m.digCnt)
    for _, v := range nPoints {
        value += v.val
    }

    atomic.AddInt64(&m._totalCnt, 1)
    if base.IsHit(30) {
        atomic.AddInt64(&m._winCnt, 1)
    } else {
        atomic.AddInt64(&m._loseCnt, 1)
    }

    return value
}
