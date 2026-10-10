// Package repository adalah SATU-SATUNYA lapisan modul R/I Risk yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services
// (pola ricommlife: satu transaksi per simpan, prosedur PEGA_* tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 K1 (MODUL.md): SATU tabel per jenis data. Tabel Pega `M_RIRISK_LIFE_SUMMARY` /
// `M_RIRISK_LIFE` BERGANTI NAMA menjadi `RIRISK_LIFE_SUMMARY` / `RIRISK_LIFE` (nama view lama yang dibuang) dan menjadi
// tabel flat tanpa JSONDATA - migrasi inti 935-940. Nol JSON di modul ini.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - `RIRISK_LIFE_SUMMARY`, kelas `ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY`; kolom = KolomRingkasan.
	TabelRingkasan = "RIRISK_LIFE_SUMMARY"
	// TabelRincian - `RIRISK_LIFE`, kelas `ASM-FW-GISFW-Int-RI_RISK_LIFE`; kolom = KolomRincian.
	TabelRincian = "RIRISK_LIFE"
)

// Objek yang DIBACA saja.
const (
	// TabelSitus - `M_SITE_DATABASE` (ID NUMBER; baris (1,'1') dan (2,'0') DEV): awalan ID ringkasan baru.
	TabelSitus = "M_SITE_DATABASE"
)

// situsAktif - nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai (prosedur PEGA_M_RIRISK_LIFE_SUMMARY);
// diikat.
const situsAktif = "1"

// Sequence warisan (K3: TETAP, nol sequence baru). Nomornya dibentuk models.BentukID (ringkasan) /
// models.BentukIDRincian (rincian).
const (
	// SeqRingkasan - `M_RIRISK_LIFE_SUMMARY_SEQ` (last_number 166 DEV), prosedur PEGA_M_RIRISK_LIFE_SUMMARY.
	SeqRingkasan = "M_RIRISK_LIFE_SUMMARY_SEQ"
	// SeqRincian - `M_RIRISK_LIFE_SEQ` (last_number 31723 DEV), prosedur PEGA_M_RIRISK_LIFE.
	SeqRincian = "M_RIRISK_LIFE_SEQ"
)

// KolomRingkasan - kolom `RIRISK_LIFE_SUMMARY` sesudah 937 = kolom view lama bernama sama (ID warisan + 936).
var KolomRingkasan = []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}

// KolomRincian - kolom `RIRISK_LIFE` sesudah 940 = kolom view lama bernama sama, urutan sama (warisan + AGE 939).
var KolomRincian = []string{"ID", "IDUSEDBY", "USEDBY", "AGE", "YEAR", "MONTH", "RISK", "CONTRACT"}

// fmtAngka - kolom NUMBER sebagai teks tanpa bergantung NLS sesi (`TM9` menulis `.5`; diurai Go AngkaOracle).
const fmtAngka = `TO_CHAR(%s, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// Kolom yang dibaca (nol `SELECT *`).
const kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelRincian}

// DaftarDibacaSaja - hanya SELECT.
var DaftarDibacaSaja = []string{TabelSitus}
