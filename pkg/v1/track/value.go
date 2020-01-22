package track

import (
	"fmt"

	"github.com/pbarker/logger"
	g "gorgonia.org/gorgonia"
)

// TrackedValue is a value being tracked.
type TrackedValue interface {
	// Name of the value.
	Name() string

	// Scalar value.
	Scalar() float64

	// Print the value.
	Print()

	// Data converts the value to a historical value.
	Data(episode, timestep int) *HistoricalValue

	// Aggregator is the aggregator for this value.
	Aggregator() Aggregator
}

// TrackedValueOpt is an option for a tracked value.
type TrackedValueOpt func(TrackedValue)

// WithIndex sets an index to use if the given value is non scalar.
// Defaults to 0.
func WithIndex(index int) func(TrackedValue) {
	return func(t TrackedValue) {
		switch val := t.(type) {
		case *TrackedNodeValue:
			val.index = index
		case *TrackedScalarValue:
			val.index = index
		}
	}
}

// WithAggregator sets an aggregator to use with the value.
// Default to MeanAggregator.
func WithAggregator(aggregator Aggregator) func(TrackedValue) {
	return func(t TrackedValue) {
		switch val := t.(type) {
		case *TrackedNodeValue:
			val.aggregator = aggregator
		case *TrackedScalarValue:
			val.aggregator = aggregator
		}
	}
}

// WithNamespace adds a namespace to the tracked value.
func WithNamespace(namespace string) func(TrackedValue) {
	return func(t TrackedValue) {
		switch val := t.(type) {
		case *TrackedNodeValue:
			val.name = namespaceValue(namespace, val.name)
		case *TrackedScalarValue:
			val.name = namespaceValue(namespace, val.name)
		}
	}
}

func namespaceValue(namespace, name string) string {
	return fmt.Sprintf("%s_%s", namespace, name)
}

// TrackedNodeValue is a tracked node value.
type TrackedNodeValue struct {
	// name of the tracked node value.
	name string

	// value of the tracked node value.
	value g.Value

	// index of the value.
	index int

	// aggregator for the value.
	aggregator Aggregator
}

// NewTrackedNodeValue returns a new tracked value.
func NewTrackedNodeValue(name string, opts ...TrackedValueOpt) *TrackedNodeValue {
	var val g.Value
	tv := &TrackedNodeValue{
		name:  name,
		value: val,
	}
	for _, opt := range opts {
		opt(tv)
	}
	if tv.aggregator == nil {
		tv.aggregator = MeanAggregator
	}
	return tv
}

// Name of the value.
func (t *TrackedNodeValue) Name() string {
	return t.name
}

// Scalar value.
func (t *TrackedNodeValue) Scalar() float64 {
	if t.value == nil {
		return 0.0
	}
	data := t.value.Data()
	return toF64(data, t.index)
}

// Print the value.
func (t *TrackedNodeValue) Print() {
	logger.Infov(t.name, t.Scalar())
}

// Data converts the value to a historical value.
func (t *TrackedNodeValue) Data(episode, timestep int) *HistoricalValue {
	return &HistoricalValue{
		Name:         t.name,
		TrackedValue: t.Scalar(),
		Timestep:     timestep,
		Episode:      episode,
	}
}

// Aggregator returns the aggregator for this value.
func (t *TrackedNodeValue) Aggregator() Aggregator {
	return t.aggregator
}

// TrackedScalarValue is a tracked value that can be convertible to float64.
type TrackedScalarValue struct {
	// name of the tracked value.
	name string

	// value of the tracked value.
	value interface{}

	// index of the scalar.
	index int

	// aggregator for the value.
	aggregator Aggregator
}

// NewTrackedScalarValue returns a new tracked value.
func NewTrackedScalarValue(name string, value interface{}, opts ...TrackedValueOpt) *TrackedScalarValue {
	tv := &TrackedScalarValue{
		name:  name,
		value: value,
	}
	for _, opt := range opts {
		opt(tv)
	}
	if tv.aggregator == nil {
		tv.aggregator = MeanAggregator
	}
	return tv
}

// Name of the value.
func (t *TrackedScalarValue) Name() string {
	return t.name
}

// Scalar value.
func (t *TrackedScalarValue) Scalar() float64 {
	return toF64(t.value, t.index)
}

// Print the value.
func (t *TrackedScalarValue) Print() {
	logger.Infov(t.name, t.Scalar())
}

// Data takes the current tracked value and returns a historical value.
func (t *TrackedScalarValue) Data(episode, timestep int) *HistoricalValue {
	return &HistoricalValue{
		Name:         t.name,
		TrackedValue: t.Scalar(),
		Timestep:     timestep,
		Episode:      episode,
	}
}

// Aggregator returns the aggregator for this value.
func (t *TrackedScalarValue) Aggregator() Aggregator {
	return t.aggregator
}

func toF64(data interface{}, index int) float64 {
	var ret float64
	switch val := data.(type) {
	case float64:
		ret = val
	case []float64:
		ret = val[index]
	case float32:
		ret = float64(val)
	case []float32:
		ret = float64(val[index])
	case int:
		ret = float64(val)
	case []int:
		ret = float64(val[index])
	case int32:
		ret = float64(val)
	case []int32:
		ret = float64(val[index])
	case int64:
		ret = float64(val)
	case []int64:
		ret = float64(val[index])
	case []interface{}:
		ret = toF64(val[index], index)
	default:
		logger.Fatalf("unknown type %T %v could not cast to float64", val, val)
	}
	return ret
}
