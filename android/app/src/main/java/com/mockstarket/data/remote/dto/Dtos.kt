package com.mockstarket.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class UserDto(
    val id: String,
    @Json(name = "firebase_uid") val firebaseUid: String,
    @Json(name = "display_name") val displayName: String,
    @Json(name = "avatar_url") val avatarUrl: String?,
    @Json(name = "is_guest") val isGuest: Boolean,
    @Json(name = "login_streak") val loginStreak: Int,
    @Json(name = "longest_streak") val longestStreak: Int,
)

@JsonClass(generateAdapter = true)
data class StockDto(
    val ticker: String,
    val name: String,
    val sector: String,
    @Json(name = "base_price") val basePrice: String,
    @Json(name = "current_price") val currentPrice: String,
    @Json(name = "day_open") val dayOpen: String,
    @Json(name = "day_high") val dayHigh: String,
    @Json(name = "day_low") val dayLow: String,
    @Json(name = "prev_close") val prevClose: String,
    val volume: Long,
    val volatility: String,
    val description: String?,
)

@JsonClass(generateAdapter = true)
data class PricePointDto(
    val id: Long,
    val ticker: String,
    val price: String,
    val open: String,
    val high: String,
    val low: String,
    val close: String,
    val volume: Long,
    val interval: String,
    @Json(name = "recorded_at") val recordedAt: String,
)

@JsonClass(generateAdapter = true)
data class MarketSummaryDto(
    @Json(name = "index_value") val indexValue: String,
    @Json(name = "index_change_pct") val indexChangePct: String,
    @Json(name = "total_stocks") val totalStocks: Int,
    val gainers: Int,
    val losers: Int,
)

@JsonClass(generateAdapter = true)
data class TradeDto(
    val id: String,
    val ticker: String,
    val side: String,
    val shares: Int,
    val price: String,
    val total: String,
    @Json(name = "created_at") val createdAt: String,
)

@JsonClass(generateAdapter = true)
data class OrderDto(
    val id: String,
    val ticker: String,
    val side: String,
    @Json(name = "order_type") val orderType: String,
    val shares: Int,
    @Json(name = "limit_price") val limitPrice: String?,
    @Json(name = "stop_price") val stopPrice: String?,
    val status: String,
    @Json(name = "created_at") val createdAt: String,
)

@JsonClass(generateAdapter = true)
data class PortfolioDto(
    val id: String,
    @Json(name = "user_id") val userId: String,
    val cash: String,
    @Json(name = "net_worth") val netWorth: String,
)

@JsonClass(generateAdapter = true)
data class PositionDto(
    val id: String,
    val ticker: String,
    val shares: Int,
    @Json(name = "avg_cost") val avgCost: String,
    @Json(name = "current_price") val currentPrice: String,
    @Json(name = "market_value") val marketValue: String,
    val pnl: String,
    @Json(name = "pnl_pct") val pnlPct: String,
)

@JsonClass(generateAdapter = true)
data class PortfolioResponseDto(
    val portfolio: PortfolioDto,
    val positions: List<PositionDto>,
    @Json(name = "net_worth") val netWorth: String,
    val invested: String,
)

@JsonClass(generateAdapter = true)
data class PortfolioHistoryDto(
    val id: Long,
    @Json(name = "net_worth") val netWorth: String,
    val cash: String,
    @Json(name = "recorded_at") val recordedAt: String,
)

@JsonClass(generateAdapter = true)
data class LeaderboardEntryDto(
    val id: Long,
    @Json(name = "user_id") val userId: String,
    @Json(name = "display_name") val displayName: String,
    @Json(name = "net_worth") val netWorth: String,
    @Json(name = "total_return") val totalReturn: String,
    val rank: Int,
    val period: String,
)

@JsonClass(generateAdapter = true)
data class AchievementDto(
    val id: String,
    val name: String,
    val description: String,
    val icon: String,
    val category: String,
)

@JsonClass(generateAdapter = true)
data class UserAchievementDto(
    val id: String,
    @Json(name = "achievement_id") val achievementId: String,
    @Json(name = "earned_at") val earnedAt: String,
)

@JsonClass(generateAdapter = true)
data class AlertDto(
    val id: String,
    val ticker: String,
    val condition: String,
    @Json(name = "target_price") val targetPrice: String,
    val triggered: Boolean,
    @Json(name = "created_at") val createdAt: String,
)
