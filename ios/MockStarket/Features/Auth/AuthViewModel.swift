import SwiftUI
import Observation

@MainActor @Observable
final class AuthViewModel {
    var isLoading = false
    var errorMessage: String?

    private let apiClient = APIClient.shared

    func signInAsGuest() async throws -> String {
        isLoading = true
        defer { isLoading = false }
        errorMessage = nil

        // In dev mode, the token IS the Firebase UID.
        // The backend treats Bearer tokens as UIDs when DEV_MODE=true.
        let guestUID = "guest_\(UUID().uuidString.prefix(8))"

        // Create the guest account on the backend
        let _: User = try await apiClient.request(
            .createGuest,
            token: guestUID
        )

        return guestUID
    }
}
