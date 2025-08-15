package seeds

import (
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

func holidaySeed(db *bun.DB) error {
	holidays := []model.Holiday{
		{
			Name:        "วันขึ้นปีใหม่",
			Description: "วันขึ้นปีใหม่ 2568",
			IsActive:    true,
			StartDate:   1735689600, // 1 January 2025 00:00:00 +07:00
			EndDate:     1735775999, // 1 January 2025 23:59:59 +07:00
		},
		{
			Name:        "วันตรุษจีน",
			Description: "วันตรุษจีน 2568",
			IsActive:    true,
			StartDate:   1738202400, // 29 January 2025 00:00:00 +07:00
			EndDate:     1738288799, // 29 January 2025 23:59:59 +07:00
		},
		{
			Name:        "วันมาฆบูชา",
			Description: "วันมาฆบูชา 2568",
			IsActive:    true,
			StartDate:   1739343600, // 12 February 2025 00:00:00 +07:00
			EndDate:     1739429999, // 12 February 2025 23:59:59 +07:00
		},
		{
			Name:        "จุดจบของเราะมะฎอน",
			Description: "จุดจบของเราะมะฎอน 2568",
			IsActive:    true,
			StartDate:   1743404400, // 30 March 2025 00:00:00 +07:00
			EndDate:     1743490799, // 30 March 2025 23:59:59 +07:00
		},
		{
			Name:        "วันจักรี",
			Description: "วันจักรี 2568",
			IsActive:    true,
			StartDate:   1744009200, // 6 April 2025 00:00:00 +07:00
			EndDate:     1744095599, // 6 April 2025 23:59:59 +07:00
		},
		{
			Name:        "วันหยุดชดเชยวันจักรี",
			Description: "วันหยุดชดเชยวันจักรี 2568",
			IsActive:    true,
			StartDate:   1744095600, // 7 April 2025 00:00:00 +07:00
			EndDate:     1744181999, // 7 April 2025 23:59:59 +07:00
		},
		{
			Name:        "วันสงกรานต์",
			Description: "วันสงกรานต์ 2568 (13-15 เมษายน)",
			IsActive:    true,
			StartDate:   1744614000, // 13 April 2025 00:00:00 +07:00
			EndDate:     1744872799, // 15 April 2025 23:59:59 +07:00
		},
		{
			Name:        "วันแรงงาน",
			Description: "วันแรงงาน 2568",
			IsActive:    true,
			StartDate:   1746169200, // 1 May 2025 00:00:00 +07:00
			EndDate:     1746255599, // 1 May 2025 23:59:59 +07:00
		},
		{
			Name:        "พระราชพิธีบรมราชาภิเษกของพระมหากษัตริย์ไทย",
			Description: "พระราชพิธีบรมราชาภิเษกของพระมหากษัตริย์ไทย 2568",
			IsActive:    true,
			StartDate:   1746428400, // 4 May 2025 00:00:00 +07:00
			EndDate:     1746514799, // 4 May 2025 23:59:59 +07:00
		},
		{
			Name:        "วันหยุดชดเชยพระราชพิธีบรมราชาภิเษกของพระมหากษัตริย์ไทย",
			Description: "วันหยุดชดเชยพระราชพิธีบรมราชาภิเษกของพระมหากษัตริย์ไทย 2568",
			IsActive:    true,
			StartDate:   1746514800, // 5 May 2025 00:00:00 +07:00
			EndDate:     1746601199, // 5 May 2025 23:59:59 +07:00
		},
		{
			Name:        "วันพืชมงคล",
			Description: "วันพืชมงคล 2568",
			IsActive:    true,
			StartDate:   1746860400, // 9 May 2025 00:00:00 +07:00
			EndDate:     1746946799, // 9 May 2025 23:59:59 +07:00
		},
		{
			Name:        "วันวิสาขบูชา",
			Description: "วันวิสาขบูชา 2568",
			IsActive:    true,
			StartDate:   1747033200, // 11 May 2025 00:00:00 +07:00
			EndDate:     1747119599, // 11 May 2025 23:59:59 +07:00
		},
		{
			Name:        "วันหยุดชดเชยวันวิสาขบูชา",
			Description: "วันหยุดชดเชยวันวิสาขบูชา 2568",
			IsActive:    true,
			StartDate:   1747119600, // 12 May 2025 00:00:00 +07:00
			EndDate:     1747205999, // 12 May 2025 23:59:59 +07:00
		},
		{
			Name:        "วันหยุดชดเชยวันเฉลิมพระชนมพรรษา พระราชินี",
			Description: "วันหยุดชดเชยวันเฉลิมพระชนมพรรษา พระราชินี 2568",
			IsActive:    true,
			StartDate:   1748934000, // 2 June 2025 00:00:00 +07:00
			EndDate:     1749020399, // 2 June 2025 23:59:59 +07:00
		},
		{
			Name:        "วันเฉลิมพระชนมพรรษา พระราชินี",
			Description: "วันเฉลิมพระชนมพรรษา พระราชินี 2568",
			IsActive:    true,
			StartDate:   1749020400, // 3 June 2025 00:00:00 +07:00
			EndDate:     1749106799, // 3 June 2025 23:59:59 +07:00
		},
		{
			Name:        "วันอาสาฬหบูชา",
			Description: "วันอาสาฬหบูชา 2568",
			IsActive:    true,
			StartDate:   1752102000, // 10 July 2025 00:00:00 +07:00
			EndDate:     1752188399, // 10 July 2025 23:59:59 +07:00
		},
		{
			Name:        "วันเข้าพรรษา",
			Description: "วันเข้าพรรษา 2568",
			IsActive:    true,
			StartDate:   1752188400, // 11 July 2025 00:00:00 +07:00
			EndDate:     1752274799, // 11 July 2025 23:59:59 +07:00
		},
		{
			Name:        "วันเกิดของพระบาทสมเด็จพระเจ้าอยู่หัว",
			Description: "วันเกิดของพระบาทสมเด็จพระเจ้าอยู่หัว 2568",
			IsActive:    true,
			StartDate:   1753657200, // 28 July 2025 00:00:00 +07:00
			EndDate:     1753743599, // 28 July 2025 23:59:59 +07:00
		},
		{
			Name:        "วันหยุดชดเชยวันแม่แห่งชาติ",
			Description: "วันหยุดชดเชยวันแม่แห่งชาติ 2568",
			IsActive:    true,
			StartDate:   1754866800, // 11 August 2025 00:00:00 +07:00
			EndDate:     1754953199, // 11 August 2025 23:59:59 +07:00
		},
		{
			Name:        "วันแม่แห่งชาติ",
			Description: "วันแม่แห่งชาติ 2568",
			IsActive:    true,
			StartDate:   1754953200, // 12 August 2025 00:00:00 +07:00
			EndDate:     1755039599, // 12 August 2025 23:59:59 +07:00
		},
		{
			Name:        "วันคล้ายวันสวรรคต พระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช",
			Description: "วันคล้ายวันสวรรคต พระบาทสมเด็จพระปรมินทรมหาภูมิพลอดุลยเดช 2568",
			IsActive:    true,
			StartDate:   1760252400, // 13 October 2025 00:00:00 +07:00
			EndDate:     1760338799, // 13 October 2025 23:59:59 +07:00
		},
		{
			Name:        "วันปิยมหาราช",
			Description: "วันปิยมหาราช 2568",
			IsActive:    true,
			StartDate:   1761116400, // 23 October 2025 00:00:00 +07:00
			EndDate:     1761202799, // 23 October 2025 23:59:59 +07:00
		},
		{
			Name:        "วันคล้ายวันพระราชสมภพ รัชกาลที่ 9",
			Description: "วันคล้ายวันพระราชสมภพ รัชกาลที่ 9 ปี 2568",
			IsActive:    true,
			StartDate:   1764658800, // 5 December 2025 00:00:00 +07:00
			EndDate:     1764745199, // 5 December 2025 23:59:59 +07:00
		},
		{
			Name:        "วันรัฐธรรมนูญ",
			Description: "วันรัฐธรรมนูญ 2568",
			IsActive:    true,
			StartDate:   1765090800, // 10 December 2025 00:00:00 +07:00
			EndDate:     1765177199, // 10 December 2025 23:59:59 +07:00
		},
		{
			Name:        "คริสต์มาส",
			Description: "คริสต์มาส 2568",
			IsActive:    true,
			StartDate:   1766386800, // 25 December 2025 00:00:00 +07:00
			EndDate:     1766473199, // 25 December 2025 23:59:59 +07:00
		},
		{
			Name:        "วันส่งท้ายปีเก่า",
			Description: "วันส่งท้ายปีเก่า 2568",
			IsActive:    true,
			StartDate:   1766905200, // 31 December 2025 00:00:00 +07:00
			EndDate:     1766991599, // 31 December 2025 23:59:59 +07:00
		},
	}
	_, err := db.NewInsert().Model(&holidays).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
