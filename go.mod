module nusantarare

go 1.22

require (
	// Desimal untuk seluruh jalur uang. Keputusan DECIDED-TEKNIS 18-09
	// (dastin\_migration-docs\claim-non-prop\3-to-tickets\TICKETS.md, bab
	// T-keputusan): cockroachdb/apd dipakai, shopspring/decimal ditolak.
	// Konteks presisi 38 dinyatakan di satu tempat: pkg/utils/decimal.go.
	// ADR-U-0003 - ADR-U-0016: uang tidak pernah float.
	github.com/cockroachdb/apd/v3 v3.2.1

	// [usulan] Driver Oracle murni Go, tanpa Instant Client. BELUM pernah
	// diputuskan lewat ADR; boleh diganti lewat keputusan tertulis.
	github.com/sijms/go-ora/v2 v2.8.19
)

// [usulan] Router memakai net/http bawaan Go 1.22 (pola "GET /path"),
// sehingga tidak ada ketergantungan router pihak ketiga sama sekali.
