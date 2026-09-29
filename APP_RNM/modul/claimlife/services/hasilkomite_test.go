package services_test

// Membaca hasil keputusan Komite dan memulai putaran berikutnya - TANPA Oracle.
//
// Pemilik: tiket 11. Dibaca sesudah: hasilkomite.go.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/kontrak"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/services"
)

// TestBarisLanjutanMewarisiDelapanKolomTanpaStatus - AC 5 spec.
func TestBarisLanjutanMewarisiDelapanKolomTanpaStatus(t *testing.T) {
	pertama := models.BarisAdjustment{
		ID: "A-1", KodeStatus: kontrak.KodeDitolak,
		CurrencyID: "IDR", NamaBank: "UJI-BANK", IDBank: "UJI-006",
		NomorRekening: "0012345", NomorAkseptasi: "UJI-AKS-1",
		JumlahKlaim:      uang(t, "1500000", "IDR"),
		ShareNusantaraRe: uang(t, "100", "IDR"),
		SumInsured:       uang(t, "200", "IDR"),
		KomiteID:         "KMT-000001",
	}
	p := models.Peserta{ID: "P-1", MataUang: "IDR",
		Baris: []models.BarisAdjustment{pertama}}

	baru, err := services.BarisLanjutan(p)
	if err != nil {
		t.Fatalf("BarisLanjutan: %v", err)
	}
	// Status Outstanding, BUKAN warisan.
	if baru.KodeStatus != kontrak.KodeOutstanding {
		t.Errorf("kode status = %q, mau Outstanding", baru.KodeStatus)
	}
	// ⛔ Pengenal, nomor akseptasi, dan tautan Komite TIDAK ikut. Baris baru
	// yang membawa KOMITE_ID barisnya sendiri akan tampak sudah diserahkan.
	for _, k := range []struct{ nama, isi string }{
		{"ID", baru.ID},
		{"NomorAkseptasi", baru.NomorAkseptasi},
		{"KomiteID", baru.KomiteID},
	} {
		if k.isi != "" {
			t.Errorf("%s ikut terbawa ke baris lanjutan: %q", k.nama, k.isi)
		}
	}
	// Kedelapan kolom warisan memang ikut: enam kolom uang + CURRENCY_ID +
	// mata uang klaim.
	if baru.CurrencyID != "IDR" || baru.JumlahKlaim.Currency != "IDR" {
		t.Errorf("mata uang tidak diwarisi: %+v", baru)
	}
	if baru.ShareNusantaraRe.Kosong() || baru.SumInsured.Kosong() {
		t.Error("kolom uang warisan kosong")
	}

	// ⛔ Medan bank TIDAK diwarisi - tiap baris punya tujuan pembayarannya
	// sendiri, dan `WarisiKolom` memang tidak menyalinnya.
	//
	// ⭐ Akibatnya nyata lintas tiket: baris lanjutan TIDAK lolos gerbang
	// rekening tiket 10 sampai seseorang mengisinya. Itu benar - rekening
	// putaran lama tidak otomatis menjadi rekening putaran baru - dan
	// dikunci di sini supaya tidak "diperbaiki" diam-diam kelak.
	for _, k := range []struct{ nama, isi string }{
		{"NamaBank", baru.NamaBank},
		{"IDBank", baru.IDBank},
		{"NomorRekening", baru.NomorRekening},
	} {
		if k.isi != "" {
			t.Errorf("%s diwarisi ke baris lanjutan: %q", k.nama, k.isi)
		}
	}
	if err := services.PeriksaRekening(baru); !errors.Is(
		err, services.ErrRekeningBelumLengkap) {
		t.Errorf("baris lanjutan lolos gerbang rekening: %v", err)
	}
	if !baru.TanggalAkseptasi.IsZero() {
		t.Error("tanggal akseptasi ikut terbawa")
	}
}

// TestBarisLanjutanMenolakPesertaTanpaBarisDitolak - putaran berikutnya lahir
// dari penolakan, bukan dari apa pun.
func TestBarisLanjutanMenolakPesertaTanpaBarisDitolak(t *testing.T) {
	for _, k := range []struct {
		apa string
		p   models.Peserta
	}{
		{"peserta tanpa baris", models.Peserta{ID: "P-1"}},
		{"baris terakhir masih Outstanding", models.Peserta{ID: "P-1",
			Baris: []models.BarisAdjustment{{ID: "A-1", KodeStatus: kontrak.KodeOutstanding}}}},
		{"baris terakhir sudah Aksep", models.Peserta{ID: "P-1",
			Baris: []models.BarisAdjustment{{ID: "A-1", KodeStatus: kontrak.KodeAksep}}}},
	} {
		if _, err := services.BarisLanjutan(k.p); !errors.Is(
			err, services.ErrBukanPenolakan) {
			t.Errorf("%s: galat = %v, mau ErrBukanPenolakan", k.apa, err)
		}
	}
}

// TestKlaimTidakTerminalSetelahPenolakan - ADR-U-0011.
//
// Penolakan menghasilkan putaran berikutnya, bukan akhir: status klaim adalah
// keadaan TURUNAN dari kumpulan barisnya, dan satu baris ditolak tidak
// menutup klaimnya.
func TestKlaimTidakTerminalSetelahPenolakan(t *testing.T) {
	k := models.Klaim{Peserta: []models.Peserta{{ID: "P-1",
		Baris: []models.BarisAdjustment{{ID: "A-1", KodeStatus: kontrak.KodeDitolak}}}}}
	ditolak := k.StatusTurunan()

	// Baris lanjutan lahir, dan klaimnya kembali berjalan.
	baru, err := services.BarisLanjutan(k.Peserta[0])
	if err != nil {
		t.Fatalf("BarisLanjutan: %v", err)
	}
	k.Peserta[0].Baris = append(k.Peserta[0].Baris, baru)
	if k.StatusTurunan() == ditolak {
		t.Error("status klaim tidak berubah sesudah baris lanjutan lahir; " +
			"penolakan menjadi akhir, padahal klaim tidak terminal (ADR-U-0011)")
	}
}

// TestPutaranMenjagaPagarnya - identitas, peran, pengenal, Oracle.
func TestPutaranMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if err := svc.Putaran().Tambah(ctx, inti.Pelaku{}, "CLM-1", "P-1",
		saatUji); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	// ⛔ Menambah baris lanjutan adalah menyimpan ke Outstanding - perannya
	// sama dengan penyimpanan Outstanding, bukan peran penolak.
	if err := svc.Putaran().Tambah(ctx,
		pelakuBerperan(inti.PeranMedicalAdvisor), "CLM-1", "P-1",
		saatUji); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("Medical menambah baris: galat = %v, mau ErrTanpaWewenang", err)
	}
	if err := svc.Putaran().Tambah(ctx,
		pelakuBerperan(services.PeranSimpanOutstanding), " ", "P-1",
		saatUji); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
}
