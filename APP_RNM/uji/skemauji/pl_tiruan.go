package skemauji

// Tiruan tabel WARISAN rekap PremiumList Life untuk uji bertag db - PL-09
// (GILIRAN-18): `M_LIFE_PREMIUM_SUMMARY` dan sequence PK-nya.
//
// Daftar kolom TIDAK ditulis ulang di sini: ia datang dari penulisnya
// (`repository.NamaKolomSummaryWarisan`, VERBATIM badan prosedur
// `PEGA_M_LIFE_PREMIUM_SUMMARY`), sehingga tiruan tidak mungkin tertinggal.
//
// ⚠️ `[belum terverifikasi]` TIPE kolomnya: katalog DEV menyebut 38 kolom,
// tipenya tidak tercatat di repo. Uang dibuat NUMBER - seperti kolom uang
// `M_LIFE_PREMIUM_DETAIL` di katalog - supaya jalur "teks desimal ->
// konversi Oracle" yang dipakai parameter `VARCHAR2` prosedur itu ikut
// teruji; teks VARCHAR2(255).

import (
	"fmt"
	"strings"

	plrepo "nusantarare/modul/premiumlistlife/backend/repository"
)

// namaTabelSummaryPolis adalah tabel warisan rekap polis yang ditiru.
const namaTabelSummaryPolis = "M_LIFE_PREMIUM_SUMMARY"

// namaSequenceSummaryPolis adalah sequence PK tabel itu.
const namaSequenceSummaryPolis = "M_LIFE_PREMIUM_SUMMARY_SEQ"

// ddlTiruanSummaryPolis membuat tiruan tabel dan sequence-nya.
func ddlTiruanSummaryPolis(skema string) []string {
	kolom := []string{"ID VARCHAR2(50)"}
	for _, k := range plrepo.NamaKolomSummaryWarisan() {
		tipe := "VARCHAR2(255)"
		if plrepo.KolomUangSummaryWarisan(k) {
			tipe = "NUMBER"
		}
		kolom = append(kolom, k+" "+tipe)
	}
	return []string{
		fmt.Sprintf("CREATE TABLE %s.%s (%s)", skema, namaTabelSummaryPolis, strings.Join(kolom, ", ")),
		fmt.Sprintf("CREATE SEQUENCE %s.%s", skema, namaSequenceSummaryPolis),
	}
}
