package gae_mail

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"google.golang.org/appengine/mail"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/home.html")
		if err != nil {
			log.Errorf(c, "Failed to load the home template: %v", err)
			fmt.Println(err)
			return
		}
		t.Execute(w, map[string]string{"Title": "Upload and email file", "EmailAddress": "test@test.com"})
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Errorf(c, "Failed to retrieve the form file: %v", err)
			fmt.Println(err)
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf(c, "Failed read all the file: %v", err)
			fmt.Println(err)
		}

		addr := r.FormValue("emailaddress")

		msg := &mail.Message{
			Sender:  "Example.com Support <support@example.com>",
			To:      []string{addr},
			Subject: "Email by GAE",
			Body:    fmt.Sprintf("File %s attached", handler.Filename),
			Attachments: []mail.Attachment{
				{
					Name:      handler.Filename,
					Data:      data,
					ContentID: "ID_" + handler.Filename,
				},
			},
		}

		if err := mail.Send(c, msg); err != nil {
			log.Errorf(c, "Couldn't send email: %v", err)
		}

		t, err := template.ParseFiles("templates/email-result.html")
		if err != nil {
			log.Errorf(c, "Failed to load the email-result template: %v", err)
			fmt.Println(err)
			return
		}
		t.Execute(w, map[string]string{"Title": "Email result", "StatusCode": "200"})
	}
}
