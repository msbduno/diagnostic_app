//
//  APIClient.swift
//  frontend
//
//  Created by Mathis Bossard on 09/12/2025.
//

import Foundation

class APIClient {
    
    // MARK: - Configuration
    static let shared = APIClient()
    
    private let baseURL = "http://localhost:8080/api/v1"
    private let session: URLSession
    
    private init() {
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        self.session = URLSession(configuration: configuration)
    }
    
    // MARK: - Envoyer un diagnostic (POST)
    func sendDiagnostic(_ diagnostic: DiagnosticData) async throws -> DiagnosticResponse {
        guard let url = URL(string: "\(baseURL)/diagnostics") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        // Encoder les données
        let encoder = JSONEncoder()
        encoder.keyEncodingStrategy = .convertToSnakeCase
        request.httpBody = try encoder.encode(diagnostic)
        
        // Log de la requête (debug)
        if let jsonString = String(data: request.httpBody!, encoding: .utf8) {
            print("Envoi diagnostic:")
            print(jsonString)
        }
        
        // Envoyer la requête
        let (data, response) = try await session.data(for: request)
        
        // Vérifier le status code
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }
        
        print("Status code: \(httpResponse.statusCode)")
        
        guard (200...299).contains(httpResponse.statusCode) else {
            // Tenter de décoder l'erreur du backend
            if let errorResponse = try? JSONDecoder().decode(BackendError.self, from: data) {
                throw APIError.serverError(errorResponse.error)
            }
            throw APIError.httpError(httpResponse.statusCode)
        }
        
        // Décoder la réponse
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let diagnosticResponse = try decoder.decode(DiagnosticResponse.self, from: data)
        
        print("Diagnostic envoyé avec succès - ID: \(diagnosticResponse.id)")
        
        return diagnosticResponse
    }
    
    // MARK: - Récupérer tous les diagnostics (GET)
    func getAllDiagnostics() async throws -> [DiagnosticRecord] {
        guard let url = URL(string: "\(baseURL)/diagnostics") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError.invalidResponse
        }
        
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let diagnostics = try decoder.decode([DiagnosticRecord].self, from: data)
        
        print("📋 \(diagnostics.count) diagnostics récupérés")
        
        return diagnostics
    }
    
    // MARK: - Récupérer diagnostics par numéro de série (GET)
    func getDiagnosticsBySerial(_ serialNumber: String) async throws -> [DiagnosticRecord] {
        guard let url = URL(string: "\(baseURL)/diagnostics/serial/\(serialNumber)") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError.invalidResponse
        }
        
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let diagnostics = try decoder.decode([DiagnosticRecord].self, from: data)
        
        print("\(diagnostics.count) diagnostics trouvés pour \(serialNumber)")
        
        return diagnostics
    }
    
    // MARK: - Récupérer les statistiques (GET)
    func getStatistics() async throws -> Statistics {
        guard let url = URL(string: "\(baseURL)/statistics") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError.invalidResponse
        }
        
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let stats = try decoder.decode(Statistics.self, from: data)
        
        print("Statistiques récupérées")
        
        return stats
    }
    
    // MARK: - Health Check (GET)
    func healthCheck() async throws -> Bool {
        guard let url = URL(string: "\(baseURL)/health") else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.timeoutInterval = 5 // Timeout court pour health check
        
        let (_, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            return false
        }
        
        let isHealthy = (200...299).contains(httpResponse.statusCode)
        print(isHealthy ? "Backend connecté" : "Backend non disponible")
        
        return isHealthy
    }

}
