package controllers

import (
	"github.com/andreweggleston/GoSeniorAssassin/controllers/admin"
	"html/template"
)

func InitTemplates() {
	admin.InitAdminTemplates()

	mainTempl = template.Must(template.ParseFiles("views/index.html"))
}