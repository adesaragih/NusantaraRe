package services

// Jalur Life endorsement - tiket E18.
//
// Untuk apa berkas ini: dua titik cabang Life yang tidak butuh Oracle - rute
// flow sesudah Marketing menyetujui, dan rencana + susunan nomor endorsement.
// Titik cabang lain sudah di berkasnya sendiri: lapis B keluar untuk Life
// (nilailama.go, langkah 1) dan varian lapis C `_LIFE` (lapisc.go). Jalur
// produksi Life (`SaveFacinLive_Act`, `SaveFacinSpreadLife_Sql`) milik E21
// yang ter-block tabel flat.
//
// Dibaca sesudah: lapisc.go, nilailama.go.
//
// ⛔ `[terverifikasi]` Flow `InputEDMLife` yang disebut tiket E18 TIDAK ADA di
// korpus (`grep -rl InputEDMLife "Endorsment Fac In"` → nol berkas). Life bercabang di DALAM flow
// `Flow/InputAddendumFacultativeIn.xml`, di shape `Decision19` ("IS LIFE?").
// Yang dimodelkan di sini: cabang itu, bukan flow tersendiri (keputusan A03).
//
// Predikat `IsLife` sendiri (`When/IsLife.xml`: `pyWorkPage.Quotation.
// BusinessOldId = "L1"` … `"L16"`, `pyLogic` A OR … OR P) dinilai registry NB
// lewat `kontrak.PenilaiPredikatFacIn`; di sini ia masukan.

import "strings"

// LangkahFlow - shape tujuan di `Flow/InputAddendumFacultativeIn.xml`.
type LangkahFlow string

const (
	// LangkahCekGrup - `Decision5` "Is it group?": awal cabang non-Life yang
	// berujung tangga akseptasi (`GetLimitAkseptasi_*`, assignment UW).
	LangkahCekGrup LangkahFlow = "Decision5"
	// LangkahCekGalatKonversi - `Decision39` "Err Konversi?", lalu `Utility3`
	// SaveEDMToJsonPolicy_Act. ⚠️ Shape ini TITIK TEMU: jalur non-Life juga
	// tiba di sini, tetapi SESUDAH tangga akseptasi (dari `Decision11`, `29`,
	// `31`). Life tiba di sini langsung dari `Decision19`, tanpa satu pun shape
	// tangga di antaranya.
	LangkahCekGalatKonversi LangkahFlow = "Decision39"
)

// MembukaTanggaAkseptasi - apakah langkah ini menuju tangga akseptasi.
//
// ⛔ K-044: Life melewati tangga akseptasi; korpus tidak memuat tangga untuk
// Life. `[pertanyaan terbuka]` apakah itu memang tanpa persetujuan - milik
// work owner dan Underwriting. Diport apa adanya.
func (l LangkahFlow) MembukaTanggaAkseptasi() bool { return l == LangkahCekGrup }

// LangkahSesudahKonfirmasiMarketing - konektor `Decision16` (Accept?, hasil
// `confirm`) → `Decision19` "IS LIFE?":
//
//	Decision19 -> Decision39   [IsLife]
//	Decision19 -> Decision5    [Else]
func LangkahSesudahKonfirmasiMarketing(isLife bool) LangkahFlow {
	if isLife {
		return LangkahCekGalatKonversi
	}
	return LangkahCekGrup
}

// PermintaanData - `RequestType` sebuah langkah RDB-List. Pelaksanaannya milik
// repository (E17, ter-block).
type PermintaanData string

// ⛔ Langkah 12 (`GenerateEndorsementNo`) dan 13 (`GenerateEDMNoLife`)
// berlabel `//` - di-remark, tidak dipakai lagi (keputusan work owner
// 01-10-2026 butir 43) - jadi TIDAK diport. `[terverifikasi]`
// `grep -c '<pyStepsBlockName>//</pyStepsBlockName>' SaveEDMToJsonPolicy_Act.xml`
// = 4; cara kedua `py docs/alat/langkah.py SaveEDMToJsonPolicy_Act.xml | grep -c "label='//'"`
// = 4 (langkah 4, 5, 12, 13). Itulah sebab kedua rule itu tidak ada di
// `RDBList\`: nomor endorsement disusun langkah 14-17.
const (
	// MintaKodeProdNonLife - langkah 14.1, `RDBList/GetKodeProdNonLife_SQL`.
	MintaKodeProdNonLife PermintaanData = "GetKodeProdNonLife_SQL"
	// MintaKodeProdLife - langkah 15.1, `RDBList/GetKodeProdLife_SQL`.
	MintaKodeProdLife PermintaanData = "GetKodeProdLife_SQL"
	// MintaNomorUrut - langkah 16, `RDBList/GetSequenceNumber_SQL` ("generate
	// MM.YYYY DAN SEQUENCE"; memanggil prosedur tersimpan, isi: D-1 DBA).
	MintaNomorUrut PermintaanData = "GetSequenceNumber_SQL"
)

// RencanaNomorEndorsemen - langkah 14-16 `Activity/SaveEDMToJsonPolicy_Act.xml`
// yang berjalan, berurutan (12-13 di-remark).
//
// Gerbang tiap langkah (dua baris prakondisi, `[terverifikasi]`):
//
//	14  EndorsementNo=="" T=2 F=3 · IsLife T=3 F=(kosong)
//	15  EndorsementNo=="" T=2 F=3 · IsLife T=2 F=3
//	16  EndorsementNo=="" T=2 F=3
//
// ⚠️ `[dugaan]` `WhenFalse` kosong (langkah 14, baris IsLife) dibaca "lanjut"
// - perilaku yang konsisten dengan pasangan komplementer langkah 15.
// Langkah 14.1/15.1 ber-`pyStepsPreCondition=false`: tetap jalan (P-11).
func RencanaNomorEndorsemen(isLife bool, endorsementNo string) []PermintaanData {
	if endorsementNo != "" {
		return nil
	}
	if isLife {
		return []PermintaanData{MintaKodeProdLife, MintaNomorUrut}
	}
	return []PermintaanData{MintaKodeProdNonLife, MintaNomorUrut}
}

// HasilNomorEndorsemen - nilai yang dikembalikan query penomoran.
type HasilNomorEndorsemen struct {
	// KodeProd - `ParamSeq.HASIL3` dari query kode produk (14.1/15.1).
	KodeProd string
	// BusinessOldId - `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId`.
	BusinessOldId string
	// Hasil1, Hasil2 - `ParamSeq.HASIL1`/`HASIL2` dari `GetSequenceNumber_SQL`
	// (deskripsi langkah: "MM.YYYY DAN SEQUENCE"; bentuknya belum
	// terverifikasi).
	Hasil1, Hasil2 string
}

// SusunNomorEndorsemen - langkah 14.2/15.2 `ParamSeq.CARI2 = ParamSeq.HASIL3+"E"`
// lalu langkah 17:
//
//	EndorsementNo = ParamSeq.CARI2 + QuotationData.BusinessOldId + "." + HASIL1 + "." + HASIL2
//
// Rumus yang sama untuk Life dan non-Life; yang berbeda hanya query kode
// produknya.
func SusunNomorEndorsemen(h HasilNomorEndorsemen) string {
	return strings.Join([]string{h.KodeProd + "E" + h.BusinessOldId, h.Hasil1, h.Hasil2}, ".")
}
