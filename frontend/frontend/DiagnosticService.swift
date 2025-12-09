//
//  DiagnosticService.swift
//  frontend
//
//  Created by Mathis Bossard on 09/12/2025.
//

import Foundation
import IOKit
import IOKit.ps
import Combine

class DiagnosticService: ObservableObject {
    
    // MARK: - Published Properties (pour l'UI)
    @Published var currentStep: DiagnosticStep = .cpu
    @Published var progress: Double = 0.0
    @Published var isRunning: Bool = false
    @Published var diagnosticResult: DiagnosticData?
    @Published var errorMessage: String?
    
    // MARK: - Nom de la machine
    private func getMachineName() -> String {
        return Host.current().localizedName ?? "Unknown Mac"
    }
    
    // MARK: - Lancer le diagnostic complet
    func runDiagnostic() async {
        await MainActor.run {
            isRunning = true
            progress = 0.0
            errorMessage = nil
        }
        
        let startTime = Date()
        
        do {
            // Étape 1: CPU
            await updateStep(.cpu)
            let cpuInfo = getCPUInfo()
            try await Task.sleep(nanoseconds: 500_000_000) // Simulation 0.5s
            
            // Étape 2: RAM
            await updateStep(.ram)
            let ramInfo = getRAMInfo()
            try await Task.sleep(nanoseconds: 500_000_000)
            
            // Étape 3: Stockage
            await updateStep(.storage)
            let storageInfo = getStorageInfo()
            try await Task.sleep(nanoseconds: 500_000_000)
            
            // Étape 4: Batterie
            await updateStep(.battery)
            let batteryInfo = getBatteryInfo()
            try await Task.sleep(nanoseconds: 500_000_000)
            
            // Étape 5: Informations système
            await updateStep(.system)
            let serialNumber = getSerialNumber()
            try await Task.sleep(nanoseconds: 500_000_000)
            
            // Calculer la durée du test
            let duration = Int(Date().timeIntervalSince(startTime))
            
            // Construire le résultat
            let result = DiagnosticData(
                machineName: getMachineName(),
                serialNumber: serialNumber,
                cpuModel: cpuInfo.model,
                cpuCores: cpuInfo.cores,
                ramTotalGB: ramInfo.total,
                ramUsedGB: ramInfo.used,
                storageTotalGB: storageInfo.total,
                storageUsedGB: storageInfo.used,
                batteryHealth: batteryInfo.health,
                batteryCycleCount: batteryInfo.cycles,
                batteryPercentage: batteryInfo.percentage,
                testDurationSeconds: duration,
                status: "completed"
            )
            
            await MainActor.run {
                self.diagnosticResult = result
                self.progress = 1.0
                self.isRunning = false
            }
            
        } catch {
            await MainActor.run {
                self.errorMessage = "Erreur pendant le diagnostic: \(error.localizedDescription)"
                self.isRunning = false
            }
        }
    }
    
    // MARK: - Mise à jour de l'étape
    private func updateStep(_ step: DiagnosticStep) async {
        await MainActor.run {
            self.currentStep = step
            self.progress = step.progress
        }
    }
    
    // MARK: - Collecte CPU
    private func getCPUInfo() -> (model: String, cores: Int) {
        let processInfo = ProcessInfo.processInfo
        
        // Nombre de cœurs
        let cores = processInfo.processorCount
        
        // Modèle du CPU (via sysctl)
        var size = 0
        sysctlbyname("machdep.cpu.brand_string", nil, &size, nil, 0)
        var machine = [CChar](repeating: 0, count: size)
        sysctlbyname("machdep.cpu.brand_string", &machine, &size, nil, 0)
        let cpuModel = String(cString: machine).trimmingCharacters(in: .whitespacesAndNewlines)
        
        return (model: cpuModel.isEmpty ? "Unknown CPU" : cpuModel, cores: cores)
    }
    
    // MARK: - Collecte RAM
    private func getRAMInfo() -> (total: Double, used: Double) {
        let processInfo = ProcessInfo.processInfo
        let totalRAM = Double(processInfo.physicalMemory) / 1_073_741_824.0 // Conversion en GB
        
        // RAM utilisée (approximation via vm_stat)
        var vmStats = vm_statistics64()
        var count = mach_msg_type_number_t(MemoryLayout<vm_statistics64>.size / MemoryLayout<integer_t>.size)
        
        let result = withUnsafeMutablePointer(to: &vmStats) {
            $0.withMemoryRebound(to: integer_t.self, capacity: Int(count)) {
                host_statistics64(mach_host_self(), HOST_VM_INFO64, $0, &count)
            }
        }
        
        var usedRAM = totalRAM * 0.5 // Valeur par défaut
        
        if result == KERN_SUCCESS {
            let pageSize = Double(vm_kernel_page_size)
            let active = Double(vmStats.active_count) * pageSize
            let wired = Double(vmStats.wire_count) * pageSize
            let compressed = Double(vmStats.compressor_page_count) * pageSize
            usedRAM = (active + wired + compressed) / 1_073_741_824.0
        }
        
        return (total: totalRAM, used: usedRAM)
    }
    
    // MARK: - Collecte Stockage
    private func getStorageInfo() -> (total: Double, used: Double) {
        let fileURL = URL(fileURLWithPath: "/")
        
        do {
            let values = try fileURL.resourceValues(forKeys: [
                .volumeTotalCapacityKey,
                .volumeAvailableCapacityKey
            ])
            
            let totalSpace = Double(values.volumeTotalCapacity ?? 0) / 1_073_741_824.0
            let availableSpace = Double(values.volumeAvailableCapacity ?? 0) / 1_073_741_824.0
            let usedSpace = totalSpace - availableSpace
            
            return (total: totalSpace, used: usedSpace)
        } catch {
            print("Erreur collecte stockage: \(error)")
            return (total: 0, used: 0)
        }
    }
    
    // MARK: - Collecte Batterie
    private func getBatteryInfo() -> (health: String, cycles: Int, percentage: Int) {
        let snapshot = IOPSCopyPowerSourcesInfo()?.takeRetainedValue()
        let sources = IOPSCopyPowerSourcesList(snapshot)?.takeRetainedValue() as? [CFTypeRef]
        
        guard let sources = sources, !sources.isEmpty else {
            return (health: "No Battery", cycles: 0, percentage: 0)
        }
        
        guard let powerSource = sources.first,
              let info = IOPSGetPowerSourceDescription(snapshot, powerSource)?.takeUnretainedValue() as? [String: Any] else {
            return (health: "Unknown", cycles: 0, percentage: 0)
        }
        
        // Pourcentage de batterie
        let percentage = info[kIOPSCurrentCapacityKey] as? Int ?? 0
        
        // Cycles de charge
        let cycles = info["Cycle Count"] as? Int ?? 0
        
        // État de santé
        let maxCapacity = info[kIOPSMaxCapacityKey] as? Int ?? 100
        let health: String
        if maxCapacity >= 80 {
            health = "Good"
        } else if maxCapacity >= 60 {
            health = "Fair"
        } else {
            health = "Poor"
        }
        
        return (health: health, cycles: cycles, percentage: percentage)
    }
    
    // MARK: - Numéro de série
    private func getSerialNumber() -> String {
        let platformExpert = IOServiceGetMatchingService(
            kIOMainPortDefault,
            IOServiceMatching("IOPlatformExpertDevice")
        )
        
        guard platformExpert > 0 else {
            return "UNKNOWN_SERIAL"
        }
        
        defer { IOObjectRelease(platformExpert) }
        
        guard let serialNumberAsCFString = IORegistryEntryCreateCFProperty(
            platformExpert,
            kIOPlatformSerialNumberKey as CFString,
            kCFAllocatorDefault,
            0
        )?.takeUnretainedValue() as? String else {
            return "UNKNOWN_SERIAL"
        }
        
        return serialNumberAsCFString
    }
}
