package services_test

// Penjaga statik tiket 11 - hasil Komite dibaca, tidak disimpan ulang.
//
// Pemilik: tiket 11. Dibaca sesudah: hasilkomite.go.
//
// Dua aturan dijaga:
//
//  1. ADR-U-0001: `AcceptStatus` dipetakan ke `STS_REJECT` DI BATAS, dan tidak
//     pernah disimpan sebagai status kedua. Satu baris satu status.
//  2. `[keputusan work owner 2026-09-15]`: Komite MENULIS, Claim Life MEMBACA.
//     Konteks ini tidak pernah menerapkan keputusan Komite ke baris lama.
//
// ⛔ Yang diperiksa TEMPAT PENYIMPANANNYA, bukan namanya di kode. Pelajaran
// tiket 12: pola atas teks sumber dapat dielakkan penulisan ulang yang setara,
// sedangkan sebuah kolom yang tidak ada di DDL tidak dapat menyimpan apa pun
// betapa pun kodenya ditulis ulang.

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// polaKolomStatusKedua mencocokkan kolom yang akan menyimpan keputusan Komite
// sebagai status kedua.
// ⛔ TANPA `\b`. Di Go, `\b` memperlakukan garis bawah sebagai huruf,
// sehingga `KOMITE_ACCEPT_STATUS` LOLOS dari `\bACCEPT_?STATUS\b` - elakan
// yang dibangun dan terbukti hijau. Kini dicocokkan sebagai POTONGAN, bukan
// sebagai kata utuh.
var polaKolomStatusKedua = regexp.MustCompile(
	`(?i)(ACCEPT_?STATUS|APROVAL|APPROVAL|DATE_?APPROVE|APPROVE_?DATE|KOMITE_COMMENT|COMMENT_?TEXT)`)

// TestNolKolomStatusKeduaDiSkema - ADR-U-0001, diperiksa di DDL.
func TestNolKolomStatusKeduaDiSkema(t *testing.T) {
	// ⛔ MENELUSURI, bukan membaca satu direktori. Ronde pertama memakai
	// `ReadDir` tanpa rekursi dan `HasSuffix(".sql")` yang peka huruf besar,
	// sehingga `migrations/komite/011.sql` dan `012.SQL` tidak pernah dibaca.
	// Keduanya dibangun sebagai elakan dan terbukti hijau.
	dir := filepath.Join("..", "repository", "migrations")
	diperiksa := 0
	err := filepath.Walk(dir, func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.EqualFold(filepath.Ext(jalur), ".sql") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		diperiksa++
		if m := polaKolomStatusKedua.FindString(string(isi)); m != "" {
			t.Errorf("%s: memuat %q; keputusan Komite dipetakan ke STS_REJECT "+
				"DI BATAS dan tidak disimpan sebagai status kedua (ADR-U-0001). "+
				"Roster dan keputusan per anggota milik konteks Komite Claim Life",
				filepath.ToSlash(jalur), m)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ Pembacanya sendiri diperiksa: penjaga yang membaca nol berkas hijau
	// selamanya, dan hijau selamanya tidak dapat dibedakan dari aman.
	if diperiksa < 10 {
		t.Fatalf("hanya %d berkas migrasi terbaca; pembacanya yang rusak", diperiksa)
	}
}

// TestBarisLanjutanSelaluMemulaiPutaranBaru - PERILAKU, bukan teks.
//
// ⛔ Apa pun bentuk penulisan ulangnya, baris yang lahir dari penolakan wajib:
// berstatus Outstanding, tanpa pengenal, tanpa nomor akseptasi, dan tanpa
// tautan Komite. Bila salah satunya bocor, baris baru akan tampak sudah
// diserahkan atau sudah diputus - dan putaran berikutnya tidak pernah mulai.
func TestBarisLanjutanSelaluMemulaiPutaranBaru(t *testing.T) {
	asal := models.Peserta{ID: "P-1", MataUang: "IDR",
		Baris: []models.BarisAdjustment{{
			ID: "A-1", KodeStatus: models.KodeDitolak,
			NomorAkseptasi: "UJI-AKS", KomiteID: "KMT-000009",
			CurrencyID: "IDR",
		}}}
	// Salinan dalam sebelum dipanggil - pembanding yang tidak ikut berubah.
	sebelum := asal.Baris[0]

	baru, err := services.BarisLanjutan(asal)
	if err != nil {
		t.Fatal(err)
	}
	if baru.KodeStatus != models.KodeOutstanding {
		t.Errorf("baris lanjutan berkode %q, mau Outstanding", baru.KodeStatus)
	}
	for _, k := range []struct{ nama, isi string }{
		{"ID", baru.ID}, {"NomorAkseptasi", baru.NomorAkseptasi},
		{"KomiteID", baru.KomiteID},
	} {
		if k.isi != "" {
			t.Errorf("%s bocor ke putaran baru: %q", k.nama, k.isi)
		}
	}
	// ⛔ Dan baris asalnya TIDAK berubah - SELURUHNYA, bukan dua medan.
	//
	// `BarisLanjutan` menerima peserta sebagai NILAI, tetapi slice berbagi
	// array yang sama dengan pemanggilnya: menulis lewat `p.Baris[i]` menimpa
	// baris ASLI meski parameternya nilai. Ronde pertama hanya memeriksa
	// `KodeStatus` dan panjangnya, dan mengaku "dilindungi parameter nilai" -
	// pengakuan yang lebih besar daripada yang diperiksanya. Elakan yang
	// menulis `NomorAkseptasi` lolos utuh.
	//
	// Kini seluruh nilai baris dibandingkan sekaligus, sehingga medan mana pun
	// yang ternoda akan tertangkap - termasuk yang belum ada hari ini.
	if !reflect.DeepEqual(asal.Baris[0], sebelum) {
		t.Errorf("baris asal ternoda:\n  sebelum: %+v\n  sesudah: %+v; "+
			"keputusan Komite ditimpa konteks yang hanya berhak membacanya",
			sebelum, asal.Baris[0])
	}
	if len(asal.Baris) != 1 {
		t.Errorf("baris asal bertambah menjadi %d; BarisLanjutan menulis ke "+
			"peserta pemanggilnya", len(asal.Baris))
	}
}

// TestKonteksIniTidakPernahMenulisKeputusan - AC terakhir tiket 11.
//
// ⛔ Bunyi AC-nya *"tiket ini tidak menulis STS_REJECT"*, tetapi harfiahnya
// tidak benar dan tidak boleh dipura-purakan: `BarisLanjutan` MENULIS
// `STS_REJECT = 0` pada baris BARU. Yang sesungguhnya dijaga lebih tajam dan
// lebih berguna: konteks ini tidak pernah menulis KEPUTUSAN - `1` (Aksep)
// maupun `2` (Ditolak). Keduanya milik Komite Claim Life
// `[keputusan work owner 2026-09-15]`.
//
// Diperiksa lewat PERILAKU: LIMA bentuk kode asal dicoba - ketiga kode yang
// dikenal, kosong, dan satu yang asing - dan tidak satu pun menghasilkan
// baris berkeputusan. Lima, bukan lebih: angkanya disebut apa adanya.
func TestKonteksIniTidakPernahMenulisKeputusan(t *testing.T) {
	dicoba := 0
	for _, kodeAsal := range []string{
		models.KodeDitolak, models.KodeOutstanding, models.KodeAksep, "", "9",
	} {
		p := models.Peserta{ID: "P-1", MataUang: "IDR",
			Baris: []models.BarisAdjustment{{ID: "A-1", KodeStatus: kodeAsal}}}
		dicoba++
		baru, err := services.BarisLanjutan(p)
		if err != nil {
			// Hanya penolakan yang membuka putaran; sisanya memang ditolak.
			continue
		}
		if baru.KodeStatus == models.KodeAksep || baru.KodeStatus == models.KodeDitolak {
			t.Errorf("asal %q menghasilkan baris berKEPUTUSAN %q; keputusan "+
				"milik Komite Claim Life, konteks ini hanya membacanya",
				kodeAsal, baru.KodeStatus)
		}
	}
	if dicoba != 5 {
		t.Fatalf("%d bentuk dicoba, mau 5", dicoba)
	}

	// ⛔ Dan berkas ini tidak boleh memanggil PENGUBAH STATUS BARIS LAMA.
	//
	// Elakan yang terbukti hijau terhadap ronde pertama: menyisipkan
	// `baca.PerbaruiStatusBaris(...)` ke dalam transaksi `Tambah`, sehingga
	// konteks ini menimpa `STS_REJECT` baris LAMA menjadi Aksep - persis yang
	// `[keputusan work owner 2026-09-15]` larang - dan seluruh suite tetap
	// hijau. Sebabnya: bagian perilaku hanya memanggil `BarisLanjutan`, tidak
	// pernah `Tambah`; bagian tekstualnya hanya mencari penugasan literal.
	//
	// ⚠️ Cakupannya DINYATAKAN: ia memeriksa teks berkas ini saja, dan
	// pemanggilan yang disembunyikan di balik variabel perantara atau di
	// berkas lain tidak tertangkap. Yang menutup sisanya pemeriksaan SAAT
	// JALAN di `repository.PeriksaBarisBaru`.
	isi, err := os.ReadFile(filepath.Join("..", "services", "hasilkomite.go"))
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentar(string(isi))
	for _, terlarang := range []string{
		"PerbaruiStatusBaris(", "CabutPenandaDipilih(", ".Ubah(", ".Tolak(",
	} {
		if strings.Contains(kode, terlarang) {
			t.Errorf("hasilkomite.go memanggil %q; konteks ini MEMBACA keputusan "+
				"Komite, ia tidak pernah mengubah status baris lama", terlarang)
		}
	}
}
