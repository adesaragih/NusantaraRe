// Package acceptance memuat mesin akseptasi underwriting siklus New Business
// Fac In. Berkas ini (tiket NB-13): domain hasil keputusan dan SATU flag fase
// untuk jalur yang belum terverifikasi. Tangga akseptasi bentuk A (NB-11):
// `tangga.go`; bentuk B (NB-12) belum diport.
package acceptance

import (
	"errors"
	"fmt"
)

// Keputusan - nilai `ProposalAcceptStatus`, hasil keputusan underwriting.
//
// [terverifikasi] Domain dikunci K-021 ke enam nilai. Dekodernya
// NB FacIn\DecisionTable\IsUWAccepted.xml (ASM-FW-GISFW-WORK / ISUWACCEPTED):
// baris L302-L306 `1, 2, 3, 4, 9` → hasil L343-L347 `confirm, reject, ask,
// banding, revise`; `<pyDefaultResult>decline</pyDefaultResult>` L91 (nilai 7
// jatuh ke default). Salinan RNW identik byte (`cmp`, 01-10-2026; identik byte
// berarti juga identik ternormalisasi); salinan EDM berbeda isi barisnya - NB
// memakai varian NB-nya (CLAUDE.md §4.6).
//
// Penulis tiap nilai: NB FacIn\DataTransform\<nama>.xml, semuanya kelas
// ASM-FW-GISFW-Work (identitas ASM-FW-GISFW-WORK!<NAMA>), baris
// `<pyPropertiesName>.ProposalAcceptStatus` disebut per konstanta.
type Keputusan string

const (
	// Accept - 1. Penulis: DataTransform SetAkseptasiProposal L139.
	Accept Keputusan = "1"
	// Reject - 2, masih dapat dibanding (K-014). Penulis: SetRejectProposal L160.
	Reject Keputusan = "2"
	// Ask - 3. Penulis: SetAskProposal_DT L143.
	Ask Keputusan = "3"
	// Banding - 4. Penulis: SetBandingProposal_DT L186 (juga IsBanding="true",
	// .LetterNo = .BandingTo).
	Banding Keputusan = "4"
	// Decline - 7, penolakan final tanpa hak banding (K-014). Penulis:
	// SetDeclineProposal_DT L176.
	Decline Keputusan = "7"
	// Revise - 9. Penulis: SetReviseProposal L150.
	Revise Keputusan = "9"
)

var hasilIsUWAccepted = map[Keputusan]string{
	Accept: "confirm", Reject: "reject", Ask: "ask", Banding: "banding", Revise: "revise",
}

// HasilIsUWAccepted - hasil DecisionTable IsUWAccepted untuk keputusan ini;
// Decline lewat default `decline`.
func (k Keputusan) HasilIsUWAccepted() string {
	if h, ada := hasilIsUWAccepted[k]; ada {
		return h
	}
	return "decline"
}

// MatikanBindingDanRISlip - keputusan ini mematikan konfirmasi binding
// (`.OfferFacIn.ConfirmBinding = 0`) dan penerimaan R/I slip
// (`.OfferFacIn.ReceivedRiSlip = false`).
//
// [terverifikasi] Reject: NB FacIn\DataTransform\SetRejectProposal.xml
// (ASM-FW-GISFW-WORK!SETREJECTPROPOSAL) L189/L219 (K-014). Revise: efek yang SAMA
// di SetReviseProposal.xml (ASM-FW-GISFW-WORK!SETREVISEPROPOSAL) L179/L209 - tidak
// disebut K-014; diport apa adanya (keputusan agent A10, menunggu konfirmasi).
//
// Efek Reject yang lain - menjadi syarat jalur banding lewat hitungan riwayat
// berstatus REJECT (Protection_Act langkah 19-20) - milik tangga NB-11, belum
// diport.
func (k Keputusan) MatikanBindingDanRISlip() bool {
	return k == Reject || k == Revise
}

// Fase - SATU flag fase (K-008, ADR-F-0002), bukan dua basis kode.
type Fase uint8

const (
	// ParalelRun - sistem baru berjalan berdampingan dengan Pega.
	ParalelRun Fase = iota + 1
	// Produksi - sistem baru menggantikan Pega.
	Produksi
)

// Pencatat - penerima catatan wajib. Nil = tidak dicatat (hanya uji).
type Pencatat func(string)

// wajibSah - Fase nol (tidak diisi) atau tak dikenal TIDAK diperlakukan diam-diam
// sebagai salah satu fase: perilaku dua fase berlawanan arah.
func (f Fase) wajibSah() {
	if f != ParalelRun && f != Produksi {
		panic(fmt.Sprintf("acceptance: fase %d tidak dikenal", uint8(f)))
	}
}

func (c Pencatat) catat(s string) {
	if c != nil {
		c(s)
	}
}

// ErrDiLuarDomain - nilai ProposalAcceptStatus di luar enam nilai K-021. Di
// produksi galat ini wajib menghentikan jalur tulis.
var ErrDiLuarDomain = errors.New("acceptance: ProposalAcceptStatus di luar domain K-021")

// PeriksaDomain - validasi domain. ⚠️ PERILAKU BARU: sistem lama tidak punya
// validasi domain pada kolom ini (selain wajib-isi), jadi nilai di luar domain
// tidak pernah ditolak diam-diam - selalu dicatat. Paralel run: panic (premis
// salah, dan paralel run ada untuk menemukannya). Produksi: ErrDiLuarDomain.
func PeriksaDomain(fase Fase, nilai string, c Pencatat) (Keputusan, error) {
	fase.wajibSah()
	k := Keputusan(nilai)
	if _, ada := hasilIsUWAccepted[k]; ada || k == Decline {
		return k, nil
	}
	pesan := fmt.Sprintf("ProposalAcceptStatus %q di luar domain K-021 {1,2,3,4,7,9}", nilai)
	c.catat(pesan)
	if fase == ParalelRun {
		panic("acceptance: " + pesan)
	}
	return "", fmt.Errorf("%w: %q", ErrDiLuarDomain, nilai)
}

// JalurBelumTerverifikasi - situasi (a): keputusan dikenali, tetapi baris
// keputusannya belum diverifikasi terhadap ekspor produksi. Paralel run: panic.
// Produksi: Decline + catatan - `decline` adalah `<pyDefaultResult>` IsUWAccepted
// yang terekam, jadi memilihnya reproduksi, bukan tebakan. Transisi suatu jalur
// dari panic ke Decline hanya setelah jalur itu diverifikasi (syarat K-008).
func JalurBelumTerverifikasi(fase Fase, k Keputusan, alasan string, c Pencatat) Keputusan {
	fase.wajibSah()
	pesan := fmt.Sprintf("jalur keputusan %q belum terverifikasi: %s", k, alasan)
	if fase == ParalelRun {
		panic("acceptance: " + pesan)
	}
	c.catat(pesan + " - diperlakukan decline")
	return Decline
}

// KondisiTakDikenali - situasi (b): kondisi tidak dikenali sama sekali → panic
// di KEDUA fase.
func KondisiTakDikenali(_ Fase, alasan string) {
	panic("acceptance: kondisi tidak dikenali: " + alasan)
}
