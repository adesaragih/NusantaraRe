package repository

// Untuk apa berkas ini: PEMBACA DOKUMEN LAMA ENDORSEMEN dan penulis TAMBAHAN pemuat - tiket EDM 10 (spec-penyimpanan
// ID-3, AC 44) dan 09 (ID-35..ID-37, AC 40-43). Asal: pola `modul/nbtreatyin/backend/repository/lama.go` (tiket NB 22).
//
// `POOLDATA.JSON_POLIS` hanya DIBACA (SELECT) - tidak dibuat, tidak diubah, tidak dihapus. Isi dokumen ditulis lewat
// antarmuka yang SAMA dengan jalur biasa (ID-3, AC 44; dijaga `TestPemuatEDMTanpaJalurTulisTerpisah`):
// `SisipKasus` (PRODKE / OLD_POLIS_ID / NOENDORS / EDM_TYPE + UNIQUE OLD_POLIS_ID), `SimpanHalaman`,
// `SetelNomorPolisSelesai`, `SimpanSelisih` (SUMBER 'PEGA'), `CatatUsulan`, `TutupKasus`. Berkas ini hanya
// menambah dua pembaruan yang jalur biasa tidak pernah tulis:
//
//	SetelKolomDatarLamaEDM  TGL_INPUT, USERNAME json_polis apa adanya (pola NB `SetelKolomDatarLama`) - di sistem
//	                        baru keduanya lahir dari sesi, bukan dari dokumen
//	SetelPenandaMigrasi     PASANGAN_BERGESER / RUMUS_BERLAPIS baris selisih SUMBER 'PEGA' (AC 43: SUMBER disaring
//	                        di SQL, baris 'GO' tidak pernah tersentuh)
//
// `[terverifikasi]` bentuk baris JSON_POLIS: procedure `PEGA_JSON_POLIS_TREATYIN` (dipanggil
// `RDBList/SavePolisTreatyInEDM_SQL`) - `IDPEGA`, `DATA_JSON` CLOB, `TGL_INPUT` (SYSDATE), `NOPOLIS`, `NOENDORS`,
// `PRODKE`, `TGL_PROD` (DATE), `USERNAME`. Dokumen lini lain disaring di Go (`models.PecahDokumenEDM`).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// sqlPmKunciJSONPolisEDM - kunci baca (ROWID) + kunci urut setiap baris generasi endorsemen (PRODKE bukan '0';
// PRODKE kosong milik pemuat NB, `sqlKunciJSONPolis` NB). Urutan GENERASI (PRODKE bilangan) disusun di Go
// (`models.UrutKunciGenerasi`) - PRODKE VARCHAR2 di DEV, ORDER BY teks salah urut ("9" > "10") dan TO_NUMBER akan
// menggagalkan seluruh kueri pada satu nilai rusak. Dokumen dibaca satu per satu lewat ROWID: tanpa kursor CLOB
// terbuka sepanjang pemuatan.
func sqlPmKunciJSONPolisEDM(t string) string {
	return fmt.Sprintf(`SELECT ROWIDTOCHAR(ROWID), NOPOLIS, TRIM(TO_CHAR(PRODKE)) FROM %s
	  WHERE PRODKE IS NOT NULL AND TRIM(TO_CHAR(PRODKE)) <> '0'
	  ORDER BY NOPOLIS, ROWID`, t)
}

// tabelKerjaPega - tabel kerja Pega (kelas ASM-FW-GISFW-Work-*, skema DATAPEGA): PZINSKEY = JSON_POLIS.IDPEGA dokumen
// Pega asli. Hanya DIBACA.
const tabelKerjaPega = "DATAPEGA.PC_ASM_FW_GISFW_WORK"

// sqlPmPembuatPega - pembuat kasus Pega dokumen lama (WO 07-10-2026 "PXCREATEOPERATOR,PXCREATEOPNAME"): satu baris
// tabel kerja Pega menurut pzInsKey = IDPEGA. Baca saja.
func sqlPmPembuatPega() string {
	return fmt.Sprintf(`SELECT PXCREATEOPERATOR, PXCREATEOPNAME FROM %s WHERE PZINSKEY = :1`, tabelKerjaPega)
}

// PembuatPega - PXCREATEOPERATOR / PXCREATEOPNAME kasus Pega ber-pzInsKey `idPega` untuk CREATE_OP / CREATE_OP_NAME
// berkas salinan. Tanpa baris Pega (alat pemuat CLI tidak menyaring tabel kerja Pega) = kosong: pembuat tidak
// dikarang, ditulis NULL.
func (g *Gudang) PembuatPega(ctx context.Context, tx *db.Tx, idPega string) (string, string, error) {
	q := sqlPmPembuatPega()
	if err := db.PeriksaSQL(q); err != nil {
		return "", "", err
	}
	var op, nama sql.NullString
	err := g.pembaca(tx).QueryRowContext(ctx, q, idPega).Scan(&op, &nama)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("repository: membaca pembuat kasus Pega: %w", err)
	}
	return strings.TrimSpace(op.String), strings.TrimSpace(nama.String), nil
}

// sqlPmKunciJSONPolisEDMCopyOld - kunci popup Copy Old. ⛔ Perintah work owner 07-10-2026: data lama dipilih dengan
// `SELECT * FROM DATAPEGA.PC_ASM_FW_GISFW_WORK a, json_polis b, treatyinproduction c WHERE a.pzinskey = b.idpega AND
// b.idpega = c.idpega` - dokumen yang kasus Pega-nya ada DAN sudah berproduksi. Ditulis EXISTS: TREATYINPRODUCTION
// berbaris banyak per IDPEGA. Sama dengan Copy Old NB (`modul/nbtreatyin` `sqlKunciJSONPolisCopyOld`); alat pemuat
// massal tetap tanpa saringan (KEPUTUSAN 23-09-2026 butir 4).
func sqlPmKunciJSONPolisEDMCopyOld(t, prod string) string {
	return fmt.Sprintf(`SELECT ROWIDTOCHAR(b.ROWID), b.NOPOLIS, TRIM(TO_CHAR(b.PRODKE)) FROM %s b
	  WHERE b.PRODKE IS NOT NULL AND TRIM(TO_CHAR(b.PRODKE)) <> '0'
	    AND EXISTS (SELECT 1 FROM %s a WHERE a.PZINSKEY = b.IDPEGA)
	    AND EXISTS (SELECT 1 FROM %s c WHERE c.IDPEGA = b.IDPEGA)
	  ORDER BY b.NOPOLIS, b.ROWID`, t, tabelKerjaPega, prod)
}

func sqlPmBacaJSONPolis(t string) string {
	return fmt.Sprintf(`SELECT IDPEGA, NOPOLIS, NOENDORS, TRIM(TO_CHAR(PRODKE)),
	        TO_CHAR(TGL_INPUT, '%s'), TO_CHAR(TGL_PROD, '%s'), USERNAME, DATA_JSON
	   FROM %s WHERE ROWID = CHARTOROWID(:1)`, fmtTanggal, fmtTanggal, t)
}

// sqlPmBacaJSONPolisNB - dokumen generasi NB satu nomor polis (pembanding generasi pertama saat uji-kering).
func sqlPmBacaJSONPolisNB(t string) string {
	return fmt.Sprintf(`SELECT IDPEGA, NOPOLIS, NOENDORS, TRIM(TO_CHAR(PRODKE)),
	        TO_CHAR(TGL_INPUT, '%s'), TO_CHAR(TGL_PROD, '%s'), USERNAME, DATA_JSON
	   FROM %s WHERE NOPOLIS = :1 AND (PRODKE IS NULL OR TRIM(TO_CHAR(PRODKE)) = '0')
	  ORDER BY ROWID`, fmtTanggal, fmtTanggal, t)
}

func sqlPmKunciGenerasi(t string) string {
	return fmt.Sprintf(`SELECT NOPOLIS, PRODKE, NOENDORS FROM %s WHERE ID = :1`, t)
}

// sqlPmSetelKolomDatarLama - kolom datar json_polis generasi endorsemen yang baru disisipkan; hanya generasi
// terbuka (`syaratTerbuka`, ID-10) ber-PRODKE >= 1, sama dengan setiap UPDATE T_GENERAL_POLIS_TREATY lain.
func sqlPmSetelKolomDatarLama(t string) string {
	return fmt.Sprintf(`UPDATE %s g SET TGL_INPUT = TO_DATE(:1, '%s'), USERNAME = :2 WHERE g.ID = :3 AND g.PRODKE >= 1 AND %s`,
		t, fmtTanggal, syaratTerbuka(t))
}

// sqlPmSetelPenanda - penanda satu baris selisih; SUMBER 'PEGA' disaring di SQL (ID-37, AC 43). `berlapis` -
// tabel ber-RUMUS_BERLAPIS (361, 362; 363 tidak).
func sqlPmSetelPenanda(anak, induk string, berlapis bool) string {
	set, n := "PASANGAN_BERGESER = :1", 2
	if berlapis {
		set, n = set+", RUMUS_BERLAPIS = :2", 3
	}
	return fmt.Sprintf(`UPDATE %s a SET %s
	  WHERE a.NOURUT = :%d AND a.DIFFERENCE_ID IN (SELECT d.ID FROM %s d WHERE d.POLIS_ID = :%d AND d.SUMBER = :%d)`,
		anak, set, n, induk, n+1, n+2)
}

// sqlPmAdaRiwayatIDPega - proteksi dobel salinan SuggestList lama (F3; perintah WO 07-10-2026 "TAMBAKAN PROTEKSI UNTUK 2 TABLE INI historyakseptasiproduction,historyakseptasiPEGA -
// SAAT COPY, JIKA UDAH ADA PADA 2 TABLE ITU JANGAN DI COPY, SUPAYA TIDAK DOUBLE"):
// cacah sumber riwayat IDPEGA itu - HISTORYAKSEPTASIPRODUCTION (SuggestList) dan HISTORYAKSEPTASIPEGA (History),
// paling banyak satu baris per tabel. > 0 = sudah ada, jangan disalin.
func sqlPmAdaRiwayatIDPega(produksi, pega string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM (
	         SELECT 1 FROM %s WHERE IDPEGA = :1 AND ROWNUM = 1
	         UNION ALL
	         SELECT 1 FROM %s WHERE ID_PEGA = :2 AND ROWNUM = 1)`, produksi, pega)
}

// KunciJSONPolisEDM - kunci setiap baris generasi endorsemen di JSON_POLIS (belum berurut generasi).
func (g *Gudang) KunciJSONPolisEDM(ctx context.Context) ([]models.KunciJSONPolis, error) {
	t, err := g.nama(tabelJSONPolis)
	if err != nil {
		return nil, err
	}
	return g.kunciJSONPolisEDM(ctx, sqlPmKunciJSONPolisEDM(t))
}

// KunciJSONPolisEDMCopyOld - kunci generasi endorsemen yang kasus Pega-nya ada dan sudah berproduksi
// (`sqlPmKunciJSONPolisEDMCopyOld`, popup Copy Old).
func (g *Gudang) KunciJSONPolisEDMCopyOld(ctx context.Context) ([]models.KunciJSONPolis, error) {
	t, err := g.nama(tabelJSONPolis)
	if err != nil {
		return nil, err
	}
	prod, err := g.nama(tabelProduksi)
	if err != nil {
		return nil, err
	}
	return g.kunciJSONPolisEDM(ctx, sqlPmKunciJSONPolisEDMCopyOld(t, prod))
}

func (g *Gudang) kunciJSONPolisEDM(ctx context.Context, q string) ([]models.KunciJSONPolis, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kunci JSON_POLIS endorsemen: %w", err)
	}
	defer rows.Close()
	var out []models.KunciJSONPolis
	for rows.Next() {
		var k, nopol, prodke sql.NullString
		if err := rows.Scan(&k, &nopol, &prodke); err != nil {
			return nil, fmt.Errorf("repository: membaca kunci JSON_POLIS endorsemen: %w", err)
		}
		out = append(out, models.KunciJSONPolis{Kunci: teks(k), NoPolis: teks(nopol), ProdKe: teks(prodke)})
	}
	return out, rows.Err()
}

// pmPindaiJSONPolis - satu baris JSON_POLIS (urutan kolom `sqlPmBacaJSONPolis`).
func pmPindaiJSONPolis(sc interface{ Scan(...any) error }) (models.BarisJSONPolis, error) {
	var idpega, nopol, noendors, prodke, tglInput, tglProd, user, dok sql.NullString
	if err := sc.Scan(&idpega, &nopol, &noendors, &prodke, &tglInput, &tglProd, &user, &dok); err != nil {
		return models.BarisJSONPolis{}, err
	}
	return models.BarisJSONPolis{
		IDPega: teks(idpega), NoPolis: teks(nopol), NoEndors: teks(noendors), ProdKe: teks(prodke),
		TglInput: teks(tglInput), TglProd: teks(tglProd), Username: teks(user), DataJSON: []byte(dok.String),
	}, nil
}

// BacaJSONPolisEDM membaca satu baris JSON_POLIS menurut kuncinya.
func (g *Gudang) BacaJSONPolisEDM(ctx context.Context, kunci string) (models.BarisJSONPolis, error) {
	t, err := g.nama(tabelJSONPolis)
	if err != nil {
		return models.BarisJSONPolis{}, err
	}
	q := sqlPmBacaJSONPolis(t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisJSONPolis{}, err
	}
	b, err := pmPindaiJSONPolis(g.db.QueryRowContext(ctx, q, kunci))
	if err != nil {
		return models.BarisJSONPolis{}, fmt.Errorf("repository: membaca JSON_POLIS %s: %w", kunci, err)
	}
	return b, nil
}

// BacaJSONPolisNB membaca dokumen generasi NB (PRODKE '0' / kosong) satu nomor polis - pembanding generasi pertama
// endorsemen saat uji-kering (tanpa membaca tabel relasional).
func (g *Gudang) BacaJSONPolisNB(ctx context.Context, nopolis string) ([]models.BarisJSONPolis, error) {
	t, err := g.nama(tabelJSONPolis)
	if err != nil {
		return nil, err
	}
	q := sqlPmBacaJSONPolisNB(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, nopolis)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca JSON_POLIS generasi NB: %w", err)
	}
	defer rows.Close()
	var out []models.BarisJSONPolis
	for rows.Next() {
		b, err := pmPindaiJSONPolis(rows)
		if err != nil {
			return nil, fmt.Errorf("repository: membaca JSON_POLIS generasi NB: %w", err)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// KunciGenerasiLama - NOPOLIS / PRODKE / NOENDORS generasi ber-ID itu (ada = false bila belum ada). Pemuat
// membedakan "dokumen sudah dimuat" (kunci sama - dilewati) dari "ID dipakai generasi lain" (galat, tidak menimpa).
func (g *Gudang) KunciGenerasiLama(ctx context.Context, tx *db.Tx, id string) (models.KunciGenerasi, bool, error) {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return models.KunciGenerasi{}, false, err
	}
	q := sqlPmKunciGenerasi(t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.KunciGenerasi{}, false, err
	}
	var nopol, noendors sql.NullString
	var prodke sql.NullInt64
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(&nopol, &prodke, &noendors)
	if errors.Is(err, sql.ErrNoRows) {
		return models.KunciGenerasi{}, false, nil
	}
	if err != nil {
		return models.KunciGenerasi{}, false, fmt.Errorf("repository: memeriksa generasi: %w", err)
	}
	return models.KunciGenerasi{NoPolis: teks(nopol), ProdKe: int(prodke.Int64), NoEndors: teks(noendors)}, true, nil
}

// SetelKolomDatarLamaEDM menulis kolom datar json_polis apa adanya ke generasi endorsemen yang baru disisipkan
// `SisipKasus` di transaksi pemanggil yang sama.
func (g *Gudang) SetelKolomDatarLamaEDM(ctx context.Context, tx *db.Tx, id string, k models.KolomDatarLama) error {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyimpan kolom datar json_polis", sqlPmSetelKolomDatarLama(t),
		db.KosongJadiNil(k.TglInput), db.KosongJadiNil(k.Username), id)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, "kolom datar json_polis")
}

// SetelPenandaMigrasi menulis kedua penanda migrasi setiap baris selisih generasi `polisID` - HANYA baris yang
// induknya SUMBER 'PEGA' (ID-37, AC 43). Setiap baris wajib tersentuh tepat satu kali: penanda untuk baris yang tidak
// ada, atau atas proyeksi 'GO', gagal dan membatalkan transaksi pemanggil. Nol kolom angka disentuh (AC 41).
func (g *Gudang) SetelPenandaMigrasi(ctx context.Context, tx *db.Tx, polisID string, p models.PenandaMigrasi) error {
	induk, err := g.nama(models.TabelSelisih.Nama)
	if err != nil {
		return err
	}
	for _, x := range []struct {
		tabel    models.Tabel
		baris    []models.PenandaBaris
		berlapis bool
	}{
		{models.TabelSelisihSpreading, p.Spreading, true},
		{models.TabelSelisihAngsuran, p.Angsuran, true},
		{models.TabelSelisihLapisan, p.Lapisan, false},
	} {
		if len(x.baris) == 0 {
			continue
		}
		anak, err := g.nama(x.tabel.Nama)
		if err != nil {
			return err
		}
		q := sqlPmSetelPenanda(anak, induk, x.berlapis)
		for i, b := range x.baris {
			args := []any{models.TeksPenanda(b.PasanganBergeser)}
			if x.berlapis {
				args = append(args, models.TeksPenanda(b.RumusBerlapis))
			}
			args = append(args, i+1, polisID, models.SumberPega)
			hasil, err := jalankan(ctx, tx, "menyimpan penanda migrasi "+x.tabel.Nama, q, args...)
			if err != nil {
				return err
			}
			if err := db.PastikanSatuBaris(hasil, fmt.Sprintf("penanda migrasi %s NOURUT %d", x.tabel.Nama, i+1)); err != nil {
				return err
			}
		}
	}
	return nil
}

// AdaRiwayatIDPega - IDPEGA itu sudah punya baris di POOLDATA.HISTORYAKSEPTASIPRODUCTION ATAU
// POOLDATA.HISTORYAKSEPTASIPEGA: proteksi dobel salinan SuggestList dokumen lama (services `salinUsulanLama`).
// Baca saja, di transaksi pemuatan dokumen itu.
func (g *Gudang) AdaRiwayatIDPega(ctx context.Context, tx *db.Tx, idPega string) (bool, error) {
	produksi, err := g.nama(tabelRiwayatProduksi)
	if err != nil {
		return false, err
	}
	pega, err := g.nama(tabelRiwayatPega)
	if err != nil {
		return false, err
	}
	q := sqlPmAdaRiwayatIDPega(produksi, pega)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := g.pembaca(tx).QueryRowContext(ctx, q, idPega, idPega).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa riwayat IDPEGA: %w", err)
	}
	return n > 0, nil
}
