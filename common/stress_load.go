package common

// StressLoad type for describe stress load params
type StressLoad struct {
	Instances  int
	Duration   int
	Type       string
	From       int
	To         int
	ShootCount int `yaml:"-"`
}

// GetShootCount return correct shoot count
func (s *StressLoad) GetShootCount() int {
	return s.ShootCount
}
