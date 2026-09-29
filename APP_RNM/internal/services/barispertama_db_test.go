//go:build db

package services_test

// Baris adjustment PERTAMA terhadap skema uji Oracle - GILIRAN-13 butir bo.
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

var pelakuSPVUji = services.Pelaku{AkunID: "UJI-SPV", Peran: []string{services.PeranSPV}}

// pohonTanpaBaris menyimpan satu klaim berpeserta TANPA baris adjustment -
// keadaan sesudah pendaftaran (`pendaftaran.go`, `Baris: []`).
func pohonTanpaBaris(t *testing.T, svc *services.Service, db *repository.DB,
	id string, tahap models.Tahap, kodePeserta string) (models.PohonKlaim, string) {

	t.Helper()
	peran, _ := models.PeranPemegangTahap(tahap)
	pohon := models.PohonKlaim{
		Work: models.WorkClaim{ID: id, Lini: models.LiniLife, Type: "QP",
			Tahap: tahap.String(), PyPosition: peran},
		Klaim: models.Klaim{
			NomorKlaim: "UJI-" + id,
			Peserta: []models.Peserta{{
				NomorSertifikat: "006", MataUang: "IDR", KodeStatus: kodePeserta,
			}},
		},
	}
	ctx := context.Background()
	if err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return repository.NewPohonKlaim(db).Simpan(ctx, tx, pohon)
	}); err != nil {
		t.Fatalf("menyiapkan pohon: %v", err)
	}
	peserta, err := repository.NewKlaimLife(db).AmbilPeserta(ctx, id)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v (%d)", err, len(peserta))
	}
	return pohon, peserta[0].ID
}

// TestAddMelahirkanBarisPertamaKosong - jalur penuh `Add` b17937 pada grid
// kosong, lalu `Add` kedua yang menjadi jalur putaran.
func TestAddMelahirkanBarisPertamaKosong(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon, pesertaID := pohonTanpaBaris(t, svc, db, "CLM-UJI730", models.TahapClaimAnalis, "")
	ctx := context.Background()
	jejak := &jejakUji{}
	if err := svc.Putaran().DenganJejak(jejak).Tambah(ctx, pelakuSPVUji,
		pohon.Work.ID, pesertaID, saatUjiKomite); err != nil {
		t.Fatalf("Add pada grid kosong: %v", err)
	}

	// Dibaca ULANG dari Oracle.
	baca := repository.NewKlaimLife(db)
	perBaris, err := baca.AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	daftar := perBaris[pesertaID]
	if len(daftar) != 1 {
		t.Fatalf("baris = %d, mau 1", len(daftar))
	}
	if b := daftar[0]; b.KodeStatus != "" || b.JumlahKlaim.Currency != "" ||
		b.JumlahKlaim.Amount != nil || b.NomorAkseptasi != "" {
		t.Errorf("baris pertama tidak kosong: %+v", b)
	}
	peserta, err := baca.AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if peserta[0].IsCheck != models.PenandaDipilih {
		t.Errorf("IS_CHECK = %q, mau %q (langkah 1 b328 activity pewaris)",
			peserta[0].IsCheck, models.PenandaDipilih)
	}
	if len(jejak.catatan) != 1 || jejak.catatan[0].AdjustmentID != daftar[0].ID {
		t.Errorf("jejak = %+v, mau satu catatan untuk baris %s", jejak.catatan, daftar[0].ID)
	}

	// `Add` kedua = jalur putaran: baris pertama belum ditolak.
	if err := svc.Putaran().DenganJejak(jejak).Tambah(ctx, pelakuSPVUji,
		pohon.Work.ID, pesertaID, saatUjiKomite); !errors.Is(err, services.ErrBukanPenolakan) {
		t.Errorf("Add kedua: galat = %v, mau ErrBukanPenolakan", err)
	}
}

// TestAddBarisPertamaMenjagaGerbangnya - tahap, lalu `.STS_REJECT` peserta.
func TestAddBarisPertamaMenjagaGerbangnya(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()
	ctx := context.Background()

	// Tahap Outstanding: `Add` tidak tampil (b18160), dan tidak dilayani.
	_, pesertaOS := pohonTanpaBaris(t, svc, db, "CLM-UJI731", models.TahapOutstanding, "")
	if err := svc.Putaran().DenganJejak(&jejakUji{}).Tambah(ctx, pelakuSPVUji,
		"CLM-UJI731", pesertaOS, saatUjiKomite); !errors.Is(
		err, services.ErrTahapTanpaAddAdjustment) {
		t.Errorf("Add di Outstanding: galat = %v, mau ErrTahapTanpaAddAdjustment", err)
	}

	// Peserta yang sudah diputus: isian layar Detail beku.
	_, pesertaDitolak := pohonTanpaBaris(t, svc, db, "CLM-UJI732", models.TahapClaimAnalis,
		models.KodeDitolak)
	if err := svc.Putaran().DenganJejak(&jejakUji{}).Tambah(ctx, pelakuSPVUji,
		"CLM-UJI732", pesertaDitolak, saatUjiKomite); !errors.Is(
		err, services.ErrPesertaSudahDiputus) {
		t.Errorf("Add pada peserta ditolak: galat = %v, mau ErrPesertaSudahDiputus", err)
	}

	// Peserta milik klaim LAIN - pengenalnya datang dari URL. Klaimnya di
	// Claim Analis, jadi gerbang tahap lolos dan yang menolak kepemilikan.
	if err := svc.Putaran().DenganJejak(&jejakUji{}).Tambah(ctx, pelakuSPVUji,
		"CLM-UJI732", pesertaOS, saatUjiKomite); !errors.Is(err, services.ErrPermintaanTidakSah) {
		t.Errorf("peserta klaim lain: galat = %v, mau ErrPermintaanTidakSah", err)
	}

	// Nol baris lahir dari ketiga penolakan di atas.
	for _, id := range []string{"CLM-UJI731", "CLM-UJI732"} {
		perBaris, err := repository.NewKlaimLife(db).AmbilBaris(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if len(perBaris) != 0 {
			t.Errorf("%s: baris lahir walau ditolak: %v", id, perBaris)
		}
	}
}
