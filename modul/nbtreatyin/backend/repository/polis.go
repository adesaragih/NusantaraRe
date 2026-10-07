package repository

// Untuk apa berkas ini: SIMPAN DAN MUAT HALAMAN KERJA - pengganti
// `SaveJsonPolisTreatyIn_Act` (JSON_POLIS) dan `Obj-Save pyWorkPage`. Halaman
// dipecah ke T_GENERAL_POLIS_TREATY + tabel anak menurut KATALOG
// (`models/katalog.go`), di dalam transaksi pemanggil (tiket 16-21; spec §5.7,
// AC 16, 29, 38; spec-penyimpanan ID-11, ID-12, ID-23..ID-31).
//
// ⛔ NB: baris anak boleh dihapus, NOURUT DINOMORI ULANG rapat 1..n setiap
// simpan (ID-12) - aman sebab di PRODKE 0 belum ada generasi untuk
// dipasangkan. Caranya hapus-lalu-sisip seluruh anak satu polis, di transaksi
// yang sama dengan baris induknya: gagal di tengah = tidak ada yang tersimpan
// (AC 29).
//
// ⛔ Generasi tertutup (ada penerus yang OLD_POLIS_ID-nya menunjuk generasi
// ini - `syaratTerbuka`) ditolak SEBELUM satu pun baris disentuh (ID-10).
// Tanpa kolom penanda: diagram grilling tidak memuatnya.
//
// ⛔ T_POLIS_CEDING anak T_POLIS_QUOTATION (diagram O39, QUOTATION_ID): anak
// dihapus SEBELUM baris quotation ditulis ulang, lalu disisip sesudahnya.
//
// ⛔ Halaman `Quotation` dan salinannya `PolicyTreatyIn.QuotationData`
// disimpan SATU kali (ID-23). Nilai yang ditulis: QuotationData bila terisi,
// selain itu Quotation - layar menyunting QuotationData (`.QuotationData.MOID`,
// `.QuotationData.NoOfferSlip`, `.QuotationData.IsSurveyReport`), aktivitas
// menulis Quotation; keduanya diselaraskan `SalinQuotation` di pra-proses.
// Saat dimuat, keduanya diisi nilai yang sama.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// anakTabel menyebut kolom penunjuk induk tabel anak dan cucu.
type anakTabel struct {
	tabel models.Tabel
	// kolomInduk - POLIS_ID untuk anak, INSTALMENT_ID / XOL_ID untuk cucu.
	kolomInduk string
	// cucu - tabel cucu (daftar bersarang di setiap baris), bila ada.
	cucu *anakTabel
}

var (
	cucuAngsuran = &anakTabel{tabel: models.TabelAngsuranRinci, kolomInduk: "INSTALMENT_ID"}
	cucuLayer    = &anakTabel{tabel: models.TabelLayerXOL, kolomInduk: "XOL_ID"}
	// tabelAnak - urutan sisip (induk sebelum cucu ada di dalam tiap entri).
	// Nilai kolom induk anak = ID polis: QUOTATION_ID adalah POLIS_ID baris
	// T_POLIS_QUOTATION (1:1, kunci utama bersama).
	tabelAnak = []anakTabel{
		{tabel: models.TabelCeding, kolomInduk: "QUOTATION_ID"},
		{tabel: models.TabelAngsuran, kolomInduk: "POLIS_ID", cucu: cucuAngsuran},
		{tabel: models.TabelSpreading, kolomInduk: "POLIS_ID"},
		{tabel: models.TabelXOL, kolomInduk: "POLIS_ID", cucu: cucuLayer},
		{tabel: models.TabelSurvei, kolomInduk: "POLIS_ID"},
	}
)

// ------------------------------------------------------------------ simpan

// SimpanHalaman menulis seluruh halaman kerja satu kasus.
func (g *Gudang) SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	if err := models.PeriksaBentukSimpan(h); err != nil {
		return fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	if err := g.tulisInduk(ctx, tx, id, h); err != nil {
		return err
	}
	if err := g.hapusAnak(ctx, tx, id); err != nil {
		return err
	}
	if err := g.tulisQuotation(ctx, tx, id, h); err != nil {
		return err
	}
	for _, a := range tabelAnak {
		if err := g.sisipAnak(ctx, tx, a, id, h.AmbilDaftar(a.tabel.Daftar), h, a.tabel.Daftar); err != nil {
			return err
		}
	}
	return nil
}

// tulisInduk memperbarui T_GENERAL_POLIS_TREATY - hanya bila generasinya terbuka.
func (g *Gudang) tulisInduk(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	var set []string
	var args []any
	n := 1
	for _, k := range models.TabelGeneralPolis.Kolom {
		eks, jml := ekspresiTulis(k, n)
		v, err := nilaiTulis(k, h.Ambil(k.Properti))
		if err != nil {
			return err
		}
		set = append(set, k.Kolom+" = "+eks)
		args = append(args, v...)
		n += jml
	}
	q := fmt.Sprintf(`UPDATE %s g SET %s WHERE g.ID = :%d AND %s`, t, strings.Join(set, ", "), n, syaratTerbuka(t))
	hasil, err := jalankan(ctx, tx, "menyimpan polis", q, append(args, id)...)
	if err != nil {
		return err
	}
	if c, _ := hasil.RowsAffected(); c == 1 {
		return nil
	}
	// Nol baris: kasus tidak ada, atau generasinya tertutup.
	var tutup int
	qc := fmt.Sprintf(`SELECT CASE WHEN %s THEN 0 ELSE 1 END FROM %s g WHERE g.ID = :1`, syaratTerbuka(t), t)
	if err := db.PeriksaSQL(qc); err != nil {
		return err
	}
	switch err := tx.QueryRowContext(ctx, qc, id).Scan(&tutup); {
	case errors.Is(err, sql.ErrNoRows):
		return ErrKasusTidakAda
	case err != nil:
		return fmt.Errorf("repository: memeriksa generasi polis: %w", err)
	case tutup == 1:
		return ErrGenerasiTertutup
	}
	return fmt.Errorf("repository: penyimpanan polis %s tidak menyentuh baris", id)
}

// syaratTerbuka - generasi beralias `g` belum punya penerus (ID-10): tidak ada
// baris yang OLD_POLIS_ID-nya menunjuk g.ID. Pengganti kolom penanda tutup yang
// tidak ada di diagram grilling.
func syaratTerbuka(t string) string {
	return fmt.Sprintf(`NOT EXISTS (SELECT 1 FROM %s s WHERE s.OLD_POLIS_ID = g.ID)`, t)
}

func (g *Gudang) tulisQuotation(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t, err := g.nama(models.TabelQuotation.Nama)
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, "menghapus quotation", fmt.Sprintf(`DELETE FROM %s WHERE POLIS_ID = :1`, t), id); err != nil {
		return err
	}
	kolom := []string{"POLIS_ID"}
	nilai := []string{":1"}
	args := []any{id}
	n := 2
	for _, k := range models.TabelQuotation.Kolom {
		eks, jml := ekspresiTulis(k, n)
		v, err := nilaiTulis(k, models.NilaiQuotation(h, k.Properti))
		if err != nil {
			return err
		}
		kolom, nilai = append(kolom, k.Kolom), append(nilai, eks)
		args = append(args, v...)
		n += jml
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	_, err = jalankan(ctx, tx, "menyimpan quotation", q, args...)
	return err
}

// hapusAnak menghapus seluruh cucu lalu anak satu polis.
func (g *Gudang) hapusAnak(ctx context.Context, tx *db.Tx, id string) error {
	for _, a := range tabelAnak {
		t, err := g.nama(a.tabel.Nama)
		if err != nil {
			return err
		}
		if a.cucu != nil {
			c, err := g.nama(a.cucu.tabel.Nama)
			if err != nil {
				return err
			}
			q := fmt.Sprintf(`DELETE FROM %s WHERE %s IN (SELECT ID FROM %s WHERE %s = :1)`, c, a.cucu.kolomInduk, t, a.kolomInduk)
			if _, err := jalankan(ctx, tx, "menghapus "+a.cucu.tabel.Nama, q, id); err != nil {
				return err
			}
		}
		if _, err := jalankan(ctx, tx, "menghapus "+a.tabel.Nama,
			fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, t, a.kolomInduk), id); err != nil {
			return err
		}
	}
	return nil
}

// sisipAnak menyisipkan baris satu daftar (dan cucunya), NOURUT 1..n.
func (g *Gudang) sisipAnak(ctx context.Context, tx *db.Tx, a anakTabel, induk string,
	baris []models.Baris, h *models.Halaman, jalurDaftar string) error {

	if len(baris) == 0 {
		return nil
	}
	t, err := g.nama(a.tabel.Nama)
	if err != nil {
		return err
	}
	kolom := []string{"ID", a.kolomInduk, "NOURUT"}
	nilai := []string{":1", ":2", ":3"}
	n := 4
	for _, k := range a.tabel.Kolom {
		eks, jml := ekspresiTulis(k, n)
		kolom, nilai = append(kolom, k.Kolom), append(nilai, eks)
		n += jml
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	for i, b := range baris {
		idBaris, err := idBaru()
		if err != nil {
			return err
		}
		args := []any{idBaris, induk, i + 1}
		for _, k := range a.tabel.Kolom {
			v, err := nilaiTulis(k, b[k.Properti])
			if err != nil {
				return fmt.Errorf("%s baris %d: %w", jalurDaftar, i+1, err)
			}
			args = append(args, v...)
		}
		if _, err := jalankan(ctx, tx, "menyimpan "+a.tabel.Nama, q, args...); err != nil {
			return err
		}
		if a.cucu != nil {
			jalurCucu := models.JalurAnak(jalurDaftar, i+1, a.cucu.tabel.Daftar)
			if err := g.sisipAnak(ctx, tx, *a.cucu, idBaris, h.AmbilDaftar(jalurCucu), h, jalurCucu); err != nil {
				return err
			}
		}
	}
	return nil
}

// ------------------------------------------------------------------ nomor polis

// SetelNomorPolis menulis NOPOLIS SEKALI (AC 74): berkas yang sudah bernomor
// ditolak, dan indeks unik (NOPOLIS, PRODKE) menolak nomor kembar (AC 31).
func (g *Gudang) SetelNomorPolis(ctx context.Context, tx *db.Tx, id, nopol string) error {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyimpan nomor polis",
		fmt.Sprintf(`UPDATE %s g SET NOPOLIS = :1 WHERE g.ID = :2 AND g.NOPOLIS IS NULL AND %s`, t, syaratTerbuka(t)), nopol, id)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrNomorPolisSudahAda
	}
	return nil
}

// ------------------------------------------------------------------ muat

// BacaHalaman memuat halaman kerja satu kasus dari T_GENERAL_POLIS_TREATY dan
// anak-anaknya. Nilai halaman `TreatyIn` TIDAK dimuat di sini - dibaca ulang
// dari view lewat TREATY_IN_ID (services).
func (g *Gudang) BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error) {
	h := models.HalamanBaru()
	if err := g.bacaInduk(ctx, tx, id, h); err != nil {
		return nil, err
	}
	if err := g.bacaQuotation(ctx, tx, id, h); err != nil {
		return nil, err
	}
	for _, a := range tabelAnak {
		if err := g.bacaAnak(ctx, tx, a, id, h); err != nil {
			return nil, err
		}
	}
	// SuggestList - dari POOLDATA.HISTORYAKSEPTASIPRODUCTION (K4, usulan.go).
	if err := g.bacaUsulan(ctx, tx, id, h); err != nil {
		return nil, err
	}
	h.Setel("pyID", id)
	return h, nil
}

func (g *Gudang) bacaInduk(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	ks := models.TabelGeneralPolis.Kolom
	q := fmt.Sprintf(`SELECT NOPOLIS, %s FROM %s WHERE ID = :1`, daftarBaca(ks, ""), t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	nilai := make([]sql.NullString, len(ks)+1)
	tujuan := make([]any, len(nilai))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKasusTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: membaca polis: %w", err)
	}
	h.Setel(models.HalamanPolis+".PolicyNo", teks(nilai[0]))
	for i, k := range ks {
		if v := nilaiBaca(k, nilai[i+1]); v != "" {
			h.Setel(k.Properti, v)
		}
	}
	return nil
}

func (g *Gudang) bacaQuotation(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t, err := g.nama(models.TabelQuotation.Nama)
	if err != nil {
		return err
	}
	ks := models.TabelQuotation.Kolom
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE POLIS_ID = :1`, daftarBaca(ks, ""), t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	nilai := make([]sql.NullString, len(ks))
	tujuan := make([]any, len(nilai))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil // kasus yang belum pernah disimpan
	}
	if err != nil {
		return fmt.Errorf("repository: membaca quotation: %w", err)
	}
	for i, k := range ks {
		if v := nilaiBaca(k, nilai[i]); v != "" {
			h.Setel(models.HalamanQuotation+"."+k.Properti, v)
			h.Setel(models.HalamanPolis+".QuotationData."+k.Properti, v)
		}
	}
	return nil
}

// bacaAnak membaca satu daftar (berurut NOURUT) beserta cucunya.
func (g *Gudang) bacaAnak(ctx context.Context, tx *db.Tx, a anakTabel, id string, h *models.Halaman) error {
	t, err := g.nama(a.tabel.Nama)
	if err != nil {
		return err
	}
	ks := a.tabel.Kolom
	q := fmt.Sprintf(`SELECT ID, %s FROM %s WHERE %s = :1 ORDER BY NOURUT`, daftarBaca(ks, ""), t, a.kolomInduk)
	baris, kunci, err := g.bacaBaris(ctx, tx, q, ks, id)
	if err != nil {
		return err
	}
	if len(baris) == 0 {
		return nil
	}
	h.SetelDaftar(a.tabel.Daftar, baris)
	if a.cucu == nil {
		return nil
	}
	c, err := g.nama(a.cucu.tabel.Nama)
	if err != nil {
		return err
	}
	for i, k := range kunci {
		qc := fmt.Sprintf(`SELECT ID, %s FROM %s WHERE %s = :1 ORDER BY NOURUT`, daftarBaca(a.cucu.tabel.Kolom, ""), c, a.cucu.kolomInduk)
		cucu, _, err := g.bacaBaris(ctx, tx, qc, a.cucu.tabel.Kolom, k)
		if err != nil {
			return err
		}
		if len(cucu) > 0 {
			h.SetelDaftar(models.JalurAnak(a.tabel.Daftar, i+1, a.cucu.tabel.Daftar), cucu)
		}
	}
	return nil
}

func (g *Gudang) bacaBaris(ctx context.Context, tx *db.Tx, q string, ks []models.Kolom, arg string) ([]models.Baris, []string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, nil, err
	}
	rows, err := g.pembaca(tx).QueryContext(ctx, q, arg)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: membaca baris anak: %w", err)
	}
	defer rows.Close()
	var out []models.Baris
	var kunci []string
	for rows.Next() {
		var idBaris string
		nilai := make([]sql.NullString, len(ks))
		tujuan := []any{&idBaris}
		for i := range nilai {
			tujuan = append(tujuan, &nilai[i])
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, nil, fmt.Errorf("repository: membaca baris anak: %w", err)
		}
		b := models.Baris{}
		for i, k := range ks {
			if v := nilaiBaca(k, nilai[i]); v != "" {
				b[k.Properti] = v
			}
		}
		out = append(out, b)
		kunci = append(kunci, idBaris)
	}
	return out, kunci, rows.Err()
}
