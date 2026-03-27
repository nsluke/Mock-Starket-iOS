import SwiftUI

struct AuthRootView: View {
    @Environment(AppState.self) private var appState
    @State private var viewModel = AuthViewModel()

    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()

            VStack(spacing: 40) {
                Spacer()

                // Logo and title
                VStack(spacing: 16) {
                    Image(systemName: "chart.line.uptrend.xyaxis")
                        .font(.system(size: 64))
                        .foregroundStyle(Theme.accent)

                    Text("Mock Starket")
                        .font(.system(size: 36, weight: .bold, design: .rounded))
                        .foregroundStyle(Theme.textPrimary)

                    Text("Learn to trade. Risk nothing.")
                        .font(.subheadline)
                        .foregroundStyle(Theme.textSecondary)
                }

                Spacer()

                // Action buttons
                VStack(spacing: 16) {
                    if viewModel.isLoading {
                        ProgressView()
                            .tint(Theme.accent)
                    } else {
                        Button {
                            Task { await signInAsGuest() }
                        } label: {
                            Label("Continue as Guest", systemImage: "person.fill")
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(Theme.accent)
                                .foregroundStyle(.black)
                                .clipShape(RoundedRectangle(cornerRadius: 14))
                                .font(.headline)
                        }

                        Button {
                            Task { await signInAsGuest() } // TODO: Replace with real auth
                        } label: {
                            Label("Sign in with Email", systemImage: "envelope.fill")
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(Theme.surfaceElevated)
                                .foregroundStyle(Theme.textPrimary)
                                .clipShape(RoundedRectangle(cornerRadius: 14))
                                .font(.headline)
                        }
                    }

                    if let error = viewModel.errorMessage {
                        Text(error)
                            .font(.caption)
                            .foregroundStyle(Theme.negative)
                    }
                }
                .padding(.horizontal, 24)

                Spacer()
                    .frame(height: 40)
            }
        }
    }

    private func signInAsGuest() async {
        do {
            let token = try await viewModel.signInAsGuest()
            try await appState.signIn(token: token)
        } catch {
            viewModel.errorMessage = error.localizedDescription
        }
    }
}

struct LaunchView: View {
    var body: some View {
        ZStack {
            Theme.background.ignoresSafeArea()
            VStack(spacing: 16) {
                Image(systemName: "chart.line.uptrend.xyaxis")
                    .font(.system(size: 48))
                    .foregroundStyle(Theme.accent)
                ProgressView()
                    .tint(Theme.accent)
            }
        }
    }
}

#Preview {
    AuthRootView()
        .environment(AppState())
}
