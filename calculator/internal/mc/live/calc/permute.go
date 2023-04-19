/*
	递归找出需要的组合
*/

package calc

import (
	"time"

	"calculator/internal/base"

	"github.com/go-kratos/kratos/v2/log"
)

const MAXCOUNT = 13

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

	start := time.Now()

	m := map[int32]int32{}
	for _, v := range cards {
		m[v]++
	}

	res := &tagResult{info: &tagComInfo{}}
	cache := &tagComInfo{}
	c.maxHand = int(c.gang*4 + c.ke*3 + c.shun*3 + c.dui*2 + c.ca*2)
	if c.maxHand > MAXCOUNT {
		c.maxHand = MAXCOUNT
	}
	backtracking(m, c, res, cache)

	if _OpenDebugLog {
		use := time.Since(start).Milliseconds()
		if use > 100 {
			log.Errorf("use:%+v (time over 100ms),cards:%+v c:%+v", use, descCards(cards), c)
		}
		log.Infof("use:%+v/ms c:%+v %+v", use, c, res.info)
		if len(res.info.cards) > 0 && !base.SliceContainAll(cards, res.info.cards...) {
			log.Errorf("=> error. !base.SliceContainAll(cards,special...) %+v %+v", descCards(cards), descCards(res.info.cards))
			res.info = &tagComInfo{}

			log.Infof("\n\n\n")
		}
	}

	return res
}

func backtracking(m map[int32]int32, c *tagCondition, res *tagResult, cache *tagComInfo) {

	// 找到一组就退出
	if len(res.info.cards) >= c.maxHand {
		return
	}

	//成功找到一组
	cnt := len(cache.gang)*4 + (len(cache.ke)+len(cache.shun))*3 + (len(cache.dui)+len(cache.ca))*2
	if cnt > MAXCOUNT {
		return
	}
	if cnt > len(res.info.cards) && cnt <= MAXCOUNT {
		tmp := cache.clone()
		res.info = tmp
	}

	for k, v := range m {
		if v <= 0 {
			continue
		}

		if c.gang > 0 && int(c.gang) > len(cache.gang) {
			//杠
			if m[k] >= 4 {
				m[k] -= 4
				cache.gang = append(cache.gang, []int32{k, k, k, k})
				backtracking(m, c, res, cache)
				m[k] += 4
				cache.gang = cache.gang[:len(cache.gang)-1]
			}
		}

		if c.ke > 0 && int(c.ke) > len(cache.ke) {
			//刻
			if m[k] >= 3 {
				m[k] -= 3
				cache.ke = append(cache.ke, []int32{k, k, k})
				backtracking(m, c, res, cache)
				m[k] += 3
				cache.ke = cache.ke[:len(cache.ke)-1]
			}
		}

		if c.shun > 0 && int(c.shun) > len(cache.shun) {
			//顺
			if m[k] >= 1 && m[k+1] >= 1 && m[k+2] >= 1 &&
				toColor(k) == toColor(k+1) && toColor(k+1) == toColor(k+2) {
				m[k]--
				m[k+1]--
				m[k+2]--
				cache.shun = append(cache.shun, []int32{k, k + 1, k + 2})
				backtracking(m, c, res, cache)
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
				backtracking(m, c, res, cache)
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
				backtracking(m, c, res, cache)
				m[k]++
				m[k-1]++
				m[k-2]++
				cache.shun = cache.shun[:len(cache.shun)-1]
			}
		}

		if c.dui > 0 && int(c.dui) > len(cache.dui) {
			//对
			if m[k] >= 2 {
				m[k] -= 2
				cache.dui = append(cache.dui, []int32{k, k})
				backtracking(m, c, res, cache)
				m[k] += 2
				cache.dui = cache.dui[:len(cache.dui)-1]
			}
		}

		if c.ca > 0 && int(c.ca) > len(cache.ca) {
			//茬
			if m[k] >= 1 && m[k+1] >= 1 &&
				toColor(k) == toColor(k+1) {
				m[k]--
				m[k+1]--
				cache.ca = append(cache.ca, []int32{k, k + 1})
				backtracking(m, c, res, cache)
				m[k]++
				m[k+1]++
				cache.ca = cache.ca[:len(cache.ca)-1]
			}

			if m[k] >= 1 && m[k+2] >= 1 &&
				toColor(k) == toColor(k+2) {
				m[k]--
				m[k+2]--
				cache.ca = append(cache.ca, []int32{k, k + 2})
				backtracking(m, c, res, cache)
				m[k]++
				m[k+2]++
				cache.ca = cache.ca[:len(cache.ca)-1]
			}

			if m[k] >= 1 && m[k-1] >= 1 &&
				toColor(k) == toColor(k-1) {
				m[k]--
				m[k-1]--
				cache.ca = append(cache.ca, []int32{k, k - 1})
				backtracking(m, c, res, cache)
				m[k]++
				m[k-1]++
				cache.ca = cache.ca[:len(cache.ca)-1]
			}

			if m[k] >= 1 && m[k-2] >= 1 &&
				toColor(k) == toColor(k-2) {
				m[k]--
				m[k-2]--
				cache.ca = append(cache.ca, []int32{k, k - 2})
				backtracking(m, c, res, cache)
				m[k]++
				m[k-2]++
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
