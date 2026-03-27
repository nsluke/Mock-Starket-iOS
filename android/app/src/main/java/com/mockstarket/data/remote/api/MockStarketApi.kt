package com.mockstarket.data.remote.api

import com.mockstarket.data.remote.dto.*
import retrofit2.http.*

interface MockStarketApi {

    // Auth
    @POST("api/v1/auth/register")
    suspend fun register(@Body body: RegisterRequest): UserDto

    @POST("api/v1/auth/guest")
    suspend fun createGuest(): UserDto

    @GET("api/v1/auth/me")
    suspend fun getMe(): UserDto

    @PUT("api/v1/auth/me")
    suspend fun updateMe(@Body body: UpdateUserRequest): Map<String, String>

    // Stocks
    @GET("api/v1/stocks")
    suspend fun getStocks(): List<StockDto>

    @GET("api/v1/stocks/{ticker}")
    suspend fun getStock(@Path("ticker") ticker: String): StockDto

    @GET("api/v1/stocks/{ticker}/history")
    suspend fun getStockHistory(
        @Path("ticker") ticker: String,
        @Query("interval") interval: String = "1m"
    ): List<PricePointDto>

    @GET("api/v1/stocks/market-summary")
    suspend fun getMarketSummary(): MarketSummaryDto

    // Trading
    @POST("api/v1/trades")
    suspend fun executeTrade(@Body body: TradeRequest): TradeDto

    @GET("api/v1/trades")
    suspend fun getTradeHistory(
        @Query("limit") limit: Int = 50,
        @Query("offset") offset: Int = 0
    ): List<TradeDto>

    // Orders
    @POST("api/v1/orders")
    suspend fun createOrder(@Body body: OrderRequest): OrderDto

    @GET("api/v1/orders")
    suspend fun getOrders(): List<OrderDto>

    @DELETE("api/v1/orders/{id}")
    suspend fun cancelOrder(@Path("id") id: String): Map<String, String>

    // Portfolio
    @GET("api/v1/portfolio")
    suspend fun getPortfolio(): PortfolioResponseDto

    @GET("api/v1/portfolio/history")
    suspend fun getPortfolioHistory(@Query("limit") limit: Int = 100): List<PortfolioHistoryDto>

    // Leaderboard
    @GET("api/v1/leaderboard")
    suspend fun getLeaderboard(@Query("period") period: String = "alltime"): List<LeaderboardEntryDto>

    // Alerts
    @POST("api/v1/alerts")
    suspend fun createAlert(@Body body: AlertRequest): AlertDto

    @GET("api/v1/alerts")
    suspend fun getAlerts(): List<AlertDto>

    @DELETE("api/v1/alerts/{id}")
    suspend fun deleteAlert(@Path("id") id: String): Map<String, String>

    // Achievements
    @GET("api/v1/achievements")
    suspend fun getAchievements(): List<AchievementDto>

    @GET("api/v1/achievements/me")
    suspend fun getMyAchievements(): List<UserAchievementDto>

    // Watchlist
    @GET("api/v1/watchlist")
    suspend fun getWatchlist(): List<String>

    @POST("api/v1/watchlist")
    suspend fun addToWatchlist(@Body body: WatchlistRequest): Map<String, String>

    @DELETE("api/v1/watchlist/{ticker}")
    suspend fun removeFromWatchlist(@Path("ticker") ticker: String): Map<String, String>
}

// Request bodies
data class RegisterRequest(val display_name: String, val is_guest: Boolean)
data class UpdateUserRequest(val display_name: String, val avatar_url: String?)
data class TradeRequest(val ticker: String, val side: String, val shares: Int)
data class OrderRequest(val ticker: String, val side: String, val order_type: String, val shares: Int, val limit_price: String? = null, val stop_price: String? = null)
data class AlertRequest(val ticker: String, val condition: String, val target_price: String)
data class WatchlistRequest(val ticker: String)
