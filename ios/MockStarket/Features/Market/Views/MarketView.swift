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

            // Asset Category Filter
            Section {
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 8) {
                        ForEach(AssetCategory.allCases) { category in
                            SectorChip(title: category.rawValue, isSelected: viewModel.selectedCategory == category) {
                                viewModel.selectedCategory = category
                                viewModel.selectedSector = nil
                            }
                        }
                    }
                    .padding(.horizontal, 4)
                }
                .listRowBackground(Theme.surface)
                .listRowSeparator(.hidden)
                .listRowInsets(EdgeInsets(top: 4, leading: 8, bottom: 4, trailing: 8))
            }

            // Sector Filter
            if viewModel.availableSectors.count > 1 {
                Section {
                    ScrollView(.horizontal, showsIndicators: false) {
                        HStack(spacing: 8) {
                            SectorChip(title: "All Sectors", isSelected: viewModel.selectedSector == nil) {
                                viewModel.selectedSector = nil
                            }
                            ForEach(viewModel.availableSectors, id: \.self) { sector in
                                SectorChip(title: sector, isSelected: viewModel.selectedSector == sector) {
                                    viewModel.selectedSector = sector
                                }
                            }
                        }
                        .padding(.horizontal, 4)
                    }
                    .listRowBackground(Theme.surface)
                    .listRowSeparator(.hidden)
                    .listRowInsets(EdgeInsets(top: 0, leading: 8, bottom: 0, trailing: 8))
                }
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
                Text(viewModel.selectedCategory == .all ? "All Assets" : viewModel.selectedCategory.rawValue)
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

struct SectorChip: View {
    let title: String
    let isSelected: Bool
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Text(title)
                .font(.caption.weight(.medium))
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(isSelected ? Theme.accent.opacity(0.2) : Theme.surfaceElevated)
                .foregroundStyle(isSelected ? Theme.accent : Theme.textSecondary)
                .clipShape(Capsule())
        }
        .buttonStyle(.plain)
    }
}

#Preview {
    NavigationStack {
        MarketView()
    }
    .preferredColorScheme(.dark)
}
