package com.mockstarket.ui.market

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.domain.model.Stock
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.math.BigDecimal
import javax.inject.Inject

data class MarketUiState(
    val stocks: List<Stock> = emptyList(),
    val isLoading: Boolean = true,
    val searchQuery: String = "",
    val error: String? = null,
    val indexValue: BigDecimal = BigDecimal.ZERO,
    val indexChangePct: BigDecimal = BigDecimal.ZERO,
    val gainers: Int = 0,
    val losers: Int = 0,
)

@HiltViewModel
class MarketViewModel @Inject constructor(
    private val api: MockStarketApi,
) : ViewModel() {

    private val _uiState = MutableStateFlow(MarketUiState())
    val uiState: StateFlow<MarketUiState> = _uiState.asStateFlow()

    init {
        loadData()
    }

    fun loadData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val stockDtos = api.getStocks()
                val summary = api.getMarketSummary()

                val stocks = stockDtos.map { dto ->
                    Stock(
                        ticker = dto.ticker,
                        name = dto.name,
                        sector = dto.sector,
                        basePrice = BigDecimal(dto.base_price),
                        currentPrice = BigDecimal(dto.current_price),
                        dayOpen = BigDecimal(dto.day_open),
                        dayHigh = BigDecimal(dto.day_high),
                        dayLow = BigDecimal(dto.day_low),
                        prevClose = BigDecimal(dto.prev_close),
                        volume = dto.volume,
                        volatility = BigDecimal(dto.volatility),
                        description = dto.description,
                    )
                }

                _uiState.value = MarketUiState(
                    stocks = stocks,
                    isLoading = false,
                    indexValue = BigDecimal(summary.index_value),
                    indexChangePct = BigDecimal(summary.index_change_pct),
                    gainers = summary.gainers,
                    losers = summary.losers,
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(isLoading = false, error = e.message)
            }
        }
    }

    fun setSearchQuery(query: String) {
        _uiState.value = _uiState.value.copy(searchQuery = query)
    }

    fun filteredStocks(): List<Stock> {
        val state = _uiState.value
        if (state.searchQuery.isBlank()) return state.stocks
        val q = state.searchQuery.lowercase()
        return state.stocks.filter {
            it.ticker.lowercase().contains(q) || it.name.lowercase().contains(q)
        }
    }
}
