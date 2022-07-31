package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/rivo/tview"
)

const trendingURL = "https://github.com/trending/"
const baseURL = "https://github.com"

var qs = []*survey.Question{
	{
		Name: "language",
		Prompt: &survey.Select{
			Message: "Which language trending?",
			Options: []string{"All", "Go", "Javascript", "Ruby", "TypeScript", "Python"},
			Default: "All",
		},
	},
}

func main() {
	var language string
	err := survey.Ask(qs, &language)
	if err != nil {
		log.Fatal(err)
	}
	var URL string
	if language == "All" {
		URL = trendingURL
	} else {
		URL = fmt.Sprintf("%s%s", trendingURL, language)
	}
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	app := tview.NewApplication()
	list := tview.NewList()
	doc.Find(".Box-row").Each(func(i int, s *goquery.Selection) {
		repo, _ := s.Find("h1.h3.lh-condensed > a").Attr("href")
		title := s.Find("p.col-9.color-fg-muted.my-1.pr-4").Text()
		repoURL := fmt.Sprintf("%s%s", baseURL, repo)
		list.AddItem(repoURL, title, 1, func() {
			err := exec.Command("open", repoURL).Run()
			if err != nil {
				log.Fatal(err)
			}
		})
	})
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
