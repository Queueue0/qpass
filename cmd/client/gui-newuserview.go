package main

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Queueue0/qpass/internal/validator"
)

func (a *Application) NewUserView(w *app.Window) (bool, error) {
	var ops op.Ops
	var userName widget.Editor
	var password widget.Editor
	var confirmPassword widget.Editor
	var addBtn widget.Clickable
	var v validator.Validator

	created := false

	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if addBtn.Clicked(gtx) {
				un, pw, cpw := userName.Text(), password.Text(), confirmPassword.Text()

				// Validate data
				v = validator.Validator{}
				v.CheckField(validator.NotBlank(un), "username", "This field cannot be blank")
				v.CheckField(validator.NotBlank(pw), "password", "This field cannot be blank")
				v.CheckField(validator.NotBlank(cpw), "confirm", "This field cannot be blank")
				v.CheckField(validator.Matches(cpw, pw), "password", "Passwords don't match")

				if v.Valid() {
					id, err := a.UserModel.Insert(un, pw)
					if err != nil {
						return false, err
					}

					UUID, err := a.UserModel.IDtoUUID(id)
					if err != nil {
						return false, err
					}

					err = a.Logs.NewLastSync(UUID)
					if err != nil {
						return false, err
					}

					created = true
					w.Perform(system.ActionClose)
				}
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if v.Valid() {
							return layout.Dimensions{}
						}
						errTxt, in := v.FieldErrors["username"]
						if !in {
							return layout.Dimensions{}
						}

						txt := material.Body1(th, errTxt)
						txt.Color = color.NRGBA{R: 244, G: 67, B: 54, A: 255}

						margins := layout.UniformInset(unit.Dp(10))
						margins.Bottom = 0
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return txt.Layout(gtx)
							},
						)

					},
				),
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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if v.Valid() {
							return layout.Dimensions{}
						}
						errTxt, in := v.FieldErrors["password"]
						if !in {
							return layout.Dimensions{}
						}

						txt := material.Body1(th, errTxt)
						txt.Color = color.NRGBA{R: 244, G: 67, B: 54, A: 255}

						margins := layout.UniformInset(unit.Dp(10))
						margins.Bottom = 0
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return txt.Layout(gtx)
							},
						)

					},
				),
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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if v.Valid() {
							return layout.Dimensions{}
						}
						errTxt, in := v.FieldErrors["confirm"]
						if !in {
							return layout.Dimensions{}
						}

						txt := material.Body1(th, errTxt)
						txt.Color = color.NRGBA{R: 244, G: 67, B: 54, A: 255}

						margins := layout.UniformInset(unit.Dp(10))
						margins.Bottom = 0
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return txt.Layout(gtx)
							},
						)

					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						txt := material.Editor(th, &confirmPassword, "Confirm Password")
						confirmPassword.SingleLine = true
						confirmPassword.Mask = '*'
						confirmPassword.Submit = true

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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.UniformInset(unit.Dp(10))
						btn := material.Button(th, &addBtn, "+ Add")
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return btn.Layout(gtx)
						})
					},
				),
			)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return created, e.Err
		}
	}
}
