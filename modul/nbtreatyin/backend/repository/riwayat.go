package repository

// Untuk apa berkas ini: RIWAYAT, NOMOR POLIS, PERAN TEMPAT, dan NAMA TAMPILAN.
//
//   - `HISTORYAKSEPTASIPEGA` - `Activity/InsertHistoryAkseptasiPega` (RDB
//     `InsertHistoryAkseptasiPega_Sql`), pasca-submit KEDUA flow action. Pega
//     menulisnya tanpa skema dan dengan COMMIT di dalam blok; di sini skema
//     eksplisit (AC 30, 90) dan di transaksi submit - gagal menulis riwayat
//     membatalkan seluruh submit (AC 83).
//   - Nomor polis - `GeneratePolicyNoTreaty_Act` lewat penomor bersama
//     `inti/backend/penomor` (padanan `PROC_GENERATE_SEQUENCE_NUMBER`).
//   - `M_NBTRIN_PERAN_TEMPAT` - pemetaan tempat -> peran (tiket 05).
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

// ------------------------------------------------------------------ peran tempat

// PeranTempat adalah satu baris M_NBTRIN_PERAN_TEMPAT.
type PeranTempat struct {
	KodeTempat string
	Peran      string
	// Arah - MUNCUL (hanya untuk peran ini) atau KECUALI (semua kecuali
	// peran ini) - TIDAK ditebak (AC 82).
	Arah string
}

// DaftarPeranTempat membaca seluruh pemetaan. Tabel kosong = setiap tempat
// tertunda (AC 81).
func (g *Gudang) DaftarPeranTempat(ctx context.Context) ([]PeranTempat, error) {
	t, err := g.nama(tabelPeranTempat)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT KODE_TEMPAT, PERAN, ARAH FROM %s ORDER BY KODE_TEMPAT, PERAN`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peran tempat: %w", err)
	}
	defer rows.Close()
	var out []PeranTempat
	for rows.Next() {
		var p PeranTempat
		if err := rows.Scan(&p.KodeTempat, &p.Peran, &p.Arah); err != nil {
			return nil, fmt.Errorf("repository: membaca peran tempat: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ------------------------------------------------------------------ nama tampilan

// NamaTampilan membaca `M_LOGIN_GO.NAME` satu akun - padanan
// `OperatorID.pyUserName`.
//
// ⛔ Tanpa jatuh-balik ke pengenal akun: `OperatorName`/`USERNAME` berarti
// NAMA TAMPILAN saja (AC 40, 42; P33). Akun tanpa nama menghasilkan teks
// kosong, bukan login ID.
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
