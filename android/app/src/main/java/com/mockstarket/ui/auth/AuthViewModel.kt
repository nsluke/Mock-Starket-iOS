package com.mockstarket.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.mockstarket.data.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val isLoading: Boolean = false,
    val isSignedIn: Boolean = false,
    val error: String? = null,
)

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository,
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    init {
        checkAuth()
    }

    private fun checkAuth() {
        viewModelScope.launch {
            val token = authRepository.getToken()
            if (token != null) {
                try {
                    authRepository.getCurrentUser()
                    _uiState.value = AuthUiState(isSignedIn = true)
                } catch (_: Exception) {
                    // Token invalid, stay on auth screen
                }
            }
        }
    }

    fun signInAsGuest() {
        viewModelScope.launch {
            _uiState.value = AuthUiState(isLoading = true)
            try {
                authRepository.signInAsGuest()
                _uiState.value = AuthUiState(isSignedIn = true)
            } catch (e: Exception) {
                _uiState.value = AuthUiState(error = e.message ?: "Sign in failed")
            }
        }
    }
}
