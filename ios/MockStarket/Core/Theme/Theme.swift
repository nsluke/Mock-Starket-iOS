import SwiftUI

enum Theme {
    // Brand colors
    static let accent = Color(hex: 0x50E3C2)       // Aquamarine green
    static let positive = Color(hex: 0x4ADE80)     // Green for gains
    static let negative = Color(hex: 0xF87171)     // Red for losses
    static let neutral = Color(hex: 0x9CA3AF)      // Gray for unchanged

    // Surface colors
    static let background = Color(hex: 0x0D1117)
    static let surface = Color(hex: 0x161B22)
    static let surfaceElevated = Color(hex: 0x21262D)
    static let border = Color(hex: 0x30363D)

    // Text
    static let textPrimary = Color(hex: 0xE6EDF3)
    static let textSecondary = Color(hex: 0x8B949E)
    static let textTertiary = Color(hex: 0x6E7681)

    // Helpers
    static func priceColor(for change: Decimal) -> Color {
        if change > 0 { return positive }
        if change < 0 { return negative }
        return neutral
    }

    static func priceArrow(for change: Decimal) -> String {
        if change > 0 { return "arrow.up.right" }
        if change < 0 { return "arrow.down.right" }
        return "minus"
    }
}

extension Color {
    init(hex: UInt, alpha: Double = 1.0) {
        self.init(
            .sRGB,
            red: Double((hex >> 16) & 0xFF) / 255.0,
            green: Double((hex >> 8) & 0xFF) / 255.0,
            blue: Double(hex & 0xFF) / 255.0,
            opacity: alpha
        )
    }
}
