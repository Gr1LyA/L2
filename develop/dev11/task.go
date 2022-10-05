package main

/*
структура проекта:
	есть три основных объекта, у каждого из них своя зона ответсвенности
		storage - потокобезопасное хранилище:
			хранит map[string]map[string]event и мъютекс
			map[string]map[string]event - ключ основного множества - user id, ключ подмножества - event id
		calendar - объект со встроенным со встроенным storage, имеет методы возвращающие ивенты для дня, недели, месяца
		server - отвечает за серверную часть


	customDate - структура для кастомного маршалинга/анмаршалинга времени
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const layout = "2006-01-02"

func main() {
	s := newServer("localhost", "8081")
	err := s.server.ListenAndServe()
	fmt.Println(err)
}

// server - основной объект
type server struct {
	server *http.Server
	*calendar
}

// Конструктор server
func newServer(host string, port string) (s *server) {
	mux := http.NewServeMux()

	s = &server{
		server: &http.Server{
			Addr:    host + ":" + port,
			Handler: middleware(mux),
		},
		calendar: newCalendar(),
	}

	// роутинг с соответсвующими методами объекта server
	mux.HandleFunc("/events_for_day", s.eventsForDay)
	mux.HandleFunc("/events_for_week", s.eventsForWeek)
	mux.HandleFunc("/events_for_month", s.eventsForMonth)
	mux.HandleFunc("/create_event", s.eventCreate)
	mux.HandleFunc("/update_event", s.eventUpdate)
	mux.HandleFunc("/delete_event", s.eventDelete)

	return
}

// Функция - прослойка для логирования запросов
func middleware(mux http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method: " + r.Method + ", url: " + r.RequestURI)
		mux.ServeHTTP(w, r)
	}
}

func (s *server) eventsForDay(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodGet {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method get"))
		return
	}

	UserID, t, err := parseGet(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.getEventsForDay(UserID, t)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	if len(events) == 0 {
		responceResult(w, http.StatusOK, "no events for this day")
	} else {
		responceResult(w, http.StatusOK, events)
	}
}

func (s *server) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodGet {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method get"))
		return
	}

	UserID, t, err := parseGet(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.getEventsForWeek(UserID, t)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	if len(events) == 0 {
		responceResult(w, http.StatusOK, "no events for this week")
	} else {
		responceResult(w, http.StatusOK, events)
	}
}

func (s *server) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodGet {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method get"))
		return
	}

	UserID, t, err := parseGet(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.getEventsForMonth(UserID, t)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	if len(events) == 0 {
		responceResult(w, http.StatusOK, "no events for this month")
	} else {
		responceResult(w, http.StatusOK, events)
	}
}

func parseGet(r *http.Request) (UserID string, t time.Time, err error) {
	if r.FormValue("UserID") == "" || r.FormValue("Date") == "" {
		return UserID, t, errors.New("expected UserID and Date in query string")
	}

	UserID = r.FormValue("UserID")

	t, err = time.Parse(layout, r.FormValue("Date"))
	return
}

func (s *server) eventCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method post"))
		return
	}

	e, err := parsePost(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	err = s.store(e.UserID, *e)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	responceResult(w, http.StatusOK, "event succesfuly add")
}

func (s *server) eventUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method post"))
		return
	}

	e, err := parsePost(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	err = s.update(e.UserID, *e)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	responceResult(w, http.StatusOK, "event succesfuly update")
}

func (s *server) eventDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		responceError(w, http.StatusServiceUnavailable, errors.New("expected method post"))
		return
	}

	e, err := parsePost(r)
	if err != nil {
		responceError(w, http.StatusBadRequest, err)
		return
	}

	err = s.delete(e.UserID, e.EventID)
	if err != nil {
		responceError(w, http.StatusInternalServerError, err)
		return
	}

	responceResult(w, http.StatusOK, "event succesfuly delete")
}

func parsePost(r *http.Request) (*event, error) {
	e := &event{}

	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

type calendar struct {
	*storage
}

func newCalendar() *calendar {
	return &calendar{newStorage()}
}

func (c *calendar) getEventsForDay(UserID string, Date time.Time) (result []event, err error) {
	events, ok := c.load(UserID)
	if !ok {
		return nil, errors.New("user id not exist")
	}
	for _, v := range events {
		if v.Date.Year() == Date.Year() && v.Date.Month() == Date.Month() && v.Date.Day() == Date.Day() {
			result = append(result, v)
		}
	}
	return
}

func (c *calendar) getEventsForWeek(UserID string, Date time.Time) (result []event, err error) {
	events, ok := c.load(UserID)
	if !ok {
		return nil, errors.New("user id not exist")
	}
	year, week := Date.ISOWeek()
	for _, v := range events {
		yearEvent, weekEvent := v.Date.ISOWeek()
		if year == yearEvent && week == weekEvent {
			result = append(result, v)
		}
	}
	return
}

func (c *calendar) getEventsForMonth(UserID string, Date time.Time) (result []event, err error) {
	events, ok := c.load(UserID)
	if !ok {
		return nil, errors.New("user id not exist")
	}
	for _, v := range events {
		if v.Date.Year() == Date.Year() && v.Date.Month() == Date.Month() {
			result = append(result, v)
		}
	}
	return
}

type storage struct {
	// Ключ user_id, значение - массив его ивентов
	m map[string]map[string]event
	sync.RWMutex
}

func newStorage() *storage {
	return &storage{m: make(map[string]map[string]event)}
}

func (c *storage) load(key string) (v map[string]event, ok bool) {
	c.RLock()
	v, ok = c.m[key]
	c.RUnlock()
	return
}

func (c *storage) store(key string, value event) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[key]; !ok {
		c.m[key] = make(map[string]event)
	}
	if _, ok := c.m[key][value.EventID]; !ok {
		c.m[key][value.EventID] = value
	} else {
		return errors.New("event already created")
	}
	return nil
}

func (c *storage) update(key string, value event) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[key]; !ok {
		return errors.New("user id not exist")
	}
	if _, ok := c.m[key][value.EventID]; !ok {
		return errors.New("event not exist")
	}
	c.m[key][value.EventID] = value
	return nil
}

func (c *storage) delete(UserID string, EventID string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[UserID]; !ok {
		return errors.New("user id not exist")
	}
	if _, ok := c.m[UserID][EventID]; !ok {
		return errors.New("event not exist")
	}
	delete(c.m[UserID], EventID)
	if len(c.m[UserID]) == 0 {
		delete(c.m, UserID)
	}
	return nil
}

type event struct {
	EventID     string     `json:"EventID"`
	Date        customDate `json:"Date"`
	Description string     `json:"Description"`
	UserID      string     `json:"UserID"`
}

type customDate struct {
	time.Time
}

func (c *customDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}

	c.Time, err = time.Parse(layout, s)
	return
}

func (c *customDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(layout))), nil
}

func responceResult(w http.ResponseWriter, code int, data any) {
	e := struct {
		Msg any `json:"result"`
	}{Msg: data}

	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
		w.Write(b)
	}
}

func responceError(w http.ResponseWriter, code int, err error) {
	e := struct {
		Msg any `json:"error"`
	}{Msg: err.Error()}

	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
		w.Write(b)
	}
}
