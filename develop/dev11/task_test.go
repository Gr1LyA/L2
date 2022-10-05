package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const post = `{
	"EventID": "2",
	"Date": "2022-10-10",
	"Description":"alo",
	"UserID": "2"
}`

func TestServer(t *testing.T) {
	s := newServer("localhost", "8081")

	// test create_event
	{
		rec := httptest.NewRecorder()

		r := strings.NewReader(post)
		req, err := http.NewRequest(http.MethodPost, "/create_event", r)
		if err != nil {
			t.Fatal(err)
		}

		s.server.Handler.ServeHTTP(rec, req)

		fmt.Println(rec.Body.String())
		if rec.Code != 200 {
			t.Error(rec.Code)
		}
	}

	//test create_event with already created event
	{
		rec := httptest.NewRecorder()

		r := strings.NewReader(post)
		req, err := http.NewRequest(http.MethodPost, "/create_event", r)
		if err != nil {
			t.Fatal(err)
		}

		s.server.Handler.ServeHTTP(rec, req)

		fmt.Println(rec.Body.String())
		if rec.Code != 500 {
			t.Error(rec.Code)
		}
	}

	//test events_for_day
	{
		rec := httptest.NewRecorder()

		r := strings.NewReader("")
		req, err := http.NewRequest(http.MethodGet, "/events_for_day?UserID=2&Date=2022-10-10", r)

		if err != nil {
			t.Fatal(err)
		}

		s.server.Handler.ServeHTTP(rec, req)

		fmt.Println(rec.Body.String())
		if rec.Code != 200 {
			t.Error(rec.Code)
		}
	}
	//  test delete_event
	{
		rec := httptest.NewRecorder()
		r := strings.NewReader(post)
		req, err := http.NewRequest(http.MethodPost, "/delete_event", r)

		if err != nil {
			t.Fatal(err)
		}

		s.server.Handler.ServeHTTP(rec, req)

		fmt.Println(rec.Body.String())
		if rec.Code != 200 {
			t.Error(rec.Code)
		}
	}
}
