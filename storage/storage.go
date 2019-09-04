package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	uuid "github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type TrackEntries struct {
	Urls []TrackEntry
}

type TrackEntry struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	Dynamic  bool      `json:"dynamic"`
	Pattern  string    `json:"pattern"`
	Beautify bool      `json:"beautify"`
}

var (
	filePath = os.Getenv("FILEMONITOR") + "/%s-%s.txt"
)

var entries TrackEntries

func AddURL(name string, dynamic bool, pattern string, beautify bool) {

	var err error
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	// make sure we got an updated entries global
	loadFile()

	var trackEntry = TrackEntry{ID: id, URL: name, Dynamic: dynamic, Pattern: pattern, Beautify: beautify}
	entries := append(entries.Urls, trackEntry)

	f, err := json.Marshal(entries)
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile("urls.json", f, 644)
	if err != nil {
		log.Println(err)
	}
}

func GetEntries() *TrackEntries {

	loadFile()

	return &entries
}

func loadFile() {
	var err error
	currentUrls, err := ioutil.ReadFile("urls.json")
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(currentUrls, &entries.Urls)
	if err != nil {
		log.Println(err)
	}

}

func SaveFile(id uuid.UUID, body []byte, url string, repo git.Worktree, beautify bool) bool {
	var err error
	var replacer = strings.NewReplacer("http://", "", "https://", "", "/", "")
	var options = jsbeautifier.DefaultOptions()

	url = replacer.Replace(url)
	if beautify == true {

		log.Println("[*] Beautifying the body content...")

		bodyBeautify := string(body[:])

		// Rewrite the file and beautify it. Maybe we should beautify the body directly
		jsbeautifier, err := jsbeautifier.Beautify(&bodyBeautify, options)
		if err != nil {
			log.Println(err)
		}

		body = []byte(jsbeautifier)
	}

	// If file does not exists, just create it and don't check the diff
	if _, err = os.Stat(fmt.Sprintf(filePath, url, id)); os.IsNotExist(err) {

		err = ioutil.WriteFile(fmt.Sprintf(filePath, url, id), body, 0644)
		if err != nil {
			log.Println("[*] We could not create this file.")
			return false
		}
		_, err := repo.Add(fmt.Sprintf(filePath, url, id))
		if err != nil {
			log.Println("[*] We could not add this file to the git repository.")
			return false
		}

		_, err = repo.Commit("Filemonitor: Boss, I am adding a new file to track.", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "FILEMONITOR",
				Email: "filemonitor@filemonitor.test",
				When:  time.Now(),
			},
		})

		if err != nil {
			log.Println("[*] We could not commit this file to the git repository.")
			return false
		}

	} else {

		dmp := diffmatchpatch.New()

		previousState, err := ioutil.ReadFile(fmt.Sprintf(filePath, url, id))
		if err != nil {
			log.Println(err)
			return false
		}

		// @TODO Use diffmanpatch to get the amount of differences of current + previous check, in order to eliminate
		// false positives (csrf tokens/timestamp changes on the webpage)
		diffs := dmp.DiffMain(string(previousState), string(body), false)

		isEqual := true
		for _, element := range diffs {
			// DiffEqual Operation = 0
			if element.Type != 0 {

				isEqual = false
				break
			}

		}

		if isEqual == false {

			err = ioutil.WriteFile(fmt.Sprintf(filePath, url, id), body, 0644)
			if err != nil {
				log.Println(err)
				return false
			}

			_, err = repo.Add(fmt.Sprintf("%s-%s.txt", url, id))
			if err != nil {
				log.Println(err)
				return false
			}

			_, err := repo.Commit("Filemonitor: Boss! Attention! There are some changes!", &git.CommitOptions{
				Author: &object.Signature{
					Name:  "FILEMONITOR",
					Email: "filemonitor@filemonitor.test",
					When:  time.Now(),
				},
			})
			if err != nil {
				log.Println("[*] We could not commit this file to the git repository.")
				log.Println(err)

				return false
			}

		}

	}

	return true

}
