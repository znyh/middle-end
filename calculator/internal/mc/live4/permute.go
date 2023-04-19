package live4

import (
    "calculator/internal/base"
)

const _MaxCountLimit = 18

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
    gang int32
    shun int32
    ke   int32
    dui  int32
    ca   int32

    maxHand int
}

type tagNode struct {
    value int32 // 值
    count int32 // 个数
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
    cache := &tagComInfo{}
    res := &tagResult{info: &tagComInfo{}}
    _data, _cards := hCardToDCard2(cards)
    c.maxHand = base.MInInt(int(c.gang*4+c.ke*3+c.shun*3+c.dui*2+c.ca*2), len(cards), _MaxCountLimit)
    split(_data, _cards, c, res, cache, 0)
    return res
}

func hCardToDCard2(cList []int32) (_data [5][10]tagNode, _cards []int32) {
    _data = [5][10]tagNode{}
    _cards = []int32{}
    for i := 0; i < 5; i++ {
        for j := 0; j < 10; j++ {
            _data[i][j] = tagNode{}
        }
    }
    m := map[int32]bool{}
    for i := 0; i < len(cList); i++ {
        ty := toColor(cList[i]) //cList[i] / 0x10
        tv := toPoint(cList[i]) //cList[i] % 0x10
        _data[ty][tv].value = cList[i]
        _data[ty][tv].count++
        if _, ok := m[cList[i]]; !ok {
            m[cList[i]] = true
            _cards = append(_cards, cList[i])
        }
    }

    _cards = base.SliceShuffle(_cards)

    return _data, _cards
}

func split(_data [5][10]tagNode, cList []int32, c *tagCondition, res *tagResult, cache *tagComInfo, index int32) {
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

    for k := index; k < int32(len(cList)); k++ {

        // 计算结果

        i := toPoint(cList[k]) //cList[k] % 0x10
        j := toColor(cList[k]) //cList[k] / 0x10

        //gang
        if _data[j][i].count >= 4 && c.gang > 0 && int(c.gang) > len(cache.gang) {
            _data[j][i].count -= 4
            cache.gang = append(cache.gang, []int32{_data[j][i].value, _data[j][i].value, _data[j][i].value, _data[j][i].value})
            if _data[j][i].count > 0 {
                split(_data, cList, c, res, cache, k)
            } else {
                split(_data, cList, c, res, cache, k+1)
            }
            _data[j][i].count += 4
            cache.gang = cache.gang[:len(cache.gang)-1]
        }

        //ke
        if _data[j][i].count >= 3 && c.ke > 0 && int(c.ke) > len(cache.ke) {
            _data[j][i].count -= 3
            cache.ke = append(cache.ke, []int32{_data[j][i].value, _data[j][i].value, _data[j][i].value})
            if _data[j][i].count > 0 {
                split(_data, cList, c, res, cache, k)
            } else {
                split(_data, cList, c, res, cache, k+1)
            }
            _data[j][i].count += 3
            cache.ke = cache.ke[:len(cache.ke)-1]
        }

        //dui
        if _data[j][i].count >= 2 && c.dui > 0 && int(c.dui) > len(cache.dui) {
            _data[j][i].count -= 2
            cache.dui = append(cache.dui, []int32{_data[j][i].value, _data[j][i].value})
            if _data[j][i].count > 0 {
                split(_data, cList, c, res, cache, k)
            } else {
                split(_data, cList, c, res, cache, k+1)
            }
            _data[j][i].count += 2
            cache.dui = cache.dui[:len(cache.dui)-1]
        }

        //shun
        if j > 0 && j < 4 && c.shun > 0 && int(c.shun) > len(cache.shun) {
            if (i >= 0 && i+2 <= 9) && _data[j][i].count > 0 && _data[j][i+1].count > 0 && _data[j][i+2].count > 0 {
                _data[j][i].count -= 1
                _data[j][i+1].count -= 1
                _data[j][i+2].count -= 1
                cache.shun = append(cache.shun, []int32{_data[j][i].value, _data[j][i+1].value, _data[j][i+2].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i+1].count += 1
                _data[j][i+2].count += 1
                cache.shun = cache.shun[:len(cache.shun)-1]
            }

            if (i-1 >= 0 && i+1 <= 9) && _data[j][i].count > 0 && _data[j][i-1].count > 0 && _data[j][i+1].count > 0 {
                _data[j][i].count -= 1
                _data[j][i-1].count -= 1
                _data[j][i+1].count -= 1
                cache.shun = append(cache.shun, []int32{_data[j][i].value, _data[j][i-1].value, _data[j][i+1].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i-1].count += 1
                _data[j][i+1].count += 1
                cache.shun = cache.shun[:len(cache.shun)-1]
            }

            if (i-2 >= 0 && i <= 9) && _data[j][i].count > 0 && _data[j][i-1].count > 0 && _data[j][i-2].count > 0 {
                _data[j][i].count -= 1
                _data[j][i-1].count -= 1
                _data[j][i-2].count -= 1
                cache.shun = append(cache.shun, []int32{_data[j][i].value, _data[j][i-1].value, _data[j][i-2].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i-1].count += 1
                _data[j][i-2].count += 1
                cache.shun = cache.shun[:len(cache.shun)-1]
            }
        }

        //ca
        if j > 0 && j < 4 && c.ca > 0 && int(c.ca) > len(cache.ca) {
            if (i >= 0 && i+1 <= 9) && _data[j][i].count > 0 && _data[j][i+1].count > 0 {
                _data[j][i].count -= 1
                _data[j][i+1].count -= 1
                cache.ca = append(cache.ca, []int32{_data[j][i].value, _data[j][i+1].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i+1].count += 1
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if (i >= 0 && i+2 <= 9) && _data[j][i].count > 0 && _data[j][i+2].count > 0 {
                _data[j][i].count -= 1
                _data[j][i+2].count -= 1
                cache.ca = append(cache.ca, []int32{_data[j][i].value, _data[j][i+2].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i+2].count += 1
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if (i-1 >= 0 && i <= 9) && _data[j][i].count > 0 && _data[j][i-1].count > 0 {
                _data[j][i].count -= 1
                _data[j][i-1].count -= 1
                cache.ca = append(cache.ca, []int32{_data[j][i].value, _data[j][i-1].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i-1].count += 1
                cache.ca = cache.ca[:len(cache.ca)-1]
            }

            if (i-2 >= 0 && i <= 9) && _data[j][i].count > 0 && _data[j][i-2].count > 0 {
                _data[j][i].count -= 1
                _data[j][i-2].count -= 1
                cache.ca = append(cache.ca, []int32{_data[j][i].value, _data[j][i-2].value})
                if _data[j][i].count > 0 {
                    split(_data, cList, c, res, cache, k)
                } else {
                    split(_data, cList, c, res, cache, k+1)
                }
                _data[j][i].count += 1
                _data[j][i-2].count += 1
                cache.ca = cache.ca[:len(cache.ca)-1]
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
