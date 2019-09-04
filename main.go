package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kapytein/filemonitor/observer"
	"github.com/kapytein/filemonitor/storage"
)

func main() {

	fmt.Println("FILEMONITOR v0.0.1 - Monitoring files at your wish.")

	if os.Getenv("FILEMONITOR") == "" || os.Getenv("GIT_TOKEN") == "" {

		log.Fatalln("You haven't set the environment variable for the local git repository or Git(Hub|Lab|Whatever) token.")

	} else {

		url := flag.String("url", "", "This is the URL of the webpage to track (or the webpage the link is on if it is a dynamic link)")
		pattern := flag.String("pattern", "", "This is the pattern where we will look for on the web page (in case of a dynamic link)")
		fetch := flag.Bool("fetch", false, "Use this option if you want to start fetching the files.")
		threads := flag.Int("threads", 5, "If you choose to fetch, this is the amount of threads (concurrency) which will be used. Default is 5.")
		beautify := flag.Bool("beautify", false, "Use this if you would like to have the file 'JS beautified' when saving.")

		flag.Parse()

		if *fetch == true {
			log.Println("[!] Starting to fetch all links..")
			observer.Start(*threads)
		} else if *url != "" {
			if *pattern != "" {
				log.Println("[!] Adding new (dynamic) entry..")
				storage.AddURL(*url, true, *pattern, *beautify)
			} else {
				log.Println("[!] Adding new entry..")
				storage.AddURL(*url, false, "", *beautify)
			}
		} else {
			fmt.Println("[!] No option selected")
		}
	}

}
