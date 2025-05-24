package config

import "testing"

func TestLoadConfig_Valid(t *testing.T) {
	const jsonCfg = `{
		"laps": 2, "lapLen": 4000, "penaltyLen": 150,
		"firingLines": 3, "start": "09:30:00", "startDelta": "00:00:30"
	}`
	var c Config
	if err := c.UnmarshalJSON([]byte(jsonCfg)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Laps != 2 || c.StartDelta.Seconds() != 30 {
		t.Fatalf("parsed values mismatch: %+v", c)
	}
}

func TestLoadConfig_InvalidFields(t *testing.T) {
	bad := []string{
		`{"laps":0,"lapLen":1,"penaltyLen":1,"firingLines":1,"start":"09:00:00","startDelta":"00:00:01"}`,
		`{"laps":1,"lapLen":-1,"penaltyLen":1,"firingLines":1,"start":"09:00:00","startDelta":"00:00:01"}`,
		`{"laps":1,"lapLen":1,"penaltyLen":1,"firingLines":0,"start":"09:00:00","startDelta":"00:00:01"}`,
		`{"laps":1,"lapLen":1,"penaltyLen":1,"firingLines":1,"start":"bad","startDelta":"00:00:01"}`,
		`{"laps":1,"lapLen":1,"penaltyLen":1,"firingLines":1,"start":"09:00:00","startDelta":"xx"}`,
	}
	for _, js := range bad {
		var c Config
		if err := c.UnmarshalJSON([]byte(js)); err == nil {
			t.Errorf("expected error for %s", js)
		}
	}
}
