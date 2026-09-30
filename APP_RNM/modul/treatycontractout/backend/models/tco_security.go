package models

// Security reinsurer - tiket 06 Treaty Contract Out.
//
// Untuk apa berkas ini: gerbang satu baris security di bawah seorang
// reinsurer (`SaveSecurityReinsurer_Act`) dengan struktur bersih
// (penyimpangan sadar 5): PK surrogate, `REAS_SECURITY` atribut biasa.
//
// Bukti gerbang (`Section/InputTreatyContractReinsType.xml`):
//
//	b19648 `Security Name` wajib (b19642/b19693), pemilih `BrowseAgentReinsSOA_RD`
//	       b19711 -> `.ID` ke `REAS_SECURITY` (b19739), `.ClientName` tampil
//	b19888 `%Share` - onchange `SetErrorMessageReinsurer` b19952
//
// ⚠️ `SetErrorMessageReinsurer` memeriksa `InputTreatyReinsurer.PctShare`
// (b518), BUKAN share security - residu salin-tempel; di Pega share security
// tidak diperiksa sama sekali. Maksud gerbangnya (0..100, `SetErrorMessageBetween`)
// ditegakkan di sini, dan AC 17 ("beserta porsinya") menjadikannya wajib -
// PENYIMPANGAN SADAR dari Pega [keputusan work owner 29-09-2026] (OQ-TCO-17, ditutup).
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// LabelSecurityNameTCO - label form b19648.
const LabelSecurityNameTCO = "Security Name"

// LabelShareSecurityTCO - label form b19888.
const LabelShareSecurityTCO = "%Share"

// LebarReasSecurityTCO - `REAS_SECURITY CHAR(10) NOT NULL` `[data DBA]`.
const LebarReasSecurityTCO = 10

var (
	// ErrSecurityMelampauiLebar - kode security tidak muat di CHAR(10) warisan
	// (temuan /code-review: tanpa gerbang ini Oracle menjawab ORA-12899 = 500).
	ErrSecurityMelampauiLebar = errors.New("models: security code exceeds the 10 characters of column REAS_SECURITY")
	// ErrSecurityKosong - `Security Name` wajib (b19642/b19693).
	ErrSecurityKosong = errors.New("models: security must be selected")
	// ErrSecurityTanpaReinsurer - security selalu menggantung pada reinsurer (AC 17).
	ErrSecurityTanpaReinsurer = errors.New("models: security must belong to a reinsurer")
)

// PeriksaSecurityTCO menjalankan gerbang wajib-isi satu baris security.
func PeriksaSecurityTCO(s SecurityReinsurer) error {
	if strings.TrimSpace(s.ReasID) == "" {
		return ErrSecurityTanpaReinsurer
	}
	if strings.TrimSpace(s.ReasSecurity) == "" {
		return fmt.Errorf("%w: %s", ErrSecurityKosong, LabelSecurityNameTCO)
	}
	if len(strings.TrimSpace(s.ReasSecurity)) > LebarReasSecurityTCO {
		return fmt.Errorf("%w: %s %q", ErrSecurityMelampauiLebar, LabelSecurityNameTCO, s.ReasSecurity)
	}
	return nil
}

// UraiShareSecurityTCO menormalkan `%Share` di batas masukan: wajib, desimal
// persis, 0..100 - aturan persen yang sama dengan share reinsurer (tiket 05).
func UraiShareSecurityTCO(teks string) (*apd.Decimal, error) {
	return UraiPersenMasukTCO(LabelShareSecurityTCO, teks)
}
