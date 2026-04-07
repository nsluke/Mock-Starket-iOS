import SwiftUI
import Observation

enum AssetCategory: String, CaseIterable, Identifiable {
    case all = "All"
    case stock = "Stocks"
    case etf = "ETFs"
    case crypto = "Crypto"
    case commodity = "Commodities"

    var id: String { rawValue }

    /// Maps to the `assetType` string on `Stock`.
    var assetTypeKey: String {
        switch self {
        case .all: return ""
        case .stock: return "stock"
        case .etf: return "etf"
        case .crypto: return "crypto"
        case .commodity: return "commodity"
        }
    }
}

@MainActor @Observable
final class MarketViewModel {
    var stocks: [Stock] = []
    var searchText = ""
    var selectedCategory: AssetCategory = .all
    var selectedSector: String?
    var isLoading = false
    var errorMessage: String?
    var marketSummary: MarketSummary?

    private let apiClient = APIClient.shared
    private let wsManager = WebSocketManager.shared

    /// Sectors available for the currently selected asset category.
    var availableSectors: [String] {
        let pool: [Stock]
        if selectedCategory == .all {
            pool = stocks
        } else {
            pool = stocks.filter { $0.assetType == selectedCategory.assetTypeKey }
        }
        let unique = Set(pool.map(\.sector))
        return unique.sorted()
    }

    var filteredStocks: [Stock] {
        var result = stocks

        if selectedCategory != .all {
            result = result.filter { $0.assetType == selectedCategory.assetTypeKey }
        }

        if let sector = selectedSector {
            result = result.filter { $0.sector == sector }
        }

        if !searchText.isEmpty {
            result = result.filter {
                $0.ticker.localizedCaseInsensitiveContains(searchText) ||
                $0.name.localizedCaseInsensitiveContains(searchText)
            }
        }

        return result
    }

    var topGainers: [Stock] {
        stocks.sorted { $0.changePct > $1.changePct }.prefix(5).map { $0 }
    }

    var topLosers: [Stock] {
        stocks.sorted { $0.changePct < $1.changePct }.prefix(5).map { $0 }
    }

    func loadStocks() async {
        isLoading = true
        defer { isLoading = false }

        do {
            let token = AuthManager.shared.currentToken
            stocks = try await apiClient.request(.listStocks, token: token)
            marketSummary = try await apiClient.request(.marketSummary, token: token)
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func subscribeToUpdates() {
        let vm = self
        wsManager.onPriceUpdate { updates in
            Task { @MainActor in
                vm.applyPriceUpdates(updates)
            }
        }
    }

    private func applyPriceUpdates(_ updates: [PriceUpdate]) {
        for update in updates {
            if let index = stocks.firstIndex(where: { $0.ticker == update.ticker }) {
                stocks[index].currentPrice = update.price
                stocks[index].dayHigh = update.high
                stocks[index].dayLow = update.low
                stocks[index].volume = update.volume
            }
        }
    }
}
