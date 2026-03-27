import SwiftUI
import Observation

@Observable
final class ProfileViewModel {
    var achievements: [Achievement] = []

    private let apiClient = APIClient.shared

    func load() async {
        do {
            achievements = try await apiClient.request(
                .listAchievements,
                token: AuthManager.shared.currentToken
            )
        } catch {
            // Handle error
        }
    }
}
