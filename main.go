// https://pkg.go.dev/github.com/asticode/go-astilectron#section-readme
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var window *astilectron.Window
var app *astilectron.Astilectron

type FormData struct {
	Service    string `json:"service"`
	Charge     string `json:"charge"`
	AccountNum string `json:"accountNum"`
	Due        string `json:"due"`
}

// handling /addData route
func addData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	file, err := os.Open("./ui/addData.html")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	addData, _ := io.ReadAll(file)
	fmt.Fprintln(w, string(addData))

	log.Println(r.Body)

	// receiving form data
	var data FormData
	temp, _ := io.ReadAll(r.Body)
	log.Println(string(temp))

	json.Unmarshal(temp, &data)

	// writing form data to a file
	if data.Due != "" {
		fileData, err := os.OpenFile("data.dat", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Println(err)
		}
		defer fileData.Close()

		v := reflect.ValueOf(data)
		log.Println(v)

		for i := 0; i < 4; i++ {
			if i == 3 {
				fileData.WriteString(v.Field(i).String())
				fileData.WriteString("\n")
			} else {
				fileData.WriteString(v.Field(i).String() + "|")
			}
		}
	}

}

func loadCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	file, err := os.Open("./ui/style.css")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	css, _ := io.ReadAll(file)
	fmt.Fprintln(w, string(css))
}

func main() {
	http.HandleFunc("/addData", addData)
	http.HandleFunc("/style.css", loadCSS)
	http.HandleFunc("/showData", func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.OpenFile("data.dat", os.O_APPEND, 0644)
		defer file.Close()
		data, _ := io.ReadAll(file)
		fmt.Println(data)
		split := strings.Split(string(data), "\n")
		log.Println(split)

		dataSlice := make([][]string, 0)
		for _, v := range split {
			if v != "" {
				dataSlice = append(dataSlice, strings.Split(v, "|"))
			}
		}

		fmt.Println(dataSlice)

		tmplt, err := template.ParseFiles("./ui/showData.html")
		if err != nil {
			fmt.Println(err)
		}
		tmplt.Execute(w, dataSlice)
	})
	go http.ListenAndServe(":8080", nil)

	app, _ = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:            "test",
		VersionAstilectron: "0.49.0",
		VersionElectron:    "20.0.0",
	})
	defer app.Close()

	app.Start()
	window, _ = app.NewWindow("./ui/index.html", &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Height:    astikit.IntPtr(600),
		Width:     astikit.IntPtr(600),
		Resizable: astikit.BoolPtr(false),
	})
	window.Create()
	defer window.Close()

	events()
	app.Wait()
}

func events() {
	window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var tempWindow *astilectron.Window
		var s string
		m.Unmarshal(&s)

		switch s {
		case "add-data":
			log.Println("adding new data")
			tempWindow, _ = app.NewWindow("http://localhost:8080/addData", &astilectron.WindowOptions{
				Center:    astikit.BoolPtr(true),
				Height:    astikit.IntPtr(600),
				Width:     astikit.IntPtr(600),
				Resizable: astikit.BoolPtr(false)})
			tempWindow.Create()

			tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
				var mess string
				mes.Unmarshal(&mess)
				if mess == "exit" {
					tempWindow.Close()
				}
				return nil
			})
		case "show-data":
			tempWindow, _ = app.NewWindow("http://localhost:8080/showData", &astilectron.WindowOptions{
				Center:    astikit.BoolPtr(true),
				Height:    astikit.IntPtr(600),
				Width:     astikit.IntPtr(600),
				Resizable: astikit.BoolPtr(false)})
			tempWindow.Create()

			tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
				var tempMess string
				mes.Unmarshal(&tempMess)
				log.Println(tempMess)

				switch tempMess {
				case "exit":
					tempWindow.Close()
				default:
					file, _ := os.OpenFile("data.dat", os.O_APPEND, 0644)
					defer file.Close()

					data, _ := io.ReadAll(file)

					log.Println(tempMess)
					line, err := strconv.Atoi(tempMess)
					if err != nil {
						log.Println(err)
					}
					log.Println(line)
					lines := bytes.Split(data, []byte("\n"))
					log.Println(lines, len(lines))
					del := removeLine(lines, line)
					log.Println(del, len(del))

					os.Truncate(file.Name(), 0)
					for _, v := range del {
						if string(v) != "" {
							file.Write(v)
							file.WriteString("\n")
						}
					}
				}

				// file.Close()
				return nil
			})
		case "exit":
			window.Close()
			os.Exit(0)
		}

		return nil
	})
}

func removeLine(s [][]byte, i int) [][]byte {
	log.Println("i: ", i, ", len(s): ", len(s))
	if i == len(s)-1 {
		return s[:len(s)-2]
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
