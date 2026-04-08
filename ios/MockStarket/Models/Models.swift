import Foundation

// MARK: - User

struct User: Codable, Identifiable, Sendable {
    let id: UUID
    let firebaseUID: String
    var displayName: String
    var avatarURL: String?
    let isGuest: Bool
    let createdAt: Date
    let updatedAt: Date
    var loginStreak: Int
    var longestStreak: Int

    enum CodingKeys: String, CodingKey {
        case id
        case firebaseUID = "firebase_uid"
        case displayName = "display_name"
        case avatarURL = "avatar_url"
        case isGuest = "is_guest"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case loginStreak = "login_streak"
        case longestStreak = "longest_streak"
    }
}

// MARK: - Stock

struct Stock: Codable, Identifiable, Hashable, Sendable {
    var id: String { ticker }
    let ticker: String
    let name: String
    let sector: String
    var assetType: String = "stock" // stock, etf, crypto, commodity
    let basePrice: Decimal
    var currentPrice: Decimal
    var dayOpen: Decimal
    var dayHigh: Decimal
    var dayLow: Decimal
    var prevClose: Decimal
    var volume: Int64
    let volatility: Decimal
    let description: String?
    let logoURL: String?

    /// Display-friendly ticker: strips "X:" prefix and "USD" suffix from crypto tickers.
    var displayTicker: String {
        if ticker.hasPrefix("X:") && ticker.hasSuffix("USD") {
            return String(ticker.dropFirst(2).dropLast(3))
        }
        return ticker
    }

    enum CodingKeys: String, CodingKey {
        case ticker, name, sector, volume, volatility, description
        case logoURL = "logo_url"
        case assetType = "asset_type"
        case basePrice = "base_price"
        case currentPrice = "current_price"
        case dayOpen = "day_open"
        case dayHigh = "day_high"
        case dayLow = "day_low"
        case prevClose = "prev_close"
    }

    init(ticker: String, name: String, sector: String, assetType: String = "stock", basePrice: Decimal, currentPrice: Decimal, dayOpen: Decimal, dayHigh: Decimal, dayLow: Decimal, prevClose: Decimal, volume: Int64, volatility: Decimal, description: String? = nil, logoURL: String? = nil) {
        self.ticker = ticker; self.name = name; self.sector = sector; self.assetType = assetType
        self.basePrice = basePrice; self.currentPrice = currentPrice; self.dayOpen = dayOpen
        self.dayHigh = dayHigh; self.dayLow = dayLow; self.prevClose = prevClose
        self.volume = volume; self.volatility = volatility; self.description = description
        self.logoURL = logoURL
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        ticker = try c.decode(String.self, forKey: .ticker)
        name = try c.decode(String.self, forKey: .name)
        sector = try c.decode(String.self, forKey: .sector)
        assetType = try c.decodeIfPresent(String.self, forKey: .assetType) ?? "stock"
        basePrice = try Self.decimalFromStringOrNumber(c, forKey: .basePrice)
        currentPrice = try Self.decimalFromStringOrNumber(c, forKey: .currentPrice)
        dayOpen = try Self.decimalFromStringOrNumber(c, forKey: .dayOpen)
        dayHigh = try Self.decimalFromStringOrNumber(c, forKey: .dayHigh)
        dayLow = try Self.decimalFromStringOrNumber(c, forKey: .dayLow)
        prevClose = try Self.decimalFromStringOrNumber(c, forKey: .prevClose)
        volume = try c.decode(Int64.self, forKey: .volume)
        volatility = try Self.decimalFromStringOrNumber(c, forKey: .volatility)
        description = try c.decodeIfPresent(String.self, forKey: .description)
        logoURL = try c.decodeIfPresent(String.self, forKey: .logoURL)
    }

    /// Decodes a Decimal that may come as a JSON string or number.
    private static func decimalFromStringOrNumber(_ container: KeyedDecodingContainer<CodingKeys>, forKey key: CodingKeys) throws -> Decimal {
        if let str = try? container.decode(String.self, forKey: key), let d = Decimal(string: str) {
            return d
        }
        return try container.decode(Decimal.self, forKey: key)
    }

    var change: Decimal { currentPrice - dayOpen }
    var changePct: Decimal {
        guard dayOpen != 0 else { return 0 }
        return (change / dayOpen) * 100
    }
    var isUp: Bool { change >= 0 }
}

// MARK: - ETF Holding

struct ETFHolding: Codable, Identifiable, Sendable {
    var id: String { ticker }
    let ticker: String
    let name: String
    let weight: String
    let price: String
}

// MARK: - Portfolio

struct Portfolio: Codable, Sendable {
    let id: UUID
    let userID: UUID
    var cash: Decimal
    var netWorth: Decimal
    let createdAt: Date
    let updatedAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case cash
        case netWorth = "net_worth"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

// MARK: - Holding / Position

struct Holding: Codable, Identifiable, Sendable {
    let id: UUID
    let portfolioID: UUID
    let ticker: String
    let shares: Int
    let avgCost: Decimal
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case portfolioID = "portfolio_id"
        case ticker, shares
        case avgCost = "avg_cost"
        case createdAt = "created_at"
    }
}

struct Position: Codable, Identifiable, Sendable {
    let id: UUID
    let portfolioID: UUID
    let ticker: String
    let shares: Int
    let avgCost: Decimal
    let currentPrice: Decimal
    let marketValue: Decimal
    let pnl: Decimal
    let pnlPct: Decimal

    enum CodingKeys: String, CodingKey {
        case id
        case portfolioID = "portfolio_id"
        case ticker, shares
        case avgCost = "avg_cost"
        case currentPrice = "current_price"
        case marketValue = "market_value"
        case pnl
        case pnlPct = "pnl_pct"
    }

    var isProfit: Bool { pnl >= 0 }
}

// MARK: - Trade

struct Trade: Codable, Identifiable, Sendable {
    let id: UUID
    let userID: UUID
    let ticker: String
    let side: String
    let shares: Int
    let price: Decimal
    let total: Decimal
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case ticker, side, shares, price, total
        case createdAt = "created_at"
    }

    var isBuy: Bool { side == "buy" }
}

// MARK: - Order

struct Order: Codable, Identifiable, Sendable {
    let id: UUID
    let ticker: String
    let side: String
    let orderType: String
    let shares: Int
    let limitPrice: Decimal?
    let stopPrice: Decimal?
    let status: String
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id, ticker, side, shares, status
        case orderType = "order_type"
        case limitPrice = "limit_price"
        case stopPrice = "stop_price"
        case createdAt = "created_at"
    }
}

// MARK: - Leaderboard

struct LeaderboardEntry: Codable, Identifiable, Sendable {
    let id: Int64
    let userID: UUID
    let displayName: String
    let netWorth: Decimal
    let totalReturn: Decimal
    let rank: Int
    let period: String

    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case displayName = "display_name"
        case netWorth = "net_worth"
        case totalReturn = "total_return"
        case rank, period
    }
}

// MARK: - Achievement

struct Achievement: Codable, Identifiable, Sendable {
    let id: String
    let name: String
    let description: String
    let icon: String
    let category: String
}

struct UserAchievement: Codable, Identifiable, Sendable {
    let id: UUID
    let achievementID: String
    let earnedAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case achievementID = "achievement_id"
        case earnedAt = "earned_at"
    }
}

// MARK: - Price Alert

struct PriceAlert: Codable, Identifiable, Sendable {
    let id: UUID
    let ticker: String
    let condition: String
    let targetPrice: Decimal
    let triggered: Bool
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id, ticker, condition, triggered
        case targetPrice = "target_price"
        case createdAt = "created_at"
    }
}

// MARK: - Daily Challenge

struct DailyChallenge: Codable, Identifiable, Sendable {
    let id: UUID
    let date: String
    let challengeType: String
    let description: String
    let targetJSON: String
    let rewardCash: Decimal
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id, date, description
        case challengeType = "challenge_type"
        case targetJSON = "target_json"
        case rewardCash = "reward_cash"
        case createdAt = "created_at"
    }
}

struct UserChallenge: Codable, Identifiable, Sendable {
    let id: UUID
    let userID: UUID
    let challengeID: UUID
    let completed: Bool
    let completedAt: Date?
    let claimed: Bool

    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case challengeID = "challenge_id"
        case completed
        case completedAt = "completed_at"
        case claimed
    }
}

struct ChallengeResponse: Codable, Sendable {
    let challenge: DailyChallenge
    let progress: UserChallenge?
}

struct ChallengeCheckResponse: Codable, Sendable {
    let completed: Bool
}

// MARK: - Options

struct OptionContract: Codable, Identifiable, Hashable, Sendable {
    let id: UUID
    let ticker: String
    let optionType: String
    let strikePrice: Decimal
    let expiration: Date
    let contractSymbol: String
    var bidPrice: Decimal
    var askPrice: Decimal
    var lastPrice: Decimal
    var markPrice: Decimal
    var openInterest: Int
    var volume: Int
    var impliedVol: Decimal
    var delta: Decimal
    var gamma: Decimal
    var theta: Decimal
    var vega: Decimal
    var rho: Decimal
    let status: String

    enum CodingKeys: String, CodingKey {
        case id, ticker, expiration, status, volume, delta, gamma, theta, vega, rho
        case optionType = "option_type"
        case strikePrice = "strike_price"
        case contractSymbol = "contract_symbol"
        case bidPrice = "bid_price"
        case askPrice = "ask_price"
        case lastPrice = "last_price"
        case markPrice = "mark_price"
        case openInterest = "open_interest"
        case impliedVol = "implied_vol"
    }

    var isCall: Bool { optionType == "call" }
    var isPut: Bool { optionType == "put" }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(UUID.self, forKey: .id)
        ticker = try c.decode(String.self, forKey: .ticker)
        optionType = try c.decode(String.self, forKey: .optionType)
        expiration = try c.decode(Date.self, forKey: .expiration)
        contractSymbol = try c.decode(String.self, forKey: .contractSymbol)
        openInterest = try c.decode(Int.self, forKey: .openInterest)
        volume = try c.decode(Int.self, forKey: .volume)
        status = try c.decode(String.self, forKey: .status)
        strikePrice = Self.dec(c, .strikePrice)
        bidPrice = Self.dec(c, .bidPrice)
        askPrice = Self.dec(c, .askPrice)
        lastPrice = Self.dec(c, .lastPrice)
        markPrice = Self.dec(c, .markPrice)
        impliedVol = Self.dec(c, .impliedVol)
        delta = Self.dec(c, .delta)
        gamma = Self.dec(c, .gamma)
        theta = Self.dec(c, .theta)
        vega = Self.dec(c, .vega)
        rho = Self.dec(c, .rho)
    }

    private static func dec(_ c: KeyedDecodingContainer<CodingKeys>, _ key: CodingKeys) -> Decimal {
        if let s = try? c.decode(String.self, forKey: key), let d = Decimal(string: s) { return d }
        return (try? c.decode(Decimal.self, forKey: key)) ?? 0
    }
}

struct OptionChainResponse: Codable, Sendable {
    let ticker: String
    let underlyingPrice: Decimal
    let calls: [OptionContract]
    let puts: [OptionContract]

    enum CodingKeys: String, CodingKey {
        case ticker, calls, puts
        case underlyingPrice = "underlying_price"
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        ticker = try c.decode(String.self, forKey: .ticker)
        calls = try c.decode([OptionContract].self, forKey: .calls)
        puts = try c.decode([OptionContract].self, forKey: .puts)
        if let s = try? c.decode(String.self, forKey: .underlyingPrice), let d = Decimal(string: s) {
            underlyingPrice = d
        } else {
            underlyingPrice = (try? c.decode(Decimal.self, forKey: .underlyingPrice)) ?? 0
        }
    }
}

struct OptionPosition: Codable, Identifiable, Sendable {
    let id: UUID
    let portfolioID: UUID
    let contractID: UUID
    let quantity: Int
    let avgCost: Decimal
    let collateral: Decimal
    let contract: OptionContract
    let marketValue: Decimal
    let pnl: Decimal
    let pnlPct: Decimal
    let isLong: Bool

    enum CodingKeys: String, CodingKey {
        case id, quantity, collateral, contract
        case portfolioID = "portfolio_id"
        case contractID = "contract_id"
        case avgCost = "avg_cost"
        case marketValue = "market_value"
        case pnl
        case pnlPct = "pnl_pct"
        case isLong = "is_long"
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(UUID.self, forKey: .id)
        portfolioID = try c.decode(UUID.self, forKey: .portfolioID)
        contractID = try c.decode(UUID.self, forKey: .contractID)
        quantity = try c.decode(Int.self, forKey: .quantity)
        contract = try c.decode(OptionContract.self, forKey: .contract)
        isLong = try c.decode(Bool.self, forKey: .isLong)
        avgCost = Self.dec(c, .avgCost)
        collateral = Self.dec(c, .collateral)
        marketValue = Self.dec(c, .marketValue)
        pnl = Self.dec(c, .pnl)
        pnlPct = Self.dec(c, .pnlPct)
    }

    private static func dec(_ c: KeyedDecodingContainer<CodingKeys>, _ key: CodingKeys) -> Decimal {
        if let s = try? c.decode(String.self, forKey: key), let d = Decimal(string: s) { return d }
        return (try? c.decode(Decimal.self, forKey: key)) ?? 0
    }
}

struct OptionTrade: Codable, Identifiable, Sendable {
    let id: UUID
    let userID: UUID
    let contractID: UUID
    let side: String
    let quantity: Int
    let price: Decimal
    let total: Decimal
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id, side, quantity, price, total
        case userID = "user_id"
        case contractID = "contract_id"
        case createdAt = "created_at"
    }
}

// MARK: - Price Point (Charts)

struct PricePoint: Codable, Identifiable, Sendable {
    let id: Int64
    let ticker: String
    let price: Decimal
    let open: Decimal
    let high: Decimal
    let low: Decimal
    let close: Decimal
    let volume: Int64
    let interval: String
    let recordedAt: Date

    enum CodingKeys: String, CodingKey {
        case id, ticker, price, open, high, low, close, volume, interval
        case recordedAt = "recorded_at"
    }
}

// MARK: - Portfolio Response

struct PortfolioResponse: Codable, Sendable {
    let portfolio: Portfolio
    let positions: [Position]
    let netWorth: Decimal
    let invested: Decimal

    enum CodingKeys: String, CodingKey {
        case portfolio, positions
        case netWorth = "net_worth"
        case invested
    }
}

// MARK: - Market Summary

struct MarketSummary: Codable, Sendable {
    let indexValue: Decimal
    let indexChangePct: Decimal
    let totalStocks: Int
    let gainers: Int
    let losers: Int

    enum CodingKeys: String, CodingKey {
        case indexValue = "index_value"
        case indexChangePct = "index_change_pct"
        case totalStocks = "total_stocks"
        case gainers, losers
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        indexValue = Self.dec(c, .indexValue)
        indexChangePct = Self.dec(c, .indexChangePct)
        totalStocks = try c.decode(Int.self, forKey: .totalStocks)
        gainers = try c.decode(Int.self, forKey: .gainers)
        losers = try c.decode(Int.self, forKey: .losers)
    }

    private static func dec(_ c: KeyedDecodingContainer<CodingKeys>, _ key: CodingKeys) -> Decimal {
        if let s = try? c.decode(String.self, forKey: key), let d = Decimal(string: s) { return d }
        return (try? c.decode(Decimal.self, forKey: key)) ?? 0
    }
}
