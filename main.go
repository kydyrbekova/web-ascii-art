package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"text/template"
)

type Vihod struct {
	Output, Fs, String string
}

func main() {

	http.HandleFunc("/", mainPage)
	fs := http.FileServer(http.Dir("css"))

	http.Handle("/css/", http.StripPrefix("/css/", fs))

	port := ":3000"
	println("Server listen on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Print("ListenAndServe", err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	// разобрать файл
	tmpl, err := template.ParseFiles("index.html")

	value := r.FormValue("text")
	valueStlye := r.FormValue("style")

	cmd := exec.Command("./ascii-art", value, valueStlye)
	cmd.Dir = "ascii"
	output, _ := cmd.Output()

	vihod := Vihod{Output: string(output), Fs: valueStlye, String: value}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, vihod); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}
