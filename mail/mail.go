package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/echlebek/erickson/assets"
)

// Nil is a Mailer that does nothing.
var Nil = new(nilMailer)

const headers = `Content-Type: text/html; charset=ISO-8859-1

`

const reviewPostedHeader = "To: %s\r\nSubject: Code review request from %s\r\n\r\n"
const reviewAnnotatedHeader = "Content-Type: text/html; charset=ISO-8859-1\r\nTo: %s\r\nSubject: %s reviewed your code.\r\n\r\n"

var reviewPostedBody = template.Must(template.New("review-posted").Parse(`
Hi {{ .Recipient }},

{{ .Sender }} is requesting that you view their commits to the {{ .Repository }} repository.

<a href="{{ .ReviewURL }}">Click here</a> to see the code review.
`))

// Mailer sends e-mail. It is purpose-built for sending specific messages.
type Mailer interface {
	NotifyReviewPosted(Message) error
	NotifyReviewAnnotated(Message) error
}

// NewMailer returns a mailer that will send mail from broker@server.
// If server == "", then Mailer will be Nil.
func NewMailer(server, broker string, auth smtp.Auth) Mailer {
	if server == "" {
		return Nil
	}
	return mailer{Server: server, Auth: auth, Broker: broker}
}

type nilMailer struct{}

func (n *nilMailer) NotifyReviewPosted(Message) error {
	return nil
}

func (n *nilMailer) NotifyReviewAnnotated(m Message) error {
	return assets.Templates["mail_comments.html"].Execute(ioutil.Discard, m)
}

type mailer struct {
	Server string
	Auth   smtp.Auth
	Broker string
}

// Annotation is a simplified representation of a review.Annotation
type Annotation struct {
	File       string
	LineNumber int
	LHS        string
	RHS        string
	Comment    string
}

// Returns the number of lines in the Comment.
func (a Annotation) CommentLines() int {
	s := strings.Split(a.Comment, "\n")
	return len(s)
}

// Message represents a message sent by the server to Recipient on behalf of Sender.
type Message struct {
	Sender     string
	Recipient  string
	Repository string
	ReviewURL  string

	Annotations []Annotation
}

func WriteMessage(w http.ResponseWriter, m Message, tmpl assets.Template) {
	tmpl.Execute(w, m)
}

func (c mailer) NotifyReviewPosted(m Message) error {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, headers)
	fmt.Fprintf(buf, reviewPostedHeader, m.Recipient, m.Sender)

	if err := reviewPostedBody.Execute(buf, m); err != nil {
		return fmt.Errorf("couldn't send mail to %s: %s", m.Recipient, err)
	}
	to := []string{m.Recipient}
	if err := smtp.SendMail(c.Server+":25", c.Auth, c.Broker, to, buf.Bytes()); err != nil {
		return fmt.Errorf("couldn't send mail to %s: %s", m.Recipient, err)
	}
	return nil
}

func (c mailer) NotifyReviewAnnotated(m Message) error {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, reviewAnnotatedHeader, m.Recipient, m.Sender)

	if err := assets.Templates["mail_comments.html"].Execute(buf, m); err != nil {
		return fmt.Errorf("couldn't send mail to %s: %s", m.Recipient, err)
	}
	to := []string{m.Recipient}
	if err := smtp.SendMail(c.Server+":25", c.Auth, c.Broker, to, buf.Bytes()); err != nil {
		return fmt.Errorf("couldn't send mail to %s: %s", m.Recipient, err)
	}
	return nil
}
