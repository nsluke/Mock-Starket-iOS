package com.mockstarket.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

val Accent = Color(0xFF50E3C2)
val Positive = Color(0xFF4ADE80)
val Negative = Color(0xFFF87171)
val Neutral = Color(0xFF9CA3AF)

val Background = Color(0xFF0D1117)
val Surface = Color(0xFF161B22)
val SurfaceElevated = Color(0xFF21262D)
val Border = Color(0xFF30363D)

val TextPrimary = Color(0xFFE6EDF3)
val TextSecondary = Color(0xFF8B949E)
val TextTertiary = Color(0xFF6E7681)

private val DarkColorScheme = darkColorScheme(
    primary = Accent,
    onPrimary = Color.Black,
    secondary = Accent,
    onSecondary = Color.Black,
    tertiary = Accent,
    background = Background,
    onBackground = TextPrimary,
    surface = Surface,
    onSurface = TextPrimary,
    surfaceVariant = SurfaceElevated,
    onSurfaceVariant = TextSecondary,
    outline = Border,
    error = Negative,
    onError = Color.White,
)

@Composable
fun MockStarketTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = DarkColorScheme,
        content = content,
    )
}

fun priceColor(isUp: Boolean): Color = if (isUp) Positive else Negative
