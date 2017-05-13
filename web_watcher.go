package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ngiyshhk/rxgo-examples/util"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

type WebWatcher struct {
	Duration  time.Duration
	Url       string
	ElementAt string
	SlackUrl  string
}

type WebWatcherFactory struct{}

func (_ WebWatcherFactory) Create() *WebWatcher {
	url := flag.String("t", "", "target url")
	elementAt := flag.String("e", "", "target element")
	duration := flag.Int("s", 10, "access duration[second]")
	slackUrl := flag.String("p", "", "slack incoming hook url")
	flag.Parse()

	return &WebWatcher{Url: *url, ElementAt: *elementAt, Duration: time.Duration(*duration) * time.Second, SlackUrl: *slackUrl}
}

func (ww WebWatcher) ErrFilter(v interface{}) bool {
	err, is_error := v.(error)
	if is_error {
		fmt.Errorf("%v\n", err)
	}
	return !is_error
}

func (ww WebWatcher) Crawl(v interface{}) interface{} {
	doc, err := goquery.NewDocument(ww.Url)
	if err != nil {
		return err
	}
	return doc
}

func (ww WebWatcher) Scrape() func(interface{}) interface{} {
	before := make([]string, 0)
	return func(v interface{}) interface{} {
		doc := v.(*goquery.Document)
		latest := doc.Find(ww.ElementAt).Map(func(_ int, s *goquery.Selection) string {
			return s.Text()
		})

		updated := util.Diff(before, latest)
		before = latest
		if len(updated) == 0 {
			return errors.New("not updated.")
		}
		return updated
	}
}

func (ww WebWatcher) Sink() observer.Observer {
	return observer.Observer{
		NextHandler: func(v interface{}) {
			items := v.([]string)
			jsonStr := fmt.Sprintf(`{"text": "%s"}`, strings.Join(items, "\n"))

			req, err := http.NewRequest(
				"POST",
				ww.SlackUrl,
				bytes.NewBuffer([]byte(jsonStr)),
			)
			if err != nil {
				return
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		},

		ErrHandler: func(err error) {
			fmt.Printf("Encountered error: %v\n", err)
		},
	}
}

func main() {
	ww := WebWatcherFactory{}.Create()

	src := observable.Interval(make(chan struct{}), ww.Duration)

	<-src.Map(ww.Crawl).Filter(ww.ErrFilter).Map(ww.Scrape()).Filter(ww.ErrFilter).Subscribe(ww.Sink())
}
