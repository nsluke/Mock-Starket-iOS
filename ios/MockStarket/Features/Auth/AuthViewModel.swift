import SwiftUI
import Observation

@Observable
final class AuthViewModel {
    var isLoading = false
    var errorMessage: String?

    private let apiClient = APIClient.shared

    func signInAsGuest() async throws -> String {
        isLoading = true
        defer { isLoading = false }
        errorMessage = nil

        // In dev mode, generate a unique guest token
        let guestUID = "guest_\(UUID().uuidString.prefix(8))"

        struct RegisterBody: Encodable {
            let display_name: String
            let is_guest: Bool
        }

        let _: User = try await apiClient.request(
            .createGuest,
            token: guestUID,
            body: RegisterBody(display_name: "Guest Trader", is_guest: true)
        )

        return guestUID
    }
}
