import SwiftUI

struct MarketView: View {
    @State private var viewModel = MarketViewModel()

    var body: some View {
        List {
            // Market Summary Card
            if let summary = viewModel.marketSummary {
                MarketSummaryCard(summary: summary)
                    .listRowBackground(Theme.surface)
            }

            // Stock List
            Section {
                ForEach(viewModel.filteredStocks) { stock in
                    NavigationLink(value: stock) {
                        StockRowView(stock: stock)
                    }
                    .listRowBackground(Theme.surface)
                }
            } header: {
                Text("All Stocks")
                    .foregroundStyle(Theme.textSecondary)
            }
        }
        .listStyle(.plain)
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Market")
        .searchable(text: $viewModel.searchText, prompt: "Search stocks...")
        .refreshable {
            await viewModel.loadStocks()
        }
        .navigationDestination(for: Stock.self) { stock in
            StockDetailView(ticker: stock.ticker)
        }
        .task {
            await viewModel.loadStocks()
            viewModel.subscribeToUpdates()
        }
    }
}

struct StockRowView: View {
    let stock: Stock

    var body: some View {
        HStack(spacing: 12) {
            // Ticker badge
            Text(stock.ticker)
                .font(.system(.caption, design: .monospaced, weight: .bold))
                .foregroundStyle(.white)
                .padding(.horizontal, 8)
                .padding(.vertical, 4)
                .background(Theme.accent.opacity(0.2))
                .clipShape(RoundedRectangle(cornerRadius: 6))

            // Name
            VStack(alignment: .leading, spacing: 2) {
                Text(stock.name)
                    .font(.subheadline.weight(.medium))
                    .foregroundStyle(Theme.textPrimary)
                    .lineLimit(1)
                Text(stock.sector)
                    .font(.caption2)
                    .foregroundStyle(Theme.textTertiary)
            }

            Spacer()

            // Price and change
            VStack(alignment: .trailing, spacing: 2) {
                Text(stock.currentPrice.currencyFormatted)
                    .font(.subheadline.weight(.semibold))
                    .foregroundStyle(Theme.textPrimary)
                    .contentTransition(.numericText())

                HStack(spacing: 2) {
                    Image(systemName: Theme.priceArrow(for: stock.change))
                        .font(.caption2)
                    Text(stock.changePct.percentFormatted)
                        .font(.caption.weight(.medium))
                }
                .foregroundStyle(Theme.priceColor(for: stock.change))
            }
        }
        .padding(.vertical, 4)
    }
}

struct MarketSummaryCard: View {
    let summary: MarketSummary

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Market Index")
                    .font(.caption.weight(.medium))
                    .foregroundStyle(Theme.textSecondary)
                Spacer()
                Text("\(summary.gainers) gainers / \(summary.losers) losers")
                    .font(.caption2)
                    .foregroundStyle(Theme.textTertiary)
            }

            HStack(alignment: .firstTextBaseline, spacing: 8) {
                Text(summary.indexValue.currencyFormatted)
                    .font(.title.weight(.bold))
                    .foregroundStyle(Theme.textPrimary)

                Text(summary.indexChangePct.percentFormatted)
                    .font(.subheadline.weight(.semibold))
                    .foregroundStyle(Theme.priceColor(for: summary.indexChangePct))
            }
        }
        .padding()
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 12))
    }
}

#Preview {
    NavigationStack {
        MarketView()
    }
    .preferredColorScheme(.dark)
}
