package faster

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"os"
	"strings"
)

const filePath string = "./data/users.txt"

// easyjson:json
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	u := &User{}

	_, _ = fmt.Fprintln(out, "found users:")

	uniqueBrowsers := make(map[string]struct{})

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

		err = easyjson.Unmarshal(line, u)
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
