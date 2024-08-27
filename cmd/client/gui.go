package main

import (
	"fmt"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var borderColor = color.NRGBA{R: 50, G: 50, B: 50, A: 255}

func (a *Application) mainView(w *app.Window) error {
	var ops op.Ops
	var serviceName widget.Editor
	var saveButton widget.Clickable
	var res string = a.ActiveUser.Username
	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			//			if saveButton.Clicked(gtx) {
			//				bres, err := crypto.GetKey(serviceName.Text())
			//				if err != nil {
			//					log.Fatal(err)
			//				}
			//				res = hex.EncodeToString(bres)
			//			}

			layout.Flex{
				Axis:    layout.Vertical,
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
						lbl := material.Label(th, unit.Sp(25), fmt.Sprintf("Logged In As: %s", res))
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

func (a *Application) loginView(w *app.Window) error {
	var ops op.Ops
	var userName widget.Editor
	var password widget.Editor
	var loginBtn widget.Clickable
	var errorTxt string
	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if loginBtn.Clicked(gtx) {
				u, err := a.UserModel.Authenticate(userName.Text(), password.Text())
				if err != nil {
					errorTxt = err.Error()
				} else {
					a.ActiveUser = &u
					// TODO: handle this error
					a.Passwords, _ = a.PasswordModel.GetAllForUser(*a.ActiveUser)
					w.Perform(system.ActionClose)
				}
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				// Error Text
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), errorTxt)
						return lbl.Layout(gtx)
					},
				),
				// Username
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						txt := material.Editor(th, &userName, "Username")
						userName.SingleLine = true
						userName.LineHeight = unit.Sp(20)
						userName.LineHeightScale = 1

						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, txt.Layout)
							},
						)
					},
				),
				// Password
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						txt := material.Editor(th, &password, "Password")
						password.SingleLine = true
						password.Mask = '*'

						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, txt.Layout)
							},
						)
					},
				),
				// Button
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &loginBtn, "Log In")

						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return btn.Layout(gtx)
							},
						)
					},
				),
			)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}
