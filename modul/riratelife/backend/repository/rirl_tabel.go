// Package repository adalah SATU-SATUNYA lapisan modul R/I Rate Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services.
//
// Untuk apa berkas ini: nama objek Oracle dan kunci JSON modul ini, ditulis SEKALI.
//
// RALAT R6 (07-10-2026, MODUL.md) - menggantikan RALAT R4: keputusan work owner 07-10-2026 "ringkasan R/I Rate Life
// cukup SATU tabel - M_RATE_LIFE_SUMMARY". Ringkasan dibaca DAN ditulis di kolom `M_RATE_LIFE_SUMMARY` (ID, USEDBY,
// TYPE, MODIFIEDDATE, OPERATORID - migrasi inti 927/928); JSONDATA ringkasan dan tabel flat `RATE_LIFE_SUMMARY` (926)
// dibuang 928.
//
// RALAT R7 (10-10-2026, MODUL.md) - keputusan work owner "view rate_life sudah tidak digunakan, sekarang gunakan tabel
// M_RATE_LIFE": view `RATE_LIFE` dibuang di DEV (ORA-00942) dan `M_RATE_LIFE` kini tabel FLAT tanpa JSONDATA (ID, IDUSEDBY,
// USEDBY, TYPE, GENDER, CONTRACT, AGE, RATE - ALL_TAB_COLUMNS DEV dikirim work owner). Rincian dibaca DAN ditulis di
// kolom `M_RATE_LIFE`; nol JSON. Tanpa migrasi (tabel diubah DBA).
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - `M_RATE_LIFE_SUMMARY`, kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`; kolom = KolomRingkasan
	// (RALAT R6). Satu-satunya tabel ringkasan: ID terpakai pun diperiksa di sini saja.
	TabelRingkasan = "M_RATE_LIFE_SUMMARY"
	// TabelRate - kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`AddToListSummary_Act` b1833); tabel flat KolomTabelRate (RALAT R7).
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

// KolomTabelRate - kolom `M_RATE_LIFE` (ALL_TAB_COLUMNS DEV 10-10-2026, RALAT R7) - sama dengan kolom view `RATE_LIFE` lama.
var KolomTabelRate = []string{"ID", "IDUSEDBY", "USEDBY", "TYPE", "GENDER", "CONTRACT", "AGE", "RATE"}

// LebarKolomRate - lebar (byte) kolom `M_RATE_LIFE` (ALL_TAB_COLUMNS DEV 10-10-2026); ID NOT NULL.
var LebarKolomRate = map[string]int{"ID": 10, "IDUSEDBY": 10, "USEDBY": 500, "TYPE": 100, "GENDER": 10, "CONTRACT": 10, "AGE": 10, "RATE": 50}

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

// Kolom `M_RATE_LIFE` yang ditulis (RALAT R7; dulu kunci JSONDATA bernama sama). `TYPE` tidak ditulis (NULL di seluruh
// baris DEV). `USEDBY` = salinan nama ringkasan.
const (
	KolomRateIDUsedBy = "IDUSEDBY"
	KolomRateUsedBy   = "USEDBY"
	KolomRateGender   = "GENDER"
	KolomRateContract = "CONTRACT"
	KolomRateAge      = "AGE"
	KolomRateRate     = "RATE"
)

// Kolom yang dibaca (nol `SELECT *`).
const (
	kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`
	kolomRate      = `ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE`
)

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRate}

// DaftarDibacaSaja - objek yang hanya di-SELECT modul ini: KOSONG sejak RALAT R7 (view `RATE_LIFE` dibuang; rincian dibaca
// dan ditulis di `M_RATE_LIFE`).
var DaftarDibacaSaja = []string{}
