package repository

// View Detail peserta POLIS di tab Policy Details layar Outstanding - keputusan work owner 11-10-2026 ("iya buat saja").
//
// Korpus `kelvin/InboxClaimLife (Done)`:
//
//	Section/DetailPolisLife.xml   tombol View Detail b631 -> LoadDataPeserta_Act b650; grid `TempDetail.pxResults` b4537:
//	                              No. .URUT b4700/b6231, POLICY NO b4838, CERTIFICATE NO b4981, NAME OF INSURED b5126,
//	                              SEX b5266, DATE OF BIRTH b5402, AGE b5538, PLAN b5674, BEGIN DATE b5810,
//	                              EFFECTIVE DATE b5946, EXPIRED DATE b6082; NextPrev ⟪ ⟨ ⟩ ⟫ b1216/b1489/b2877/b3106.
//	Activity/LoadDataPeserta_Act  halaman = PagiNasi.CARI1, 500 baris (b572 CARI11 = CARI1*500, b619 CARI10 = CARI11-500);
//	                              RDB-List GetPesertaClaim_sql b892.
//	RDBList/GetPesertaClaim_sql   `select rownum urut, a.* from m_life_premium_detail a where pl_number = {PL}
//	                              order by BEGIN_DATE, CERTIFICATE_NO, NAME_OF_INSURED desc` lalu urut > CARI10 dan <= CARI11.
//	RDBList/getMaxPagination_sql  jumlah halaman = CEIL(COUNT / 500).
//
// ⛔ PENYIMPANGAN SADAR, dicatat:
//   - Saringan versi: Pega HANYA `PL_NUMBER` (semua versi polis, termasuk baris Batal/Delete). Di sini `PL_NUMBER`
//     (berindeks INDEX4) DAN `IDPEGA` = versi polis yang SAMA dengan header klaim (arahan 11-10-2026). Tanpa saringan
//     status EDM (Batal/Delete) - LoadDataPeserta_Act tidak menyaringnya.
//   - Urutan: Pega memberi ROWNUM SEBELUM ORDER BY (Oracle), sehingga halamannya tidak benar-benar terurut. Di sini
//     urutan Pega diterapkan sungguh-sungguh, dengan ID sebagai pemutus seri supaya halaman stabil.
//   - Tanggal: Pega menambah 7 jam untuk tampilan (b1097 dst, koreksi zona waktu). Di sini DATE dibaca sebagai tanggal
//     kalender lewat TO_CHAR - hasil tampilnya sama.
//
// ⛔ Hanya MEMBACA. Nama tertanggung tidak disalin ke mana pun (kolomSalin tetap tanpa nama).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// UkuranHalamanPesertaPolis - 500 baris per halaman (LoadDataPeserta_Act b572/b619, getMaxPagination_sql).
const UkuranHalamanPesertaPolis = 500

// BarisPesertaPolis adalah satu baris grid View Detail.
type BarisPesertaPolis struct {
	Urut            int    `json:"urut"`
	NomorPolis      string `json:"nomorPolis"`
	NomorSertifikat string `json:"nomorSertifikat"`
	Nama            string `json:"namaTertanggung"`
	JenisKelamin    string `json:"jenisKelamin"`
	TanggalLahir    string `json:"tanggalLahir"`
	Umur            string `json:"umur"`
	Plan            string `json:"plan"`
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalEfektif  string `json:"tanggalEfektif"`
	TanggalExpired  string `json:"tanggalExpired"`
}

// HalamanPesertaPolis adalah satu halaman grid View Detail.
type HalamanPesertaPolis struct {
	Baris       []BarisPesertaPolis `json:"baris"`
	Halaman     int                 `json:"halaman"`
	Ukuran      int                 `json:"ukuran"`
	JumlahBaris int                 `json:"total"`
}

// sqlPesertaPolisDetail - satu halaman peserta versi polis yang sama dengan header klaim. Penyaring berindeks
// PL_NUMBER = :1, versi IDPEGA = :2, OFFSET :3, FETCH FIRST :4 (batas hasil). Jumlah baris lewat COUNT(*) OVER ().
func sqlPesertaPolisDetail(sumber string) string {
	return fmt.Sprintf(`SELECT POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, SEX, TO_CHAR(DOB, 'YYYY-MM-DD'),
	       TO_CHAR(AGE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), PLAN,
	       TO_CHAR(BEGIN_DATE, 'YYYY-MM-DD'), TO_CHAR(EFFECTIVE_DATE, 'YYYY-MM-DD'), TO_CHAR(EXPIRED_DATE, 'YYYY-MM-DD'),
	       COUNT(*) OVER ()
	  FROM %s
	 WHERE PL_NUMBER = :1 AND IDPEGA = :2
	 ORDER BY BEGIN_DATE, CERTIFICATE_NO, NAME_OF_INSURED DESC, ID
	 OFFSET :3 ROWS FETCH FIRST :4 ROWS ONLY`, sumber)
}

// PesertaPolisDetail membaca satu halaman (1-based) grid View Detail. PL atau IDPEGA kosong = halaman kosong (tidak
// memindai tabel 66,8 juta baris tanpa penyaring).
func (r *KlaimLife) PesertaPolisDetail(ctx context.Context, plNumber, idPega string, halaman int) (HalamanPesertaPolis, error) {
	if halaman < 1 {
		halaman = 1
	}
	out := HalamanPesertaPolis{Baris: []BarisPesertaPolis{}, Halaman: halaman, Ukuran: UkuranHalamanPesertaPolis}
	if strings.TrimSpace(plNumber) == "" || strings.TrimSpace(idPega) == "" {
		return out, nil
	}
	sumber, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return out, err
	}
	q := sqlPesertaPolisDetail(sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return out, err
	}
	awal := (halaman - 1) * UkuranHalamanPesertaPolis
	rows, err := r.db.QueryContext(ctx, q, plNumber, idPega, awal, UkuranHalamanPesertaPolis)
	if err != nil {
		return out, fmt.Errorf("repository: membaca peserta polis: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var k [10]sql.NullString
		var total int
		if err := rows.Scan(&k[0], &k[1], &k[2], &k[3], &k[4], &k[5], &k[6], &k[7], &k[8], &k[9], &total); err != nil {
			return out, fmt.Errorf("repository: membaca satu peserta polis: %w", err)
		}
		out.JumlahBaris = total
		out.Baris = append(out.Baris, BarisPesertaPolis{
			Urut: awal + len(out.Baris) + 1, NomorPolis: k[0].String, NomorSertifikat: k[1].String, Nama: k[2].String,
			JenisKelamin: k[3].String, TanggalLahir: k[4].String, Umur: k[5].String, Plan: k[6].String,
			TanggalMulai: k[7].String, TanggalEfektif: k[8].String, TanggalExpired: k[9].String,
		})
	}
	return out, rows.Err()
}
