import SwiftUI
import Observation

@MainActor @Observable
final class MarketViewModel {
    var stocks: [Stock] = []
    var searchText = ""
    var isLoading = false
    var errorMessage: String?
    var marketSummary: MarketSummary?

    private let apiClient = APIClient.shared
    private let wsManager = WebSocketManager.shared

    var filteredStocks: [Stock] {
        if searchText.isEmpty { return stocks }
        return stocks.filter {
            $0.ticker.localizedCaseInsensitiveContains(searchText) ||
            $0.name.localizedCaseInsensitiveContains(searchText)
        }
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
