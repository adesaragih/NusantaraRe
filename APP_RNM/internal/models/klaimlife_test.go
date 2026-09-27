package models

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Tiket 01 AC-3: status ditampilkan sebagai kata (Outstanding / Aksep / Ditolak),
// bukan nama field STS_REJECT dan bukan angka. (spec.md US-26)
//
// Jebakan yang diuji secara khusus: nilai "1" berarti DIAKSEP meskipun kolomnya
// bernama STS_REJECT (spec.md bab Problem butir 3 - "Nama menyesatkan").
func TestStatusBarisDariKode(t *testing.T) {
	kasus := []struct {
		kode      string
		mau       StatusBaris
		mauKata   string
		diketahui bool
	}{
		{"0", StatusOutstanding, "Outstanding", true},
		{"1", StatusAksep, "Aksep", true},
		{"2", StatusDitolak, "Ditolak", true},
	}
	for _, k := range kasus {
		got := StatusBarisDariKode(k.kode)
		if got != k.mau {
			t.Errorf("kode %q: dapat %v, mau %v", k.kode, got, k.mau)
		}
		if got.String() != k.mauKata {
			t.Errorf("kode %q: kata %q, mau %q", k.kode, got.String(), k.mauKata)
		}
		if got.Diketahui() != k.diketahui {
			t.Errorf("kode %q: Diketahui()=%v", k.kode, got.Diketahui())
		}
	}
}

// Kolom nullable (ADR-U-0027) dan kode di luar ketiga nilai yang tertulis di spec
// TIDAK ditebak artinya. Spec hanya menetapkan 0, 1, dan 2.
func TestStatusBarisTidakMenebak(t *testing.T) {
	for _, kode := range []string{"", "3", "9", "X", " "} {
		got := StatusBarisDariKode(kode)
		if got != StatusTidakDiketahui {
			t.Errorf("kode %q: dapat %v, mau StatusTidakDiketahui", kode, got)
		}
		if got.Diketahui() {
			t.Errorf("kode %q: Diketahui() harus false", kode)
		}
	}
}

// ADR-U-0022: kode dan penanda tetap TEKS. "006" yang kembali sebagai "6" lolos
// pulang-pergi dan memecahkan penggolong - karena itu perbandingan kode tidak
// pernah lewat bilangan.
func TestKodeStatusTidakPernahJadiBilangan(t *testing.T) {
	if StatusBarisDariKode("00") != StatusTidakDiketahui {
		t.Error(`"00" tidak boleh dianggap sama dengan "0"`)
	}
	if StatusBarisDariKode("01") != StatusTidakDiketahui {
		t.Error(`"01" tidak boleh dianggap sama dengan "1"`)
	}
	if StatusBarisDariKode("006") != StatusTidakDiketahui {
		t.Error(`"006" tidak boleh dianggap sama dengan "6"`)
	}
}

// Kode mentah tetap dibawa apa adanya, supaya nilai yang tidak dikenal dapat
// dilaporkan tanpa ditebak.
func TestBarisMembawaKodeMentah(t *testing.T) {
	b := BarisAdjustment{KodeStatus: "7"}
	if b.Status() != StatusTidakDiketahui {
		t.Fatalf("status = %v", b.Status())
	}
	if b.KodeStatus != "7" {
		t.Errorf("kode mentah hilang: %q", b.KodeStatus)
	}
}

// Tiket 01 AC-4: tidak ada nilai uang sebagai binary floating point di lapisan
// mana pun maupun di JSON respons. (spec.md AC-22)
func TestUangDiJSONAdalahTeks(t *testing.T) {
	uang, err := NewMoney("1234567890.12345678", "IDR")
	if err != nil {
		t.Fatal(err)
	}
	b := BarisAdjustment{ID: "A1", JumlahKlaim: uang, KodeStatus: "1"}
	keluar, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	s := string(keluar)
	if !strings.Contains(s, `"amount":"1234567890.12345678"`) {
		t.Errorf("jumlah uang bukan teks desimal: %s", s)
	}
	if strings.Contains(s, `"amount":1234567890`) {
		t.Errorf("jumlah uang keluar sebagai angka JSON: %s", s)
	}
	if !strings.Contains(s, `"status":"Aksep"`) {
		t.Errorf("status tidak keluar sebagai kata: %s", s)
	}
	if strings.Contains(s, "STS_REJECT") {
		t.Errorf("nama field warisan bocor ke kontrak API: %s", s)
	}
}

// Presisi tidak boleh hilang saat pulang-pergi lewat JSON. Delapan angka di
// belakang koma sejajar dengan kolom uang (ADR-U-0003, ADR-U-0016).
func TestUangPulangPergiTanpaKehilanganPresisi(t *testing.T) {
	asal := `{"amount":"0.00000001","currency":"IDR"}`
	var m Money
	if err := json.Unmarshal([]byte(asal), &m); err != nil {
		t.Fatal(err)
	}
	balik, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if string(balik) != asal {
		t.Errorf("pulang-pergi berubah: %s -> %s", asal, balik)
	}
}

// Klaim membawa peserta, dan peserta membawa barisnya sendiri. Meratakan baris
// ke header menghapus informasi peserta pemilik dan mematahkan mesin status
// (spec.md, ADR-U-0011).
func TestBarisMelekatPadaPesertanya(t *testing.T) {
	k := Klaim{
		ID: "K1",
		Peserta: []Peserta{
			{ID: "P1", Baris: []BarisAdjustment{{ID: "A1"}, {ID: "A2"}}},
			{ID: "P2", Baris: []BarisAdjustment{{ID: "A3"}}},
		},
	}
	if n := k.CacahBaris(); n != 3 {
		t.Errorf("CacahBaris() = %d, mau 3", n)
	}
	if len(k.Peserta[0].Baris) != 2 {
		t.Errorf("peserta pertama kehilangan barisnya")
	}
	// Tiket 01 AC-1: seluruh baris, bukan hanya yang terakhir.
	if k.Peserta[0].Baris[0].ID != "A1" {
		t.Errorf("baris pertama hilang - hanya baris terakhir yang tersimpan")
	}
}

// ClaimRetro ikut ke JSON, sebagai TEKS.
//
// ⛔ Tinjauan ronde 6 menemukan medan ini dibaca dari Oracle lalu dibuang saat
// serialisasi: MarshalJSON mendaftar medannya satu per satu, dan yang baru
// terlewat. Separuh butir w2 mati tanpa satu test pun gagal.
func TestKlaimMembawaClaimRetroKeJSON(t *testing.T) {
	uang, err := NewMoney("7777.7777", "IDR")
	if err != nil {
		t.Fatal(err)
	}
	k := Klaim{ID: "UJI-1", ClaimRetro: uang}
	b, err := json.Marshal(k)
	if err != nil {
		t.Fatal(err)
	}
	teks := string(b)
	if !strings.Contains(teks, `"claimRetro"`) {
		t.Fatalf("JSON tidak memuat claimRetro: %s", teks)
	}
	// ⛔ Angka JSON dilarang: setiap pustaka membacanya sebagai float64
	// (ADR-U-0003). Jumlahnya harus muncul di dalam tanda kutip.
	if !strings.Contains(teks, `"7777.7777"`) {
		t.Errorf("jumlah claimRetro bukan teks berkutip: %s", teks)
	}
	if strings.Contains(teks, `:7777.7777`) {
		t.Errorf("jumlah claimRetro ditulis sebagai angka JSON: %s", teks)
	}
}

// TanggalKejadian harus menyeberang ke layar.
//
// ⛔ Tombol `Edit Date` (`ClaimLifeDetailGCNM.xml` b14115) menampilkan tanggal
// yang SEDANG berlaku. Sebelum kelompok Detail & Tutup, medannya tidak ikut
// di JSON, sehingga kotaknya selalu terbuka kosong - dan kosong terbaca
// "belum diisi" padahal mungkin sudah. Orang lalu mengetik ulang tanggal yang
// sudah benar, dan setiap pengetikan ulang melewati ValidasiDOL lagi.
func TestPesertaJSONMembawaTanggalKejadian(t *testing.T) {
	b, err := json.Marshal(Peserta{ID: "P-1", TanggalKejadian: "2026-03-01"})
	if err != nil {
		t.Fatal(err)
	}
	var isi map[string]any
	if err := json.Unmarshal(b, &isi); err != nil {
		t.Fatal(err)
	}
	if isi["tanggalKejadian"] != "2026-03-01" {
		t.Errorf("tanggalKejadian = %v, mau teks apa adanya", isi["tanggalKejadian"])
	}
	// ⛔ HIMPUNAN KUNCINYA DIKUNCI UTUH, bukan daftar kata terlarang.
	//
	// Ronde pertama menyenaraikan tiga nama harfiah - "namaTertanggung",
	// "nameOfInsured", "pemegangPolis" - dan itu penjaga yang menenangkan
	// tanpa menjaga: medan bernama "namaPeserta" lolos utuh. Daftar larangan
	// selalu kalah dari nama yang belum terpikirkan.
	//
	// Yang dikunci karena itu kebalikannya: SETIAP kunci yang boleh ada.
	// Menambah medan apa pun ke marshaller ini menyalakan uji ini, dan orang
	// yang menambahnya harus menuliskan namanya di sini - yaitu harus
	// melihatnya.
	kunci := make([]string, 0, len(isi))
	for k := range isi {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	// "total" masuk 27-09-2026 bersama keputusan bc (enam total peserta,
	// dihitung saat dibaca). Ia melewati penjaga ini lebih dulu, dan itu
	// memang tugasnya: medan baru pada marshaller harus DILIHAT orang,
	// bukan menyelinap. ⛔ Nama orang tetap TIDAK ADA di sini.
	mau := []string{"baris", "id", "isCheck", "mataUang", "nomorPolis",
		"nomorPremiList", "nomorSertifikat", "tanggalKejadian", "total"}
	if !reflect.DeepEqual(kunci, mau) {
		t.Errorf("kunci JSON peserta = %v, mau %v; medan baru harus disengaja",
			kunci, mau)
	}
}
