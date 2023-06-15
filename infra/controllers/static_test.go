package controllers

// import (
// 	"bytes"
// 	"fmt"
// 	"html/template"
// 	"io/ioutil"
// 	"lenslocked/templates"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// )

// func TestContact(t *testing.T) {
// 	var userController Users
// 	r := NewRouterTest(userController)
// 	ts := httptest.NewServer(r)
// 	defer ts.Close()
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contact", ts.URL), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	res, err := ts.Client().Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	buf := bytes.Buffer{}
// 	templ, err := template.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = templ.Execute(&buf, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if diff := cmp.Diff(string(data), buf.String()); diff != "" {
// 		t.Errorf("Create mismatch (-want +got):\n%v", diff)
// 	}
// }

// func TestFAQ(t *testing.T) {
// 	var userController Users
// 	r := NewRouter(userController)
// 	ts := httptest.NewServer(r)
// 	defer ts.Close()
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faq", ts.URL), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	res, err := ts.Client().Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	buf := bytes.Buffer{}
// 	templ, err := template.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	questions := []struct {
// 		Question string
// 		Answer   template.HTML
// 	}{
// 		{
// 			Question: "Is there a free version?",
// 			Answer:   "yes! We offer a free trial for 30 days on any paid plans.",
// 		},
// 		{
// 			Question: "What are your support hours?",
// 			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
// 		},
// 		{
// 			Question: "How do I contact support?",
// 			Answer:   `Email us <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
// 		},
// 	}
// 	err = templ.Execute(&buf, questions)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if diff := cmp.Diff(string(data), buf.String()); diff != "" {
// 		t.Errorf("Create mismatch (-want +got):\n%v", diff)
// 	}
// }

// func TestHome(t *testing.T) {
// 	var userController Users
// 	r := NewRouter(userController)
// 	ts := httptest.NewServer(r)
// 	defer ts.Close()
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/", ts.URL), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	res, err := ts.Client().Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	buf := bytes.Buffer{}
// 	templ, err := template.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = templ.Execute(&buf, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if diff := cmp.Diff(string(data), buf.String()); diff != "" {
// 		t.Errorf("Create mismatch (-want +got):\n%v", diff)
// 	}
// }
