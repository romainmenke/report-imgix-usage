package prompt

func cost(bandwidth int, images int) float64 {
	return ((float64(bandwidth) / (1024 * 1024 * 1024)) * 0.08) + ((float64(images) / 1000) * 3)
}
