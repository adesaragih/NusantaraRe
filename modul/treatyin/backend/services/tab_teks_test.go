package services_test

// Uji pemilihan ejaan tab teks - NOL koneksi Oracle.
//
// ⛔ Yang diuji di sini adalah aturan yang paling mudah salah TANPA
// KETAHUAN: memilih ejaan yang keliru menghasilkan layar yang terisi rapi
// dengan teks milik medan lain. Tidak ada galat, tidak ada sel kosong,
// tidak ada petunjuk apa pun bagi pembacanya.
//
// Angka dalam uji ini dari sapuan 4 Oktober 2026 atas seluruh 1.854
// dokumen - bukan dikarang.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func TestEjaanDipilihMenurutCabang(t *testing.T) {
	mentah := map[string]string{
		"ExclusionsP": "milik proporsional",
		"Exclusions":  "milik non-proporsional",
	}

	prop := services.PilihTabTeks(mentah, true,
		services.EjaanPengecualianProp, services.EjaanPengecualianNonProp)
	if prop.Isi != "milik proporsional" || prop.Ejaan != "ExclusionsP" {
		t.Errorf("cabang proporsional memakai %q dari %q", prop.Isi, prop.Ejaan)
	}

	non := services.PilihTabTeks(mentah, false,
		services.EjaanPengecualianProp, services.EjaanPengecualianNonProp)
	if non.Isi != "milik non-proporsional" || non.Ejaan != "Exclusions" {
		t.Errorf("cabang non-proporsional memakai %q dari %q", non.Isi, non.Ejaan)
	}
}

// ⛔ UJI PALING PENTING DI BERKAS INI.
//
// Ejaan cabang TIDAK ADA, ejaan lain ADA. Hasilnya WAJIB kosong. Terukur:
// dari 303 dokumen yang punya lebih dari satu ejaan `SpecialConditions*`,
// NOL yang isinya identik - jadi "jatuh ke ejaan satunya" menampilkan teks
// yang SALAH, bukan salinan yang basi.
//
// Keadaan ini nyata: 2 kontrak proporsional dan 2 non-proporsional untuk
// Exclusions, 20 dan 33 untuk Special Conditions.
func TestEjaanCabangTidakAdaMenghasilkanKOSONGBukanEjaanLain(t *testing.T) {
	mentah := map[string]string{"Exclusions": "TEKS MILIK CABANG SEBERANG"}

	prop := services.PilihTabTeks(mentah, true,
		services.EjaanPengecualianProp, services.EjaanPengecualianNonProp)
	if prop.Isi != "" {
		t.Errorf("isi %q diambil dari ejaan cabang seberang; ia teks yang BERBEDA, "+
			"dan pembacanya tidak punya cara tahu", prop.Isi)
	}
	if prop.Ejaan != "" {
		t.Errorf("Ejaan %q disebut padahal nol teks dipilih", prop.Ejaan)
	}
	// ⛔ Tetapi keberadaannya DISEBUT - itu bedanya dengan menyembunyikan.
	if len(prop.EjaanLain) != 1 || prop.EjaanLain[0] != "Exclusions" {
		t.Errorf("EjaanLain = %v, mau [Exclusions]", prop.EjaanLain)
	}
}

// Dua ejaan sama-sama berisi: yang sesuai cabang dipakai, yang lain
// DISEBUT. 175 kontrak proporsional dan 126 non-proporsional ada di
// keadaan ini untuk Exclusions saja.
func TestDuaEjaanBerisiMakaYangLainDisebut(t *testing.T) {
	mentah := map[string]string{
		"ExclusionsP": "yang tampil",
		"Exclusions":  "yang TIDAK tampil",
	}
	h := services.PilihTabTeks(mentah, true,
		services.EjaanPengecualianProp, services.EjaanPengecualianNonProp)
	if h.Isi != "yang tampil" {
		t.Errorf("isi = %q", h.Isi)
	}
	if len(h.EjaanLain) != 1 || h.EjaanLain[0] != "Exclusions" {
		t.Errorf("EjaanLain = %v; pembacanya berhak tahu ada teks lain", h.EjaanLain)
	}
}

// ⚠️ Ejaan KETIGA `SpecialConditionsp` dipakai KEDUA cabang - 170
// proporsional, 122 non - jadi ia tidak pernah terpilih dan SELALU muncul
// sebagai ejaan lain.
func TestEjaanKetigaTidakPernahTerpilihTetapiSelaluDisebut(t *testing.T) {
	mentah := map[string]string{
		"SpecialConditionsP": "prop",
		"SpecialConditions":  "nonprop",
		"SpecialConditionsp": "ejaan ketiga",
	}
	for _, p := range []struct {
		prop bool
		mau  string
	}{{true, "SpecialConditionsP"}, {false, "SpecialConditions"}} {
		h := services.PilihTabTeks(mentah, p.prop,
			services.EjaanSyaratProp, services.EjaanSyaratNonProp, services.EjaanSyaratKetiga)
		if h.Ejaan != p.mau {
			t.Errorf("prop=%v memilih %q, mau %q", p.prop, h.Ejaan, p.mau)
		}
		if len(h.EjaanLain) != 2 {
			t.Errorf("prop=%v: EjaanLain = %v, mau dua", p.prop, h.EjaanLain)
		}
		ada := false
		for _, e := range h.EjaanLain {
			if e == services.EjaanSyaratKetiga {
				ada = true
			}
		}
		if !ada {
			t.Errorf("prop=%v: ejaan ketiga tidak disebut: %v", p.prop, h.EjaanLain)
		}
	}
}

// Nol ejaan: kosong, nol ejaan lain, nol galat. 163 kontrak proporsional
// dan 230 non-proporsional tidak punya satu pun ejaan Special Conditions.
func TestNolEjaanBukanGalat(t *testing.T) {
	h := services.PilihTabTeks(map[string]string{}, true,
		services.EjaanSyaratProp, services.EjaanSyaratNonProp, services.EjaanSyaratKetiga)
	if h.Isi != "" || h.Ejaan != "" {
		t.Errorf("%+v", h)
	}
	// Irisan KOSONG, bukan nil - pemanggil JSON tidak perlu membedakan
	// `null` dari `[]`.
	if h.EjaanLain == nil {
		t.Error("EjaanLain nil, mau irisan kosong")
	}
	if len(h.EjaanLain) != 0 {
		t.Errorf("EjaanLain = %v", h.EjaanLain)
	}
}

// ⛔ Cabang dibaca dari KOLOM `TREATY_IN.PROPORTIONTYPE`, bukan dari kunci
// dokumen. Terukur: kolomnya terisi pada seluruh 1.854 baris, kunci
// dokumennya TIDAK ADA pada tiga (1001854, 1001855, 1001856).
func TestCabangDibacaDariNilaiKolom(t *testing.T) {
	if !services.SifatProporsional("Proportional") {
		t.Error("Proportional tidak dikenali")
	}
	for _, lain := range []string{"NonProportional", "", "   ", "proportional", "PROPORTIONAL"} {
		if services.SifatProporsional(lain) {
			t.Errorf("%q dikenali sebagai proporsional; hanya ejaan persis yang sah", lain)
		}
	}
}

// Kunci yang ADA bernilai KOSONG tidak disebut sebagai ejaan lain - tidak
// ada teks yang pembacanya lewatkan.
//
// ⚠️ Terukur nol kali di korpus hari ini (kelima kunci tidak pernah
// bernilai kosong), dan justru karena itu ia diuji: perilaku yang tidak
// pernah terpicu adalah perilaku yang tidak pernah terbukti.
func TestKunciKosongTidakDisebutSebagaiEjaanLain(t *testing.T) {
	mentah := map[string]string{"ExclusionsP": "ada isinya", "Exclusions": ""}
	h := services.PilihTabTeks(mentah, true,
		services.EjaanPengecualianProp, services.EjaanPengecualianNonProp)
	if len(h.EjaanLain) != 0 {
		t.Errorf("EjaanLain = %v; ejaan bernilai kosong bukan teks yang terlewat", h.EjaanLain)
	}
}
