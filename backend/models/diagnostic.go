package models

import "time"

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
