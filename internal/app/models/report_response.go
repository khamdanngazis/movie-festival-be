package models

type ReportViewsResponse struct {
	MostViewedMovie MovieViewResponse `json:"most_viewed_movie"`
	GenreStats      []GenreStatItem   `json:"genre_stats"`
}

type MovieViewResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Artists     string `json:"artists"`
	Genres      string `json:"genres"`
	ViewCount   int64  `json:"view_count"`
}

type GenreStatItem struct {
	Genre      string `json:"genre"`
	TotalViews int    `json:"total_views"`
}
