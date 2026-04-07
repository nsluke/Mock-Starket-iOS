import SwiftUI
import Observation

@MainActor @Observable
final class OptionsChainViewModel {
    var expirations: [Date] = []
    var selectedExpiration: Date?
    var chain: OptionChainResponse?
    var isLoading = false
    var selectedContract: OptionContract?
    var showTradeSheet = false

    private let apiClient = APIClient.shared

    func load(ticker: String) async {
        isLoading = true
        defer { isLoading = false }

        do {
            let token = AuthManager.shared.currentToken
            let dates: [Date] = try await apiClient.request(
                .getOptionExpirations(ticker: ticker),
                token: token
            )
            expirations = dates
            if selectedExpiration == nil, let first = dates.first {
                selectedExpiration = first
            }
            await loadChain(ticker: ticker)
        } catch {
            // Silently handle — options may not be available for all assets
        }
    }

    func selectExpiration(_ date: Date, ticker: String) async {
        selectedExpiration = date
        await loadChain(ticker: ticker)
    }

    private func loadChain(ticker: String) async {
        guard let exp = selectedExpiration else { return }

        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime]
        let expStr = formatter.string(from: exp)

        do {
            let token = AuthManager.shared.currentToken
            chain = try await apiClient.request(
                .getOptionChain(ticker: ticker, expiration: expStr),
                token: token
            )
        } catch {
            // Handle silently
        }
    }

    func openTrade(contract: OptionContract) {
        selectedContract = contract
        showTradeSheet = true
    }
}
