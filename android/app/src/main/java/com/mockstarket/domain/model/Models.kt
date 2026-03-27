package com.mockstarket.domain.model

import java.math.BigDecimal
import java.time.Instant

data class Stock(
    val ticker: String,
    val name: String,
    val sector: String,
    val basePrice: BigDecimal,
    val currentPrice: BigDecimal,
    val dayOpen: BigDecimal,
    val dayHigh: BigDecimal,
    val dayLow: BigDecimal,
    val prevClose: BigDecimal,
    val volume: Long,
    val volatility: BigDecimal,
    val description: String? = null,
) {
    val change: BigDecimal get() = currentPrice - dayOpen
    val changePct: BigDecimal get() {
        if (dayOpen == BigDecimal.ZERO) return BigDecimal.ZERO
        return change.divide(dayOpen, 4, java.math.RoundingMode.HALF_UP) * BigDecimal(100)
    }
    val isUp: Boolean get() = change >= BigDecimal.ZERO
}

data class Portfolio(
    val id: String,
    val userId: String,
    val cash: BigDecimal,
    val netWorth: BigDecimal,
)

data class Position(
    val id: String,
    val ticker: String,
    val shares: Int,
    val avgCost: BigDecimal,
    val currentPrice: BigDecimal,
    val marketValue: BigDecimal,
    val pnl: BigDecimal,
    val pnlPct: BigDecimal,
) {
    val isProfit: Boolean get() = pnl >= BigDecimal.ZERO
}

data class Trade(
    val id: String,
    val ticker: String,
    val side: String,
    val shares: Int,
    val price: BigDecimal,
    val total: BigDecimal,
    val createdAt: Instant,
) {
    val isBuy: Boolean get() = side == "buy"
}

data class LeaderboardEntry(
    val id: Long,
    val userId: String,
    val displayName: String,
    val netWorth: BigDecimal,
    val totalReturn: BigDecimal,
    val rank: Int,
    val period: String,
)

data class Achievement(
    val id: String,
    val name: String,
    val description: String,
    val icon: String,
    val category: String,
)

data class User(
    val id: String,
    val displayName: String,
    val isGuest: Boolean,
    val loginStreak: Int,
    val longestStreak: Int,
)

data class PriceUpdate(
    val ticker: String,
    val price: BigDecimal,
    val change: BigDecimal,
    val changePct: BigDecimal,
    val volume: Long,
    val high: BigDecimal,
    val low: BigDecimal,
)
