package repository

// Diagnosa per peserta - butir bd, tabel `T_CLAIMLF_DIAGNOSE` (migrasi 018).
//
// Untuk apa berkas ini: baca dan tulis `.DiagnoseList`. Berbeda dengan
// `penyakit.go` - yang di sana KATALOG warisan yang dibaca saja, yang di
// sini milik klaim dan ditulis.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033), nol
// `COMMIT` (ADR-U-0029), dan identitas dari `SEQ_CLAIMLF_DIAGNOSE`
// (ADR-0006).
//
// Dibaca sesudah: klaimlife.go - `AmbilDiagnosa` mengikuti bentuk
// `AmbilDokumen` persis, sebab keduanya daftar milik PESERTA yang dibaca
// sekali untuk seluruh klaim.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/internal/models"
	"nusantarare/inti/db"
)

// ErrDiagnosaTidakAda - baris diagnosa yang diminta bukan milik klaim itu.
var ErrDiagnosaTidakAda = errors.New("repository: diagnosa tidak ada pada klaim ini")

// Diagnosa membaca dan menulis daftar diagnosa peserta.
type Diagnosa struct{ db *db.DB }

// NewDiagnosa menyusunnya.
func NewDiagnosa(db *db.DB) *Diagnosa { return &Diagnosa{db: db} }

// tabelDiagnosa dan tabelPesertaDiagnosa - dua nama yang selalu dipakai
// berpasangan, sebab setiap pertanyaan tentang diagnosa dibatasi KLAIM-nya
// lewat pesertanya.
func (r *Diagnosa) tabel() (string, string, error) {
	diag, err := r.db.Qualify("T_CLAIMLF_DIAGNOSE")
	if err != nil {
		return "", "", err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return "", "", err
	}
	return diag, pes, nil
}

// sqlAmbilDiagnosa merakit pembacaan seluruh diagnosa satu klaim.
//
// ⛔ `ORDER BY d.URUTAN, d.ID` - dan keduanya, bukan salah satu. `URUTAN`
// adalah yang dilihat pemakai; `ID` memutus seri bila dua baris sempat
// berurutan sama, sehingga daftar yang dibaca dua kali tidak pernah berbeda
// susunannya. Grid tanpa urutan stabil menampilkan baris yang berpindah
// sendiri, dan orang akan mengira datanya berubah.
func sqlAmbilDiagnosa(diag, pes string) string {
	return fmt.Sprintf(
		`SELECT d.PREMIUM_LIST_DETAIL_ID, d.ID, d.URUTAN,
		        d.ICD_CODE, d.DISEASE, d.GROUP_DIAGNOSE, d.STS_REJECT
		   FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL
		  ORDER BY d.PREMIUM_LIST_DETAIL_ID, d.URUTAN, d.ID`, diag, pes)
}

// AmbilDiagnosa membaca SELURUH diagnosa satu klaim, dikelompokkan per
// peserta.
//
// Satu query untuk seluruh klaim, bukan satu per peserta: klaim grup dapat
// berpeserta ratusan, dan N+1 query atas layar yang dibuka setiap kali orang
// menekan sebuah baris adalah cara termudah membuat layar terasa rusak.
func (r *Diagnosa) AmbilDiagnosa(ctx context.Context, klaimID string) (
	map[string][]models.Diagnosa, error) {

	diag, pes, err := r.tabel()
	if err != nil {
		return nil, err
	}
	q := sqlAmbilDiagnosa(diag, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca diagnosa klaim: %w", err)
	}
	defer baris.Close()

	keluar := map[string][]models.Diagnosa{}
	for baris.Next() {
		var (
			pesertaID            sql.NullString
			id                   sql.NullInt64
			urutan               sql.NullInt64
			icd, nama, grup, sts sql.NullString
		)
		if err := baris.Scan(&pesertaID, &id, &urutan, &icd, &nama, &grup,
			&sts); err != nil {
			return nil, fmt.Errorf("repository: memindai diagnosa: %w", err)
		}
		d := models.Diagnosa{
			ID:            id.Int64,
			PesertaID:     pesertaID.String,
			Urutan:        int(urutan.Int64),
			KodeICD:       icd.String,
			Nama:          nama.String,
			GroupDiagnose: grup.String,
			KodeStatus:    sts.String,
		}
		keluar[d.PesertaID] = append(keluar[d.PesertaID], d)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca diagnosa klaim: %w", err)
	}
	return keluar, nil
}

// sqlSisipDiagnosa merakit penyisipan satu baris kosong.
//
// ⛔ BARIS KOSONG, dan itu meniru `Add` b4700 `addRow` apa adanya: di Pega
// tombol itu menambahkan anggota PageList yang propertinya belum terisi -
// pemakai lalu menekan `Find Disease` untuk mengisinya. Menuntut ICD dan
// nama pada saat menambah akan membalik urutan kerja orang.
//
// ⚠️ `STS_REJECT` ikut diisi dari peserta SEJAK LAHIR, bukan dibiarkan NULL.
// `SetSTS_Reject` hanya berjalan ketika keputusan peserta berubah; baris yang
// lahir SESUDAH itu tidak akan pernah disentuhnya, dan akan tampak belum
// diputus padahal pesertanya sudah. (Di jalur nyata gerbang `DiagnosaTerkunci`
// mencegah kelahiran itu - tetapi penjaga yang bersandar pada penjaga lain
// adalah satu penjaga, bukan dua.)
func sqlSisipDiagnosa(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(ID, PREMIUM_LIST_DETAIL_ID, URUTAN, STS_REJECT)
		VALUES (:1,:2,:3,:4)`, tabel)
}

// SisipDiagnosa menambahkan satu baris kosong di ekor daftar peserta.
//
// Mengembalikan baris yang baru lahir, lengkap dengan pengenal dan urutannya
// - supaya layar dapat menampilkannya tanpa membaca ulang seluruh klaim.
func (r *Diagnosa) SisipDiagnosa(ctx context.Context, tx *db.Tx,
	pesertaID string, urutan int, stsPeserta string) (models.Diagnosa, error) {

	tabel, _, err := r.tabel()
	if err != nil {
		return models.Diagnosa{}, err
	}
	nomor, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_DIAGNOSE")
	if err != nil {
		return models.Diagnosa{}, err
	}
	var id int64
	if _, err := fmt.Sscan(nomor, &id); err != nil {
		return models.Diagnosa{}, fmt.Errorf(
			"repository: nomor diagnosa %q bukan angka: %w", nomor, err)
	}
	q := sqlSisipDiagnosa(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Diagnosa{}, err
	}
	hasil, err := tx.ExecContext(ctx, q, id, pesertaID, urutan,
		db.KosongJadiNil(stsPeserta))
	if err != nil {
		return models.Diagnosa{}, fmt.Errorf("repository: menyisip diagnosa: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "penyisipan diagnosa"); err != nil {
		return models.Diagnosa{}, err
	}
	return models.Diagnosa{
		ID: id, PesertaID: pesertaID, Urutan: urutan, KodeStatus: stsPeserta,
	}, nil
}

// sqlPerbaruiDiagnosa merakit pembaruan tiga kolom isi.
//
// ⛔ `URUTAN` dan `STS_REJECT` TIDAK ikut. Yang pertama milik `Add`/`Delete`,
// yang kedua milik `SetSTS_Reject` - dan rute `Choose` yang diam-diam dapat
// menulis keduanya adalah rute yang dapat membatalkan keputusan peserta
// tanpa siapa pun memintanya.
func sqlPerbaruiDiagnosa(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
		   SET ICD_CODE = :1, DISEASE = :2, GROUP_DIAGNOSE = :3
		 WHERE ID = :4`, tabel)
}

// PerbaruiDiagnosa menulis hasil `Choose` ke satu baris.
//
// `[terverifikasi]` `Activity/SetDisease.xml` b260-261 `.DISEASE =
// Param.Disease` dan b307-308 `.ICDCODE = Param.ICD_Code`. Kolom ketiga
// (`GROUP_DIAGNOSE`) BUKAN dari rule itu - ia datang dari dropdown b5860
// yang mem-`postValue` b5886 sendiri. Satu rute menulis ketiganya sebab di
// layar keduanya menyunting baris yang sama.
func (r *Diagnosa) PerbaruiDiagnosa(ctx context.Context, tx *db.Tx,
	id int64, kodeICD, nama, grup string) error {

	tabel, _, err := r.tabel()
	if err != nil {
		return err
	}
	q := sqlPerbaruiDiagnosa(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(kodeICD),
		db.KosongJadiNil(nama), db.KosongJadiNil(grup), id)
	if err != nil {
		return fmt.Errorf("repository: memperbarui diagnosa: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "pembaruan diagnosa")
}

// HapusDiagnosa menghapus satu baris.
func (r *Diagnosa) HapusDiagnosa(ctx context.Context, tx *db.Tx, id int64) error {
	tabel, _, err := r.tabel()
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: menghapus diagnosa: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penghapusan diagnosa")
}

// sqlRapatkanUrutan merakit perapatan nomor urut sesudah penghapusan.
//
// ⛔ Satu pernyataan, bukan satu per baris. `deleteRow` b6170 di Pega
// menggeser subscript anggota sesudahnya seketika; menirunya dengan N
// UPDATE berarti daftar berlubang selama N-1 pernyataan, dan pembaca yang
// datang di tengahnya melihat urutan yang tidak pernah ada.
func sqlRapatkanUrutan(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET URUTAN = URUTAN - 1
		 WHERE PREMIUM_LIST_DETAIL_ID = :1 AND URUTAN > :2`, tabel)
}

// RapatkanUrutan menutup lubang yang ditinggalkan satu penghapusan.
func (r *Diagnosa) RapatkanUrutan(ctx context.Context, tx *db.Tx,
	pesertaID string, urutanTerhapus int) error {

	tabel, _, err := r.tabel()
	if err != nil {
		return err
	}
	q := sqlRapatkanUrutan(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	// ⚠️ Nol baris terpengaruh adalah keadaan yang SAH: yang dihapus baris
	// terakhir. Karena itu `pastikanSatuBaris` sengaja tidak dipakai di sini.
	if _, err := tx.ExecContext(ctx, q, pesertaID, urutanTerhapus); err != nil {
		return fmt.Errorf("repository: merapatkan urutan diagnosa: %w", err)
	}
	return nil
}

// sqlCerminkanStsReject merakit pencerminan keputusan peserta.
func sqlCerminkanStsReject(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET STS_REJECT = :1
		 WHERE PREMIUM_LIST_DETAIL_ID = :2`, tabel)
}

// CerminkanStsReject menyalin `STS_REJECT` peserta ke SETIAP diagnosanya.
//
// `[terverifikasi]` `Activity/SetSTS_Reject.xml` - kelas b67
// `Int-LIFE_PREMIUM_DETAIL` (peserta), halaman langkah b241 `.DiagnoseList`,
// b257-258 `.STS_REJECT = Primary.STS_REJECT`, berulang b345 `EMBEDDED`,
// prasyarat b325 `WhenTrue` 2 / `WhenFalse` 2 - keduanya LANJUT, jadi tanpa
// gerbang apa pun.
//
// ⛔ SATU pernyataan untuk seluruh daftar, dan itu bukan sekadar optimasi:
// rule Pega memutarinya sebagai satu langkah, dan setengah daftar yang
// tercermin adalah keadaan yang tidak pernah ada di sistem lama.
//
// ⚠️ Nol baris terpengaruh SAH - peserta boleh tidak punya diagnosa sama
// sekali. Di Pega pun putaran atas PageList kosong berjalan nol kali tanpa
// mengeluh.
func (r *Diagnosa) CerminkanStsReject(ctx context.Context, tx *db.Tx,
	pesertaID, sts string) error {

	tabel, _, err := r.tabel()
	if err != nil {
		return err
	}
	q := sqlCerminkanStsReject(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, db.KosongJadiNil(sts), pesertaID); err != nil {
		return fmt.Errorf("repository: mencerminkan STS_REJECT ke diagnosa: %w", err)
	}
	return nil
}
