package swu

import (
	"os"
	"testing"
)

func TestNewSWU(t *testing.T) {
	api := New("key")
	if api == nil {
		t.Error("New should not return nil")
	}
}

func TestTemplates(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.Emails()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTemplate(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.GetTemplate("tem_bRKXvNLAXTG8EGxhut3gCe")
	if err != nil {
		t.Error(err)
	}
}

func TestGetTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.GetTemplateVersion("tem_bRKXvNLAXTG8EGxhut3gCe", "ver_Hh35dZhnffghidEy6VeHKL")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.UpdateTemplateVersion("tem_bRKXvNLAXTG8EGxhut3gCe", "ver_Hh35dZhnffghidEy6VeHKL",
		&Version{
			Name:    "Test",
			Subject: "Test",
			Text:    "test",
		})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTemplate(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.CreateTemplate(&Version{
		Name:    "test",
		Subject: "test",
		Text:    "ALOHA",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	_, err := api.CreateTemplateVersion("tem_nXAPFGXQXFKcibJHdm9PZ9", &Version{
		Name:    "test",
		Subject: "test",
		Text:    "ALOHA1",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestSend(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	email := &Email{
		ID: "tem_bRKXvNLAXTG8EGxhut3gCe",
		Recipient: &Recipient{
			Address: "yamil@sendgrid.com",
		},
		EmailData: make(map[string]string),
	}
	err := api.Send(email)
	if err != nil {
		t.Error(err)
	}
}

func TestRender(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	req := &RenderRequest{
		Template: "tem_Jx7kFb39gk7kTVd78Xpy8BWR",
		TemplateData: make(map[string]string),
	}
	_, err := api.Render(req)
	if err != nil {
		t.Error(err)
	}
}
