// Package repository adalah SATU-SATUNYA lapisan modul R/I Rate Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services.
//
// Untuk apa berkas ini: nama objek Oracle dan kunci JSON modul ini, ditulis SEKALI.
//
// RALAT R6 (07-10-2026, MODUL.md) - menggantikan RALAT R4: keputusan work owner 07-10-2026 "ringkasan R/I Rate Life
// cukup SATU tabel - M_RATE_LIFE_SUMMARY". Ringkasan dibaca DAN ditulis di kolom `M_RATE_LIFE_SUMMARY` (ID, USEDBY,
// TYPE, MODIFIEDDATE, OPERATORID - migrasi inti 927/928); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926)
// dibuang 928. Rincian rate (`M_RATE_LIFE` + view `RATE_LIFE`) TIDAK berubah.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - `M_RATE_LIFE_SUMMARY`, kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`; kolom = KolomRingkasan
	// (RALAT R6). Satu-satunya tabel ringkasan: ID terpakai pun diperiksa di sini saja.
	TabelRingkasan = "M_RATE_LIFE_SUMMARY"
	// TabelRate - kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`AddToListSummary_Act` b1833); `ID VARCHAR2(10)`,
	// `JSONDATA CLOB` (katalog DEV).
	TabelRate = "M_RATE_LIFE"
	// KolomJSON - kolom CLOB berisi JSON di `M_RATE_LIFE` (rincian).
	KolomJSON = "JSONDATA"
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

// KolomViewRate - kolom view `RATE_LIFE` (katalog DEV, `modul/claimlife/docs/KATALOG-TABEL-PESERTA-DAN-TREATY.md`);
// selain ID = kunci `M_RATE_LIFE.JSONDATA`.
var KolomViewRate = []string{"ID", "IDUSEDBY", "USEDBY", "TYPE", "GENDER", "CONTRACT", "AGE", "RATE"}

// ViewRate - view warisan yang DIBACA (tidak pernah ditulis, tidak di-DROP): `SELECT a.ID, a.JSONDATA.IDUSEDBY,
// a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER, a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM
// M_RATE_LIFE a` (katalog DEV).
const ViewRate = "RATE_LIFE"

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

// Kunci JSON `M_RATE_LIFE.JSONDATA` (terbukti dari definisi view `RATE_LIFE`). `TYPE` tidak ditulis (NULL di
// seluruh baris DEV). `USEDBY` = kunci salinan nama.
const (
	JSONUsedBy   = "USEDBY"
	JSONIDUsedBy = "IDUSEDBY"
	JSONGender   = "GENDER"
	JSONContract = "CONTRACT"
	JSONAge      = "AGE"
	JSONRate     = "RATE"
)

// Kolom yang dibaca (nol `SELECT *`).
const (
	kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`
	kolomRate      = `ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE`
)

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRate}

// DaftarDibacaSaja - hanya SELECT: view rincian.
var DaftarDibacaSaja = []string{ViewRate}
