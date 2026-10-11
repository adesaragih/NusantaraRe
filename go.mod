module nusantarare

go 1.22

require (
	// Desimal untuk seluruh jalur uang. Keputusan DECIDED-TEKNIS 18 September
	// 2026, dastin\_migration-docs\claim-non-prop\3-to-tickets\TICKETS.md
	// bab "T-38 - Tipe desimal di sisi Golang" (baris 654, tabel putusan
	// baris 661-666): cockroachdb/apd dipakai, shopspring/decimal ditolak.
	// Konteks presisi 38 dinyatakan di satu tempat: inti/utils/decimal.go.
	// ADR-U-0003 - ADR-U-0016: uang tidak pernah float.
	github.com/cockroachdb/apd/v3 v3.2.1

	// PDF dokumen akseptasi Komite Claim Prop (PrintFileAcceptance_TKMT S11
	// HTMLToPDF - mesin PDF platform Pega tanpa padanan). Keputusan work owner
	// 08-10-2026 (OQ-KCP-07 "A"). Go murni, lisensi MIT; paket utamanya hanya
	// memakai pustaka standar dan menuntut go 1.20 (baris `go` tidak naik).
	github.com/go-pdf/fpdf v0.9.0

	// [usulan] Driver Oracle murni Go, tanpa Instant Client. BELUM pernah
	// diputuskan lewat ADR; boleh diganti lewat keputusan tertulis.
	github.com/sijms/go-ora/v2 v2.8.19

	// bcrypt untuk hash sandi login M_LOGIN_GO (keputusan work owner
	// 01-10-2026). v0.33.0: rilis terakhir yang menuntut go 1.20 - rilis
	// sesudahnya menuntut go 1.23 dan akan menaikkan baris `go` di atas.
	golang.org/x/crypto v0.33.0
)

// [usulan] Router memakai net/http bawaan Go 1.22 (pola "GET /path"),
// sehingga tidak ada ketergantungan router pihak ketiga sama sekali.
