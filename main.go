package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/i-b8o/chrome"
)

func main() {
	var rootUrl string
	flag.StringVar(&rootUrl, "u", "", "number of lines to read from the file")
	flag.Parse()

	fmt.Println(rootUrl)
	ctxt, cancel := chrome.Init("~/.config/chromium/Default/")
	defer cancel()

	err := chrome.OpenURL(ctxt, rootUrl, true)
	if err != nil {
		log.Fatal(err)
	}
	err = chrome.WaitLoaded(ctxt)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
}
