package planning

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"math"
	"sort"
)

type Point struct{ Lat, Lon float64 }

func (p Point) Valid() bool { return p.Lat >= -90 && p.Lat <= 90 && p.Lon >= -180 && p.Lon <= 180 }
func Distance(a, b Point) float64 {
	if !a.Valid() || !b.Valid() {
		return 0
	}
	dx := (a.Lat - b.Lat) * 111
	dy := (a.Lon - b.Lon) * 111 * math.Cos(a.Lat*math.Pi/180)
	return math.Sqrt(dx*dx + dy*dy)
}

type Stop struct {
	ID       string
	Point    Point
	Load     int
	Priority int
}

func SortStops(stops []Stop) []Stop {
	out := append([]Stop(nil), stops...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Priority > out[j].Priority })
	return out
}
func ValidateStops(stops []Stop) error {
	if len(stops) == 0 {
		return common.ErrInvalid
	}
	seen := map[string]bool{}
	for _, s := range stops {
		if s.ID == "" || seen[s.ID] || !s.Point.Valid() || s.Load <= 0 {
			return fmt.Errorf("%w: stop", common.ErrInvalid)
		}
		seen[s.ID] = true
	}
	return nil
}
func TotalLoad(stops []Stop) int {
	n := 0
	for _, s := range stops {
		n += s.Load
	}
	return n
}
