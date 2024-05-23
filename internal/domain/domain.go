package domain

import (
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser"
	"github.com/DedAzaMarks/ABS/internal/server/statemachine"
	"github.com/google/uuid"
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

type TGUser struct {
	UserID        string
	State         *statemachine.StateMachine
	SearchResults []SignedSearchResult
	FilmResults   []SignedFilmResult
	Client        uuid.UUID
}

func NewTGUser(userID string) *TGUser {
	return &TGUser{
		UserID: userID,
		State:  statemachine.NewStateMachine(),
	}
}
