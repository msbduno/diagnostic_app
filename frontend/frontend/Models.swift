//
//  Models.swift
//  frontend
//
//  Created by Mathis Bossard on 09/12/2025.
//

import Foundation

// MARK: - Diagnostic Data (à envoyer au backend)
struct DiagnosticData: Codable {
    let machineName: String
    let serialNumber: String
    let cpuModel: String
    let cpuCores: Int
    let ramTotalGB: Double
    let ramUsedGB: Double
    let storageTotalGB: Double
    let storageUsedGB: Double
    let batteryHealth: String
    let batteryCycleCount: Int
    let batteryPercentage: Int
    let testDurationSeconds: Int
    let status: String
    
    enum CodingKeys: String, CodingKey {
        case machineName = "machine_name"
        case serialNumber = "serial_number"
        case cpuModel = "cpu_model"
        case cpuCores = "cpu_cores"
        case ramTotalGB = "ram_total_gb"
        case ramUsedGB = "ram_used_gb"
        case storageTotalGB = "storage_total_gb"
        case storageUsedGB = "storage_used_gb"
        case batteryHealth = "battery_health"
        case batteryCycleCount = "battery_cycle_count"
        case batteryPercentage = "battery_percentage"
        case testDurationSeconds = "test_duration_seconds"
        case status
    }
}

// MARK: - Réponse du backend après POST
struct DiagnosticResponse: Codable {
    let id: Int
    let serialNumber: String
    let timestamp: String
    let status: String
    
    enum CodingKeys: String, CodingKey {
        case id
        case serialNumber = "serial_number"
        case timestamp
        case status
    }
}

// MARK: - Diagnostic complet (réponse GET)
struct DiagnosticRecord: Codable, Identifiable {
    let id: Int
    let serialNumber: String
    let timestamp: String
    let cpuModel: String
    let cpuCores: Int
    let ramTotalGB: Double
    let ramUsedGB: Double
    let storageTotalGB: Double
    let storageUsedGB: Double
    let batteryHealth: String
    let batteryCycleCount: Int
    let batteryPercentage: Int
    let testDurationSeconds: Int
    let status: String
    
    enum CodingKeys: String, CodingKey {
        case id
        case serialNumber = "serial_number"
        case timestamp
        case cpuModel = "cpu_model"
        case cpuCores = "cpu_cores"
        case ramTotalGB = "ram_total_gb"
        case ramUsedGB = "ram_used_gb"
        case storageTotalGB = "storage_total_gb"
        case storageUsedGB = "storage_used_gb"
        case batteryHealth = "battery_health"
        case batteryCycleCount = "battery_cycle_count"
        case batteryPercentage = "battery_percentage"
        case testDurationSeconds = "test_duration_seconds"
        case status
    }
}

// MARK: - Statistiques (réponse GET /statistics)
struct Statistics: Codable {
    let totalTests: Int
    let completedTests: Int
    let failedTests: Int
    let averageTestDuration: Double
    let uniqueMachines: Int
    
    enum CodingKeys: String, CodingKey {
        case totalTests = "total_tests"
        case completedTests = "completed_tests"
        case failedTests = "failed_tests"
        case averageTestDuration = "average_test_duration"
        case uniqueMachines = "unique_machines"
    }
}

// MARK: - État du test (pour l'UI)
enum TestStatus {
    case idle           // Pas encore lancé
    case running        // En cours
    case completed      // Terminé avec succès
    case failed         // Échoué
}

// MARK: - Étape du diagnostic (pour la progression)
enum DiagnosticStep: String, CaseIterable {
    case cpu = "Test CPU"
    case ram = "Test RAM"
    case storage = "Test Stockage"
    case battery = "Test Batterie"
    case system = "Informations Système"
    case upload = "Envoi au serveur"
    
    var progress: Double {
        switch self {
        case .cpu: return 0.15
        case .ram: return 0.30
        case .storage: return 0.50
        case .battery: return 0.70
        case .system: return 0.85
        case .upload: return 1.0
        }
    }
}

// MARK: - Erreurs API
enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(Int)
    case serverError(String)
    case encodingError
    case decodingError
    case networkError
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "URL invalide"
        case .invalidResponse:
            return "Réponse du serveur invalide"
        case .httpError(let code):
            return "Erreur HTTP \(code)"
        case .serverError(let message):
            return "Erreur serveur: \(message)"
        case .encodingError:
            return "Erreur d'encodage des données"
        case .decodingError:
            return "Erreur de décodage de la réponse"
        case .networkError:
            return "Erreur réseau - Vérifiez que le backend est démarré"
        }
    }
}

// MARK: - Erreur Backend
struct BackendError: Codable {
    let error: String
    let details: String?
}
