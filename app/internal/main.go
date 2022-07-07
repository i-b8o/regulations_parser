package main

import (
	"flag"
	chrwr "reg_parser/pkg/chromedp_wrapper"
	"reg_parser/pkg/chromedp_wrapper/logger"
)

func main() {
	logger := logger.NewLogger()

	var rootUrl, abbreviation string
	flag.StringVar(&rootUrl, "u", "", "start url")
	flag.StringVar(&abbreviation, "a", "", "abbreviation")
	flag.Parse()

	if len(rootUrl) == 0 {
		logger.Error("url is empty")
		return
	} else if len(abbreviation) == 0 {
		logger.Error("abbreviation is empty")
		return
	}

	ctx, cancel := chrwr.Init()
	defer cancel()

	logger.Info("Chrome wrapper initialisation")
	c := chrwr.NewChromeWrapper()

	logger.Info("openning url %s", rootUrl)
	err := c.OpenURL(ctx, rootUrl)
	if err != nil {
		logger.Error(err)
		return
	}

}
