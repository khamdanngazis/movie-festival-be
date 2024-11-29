package models

type SaveMovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    int    `json:"duration" validate:"required"` // in minutes
	Artists     string `json:"artists" validate:"required"`
	Genres      string `json:"genres" validate:"required"`
	WatchURL    string `json:"watch_url" validate:"required"`
}
