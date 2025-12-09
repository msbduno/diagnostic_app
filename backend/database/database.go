package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"diagnostic-backend/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initialise la connexion à la base de données SQLite
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur d'ouverture de la base de données: %v", err)
	}

	// Vérifier la connexion
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("erreur de connexion à la base de données: %v", err)
	}

	log.Println("Connexion à la base de données établie")

	// Créer les tables
	if err = createTables(); err != nil {
		return fmt.Errorf("erreur de création des tables: %v", err)
	}

	return nil
}

// createTables crée les tables nécessaires
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS diagnostics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		machine_name TEXT NOT NULL,
		serial_number TEXT NOT NULL,
		model TEXT NOT NULL,
		os_version TEXT NOT NULL,
		macos_version TEXT,
		
		cpu_model TEXT NOT NULL,
		cpu_cores INTEGER NOT NULL,
		cpu_frequency TEXT NOT NULL,
		cpu_temperature TEXT,
		
		ram_total TEXT NOT NULL,
		ram_used TEXT NOT NULL,
		ram_available TEXT NOT NULL,
		ram_type TEXT,
		
		storage_type TEXT NOT NULL,
		storage_capacity TEXT NOT NULL,
		storage_used TEXT NOT NULL,
		storage_available TEXT NOT NULL,
		storage_health TEXT,
		storage_device_name TEXT,
		
		battery_cycle_count INTEGER NOT NULL,
		battery_health TEXT NOT NULL,
		battery_capacity TEXT NOT NULL,
		battery_max_capacity TEXT,
		battery_condition TEXT,
		battery_is_charging BOOLEAN NOT NULL,
		battery_power_adapter TEXT,
		
		status TEXT NOT NULL,
		duration REAL NOT NULL,
		timestamp DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_serial_number ON diagnostics(serial_number);
	CREATE INDEX IF NOT EXISTS idx_created_at ON diagnostics(created_at);
	CREATE INDEX IF NOT EXISTS idx_status ON diagnostics(status);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Tables créées ou déjà existantes")
	return nil
}

// CreateDiagnostic insère un nouveau diagnostic dans la base de données
func CreateDiagnostic(diag models.DiagnosticRequest) (int64, error) {
	query := `
	INSERT INTO diagnostics (
		machine_name, serial_number, model, os_version, macos_version,
		cpu_model, cpu_cores, cpu_frequency, cpu_temperature,
		ram_total, ram_used, ram_available, ram_type,
		storage_type, storage_capacity, storage_used, storage_available, storage_health, storage_device_name,
		battery_cycle_count, battery_health, battery_capacity, battery_max_capacity, 
		battery_condition, battery_is_charging, battery_power_adapter,
		status, duration, timestamp
	) VALUES (
		?, ?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?,
		?, ?, ?
	)
	`

	result, err := DB.Exec(query,
		diag.SystemInfo.MachineName, diag.SystemInfo.SerialNumber, diag.SystemInfo.Model,
		diag.SystemInfo.OSVersion, diag.SystemInfo.MacOSVersion,
		diag.CPU.Model, diag.CPU.Cores, diag.CPU.Frequency, diag.CPU.Temperature,
		diag.RAM.Total, diag.RAM.Used, diag.RAM.Available, diag.RAM.Type,
		diag.Storage.Type, diag.Storage.Capacity, diag.Storage.Used, diag.Storage.Available,
		diag.Storage.Health, diag.Storage.DeviceName,
		diag.Battery.CycleCount, diag.Battery.Health, diag.Battery.Capacity,
		diag.Battery.MaxCapacity, diag.Battery.Condition, diag.Battery.IsCharging,
		diag.Battery.PowerAdapter,
		diag.Status, diag.Duration, time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("Diagnostic créé avec l'ID: %d (Machine: %s)", id, diag.SystemInfo.MachineName)
	return id, nil
}

// GetAllDiagnostics récupère tous les diagnostics
func GetAllDiagnostics(limit int) ([]models.Diagnostic, error) {
	query := `
	SELECT 
		id, machine_name, serial_number, model, os_version, macos_version,
		cpu_model, cpu_cores, cpu_frequency, cpu_temperature,
		ram_total, ram_used, ram_available, ram_type,
		storage_type, storage_capacity, storage_used, storage_available, storage_health, storage_device_name,
		battery_cycle_count, battery_health, battery_capacity, battery_max_capacity, 
		battery_condition, battery_is_charging, battery_power_adapter,
		status, duration, timestamp, created_at
	FROM diagnostics
	ORDER BY created_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diagnostics []models.Diagnostic

	for rows.Next() {
		var d models.Diagnostic
		var macosVersion, cpuTemp, ramType, storageHealth, storageDevice sql.NullString
		var batteryMaxCapacity, batteryCondition, batteryPowerAdapter sql.NullString

		err := rows.Scan(
			&d.ID, &d.SystemInfo.MachineName, &d.SystemInfo.SerialNumber, &d.SystemInfo.Model,
			&d.SystemInfo.OSVersion, &macosVersion,
			&d.CPU.Model, &d.CPU.Cores, &d.CPU.Frequency, &cpuTemp,
			&d.RAM.Total, &d.RAM.Used, &d.RAM.Available, &ramType,
			&d.Storage.Type, &d.Storage.Capacity, &d.Storage.Used, &d.Storage.Available,
			&storageHealth, &storageDevice,
			&d.Battery.CycleCount, &d.Battery.Health, &d.Battery.Capacity,
			&batteryMaxCapacity, &batteryCondition, &d.Battery.IsCharging, &batteryPowerAdapter,
			&d.Status, &d.Duration, &d.Timestamp, &d.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Gérer les valeurs NULL
		if macosVersion.Valid {
			d.SystemInfo.MacOSVersion = macosVersion.String
		}
		if cpuTemp.Valid {
			d.CPU.Temperature = cpuTemp.String
		}
		if ramType.Valid {
			d.RAM.Type = ramType.String
		}
		if storageHealth.Valid {
			d.Storage.Health = storageHealth.String
		}
		if storageDevice.Valid {
			d.Storage.DeviceName = storageDevice.String
		}
		if batteryMaxCapacity.Valid {
			d.Battery.MaxCapacity = batteryMaxCapacity.String
		}
		if batteryCondition.Valid {
			d.Battery.Condition = batteryCondition.String
		}
		if batteryPowerAdapter.Valid {
			d.Battery.PowerAdapter = batteryPowerAdapter.String
		}

		diagnostics = append(diagnostics, d)
	}

	return diagnostics, nil
}

// GetDiagnosticByID récupère un diagnostic par son ID
func GetDiagnosticByID(id int64) (*models.Diagnostic, error) {
	query := `
	SELECT 
		id, machine_name, serial_number, model, os_version, macos_version,
		cpu_model, cpu_cores, cpu_frequency, cpu_temperature,
		ram_total, ram_used, ram_available, ram_type,
		storage_type, storage_capacity, storage_used, storage_available, storage_health, storage_device_name,
		battery_cycle_count, battery_health, battery_capacity, battery_max_capacity, 
		battery_condition, battery_is_charging, battery_power_adapter,
		status, duration, timestamp, created_at
	FROM diagnostics
	WHERE id = ?
	`

	var d models.Diagnostic
	var macosVersion, cpuTemp, ramType, storageHealth, storageDevice sql.NullString
	var batteryMaxCapacity, batteryCondition, batteryPowerAdapter sql.NullString

	err := DB.QueryRow(query, id).Scan(
		&d.ID, &d.SystemInfo.MachineName, &d.SystemInfo.SerialNumber, &d.SystemInfo.Model,
		&d.SystemInfo.OSVersion, &macosVersion,
		&d.CPU.Model, &d.CPU.Cores, &d.CPU.Frequency, &cpuTemp,
		&d.RAM.Total, &d.RAM.Used, &d.RAM.Available, &ramType,
		&d.Storage.Type, &d.Storage.Capacity, &d.Storage.Used, &d.Storage.Available,
		&storageHealth, &storageDevice,
		&d.Battery.CycleCount, &d.Battery.Health, &d.Battery.Capacity,
		&batteryMaxCapacity, &batteryCondition, &d.Battery.IsCharging, &batteryPowerAdapter,
		&d.Status, &d.Duration, &d.Timestamp, &d.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("diagnostic non trouvé")
	}
	if err != nil {
		return nil, err
	}

	// Gérer les valeurs NULL
	if macosVersion.Valid {
		d.SystemInfo.MacOSVersion = macosVersion.String
	}
	if cpuTemp.Valid {
		d.CPU.Temperature = cpuTemp.String
	}
	if ramType.Valid {
		d.RAM.Type = ramType.String
	}
	if storageHealth.Valid {
		d.Storage.Health = storageHealth.String
	}
	if storageDevice.Valid {
		d.Storage.DeviceName = storageDevice.String
	}
	if batteryMaxCapacity.Valid {
		d.Battery.MaxCapacity = batteryMaxCapacity.String
	}
	if batteryCondition.Valid {
		d.Battery.Condition = batteryCondition.String
	}
	if batteryPowerAdapter.Valid {
		d.Battery.PowerAdapter = batteryPowerAdapter.String
	}

	return &d, nil
}

// GetDiagnosticsBySerialNumber récupère tous les diagnostics d'une machine
func GetDiagnosticsBySerialNumber(serialNumber string) ([]models.Diagnostic, error) {
	query := `
	SELECT 
		id, machine_name, serial_number, model, os_version, macos_version,
		cpu_model, cpu_cores, cpu_frequency, cpu_temperature,
		ram_total, ram_used, ram_available, ram_type,
		storage_type, storage_capacity, storage_used, storage_available, storage_health, storage_device_name,
		battery_cycle_count, battery_health, battery_capacity, battery_max_capacity, 
		battery_condition, battery_is_charging, battery_power_adapter,
		status, duration, timestamp, created_at
	FROM diagnostics
	WHERE serial_number = ?
	ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, serialNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diagnostics []models.Diagnostic

	for rows.Next() {
		var d models.Diagnostic
		var macosVersion, cpuTemp, ramType, storageHealth, storageDevice sql.NullString
		var batteryMaxCapacity, batteryCondition, batteryPowerAdapter sql.NullString

		err := rows.Scan(
			&d.ID, &d.SystemInfo.MachineName, &d.SystemInfo.SerialNumber, &d.SystemInfo.Model,
			&d.SystemInfo.OSVersion, &macosVersion,
			&d.CPU.Model, &d.CPU.Cores, &d.CPU.Frequency, &cpuTemp,
			&d.RAM.Total, &d.RAM.Used, &d.RAM.Available, &ramType,
			&d.Storage.Type, &d.Storage.Capacity, &d.Storage.Used, &d.Storage.Available,
			&storageHealth, &storageDevice,
			&d.Battery.CycleCount, &d.Battery.Health, &d.Battery.Capacity,
			&batteryMaxCapacity, &batteryCondition, &d.Battery.IsCharging, &batteryPowerAdapter,
			&d.Status, &d.Duration, &d.Timestamp, &d.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Gérer les valeurs NULL
		if macosVersion.Valid {
			d.SystemInfo.MacOSVersion = macosVersion.String
		}
		if cpuTemp.Valid {
			d.CPU.Temperature = cpuTemp.String
		}
		if ramType.Valid {
			d.RAM.Type = ramType.String
		}
		if storageHealth.Valid {
			d.Storage.Health = storageHealth.String
		}
		if storageDevice.Valid {
			d.Storage.DeviceName = storageDevice.String
		}
		if batteryMaxCapacity.Valid {
			d.Battery.MaxCapacity = batteryMaxCapacity.String
		}
		if batteryCondition.Valid {
			d.Battery.Condition = batteryCondition.String
		}
		if batteryPowerAdapter.Valid {
			d.Battery.PowerAdapter = batteryPowerAdapter.String
		}

		diagnostics = append(diagnostics, d)
	}

	return diagnostics, nil
}

// GetStatistics récupère des statistiques générales
func GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Nombre total de diagnostics
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM diagnostics").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_diagnostics"] = total

	// Nombre de machines uniques
	var uniqueMachines int
	err = DB.QueryRow("SELECT COUNT(DISTINCT serial_number) FROM diagnostics").Scan(&uniqueMachines)
	if err != nil {
		return nil, err
	}
	stats["unique_machines"] = uniqueMachines

	// Répartition par statut
	rows, err := DB.Query("SELECT status, COUNT(*) as count FROM diagnostics GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusCounts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		statusCounts[status] = count
	}
	stats["status_distribution"] = statusCounts

	// Dernier diagnostic
	var lastDiag sql.NullTime
	err = DB.QueryRow("SELECT MAX(created_at) FROM diagnostics").Scan(&lastDiag)
	if err != nil {
		return nil, err
	}
	if lastDiag.Valid {
		stats["last_diagnostic"] = lastDiag.Time
	}

	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	log.Printf("📊 Statistiques: %s", string(statsJSON))

	return stats, nil
}

// CloseDB ferme la connexion à la base de données
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
