// Package handlers contiene los manejadores HTTP para las rutas de la API.
// Cada función aquí corresponde a un endpoint específico y maneja la lógica de recibir solicitudes, interactuar con el repositorio y enviar respuestas.
package handlers

import (
	"encoding/json"
	"errors" // Para usar gorm.ErrRecordNotFound
	"net/http"
	"strconv" // Para convertir ID de string a int

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm" // Para gorm.Expr y gorm.ErrRecordNotFound

	"lab6/models"     // Asegúrate que la ruta de importación sea correcta
	"lab6/repository" // Asegúrate que la ruta de importación sea correcta
)

// ErrorResponse define una estructura estándar para respuestas de error JSON.
// Se utiliza para dar formato consistente a los errores de la API en las respuestas HTTP.
// @Description Estructura estándar para errores de la API con un mensaje descriptivo.
type ErrorResponse struct {
	// Message contiene el mensaje descriptivo del error ocurrido.
	// example: "Serie no encontrada"
	Message string `json:"message"`
}

// writeError es una función helper para escribir errores JSON estandarizados.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// --- Handlers ---

// GetAllSeries godoc
// @Summary      Listar todas las series
// @Description  Obtiene una lista completa de todas las series almacenadas en la base de datos.
// @Tags         Series
// @Accept       json
// @Produce      json
// @Success      200 {array}  models.Series "Lista de series recuperada exitosamente"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al buscar series"
// @Router       /series [get]
func GetAllSeries(w http.ResponseWriter, r *http.Request) {
	var series []models.Series
	if err := repository.DB.Find(&series).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error buscando series: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Explícito OK
	json.NewEncoder(w).Encode(series)
}

// GetSeriesByID godoc
// @Summary      Obtener una serie por ID
// @Description  Obtiene los detalles de una serie específica usando su ID numérico proporcionado en la URL.
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie a buscar" Format(int64) example(1)
// @Success      200 {object} models.Series "Detalles de la serie encontrados"
// @Failure      400 {object} ErrorResponse "ID proporcionado inválido (no es un número)"
// @Failure      404 {object} ErrorResponse "Serie no encontrada con el ID proporcionado"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al buscar la serie"
// @Router       /series/{id} [get]
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	var serie models.Series
	result := repository.DB.First(&serie, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+result.Error.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie)
}

// CreateSeries godoc
// @Summary      Crear una nueva serie
// @Description  Añade una nueva serie a la base de datos utilizando los datos proporcionados en el cuerpo de la solicitud. El ID es auto-generado por la base de datos.
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        series body models.Series true "Datos de la nueva serie a crear (el campo ID será ignorado)"
// @Success      201 {object} models.Series "Serie creada exitosamente (devuelve el objeto completo con el nuevo ID)"
// @Failure      400 {object} ErrorResponse "Entrada inválida (ej. JSON mal formado, falta título)"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al guardar la serie"
// @Router       /series [post]
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var newSeries models.Series
	// Decodificar el cuerpo de la solicitud en la estructura newSeries
	if err := json.NewDecoder(r.Body).Decode(&newSeries); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de la solicitud inválido: "+err.Error())
		return
	}

	// Validación básica (GORM podría manejar 'not null', pero una validación explícita es buena)
	if newSeries.Title == "" {
		writeError(w, http.StatusBadRequest, "El campo 'title' es obligatorio")
		return
	}

	// Crear el registro en la base de datos
	// GORM asignará el ID automáticamente si la creación es exitosa.
	if err := repository.DB.Create(&newSeries).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error creando la serie: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(newSeries) // Devolver el objeto creado con su ID
}

// UpdateSeries godoc
// @Summary      Actualizar una serie existente
// @Description  Actualiza todos los campos de una serie existente identificada por su ID, utilizando los datos proporcionados en el cuerpo de la solicitud.
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie a actualizar" example(1)
// @Param        series body models.Series true "Nuevos datos completos para la serie (se usará el ID de la URL, no el del cuerpo si existe)"
// @Success      200 {object} models.Series "Serie actualizada exitosamente"
// @Failure      400 {object} ErrorResponse "Entrada inválida (ej. JSON mal formado, ID inválido en URL)"
// @Failure      404 {object} ErrorResponse "Serie no encontrada con el ID proporcionado"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al actualizar la serie"
// @Router       /series/{id} [put]
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	// Verificar primero si la serie existe
	var serie models.Series
	result := repository.DB.First(&serie, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada para actualizar")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+result.Error.Error())
		}
		return
	}

	// Decodificar los datos actualizados del cuerpo de la solicitud
	var updatedData models.Series
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de la solicitud inválido: "+err.Error())
		return
	}

	// Actualizar los campos del objeto 'serie' existente con los 'updatedData'
	// Es importante usar el ID de la URL, no el del cuerpo (si lo tuviera)
	serie.Title = updatedData.Title
	serie.Status = updatedData.Status
	serie.LastEpisodeWatched = updatedData.LastEpisodeWatched
	serie.TotalEpisodes = updatedData.TotalEpisodes
	serie.Ranking = updatedData.Ranking

	// Guardar los cambios en la base de datos
	if err := repository.DB.Save(&serie).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error actualizando la serie: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie) // Devolver el objeto actualizado
}

// DeleteSeries godoc
// @Summary      Eliminar una serie
// @Description  Elimina permanentemente una serie de la base de datos utilizando su ID.
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie a eliminar" example(1)
// @Success      204 "Sin contenido (eliminado exitosamente)"
// @Failure      400 {object} ErrorResponse "ID proporcionado inválido (no es un número)"
// @Failure      404 {object} ErrorResponse "Serie no encontrada para eliminar"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al eliminar la serie"
// @Router       /series/{id} [delete]
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	// Intentar eliminar la serie directamente por ID
	// Usamos un modelo vacío para que GORM sepa en qué tabla buscar
	result := repository.DB.Delete(&models.Series{}, id)

	if result.Error != nil {
		// GORM no devuelve ErrRecordNotFound en Delete si no afecta filas,
		// por eso verificamos RowsAffected después.
		writeError(w, http.StatusInternalServerError, "Error eliminando la serie: "+result.Error.Error())
		return
	}

	// Verificar si alguna fila fue realmente eliminada
	if result.RowsAffected == 0 {
		writeError(w, http.StatusNotFound, "Serie no encontrada para eliminar")
		return
	}

	// Éxito, no devolver cuerpo
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// --- Handlers de Acciones Específicas (PATCH) ---

// UpdateSeriesStatus godoc
// @Summary      Actualizar estado de una serie (parcial)
// @Description  Actualiza únicamente el campo 'status' de una serie existente identificada por su ID.
// @Tags         Series Actions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie cuyo estado se actualizará" example(1)
// @Param        status body models.StatusUpdate true "Objeto JSON con el nuevo estado"
// @Success      200 {object} models.Series "Estado actualizado, devuelve la serie completa"
// @Failure      400 {object} ErrorResponse "Entrada inválida (ej. JSON mal formado, falta status, ID inválido)"
// @Failure      404 {object} ErrorResponse "Serie no encontrada"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al actualizar el estado"
// @Router       /series/{id}/status [patch]
func UpdateSeriesStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	// Verificar si la serie existe
	var serie models.Series
	if err := repository.DB.First(&serie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+err.Error())
		}
		return
	}

	// Decodificar el cuerpo de la solicitud que contiene solo el estado
	var statusUpdate models.StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		writeError(w, http.StatusBadRequest, "Cuerpo de la solicitud inválido: "+err.Error())
		return
	}

	// Validar que el status no esté vacío
	if statusUpdate.Status == "" {
		writeError(w, http.StatusBadRequest, "El campo 'status' no puede estar vacío")
		return
	}

	// Actualizar solo el campo 'Status' usando Update para eficiencia
	if err := repository.DB.Model(&serie).Update("status", statusUpdate.Status).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error actualizando el estado de la serie: "+err.Error())
		return
	}

	// Refrescar los datos de la serie (Update no actualiza el struct original por defecto)
	repository.DB.First(&serie, id) // Volver a leer para obtener el estado actualizado

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// IncrementSeriesEpisode godoc
// @Summary      Incrementar episodio visto
// @Description  Incrementa en 1 el campo 'lastEpisodeWatched' de la serie identificada por ID. No realiza cambios si el último episodio visto ya es igual o mayor al total de episodios.
// @Tags         Series Actions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie cuyo episodio se incrementará" example(1)
// @Success      200 {object} models.Series "Episodio incrementado, devuelve la serie actualizada"
// @Failure      400 {object} ErrorResponse "ID proporcionado inválido"
// @Failure      404 {object} ErrorResponse "Serie no encontrada"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al incrementar el episodio"
// @Router       /series/{id}/episode [patch]
func IncrementSeriesEpisode(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	var serie models.Series
	if err := repository.DB.First(&serie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+err.Error())
		}
		return
	}

	// No incrementar si ya se alcanzó o superó el total (si el total es > 0)
	if serie.TotalEpisodes > 0 && serie.LastEpisodeWatched >= serie.TotalEpisodes {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // OK, pero no se hizo cambio
		json.NewEncoder(w).Encode(serie)
		return
	}

	// Incrementar el contador de episodios vistos usando expresión SQL para atomicidad
	newEpisodeCount := serie.LastEpisodeWatched + 1
	if err := repository.DB.Model(&serie).Update("last_episode_watched", gorm.Expr("last_episode_watched + ?", 1)).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error incrementando el episodio: "+err.Error())
		return
	}

	// Actualizar el valor en la struct local para la respuesta
	serie.LastEpisodeWatched = newEpisodeCount

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// UpvoteSeries godoc
// @Summary      Votar positivamente (Upvote) una serie
// @Description  Incrementa en 1 el campo 'ranking' de la serie identificada por su ID.
// @Tags         Series Actions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie a votar positivamente" example(1)
// @Success      200 {object} models.Series "Ranking incrementado, devuelve la serie actualizada"
// @Failure      400 {object} ErrorResponse "ID proporcionado inválido"
// @Failure      404 {object} ErrorResponse "Serie no encontrada"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al actualizar el ranking"
// @Router       /series/{id}/upvote [patch]
func UpvoteSeries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	var serie models.Series
	// Verificar si existe antes de intentar actualizar
	if err := repository.DB.First(&serie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+err.Error())
		}
		return
	}

	// Incrementar el ranking usando una expresión SQL para atomicidad
	if err := repository.DB.Model(&serie).Update("ranking", gorm.Expr("ranking + 1")).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error al votar positivamente (upvote): "+err.Error())
		return
	}

	// Refrescar la struct para obtener el nuevo valor del ranking
	repository.DB.First(&serie, id) // Volver a leer de la DB

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}

// DownvoteSeries godoc
// @Summary      Votar negativamente (Downvote) una serie
// @Description  Decrementa en 1 el campo 'ranking' de la serie identificada por su ID.
// @Tags         Series Actions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la Serie a votar negativamente" example(1)
// @Success      200 {object} models.Series "Ranking decrementado, devuelve la serie actualizada"
// @Failure      400 {object} ErrorResponse "ID proporcionado inválido"
// @Failure      404 {object} ErrorResponse "Serie no encontrada"
// @Failure      500 {object} ErrorResponse "Error interno del servidor al actualizar el ranking"
// @Router       /series/{id}/downvote [patch]
func DownvoteSeries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido: "+idStr)
		return
	}

	var serie models.Series
	// Verificar si existe antes de intentar actualizar
	if err := repository.DB.First(&serie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "Serie no encontrada")
		} else {
			writeError(w, http.StatusInternalServerError, "Error buscando la serie: "+err.Error())
		}
		return
	}

	// Decrementar el ranking usando una expresión SQL
	// Considerar si se quiere prevenir que el ranking sea negativo
	if err := repository.DB.Model(&serie).Update("ranking", gorm.Expr("ranking - 1")).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Error al votar negativamente (downvote): "+err.Error())
		return
	}

	// Refrescar la struct para obtener el nuevo valor del ranking
	repository.DB.First(&serie, id) // Volver a leer de la DB

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie) // Devuelve la serie actualizada
}
