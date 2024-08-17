package main

import (
	"database/sql"
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Queueue0/qpass/internal/models"
)

type Application struct {
	UserModel  *models.UserModel
	ActiveUser *models.User
}

func main() {
	qpassHome, err := getQpassHome()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("file:%s/pwdb.sqlite?mode=rwc", qpassHome)
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = initializeDB(db)
	if err != nil {
		log.Fatal(err)
	}

	um := models.UserModel{
		DB: db,
	}

	a := Application{
		UserModel: &um,
	}

	go func() {
		lw := new(app.Window)
		lw.Option(app.Title("Login"))
		lw.Option(app.Size(unit.Dp(500), unit.Dp(200)))
		if err := a.login(lw); err != nil {
			log.Fatal(err)
		}

		if a.ActiveUser != nil {
			w := new(app.Window)
			w.Option(app.Title("QPass"))
			w.Option(app.Size(unit.Dp(1280), unit.Dp(720)))
			if err := a.draw(w); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}()

	app.Main()
}

func getQpassHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	qpassHome := fmt.Sprintf("%s/.qpass", homeDir)
	if _, err := os.Stat(qpassHome); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(qpassHome, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return qpassHome, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func initializeDB(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, salt TEXT)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, userId INT, service TEXT, username TEXT, password TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) draw(w *app.Window) error {
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

func (a *Application) login(w *app.Window) error {
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

						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}

						border := widget.Border{
							Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
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
							Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
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
