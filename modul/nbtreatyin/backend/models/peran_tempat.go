package models

// Untuk apa berkas ini: TEMPAT BERPERAN - tiket 05 (spec §5.4; AC 12, 13, 81,
// 82, 91). Di Pega beberapa bagian layar muncul atau wajib hanya untuk orang
// tertentu (`OperatorID.pyUserIdentifier == '<ID-operator-N>'`). Pemetaan
// nama -> peran TIDAK ada di korpus; sistem baru memetakan TEMPAT -> PERAN ->
// ARAH, dan peran pengguna dibaca dari `inti.Pelaku.Peran` (workbasket akun).
//
// ⛔ KONSTANTA KODE, BUKAN TABEL. `[keputusan work owner]` K16 (03-10-2026,
// PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2): tabel `M_NBTRIN_PERAN_TEMPAT` (migrasi
// 330 putaran 1) tidak ada di diagram grilling dan DIHAPUS; pemetaannya hidup
// di `PemetaanPeranTempat` di bawah - KOSONG sampai IAM menjawab (K12). Tiket
// 05 tetap needs-info. Perilaku saat kosong = perilaku tabel kosong putaran 1:
//
//	tempat tanpa baris  -> TERTUNDA: bagian itu tidak tampil, syaratnya tidak
//	                       dianggap terpenuhi (AC 81) - tidak ditebak
//	ARAH MUNCUL         -> tampil hanya bagi pemegang peran itu
//	ARAH KECUALI        -> tampil bagi semua KECUALI pemegang peran itu
//	kedua arah / arah asing di satu tempat -> tetap tertunda (AC 82)
//
// ⛔ Nol peran karangan (AC 91): uji memakai peran fiktif berawalan UJI-.
//
// Tempat identitas di rule TERJANGKAU dan nasibnya (tiket 05):
//
//	ListSuggest .ProductionDate (tampil + wajib, <ID-operator-3>/<ID-operator-4>)
//	    -> TempatTanggalProduksi, lewat pemetaan ini
//	DetailDeptHeadTreatyIn_UW tiga tombol Submit (<ID-operator-1>)
//	    -> diganti POSISI kasus (TombolUntuk; AC 8, P13) - catatan tiket 05
//	DetailPoliciesNonProportional label "NON EDM" / "EDM" (<ID-operator-2>)
//	    -> milik jalur NonProp (paket XOL); tempatnya belum dipetakan
//	When IsSPVCreate (<ID-operator-5>/<ID-operator-6>), IsTreaty1, IsSPVTreaty1
//	    -> hanya memilih Assignment4 atau 6 (posisi sama); tidak dibangun

// PeranTempat adalah satu baris pemetaan tempat -> peran -> arah.
type PeranTempat struct {
	KodeTempat string
	Peran      string
	// Arah - ArahMuncul atau ArahKecuali; TIDAK ditebak (AC 82).
	Arah string
}

// TempatTanggalProduksi - `Section/ListSuggest` medan `.ProductionDate`.
const TempatTanggalProduksi = "LISTSUGGEST_PRODUCTIONDATE"

// SemuaTempat - tempat yang dibaca layanan.
var SemuaTempat = []string{TempatTanggalProduksi}

// Arah pemeriksaan.
const (
	ArahMuncul  = "MUNCUL"
	ArahKecuali = "KECUALI"
)

// PemetaanPeranTempat - pemetaan yang berlaku. ⛔ KOSONG sampai IAM bersama
// work owner menyerahkan peran dan arah ke-12 tempat (tiket 05, K12). Mengisi
// baris di sini adalah keputusan wewenang - bukan pekerjaan agen.
var PemetaanPeranTempat = []PeranTempat{}

// TempatTampil menghitung tampil/tidaknya setiap tempat `SemuaTempat` bagi
// pemegang peran `punyaPeran` (lazimnya `inti.Pelaku.PunyaPeran`).
func TempatTampil(pemetaan []PeranTempat, punyaPeran func(string) bool) map[string]bool {
	out := map[string]bool{}
	for _, t := range SemuaTempat {
		out[t] = false
	}
	type keadaan struct{ muncul, kecuali, asing, punyaMuncul, punyaKecuali bool }
	per := map[string]*keadaan{}
	for _, b := range pemetaan {
		k := per[b.KodeTempat]
		if k == nil {
			k = &keadaan{}
			per[b.KodeTempat] = k
		}
		switch b.Arah {
		case ArahMuncul:
			k.muncul = true
			k.punyaMuncul = k.punyaMuncul || punyaPeran(b.Peran)
		case ArahKecuali:
			k.kecuali = true
			k.punyaKecuali = k.punyaKecuali || punyaPeran(b.Peran)
		default:
			k.asing = true
		}
	}
	for kode, k := range per {
		if _, dikenal := out[kode]; !dikenal || k.asing {
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
	return out
}
