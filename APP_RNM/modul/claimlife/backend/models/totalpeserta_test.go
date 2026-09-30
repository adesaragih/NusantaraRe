package models

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
)

// uang membangun nilai uang untuk uji; gagal keras bila teksnya salah.
func uangPeserta(t *testing.T, jumlah string) uang.Money {
	t.Helper()
	m, err := uang.NewMoney(jumlah, "IDR")
	if err != nil {
		t.Fatalf("NewMoney(%q): %v", jumlah, err)
	}
	return m
}

// barisUji membangun satu baris adjustment dengan ENAM nilai yang
// BERBEDA-BEDA.
//
// ⛔ Sengaja berbeda semua. Bila keenamnya bernilai sama, dua kolom yang
// tertukar di kolomTotal akan menghasilkan angka yang persis benar dan uji
// ini akan diam - padahal yang tertukar adalah dua ANGKA UANG.
func barisUji(t *testing.T, kelipatan string, status string) BarisAdjustment {
	t.Helper()
	k := func(dasar string) uang.Money { return uangPeserta(t, dasar+kelipatan) }
	return BarisAdjustment{
		KodeStatus:       status,
		CurrencyID:       "IDR",
		CedingRetention:  k("1"),
		ShareNusantaraRe: k("2"),
		SumInsured:       k("3"),
		SumReasured:      k("4"),
		ShareRetro:       k("5"),
		JumlahKlaim:      k("6"),
	}
}

// ambilEnam membaca keenam total sebagai teks, berurutan seperti kolomTotal.
func ambilEnam(tp TotalPeserta) []string {
	return []string{
		tp.CedingRetention.String(),
		tp.ShareNusantaraRe.String(),
		tp.SumInsured.String(),
		tp.SumReasured.String(),
		tp.ShareRetro.String(),
		tp.JumlahKlaim.String(),
	}
}

func TestHitungTotalPeserta(t *testing.T) {
	t.Run("contoh terhitung tangan, keenam kolom terpisah", func(t *testing.T) {
		// Dua baris. Nilainya dipilih supaya tiap kolom berjumlah BEDA, dan
		// jumlah yang diharapkan ditulis sebagai LITERAL - bukan dihitung
		// ulang dengan rumus yang sedang diuji.
		//
		//   baris 1: 10 20 30 40 50 60
		//   baris 2:  1  2  3  4  5  6
		//   jumlah : 11 22 33 44 55 66
		baris := []BarisAdjustment{
			barisUji(t, "0", kontrak.KodeAksep),
			barisUji(t, "", kontrak.KodeAksep),
		}
		total, err := HitungTotalPeserta(baris, "IDR")
		if err != nil {
			t.Fatalf("HitungTotalPeserta: %v", err)
		}
		mau := []string{"11 IDR", "22 IDR", "33 IDR", "44 IDR", "55 IDR", "66 IDR"}
		for i, got := range ambilEnam(total) {
			if got != mau[i] {
				t.Errorf("total %s = %q, mau %q", kolomTotal[i].nama, got, mau[i])
			}
		}
	})

	t.Run("baris DITOLAK ikut dijumlah", func(t *testing.T) {
		// ⛔ Ini yang XML lakukan, dan ia berlawanan dengan dugaan wajar.
		// Langkah 8.1 b4221 dan 23.1 b10841 sama-sama berprasyarat KOSONG:
		// tidak ada penyaringan atas STS_REJECT sama sekali. Menyaringnya
		// menghasilkan angka yang lebih masuk akal dan bukan angka Pega.
		diaksep := []BarisAdjustment{barisUji(t, "", kontrak.KodeAksep)}
		campur := []BarisAdjustment{
			barisUji(t, "", kontrak.KodeAksep),
			barisUji(t, "", kontrak.KodeDitolak),
		}
		satu, err := HitungTotalPeserta(diaksep, "IDR")
		if err != nil {
			t.Fatalf("satu baris: %v", err)
		}
		dua, err := HitungTotalPeserta(campur, "IDR")
		if err != nil {
			t.Fatalf("dua baris: %v", err)
		}
		if satu.JumlahKlaim.String() == dua.JumlahKlaim.String() {
			t.Errorf("baris ditolak TIDAK ikut dijumlah: %q = %q",
				satu.JumlahKlaim.String(), dua.JumlahKlaim.String())
		}
		if dua.JumlahKlaim.String() != "12 IDR" {
			t.Errorf("JumlahKlaim = %q, mau %q - baris ditolak harus ikut",
				dua.JumlahKlaim.String(), "12 IDR")
		}
	})

	t.Run("peserta tanpa baris: NOL, bukan kosong", func(t *testing.T) {
		// Penampungnya mulai dari literal 0 (b4027..b4159) dan langkah 8.2
		// menulisnya apa adanya, jadi Pega menuliskan 0. Mengembalikan
		// kosong akan membuat layar berkata "belum ada datanya" untuk
		// peserta yang datanya lengkap dan memang berjumlah nol.
		total, err := HitungTotalPeserta(nil, "IDR")
		if err != nil {
			t.Fatalf("HitungTotalPeserta(nil): %v", err)
		}
		for i, got := range ambilEnam(total) {
			if got != "0 IDR" {
				t.Errorf("total %s = %q, mau %q", kolomTotal[i].nama, got, "0 IDR")
			}
		}
		if total.JumlahKlaim.Kosong() {
			t.Error("total kosong; Pega menulis 0")
		}
	})

	t.Run("kolom kosong tidak menyumbang dan tidak menggagalkan", func(t *testing.T) {
		// Money.Add MENOLAK nilai kosong, jadi tanpa penjagaan ini satu
		// kolom kosong akan menggagalkan SELURUH pembacaan Detail - layar
		// putih karena satu sel yang memang boleh kosong.
		penuh := barisUji(t, "", kontrak.KodeAksep)
		bolong := barisUji(t, "", kontrak.KodeAksep)
		bolong.SumInsured = uang.Money{}
		total, err := HitungTotalPeserta([]BarisAdjustment{penuh, bolong}, "IDR")
		if err != nil {
			t.Fatalf("HitungTotalPeserta: %v", err)
		}
		if got := total.SumInsured.String(); got != "3 IDR" {
			t.Errorf("SumInsured = %q, mau %q (hanya baris yang terisi)", got, "3 IDR")
		}
		if got := total.SumReasured.String(); got != "8 IDR" {
			t.Errorf("SumReasured = %q, mau %q (tak terpengaruh)", got, "8 IDR")
		}
	})

	t.Run("mata uang beragam DITOLAK", func(t *testing.T) {
		a := barisUji(t, "", kontrak.KodeAksep)
		b := barisUji(t, "", kontrak.KodeAksep)
		b.CurrencyID = "USD"
		_, err := HitungTotalPeserta([]BarisAdjustment{a, b}, "IDR")
		if !errors.Is(err, ErrTotalMataUangBeragam) {
			t.Errorf("galat = %v, mau ErrTotalMataUangBeragam", err)
		}
	})

	t.Run("kolom bermata uang lain DITOLAK, bukan dilabeli ulang", func(t *testing.T) {
		// ⛔ Pemeriksaan di atas hanya membaca CurrencyID BARIS. Satu kolom
		// yang nilai uangnya bermata uang lain karena itu lolos pemeriksaan
		// itu - dan ronde pertama melabelinya ulang menjadi mata uang
		// penampung, sehingga ia terjumlah diam-diam di bawah mata uang yang
		// salah. Money.Add yang seharusnya menolaknya justru dilucuti.
		b := barisUji(t, "", kontrak.KodeAksep)
		usd, err := uang.NewMoney("5", "USD")
		if err != nil {
			t.Fatal(err)
		}
		b.SumInsured = usd
		if _, err := HitungTotalPeserta([]BarisAdjustment{b}, "IDR"); !errors.Is(err, ErrTotalMataUangBeragam) {
			t.Errorf("galat = %v, mau ErrTotalMataUangBeragam", err)
		}
	})

	t.Run("nilai uang TANPA label mata uang mengikut penampung", func(t *testing.T) {
		// Ia bukan mata uang lain; ia hanya belum berlabel. Menolaknya akan
		// menggagalkan seluruh pembacaan Detail untuk data yang sah.
		b := barisUji(t, "", kontrak.KodeAksep)
		tanpaLabel, err := uang.NewMoney("5", "")
		if err != nil {
			t.Fatal(err)
		}
		b.SumInsured = tanpaLabel
		total, err := HitungTotalPeserta([]BarisAdjustment{b}, "IDR")
		if err != nil {
			t.Fatalf("HitungTotalPeserta: %v", err)
		}
		if got := total.SumInsured.String(); got != "5 IDR" {
			t.Errorf("SumInsured = %q, mau %q", got, "5 IDR")
		}
	})

	t.Run("mata uang peserta kosong: diisi dari baris", func(t *testing.T) {
		total, err := HitungTotalPeserta([]BarisAdjustment{barisUji(t, "", kontrak.KodeAksep)}, "")
		if err != nil {
			t.Fatalf("HitungTotalPeserta: %v", err)
		}
		if got := total.JumlahKlaim.String(); got != "6 IDR" {
			t.Errorf("JumlahKlaim = %q, mau %q", got, "6 IDR")
		}
	})
}

// TestEnamKolomTotalDikunci mengunci CACAH dan URUTAN kolom total.
//
// ⛔ Sebabnya tercatat: ronde sebelumnya menghitung LIMA total, sebab
// pencacahannya memakai rujukan CheckTotalAdjustmentClaim - dan
// Total Ceding Retention (b20629/b20636) satu-satunya yang TIDAK punya
// aksi refresh, sehingga ia tidak ikut tercacah. Total yang hilang dari layar
// tidak berbunyi sendiri; uji inilah yang berbunyi.
func TestEnamKolomTotalDikunci(t *testing.T) {
	if len(kolomTotal) != JumlahKolomTotal {
		t.Fatalf("kolomTotal = %d, mau %d", len(kolomTotal), JumlahKolomTotal)
	}
	mau := []string{
		"CEDING_RETENTION", "SHARE_NUSANTARA_RE", "SUM_INSURED",
		"SUM_REASURED", "SHARE_RETRO", "CLAIM_AMOUNT",
	}
	for i, k := range kolomTotal {
		if k.nama != mau[i] {
			t.Errorf("kolomTotal[%d] = %q, mau %q", i, k.nama, mau[i])
		}
	}
}

// TestTiapKolomTotalMembacaKolomnyaSendiri membuktikan tidak ada pengambil
// yang tertukar: satu baris, satu kolom terisi, hanya satu total berubah.
func TestTiapKolomTotalMembacaKolomnyaSendiri(t *testing.T) {
	for i, kolom := range kolomTotal {
		b := BarisAdjustment{CurrencyID: "IDR"}
		// Isi HANYA kolom ke-i, lewat pemetaan yang ditulis TANGAN di bawah.
		isiKolom(t, &b, kolom.nama, uangPeserta(t, "7"))
		total, err := HitungTotalPeserta([]BarisAdjustment{b}, "IDR")
		if err != nil {
			t.Fatalf("%s: %v", kolom.nama, err)
		}
		for j, got := range ambilEnam(total) {
			mau := "0 IDR"
			if j == i {
				mau = "7 IDR"
			}
			if got != mau {
				t.Errorf("isi %s -> total %s = %q, mau %q",
					kolom.nama, kolomTotal[j].nama, got, mau)
			}
		}
	}
}

// isiKolom menulis satu nilai ke medan yang bernama kolom Pega itu.
//
// ⚠️ Ditulis TERPISAH dari kolomTotal dengan sengaja. Memakai kolomTotal.ambil
// untuk menyiapkan data lalu mengujinya dengan kolomTotal.ambil hanya
// membuktikan daftar itu konsisten dengan dirinya sendiri; pemetaan di sini
// ditulis tangan sehingga dua kolom yang tertukar di sana akan berbunyi.
func isiKolom(t *testing.T, b *BarisAdjustment, nama string, nilai uang.Money) {
	t.Helper()
	switch nama {
	case "CEDING_RETENTION":
		b.CedingRetention = nilai
	case "SHARE_NUSANTARA_RE":
		b.ShareNusantaraRe = nilai
	case "SUM_INSURED":
		b.SumInsured = nilai
	case "SUM_REASURED":
		b.SumReasured = nilai
	case "SHARE_RETRO":
		b.ShareRetro = nilai
	case "CLAIM_AMOUNT":
		b.JumlahKlaim = nilai
	default:
		t.Fatalf("kolom %q tidak dikenal uji ini - tambahkan pemetaannya", nama)
	}
}

// TestNamaJSONTotalPesertaDikunci mengunci keenam nama medan JSON.
//
// ⛔ SISI GO DARI KONTRAK DUA SISI. Pasangannya ada di
// `modul/claimlife/frontend/components/PanelTotalPeserta.test.ts`, blok
// "kontrak JSON TotalPeserta", dan keduanya memuat daftar nama yang sama.
//
// Sebabnya tercatat: tiga cacat berbentuk sama sudah terjadi di modul ini -
// envelope galat vs error, rute tanpa pemanggil, dan IsCheck "1" vs "true".
// Ketiganya lolos karena TIAP SISI HIJAU SENDIRIAN. Nama medan uang yang
// berganti di satu sisi membuat React membaca undefined dan menampilkan sel
// KOSONG, tanpa satu pun galat.
//
// ⚠️ Dikunci sebagai HIMPUNAN UTUH, bukan daftar nama terlarang. Pelajaran
// penjaga nama peserta: daftar tiga nama terlarang meloloskan nama keempat,
// dan penjaga yang meloloskan diam-diam lebih buruk daripada tidak ada.
func TestNamaJSONTotalPesertaDikunci(t *testing.T) {
	total, err := HitungTotalPeserta(nil, "IDR")
	if err != nil {
		t.Fatalf("HitungTotalPeserta: %v", err)
	}
	mentah, err := json.Marshal(total)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var peta map[string]json.RawMessage
	if err := json.Unmarshal(mentah, &peta); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	nama := make([]string, 0, len(peta))
	for k := range peta {
		nama = append(nama, k)
	}
	sort.Strings(nama)
	mau := []string{
		"cedingRetention", "jumlahKlaim", "shareNusantaraRe",
		"shareRetro", "sumInsured", "sumReasured",
	}
	if !reflect.DeepEqual(nama, mau) {
		t.Errorf("nama medan JSON = %v, mau %v.\n"+
			"Ubah JUGA PanelTotalPeserta.test.ts blok \"kontrak JSON TotalPeserta\" "+
			"dan tipe TotalPeserta di api.ts - ketiganya satu kontrak.", nama, mau)
	}
}

// TestTotalPesertaMenyeberangDiJSONPeserta membuktikan totalnya benar-benar
// sampai ke layar: medan yang dihitung tetapi tidak dikirim adalah persis
// cacat "rute tanpa pemanggil" dalam bentuk lain.
func TestTotalPesertaMenyeberangDiJSONPeserta(t *testing.T) {
	p := Peserta{NomorSertifikat: "UJI-001", MataUang: "IDR"}
	total, err := HitungTotalPeserta([]BarisAdjustment{barisUji(t, "", kontrak.KodeAksep)}, "IDR")
	if err != nil {
		t.Fatalf("HitungTotalPeserta: %v", err)
	}
	p.Total = total
	mentah, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var keluar struct {
		Total map[string]struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"total"`
	}
	if err := json.Unmarshal(mentah, &keluar); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(keluar.Total) != JumlahKolomTotal {
		t.Fatalf("medan total di JSON = %d, mau %d", len(keluar.Total), JumlahKolomTotal)
	}
	if got := keluar.Total["jumlahKlaim"].Amount; got != "6" {
		t.Errorf("total.jumlahKlaim.amount = %q, mau %q", got, "6")
	}
	// ⛔ TEKS, tidak pernah angka JSON (ADR-U-0016). Angka JSON dibaca
	// sebagai float64 oleh pustaka mana pun, dan itu uang yang rusak diam.
	if !strings.Contains(string(mentah), `"jumlahKlaim":{"amount":"6"`) {
		t.Errorf("jumlah uang tidak dikirim sebagai teks: %s", mentah)
	}
}
