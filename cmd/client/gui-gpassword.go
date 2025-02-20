package main

import (
	"image/color"
	"sort"

	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/Queueue0/qpass/internal/models"
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
