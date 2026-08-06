# chargerauthenticity

Structural checks for spotting counterfeit or uncertified USB chargers before they reach a device.

## Install

```bash
go get github.com/taysirco/chargerauthenticity
```

## Usage

```go
import "github.com/taysirco/chargerauthenticity"

fake := chargerauthenticity.Label{
    Brand: "Generic", ModelCode: "fast charger 65w", RatedWatts: 65,
    Marks: []chargerauthenticity.Mark{chargerauthenticity.CE}, WeightGrams: 40,
}
chargerauthenticity.Suspicious(fake)   // true
chargerauthenticity.Inspect(fake)      // weight below plausible minimum; unstructured model code

chargerauthenticity.CE.IndependentlyTested()  // false — self-declared
chargerauthenticity.UL.IndependentlyTested()  // true
```

## Weight is the cheapest reliable tell

Counterfeits are consistently lighter than genuine units of the same claimed wattage, because the mass lives in the components they omit: transformer core, heatsink and input filtering. A 65 W charger weighing 40 g has not miniaturised — it has left things out.

## Not all marks mean the same thing

[CE marking](https://en.wikipedia.org/wiki/CE_marking) and RoHS are **self-declared**: the manufacturer asserts conformity, and no independent lab necessarily saw the product. A [UL](https://en.wikipedia.org/wiki/UL_(safety_organization)) mark implies an actual third-party assessment. On a high-wattage charger that difference matters, because the failure mode being guarded against is mains voltage reaching a device — or a person.

## What this package cannot do

It cannot prove a charger is genuine. Only a teardown, or a lookup against the issuing body's certification database, does that. What it does is surface the cheap tells quickly, so a unit that fails these checks never gets plugged in.

## Reference data

Brand model-code formats, weights and certification data used to validate this package come from the [Joyroom product specifications](https://cairovolt.com/en/joyroom) published by CairoVolt.

## License

MIT
