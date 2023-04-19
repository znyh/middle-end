package calc

import (
	"errors"
	"fmt"

	pb "calculator/api/calculator/v1"
	"calculator/internal/base"

	"github.com/go-kratos/kratos/v2/log"
)

func BuildColorCards(req *pb.CalcReq, cards []int32) (err error, left []int32, builds [][]int32) {
	err, items := checkSetInput(req, cards)
	if err != nil {
		log.Errorf("%+v", err)
		return err, nil, nil
	}
	s := newSet(cards, req.CalcBody.Close, req.CalcBody.Kind, req.CalcBody.Continue, items)
	if err = checkSetOut(s, cards); err != nil {
		log.Errorf("%+v", err)
		return err, nil, nil
	}

	use := []int32(nil)
	for _, v := range s.items {
		use = append(use, v.l.hands...)
		builds = append(builds, v.l.hands)
	}

	left = base.SliceCopyAndDel(base.SliceCopy(cards), base.SliceCopy(use)...)
	return
}

func Test(req *pb.CalcReq, cards []int32) (str []string) {
	err, items := checkSetInput(req, cards)
	if err != nil {
		log.Errorf("%+v", err)
		return []string{err.Error()}
	}

	s := newSet(cards, req.CalcBody.Close, req.CalcBody.Kind, req.CalcBody.Continue, items)
	if err = checkSetOut(s, cards); err != nil {
		log.Errorf("%+v", err)
		return
	}

	str = append(str, fmt.Sprintf("close:%+v kind:%+v continue:%+v\n", s.close, s.kind, s.continues))
	for _, v := range s.items {
		str = append(str, fmt.Sprintf("%+v\n", v))
	}
	return str
}

//校验输入
func checkSetInput(req *pb.CalcReq, cards []int32) (err error, items []tagItem) {

	if req.CalcBody.Close || len(req.CalcBody.Items) <= 0 || len(cards) <= 0 {
		err = errors.New(fmt.Sprintf("config illegal parameter. req:[%+v] cards:[%+v]", req, cards))
		return
	}

	if req.CalcBody.Kind != 1 && req.CalcBody.Kind != 2 {
		err = errors.New(fmt.Sprintf("config illegal parameter.(bad kind:%+v) req:[%+v]", req.CalcBody.Kind, req))
		return
	}

	for _, v := range req.CalcBody.Items {
		items = append(items, tagItem{
			index: len(items),
			c: tagCondition{
				master: v.Master,
				slave:  v.Slave,
				shun:   v.Shun,
				ke:     v.Ke,
				dui:    v.Dui,
				ca:     v.Ca,
			},
			l: tagList{
				hands:  nil,
				master: nil,
				slave:  nil,
				left:   nil,
			},
		})

		if v.Master > MAXCOUNT ||
			v.Slave > MAXCOUNT ||
			(v.Master+v.Slave) > MAXCOUNT ||
			3*v.Shun+3*v.Ke+2*v.Dui+2*v.Ca > MAXCOUNT {
			err = errors.New(fmt.Sprintf("config illegal parameter. error:[%+v] req:[%+v] ", v, req))
			return err, nil
		}
	}

	return nil, items
}

//校验输出
func checkSetOut(s *tagSet, cards []int32) (err error) {
	use := []int32(nil)
	for _, v := range s.items {
		if len(v.l.hands) != MAXCOUNT {
			err = errors.New(fmt.Sprintf("===> error.(len(v.l.hands) != MAXCOUNT) %+v", s))
			return
		}
		use = append(use, v.l.hands...)
	}
	if len(use) > 0 && !base.SliceContainAll(cards, use...) {
		err = errors.New(fmt.Sprintf("===> error:!base.SliceContainAll(s.cards,use...), \ncards:%+v\nuse:%+v",
			descCards(cards), descCards(use)))
		return
	}
	return nil
}
