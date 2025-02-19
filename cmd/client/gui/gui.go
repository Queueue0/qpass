package gui

import (
	"fmt"
	"image/color"
	"sort"

	"gioui.org/app"
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

func LoginView(w *app.Window, um *models.UserModel, pm *models.PasswordModel, lm *models.LogModel, au *models.User, pwl models.PasswordList) error {
	var ops op.Ops
	var userName widget.Editor
	var password widget.Editor
	var loginBtn widget.Clickable
	var newUserBtn widget.Clickable
	var errorTxt string
	th := material.NewTheme()

	var login = func() {
		u, err := um.Authenticate(userName.Text(), password.Text())
		if err != nil {
			errorTxt = err.Error()
		} else {
			au = &u
			// TODO: handle this error
			pwl, _ = pm.GetAllForUser(*au)
			w.Perform(system.ActionClose)
		}
	}

	var addUser = func() {
		go func() {
			aw := new(app.Window)
			aw.Option(app.Title("New User"))
			aw.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			created, err := NewUserView(aw, um, lm)
			if err != nil {
				fmt.Println(err.Error())
			}

			if created {
				errorTxt = "Successfully added new user! Please log in."
				w.Invalidate()
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

func AddView(w *app.Window, pm *models.PasswordModel, au models.User) (models.Password, error) {
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
					id, err := pm.Insert(au, sn, un, pw)
					if err != nil {
						return models.Password{}, err
					}
					np = models.Password{
						ID:          id,
						UserID:      au.ID,
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

func NewUserView(w *app.Window, um *models.UserModel, lm *models.LogModel) (bool, error) {
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
					id, err := um.Insert(un, pw)
					if err != nil {
						return false, err
					}

					UUID, err := um.IDtoUUID(id)
					if err != nil {
						return false, err
					}

					err = lm.NewLastSync(UUID)
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
