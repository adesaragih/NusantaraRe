package models_test

// Penjaga statik tiket 06 AC 5: `Type` klaim punya SATU rumah tersimpan.
//
// Pemilik: tiket 06.
//
// Dibaca sesudah: pohonklaim.go.
//
// Di Pega ia ada dua salinan - `LoadDataPeserta_Act` menyalin
// `pyWorkPage.PolicyDataLife.Type` ke `pyWorkPage.Type`, dan sejak itu dua
// halaman memegang nilai yang sama. Dua salinan berarti dua kesempatan untuk
// berbeda, dan yang membaca salinan basi tidak akan pernah tahu.
//
// ⚠️ Penjaga ini menelusuri SELURUH Claim Life (dulu `internal/`), bukan hanya paket ini - dan sejak 10-10-2026
// berkas modul lain HANYA bila kodenya menyentuh tabel Claim Life (pindai_bantu_test.go): `Type` yang dijaga adalah
// `Type` KLAIM; medan `Type` bordereaux, nbfacin, atau polis PremiumList bukan salinannya.
// Percobaan pertama hanya memindai `internal/models` sambil mengaku memeriksa
// seluruhnya - pengakuan yang lebih luas daripada yang diperiksanya, dan itu
// jenis penjaga yang paling berbahaya: ia menenangkan tanpa menjaga.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaMedanType mencocokkan deklarasi medan struct yang namanya PERSIS Type.
// Nama seperti TreatyTypeID atau TypeName tidak ikut.
var polaMedanType = regexp.MustCompile(`(?m)^\t+Type\s+string`)

// medanTypeYangSah adalah daftar tempat medan `Type string` boleh ada, masing
// masing dengan sebabnya. Menambah baris ke sini adalah tindakan SADAR, dan
// itulah gunanya daftar ini - bukan untuk melonggarkan, melainkan untuk
// membuat setiap salinan baru harus dipertanggungjawabkan.
var medanTypeYangSah = map[string]string{
	// Satu-satunya rumah TERSIMPAN: kolom T_GENERAL_CLAIM.TYPE - pindah dari
	// T_WORK_CLAIM.TYPE di migrasi 023 (keputusan work owner 01-10-2026).
	"models/klaimlife.go": "Klaim - satu-satunya salinan yang tersimpan",
	// ⛔ 10-10-2026: muatan permintaan Register (`PermintaanDaftar.Type`, badan JSON `type`) DICABUT - Type case kini
	// dibaca dari polis PL-nya (services/kuncipolis.go), tidak lagi dikirim klien.
	// ⛔ `MuatanKomite.Type` (A2, butir af) DIBUANG di migrasi 023: baris
	// kerja Komite tidak lagi menyimpan TYPE - salinannya dulu tidak pernah
	// dibaca; wewenang Komite membaca `TypeKlaim` klaim induk.
	// ⛔ Lingkup 10-10-2026: entri PremiumList Life (`PenawaranPolis`, `BarisInboxPolis`, `IsianDataPolis`,
	// `JawabanDataPolis`, `isiDataPolis`) dan `inti/backend/kontrak/polis.go` (`PolisRingkas`) DICABUT dari daftar ini:
	// semuanya `Type` POLIS (`T_PREMIUM_LIST.TYPE`), berkasnya tidak menyentuh tabel Claim Life, sehingga kini di luar
	// pindaian - bukan "diizinkan", melainkan bukan urusan penjaga Type KLAIM. Bila salah satunya kelak menyentuh
	// tabel Claim Life, penjaga ini menagihnya lagi.
}

func TestTypeKlaimHanyaSatuRumahTersimpan(t *testing.T) {
	akar := akarAplikasiPindai
	var temuan []string
	err := filepath.Walk(akar, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && lewatiFolderPindai(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") ||
			strings.HasSuffix(jalur, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		if !polaMedanType.MatchString(string(isi)) {
			return nil
		}
		if !dalamLingkupClaimLife(jalur, polaKomentar.ReplaceAllString(string(isi), "")) {
			return nil
		}
		rel := relLapisan(jalur)
		if !milikClaimLife(jalur) {
			rel = relAkar(jalur) // modul lain: jalur penuh, tidak bertabrakan dengan nama lapisan Claim Life
		}
		temuan = append(temuan, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(temuan) == 0 {
		t.Fatal("nol medan Type ditemukan; pembacanya yang rusak, bukan kodenya")
	}
	for _, jalur := range temuan {
		if _, sah := medanTypeYangSah[jalur]; !sah {
			t.Errorf("medan `Type string` baru di %s. Bila ia muatan permintaan, "+
				"daftarkan di medanTypeYangSah beserta sebabnya; bila ia salinan "+
				"TERSIMPAN kedua, ia melanggar tiket 06 AC 5 - `Type` dibaca dari "+
				"satu tempat, bukan dua", jalur)
		}
	}
	// Rumah tersimpannya harus benar-benar ada, bukan sekadar tidak dilanggar.
	var adaRumah bool
	for _, jalur := range temuan {
		if jalur == "models/klaimlife.go" {
			adaRumah = true
		}
	}
	if !adaRumah {
		t.Error("medan Type di Klaim hilang; penjaga ini kehilangan yang dijaganya")
	}
	// Daftar sah tidak boleh memuat entri basi: setiap baris harus menunjuk medan yang benar-benar ditemukan.
	for jalur := range medanTypeYangSah {
		ada := false
		for _, x := range temuan {
			ada = ada || x == jalur
		}
		if !ada {
			t.Errorf("medanTypeYangSah memuat %s, tetapi medan Type di sana tidak ditemukan dalam lingkup - cabut barisnya", jalur)
		}
	}
}
