package polygon

import "strconv"

// SectorFromSIC maps a SIC code to an investor-friendly sector name.
func SectorFromSIC(sicCode string) string {
	code, err := strconv.Atoi(sicCode)
	if err != nil || code == 0 {
		return "Other"
	}

	switch {
	// Technology
	case code >= 3570 && code <= 3579: // Computer & Office Equipment
		return "Technology"
	case code >= 3660 && code <= 3669: // Communications Equipment
		return "Technology"
	case code >= 3670 && code <= 3679: // Electronic Components & Semiconductors
		return "Technology"
	case code >= 3810 && code <= 3829: // Measuring & Controlling Instruments
		return "Technology"
	case code >= 7370 && code <= 7379: // Computer Services & Software
		return "Technology"
	case code >= 3812 && code <= 3812: // Defense Electronics
		return "Technology"
	case code == 3674: // Semiconductors
		return "Technology"
	case code == 5045 || code == 5065: // Computers/Electronics Wholesale
		return "Technology"
	case code == 5961: // Catalog & Mail-Order (e-commerce like Amazon)
		return "Technology"
	case code == 4813: // Telephone Communications
		return "Technology"
	case code >= 4810 && code <= 4899: // Communications
		return "Technology"

	// Healthcare
	case code >= 2830 && code <= 2836: // Drugs
		return "Healthcare"
	case code >= 2800 && code <= 2899: // Chemicals (pharma related)
		return "Healthcare"
	case code >= 3840 && code <= 3851: // Medical Instruments
		return "Healthcare"
	case code >= 5912 && code <= 5912: // Drug Stores
		return "Healthcare"
	case code >= 8000 && code <= 8099: // Health Services
		return "Healthcare"
	case code == 6324: // Hospital & Medical Service Plans (insurers like UNH)
		return "Healthcare"
	case code >= 6320 && code <= 6329: // Medical Insurance
		return "Healthcare"

	// Financial
	case code >= 6000 && code <= 6099: // Banks
		return "Financial"
	case code >= 6100 && code <= 6199: // Credit Institutions
		return "Financial"
	case code >= 6200 && code <= 6299: // Security Brokers/Dealers
		return "Financial"
	case code >= 6300 && code <= 6399: // Insurance
		return "Financial"
	case code >= 6400 && code <= 6499: // Insurance Agents
		return "Financial"
	case code >= 6500 && code <= 6599: // Real Estate
		return "Financial"
	case code >= 6700 && code <= 6799: // Holding Companies
		return "Financial"
	case code == 7389: // Services-Misc Business Services (includes payment processors like V, MA)
		return "Financial"

	// Energy
	case code >= 1300 && code <= 1399: // Oil & Gas Extraction
		return "Energy"
	case code >= 1400 && code <= 1499: // Mining (non-metallic)
		return "Energy"
	case code >= 2900 && code <= 2999: // Petroleum Refining
		return "Energy"
	case code >= 4900 && code <= 4999: // Electric, Gas, Sanitary Services
		return "Energy"
	case code == 1389: // Oil & Gas Field Services
		return "Energy"

	// Consumer
	case code >= 2000 && code <= 2111: // Food & Beverages
		return "Consumer"
	case code >= 5400 && code <= 5499: // Food Stores
		return "Consumer"
	case code >= 5800 && code <= 5899: // Eating & Drinking Places
		return "Consumer"
	case code >= 5300 && code <= 5399: // General Merchandise Stores
		return "Consumer"
	case code >= 5200 && code <= 5299: // Building Materials Retail
		return "Consumer"
	case code >= 5600 && code <= 5699: // Apparel & Accessory Stores
		return "Consumer"
	case code >= 5900 && code <= 5999: // Retail Stores
		return "Consumer"
	case code >= 5100 && code <= 5199: // Wholesale - Nondurable Goods
		return "Consumer"
	case code >= 7800 && code <= 7999: // Amusement & Recreation
		return "Consumer"
	case code >= 2100 && code <= 2199: // Tobacco Products
		return "Consumer"
	case code >= 2200 && code <= 2399: // Textile/Apparel
		return "Consumer"
	case code == 3021 || code == 3140: // Rubber/Footwear (Nike)
		return "Consumer"
	case code >= 3940 && code <= 3949: // Toys & Sporting Goods
		return "Consumer"
	case code == 3944: // Games & Toys
		return "Consumer"

	// Industrial
	case code >= 3500 && code <= 3559: // Industrial Machinery
		return "Industrial"
	case code >= 3700 && code <= 3799: // Transportation Equipment
		return "Industrial"
	case code >= 3720 && code <= 3729: // Aircraft & Parts
		return "Industrial"
	case code >= 3760 && code <= 3769: // Guided Missiles/Space Vehicles
		return "Industrial"
	case code >= 4500 && code <= 4599: // Air Transportation
		return "Industrial"
	case code >= 4200 && code <= 4299: // Trucking & Warehousing
		return "Industrial"
	case code >= 4400 && code <= 4499: // Water Transportation
		return "Industrial"
	case code >= 3400 && code <= 3499: // Fabricated Metal Products
		return "Industrial"
	case code >= 3300 && code <= 3399: // Primary Metal Industries
		return "Industrial"
	case code >= 1500 && code <= 1799: // Construction
		return "Industrial"
	case code >= 3600 && code <= 3659: // Electrical Equipment (not semiconductors)
		return "Industrial"

	default:
		return "Other"
	}
}

// SectorFromTickerDetail determines sector from Polygon ticker detail.
func SectorFromTickerDetail(detail *TickerDetail) string {
	if detail == nil {
		return "Other"
	}

	// Crypto
	if detail.Market == "crypto" {
		return "Crypto"
	}

	// ETFs
	if detail.Type == "ETF" || detail.Type == "ETN" {
		return "ETF"
	}

	// Use SIC code for stocks
	if detail.SICCode != "" {
		return SectorFromSIC(detail.SICCode)
	}

	return "Other"
}
