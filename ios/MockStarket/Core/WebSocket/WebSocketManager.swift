import Foundation
import Observation

@Observable
final class WebSocketManager: @unchecked Sendable {
    static let shared = WebSocketManager()

    enum ConnectionState {
        case disconnected
        case connecting
        case connected
    }

    private(set) var state: ConnectionState = .disconnected
    private var webSocketTask: URLSessionWebSocketTask?
    private var session: URLSession
    private var baseURL: String
    private var reconnectAttempts = 0
    private let maxReconnectAttempts = 10

    private var priceUpdateHandler: (([PriceUpdate]) -> Void)?
    private var portfolioUpdateHandler: ((PortfolioWSUpdate) -> Void)?
    private var tradeExecutedHandler: ((Trade) -> Void)?
    private var marketEventHandler: ((MarketEventWS) -> Void)?
    private var achievementHandler: ((Achievement) -> Void)?
    private var alertHandler: ((PriceAlertTrigger) -> Void)?

    init(baseURL: String = "ws://localhost:8080") {
        self.baseURL = baseURL
        self.session = URLSession(configuration: .default)
    }

    func setBaseURL(_ url: String) {
        self.baseURL = url
    }

    func connect(token: String) {
        guard state == .disconnected else { return }

        state = .connecting
        guard let url = URL(string: "\(baseURL)/ws?user_id=\(token)") else { return }

        webSocketTask = session.webSocketTask(with: url)
        webSocketTask?.resume()
        state = .connected
        reconnectAttempts = 0

        subscribe(to: "market")
        subscribe(to: "portfolio")

        receiveMessages()
        startPing()
    }

    func disconnect() {
        webSocketTask?.cancel(with: .normalClosure, reason: nil)
        webSocketTask = nil
        state = .disconnected
    }

    // MARK: - Subscriptions

    func subscribe(to channel: String) {
        let message: [String: String] = ["type": "subscribe", "channel": channel]
        send(message)
    }

    func unsubscribe(from channel: String) {
        let message: [String: String] = ["type": "unsubscribe", "channel": channel]
        send(message)
    }

    // MARK: - Handlers

    func onPriceUpdate(_ handler: @escaping ([PriceUpdate]) -> Void) {
        priceUpdateHandler = handler
    }

    func onPortfolioUpdate(_ handler: @escaping (PortfolioWSUpdate) -> Void) {
        portfolioUpdateHandler = handler
    }

    func onTradeExecuted(_ handler: @escaping (Trade) -> Void) {
        tradeExecutedHandler = handler
    }

    func onMarketEvent(_ handler: @escaping (MarketEventWS) -> Void) {
        marketEventHandler = handler
    }

    func onAchievement(_ handler: @escaping (Achievement) -> Void) {
        achievementHandler = handler
    }

    func onAlertTriggered(_ handler: @escaping (PriceAlertTrigger) -> Void) {
        alertHandler = handler
    }

    // MARK: - Private

    private func send(_ message: [String: String]) {
        guard let data = try? JSONEncoder().encode(message),
              let string = String(data: data, encoding: .utf8) else { return }
        webSocketTask?.send(.string(string)) { _ in }
    }

    private func receiveMessages() {
        webSocketTask?.receive { [weak self] result in
            switch result {
            case .success(let message):
                self?.handleMessage(message)
                self?.receiveMessages()
            case .failure:
                self?.handleDisconnect()
            }
        }
    }

    private func handleMessage(_ message: URLSessionWebSocketTask.Message) {
        guard case .string(let text) = message,
              let data = text.data(using: .utf8) else { return }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601

        guard let envelope = try? decoder.decode(WSEnvelope.self, from: data) else { return }

        switch envelope.type {
        case "price_batch":
            if let updates = try? decoder.decode([PriceUpdate].self, from: envelope.data ?? Data()) {
                priceUpdateHandler?(updates)
            }
        case "portfolio_update":
            if let update = try? decoder.decode(PortfolioWSUpdate.self, from: envelope.data ?? Data()) {
                portfolioUpdateHandler?(update)
            }
        case "trade_executed":
            if let trade = try? decoder.decode(Trade.self, from: envelope.data ?? Data()) {
                tradeExecutedHandler?(trade)
            }
        case "market_event":
            if let event = try? decoder.decode(MarketEventWS.self, from: envelope.data ?? Data()) {
                marketEventHandler?(event)
            }
        case "achievement_unlocked":
            if let achievement = try? decoder.decode(Achievement.self, from: envelope.data ?? Data()) {
                achievementHandler?(achievement)
            }
        case "alert_triggered":
            if let alert = try? decoder.decode(PriceAlertTrigger.self, from: envelope.data ?? Data()) {
                alertHandler?(alert)
            }
        default:
            break
        }
    }

    private func handleDisconnect() {
        state = .disconnected
        attemptReconnect()
    }

    private func attemptReconnect() {
        guard reconnectAttempts < maxReconnectAttempts else { return }
        reconnectAttempts += 1
        let delay = min(pow(2.0, Double(reconnectAttempts)), 30.0)

        Task {
            try? await Task.sleep(for: .seconds(delay))
            // Would need stored token to reconnect
        }
    }

    private func startPing() {
        Task {
            while state == .connected {
                try? await Task.sleep(for: .seconds(30))
                webSocketTask?.sendPing { _ in }
            }
        }
    }
}

// MARK: - WebSocket Message Types

struct WSEnvelope: Codable {
    let type: String
    let timestamp: Int64?
    let data: Data?

    enum CodingKeys: String, CodingKey {
        case type, timestamp, data
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        type = try container.decode(String.self, forKey: .type)
        timestamp = try container.decodeIfPresent(Int64.self, forKey: .timestamp)
        // Decode data as raw JSON
        if container.contains(.data) {
            let rawData = try container.decode(AnyCodable.self, forKey: .data)
            data = try JSONEncoder().encode(rawData)
        } else {
            data = nil
        }
    }
}

// Helper for raw JSON passthrough
struct AnyCodable: Codable {
    let value: Any

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let dict = try? container.decode([String: AnyCodable].self) {
            value = dict
        } else if let array = try? container.decode([AnyCodable].self) {
            value = array
        } else if let string = try? container.decode(String.self) {
            value = string
        } else if let number = try? container.decode(Double.self) {
            value = number
        } else if let bool = try? container.decode(Bool.self) {
            value = bool
        } else {
            value = NSNull()
        }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        switch value {
        case let dict as [String: AnyCodable]:
            try container.encode(dict)
        case let array as [AnyCodable]:
            try container.encode(array)
        case let string as String:
            try container.encode(string)
        case let number as Double:
            try container.encode(number)
        case let bool as Bool:
            try container.encode(bool)
        default:
            try container.encodeNil()
        }
    }
}

struct PriceUpdate: Codable, Sendable {
    let ticker: String
    let price: Decimal
    let change: Decimal
    let changePct: Decimal
    let volume: Int64
    let high: Decimal
    let low: Decimal

    enum CodingKeys: String, CodingKey {
        case ticker, price, change, volume, high, low
        case changePct = "change_pct"
    }
}

struct PortfolioWSUpdate: Codable, Sendable {
    let cash: Decimal
    let netWorth: Decimal

    enum CodingKeys: String, CodingKey {
        case cash
        case netWorth = "net_worth"
    }
}

struct MarketEventWS: Codable, Sendable {
    let event: String
    let ticker: String?
    let sector: String?
    let headline: String
    let impact: String
    let magnitude: String
}

struct PriceAlertTrigger: Codable, Sendable {
    let alertID: UUID
    let ticker: String
    let condition: String
    let targetPrice: Decimal
    let currentPrice: Decimal

    enum CodingKeys: String, CodingKey {
        case alertID = "alert_id"
        case ticker, condition
        case targetPrice = "target_price"
        case currentPrice = "current_price"
    }
}
