package settlement

import "github.com/filspay/settlement/internal/money"

// NetPayableFils derives the net payout from frozen gross and fee totals.
//
// Fix: BUG-002 — Replaced float32 arithmetic with integer subtraction.
// float32 division (grossFils/1000 - feeFils/1000) then math.Round can lose
// precision on monetary values; since both inputs are already integer fils,
// net is simply gross - fee with no rounding needed.
func NetPayableFils(grossFils, feeFils money.Fils) money.Fils {
	return grossFils - feeFils
}
