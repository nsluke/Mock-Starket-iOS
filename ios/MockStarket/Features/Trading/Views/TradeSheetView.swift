import SwiftUI

struct TradeSheetView: View {
    let ticker: String
    let side: String
    let currentPrice: Decimal
    var currentShares: Int = 0

    @State private var shares = 1
    @State private var useDollars = false
    @State private var dollarAmount = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var success = false
    @State private var confirming = false
    @Environment(\.dismiss) private var dismiss

    private let apiClient = APIClient.shared

    private var estimatedTotal: Decimal {
        currentPrice * Decimal(shares)
    }

    private var estimatedSharesFromDollars: Int {
        guard let amount = Decimal(string: dollarAmount), amount > 0, currentPrice > 0 else { return 0 }
        return NSDecimalNumber(decimal: amount / currentPrice).intValue
    }

    private var effectiveShares: Int {
        useDollars ? estimatedSharesFromDollars : shares
    }

    var body: some View {
        NavigationStack {
            ZStack {
                Theme.background.ignoresSafeArea()

                if success {
                    successView
                } else if confirming {
                    confirmView
                } else {
                    inputView
                }
            }
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button { dismiss() } label: {
                        Image(systemName: "xmark")
                            .foregroundStyle(Theme.textSecondary)
                    }
                }
            }
        }
    }

    // MARK: - Input View

    private var inputView: some View {
        VStack(spacing: 0) {
            Spacer()

            // Side indicator
            Text(side == "buy" ? "Buy \(ticker)" : "Sell \(ticker)")
                .font(.title3.weight(.semibold))
                .foregroundStyle(Theme.textPrimary)

            Text(currentPrice.formatted(.currency(code: "USD")))
                .font(.caption)
                .foregroundStyle(Theme.textSecondary)
                .padding(.top, 2)

            Spacer()

            // Amount display
            if useDollars {
                VStack(spacing: 4) {
                    Text("$\(dollarAmount.isEmpty ? "0" : dollarAmount)")
                        .font(.system(size: 56, weight: .bold, design: .rounded))
                        .foregroundStyle(Theme.textPrimary)

                    if estimatedSharesFromDollars > 0 {
                        Text("≈ \(estimatedSharesFromDollars) share\(estimatedSharesFromDollars == 1 ? "" : "s")")
                            .font(.subheadline)
                            .foregroundStyle(Theme.textTertiary)
                    }
                }
            } else {
                VStack(spacing: 4) {
                    HStack(spacing: 24) {
                        Button {
                            if shares > 1 { shares -= 1 }
                        } label: {
                            Image(systemName: "minus.circle.fill")
                                .font(.title)
                                .foregroundStyle(Theme.textSecondary)
                        }

                        Text("\(shares)")
                            .font(.system(size: 56, weight: .bold, design: .rounded))
                            .foregroundStyle(Theme.textPrimary)
                            .frame(minWidth: 80)

                        Button {
                            shares += 1
                        } label: {
                            Image(systemName: "plus.circle.fill")
                                .font(.title)
                                .foregroundStyle(Theme.accent)
                        }
                    }

                    Text("share\(shares == 1 ? "" : "s")")
                        .font(.subheadline)
                        .foregroundStyle(Theme.textTertiary)
                }
            }

            Spacer()

            // Quick buttons
            HStack(spacing: 8) {
                ForEach([1, 5, 10, 25], id: \.self) { qty in
                    quickButton("\(qty)") { shares = qty; useDollars = false }
                }

                if side == "sell" && currentShares > 0 {
                    quickButton("All") {
                        shares = currentShares
                        useDollars = false
                    }
                }

                quickButton("$") {
                    useDollars.toggle()
                    dollarAmount = ""
                }
            }
            .padding(.horizontal, 24)

            // Dollar keypad (when in dollar mode)
            if useDollars {
                dollarKeypad
                    .padding(.top, 16)
            }

            Spacer()

            // Estimated total
            if effectiveShares > 0 {
                HStack {
                    Text("Estimated cost")
                        .foregroundStyle(Theme.textSecondary)
                    Spacer()
                    Text((currentPrice * Decimal(effectiveShares)).formatted(.currency(code: "USD")))
                        .fontWeight(.semibold)
                }
                .font(.subheadline)
                .padding(.horizontal, 24)
                .padding(.bottom, 12)
            }

            // Review button
            Button {
                confirming = true
            } label: {
                Text("Review Order")
                    .font(.headline)
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 18)
                    .background(side == "buy" ? Theme.positive : Theme.negative)
                    .foregroundStyle(.white)
                    .clipShape(RoundedRectangle(cornerRadius: 16))
            }
            .disabled(effectiveShares <= 0)
            .opacity(effectiveShares > 0 ? 1 : 0.4)
            .padding(.horizontal, 24)
            .padding(.bottom, 16)

            if let error = errorMessage {
                Text(error)
                    .font(.caption)
                    .foregroundStyle(Theme.negative)
                    .padding(.bottom, 8)
            }
        }
    }

    // MARK: - Confirm View

    private var confirmView: some View {
        VStack(spacing: 24) {
            Spacer()

            Image(systemName: side == "buy" ? "arrow.up.circle.fill" : "arrow.down.circle.fill")
                .font(.system(size: 64))
                .foregroundStyle(side == "buy" ? Theme.positive : Theme.negative)

            VStack(spacing: 8) {
                Text("\(side == "buy" ? "Buy" : "Sell") \(effectiveShares) share\(effectiveShares == 1 ? "" : "s") of \(ticker)")
                    .font(.title3.weight(.semibold))

                Text("at \(currentPrice.formatted(.currency(code: "USD"))) per share")
                    .foregroundStyle(Theme.textSecondary)

                Text("Total: \((currentPrice * Decimal(effectiveShares)).formatted(.currency(code: "USD")))")
                    .font(.title2.weight(.bold))
                    .padding(.top, 8)
            }

            Spacer()

            // Swipe to confirm (simplified as button)
            Button {
                Task { await executeTrade() }
            } label: {
                Group {
                    if isLoading {
                        ProgressView().tint(.white)
                    } else {
                        Label("Confirm \(side == "buy" ? "Purchase" : "Sale")", systemImage: "checkmark.circle.fill")
                            .font(.headline)
                    }
                }
                .frame(maxWidth: .infinity)
                .padding(.vertical, 18)
                .background(side == "buy" ? Theme.positive : Theme.negative)
                .foregroundStyle(.white)
                .clipShape(RoundedRectangle(cornerRadius: 16))
            }
            .disabled(isLoading)
            .padding(.horizontal, 24)

            Button("Go Back") {
                confirming = false
            }
            .foregroundStyle(Theme.textSecondary)
            .padding(.bottom, 16)
        }
    }

    // MARK: - Success View

    private var successView: some View {
        VStack(spacing: 24) {
            Spacer()

            Image(systemName: "checkmark.circle.fill")
                .font(.system(size: 80))
                .foregroundStyle(Theme.positive)

            Text("Order Placed!")
                .font(.title.weight(.bold))

            Text("\(side == "buy" ? "Bought" : "Sold") \(effectiveShares) share\(effectiveShares == 1 ? "" : "s") of \(ticker)")
                .foregroundStyle(Theme.textSecondary)

            Spacer()

            Button {
                dismiss()
            } label: {
                Text("Done")
                    .font(.headline)
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 18)
                    .background(Theme.accent)
                    .foregroundStyle(.black)
                    .clipShape(RoundedRectangle(cornerRadius: 16))
            }
            .padding(.horizontal, 24)
            .padding(.bottom, 16)
        }
    }

    // MARK: - Components

    private func quickButton(_ label: String, action: @escaping () -> Void) -> some View {
        Button(action: action) {
            Text(label)
                .font(.caption.weight(.semibold))
                .padding(.horizontal, 14)
                .padding(.vertical, 8)
                .background(Theme.surfaceElevated)
                .foregroundStyle(Theme.textSecondary)
                .clipShape(Capsule())
        }
    }

    private var dollarKeypad: some View {
        let keys = [
            ["1", "2", "3"],
            ["4", "5", "6"],
            ["7", "8", "9"],
            [".", "0", "⌫"],
        ]

        return VStack(spacing: 8) {
            ForEach(keys, id: \.self) { row in
                HStack(spacing: 8) {
                    ForEach(row, id: \.self) { key in
                        Button {
                            if key == "⌫" {
                                if !dollarAmount.isEmpty { dollarAmount.removeLast() }
                            } else if key == "." {
                                if !dollarAmount.contains(".") { dollarAmount += "." }
                            } else {
                                dollarAmount += key
                            }
                        } label: {
                            Text(key)
                                .font(.title3.weight(.medium))
                                .frame(maxWidth: .infinity)
                                .frame(height: 44)
                                .background(Theme.surfaceElevated)
                                .foregroundStyle(Theme.textPrimary)
                                .clipShape(RoundedRectangle(cornerRadius: 8))
                        }
                    }
                }
            }
        }
        .padding(.horizontal, 24)
    }

    // MARK: - Trade Execution

    private func executeTrade() async {
        isLoading = true
        errorMessage = nil

        struct TradeBody: Encodable {
            let ticker: String
            let side: String
            let shares: Int
        }

        do {
            let _: Trade = try await apiClient.request(
                .executeTrade,
                token: AuthManager.shared.currentToken,
                body: TradeBody(ticker: ticker, side: side, shares: effectiveShares)
            )
            success = true
        } catch {
            errorMessage = error.localizedDescription
            confirming = false
        }

        isLoading = false
    }
}

#Preview {
    TradeSheetView(ticker: "PIPE", side: "buy", currentPrice: 267.00, currentShares: 5)
        .preferredColorScheme(.dark)
}
