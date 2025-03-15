package main

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/validator"
)

func (a *Application) EditView(w *app.Window, p *models.Password) error {
	var (
		ops         op.Ops
		serviceName widget.Editor
		userName    widget.Editor
		password    widget.Editor
		editBtn      widget.Clickable
		v           validator.Validator
	)

	serviceName.SetText(p.ServiceName)
	userName.SetText(p.Username)
	password.SetText(p.Password)

	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if editBtn.Clicked(gtx) {
				sn, un, pw := serviceName.Text(), userName.Text(), password.Text()

				// Validate data
				v = validator.Validator{}
				v.CheckField(validator.NotBlank(sn), "service", "Service Name cannot be blank")

				if v.Valid() {
					err := a.PasswordModel.Update(*a.ActiveUser, p.UUID.String(), sn, un, pw)
					if err != nil {
						return err
					}

					npw, err := a.PasswordModel.Get(p.ID, *a.ActiveUser)
					if err != nil {
						return err
					}

					p.ServiceName = npw.ServiceName
					p.Username = npw.Username
					p.Password = npw.Password
					return nil
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
						errTxt, in := v.FieldErrors["service"]
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
						txt := material.Editor(th, &serviceName, "Service Name")
						serviceName.SingleLine = true
						serviceName.Submit = true

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
						txt := material.Editor(th, &password, "Password")
						password.SingleLine = true
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
						margins := layout.UniformInset(unit.Dp(10))
						btn := material.Button(th, &editBtn, "Submit")
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return btn.Layout(gtx)
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
