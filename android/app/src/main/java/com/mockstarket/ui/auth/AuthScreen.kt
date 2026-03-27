package com.mockstarket.ui.auth

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ShowChart
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.mockstarket.ui.theme.Accent
import com.mockstarket.ui.theme.TextSecondary
import com.mockstarket.ui.theme.TextTertiary

@Composable
fun AuthScreen(
    onSignedIn: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(uiState.isSignedIn) {
        if (uiState.isSignedIn) onSignedIn()
    }

    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center,
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier.padding(32.dp),
        ) {
            Spacer(Modifier.weight(1f))

            Icon(
                Icons.Default.ShowChart,
                contentDescription = null,
                modifier = Modifier.size(80.dp),
                tint = Accent,
            )
            Spacer(Modifier.height(24.dp))
            Text(
                "Mock Starket",
                fontSize = 32.sp,
                fontWeight = FontWeight.Bold,
            )
            Spacer(Modifier.height(8.dp))
            Text(
                "Learn to trade. Risk nothing.",
                color = TextSecondary,
                textAlign = TextAlign.Center,
            )
            Spacer(Modifier.height(8.dp))
            Text(
                "$100,000 starting cash",
                color = TextTertiary,
                fontSize = 13.sp,
            )

            Spacer(Modifier.weight(1f))

            if (uiState.isLoading) {
                CircularProgressIndicator(color = Accent)
            } else {
                Button(
                    onClick = { viewModel.signInAsGuest() },
                    modifier = Modifier.fillMaxWidth().height(56.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Accent),
                ) {
                    Text("Continue as Guest", fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.onPrimary)
                }

                Spacer(Modifier.height(12.dp))

                OutlinedButton(
                    onClick = { viewModel.signInAsGuest() },
                    modifier = Modifier.fillMaxWidth().height(56.dp),
                ) {
                    Text("Sign in with Email")
                }
            }

            if (uiState.error != null) {
                Spacer(Modifier.height(12.dp))
                Text(uiState.error!!, color = MaterialTheme.colorScheme.error, fontSize = 13.sp)
            }

            Spacer(Modifier.height(40.dp))
        }
    }
}
