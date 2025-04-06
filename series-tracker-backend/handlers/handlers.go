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
	repository.DB.Find(&series)
	json.NewEncoder(w).Encode(series)
}

// Obtener una serie por ID
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

// Crear una nueva serie
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var newSeries models.Series
	json.NewDecoder(r.Body).Decode(&newSeries)

	repository.DB.Create(&newSeries)
	json.NewEncoder(w).Encode(newSeries)
}

// Actualizar una serie existente
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(&serie)
	repository.DB.Save(&serie)
	json.NewEncoder(w).Encode(serie)
}

// Eliminar una serie
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var serie models.Series

	result := repository.DB.First(&serie, id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	repository.DB.Delete(&serie)
	w.WriteHeader(http.StatusNoContent)
}

