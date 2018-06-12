package algorithms

import (
	"fmt"
	"math"
	"time"

	"github.com/patrickmn/go-cache"
)

var epoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

var c *cache.Cache

func init() {
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// Hot calculate hot score
func Hot(ups float64, downs float64, date time.Time) float64 {
	return _hot(ups, downs, epochSeconds(date))
}

// Controversy calculate controversy score
func Controversy(ups, downs float64) float64 {
	if downs <= 0 || ups <= 0 {
		return 0
	}
	magnitude := ups + downs
	var balance float64
	if ups > downs {
		balance = downs / ups
	} else {
		balance = ups / downs
	}
	return math.Pow(magnitude, balance)
}

// Confidence calculate confidence score
func Confidence(ups, downs float64) float64 {
	if ups+downs == 0 {
		return 0
	}
	key := fmt.Sprintf("%f-%f", ups, downs)
	result, found := c.Get(key)
	if found {
		return result.(float64)
	}

	confidenceScore := _confidence(ups, downs)
	c.Set(key, confidenceScore, -1)
	return confidenceScore
}

func epochSeconds(date time.Time) float64 {
	td := date.Sub(epoch)
	return float64(td.Nanoseconds()) / 1000000000
}

func score(ups float64, downs float64) float64 {
	return ups - downs
}

func _hot(ups, downs, date float64) float64 {
	s := score(ups, downs)
	var sign float64
	order := math.Log10(math.Max(math.Abs(s), 1))
	if s > 0 {
		sign = 1
	} else if s < 0 {
		sign = -1
	} else {
		sign = 0
	}
	seconds := date - 1134028003
	return round(sign*order+seconds/45000, 7)
}

func round(val float64, prec int) float64 {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}

func _confidence(ups, downs float64) float64 {
	n := ups + downs
	if n == 0 {
		return 0
	}

	z := 1.281551565545 // 80% confidence
	p := ups / n

	left := p + 1/(2*n)*z*z
	right := z * math.Sqrt(p*(1-p)/n+z*z/(4*n*n))
	under := 1 + 1/n*z*z

	return (left - right) / under
}
