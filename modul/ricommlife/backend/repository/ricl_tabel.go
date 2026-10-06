// Package repository adalah SATU-SATUNYA lapisan modul R/I Comm Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services
// (keputusan work owner 06-10-2026 butir 5: satu transaksi per simpan, prosedur PEGA_* tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle, kolom, dan kunci JSON modul ini, ditulis SEKALI.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - JSON warisan Pega, kelas `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY` (1 baris DEV, ID 1000003). TIDAK
	// diubah strukturnya (butir 6): sisip `JSON_OBJECT`, ubah baca-ubah-tulis di Go (ricl_json.go), nol JSON_MERGEPATCH.
	TabelRingkasan = "M_RICOMM_LIFE_SUMMARY"
	// TabelKomisi - tabel FLAT rincian (migrasi inti 924), menggantikan VIEW warisan bernama sama (butir 1).
	TabelKomisi = "RICOMM_LIFE"
	// KolomJSON - kolom CLOB berisi JSON di tabel warisan.
	KolomJSON = "JSONDATA"
)

// Objek yang DIBACA saja.
const (
	// ViewRingkasan - `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM
	// POOLDATA.M_RICOMM_LIFE_SUMMARY a` (dicek work owner 06-10-2026).
	ViewRingkasan = "RICOMM_LIFE_SUMMARY"
	// TabelJSONLama - `M_RICOMM_LIFE` (JSON, 0 baris DEV), sumber alat pindahflat; TIDAK pernah ditulis (butir 3).
	TabelJSONLama = "M_RICOMM_LIFE"
	// TabelSitus - `M_SITE_DATABASE` (ID NUMBER; baris (1,'1') dan (2,'0') DEV): awalan ID baru (butir 4).
	TabelSitus = "M_SITE_DATABASE"
)

// situsAktif - nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai (prosedur PEGA_M_RICOMM_LIFE b13); diikat.
const situsAktif = "1"

// Sequence warisan (butir 4: nol sequence baru). Nomornya dibentuk models.BentukID.
const (
	// SeqRingkasan - `M_RICOMM_LIFE_SUMMARY_SEQ` (last 4 DEV), prosedur PEGA_M_RICOMM_LIFE_SUMMARY.
	SeqRingkasan = "M_RICOMM_LIFE_SUMMARY_SEQ"
	// SeqKomisi - `M_RICOMM_LIFE_SEQ` (last 43 DEV), prosedur PEGA_M_RICOMM_LIFE b21.
	SeqKomisi = "M_RICOMM_LIFE_SEQ"
)

// KolomViewRingkasan - kolom view `RICOMM_LIFE_SUMMARY` menurut definisinya; selain ID = kunci JSONDATA ringkasan.
var KolomViewRingkasan = []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}

// KolomViewKomisiLama - kolom VIEW warisan `RICOMM_LIFE` (sebelum 924) = kunci `M_RICOMM_LIFE.JSONDATA` + ID. Tabel flat
// memuat PERSIS kolom ini (butir 1: "kolom hanya yang ada di view/XML").
var KolomViewKomisiLama = []string{"ID", "IDUSEDBY", "USEDBY", "CONTRACT", "YEAR", "COMM"}

// Kunci JSON ringkasan yang ditulis (butir 6) - `pxObjClass` seperti satu-satunya baris DEV.
const (
	JSONUsedBy     = "USEDBY"
	JSONOperatorID = "OPERATORID"
	JSONModified   = "MODIFIEDDATE"
	JSONKelas      = "pxObjClass"
)

// Kunci JSON `M_RICOMM_LIFE.JSONDATA` (dibaca alat pindah saja).
const (
	JSONIDUsedBy = "IDUSEDBY"
	JSONContract = "CONTRACT"
	JSONYear     = "YEAR"
	JSONComm     = "COMM"
)

// fmtAngka - kolom NUMBER sebagai teks tanpa bergantung NLS sesi (`TM9` menulis `.5`; dirapikan angkaBaca).
const fmtAngka = `TO_CHAR(%s, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// Kolom yang dibaca (nol `SELECT *`).
const kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE/LOCK.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelKomisi}

// DaftarDibacaSaja - hanya SELECT.
var DaftarDibacaSaja = []string{ViewRingkasan, TabelJSONLama, TabelSitus}
