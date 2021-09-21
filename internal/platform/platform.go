package platform

// Platform lists supported arch/os combinations
type Platform struct {
	OS   string `yaml:"os"`
	Arch string `yaml:"arch"`
}
