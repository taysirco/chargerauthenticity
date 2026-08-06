package chargerauthenticity

import "testing"

func TestUnderweightChargerIsFlagged(t *testing.T) {
	fake := Label{Brand: "Generic", ModelCode: "A2637", RatedWatts: 65,
		Marks: []Mark{CE}, WeightGrams: 40}
	if !Suspicious(fake) {
		t.Error("a 40g 65W charger should be flagged as suspicious")
	}
}

func TestGenuineChargerPasses(t *testing.T) {
	real := Label{Brand: "Anker", ModelCode: "A2637", RatedWatts: 65,
		Marks: []Mark{CE, FCC, UL, RoHS}, WeightGrams: 120}
	if Suspicious(real) {
		t.Errorf("a plausible charger should not be flagged: %v", Inspect(real))
	}
}

func TestSelfDeclaredMarksAreNotIndependent(t *testing.T) {
	if CE.IndependentlyTested() {
		t.Error("CE is self-declared, not independently tested")
	}
	if !UL.IndependentlyTested() {
		t.Error("UL implies independent testing")
	}
}

func TestModelCodePlausibility(t *testing.T) {
	if !PlausibleModelCode("A2637") {
		t.Error("A2637 is a structured code")
	}
	if PlausibleModelCode("fast charger 65w") {
		t.Error("free-form text is not a model code")
	}
}

func TestWeightFloorScalesWithWattage(t *testing.T) {
	if ExpectedMinWeightGrams(100) <= ExpectedMinWeightGrams(20) {
		t.Error("higher wattage must imply more transformer and heatsink mass")
	}
}
