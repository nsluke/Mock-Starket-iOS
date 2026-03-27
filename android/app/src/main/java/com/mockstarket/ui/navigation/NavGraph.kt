package com.mockstarket.ui.navigation

import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.mockstarket.ui.market.MarketScreen
import com.mockstarket.ui.portfolio.PortfolioScreen
import com.mockstarket.ui.leaderboard.LeaderboardScreen
import com.mockstarket.ui.profile.ProfileScreen
import com.mockstarket.ui.stockdetail.StockDetailScreen

sealed class Screen(val route: String, val label: String, val icon: ImageVector) {
    object Market : Screen("market", "Market", Icons.Default.ShowChart)
    object Portfolio : Screen("portfolio", "Portfolio", Icons.Default.BusinessCenter)
    object Leaderboard : Screen("leaderboard", "Leaderboard", Icons.Default.EmojiEvents)
    object Profile : Screen("profile", "Profile", Icons.Default.Person)
}

val bottomNavScreens = listOf(Screen.Market, Screen.Portfolio, Screen.Leaderboard, Screen.Profile)

@Composable
fun MockStarketNavHost(onSignedOut: () -> Unit = {}) {
    val navController = rememberNavController()

    Scaffold(
        bottomBar = {
            NavigationBar {
                val navBackStackEntry by navController.currentBackStackEntryAsState()
                val currentDestination = navBackStackEntry?.destination

                bottomNavScreens.forEach { screen ->
                    NavigationBarItem(
                        icon = { Icon(screen.icon, contentDescription = screen.label) },
                        label = { Text(screen.label) },
                        selected = currentDestination?.hierarchy?.any { it.route == screen.route } == true,
                        onClick = {
                            navController.navigate(screen.route) {
                                popUpTo(navController.graph.findStartDestination().id) {
                                    saveState = true
                                }
                                launchSingleTop = true
                                restoreState = true
                            }
                        }
                    )
                }
            }
        }
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = Screen.Market.route,
            modifier = Modifier.padding(innerPadding)
        ) {
            composable(Screen.Market.route) { MarketScreen(navController) }
            composable(Screen.Portfolio.route) { PortfolioScreen(navController) }
            composable(Screen.Leaderboard.route) { LeaderboardScreen() }
            composable(Screen.Profile.route) { ProfileScreen(onSignedOut = onSignedOut) }

            composable(
                route = "stock/{ticker}",
                arguments = listOf(navArgument("ticker") { type = NavType.StringType })
            ) { backStackEntry ->
                val ticker = backStackEntry.arguments?.getString("ticker") ?: return@composable
                StockDetailScreen(ticker = ticker, navController = navController)
            }
        }
    }
}
