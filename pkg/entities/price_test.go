package entities_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MaKcm14/best-price-service/price-service/pkg/entities"
)

func TestNewPricePositiveCases(t *testing.T) {
	t.Run("Positive: price div's remainder != 0", func(t *testing.T) {
		var testPriceObj = entities.NewPrice(15000, 5000)
		assert.Equal(t, entities.Price{
			BasePrice:     15000,
			DiscountPrice: 5000,
			Discount:      66,
		}, testPriceObj)
	})

	t.Run("Positive: price div's remainder = 0", func(t *testing.T) {
		var testPriceObj = entities.NewPrice(20000, 5000)
		assert.Equal(t, entities.Price{
			BasePrice:     20000,
			DiscountPrice: 5000,
			Discount:      75,
		}, testPriceObj)
	})
}

func TestNewPriceExtremeCases(t *testing.T) {
	t.Run("Extreme: the base price equals discount price", func(t *testing.T) {
		testPriceObj := entities.NewPrice(10, 10)

		assert.Equal(t, entities.Price{
			BasePrice:     10,
			DiscountPrice: 10,
			Discount:      0,
		}, testPriceObj)
	})
}

func TestNewPriceNegativeAndExtremeCases(t *testing.T) {
	type args struct {
		basePrice int
		discPrice int
	}

	tests := []struct {
		name string
		args args
		want entities.Price
	}{
		{"Negative: the wrong base price and the correct discount price", args{0, 10}, entities.Price{}},
		{"Negative: the correct base price and the wrong discount price", args{10, 0}, entities.Price{}},
		{"Negative: the wrong base and discount prices", args{0, 0}, entities.Price{}},
		{"Negative: the base price is less than the discount price", args{10, 20}, entities.Price{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.NewPrice(tt.args.basePrice, tt.args.discPrice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
