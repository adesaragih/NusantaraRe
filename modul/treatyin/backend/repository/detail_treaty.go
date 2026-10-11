package repository

// Tulisan DETAIL kontrak tuntas — `TREATYINDETAIL` (Treaty In) dan
// `TREATYINDETAILEDM` (Adjustment). Perintah WO 9 Oktober 2026; susunan
// barisnya di `services/detail_treaty.go`.
//
// Padanan Pega, di dalam transaksi tombol yang SAMA (kepala sudah ditulis):
//
//	RemoveTreatyInDetail(Edm)   DELETE FROM <tabel> WHERE TREATYID = :id
//	PEGA_M_TREATY_IN_DETAIL     INSERT satu baris per sisipan:
//	(_EDM), per baris             ID           situs || LPAD(n, 6, '0')
//	                              COMMENCEMENT TO_DATE(kepala.COMMENCEMENT, 'YYYYMMDD')
//	                              TERMINATION  TO_DATE(kepala.TERMINATION, 'YYYYMMDD')
//	                              kolom lain   isi SaveData
//
// ⛔ Yang TIDAK ditiru, dan sebabnya:
//   - Dokumen JSON detail (`INSERT … (ID, JSONDATA)` dan `DELETE` padanannya
//     di `RemoveTreatyInDetail`): tabel dokumen itu dilarang keras di
//     aplikasi (`TestAplikasiTidakMenyebutMTreatyIn`).
//   - Sequence prosedur untuk `n`: namanya berawalan tabel dokumen yang sama.
//     `n` = angka TERTINGGI berawalan situs di KEDUA tabel detail (keduanya
//     berbagi satu sequence di Pega) + 1, di bawah kunci baris situs — pola
//     `idKontrakBaru`. Di DEV angka itu 719181 dan sequence Pega berikutnya
//     719182: deretnya bersambung.
//   - Galat per baris yang ditelan (`StsSave := 0`, baris hilang): di sini
//     galat menggagalkan seluruh tombol, dan transaksinya dibatalkan.
//
// ⭐ Tanggal dibaca prosedur dengan `SELECT … INTO` yang menjadi NULL bila
// barisnya nol, lebih dari satu, atau tanggalnya tak terbaca — ditiru dengan
// `COUNT(*) = 1` dan `DEFAULT NULL ON CONVERSION ERROR`.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

const (
	// TabelDetailTreatyIn - detail kontrak Treaty In tuntas.
	TabelDetailTreatyIn = TabelDetailWarisan
	// TabelDetailTreatyInEDM - detail penyesuaian tuntas.
	TabelDetailTreatyInEDM = "TREATYINDETAILEDM"
)

// SQLDetailTreaty - perintah tulisan detail satu tabel (hapus, angka
// pengenal tertinggi, sisip) atas nama yang SUDAH berkualifikasi. Murni.
func SQLDetailTreaty(nama, kepala, ti, edm string, kolom []models.KolomDetail) (hapus, tertinggi, sisip string, err error) {
	hapus = fmt.Sprintf("DELETE FROM %s WHERE TREATYID = :1", nama)
	tertinggi = fmt.Sprintf(`SELECT NVL(MAX(N), 0) FROM (
		SELECT TO_NUMBER(SUBSTR(ID, :1)) N FROM %s WHERE REGEXP_LIKE(ID, :2)
		UNION ALL
		SELECT TO_NUMBER(SUBSTR(ID, :3)) N FROM %s WHERE REGEXP_LIKE(ID, :4))`, ti, edm)

	// Urut bind = urut kemunculan (Oracle tidak memakai nomor :n).
	nomor := 0
	bind := func() string { nomor++; return fmt.Sprintf(":%d", nomor) }
	namaKolom := []string{"ID"}
	nilai := []string{bind()}
	for _, k := range kolom {
		namaKolom = append(namaKolom, k.Nama)
		if k.Angka {
			// Koefisien bulat / pangkat sepuluh — nol NLS (`pecahAngka`);
			// koefisien NULL menghasilkan NULL.
			koef := bind()
			nilai = append(nilai, fmt.Sprintf("(TO_NUMBER(%s) / POWER(10, %s))", koef, bind()))
			continue
		}
		nilai = append(nilai, bind())
	}
	for _, t := range []string{"COMMENCEMENT", "TERMINATION"} {
		namaKolom = append(namaKolom, t)
		nilai = append(nilai, fmt.Sprintf(
			"(SELECT CASE WHEN COUNT(*) = 1 THEN MAX(TO_DATE(%s DEFAULT NULL ON CONVERSION ERROR, 'YYYYMMDD')) END FROM %s WHERE ID = %s)",
			t, kepala, bind()))
	}
	sisip = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", nama, strings.Join(namaKolom, ", "), strings.Join(nilai, ", "))
	for _, q := range []string{hapus, tertinggi, sisip} {
		if err := db.PeriksaSQL(q); err != nil {
			return "", "", "", err
		}
	}
	return hapus, tertinggi, sisip, nil
}

// cacahDetail - cacah baris detail satu kontrak di dalam transaksi (uji `db`).
func (g *Gudang) cacahDetail(ctx context.Context, tx *db.Tx, tabel, treatyID string) (int, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE TREATYID = :1", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, treatyID).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: mencacah %s %s: %w", tabel, treatyID, err)
	}
	return n, nil
}

// tulisDetailTreaty - hapus baris detail `treatyID`, lalu sisipkan `baris`.
func (g *Gudang) tulisDetailTreaty(ctx context.Context, tx *db.Tx, tabel, tabelKepala, treatyID string,
	kolom []models.KolomDetail, baris []models.BarisDetailTreaty) error {
	nama := map[string]string{}
	for _, t := range []string{tabel, tabelKepala, TabelDetailTreatyIn, TabelDetailTreatyInEDM} {
		q, err := g.db.Qualify(t)
		if err != nil {
			return err
		}
		nama[t] = q
	}
	hapus, tertinggi, sisip, err := SQLDetailTreaty(nama[tabel], nama[tabelKepala],
		nama[TabelDetailTreatyIn], nama[TabelDetailTreatyInEDM], kolom)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, hapus, treatyID); err != nil {
		return fmt.Errorf("repository: menghapus %s %s: %w", tabel, treatyID, err)
	}
	if len(baris) == 0 {
		return nil
	}
	situs, err := g.kodeSitusTerkunci(ctx, tx)
	if err != nil {
		return err
	}
	pola := fmt.Sprintf("^%s[0-9]{%d,}$", situs, panjangNomorKontrak)
	mulai := len(situs) + 1
	var n int64
	if err := tx.QueryRowContext(ctx, tertinggi, mulai, pola, mulai, pola).Scan(&n); err != nil {
		return fmt.Errorf("repository: membaca pengenal %s tertinggi: %w", tabel, err)
	}
	for i, b := range baris {
		n++
		args := []any{situs + kiriNol(fmt.Sprint(n), panjangNomorKontrak)}
		for _, k := range kolom {
			switch {
			case k.Angka:
				koef, skala := pecahAngka(b.Angka[k.Nama])
				args = append(args, koef, skala)
			case k.Nama == "TREATYID":
				// Kontrak yang barisnya baru DIHAPUS — hapus dan sisip tidak
				// boleh berselisih pengenal.
				args = append(args, treatyID)
			default:
				args = append(args, kosongNil(b.Teks[k.Nama]))
			}
		}
		args = append(args, treatyID, treatyID)
		if _, err := tx.ExecContext(ctx, sisip, args...); err != nil {
			return fmt.Errorf("repository: menyisipkan %s %s baris ke-%d: %w", tabel, treatyID, i+1, err)
		}
	}
	return nil
}
