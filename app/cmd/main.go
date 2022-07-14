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
	"strconv"
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
	i := 0
	for btnNextExists {
		time.Sleep(3 * time.Second)
		netError, err := netError(ctx, c)
		if netError {
			fmt.Scanln()
		}
		// Chapter
		i++
		iStr := strconv.FormatInt(int64(i), 10)
		s, err = getChapterInfo(ctx, regulationID, iStr, c)
		if err != nil {
			log.Error(err)
		}
		log.Info(s)

		fmt.Println(s)
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

func getChapterInfo(ctx context.Context, regulationID, chapterOrderNum string, c *chrwr.Chrome) (string, error) {
	err := c.WaitLoaded(ctx)
	if err != nil {
		return "", err
	}

	chapterExist, err := c.GetBool(ctx, script.JSCheckChapter)

	if err != nil {
		return "", err
	}

	if !chapterExist {
		fmt.Println("Where is h1?")
		fmt.Scanln()
	}

	paragraphExists, err := c.GetBool(ctx, script.JSCheckParagraphs)
	if err != nil {
		return "", err
	}

	if !paragraphExists {
		fmt.Println("Where are paragraphs?")
		fmt.Scanln()
	}

	return c.GetString(ctx, script.JSChapter(regulationID, chapterOrderNum))
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

func netError(ctx context.Context, c *chrwr.Chrome) (netErr bool, err error) {

	err = c.WaitLoaded(ctx)
	if err != nil {
		return true, err
	}
	return c.GetBool(ctx, `document.getElementsByClassName("neterror").length > 0`)
}
