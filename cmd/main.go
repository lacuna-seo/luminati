// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/enescakir/emoji"
	"github.com/gookit/color"
	"github.com/lacuna-seo/luminati"
	"github.com/lacuna-seo/luminati/cmd/prompts"
	"time"
)

func main() {
	fmt.Printf("\n %v Welcome to the Lacuna Luminati Client...\n\n", emoji.WavingHand)

	keyword := prompts.KeywordPrompt.Get()
	url := prompts.URLPrompt.Get()
	country := prompts.CountryPrompt.Get()
	color.Yellow.Print("Use K and J to Navigate")
	output := prompts.JSONPrompt.Get()

	if country == "" {
		country = "uk"
	}

	fmt.Printf("%v Keyword: %s\n", emoji.OpenBook, keyword)
	fmt.Printf("%v URL: %s\n", emoji.Link, url)
	fmt.Printf("%v Country: %s\n", emoji.GlobeShowingEuropeAfrica, country)
	fmt.Printf("%v Output: %s\n\n", emoji.Information, output)

	s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	s.Suffix = " Getting results from Luminati"
	s.Start()
	fmt.Println()

	time.Sleep(1 * time.Second)

	client, err := luminati.New(cache, time.Hour*1)
	if err != nil {
		prompts.Exit(err)
	}

	opts := luminati.Options{
		Keyword: keyword,
		Country: country,
	}

	if output == "HTML" {
		pp.Fatalln(client.HTML(opts))
	}

	serps, meta, err := client.JSON(opts)
	if err != nil {
		prompts.Exit(err)
	}

	pp.Fatalln(serps)
	pp.Fatalln(serps.CheckURL(url))
	pp.Fatalln(meta)
}
