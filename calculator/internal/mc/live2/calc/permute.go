/*
	递归找出需要的组合
*/

package calc

import (
    "calculator/internal/base"

    log "github.com/sirupsen/logrus"
)

const MAXCOUNT = 13
const _MaxCountLimit = 20 //组合最多20张 (限制数量减少更多的计算)

var _OpenDebugLog = true

type tagResult struct {
    info *tagComInfo
}
type tagComInfo struct {
    //序列
    cards []int32
    //组合
    gang [][]int32
    shun [][]int32
    ke   [][]int32
    dui  [][]int32
    ca   [][]int32 //茬
}

//配置参数
type tagCondition struct {
    gang int32 //杠个数
    shun int32 //顺子个数
    ke   int32 //刻子个数
    dui  int32 //对子个数
    ca   int32 //茬个数

    maxHand int //组合最大张数 (用于递归剪枝)
}

func (c *tagComInfo) clone() (t *tagComInfo) {
    t = &tagComInfo{}
    t.gang = sliceCopy(c.gang)
    t.ke = sliceCopy(c.ke)
    t.shun = sliceCopy(c.shun)
    t.dui = sliceCopy(c.dui)
    t.ca = sliceCopy(c.ca)
    t.cards = sliceMerge(t.gang, t.shun, t.ke, t.dui, t.ca)
    return t
}

func permute(cards []int32, c *tagCondition) *tagResult {

    m := map[int32]int32{}
    for _, v := range cards {
        m[v]++
    }

    cList := []int32(nil)
    for card, _ := range m {
        cList = append(cList, card)
    }
    cList = base.SliceSortRise(cList)

    c.maxHand = base.MInInt(int(c.gang*4+c.ke*3+c.shun*3+c.dui*2+c.ca*2), len(cards), _MaxCountLimit)
    ////剪枝 剪掉不合理的条件判断
    //{
    //
    //
    //    //c.gang = base.MInInt32(c.gang, _MaxCountLimit/4, cnt[4])
    //    //c.ke = base.MInInt32(c.ke, _MaxCountLimit/3, cnt[4]+cnt[3])
    //    //c.shun = base.MInInt32(c.shun, _MaxCountLimit/3, int32(len(cards))/3)
    //    //c.dui = base.MInInt32(c.dui, _MaxCountLimit/2, 2*cnt[4]+cnt[3]+cnt[2])
    //    //c.ca = base.MInInt32(c.ca, _MaxCountLimit/2, int32(len(cards))/2)
    //    c.maxHand = base.MInInt(int(c.gang*4+c.ke*3+c.shun*3+c.dui*2+c.ca*2), len(cards), _MaxCountLimit)
    //}

    res := &tagResult{info: &tagComInfo{}}
    cache := &tagComInfo{}

    //colorMap := map[int32][]int32{}
    //for _, v := range cards {
    //    color := toColor(v)
    //    colorMap[color] = append(colorMap[color], v)
    //}
    //for _, colorcards := range colorMap {
    //    tm := map[int32]int32{}
    //    for _, v := range colorcards {
    //        tm[v]++
    //    }
    //    backtracking(tm, c, res, cache)
    //}

    backtracking(m, cList, c, res, cache, 0)

    if len(res.info.cards) > 0 && !base.SliceContainAll(cards, res.info.cards...) {
        log.Errorf("=> error. !base.SliceContainAll(cards,special...) %+v %+v", descCards(cards), descCards(res.info.cards))
        res.info = &tagComInfo{}
        return res
    }

    return res
}

func backtracking(m map[int32]int32, cList []int32, c *tagCondition, res *tagResult, cache *tagComInfo, index int32) {

    if index >= int32(len(cList)) {
        return
    }
    // 找到一组就退出
    if len(res.info.cards) >= c.maxHand {
        return
    }

    //成功找到一组
    cnt := len(cache.gang)*4 + (len(cache.ke)+len(cache.shun))*3 + (len(cache.dui)+len(cache.ca))*2
    if cnt > len(res.info.cards) {
        tmp := cache.clone()
        res.info = tmp
    }

    for i := index; i < int32(len(cList)); i++ {

        k := cList[i]

        //杠
        if m[k] >= 4 && c.gang > 0 && int(c.gang) > len(cache.gang) {
            m[k] -= 4
            cache.gang = append(cache.gang, []int32{k, k, k, k})
            if m[k] > 0 {
                backtracking(m, cList, c, res, cache, i)
            } else {
                backtracking(m, cList, c, res, cache, i+1)
            }

            m[k] += 4
            cache.gang = cache.gang[:len(cache.gang)-1]
        }

        //刻
        if m[k] >= 3 && c.ke > 0 && int(c.ke) > len(cache.ke) {
            m[k] -= 3
            cache.ke = append(cache.ke, []int32{k, k, k})
            if m[k] > 0 {
                backtracking(m, cList, c, res, cache, i)
            } else {
                backtracking(m, cList, c, res, cache, i+1)
            }
            m[k] += 3
            cache.ke = cache.ke[:len(cache.ke)-1]
        }

        //对
        if m[k] >= 2 && c.dui > 0 && int(c.dui) > len(cache.dui) {
            m[k] -= 2
            cache.dui = append(cache.dui, []int32{k, k})
            if m[k] > 0 {
                backtracking(m, cList, c, res, cache, i)
            } else {
                backtracking(m, cList, c, res, cache, i+1)
            }
            m[k] += 2
            cache.dui = cache.dui[:len(cache.dui)-1]
        }

        if m[k] >= 1 && c.shun > 0 && int(c.shun) > len(cache.shun) {
            //顺
            if m[k] >= 1 && m[k+1] >= 1 && m[k+2] >= 1 &&
                    toColor(k) == toColor(k+1) && toColor(k+1) == toColor(k+2) {
                m[k]--
                m[k+1]--
                m[k+2]--
                cache.shun = append(cache.shun, []int32{k, k + 1, k + 2})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k+1]++
                m[k+2]++
                cache.shun = cache.shun[:len(cache.shun)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 && m[k+1] >= 1 &&
                    toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k+1) {
                m[k]--
                m[k-1]--
                m[k+1]--
                cache.shun = append(cache.shun, []int32{k, k - 1, k + 1})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k-1]++
                m[k+1]++
                cache.shun = cache.shun[:len(cache.shun)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 && m[k-2] >= 1 &&
                    toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k-2) {
                m[k]--
                m[k-1]--
                m[k-2]--
                cache.shun = append(cache.shun, []int32{k, k - 1, k - 2})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k-1]++
                m[k-2]++
                cache.shun = cache.shun[:len(cache.shun)-1]
            }
        }

        if m[k] >= 1 && c.ca > 0 && int(c.ca) > len(cache.ca) {
            //茬
            if m[k] >= 1 && m[k+1] >= 1 &&
                    toColor(k) == toColor(k+1) {
                m[k]--
                m[k+1]--
                cache.ca = append(cache.ca, []int32{k, k + 1})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k+1]++
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if m[k] >= 1 && m[k+2] >= 1 &&
                    toColor(k) == toColor(k+2) {
                m[k]--
                m[k+2]--
                cache.ca = append(cache.ca, []int32{k, k + 2})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k+2]++
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 &&
                    toColor(k) == toColor(k-1) {
                m[k]--
                m[k-1]--
                cache.ca = append(cache.ca, []int32{k, k - 1})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k-1]++
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if m[k] >= 1 && m[k-2] >= 1 &&
                    toColor(k) == toColor(k-2) {
                m[k]--
                m[k-2]--
                cache.ca = append(cache.ca, []int32{k, k - 2})
                if m[k] > 0 {
                    backtracking(m, cList, c, res, cache, i)
                } else {
                    backtracking(m, cList, c, res, cache, i+1)
                }
                m[k]++
                m[k-2]++
                cache.ca = cache.ca[:len(cache.ca)-1]
            }
        }
    }

    //for k, v := range m {
    //    if v <= 0 {
    //        continue
    //    }
    //
    //    //杠
    //    if m[k] >= 4 && c.gang > 0 && int(c.gang) > len(cache.gang) {
    //        m[k] -= 4
    //        cache.gang = append(cache.gang, []int32{k, k, k, k})
    //        backtracking(m, c, res, cache)
    //        m[k] += 4
    //        cache.gang = cache.gang[:len(cache.gang)-1]
    //    }
    //
    //    //刻
    //    if m[k] >= 3 && c.ke > 0 && int(c.ke) > len(cache.ke) {
    //        m[k] -= 3
    //        cache.ke = append(cache.ke, []int32{k, k, k})
    //        backtracking(m, c, res, cache)
    //        m[k] += 3
    //        cache.ke = cache.ke[:len(cache.ke)-1]
    //    }
    //
    //    //对
    //    if m[k] >= 2 && c.dui > 0 && int(c.dui) > len(cache.dui) {
    //        m[k] -= 2
    //        cache.dui = append(cache.dui, []int32{k, k})
    //        backtracking(m, c, res, cache)
    //        m[k] += 2
    //        cache.dui = cache.dui[:len(cache.dui)-1]
    //    }
    //
    //    if m[k] >= 1 && c.shun > 0 && int(c.shun) > len(cache.shun) {
    //        //顺
    //        if m[k] >= 1 && m[k+1] >= 1 && m[k+2] >= 1 &&
    //                toColor(k) == toColor(k+1) && toColor(k+1) == toColor(k+2) {
    //            m[k]--
    //            m[k+1]--
    //            m[k+2]--
    //            cache.shun = append(cache.shun, []int32{k, k + 1, k + 2})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k+1]++
    //            m[k+2]++
    //            cache.shun = cache.shun[:len(cache.shun)-1]
    //        }
    //
    //        if m[k] >= 1 && m[k-1] >= 1 && m[k+1] >= 1 &&
    //                toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k+1) {
    //            m[k]--
    //            m[k-1]--
    //            m[k+1]--
    //            cache.shun = append(cache.shun, []int32{k, k - 1, k + 1})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k-1]++
    //            m[k+1]++
    //            cache.shun = cache.shun[:len(cache.shun)-1]
    //        }
    //
    //        if m[k] >= 1 && m[k-1] >= 1 && m[k-2] >= 1 &&
    //                toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k-2) {
    //            m[k]--
    //            m[k-1]--
    //            m[k-2]--
    //            cache.shun = append(cache.shun, []int32{k, k - 1, k - 2})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k-1]++
    //            m[k-2]++
    //            cache.shun = cache.shun[:len(cache.shun)-1]
    //        }
    //    }
    //
    //    if m[k] >= 1 && c.ca > 0 && int(c.ca) > len(cache.ca) {
    //        //茬
    //        if m[k] >= 1 && m[k+1] >= 1 &&
    //                toColor(k) == toColor(k+1) {
    //            m[k]--
    //            m[k+1]--
    //            cache.ca = append(cache.ca, []int32{k, k + 1})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k+1]++
    //            cache.ca = cache.ca[:len(cache.ca)-1]
    //        }
    //
    //        if m[k] >= 1 && m[k+2] >= 1 &&
    //                toColor(k) == toColor(k+2) {
    //            m[k]--
    //            m[k+2]--
    //            cache.ca = append(cache.ca, []int32{k, k + 2})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k+2]++
    //            cache.ca = cache.ca[:len(cache.ca)-1]
    //        }
    //
    //        if m[k] >= 1 && m[k-1] >= 1 &&
    //                toColor(k) == toColor(k-1) {
    //            m[k]--
    //            m[k-1]--
    //            cache.ca = append(cache.ca, []int32{k, k - 1})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k-1]++
    //            cache.ca = cache.ca[:len(cache.ca)-1]
    //        }
    //
    //        if m[k] >= 1 && m[k-2] >= 1 &&
    //                toColor(k) == toColor(k-2) {
    //            m[k]--
    //            m[k-2]--
    //            cache.ca = append(cache.ca, []int32{k, k - 2})
    //            backtracking(m, c, res, cache)
    //            m[k]++
    //            m[k-2]++
    //            cache.ca = cache.ca[:len(cache.ca)-1]
    //        }
    //    }
    //
    //}
}

func sliceCopy(s [][]int32) (result [][]int32) {
    result = make([][]int32, len(s))
    copy(result, s)
    return
}
func sliceMerge(s ...[][]int32) (result []int32) {
    for _, v := range s {
        for _, vv := range v {
            result = append(result, vv...)
        }
    }
    return
}
