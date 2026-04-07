import SwiftUI

struct OptionsChainView: View {
    let ticker: String
    @State private var viewModel = OptionsChainViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                // Educational banner
                OptionsEducationalBanner()
                    .padding(.horizontal)

                // Expiration picker
                if !viewModel.expirations.isEmpty {
                    ScrollView(.horizontal, showsIndicators: false) {
                        HStack(spacing: 8) {
                            ForEach(viewModel.expirations, id: \.self) { date in
                                let isSelected = viewModel.selectedExpiration == date
                                Button {
                                    Task { await viewModel.selectExpiration(date, ticker: ticker) }
                                } label: {
                                    Text(date.formatted(.dateTime.month(.abbreviated).day()))
                                        .font(.caption.weight(.medium))
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 6)
                                        .background(isSelected ? Theme.accent.opacity(0.2) : Theme.surfaceElevated)
                                        .foregroundStyle(isSelected ? Theme.accent : Theme.textSecondary)
                                        .clipShape(Capsule())
                                }
                            }
                        }
                        .padding(.horizontal)
                    }
                }

                // Chain table
                if let chain = viewModel.chain {
                    chainTable(chain)
                } else if viewModel.isLoading {
                    ProgressView()
                        .padding(.vertical, 40)
                } else {
                    Text("No options available for this stock.")
                        .font(.subheadline)
                        .foregroundStyle(Theme.textTertiary)
                        .padding(.vertical, 40)
                }
            }
            .padding(.vertical)
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("\(ticker) Options")
        .navigationBarTitleDisplayMode(.inline)
        .task {
            await viewModel.load(ticker: ticker)
        }
        .sheet(isPresented: $viewModel.showTradeSheet) {
            if let contract = viewModel.selectedContract, let chain = viewModel.chain {
                OptionsTradeSheetView(
                    contract: contract,
                    underlyingPrice: chain.underlyingPrice
                )
            }
        }
    }

    @ViewBuilder
    private func chainTable(_ chain: OptionChainResponse) -> some View {
        let strikes = allStrikes(chain)
        let callMap = Dictionary(uniqueKeysWithValues: chain.calls.map { ($0.strikePrice, $0) })
        let putMap = Dictionary(uniqueKeysWithValues: chain.puts.map { ($0.strikePrice, $0) })

        VStack(spacing: 0) {
            // Header
            HStack {
                Text("CALLS")
                    .font(.caption.weight(.bold))
                    .foregroundStyle(.green)
                    .frame(maxWidth: .infinity, alignment: .leading)
                Text("Strike")
                    .font(.caption.weight(.bold))
                    .foregroundStyle(Theme.textSecondary)
                    .frame(width: 70)
                Text("PUTS")
                    .font(.caption.weight(.bold))
                    .foregroundStyle(.red)
                    .frame(maxWidth: .infinity, alignment: .trailing)
            }
            .padding(.horizontal)
            .padding(.vertical, 8)
            .background(Theme.surfaceElevated)

            Divider().background(Theme.surfaceElevated)

            // Rows
            ForEach(strikes, id: \.self) { strike in
                let call = callMap[strike]
                let put = putMap[strike]
                let isATM = abs(chain.underlyingPrice - strike) / chain.underlyingPrice < Decimal(string: "0.01")!
                let callITM = chain.underlyingPrice > strike
                let putITM = chain.underlyingPrice < strike

                HStack(spacing: 0) {
                    // Call side
                    Button {
                        if let call { viewModel.openTrade(contract: call) }
                    } label: {
                        contractCell(contract: call, itm: callITM)
                            .frame(maxWidth: .infinity)
                    }
                    .buttonStyle(.plain)

                    // Strike
                    VStack(spacing: 2) {
                        Text(strike.currencyFormatted)
                            .font(.caption.weight(.bold))
                            .foregroundStyle(Theme.textPrimary)
                        if isATM {
                            Text("ATM")
                                .font(.system(size: 8, weight: .bold))
                                .foregroundStyle(.yellow)
                        }
                    }
                    .frame(width: 70)

                    // Put side
                    Button {
                        if let put { viewModel.openTrade(contract: put) }
                    } label: {
                        contractCell(contract: put, itm: putITM)
                            .frame(maxWidth: .infinity)
                    }
                    .buttonStyle(.plain)
                }
                .padding(.vertical, 6)
                .padding(.horizontal)
                .background(isATM ? Color.yellow.opacity(0.05) : .clear)

                Divider().background(Theme.surfaceElevated.opacity(0.5))
            }
        }
        .background(Theme.surface)
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .padding(.horizontal)
    }

    private func contractCell(contract: OptionContract?, itm: Bool) -> some View {
        Group {
            if let c = contract {
                VStack(spacing: 2) {
                    HStack {
                        Text(c.bidPrice.currencyFormatted)
                            .font(.caption2.monospacedDigit())
                            .foregroundStyle(Theme.textSecondary)
                        Spacer()
                        Text(c.askPrice.currencyFormatted)
                            .font(.caption2.monospacedDigit())
                            .foregroundStyle(Theme.textPrimary)
                    }
                    HStack {
                        Text("Δ \(c.delta.formatted(.number.precision(.fractionLength(2))))")
                            .font(.system(size: 9))
                            .foregroundStyle(Theme.textTertiary)
                        Spacer()
                        Text("\((c.impliedVol * 100).formatted(.number.precision(.fractionLength(1))))%")
                            .font(.system(size: 9))
                            .foregroundStyle(Theme.textTertiary)
                    }
                }
                .padding(.horizontal, 6)
                .padding(.vertical, 4)
                .background(itm ? Color.green.opacity(0.05) : .clear)
                .clipShape(RoundedRectangle(cornerRadius: 6))
            } else {
                Text("—")
                    .font(.caption2)
                    .foregroundStyle(Theme.textTertiary)
            }
        }
    }

    private func allStrikes(_ chain: OptionChainResponse) -> [Decimal] {
        var set = Set<Decimal>()
        chain.calls.forEach { set.insert($0.strikePrice) }
        chain.puts.forEach { set.insert($0.strikePrice) }
        return set.sorted()
    }
}

// MARK: - Educational Banner

struct OptionsEducationalBanner: View {
    @State private var expanded = false

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                Image(systemName: "lightbulb.fill")
                    .foregroundStyle(Theme.accent)
                    .font(.caption)
                Text("New to Options?")
                    .font(.caption.weight(.semibold))
                    .foregroundStyle(Theme.accent)
                Spacer()
            }

            Text("Options let you bet on a stock's direction without owning it. Each contract = 100 shares. Tap any row to trade.")
                .font(.caption2)
                .foregroundStyle(Theme.textSecondary)

            if expanded {
                VStack(alignment: .leading, spacing: 6) {
                    infoRow(title: "Calls", text: "Profit when stock goes up. Pay premium for right to buy at strike.", color: .green)
                    infoRow(title: "Puts", text: "Profit when stock goes down. Pay premium for right to sell at strike.", color: .red)
                    infoRow(title: "ITM", text: "In-the-Money — has intrinsic value.", color: .green)
                    infoRow(title: "OTM", text: "Out-of-the-Money — no intrinsic value yet.", color: .gray)
                    infoRow(title: "ATM", text: "At-the-Money — strike ≈ current price.", color: .yellow)
                }
            }

            Button(expanded ? "Show less" : "Learn more") {
                withAnimation { expanded.toggle() }
            }
            .font(.caption2.weight(.medium))
            .foregroundStyle(Theme.accent)
        }
        .padding()
        .background(Theme.accent.opacity(0.05))
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .overlay(RoundedRectangle(cornerRadius: 12).stroke(Theme.accent.opacity(0.2), lineWidth: 1))
    }

    private func infoRow(title: String, text: String, color: Color) -> some View {
        HStack(alignment: .top, spacing: 8) {
            Text(title)
                .font(.caption2.weight(.bold))
                .foregroundStyle(color)
                .frame(width: 40, alignment: .leading)
            Text(text)
                .font(.caption2)
                .foregroundStyle(Theme.textSecondary)
        }
    }
}
