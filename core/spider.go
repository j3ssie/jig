package core

import (
	"fmt"
	"time"

	"github.com/ysmood/kit"
	"github.com/ysmood/rod"
	"github.com/ysmood/rod/lib/input"
	"github.com/ysmood/rod/lib/launcher"
)

// Rod provides a lot of debug options, you can use set methods to enable them or use environment variables
// list at "lib/defaults".
func DebugMode() {
	url := launcher.New().
		Headless(false). // run chrome on foreground, you can also use env "rod=show"
		Devtools(true). // open devtools for each new tab
		Launch()

	browser := rod.New().
		ControlURL(url).
		Trace(true). // show trace of each input action
		Slowmotion(2 * time.Second). // each input action will take 2 second
		Connect().
		Timeout(time.Minute)

	// the monitor server that plays the screenshots of each tab, useful when debugging headlee mode
	browser.ServeMonitor(":9777")

	defer browser.Close()

	p1 := browser.Page("https://www.duckduckgo.com/")
	page := browser.Page("https://www.wikipedia.org/")
	fmt.Println(p1.Element(".html").Text())

	page.Element("#searchLanguage").Select("[lang=zh]")
	page.Element("#searchInput").Input("热干面")
	page.Keyboard.Press(input.Enter)

	fmt.Println(page.Element("#firstHeading").Text())

	// get the image binary
	img := page.Element(`[alt="Hot Dry Noodles.jpg"]`)
	_ = kit.OutputFile("tmp/img.jpg", img.Resource(), nil)

	// pause the js execution
	// you can resume by open the devtools and click the resume button on source tab
	page.Pause()

	// Skip
	// Output: 热干面
}

// Open wikipedia, search for "idempotent", and print the title of result page
func Basic() {
	// launch and connect to a browser
	url := launcher.New().
		Headless(false). // run chrome on foreground, you can also use env "rod=show"
		// Devtools(true).  // open devtools for each new tab
		Launch()

	browser := rod.New().
		ControlURL(url).
		Trace(true). // show trace of each input action
		// Slowmotion(2 * time.Second). // each input action will take 2 second
		Connect().
		Timeout(time.Minute)
	// browser := rod.New().Connect()

	// Even you forget to close, rod will close it after main process ends
	defer browser.Close()

	// timeout will be passed to chained function calls
	page := browser.Timeout(time.Minute).Page("https://github.com")

	// make sure windows size is consistent
	page.Window(0, 0, 1200, 600)

	// use css selector to get the search input element and input "git"
	page.Element("input").Input("git").Press(input.Enter)

	// wait until css selector get the element then get the text content of it
	text := page.Element(".codesearch-results p").Text()
	html := page.Element("html").HTML()

	fmt.Println(text)
	fmt.Println(html)

	// Output: Git is the most widely used version control system.
}
