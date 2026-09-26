//go:build db

package services_test

// Penyerahan ke Komite terhadap skema uji Oracle - tiket 10.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus.
//
// Dibaca sesudah: services/komite.go.

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
	"nusantarare/pkg/utils"
)

// rosterUji menggantikan `POOLDATA.EMAILKOMITE` yang belum disahkan (butir af).
type rosterUji struct {
	tingkat int
	ambang  []models.Money
}

func (r *rosterUji) CacahTingkat(_ context.Context, ambang models.Money,
	_ string) (int, error) {
	r.ambang = append(r.ambang, ambang)
	return r.tingkat, nil
}

// kasusUji menggantikan penyimpan kasus Komite yang belum disahkan (butir af).
type kasusUji struct {
	muatan []services.MuatanKomite
	id     string
}

func (k *kasusUji) Buat(_ context.Context, _ *repository.Tx,
	m services.MuatanKomite) (string, error) {
	k.muatan = append(k.muatan, m)
	return k.id, nil
}

// pohonUjiKomite membuat satu klaim Type TP - bebas peran menurut tiket 07 -
// dengan satu baris Outstanding yang rekeningnya LENGKAP.
func pohonUjiKomite(t *testing.T, svc *services.Service, db *repository.DB,
	workID string, bank, idBank, rekening string) models.PohonKlaim {

	t.Helper()
	pohon := models.PohonKlaim{
		Work: models.WorkClaim{ID: workID, Lini: models.LiniLife, Type: "TP"},
		Klaim: models.Klaim{
			NomorKlaim: "UJI-CLM-" + workID,
			Peserta: []models.Peserta{{
				NomorSertifikat: "006", MataUang: "IDR",
				Baris: []models.BarisAdjustment{{
					KodeStatus: models.KodeOutstanding,
					// ⛔ KEDUA kolom mata uang diisi. Ronde pertama hanya
					// mengisi CURRENCY_ID dan membiarkan JumlahKlaim kosong -
					// titik buta yang persis meloloskan cacat di
					// PeriksaSatuMataUang: yang menyeberang justru CURRENCY.
					CurrencyID:    "IDR",
					JumlahKlaim:   models.Money{Amount: apd.New(1500000, 0), Currency: "IDR"},
					NamaBank:      bank,
					IDBank:        idBank,
					NomorRekening: rekening,
				}},
			}},
		},
	}
	ctx := context.Background()
	err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return repository.NewPohonKlaim(db).Simpan(ctx, tx, pohon)
	})
	if err != nil {
		t.Fatalf("menyiapkan pohon: %v", err)
	}
	return pohon
}

// barisPertama membaca peserta dan baris pertama sebuah klaim.
func barisPertama(t *testing.T, db *repository.DB, workID string) (string, models.BarisAdjustment) {
	t.Helper()
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)
	peserta, err := baca.AmbilPeserta(ctx, workID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v (%d peserta)", err, len(peserta))
	}
	perBaris, err := baca.AmbilBaris(ctx, workID)
	if err != nil {
		t.Fatalf("membaca baris: %v", err)
	}
	adj := perBaris[peserta[0].ID]
	if len(adj) != 1 {
		t.Fatalf("baris = %d, mau 1", len(adj))
	}
	return peserta[0].ID, adj[0]
}

// TestGerbangRekeningDitegakkanDiLayananBukanDiLayar - AC tiket 10.
//
// ⛔ Test ini memanggil LAYANAN LANGSUNG, tanpa melewati satu pun kontrol
// layar. Layar memang menampilkan tombol "Send ke Komite" untuk baris ini -
// kondisi tampilnya hanya status Outstanding dan KOMITE_ID kosong, dan ia
// TIDAK memeriksa medan bank sama sekali. Jadi bila gerbangnya tidak ada di
// layanan, ia tidak ada di mana pun.
func TestGerbangRekeningDitegakkanDiLayananBukanDiLayar(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	// Rekening kosong - persis keadaan yang layar biarkan lewat.
	pohon := pohonUjiKomite(t, svc, db, "CLM-UJI710", "", "", "")
	pesertaID, baris := barisPertama(t, db, pohon.Work.ID)

	roster := &rosterUji{tingkat: 2}
	kasus := &kasusUji{id: "KMT-UJI-1"}
	err := svc.Komite().DenganRoster(roster).DenganKasus(kasus).
		DenganJejak(&jejakUji{}).Serahkan(context.Background(),
		services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranAdmin}},
		pohon.Work.ID, pesertaID, baris.ID, saatUjiKomite)

	if !errors.Is(err, services.ErrRekeningBelumLengkap) {
		t.Fatalf("galat = %v, mau ErrRekeningBelumLengkap", err)
	}
	// ⛔ Gerbangnya berhenti SEBELUM roster disentuh. Memanggil roster untuk
	// baris yang tidak boleh diserahkan berarti menyentuh tabel produksi
	// tanpa alasan.
	if len(roster.ambang) != 0 {
		t.Errorf("roster dipanggil %d kali padahal rekening kosong", len(roster.ambang))
	}
	if len(kasus.muatan) != 0 {
		t.Errorf("kasus komite dibuat %d kali padahal rekening kosong", len(kasus.muatan))
	}
}

var saatUjiKomite = time.Date(2026, 9, 27, 10, 0, 0, 0, time.UTC)

// TestPenyerahanMenautkanBarisDanMembawaMuatannya - jalur lengkap.
func TestPenyerahanMenautkanBarisDanMembawaMuatannya(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiKomite(t, svc, db, "CLM-UJI711", "UJI-BANK", "UJI-006", "0012345")
	pesertaID, baris := barisPertama(t, db, pohon.Work.ID)

	roster := &rosterUji{tingkat: 3}
	kasus := &kasusUji{id: "KMT-UJI-711"}
	jejak := &jejakUji{}
	ctx := context.Background()
	err := svc.Komite().DenganRoster(roster).DenganKasus(kasus).DenganJejak(jejak).
		Serahkan(ctx, services.Pelaku{AkunID: "UJI-AKUN",
			Peran: []string{services.PeranAdmin}},
			pohon.Work.ID, pesertaID, baris.ID, saatUjiKomite)
	if err != nil {
		t.Fatalf("menyerahkan: %v", err)
	}

	// ⭐ AC 24: nilai klaim, mata uangnya, dan status baris SAAT penyerahan -
	// tiga hal yang tidak ada di sistem lama.
	if len(kasus.muatan) != 1 {
		t.Fatalf("muatan = %d, mau 1", len(kasus.muatan))
	}
	m := kasus.muatan[0]
	if m.KodeStatus != models.KodeOutstanding {
		t.Errorf("muatan membawa kode %q, mau Outstanding", m.KodeStatus)
	}
	if m.TingkatKomite != 3 {
		t.Errorf("tingkat = %d, mau 3", m.TingkatKomite)
	}
	if m.AkunID != "UJI-AKUN" {
		t.Errorf("muatan tidak membawa pelaku: %q", m.AkunID)
	}
	// ⭐ Nilai klaim BESERTA mata uangnya - AC 24. Mata uang yang menyeberang
	// adalah CURRENCY, kolom yang dulu tidak dijaga.
	if m.JumlahKlaim.Currency != "IDR" {
		t.Errorf("muatan bermata uang %q, mau IDR", m.JumlahKlaim.Currency)
	}
	if utils.FormatDecimal(m.JumlahKlaim.Amount) != "1500000" {
		t.Errorf("muatan berjumlah %q, mau 1500000",
			utils.FormatDecimal(m.JumlahKlaim.Amount))
	}
	// Ambang yang dicarikan roster memakai mata uang yang sama.
	if len(roster.ambang) != 1 || roster.ambang[0].Currency != "IDR" {
		t.Errorf("ambang roster = %+v, mau satu ambang IDR", roster.ambang)
	}

	// Penautan benar-benar tersimpan, dan dibaca ULANG dari Oracle - bukan
	// dipercaya dari nilai yang baru saja ditulis.
	_, sesudah := barisPertama(t, db, pohon.Work.ID)
	if sesudah.KomiteID != "KMT-UJI-711" {
		t.Errorf("KOMITE_ID tersimpan = %q, mau %q", sesudah.KomiteID, "KMT-UJI-711")
	}
	if len(jejak.catatan) != 1 {
		t.Errorf("jejak = %d catatan, mau 1", len(jejak.catatan))
	}

	// ⛔ Penyerahan kedua atas baris yang sama DITOLAK - dan ditolak oleh
	// keadaan yang tersimpan, bukan oleh ingatan proses ini.
	err = svc.Komite().DenganRoster(roster).DenganKasus(kasus).DenganJejak(jejak).
		Serahkan(ctx, services.Pelaku{AkunID: "UJI-AKUN",
			Peran: []string{services.PeranAdmin}},
			pohon.Work.ID, pesertaID, baris.ID, saatUjiKomite)
	if !errors.Is(err, services.ErrBarisSudahDiserahkan) {
		t.Errorf("penyerahan kedua: galat = %v, mau ErrBarisSudahDiserahkan", err)
	}
}

// TestPenyerahanGagalTidakMeninggalkanPenautanSeparuh - ADR-U-0029 Akibat 2.
func TestPenyerahanGagalTidakMeninggalkanPenautanSeparuh(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()

	pohon := pohonUjiKomite(t, svc, db, "CLM-UJI712", "UJI-BANK", "UJI-006", "0012345")
	pesertaID, baris := barisPertama(t, db, pohon.Work.ID)

	// Kasus komite lahir dan baris tertaut, LALU jejaknya gagal. Bila
	// transaksinya tidak utuh, KOMITE_ID akan tertinggal terisi.
	err := svc.Komite().DenganRoster(&rosterUji{tingkat: 1}).
		DenganKasus(&kasusUji{id: "KMT-UJI-712"}).DenganJejak(jejakGagal{}).
		Serahkan(context.Background(), services.Pelaku{AkunID: "UJI-AKUN",
			Peran: []string{services.PeranAdmin}},
			pohon.Work.ID, pesertaID, baris.ID, saatUjiKomite)
	if !errors.Is(err, errJejakSengaja) {
		t.Fatalf("galat = %v, mau errJejakSengaja", err)
	}

	_, sesudah := barisPertama(t, db, pohon.Work.ID)
	if sesudah.KomiteID != "" {
		t.Errorf("KOMITE_ID tertinggal %q sesudah transaksi gagal; "+
			"penautan separuh jadi", sesudah.KomiteID)
	}
}
