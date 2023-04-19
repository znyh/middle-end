package calc

import (
	"testing"
	"time"

	pb "calculator/api/calculator/v1"
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

func TestToMap(t *testing.T) {
	cards := base.SliceCopy(oneCards)
	cm := toColorMap(cards)
	log.Infof("\ncc:%+v\ncm:%+v\n", cm, descCards(cards))
}

func TestPermute(t *testing.T) {

	var (
		count = 1

		calculate = 0

		start = time.Now()

		cards = []int32{
			11, 12, 13, 14, 15, 16, 17, 18, 19,
			11, 12, 13, 14, 15, 16, 17, 18, 19,
			11, 12, 13, 14, 15, 16, 17, 18, 19,
			11, 12, 13, 14, 15, 16, 17, 18, 19,
		}
	)

	for i := 0; i < count; i++ {

		for shun := int32(0); shun <= 4; shun++ {
			for ke := int32(0); ke <= 4; ke++ {
				for dui := int32(0); dui <= 6; dui++ {
					for ca := int32(0); ca <= 6; ca++ {

						if 3*shun+3*ke+2*dui+2*ca > MAXCOUNT {
							continue
						}

						calculate++

						permute(cards, &tagCondition{master: 0, slave: 0, shun: shun, ke: ke, dui: dui, ca: ca})

					}
				}
			}
		}

	}

	log.Infof("count:%+v calculate:%+v use:%+v/ms", count, calculate, time.Since(start).Milliseconds())
	log.Infof("------------------")
}

func TestPermute2(t *testing.T) {
	start := time.Now()
	cards := []int32{
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,
		11, 12, 13, 14, 15, 16, 17, 18, 19,
	}

	c := &tagCondition{master: 0, slave: 0, shun: 0, ke: 1, dui: 1, ca: 2}

	permute(cards, c)

	log.Infof("use:%+v/ms", time.Since(start).Milliseconds())
}

func TestBuildColorCards(t *testing.T) {

	cards := base.SliceCopy(oneCards)
	start := time.Now()
	req := &pb.CalcReq{
		CalcBody: &pb.CalcReq_CalcBody{
			Close:    false,
			Kind:     int32(2),
			Continue: false,
			Items: []*pb.Item{
				{Master: 0, Slave: 0, Shun: 2, Ke: 0, Dui: 2, Ca: 0},
				{Master: 0, Slave: 0, Shun: 0, Ke: 2, Dui: 0, Ca: 1},
				{Master: 0, Slave: 0, Shun: 1, Ke: 0, Dui: 1, Ca: 0},
				{Master: 0, Slave: 0, Shun: 1, Ke: 1, Dui: 0, Ca: 0},
			},
		},
	}
	_, left, builds := BuildColorCards(req, cards)
	log.Infof("cost:%+v/ms cards:%+v use:%+v left:%+v", time.Since(start).Milliseconds(), len(cards), len(cards)-len(left), len(left))
	for _, v := range builds {
		log.Infof("%+v", v)
	}

}

func TestNewSet(t *testing.T) {

	var (
		count  = 0
		errCnt = 0
		cards  = base.SliceCopy(oneCards)
		items  = []*pb.Item(nil)
	)

	for j := 0; j < 4; j++ {

		for master := int32(0); master <= MAXCOUNT; master++ {
			for slave := int32(0); slave <= MAXCOUNT; slave++ {

				if master+slave > MAXCOUNT {
					continue
				}

				for kind := 1; kind <= 2; kind++ {
					for continues := 1; continues <= 2; continues++ {
						for shun := int32(0); shun <= 4; shun++ {
							for ke := int32(0); ke <= 4; ke++ {
								for dui := int32(0); dui <= 6; dui++ {
									for ca := int32(0); ca <= 6; ca++ {

										if 3*shun+3*ke+2*dui+2*ca > MAXCOUNT {
											continue
										}

										if len(items) < 4 {
											items = append(items, &pb.Item{
												Master: master,
												Slave:  slave,
												Shun:   shun,
												Ke:     ke,
												Dui:    dui,
												Ca:     ca,
											})
											continue
										}

										req := &pb.CalcReq{
											CalcBody: &pb.CalcReq_CalcBody{
												Close:    false,
												Kind:     int32(kind),
												Continue: continues%2 == 0,
												Items:    items,
											},
										}

										if err, _, _ := BuildColorCards(req, cards); err != nil {
											errCnt++
										}

										count++
										items = []*pb.Item(nil)
									}
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
