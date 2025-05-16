package utility

// Check if price is retesting imbalance or liquidity pool
func PriceInZone(price float64, top float64, bottom float64) bool {
	return price <= top && price >= bottom
}
