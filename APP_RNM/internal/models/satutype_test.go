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
// ⚠️ Penjaga ini menelusuri SELURUH `internal/`, bukan hanya paket ini.
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
	// Satu-satunya rumah TERSIMPAN: kolom T_WORK_CLAIM.TYPE.
	"models/pohonklaim.go": "WorkClaim - satu-satunya salinan yang tersimpan",
	// Muatan permintaan, hidup sepanjang satu permintaan HTTP lalu hilang.
	// Ia bukan salinan kedua melainkan jalan masuk menuju yang pertama.
	"services/pendaftaran.go": "PermintaanDaftar - muatan permintaan, tidak tersimpan",
	"handlers/register.go":    "badan JSON - muatan permintaan, tidak tersimpan",
}

func TestTypeKlaimHanyaSatuRumahTersimpan(t *testing.T) {
	akar := ".."
	var temuan []string
	err := filepath.Walk(akar, func(jalur string, info os.FileInfo, err error) error {
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
		rel := filepath.ToSlash(strings.TrimPrefix(filepath.ToSlash(jalur), "../"))
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
		if jalur == "models/pohonklaim.go" {
			adaRumah = true
		}
	}
	if !adaRumah {
		t.Error("medan Type di WorkClaim hilang; penjaga ini kehilangan yang dijaganya")
	}
}
