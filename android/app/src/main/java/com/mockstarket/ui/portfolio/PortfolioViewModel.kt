package com.mockstarket.ui.portfolio

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.domain.model.Position
import com.mockstarket.domain.model.Trade
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.math.BigDecimal
import java.time.Instant
import javax.inject.Inject

data class PortfolioUiState(
    val netWorth: BigDecimal = BigDecimal.ZERO,
    val cash: BigDecimal = BigDecimal.ZERO,
    val invested: BigDecimal = BigDecimal.ZERO,
    val positions: List<Position> = emptyList(),
    val trades: List<Trade> = emptyList(),
    val isLoading: Boolean = true,
    val error: String? = null,
)

@HiltViewModel
class PortfolioViewModel @Inject constructor(
    private val api: MockStarketApi,
) : ViewModel() {

    private val _uiState = MutableStateFlow(PortfolioUiState())
    val uiState: StateFlow<PortfolioUiState> = _uiState.asStateFlow()

    init { load() }

    fun load() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val resp = api.getPortfolio()
                val trades = api.getTradeHistory(limit = 20, offset = 0)

                val positions = resp.positions.map {
                    Position(
                        id = it.id,
                        ticker = it.ticker,
                        shares = it.shares,
                        avgCost = BigDecimal(it.avg_cost),
                        currentPrice = BigDecimal(it.current_price),
                        marketValue = BigDecimal(it.market_value),
                        pnl = BigDecimal(it.pnl),
                        pnlPct = BigDecimal(it.pnl_pct),
                    )
                }

                val tradeList = trades.map {
                    Trade(
                        id = it.id,
                        ticker = it.ticker,
                        side = it.side,
                        shares = it.shares,
                        price = BigDecimal(it.price),
                        total = BigDecimal(it.total),
                        createdAt = Instant.parse(it.created_at),
                    )
                }

                _uiState.value = PortfolioUiState(
                    netWorth = BigDecimal(resp.net_worth),
                    cash = BigDecimal(resp.portfolio.cash),
                    invested = BigDecimal(resp.invested),
                    positions = positions,
                    trades = tradeList,
                    isLoading = false,
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(isLoading = false, error = e.message)
            }
        }
    }
}
