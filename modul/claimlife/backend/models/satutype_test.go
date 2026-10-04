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
// ⚠️ Penjaga ini menelusuri SELURUH aplikasi (dulu `internal/`), bukan hanya paket ini.
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
	// Muatan permintaan, hidup sepanjang satu permintaan HTTP lalu hilang.
	// Ia bukan salinan kedua melainkan jalan masuk menuju yang pertama.
	"services/pendaftaran.go": "PermintaanDaftar - muatan permintaan, tidak tersimpan",
	"handlers/register.go":    "badan JSON - muatan permintaan, tidak tersimpan",
	// ⛔ `MuatanKomite.Type` (A2, butir af) DIBUANG di migrasi 023: baris
	// kerja Komite tidak lagi menyimpan TYPE - salinannya dulu tidak pernah
	// dibaca; wewenang Komite membaca `TypeKlaim` klaim induk.
	// ⛔ MODUL LAIN, TABEL LAIN. `PenawaranPolis.Type` adalah `Type`
	// sebuah POLIS (`pyWorkPage.Type` di `InputPolicyHolder`), bukan `Type`
	// sebuah KLAIM. Rumah tersimpannya `T_PREMIUM_LIST.TYPE` (migrasi 051),
	// tabel yang berbeda dari `T_WORK_CLAIM`.
	//
	// Yang di berkas ini muatan gerbang `ProtectAccept`: ia hidup sepanjang
	// satu pemeriksaan lalu hilang, dan dibaca untuk memilih cabang `TP`
	// b4121 atau `TR` b4353. Bukan salinan kedua dari apa pun.
	//
	// ⚠️ Penjaga ini menagihnya dan itu BENAR: dua modul yang keduanya
	// punya medan bernama `Type` adalah persis keadaan yang membuat orang
	// membaca yang satu sebagai yang lain.
	"models/polis_validasi.go": "PenawaranPolis - Type POLIS (T_PREMIUM_LIST), muatan gerbang ProtectAccept",
	// ⛔ Kotak masuk PremiumList: `Type` salah satu dari tiga belas
	// kolom `InboxPremiumList.xml` (b836 `A.Type`). Ia DIBACA dari
	// `T_PREMIUM_LIST.TYPE` - rumah tersimpannya - lalu menyeberang ke
	// layar. Nol baris disimpan dari sini; keduanya jalur BACA.
	"repository/polis_inbox.go": "BarisInboxPolis - kolom baca InboxPremiumList b836",
	"services/polis_inbox.go":   "BarisInboxPolis - bentuk layar, jalur baca",
	// ⛔ Butir pl4/av. `PolisRingkas.Type` adalah `Type` POLIS yang DIBACA
	// dari rumah tersimpannya - `T_PREMIUM_LIST.TYPE`, migrasi 051 - lalu
	// diserahkan ke Claim Life sebagai `PolicyDataLife.Type`.
	//
	// ⛔ INI YANG MEMBUAT PENJAGA INI BERHARGA DI SINI: di Pega, layar
	// Register Claim Life menampilkan `.PolicyDataLife.Type` sebagai medan
	// `pyReadOnly` - DIBACA dari polis, tidak pernah diketik dan tidak
	// pernah disimpan ke tabel klaim. Menyalinnya ke `T_GENERAL_CLAIM.TYPE`
	// akan membuat dua `Type` untuk satu klaim, dan yang satu akan basi
	// begitu polisnya di-endorse.
	// Refactor bentuk B (30-09-2026): tipe `PolisRingkas` kini kontrak lintas
	// modul (dulu repository/polis_ringkas.go).
	"inti/backend/kontrak/polis.go": "PolisRingkas - Type POLIS dibaca dari T_PREMIUM_LIST untuk PolicyDataLife",
	"services/polis_ringkas.go":     "bentuk layar PolicyDataLife, jalur baca",
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
		rel := relLapisan(jalur)
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
}
