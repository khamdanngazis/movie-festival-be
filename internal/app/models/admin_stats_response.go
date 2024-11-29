package models

type AdminStatsResponse struct {
	MostVotedMovie  MovieStatsResponse `json:"most_voted_movie"`
	MostViewedGenre GenreStatsResponse `json:"most_viewed_genre"`
}

type MovieStatsResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Artists     string `json:"artists"`
	Genres      string `json:"genres"`
	VotesCount  int64  `json:"votes_count"`
}

type GenreStatsResponse struct {
	Genre      string `json:"genre"`
	TotalViews int64  `json:"total_views"`
}
