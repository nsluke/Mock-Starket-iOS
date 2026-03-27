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

struct Stock: Codable, Identifiable, Sendable {
    var id: String { ticker }
    let ticker: String
    let name: String
    let sector: String
    let basePrice: Decimal
    var currentPrice: Decimal
    var dayOpen: Decimal
    var dayHigh: Decimal
    var dayLow: Decimal
    var prevClose: Decimal
    var volume: Int64
    let volatility: Decimal
    let description: String?

    enum CodingKeys: String, CodingKey {
        case ticker, name, sector, volume, volatility, description
        case basePrice = "base_price"
        case currentPrice = "current_price"
        case dayOpen = "day_open"
        case dayHigh = "day_high"
        case dayLow = "day_low"
        case prevClose = "prev_close"
    }

    var change: Decimal { currentPrice - dayOpen }
    var changePct: Decimal {
        guard dayOpen != 0 else { return 0 }
        return (change / dayOpen) * 100
    }
    var isUp: Bool { change >= 0 }
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
}
