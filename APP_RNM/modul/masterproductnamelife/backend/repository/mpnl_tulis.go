package repository

// Penulis produk (paket 3: sisi umum; paket 4: sisi inward, satu transaksi).
//
//	SaveProductName_Act 8 b1623 `·` PRE=false  DATAPEGA ← @GetPageJSONString() halaman ProductName
//	                    9 b1831 `·` PRE=false  RDB SaveProductNameLIfe → PEGA_M_PRODUCT_LIFE (TIDAK dipanggil)
//	                   10 b2019 `·` PRE=false  RDB SaveProductNameLIfeFlat: UPDATE … SET RIRISKID, RIRISK WHERE ID
//
// ⛔ Prosedur ditiru: upsert dikunci `ID`, `ID` baru dari sequence; kolom datar
// `RIRISKID`, `RIRISK` (langkah 10, hidup - RALAT R8) ditulis di pernyataan YANG
// SAMA dengan `JSONDATA`, di transaksi pemanggil. Nol COMMIT. Hanya dua kolom
// datar itu - persis `SaveProductNameLIfeFlat` b84 dan katalog DEV (lanjutan 1 L1).
// ⛔ `JSONDATA` diikat sebagai CLOB (`go_ora.Clob`): daftar komentar tumbuh
// setiap simpan dan dapat melampaui batas VARCHAR2 bind.
// ⛔ go-ora mengikat menurut URUTAN KEMUNCULAN placeholder - nomor placeholder
// di setiap SQL di bawah = urutan argumennya.

import (
	"context"
	"errors"
	"fmt"

	go_ora "github.com/sijms/go-ora/v2"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

func sqlSisipUmum(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, JSONDATA, RIRISKID, RIRISK) VALUES (:1, :2, :3, :4)`, tabel)
}

func sqlPerbaruiUmum(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET JSONDATA = :1, RIRISKID = :2, RIRISK = :3 WHERE ID = :4`, tabel)
}

func sqlSisipInward(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, JSONDATA) VALUES (:1, :2)`, tabel)
}

func sqlPerbaruiInward(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET JSONDATA = :1 WHERE ID = :2`, tabel)
}

func argPerbaruiInward(id, jsonInward string) []any { return []any{clob(jsonInward), id} }

func clob(teks string) go_ora.Clob { return go_ora.Clob{String: teks, Valid: true} }

// datar - kolom datar `M_PRODUCT_LIFE` dari produk (`SaveProductNameLIfeFlat` b84).
func datar(p models.Produk) []any {
	return []any{db.KosongJadiNil(p.Umum.RIRiskID), db.KosongJadiNil(p.Umum.RIRisk)}
}

func argSisipUmum(p models.Produk, jsonUmum string) []any {
	return append([]any{p.ID, clob(jsonUmum)}, datar(p)...)
}

func argPerbaruiUmum(p models.Produk, jsonUmum string) []any {
	return append(append([]any{clob(jsonUmum)}, datar(p)...), p.ID)
}

// KunciProduk - produk utuh, barisnya dikunci `FOR UPDATE` (di dalam simpan).
func (g *Gudang) KunciProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	s, err := g.AmbilSimpanan(ctx, tx, id, true)
	if err != nil {
		return models.Produk{}, err
	}
	return UraiProduk(s.ID, s.JSONUmum, s.IDInward, s.JSONInward)
}

func (g *Gudang) exec(ctx context.Context, tx *db.Tx, objek string, susun func(string) string, args ...any) error {
	if !tx.Terisi() {
		return errors.New("repository: writing a product requires a transaction")
	}
	q, err := g.siapkan(objek, susun)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: writing %s: %w", objek, err)
	}
	return db.PastikanSatuBaris(hasil, "writing "+objek)
}

// SisipProduk menerbitkan ID baru (sequence) dan menulis produk baru - KEDUA
// tabel di transaksi pemanggil (P4). Baris inward ber-`ID` = `PRODUCTID` = ID
// produk (R14). Mengembalikan ID itu.
//
// Salinan (`Copy`, `p.SalinanDari`): JSON tersimpan produk asal menjadi dasar kedua
// sisi - `CopyProduct` menyalin halaman utuh, termasuk kunci yang tidak dikelola layar.
func (g *Gudang) SisipProduk(ctx context.Context, tx *db.Tx, p models.Produk) (string, error) {
	var dasar SimpananProduk
	if p.SalinanDari != "" {
		s, err := g.AmbilSimpanan(ctx, tx, p.SalinanDari, false)
		if err != nil {
			return "", err
		}
		dasar = s
	}
	id, err := g.identitasBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	p.ID = id
	p.Inward.ID, p.Inward.ProductID = id, id
	umum, err := RakitUmum(p, dasar.JSONUmum, true)
	if err != nil {
		return "", err
	}
	if err := g.exec(ctx, tx, TabelProduk, sqlSisipUmum, argSisipUmum(p, umum)...); err != nil {
		return "", err
	}
	inward, err := RakitInward(p, dasar.JSONInward)
	if err != nil {
		return "", err
	}
	if err := g.exec(ctx, tx, TabelInward, sqlSisipInward, id, clob(inward)); err != nil {
		return "", err
	}
	return id, nil
}

// PerbaruiProduk menulis ulang produk yang ada - kedua tabel; kunci JSON yang
// tidak dikelola dipertahankan dari JSON tersimpan (dibaca di transaksi yang
// sama). Baris inward lama diperbarui menurut ID-NYA sendiri (data lama dapat
// ber-ID sequence inward); produk tanpa baris inward mendapat baris baru ber-ID
// produk.
func (g *Gudang) PerbaruiProduk(ctx context.Context, tx *db.Tx, p models.Produk) error {
	s, err := g.AmbilSimpanan(ctx, tx, p.ID, true)
	if err != nil {
		return err
	}
	umum, err := RakitUmum(p, s.JSONUmum, false)
	if err != nil {
		return err
	}
	if err := g.exec(ctx, tx, TabelProduk, sqlPerbaruiUmum, argPerbaruiUmum(p, umum)...); err != nil {
		return err
	}
	p.Inward.ProductID = p.ID
	if !s.AdaInward {
		if s.InwardMilikLain != "" {
			// Menyisipkan baris ber-ID sama akan menggandakan ID (nol PK) di samping baris produk lain.
			return fmt.Errorf("%w: inward row %s in %s belongs to product %s", ErrIdentitasBentrok, p.ID, TabelInward,
				s.InwardMilikLain)
		}
		p.Inward.ID = p.ID
		inward, err := RakitInward(p, "")
		if err != nil {
			return err
		}
		return g.exec(ctx, tx, TabelInward, sqlSisipInward, p.ID, clob(inward))
	}
	p.Inward.ID = s.IDInward
	inward, err := RakitInward(p, s.JSONInward)
	if err != nil {
		return err
	}
	return g.exec(ctx, tx, TabelInward, sqlPerbaruiInward, argPerbaruiInward(s.IDInward, inward)...)
}
