// Package repository adalah SATU-SATUNYA lapisan modul R/I Rate Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services.
//
// Untuk apa berkas ini: nama objek Oracle dan kunci JSON modul ini, ditulis SEKALI.
//
// K1 keputusan work owner 05-10-2026: CRUD menulis tabel fisik `M_RATE_LIFE_SUMMARY`; view `RATE_LIFE_SUMMARY` TIDAK
// di-DROP dan dipakai membaca grid; nol DDL pada tabel/view warisan.
//
// ⚠️ ASUMSI (MODUL.md "Asumsi terbuka" A1): definisi view `RATE_LIFE_SUMMARY` dan kolom `M_RATE_LIFE_SUMMARY` BELUM
// terbukti. Dianggap `M_RATE_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB)` dengan kunci JSON `USEDBY`, `OPERATORID`,
// `MODIFIEDDATE` - pola yang terbukti untuk `M_RATE_LIFE` / `RATE_LIFE`. WO/DBA wajib memeriksa
// `SELECT TEXT FROM ALL_VIEWS WHERE VIEW_NAME='RATE_LIFE_SUMMARY'` sebelum menyalakan menu; bila berbeda, ubah HANYA
// konstanta di berkas ini.
package repository

// Tabel fisik yang DITULIS (JSON Pega; kunci lain milik Pega dipertahankan `JSON_MERGEPATCH`).
const (
	// TabelRingkasan - kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY` (339 baris DEV).
	TabelRingkasan = "M_RATE_LIFE_SUMMARY"
	// TabelRate - kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`AddToListSummary_Act` b1833); `ID VARCHAR2(10)`,
	// `JSONDATA CLOB` (katalog DEV).
	TabelRate = "M_RATE_LIFE"
	// KolomJSON - kolom CLOB berisi JSON di kedua tabel.
	KolomJSON = "JSONDATA"
)

// View warisan yang DIBACA (tidak pernah ditulis, tidak di-DROP).
const (
	// ViewRingkasan - `RATE_LIFE_SUMMARY` (6 kolom; dibaca juga `mastercontractretrolife`, `masterproductnamelife`).
	ViewRingkasan = "RATE_LIFE_SUMMARY"
	// ViewRate - `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
	// a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM M_RATE_LIFE a` (katalog DEV).
	ViewRate = "RATE_LIFE"
)

// Sequence ID baru (K2 keputusan work owner 05-10-2026: dari sequence Oracle) - dibuat migrasi inti 923.
const (
	SeqRingkasan = "SEQ_M_RATE_LIFE_SUMMARY"
	SeqRate      = "SEQ_M_RATE_LIFE"
)

// Kunci JSON `M_RATE_LIFE_SUMMARY.JSONDATA` (ASUMSI A1) = nama kolom view `RATE_LIFE_SUMMARY` yang dibaca grid XML.
const (
	JSONUsedBy     = "USEDBY"
	JSONOperatorID = "OPERATORID"
	JSONModified   = "MODIFIEDDATE"
)

// Kunci JSON `M_RATE_LIFE.JSONDATA` (terbukti dari definisi view `RATE_LIFE`). `TYPE` tidak ditulis (NULL di
// seluruh baris DEV).
const (
	JSONIDUsedBy = "IDUSEDBY"
	JSONGender   = "GENDER"
	JSONContract = "CONTRACT"
	JSONAge      = "AGE"
	JSONRate     = "RATE"
)

// Kolom view yang dibaca (nol `SELECT *`).
const (
	kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`
	kolomRate      = `ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE`
)

// DaftarTabelDitulis - penjaga modul: hanya dua tabel fisik ini yang menerima INSERT/UPDATE/DELETE.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRate}

// DaftarViewDibacaSaja - hanya SELECT.
var DaftarViewDibacaSaja = []string{ViewRingkasan, ViewRate}
