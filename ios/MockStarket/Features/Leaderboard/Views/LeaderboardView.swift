import SwiftUI

struct LeaderboardView: View {
    @State private var viewModel = LeaderboardViewModel()

    var body: some View {
        VStack(spacing: 0) {
            // Period picker
            Picker("Period", selection: $viewModel.selectedPeriod) {
                Text("Daily").tag("daily")
                Text("Weekly").tag("weekly")
                Text("All Time").tag("alltime")
            }
            .pickerStyle(.segmented)
            .padding()

            List {
                ForEach(viewModel.entries) { entry in
                    LeaderboardRowView(entry: entry)
                        .listRowBackground(Theme.surface)
                }
            }
            .listStyle(.plain)
            .scrollContentBackground(.hidden)
        }
        .background(Theme.background)
        .navigationTitle("Leaderboard")
        .refreshable {
            await viewModel.load()
        }
        .task {
            await viewModel.load()
        }
        .onChange(of: viewModel.selectedPeriod) {
            Task { await viewModel.load() }
        }
    }
}

struct LeaderboardRowView: View {
    let entry: LeaderboardEntry

    var body: some View {
        HStack(spacing: 12) {
            // Rank
            ZStack {
                Circle()
                    .fill(rankColor)
                    .frame(width: 36, height: 36)
                Text("#\(entry.rank)")
                    .font(.caption.weight(.bold))
                    .foregroundStyle(.white)
            }

            // Name
            VStack(alignment: .leading, spacing: 2) {
                Text(entry.displayName)
                    .font(.subheadline.weight(.medium))
                    .foregroundStyle(Theme.textPrimary)
                Text(entry.totalReturn.percentFormatted)
                    .font(.caption)
                    .foregroundStyle(Theme.priceColor(for: entry.totalReturn))
            }

            Spacer()

            // Net worth
            Text(entry.netWorth.currencyFormatted)
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(Theme.textPrimary)
        }
        .padding(.vertical, 4)
    }

    private var rankColor: Color {
        switch entry.rank {
        case 1: return .yellow
        case 2: return .gray
        case 3: return .orange
        default: return Theme.surfaceElevated
        }
    }
}

#Preview {
    NavigationStack {
        LeaderboardView()
    }
    .preferredColorScheme(.dark)
}
