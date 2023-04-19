/*
	递归找出需要的组合
*/

package calc

import (
	"time"

	"calculator/internal/base"

	"github.com/go-kratos/kratos/v2/log"
)

var _OpenDebugLog = false

type tagResult struct {
	infos []*tagComInfo
}
type tagComInfo struct {
	//序列
	cards []int32
	//组合
	shun [][]int32
	ke   [][]int32
	dui  [][]int32
	ca   [][]int32
}

func (c *tagComInfo) clone() (t *tagComInfo) {
	t = &tagComInfo{}
	t.ke = sliceCopy(c.ke)
	t.shun = sliceCopy(c.shun)
	t.dui = sliceCopy(c.dui)
	t.ca = sliceCopy(c.ca)
	t.cards = sliceMerge(t.shun, t.ke, t.dui, t.ca)
	return t
}

func permute(cards []int32, c *tagCondition) *tagResult {

	start := time.Now()

	m := map[int32]int32{}
	for _, v := range cards {
		m[v]++
	}

	res := &tagResult{}
	cache := &tagComInfo{}
	backtracking(m, c, res, cache)

	if _OpenDebugLog {
		use := time.Since(start).Milliseconds()
		if use > 100 {
			log.Errorf("use:%+v (time over 100ms),cards:%+v c:%+v", use, descCards(cards), c)
		}
		for k, v := range res.infos {
			log.Infof("use:%+v/ms c:%+v %+v", use, c, v)
			if len(v.cards) > 0 && !base.SliceContainAll(cards, v.cards...) {
				log.Errorf("=> error. !base.SliceContainAll(cards,special...) %+v %+v", descCards(cards), descCards(v.cards))
				res.infos[k] = &tagComInfo{}
			}
		}
	}

	return res
}

func backtracking(m map[int32]int32, c *tagCondition, res *tagResult, cache *tagComInfo) {
	//找到一组就退出
	if len(res.infos) >= 1 {
		return
	}

	// 成功找到一组
	if len(cache.shun) == int(c.shun) &&
		len(cache.ke) == int(c.ke) &&
		len(cache.dui) == int(c.dui) &&
		len(cache.ca) == int(c.ca) {

		tmp := cache.clone()

		//不要求连号的
		if !c.mustConsecutive {
			res.infos = append(res.infos, tmp)
			return
		}

		//要求是连号的
		if c.mustConsecutive && checkConsecutive(tmp.cards) {
			res.infos = append(res.infos, tmp)
			return
		}

	}

	for k, v := range m {
		if v <= 0 {
			continue
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
