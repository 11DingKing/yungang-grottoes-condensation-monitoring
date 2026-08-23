package monitoring

import (
	"context"
	"errors"
	"math"
	"sort"
	"sync"
	"time"
)

var (
	ErrInvalidSample = errors.New("invalid environmental sample")
	ErrStaleSample   = errors.New("environmental sample is stale")
)

type State string

const (
	StateStable       State = "stable"
	StateCondensation State = "condensation_risk"
	StateTreatment    State = "treatment_active"
	StateSensorFault  State = "sensor_fault"
)

type Sample struct {
	ID, CaveCode, SensorID string
	WallTemperature        float64
	RelativeHumidity       float64
	OutsideTemperature     float64
	OutsideHumidity        float64
	RecordedAt             time.Time
}

type Decision struct {
	State       State
	Action      string
	DewPoint    float64
	Margin      float64
	EvaluatedAt time.Time
}

func DewPoint(temperature, humidity float64) float64 {
	const a, b = 17.62, 243.12
	if humidity <= 0 {
		return math.Inf(-1)
	}
	gamma := math.Log(humidity/100) + a*temperature/(b+temperature)
	return b * gamma / (a - gamma)
}

func Evaluate(ctx context.Context, sample Sample, now time.Time, maxAge time.Duration) (Decision, error) {
	if err := ctx.Err(); err != nil {
		return Decision{}, err
	}
	if sample.CaveCode == "" || sample.SensorID == "" || sample.RecordedAt.IsZero() ||
		sample.RelativeHumidity <= 0 || sample.RelativeHumidity > 100 || sample.OutsideHumidity <= 0 || sample.OutsideHumidity > 100 {
		return Decision{}, ErrInvalidSample
	}
	if maxAge > 0 && now.Sub(sample.RecordedAt) > maxAge {
		return Decision{}, ErrStaleSample
	}
	dew := DewPoint(sample.WallTemperature, sample.RelativeHumidity)
	margin := sample.WallTemperature - dew
	state := StateStable
	action := "observe"
	if margin <= 0.2 {
		state = StateCondensation
		action = "enable_adaptive_dehumidification"
		if sample.OutsideTemperature < sample.WallTemperature && sample.OutsideHumidity > sample.RelativeHumidity {
			action = "isolate_inlet_and_dehumidify"
		}
	}
	return Decision{State: state, Action: action, DewPoint: dew, Margin: margin, EvaluatedAt: now.UTC()}, nil
}

type Ledger struct {
	mu        sync.RWMutex
	latest    map[string]Sample
	decisions map[string]Decision
}

func NewLedger() *Ledger {
	return &Ledger{latest: make(map[string]Sample), decisions: make(map[string]Decision)}
}

func (l *Ledger) Record(ctx context.Context, sample Sample, now time.Time, maxAge time.Duration) (Decision, error) {
	d, err := Evaluate(ctx, sample, now, maxAge)
	if err != nil {
		return Decision{}, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if old, ok := l.latest[sample.SensorID]; ok && old.RecordedAt.After(sample.RecordedAt) {
		return Decision{}, errors.New("out-of-order sample")
	}
	l.latest[sample.SensorID] = sample
	l.decisions[sample.SensorID] = d
	return d, nil
}

func (l *Ledger) Snapshot() []Sample {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Sample, 0, len(l.latest))
	for _, sample := range l.latest {
		out = append(out, sample)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SensorID < out[j].SensorID })
	return out
}

func (l *Ledger) Decision(sensorID string) (Decision, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	d, ok := l.decisions[sensorID]
	return d, ok
}
