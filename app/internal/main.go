package main

import (
	"bufio"
	"flag"
	"os"
	chrwr "reg_parser/pkg/chromedp_wrapper"
	"reg_parser/pkg/chromedp_wrapper/logger"
)

func main() {
	var rootUrl string
	flag.StringVar(&rootUrl, "u", "https://www.google.com/", "number of lines to read from the file")
	flag.Parse()

	// logger

	// create context
	ctx, cancel := chrwr.Init()
	defer cancel()
	logger := logger.NewLogger()
	logger.Info("Chrome wrapper initialisation")
	c := chrwr.NewChromeWrapper()
	logger.Info("openning url %s", rootUrl)
	err := c.OpenURL(ctx, rootUrl)
	if err != nil {
		logger.Error(err)
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return
	}

}
