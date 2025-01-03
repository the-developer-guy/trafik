package main

import "fmt"

func magnitude(num int64) string {
	if num < 0 {
		num *= -1
	}

	if num > 1_000_000_000_000 {
		return fmt.Sprintf("%dT", num/1_000_000_000_000)
	} else if num > 1_000_000_000 {
		return fmt.Sprintf("%dG", num/1_000_000_000)
	} else if num > 1_000_000 {
		return fmt.Sprintf("%dM", num/1_000_000)
	} else if num > 1_000 {
		return fmt.Sprintf("%dk", num/1_000)
	} else {
		return fmt.Sprintf("%d", num)
	}
}

func magnitudeWithPrecision(num int64) string {
	if num < 0 {
		num *= -1
	}

	if num > 1_000_000_000_000 {
		return fmt.Sprintf("%.3fT", float64(num)/1_000_000_000_000)
	} else if num > 1_000_000_000 {
		return fmt.Sprintf("%.3fG", float64(num)/1_000_000_000)
	} else if num > 1_000_000 {
		return fmt.Sprintf("%.3fM", float64(num)/1_000_000)
	} else if num > 1_000 {
		return fmt.Sprintf("%.3fk", float64(num)/1_000)
	} else {
		return fmt.Sprintf("%d", num)
	}
}
