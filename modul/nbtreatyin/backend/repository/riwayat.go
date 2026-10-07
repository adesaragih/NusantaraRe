package repository

// Untuk apa berkas ini: RIWAYAT, NOMOR POLIS, dan NAMA TAMPILAN.
//
//   - `HISTORYAKSEPTASIPEGA` - `Activity/InsertHistoryAkseptasiPega` (RDB
//     `InsertHistoryAkseptasiPega_Sql`), pasca-submit KEDUA flow action. Pega
//     menulisnya tanpa skema dan dengan COMMIT di dalam blok; di sini skema
//     eksplisit (AC 30, 90) dan di transaksi submit - gagal menulis riwayat
//     membatalkan seluruh submit (AC 83).
//   - Nomor polis - `GeneratePolicyNoTreaty_Act` lewat penomor bersama
//     `inti/backend/penomor` (padanan `PROC_GENERATE_SEQUENCE_NUMBER`).
//   - `M_LOGIN_GO.NAME` - nama tampilan (`OperatorID.pyUserName`), dibaca saja
//     seperti modul marketingofficer.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/nbtreatyin/backend/models"
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

// DaftarRiwayat membaca riwayat satu kasus, berurut waktu (AC 72).
func (g *Gudang) DaftarRiwayat(ctx context.Context, idPega string) ([]models.Riwayat, error) {
	t, err := g.nama(tabelRiwayatPega)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID_PEGA, STATUS, USERNAME, WORKBASKET, OPERATORID, TO_CHAR(TGL_TRANSFER, '%s')
	  FROM %s WHERE ID_PEGA = :1 ORDER BY TGL_TRANSFER, ROWID`, fmtTanggal, t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, idPega)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca riwayat: %w", err)
	}
	defer rows.Close()
	var out []models.Riwayat
	for rows.Next() {
		var v [6]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5]); err != nil {
			return nil, fmt.Errorf("repository: membaca riwayat: %w", err)
		}
		out = append(out, models.Riwayat{IDPega: v[0].String, Status: v[1].String, Username: v[2].String,
			Workbasket: v[3].String, OperatorID: v[4].String, TglTransfer: v[5].String})
	}
	return out, rows.Err()
}

// ------------------------------------------------------------------ nomor polis

// BahanNomor - hasil `GeneratePolicyNoTreaty_Act` langkah 5-26.
type BahanNomor struct {
	NoPolis        string
	ProductionDate time.Time
}

// TerbitkanNomorPolis = `GeneratePolicyNoTreaty_Act` langkah efektif
// (langkah 14-24 berlabel `//`):
//
//	5    tanggal produksi: sekarang / StatementDate di depan / geser closing
//	     (RDB `GETTanggalClosing_SQL`)
//	7-9  tipe QR/QP/TP (`models.TipeNomorPolis`)
//	25   awalan `KODE_PRODUKSI WHERE TYPE='NONLIFE'` (`GetKodeProdNonLife_SQL`)
//	26   urut dari `GENERATE_SEQUENCE_NUMBER` (CLASS pxObjClass, JENIS
//	     awalan+"QR/QP/TP", tanggal produksi) - di-lookup penomor bersama
//	28   PolicyNo = awalan + tipe + ".T" + OJKBusinessID + "." + MM.YYYY + "." + urut5
//
// ⛔ Tipe kosong (DueTo bukan 1/0 dan bukan XOL Retro) DITOLAK - Pega
// menerbitkan nomor tanpa huruf tipe yang tampak sah (pola penolakan
// `penomor.ErrTipePLTanpaCabang`).
func (g *Gudang) TerbitkanNomorPolis(ctx context.Context, tx *db.Tx, h *models.Halaman, sekarang time.Time) (BahanNomor, error) {
	tipe := models.TipeNomorPolis(h)
	if tipe == "" {
		return BahanNomor{}, ErrTipeNomorKosong
	}
	ojk := strings.TrimSpace(h.Ambil(models.HalamanPolis + ".OJKBusinessID"))
	if ojk == "" {
		return BahanNomor{}, ErrOJKKosong
	}
	hari, err := g.nomor.HariClosing(ctx, tx)
	if err != nil {
		return BahanNomor{}, err
	}
	var statement time.Time
	if s := strings.TrimSpace(h.Ambil(models.HalamanPolis + ".StatementDate")); s != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, sekarang.Location()); err == nil {
			statement = t
		}
	}
	prod := models.TanggalProduksiNomor(sekarang, statement, hari)
	periode, err := penomor.HitungPeriodeNomor(prod, hari)
	if err != nil {
		return BahanNomor{}, err
	}
	awalan, err := g.nomor.AwalanProduksi(ctx, tx, models.LiniKasus)
	if err != nil {
		return BahanNomor{}, err
	}
	urut, err := g.nomor.UrutNomorBerikut(ctx, tx, models.KelasDeret, models.JenisDeret(awalan), periode, prod)
	if err != nil {
		return BahanNomor{}, err
	}
	return BahanNomor{
		NoPolis:        models.RakitNomorPolis(awalan, tipe, ojk, periode.MMYYYY, urut),
		ProductionDate: prod,
	}, nil
}

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
