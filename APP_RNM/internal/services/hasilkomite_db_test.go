//go:build db

package services_test

// Putaran adjustment berikutnya terhadap skema uji Oracle - tiket 11.
//
// ⛔ Seluruh test di sini MELEWATI bila ORACLE_DSN belum dikonfigurasi.
// Melewati bukan lulus.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

// TestPutaranMelahirkanBarisOutstandingBaru - AC 5 spec, jalur penuh.
func TestPutaranMelahirkanBarisOutstandingBaru(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiKomite(t, svc, db, "CLM-UJI720", "UJI-BANK", "UJI-006", "0012345")
	pesertaID, baris := barisPertama(t, db, pohon.Work.ID)
	ctx := context.Background()

	// Baris pertama ditolak dulu - putaran berikutnya lahir dari penolakan.
	jejak := &jejakUji{}
	if err := svc.Status().DenganJejak(jejak).Tolak(ctx,
		services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranAdmin}},
		pohon.Work.ID, baris.ID, saatUjiKomite); err != nil {
		t.Fatalf("menolak baris pertama: %v", err)
	}

	// ⛔ Putaran berikutnya TIDAK boleh dibuka oleh peran penolak.
	if err := svc.Putaran().DenganJejak(jejak).Tambah(ctx,
		services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranAdmin}},
		pohon.Work.ID, pesertaID, saatUjiKomite); !errors.Is(
		err, services.ErrTanpaWewenang) {
		t.Errorf("Admin membuka putaran: galat = %v, mau ErrTanpaWewenang", err)
	}

	if err := svc.Putaran().DenganJejak(jejak).Tambah(ctx,
		services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranSPV}},
		pohon.Work.ID, pesertaID, saatUjiKomite); err != nil {
		t.Fatalf("membuka putaran: %v", err)
	}

	// Dibaca ULANG dari Oracle - bukan dipercaya dari nilai yang baru ditulis.
	perBaris, err := repository.NewKlaimLife(db).AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	daftar := perBaris[pesertaID]
	if len(daftar) != 2 {
		t.Fatalf("baris = %d, mau 2", len(daftar))
	}
	lama, baru := daftar[0], daftar[1]
	// ⛔ Keputusan baris lama TIDAK tersentuh. Konteks ini membaca keputusan
	// Komite, ia tidak menimpanya.
	if lama.KodeStatus != models.KodeDitolak {
		t.Errorf("baris lama berkode %q; keputusannya tertimpa", lama.KodeStatus)
	}
	if baru.KodeStatus != models.KodeOutstanding {
		t.Errorf("baris baru berkode %q, mau Outstanding", baru.KodeStatus)
	}
	if baru.ID == lama.ID || baru.ID == "" {
		t.Errorf("pengenal baris baru = %q (lama %q)", baru.ID, lama.ID)
	}
	// Kolom warisan ikut; tautan Komite dan nomor akseptasi TIDAK.
	if baru.CurrencyID != lama.CurrencyID {
		t.Errorf("CURRENCY_ID tidak diwarisi: %q vs %q", baru.CurrencyID, lama.CurrencyID)
	}
	if baru.KomiteID != "" || baru.NomorAkseptasi != "" {
		t.Errorf("putaran baru membawa sisa putaran lama: KOMITE_ID=%q ACCEPTED_NO=%q",
			baru.KomiteID, baru.NomorAkseptasi)
	}
	// ⛔ Medan bank TIDAK diwarisi, sehingga baris baru belum dapat diserahkan
	// ke Komite sampai rekeningnya diisi - gerbang tiket 10.
	if baru.NamaBank != "" {
		t.Errorf("medan bank diwarisi: %q", baru.NamaBank)
	}
	// ⛔ PENANDA DIPILIH DIPULIHKAN. `Tolak` mencabutnya (IS_CHECK='false'),
	// dan tanpa memulihkannya baris lanjutan TIDAK AKAN PERNAH diambil
	// Komite: `[terverifikasi]` `Komite Claim Life/Activity/
	// KomitePostAdjustment.xml` pecahan baris 1900, 2117, 4896, 5341, 7622,
	// dan 8067 - keenam prasyaratnya menuntut `.IsCheck = true`.
	// Jendelanya "langkah prasyarat yang BERBEDA" - lihat sensus bernama di
	// `hasilkomite.go`, yang juga menyebut angka 8 dan 13 untuk jendela lain.
	//
	// Putaran yang lahir tetapi tidak dapat diambil siapa pun bukan putaran.
	peserta, err := repository.NewKlaimLife(db).AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v (%d)", err, len(peserta))
	}
	if peserta[0].IsCheck != "true" {
		t.Errorf("IS_CHECK peserta = %q sesudah putaran baru, mau \"true\"; "+
			"baris lanjutan tidak akan pernah diambil Komite", peserta[0].IsCheck)
	}

	// Header mencerminkan baris TERAKHIR.
	hdr, err := repository.NewKlaimLife(db).AmbilHeader(ctx, pohon.Work.ID)
	if err != nil || hdr == nil {
		t.Fatalf("membaca header: %v", err)
	}
	if hdr.KodeStatus != models.KodeOutstanding {
		t.Errorf("header berkode %q, mau Outstanding - ia mencerminkan baris "+
			"terakhir, dan baris terakhir kini yang baru", hdr.KodeStatus)
	}
}

// TestPutaranKeduaDitolakSelamaBarisTerakhirBelumDiputus - satu putaran
// berjalan pada satu waktu.
func TestPutaranKeduaDitolakSelamaBarisTerakhirBelumDiputus(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiKomite(t, svc, db, "CLM-UJI721", "UJI-BANK", "UJI-006", "0012345")
	pesertaID, baris := barisPertama(t, db, pohon.Work.ID)
	ctx := context.Background()
	spv := services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranSPV}}

	// Baris masih Outstanding: belum ada penolakan, jadi belum ada putaran baru.
	if err := svc.Putaran().DenganJejak(&jejakUji{}).Tambah(ctx, spv,
		pohon.Work.ID, pesertaID, saatUjiKomite); !errors.Is(
		err, services.ErrBukanPenolakan) {
		t.Fatalf("galat = %v, mau ErrBukanPenolakan", err)
	}
	_ = baris
}
