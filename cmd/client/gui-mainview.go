package main

import (
	"fmt"
	"io"
	"strings"

	"gio.tools/icons"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (a *Application) MainView(w *app.Window) error {
	pws := gpwList{}
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
				go func() {
					aw := new(app.Window)
					aw.Option(app.Title("New Password"))
					aw.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
					np, err := a.AddView(aw)
					if err != nil {
						fmt.Println(err.Error())
					}

					if np.ID > 0 {
						a.Passwords = append(a.Passwords, np)
						gp := newGPassword(np)
						pws = append(pws, &gp)
						pws.sort()
						w.Invalidate()
					}
				}()
			}

			for i := range pws {
				p := pws[i]
				if p.CopyBtn.Clicked(gtx) {
					gtx.Execute(clipboard.WriteCmd{Type: "application/text", Data: io.NopCloser(strings.NewReader(p.Password))})
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
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
													layout.Flexed(4,
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
													layout.Rigid(
														func(gtx layout.Context) layout.Dimensions {
															var bi *widget.Icon
															var bd string
															if p.Shown {
																bi = icons.ActionVisibilityOff
																bd = "Hide"
															} else {
																bi = icons.ActionVisibility
																bd = "Show"
															}
															btn := material.IconButton(th, p.ShowBtn, bi, bd)
															return btn.Layout(gtx)
														},
													),
													layout.Rigid(
														func(gtx layout.Context) layout.Dimensions {
															btn := material.IconButton(th, p.CopyBtn, icons.ContentContentCopy, "Copy Password")
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
