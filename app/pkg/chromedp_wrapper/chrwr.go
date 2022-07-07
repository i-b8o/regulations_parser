package chromedp_wrapper

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Chrome struct {
	timeOut int
}

func Init() (context.Context, context.CancelFunc) {
	opts := []chromedp.ExecAllocatorOption{chromedp.Flag("no-sandbox", true)}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(allocContext)
}

func NewChromeWrapper() *Chrome {
	return &Chrome{timeOut: 60}
}

func openURL(url string, message *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptOpenURL(url), message),
	}
}

func (c *Chrome) OpenURL(ctxt context.Context, url string) error {
	var message string
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, openURL(url, &message)))
	return err
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func waitVisible(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	}
}

func (c *Chrome) WaitVisible(ctxt context.Context, selector string) error {
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitVisible(selector)))
	return err
}

func waitReady(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitReady(selector, chromedp.ByQuery),
	}
}

func (c *Chrome) WaitReady(ctxt context.Context, selector string) error {
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitReady(selector)))
	return err
}

func getString(jsString string, resultString *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetString(jsString), resultString),
	}
}

func (c *Chrome) GetString(ctxt context.Context, jsString string) (string, error) {
	var resultString string
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getString(jsString, &resultString)))
	return resultString, err
}

func getStringsSlice(jsString string, resultSlice *[]string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetStringsSlice(jsString), resultSlice),
	}
}

func GetStringsSlice(ctxt context.Context, jsString string) ([]string, error) {
	var stringSlice []string
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getStringsSlice(jsString, &stringSlice)))
	return stringSlice, err
}

func getBool(jsBool string, resultBool *bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(scriptGetBool(jsBool), resultBool),
	}
}

func GetBool(ctxt context.Context, jsBool string) (bool, error) {
	var resultBool bool
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, getBool(jsBool, &resultBool)))
	return resultBool, err
}

func click(selector string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(selector, chromedp.ByQuery),
	}
}

func Click(ctxt context.Context, selector string) error {
	err := chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, waitVisible(selector)))
	if err != nil {
		return err
	}
	return chromedp.Run(ctxt, RunWithTimeOut(&ctxt, 60, click(selector)))

}
