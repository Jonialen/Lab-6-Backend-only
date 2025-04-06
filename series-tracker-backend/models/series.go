// Package models define las estructuras de datos utilizadas en la aplicación,
// representando las entidades principales como las Series de TV.
package models

// Series representa la estructura de una serie de TV en la base de datos.
// Contiene información sobre el título, estado de visualización, progreso y ranking.
// @Description Estructura de datos para una Serie de TV.
type Series struct {
	// ID es el identificador único de la serie (Clave primaria, autoincremental).
	// example: 1
	ID int `json:"id" gorm:"primaryKey"`

	// Title es el título de la serie. Es un campo obligatorio.
	// example: "Attack on Titan"
	// required: true
	Title string `json:"title" binding:"required" gorm:"not null"`

	// Status indica el estado actual de visualización de la serie.
	// Debe ser uno de: 'Plan to Watch', 'Watching', 'Completed', 'Dropped'.
	// example: "Watching"
	Status string `json:"status"`

	// LastEpisodeWatched es el número del último episodio que el usuario ha visto.
	// example: 10
	LastEpisodeWatched int `json:"lastEpisodeWatched"`

	// TotalEpisodes es el número total de episodios que tiene la serie.
	// example: 24
	TotalEpisodes int `json:"totalEpisodes"`

	// Ranking es una puntuación o valoración asignada a la serie por el usuario.
	// Puede ser modificada mediante los endpoints de upvote/downvote.
	// example: 8
	Ranking int `json:"ranking"`
}

// StatusUpdate se usa específicamente para el endpoint "PATCH /api/series/{id}/status" para actualizar únicamente el estado de la serie de forma parcial.
// @Description Estructura para la actualización parcial del estado de una serie.
type StatusUpdate struct {
	// Status es el nuevo estado que se asignará a la serie. Campo obligatorio.
	// example: "Completed"
	// required: true
	Status string `json:"status" binding:"required"`
}
