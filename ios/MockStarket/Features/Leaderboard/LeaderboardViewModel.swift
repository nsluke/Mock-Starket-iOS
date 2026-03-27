import SwiftUI
import Observation

@Observable
final class LeaderboardViewModel {
    var entries: [LeaderboardEntry] = []
    var selectedPeriod = "alltime"
    var isLoading = false

    private let apiClient = APIClient.shared

    func load() async {
        isLoading = true
        defer { isLoading = false }

        do {
            entries = try await apiClient.request(
                .getLeaderboard(period: selectedPeriod),
                token: AuthManager.shared.currentToken
            )
        } catch {
            // Handle error
        }
    }
}
