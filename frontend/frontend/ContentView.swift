import SwiftUI

struct ContentView: View {
    
    // MARK: - State Objects
    @StateObject private var diagnosticService = DiagnosticService()
    
    // MARK: - State Variables
    @State private var showResults = false
    
    var body: some View {
        ZStack {
            // Background gradient
            LinearGradient(
                gradient: Gradient(colors: [Color.blue.opacity(0.1), Color.purple.opacity(0.1)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()
            
            if diagnosticService.isRunning {
                // Vue pendant le test avec sidebar
                runningTestView
            } else if showResults, let result = diagnosticService.diagnosticResult {
                // Vue Dashboard avec résultats
                dashboardView(result: result)
            } else {
                // Vue initiale
                VStack(spacing: 30) {
                    // MARK: - Header
                    VStack(spacing: 10) {
                        Image(systemName: "desktopcomputer")
                            .font(.system(size: 60))
                            .foregroundColor(.blue)
                        
                        Text("Diagnostic Hardware")
                            .font(.largeTitle)
                            .fontWeight(.bold)
                        
                        Text("Test automatique des composants Mac")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.top, 40)
                    
                    Spacer()
                    
                    initialView
                    
                    Spacer()
                    
                    // MARK: - Error Message
                    if let error = diagnosticService.errorMessage {
                        Text(error)
                            .foregroundColor(.red)
                            .font(.caption)
                            .padding()
                            .background(Color.red.opacity(0.1))
                            .cornerRadius(8)
                            .padding(.horizontal)
                    }
                    
                    // MARK: - Action Button
                    actionButton
                        .padding(.bottom, 40)
                }
                .frame(maxWidth: 600)
            }
        }
        .frame(minWidth: 900, minHeight: 650)
    }
    
    // MARK: - Initial View
    private var initialView: some View {
        VStack(spacing: 20) {
            Image(systemName: "waveform.path.ecg")
                .font(.system(size: 80))
                .foregroundColor(.blue.opacity(0.5))
            
            Text("Prêt à démarrer le diagnostic")
                .font(.title2)
                .fontWeight(.semibold)
        }
        .padding()
    }
    
    // MARK: - Running Test View avec Sidebar
    private var runningTestView: some View {
        HStack(spacing: 0) {
            // Sidebar de progression
            VStack(spacing: 20) {
                Text("Progression")
                    .font(.headline)
                    .padding(.top, 30)
                
                ForEach(DiagnosticStep.allCases, id: \.rawValue) { step in
                    HStack {
                        Image(systemName: diagnosticService.currentStep == step ? "circle.fill" :
                                         diagnosticService.progress >= step.progress ? "checkmark.circle.fill" : "circle")
                            .foregroundColor(diagnosticService.currentStep == step ? .blue :
                                           diagnosticService.progress >= step.progress ? .green : .gray)
                        
                        Text(step.rawValue)
                            .font(.subheadline)
                            .foregroundColor(diagnosticService.currentStep == step ? .primary : .secondary)
                        
                        Spacer()
                    }
                    .padding(.horizontal, 20)
                }
                
                Spacer()
                
                // Barre de progression globale
                VStack(spacing: 10) {
                    Text("\(Int(diagnosticService.progress * 100))%")
                        .font(.title)
                        .fontWeight(.bold)
                        .foregroundColor(.blue)
                    
                    ProgressView(value: diagnosticService.progress)
                        .progressViewStyle(LinearProgressViewStyle(tint: .blue))
                        .frame(height: 8)
                        .scaleEffect(x: 1, y: 1.5, anchor: .center)
                }
                .padding(20)
            }
            .frame(width: 280)
            .background(Color.white.opacity(0.3))
            
            // Zone principale
            VStack {
                Spacer()
                
                ProgressView()
                    .scaleEffect(2.5)
                    .padding()
                
                Text(diagnosticService.currentStep.rawValue)
                    .font(.title)
                    .fontWeight(.semibold)
                    .padding()
                
                Text("Test en cours...")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                
                Spacer()
            }
            .frame(maxWidth: .infinity)
        }
    }
    
    // MARK: - Dashboard View
    private func dashboardView(result: DiagnosticData) -> some View {
        HStack(spacing: 0) {
            // Sidebar de progression (complétée)
            VStack(spacing: 20) {
                Text("Tests effectués")
                    .font(.headline)
                    .padding(.top, 30)
                
                ForEach(DiagnosticStep.allCases, id: \.rawValue) { step in
                    HStack {
                        Image(systemName: "checkmark.circle.fill")
                            .foregroundColor(.green)
                        
                        Text(step.rawValue)
                            .font(.subheadline)
                        
                        Spacer()
                    }
                    .padding(.horizontal, 20)
                }
                
                Spacer()
                
                // Durée du test
                VStack(spacing: 5) {
                    Text("Durée totale")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    Text("\(result.testDurationSeconds)s")
                        .font(.title2)
                        .fontWeight(.bold)
                        .foregroundColor(.blue)
                }
                .padding(20)
            }
            .frame(width: 280)
            .background(Color.white.opacity(0.3))
            
            // Zone principale - Dashboard
            ScrollView {
                VStack(spacing: 25) {
                    // Header avec titre et numéro de série
                    VStack(spacing: 10) {
                        HStack {
                            Image(systemName: "checkmark.circle.fill")
                                .font(.system(size: 40))
                                .foregroundColor(.green)
                            
                            VStack(alignment: .leading, spacing: 5) {
                                Text("Test terminé pour Mac")
                                    .font(.title)
                                    .fontWeight(.bold)
                                
                                Text(result.serialNumber)
                                    .font(.title3)
                                    .foregroundColor(.secondary)
                            }
                            
                            Spacer()
                        }
                        .padding(.horizontal, 30)
                        .padding(.top, 30)
                    }
                    
                    Divider()
                        .padding(.horizontal, 30)
                    
                    // Grille de résultats
                    VStack(spacing: 20) {
                        // Ligne 1: CPU et RAM
                        HStack(spacing: 20) {
                            DashboardCard(
                                icon: "cpu",
                                title: "Processeur",
                                value: result.cpuModel,
                                detail: "\(result.cpuCores) cœurs",
                                color: .blue
                            )
                            
                            DashboardCard(
                                icon: "memorychip",
                                title: "Mémoire RAM",
                                value: String(format: "%.1f GB / %.1f GB", result.ramUsedGB, result.ramTotalGB),
                                detail: String(format: "%.0f%% utilisé", (result.ramUsedGB / result.ramTotalGB) * 100),
                                color: .purple
                            )
                        }
                        
                        // Ligne 2: Stockage et Batterie
                        HStack(spacing: 20) {
                            DashboardCard(
                                icon: "internaldrive",
                                title: "Stockage",
                                value: String(format: "%.0f GB / %.0f GB", result.storageUsedGB, result.storageTotalGB),
                                detail: String(format: "%.0f%% utilisé", (result.storageUsedGB / result.storageTotalGB) * 100),
                                color: .orange
                            )
                            
                            DashboardCard(
                                icon: "battery.100",
                                title: "Batterie",
                                value: "\(result.batteryPercentage)%",
                                detail: "\(result.batteryHealth) • \(result.batteryCycleCount) cycles",
                                color: batteryColor(result.batteryHealth)
                            )
                        }
                    }
                    .padding(.horizontal, 30)
                    
                    // Bouton nouveau test
                    Button(action: {
                        showResults = false
                        diagnosticService.diagnosticResult = nil
                        diagnosticService.errorMessage = nil
                    }) {
                        HStack {
                            Image(systemName: "arrow.counterclockwise")
                            Text("Nouveau Test")
                        }
                        .font(.headline)
                        .foregroundColor(.white)
                        .frame(width: 250, height: 50)
                        .background(Color.blue)
                        .cornerRadius(12)
                        .shadow(radius: 5)
                    }
                    .buttonStyle(PlainButtonStyle())
                    .padding(.top, 20)
                    .padding(.bottom, 30)
                }
            }
            .frame(maxWidth: .infinity)
        }
    }
    
    // MARK: - Action Button
    private var actionButton: some View {
        Button(action: {
            launchDiagnostic()
        }) {
            HStack {
                Image(systemName: "play.fill")
                Text("Lancer le Test")
            }
            .font(.headline)
            .foregroundColor(.white)
            .frame(width: 250, height: 50)
            .background(Color.blue)
            .cornerRadius(12)
            .shadow(radius: 5)
        }
        .buttonStyle(PlainButtonStyle())
        .disabled(diagnosticService.isRunning)
    }
    
    // MARK: - Helper Functions
    
    private func batteryColor(_ health: String) -> Color {
        switch health {
        case "Good": return .green
        case "Fair": return .orange
        default: return .red
        }
    }
    
    private func launchDiagnostic() {
        Task {
            // Lancer le diagnostic
            await diagnosticService.runDiagnostic()
            
            // Si succès, envoyer au backend
            if let result = diagnosticService.diagnosticResult {
                do {
                    let response = try await APIClient.shared.sendDiagnostic(result)
                    print("Diagnostic enregistré - ID: \(response.id)")
                    
                    await MainActor.run {
                        showResults = true
                    }
                } catch {
                    print("Erreur d'envoi: \(error)")
                    await MainActor.run {
                        diagnosticService.errorMessage = "Erreur d'envoi: \(error.localizedDescription)"
                        showResults = true // Afficher quand même les résultats
                    }
                }
            }
        }
    }
}

// MARK: - Dashboard Card Component
struct DashboardCard: View {
    let icon: String
    let title: String
    let value: String
    let detail: String
    let color: Color
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Image(systemName: icon)
                    .font(.title2)
                    .foregroundColor(color)
                
                Text(title)
                    .font(.headline)
                    .foregroundColor(.secondary)
                
                Spacer()
            }
            
            Text(value)
                .font(.title3)
                .fontWeight(.semibold)
                .lineLimit(2)
                .minimumScaleFactor(0.8)
            
            Text(detail)
                .font(.caption)
                .foregroundColor(.secondary)
        }
        .padding(20)
        .frame(maxWidth: .infinity, minHeight: 120)
        .background(Color.white.opacity(0.8))
        .cornerRadius(15)
        .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
    }
}

// MARK: - Preview
#Preview {
    ContentView()
}
