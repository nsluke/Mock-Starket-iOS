package com.mockstarket.ui

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.runtime.*
import com.mockstarket.ui.auth.AuthScreen
import com.mockstarket.ui.navigation.MockStarketNavHost
import com.mockstarket.ui.theme.MockStarketTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            MockStarketTheme {
                var isAuthenticated by remember { mutableStateOf(false) }

                if (isAuthenticated) {
                    MockStarketNavHost(
                        onSignedOut = { isAuthenticated = false }
                    )
                } else {
                    AuthScreen(
                        onSignedIn = { isAuthenticated = true }
                    )
                }
            }
        }
    }
}
