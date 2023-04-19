package calc

import (
    "fmt"
    "testing"
    "time"

    "calculator/internal/base"

    log "github.com/sirupsen/logrus"
)

var oneCards = []int32{
    11, 12, 13, 14, 15, 16, 17, 18, 19,
    11, 12, 13, 14, 15, 16, 17, 18, 19,
    11, 12, 13, 14, 15, 16, 17, 18, 19,
    11, 12, 13, 14, 15, 16, 17, 18, 19,

    21, 22, 23, 24, 25, 26, 27, 28, 29,
    21, 22, 23, 24, 25, 26, 27, 28, 29,
    21, 22, 23, 24, 25, 26, 27, 28, 29,
    21, 22, 23, 24, 25, 26, 27, 28, 29,

    31, 32, 33, 34, 35, 36, 37, 38, 39,
    31, 32, 33, 34, 35, 36, 37, 38, 39,
    31, 32, 33, 34, 35, 36, 37, 38, 39,
    31, 32, 33, 34, 35, 36, 37, 38, 39,

    41, 42, 43, 44, 45, 46, 47,
    41, 42, 43, 44, 45, 46, 47,
    41, 42, 43, 44, 45, 46, 47,
    41, 42, 43, 44, 45, 46, 47,
}

func TestDfs(t *testing.T) {
    var (
        m = 2
        a = []int32{0, 1, 2, 3}
    )

    res := base.Dfs(a, m)
    //INFO msg=[[0 1] [0 2] [0 3] [1 2] [1 3] [2 3]]
    log.Infof("%+v", res)
}

func TestPermute(t *testing.T) {
    start := time.Now()
    max := int32(10)

    cards := []int32{
        11, 12, 13, 18, 19,
        11, 13, 14, 16, 17, 19,
    }
    //cards := base.SliceCopy(oneCards)
    //cards = getColorPointsCards(cards, []int32{1}, []int32{1, 2, 3, 4, 5, 6, 7, 8, 9})
    fmt.Printf("%+v", descCards(cards))

    c := &tagCondition{gang: max, shun: max, ke: max, dui: max, ca: max, maxHand: 0}

    res := permute(cards, c)

    fmt.Printf("use:%+v/ms c:%+v len(cards):%+v \ncombine:%+v\n", time.Since(start).Milliseconds(), c, len(res.cards), res.info)
}

func TestPermute2(t *testing.T) {

    var (
        calculate = 0
        cards     = base.SliceCopy(oneCards)
        max       = int32(10)
    )

    for gang := int32(0); gang <= max; gang++ {
        for shun := int32(0); shun <= max; shun++ {
            for ke := int32(0); ke <= max; ke++ {
                for dui := int32(0); dui <= max; dui++ {
                    for ca := int32(0); ca <= max; ca++ {

                        //if 4*gang+3*shun+3*ke+2*dui+2*ca > _MAXCOUNTLimit {
                        //	continue
                        //}

                        calculate++
                        start := time.Now()
                        c := &tagCondition{gang: gang, shun: shun, ke: ke, dui: dui, ca: ca, maxHand: 0}
                        res := permute(cards, c)
                        if use := time.Since(start).Milliseconds(); use > 100 {
                            log.Errorf("use:%+v/ms (time over 100ms),c:%+v cards:%+v ", use, c, res.cards)
                        }

                    }
                }
            }
        }
    }

    log.Infof("count:%+v %+v", calculate, descCards(cards))
}

func TestBuild(T *testing.T) {
    var (
        count      = 0
        errCnt     = 0
        cards      = base.SliceCopy(oneCards)
        rules      = []TagLegalRule(nil)
        pointLimit = [][]int32{{1, 2, 3, 4, 5, 6, 7, 8, 9}}
    )

    c := &TagLegalConfig{
        Cards: cards,
        Rules: rules,
    }

    c.Rules = []TagLegalRule{
        {Gang: 0, Ke: 0, Shun: 0, Dui: 0, DoubleCha: 2, SingleCha: 3, LzCount: 0, BadCount: 0, RandCount: 0, ColorLimit: 1, PointLimit: pointLimit},
        {Gang: 0, Ke: 0, Shun: 0, Dui: 0, DoubleCha: 2, SingleCha: 4, LzCount: 0, BadCount: 0, RandCount: 0, ColorLimit: 1, PointLimit: pointLimit},
        {Gang: 0, Ke: 0, Shun: 0, Dui: 0, DoubleCha: 2, SingleCha: 5, LzCount: 0, BadCount: 0, RandCount: 0, ColorLimit: 1, PointLimit: pointLimit},
        {Gang: 0, Ke: 0, Shun: 0, Dui: 0, DoubleCha: 2, SingleCha: 6, LzCount: 0, BadCount: 0, RandCount: 0, ColorLimit: 1, PointLimit: pointLimit},
    }

    start := time.Now()

    err, _, _, builds := build(c)
    if err != nil {
        errCnt++
    }

    count++
    rules = []TagLegalRule(nil)

    if use := time.Since(start).Milliseconds(); use > 100 {
        fmt.Printf("\n\n\n")
        for _, v := range rules {
            need := 4*v.Gang + 3*v.Shun + 3*v.Ke + 2*v.Dui + 2*v.DoubleCha
            log.Infof("{[gang:%+v ke:%+v dui:%+v shun:%+v ca:%+v color:%+v] need:%+v len:%+v cards:%+v}",
                v.Gang, v.Ke, v.Dui, v.Shun, v.DoubleCha, v.ColorLimit, need, len(builds[v.Pid]), builds[v.Pid])
        }
        log.Errorf("use:%+v/ms (time over 100ms),c:%+v ", use, c.Rules)
    }
}

func TestBuild2(t *testing.T) {

    var (
        count  = 0
        errCnt = 0
        cards  = base.SliceCopy(oneCards)
        rules  = []TagLegalRule(nil)
        limit  = [][]int32{{1, 2, 3, 4, 5, 6, 7, 8, 9}}
        max    = int32(6)
    )

    for j := 0; j < 4; j++ {

        for colorLimit := int32(1); colorLimit <= 3; colorLimit++ {

            for gang := int32(0); gang <= max; gang++ {
                for shun := int32(0); shun <= max; shun++ {
                    for ke := int32(0); ke <= max; ke++ {
                        for dui := int32(0); dui <= max; dui++ {
                            for ca := int32(0); ca <= max; ca++ {

                                if 4*gang+3*shun+3*ke+2*dui+2*ca == 0 {
                                    continue
                                }

                                if len(rules) < 4 {
                                    rules = append(rules, TagLegalRule{
                                        Pid:        uint64(len(rules)),
                                        Seq:        0,
                                        Gang:       gang,
                                        Ke:         ke,
                                        Shun:       shun,
                                        Dui:        dui,
                                        DoubleCha:  ca,
                                        SingleCha:  0,
                                        LzCount:    0,
                                        BadCount:   0,
                                        RandCount:  0,
                                        ColorLimit: colorLimit,
                                        PointLimit: limit,
                                    })
                                    continue
                                }

                                c := &TagLegalConfig{
                                    Cards: cards,
                                    Rules: rules,
                                }

                                start := time.Now()

                                err, _, _, builds := build(c)
                                if err != nil {
                                    errCnt++
                                }

                                //fmt.Printf("\n\n\n")
                                //for _, v := range rules {
                                //	need := 4*v.Gang + 3*v.Shun + 3*v.Ke + 2*v.Dui + 2*v.DoubleCha
                                //	log.Infof("{[gang:%+v ke:%+v dui:%+v shun:%+v ca:%+v color:%+v] need:%+v len:%+v cards:%+v}",
                                //		v.Gang, v.Ke, v.Dui, v.Shun, v.DoubleCha, v.ColorLimit, need, len(builds[v.Pid]), builds[v.Pid])
                                //}

                                if use := time.Since(start).Milliseconds(); use > 100 {
                                    errCnt++
                                    fmt.Printf("\n\n\n")
                                    for _, v := range rules {
                                        need := 4*v.Gang + 3*v.Shun + 3*v.Ke + 2*v.Dui + 2*v.DoubleCha
                                        log.Infof("{[gang:%+v ke:%+v dui:%+v shun:%+v ca:%+v color:%+v] need:%+v len:%+v cards:%+v}",
                                            v.Gang, v.Ke, v.Dui, v.Shun, v.DoubleCha, v.ColorLimit, need, len(builds[v.Pid]), builds[v.Pid])
                                    }
                                    log.Errorf("use:%+v/ms count:%+v errCnt:%+v", use, count, errCnt)
                                }

                                //if count > 100 {
                                //	return
                                //}

                                count++
                                rules = []TagLegalRule(nil)

                            }
                        }
                    }
                }

            }
        }

    }

    defer log.Infof("====> count:%+v errCnt:%+v", count, errCnt)
}
