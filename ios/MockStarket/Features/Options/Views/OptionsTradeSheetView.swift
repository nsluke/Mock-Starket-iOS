import SwiftUI

struct OptionsTradeSheetView: View {
    let contract: OptionContract
    let underlyingPrice: Decimal
    @Environment(\.dismiss) private var dismiss

    @State private var side: String = "buy_to_open"
    @State private var quantity: Int = 1
    @State private var isReview = false
    @State private var isSubmitting = false
    @State private var errorMessage: String?

    private let apiClient = APIClient.shared

    private var price: Decimal {
        side.hasPrefix("buy") ? contract.askPrice : contract.bidPrice
    }

    private var totalCost: Decimal {
        price * Decimal(quantity) * 100
    }

    private var breakEven: Decimal {
        contract.isCall ? contract.strikePrice + price : contract.strikePrice - price
    }

    private var isLong: Bool {
        side.hasPrefix("buy")
    }

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(spacing: 20) {
                    // Contract info header
                    contractHeader

                    if !isReview {
                        configureView
                    } else {
                        reviewView
                    }
                }
                .padding()
            }
            .scrollContentBackground(.hidden)
            .background(Theme.background)
            .navigationTitle("Trade Option")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") { dismiss() }
                }
            }
        }
    }

    // MARK: - Contract Header

    private var contractHeader: some View {
        VStack(spacing: 8) {
            HStack {
                Text(contract.ticker)
                    .font(.title3.weight(.bold))
                Text(contract.isCall ? "CALL" : "PUT")
                    .font(.caption.weight(.bold))
                    .padding(.horizontal, 8)
                    .padding(.vertical, 3)
                    .background((contract.isCall ? Color.green : Color.red).opacity(0.15))
                    .foregroundStyle(contract.isCall ? .green : .red)
                    .clipShape(Capsule())
                Spacer()
                moneynessBadge
            }

            HStack {
                Text("$\(contract.strikePrice.currencyFormatted) strike")
                    .font(.subheadline)
                    .foregroundStyle(Theme.textSecondary)
                Spacer()
                Text("Exp \(contract.expiration.formatted(.dateTime.month(.abbreviated).day()))")
                    .font(.subheadline)
                    .foregroundStyle(Theme.textSecondary)
            }

            // Bid / Mark / Ask
            HStack(spacing: 12) {
                priceCell(label: "Bid", value: contract.bidPrice)
                priceCell(label: "Mark", value: contract.markPrice, accent: true)
                priceCell(label: "Ask", value: contract.askPrice)
            }
        }
        .padding()
        .background(Theme.surfaceElevated)
        .clipShape(RoundedRectangle(cornerRadius: 12))
    }

    private var moneynessBadge: some View {
        let diff = abs(underlyingPrice - contract.strikePrice) / underlyingPrice
        let isATM = diff < Decimal(string: "0.01")!
        let isITM = contract.isCall ? underlyingPrice > contract.strikePrice : underlyingPrice < contract.strikePrice
        let label = isATM ? "ATM" : isITM ? "ITM" : "OTM"
        let color: Color = isATM ? .yellow : isITM ? .green : .gray

        return Text(label)
            .font(.caption2.weight(.bold))
            .padding(.horizontal, 6)
            .padding(.vertical, 2)
            .background(color.opacity(0.15))
            .foregroundStyle(color)
            .clipShape(Capsule())
    }

    private func priceCell(label: String, value: Decimal, accent: Bool = false) -> some View {
        VStack(spacing: 2) {
            Text(label)
                .font(.caption2)
                .foregroundStyle(Theme.textTertiary)
            Text(value.currencyFormatted)
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(accent ? Theme.accent : Theme.textPrimary)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 8)
        .background(accent ? Theme.accent.opacity(0.1) : Theme.surface)
        .clipShape(RoundedRectangle(cornerRadius: 8))
    }

    // MARK: - Configure View

    private var configureView: some View {
        VStack(spacing: 16) {
            // Side picker
            VStack(alignment: .leading, spacing: 8) {
                Text("Order Side")
                    .font(.caption.weight(.medium))
                    .foregroundStyle(Theme.textTertiary)

                HStack(spacing: 8) {
                    sideButton("buy_to_open", label: "Buy to Open", color: .green)
                    sideButton("sell_to_open", label: "Sell to Open", color: .orange)
                }
            }

            // Quantity
            VStack(alignment: .leading, spacing: 8) {
                Text("Contracts")
                    .font(.caption.weight(.medium))
                    .foregroundStyle(Theme.textTertiary)
                HStack {
                    Button { if quantity > 1 { quantity -= 1 } } label: {
                        Image(systemName: "minus.circle.fill")
                            .font(.title2)
                            .foregroundStyle(Theme.textSecondary)
                    }
                    Text("\(quantity)")
                        .font(.title2.weight(.bold).monospacedDigit())
                        .frame(minWidth: 60)
                    Button { if quantity < 100 { quantity += 1 } } label: {
                        Image(systemName: "plus.circle.fill")
                            .font(.title2)
                            .foregroundStyle(Theme.textSecondary)
                    }
                    Spacer()
                    Text("= \(quantity * 100) shares")
                        .font(.caption)
                        .foregroundStyle(Theme.textTertiary)
                }
            }

            // Cost summary
            VStack(spacing: 8) {
                summaryRow(label: isLong ? "Total Cost" : "Premium Received", value: totalCost.currencyFormatted, color: isLong ? .red : .green)
                summaryRow(label: "Break-even at expiry", value: breakEven.currencyFormatted)
            }
            .padding()
            .background(Theme.surface)
            .clipShape(RoundedRectangle(cornerRadius: 12))

            // Risk warning for writing
            if !isLong {
                HStack(alignment: .top, spacing: 8) {
                    Image(systemName: "exclamationmark.triangle.fill")
                        .foregroundStyle(.orange)
                        .font(.caption)
                    Text("Writing options can result in losses greater than the premium received. Collateral is required.")
                        .font(.caption2)
                        .foregroundStyle(Theme.textSecondary)
                }
                .padding()
                .background(Color.orange.opacity(0.05))
                .clipShape(RoundedRectangle(cornerRadius: 12))
            }

            // Greeks
            GreeksView(contract: contract)

            Button("Review Order") {
                withAnimation { isReview = true }
            }
            .buttonStyle(.borderedProminent)
            .tint(Theme.accent)
            .controlSize(.large)
        }
    }

    // MARK: - Review View

    private var reviewView: some View {
        VStack(spacing: 16) {
            VStack(spacing: 8) {
                summaryRow(label: "Action", value: sideLabel(side))
                summaryRow(label: "Contract", value: "\(contract.ticker) $\(contract.strikePrice.currencyFormatted) \(contract.optionType.uppercased())")
                summaryRow(label: "Expiration", value: contract.expiration.formatted(.dateTime.month(.abbreviated).day(.twoDigits).year()))
                summaryRow(label: "Quantity", value: "\(quantity) contract\(quantity > 1 ? "s" : "") (\(quantity * 100) shares)")
                summaryRow(label: "Price/contract", value: price.currencyFormatted)
                Divider().background(Theme.surfaceElevated)
                summaryRow(label: isLong ? "Total Debit" : "Total Credit", value: totalCost.currencyFormatted, color: isLong ? .red : .green)
                summaryRow(label: "Break-even", value: breakEven.currencyFormatted, color: .yellow)
            }
            .padding()
            .background(Theme.surface)
            .clipShape(RoundedRectangle(cornerRadius: 12))

            if let error = errorMessage {
                Text(error)
                    .font(.caption)
                    .foregroundStyle(.red)
                    .multilineTextAlignment(.center)
            }

            Text("This is a simulated trade for educational purposes.")
                .font(.caption2)
                .foregroundStyle(Theme.textTertiary)
                .multilineTextAlignment(.center)

            HStack(spacing: 12) {
                Button("Back") {
                    withAnimation { isReview = false }
                }
                .buttonStyle(.bordered)
                .controlSize(.large)

                Button(isSubmitting ? "Executing..." : "Confirm Trade") {
                    Task { await submitTrade() }
                }
                .buttonStyle(.borderedProminent)
                .tint(Theme.accent)
                .controlSize(.large)
                .disabled(isSubmitting)
            }
        }
    }

    // MARK: - Helpers

    private func sideButton(_ value: String, label: String, color: Color) -> some View {
        Button {
            side = value
        } label: {
            Text(label)
                .font(.caption.weight(.medium))
                .frame(maxWidth: .infinity)
                .padding(.vertical, 10)
                .background(side == value ? color.opacity(0.15) : Theme.surface)
                .foregroundStyle(side == value ? color : Theme.textSecondary)
                .clipShape(RoundedRectangle(cornerRadius: 8))
                .overlay(RoundedRectangle(cornerRadius: 8).stroke(side == value ? color.opacity(0.5) : .clear, lineWidth: 1))
        }
        .buttonStyle(.plain)
    }

    private func summaryRow(label: String, value: String, color: Color? = nil) -> some View {
        HStack {
            Text(label)
                .font(.subheadline)
                .foregroundStyle(Theme.textSecondary)
            Spacer()
            Text(value)
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(color ?? Theme.textPrimary)
        }
    }

    private func sideLabel(_ side: String) -> String {
        switch side {
        case "buy_to_open": return "Buy to Open"
        case "sell_to_open": return "Sell to Open"
        case "buy_to_close": return "Buy to Close"
        case "sell_to_close": return "Sell to Close"
        default: return side
        }
    }

    private func submitTrade() async {
        isSubmitting = true
        errorMessage = nil

        struct TradeBody: Encodable {
            let contract_id: UUID
            let side: String
            let quantity: Int
        }

        do {
            let token = AuthManager.shared.currentToken
            let _: OptionTrade = try await apiClient.request(
                .executeOptionsTrade,
                token: token,
                body: TradeBody(contract_id: contract.id, side: side, quantity: quantity)
            )
            dismiss()
        } catch {
            errorMessage = error.localizedDescription
        }

        isSubmitting = false
    }
}

// MARK: - Greeks View

struct GreeksView: View {
    let contract: OptionContract
    @State private var expandedGreek: String?

    private let greekInfo: [(key: String, label: String, symbol: String, description: String)] = [
        ("delta", "Delta", "Δ", "How much the option price moves per $1 change in the stock. Calls: 0 to 1, Puts: -1 to 0."),
        ("gamma", "Gamma", "Γ", "Rate of change of delta. Higher gamma = delta changes faster as stock moves."),
        ("theta", "Theta", "Θ", "Time decay — how much value lost per day. Hurts buyers, helps sellers."),
        ("vega", "Vega", "ν", "Sensitivity to volatility. Higher vega = more affected by IV changes."),
        ("rho", "Rho", "ρ", "Sensitivity to interest rates. Usually the smallest greek."),
    ]

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("The Greeks")
                .font(.caption.weight(.medium))
                .foregroundStyle(Theme.textTertiary)

            LazyVGrid(columns: Array(repeating: GridItem(.flexible(), spacing: 8), count: 5), spacing: 8) {
                ForEach(greekInfo, id: \.key) { greek in
                    let value = greekValue(greek.key)
                    Button {
                        withAnimation { expandedGreek = expandedGreek == greek.key ? nil : greek.key }
                    } label: {
                        VStack(spacing: 4) {
                            Text(greek.symbol)
                                .font(.caption2)
                                .foregroundStyle(Theme.textTertiary)
                            Text(value.formatted(.number.precision(.fractionLength(4))))
                                .font(.caption.monospacedDigit().weight(.semibold))
                                .foregroundStyle(greekColor(greek.key, value: value))
                        }
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 8)
                        .background(Theme.surface)
                        .clipShape(RoundedRectangle(cornerRadius: 8))
                    }
                    .buttonStyle(.plain)
                }
            }

            if let key = expandedGreek, let info = greekInfo.first(where: { $0.key == key }) {
                HStack(alignment: .top, spacing: 8) {
                    Image(systemName: "info.circle.fill")
                        .foregroundStyle(Theme.accent)
                        .font(.caption)
                    VStack(alignment: .leading, spacing: 2) {
                        Text(info.label)
                            .font(.caption.weight(.semibold))
                            .foregroundStyle(Theme.textPrimary)
                        Text(info.description)
                            .font(.caption2)
                            .foregroundStyle(Theme.textSecondary)
                    }
                }
                .padding()
                .background(Theme.accent.opacity(0.05))
                .clipShape(RoundedRectangle(cornerRadius: 8))
                .transition(.opacity.combined(with: .move(edge: .top)))
            }
        }
    }

    private func greekValue(_ key: String) -> Decimal {
        switch key {
        case "delta": return contract.delta
        case "gamma": return contract.gamma
        case "theta": return contract.theta
        case "vega": return contract.vega
        case "rho": return contract.rho
        default: return 0
        }
    }

    private func greekColor(_ key: String, value: Decimal) -> Color {
        if key == "theta" && value < 0 { return .red }
        if key == "delta" { return value >= 0 ? .green : .red }
        return Theme.textPrimary
    }
}
