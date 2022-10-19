// https://pkg.go.dev/github.com/asticode/go-astilectron#section-readme
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

//mby split a code to files for clarity

var window *astilectron.Window
var app *astilectron.Astilectron

type FormData struct {
	Service    string `json:"service"`
	Charge     string `json:"charge"`
	AccountNum string `json:"accountNum"`
	Due        string `json:"due"`
}

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

	var data FormData
	temp, _ := io.ReadAll(r.Body)
	log.Println(string(temp))

	json.Unmarshal(temp, &data)

	fileData, err := os.OpenFile("data.dat", os.O_CREATE|os.O_WRONLY, 0644)
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

func main() {
	http.HandleFunc("/addData", addData)
	go http.ListenAndServe(":8080", nil)

	app, _ = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:            "test",
		VersionAstilectron: "0.49.0",
		VersionElectron:    "20.0.0",
	})
	defer app.Close()

	app.Start()
	// http://localhost:8080/index
	window, _ = app.NewWindow("./ui/index.html", &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Height:    astikit.IntPtr(600),
		Width:     astikit.IntPtr(600),
		Resizable: astikit.BoolPtr(false),
	})
	window.Create()
	defer window.Close()

	listen()

	// window.OpenDevTools()

	app.Wait()
}

func listen() {
	window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var tempWindow *astilectron.Window
		// Unmarshal
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
			// tempWindow.OpenDevTools()
			tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
				var mess string
				mes.Unmarshal(&mess)
				if mess == "exit" {
					tempWindow.Close()
				}
				return nil
			})
		case "show-data":
			log.Println("showing data")
			HTML, err := os.OpenFile("./ui/showData.html", os.O_CREATE, 0644)
			if err != nil {
				log.Println(err)
			}
			defer HTML.Close()

			// firstPart := `
			// <!DOCTYPE html>
			// <html lang="en">
			// <head>
			// <meta charset="UTF-8">
			// <meta http-equiv="X-UA-Compatible" content="IE=edge">
			// <meta name="viewport" content="width=device-width, initial-scale=1.0">
			// <link rel="stylesheet" href="./style.css">
			// 	<title>Document</title>
			// </head>
			// <body>
			// <table id="dataTable">
			// <tr>
			// <th>ZA CO</th>
			// <th>ILE</th>
			// <th>NR KONTA</th>
			// <th>DO KIEDY</th>
			// </tr>`

			// lastPart := `
			// </table>
			// <input type="button" value="exit" id="btnExit">
			// <input type="button" value="save" id="btnSave">
			// <script>
			// 	document.addEventListener('astilectron-ready', function(){
			// 		btnExit.addEventListener('click', function(){
			// 			astilectron.sendMessage("exit");
			// 		})
			// 		btnSave.addEventListener('click', function(){
			// 			astilectron.sendMessage("save");
			// 			console.log("save");

			// 			let rows = dataTable.rows.length - 1;

			// 			console.log(rows);

			// 			astilectron.sendMessage("rows")

			// 		})
			// 	})
			// </script>
			// <input type=button onClick=window.location.reload()>
			// </body>
			// </html>`

			file, _ := os.OpenFile("data.dat", os.O_RDONLY, 0644)
			defer file.Close()

			data, _ := io.ReadAll(file)

			log.Println(string(data))

			split := strings.Split(string(data), "\n")

			log.Println(split)

			// toRead := make([]string, 0)
			// for _, v := range split {
			// 	temp := strings.Split(v, "/")
			// 	toRead = append(toRead, temp...)
			// }
			// log.Println(toRead)

			// HTML.Write([]byte("<tr>"))
			// for i, v := range toRead {
			// 	HTML.Write([]byte("<td>"))
			// 	HTML.Write([]byte(v))
			// 	HTML.Write([]byte("</td>"))
			// 	if i%4 == 0 && i != 0 {
			// 		HTML.Write([]byte("</tr><tr>"))
			// 	}
			// }
			// HTML.Write([]byte("</tr>"))

			// HTML.Write([]byte(firstPart))
			// HTML.Write([]byte((data)))
			// HTML.Write([]byte(lastPart))

			// tempWindow, _ = app.NewWindow("./ui/showData.html", &astilectron.WindowOptions{
			// 	Center:    astikit.BoolPtr(true),
			// 	Height:    astikit.IntPtr(600),
			// 	Width:     astikit.IntPtr(600),
			// 	Resizable: astikit.BoolPtr(false)})
			// tempWindow.Create()

			// // tempWindow.OpenDevTools()

			// tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
			// 	var tempMess string
			// 	mes.Unmarshal(&tempMess)
			// 	switch tempMess {
			// 	case "exit":
			// 		file.Close()
			// 		HTML.Close()
			// 		tempWindow.Close()
			// 	case "save":
			// 		tempWindow.On("window.event.message", func(e astilectron.Event) (deleteListener bool) {
			// 			var t string
			// 			e.Message.Unmarshal(&t)
			// 			log.Println(t)
			// 			return true
			// 		})

			//add deleting certain <tr></tr> from a file

			//cant use loop

			// offset := len(firstPart)
			// toSave := make([]byte, len(data))
			// HTML.ReadAt(toSave, int64(offset))
			// log.Println(string(toSave))

			// HTML.Close()
			// file.Close()
			// }
			// return nil
			// })
		case "exit":
			window.Close()
			os.Exit(0)
		}

		return nil
	})
}
