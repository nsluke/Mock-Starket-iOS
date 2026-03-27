package com.mockstarket.ui.stockdetail

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.data.remote.api.TradeRequest
import com.mockstarket.domain.model.Stock
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.math.BigDecimal
import java.text.NumberFormat
import java.util.Locale
import javax.inject.Inject

data class StockDetailUiState(
    val stock: Stock? = null,
    val isLoading: Boolean = true,
    val side: String = "buy",
    val shares: String = "",
    val isTrading: Boolean = false,
    val tradeMessage: String? = null,
    val tradeSuccess: Boolean = false,
)

@HiltViewModel
class StockDetailViewModel @Inject constructor(
    private val api: MockStarketApi,
) : ViewModel() {

    private val _uiState = MutableStateFlow(StockDetailUiState())
    val uiState: StateFlow<StockDetailUiState> = _uiState.asStateFlow()

    private var ticker = ""

    fun load(ticker: String) {
        this.ticker = ticker
        viewModelScope.launch {
            _uiState.value = StockDetailUiState(isLoading = true)
            try {
                val dto = api.getStock(ticker)
                val stock = Stock(
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
                _uiState.value = StockDetailUiState(stock = stock, isLoading = false)
            } catch (e: Exception) {
                _uiState.value = StockDetailUiState(isLoading = false, tradeMessage = e.message)
            }
        }
    }

    fun setSide(side: String) {
        _uiState.value = _uiState.value.copy(side = side, tradeMessage = null)
    }

    fun setShares(shares: String) {
        _uiState.value = _uiState.value.copy(shares = shares, tradeMessage = null)
    }

    fun executeTrade() {
        val qty = _uiState.value.shares.toIntOrNull() ?: return
        if (qty <= 0) return

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isTrading = true, tradeMessage = null)
            try {
                api.executeTrade(TradeRequest(ticker, _uiState.value.side, qty))
                val fmt = NumberFormat.getCurrencyInstance(Locale.US)
                val price = _uiState.value.stock?.currentPrice ?: BigDecimal.ZERO
                val verb = if (_uiState.value.side == "buy") "Bought" else "Sold"
                _uiState.value = _uiState.value.copy(
                    isTrading = false,
                    shares = "",
                    tradeMessage = "$verb $qty shares of $ticker at ${fmt.format(price)}",
                    tradeSuccess = true,
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isTrading = false,
                    tradeMessage = e.message ?: "Trade failed",
                    tradeSuccess = false,
                )
            }
        }
    }
}
