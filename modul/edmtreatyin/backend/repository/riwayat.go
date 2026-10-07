package repository

// Untuk apa berkas ini: RIWAYAT, HARI CLOSING, dan NAMA TAMPILAN.
//
//   - `HISTORYAKSEPTASIPEGA` - `Activity/InsertHistoryAkseptasiPega` (RDB
//     `InsertHistoryAkseptasiPega_Sql`), pasca-submit KEDUA flow action. Pega
//     menulisnya tanpa skema dan dengan COMMIT di dalam blok; di sini skema
//     eksplisit (AC 30, 90) dan di transaksi submit - gagal menulis riwayat
//     membatalkan seluruh submit (AC 83).
//   - Hari closing - RDB `GETTanggalClosing_SQL` lewat penomor bersama
//     `inti/backend/penomor` (BACA saja). EDM TIDAK menerbitkan nomor polis:
//     generasi endorsemen memakai nomor polis induk + EDMNo
//     (`models.NomorEDM`, `Activity/SetEDMTNoPolis` langkah 3).
//   - `M_LOGIN_GO.NAME` - nama tampilan (`OperatorID.pyUserName`), dibaca saja
//     seperti modul marketingofficer.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// CatatRiwayat menulis satu baris riwayat di transaksi submit (AC 43).
// `ID_KOMITE` (InsertHistory.CARI6) tidak pernah diisi rule ini - NULL.
func (g *Gudang) CatatRiwayat(ctx context.Context, tx *db.Tx, r models.Riwayat) error {
	t, err := g.nama(tabelRiwayatPega)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`INSERT INTO %s (ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET, ID_KOMITE, OPERATORID)
		VALUES (:1, SYSDATE, :2, :3, :4, NULL, :5)`, t)
	hasil, err := jalankan(ctx, tx, "menulis riwayat", q, r.IDPega, db.KosongJadiNil(r.Status),
		db.KosongJadiNil(r.Username), db.KosongJadiNil(r.Workbasket), r.OperatorID)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, "penulisan riwayat")
}

// ------------------------------------------------------------------ nomor polis

// HariClosing = RDB `GETTanggalClosing_SQL` (`POOLDATA.TANGGAL_CLOSING`), lewat
// penomor bersama.
func (g *Gudang) HariClosing(ctx context.Context, tx *db.Tx) (int, error) {
	return g.nomor.HariClosing(ctx, tx)
}

// ------------------------------------------------------------------ nama tampilan

// NamaTampilan membaca `M_LOGIN_GO.NAME` satu akun - padanan
// `OperatorID.pyUserName`.
//
// ⛔ Tanpa jatuh-balik ke pengenal akun: `OperatorName`/`USERNAME` berarti
// NAMA TAMPILAN saja (AC 40, 42; P33). Akun tanpa nama menghasilkan teks
// kosong, bukan login ID.
// Tabel login inti yang dibaca untuk nama kotak masuk NBStatus (`PemegangKotakMasuk`).
const (
	tabelLoginWorkbasket = "M_LOGIN_GO_WORKBASKET"
	tabelWorkbasket      = "M_WORKBASKET"
	// benderaAkunAktif - M_LOGIN_GO.IS_ACTIVE akun aktif (bendera '1'/'0' tabel login inti).
	benderaAkunAktif = "1"
	// divisiTidakDihitung - M_LOGIN_GO.DIVISION_CODE akun yang TIDAK dihitung sebagai pemegang workbasket
	// (M_DIVISION "IT" = Information Technology; keputusan work owner 06-10-2026: "kalo yg IT jangan ikut dihitung").
	divisiTidakDihitung = "IT"
)

// sqlPemegangKotakMasuk - satu baris agregat dari master workbasket :3: jumlah akun AKTIF (:1) di luar divisi IT
// (:2; divisi kosong tetap dihitung) pemegangnya, nama
// terbesar di antaranya (= satu-satunya nama bila jumlahnya 1), dan nama workbasket. LEFT JOIN, bukan subquery
// skalar di samping COUNT/MAX (ORA-00937 di DEV 06-10-2026); workbasket tak terdaftar = 0 pemegang tanpa nama.
// Penampung ditulis menurut urutan kemunculan = urutan argumen.
func sqlPemegangKotakMasuk(loginWB, login, wb string) string {
	return fmt.Sprintf(`SELECT COUNT(g.LOGIN_ID), MAX(g.NAME), MAX(w.NAME)
	  FROM %s w
	  LEFT JOIN %s l ON l.WORKBASKET_ID = w.WORKBASKET_ID
	  LEFT JOIN %s g ON g.LOGIN_ID = l.LOGIN_ID AND g.IS_ACTIVE = :1
	        AND (g.DIVISION_CODE IS NULL OR g.DIVISION_CODE <> :2)
	 WHERE w.WORKBASKET_ID = :3`, wb, loginWB, login)
}

// PemegangKotakMasuk membaca pemegang aktif workbasket tujuan Submit (NBStatus, `[keputusan work owner
// 06-10-2026]`).
func (g *Gudang) PemegangKotakMasuk(ctx context.Context, workbasket string) (models.PemegangKotakMasuk, error) {
	lwb, err := g.nama(tabelLoginWorkbasket)
	if err != nil {
		return models.PemegangKotakMasuk{}, err
	}
	akun, err := g.nama(tabelAkunLogin)
	if err != nil {
		return models.PemegangKotakMasuk{}, err
	}
	wb, err := g.nama(tabelWorkbasket)
	if err != nil {
		return models.PemegangKotakMasuk{}, err
	}
	q := sqlPemegangKotakMasuk(lwb, akun, wb)
	if err := db.PeriksaSQL(q); err != nil {
		return models.PemegangKotakMasuk{}, err
	}
	var n int
	var namaAkun, namaWB sql.NullString
	if err := g.db.QueryRowContext(ctx, q, benderaAkunAktif, divisiTidakDihitung, workbasket).Scan(&n, &namaAkun, &namaWB); err != nil {
		return models.PemegangKotakMasuk{}, fmt.Errorf("repository: membaca pemegang workbasket: %w", err)
	}
	return models.PemegangKotakMasuk{Jumlah: n, NamaAkun: namaAkun.String, NamaWorkbasket: namaWB.String}, nil
}

func (g *Gudang) NamaTampilan(ctx context.Context, loginID string) (string, error) {
	if loginID == "" {
		return "", nil
	}
	t, err := g.nama(tabelAkunLogin)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca nama tampilan", fmt.Sprintf(`SELECT NAME FROM %s WHERE LOGIN_ID = :1`, t), loginID)
}
