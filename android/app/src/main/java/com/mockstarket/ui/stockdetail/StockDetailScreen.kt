package com.mockstarket.ui.stockdetail

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavController
import com.mockstarket.ui.theme.*
import java.text.NumberFormat
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StockDetailScreen(
    ticker: String,
    navController: NavController,
    viewModel: StockDetailViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()
    val fmt = remember { NumberFormat.getCurrencyInstance(Locale.US) }

    LaunchedEffect(ticker) { viewModel.load(ticker) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(ticker, fontFamily = FontFamily.Monospace, fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = { navController.popBackStack() }) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                },
            )
        }
    ) { padding ->
        if (uiState.isLoading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                CircularProgressIndicator(color = Accent)
            }
            return@Scaffold
        }

        val stock = uiState.stock ?: return@Scaffold

        Column(
            modifier = Modifier
                .padding(padding)
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp),
        ) {
            // Price header
            Card(
                colors = CardDefaults.cardColors(containerColor = Surface),
                shape = RoundedCornerShape(16.dp),
            ) {
                Column(Modifier.padding(20.dp)) {
                    Text(stock.name, fontWeight = FontWeight.Medium, color = TextSecondary)
                    Spacer(Modifier.height(4.dp))
                    Row(verticalAlignment = Alignment.Bottom) {
                        Text(fmt.format(stock.currentPrice), fontSize = 32.sp, fontWeight = FontWeight.Bold)
                        Spacer(Modifier.width(12.dp))
                        Text(
                            "${if (stock.isUp) "+" else ""}${fmt.format(stock.change)} (${stock.changePct.setScale(2)}%)",
                            color = priceColor(stock.isUp),
                            fontWeight = FontWeight.SemiBold,
                        )
                    }
                    Spacer(Modifier.height(16.dp))

                    Row(horizontalArrangement = Arrangement.spacedBy(24.dp)) {
                        StatItem("Open", fmt.format(stock.dayOpen))
                        StatItem("High", fmt.format(stock.dayHigh), Positive)
                        StatItem("Low", fmt.format(stock.dayLow), Negative)
                        StatItem("Vol", formatVolume(stock.volume))
                    }
                }
            }

            // Trade panel
            Card(
                colors = CardDefaults.cardColors(containerColor = Surface),
                shape = RoundedCornerShape(16.dp),
            ) {
                Column(Modifier.padding(20.dp)) {
                    Text("Trade $ticker", fontWeight = FontWeight.SemiBold)
                    Spacer(Modifier.height(12.dp))

                    // Buy/Sell toggle
                    Row(Modifier.fillMaxWidth()) {
                        FilterChip(
                            selected = uiState.side == "buy",
                            onClick = { viewModel.setSide("buy") },
                            label = { Text("Buy") },
                            modifier = Modifier.weight(1f).padding(end = 4.dp),
                            colors = FilterChipDefaults.filterChipColors(
                                selectedContainerColor = Positive.copy(alpha = 0.2f),
                                selectedLabelColor = Positive,
                            ),
                        )
                        FilterChip(
                            selected = uiState.side == "sell",
                            onClick = { viewModel.setSide("sell") },
                            label = { Text("Sell") },
                            modifier = Modifier.weight(1f).padding(start = 4.dp),
                            colors = FilterChipDefaults.filterChipColors(
                                selectedContainerColor = Negative.copy(alpha = 0.2f),
                                selectedLabelColor = Negative,
                            ),
                        )
                    }
                    Spacer(Modifier.height(12.dp))

                    OutlinedTextField(
                        value = uiState.shares,
                        onValueChange = { viewModel.setShares(it) },
                        label = { Text("Shares") },
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        shape = RoundedCornerShape(12.dp),
                    )

                    val qty = uiState.shares.toIntOrNull() ?: 0
                    if (qty > 0) {
                        Spacer(Modifier.height(8.dp))
                        Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Estimated Total", color = TextSecondary, fontSize = 14.sp)
                            Text(
                                fmt.format(stock.currentPrice * java.math.BigDecimal(qty)),
                                fontWeight = FontWeight.SemiBold,
                            )
                        }
                    }

                    Spacer(Modifier.height(16.dp))

                    Button(
                        onClick = { viewModel.executeTrade() },
                        modifier = Modifier.fillMaxWidth().height(48.dp),
                        enabled = !uiState.isTrading && qty > 0,
                        colors = ButtonDefaults.buttonColors(
                            containerColor = if (uiState.side == "buy") Positive else Negative,
                        ),
                        shape = RoundedCornerShape(12.dp),
                    ) {
                        Text(
                            if (uiState.isTrading) "Processing..."
                            else "${if (uiState.side == "buy") "Buy" else "Sell"} $ticker",
                            fontWeight = FontWeight.Bold,
                        )
                    }

                    // Result message
                    uiState.tradeMessage?.let { msg ->
                        Spacer(Modifier.height(8.dp))
                        Surface(
                            color = if (uiState.tradeSuccess) Positive.copy(alpha = 0.1f) else Negative.copy(alpha = 0.1f),
                            shape = RoundedCornerShape(8.dp),
                        ) {
                            Text(
                                msg,
                                modifier = Modifier.padding(12.dp),
                                color = if (uiState.tradeSuccess) Positive else Negative,
                                fontSize = 13.sp,
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun StatItem(label: String, value: String, color: androidx.compose.ui.graphics.Color = TextPrimary) {
    Column {
        Text(label, color = TextTertiary, fontSize = 12.sp)
        Text(value, fontWeight = FontWeight.Medium, fontSize = 14.sp, color = color)
    }
}

private fun formatVolume(v: Long): String = when {
    v >= 1_000_000 -> "${v / 1_000_000}M"
    v >= 1_000 -> "${v / 1_000}K"
    else -> v.toString()
}
