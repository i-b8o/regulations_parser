package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reg_parser/internal/script"
	chrwr "reg_parser/pkg/chromedp_wrapper"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type Response struct {
	ID       string   `json:"id,omitempty"`
	Err      error    `json:"err,omitempty"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Message  string   `json:"message,omitempty"`
}

func main() {
	log := logrus.New()
	f, err := os.OpenFile("testlogrus.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
	})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)
	log.SetLevel(logrus.WarnLevel)

	// log.Formatter = &logrus.TextFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	// 		filename := path.Base(f.File)
	// 		return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
	// 	},
	// 	DisableColors: false,
	// 	FullTimestamp: true,
	// }

	// log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	var rootUrl, abbreviation string
	flag.StringVar(&rootUrl, "u", "", "start url")
	flag.StringVar(&abbreviation, "a", "", "abbreviation")
	flag.Parse()

	if len(rootUrl) == 0 {
		log.Error("url is empty")
		return
	} else if len(abbreviation) == 0 {
		log.Error("abbreviation is empty")
		return
	}

	ctx, cancel := chrwr.Init()
	defer cancel()

	log.Info("Chrome wrapper initialisation")
	c := chrwr.NewChromeWrapper()

	log.Info("openning url ", rootUrl)
	err = c.OpenURL(ctx, rootUrl)
	if err != nil {
		log.Error(err)
		return
	}

	// Regulation
	s, err := getRegulationName(ctx, abbreviation, c)
	if err != nil {
		log.Error(err)
	}
	log.Info(s)
	response, err := sendPOST("http://localhost:10000/r", s, log)
	if err != nil {
		log.Error(err)
	}
	if len(response.Errors) > 0 {
		log.Error(response.Errors)
	}
	// Pause
	fmt.Scanln()

	btnNextExists, err := btnNextExistsFunc(ctx, c)
	if err != nil {
		log.Error(err)
	}
	regulationID := response.ID
	for btnNextExists {
		time.Sleep(3 * time.Second)
		// Chapter
		s, err = getChapterInfo(ctx, regulationID, c)
		if err != nil {
			log.Error(err)
		}
		log.Info(s)

		response, err = sendPOST("http://localhost:10000/c", s, log)
		if err != nil {
			log.Error(err)
		}
		if len(response.Errors) > 0 {
			log.Error(response.Errors)
		}

		if len(response.ID) == 0 {
			log.Error("chapter.ID is empty")
		}
		// Paragrahs
		s, err = getParagraphs(ctx, response.ID, c)
		if err != nil {
			log.Error(err)
		}
		log.Info(s)

		response, err = sendPOST("http://localhost:10000/p", s, log)
		if err != nil {
			log.Error(err)
		}
		if len(response.Errors) > 0 {
			log.Error(response.Errors)
		}

		if len(response.Warnings) > 0 {
			log.Warn(response.Warnings)
		}

		log.Info("\nErrors: ", response.Errors)

		btnNextExists, err = btnNextExistsFunc(ctx, c)
		if err != nil {
			log.Error(err)
		}

		if btnNextExists {
			err = btnNextClick(ctx, c)
			if err != nil {
				log.Error(err)
			}
		}

	}

}

func sendPOST(url, payload string, log *logrus.Logger) (Response, error) {
	var response Response
	var jsonStr = []byte(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return response, err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	log.Info(string(respBody))
	err = json.Unmarshal(respBody, &response)

	return response, err

}

func getRegulationName(ctx context.Context, abbreviation string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, script.JSRegulation(abbreviation))
}

func getChapterInfo(ctx context.Context, regulationID string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, script.JSChapter(regulationID))
}

func getParagraphs(ctx context.Context, chapterID string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}
	return c.GetString(ctx, script.JSParagraphs(chapterID))
}

func btnNextExistsFunc(ctx context.Context, c *chrwr.Chrome) (bool, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return false, err
	}
	return c.GetBool(ctx, `document.querySelectorAll(".pages__right").length > 0`)
}

func btnNextClick(ctx context.Context, c *chrwr.Chrome) error {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return err
	}
	// return c.GetBool(ctx, `document.querySelectorAll(".pages__right").length > 0`)
	return c.Click(ctx, ".pages__right")
}
