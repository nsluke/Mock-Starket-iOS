import SwiftUI
import Observation

@MainActor @Observable
final class SettingsViewModel {
    var displayName = ""
    var isSaving = false
    var saveSuccess = false
    var errorMessage: String?

    private let apiClient = APIClient.shared
    private let authManager = AuthManager.shared

    func load(user: User?) {
        displayName = user?.displayName ?? ""
    }

    func updateProfile() async -> Bool {
        guard !displayName.trimmingCharacters(in: .whitespaces).isEmpty else {
            errorMessage = "Display name cannot be empty"
            return false
        }

        isSaving = true
        defer { isSaving = false }

        guard let token = authManager.currentToken else { return false }

        struct UpdateRequest: Encodable {
            let display_name: String
        }

        struct StatusResponse: Decodable { let status: String }

        do {
            let _: StatusResponse = try await apiClient.request(
                .updateMe,
                token: token,
                body: UpdateRequest(display_name: displayName)
            )
            saveSuccess = true
            return true
        } catch {
            errorMessage = error.localizedDescription
            return false
        }
    }

    func deleteAccount() async -> Bool {
        guard let token = authManager.currentToken else { return false }

        struct StatusResponse: Decodable { let status: String }

        do {
            let _: StatusResponse = try await apiClient.request(.deleteMe, token: token)
            return true
        } catch {
            errorMessage = error.localizedDescription
            return false
        }
    }
}
