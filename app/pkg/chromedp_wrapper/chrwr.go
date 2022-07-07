package chromedp_wrapper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
)

var (
	badColor = color.New(color.FgHiRed, color.Bold)
	okColor  = color.New(color.FgHiCyan)
)

type Chrome struct {
	needLog bool
	timeOut int
}

func Init() (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{chromedp.Flag("no-sandbox", true)}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(allocContext)
}

func NewChromeWrapper() *Chrome {
	return &Chrome{needLog: true, timeOut: 60}
}

func openURL(url string, message *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptOpenURL(url), message),
	}
}

func (c *Chrome) OpenURL(ctxt context.Context, url string) {
	if c.needLog {
		c := color.New(color.FgGreen)
		c.Printf("Opening page url %s - ", url)
	}

	var message string
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, openURL(url, &message)))
	if err != nil {

		log.Fatal(err)
	}

	if c.needLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Println("Ok!.")
	}
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}
