package meetings

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Meetings struct {
	Id          int64  `json:"id"`
	RoomId      int64  `json:"room_id"`
	UserLogin   string `json:"user_login"`
	NameMeeting string `json:"name_meeting"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Result      string `json:"result"`
}

func (s *Service) AddNewMeeting(meetings Meetings, pool *pgxpool.Pool) (err error) {
	var first, second int64
	err = pool.QueryRow(context.Background(), `SELECT id FROM meetings WHERE start_time <= $1 and end_time >= $1 and room_id = $2`, meetings.StartTime, meetings.RoomId).Scan(&first)
	if err == nil{
		return
	}
	err = pool.QueryRow(context.Background(), `SELECT id FROM meetings WHERE start_time <= $1 and end_time >= $1 and room_id = $2`, meetings.StartTime, meetings.RoomId).Scan(&second)
	if err == nil {
		return
	}
	if first == 0 && second == 0{
	_, err = pool.Exec(context.Background(), `INSERT INTO meetings(room_id, user_login, name_meeting, start_time, end_time)
VALUES ($1, $2, $3, $4, $5);`, meetings.RoomId, meetings.UserLogin,  meetings.NameMeeting, meetings.StartTime, meetings.EndTime)
	if err != nil {
		return
	}
	return nil}else{
		return errors.New("In this time meeting has have")
	}
}

func (s *Service) AllMeeting(pool *pgxpool.Pool) (list []Meetings ,err error)  {
	rows, err := pool.Query(context.Background(), `SELECT id, room_id, user_login, name_meeting, start_time, end_time, result FROM meetings;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		meetings := Meetings{}
		err := rows.Scan(&meetings.Id, &meetings.RoomId, &meetings.UserLogin, &meetings.NameMeeting, &meetings.StartTime, &meetings.EndTime, &meetings.Result)
		if err != nil {
			return nil, err
		}
		list = append(list, meetings)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) MeetingByRoomID(id int64, pool *pgxpool.Pool) (list []Meetings, err error) {
	query, err := pool.Query(context.Background(), `select id,room_id, user_login, name_meeting, start_time, end_time, result  from meetings where room_id=$1;`,
		id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't list from database meetings (id: %d)!", id))
	}
	for query.Next() {
		meetings := Meetings{}
		err := query.Scan(&meetings.Id, &meetings.RoomId, &meetings.UserLogin, &meetings.NameMeeting, &meetings.StartTime, &meetings.EndTime, &meetings.Result)
		if err != nil {
			return nil, err
		}
		list = append(list, meetings)
	}
	err = query.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) AddResultById(meetings Meetings, id int64, pool *pgxpool.Pool) (err error) {
	timestamp := time.Now().Unix()
	var first int64
	err = pool.QueryRow(context.Background(), `SELECT id FROM meetings WHERE end_time <= $1 and id = $2`, timestamp, id).Scan(&first)
	if err != nil{
		return
	}
	if first != 0{
	_, err = pool.Exec(context.Background(), `update meetings set result = $1 where id = $2`, meetings.Result, id)
	if err != nil {
		return
	}
	return nil}else{
		return errors.New("meeting have now end")
	}
}

func (s *Service) MeetingsCurrentlyAndInThisRoom(id int64, pool *pgxpool.Pool) (meetings Meetings, err error) {
	timestamp := time.Now().Unix()
	err = pool.QueryRow(context.Background(),
		`SELECT id, room_id, user_login, name_meeting, start_time, end_time, result FROM meetings WHERE start_time <= $1 and end_time >= $1 and room_id = $2`,
		timestamp, id).Scan(&meetings.Id, &meetings.RoomId, &meetings.UserLogin, &meetings.NameMeeting, &meetings.StartTime, &meetings.EndTime, &meetings.Result)
	if err != nil {
		return Meetings{}, errors.New(fmt.Sprintf("can't list from database meetings (id: %d)!", id))
	}
	return

}

