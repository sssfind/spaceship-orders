package env

import "os"

type demoFault struct {
	mode string
}

func NewDemoFaultConfig() *demoFault {
	mode := os.Getenv("DEMO_FAULT")
	if mode == "" {
		mode = "off"
	}
	return &demoFault{mode: mode}
}

func (cfg *demoFault) Mode() string {
	return cfg.mode
}
