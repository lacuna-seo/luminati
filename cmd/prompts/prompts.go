// Copyright 2021 Reddico. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prompts

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"os"
)

// Prompt defines a prompt for obtaining information
// to pass to the Runner
type Prompt struct {
	Label        string
	ErrMsg       string
	Instructions string
	Items []string
	Validate bool
	PostProcess  func(s string) string
}

var (
	// KeywordPrompt defines the prompt for obtaining
	// the keyword data.
	KeywordPrompt = Prompt{
		Label:        "Keyword",
		ErrMsg:       "Enter a keyword",
		Instructions: "Enter a keyword for the domain.",
		Validate: true,
		PostProcess:  nil,
	}
	// URLPrompt defines the prompt for obtaining
	// the client name.
	URLPrompt = Prompt{
		Label:        "URL",
		Instructions: "Optional: Enter a URL to compare against, defaults to ''",
		PostProcess:  nil,
	}
	// CountryPrompt defines the prompt for obtaining
	// a country to search for.
	CountryPrompt = Prompt{
		Label:        "Country",
		Instructions: "Optional: Enter a country, defaults to 'uk'",
	}
	// JSONPrompt defines the prompt for obtaining
	// the output data.
	JSONPrompt = Prompt{
		Label: "Select HTML or JSON output",
		Items: []string{"JSON", "HTML"},
	}
)

// Get retrieves a prompt and prints out any instructions
// and the runs the post process hook. Returns the
// prompt value upon success.
func (p *Prompt) Get() string {
	fmt.Println(p.Instructions)
	if len(p.Items) != 0 {
		prompt := promptui.Select{
			Label: p.Label,
			Items: p.Items,
		}
		_, result, err := prompt.Run()
		if err != nil {
			Exit(err)
		}
		return result
	}
	prompt := promptui.Prompt{
		Label: p.Label,
		Validate: func(s string) error {
			if p.Validate && s == "" {
				return errors.New(p.ErrMsg)
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		Exit(err)
	}
	fmt.Printf("\n")
	if p.PostProcess != nil {
		return p.PostProcess(result)
	}
	return result
}

// Exit exits' the application with an error.
func Exit(err error) {
	color.Red.Println("Error: " + err.Error())
	os.Exit(1)
}
