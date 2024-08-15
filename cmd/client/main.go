package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("QPass"))
		w.Option(app.Size(unit.Dp(1280), unit.Dp(700)))
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func draw(w *app.Window) error {
	var ops op.Ops
	var serviceName widget.Editor
	var saveButton widget.Clickable
	var res string
	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if saveButton.Clicked(gtx) {
				res = send(serviceName.Text())
			}

			layout.Flex{
				Axis: layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
			// Textbox
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					txt := material.Editor(th, &serviceName, "Text to send")
					return txt.Layout(gtx)
				},
			),
			// Button
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &saveButton, "Send Request")
					return btn.Layout(gtx)
				},
			),
			// Result Text
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(25), fmt.Sprintf("Result: %s", res))
					return lbl.Layout(gtx)
				},
			),
		)
		e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
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
