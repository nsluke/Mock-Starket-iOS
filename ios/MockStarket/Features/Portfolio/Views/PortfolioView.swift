import SwiftUI
import Charts

struct PortfolioView: View {
    @State private var viewModel = PortfolioViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                // Net Worth Hero
                if let response = viewModel.portfolioResponse {
                    netWorthCard(response)
                }

                // Breakdown
                if let response = viewModel.portfolioResponse {
                    breakdownRow(response)
                }

                // Holdings
                if !viewModel.positions.isEmpty {
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Holdings")
                            .font(.headline)
                            .foregroundStyle(Theme.textPrimary)
                            .padding(.horizontal)

                        ForEach(viewModel.positions) { position in
                            NavigationLink(value: Stock(ticker: position.ticker, name: "", sector: "", basePrice: 0, currentPrice: position.currentPrice, dayOpen: 0, dayHigh: 0, dayLow: 0, prevClose: 0, volume: 0, volatility: 0, description: nil)) {
                                PositionRowView(position: position)
                            }
                        }
                    }
                } else if !viewModel.isLoading {
                    VStack(spacing: 12) {
                        Image(systemName: "chart.pie")
                            .font(.system(size: 48))
                            .foregroundStyle(Theme.textTertiary)
                        Text("No holdings yet")
                            .font(.headline)
                            .foregroundStyle(Theme.textSecondary)
                        Text("Buy some stocks to get started!")
                            .font(.subheadline)
                            .foregroundStyle(Theme.textTertiary)
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 48)
                }
            }
            .padding(.vertical)
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Portfolio")
        .refreshable {
            await viewModel.load()
        }
        .task {
            await viewModel.load()
        }
    }

    private func netWorthCard(_ response: PortfolioResponse) -> some View {
        VStack(spacing: 8) {
            Text("Net Worth")
                .font(.caption.weight(.medium))
                .foregroundStyle(Theme.textSecondary)

            Text(response.netWorth.currencyFormatted)
                .font(.system(size: 40, weight: .bold, design: .rounded))
                .foregroundStyle(Theme.textPrimary)
                .contentTransition(.numericText())

            let returnAmount = response.netWorth - 100000  // Starting cash
            let returnPct = returnAmount / 100000 * 100
            HStack(spacing: 4) {
                Image(systemName: Theme.priceArrow(for: returnAmount))
                Text("\(returnAmount.currencyFormatted) (\(returnPct.percentFormatted))")
            }
            .font(.subheadline.weight(.semibold))
            .foregroundStyle(Theme.priceColor(for: returnAmount))
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 24)
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 16))
        .padding(.horizontal)
    }

    private func breakdownRow(_ response: PortfolioResponse) -> some View {
        HStack(spacing: 12) {
            breakdownCell(title: "Cash", value: response.portfolio.cash.currencyFormatted, icon: "banknote")
            breakdownCell(title: "Invested", value: response.invested.currencyFormatted, icon: "chart.bar")
        }
        .padding(.horizontal)
    }

    private func breakdownCell(title: String, value: String, icon: String) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Label(title, systemImage: icon)
                .font(.caption)
                .foregroundStyle(Theme.textTertiary)
            Text(value)
                .font(.headline)
                .foregroundStyle(Theme.textPrimary)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(16)
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 12))
    }
}

struct PositionRowView: View {
    let position: Position

    var body: some View {
        HStack(spacing: 12) {
            Text(position.ticker)
                .font(.system(.caption, design: .monospaced, weight: .bold))
                .foregroundStyle(.white)
                .padding(.horizontal, 8)
                .padding(.vertical, 4)
                .background(Theme.accent.opacity(0.2))
                .clipShape(RoundedRectangle(cornerRadius: 6))

            VStack(alignment: .leading, spacing: 2) {
                Text("\(position.shares) shares")
                    .font(.subheadline.weight(.medium))
                    .foregroundStyle(Theme.textPrimary)
                Text("Avg \(position.avgCost.currencyFormatted)")
                    .font(.caption2)
                    .foregroundStyle(Theme.textTertiary)
            }

            Spacer()

            VStack(alignment: .trailing, spacing: 2) {
                Text(position.marketValue.currencyFormatted)
                    .font(.subheadline.weight(.semibold))
                    .foregroundStyle(Theme.textPrimary)
                Text(position.pnl.currencyFormatted)
                    .font(.caption.weight(.medium))
                    .foregroundStyle(Theme.priceColor(for: position.pnl))
            }
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
    }
}

#Preview {
    NavigationStack {
        PortfolioView()
    }
    .preferredColorScheme(.dark)
}
