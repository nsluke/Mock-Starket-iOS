import SwiftUI
import Observation

@MainActor @Observable
final class AppState {
    enum AuthState {
        case loading
        case unauthenticated
        case authenticated
    }

    var authState: AuthState = .loading
    var currentUser: User?
    var isConnected: Bool = false

    private let apiClient = APIClient.shared
    private let webSocketManager = WebSocketManager.shared
    private let authManager = AuthManager.shared

    func checkAuth() async {
        if let token = authManager.currentToken {
            do {
                let user: User = try await apiClient.request(.getMe, token: token)
                self.currentUser = user
                self.authState = .authenticated
                await connectWebSocket()
            } catch {
                self.authState = .unauthenticated
            }
        } else {
            self.authState = .unauthenticated
        }
    }

    func signIn(token: String) async throws {
        authManager.saveToken(token)
        let user: User = try await apiClient.request(.getMe, token: token)
        self.currentUser = user
        self.authState = .authenticated
        await connectWebSocket()
    }

    func signOut() {
        authManager.clearToken()
        webSocketManager.disconnect()
        currentUser = nil
        authState = .unauthenticated
    }

    private func connectWebSocket() async {
        guard let token = authManager.currentToken else { return }
        webSocketManager.connect(token: token)
    }
}
