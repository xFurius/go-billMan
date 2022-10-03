package main

import (
	"io"
	"log"
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var window *astilectron.Window
var app *astilectron.Astilectron

// https://pkg.go.dev/github.com/asticode/go-astilectron#section-readme
func main() {
	app, _ = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:            "test",
		VersionAstilectron: "0.49.0",
		VersionElectron:    "20.0.0",
	})
	defer app.Close()

	app.Start()

	//opening .html file
	//appending everything to <body>
	//appending data
	//appending everything after </body>

	// page := `<!DOCTYPE html>
	// <html lang="en">
	// <head>
	// <meta charset="UTF-8">
	// <meta http-equiv="X-UA-Compatible" content="IE=edge">
	// <meta name="viewport" content="width=device-width, initial-scale=1.0">
	// <link rel="stylesheet" href="./style.css">
	// <title>Document</title>
	// </head>
	// <body>
	// <h3>Guibuibuioh</h3>
	// </body>
	// </html>`

	// //file, err := os.OpenFile("./ui/index.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// file, err := os.OpenFile("./ui/index.html", os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer file.Close()

	// file.Write([]byte(page))

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
			tempWindow, _ = app.NewWindow("./ui/addData.html", &astilectron.WindowOptions{
				Center:    astikit.BoolPtr(true),
				Height:    astikit.IntPtr(600),
				Width:     astikit.IntPtr(600),
				Resizable: astikit.BoolPtr(false)})
			tempWindow.Create()
			tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
				var tempMess string
				mes.Unmarshal(&tempMess)
				if tempMess == "exit" {
					tempWindow.Close()
				} else {
					file, err := os.OpenFile("data.dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						log.Println(err)
					}
					defer file.Close()

					file.Write([]byte(tempMess))
					file.Write([]byte(`<td><input type="button" value="usun" id="btnDel"></td></tr>`))
					file.Write([]byte("\n"))
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

			firstPart := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<link rel="stylesheet" href="./style.css">
				<title>Document</title>
			</head>
			<body>
			<table>
			<tr>
			<th>ZA CO</th>
			<th>ILE</th>
			<th>NR KONTA</th>
			<th>DO KIEDY</th>
			</tr>`

			lastPart := `
			</table>
			<input type="button" value="exit" id="btnExit">
			<input type="button" value="save" id="btnSave">		
			<script>
				document.addEventListener('astilectron-ready', function(){
					btnExit.addEventListener('click', function(){
						astilectron.sendMessage("exit");
					})
					btnSave.addEventListener('click', function(){
						astilectron.sendMessage("save");
					})
				})
			</script>
			<input type=button onClick=window.location.reload()>
			</body>
			</html>`

			file, _ := os.Open("data.dat")
			defer file.Close()

			data, _ := io.ReadAll(file)

			log.Println(data)

			HTML.Write([]byte(firstPart))
			HTML.Write([]byte((data)))
			HTML.Write([]byte(lastPart))

			tempWindow, _ = app.NewWindow("./ui/showData.html", &astilectron.WindowOptions{
				Center:    astikit.BoolPtr(true),
				Height:    astikit.IntPtr(600),
				Width:     astikit.IntPtr(600),
				Resizable: astikit.BoolPtr(false)})
			tempWindow.Create()

			tempWindow.OnMessage(func(mes *astilectron.EventMessage) interface{} {
				var tempMess string
				mes.Unmarshal(&tempMess)
				switch tempMess {
				case "exit":
					file.Close()
					HTML.Close()
					tempWindow.Close()
				case "save":
					// toSave := make([]byte, sB.Len())
					// log.Println(toSave)
					// HTML.ReadAt(toSave, int64(len(firstPart)))
					// log.Println(string(toSave))

					// var newString string
					// replacer := strings.NewReplacer("<tr></td>", "|", "</td><td>", "|")
					// newString = replacer.Replace(string(toSave))
					// log.Println(newString)
					HTML.Close()
					file.Close()
				}
				return nil
			})
		case "exit":
			window.Close()
			os.Exit(0)
		}

		return nil
	})
}
