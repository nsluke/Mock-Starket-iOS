import SwiftUI
import FirebaseCore

@main
struct MockStarketApp: App {
    @State private var appState = AppState()

    init() {
        FirebaseApp.configure()
    }

    var body: some Scene {
        WindowGroup {
            Group {
                switch appState.authState {
                case .loading:
                    LaunchView()
                case .unauthenticated:
                    AuthRootView()
                case .authenticated:
                    MainTabView()
                }
            }
            .environment(appState)
            .task {
                await appState.checkAuth()
            }
        }
    }
}
