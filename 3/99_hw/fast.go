package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
			if strings.Contains(b, "Android") {
				isAndroid = true
				uniqueBrowsers[b] = struct{}{}
			} else if strings.Contains(b, "MSIE") {
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

// Только регулярки заменены на strings
func FastSearchV0_1(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	users := make([]map[string]interface{}, 0)
	for _, line := range lines {
		user := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {
		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		if mail, ok := user["email"].(string); ok {
			mail := strings.Replace(mail, "@", " [at] ", -1)
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], mail)
		}

	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
