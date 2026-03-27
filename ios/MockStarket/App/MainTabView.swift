import SwiftUI

struct MainTabView: View {
    var body: some View {
        TabView {
            Tab("Market", systemImage: "chart.line.uptrend.xyaxis") {
                NavigationStack {
                    MarketView()
                }
            }

            Tab("Portfolio", systemImage: "briefcase.fill") {
                NavigationStack {
                    PortfolioView()
                }
            }

            Tab("Leaderboard", systemImage: "trophy.fill") {
                NavigationStack {
                    LeaderboardView()
                }
            }

            Tab("Profile", systemImage: "person.circle.fill") {
                NavigationStack {
                    ProfileView()
                }
            }
        }
        .tint(Theme.accent)
    }
}

#Preview {
    MainTabView()
        .environment(AppState())
        .preferredColorScheme(.dark)
}
