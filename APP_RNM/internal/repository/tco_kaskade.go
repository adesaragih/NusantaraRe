package repository

// Kaskade hapus kontrak dan reinsurer - tiket 10 Treaty Contract Out.
//
// Untuk apa berkas ini: hitung dampak dan jalankan kaskade yang ditiru dari
// `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` b80-b93 (empat DELETE berurut)
// dan `RDBList/DeleteFromTreatyReinsurer_Act.xml` b60-b64 (security lalu
// reinsurer), di tabel `T_` - SATU keluarga tabel (AC 63/64).
//
// ⛔ `PROPORTIONALARRG` (klausul) TIDAK disentuh (AC 44; penyimpangan sadar
// 4): klausul milik tahun/grup/jenis, dipakai lintas kontrak. Ia hanya
// DIHITUNG, supaya popup dapat menyatakannya tetap hidup.
//
// ⛔ Bisnis dicocokkan `(TREATYYEARID = :x OR TREATYYEARID IS NULL)` VERBATIM -
// baris lama tanpa TREATYYEARID tidak boleh tertinggal.
//
// ⛔ Hitungan dan penghapusan memakai saringan yang SAMA (fungsi yang sama
// merakit WHERE-nya), sehingga angka popup = angka yang terhapus.
//
// Dibaca sesudah: tco_business.go, tco_security.go.

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/internal/models"
)

// DampakHapusTCO adalah jumlah baris tiap jenis yang ikut terhapus - dan
// jumlah klausul yang TIDAK terhapus.
type DampakHapusTCO struct {
	Kontrak, Reinsurer, Security, Business, KlausulTetap int64
	// Bersama - kontrak LAIN (tahun lain, teks tahun + grup sama) yang memakai
	// kombinasi yang sama. Reinsurer/security mereka IKUT terhapus, seperti
	// Pega [keputusan work owner 29-09-2026, OQ-TCO-21]; business hanya yang
	// `TREATYYEARID` tahun ini atau kosong (`saringBusinessKaskadeTCO`). Cacahnya
	// wajib disebut popup dan dikonfirmasi pemakai.
	Bersama int64
}

// saringan WHERE bersama hitung & hapus.
const (
	saringBusinessKaskadeTCO  = `TREATYYEAR = :1 AND (TREATYYEARID = :2 OR TREATYYEARID IS NULL) AND TREATYGROUPID = :3 AND REINSTYPEID = :4`
	saringReinsurerKaskadeTCO = `TREATYYEAR = :1 AND TREATYGROUPID = :2 AND REINSTYPEID = :3`
	saringSecurityKaskadeTCO  = `REAS_ID IN (SELECT ID FROM %s WHERE TREATYYEAR = :1 AND TREATYGROUPID = :2 AND REINSTYPEID = :3)`
	// Klausul milik kontrak itu, dihitung DARI INDUKNYA [keputusan work owner
	// 29-09-2026, OQ-TCO-20]: baris induk (`PARENTREINSTYPEID` = sentinel "00")
	// berjenis reasuransi kontrak, dan baris anak yang `PARENTREINSTYPEID`-nya
	// jenis kontrak - bukan anak yang kebetulan ber-REINSTYPEID sama di bawah
	// induk lain. Hanya DIHITUNG; klausul tidak dihapus (Pega, AC 44).
	saringKlausulTetapTCO = `TREATYYEARID = :1 AND ((PARENTREINSTYPEID = :2 AND REINSTYPEID = :3) OR PARENTREINSTYPEID = :4)`
)

// sqlKontrakBersamaTCO - kontrak lain yang kombinasinya sama.
func sqlKontrakBersamaTCO(kontrak, tahun string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s c, %s y WHERE y.ID = c.IDTREATYYEAR
	   AND y.TREATYYEAR = :1 AND y.TREATYGROUPID = :2 AND c.REINSTYPEID = :3 AND c.ID <> :4`, kontrak, tahun)
}

func sqlHitungTCO(tabel, saring string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, tabel, saring)
}

func sqlHapusKaskadeTCO(tabel, saring string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s`, tabel, saring)
}

func sqlHapusKontrakKaskadeTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND IDTREATYYEAR = :2`, tabel)
}

func sqlHapusSecurityReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE REAS_ID = :1`, tabel)
}

func sqlHapusReinsurerSatuTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND %s`, tabel,
		`TREATYYEAR = :2 AND TREATYGROUPID = :3 AND REINSTYPEID = :4`)
}

func sqlHitungSecurityReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE REAS_ID = :1`, tabel)
}

// ErrKaskadeTidakUtuh - kontrak/reinsurer yang dihapus tidak tepat satu baris.
var ErrKaskadeTidakUtuh = errors.New("repository: kaskade hapus tidak mengenai tepat satu baris induk")

// KaskadeTCO menghitung dan menjalankan kaskade hapus.
type KaskadeTCO struct{ db *DB }

// NewKaskadeTCO menyusunnya.
func NewKaskadeTCO(db *DB) *KaskadeTCO { return &KaskadeTCO{db: db} }

type tabelKaskadeTCO struct{ kontrak, reinsurer, security, business, klausul, tahun string }

func (k *KaskadeTCO) tabel() (tabelKaskadeTCO, error) {
	var t tabelKaskadeTCO
	for _, p := range []struct {
		logis string
		isi   *string
	}{{TabelKontrakTCO, &t.kontrak}, {TabelReinsurerTCO, &t.reinsurer}, {TabelSecurityTCO, &t.security},
		{TabelBusinessTCO, &t.business}, {TabelKlausulTCO, &t.klausul}, {TabelTahunTCO, &t.tahun}} {
		q, err := k.db.Qualify(p.logis)
		if err != nil {
			return t, err
		}
		*p.isi = q
	}
	return t, nil
}

func (k *KaskadeTCO) hitung(ctx context.Context, q string, args ...any) (int64, error) {
	if err := PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int64
	if err := k.db.bacaTCO(ctx).QueryRowContext(ctx, q, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: menghitung dampak hapus: %w", err)
	}
	return n, nil
}

// DampakKontrak menghitung baris yang ikut terhapus bersama kontrak itu.
func (k *KaskadeTCO) DampakKontrak(ctx context.Context, kom models.KombinasiTCO, tahunID, kontrakID string) (DampakHapusTCO, error) {
	t, err := k.tabel()
	if err != nil {
		return DampakHapusTCO{}, err
	}
	d := DampakHapusTCO{Kontrak: 1}
	if d.Bersama, err = k.hitung(ctx, sqlKontrakBersamaTCO(t.kontrak, t.tahun),
		kom.TreatyYear, kom.TreatyGroupID, kom.ReinsTypeID, kontrakID); err != nil {
		return DampakHapusTCO{}, err
	}
	for _, l := range langkahHapusKontrakTCO(t, kom, tahunID, kontrakID, &d) {
		if l.hitung == "" {
			continue
		}
		if *l.isi, err = k.hitung(ctx, l.hitung, l.args...); err != nil {
			return DampakHapusTCO{}, err
		}
	}
	if d.KlausulTetap, err = k.hitung(ctx, sqlHitungTCO(t.klausul, saringKlausulTetapTCO),
		tahunID, models.ParentReinsTypeTanpaInduk, kom.ReinsTypeID, kom.ReinsTypeID); err != nil {
		return DampakHapusTCO{}, err
	}
	return d, nil
}

// HapusKontrak menjalankan kaskade dalam tx: kontrak, business, security,
// reinsurer - urutan `DeleteFromTREATYCONTRACT_SQL`. Mengembalikan jumlah
// yang BENAR-BENAR terhapus. Klausul tidak disentuh.
func (k *KaskadeTCO) HapusKontrak(ctx context.Context, tx *Tx, kom models.KombinasiTCO, tahunID, kontrakID string) (
	DampakHapusTCO, error) {

	if tx == nil {
		return DampakHapusTCO{}, errors.New("repository: kaskade hapus kontrak menuntut transaksi")
	}
	t, err := k.tabel()
	if err != nil {
		return DampakHapusTCO{}, err
	}
	var d DampakHapusTCO
	if d.Bersama, err = k.hitung(ctx, sqlKontrakBersamaTCO(t.kontrak, t.tahun),
		kom.TreatyYear, kom.TreatyGroupID, kom.ReinsTypeID, kontrakID); err != nil {
		return DampakHapusTCO{}, err
	}
	for _, l := range langkahHapusKontrakTCO(t, kom, tahunID, kontrakID, &d) {
		if err := PeriksaSQL(l.q); err != nil {
			return DampakHapusTCO{}, err
		}
		hasil, err := tx.tx.ExecContext(ctx, l.q, l.args...)
		if err != nil {
			return DampakHapusTCO{}, fmt.Errorf("repository: kaskade hapus kontrak %s: %w", kontrakID, err)
		}
		if *l.isi, err = hasil.RowsAffected(); err != nil {
			return DampakHapusTCO{}, err
		}
	}
	if d.Kontrak != 1 {
		return DampakHapusTCO{}, fmt.Errorf("%w: kontrak %s terhapus %d baris", ErrKaskadeTidakUtuh, kontrakID, d.Kontrak)
	}
	return d, nil
}

// langkahKaskadeTCO adalah satu DELETE kaskade dan tempat jumlahnya.
type langkahKaskadeTCO struct {
	tabel  string
	q      string
	hitung string
	args   []any
	isi    *int64
}

// langkahHapusKontrakTCO - urutan `DeleteFromTREATYCONTRACT_SQL` b80-b93:
// kontrak, business, security, reinsurer. ⛔ Tidak ada langkah klausul.
//
// ⛔ SEPERTI PEGA walau kombinasinya dipakai kontrak lain (`d.Bersama` > 0):
// seluruh anak kombinasi ikut terhapus [keputusan work owner 29-09-2026,
// OQ-TCO-21 - "hapus saja, samain dengan pega"]. Yang mencegah penghapusan
// diam adalah popup (cacah kontrak lain) dan konfirmasinya, bukan kaskade ini.
func langkahHapusKontrakTCO(t tabelKaskadeTCO, kom models.KombinasiTCO, tahunID, kontrakID string, d *DampakHapusTCO) []langkahKaskadeTCO {
	argBiz := []any{kom.TreatyYear, tahunID, kom.TreatyGroupID, kom.ReinsTypeID}
	saringSec := fmt.Sprintf(saringSecurityKaskadeTCO, t.reinsurer)
	return []langkahKaskadeTCO{
		{t.kontrak, sqlHapusKontrakKaskadeTCO(t.kontrak), "", []any{kontrakID, tahunID}, &d.Kontrak},
		{t.business, sqlHapusKaskadeTCO(t.business, saringBusinessKaskadeTCO), sqlHitungTCO(t.business, saringBusinessKaskadeTCO), argBiz, &d.Business},
		{t.security, sqlHapusKaskadeTCO(t.security, saringSec), sqlHitungTCO(t.security, saringSec), argKombinasi(kom), &d.Security},
		{t.reinsurer, sqlHapusKaskadeTCO(t.reinsurer, saringReinsurerKaskadeTCO),
			sqlHitungTCO(t.reinsurer, saringReinsurerKaskadeTCO), argKombinasi(kom), &d.Reinsurer},
	}
}

// DampakReinsurer menghitung security yang ikut terhapus bersama reinsurer itu.
func (k *KaskadeTCO) DampakReinsurer(ctx context.Context, reinsurerID string) (DampakHapusTCO, error) {
	t, err := k.tabel()
	if err != nil {
		return DampakHapusTCO{}, err
	}
	n, err := k.hitung(ctx, sqlHitungSecurityReinsurerTCO(t.security), reinsurerID)
	return DampakHapusTCO{Reinsurer: 1, Security: n}, err
}

// HapusReinsurer - `DeleteFromTreatyReinsurer_Act`: security lalu reinsurer,
// reinsurer dibatasi kombinasinya.
func (k *KaskadeTCO) HapusReinsurer(ctx context.Context, tx *Tx, kom models.KombinasiTCO, reinsurerID string) (DampakHapusTCO, error) {
	if tx == nil {
		return DampakHapusTCO{}, errors.New("repository: hapus reinsurer menuntut transaksi")
	}
	t, err := k.tabel()
	if err != nil {
		return DampakHapusTCO{}, err
	}
	var d DampakHapusTCO
	for _, l := range []struct {
		q    string
		args []any
		isi  *int64
	}{
		{sqlHapusSecurityReinsurerTCO(t.security), []any{reinsurerID}, &d.Security},
		{sqlHapusReinsurerSatuTCO(t.reinsurer), append([]any{reinsurerID}, argKombinasi(kom)...), &d.Reinsurer},
	} {
		if err := PeriksaSQL(l.q); err != nil {
			return DampakHapusTCO{}, err
		}
		hasil, err := tx.tx.ExecContext(ctx, l.q, l.args...)
		if err != nil {
			return DampakHapusTCO{}, fmt.Errorf("repository: hapus reinsurer %s: %w", reinsurerID, err)
		}
		if *l.isi, err = hasil.RowsAffected(); err != nil {
			return DampakHapusTCO{}, err
		}
	}
	if d.Reinsurer != 1 {
		return DampakHapusTCO{}, fmt.Errorf("%w: reinsurer %s terhapus %d baris", ErrKaskadeTidakUtuh, reinsurerID, d.Reinsurer)
	}
	return d, nil
}
