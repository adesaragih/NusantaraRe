package services

// Untuk apa berkas ini: TEMPAT BERPERAN - tiket 05 (spec §5.4; AC 12, 13, 81,
// 82, 91). Di Pega beberapa bagian layar muncul atau wajib hanya untuk orang
// tertentu (`OperatorID.pyUserIdentifier == '<ID-operator-N>'`). Pemetaan
// nama -> peran TIDAK ada di korpus; sistem baru membaca pemetaan
// TEMPAT -> PERAN -> ARAH dari `M_NBTRIN_PERAN_TEMPAT` (migrasi 330), yang
// diisi IAM bersama work owner KEMUDIAN.
//
//	tempat tanpa baris  -> TERTUNDA: bagian itu tidak tampil, syaratnya tidak
//	                       dianggap terpenuhi (AC 81) - tidak ditebak
//	ARAH MUNCUL         -> tampil hanya bagi pemegang peran itu
//	ARAH KECUALI        -> tampil bagi semua KECUALI pemegang peran itu
//
// ⛔ Arah tidak pernah ditebak (AC 82): baris tanpa arah yang dikenal
// diabaikan, dan tempatnya tetap tertunda. ⛔ Nol peran karangan (AC 91):
// uji memakai peran fiktif berawalan UJI-.
//
// Tempat identitas di rule TERJANGKAU dan nasibnya:
//
//	ListSuggest .ProductionDate (tampil + wajib, <ID-operator-3>/<ID-operator-4>)
//	    -> TempatTanggalProduksi, lewat tabel ini
//	DetailDeptHeadTreatyIn_UW tiga tombol Submit (<ID-operator-1>)
//	    -> diganti POSISI kasus (models.TombolUntuk; AC 8, P13) - catatan tiket 05
//	DetailPoliciesNonProportional label "NON EDM" / "EDM" (<ID-operator-2>)
//	    -> bagian XOL non-proporsional tidak dibangun (JSON master, P29)
//	When IsSPVCreate (<ID-operator-5>/<ID-operator-6>), IsTreaty1, IsSPVTreaty1
//	    -> hanya memilih Assignment4 atau 6 (posisi sama); tidak dibangun

import (
	"context"

	inti "nusantarare/inti/backend"
)

// TempatTanggalProduksi - `Section/ListSuggest` medan `.ProductionDate`.
const TempatTanggalProduksi = "LISTSUGGEST_PRODUCTIONDATE"

// SemuaTempat - tempat yang dibaca layanan ini.
var SemuaTempat = []string{TempatTanggalProduksi}

// Arah pemeriksaan - VERBATIM nilai CHECK migrasi 330.
const (
	ArahMuncul  = "MUNCUL"
	ArahKecuali = "KECUALI"
)

// tempat menghitung tampil/tidaknya setiap tempat bagi pelaku.
func (l *Layanan) tempat(ctx context.Context, p inti.Pelaku) (map[string]bool, error) {
	baris, err := l.g.DaftarPeranTempat(ctx)
	if err != nil {
		return nil, err
	}
	out := map[string]bool{}
	for _, t := range SemuaTempat {
		out[t] = false
	}
	type keadaan struct{ muncul, kecuali, punyaMuncul, punyaKecuali bool }
	per := map[string]*keadaan{}
	for _, b := range baris {
		k := per[b.KodeTempat]
		if k == nil {
			k = &keadaan{}
			per[b.KodeTempat] = k
		}
		switch b.Arah {
		case ArahMuncul:
			k.muncul = true
			k.punyaMuncul = k.punyaMuncul || p.PunyaPeran(b.Peran)
		case ArahKecuali:
			k.kecuali = true
			k.punyaKecuali = k.punyaKecuali || p.PunyaPeran(b.Peran)
		}
	}
	for kode, k := range per {
		if _, dikenal := out[kode]; !dikenal {
			continue
		}
		switch {
		case k.muncul && !k.kecuali:
			out[kode] = k.punyaMuncul
		case k.kecuali && !k.muncul:
			out[kode] = !k.punyaKecuali
		}
		// kedua arah di satu tempat = bertentangan: tetap tertunda.
	}
	return out, nil
}
