// Package repository adalah SATU-SATUNYA lapisan modul R/I Rate Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services.
//
// Untuk apa berkas ini: nama objek Oracle dan kunci JSON modul ini, ditulis SEKALI.
//
// RALAT R6 (07-10-2026, MODUL.md) - menggantikan RALAT R4: keputusan work owner 07-10-2026 "ringkasan R/I Rate Life
// cukup SATU tabel - M_RATE_LIFE_SUMMARY". Ringkasan dibaca DAN ditulis di kolom `M_RATE_LIFE_SUMMARY` (ID, USEDBY,
// TYPE, MODIFIEDDATE, OPERATORID - migrasi inti 927/928); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926)
// dibuang 928. RALAT R7 (keputusan work owner 07-10-2026): rincian rate pun SATU tabel flat `M_RATE_LIFE` (kolom
// `KolomRate`, migrasi inti 929/930); JSONDATA rincian dan view `RATE_LIFE` dibuang 930.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - `M_RATE_LIFE_SUMMARY`, kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`; kolom = KolomRingkasan
	// (RALAT R6). Satu-satunya tabel ringkasan: ID terpakai pun diperiksa di sini saja.
	TabelRingkasan = "M_RATE_LIFE_SUMMARY"
	// TabelRate - `M_RATE_LIFE`, kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`AddToListSummary_Act` b1833); tabel flat
	// sesudah migrasi inti 929/930 (RALAT R7): kolom = KolomRate, nol JSONDATA.
	TabelRate = "M_RATE_LIFE"
)

// KolomViewRingkasan - SEJARAH: kolom VIEW warisan `RATE_LIFE_SUMMARY` (ALL_VIEWS, owner POOLDATA, dibaca WO
// 06-10-2026): `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID,
// a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`. Dipulihkan jalur mundur 926.
var KolomViewRingkasan = []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID", "FLAG"}

// KolomRingkasan - kolom `M_RATE_LIFE_SUMMARY` sesudah 927/928 = kolom view TANPA FLAG (RALAT R5: FLAG tidak digunakan;
// nilai lamanya hanya di cadangan CSV LANGKAH-WO (a)). Lebar = 926.
var KolomRingkasan = []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID"}

// LebarKolomRingkasan - lebar (byte) kolom ringkasan; diikat uji ke DDL 926, 927, dan 928_down.
var LebarKolomRingkasan = map[string]int{"ID": 10, "USEDBY": 500, "TYPE": 100, "MODIFIEDDATE": 50, "OPERATORID": 200}

// KolomRate - kolom `M_RATE_LIFE` sesudah 929/930 = PERSIS kolom view warisan `RATE_LIFE` (riwayat: `SELECT a.ID,
// a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER, a.JSONDATA.CONTRACT, a.JSONDATA.AGE,
// a.JSONDATA.RATE FROM M_RATE_LIFE a`, dibuang 930, dipulihkan 930_down). Semua TEKS seperti view (RATE berkoma atau
// bertitik desimal apa adanya).
var KolomRate = []string{"ID", "IDUSEDBY", "USEDBY", "TYPE", "GENDER", "CONTRACT", "AGE", "RATE"}

// LebarKolomRate - lebar (byte) kolom rincian (929); [terverifikasi data DEV 07-10-2026] panjang maksimum isi ID 7,
// IDUSEDBY 7, USEDBY 95, TYPE kosong, GENDER 1, CONTRACT 3, AGE 3, RATE 20.
var LebarKolomRate = map[string]int{"ID": 10, "IDUSEDBY": 10, "USEDBY": 500, "TYPE": 100, "GENDER": 10, "CONTRACT": 10,
	"AGE": 10, "RATE": 50}

// Sequence ID baru (K2 keputusan work owner 05-10-2026: dari sequence Oracle) - dibuat migrasi inti 923.
const (
	SeqRingkasan = "SEQ_M_RATE_LIFE_SUMMARY"
	SeqRate      = "SEQ_M_RATE_LIFE"
)

// Kolom ringkasan yang ditulis modul (XML `InboxSummaryRIRate`). `TYPE` TIDAK ditulis (tidak ada di XML): baris baru
// NULL, Edit tidak menimpa.
const (
	KolomUsedBy     = "USEDBY"
	KolomOperatorID = "OPERATORID"
	KolomModified   = "MODIFIEDDATE"
)

// Kolom yang dibaca (nol `SELECT *`).
const (
	kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`
	kolomRate      = `ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE`
)

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE. Nol objek baca-saja lain:
// view `RATE_LIFE` dibuang 930 (RALAT R7).
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRate}
