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
            Text(stock.name)
                .font(.title3.weight(.semibold))
                .foregroundStyle(Theme.textSecondary)

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
        .sheet(isPresented: $viewModel.showTradeSheet) {
            TradeSheetView(ticker: ticker, side: viewModel.tradeSide)
        }
    }
}

#Preview {
    NavigationStack {
        StockDetailView(ticker: "PLNX")
    }
    .preferredColorScheme(.dark)
}
