package shelter

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Status string

const (
	Open    Status = "open"
	Paused  Status = "paused"
	Full    Status = "full"
	Retired Status = "retired"
)

type Shelter struct {
	ID, ZoneID, Name            string
	Capacity, Occupied, Version int
	Status                      Status
}

func (s Shelter) Available() int { return s.Capacity - s.Occupied }
func (s *Shelter) Assign(people int) error {
	if people <= 0 {
		return common.ErrInvalid
	}
	if s.Status == Retired || s.Status == Paused {
		return fmt.Errorf("%w: shelter unavailable", common.ErrConflict)
	}
	if s.Available() < people {
		return fmt.Errorf("%w: shelter capacity", common.ErrConflict)
	}
	s.Occupied += people
	s.Version++
	if s.Occupied == s.Capacity {
		s.Status = Full
	}
	return nil
}
func (s *Shelter) Release(people int) error {
	if people <= 0 || people > s.Occupied {
		return common.ErrInvalid
	}
	s.Occupied -= people
	s.Version++
	if s.Status == Full {
		s.Status = Open
	}
	return nil
}
