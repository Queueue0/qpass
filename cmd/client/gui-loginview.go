package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (a *Application) LoginView(w *app.Window) error {
	var (
		ops        op.Ops
		username   widget.Editor
		password   widget.Editor
		loginBtn   widget.Clickable
		syncBtn    widget.Clickable
		newUserBtn widget.Clickable
		optionsBtn widget.Clickable
		errorTxt   string
	)

	th := material.NewTheme()

	var login = func() {
		u, err := a.UserModel.Authenticate(username.Text(), password.Text())
		if err != nil {
			errorTxt = err.Error()
		} else {
			a.ActiveUser = &u
			w.Perform(system.ActionClose)
		}
	}

	var sync = func() {
		err := a.loginSync(username.Text(), password.Text())
		if err != nil {
			errorTxt = err.Error()
		} else {
			w.Perform(system.ActionClose)
		}
	}

	var addUser = func() {
		go func() {
			aw := new(app.Window)
			aw.Option(app.Title("New User"))
			aw.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			created, err := a.NewUserView(aw)
			if err != nil {
				fmt.Println(err.Error())
			}

			if created {
				errorTxt = "Successfully added new user! Please log in."
				w.Invalidate()
			}
		}()
	}

	var options = func() {
		go func() {
			ow := new(app.Window)
			ow.Option(app.Title("Settings"))
			ow.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			err := a.OptionView(ow)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if loginBtn.Clicked(gtx) {
				login()
			}

			if syncBtn.Clicked(gtx) {
				sync()
			}

			if newUserBtn.Clicked(gtx) {
				addUser()
			}

			if optionsBtn.Clicked(gtx) {
				options()
			}

			we, ok := username.Update(gtx)
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
						txt := material.Editor(th, &username, "Username")
						username.SingleLine = true
						username.Submit = true

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
				// Buttons
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbtn := material.Button(th, &loginBtn, "Log In")
						sbtn := material.Button(th, &syncBtn, "Sync")
						nubtn := material.Button(th, &newUserBtn, "New User")
						obtn := material.Button(th, &optionsBtn, "⚙️")

						margins := layout.UniformInset(unit.Dp(5))

						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceSides,
						}.Layout(gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return lbtn.Layout(gtx)
										},
									)
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return sbtn.Layout(gtx)
										},
									)
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return nubtn.Layout(gtx)
										},
									)
								},
							),
							layout.Rigid(
								func (gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func (gtx layout.Context) layout.Dimensions {
											return obtn.Layout(gtx)
										},
									)
								},
							),
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
