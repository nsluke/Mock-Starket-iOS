import SwiftUI

struct SettingsView: View {
    @Environment(AppState.self) private var appState
    @State private var viewModel = SettingsViewModel()
    @State private var showDeleteConfirm = false

    var body: some View {
        Form {
            // Profile
            Section("Profile") {
                TextField("Display Name", text: $viewModel.displayName)
                    .foregroundStyle(Theme.textPrimary)

                Button {
                    Task {
                        if await viewModel.updateProfile() {
                            // Refresh user data
                            await appState.checkAuth()
                        }
                    }
                } label: {
                    HStack {
                        Text(viewModel.isSaving ? "Saving..." : viewModel.saveSuccess ? "Saved!" : "Save Changes")
                        if viewModel.saveSuccess {
                            Image(systemName: "checkmark.circle.fill")
                                .foregroundStyle(Theme.positive)
                        }
                    }
                }
                .disabled(viewModel.isSaving || viewModel.displayName == appState.currentUser?.displayName)
            }

            // Account Info
            Section("Account") {
                if let user = appState.currentUser {
                    LabeledContent("Account Type", value: user.isGuest ? "Guest" : "Registered")
                    LabeledContent("Member Since", value: user.createdAt.formatted(date: .abbreviated, time: .omitted))
                    LabeledContent("Login Streak", value: "\(user.loginStreak) days")
                    LabeledContent("Longest Streak", value: "\(user.longestStreak) days")
                }
            }

            // App Info
            Section("App") {
                LabeledContent("Version", value: "1.0.0")
                LabeledContent("Build", value: "1")
            }

            // Danger Zone
            Section {
                Button(role: .destructive) {
                    showDeleteConfirm = true
                } label: {
                    Label("Delete Account", systemImage: "trash")
                }
            } footer: {
                Text("This will permanently delete your account and all trading data. This cannot be undone.")
            }

            if let error = viewModel.errorMessage {
                Section {
                    Text(error)
                        .foregroundStyle(Theme.negative)
                        .font(.caption)
                }
            }
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Settings")
        .onAppear {
            viewModel.load(user: appState.currentUser)
        }
        .confirmationDialog(
            "Delete Account",
            isPresented: $showDeleteConfirm,
            titleVisibility: .visible
        ) {
            Button("Delete Account", role: .destructive) {
                Task {
                    if await viewModel.deleteAccount() {
                        appState.signOut()
                    }
                }
            }
            Button("Cancel", role: .cancel) {}
        } message: {
            Text("This will permanently delete your account and all data. This cannot be undone.")
        }
    }
}

#Preview {
    NavigationStack {
        SettingsView()
    }
    .environment(AppState())
    .preferredColorScheme(.dark)
}
