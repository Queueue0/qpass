package main

import (
	"gioui.org/app"
	"github.com/Queueue0/qpass/cmd/client/gui"
)

func (a *Application) newUserView(w *app.Window) (bool, error) {
	return gui.NewUserView(w, a.UserModel, a.Logs)
}


func (a *Application) loginView(w *app.Window) (error) {
	return gui.LoginView(w, a.UserModel, a.PasswordModel, a.Logs, a.ActiveUser, a.Passwords)
}

func (a *Application) mainView(w *app.Window) error {
	return gui.MainView(w, a.Passwords, a.PasswordModel, *a.ActiveUser)
}
