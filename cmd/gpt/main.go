package main

import (
	"gpt"
)

func main() {
	gpt.Generate("testdata/src/a/a.go", "testdata/src/a/lib")
}
