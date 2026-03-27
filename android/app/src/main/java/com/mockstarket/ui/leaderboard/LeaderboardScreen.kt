package com.mockstarket.ui.leaderboard

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.mockstarket.ui.theme.*
import java.text.NumberFormat
import java.util.Locale

private val periods = listOf("alltime" to "All Time", "weekly" to "Weekly", "daily" to "Daily")

@Composable
fun LeaderboardScreen(
    viewModel: LeaderboardViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()
    val fmt = remember { NumberFormat.getCurrencyInstance(Locale.US) }

    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        item { Text("Leaderboard", fontSize = 24.sp, fontWeight = FontWeight.Bold) }

        item {
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                periods.forEach { (value, label) ->
                    FilterChip(selected = uiState.period == value, onClick = { viewModel.setPeriod(value) }, label = { Text(label) })
                }
            }
        }

        if (uiState.isLoading) {
            item { Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) { CircularProgressIndicator(color = Accent) } }
        } else if (uiState.entries.isEmpty()) {
            item {
                Card(colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
                    Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) { Text("No rankings yet.", color = TextSecondary) }
                }
            }
        } else {
            items(uiState.entries, key = { it.id }) { entry ->
                val badgeColor = when (entry.rank) { 1 -> Color(0xFFFFD700); 2 -> Color(0xFFC0C0C0); 3 -> Color(0xFFCD7F32); else -> TextTertiary }
                Card(colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
                    Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                        Surface(Modifier.size(36.dp), shape = CircleShape, color = badgeColor.copy(alpha = 0.2f)) {
                            Box(contentAlignment = Alignment.Center) { Text("${entry.rank}", fontWeight = FontWeight.Bold, color = badgeColor) }
                        }
                        Spacer(Modifier.width(12.dp))
                        Text(entry.displayName, fontWeight = FontWeight.Medium, modifier = Modifier.weight(1f))
                        Column(horizontalAlignment = Alignment.End) {
                            Text(fmt.format(entry.netWorth), fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                            Text("${if (entry.totalReturn >= java.math.BigDecimal.ZERO) "+" else ""}${entry.totalReturn.setScale(2)}%", color = priceColor(entry.totalReturn >= java.math.BigDecimal.ZERO), fontSize = 13.sp)
                        }
                    }
                }
            }
        }
    }
}
