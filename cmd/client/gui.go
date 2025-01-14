package main

import (
	"fmt"
	"image/color"
	"io"
	"sort"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Queueue0/qpass/internal/models"
	"github.com/Queueue0/qpass/internal/validator"
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

func (p *gPassword) toggleShow() {
	p.Shown = !p.Shown
}

type gpwList []*gPassword

func (gl gpwList) sort() {
	sort.Slice(gl, func(i, j int) bool {
		return gl[i].ServiceName < gl[j].ServiceName
	})
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
					np, err := a.addView(aw)
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
	var newUserBtn widget.Clickable
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

	var addUser = func() {
		go func() {
			aw := new(app.Window)
			aw.Option(app.Title("New User"))
			aw.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			created, err := a.newUserView(aw)
			if err != nil {
				fmt.Println(err.Error())
			}
			
			if created {
				errorTxt = "Successfully added new user! Please log in."
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

			if newUserBtn.Clicked(gtx) {
				addUser()
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
				// Buttons
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbtn := material.Button(th, &loginBtn, "Log In")
						nubtn := material.Button(th, &newUserBtn, "New User")

						margins := layout.UniformInset(unit.Dp(10))

						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceAround,
						}.Layout(gtx,
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return lbtn.Layout(gtx)
										},
									)
								},
							),
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return margins.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return nubtn.Layout(gtx)
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

func (a *Application) addView(w *app.Window) (models.Password, error) {
	var ops op.Ops
	var serviceName widget.Editor
	var userName widget.Editor
	var password widget.Editor
	var addBtn widget.Clickable
	var v validator.Validator

	var np models.Password

	th := material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if addBtn.Clicked(gtx) {
				sn, un, pw := serviceName.Text(), userName.Text(), password.Text()

				// Validate data
				v = validator.Validator{}
				v.CheckField(validator.NotBlank(sn), "service", "Service Name cannot be blank")

				if v.Valid() {
					id, err := a.PasswordModel.Insert(*a.ActiveUser, sn, un, pw)
					if err != nil {
						return models.Password{}, err
					}
					np = models.Password{
						ID:          id,
						UserID:      a.ActiveUser.ID,
						ServiceName: sn,
						Username:    un,
						Password:    pw,
					}

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
						btn := material.Button(th, &addBtn, "+ Add")
						return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return btn.Layout(gtx)
						})
					},
				),
			)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return np, e.Err
		}
	}
}

func (a *Application) newUserView(w *app.Window) (bool, error) {
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
					_, err := a.UserModel.Insert(un, pw)
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
