package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"lab6/models"
	"lab6/repository"
	"gorm.io/gorm"
)

// Obtener todas las series
func GetAllSeries(w http.ResponseWriter, r *http.Request) {
	var series []models.Series
	// Añadir manejo de errores básico
	if err := repository.DB.Find(&series).Error; err != nil {
		http.Error(w, "Error fetching series: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

// Obtener una serie por ID
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie)
}

// Crear una nueva serie
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var newSeries models.Series
	if err := json.NewDecoder(r.Body).Decode(&newSeries); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validación básica 
	if newSeries.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := repository.DB.Create(&newSeries).Error; err != nil {
		http.Error(w, "Error creating series: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(newSeries)
}

// Actualizar una serie existente 
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id) // Primero busca si existe

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updatedData models.Series
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	serie.Title = updatedData.Title
	serie.Status = updatedData.Status
	serie.LastEpisodeWatched = updatedData.LastEpisodeWatched
	serie.TotalEpisodes = updatedData.TotalEpisodes
	serie.Ranking = updatedData.Ranking

	if err := repository.DB.Save(&serie).Error; err != nil {
		http.Error(w, "Error updating series: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie)
}

// Eliminar una serie
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series 

	result := repository.DB.Delete(&serie, id)

	if result.Error != nil {
		http.Error(w, "Error deleting series: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Serie no encontrada para eliminar", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}


// UpdateSeriesStatus actualiza solo el estado de una serie
func UpdateSeriesStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Decodificar el cuerpo de la solicitud que contiene solo el estado
	var statusUpdate models.StatusUpdate 
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Actualizar solo el campo 'Status' usando Update para eficiencia
	if err := repository.DB.Model(&serie).Update("status", statusUpdate.Status).Error; err != nil {
		http.Error(w, "Error updating series status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// IncrementSeriesEpisode incrementa el último episodio visto
func IncrementSeriesEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Incrementar el contador de episodios vistos
	if serie.TotalEpisodes > 0 && serie.LastEpisodeWatched >= serie.TotalEpisodes {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serie)
		return
	}


	// Actualizar usando Update para la operación de incremento
	newEpisodeCount := serie.LastEpisodeWatched + 1
	if err := repository.DB.Model(&serie).Update("last_episode_watched", newEpisodeCount).Error; err != nil {
		http.Error(w, "Error incrementing episode: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Actualizar el valor en la struct local para la respuesta
	serie.LastEpisodeWatched = newEpisodeCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// UpvoteSeries incrementa el ranking de una serie
func UpvoteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Incrementar el ranking usando una expresión SQL para atomicidad
	if err := repository.DB.Model(&serie).Update("ranking", gorm.Expr("ranking + 1")).Error; err != nil {
		http.Error(w, "Error upvoting series: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Refrescar la struct para obtener el nuevo valor del ranking
	repository.DB.First(&serie, id) // Volver a leer de la DB

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// DownvoteSeries decrementa el ranking de una serie
func DownvoteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding series: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := repository.DB.Model(&serie).Update("ranking", gorm.Expr("ranking - 1")).Error; err != nil {
		http.Error(w, "Error downvoting series: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Refrescar la struct para obtener el nuevo valor del ranking
	repository.DB.First(&serie, id) // Volver a leer de la DB

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}
