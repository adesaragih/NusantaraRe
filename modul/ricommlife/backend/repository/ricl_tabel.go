// Package repository adalah SATU-SATUNYA lapisan modul R/I Comm Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services
// (keputusan work owner 06-10-2026 butir 5: satu transaksi per simpan, prosedur PEGA_* tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// RALAT R1 (keputusan work owner 08-10-2026, MODUL.md): SATU tabel per jenis data, pola R/I Rate Life. Ringkasan =
// kolom `M_RICOMM_LIFE_SUMMARY` (migrasi inti 931/932; JSONDATA dan view `RICOMM_LIFE_SUMMARY` dibuang 932). Rincian =
// kolom `M_RICOMM_LIFE` (933/934; JSONDATA dibuang, tabel flat `RICOMM_LIFE` 924 dibuang 934). Nol JSON di modul ini.
package repository

// Tabel yang DITULIS.
const (
	// TabelRingkasan - `M_RICOMM_LIFE_SUMMARY`, kelas `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY`; kolom = KolomRingkasan.
	TabelRingkasan = "M_RICOMM_LIFE_SUMMARY"
	// TabelKomisi - `M_RICOMM_LIFE`, kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`; kolom = KolomKomisi.
	TabelKomisi = "M_RICOMM_LIFE"
)

// Objek yang DIBACA saja.
const (
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

// KolomRingkasan - kolom `M_RICOMM_LIFE_SUMMARY` sesudah 932 = kolom view `RICOMM_LIFE_SUMMARY` lama (ID warisan + 931).
var KolomRingkasan = []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}

// KolomKomisi - kolom `M_RICOMM_LIFE` sesudah 934 = kolom tabel `RICOMM_LIFE` 924 (ID warisan + 933).
var KolomKomisi = []string{"ID", "IDUSEDBY", "USEDBY", "CONTRACT", "YEAR", "COMM"}

// fmtAngka - kolom NUMBER sebagai teks tanpa bergantung NLS sesi (`TM9` menulis `.5`; diurai Go AngkaOracle).
const fmtAngka = `TO_CHAR(%s, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// Kolom yang dibaca (nol `SELECT *`).
const kolomRingkasan = `ID, USEDBY, OPERATORID, MODIFIEDDATE`

// DaftarTabelDitulis - penjaga modul: hanya dua tabel ini yang menerima INSERT/UPDATE/DELETE.
var DaftarTabelDitulis = []string{TabelRingkasan, TabelKomisi}

// DaftarDibacaSaja - hanya SELECT.
var DaftarDibacaSaja = []string{TabelSitus}
