import SwiftUI
import Observation

@MainActor @Observable
final class StockDetailViewModel {
    enum TimeRange: CaseIterable {
        case oneDay, oneWeek, oneMonth, threeMonths, all

        var label: String {
            switch self {
            case .oneDay: return "1D"
            case .oneWeek: return "1W"
            case .oneMonth: return "1M"
            case .threeMonths: return "3M"
            case .all: return "ALL"
            }
        }

        var interval: String {
            switch self {
            case .oneDay: return "1m"
            case .oneWeek: return "5m"
            case .oneMonth: return "1h"
            case .threeMonths, .all: return "1d"
            }
        }
    }

    var stock: Stock?
    var priceHistory: [PricePoint] = []
    var selectedRange: TimeRange = .oneDay
    var showTradeSheet = false
    var tradeSide = "buy"
    var isLoading = false

    private let apiClient = APIClient.shared
    private let wsManager = WebSocketManager.shared

    func load(ticker: String) async {
        isLoading = true
        defer { isLoading = false }

        let token = AuthManager.shared.currentToken

        do {
            stock = try await apiClient.request(.getStock(ticker: ticker), token: token)
            priceHistory = try await apiClient.request(
                .getStockHistory(ticker: ticker, interval: selectedRange.interval),
                token: token
            )
        } catch {
            // Handle error
        }

        // Subscribe to live updates
        let vm = self
        wsManager.onPriceUpdate { updates in
            guard let update = updates.first(where: { $0.ticker == ticker }) else { return }
            Task { @MainActor in
                vm.stock?.currentPrice = update.price
                vm.stock?.dayHigh = update.high
                vm.stock?.dayLow = update.low
                vm.stock?.volume = update.volume
            }
        }
    }

    func selectTimeRange(_ range: TimeRange) async {
        selectedRange = range
        guard let ticker = stock?.ticker else { return }

        let token = AuthManager.shared.currentToken
        do {
            priceHistory = try await apiClient.request(
                .getStockHistory(ticker: ticker, interval: range.interval),
                token: token
            )
        } catch {
            // Handle error
        }
    }
}
