package mapping

// Remapper contains list definition
type Remapper struct {
	Entries map[string]string `yaml:"map"`
}

// Mapper should return a list of available release versions
type Mapper interface {
	MustInterpolate(key string) string
	IsZero() bool
}

// MustInterpolate interpolates a key
func (r Remapper) MustInterpolate(k string) string {
	if v, ok := r.Entries[k]; ok {
		return v
	}
	return k
}

// IsZero returns true if the map is empty
func (r Remapper) IsZero() bool {
	if r.Entries == nil {
		return true
	}
	return false
}
