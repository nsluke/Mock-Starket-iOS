import SwiftUI

struct ProfileView: View {
    @Environment(AppState.self) private var appState
    @State private var viewModel = ProfileViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: 24) {
                // Avatar and name
                VStack(spacing: 12) {
                    Image(systemName: "person.circle.fill")
                        .font(.system(size: 80))
                        .foregroundStyle(Theme.accent)

                    if let user = appState.currentUser {
                        Text(user.displayName)
                            .font(.title2.weight(.bold))
                            .foregroundStyle(Theme.textPrimary)

                        Text("Member since \(user.createdAt.formatted(date: .abbreviated, time: .omitted))")
                            .font(.caption)
                            .foregroundStyle(Theme.textTertiary)
                    }
                }
                .padding(.top, 20)

                // Streak
                if let user = appState.currentUser, user.loginStreak > 0 {
                    HStack(spacing: 8) {
                        Image(systemName: "flame.fill")
                            .foregroundStyle(.orange)
                        Text("\(user.loginStreak) day streak")
                            .font(.subheadline.weight(.semibold))
                            .foregroundStyle(Theme.textPrimary)
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 10)
                    .background(Theme.surfaceElevated)
                    .clipShape(Capsule())
                }

                // Achievements
                if !viewModel.achievements.isEmpty {
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Achievements")
                            .font(.headline)
                            .foregroundStyle(Theme.textPrimary)
                            .padding(.horizontal)

                        LazyVGrid(columns: [
                            GridItem(.adaptive(minimum: 80))
                        ], spacing: 12) {
                            ForEach(viewModel.achievements) { achievement in
                                VStack(spacing: 6) {
                                    Image(systemName: achievement.icon)
                                        .font(.title2)
                                        .foregroundStyle(Theme.accent)
                                    Text(achievement.name)
                                        .font(.caption2)
                                        .foregroundStyle(Theme.textSecondary)
                                        .lineLimit(1)
                                }
                                .frame(width: 80, height: 80)
                                .background(Theme.surfaceElevated)
                                .clipShape(RoundedRectangle(cornerRadius: 12))
                            }
                        }
                        .padding(.horizontal)
                    }
                }

                // Quick actions
                VStack(spacing: 0) {
                    NavigationLink {
                        ChallengesView()
                    } label: {
                        settingsRow(icon: "star.fill", title: "Daily Challenge")
                    }
                    Divider().background(Theme.border)

                    NavigationLink {
                        AlertsView()
                    } label: {
                        settingsRow(icon: "bell.fill", title: "Price Alerts")
                    }
                    Divider().background(Theme.border)

                    NavigationLink {
                        SettingsView()
                    } label: {
                        settingsRow(icon: "gear", title: "Settings")
                    }
                    Divider().background(Theme.border)

                    Button {
                        appState.signOut()
                    } label: {
                        HStack {
                            Image(systemName: "rectangle.portrait.and.arrow.right")
                                .foregroundStyle(Theme.negative)
                            Text("Sign Out")
                                .foregroundStyle(Theme.negative)
                            Spacer()
                        }
                        .padding(.horizontal, 16)
                        .padding(.vertical, 14)
                    }
                }
                .background(Theme.surfaceElevated)
                .clipShape(RoundedRectangle(cornerRadius: 12))
                .padding(.horizontal)
            }
            .padding(.bottom, 40)
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Profile")
        .task {
            await viewModel.load()
        }
    }

    private func settingsRow(icon: String, title: String) -> some View {
        HStack {
            Image(systemName: icon)
                .foregroundStyle(Theme.accent)
            Text(title)
                .foregroundStyle(Theme.textPrimary)
            Spacer()
            Image(systemName: "chevron.right")
                .foregroundStyle(Theme.textTertiary)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
    }
}

#Preview {
    NavigationStack {
        ProfileView()
    }
    .environment(AppState())
    .preferredColorScheme(.dark)
}
