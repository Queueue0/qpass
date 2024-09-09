package main

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var borderColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
var inputPadding = unit.Dp(10)

type gPassword struct {
	ServiceName string
	Username    string
	Password    string
	Shown       bool
	ShowBtn     *widget.Clickable
	CopyBtn     *widget.Clickable
}

func (p *gPassword) copy() {
	c := clipboard.New()
	err := c.CopyText(p.Password)
	if err != nil {
		panic(err)
	}
}

func (p *gPassword) toggleShow() {
	p.Shown = !p.Shown
}

func newGPassword(p models.Password) gPassword {
	return gPassword{
		ServiceName: p.ServiceName,
		Username:    p.Username,
		Password:    p.Password,
		Shown:       false,
		ShowBtn:     &widget.Clickable{},
		CopyBtn:     &widget.Clickable{},
	}
}

func (a *Application) mainView(w *app.Window) error {
	pws := []*gPassword{}
	for _, p := range a.Passwords {
		gp := newGPassword(p)
		pws = append(pws, &gp)
	}

	var ops op.Ops
	var addBtn widget.Clickable
	var pwlist widget.List
	th := material.NewTheme()

	pwlist.List.Axis = layout.Vertical

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if addBtn.Clicked(gtx) {
				fmt.Println("Add btn clicked")
			}

			for i := range pws {
				p := pws[i]
				if p.CopyBtn.Clicked(gtx) {
					p.copy()
				}

				if p.ShowBtn.Clicked(gtx) {
					p.toggleShow()
				}
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				// Button
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						inset := layout.UniformInset(unit.Dp(10))
						btn := material.Button(th, &addBtn, "+ Add New")
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions{
							return btn.Layout(gtx)
						})
					},
				),
				// Header
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						inset := layout.UniformInset(unit.Dp(2))
						inset.Right = unit.Dp(24)
						text := material.Body1(th, "")
						text.MaxLines = 1

						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceEnd,
							}.Layout(gtx,
								layout.Flexed(1,
									func(gtx layout.Context) layout.Dimensions {
										text.Text = ""
										text.Font.Weight = font.Bold
										text.Text = "Service Name"
										return text.Layout(gtx)
									},
								),
								layout.Flexed(1,
									func(gtx layout.Context) layout.Dimensions {
										text.Text = ""
										text.Font.Weight = font.Bold
										text.Text = "Username"
										return text.Layout(gtx)
									},
								),
								layout.Flexed(1,
									func(gtx layout.Context) layout.Dimensions {
										text.Text = ""
										text.Font.Weight = font.Bold
										text.Text = " Password"
										return text.Layout(gtx)
									},
								),
							)
						})
					},
				),
				// Password list
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						// var grid outlay.Grid
						text := material.Body1(th, "")
						text.MaxLines = 1

						inset := layout.UniformInset(unit.Dp(2))
						// dims := inset.Layout(gtx, text.Layout)

						list := material.List(th, &pwlist)
						return list.Layout(gtx, len(pws),
							// Rows
							func(gtx layout.Context, i int) layout.Dimensions {
								p := pws[i]
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Axis:    layout.Horizontal,
										Spacing: layout.SpaceEnd,
									}.Layout(gtx,
										layout.Flexed(1,
											func(gtx layout.Context) layout.Dimensions {
												text.Text = ""
												text.Font.Weight = font.Normal
												text.Text = p.ServiceName
												return text.Layout(gtx)
											},
										),
										layout.Flexed(1,
											func(gtx layout.Context) layout.Dimensions {
												text.Text = ""
												text.Font.Weight = font.Normal
												text.Text = p.Username
												return text.Layout(gtx)
											},
										),
										layout.Flexed(1,
											func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{
													Axis:    layout.Horizontal,
													Spacing: layout.SpaceEnd,
												}.Layout(gtx,
													layout.Flexed(2,
														func(gtx layout.Context) layout.Dimensions {
															text.Text = ""
															text.Font.Weight = font.Normal
															if p.Shown {
																text.Text = p.Password
															} else {
																text.Text = strings.Repeat("*", len(p.Password))
															}
															return text.Layout(gtx)
														},
													),
													layout.Flexed(1,
														func(gtx layout.Context) layout.Dimensions {
															var bt string
															if p.Shown {
																bt = "Hide"
															} else {
																bt = "Show"
															}
															btn := material.Button(th, p.ShowBtn, bt)
															return btn.Layout(gtx)
														},
													),
													layout.Flexed(1,
														func(gtx layout.Context) layout.Dimensions {
															btn := material.Button(th, p.CopyBtn, "Copy")
															return btn.Layout(gtx)
														},
													),
												)
											},
										),
									)
								})
							})
					}))
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
