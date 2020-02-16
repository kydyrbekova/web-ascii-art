package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var text, fonts, Result string

func main() {
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/html/", http.StripPrefix("/html/", fs))
	http.HandleFunc("/", handlefunc)
	fmt.Println("Listening to port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handlefunc(w http.ResponseWriter, request *http.Request) {

	if request.URL.Path != "/" {
		http.Error(w, "Error 404\n Page not Found", http.StatusNotFound)
		return
	}

	switch request.Method {
	case "GET":
		http.ServeFile(w, request, "html")
	case "POST":

		fonts := request.FormValue("fonts")
		text := request.FormValue("text")

		arr := []byte(text)
		var newArr []byte
		for _, ch := range arr {
			if ch != 13 {
				newArr = append(newArr, ch)
			}
		}

		text = string(newArr)

		for _, e := range text {
			if (e < 32 || e > 126) && e != 10 {
				http.Error(w, "Error 400\n Bad Request,"+string(e), http.StatusNotFound)
				return
			}
		}
		file, err := os.Open("./" + fonts + ".txt")

		if err != nil {
			http.Error(w, "Internal Server error - 500", http.StatusNotFound)
			return
		}



		fileContent := ScanFile(file)
		arg := strings.Split(text+" ", "\\n")

		for _, val := range arg {
			asciiResult(string(val), fileContent)
		}

		temp, _ := template.ParseFiles("html/index.html")
		temp.Execute(w, Result)
	}
}

func asciiResult(s string, fileContent []string) string {
	for i := 1; i <= 8; i++ {
		for _, arg := range s {
			indexBase := int(rune(arg)-32) * 9
			Result = Result + string(fileContent[indexBase+i])
		}
		Result = Result + "\n"
	}
	return Result
}

func ScanFile(font *os.File) []string {
	Result = ""
	var fileContent []string
	scanner := bufio.NewScanner(font)
	for scanner.Scan() {
		lines := scanner.Text()
		fileContent = append(fileContent, lines)
	}
	return fileContent
}

// func handleDownload() {

// }
