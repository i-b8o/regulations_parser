package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	chrwr "reg_parser/pkg/chromedp_wrapper"
	"reg_parser/pkg/logger"
)

type RegulationResponse struct {
	RegulationID string `json:"regulation_id"`
}
type ChapterResponse struct {
	ChapterID string `json:"chapter_id"`
}

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

	// Regulation
	s, err := getRegulationName(ctx, abbreviation, c)
	if err != nil {
		logger.Error(err)

	}
	logger.Info(s)
	respBody, err := sendPOST("http://localhost:10000/r", s)
	if err != nil {
		logger.Error(err)
	}

	var r RegulationResponse

	err = json.Unmarshal(respBody, &r)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Pause
	fmt.Scanln()

	// Chapter
	s, err = getChapterInfo(ctx, r.RegulationID, c)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(s)

	respBody, err = sendPOST("http://localhost:10000/c", s)
	if err != nil {
		logger.Error(err)
	}

	var ch ChapterResponse

	err = json.Unmarshal(respBody, &ch)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("\nChapter ID: ", ch.ChapterID)

	// Paragrahs
	s, err = getParagraphs(ctx, ch.ChapterID, c)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(s)

	// respBody, err = sendPOST("http://localhost:10000/p", s)
	// if err != nil {
	// 	logger.Error(err)
	// }

	// var chapterResponse ChapterResponse

	// err = json.Unmarshal(respBody, &chapterResponse)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// logger.Info("\nChapter ID: ", chapterResponse.ChapterID)
}

func sendPOST(url, payload string) ([]byte, error) {

	var jsonStr = []byte(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}

func getRegulationName(ctx context.Context, abbreviation string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, jsRegulation(abbreviation))
}

func getChapterInfo(ctx context.Context, regulationID string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, jsChapter(regulationID))
}

func getParagraphs(ctx context.Context, chapterID string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, jsParagraphs(chapterID))
}
