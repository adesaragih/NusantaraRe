// Berkas ini adalah Seam 2 paket premium (tiket NB-03): satu-satunya tempat yang
// menetapkan SATUAN sebuah rasio - per mille atau persen - beserta pembaginya.
// Satuan `.Rate` datang dari lini bisnis (`SatuanRate`), satuan `ProRatePercent`
// dari medannya sendiri (`SatuanProRata`, dan satuan persen medan lain di
// bawah); `rasio` hanya dibangun lewat fungsi pembangun di berkas ini (`rasioRate`,
// `rasioProRata`, `rasioPeriodePendek`, `rasioDiskon`, `rasioPersen`, `rasio.faktor`), dan
// `TestRasioHanyaDibangunDiResolver` menjaganya. Keputusan:
// `docs/KEPUTUSAN-30-09-2026.md` butir 26-27.

package premium

import (
	"fmt"

	"nusantarare/inti/backend/uang"

	"github.com/cockroachdb/apd/v3"
)

// LiniBisnis - lini bisnis (COB) pada peta skala K-018. Nilainya label K-018
// apa adanya.
//
// [pertanyaan terbuka] Kode lini bisnis di data produksi: predikat `IsPA`,
// `IsMBU`, `IsFire`, … yang mengisinya diport di tiket 10 (jembatan predikat →
// resolver). Sampai itu, pemanggil memakai konstanta di bawah, dan ejaan lain -
// termasuk beda huruf besar-kecil seperti `Fire` - ditolak `SatuanRate`.
type LiniBisnis string

// Tujuh lini bisnis K-018. Layering dicabut: tidak dipakai di sistem baru
// (keputusan work owner 01-10-2026, butir 30).
const (
	LiniPA          LiniBisnis = "PA"
	LiniFire        LiniBisnis = "FIRE"
	LiniMBU         LiniBisnis = "MBU"
	LiniAneka       LiniBisnis = "ANEKA"
	LiniBonding     LiniBisnis = "BONDING"
	LiniGolf        LiniBisnis = "GOLF"
	LiniMarineCargo LiniBisnis = "MARINE CARGO"
)

// LiniDikenal - lini ada di peta skala K-018 (SatuanRate tidak akan panic). Untuk
// pemanggil yang menerima lini dari luar (HTTP): ejaan lain ditolak sebelum mesin
// dipanggil, bukan ditebak.
func LiniDikenal(l LiniBisnis) bool {
	switch l {
	case LiniPA, LiniFire, LiniMBU, LiniAneka, LiniBonding, LiniGolf, LiniMarineCargo:
		return true
	}
	return false
}

// Satuan - satuan rasio yang MELEKAT pada nilainya (K-018, ADR-F-0004). Tidak ada
// nilai bawaan: nol bukan satuan.
type Satuan uint8

const (
	// PerMil - ‰.
	PerMil Satuan = iota + 1
	// Persen - %.
	Persen
)

// dataSatuan - simbol dan pembagi tiap satuan, di SATU tabel (K-018 aturan
// penguraian: ‰ menyumbang 1.000, % menyumbang 100 ke pembagi komposit).
var dataSatuan = map[Satuan]struct {
	simbol  string
	pembagi int64
}{
	PerMil:       {"‰", 1000},
	Persen:       {"%", 100},
	satuanFaktor: {"×", 1},
}

func (s Satuan) String() string {
	if d, ada := dataSatuan[s]; ada {
		return d.simbol
	}
	return fmt.Sprintf("Satuan(%d)", uint8(s))
}

// Pembagi - sumbangan satuan ini ke pembagi komposit. Pembagi gabungan satu
// `@Math.divide` adalah HASIL KALI sumbangan tiap faktor bersatuan - mis. 100000
// = 1.000 (rate ‰) × 100 (ProRatePercent %) - dan satuan rate dibaca dari faktor
// yang menempel padanya saja, tidak pernah dari pembagi total (K-018).
func (s Satuan) Pembagi() int64 {
	if d, ada := dataSatuan[s]; ada {
		return d.pembagi
	}
	panic(fmt.Sprintf("premium: satuan %d tidak dikenal", uint8(s)))
}

// SatuanRate - SATU-SATUNYA pengisi satuan `.Rate` menurut lini bisnis. Ia
// mengikuti RUMUS, bukan label layar: [terverifikasi] K-018 Konsekuensi 3-4
// mencatat tiga label layar yang bertentangan dengan rumusnya (MBU `Standard
// Rate (‰)`, tampilan PA `Rate (%)`, `LayerListDtl` `(%)`), dan rumus yang menang.
//
// Peta ini fakta struktural korpus, bukan parameter bisnis: ia tidak dibaca dari
// tabel konfigurasi (ADR-F-0004). Lini yang belum ada di peta → panic, bukan
// satuan bawaan (spec Modul 2).
//
// [keputusan work owner] Seluruh peta: K-018 ("DIKUNCI"). Bukti rumus yang
// dikutip K-018 sendiri hanya untuk PA, MBU, dan Layering (dicabut); lima lini lain
// bersandar pada keputusan itu, dengan penguat [dugaan] dari pembacaan sekilas
// korpus 01-10-2026. Semua berkas di bawah: `NB FacIn\Activity\`.
func SatuanRate(lini LiniBisnis) Satuan {
	switch lini {
	case LiniPA, LiniFire:
		// [terverifikasi] PA: CalculatePremiPA_FacIn.xml
		// (ASM-FW-GISFW-DATA-COVERAGE / CALCULATEPREMIPA_FACIN) L713, L858, L1003,
		// L1146 - keempatnya ‰, dikutip K-018. Terpisah dari itu: rumus yang
		// BENAR-BENAR dipakai, L713 (metode '1') dan L1003 (metode '3'), cocok
		// eksak dengan empat kasus PA nyata (`rekonsiliasi_test.go`); L1146 tidak
		// dipakai - langkahnya berlabel `//`, di-remark (butir 13, 43).
		// [keputusan work owner] FIRE: K-018. [dugaan] Konsisten dengan pembagi
		// 100000 di CopyAllObjFacOutFireAneka_ACT.xml
		// (ASM-FW-GISFW-DATA-FACOFFER / COPYALLOBJFACOUTFIREANEKA_ACT); cabang mana
		// untuk lini mana belum dibaca.
		return PerMil
	case LiniMBU, LiniAneka, LiniBonding, LiniGolf, LiniMarineCargo:
		// [terverifikasi] MBU: FillPremiMBU_FacIn.xml
		// (ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIMBU_FACIN) L1144, L1626, bentuk
		// bersarang, dikutip K-018.
		// [keputusan work owner] ANEKA, BONDING, GOLF, MARINE CARGO: K-018.
		// [dugaan] Konsisten dengan FillPremiGolf.xml
		// (ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIGOLF) L1643 ÷10000 = rate % ×
		// PctShortPeriod %, dan CountCoverageMarine.xml
		// (ASM-FW-GISFW-DATA-COVERAGE / COUNTCOVERAGEMARINE) L999 ÷100. BONDING:
		// belum ditemukan rumus pembanding.
		return Persen
	}
	panic(PanikLini{Lini: string(lini)})
}

// PanikLini - nilai panic SatuanRate untuk lini di luar peta skala: keadaan masukan
// yang diketahui (lini dari luar), bukan bug - pemanggil HTTP memilahnya lewat
// errors.As dan menjawab 422. Panic lain paket ini (satuan tak dikenal) tetap bug.
type PanikLini struct{ Lini string }

func (p PanikLini) Error() string {
	return fmt.Sprintf("premium: lini bisnis %q belum ada di peta skala", p.Lini)
}

// SatuanProRata - satuan `pyWorkPage.OfferFacIn.ProRatePercent`. Ia TIDAK
// bergantung lini bisnis: [terverifikasi] K-018 tabel sumbangan faktor
// menetapkan `ProRatePercent (persen) | 100` sebagai faktor tersendiri, terlepas
// dari satuan rate (butir 27).
func SatuanProRata() Satuan {
	return Persen
}

// rasioRate - `.Rate` coverage beserta satuannya menurut lini bisnis.
func rasioRate(lini LiniBisnis, nilai *apd.Decimal) rasio {
	return rasio{nilai: uang.Ratio{Value: nilai}, satuan: SatuanRate(lini)}
}

// rasioProRata - `ProRatePercent` beserta satuannya.
func rasioProRata(nilai *apd.Decimal) rasio {
	return rasio{nilai: uang.Ratio{Value: nilai}, satuan: SatuanProRata()}
}

// rasioPeriodePendek - `.PctShortPeriod`, persen: [terverifikasi]
// CalculatePremiPA_FacIn langkah 5 L858 membaginya dengan 100
// (`@Math.divide((.PctShortPeriod),100,4)`).
func rasioPeriodePendek(nilai *apd.Decimal) rasio {
	return rasio{nilai: uang.Ratio{Value: nilai}, satuan: Persen}
}

// rasioDiskon - `.DiscountPercentage`, persen: [terverifikasi]
// CalculatePremiPA_FacIn langkah 2 membagi `.DiscountPercentage * .Premium`
// dengan 100.
func rasioDiskon(nilai *apd.Decimal) rasio {
	return rasio{nilai: uang.Ratio{Value: nilai}, satuan: Persen}
}

// rasioPersen - medan persen lain yang menempel pada perkalian premi tiket 18:
// `.IndemnityPercentage`, `.LostLimit` / `Local.LossLimit`, `.PctAdjustment`,
// `.FirstScale`, `.Loading` (ANEKA/GOLF), dan pro-rata tetap 100 metode 3.
// [terverifikasi] tiap pembagi komposit di CountPremi_ACT langkah 19-50 dan
// CountPremiCoverageAneka langkah 14-17 terurai menjadi rate × 100 per medan ini
// (lihat lini_lain.go).
func rasioPersen(nilai *apd.Decimal) rasio {
	return rasio{nilai: uang.Ratio{Value: nilai}, satuan: Persen}
}

// satuanFaktor - rasio yang SUDAH dibagi pembagi satuannya: pecahan murni,
// pembagi 1. Hanya lahir dari `rasio.faktor`.
const satuanFaktor Satuan = 99

// faktor - `@Math.divide(rasio, pembagiSatuan, desimal)`: rasio dijadikan
// pecahan murni, dibulatkan ke `desimal` (setengah-ke-atas, seperti bagiBulat).
func (r rasio) faktor(desimal int32) (rasio, error) {
	if r.nilai.Kosong() {
		return rasio{}, ErrRasioKosong
	}
	d, err := bagiBulatDesimal(r.nilai.Value, r.satuan.Pembagi(), desimal)
	if err != nil {
		return rasio{}, err
	}
	return rasio{nilai: uang.Ratio{Value: d}, satuan: satuanFaktor}, nil
}
