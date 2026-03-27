package com.mockstarket.di

import com.mockstarket.data.local.datastore.TokenManager
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthInterceptor @Inject constructor(
    private val tokenManager: TokenManager,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        val token = runBlocking { tokenManager.getToken() }

        return if (token != null) {
            val authed = request.newBuilder()
                .addHeader("Authorization", "Bearer $token")
                .build()
            chain.proceed(authed)
        } else {
            chain.proceed(request)
        }
    }
}
