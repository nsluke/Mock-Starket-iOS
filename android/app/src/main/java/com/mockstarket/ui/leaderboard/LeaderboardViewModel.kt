package com.mockstarket.ui.leaderboard

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.domain.model.LeaderboardEntry
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.math.BigDecimal
import javax.inject.Inject

data class LeaderboardUiState(
    val entries: List<LeaderboardEntry> = emptyList(),
    val period: String = "alltime",
    val isLoading: Boolean = true,
)

@HiltViewModel
class LeaderboardViewModel @Inject constructor(
    private val api: MockStarketApi,
) : ViewModel() {

    private val _uiState = MutableStateFlow(LeaderboardUiState())
    val uiState: StateFlow<LeaderboardUiState> = _uiState.asStateFlow()

    init { load() }

    fun setPeriod(period: String) {
        _uiState.value = _uiState.value.copy(period = period)
        load()
    }

    fun load() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val dtos = api.getLeaderboard(_uiState.value.period)
                val entries = dtos.map {
                    LeaderboardEntry(
                        id = it.id,
                        userId = it.user_id,
                        displayName = it.display_name,
                        netWorth = BigDecimal(it.net_worth),
                        totalReturn = BigDecimal(it.total_return),
                        rank = it.rank,
                        period = it.period,
                    )
                }
                _uiState.value = _uiState.value.copy(entries = entries, isLoading = false)
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoading = false)
            }
        }
    }
}
