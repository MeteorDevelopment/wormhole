package email

import (
	"fmt"
	"github.com/gobuffalo/flect"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"wormhole/pkg/config"
)

type Template struct {
	Subject   string
	Content   string
	Formatter *template.Template
}

var (
	baseTemplate *template.Template

	ConfirmRegister       *Template
	ConfirmChangePassword *Template
	ConfirmChangeEmail    *Template
	ConfirmDeleteAccount  *Template
)

func initTemplates() {
	initBaseTemplate()

	ConfirmRegister = NewTemplate("confirm-email", ConfirmLink("email", "registerConfirm", "", false))
	ConfirmChangePassword = NewTemplate("confirm-password-change", ConfirmLink("change of password", "change/passwordConfirm", "", true))
	ConfirmChangeEmail = NewTemplate("confirm-email-change", ConfirmLink("change of email", "change/emailConfirm", "Your email address will be changed to: {{.NewEmail}}.", true))
	ConfirmDeleteAccount = NewTemplate("confirm-account-deletion", ConfirmLink("account deletion", "deleteConfirm", "Your account will be permanently deleted.", true))
}

func NewTemplate(name, content string) *Template {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		log.Fatalf("Failed to parse email template: %v", err)
	}

	return &Template{
		Subject:   flect.Titleize(name),
		Content:   content,
		Formatter: tmpl,
	}
}

func ConfirmLink(action, route, info string, warn bool) string {
	str := fmt.Sprintf("Please use the link below to confirm your %s.\nThis link will be valid for %d minutes.\n", action, config.Get().JwtExpirationTime)
	if info != "" {
		str += info + "\n"
	}
	if warn {
		str += "If you did not request this, change your password immediately. Do not click this link!\n"
	}
	str += fmt.Sprintf("\n{{.ApiBase}}/auth/%s?token={{.Token}}", route)
	return str
}

func initBaseTemplate() {
	dir := "templates"
	fileName := "base.txt"

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// find base.txt
	var baseFile os.DirEntry
	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName {
			baseFile = file
			break
		}
	}

	if baseFile == nil {
		log.Fatal(errors.New("base.txt not found"))
	}

	// parse base.txt
	content, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("base").Parse(string(content))
	if err != nil {
		log.Fatal(err)
	}

	baseTemplate = tmpl
}
