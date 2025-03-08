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

func (a *Application) OptionView(w *app.Window) error {
	var (
		ops       op.Ops
		v         validator.Validator
		addressEd widget.Editor
		portEd    widget.Editor
		saveBtn   widget.Clickable
		cancelBtn widget.Clickable

		th *material.Theme = material.NewTheme()
	)

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if saveBtn.Clicked(gtx) {
				// Validate input, save, and reload config
			}

			if cancelBtn.Clicked(gtx) {
				w.Perform(system.ActionClose)
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
						lbl := material.Body1(th, "Sync Server Address")
						txt := material.Editor(th, &addressEd, "127.0.0.1")
						addressEd.SingleLine = true
						addressEd.Submit = true

						margins := layout.UniformInset(unit.Dp(10))
						lblMargins := layout.UniformInset(unit.Dp(10))
						lblMargins.Right = unit.Dp(50)
						padding := layout.UniformInset(inputPadding)

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return lblMargins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return padding.Layout(gtx, lbl.Layout)
									})
								},
							),
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return border.Layout(gtx,
											func(gtx layout.Context) layout.Dimensions {
												return padding.Layout(gtx, txt.Layout)
											},
										)
									})
								},
							),
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(th, "Sync Server Port")
						txt := material.Editor(th, &portEd, "10448")
						portEd.SingleLine = true
						portEd.Submit = true

						margins := layout.UniformInset(unit.Dp(10))
						lblMargins := layout.UniformInset(unit.Dp(10))
						lblMargins.Right = unit.Dp(50)
						padding := layout.UniformInset(inputPadding)

						border := widget.Border{
							Color:        borderColor,
							CornerRadius: unit.Dp(1),
							Width:        unit.Dp(2),
						}

						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return lblMargins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return padding.Layout(gtx, lbl.Layout)
									})
								},
							),
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return border.Layout(gtx,
											func(gtx layout.Context) layout.Dimensions {
												return padding.Layout(gtx, txt.Layout)
											},
										)
									})
								},
							),
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceSides,
						}.Layout(gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									margins := layout.UniformInset(unit.Dp(10))
									btn := material.Button(th, &saveBtn, "Save")
									return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return btn.Layout(gtx)
									})
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									margins := layout.UniformInset(unit.Dp(10))
									btn := material.Button(th, &cancelBtn, "Cancel")
									return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return btn.Layout(gtx)
									})
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
