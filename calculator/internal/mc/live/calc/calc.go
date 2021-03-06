package calc

import (
	"fmt"

	"calculator/internal/base"

	"github.com/go-kratos/kratos/v2/log"
)

type TagLegalConfig struct {
	//一副牌
	Cards []int32

	//玩家配牌规则
	Rules []TagLegalRule
}
type TagLegalRule struct {
	Pid uint64 //玩家PlayerID
	Seq int32  //牌池优先级

	Gang      int32 //杠数量
	Ke        int32 //刻子数量
	Shun      int32 //顺子数量
	Dui       int32 //对子数量
	DoubleCha int32 //双茬数量
	SingleCha int32 //单茬数量
	LzCount   int32 //赖子数量
	BadCount  int32 //破坏张数
	RandCount int32 //随机张数

	ColorLimit int32     //花色限制
	PointLimit [][]int32 //点数限制
}

//文档需求：
//1.根据花色配置确定花色最大⻔数or花色，若配置的为限制花色最大⻔数为2⻔，则需在游戏内n个花色中排列组合花色如：3⻔花色的游戏，限制了花色限制为2，则需要试万筒、万条、筒条 3种组合
//2.根据点数配置确定点数，若配置的点数为“，”隔开的n组，则需要遍历n种情况
//3.根据牌型配置，随机各个牌型的组数（此刻已经确定完对转刻、茬转顺的组数了）
//4.按优先级依次尝试在3*n种花色、点数配置下发出（步骤3）确认的牌型并记录满足条件的张数x，直到成功发出注：x为张数，3杠2对2单张的张数=3*3+2*2+2*1=15张
//5.若（步骤4）最终未成功发出，则发出x最大的 花色、点数、牌型组合，多组x张数一致则在其中随机一个（注：研发同学要严格评估此流程对于服务器性能的压力，以防出现隐患）
//
//②开始生成当前“牌池优先级”最大玩家的「目标牌」，优先级相同随机一名开始（此时目标牌包含 牌型配置生成的牌+赖子数生成的牌）③从当前“牌池优先级”最大玩家目标「释放破环张数」至「总可发牌池」，剩余的牌称为「待修正牌」④下一牌池优先级最大玩家进行（1）～（2）步骤直到全部玩家、机器人处理完毕牌型点数限制花色限制全绿23468条全大中小123 ， 456 ， 7891清⺓九192小于5/大于512345，67891清一色1234567891无牌型不限制2总牌池玩家1目标牌总牌池玩家1目标牌缺张总牌池玩家1待修正牌玩家2待修正牌玩家3待修正牌玩家4待修正牌/46
func build(c *TagLegalConfig) (err error, use, left []int32, builds map[uint64][]int32) {

	m := map[int][]int32{}
	total := base.SliceCopy(c.Cards)

	for index, p := range c.Rules {

		tc := &tagCondition{
			gang:    p.Gang,
			shun:    p.Shun,
			ke:      p.Ke,
			dui:     p.Dui,
			ca:      p.SingleCha + p.DoubleCha,
			maxHand: 0,
		}

		cs := getColorCombinePointsCards(total, tc, p.ColorLimit, p.PointLimit)
		if len(cs) == 0 {
			m[index] = []int32{}
			continue
		}
		m[index] = append(m[index], cs...)

		total = base.SliceDel(total, cs...)
	}

	log.Infof("m:%+v", m)
	return
}

//递归找组合
func getColorCombinePointsCards(total []int32, c *tagCondition, colorLimit int32, pointsLimit [][]int32) (cs []int32) {
	colorList := getColorList(total)                 //花色列表
	combines := base.Dfs(colorList, int(colorLimit)) //排列组合花色：递归找花色组合
	for _, cc := range combines {
		for _, points := range pointsLimit {

			cards := getColorPointsCards(total, cc, points)
			if len(cards) <= 0 {
				continue
			}
			res := permute(cards, c)
			if len(res.info.cards) <= 0 {
				continue
			}
			return base.SliceCopy(res.info.cards)
		}
	}
	return nil
}

func getColorList(cards []int32) (colors []int32) {
	m := map[int32]bool{}
	for _, c := range cards {
		color := toColor(c)
		if m[color] {
			continue
		}
		m[color] = true
		colors = append(colors, color)
	}
	return
}

func getColorPointsCards(cards []int32, colors []int32, points []int32) (colorPointCards []int32) {
	for _, c := range cards {
		color := toColor(c)
		point := toPoint(c)
		if !base.SliceContain(colors, color) {
			continue
		}
		if !base.SliceContain(points, point) {
			continue
		}
		colorPointCards = append(colorPointCards, c)
	}
	return colorPointCards
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
