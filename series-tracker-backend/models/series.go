package models

// Series representa la estructura de una serie de TV en la base de datos.
type Series struct {
	ID                int    `json:"id"`
	Title             string `json:"title" binding:"required"` // binding:"required" para validación con Gin
	Status            string `json:"status"`                 // Valores: 'Plan to Watch', 'Watching', 'Completed', 'Dropped'
	LastEpisodeWatched int    `json:"lastEpisodeWatched"`
	TotalEpisodes     int    `json:"totalEpisodes"`
	Ranking           int    `json:"ranking"`
}

// StatusUpdate se usa específicamente para el endpoint PATCH /status
type StatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
