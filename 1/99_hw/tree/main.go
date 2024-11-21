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
	// TODO: Вернуть ошибки в корень
	printDir(w, path, "", pf)
	return nil
}

func printDir(w io.Writer, path string, preffix string, pf bool) {
	dirs, err := os.ReadDir(path)
	if !pf {
		filtered := make([]os.DirEntry, 0)
		for _, entry := range dirs {
			if entry.IsDir() {
				filtered = append(filtered, entry)
			}
		}
		dirs = filtered
	}

	if err != nil {
		panic(err)
	}

	for i, enity := range dirs {
		connector := "├───"

		isLast := i+1 == len(dirs)

		if isLast {
			connector = "└───"
		}

		if !enity.IsDir() {
			size, err := fileSize(path + string(os.PathSeparator) + enity.Name())
			if err != nil {
				panic(err)
			}

			io.WriteString(w, preffix+connector+enity.Name()+" "+size+"\n")
		} else {
			io.WriteString(w, preffix+connector+enity.Name()+"\n")
		}

		if enity.IsDir() {
			var newPreffix string

			if isLast {
				newPreffix = preffix + "\t"
			} else {
				newPreffix = preffix + "│\t"
			}

			printDir(w, path+string(os.PathSeparator)+enity.Name(), newPreffix, pf)
		}
	}
}

// Возвращает размер файла в байтах
func fileSize(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fs, err := file.Stat()
	if err != nil {
		return "", err
	}

	if size := fs.Size(); size == 0 {
		return "(empty)", nil
	} else {
		return fmt.Sprintf("(%db)", fs.Size()), nil
	}
}
