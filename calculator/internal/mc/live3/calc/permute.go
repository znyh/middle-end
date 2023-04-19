/*
	递归找出需要的组合
*/

package calc

import (
    "calculator/internal/base"

    log "github.com/sirupsen/logrus"
)

const _MaxCountLimit = 20 //组合最多20张 (限制数量减少更多的计算)

var _OpenDebugLog = true

type tagResult struct {
    cards []int32
    info  [][]int32
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

func permute(cards []int32, c *tagCondition) *tagResult {

    _use := []int32(nil)
    _com := [][]int32(nil)
    _total := base.SliceCopy(cards)
    _maxHand := base.MInInt(int(c.gang*4+c.ke*3+c.shun*3+c.dui*2+c.ca*2), len(cards), _MaxCountLimit)

    if c.gang > 0 {
        m, cnt := toMap(_total)
        result := &tagResult{}
        max := base.MInInt32(c.gang, _MaxCountLimit/4, int32(len(_total)/4), cnt[4])
        backtracking(m, &tagCondition{gang: max, maxHand: int(max * 4)}, result, [][]int32(nil))
        if cs := result.cards; len(cs) > 0 {
            _use = append(_use, cs...)
            _com = append(_com, result.info...)
            _total = base.SliceDel(_total, cs...)
        }
        if len(_use) >= _maxHand {
            return &tagResult{cards: _use, info: _com}
        }
    }
    if c.ke > 0 {
        m, cnt := toMap(_total)
        result := &tagResult{}
        cache := [][]int32(nil)
        max := base.MInInt32(c.ke, _MaxCountLimit/3, int32(len(_total)/3), cnt[4]+cnt[3])
        backtracking(m, &tagCondition{ke: max, maxHand: int(max * 3)}, result, cache)
        if cs := result.cards; len(cs) > 0 {
            _use = append(_use, cs...)
            _com = append(_com, result.info...)
            _total = base.SliceDel(_total, cs...)
        }
        if len(_use) >= _maxHand {
            return &tagResult{cards: _use, info: _com}
        }
    }
    if c.shun > 0 {
        m, cnt := toMap(_total)
        result := &tagResult{}
        cache := [][]int32(nil)
        max := base.MInInt32(c.shun, _MaxCountLimit/3, int32(len(_total)/3), cnt[3]+cnt[4])
        backtracking(m, &tagCondition{shun: max, maxHand: int(max * 3)}, result, cache)
        if cs := result.cards; len(cs) > 0 {
            _use = append(_use, cs...)
            _com = append(_com, result.info...)
            _total = base.SliceDel(_total, cs...)
        }
        if len(_use) >= _maxHand {
            return &tagResult{cards: _use, info: _com}
        }
    }
    if c.dui > 0 {
        m, cnt := toMap(_total)
        result := &tagResult{}
        cache := [][]int32(nil)
        max := base.MInInt32(c.dui, _MaxCountLimit/2, int32(len(_total)/2), 2*cnt[4]+cnt[3]+cnt[2])
        backtracking(m, &tagCondition{dui: max, maxHand: int(max * 2)}, result, cache)
        if cs := result.cards; len(cs) > 0 {
            _use = append(_use, cs...)
            _com = append(_com, result.info...)
            _total = base.SliceDel(_total, cs...)
        }
        if len(_use) >= _maxHand {
            return &tagResult{cards: _use, info: _com}
        }
    }
    if c.ca > 0 {
        m, cnt := toMap(_total)
        result := &tagResult{}
        cache := [][]int32(nil)
        max := base.MInInt32(c.ca, _MaxCountLimit/2, int32(len(_total)/2), 2*cnt[4]+cnt[3]+cnt[2])
        backtracking(m, &tagCondition{ca: max, maxHand: int(max * 2)}, result, cache)
        if cs := result.cards; len(cs) > 0 {
            _use = append(_use, cs...)
            _com = append(_com, result.info...)
            _total = base.SliceDel(_total, cs...)
        }
        if len(_use) >= _maxHand {
            return &tagResult{cards: _use, info: _com}
        }
    }

    if len(_use) > 0 && !base.SliceContainAll(cards, _use...) {
        log.Errorf("=> error. !base.SliceContainAll(cards,special...) %+v %+v", descCards(cards), descCards(_use))
        return &tagResult{}
    }

    return &tagResult{cards: _use, info: _com}
}

func permute2(cards []int32, c *tagCondition) *tagResult {

    //剪掉不合理的条件判断
    c.maxHand = int(c.gang*4 + c.ke*3 + c.shun*3 + c.dui*2 + c.ca*2)
    c.maxHand = base.MInInt(c.maxHand, len(cards), _MaxCountLimit)

    m := map[int32]int32{}
    for _, v := range cards {
        m[v]++
    }

    res := &tagResult{}
    cache := [][]int32(nil)
    backtracking(m, c, res, cache)

    return res
}

func backtracking(m map[int32]int32, c *tagCondition, res *tagResult, cache [][]int32) {
    if len(res.cards) >= c.maxHand {
        return
    }

    if merge := sliceMerge(cache); len(merge) > len(res.cards) {
        res.cards = base.SliceCopy(merge)
        res.info = sliceCopy(cache)
    }

    for k, v := range m {
        if v <= 0 {
            continue
        }

        //杠
        if m[k] >= 4 && c.gang > 0 {
            m[k] -= 4
            cache = append(cache, []int32{k, k, k, k})
            backtracking(m, c, res, cache)
            m[k] += 4
            cache = cache[:len(cache)-1]
        }

        //刻
        if m[k] >= 3 && c.ke > 0 {
            m[k] -= 3
            cache = append(cache, []int32{k, k, k})
            backtracking(m, c, res, cache)
            m[k] += 3
            cache = cache[:len(cache)-1]
        }

        //对
        if m[k] >= 2 && c.dui > 0 {
            m[k] -= 2
            cache = append(cache, []int32{k, k})
            backtracking(m, c, res, cache)
            m[k] += 2
            cache = cache[:len(cache)-1]
        }

        if m[k] >= 1 && c.shun > 0 {
            //顺
            if m[k] >= 1 && m[k+1] >= 1 && m[k+2] >= 1 &&
                    toColor(k) == toColor(k+1) && toColor(k+1) == toColor(k+2) {
                m[k]--
                m[k+1]--
                m[k+2]--
                cache = append(cache, []int32{k, k + 1, k + 2})
                backtracking(m, c, res, cache)
                m[k]++
                m[k+1]++
                m[k+2]++
                cache = cache[:len(cache)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 && m[k+1] >= 1 &&
                    toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k+1) {
                m[k]--
                m[k-1]--
                m[k+1]--
                cache = append(cache, []int32{k, k - 1, k + 1})
                backtracking(m, c, res, cache)
                m[k]++
                m[k-1]++
                m[k+1]++
                cache = cache[:len(cache)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 && m[k-2] >= 1 &&
                    toColor(k) == toColor(k-1) && toColor(k-1) == toColor(k-2) {
                m[k]--
                m[k-1]--
                m[k-2]--
                cache = append(cache, []int32{k, k - 1, k - 2})
                backtracking(m, c, res, cache)
                m[k]++
                m[k-1]++
                m[k-2]++
                cache = cache[:len(cache)-1]
            }
        }

        if m[k] >= 1 && c.ca > 0 {
            //茬
            if m[k] >= 1 && m[k+1] >= 1 &&
                    toColor(k) == toColor(k+1) {
                m[k]--
                m[k+1]--
                cache = append(cache, []int32{k, k + 1})
                backtracking(m, c, res, cache)
                m[k]++
                m[k+1]++
                cache = cache[:len(cache)-1]
            }

            if m[k] >= 1 && m[k+2] >= 1 &&
                    toColor(k) == toColor(k+2) {
                m[k]--
                m[k+2]--
                cache = append(cache, []int32{k, k + 2})
                backtracking(m, c, res, cache)
                m[k]++
                m[k+2]++
                cache = cache[:len(cache)-1]
            }

            if m[k] >= 1 && m[k-1] >= 1 &&
                    toColor(k) == toColor(k-1) {
                m[k]--
                m[k-1]--
                cache = append(cache, []int32{k, k - 1})
                backtracking(m, c, res, cache)
                m[k]++
                m[k-1]++
                cache = cache[:len(cache)-1]
            }

            if m[k] >= 1 && m[k-2] >= 1 &&
                    toColor(k) == toColor(k-2) {
                m[k]--
                m[k-2]--
                cache = append(cache, []int32{k, k - 2})
                backtracking(m, c, res, cache)
                m[k]++
                m[k-2]++
                cache = cache[:len(cache)-1]
            }
        }

    }
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
