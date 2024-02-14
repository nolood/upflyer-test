package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/nolood/upflyer-test.git/internal/config"
)

func DateOfFirstPost(channelName string) time.Time {
	c := colly.NewCollector()

	c.SetRequestTimeout(10 * time.Second)

	c.OnError(func(_ *colly.Response, err error) {
		config.Logger.Error(err.Error())
	})

	var dates []time.Time

	c.OnHTML("a.tgme_widget_message_date time", func(e *colly.HTMLElement) {
		date, err := time.Parse(time.RFC3339, e.Attr("datetime"))
		if err != nil {
			config.Logger.Error(err.Error())
		}
		dates = append(dates, date)
	})

	c.Visit(fmt.Sprintf("https://t.me/s/%s/1", channelName))

	return dates[1]
}

func SubscribersCount(channelName string) int {
	var err error

	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		config.Logger.Error(err.Error())
	})

	var subsCount int

	c.OnHTML("div.tgme_page_extra", func(e *colly.HTMLElement) {
		countWithString := strings.Split(e.Text, " ")
		stringCount := countWithString[:len(countWithString)-1]
		subsCount, err = strconv.Atoi(strings.Join(stringCount, ""))
		if err != nil {
			config.Logger.Error(err.Error())
		}
	})

	c.Visit(fmt.Sprintf("https://t.me/%s", channelName))

	return subsCount
}

func ViewsAndCountFromLastWeekPosts(channelName string) (int, []int, error) {
	c := colly.NewCollector()

	c.SetRequestTimeout(10 * time.Second)

	c.OnError(func(_ *colly.Response, err error) {
		config.Logger.Error(err.Error())
	})

	var postViews []int

	c.OnHTML("div.tgme_widget_message_info", func(e *colly.HTMLElement) {
		dateString := e.ChildAttr("span.tgme_widget_message_meta a.tgme_widget_message_date time", "datetime")
		date, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			config.Logger.Error(err.Error())
			return
		}

		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		if date.After(sevenDaysAgo) {
			viewCountString := e.ChildText("span.tgme_widget_message_views")
			viewCount, err := convertViewCount(viewCountString)
			if err != nil {
				config.Logger.Error(err.Error())
				return
			}
			postViews = append(postViews, viewCount)
		}
	})

	c.Visit(fmt.Sprintf("https://t.me/s/%s", channelName))

	return len(postViews), postViews, nil
}

func convertViewCount(viewCountString string) (int, error) {
	viewCountString = strings.TrimSpace(strings.Replace(viewCountString, ",", "", -1))

	if strings.Contains(viewCountString, "K") {
		viewCountString = strings.Replace(viewCountString, "K", "", -1)
		viewCountFloat, err := strconv.ParseFloat(viewCountString, 64)
		if err != nil {
			return 0, err
		}
		viewCount := int(viewCountFloat * 1000)
		return viewCount, nil
	}

	if strings.Contains(viewCountString, "M") {
		viewCountString = strings.Replace(viewCountString, "M", "", -1)
		viewCountFloat, err := strconv.ParseFloat(viewCountString, 64)
		if err != nil {
			return 0, err
		}
		viewCount := int(viewCountFloat * 1000000)
		return viewCount, nil
	}

	viewCount, err := strconv.Atoi(viewCountString)
	if err != nil {
		return 0, err
	}
	return viewCount, nil
}
