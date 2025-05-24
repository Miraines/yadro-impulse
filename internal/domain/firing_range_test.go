package domain

import "testing"

func TestFiringRange_ShotAndMissLogic(t *testing.T) {
	fr := FiringRange{}
	for i := 0; i < 4; i++ {
		fr.Shots = append(fr.Shots, Shot{Hit: true})
	}
	if fr.ShotsCount() != 5 {
		t.Fatalf("expected 5 shots, got %d", fr.ShotsCount())
	}
	if fr.MissesCount() != 1 {
		t.Fatalf("expected 1 miss, got %d", fr.MissesCount())
	}
}
