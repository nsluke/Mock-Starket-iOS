import SwiftUI

struct TradeSheetView: View {
    let ticker: String
    let side: String

    @State private var shares = 1
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var success = false
    @Environment(\.dismiss) private var dismiss

    private let apiClient = APIClient.shared

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                // Header
                VStack(spacing: 4) {
                    Text(side == "buy" ? "Buy" : "Sell")
                        .font(.title2.weight(.bold))
                        .foregroundStyle(side == "buy" ? Theme.positive : Theme.negative)
                    Text(ticker)
                        .font(.headline)
                        .foregroundStyle(Theme.textSecondary)
                }

                Spacer()

                // Share quantity
                VStack(spacing: 16) {
                    Text("Shares")
                        .font(.caption)
                        .foregroundStyle(Theme.textTertiary)

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

                    // Quick quantity buttons
                    HStack(spacing: 8) {
                        ForEach([1, 5, 10, 25, 50], id: \.self) { qty in
                            Button("\(qty)") {
                                shares = qty
                            }
                            .font(.caption.weight(.medium))
                            .padding(.horizontal, 12)
                            .padding(.vertical, 6)
                            .background(shares == qty ? Theme.accent.opacity(0.2) : Theme.surfaceElevated)
                            .foregroundStyle(shares == qty ? Theme.accent : Theme.textSecondary)
                            .clipShape(Capsule())
                        }
                    }
                }

                Spacer()

                if let error = errorMessage {
                    Text(error)
                        .font(.caption)
                        .foregroundStyle(Theme.negative)
                        .multilineTextAlignment(.center)
                }

                if success {
                    Label("Trade executed!", systemImage: "checkmark.circle.fill")
                        .font(.headline)
                        .foregroundStyle(Theme.positive)
                } else {
                    // Execute button
                    Button {
                        Task { await executeTrade() }
                    } label: {
                        Group {
                            if isLoading {
                                ProgressView()
                                    .tint(.white)
                            } else {
                                Text("\(side == "buy" ? "Buy" : "Sell") \(shares) Share\(shares == 1 ? "" : "s")")
                                    .font(.headline)
                            }
                        }
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 16)
                        .background(side == "buy" ? Theme.positive : Theme.negative)
                        .foregroundStyle(.white)
                        .clipShape(RoundedRectangle(cornerRadius: 14))
                    }
                    .disabled(isLoading)
                }
            }
            .padding(24)
            .background(Theme.background)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") { dismiss() }
                        .foregroundStyle(Theme.textSecondary)
                }
            }
        }
        .presentationDetents([.medium])
    }

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
                body: TradeBody(ticker: ticker, side: side, shares: shares)
            )
            success = true
            try? await Task.sleep(for: .seconds(1.5))
            dismiss()
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }
}

#Preview {
    TradeSheetView(ticker: "PLNX", side: "buy")
        .preferredColorScheme(.dark)
}
