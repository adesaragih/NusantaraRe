package repository

// Untuk apa berkas ini: PEMBACA DOKUMEN LAMA - tiket 22 (spec-penyimpanan
// ID-3, ID-21; KEPUTUSAN-RONDE-12 butir 5).
//
// `POOLDATA.JSON_POLIS` hanya DIBACA - tidak dibuat, tidak diubah, tidak
// dihapus (nol penulisan JSON). Isi dokumen ditulis ke 8 tabel diagram lewat
// antarmuka yang SAMA dengan jalur biasa - `SisipKasus`, `SimpanHalaman`,
// `SetelNomorPolis`, `TutupKasus` (ID-3, AC 56; dijaga
// `TestPemuatTanpaJalurTulisTerpisah`). Berkas ini hanya menambah satu
// pembaruan: empat kolom datar json_polis (ID-21) yang jalur biasa tidak
// pernah tulis karena di sistem baru ia lahir dari sesi, bukan dari dokumen.
// Salinan SuggestList dokumen lama (F3) ditulis lewat `CatatUsulan` jalur
// biasa; di sini hanya penjaga dobelnya (`SalinUsulanLama`, baca IDPEGA).
//
// `[terverifikasi]` bentuk baris: procedure `PEGA_JSON_POLIS_TREATYIN`
// (PERTANYAAN-untuk-DBA P1 §2) - `IDPEGA`, `DATA_JSON` CLOB, `TGL_INPUT`
// (SYSDATE), `NOPOLIS`, `NOENDORS`, `PRODKE`, `TGL_PROD` (DATE), `USERNAME`.
// Tabel dipakai bersama lini lain: dokumen non-Treaty In disaring di Go
// (`models.PecahDokumenLama`, `pxObjClass` akar).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// tabelPolisJSON - tabel warisan dokumen polis (baca saja).
const tabelPolisJSON = "JSON_POLIS"

// sqlKunciJSONPolis - ROWID baris generasi NB (PRODKE '0', `SavePolisTreatyIn_SQL`)
// beserta yang PRODKE-nya kosong (dilaporkan pemecah, tidak ditebak).
// Dokumen dibaca satu per satu lewat ROWID supaya tidak ada kursor CLOB yang
// terbuka sepanjang pemuatan.
func sqlKunciJSONPolis(t string) string {
	return fmt.Sprintf(`SELECT ROWIDTOCHAR(ROWID) FROM %s
	  WHERE PRODKE IS NULL OR TRIM(PRODKE) = '0'
	  ORDER BY IDPEGA, ROWID`, t)
}

// sqlHitungJSONPolisLain - baris generasi endorsemen (seluruh lini), hanya dihitung.
func sqlHitungJSONPolisLain(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE PRODKE IS NOT NULL AND TRIM(PRODKE) <> '0'`, t)
}

func sqlBacaJSONPolis(t string) string {
	return fmt.Sprintf(`SELECT IDPEGA, NOPOLIS, NOENDORS, TO_CHAR(PRODKE),
	        TO_CHAR(TGL_INPUT, '%s'), TO_CHAR(TGL_PROD, '%s'), USERNAME, DATA_JSON
	   FROM %s WHERE ROWID = CHARTOROWID(:1)`, fmtTanggal, fmtTanggal, t)
}

func sqlIDPegaKasus(t string) string {
	return fmt.Sprintf(`SELECT IDPEGA FROM %s WHERE ID = :1`, t)
}

func sqlSetelKolomDatarLama(t string) string {
	return fmt.Sprintf(`UPDATE %s SET IDPEGA = :1, NOENDORS = :2, TGL_INPUT = TO_DATE(:3, '%s'), USERNAME = :4 WHERE ID = :5`,
		t, fmtTanggal)
}

// sqlAdaUsulanIDPega - penjaga dobel salinan SuggestList lama (F3): cacah
// baris riwayat produksi ber-IDPEGA itu.
func sqlAdaUsulanIDPega(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, t)
}

// SalinUsulanLama menyalin baris SuggestList dokumen lama ke
// POOLDATA.HISTORYAKSEPTASIPRODUCTION (`[keputusan work owner]` F3
// 04-10-2026) lewat penulis yang SAMA dengan jalur biasa (`CatatUsulan`,
// pengganti `InsertViewSuggest_SQL`; ID-3, AC 56), di transaksi pemuatan
// dokumen itu. Penjaga dobel menurut IDPEGA: bila IDPEGA itu sudah punya
// baris riwayat produksi, tidak ada yang ditulis dan hasilnya false -
// pemuat yang diulang tidak menggandakan baris. NOURUT = 1..n berurut baris
// dokumen (`CatatUsulan` MAX+1 dari nol = `.pxListSubscript` langkah 2.1.2).
func (g *Gudang) SalinUsulanLama(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) (bool, error) {
	if len(baris) == 0 {
		return true, nil
	}
	t, err := g.nama(tabelRiwayatProduksi)
	if err != nil {
		return false, err
	}
	q := sqlAdaUsulanIDPega(t)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := g.pembaca(tx).QueryRowContext(ctx, q, idPega).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa riwayat produksi IDPEGA: %w", err)
	}
	if n > 0 {
		return false, nil
	}
	return true, g.CatatUsulan(ctx, tx, idPega, baris)
}

// KunciJSONPolis - kunci baca (ROWID) setiap baris generasi NB di JSON_POLIS.
func (g *Gudang) KunciJSONPolis(ctx context.Context) ([]string, error) {
	t, err := g.nama(tabelPolisJSON)
	if err != nil {
		return nil, err
	}
	q := sqlKunciJSONPolis(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kunci JSON_POLIS: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, fmt.Errorf("repository: membaca kunci JSON_POLIS: %w", err)
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// HitungJSONPolisLain - cacah baris generasi endorsemen yang tidak dibaca.
func (g *Gudang) HitungJSONPolisLain(ctx context.Context) (int, error) {
	t, err := g.nama(tabelPolisJSON)
	if err != nil {
		return 0, err
	}
	q := sqlHitungJSONPolisLain(t)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: menghitung JSON_POLIS: %w", err)
	}
	return n, nil
}

// BacaJSONPolis membaca satu baris JSON_POLIS menurut kuncinya.
func (g *Gudang) BacaJSONPolis(ctx context.Context, kunci string) (models.BarisJSONPolis, error) {
	t, err := g.nama(tabelPolisJSON)
	if err != nil {
		return models.BarisJSONPolis{}, err
	}
	q := sqlBacaJSONPolis(t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisJSONPolis{}, err
	}
	var idpega, nopol, noendors, prodke, tglInput, tglProd, user, dok sql.NullString
	err = g.db.QueryRowContext(ctx, q, kunci).Scan(&idpega, &nopol, &noendors, &prodke, &tglInput, &tglProd, &user, &dok)
	if err != nil {
		return models.BarisJSONPolis{}, fmt.Errorf("repository: membaca JSON_POLIS %s: %w", kunci, err)
	}
	return models.BarisJSONPolis{
		IDPega: teks(idpega), NoPolis: teks(nopol), NoEndors: teks(noendors), ProdKe: teks(prodke),
		TglInput: teks(tglInput), TglProd: teks(tglProd), Username: teks(user), DataJSON: []byte(dok.String),
	}, nil
}

// IDPegaKasus - IDPEGA generasi ber-ID itu; ErrKasusTidakAda bila belum ada.
// Pemuat memakainya untuk membedakan jalankan ulang (IDPEGA sama: dilewati)
// dari ID yang dipakai kasus lain (IDPEGA berbeda: galat, tidak ditimpa).
func (g *Gudang) IDPegaKasus(ctx context.Context, tx *db.Tx, id string) (string, error) {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return "", err
	}
	q := sqlIDPegaKasus(t)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrKasusTidakAda
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca IDPEGA kasus: %w", err)
	}
	return teks(v), nil
}

// SetelKolomDatarLama menulis kolom datar json_polis apa adanya (ID-21) ke
// generasi yang baru disisipkan `SisipKasus` di transaksi pemanggil yang
// sama - baris itu belum pernah dilihat siapa pun, jadi tanpa syarat
// generasi terbuka.
func (g *Gudang) SetelKolomDatarLama(ctx context.Context, tx *db.Tx, id string, k models.KolomDatarLama) error {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyimpan kolom datar json_polis", sqlSetelKolomDatarLama(t),
		db.KosongJadiNil(k.IDPega), db.KosongJadiNil(k.NoEndors), db.KosongJadiNil(k.TglInput),
		db.KosongJadiNil(k.Username), id)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, "kolom datar json_polis")
}
