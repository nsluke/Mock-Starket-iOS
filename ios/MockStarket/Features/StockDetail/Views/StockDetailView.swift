import SwiftUI
import Charts

struct StockDetailView: View {
    let ticker: String
    @State private var viewModel = StockDetailViewModel()

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                // Price Header
                if let stock = viewModel.stock {
                    priceHeader(stock)
                }

                // Chart
                if !viewModel.priceHistory.isEmpty {
                    priceChart
                }

                // Time range picker
                timeRangePicker

                // Stats
                if let stock = viewModel.stock {
                    statsGrid(stock)
                }

                // ETF Holdings (if this is an ETF)
                if viewModel.stock?.assetType == "etf" && !viewModel.etfHoldings.isEmpty {
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Holdings")
                            .font(.headline)
                            .foregroundStyle(Theme.textPrimary)

                        ForEach(viewModel.etfHoldings) { holding in
                            HStack {
                                Text(holding.ticker)
                                    .font(.caption.monospaced().weight(.bold))
                                    .foregroundStyle(Theme.accent)
                                    .padding(.horizontal, 6)
                                    .padding(.vertical, 3)
                                    .background(Theme.accent.opacity(0.1))
                                    .clipShape(RoundedRectangle(cornerRadius: 4))

                                Text(holding.name)
                                    .font(.subheadline)
                                    .foregroundStyle(Theme.textSecondary)

                                Spacer()

                                Text("\((Decimal(string: holding.weight) ?? 0) * 100)%")
                                    .font(.subheadline.weight(.semibold))
                                    .foregroundStyle(Theme.textPrimary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                    .padding(.horizontal)
                }

                // User position
                if viewModel.userShares > 0 {
                    HStack {
                        VStack(alignment: .leading, spacing: 2) {
                            Text("Your Position")
                                .font(.caption)
                                .foregroundStyle(Theme.textTertiary)
                            Text("\(viewModel.userShares) share\(viewModel.userShares == 1 ? "" : "s")")
                                .font(.headline)
                        }
                        Spacer()
                        VStack(alignment: .trailing, spacing: 2) {
                            Text("Market Value")
                                .font(.caption)
                                .foregroundStyle(Theme.textTertiary)
                            Text(((viewModel.stock?.currentPrice ?? 0) * Decimal(viewModel.userShares)).formatted(.currency(code: "USD")))
                                .font(.headline)
                        }
                    }
                    .padding()
                    .background(Theme.surface)
                    .clipShape(RoundedRectangle(cornerRadius: 12))
                    .padding(.horizontal)
                }

                // Options chain link (stocks only)
                if viewModel.stock?.assetType == "stock" {
                    NavigationLink {
                        OptionsChainView(ticker: ticker)
                    } label: {
                        HStack {
                            Image(systemName: "chart.bar.doc.horizontal")
                            Text("View Options Chain")
                            Spacer()
                            Image(systemName: "chevron.right")
                                .font(.caption)
                                .foregroundStyle(Theme.textTertiary)
                        }
                        .font(.subheadline.weight(.semibold))
                        .foregroundStyle(Theme.accent)
                        .padding()
                        .background(Theme.accent.opacity(0.1))
                        .clipShape(RoundedRectangle(cornerRadius: 12))
                    }
                    .padding(.horizontal)
                }

                // Trade buttons
                tradeButtons

                // Description
                if let desc = viewModel.stock?.description {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("About")
                            .font(.headline)
                            .foregroundStyle(Theme.textPrimary)
                        Text(desc)
                            .font(.subheadline)
                            .foregroundStyle(Theme.textSecondary)
                    }
                    .padding(.horizontal)
                }
            }
            .padding(.vertical)
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle(ticker)
        .navigationBarTitleDisplayMode(.inline)
        .task {
            await viewModel.load(ticker: ticker)
        }
    }

    private func priceHeader(_ stock: Stock) -> some View {
        VStack(alignment: .leading, spacing: 4) {
            HStack(spacing: 8) {
                Text(stock.name)
                    .font(.title3.weight(.semibold))
                    .foregroundStyle(Theme.textSecondary)

                Text(stock.assetType.uppercased())
                    .font(.caption2.weight(.bold))
                    .padding(.horizontal, 6)
                    .padding(.vertical, 2)
                    .background(assetTypeColor(stock.assetType).opacity(0.15))
                    .foregroundStyle(assetTypeColor(stock.assetType))
                    .clipShape(RoundedRectangle(cornerRadius: 4))
            }

            HStack(alignment: .firstTextBaseline, spacing: 12) {
                Text(stock.currentPrice.currencyFormatted)
                    .font(.system(size: 40, weight: .bold, design: .rounded))
                    .foregroundStyle(Theme.textPrimary)
                    .contentTransition(.numericText())

                VStack(alignment: .leading) {
                    HStack(spacing: 4) {
                        Image(systemName: Theme.priceArrow(for: stock.change))
                        Text(stock.change.currencyFormatted)
                    }
                    Text(stock.changePct.percentFormatted)
                }
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(Theme.priceColor(for: stock.change))
            }
        }
        .padding(.horizontal)
    }

    private var priceChart: some View {
        Chart(viewModel.priceHistory) { point in
            LineMark(
                x: .value("Time", point.recordedAt),
                y: .value("Price", point.close)
            )
            .foregroundStyle(chartColor)

            AreaMark(
                x: .value("Time", point.recordedAt),
                y: .value("Price", point.close)
            )
            .foregroundStyle(
                LinearGradient(
                    colors: [chartColor.opacity(0.3), chartColor.opacity(0.0)],
                    startPoint: .top,
                    endPoint: .bottom
                )
            )
        }
        .chartYScale(domain: .automatic(includesZero: false))
        .chartXAxis(.hidden)
        .chartYAxis {
            AxisMarks(position: .trailing) { value in
                AxisValueLabel {
                    if let decimal = value.as(Decimal.self) {
                        Text(decimal.currencyFormatted)
                            .font(.caption2)
                            .foregroundStyle(Theme.textTertiary)
                    }
                }
            }
        }
        .frame(height: 240)
        .padding(.horizontal)
    }

    private var chartColor: Color {
        guard let stock = viewModel.stock else { return Theme.accent }
        return Theme.priceColor(for: stock.change)
    }

    private var timeRangePicker: some View {
        HStack(spacing: 0) {
            ForEach(StockDetailViewModel.TimeRange.allCases, id: \.self) { range in
                Button {
                    Task { await viewModel.selectTimeRange(range) }
                } label: {
                    Text(range.label)
                        .font(.caption.weight(.semibold))
                        .padding(.vertical, 8)
                        .frame(maxWidth: .infinity)
                        .background(viewModel.selectedRange == range ? Theme.accent.opacity(0.2) : .clear)
                        .foregroundStyle(viewModel.selectedRange == range ? Theme.accent : Theme.textSecondary)
                }
            }
        }
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 8))
        .padding(.horizontal)
    }

    private func statsGrid(_ stock: Stock) -> some View {
        LazyVGrid(columns: [
            GridItem(.flexible()),
            GridItem(.flexible()),
        ], spacing: 12) {
            statCell(title: "Open", value: stock.dayOpen.currencyFormatted)
            statCell(title: "Previous Close", value: stock.prevClose.currencyFormatted)
            statCell(title: "Day High", value: stock.dayHigh.currencyFormatted)
            statCell(title: "Day Low", value: stock.dayLow.currencyFormatted)
            statCell(title: "Volume", value: Decimal(stock.volume).volumeFormatted)
            statCell(title: "Sector", value: stock.sector)
        }
        .padding(.horizontal)
    }

    private func statCell(title: String, value: String) -> some View {
        VStack(alignment: .leading, spacing: 4) {
            Text(title)
                .font(.caption)
                .foregroundStyle(Theme.textTertiary)
            Text(value)
                .font(.subheadline.weight(.medium))
                .foregroundStyle(Theme.textPrimary)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(12)
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 8))
    }

    private var tradeButtons: some View {
        HStack(spacing: 12) {
            Button {
                viewModel.showTradeSheet = true
                viewModel.tradeSide = "buy"
            } label: {
                Label("Buy", systemImage: "plus.circle.fill")
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 14)
                    .background(Theme.positive)
                    .foregroundStyle(.white)
                    .clipShape(RoundedRectangle(cornerRadius: 12))
                    .font(.headline)
            }

            Button {
                viewModel.showTradeSheet = true
                viewModel.tradeSide = "sell"
            } label: {
                Label("Sell", systemImage: "minus.circle.fill")
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 14)
                    .background(Theme.negative)
                    .foregroundStyle(.white)
                    .clipShape(RoundedRectangle(cornerRadius: 12))
                    .font(.headline)
            }
        }
        .padding(.horizontal)
        .fullScreenCover(isPresented: $viewModel.showTradeSheet) {
            TradeSheetView(
                ticker: ticker,
                side: viewModel.tradeSide,
                currentPrice: viewModel.stock?.currentPrice ?? 0,
                currentShares: viewModel.userShares
            )
        }
    }

    private func assetTypeColor(_ type: String) -> Color {
        switch type {
        case "crypto": return .orange
        case "commodity": return .yellow
        case "etf": return .purple
        default: return Theme.accent
        }
    }
}

#Preview {
    NavigationStack {
        StockDetailView(ticker: "PLNX")
    }
    .preferredColorScheme(.dark)
}
