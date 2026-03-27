package com.mockstarket.ui.profile

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ExitToApp
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.mockstarket.ui.theme.*

@Composable
fun ProfileScreen(
    onSignedOut: () -> Unit = {},
    viewModel: ProfileViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(uiState.isSignedOut) {
        if (uiState.isSignedOut) onSignedOut()
    }

    if (uiState.isLoading) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) { CircularProgressIndicator(color = Accent) }
        return
    }

    val user = uiState.user ?: return

    LazyColumn(Modifier.fillMaxSize(), contentPadding = PaddingValues(16.dp), verticalArrangement = Arrangement.spacedBy(16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
        item {
            Column(Modifier.fillMaxWidth().padding(top = 16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.Person, null, Modifier.size(80.dp), tint = Accent)
                Spacer(Modifier.height(12.dp))
                Text(user.displayName, fontSize = 20.sp, fontWeight = FontWeight.Bold)
                if (user.loginStreak > 0) {
                    Spacer(Modifier.height(8.dp))
                    Surface(color = SurfaceElevated, shape = RoundedCornerShape(20.dp)) {
                        Text("🔥 ${user.loginStreak} day streak", Modifier.padding(horizontal = 16.dp, vertical = 8.dp), fontWeight = FontWeight.SemiBold, fontSize = 14.sp)
                    }
                }
            }
        }

        item {
            Card(colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
                Column(Modifier.padding(16.dp)) {
                    Text("Account", fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(bottom = 8.dp))
                    Row(Modifier.fillMaxWidth().padding(vertical = 4.dp), horizontalArrangement = Arrangement.SpaceBetween) { Text("Type", color = TextSecondary, fontSize = 14.sp); Text(if (user.isGuest) "Guest" else "Registered", fontWeight = FontWeight.Medium, fontSize = 14.sp) }
                    Row(Modifier.fillMaxWidth().padding(vertical = 4.dp), horizontalArrangement = Arrangement.SpaceBetween) { Text("Longest Streak", color = TextSecondary, fontSize = 14.sp); Text("${user.longestStreak} days", fontWeight = FontWeight.Medium, fontSize = 14.sp) }
                }
            }
        }

        if (uiState.achievements.isNotEmpty()) {
            item { Text("Achievements (${uiState.earnedIds.size}/${uiState.achievements.size})", fontWeight = FontWeight.SemiBold) }
            items(uiState.achievements, key = { it.id }) { a ->
                val earned = a.id in uiState.earnedIds
                Card(colors = CardDefaults.cardColors(containerColor = if (earned) Surface else Surface.copy(alpha = 0.5f)), shape = RoundedCornerShape(12.dp)) {
                    Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                        Text(if (earned) "✅" else "🔒", fontSize = 24.sp)
                        Spacer(Modifier.width(12.dp))
                        Column(Modifier.weight(1f)) {
                            Text(a.name, fontWeight = FontWeight.Medium, color = if (earned) TextPrimary else TextTertiary)
                            Text(a.description, fontSize = 13.sp, color = TextSecondary)
                        }
                    }
                }
            }
        }

        item {
            OutlinedButton(onClick = { viewModel.signOut() }, Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp), colors = ButtonDefaults.outlinedButtonColors(contentColor = Negative)) {
                Icon(Icons.AutoMirrored.Filled.ExitToApp, null); Spacer(Modifier.width(8.dp)); Text("Sign Out")
            }
        }
    }
}
