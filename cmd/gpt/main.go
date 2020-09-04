package main

import (
	"flag"
	"gpt"
	"log"
)

func main() {
	mainFilePath := flag.String("main", "", "main file path")
	libDirPath := flag.String("lib", "", "lib dir path")
	flag.Parse()

	if *mainFilePath == "" || *libDirPath == "" {
		log.Print("Usage: ./gpt -main <main_file_path> -lib <library_dir_path>")
		return
	}

	gpt.Generate(*mainFilePath, *libDirPath)
}
