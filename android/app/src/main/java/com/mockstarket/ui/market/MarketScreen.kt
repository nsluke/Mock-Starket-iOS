package com.mockstarket.ui.market

import androidx.compose.foundation.clickable
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
import com.mockstarket.domain.model.Stock
import com.mockstarket.ui.theme.*
import java.text.NumberFormat
import java.util.Locale

@Composable
fun MarketScreen(
    navController: NavController,
    viewModel: MarketViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()
    val displayed = viewModel.filteredStocks()
    val fmt = remember { NumberFormat.getCurrencyInstance(Locale.US) }

    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        item {
            Text("Market", fontSize = 24.sp, fontWeight = FontWeight.Bold)
        }

        // Market summary
        if (!uiState.isLoading) {
            item {
                Card(
                    colors = CardDefaults.cardColors(containerColor = Surface),
                    shape = RoundedCornerShape(16.dp),
                ) {
                    Column(Modifier.padding(20.dp)) {
                        Text("Market Index", color = TextSecondary, fontSize = 13.sp)
                        Spacer(Modifier.height(4.dp))
                        Row(verticalAlignment = Alignment.Bottom) {
                            Text(
                                fmt.format(uiState.indexValue),
                                fontSize = 28.sp,
                                fontWeight = FontWeight.Bold,
                            )
                            Spacer(Modifier.width(12.dp))
                            Text(
                                "${if (uiState.indexChangePct >= java.math.BigDecimal.ZERO) "+" else ""}${uiState.indexChangePct}%",
                                color = priceColor(uiState.indexChangePct >= java.math.BigDecimal.ZERO),
                                fontWeight = FontWeight.SemiBold,
                            )
                        }
                        Spacer(Modifier.height(8.dp))
                        Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                            Text("${uiState.gainers} gainers", color = Positive, fontSize = 13.sp)
                            Text("${uiState.losers} losers", color = Negative, fontSize = 13.sp)
                            Text("${uiState.stocks.size} stocks", color = TextSecondary, fontSize = 13.sp)
                        }
                    }
                }
            }
        }

        // Search
        item {
            OutlinedTextField(
                value = uiState.searchQuery,
                onValueChange = { viewModel.setSearchQuery(it) },
                placeholder = { Text("Search stocks...", color = TextTertiary) },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                shape = RoundedCornerShape(12.dp),
            )
        }

        if (uiState.isLoading) {
            item {
                Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator(color = Accent)
                }
            }
        } else {
            items(displayed, key = { it.ticker }) { stock ->
                StockRow(stock = stock, fmt = fmt) {
                    navController.navigate("stock/${stock.ticker}")
                }
            }
        }
    }
}

@Composable
private fun StockRow(stock: Stock, fmt: NumberFormat, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().clickable(onClick = onClick),
        colors = CardDefaults.cardColors(containerColor = Surface),
        shape = RoundedCornerShape(12.dp),
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Surface(
                        color = Accent.copy(alpha = 0.1f),
                        shape = RoundedCornerShape(6.dp),
                    ) {
                        Text(
                            stock.ticker,
                            modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                            fontFamily = FontFamily.Monospace,
                            fontWeight = FontWeight.Bold,
                            fontSize = 13.sp,
                            color = Accent,
                        )
                    }
                    Spacer(Modifier.width(8.dp))
                    Text(stock.sector, color = TextTertiary, fontSize = 12.sp)
                }
                Spacer(Modifier.height(4.dp))
                Text(stock.name, fontWeight = FontWeight.Medium, fontSize = 14.sp)
            }

            Column(horizontalAlignment = Alignment.End) {
                Text(fmt.format(stock.currentPrice), fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                Spacer(Modifier.height(2.dp))
                Surface(
                    color = priceColor(stock.isUp).copy(alpha = 0.1f),
                    shape = RoundedCornerShape(6.dp),
                ) {
                    Text(
                        "${if (stock.isUp) "+" else ""}${stock.changePct.setScale(2)}%",
                        modifier = Modifier.padding(horizontal = 8.dp, vertical = 3.dp),
                        color = priceColor(stock.isUp),
                        fontWeight = FontWeight.SemiBold,
                        fontSize = 12.sp,
                    )
                }
            }
        }
    }
}
