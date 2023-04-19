package calc

import (
    "testing"
    "time"

    "calculator/internal/base"

    "github.com/go-kratos/kratos/v2/log"
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
    cards := []int32{
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

    cards = getColorPointsCards(cards, []int32{1, 2}, []int32{1, 2, 3})

    c := &tagCondition{gang: 3, shun: 4, ke: 4, dui: 4, ca: 4, maxHand: 0}

    permute(cards, c)

    log.Infof("use:%+v/ms", time.Since(start).Milliseconds())
}

func TestPermute2(t *testing.T) {

    var (
        count = 1

        calculate = 0

        start = time.Now()

        cards = []int32{
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
    )

    for i := 0; i < count; i++ {

        for gang := int32(0); gang <= 4; gang++ {
            for shun := int32(0); shun <= 4; shun++ {
                for ke := int32(0); ke <= 4; ke++ {
                    for dui := int32(0); dui <= 6; dui++ {
                        for ca := int32(0); ca <= 6; ca++ {

                            if 4*gang+3*shun+3*ke+2*dui+2*ca > MAXCOUNT {
                                continue
                            }

                            calculate++

                            c := &tagCondition{gang: gang, shun: shun, ke: ke, dui: dui, ca: ca, maxHand: 0}
                            permute(cards, c)

                        }
                    }
                }
            }
        }

    }

    log.Infof("count:%+v calculate:%+v use:%+v/ms", count, calculate, time.Since(start).Milliseconds())
    log.Infof("------------------")
}

func TestBuild(t *testing.T) {

    var (
        count  = 0
        errCnt = 0
        cards  = base.SliceCopy(oneCards)
        rules  = []TagLegalRule(nil)
    )

    for j := 0; j < 4; j++ {

        for colorLimit := int32(0); colorLimit <= 3; colorLimit++ {

            for gang := int32(0); gang <= 4; gang++ {
                for shun := int32(0); shun <= 4; shun++ {
                    for ke := int32(0); ke <= 4; ke++ {
                        for dui := int32(0); dui <= 6; dui++ {
                            for ca1 := int32(0); ca1 <= 6; ca1++ {

                                for ca2 := int32(0); ca2 <= 6; ca2++ {

                                    if 4*gang+3*shun+3*ke+2*dui+2*ca1+2*ca2 == 0 {
                                        continue
                                    }

                                    if 4*gang+3*shun+3*ke+2*dui+2*ca1+2*ca2 > MAXCOUNT {
                                        continue
                                    }

                                    if len(rules) < 4 {
                                        rules = append(rules, TagLegalRule{
                                            Pid:        0,
                                            Seq:        0,
                                            Gang:       gang,
                                            Ke:         ke,
                                            Shun:       shun,
                                            Dui:        dui,
                                            DoubleCha:  ca1,
                                            SingleCha:  ca2,
                                            LzCount:    0,
                                            BadCount:   0,
                                            RandCount:  0,
                                            ColorLimit: colorLimit,
                                            PointLimit: [][]int32{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
                                        })
                                        continue
                                    }

                                    c := &TagLegalConfig{
                                        Cards: cards,
                                        Rules: rules,
                                    }

                                    if err, _, _, _ := build(c); err != nil {
                                        errCnt++
                                    }

                                    count++
                                    rules = []TagLegalRule(nil)
                                }
                            }
                        }
                    }
                }

            }
        }

    }

    log.Infof("count:%+v errCnt:%+v", count, errCnt)
}
