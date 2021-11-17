// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/lacuna-seo/luminati-client"
	"log"
	"sync"
	"testing"
)

// Amount is the amount of go routines to spawn to retrieve
// from the luminati API.
const Amount = 50

// Hammer the Luminati API.
func Test_Hammer(t *testing.T) {
	client, err := luminati.New(&Cache{})
	if err != nil {
		log.Fatalln(err)
	}

	opts := luminati.Options{
		Keyword:  "macbook",
		CheckURL: "https://apple.com",
		Country:  "us",
	}

	wg := sync.WaitGroup{}

	var (
		errors = 0
		success = 0
	)
	for i := 0; i < Amount; i++ {
		wg.Add(1)
		go run(client, opts, &wg, errors, success)
	}
	wg.Wait()

	fmt.Printf("\n%v Finsihed SERP Hammer Test\n\n", emoji.ChartIncreasing)
	fmt.Printf("%v Total Errors: %d\n", emoji.CrossMark, errors)
	fmt.Printf("%v Total Success: %d", emoji.CheckMarkButton, success)
}

// Runs and prints the response.
func run(client luminati.KeywordFinder, opts luminati.Options, wg *sync.WaitGroup, errors int, success int) {
	defer wg.Done()
	_, err := client.JSON(opts)
	if err != nil {
		fmt.Printf("%v Error: %s\n", emoji.CrossMark, err.Error())
		errors++
		return
	}
	fmt.Printf("%v Sucess:  %s\n", emoji.CheckMarkButton, opts.Keyword)
	success++
}


