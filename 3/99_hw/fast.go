package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func FastSearch(out io.Writer) {
	/*
		!!! !!! !!!
		обратите внимание - в задании обязательно нужен отчет
		делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировали их
		так же обратите внимание на команду в параметре -http
		перечитайте еще раз задание
		!!! !!! !!!
	*/
	//SlowSearch(out)
	FastSearchV1(out)
}

func FastSearchV1(out io.Writer) {
	/*
		- Читаем файл построчно;
		- Каждая строка - json, а значит структура
	*/
	_, _ = fmt.Fprintln(out, "found users:")

	type user struct {
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Browsers []string `json:"browsers"`
	}

	uniqueBrowsers := make(map[string]struct{})

	u := &user{}

	file, err := os.Open(filePath)
	if err != nil {
		_, _ = fmt.Fprintln(out, err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic("error close file: " + err.Error())
		}
	}(file)

	rd := bufio.NewReader(file)
	for i := 0; true; i++ {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			_, _ = fmt.Fprintf(out, "read file line error: %v", err)
			return
		}

		err = json.Unmarshal(line, &u)
		if err != nil {
			_, _ = fmt.Fprintf(out, "unmarshal json error: %v", err)
		}

		isAndroid := false
		isMSIE := false

		for _, b := range u.Browsers {
			if strings.Contains(strings.ToLower(b), "android") {
				isAndroid = true
				uniqueBrowsers[b] = struct{}{}
			} else if strings.Contains(strings.ToLower(b), "msie") {
				uniqueBrowsers[b] = struct{}{}
				isMSIE = true
			}

		}

		if isAndroid && isMSIE {
			u.Email = fmt.Sprintf("<" + strings.Replace(u.Email, "@", " [at] ", -1) + ">")
			_, _ = fmt.Fprintf(out, "[%d] %s %s\n", i, u.Name, u.Email)
		}
	}
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Total unique browsers", len(uniqueBrowsers))

}
