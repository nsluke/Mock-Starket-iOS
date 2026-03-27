import Foundation

extension Decimal {
    var currencyFormatted: String {
        let formatter = NumberFormatter()
        formatter.numberStyle = .currency
        formatter.currencyCode = "USD"
        formatter.maximumFractionDigits = 2
        formatter.minimumFractionDigits = 2
        return formatter.string(from: self as NSDecimalNumber) ?? "$0.00"
    }

    var percentFormatted: String {
        let formatter = NumberFormatter()
        formatter.numberStyle = .decimal
        formatter.maximumFractionDigits = 2
        formatter.minimumFractionDigits = 2
        formatter.positivePrefix = "+"
        return "\(formatter.string(from: self as NSDecimalNumber) ?? "0.00")%"
    }

    var compactFormatted: String {
        let number = NSDecimalNumber(decimal: self).doubleValue
        switch abs(number) {
        case 1_000_000_000...:
            return String(format: "$%.1fB", number / 1_000_000_000)
        case 1_000_000...:
            return String(format: "$%.1fM", number / 1_000_000)
        case 1_000...:
            return String(format: "$%.1fK", number / 1_000)
        default:
            return currencyFormatted
        }
    }

    var volumeFormatted: String {
        let number = NSDecimalNumber(decimal: self).int64Value
        switch abs(number) {
        case 1_000_000...:
            return String(format: "%.1fM", Double(number) / 1_000_000)
        case 1_000...:
            return String(format: "%.1fK", Double(number) / 1_000)
        default:
            return "\(number)"
        }
    }
}
