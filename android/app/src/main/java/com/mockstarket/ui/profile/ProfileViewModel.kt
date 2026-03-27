package com.mockstarket.ui.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.data.repository.AuthRepository
import com.mockstarket.domain.model.Achievement
import com.mockstarket.domain.model.User
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ProfileUiState(
    val user: User? = null,
    val achievements: List<Achievement> = emptyList(),
    val earnedIds: Set<String> = emptySet(),
    val isLoading: Boolean = true,
    val isSignedOut: Boolean = false,
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val api: MockStarketApi,
    private val authRepository: AuthRepository,
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    init { load() }

    fun load() {
        viewModelScope.launch {
            try {
                val user = authRepository.getCurrentUser()
                val achievements = api.getAchievements()
                val earned = api.getMyAchievements()

                _uiState.value = ProfileUiState(
                    user = user,
                    achievements = achievements.map {
                        Achievement(it.id, it.name, it.description, it.icon, it.category)
                    },
                    earnedIds = earned.map { it.achievement_id }.toSet(),
                    isLoading = false,
                )
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoading = false)
            }
        }
    }

    fun signOut() {
        viewModelScope.launch {
            authRepository.signOut()
            _uiState.value = _uiState.value.copy(isSignedOut = true)
        }
    }
}
