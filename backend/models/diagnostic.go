package models

import (
	"fmt"
	"time"
)

// CPUInfo représente les informations du processeur
type CPUInfo struct {
	Model       string `json:"model"`
	Cores       int    `json:"cores"`
	Frequency   string `json:"frequency"`
	Temperature string `json:"temperature,omitempty"`
}

// RAMInfo représente les informations de la mémoire
type RAMInfo struct {
	Total     string `json:"total"`
	Used      string `json:"used"`
	Available string `json:"available"`
	Type      string `json:"type,omitempty"`
}

// StorageInfo représente les informations du stockage
type StorageInfo struct {
	Type       string `json:"type"` // SSD, HDD
	Capacity   string `json:"capacity"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	Health     string `json:"health,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
}

// BatteryInfo représente les informations de la batterie
type BatteryInfo struct {
	CycleCount   int    `json:"cycle_count"`
	Health       string `json:"health"`
	Capacity     string `json:"capacity"`
	MaxCapacity  string `json:"max_capacity,omitempty"`
	Condition    string `json:"condition,omitempty"`
	IsCharging   bool   `json:"is_charging"`
	PowerAdapter string `json:"power_adapter,omitempty"`
}

// SystemInfo représente les informations générales du système
type SystemInfo struct {
	MachineName  string `json:"machine_name"`
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
	OSVersion    string `json:"os_version"`
	MacOSVersion string `json:"macos_version,omitempty"`
}

// Diagnostic représente le diagnostic complet d'une machine
type Diagnostic struct {
	ID         int64       `json:"id"`
	SystemInfo SystemInfo  `json:"system_info"`
	CPU        CPUInfo     `json:"cpu"`
	RAM        RAMInfo     `json:"ram"`
	Storage    StorageInfo `json:"storage"`
	Battery    BatteryInfo `json:"battery"`
	Status     string      `json:"status"`   // success, partial, failed
	Duration   float64     `json:"duration"` // en secondes
	Timestamp  time.Time   `json:"timestamp"`
	CreatedAt  time.Time   `json:"created_at"`
}

// DiagnosticRequest représente la requête pour créer un diagnostic
type DiagnosticRequest struct {
	SystemInfo SystemInfo  `json:"system_info"`
	CPU        CPUInfo     `json:"cpu"`
	RAM        RAMInfo     `json:"ram"`
	Storage    StorageInfo `json:"storage"`
	Battery    BatteryInfo `json:"battery"`
	Status     string      `json:"status"`
	Duration   float64     `json:"duration"`
}

// DiagnosticResponse représente la réponse après création d'un diagnostic
type DiagnosticResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      int64  `json:"id,omitempty"`
}

// DiagnosticsListResponse représente la liste des diagnostics
type DiagnosticsListResponse struct {
	Success     bool         `json:"success"`
	Count       int          `json:"count"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// SwiftDiagnosticRequest représente le format envoyé par l'application Swift
type SwiftDiagnosticRequest struct {
	MachineName         string  `json:"machine_name"`
	SerialNumber        string  `json:"serial_number"`
	CPUModel            string  `json:"cpu_model"`
	CPUCores            int     `json:"cpu_cores"`
	RAMTotalGB          float64 `json:"ram_total_gb"`
	RAMUsedGB           float64 `json:"ram_used_gb"`
	StorageTotalGB      float64 `json:"storage_total_gb"`
	StorageUsedGB       float64 `json:"storage_used_gb"`
	BatteryCycleCount   int     `json:"battery_cycle_count"`
	BatteryPercentage   int     `json:"battery_percentage"`
	BatteryHealth       string  `json:"battery_health"`
	TestDurationSeconds float64 `json:"test_duration_seconds"`
	Status              string  `json:"status"`
}

// ToStandardRequest convertit le format Swift vers le format standard
func (s *SwiftDiagnosticRequest) ToStandardRequest() DiagnosticRequest {
	ramAvailable := s.RAMTotalGB - s.RAMUsedGB
	storageAvailable := s.StorageTotalGB - s.StorageUsedGB

	return DiagnosticRequest{
		SystemInfo: SystemInfo{
			MachineName:  s.MachineName,
			SerialNumber: s.SerialNumber,
			Model:        "Unknown",
			OSVersion:    "macOS",
		},
		CPU: CPUInfo{
			Model:     s.CPUModel,
			Cores:     s.CPUCores,
			Frequency: "N/A",
		},
		RAM: RAMInfo{
			Total:     formatGB(s.RAMTotalGB),
			Used:      formatGB(s.RAMUsedGB),
			Available: formatGB(ramAvailable),
		},
		Storage: StorageInfo{
			Type:      "SSD",
			Capacity:  formatGB(s.StorageTotalGB),
			Used:      formatGB(s.StorageUsedGB),
			Available: formatGB(storageAvailable),
		},
		Battery: BatteryInfo{
			CycleCount: s.BatteryCycleCount,
			Health:     s.BatteryHealth,
			Capacity:   formatPercent(s.BatteryPercentage),
			IsCharging: false,
		},
		Status:   mapStatus(s.Status),
		Duration: s.TestDurationSeconds,
	}
}

// formatGB formate un nombre en string avec GB
func formatGB(gb float64) string {
	return fmt.Sprintf("%.2f GB", gb)
}

// formatPercent formate un nombre en pourcentage
func formatPercent(percent int) string {
	return fmt.Sprintf("%d%%", percent)
}

// mapStatus convertit le statut Swift vers le format backend
func mapStatus(swiftStatus string) string {
	switch swiftStatus {
	case "completed":
		return "success"
	case "failed":
		return "failed"
	default:
		return "partial"
	}
}
