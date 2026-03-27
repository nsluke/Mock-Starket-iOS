package com.mockstarket.data.repository

import com.mockstarket.data.local.datastore.TokenManager
import com.mockstarket.data.remote.api.MockStarketApi
import com.mockstarket.data.remote.api.RegisterRequest
import com.mockstarket.domain.model.User
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val api: MockStarketApi,
    private val tokenManager: TokenManager,
) {
    suspend fun getToken(): String? = tokenManager.getToken()

    suspend fun signInAsGuest(): User {
        val uid = "guest-android-${System.currentTimeMillis()}"
        tokenManager.saveToken(uid)

        val dto = api.register(RegisterRequest(display_name = "Guest Trader", is_guest = true))
        return User(
            id = dto.id,
            displayName = dto.display_name,
            isGuest = dto.is_guest,
            loginStreak = dto.login_streak,
            longestStreak = dto.longest_streak,
        )
    }

    suspend fun getCurrentUser(): User {
        val dto = api.getMe()
        return User(
            id = dto.id,
            displayName = dto.display_name,
            isGuest = dto.is_guest,
            loginStreak = dto.login_streak,
            longestStreak = dto.longest_streak,
        )
    }

    suspend fun signOut() {
        tokenManager.clearToken()
    }
}
