package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/new-mux/pkg/mux"
	"github.com/shuhrat-shokirov/rest/pkg/rest"
	"log"
	"mitings-service/pkg/core/meetings"
	"net/http"
	"strconv"
)

type Server struct {
	router      *mux.ExactMux
	pool        *pgxpool.Pool
	meetingsSvc *meetings.Service
	secret      jwt.Secret
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, mitingsSvc *meetings.Service, secret jwt.Secret) *Server {
	return &Server{router: router, pool: pool, meetingsSvc: mitingsSvc, secret: secret}
}

func (s Server) ServeHTTP(writer http.ResponseWriter,request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s Server) handleNewMeeting() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		meetings := meetings.Meetings{}
		err := rest.ReadJSONBody(request, &meetings)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(meetings)
		err = s.meetingsSvc.AddNewMeeting(meetings, s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New History Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleMitingsList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.meetingsSvc.AllMeeting(s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleMeetingByRoomID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.meetingsSvc.MeetingByRoomID(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleAddResultById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		meetings := meetings.Meetings{}
		err = rest.ReadJSONBody(request, &meetings)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(meetings)
		err = s.meetingsSvc.AddResultById(meetings, int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New Result by id Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleMeetingCurrentlyAndInThisRoom() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.meetingsSvc.MeetingsCurrentlyAndInThisRoom(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}