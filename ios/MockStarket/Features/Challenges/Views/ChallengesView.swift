import SwiftUI

struct ChallengesView: View {
    @State private var viewModel = ChallengesViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                if viewModel.isLoading {
                    ProgressView()
                        .tint(Theme.accent)
                        .padding(.top, 40)
                } else if let challenge = viewModel.challenge {
                    challengeCard(challenge)
                } else {
                    ContentUnavailableView(
                        "No Challenge Today",
                        systemImage: "star.slash",
                        description: Text("Check back later for today's daily challenge.")
                    )
                    .padding(.top, 40)
                }
            }
            .padding()
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Daily Challenge")
        .task {
            await viewModel.load()
        }
    }

    private func challengeCard(_ challenge: DailyChallenge) -> some View {
        VStack(spacing: 20) {
            // Challenge icon and type
            VStack(spacing: 12) {
                Image(systemName: challengeIcon(for: challenge.challengeType))
                    .font(.system(size: 48))
                    .foregroundStyle(Theme.accent)

                Text(challenge.description)
                    .font(.title3.weight(.semibold))
                    .foregroundStyle(Theme.textPrimary)
                    .multilineTextAlignment(.center)
            }
            .padding(.top, 8)

            // Reward
            HStack(spacing: 6) {
                Image(systemName: "dollarsign.circle.fill")
                    .foregroundStyle(Theme.positive)
                Text("Reward: \(challenge.rewardCash.formatted(.currency(code: "USD")))")
                    .font(.subheadline.weight(.medium))
                    .foregroundStyle(Theme.textSecondary)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .background(Theme.positive.opacity(0.1))
            .clipShape(Capsule())

            // Status and actions
            if viewModel.isClaimed {
                Label("Reward Claimed!", systemImage: "checkmark.seal.fill")
                    .font(.headline)
                    .foregroundStyle(Theme.positive)
                    .padding()
                    .frame(maxWidth: .infinity)
                    .background(Theme.positive.opacity(0.1))
                    .clipShape(RoundedRectangle(cornerRadius: 14))

            } else if viewModel.isCompleted {
                Button {
                    Task { await viewModel.claimReward() }
                } label: {
                    Label(
                        viewModel.isClaiming ? "Claiming..." : "Claim Reward",
                        systemImage: "gift.fill"
                    )
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 16)
                    .background(Theme.accent)
                    .foregroundStyle(.black)
                    .clipShape(RoundedRectangle(cornerRadius: 14))
                    .font(.headline)
                }
                .disabled(viewModel.isClaiming)

            } else {
                VStack(spacing: 12) {
                    // Progress indicator
                    HStack {
                        Circle()
                            .fill(Theme.surfaceElevated)
                            .frame(width: 12, height: 12)
                        Text("In Progress")
                            .font(.subheadline)
                            .foregroundStyle(Theme.textSecondary)
                    }

                    Button {
                        Task { await viewModel.checkProgress() }
                    } label: {
                        Label(
                            viewModel.isChecking ? "Checking..." : "Check Progress",
                            systemImage: "arrow.clockwise"
                        )
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 16)
                        .background(Theme.surfaceElevated)
                        .foregroundStyle(Theme.textPrimary)
                        .clipShape(RoundedRectangle(cornerRadius: 14))
                        .font(.headline)
                    }
                    .disabled(viewModel.isChecking)
                }
            }

            if let error = viewModel.errorMessage {
                Text(error)
                    .font(.caption)
                    .foregroundStyle(Theme.negative)
            }
        }
        .padding(24)
        .background(Theme.surface)
        .clipShape(RoundedRectangle(cornerRadius: 16))
        .overlay(
            RoundedRectangle(cornerRadius: 16)
                .strokeBorder(Theme.border, lineWidth: 1)
        )
    }

    private func challengeIcon(for type: String) -> String {
        switch type {
        case "trade_count": return "repeat"
        case "buy_stock": return "cart.fill"
        case "sell_stock": return "banknote"
        case "profit_target": return "chart.line.uptrend.xyaxis"
        case "diversify": return "square.grid.3x3.fill"
        case "volume_trader": return "chart.bar.fill"
        default: return "star.fill"
        }
    }
}

#Preview {
    NavigationStack {
        ChallengesView()
    }
    .environment(AppState())
    .preferredColorScheme(.dark)
}
