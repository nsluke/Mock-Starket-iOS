import SwiftUI

@main
struct MockStarketApp: App {
    @State private var appState = AppState()

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
