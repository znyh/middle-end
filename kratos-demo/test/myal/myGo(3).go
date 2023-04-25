package main

import (
	"fmt"
	"math/rand"
	"time"
)

//optional int32 Value = 1;  // 1-13 14:okey牌
//optional int32 Type = 2;   // 0-4  红黄蓝黑okey
var CardList = []int32{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, // 1-13 红
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, // 1-13 黄
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x18, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, // 1-13 蓝
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, // 1-13 黑
	//0x4d, 0x4d,
}

func show(cList []int32) {
	fmt.Printf("HandCard:")
	for i := 0; i < len(cList); i++ {
		fmt.Printf("0x%02x ", cList[i])
	}
	fmt.Printf("\n")
}

type Node struct {
	Value int32 // 值
	Count int32 // 个数
}

type ComInfo struct {
	Value []int32 // 值
	Type  int32   // 0:群组 1:顺子
	LaiZi []int32 // 癞子值
}

type ResultInfo struct {
	Value int32     //最大值
	Info  []ComInfo // 组合
	Card  []int32   // 剩余手牌
}

// 手牌转 红黄蓝黑列表
func hCardToDCard(cList []int32) ([4][14]*Node, []int32) {
	_Data := [4][14]*Node{}
	_Cards := []int32{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 14; j++ {
			_Data[i][j] = &Node{}
		}
	}
	for i := 0; i < len(cList); i++ {
		ty := cList[i] / 0x10
		tv := cList[i] % 0x10
		_Data[ty][tv].Value = cList[i]
		_Data[ty][tv].Count++
		flag := false
		for _, v := range _Cards {
			if v == cList[i] {
				flag = true
				break
			}
		}
		if flag == false {
			_Cards = append(_Cards, cList[i])
		}
	}
	return _Data, _Cards
}

func showData(_data [4][14]*Node) {
	for i := 0; i < 4; i++ {
		fmt.Printf("[%d]:", i)
		for _, v := range _data[i] {
			fmt.Printf("0x%02x_%d ", v.Value, v.Count)
		}
		fmt.Printf("\n")
	}
}

// 对子
func DuiZi(cList []int32) {
	// // 打印手牌
	// show(cList)
	// // 手牌转换为 红黄蓝黑 列表
	// _data := hCardToDCard(cList)
	// showData(_data)
}

func showResult(_data [4][14]*Node, _ret []*ComInfo) {
	// 打印
	fmt.Printf("组合(%d):\n", len(_ret))
	for k, v := range _ret {
		fmt.Printf("組合(%d) %d ", k, v.Type)
		for i := 0; i < len(v.Value); i++ {
			fmt.Printf("0x%02x ", v.Value[i])
		}
		fmt.Println("")
	}
	fmt.Printf("剩余手牌:")
	for i := 0; i < 14; i++ {
		for j := 0; j < 4; j++ {
			if _data[j][i].Count > 0 {
				for l := 0; l < int(_data[j][i].Count); l++ {
					fmt.Printf("0x%02x ", _data[j][i].Value)
				}
			}
		}
	}
	fmt.Println("")
}

func finish(_data [4][14]*Node, rInfo *ResultInfo, _ret []*ComInfo) {

	bCard := []int32{}
	sum := int32(0)
	//补组合 因为当前最长组合为 3
	for _, v := range _ret {
		// 补群组
		if v.Type == 0 {
			ty := v.Value[2] / 0x10
			tv := v.Value[2] % 0x10
			for i := ty + 1; i < 4; i++ {
				if _data[i][tv].Count > 0 {
					v.Value = append(v.Value, _data[i][tv].Value)
					bCard = append(bCard, _data[i][tv].Value)
					_data[i][tv].Count -= 1
				}
			}
			// 统计值
			sum += tv * int32(len(v.Value))
		} else { // 补顺子
			ty := v.Value[2] / 0x10
			tv := v.Value[2] % 0x10
			for i := tv + 1; i < 14; i++ {
				if _data[ty][i].Count > 0 {
					v.Value = append(v.Value, _data[ty][i].Value)
					bCard = append(bCard, _data[ty][i].Value)
					_data[ty][i].Count -= 1
				} else {
					break
				}
			}
			// 统计值
			for i := 0; i < len(v.Value); i++ {
				sum += v.Value[i] % 0x10
			}
		}
	}

	// 打印组合
	//showResult(_data, _ret)

	// 保存最大值
	if sum > rInfo.Value {
		rInfo.Value = sum
		rInfo.Info = []ComInfo{}
		for _, v := range _ret {
			tNode := ComInfo{}
			tNode.Type = v.Type
			for _, v1 := range v.Value {
				tNode.Value = append(tNode.Value, v1)
			}
			for _, v1 := range v.LaiZi {
				tNode.LaiZi = append(tNode.LaiZi, v1)
			}
			rInfo.Info = append(rInfo.Info, tNode)
		}
		rInfo.Card = []int32{}
		for i := 1; i < 14; i++ {
			for j := 0; j < 4; j++ {
				if _data[j][i].Count > 0 {
					for l := 0; l < int(_data[j][i].Count); l++ {
						rInfo.Card = append(rInfo.Card, _data[j][i].Value)
					}
				}
			}
		}
	}

	// 还原
	for _, v := range bCard {
		ty := v / 0x10
		tv := v % 0x10
		_data[ty][tv].Count++
	}

	// 还原 ret
	for _, v := range _ret {
		v.Value = v.Value[:3]
	}
}

var AllCount int32

func diGuiZuPai(_data [4][14]*Node, cList []int32, rInfo *ResultInfo, _ret []*ComInfo, index int32, laiZi int32) {
	if index >= int32(len(cList)) {
		if len(_ret) == 0 {
			return
		}
		AllCount++
		finish(_data, rInfo, _ret)
		return
	}

	for k := index; k < int32(len(cList)); k++ {
		j := cList[k] / 0x10
		i := cList[k] % 0x10
		// 一张牌 2种情况 要么做顺子 要么做群组
		// 找顺子 3
		if i <= 11 && _data[j][i].Count > 0 && _data[j][i+1].Count > 0 && _data[j][i+2].Count > 0 {
			_data[j][i].Count -= 1
			_data[j][i+1].Count -= 1
			_data[j][i+2].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j][i+1].Value, _data[j][i+2].Value},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j][i+1].Count += 1
			_data[j][i+2].Count += 1
		}
		// 找顺子 2+1籁 左中籁
		if i <= 11 && _data[j][i].Count > 0 && _data[j][i+1].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j][i+1].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i+1].Value + 1},
				Value: []int32{_data[j][i].Value, _data[j][i+1].Value, _data[j][i+1].Value + 1},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j][i+1].Count += 1
		}
		// 找顺子 2+1籁 左籁右
		if i <= 11 && _data[j][i].Count > 0 && _data[j][i+2].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j][i+2].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i].Value + 1},
				Value: []int32{_data[j][i].Value, _data[j][i].Value + 1, _data[j][i+2].Value},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j][i+2].Count += 1
		}
		// 找顺子 2+1籁 籁中右
		if i > 1 && i <= 12 && _data[j][i].Count > 0 && _data[j][i+1].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j][i+1].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i].Value - 1},
				Value: []int32{_data[j][i].Value - 1, _data[j][i].Value, _data[j][i+1].Value},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j][i+1].Count += 1
		}
		// 找顺子 1+2籁 左籁籁
		if i <= 11 && _data[j][i].Count > 0 && laiZi > 1 {
			_data[j][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i].Value + 1, _data[j][i].Value + 2},
				Value: []int32{_data[j][i].Value, _data[j][i].Value + 1, _data[j][i].Value + 2},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-2)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-2)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
		}
		// 找顺子 1+2籁 籁中籁
		if i > 1 && i < 13 && _data[j][i].Count > 0 && laiZi > 1 {
			_data[j][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i].Value - 1, _data[j][i].Value + 1},
				Value: []int32{_data[j][i].Value - 1, _data[j][i].Value, _data[j][i].Value + 1},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-2)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-2)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
		}
		// 找顺子 1+2籁 籁籁右
		if i > 2 && _data[j][i].Count > 0 && laiZi > 1 {
			_data[j][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				LaiZi: []int32{_data[j][i].Value - 2, _data[j][i].Value - 1},
				Value: []int32{_data[j][i].Value - 2, _data[j][i].Value - 1, _data[j][i].Value},
				Type:  1,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-2)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-2)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
		}
		// 找群组 3
		if j == 0 && _data[j][i].Count > 0 && _data[j+1][i].Count > 0 && _data[j+3][i].Count > 0 {
			_data[j][i].Count -= 1
			_data[j+1][i].Count -= 1
			_data[j+3][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+1][i].Value, _data[j+3][i].Value},
				Type:  0,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+1][i].Count += 1
			_data[j+3][i].Count += 1
		}
		if j == 0 && _data[j][i].Count > 0 && _data[j+2][i].Count > 0 && _data[j+3][i].Count > 0 {
			_data[j][i].Count -= 1
			_data[j+2][i].Count -= 1
			_data[j+3][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+2][i].Value, _data[j+3][i].Value},
				Type:  0,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+2][i].Count += 1
			_data[j+3][i].Count += 1
		}
		if j <= 1 && _data[j][i].Count > 0 && _data[j+1][i].Count > 0 && _data[j+2][i].Count > 0 {
			_data[j][i].Count -= 1
			_data[j+1][i].Count -= 1
			_data[j+2][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+1][i].Value, _data[j+2][i].Value},
				Type:  0,
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+1][i].Count += 1
			_data[j+2][i].Count += 1
		}
		// 找群组 2+1籁
		if j == 0 && _data[j][i].Count > 0 && _data[j+1][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+1][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+1][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+1][i].Count += 1
		}
		if j == 0 && _data[j][i].Count > 0 && _data[j+2][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+2][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+2][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+2][i].Count += 1
		}
		if j == 0 && _data[j][i].Count > 0 && _data[j+3][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+3][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+3][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+3][i].Count += 1
		}
		if j == 1 && _data[j][i].Count > 0 && _data[j+1][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+1][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+1][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+1][i].Count += 1
		}
		if j == 1 && _data[j][i].Count > 0 && _data[j+2][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+2][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+2][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+2][i].Count += 1
		}
		if j == 2 && _data[j][i].Count > 0 && _data[j+1][i].Count > 0 && laiZi > 0 {
			_data[j][i].Count -= 1
			_data[j+1][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j+1][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-1)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-1)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
			_data[j+1][i].Count += 1
		}
		// 找群组 1+2籁
		if _data[j][i].Count > 0 && laiZi > 1 {
			_data[j][i].Count -= 1
			_ret = append(_ret, &ComInfo{
				Value: []int32{_data[j][i].Value, _data[j][i].Value, _data[j][i].Value},
				Type:  0,
				LaiZi: []int32{_data[j][i].Value, _data[j][i].Value},
			})

			if _data[j][i].Count > 0 {
				diGuiZuPai(_data, cList, rInfo, _ret, k, laiZi-2)
			} else {
				diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi-2)
			}

			_ret = _ret[:len(_ret)-1]
			_data[j][i].Count += 1
		}
		// 计算结果
		if k == int32(len(cList)-1) {
			diGuiZuPai(_data, cList, rInfo, _ret, k+1, laiZi)
		}
	}
}

func showResultInfo(_rInfo *ResultInfo) {
	if _rInfo.Value != 0 {
		fmt.Printf("组合总数:%d 总值:%d\n", len(_rInfo.Info), _rInfo.Value)
		for k, v := range _rInfo.Info {
			fmt.Printf("組合(%d) %d ", k, v.Type)
			for i := 0; i < len(v.Value); i++ {
				fmt.Printf("0x%02x ", v.Value[i])
			}
			for i := 0; i < len(v.LaiZi); i++ {
				fmt.Printf("癞子 0x%02x ", v.LaiZi[i])
			}
			fmt.Println("")
		}
		fmt.Printf("剩余手牌:")
		for i := 0; i < len(_rInfo.Card); i++ {
			fmt.Printf("0x%02x ", _rInfo.Card[i])
		}
		fmt.Println("")
	}
}

// 组牌
func ZuPai(cList []int32, laiZi int32) *ResultInfo {
	// 打印手牌
	//show(cList)
	// 手牌转换为 红黄蓝黑 列表
	_data, _cards := hCardToDCard(cList)
	//showData(_data)
	// 找出所有能组的牌
	_ret := []*ComInfo{}
	_rInfo := &ResultInfo{}
	diGuiZuPai(_data, _cards, _rInfo, _ret, 0, laiZi)
	// 最大组合
	//showResultInfo(_rInfo)
	return _rInfo
}

func QunZu(handCard []int32) {
	//handCard = []int32{0x01, 0x11, 0x21, 0x01, 0x11, 0x21, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	_info := &ResultInfo{}
	start := time.Now().UnixNano() / 1e3
	_info = ZuPai(handCard, 0)
	showResultInfo(_info)
	end := time.Now().UnixNano() / 1e3
	fmt.Printf("耗时：%0.3fms\n", float32(end-start)/1000)
	fmt.Println(AllCount)
}

func QunZuLaiZi(handCard []int32) {
	//handCard = []int32{0x01, 0x11, 0x21, 0x01, 0x11, 0x21, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	// 组牌(群组 顺组)
	_info := &ResultInfo{}
	start := time.Now().UnixNano() / 1e3
	_info = ZuPai(handCard, 2)
	// 打印最大结果
	showResultInfo(_info)
	end := time.Now().UnixNano() / 1e3
	fmt.Printf("耗时：%0.3fms\n", float32(end-start)/1000)
	fmt.Println(AllCount)
}

func main() {
	// 1.随机一副牌
	cList := []int32{}
	cList = append(cList, CardList...)
	cList = append(cList, CardList...)
	rand.Seed(time.Now().Unix())
	for i := 0; i < len(cList); i++ {
		rd := rand.Intn(len(cList))
		cList[i], cList[rd] = cList[rd], cList[i]
	}
	// 2.获取出 21张牌
	handCard := cList[:22]
	// 3.排序 降序
	for i := 0; i < len(handCard)-1; i++ {
		for j := i + 1; j < len(handCard); j++ {
			if handCard[j] > handCard[i] {
				handCard[i], handCard[j] = handCard[j], handCard[i]
			}
		}
	}
	QunZu(handCard)
	QunZuLaiZi(handCard)
}
