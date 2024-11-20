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
	printDir(w, path, "")

	return nil
}

func printDir(w io.Writer, path string, preffix string) {
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

		if !enity.IsDir() {
			// TODO: size: (<size>b)
			size, err := fileSize(path + string(os.PathSeparator) + enity.Name())
			if err != nil {
				panic(err)
			}

			printer(w, preffix+connector+enity.Name()+" "+size+"\n")
		} else {
			printer(w, preffix+connector+enity.Name()+"\n")
		}

		if enity.IsDir() {
			var newPreffix string

			if isLast {
				newPreffix = preffix + "\t"
			} else {
				newPreffix = preffix + "│\t"
			}

			printDir(w, path+string(os.PathSeparator)+enity.Name(), newPreffix)
		}
	}
}

func printer(w io.Writer, s string) {
	w.Write([]byte(s))
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
