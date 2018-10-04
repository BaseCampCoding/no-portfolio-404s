package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/sendgrid/sendgrid-go"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Student struct {
	name         string
	email        string
	portfolioURL string
}

var students = [...]Student{
	Student{
		name:         "Cole Anderson",
		email:        "canderson@basecampcodingacademy.org",
		portfolioURL: "https://coleanderson7278.github.io/",
	},
	Student{
		name:         "Andrew Wheeler",
		email:        "awheeler@basecampcodingacademy.org",
		portfolioURL: "http://andrewwheeler315.github.io/",
	},
	Student{
		name:         "Cody van der Poel",
		email:        "cvanderpoel@basecampcodingacademy.org",
		portfolioURL: "https://codyvanderpoel.github.io/",
	},
	Student{
		name:         "Ginger Keys",
		email:        "gkeys@basecampcodingacademy.org",
		portfolioURL: "https://ginggk.github.io/",
	},
	Student{
		name:         "Henry Moore",
		email:        "hmoore@basecampcodingacademy.org",
		portfolioURL: "https://henrymoore13.github.io/",
	},
	Student{
		name:         "Ray Turner",
		email:        "rturner@basecampcodingacademy.org",
		portfolioURL: "https://rayturner677.github.io/",
	},
	Student{
		name:         "John Morgan",
		email:        "jmorgan@basecampcodingacademy.org",
		portfolioURL: "https://johnmorgan2000.github.io/",
	},
	Student{
		name:         "Timothy Bowling",
		email:        "tbowling@basecampcodingacademy.org",
		portfolioURL: "https://timothyb9526.github.io/",
	},
	Student{
		name:         "Myeisha Madkins",
		email:        "mmadkins@basecampcodingacademy.org",
		portfolioURL: "https://myeishamadkins.github.io",
	},
	Student{
		name:         "Matt Lipsey",
		email:        "mlipsey@basecampcodingacademy.org",
		portfolioURL: "https://matt2tech.github.io/",
	},
	Student{
		name:         "Danny Peterson",
		email:        "dpeterson@basecampcodingacademy.org",
		portfolioURL: "https://dannyp123.github.io/",
	},
	Student{
		name:         "Justice Taylor",
		email:        "jtaylor@basecampcodingacademy.org",
		portfolioURL: "https://jtaylor99.github.io/",
	},
	Student{
		name:         "Jakylan Standifer",
		email:        "jstandifer@basecampcodingacademy.org",
		portfolioURL: "https://jakylan.github.io/",
	},
	Student{
		name:         "Logan Harrell",
		email:        "lharrell@basecampcodingacademy.org",
		portfolioURL: "https://laharrell20xx.github.io/",
	},
	Student{
		name:         "Irma Patton",
		email:        "ipatton@basecampcodingacademy.org",
		portfolioURL: "https://irmapatton.github.io/",
	},
}

var sendGridClient = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
var fromEmail = mail.NewEmail("No-Portfolio-404-bot", "bots@basecampcodingacademy.org")
var nateEmail = mail.NewEmail("Nate", "nate@basecampcodingacademy.org")

func main() {
	var wg sync.WaitGroup
	for _, student := range students {
		wg.Add(1)
		go func(student Student) {
			defer wg.Done()
			response, err := http.Get(student.portfolioURL)
			if err != nil {
				notifyNateOfErr(student, err)
			} else if response.StatusCode != 200 {
				notifyNateAndStudentOfNon200(student, response)
			} else {
				logSuccess(student)
			}
		}(student)
	}
	wg.Wait()
}

func notifyNateOfErr(student Student, err error) {
	log.Printf("Error: %s - %s\n", student.name, err)
	sendEmail([]*mail.Email{fromEmail}, "Portfolio-Checker: Error", fmt.Sprintf("%v", err))
}

func notifyNateAndStudentOfNon200(student Student, response *http.Response) {
	log.Printf("Non-200: %s - %d\n", student.name, response.StatusCode)
	sendEmail([]*mail.Email{nateEmail, mail.NewEmail(student.name, student.email)}, "Portfolio-Checker: Error Accessing Portfolio", "Unable to access portfolio at "+student.portfolioURL)
}

func logSuccess(student Student) {
	log.Printf("Success: %s", student.name)
}

func sendEmail(toEmails []*mail.Email, subject, body string) {
	m := new(mail.SGMailV3)
	m.SetFrom(fromEmail)
	m.Subject = subject
	p := mail.NewPersonalization()
	p.AddTos(toEmails...)
	m.AddPersonalizations(p)
	m.AddContent(mail.NewContent("text/plain", body))
	sendGridClient.Send(m)
}
