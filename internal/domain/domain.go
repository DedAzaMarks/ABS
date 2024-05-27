package domain

import (
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser"
	"github.com/DedAzaMarks/ABS/internal/server/statemachine"
	"github.com/google/uuid"
	"pgregory.net/rand"
	"strings"
)

type SignedSearchResult struct {
	ID           uuid.UUID
	SearchResult parser.SearchResult
}

type SignedFilmResult struct {
	ID         uuid.UUID
	FilmResult parser.FilmResult
}

func (fr *SignedFilmResult) String() string {
	var s strings.Builder
	s.WriteString("{ID:")
	s.WriteString(fr.ID.String())
	s.WriteString(";FilmResult:")
	s.WriteString(fr.FilmResult.Quality)
	s.WriteString(fr.FilmResult.TranslationVoiceover)
	s.WriteString(fr.FilmResult.Author)
	s.WriteString(fr.FilmResult.FileFormat)
	s.WriteString(fr.FilmResult.Size)
	s.WriteByte('}')
	return s.String()
}

type Device struct {
	ID   uuid.UUID
	Name string
}

type User struct {
	UserID     int64
	SessionKey string
	Devices    []Device

	State         *statemachine.StateMachine
	SearchResults []SignedSearchResult
	FilmResults   []SignedFilmResult
}

func TGUser2PB(user *User) *pb.TgUser {
	return &pb.TgUser{
		UserID:     user.UserID,
		SessionKey: user.SessionKey,
		Devices: func() []*pb.TgUser_Device {
			var devices []*pb.TgUser_Device
			for _, device := range user.Devices {
				devices = append(devices, &pb.TgUser_Device{
					Id:   device.ID.String(),
					Name: device.Name,
				})
			}
			return devices
		}(),
		State: string(user.State.CurrentState()),
		SearchResult: func() []*pb.TgUser_SignedSearchResult {
			var searchResults []*pb.TgUser_SignedSearchResult
			for _, searchResult := range user.SearchResults {
				searchResults = append(searchResults, &pb.TgUser_SignedSearchResult{
					Id: searchResult.ID.String(),
					SearchResult: &pb.TgUser_SignedSearchResult_SearchResult{
						Title: searchResult.SearchResult.Title,
						Href:  searchResult.SearchResult.Href,
					},
				})
			}
			return searchResults
		}(),
		FilmResult: func() []*pb.TgUser_SignedFilmResult {
			var filmResults []*pb.TgUser_SignedFilmResult
			for _, filmResult := range user.FilmResults {
				filmResults = append(filmResults, &pb.TgUser_SignedFilmResult{
					Id: filmResult.ID.String(),
					FilmResult: &pb.TgUser_SignedFilmResult_FilmResult{
						Quality:              filmResult.FilmResult.Quality,
						TranslationVoiceover: filmResult.FilmResult.TranslationVoiceover,
						Author:               filmResult.FilmResult.Author,
						FileFormat:           filmResult.FilmResult.FileFormat,
						Size:                 filmResult.FilmResult.Size,
						Magnet:               filmResult.FilmResult.Magnet,
					},
				})
			}
			return filmResults
		}(),
	}

}

func PB2TGUser(user *pb.TgUser) *User {
	return &User{
		UserID:     user.UserID,
		SessionKey: user.SessionKey,
		Devices: func() []Device {
			var devices []Device
			for _, device := range user.Devices {
				devices = append(devices, Device{
					ID:   uuid.MustParse(device.Id),
					Name: device.Name,
				})
			}
			return devices
		}(),
		State: statemachine.SetState(statemachine.State(user.State)),
		SearchResults: func() []SignedSearchResult {
			var searchResults []SignedSearchResult
			for _, ssr := range user.SearchResult {
				searchResults = append(searchResults, SignedSearchResult{
					ID: uuid.MustParse(ssr.Id),
					SearchResult: parser.SearchResult{
						Title: ssr.SearchResult.Title,
						Href:  ssr.SearchResult.Href,
					},
				})
			}
			return searchResults
		}(),
		FilmResults: func() []SignedFilmResult {
			var filmResults []SignedFilmResult
			for _, sfr := range user.FilmResult {
				filmResults = append(filmResults, SignedFilmResult{
					ID: uuid.MustParse(sfr.Id),
					FilmResult: parser.FilmResult{
						Quality:              sfr.FilmResult.Quality,
						TranslationVoiceover: sfr.FilmResult.TranslationVoiceover,
						Author:               sfr.FilmResult.Author,
						FileFormat:           sfr.FilmResult.FileFormat,
						Size:                 sfr.FilmResult.Size,
						Magnet:               sfr.FilmResult.Magnet,
					},
				})
			}
			return filmResults
		}(),
	}
}

func TGUser2DTO(user *User) *UserDTO {
	return &UserDTO{
		ID:         user.UserID,
		SessionKey: user.SessionKey,
		State:      string(user.State.CurrentState()),
		Devices: func() []DeviceDTO {
			var devices []DeviceDTO
			for _, device := range user.Devices {
				devices = append(devices, DeviceDTO{
					ID:   device.ID,
					Name: device.Name,
				})
			}
			return devices
		}(),
	}
}

func DTO2TGUser(dto *UserDTO) *User {
	return &User{
		UserID:     dto.ID,
		SessionKey: dto.SessionKey,
		State:      statemachine.SetState(statemachine.State(dto.State)),
		Devices: func() []Device {
			var devices []Device
			for _, device := range dto.Devices {
				devices = append(devices, Device{
					ID:   device.ID,
					Name: device.Name,
				})
			}
			return devices
		}(),
	}
}

func NewTGUser(userID int64) *User {
	return &User{
		UserID:     userID,
		SessionKey: newSession(),
		State:      statemachine.NewStateMachine(),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func newSession() string {
	var res [8]byte
	for i := 0; i < 8; i++ {
		res[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(res[:])
}
