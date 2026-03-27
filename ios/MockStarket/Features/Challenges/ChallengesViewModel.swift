import SwiftUI
import Observation

@MainActor @Observable
final class ChallengesViewModel {
    var challenge: DailyChallenge?
    var progress: UserChallenge?
    var isLoading = false
    var isChecking = false
    var isClaiming = false
    var errorMessage: String?

    private let apiClient = APIClient.shared
    private let authManager = AuthManager.shared

    var isCompleted: Bool { progress?.completed ?? false }
    var isClaimed: Bool { progress?.claimed ?? false }

    func load() async {
        isLoading = true
        defer { isLoading = false }

        guard let token = authManager.currentToken else { return }

        do {
            let response: ChallengeResponse = try await apiClient.request(.getTodaysChallenge, token: token)
            challenge = response.challenge
            progress = response.progress
        } catch {
            // No challenge today is not an error
            challenge = nil
            progress = nil
        }
    }

    func checkProgress() async {
        isChecking = true
        defer { isChecking = false }

        guard let token = authManager.currentToken else { return }

        do {
            let response: ChallengeCheckResponse = try await apiClient.request(.checkChallenge, token: token)
            if response.completed {
                await load() // Refresh to get updated progress
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func claimReward() async {
        guard let challenge else { return }

        isClaiming = true
        defer { isClaiming = false }

        guard let token = authManager.currentToken else { return }

        struct StatusResponse: Decodable { let status: String }

        do {
            let _: StatusResponse = try await apiClient.request(.claimChallenge(id: challenge.id), token: token)
            await load()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
