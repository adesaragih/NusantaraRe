package repository

// Baca SATU kontrak warisan - `TREATY_IN` berpasangan `M_TREATY_IN`.
//
// ⛔ BACA SAJA, kedua tabel. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL,
// nol penyebutan di migrasi. Dijaga `TestWarisanHanyaDibaca`.
//
// ⛔ Jalur ini TERPISAH dari `kontrak.go` yang membaca `KONTRAK` +
// `VERSI_KONTRAK`. Keduanya menjawab pertanyaan yang berbeda, dan
// `GET /kontrak/{id}` tidak disentuh.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// ErrWarisanTidakAda - pengenalnya tidak menunjuk baris mana pun.
var ErrWarisanTidakAda = errors.New("kontrak warisan tidak ada")

// ErrJSONWarisanRusak - dokumennya ada tetapi tidak dapat diurai.
var ErrJSONWarisanRusak = errors.New("dokumen JSON warisan tidak dapat diurai")

// teks membaca penunjuk yang boleh nil menjadi string kosong.
func teks(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// BacaKontrakWarisan membaca satu kontrak beserta dokumen aslinya.
//
// ⚠️ `LEFT JOIN` ke `M_TREATY_IN`: sapuan menemukan 1.854 lawan 1.854 dan
// nol dokumen NULL hari ini, tetapi kontrak yang kehilangan dokumennya harus
// tetap TERBUKA - delapan medan kolomnya masih dapat dibaca, dan layar yang
// menolak membuka kontrak karena dokumennya hilang menyembunyikan justru
// kontrak yang paling perlu dilihat.
func (g *Gudang) BacaKontrakWarisan(ctx context.Context, id string) (models.KontrakWarisan, error) {
	tKontrak, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	q := fmt.Sprintf(`SELECT ID, TREATYCONTRACTNAME, TERITORIALSCOPE, TREATYYEAR,
		CEDING, CEDINGID, LEADINGREINSSOURCE, LEADINGREINSSOURCEID,
		PROPORTIONTYPE, COMMENCEMENT, TERMINATION
		FROM %s WHERE ID = :1`, tKontrak)

	var k models.KontrakWarisan
	var sid, nk, ts, ty, cd, cdid, lrs, lrsid, pt, cm, tm sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&sid, &nk, &ts, &ty, &cd, &cdid,
		&lrs, &lrsid, &pt, &cm, &tm)
	if errors.Is(err, sql.ErrNoRows) {
		return models.KontrakWarisan{}, fmt.Errorf("%w: %s", ErrWarisanTidakAda, id)
	}
	if err != nil {
		return models.KontrakWarisan{}, fmt.Errorf("repository: membaca kontrak warisan %s: %w", id, err)
	}

	k.ID = sid.String
	k.NamaKontrak = nk.String
	// ⭐ `TeritorialScope` diambil dari KOLOM. Keduanya ada di sistem lama
	// dan 1.844 dari 1.854 identik begitu `\n` di-unescape; kolomnya sudah
	// berbentuk teks siap tampil.
	k.LingkupWilayah = ts.String
	k.TahunTreaty = ty.String
	k.Cedant = cd.String
	k.IDCedant = cdid.String
	k.AsalBisnis = lrs.String
	k.IDAsalBisnis = lrsid.String
	k.SifatProporsiAsli = pt.String
	k.TanggalMulaiAsli = cm.String
	k.TanggalBerakhirAsli = tm.String
	k.AdaDiJSON = map[string]bool{}
	k.TeksMentah = map[string]string{}

	// ⛔⛔ `LEFT JOIN M_TREATY_IN` DAN SELURUH PENGURAIAN DOKUMEN DICABUT,
	// 6 Oktober 2026.
	//
	// Keputusan pemilik proses: nilai yang ditarik dari `JSONDATA` dilarang
	// keras. Kesembilan medan kepala, kelima ejaan tab teks, grid Rate of
	// Exchange, dan keempat tab Limits kini datang dari TABEL PENDARATAN;
	// services merangkainya lewat `BacaRevisiPendaratan`,
	// `BacaKursPendaratan`, `BacaLayerPendaratan`, dan
	// `BacaPohonLimitsPendaratan`.
	//
	// ⭐ Fungsi ini karena itu TIDAK LAGI menyentuh `M_TREATY_IN` sama
	// sekali — dan itu juga membuang pengurai 28 KB rata-rata per pembukaan
	// kontrak.
	return k, nil
}
