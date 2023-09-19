package lib

import "fmt"

func DisplayPrice(priceInCents int) string {
	dollars := priceInCents / 100
	cents := priceInCents % 100

	return fmt.Sprintf("%d.%02d", dollars, cents)
}
