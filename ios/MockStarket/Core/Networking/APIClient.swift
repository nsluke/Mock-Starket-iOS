import Foundation

enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(Int, String)
    case decodingError(Error)
    case networkError(Error)

    var errorDescription: String? {
        switch self {
        case .invalidURL: return "Invalid URL"
        case .invalidResponse: return "Invalid response"
        case .httpError(let code, let message): return "HTTP \(code): \(message)"
        case .decodingError(let error): return "Decoding error: \(error.localizedDescription)"
        case .networkError(let error): return "Network error: \(error.localizedDescription)"
        }
    }
}

enum APIEndpoint {
    // Auth
    case register
    case createGuest
    case getMe
    case updateMe
    case deleteMe

    // Stocks
    case listStocks
    case getStock(ticker: String)
    case getStockHistory(ticker: String, interval: String)
    case getETFHoldings(ticker: String)
    case marketSummary

    // Trading
    case executeTrade
    case getTradeHistory(limit: Int, offset: Int)

    // Orders
    case createOrder
    case listOrders
    case cancelOrder(id: UUID)

    // Portfolio
    case getPortfolio
    case getPortfolioHistory

    // Leaderboard
    case getLeaderboard(period: String)

    // Alerts
    case createAlert
    case listAlerts
    case deleteAlert(id: UUID)

    // Achievements
    case listAchievements
    case getMyAchievements

    // Watchlist
    case getWatchlist
    case addToWatchlist
    case removeFromWatchlist(ticker: String)

    // Challenges
    case getTodaysChallenge
    case checkChallenge
    case claimChallenge(id: UUID)

    // System
    case health

    var path: String {
        switch self {
        case .register: return "/api/v1/auth/register"
        case .createGuest: return "/api/v1/auth/guest"
        case .getMe, .updateMe, .deleteMe: return "/api/v1/auth/me"
        case .listStocks: return "/api/v1/stocks"
        case .getStock(let ticker): return "/api/v1/stocks/\(ticker)"
        case .getStockHistory(let ticker, let interval): return "/api/v1/stocks/\(ticker)/history?interval=\(interval)"
        case .getETFHoldings(let ticker): return "/api/v1/stocks/\(ticker)/holdings"
        case .marketSummary: return "/api/v1/stocks/market-summary"
        case .executeTrade, .getTradeHistory: return "/api/v1/trades"
        case .createOrder, .listOrders: return "/api/v1/orders"
        case .cancelOrder(let id): return "/api/v1/orders/\(id)"
        case .getPortfolio: return "/api/v1/portfolio"
        case .getPortfolioHistory: return "/api/v1/portfolio/history"
        case .getLeaderboard(let period): return "/api/v1/leaderboard?period=\(period)"
        case .createAlert, .listAlerts: return "/api/v1/alerts"
        case .deleteAlert(let id): return "/api/v1/alerts/\(id)"
        case .listAchievements: return "/api/v1/achievements"
        case .getMyAchievements: return "/api/v1/achievements/me"
        case .getWatchlist, .addToWatchlist: return "/api/v1/watchlist"
        case .removeFromWatchlist(let ticker): return "/api/v1/watchlist/\(ticker)"
        case .getTodaysChallenge: return "/api/v1/challenges/today"
        case .checkChallenge: return "/api/v1/challenges/check"
        case .claimChallenge(let id): return "/api/v1/challenges/\(id)/claim"
        case .health: return "/api/v1/system/health"
        }
    }

    var method: String {
        switch self {
        case .register, .createGuest, .executeTrade, .createOrder, .createAlert, .addToWatchlist, .checkChallenge, .claimChallenge:
            return "POST"
        case .updateMe:
            return "PUT"
        case .deleteMe, .cancelOrder, .deleteAlert, .removeFromWatchlist:
            return "DELETE"
        default:
            return "GET"
        }
    }
}

actor APIClient {
    static let shared = APIClient()

    private let session: URLSession
    private let decoder: JSONDecoder
    private var baseURL: String

    init(baseURL: String = "http://localhost:8080") {
        self.baseURL = baseURL
        self.session = URLSession.shared

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        self.decoder = decoder
    }

    func setBaseURL(_ url: String) {
        self.baseURL = url
    }

    func request<T: Decodable>(_ endpoint: APIEndpoint, token: String? = nil, body: Encodable? = nil) async throws -> T {
        guard let url = URL(string: baseURL + endpoint.path) else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = endpoint.method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if let token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body {
            request.httpBody = try JSONEncoder().encode(body)
        }

        let (data, response): (Data, URLResponse)
        do {
            (data, response) = try await session.data(for: request)
        } catch {
            throw APIError.networkError(error)
        }

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            let errorMessage = String(data: data, encoding: .utf8) ?? "Unknown error"
            throw APIError.httpError(httpResponse.statusCode, errorMessage)
        }

        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError(error)
        }
    }
}
