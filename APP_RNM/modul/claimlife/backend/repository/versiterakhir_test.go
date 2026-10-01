package repository

// Peserta VERSI TERAKHIR - keputusan work owner 01-10-2026 (OQ-CL-VERSI, brief
// `PROMPT-PERBAIKAN-CLAIM-LIFE-PESERTA-VERSI-TERAKHIR.md`): sertifikat yang versi polis terakhirnya
// `Delete` atau `Batal` lewat endorsement TIDAK tampil dan TIDAK dapat diklaim. Penyimpangan sadar dari
// `Claim Life/RDBList/GetPesertaClaim_sql1.xml` (nol saringan versi maupun `EDMSTATUS`).
//
// Inventori pembaca `M_LIFE_PREMIUM_DETAIL` di modul ini (paket 1):
//
//	pesertapolis.go  Cari             WAJIB - daftar `Find Insured`: menentukan sertifikat yang dapat dipilih
//	pesertapolis.go  AmbilUntukKlaim  WAJIB - baris yang disalin ke klaim saat pendaftaran
//	gandawarisan.go  DOBSumberKosong, StatusWarisanTerakhir, AdaWarisanSamaDOL, IsiTertanggungCermin
//	                                  TIDAK - berkunci `ID` baris sumber (`SOURCE_ID`) yang SUDAH dipilih
//	                                  AmbilUntukKlaim saat pendaftaran; membaca riwayat baris itu, bukan memutuskan
//	                                  dapat-tidaknya diklaim
//	kolompeserta.go, models/klaimlife.go, services/simpanrnm.go
//	                                  TIDAK - hanya komentar; nol SQL ke tabel itu
//
// Tabel kasus di bawah dipakai DUA kali: oleh acuan aturan Go di berkas ini (selalu jalan) dan oleh uji
// `db` atas SQL sungguhan (`versiterakhir_db_test.go`, MELEWATI tanpa ORACLE_DSN).

import (
	"strconv"
	"strings"
	"testing"
)

// versiTerakhirDibangun - saklar paket 1 -> paket 2: uji bentuk SQL di bawah sudah ditulis dan terbukti MERAH
// atas SQL lama; ia dilewati sampai aturannya dibangun (paket 2 menghapus saklar ini).
const versiTerakhirDibangun = false

// barisUjiVersi - satu baris peserta tiruan (`UJI-`, nol data orang).
type barisUjiVersi struct {
	ID, PLNumberEDM, EDMStatus, TglInput, Nama string // TglInput `YYYY-MM-DD HH24:MI:SS`, kosong = NULL
}

// kasusVersi - satu sertifikat di satu PL_NUMBER dan jawaban yang dituntut.
type kasusVersi struct {
	Kode, Nama, PL, Sertifikat string
	Baris                      []barisUjiVersi
	Tampil                     bool   // tampil di Cari (dengan CariNama) - AmbilUntukKlaim mengikuti PilihAcuan
	IDTerpilih                 string // ID baris yang dipakai bila tampil
	CariNama                   string // kotak nama Find Insured (kosong = tanpa saringan nama)
}

func kasusVersiTerakhir() []kasusVersi {
	return []kasusVersi{
		{Kode: "a", Nama: "NB saja", PL: "UJI-PLV-A", Sertifikat: "UJI-C1",
			Baris:  []barisUjiVersi{{ID: "101", EDMStatus: ""}},
			Tampil: true, IDTerpilih: "101"},
		{Kode: "b", Nama: "NB + /01 Old: tampil sekali, nilai dari /01", PL: "UJI-PLV-B", Sertifikat: "UJI-C1",
			Baris:  []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-B/01", EDMStatus: "Old"}},
			Tampil: true, IDTerpilih: "102"},
		{Kode: "c", Nama: "NB + /01 Delete: tidak tampil", PL: "UJI-PLV-C", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-C/01", EDMStatus: "Delete"}}},
		{Kode: "d", Nama: "NB + /01 Batal: tidak tampil", PL: "UJI-PLV-D", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-D/01", EDMStatus: "Batal"}}},
		{Kode: "e", Nama: "NB + /01 Delete + /02 New (masuk lagi): tampil dari /02", PL: "UJI-PLV-E", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-E/01", EDMStatus: "Delete"},
				{ID: "103", PLNumberEDM: "UJI-PLV-E/02", EDMStatus: "New"}},
			Tampil: true, IDTerpilih: "103"},
		{Kode: "f", Nama: `EDMSTATUS berspasi " Delete ": tidak tampil`, PL: "UJI-PLV-F", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-F/01", EDMStatus: " Delete "}}},
		{Kode: "g", Nama: "dua baris versi sama: TGL_INPUT terbaru menang (bukan ID terbesar)", PL: "UJI-PLV-G", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"},
				{ID: "102", PLNumberEDM: "UJI-PLV-G/01", EDMStatus: "Old", TglInput: "2026-02-01 08:00:00"},
				{ID: "103", PLNumberEDM: "UJI-PLV-G/01", EDMStatus: "Delete", TglInput: "2026-01-01 08:00:00"}},
			Tampil: true, IDTerpilih: "102"},
		{Kode: "h", Nama: "PL_NUMBER berakhiran angka, baris NB ber-PL_NUMBER_EDM = PL_NUMBER: versi NB tetap 0", PL: "UJI-PLV-2024", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101", PLNumberEDM: "UJI-PLV-2024"},
				{ID: "102", PLNumberEDM: "UJI-PLV-2024/01", EDMStatus: "Delete"}}},
		{Kode: "i", Nama: "seri versi dan TGL_INPUT kosong: ID terbesar secara ANGKA (100 > 99)", PL: "UJI-PLV-I", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "99", PLNumberEDM: "UJI-PLV-I/01", EDMStatus: "Delete"},
				{ID: "100", PLNumberEDM: "UJI-PLV-I/01", EDMStatus: "Old"}},
			Tampil: true, IDTerpilih: "100"},
		{Kode: "j", Nama: "saringan nama Find Insured atas versi TERAKHIR, bukan versi lama", PL: "UJI-PLV-J", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101", Nama: "UJI NAMA LAMA"},
				{ID: "102", PLNumberEDM: "UJI-PLV-J/01", EDMStatus: "Old", Nama: "UJI NAMA BARU"}},
			CariNama: "lama"},
	}
}

// versiAcuan - ACUAN Go aturan §1 butir 1 (bukan kode produksi): angka sesudah `<PL_NUMBER>/` di
// `PL_NUMBER_EDM`; selain itu 0 (new business). Diperketat dari `REGEXP_SUBSTR(PL_NUMBER_EDM, '[0-9]+$')`
// brief: angka di akhir PL_NUMBER sendiri tidak pernah terbaca sebagai versi (kasus h).
func versiAcuan(pl, edm string) int {
	e := strings.TrimSpace(edm)
	if !strings.HasPrefix(e, pl+"/") {
		return 0
	}
	n, err := strconv.Atoi(e[len(pl)+1:])
	if err != nil || strings.Trim(e[len(pl)+1:], "0123456789") != "" {
		return 0
	}
	return n
}

// lebihBaruAcuan - urutan §1: versi DESC, TGL_INPUT DESC (kosong terakhir), ID DESC secara angka.
func lebihBaruAcuan(pl string, a, b barisUjiVersi) bool {
	if va, vb := versiAcuan(pl, a.PLNumberEDM), versiAcuan(pl, b.PLNumberEDM); va != vb {
		return va > vb
	}
	if a.TglInput != b.TglInput {
		if a.TglInput == "" || b.TglInput == "" {
			return b.TglInput == ""
		}
		return a.TglInput > b.TglInput
	}
	if len(a.ID) != len(b.ID) {
		return len(a.ID) > len(b.ID)
	}
	return a.ID > b.ID
}

// pilihAcuan - baris versi terakhir satu kasus menurut acuan, dan apakah ia hidup (PesertaHidup PRODUKSI).
// Tanpa saringan Nama: itulah jawaban AmbilUntukKlaim.
func pilihAcuan(k kasusVersi) (barisUjiVersi, bool) {
	terpilih := k.Baris[0]
	for _, b := range k.Baris[1:] {
		if lebihBaruAcuan(k.PL, b, terpilih) {
			terpilih = b
		}
	}
	return terpilih, PesertaHidup(terpilih.EDMStatus)
}

// TestAturanVersiTerakhirAcuan - jawaban kasus (a)-(j) menurut acuan aturan; PesertaHidup PRODUKSI menilai
// baris terpilih. Tabel kasusnya juga yang dijalankan uji `db` atas SQL sungguhan.
func TestAturanVersiTerakhirAcuan(t *testing.T) {
	for _, k := range kasusVersiTerakhir() {
		t.Run(k.Kode, func(t *testing.T) {
			terpilih, tampil := pilihAcuan(k)
			if k.CariNama != "" && !strings.Contains(strings.ToUpper(terpilih.Nama), strings.ToUpper(k.CariNama)) {
				tampil = false
			}
			if tampil != k.Tampil || (tampil && terpilih.ID != k.IDTerpilih) {
				t.Errorf("%s: tampil %v ID %s, mau %v ID %s", k.Nama, tampil, terpilih.ID, k.Tampil, k.IDTerpilih)
			}
		})
	}
}

// rataSQL - SQL dirapikan satu spasi supaya pencocokan tidak bergantung indentasi.
func rataSQL(q string) string { return strings.Join(strings.Fields(q), " ") }

// Potongan SQL yang DITUNTUT aturan - teks harfiah, bukan konstanta produksi: mengubah konstanta produksi
// tidak boleh ikut mengubah tuntutannya.
const (
	tuntutPeringkat = "ROW_NUMBER() OVER (PARTITION BY CERTIFICATE_NO ORDER BY"
	tuntutVersi     = "NVL(TO_NUMBER(CASE WHEN SUBSTR(TRIM(PL_NUMBER_EDM), 1, LENGTH(PL_NUMBER) + 1) = PL_NUMBER || '/' " +
		"THEN REGEXP_SUBSTR(SUBSTR(TRIM(PL_NUMBER_EDM), LENGTH(PL_NUMBER) + 2), '^[0-9]+$') END), 0) DESC"
	tuntutSeri  = "TGL_INPUT DESC NULLS LAST, LENGTH(ID) DESC, ID DESC)"
	tuntutPilih = "RN_VERSI = 1 AND (EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))"
)

// TestSQLCariPesertaVersiTerakhir - bentuk SQL `Find Insured` per kasus: setiap kasus menuntut potongan yang
// memutuskannya. Menghapus ROW_NUMBER membuat (b), (c), (d), (e), (g), (i) merah.
func TestSQLCariPesertaVersiTerakhir(t *testing.T) {
	if !versiTerakhirDibangun {
		t.Skip("paket 1: terbukti merah atas SQL lama; aturan dibangun paket 2")
	}
	q := rataSQL(func() string { s, _ := sqlCariPeserta("S.M", "UJI-PLV", "", "", 10); return s }())
	qn := rataSQL(func() string { s, _ := sqlCariPeserta("S.M", "UJI-PLV", "UJI-C", "lama", 10); return s }())
	tuntut := map[string][]string{
		"a": {tuntutPilih},
		"b": {tuntutPeringkat, tuntutPilih},
		"c": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"d": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"e": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"f": {"TRIM(EDMSTATUS) NOT IN ('Batal','Delete')"},
		"g": {tuntutPeringkat, tuntutSeri},
		"h": {tuntutVersi},
		"i": {tuntutSeri},
	}
	for kode, potongan := range tuntut {
		for _, p := range potongan {
			if !strings.Contains(q, p) {
				t.Errorf("(%s) SQL tanpa %q:\n%s", kode, p, q)
			}
		}
	}
	// (j) saringan nama di LUAR jendela (atas versi terakhir); saringan sertifikat boleh di dalam (seluruh
	// kelompok satu sertifikat ikut atau tidak, peringkat di dalamnya tidak berubah).
	if i := strings.Index(qn, "RN_VERSI = 1"); i < 0 || !strings.Contains(qn[i:], "UPPER(NAME_OF_INSURED) LIKE") {
		t.Errorf("(j) saringan nama tidak berada sesudah peringkat versi:\n%s", qn)
	}
	// §1 butir 4: jendela hanya menyentuh SATU polis - PL_NUMBER = :1 di DALAM subkueri peringkat.
	if a, b := strings.Index(q, "FROM ("), strings.Index(q, "RN_VERSI = 1"); a < 0 || b < a ||
		!strings.Contains(q[a:b], "WHERE PL_NUMBER = :1") {
		t.Errorf("jendela peringkat tidak dikurung PL_NUMBER (pemindaian penuh):\n%s", q)
	}
}
