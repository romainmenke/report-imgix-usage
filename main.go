package main

import (
	"fmt"

	"github.com/romainmenke/report-imgix-usage/prompt"
)

func main() {
	sources := prompt.PromptAuthAndGetSources()
	if sources == nil {
		fmt.Println("no sources found or invalid auth")
		return
	}

	for prompt.PromptMode(sources) {
		// doing the interactive cli thingy here
	}
}
