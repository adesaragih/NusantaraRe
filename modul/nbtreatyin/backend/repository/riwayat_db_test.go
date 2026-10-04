//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN untuk RIWAYAT AKSEPTASI
// `HISTORYAKSEPTASIPEGA` (`Activity/InsertHistoryAkseptasiPega` ->
// `RDBList/InsertHistoryAkseptasiPega_Sql`): urutan baca (spec AC 72, tiket 10)
// dan kegagalan menulis riwayat yang membatalkan seluruh submit (spec AC 83,
// tiket 08). ⛔ Belum pernah dijalankan: skema uji K11 kosong (PERMINTAAN C4).
//
// Tabel itu WARISAN `POOLDATA`, tidak dibuat migrasi modul ini. Bila skema uji
// tidak memuatnya, uji membuat TIRUAN (`buatTiruan`, keputusan WO U1) tujuh kolom yang ditulis
// `CatatRiwayat` (`ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET,
// ID_KOMITE` = INSERT rule XML + `OPERATORID`, kolom yang ada di katalog -
// PROMPT putaran 2 bab 1) dan membuangnya lagi. `TGL_TRANSFER` = `sysdate` di
// rule ⇒ DATE; panjang kolom teks tidak diketahui (PERMINTAAN C8 untuk tabel
// saudaranya) ⇒ VARCHAR2(1000) seperti kolom view warisan. Pemanggil
// `skemauji.Buka()` tetap satu (`pasang`). Fixture berawalan UJI-.

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya di skema uji):
// riwayat akseptasi warisan `POOLDATA.HISTORYAKSEPTASIPEGA`.
const tabelRiwayatPega = "HISTORYAKSEPTASIPEGA"

// siapkanRiwayatPega memastikan tabel riwayat ada di skema uji: tabel yang
// disediakan DBA dipakai apa adanya (baris UJI- dihapus sesudah uji), selain
// itu tiruan dibuat lewat `buatTiruan` (U1: berpagar POOLDATA, dibuang
// sesudah uji).
func siapkanRiwayatPega(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	kolom := []string{"ID_PEGA VARCHAR2(1000)", "TGL_TRANSFER DATE", "STATUS VARCHAR2(1000)", "USERNAME VARCHAR2(1000)",
		"WORKBASKET VARCHAR2(1000)", "ID_KOMITE VARCHAR2(1000)", "OPERATORID VARCHAR2(1000)"}
	if !buatTiruan(t, ctx, sqlDB, skema, tabelRiwayatPega, kolom) {
		t.Cleanup(func() {
			_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s.%s WHERE ID_PEGA LIKE '%%UJI-NB-%%'`, skema, tabelRiwayatPega))
		})
	}
}

// AC 72: "Riwayat dapat dibaca berurutan waktu. Test yang menemukan urutan
// acak gagal." Tiga baris ditulis `CatatRiwayat` (TGL_TRANSFER = SYSDATE, satu
// detik yang sama), lalu waktunya digeser sehingga urutan TULIS berlawanan
// dengan urutan WAKTU - pembacaan berurut ROWID/sisip saja akan gagal.
func TestDaftarRiwayatBerurutWaktu(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	siapkanRiwayatPega(t, ctx, sqlDB, skema)
	g := repository.Baru(d)
	idPega := models.KunciInstans("UJI-NB-RIW-1")
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		for _, s := range []string{"UJI-PERTAMA", "UJI-KEDUA", "UJI-KETIGA"} {
			if err := g.CatatRiwayat(ctx, tx, models.Riwayat{IDPega: idPega, Status: s, Username: "UJI NAMA",
				Workbasket: models.PosisiAdmin, OperatorID: "UJI-AKUN"}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for status, waktu := range map[string]string{
		"UJI-PERTAMA": "2026-10-03 10:00:00", // ditulis pertama, terakhir menurut waktu
		"UJI-KEDUA":   "2026-10-01 09:00:00",
		"UJI-KETIGA":  "2026-10-02 08:30:00",
	} {
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s.%s SET TGL_TRANSFER = TO_DATE(:1, 'YYYY-MM-DD HH24:MI:SS')
			WHERE ID_PEGA = :2 AND STATUS = :3`, skema, tabelRiwayatPega), waktu, idPega, status); err != nil {
			t.Fatal(err)
		}
	}
	r, err := g.DaftarRiwayat(ctx, idPega)
	if err != nil {
		t.Fatal(err)
	}
	var urut []string
	for _, b := range r {
		urut = append(urut, b.Status+"@"+b.TglTransfer)
	}
	harap := "[UJI-KEDUA@2026-10-01 09:00:00 UJI-KETIGA@2026-10-02 08:30:00 UJI-PERTAMA@2026-10-03 10:00:00]"
	if fmt.Sprint(urut) != harap {
		t.Fatalf("urutan riwayat %v, harap %s", urut, harap)
	}
	if r[0].OperatorID != "UJI-AKUN" || r[0].Username != "UJI NAMA" || r[0].Workbasket != models.PosisiAdmin {
		t.Errorf("isi baris riwayat %+v", r[0])
	}
}

// AC 83: "Kegagalan menyimpan riwayat membatalkan seluruh transaksi. Test yang
// menemukan berkas tersimpan tanpa riwayat gagal."
//
// Kegagalan disuntikkan di ORACLE, bukan di Go: constraint CHECK sementara
// menolak baris riwayat ber-OPERATORID `UJI-TOLAK` (ORA-02290). Urutan tulis
// submit admin yang disetujui (`services.Kirim`: halaman, posisi, riwayat -
// riwayat diletakkan TERAKHIR supaya penulisan halaman dan posisi sudah
// terjadi ketika riwayat gagal) dijalankan dalam satu transaksi, lalu keadaan
// dibaca ulang: tidak satu pun perubahan boleh tersisa.
func TestGagalCatatRiwayatMembatalkanSubmit(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	siapkanRiwayatPega(t, ctx, sqlDB, skema)
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s.%s ADD CONSTRAINT UJI_NB_TOLAK_RIWAYAT
		CHECK (OPERATORID IS NULL OR OPERATORID <> 'UJI-TOLAK') ENABLE NOVALIDATE`, skema, tabelRiwayatPega)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`ALTER TABLE %s.%s DROP CONSTRAINT UJI_NB_TOLAK_RIWAYAT`, skema, tabelRiwayatPega))
	})
	g := repository.Baru(d)
	const id = "UJI-NB-RIW-GAGAL"
	awal := models.HalamanBaru()
	awal.Setel("PositionNote", models.PosisiAdmin)
	awal.Setel("NBStatus", "UJI-AWAL")
	awal.Setel("PolicyTreatyIn.PremiOgp", "100")
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, awal)
	}); err != nil {
		t.Fatal(err)
	}
	sebelum, err := g.Keadaan(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}

	kirim := models.HalamanBaru()
	kirim.Setel("PositionNote", models.PosisiSecHead)
	kirim.Setel("NBStatus", models.TeksNBStatusKotakMasuk(models.PosisiSecHead))
	kirim.Setel("PolicyTreatyIn.PremiOgp", "999")
	kirim.Setel("PolicyTreatyIn.IsApproved", "1")
	err = dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SimpanHalaman(ctx, tx, id, kirim); err != nil {
			return err
		}
		if err := g.PindahPosisi(ctx, tx, id, sebelum.StatusWork, models.PosisiSecHead); err != nil {
			return err
		}
		return g.CatatRiwayat(ctx, tx, models.Riwayat{IDPega: models.KunciInstans(id), Status: models.StatusRiwayat("1"),
			Username: "UJI NAMA", Workbasket: models.PosisiAdmin, OperatorID: "UJI-TOLAK"})
	})
	if err == nil {
		t.Fatal("riwayat yang ditolak Oracle harus menggagalkan submit")
	}

	sesudah, err := g.Keadaan(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	if sesudah.PositionNote != models.PosisiAdmin || sesudah.StatusWork != sebelum.StatusWork || sesudah.Position != sebelum.Position {
		t.Errorf("posisi tersimpan tanpa riwayat: sebelum %+v, sesudah %+v", sebelum, sesudah)
	}
	h, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	if h.Ambil("PolicyTreatyIn.PremiOgp") != "100" || h.Ambil("NBStatus") != "UJI-AWAL" || h.Ambil("PolicyTreatyIn.IsApproved") != "" {
		t.Errorf("halaman tersimpan tanpa riwayat: PremiOgp %q NBStatus %q IsApproved %q",
			h.Ambil("PolicyTreatyIn.PremiOgp"), h.Ambil("NBStatus"), h.Ambil("PolicyTreatyIn.IsApproved"))
	}
	if r, err := g.DaftarRiwayat(ctx, models.KunciInstans(id)); err != nil || len(r) != 0 {
		t.Errorf("riwayat tersisa %+v %v", r, err)
	}
}
