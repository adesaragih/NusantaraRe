// Package repository adalah SATU-SATUNYA lapisan modul R/I Rate Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services.
//
// Untuk apa berkas ini: nama objek Oracle dan kunci JSON modul ini, ditulis SEKALI.
//
// RALAT R4 (06-10-2026, MODUL.md) - menggantikan K1 05-10-2026 ("CRUD menulis M_RATE_LIFE_SUMMARY, view TIDAK
// di-DROP"): keputusan work owner K-F1/K-F2 06-10-2026 - VIEW `RATE_LIFE_SUMMARY` diganti TABEL FLAT bernama sama
// (migrasi inti 926, enam kolom view). Ringkasan kini dibaca DAN ditulis di tabel flat itu (kolom bernama);
// `M_RATE_LIFE_SUMMARY` (JSON) tidak pernah ditulis lagi - cadangan, sumber alat pindahflat, dan pemeriksa ID terpakai.
// Rincian rate (`M_RATE_LIFE` + view `RATE_LIFE`) TIDAK berubah.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - TABEL FLAT `RATE_LIFE_SUMMARY` (migrasi inti 926, K-F1), kelas
	// `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`. Kolom = KolomViewRingkasan.
	TabelRingkasan = "RATE_LIFE_SUMMARY"
	// TabelRate - kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`AddToListSummary_Act` b1833); `ID VARCHAR2(10)`,
	// `JSONDATA CLOB` (katalog DEV).
	TabelRate = "M_RATE_LIFE"
	// KolomJSON - kolom CLOB berisi JSON di tabel warisan.
	KolomJSON = "JSONDATA"
)

// TabelRingkasanJSON - `M_RATE_LIFE_SUMMARY` (JSON warisan, 339 baris DEV): DIBACA saja (K-F1) - sumber alat
// pindahflat dan pemeriksa ID terpakai (MaksID/AdaID), supaya ID baru tidak bertabrakan dengan ringkasan Pega yang
// belum / akan dipindah.
const TabelRingkasanJSON = "M_RATE_LIFE_SUMMARY"

// KolomViewRingkasan - kolom VIEW warisan `RATE_LIFE_SUMMARY` menurut definisinya (ALL_VIEWS, owner POOLDATA, dibaca WO
// 06-10-2026): `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID,
// a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`. Tabel flat 926 memuat PERSIS kolom ini (K-F2); selain ID = kunci
// `M_RATE_LIFE_SUMMARY.JSONDATA` yang dibaca alat pindah.
var KolomViewRingkasan = []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID", "FLAG"}

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

// Kolom tabel flat ringkasan yang ditulis modul (XML `InboxSummaryRIRate`). `TYPE` dan `FLAG` TIDAK ditulis (tidak ada
// di XML, K-F2): baris baru NULL, baris pindahan apa adanya.
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

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE/LOCK.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRate}

// DaftarDibacaSaja - hanya SELECT: view rincian dan JSON ringkasan warisan.
var DaftarDibacaSaja = []string{ViewRate, TabelRingkasanJSON}
