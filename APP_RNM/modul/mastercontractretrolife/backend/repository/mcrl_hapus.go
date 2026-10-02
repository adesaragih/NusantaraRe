package repository

// Kaskade hapus (paket 7, K2) - DEV nol FK, jadi anak dihapus DI SINI, lebih
// dulu, di dalam transaksi pemanggil. Relasinya = keempat FK yang spec §3
// gambarkan (dan DEV tidak punya):
//
//	security.TREATYREINSURERID -> reinsurer.ID
//	reinsurer.TREATYCONTRACTID -> kontrak.ID
//	business.TREATYCONTRACTID  -> kontrak.ID
//
// Padanan Pega (datar): `DeleteTreatyLimit_SQL` b84, `DeleteSecurityReinsurer_SQL`
// b85, `DeleteSecurityReinsurerLife_SQL` b85, `DeleteRowBusinessList` b84.
// ⛔ Pencacah (popup) dan penghapus memakai predikat yang SAMA.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// predikat anak per relasi; `%s` pertama = tabel anak, `%s` kedua (bila ada) = reinsurer.
const (
	anakSecurityDariKontrak   = `TREATYREINSURERID IN (SELECT ID FROM %s WHERE TREATYCONTRACTID = :1)`
	anakReinsurerDariKontrak  = `TREATYCONTRACTID = :1`
	anakBusinessDariKontrak   = `TREATYCONTRACTID = :1`
	anakSecurityDariReinsurer = `TREATYREINSURERID = :1`
)

func sqlCacah(tabel, predikat string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, tabel, predikat)
}

func sqlHapus(tabel, predikat string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s`, tabel, predikat)
}

// langkah adalah satu relasi anak: tabelnya dan predikatnya.
type langkah struct {
	tabel    string
	predikat func(reinsurer string) string
	ke       func(d *models.Dampak) *int64
}

func (g *Gudang) langkahKontrak() []langkah {
	return []langkah{
		{TabelSecurity, func(r string) string { return fmt.Sprintf(anakSecurityDariKontrak, r) },
			func(d *models.Dampak) *int64 { return &d.Security }},
		{TabelReinsurer, func(string) string { return anakReinsurerDariKontrak },
			func(d *models.Dampak) *int64 { return &d.Reinsurer }},
		{TabelBusiness, func(string) string { return anakBusinessDariKontrak },
			func(d *models.Dampak) *int64 { return &d.Business }},
	}
}

func (g *Gudang) langkahReinsurer() []langkah {
	return []langkah{{TabelSecurity, func(string) string { return anakSecurityDariReinsurer },
		func(d *models.Dampak) *int64 { return &d.Security }}}
}

// jalankan mencacah (hapus=false) atau menghapus (hapus=true) seluruh anak,
// dengan urutan daftar langkah (anak terdalam lebih dulu).
func (g *Gudang) jalankan(ctx context.Context, tx *db.Tx, id string, daftar []langkah, hapus bool) (models.Dampak, error) {
	var d models.Dampak
	reinsurer, err := g.db.Qualify(TabelReinsurer)
	if err != nil {
		return d, err
	}
	for _, l := range daftar {
		tabel, err := g.db.Qualify(l.tabel)
		if err != nil {
			return d, err
		}
		pred := l.predikat(reinsurer)
		if hapus {
			q := sqlHapus(tabel, pred)
			if err := db.PeriksaSQL(q); err != nil {
				return d, err
			}
			h, err := tx.ExecContext(ctx, q, id)
			if err != nil {
				return d, fmt.Errorf("repository: deleting %s rows under %s: %w", l.tabel, id, err)
			}
			n, err := h.RowsAffected()
			if err != nil {
				return d, fmt.Errorf("repository: counting deleted %s rows: %w", l.tabel, err)
			}
			*l.ke(&d) = n
			continue
		}
		q := sqlCacah(tabel, pred)
		if err := db.PeriksaSQL(q); err != nil {
			return d, err
		}
		if err := g.kueri(tx).QueryRowContext(ctx, q, id).Scan(l.ke(&d)); err != nil {
			return d, fmt.Errorf("repository: counting %s rows under %s: %w", l.tabel, id, err)
		}
	}
	return d, nil
}

// hapusSatu menghapus satu baris induk/daun menurut ID; tepat satu baris.
func (g *Gudang) hapusSatu(ctx context.Context, tx *db.Tx, objek, nama, id string) error {
	q, err := g.siapkan(objek, func(t string) string { return sqlHapus(t, `ID = :1`) })
	if err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: deleting %s %s: %w", nama, id, err)
	}
	return adaSatu(h, nama, id)
}

// DampakHapusKontrak - anak yang ikut terhapus bersama kontrak id.
func (g *Gudang) DampakHapusKontrak(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error) {
	return g.jalankan(ctx, tx, id, g.langkahKontrak(), false)
}

// HapusKontrak - (1) security, (2) reinsurer dan business, (3) kontrak.
func (g *Gudang) HapusKontrak(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error) {
	d, err := g.jalankan(ctx, tx, id, g.langkahKontrak(), true)
	if err != nil {
		return d, err
	}
	return d, g.hapusSatu(ctx, tx, TabelKontrak, "treaty contract", id)
}

// DampakHapusReinsurer - security yang ikut terhapus bersama reinsurer id.
func (g *Gudang) DampakHapusReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error) {
	return g.jalankan(ctx, tx, id, g.langkahReinsurer(), false)
}

// HapusReinsurer - (1) security, (2) reinsurer.
func (g *Gudang) HapusReinsurer(ctx context.Context, tx *db.Tx, id string) (models.Dampak, error) {
	d, err := g.jalankan(ctx, tx, id, g.langkahReinsurer(), true)
	if err != nil {
		return d, err
	}
	return d, g.hapusSatu(ctx, tx, TabelReinsurer, "reinsurer", id)
}

// HapusSecurity - daun.
func (g *Gudang) HapusSecurity(ctx context.Context, tx *db.Tx, id string) error {
	return g.hapusSatu(ctx, tx, TabelSecurity, "security reinsurer", id)
}

// HapusBusiness - daun.
func (g *Gudang) HapusBusiness(ctx context.Context, tx *db.Tx, id string) error {
	return g.hapusSatu(ctx, tx, TabelBusiness, "business", id)
}
