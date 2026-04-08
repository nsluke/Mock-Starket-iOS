import SwiftUI
import Charts

struct MarketView: View {
    @State private var viewModel = MarketViewModel()

    var body: some View {
        List {
            marketSummarySection
            categoryFilterSection
            sectorFilterSection
            stockListSection
        }
        .listStyle(.plain)
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Market")
        .toolbar { cashToolbar }
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

    @ViewBuilder
    private var marketSummarySection: some View {
        if let summary = viewModel.marketSummary {
            MarketSummaryCard(summary: summary)
                .listRowBackground(Theme.surface)
        }
    }

    private var categoryFilterSection: some View {
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
    }

    @ViewBuilder
    private var sectorFilterSection: some View {
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
    }

    private var stockListSection: some View {
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

    @ToolbarContentBuilder
    private var cashToolbar: some ToolbarContent {
        if let portfolio = viewModel.portfolio {
            ToolbarItem(placement: .topBarTrailing) {
                VStack(alignment: .trailing, spacing: 1) {
                    Text(portfolio.portfolio.cash.currencyFormatted)
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(Theme.textPrimary)
                    Text("Cash")
                        .font(.caption2)
                        .foregroundStyle(Theme.textTertiary)
                }
            }
        }
    }
}

struct StockRowView: View {
    let stock: Stock

    private var showsSparkline: Bool {
        stock.assetType == "stock" || stock.assetType == "commodity"
    }

    var body: some View {
        HStack(spacing: 12) {
            // Logo or ticker badge
            if let logoURL = stock.logoURL, let url = URL(string: logoURL) {
                AsyncImage(url: url) { image in
                    image.resizable().aspectRatio(contentMode: .fit)
                } placeholder: {
                    Text(stock.displayTicker.prefix(2))
                        .font(.system(.caption2, design: .monospaced, weight: .bold))
                        .foregroundStyle(Theme.accent)
                }
                .frame(width: 32, height: 32)
                .clipShape(Circle())
                .background(Circle().fill(Theme.surfaceElevated))
            } else {
                Text(stock.displayTicker.prefix(2))
                    .font(.system(.caption, design: .monospaced, weight: .bold))
                    .foregroundStyle(Theme.accent)
                    .frame(width: 32, height: 32)
                    .background(Theme.accent.opacity(0.15))
                    .clipShape(Circle())
            }

            // Name & ticker
            VStack(alignment: .leading, spacing: 2) {
                Text(stock.displayTicker)
                    .font(.system(.caption, design: .monospaced, weight: .bold))
                    .foregroundStyle(Theme.accent)
                Text(stock.name)
                    .font(.subheadline.weight(.medium))
                    .foregroundStyle(Theme.textPrimary)
                    .lineLimit(1)
            }

            Spacer()

            // Sparkline (stocks & commodities only)
            if showsSparkline {
                SparklineView(ticker: stock.ticker, isUp: stock.isUp)
                    .frame(width: 60, height: 30)
            }

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

// MARK: - Sparkline

struct SparklineView: View {
    let ticker: String
    let isUp: Bool

    @State private var points: [Decimal] = []

    private var lineColor: Color {
        isUp ? Theme.positive : Theme.negative
    }

    var body: some View {
        Group {
            if points.count >= 2 {
                Chart(Array(points.enumerated()), id: \.offset) { index, price in
                    LineMark(
                        x: .value("T", index),
                        y: .value("P", price)
                    )
                    .foregroundStyle(lineColor)
                    .interpolationMethod(.catmullRom)
                }
                .chartXAxis(.hidden)
                .chartYAxis(.hidden)
                .chartYScale(domain: .automatic(includesZero: false))
                .chartLegend(.hidden)
            } else {
                Color.clear
            }
        }
        .task(id: ticker) {
            await loadSparkline()
        }
    }

    private func loadSparkline() async {
        if let cached = SparklineCache.shared.get(ticker) {
            points = cached
            return
        }
        let token = AuthManager.shared.currentToken
        do {
            let history: [PricePoint] = try await APIClient.shared.request(
                .getStockHistory(ticker: ticker, interval: "5m"),
                token: token
            )
            let prices = history.map(\.close)
            SparklineCache.shared.set(ticker, prices: prices)
            points = prices
        } catch {
            // Silently fail — row still shows price data
        }
    }
}

@MainActor
final class SparklineCache {
    static let shared = SparklineCache()
    private var cache: [String: [Decimal]] = [:]

    func get(_ ticker: String) -> [Decimal]? { cache[ticker] }
    func set(_ ticker: String, prices: [Decimal]) { cache[ticker] = prices }
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
