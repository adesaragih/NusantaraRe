//go:build db

package services_test

// Grid diagnosa terhadap skema uji Oracle - butir bd.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus.
//
// Dibaca sesudah: services/diagnosa.go dan statusbaris_db_test.go - fixture
// pohonnya dipakai ulang apa adanya.

import (
	"context"
	"errors"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/modul/claimlife/backend/services"
)

// pelakuAdminUji memegang tahap Outstanding, yaitu tahap awal pohon uji.
var pelakuAdminUji = inti.Pelaku{
	AkunID: "UJI-DIAG", Peran: []string{models.PeranAdminLife},
}

// TestDiagnosaTambahUbahHapus menempuh ketiga tombol grid berurutan.
func TestDiagnosaTambahUbahHapus(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiStatus(t, svc, db)
	ctx := context.Background()
	peserta, err := repository.NewKlaimLife(db).AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v (%d)", err, len(peserta))
	}
	pesertaID := peserta[0].ID
	diag := svc.Diagnosa()

	// `Add` b4690 - tiga kali, dan URUTANNYA 1, 2, 3.
	var lahir []models.Diagnosa
	for i := 1; i <= 3; i++ {
		d, err := diag.Tambah(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID)
		if err != nil {
			t.Fatalf("Tambah ke-%d: %v", i, err)
		}
		if d.Urutan != i {
			t.Errorf("Tambah ke-%d berurut %d, mau %d", i, d.Urutan, i)
		}
		lahir = append(lahir, d)
	}
	// ⛔ Baris baru KOSONG - padanan `addRow` b4700, yang menambah anggota
	// PageList yang propertinya belum terisi.
	if lahir[0].Nama != "" || lahir[0].KodeICD != "" {
		t.Errorf("baris baru sudah berisi: %+v", lahir[0])
	}

	// `Choose` b2509 -> `SetDisease` b260/b307.
	if err := diag.Ubah(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID,
		lahir[1].ID, "E11", "Diabetes mellitus tipe 2", "Metabolik"); err != nil {
		t.Fatalf("Ubah: %v", err)
	}
	perPeserta, err := repository.NewDiagnosa(db).AmbilDiagnosa(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca diagnosa: %v", err)
	}
	daftar := perPeserta[pesertaID]
	if len(daftar) != 3 {
		t.Fatalf("diagnosa = %d, mau 3", len(daftar))
	}
	if daftar[1].Nama != "Diabetes mellitus tipe 2" || daftar[1].KodeICD != "E11" {
		t.Errorf("baris kedua = %+v", daftar[1])
	}
	if daftar[1].GroupDiagnose != "Metabolik" {
		t.Errorf("GROUP_DIAGNOSE = %q, mau %q", daftar[1].GroupDiagnose, "Metabolik")
	}
	// ⛔ Dan HANYA baris itu yang berubah. Pembaruan yang menyentuh baris lain
	// tidak akan terlihat dari uji satu-baris mana pun.
	if daftar[0].Nama != "" || daftar[2].Nama != "" {
		t.Errorf("baris lain ikut berubah: %+v / %+v", daftar[0], daftar[2])
	}

	// `Delete` b6160 atas baris TENGAH - dan urutannya dirapatkan.
	if err := diag.Hapus(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID,
		lahir[0].ID); err != nil {
		t.Fatalf("Hapus: %v", err)
	}
	perPeserta, err = repository.NewDiagnosa(db).AmbilDiagnosa(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca ulang diagnosa: %v", err)
	}
	daftar = perPeserta[pesertaID]
	if len(daftar) != 2 {
		t.Fatalf("sesudah hapus = %d baris, mau 2", len(daftar))
	}
	// ⛔ 1 dan 2, tanpa lubang. Daftar berlubang akan memberi nomor yang
	// sudah dipakai kepada baris berikutnya.
	for i, d := range daftar {
		if d.Urutan != i+1 {
			t.Errorf("sesudah hapus, baris %d berurut %d, mau %d", d.ID, d.Urutan, i+1)
		}
	}
}

// TestKeputusanPesertaMencerminKeSeluruhDiagnosa - `SetSTS_Reject` b241/b257.
func TestKeputusanPesertaMencerminKeSeluruhDiagnosa(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiStatus(t, svc, db)
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)
	peserta, err := baca.AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v", err)
	}
	pesertaID := peserta[0].ID
	diag := svc.Diagnosa()
	for i := 0; i < 2; i++ {
		if _, err := diag.Tambah(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID); err != nil {
			t.Fatalf("Tambah: %v", err)
		}
	}

	perBaris, err := baca.AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca baris: %v", err)
	}
	adj := perBaris[pesertaID]
	if len(adj) != 1 {
		t.Fatalf("baris adjustment = %d, mau 1", len(adj))
	}

	// Keputusan peserta: aksep.
	saat := time.Date(2026, 9, 27, 10, 0, 0, 0, time.UTC)
	if err := svc.Status().DenganJejak(&jejakUji{}).Ubah(ctx, pelakuAdminUji,
		pohon.Work.ID, pesertaID, adj[0].ID, kontrak.StatusAksep, saat); err != nil {
		t.Fatalf("Ubah status: %v", err)
	}

	// ⛔ SETIAP diagnosa ikut, bukan yang pertama saja - padanan putaran
	// `EMBEDDED` b345 atas seluruh `.DiagnoseList`.
	perPeserta, err := repository.NewDiagnosa(db).AmbilDiagnosa(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca diagnosa: %v", err)
	}
	daftar := perPeserta[pesertaID]
	if len(daftar) != 2 {
		t.Fatalf("diagnosa = %d, mau 2", len(daftar))
	}
	for _, d := range daftar {
		if d.KodeStatus != kontrak.KodeAksep {
			t.Errorf("diagnosa %d berkode %q, mau %q", d.ID, d.KodeStatus, kontrak.KodeAksep)
		}
	}

	// ⛔ Dan SESUDAH itu ketiga tombolnya tertutup - gerbang b4682/b5059/
	// b5870/b6152. Di Pega tombolnya mati; di sini permintaannya ditolak.
	if _, err := diag.Tambah(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID); !errors.Is(
		err, services.ErrDiagnosaTerkunci) {
		t.Errorf("Tambah sesudah diputus: %v, mau ErrDiagnosaTerkunci", err)
	}
	if err := diag.Ubah(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID,
		daftar[0].ID, "A00", "x", ""); !errors.Is(err, services.ErrDiagnosaTerkunci) {
		t.Errorf("Ubah sesudah diputus: %v, mau ErrDiagnosaTerkunci", err)
	}
	if err := diag.Hapus(ctx, pelakuAdminUji, pohon.Work.ID, pesertaID,
		daftar[0].ID); !errors.Is(err, services.ErrDiagnosaTerkunci) {
		t.Errorf("Hapus sesudah diputus: %v, mau ErrDiagnosaTerkunci", err)
	}
}

// TestDiagnosaBukanMilikPesertaDitolak menjaga sarang jalurnya.
//
// ⛔ Pengenal diagnosa datang dari JALUR URL, dan jalur URL datang dari siapa
// saja. Gerbang tahap dan gerbang `STS_REJECT` memeriksa PESERTA-nya, bukan
// baris yang disebut - tanpa pemeriksaan kepemilikan, `DELETE
// .../peserta/B/diagnosa/9` menghapus baris 9 milik siapa pun.
func TestDiagnosaBukanMilikPesertaDitolak(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiStatus(t, svc, db)
	ctx := context.Background()
	peserta, err := repository.NewKlaimLife(db).AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v", err)
	}
	diag := svc.Diagnosa()
	if _, err := diag.Tambah(ctx, pelakuAdminUji, pohon.Work.ID, peserta[0].ID); err != nil {
		t.Fatalf("Tambah: %v", err)
	}
	// Pengenal yang tidak pernah lahir pada peserta ini.
	for _, uji := range []struct {
		nama  string
		jalan func() error
	}{
		{"ubah", func() error {
			return diag.Ubah(ctx, pelakuAdminUji, pohon.Work.ID, peserta[0].ID,
				999999999, "A00", "x", "")
		}},
		{"hapus", func() error {
			return diag.Hapus(ctx, pelakuAdminUji, pohon.Work.ID, peserta[0].ID, 999999999)
		}},
	} {
		if err := uji.jalan(); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%s baris asing: %v, mau ErrPermintaanTidakSah", uji.nama, err)
		}
	}
}
