package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World!")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter password...")

	submitBtn := widget.NewButton("Submit", func() {
		fmt.Println(send(input.Text))
	})

	content := container.NewVBox(input, submitBtn)

	w.SetContent(content)
	w.ShowAndRun()
}

func send(s string) string {
	conf := tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := tls.Dial("tcp", "localhost:1717", &conf)
	if err != nil {
		log.Fatal(err)
	}

	c.Write([]byte(s))

	buffer := make([]byte, 1024)
	_, err = c.Read(buffer)

	return string(buffer[:])
}
