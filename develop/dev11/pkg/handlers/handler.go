package handlers

import (
	"L2/develop/dev11/internal/event"
	"L2/develop/dev11/internal/setter"
	"L2/develop/dev11/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Handler struct {
	eventManager *event.EventManager
	validator    *validator.Validator
	logger       *logrus.Logger
	setter       setter.Setter
}

const formatDate string = "2006.02.01"

func NewHandler(v *validator.Validator, l *logrus.Logger) *Handler {
	ev := event.NewEventManager()
	s := setter.NewSetter()
	return &Handler{ev, v, l, s}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.POST("/create_event", h.createEvent())
	router.POST("/update_event", h.updateEvent())
	router.POST("/delete_event", h.deleteEvent())
	router.GET("/events_for_day", h.dayEvents())
	router.GET("/events_for_week", h.weekEvents())
	router.GET("/events_for_month", h.monthEvents())
}

func (h *Handler) debugEventManager() {
	for key, value := range *h.eventManager {
		fmt.Println("user id: ", key, " value", *value)
	}
	fmt.Println("---------------------------------------------------------")
}

func (h *Handler) postEvent(r *http.Request) (event.Event, error) {
	event := event.Event{}

	if isCorrect, err := h.validator.IsFormCorrect(r.Form); !isCorrect || err != nil {
		return event, errors.New("bad request")
	}
	if h.setter.SetFields(&event, r.Form) != nil {
		return event, errors.New("panic set")
	}

	return event, nil
}

func (h *Handler) postFormGet(param string, r *http.Request) (int, error) {
	if !r.Form.Has(param) {
		return 0, errors.New("bad request")
	}

	userID, err := strconv.Atoi(r.Form.Get(param))
	if err != nil {
		return 0, errors.New(fmt.Sprintf("bad %s", param))
	}

	return userID, nil
}

func (h *Handler) responseError(rw http.ResponseWriter, response []byte, code int) {
	h.logger.Warn(string(response), ", response code ", code)
	rw.Write(response)
	rw.WriteHeader(code)
}

func (h *Handler) createEvent() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		err := r.ParseForm()
		if err != nil {
			h.responseError(rw, []byte("bad request"), http.StatusBadRequest)
			return
		}

		event, err := h.postEvent(r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		userID, err := h.postFormGet("user_id", r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		h.eventManager.SetEvent(userID, event)
	}
}

func (h *Handler) updateEvent() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		err := r.ParseForm()
		if err != nil {
			h.responseError(rw, []byte("bad request"), http.StatusBadRequest)
			return
		}

		event, err := h.postEvent(r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		userID, err := h.postFormGet("user_id", r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		id, err := h.postFormGet("id", r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		h.eventManager.UpdateEvent(userID, id, event)
	}
}

func (h *Handler) deleteEvent() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		err := r.ParseForm()
		if err != nil {
			h.responseError(rw, []byte("bad request"), http.StatusBadRequest)
			return
		}

		userID, err := h.postFormGet("user_id", r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		id, err := h.postFormGet("id", r)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		h.eventManager.DeleteEvent(userID, id)
	}
}

type queryParams struct {
	userID int
	date   time.Time
}

func hasQuery(values url.Values, key string) bool {
	if _, found := values[key]; !found {
		return false
	}
	if len(values[key]) != 1 {
		return false
	}
	return true
}

func getQuery(values url.Values, key string) string {
	if _, found := values[key]; !found {
		return ""
	}
	if len(values[key]) != 1 {
		return ""
	}
	return values[key][0]
}

func getParams(query url.Values) (queryParams, error) {
	if !hasQuery(query, "user_id") || !hasQuery(query, "date") {
		return queryParams{}, errors.New("Bad params")
	}

	userID, err := strconv.Atoi(getQuery(query, "user_id"))
	if err != nil {
		return queryParams{}, errors.New("bad user id")
	}

	date, err := time.Parse(formatDate, getQuery(query, "date"))
	if err != nil {
		return queryParams{}, errors.New("bad date")
	}

	return queryParams{userID: userID, date: date}, nil
}

func (h *Handler) dayEvents() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		query := r.URL.Query()

		eventsParams, err := getParams(query)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		events := h.eventManager.GetEvents(eventsParams.userID, eventsParams.date, eventsParams.date.Add(time.Hour*24))

		bytes, err := marshalJSON(events)
		if err != nil {
			h.responseError(rw, []byte("result error"), http.StatusServiceUnavailable)
			return
		}

		rw.Write(bytes)
	}
}

func (h *Handler) weekEvents() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		query := r.URL.Query()

		eventsParams, err := getParams(query)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		events := h.eventManager.GetEvents(eventsParams.userID, eventsParams.date, eventsParams.date.Add(time.Hour*24*7))

		bytes, err := marshalJSON(events)
		if err != nil {
			h.responseError(rw, []byte("result error"), http.StatusServiceUnavailable)
			return
		}

		rw.Write(bytes)
	}
}

func (h *Handler) monthEvents() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.logger.Info(r.URL.Path)

		query := r.URL.Query()

		eventsParams, err := getParams(query)
		if err != nil {
			h.responseError(rw, []byte(err.Error()), http.StatusBadRequest)
			return
		}

		events := h.eventManager.GetEvents(eventsParams.userID, eventsParams.date, eventsParams.date.Add(time.Hour*24*30))

		bytes, err := marshalJSON(events)
		if err != nil {
			h.responseError(rw, []byte("result error"), http.StatusServiceUnavailable)
			return
		}

		rw.Write(bytes)
	}
}

func marshalJSON(events []event.Event) ([]byte, error) {
	response := struct {
		Result []event.Event `json:"result"`
	}{events}
	marshal, err := json.Marshal(&response)
	return marshal, err
}
