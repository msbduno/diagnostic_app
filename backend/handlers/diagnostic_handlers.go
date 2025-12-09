package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"diagnostic-backend/database"
	"diagnostic-backend/models"

	"github.com/gorilla/mux"
)

// CreateDiagnostic gère la création d'un nouveau diagnostic
func CreateDiagnostic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lire le body en bytes pour pouvoir le parser deux fois si nécessaire
	var rawData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		log.Printf("Erreur de décodage JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.DiagnosticResponse{
			Success: false,
			Message: "Format JSON invalide: " + err.Error(),
		})
		return
	}

	var diagReq models.DiagnosticRequest

	// Vérifier si c'est le format Swift (plat) ou le format standard (imbriqué)
	if _, hasSystemInfo := rawData["system_info"]; hasSystemInfo {
		// Format standard (imbriqué)
		jsonBytes, _ := json.Marshal(rawData)
		if err := json.Unmarshal(jsonBytes, &diagReq); err != nil {
			log.Printf("Erreur de parsing format standard: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.DiagnosticResponse{
				Success: false,
				Message: "Format JSON invalide: " + err.Error(),
			})
			return
		}
	} else {
		// Format Swift (plat)
		var swiftReq models.SwiftDiagnosticRequest
		jsonBytes, _ := json.Marshal(rawData)
		if err := json.Unmarshal(jsonBytes, &swiftReq); err != nil {
			log.Printf("Erreur de parsing format Swift: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.DiagnosticResponse{
				Success: false,
				Message: "Format JSON invalide: " + err.Error(),
			})
			return
		}
		// Convertir au format standard
		diagReq = swiftReq.ToStandardRequest()
		log.Printf("📱 Format Swift détecté et converti")
	}

	// Valider les données
	if err := validateDiagnostic(diagReq); err != nil {
		log.Printf("Validation échouée: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.DiagnosticResponse{
			Success: false,
			Message: "Validation échouée: " + err.Error(),
		})
		return
	}

	// Insérer dans la base de données
	id, err := database.CreateDiagnostic(diagReq)
	if err != nil {
		log.Printf("Erreur de base de données: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.DiagnosticResponse{
			Success: false,
			Message: "Erreur lors de la sauvegarde: " + err.Error(),
		})
		return
	}

	log.Printf("Diagnostic créé avec succès - ID: %d, Machine: %s, Serial: %s",
		id, diagReq.SystemInfo.MachineName, diagReq.SystemInfo.SerialNumber)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.DiagnosticResponse{
		Success: true,
		Message: "Diagnostic enregistré avec succès",
		ID:      id,
	})
}

// GetDiagnostics récupère tous les diagnostics
func GetDiagnostics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Paramètre optionnel: limit
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	diagnostics, err := database.GetAllDiagnostics(limit)
	if err != nil {
		log.Printf("Erreur de récupération: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.DiagnosticsListResponse{
			Success: false,
			Count:   0,
		})
		return
	}

	if diagnostics == nil {
		diagnostics = []models.Diagnostic{}
	}

	log.Printf(" Récupération de %d diagnostics", len(diagnostics))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DiagnosticsListResponse{
		Success:     true,
		Count:       len(diagnostics),
		Diagnostics: diagnostics,
	})
}

// GetDiagnosticByID récupère un diagnostic spécifique
func GetDiagnosticByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID invalide",
		})
		return
	}

	diagnostic, err := database.GetDiagnosticByID(id)
	if err != nil {
		log.Printf("Diagnostic non trouvé (ID: %d): %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Diagnostic non trouvé",
		})
		return
	}

	log.Printf(" Récupération du diagnostic ID: %d", id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"diagnostic": diagnostic,
	})
}

// GetDiagnosticsBySerial récupère tous les diagnostics d'une machine
func GetDiagnosticsBySerial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	serialNumber := vars["serial"]

	diagnostics, err := database.GetDiagnosticsBySerialNumber(serialNumber)
	if err != nil {
		log.Printf("Erreur de récupération: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.DiagnosticsListResponse{
			Success: false,
			Count:   0,
		})
		return
	}

	if diagnostics == nil {
		diagnostics = []models.Diagnostic{}
	}

	log.Printf(" Récupération de %d diagnostics pour la machine %s", len(diagnostics), serialNumber)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DiagnosticsListResponse{
		Success:     true,
		Count:       len(diagnostics),
		Diagnostics: diagnostics,
	})
}

// GetStatistics récupère les statistiques générales
func GetStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := database.GetStatistics()
	if err != nil {
		log.Printf("Erreur de récupération des statistiques: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Erreur lors de la récupération des statistiques",
		})
		return
	}

	log.Println("Récupération des statistiques")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"statistics": stats,
	})
}

// HealthCheck vérifie que l'API est fonctionnelle
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"message": "Backend de diagnostic opérationnel",
		"version": "1.0.0",
	})
}

// validateDiagnostic valide les données du diagnostic
func validateDiagnostic(diag models.DiagnosticRequest) error {
	if diag.SystemInfo.MachineName == "" {
		return &ValidationError{"machine_name est requis"}
	}
	if diag.SystemInfo.SerialNumber == "" {
		return &ValidationError{"serial_number est requis"}
	}
	if diag.SystemInfo.Model == "" {
		return &ValidationError{"model est requis"}
	}
	if diag.CPU.Model == "" {
		return &ValidationError{"cpu.model est requis"}
	}
	if diag.CPU.Cores <= 0 {
		return &ValidationError{"cpu.cores doit être supérieur à 0"}
	}
	if diag.RAM.Total == "" {
		return &ValidationError{"ram.total est requis"}
	}
	if diag.Storage.Type == "" {
		return &ValidationError{"storage.type est requis"}
	}
	if diag.Battery.CycleCount < 0 {
		return &ValidationError{"battery.cycle_count ne peut pas être négatif"}
	}
	if diag.Status == "" {
		return &ValidationError{"status est requis"}
	}
	return nil
}

// ValidationError représente une erreur de validation
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
