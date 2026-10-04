//go:build db

package services_test

// GILIRAN-18 paket 2 terhadap skema uji Oracle - OQ-N13.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus.

import (
	"context"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// OQ-N13: sesudah Save to RNM klaim A, cermin A berstatus '0' - dan klaim
// ganda B atas tertanggung yang sama TERTANGKAP tanpa menunggu Komite.
func TestSaveRNMMenyetelCerminSehinggaKlaimGandaTertangkap(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)

	// Polis tiruan: CEDINGCO cermin dan pemeriksa ganda datang dari sini.
	for _, tabel := range []string{"T_WORK_POLIS", "T_PREMIUM_LIST"} {
		nama, err := db.Qualify(tabel)
		if err != nil {
			t.Fatal(err)
		}
		// ⚠️ Pengenal berbeda dari uji db PremiumList (`UJI-POLIS-1801`) -
		// temuan /code-review: dua paket memakai skema uji yang sama.
		q := `INSERT INTO ` + nama + ` (ID) VALUES ('UJI-POLIS-1811')`
		if tabel == "T_PREMIUM_LIST" {
			q = `INSERT INTO ` + nama + ` (ID, NO_POLIS, CEDING_CO, PROD_KE)
			     VALUES ('UJI-POLIS-1811', 'UJI-POL-0001', 'UJI-CEDING-1', 0)`
		}
		if _, err := db.ExecContext(ctx, q); err != nil {
			t.Fatalf("polis tiruan: %v", err)
		}
	}
	lama, err := db.Qualify("OS_AKSEPTASI_KLAIM_LIFE")
	if err != nil {
		t.Fatal(err)
	}
	daftar := func(id string) string {
		t.Helper()
		p := models.PohonKlaim{
			Work: models.WorkClaim{ID: id, Lini: inti.LiniLife},
			Klaim: models.Klaim{
				NomorKlaim: "UJI-" + id, NomorPolis: "UJI-POL-0001", Type: "QP",
				Peserta: []models.Peserta{{
					NomorPremiList: "UJI-PL-1", NomorSertifikat: "006", SumberID: "UJI-SRC-1", MataUang: "IDR",
					Baris: []models.BarisAdjustment{{}},
				}},
			},
		}
		if err := svc.DalamTransaksi(ctx, func(tx *intidb.Tx) error {
			return repository.NewPohonKlaim(db).Simpan(ctx, tx, p)
		}); err != nil {
			t.Fatalf("mendaftar %s: %v", id, err)
		}
		var adj string
		if err := db.QueryRowContext(ctx, `SELECT ID FROM `+lama+` WHERE CASEID = :1`, id).Scan(&adj); err != nil {
			t.Fatalf("baris cermin %s: %v", id, err)
		}
		return adj
	}
	statusCermin := func(adj string) string {
		t.Helper()
		var s *string
		if err := db.QueryRowContext(ctx, `SELECT TO_CHAR(STS_REJECT) FROM `+lama+` WHERE ID = :1`, adj).Scan(&s); err != nil {
			t.Fatal(err)
		}
		if s == nil {
			return "<NULL>"
		}
		return *s
	}
	k := repository.KunciPesertaSumber{PLNumber: "UJI-PL-1", Sertifikat: "006", SumberID: "UJI-SRC-1"}

	adjA := daftar("CLM-UJI1801")
	adjB := daftar("CLM-UJI1802")
	// Sebelum Save to RNM: cermin A NULL, jadi B belum melihatnya.
	if s, err := baca.StatusWarisanTerakhir(ctx, k, "UJI-CEDING-1", "CLM-UJI1802"); err != nil || s != "" {
		t.Errorf("sebelum Save to RNM A: status terlihat B = %q, %v; mau kosong", s, err)
	}

	// Save to RNM A - langkah 22.1.3.2 + b176, satu transaksi.
	var disetel int
	if err := svc.DalamTransaksi(ctx, func(tx *intidb.Tx) error {
		if err := baca.TandaiBarisOutstanding(ctx, tx, adjA); err != nil {
			return err
		}
		n, err := baca.SetelCerminOutstanding(ctx, tx, adjA, "CLM-UJI1801")
		disetel = n
		return err
	}); err != nil {
		t.Fatalf("Save to RNM A: %v", err)
	}
	if disetel != 1 || statusCermin(adjA) != "0" {
		t.Errorf("cermin A: disetel %d, status %s; mau 1, '0'", disetel, statusCermin(adjA))
	}
	// Klaim ganda B kini TERTANGKAP; A sendiri tetap dikecualikan.
	if s, err := baca.StatusWarisanTerakhir(ctx, k, "UJI-CEDING-1", "CLM-UJI1802"); err != nil || s != "0" {
		t.Errorf("sesudah Save to RNM A: status terlihat B = %q, %v; mau '0'", s, err)
	}
	if s, err := baca.StatusWarisanTerakhir(ctx, k, "UJI-CEDING-1", "CLM-UJI1801"); err != nil || s != "" {
		t.Errorf("klaim A melihat cerminnya sendiri: %q, %v", s, err)
	}

	// Temuan /code-review: Admin MENOLAK baris A (RejectOSClaimLife_Act langkah
	// 5) - cermin mengikuti '2', dan klaim B tidak lagi terblokir (11.4 `==2`
	// lewati). Tanpa penyelarasan, cermin tertahan '0' selamanya.
	pesertaA, err := baca.AmbilPeserta(ctx, "CLM-UJI1801")
	if err != nil || len(pesertaA) != 1 {
		t.Fatalf("peserta A: %v (%d)", err, len(pesertaA))
	}
	if err := svc.DalamTransaksi(ctx, func(tx *intidb.Tx) error {
		return baca.PerbaruiStatusBaris(ctx, tx, pesertaA[0].ID, adjA, "0", "2", "", time.Time{})
	}); err != nil {
		t.Fatalf("menolak baris A: %v", err)
	}
	if statusCermin(adjA) != "2" {
		t.Errorf("cermin A sesudah ditolak: %s, mau '2'", statusCermin(adjA))
	}
	if s, err := baca.StatusWarisanTerakhir(ctx, k, "UJI-CEDING-1", "CLM-UJI1802"); err != nil || s != "2" {
		t.Errorf("sesudah A ditolak: status terlihat B = %q, %v; mau '2' (lolos 11.4)", s, err)
	}

	// Baris yang SUDAH berkeputusan tidak ditimpa, dan CASEID lain tidak tersentuh.
	if _, err := db.ExecContext(ctx, `UPDATE `+lama+` SET STS_REJECT = '1' WHERE ID = :1`, adjB); err != nil {
		t.Fatal(err)
	}
	if err := svc.DalamTransaksi(ctx, func(tx *intidb.Tx) error {
		for _, u := range [][2]string{{adjB, "CLM-UJI1802"}, {adjA, "CLM-UJI1802"}} {
			n, err := baca.SetelCerminOutstanding(ctx, tx, u[0], u[1])
			if err != nil {
				return err
			}
			if n != 0 {
				t.Errorf("SetelCerminOutstanding(%s, %s) menyentuh %d baris, mau 0", u[0], u[1], n)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if statusCermin(adjB) != "1" {
		t.Errorf("cermin berkeputusan '1' ditimpa: %s", statusCermin(adjB))
	}
}
