package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	chrwr "reg_parser/pkg/chromedp_wrapper"
)

func main() {
	var rootUrl string
	flag.StringVar(&rootUrl, "u", "https://www.google.com/", "number of lines to read from the file")
	flag.Parse()

	fmt.Println(rootUrl)

	// create context
	ctx, cancel := chrwr.Init()
	defer cancel()

	c := chrwr.NewChromeWrapper()
	c.OpenURL(ctx, rootUrl)

	// chr := chromedp_wrapper.NewChromeWrapper()

	// err := chr.OpenURL(ctx, rootUrl)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = chr.WaitLoaded(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	// if err != nil {
	// 	log.Fatal(err)
	// }
}
