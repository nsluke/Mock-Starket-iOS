import SwiftUI
import Observation

@Observable
final class PortfolioViewModel {
    var portfolioResponse: PortfolioResponse?
    var positions: [Position] { portfolioResponse?.positions ?? [] }
    var isLoading = false

    private let apiClient = APIClient.shared

    func load() async {
        isLoading = true
        defer { isLoading = false }

        do {
            portfolioResponse = try await apiClient.request(
                .getPortfolio,
                token: AuthManager.shared.currentToken
            )
        } catch {
            // Handle error
        }
    }
}
