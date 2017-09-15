package parser

import (
	"log"
)

type pageView struct {
	Url             string
	PageViews       int
	UniquePageViews int
	AvgTime         string
	BounceRate      float32
}

func (p *pageView) isDefault() bool {
	empty := pageView{}
	if p.Url != empty.Url {
		return false
	}
	if p.PageViews != empty.PageViews {
		return false
	}
	if p.UniquePageViews != empty.UniquePageViews {
		return false
	}
	if p.AvgTime != empty.AvgTime {
		return false
	}
	if p.BounceRate != empty.BounceRate {
		return false
	}
	return true
}

func ParsePageViewsFile() []pageView {
	const filename = "data/page-views.csv"
	records, _ := ImportCSV(filename)

	pageViews := make([]pageView, 1)

	//skip headers
	for _, record := range records[1:] {
		newPageView := pageView{}
		var err error

		newPageView.Url, err = NormalizeUrl(record[0])
		if err != nil {
			log.Print(err)
			continue
		}

		newPageView.PageViews, err = StringToInt(record[1])
		if err != nil {
			log.Print(err)
			continue
		}

		newPageView.UniquePageViews, err = StringToInt(record[2])
		if err != nil {
			log.Print(err)
			continue
		}
		newPageView.AvgTime = record[3]
		newPageView.BounceRate, err = StringPercentToFloat(record[5])
		if err != nil {
			log.Print(err)
			continue
		}

		if !newPageView.isDefault() {
			pageViews = append(pageViews, newPageView)
		}
	}
	return pageViews
}
