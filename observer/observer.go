package observer

import (
	"bytes"
	"filemonitor/storage"
	pool "filemonitor/util"
	util "filemonitor/util"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/match"
	"golang.org/x/net/html"

	uuid "github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
	httpAuth "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var r, err = git.PlainOpen(os.Getenv("FILEMONITOR"))

// The dynamic link function only supports src elements so far
func getDynamicLink(a goquery.Document, pattern string) string {

	link := "none"
	a.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")

		if match.Match(src, pattern) == true {
			link = src
		}

	})

	return link

}

func work(args ...interface{}) interface{} {

	URL := args[0].(string)
	dynamic := args[1].(bool)
	id := args[2].(uuid.UUID)
	pattern := args[3].(string)
	beautify := args[4].(bool)

	resp, err := util.CreateRequest(URL)
	if err != nil {
		log.Println("[*] This link seems to be offline.")

		return false
	}

	body := util.ReadBuffer(resp)

	// goquery wants a io reader
	parse, err := html.Parse(io.Reader(bytes.NewReader(body)))
	doc := goquery.NewDocumentFromNode(parse)
	if err != nil {
		log.Println("[*] We could not parse the HTML of this page.")
		return false
	}

	w, err := r.Worktree()
	if err != nil {

		log.Println(err)

	}

	if dynamic == false {

		log.Println("[*] Fetching " + URL + " ....")
		storage.SaveFile(id, body, URL, *w, beautify)

	} else {

		finalPath := getDynamicLink(*doc, pattern)

		if finalPath != "none" {

			var finalURL string

			if strings.HasPrefix(finalPath, "https") || strings.HasPrefix(finalPath, "http") {

				finalURL = finalPath

			} else {

				// We have only a path, it seems. We need the hostname to fetch it (thus we take it from the URL given and remove the path).
				u, err := url.Parse(URL)
				if err != nil {

					log.Println("[*] DYNAMIC ANALYSIS: I could not read this URL.")
					return false
				}

				URL = u.Scheme + "://" + u.Host
				finalURL = fmt.Sprintf("%s/%s", URL, finalPath)

			}

			log.Println("[*] DYNAMIC ANALYSIS: Fetching " + finalURL + "....")

			// In case of a relative path, we're just adding a / to be safe. We might want to change this later on
			resp, err := util.CreateRequest(finalURL)

			if err != nil {
				log.Println("[*] DYNAMIC ANALYSIS:  Dynamic analysis could not find the link. This link is possibly offline or removed from the page.")

				return false
			}

			storage.SaveFile(id, util.ReadBuffer(resp), URL, *w, beautify)

		} else {
			log.Println("[*] DYNAMIC ANALYSIS: We could not get the dynamic link of " + URL)
		}

	}

	return true
}

func Start(concurrency int) bool {

	var urlsToFetch = storage.GetEntries()

	// Set up the worker pool with X amount of workers
	fetchPool := pool.NewPool(concurrency)
	fetchPool.Run()

	//Add the jobs
	for _, element := range urlsToFetch.Urls {
		fetchPool.Add(work, element.URL, element.Dynamic, element.ID, element.Pattern, element.Beautify)
	}

	fetchPool.Wait()

	fetchPool.Stop()

	err = r.Push(&git.PushOptions{Auth: &httpAuth.BasicAuth{
		Username: "filemonitor_token_leak_and_get_pwned",
		Password: os.Getenv("GIT_TOKEN"),
	}, Progress: os.Stdout})
	if err != nil {
		log.Println("[*] We're done, but we could not push to the remote git repository. There's possibly nothing to push.")
		log.Println(err)
	}

	return true

}
