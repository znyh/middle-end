package calc

import (
	"fmt"
	"sort"

	"calculator/internal/base"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	MAXCOUNT = 13
)

type (
	tagSet struct {
		//输入
		close     bool  //是否关闭
		kind      int32 //1,2类型配牌
		continues bool  //是否连号

		//输出
		items []tagItem //
	}

	tagItem struct {
		index int          //索引
		c     tagCondition //条件
		l     tagList      //计算列表
	}

	//配置参数
	tagCondition struct {
		master int32 //主花色张数
		slave  int32 //次花色张数
		shun   int32 //顺子个数
		ke     int32 //刻子个数
		dui    int32 //对子个数
		ca     int32 //茬个数

		mustConsecutive bool //要求在区间上连号 区间:[1，3],[4,6],[7,9] （m个对子,n个对子 且m+n>=3）
	}

	//输出结果列表
	tagList struct {
		hands  []int32 //手牌
		master []int32 //主花色牌堆
		slave  []int32 //次花色牌堆
		left   []int32 //手牌剩余牌堆
	}
)

func newSet(cards []int32, close bool, kind int32, continues bool, items []tagItem) (s *tagSet) {
	s = &tagSet{
		close:     close,
		kind:      kind,
		continues: continues,
		items:     items,
	}

	//一副牌
	total := base.SliceShuffle(base.SliceCopy(cards))

	switch kind {
	case 1:
		//1:【根据花色在目前可发张数的权重来随机配牌】
		s.rand1(total)
	case 2:
		//2:【根据牌池中是否存在至少1门牌可发张数&&要求牌型配牌】
		s.rand2(total)
	}

	sort.Slice(s.items, func(i, j int) bool {
		return s.items[i].index < s.items[j].index
	})

	return s
}

//1:【根据花色在目前可发张数的权重来随机配牌】
//
//非剧本局的前提下，当本局桌内玩家or机器人配牌出现了花色配置时，按以下优先级处理发牌
//1）轮流 按需求张数最多处理配置了最多“主花色”张数玩家的“主花色”牌
//2）轮流 按需求张数最多处理配置了最多“次花色”张数玩家的“次花色”牌
//3）再处理未配置“花色”玩家的手牌
//
//当有多个玩家配置了主花色时，按以下逻辑处理：
//
//轮流处理主花色优先级为：主花色要求张数多到少，存在多名玩家同时最多时在同样多的人中随机一个优先处理
//每名玩家处理花色流程如下：
//判断是否牌池中是否存在至少1门牌可发张数>该名玩家本局要求的花色张数
//是：根据可发牌中满足张数要求的花色在目前可发张数的权重来随机花色（见下面例子）
//否：该名玩家此项配牌不生效，走纯随机，该名玩家已确定的牌照常发（比如该玩家在确定副花色的时候发现每门牌都打不到张数要求，但这个玩家的主花色已经发过了，这时该玩家的副花色就随机发）
func (s *tagSet) rand1(total []int32) {

	//处理主花色 master
	{
		//主花色排序
		sort.Slice(s.items, func(i, j int) bool {
			return s.items[i].c.master > s.items[j].c.master
		})

		//拼凑主花色
		for k, v := range s.items {
			//按权重获取主花色
			cm := toColorMap(total)
			maxColor, maxCnt := colorWeightRand(cm)
			if maxColor == -1 || maxCnt == -1 {
				continue
			}
			if int(v.c.master) > maxCnt {
				continue
			}

			cs := cm[maxColor][:int(v.c.master)]
			s.items[k].l.hands = append(s.items[k].l.hands, cs...)
			s.items[k].l.master = append(s.items[k].l.master, cs...)
			total = base.SliceDel(total, cs...)
		}
	}

	//处理次花色 slave
	{
		//次花色排序
		sort.Slice(s.items, func(i, j int) bool {
			return s.items[i].c.slave > s.items[j].c.slave
		})

		//拼凑次花色
		for k, v := range s.items {
			if len(v.l.hands) >= MAXCOUNT {
				continue
			}

			//次花色与主花色花色不一样
			masterColor := int32(-1)
			if len(s.items[k].l.master) > 0 {
				masterColor = toColor(s.items[k].l.master[0])
			}

			//去掉主花色的牌堆
			cm := toColorMap(total)
			cm[masterColor] = []int32{}

			//按权重获取次花色
			maxColor, maxCnt := colorWeightRand(cm)
			if maxColor == -1 || maxCnt == -1 {
				continue
			}
			if int(v.c.slave) > maxCnt {
				continue
			}

			need := base.MInInt(int(v.c.slave), MAXCOUNT-len(v.l.hands))
			cs := cm[maxColor][:need]
			s.items[k].l.hands = append(s.items[k].l.hands, cs...)
			s.items[k].l.slave = append(s.items[k].l.slave, cs...)
			total = base.SliceDel(total, cs...)
		}
	}

	//left
	{
		for k, v := range s.items {
			if len(v.l.hands) >= MAXCOUNT {
				continue
			}
			need := MAXCOUNT - len(v.l.hands)
			if need > len(total) {
				need = len(total)
			}

			cs := total[:need]
			s.items[k].l.hands = append(s.items[k].l.hands, cs...)
			s.items[k].l.left = append(s.items[k].l.left, cs...)
			total = base.SliceDel(total, cs...)
		}
	}

}

//根据花色权重随机 maxColor，maxCnt:随机到权重较大的花色和该花色对应的剩余张数
func colorWeightRand(cm map[int32][]int32) (maxColor int32, maxCnt int) {
	if len(cm) <= 0 {
		return -1, -1
	}

	type tagWeight struct {
		color  int32
		weight int
	}

	sum := 0
	weights := []tagWeight(nil)
	for color, cards := range cm {
		cnt := len(cards)
		sum += cnt
		weights = append(weights, tagWeight{
			color:  color,
			weight: cnt,
		})
	}

	r := base.RandRange(0, sum)

	start, end := 0, 0
	for _, v := range weights {
		end = start + v.weight - 1
		if r >= start && r <= end {
			return v.color, v.weight

		}
		start += v.weight
	}

	return -1, -1
}

//2:【根据牌池中是否存在至少1门牌可发张数&&要求牌型配牌】 特殊开关：是否连号
//
//非剧本局的前提下，当本局桌内玩家or机器人配牌出现了花色配置时，按以下优先级处理发牌
//1）轮流 按需求张数最多处理配置了最多“主花色”张数玩家的“主花色”牌
//2）轮流 按需求张数最多处理配置了最多“次花色”张数玩家的“次花色”牌
//3）再处理未配置“花色”玩家的手牌
//
//轮流处理主花色优先级为：主花色要求张数多到少，存在多名玩家同时最多时在同样多的人中随机一个优先处理
//每名玩家处理花色流程如下：
//判断是否牌池中是否存在至少1门牌可发张数&&要求牌型（要求牌型就是看对子、顺子、刻子、茬那些够不够）足够 该名玩家本局要求的花色张数与牌型
//是：选择可发牌中满足张数要求的花色张数在目前可发张数的张数最多的一门牌来确认花色，多门最大时在花色张数最多中随机一门
//否：该名玩家此项配牌不生效，走纯随机，该名玩家已确定的牌照常发（比如该玩家在确定副花色的时候发现每门牌都打不到张数要求，但这个玩家的主花色已经发过了，这时该玩家的副花色就随机发）
func (s *tagSet) rand2(total []int32) {

	//处理主花色 master
	{
		//主花色排序
		sort.Slice(s.items, func(i, j int) bool {
			return s.items[i].c.master > s.items[j].c.master
		})

		//拼凑主花色 根据要求的牌型
		for k, v := range s.items {
			//获取主花色（剩余花色张数最多的为主花色）
			cm := toColorMap(total)
			maxColor, maxCnt := colorMaxRand(cm)
			if maxColor == -1 || maxCnt == -1 {
				continue
			}
			if int(v.c.master) > maxCnt {
				continue
			}

			if !s.continues {
				//非连号的情况
				colorCards := base.SliceCopy(cm[maxColor])
				cs := getSpecialCards(colorCards, v.c)
				if len(cs) > 0 {
					s.items[k].l.hands = append(s.items[k].l.hands, cs...)
					s.items[k].l.master = append(s.items[k].l.master, cs...)
					total = base.SliceDel(total, cs...)
				}

			} else {
				//连号情况
				colorCards := base.SliceCopy(cm[maxColor])
				cs := getSpecialContinuesCards(colorCards, v.c)
				if len(cs) == 0 {
					//主花张数和牌型满足条件，连号不满足条件，得给我发满足主花张数和牌型条件的牌
					cs = getSpecialCards(colorCards, v.c)
				}
				if len(cs) > 0 {
					s.items[k].l.hands = append(s.items[k].l.hands, cs...)
					s.items[k].l.master = append(s.items[k].l.master, cs...)
					total = base.SliceDel(total, cs...)
				}
			}

		}
	}

	//处理次花色 slave
	{
		//次花色排序
		sort.Slice(s.items, func(i, j int) bool {
			return s.items[i].c.slave > s.items[j].c.slave
		})

		//拼凑次花色
		for k, v := range s.items {
			if len(v.l.hands) >= MAXCOUNT {
				continue
			}

			//次花色与主花色花色不一样
			masterColor := int32(-1)
			if len(s.items[k].l.master) > 0 {
				masterColor = toColor(s.items[k].l.master[0])
			}
			cm := toColorMap(total)
			cm[masterColor] = []int32{}
			maxColor, maxCnt := colorMaxRand(cm)
			if maxColor == -1 || maxCnt == -1 {
				continue
			}
			if int(v.c.slave) > maxCnt {
				continue
			}

			need := base.MInInt(int(v.c.slave), MAXCOUNT-len(v.l.hands))
			cs := cm[maxColor][:need]
			s.items[k].l.hands = append(s.items[k].l.hands, cs...)
			s.items[k].l.slave = append(s.items[k].l.slave, cs...)
			total = base.SliceDel(total, cs...)
		}
	}

	//left
	{
		for k, v := range s.items {
			if len(v.l.hands) >= MAXCOUNT {
				continue
			}
			need := MAXCOUNT - len(v.l.hands)
			if need > len(total) {
				need = len(total)
			}

			cs := total[:need]
			s.items[k].l.hands = append(s.items[k].l.hands, cs...)
			s.items[k].l.left = append(s.items[k].l.left, cs...)
			total = base.SliceDel(total, cs...)
		}
	}

}

func colorMaxRand(cm map[int32][]int32) (maxColor int32, maxCnt int) {
	if len(cm) <= 0 {
		return -1, -1
	}
	for color, cards := range cm {
		if cnt := len(cards); cnt > maxCnt {
			maxCnt = cnt
			maxColor = color
		}
	}
	return
}

//判断是否牌池中是否存在至少1门牌可发张数&&要求牌型（要求牌型就是看对子、顺子、刻子、茬那些够不够）足够 该名玩家本局要求的花色张数与牌型
func getSpecialCards(cards []int32, c tagCondition) (cache []int32) {
	//张数够?
	cnt := int32(len(cards))
	if cnt <= 0 {
		return
	}
	if cnt < c.shun*3+c.ke*3+c.dui*2+c.ca*2 {
		return
	}

	//计算牌型
	ret := permute(cards, &c)
	if len(ret.infos) <= 0 || len(ret.infos[0].cards) <= 0 {
		return
	}

	special := base.SliceCopy(ret.infos[0].cards)
	if len(special) >= int(c.master) {
		return special
	}
	//补充剩余的
	left := base.SliceDel(base.SliceCopy(cards), special...)
	special = append(special, left...)
	if len(special) < int(c.master) {
		return special
	}
	return special[:int(c.master)]
}

//在配牌配置满足以下条件时增加【是否连号】开关（副花色无此配置）：
//
//配置了主花色
//配置的刻子+对子组数>=3
//（注：根据之前规则对刻一定在主花色里）
//
//主花色内 对子、刻子中的必有任意3组为1～3、4～6、7～9三个类型中的连号
//主花色中-未被选中连号的 对子、刻子不可发出与连号对、刻相同的牌
func getSpecialContinuesCards(cards []int32, c tagCondition) (cache []int32) {

	//张数够?
	if len(cards) <= 0 {
		return
	}
	//配置了主花色
	if c.master <= 0 {
		return
	}
	//配置的刻子+对子组数>=3
	if c.ke+c.dui < 3 {
		return
	}

	//分区间统计 [1,3] [4,6] [7,9]
	m := map[int32][]int32{}
	start := toColor(cards[0]) * 10
	for _, v := range cards {
		if v >= start+1 && v <= start+3 {
			m[0] = append(m[0], v)
		} else if v >= start+4 && v <= start+6 {
			m[1] = append(m[1], v)
		} else if v >= start+7 && v <= start+9 {
			m[2] = append(m[2], v)
		}
	}

	//找连号的[1,3] [4,6] [7,9]
	for _, area := range m {

		//至少需要3*2=6张 要求刻子+对子组数>=3
		if cnt := int32(len(area)); cnt < 3*2 {
			continue
		}

		//判断是否有连号区间
		if ok := checkConsecutive(area); !ok {
			continue
		}

		//m个对子 n个对子 m+n == 3
		minKe := base.MInInt32(3, c.ke)
		minDui := base.MInInt32(3, c.dui)
		for ke := int32(0); ke <= minKe; ke++ {
			for dui := int32(0); dui <= minDui; dui++ {
				if ke+dui < 3 {
					continue
				}

				part1, part2 := []int32(nil), []int32(nil)

				//找连号的m个对子 n个对子
				c1 := tagCondition{
					ke:              ke,
					dui:             dui,
					mustConsecutive: true, //要求在区间[1，3],[4,6],[7,9]连号
				}
				if c1.ke+c1.dui+c1.shun+c1.ca > 0 {
					ret := permute(area, &c1)
					if len(ret.infos) <= 0 || len(ret.infos[0].cards) <= 0 {
						continue
					}
					//是否有连号区间
					ok := checkConsecutive(ret.infos[0].cards)
					if !ok {
						log.Errorf("checkAreaContinue error...(should be continued) part:%+v area:%+v", descCards(ret.infos[0].cards), descCards(area))
						continue
					}
					part1 = ret.infos[0].cards
				}

				//补充剩余的牌型 c中去掉连号的m个对子 n个对子
				c2 := tagCondition{
					master: c.master,
					slave:  c.slave,
					shun:   c.shun,
					ke:     c.ke - ke,
					dui:    c.dui - dui,
					ca:     c.ca,
				}
				if c2.ke+c2.dui+c2.shun+c2.ca > 0 {
					ret := permute(base.SliceDel(base.SliceCopy(cards), base.SliceCopy(area)...), &c2)
					if len(ret.infos) <= 0 || len(ret.infos[0].cards) <= 0 {
						continue
					}
					part2 = ret.infos[0].cards
				}

				//按主花色张数补充剩余的 (与选中的区间k不一致的牌做填充)
				special := base.SliceCopy(part1)
				special = append(special, part2...)
				if len(special) >= int(c.master) {
					return special
				}
				left := base.SliceDel(base.SliceCopy(cards), base.SliceCopy(special)...)
				needCnt := base.MInInt(int(c.master)-len(special), len(left))
				special = append(special, left[:needCnt]...)
				return special
			}
		}

	}

	return nil
}

//区间是否连号的 且张数至少为2
func checkConsecutive(area []int32) (ok bool) {
	arr := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, v := range area {
		point := toPoint(v)
		arr[point]++
	}
	for i := int32(1); i <= 9; i += 3 {
		if arr[i] >= 2 && arr[i+1] >= 2 && arr[i+2] >= 2 {
			return true
		}
	}
	return false
}

func toColorMap(cards []int32) (cm map[int32][]int32 /*color:cards*/) {
	cm = map[int32][]int32{}
	for _, c := range cards {
		color := toColor(c)
		cm[color] = append(cm[color], c)
	}
	return
}
func descCards(cards []int32) (str string) {
	maxsize := 34
	m := toIndexes(cards)

	str += "\n"
	for i := 0; i < maxsize; i++ {
		str += fmt.Sprintf("%+v ", m[i])
		if (i+1)%9 == 0 {
			str += "\n"
		}
	}
	return
}
func toIndexes(cards []int32) (arr []int) {
	maxsize := 34
	m := make([]int, maxsize)
	for _, c := range cards {
		index := toIndex(c)
		if index >= maxsize {
			continue
		}
		m[index]++
	}
	return m
}
func toIndex(card int32) int {
	color := toColor(card)
	point := toPoint(card)
	index := (color-1)*9 + point - 1
	return int(index)
}
func toColor(c int32) (color int32) {
	color = c / 10
	return
}
func toPoint(c int32) (point int32) {
	point = c % 10
	return
}
