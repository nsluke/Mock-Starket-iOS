import SwiftUI
import Observation

@MainActor @Observable
final class AlertsViewModel {
    var alerts: [PriceAlert] = []
    var isLoading = false
    var errorMessage: String?

    // Create alert form
    var ticker = ""
    var condition = "above"
    var targetPrice = ""
    var isCreating = false

    private let apiClient = APIClient.shared
    private let authManager = AuthManager.shared

    func load() async {
        isLoading = true
        defer { isLoading = false }

        guard let token = authManager.currentToken else { return }

        do {
            alerts = try await apiClient.request(.listAlerts, token: token)
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func createAlert() async -> Bool {
        guard !ticker.isEmpty,
              let price = Decimal(string: targetPrice),
              price > 0 else {
            errorMessage = "Enter a valid ticker and price"
            return false
        }

        isCreating = true
        defer { isCreating = false }

        guard let token = authManager.currentToken else { return false }

        struct AlertRequest: Encodable {
            let ticker: String
            let condition: String
            let target_price: Decimal
        }

        do {
            let _: PriceAlert = try await apiClient.request(
                .createAlert,
                token: token,
                body: AlertRequest(ticker: ticker, condition: condition, target_price: price)
            )
            ticker = ""
            targetPrice = ""
            await load()
            return true
        } catch {
            errorMessage = error.localizedDescription
            return false
        }
    }

    func deleteAlert(_ alert: PriceAlert) async {
        guard let token = authManager.currentToken else { return }

        struct StatusResponse: Decodable { let status: String }

        do {
            let _: StatusResponse = try await apiClient.request(.deleteAlert(id: alert.id), token: token)
            alerts.removeAll { $0.id == alert.id }
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
