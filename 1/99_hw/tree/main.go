package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(w io.Writer, path string, pf bool) error {
	printDir(path, "")

	return nil
}

func printDir(path string, preffix string) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for i, enity := range dirs {
		connector := "├───"

		isLast := i+1 == len(dirs)

		if isLast {
			connector = "└───"
		}

		printer(preffix + connector + enity.Name() + "\n")

		if enity.IsDir() {
			var newPreffix string

			if isLast {
				newPreffix = preffix + "\t"
			} else {
				newPreffix = preffix + "│\t"
			}

			printDir(path+string(os.PathSeparator)+enity.Name(), newPreffix)
		}
	}
}

func printer(s string) {
	fmt.Print(s)
}
