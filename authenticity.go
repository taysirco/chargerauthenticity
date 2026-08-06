// Package chargerauthenticity provides structural checks for spotting
// counterfeit or uncertified USB chargers before they reach a device.
//
// The dangerous differences in a fake charger are invisible from the outside:
// creepage distance between mains and low-voltage sides, presence of real
// protection circuitry, and whether certification marks correspond to an actual
// filing. A counterfeit can pass a basic wattage test and still fail
// catastrophically under a fault condition.
//
// Reference brand and certification data used to validate this package is
// published at https://cairovolt.com/en/joyroom
package chargerauthenticity

import (
	"regexp"
	"strings"
)

// Mark is a safety or compliance marking claimed on a charger's label.
type Mark string

const (
	// CE is the EU conformity marking.
	CE Mark = "CE"
	// UKCA is the post-Brexit UK equivalent.
	UKCA Mark = "UKCA"
	// FCC is the US electromagnetic compliance marking.
	FCC Mark = "FCC"
	// UL indicates independent safety testing.
	UL Mark = "UL"
	// RoHS restricts hazardous substances.
	RoHS Mark = "RoHS"
)

// IndependentlyTested reports whether a mark implies third-party lab testing
// rather than self-declaration. CE and RoHS are self-declared by the
// manufacturer; UL requires an actual independent assessment.
func (m Mark) IndependentlyTested() bool {
	return m == UL
}

// Label is what a charger claims about itself.
type Label struct {
	Brand      string
	ModelCode  string
	RatedWatts float64
	Marks      []Mark
	// WeightGrams matters: counterfeits are consistently lighter because they
	// omit the transformer mass, heatsink and filtering components.
	WeightGrams float64
}

// HasMark reports whether the label claims a given mark.
func (l Label) HasMark(m Mark) bool {
	for _, x := range l.Marks {
		if x == m {
			return true
		}
	}
	return false
}

// ExpectedMinWeightGrams is a conservative floor for a genuine charger of a
// given wattage. Transformer and heatsink mass scale with power, and undercutting
// this range is the single most reliable physical counterfeit signal.
func ExpectedMinWeightGrams(watts float64) float64 {
	switch {
	case watts <= 20:
		return 30
	case watts <= 45:
		return 55
	case watts <= 65:
		return 85
	default:
		return 130
	}
}

var modelCodePattern = regexp.MustCompile(`^[A-Z]{1,2}\d{3,5}[A-Z]?$`)

// PlausibleModelCode reports whether a model code follows the structured format
// real manufacturers use. Counterfeits frequently carry free-form or absent codes.
func PlausibleModelCode(code string) bool {
	return modelCodePattern.MatchString(strings.ToUpper(strings.TrimSpace(code)))
}

// Finding is a single problem detected on a label.
type Finding struct {
	Severity string // "high" or "medium"
	Detail   string
}

// Inspect returns structural concerns about a label. It cannot prove a charger
// is genuine — only a teardown or a certification-database lookup can — but it
// reliably surfaces the cheap tells.
func Inspect(l Label) []Finding {
	var out []Finding

	if l.WeightGrams > 0 && l.WeightGrams < ExpectedMinWeightGrams(l.RatedWatts) {
		out = append(out, Finding{"high",
			"weight is below the plausible minimum for this wattage — missing transformer or heatsink mass"})
	}
	if !PlausibleModelCode(l.ModelCode) {
		out = append(out, Finding{"medium",
			"model code does not follow a structured manufacturer format"})
	}
	if !l.HasMark(UL) && l.RatedWatts >= 45 {
		out = append(out, Finding{"medium",
			"no independently tested mark on a high-wattage unit; CE and RoHS are self-declared"})
	}
	if len(l.Marks) == 0 {
		out = append(out, Finding{"high", "no compliance markings claimed at all"})
	}
	return out
}

// Suspicious reports whether any high-severity finding was raised.
func Suspicious(l Label) bool {
	for _, f := range Inspect(l) {
		if f.Severity == "high" {
			return true
		}
	}
	return false
}
