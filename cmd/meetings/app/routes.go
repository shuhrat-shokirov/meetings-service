package app

import (
	"mitings-service/pkg/mux/middleware/logger"
)

func (s Server) InitRoutes() {
	s.router.POST(
		"/api/meetings/0",
		s.handleNewMeeting(),
		logger.Logger("add new history"),
	)
	s.router.GET(
		"/api/meetings",
		s.handleMitingsList(),
		logger.Logger("get list"),
	)
	s.router.GET(
		"/api/meetings/{id}",
		s.handleMeetingByRoomID(),
		logger.Logger("get history by room_id"),
	)
	s.router.POST(
		"/api/meetings/add/result/{id}",
		s.handleAddResultById(),
		logger.Logger("add result in history by id"),
	)
	s.router.GET(
		"/api/meetings/room/{id}",
		s.handleMeetingCurrentlyAndInThisRoom(),
		logger.Logger("get history by room_id"),
	)
}