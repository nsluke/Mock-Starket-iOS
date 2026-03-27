package com.mockstarket.ui.portfolio

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.mockstarket.ui.theme.*
import java.math.BigDecimal
import java.text.NumberFormat
import java.util.Locale

@Composable
fun PortfolioScreen(
    navController: NavController,
    viewModel: PortfolioViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()
    val fmt = remember { NumberFormat.getCurrencyInstance(Locale.US) }
    val startingCash = BigDecimal(100_000)

    if (uiState.isLoading) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator(color = Accent)
        }
        return
    }

    val totalPnl = uiState.netWorth - startingCash
    val totalPnlPct = if (startingCash > BigDecimal.ZERO) {
        totalPnl.divide(startingCash, 4, java.math.RoundingMode.HALF_UP) * BigDecimal(100)
    } else BigDecimal.ZERO

    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        item { Text("Portfolio", fontSize = 24.sp, fontWeight = FontWeight.Bold) }

        // Summary cards
        item {
            Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                SummaryCard("Net Worth", fmt.format(uiState.netWorth), "${totalPnlPct.setScale(2)}%", totalPnl >= BigDecimal.ZERO, Modifier.weight(1f))
                SummaryCard("Cash", fmt.format(uiState.cash), null, true, Modifier.weight(1f))
            }
        }
        item {
            Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                SummaryCard("Invested", fmt.format(uiState.invested), "${uiState.positions.size} positions", true, Modifier.weight(1f))
                SummaryCard("Total P&L", "${if (totalPnl >= BigDecimal.ZERO) "+" else ""}${fmt.format(totalPnl)}", null, totalPnl >= BigDecimal.ZERO, Modifier.weight(1f))
            }
        }

        item {
            Text("Holdings", fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(top = 8.dp))
        }

        if (uiState.positions.isEmpty()) {
            item {
                Card(colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
                    Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) {
                        Text("No holdings yet. Start trading!", color = TextSecondary)
                    }
                }
            }
        } else {
            items(uiState.positions, key = { it.id }) { pos ->
                Card(colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
                    Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                        Column(Modifier.weight(1f)) {
                            Surface(color = Accent.copy(alpha = 0.1f), shape = RoundedCornerShape(6.dp)) {
                                Text(pos.ticker, Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Bold, fontSize = 13.sp, color = Accent)
                            }
                            Spacer(Modifier.height(4.dp))
                            Text("${pos.shares} shares @ ${fmt.format(pos.avgCost)}", color = TextSecondary, fontSize = 13.sp)
                        }
                        Column(horizontalAlignment = Alignment.End) {
                            Text(fmt.format(pos.marketValue), fontWeight = FontWeight.SemiBold)
                            Text("${if (pos.isProfit) "+" else ""}${fmt.format(pos.pnl)} (${pos.pnlPct.setScale(2)}%)", color = priceColor(pos.isProfit), fontSize = 13.sp)
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun SummaryCard(label: String, value: String, subtitle: String?, isPositive: Boolean, modifier: Modifier) {
    Card(modifier = modifier, colors = CardDefaults.cardColors(containerColor = Surface), shape = RoundedCornerShape(12.dp)) {
        Column(Modifier.padding(16.dp)) {
            Text(label, color = TextTertiary, fontSize = 12.sp)
            Spacer(Modifier.height(4.dp))
            Text(value, fontWeight = FontWeight.Bold, fontSize = 18.sp)
            if (subtitle != null) {
                Spacer(Modifier.height(2.dp))
                Text(subtitle, color = TextSecondary, fontSize = 12.sp)
            }
        }
    }
}
