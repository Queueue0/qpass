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

var borderColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
var inputPadding = unit.Dp(10)

func (a *Application) mainView(w *app.Window) error {
	var ops op.Ops
	var serviceName widget.Editor
	var saveButton widget.Clickable
	var pwlist widget.List
	var res string = a.ActiveUser.Username
	th := material.NewTheme()

	pwlist.List.Axis = layout.Vertical

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
				// Password list
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						list := material.List(th, &pwlist)
						return list.Layout(gtx, len(a.Passwords), func(gtx layout.Context, i int) layout.Dimensions {
							p := a.Passwords[i]
							lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("%s\t%s\t%s", p.ServiceName, p.Username, p.Password))
							return lbl.Layout(gtx)
						})
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

	var login = func() {
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

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if loginBtn.Clicked(gtx) {
				login()
			}

			we, ok := userName.Update(gtx)
			if ok {
				switch we.(type) {
				case widget.SubmitEvent:
					login()
				}
			}

			we, ok = password.Update(gtx)
			if ok {
				switch we.(type) {
				case widget.SubmitEvent:
					login()
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
						userName.Submit = true

						margins := layout.UniformInset(unit.Dp(10))
						padding := layout.UniformInset(inputPadding)

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx,
									func(gtx layout.Context) layout.Dimensions {
										return padding.Layout(gtx, txt.Layout)
									},
								)
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
						password.Submit = true

						margins := layout.UniformInset(unit.Dp(10))
						padding := layout.UniformInset(inputPadding)

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx,
									func(gtx layout.Context) layout.Dimensions {
										return padding.Layout(gtx, txt.Layout)
									},
								)
							},
						)
					},
				),
				// Button
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &loginBtn, "Log In")

						margins := layout.UniformInset(unit.Dp(10))

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
