import SwiftUI

struct AlertsView: View {
    @State private var viewModel = AlertsViewModel()
    @State private var showCreateSheet = false

    var body: some View {
        List {
            if viewModel.alerts.isEmpty && !viewModel.isLoading {
                ContentUnavailableView(
                    "No Alerts",
                    systemImage: "bell.slash",
                    description: Text("Create a price alert to get notified when a stock hits your target price.")
                )
                .listRowBackground(Color.clear)
            }

            ForEach(viewModel.alerts) { alert in
                alertRow(alert)
                    .listRowBackground(Theme.surface)
                    .swipeActions(edge: .trailing) {
                        Button(role: .destructive) {
                            Task { await viewModel.deleteAlert(alert) }
                        } label: {
                            Label("Delete", systemImage: "trash")
                        }
                    }
            }
        }
        .scrollContentBackground(.hidden)
        .background(Theme.background)
        .navigationTitle("Price Alerts")
        .toolbar {
            ToolbarItem(placement: .topBarTrailing) {
                Button {
                    showCreateSheet = true
                } label: {
                    Image(systemName: "plus")
                }
            }
        }
        .sheet(isPresented: $showCreateSheet) {
            createAlertSheet
        }
        .task {
            await viewModel.load()
        }
    }

    private func alertRow(_ alert: PriceAlert) -> some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                HStack(spacing: 6) {
                    Text(alert.ticker)
                        .font(.headline.monospaced())
                        .foregroundStyle(Theme.accent)

                    Image(systemName: alert.condition == "above" ? "arrow.up" : "arrow.down")
                        .font(.caption)
                        .foregroundStyle(alert.condition == "above" ? Theme.positive : Theme.negative)
                }

                Text("\(alert.condition == "above" ? "Above" : "Below") \(alert.targetPrice.formatted(.currency(code: "USD")))")
                    .font(.subheadline)
                    .foregroundStyle(Theme.textSecondary)
            }

            Spacer()

            if alert.triggered {
                Label("Triggered", systemImage: "bell.fill")
                    .font(.caption)
                    .foregroundStyle(Theme.positive)
                    .padding(.horizontal, 10)
                    .padding(.vertical, 5)
                    .background(Theme.positive.opacity(0.1))
                    .clipShape(Capsule())
            } else {
                Label("Active", systemImage: "bell")
                    .font(.caption)
                    .foregroundStyle(Theme.textTertiary)
            }
        }
        .padding(.vertical, 4)
    }

    private var createAlertSheet: some View {
        NavigationStack {
            Form {
                Section("Stock") {
                    TextField("Ticker (e.g. PIPE)", text: $viewModel.ticker)
                        .textInputAutocapitalization(.characters)
                }

                Section("Condition") {
                    Picker("Trigger when price is", selection: $viewModel.condition) {
                        Text("Above").tag("above")
                        Text("Below").tag("below")
                    }
                    .pickerStyle(.segmented)

                    TextField("Target Price", text: $viewModel.targetPrice)
                        .keyboardType(.decimalPad)
                }

                if let error = viewModel.errorMessage {
                    Section {
                        Text(error)
                            .foregroundStyle(Theme.negative)
                            .font(.caption)
                    }
                }
            }
            .scrollContentBackground(.hidden)
            .background(Theme.background)
            .navigationTitle("New Alert")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") { showCreateSheet = false }
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("Create") {
                        Task {
                            if await viewModel.createAlert() {
                                showCreateSheet = false
                            }
                        }
                    }
                    .disabled(viewModel.isCreating)
                }
            }
        }
        .presentationDetents([.medium])
    }
}

#Preview {
    NavigationStack {
        AlertsView()
    }
    .environment(AppState())
    .preferredColorScheme(.dark)
}
