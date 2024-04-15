package mapping

import "testing"

func TestRemapper_MustInterpolate(t *testing.T) {

	tests := []struct {
		name    string
		entries map[string]string
		k       string
		want    string
	}{
		{name: "arch", entries: map[string]string{"amd64": "x86_64"}, k: "amd64", want: "x86_64"},
		{name: "same when empty", entries: map[string]string{}, k: "amd64", want: "amd64"},
		{name: "os+arch", entries: map[string]string{"linux": "Linux", "amd64": "x86_64"}, k: "amd64", want: "x86_64"},
		{name: "os+arch", entries: map[string]string{"linux": "Linux", "amd64": "x86_64"}, k: "linux", want: "Linux"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Remapper(tt.entries)

			if got := r.MustInterpolate(tt.k); got != tt.want {
				t.Errorf("Remapper.MustInterpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemapper_IsZero(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string
		want    bool
	}{
		{name: "empty", entries: map[string]string{}, want: true},
		{name: "nil", want: true},
		{name: "not empty", entries: map[string]string{"amd64": "AMD64"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Remapper(tt.entries)

			if got := r.IsZero(); got != tt.want {
				t.Errorf("Remapper.MustInterpolate() = %v, want %v", got, tt.want)
			}
		})
	}
}
