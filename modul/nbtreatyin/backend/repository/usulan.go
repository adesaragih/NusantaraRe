package repository

// Untuk apa berkas ini: CATATAN USULAN -> tabel lama
// `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (keputusan work owner K4 03-10-2026;
// spec-penyimpanan ID-5, ID-31, AC 39-44; diagram grilling "1:N IDPEGA ·
// 15 kolom · SUDAH datar"). Pengganti `RDBList/InsertViewSuggest_SQL`:
//
//	INSERT INTO POOLDATA.historyakseptasiproduction (IDPEGA, TYPE_POLIS, NOURUT,
//	  POSISI, PIC, TGL_INP, DIV, TYPE, PUTARAN, APPROVAL, KETERANGAN,
//	  AKSES_LOGIN, B2B, BUSINESS_CODE, PERCENT_RNM)
//	VALUES ({pyWorkPage.pzInsKey}, {CARI1}, {CARI2}, {CARI3}, {CARI4},
//	  To_date({CARI5}, 'DD/MM/YYYY HH24:MI:SS'), {CARI6}, {CARI7}, {CARI8},
//	  {CARI9}, substr({CARI10},0,3990), {OperatorID.pyUserIdentifier},
//	  {OfferFacIn.IsB2B}, {Quotation.BusinessCode}, {OfferFacIn.PercentShare});
//	COMMIT;
//
// ⛔ TABEL WARISAN - tidak dibuat, tidak diubah strukturnya (MODUL.md "Tabel
// warisan"). ⛔ Tanpa COMMIT: ditulis di transaksi submit (AC 83); skema
// eksplisit lewat `Qualify` (AC 47). Pemetaan nilai: `models/usulan.go`.
//
// ⚠️ Tipe kolom NOURUT belum dicek ke katalog Oracle (butir terbuka): dibaca
// dan diurutkan lewat `TO_NUMBER(NOURUT)` - benar untuk NUMBER maupun
// VARCHAR2 berisi angka; modul ini satu-satunya penulis baris kasus treaty
// (kalang Pega bersyarat BusinessFac "F", tidak pernah menulis treaty).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

const tabelRiwayatProduksi = "HISTORYAKSEPTASIPRODUCTION"

func sqlNourutUsulan(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(NVL(MAX(TO_NUMBER(NOURUT)), 0)) FROM %s WHERE IDPEGA = :1`, t)
}

func sqlSisipUsulan(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (IDPEGA, TYPE_POLIS, NOURUT, POSISI, PIC, TGL_INP, DIV, TYPE, PUTARAN,
		APPROVAL, KETERANGAN, AKSES_LOGIN, B2B, BUSINESS_CODE, PERCENT_RNM)
		VALUES (:1, :2, :3, :4, :5, TO_DATE(:6, '%s'), :7, :8, :9, :10, SUBSTR(:11, 1, %d), :12, :13, :14, :15)`,
		t, fmtTanggal, models.PanjangKeterangan)
}

func sqlBacaUsulan(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(NOURUT), PIC, TO_CHAR(TGL_INP, '%s'), APPROVAL, KETERANGAN, AKSES_LOGIN
	  FROM %s WHERE IDPEGA = :1 ORDER BY TO_NUMBER(NOURUT)`, fmtTanggal, t)
}

// CatatUsulan menulis catatan baru satu kasus. NOURUT = berikutnya per IDPEGA,
// dibaca di transaksi yang SAMA sesudah kasus dikunci (`KunciKasus` FOR
// UPDATE) - dua submit serentak tidak berbagi nomor.
func (g *Gudang) CatatUsulan(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) error {
	if len(baris) == 0 {
		return nil
	}
	t, err := g.nama(tabelRiwayatProduksi)
	if err != nil {
		return err
	}
	qn := sqlNourutUsulan(t)
	if err := db.PeriksaSQL(qn); err != nil {
		return err
	}
	var akhir sql.NullString
	if err := tx.QueryRowContext(ctx, qn, idPega).Scan(&akhir); err != nil {
		return fmt.Errorf("repository: membaca NOURUT riwayat produksi: %w", err)
	}
	n, err := strconv.Atoi(akhir.String)
	if err != nil {
		return fmt.Errorf("repository: NOURUT riwayat produksi %q bukan bilangan: %w", akhir.String, err)
	}
	q := sqlSisipUsulan(t)
	for i, u := range baris {
		hasil, err := jalankan(ctx, tx, "menulis riwayat produksi", q,
			idPega, db.KosongJadiNil(u.TypePolis), n+i+1, db.KosongJadiNil(u.Posisi), db.KosongJadiNil(u.PIC),
			db.KosongJadiNil(u.TglInp), db.KosongJadiNil(u.Div), db.KosongJadiNil(u.Type), db.KosongJadiNil(u.Putaran),
			db.KosongJadiNil(u.Approval), db.KosongJadiNil(u.Keterangan), db.KosongJadiNil(u.AksesLogin),
			db.KosongJadiNil(u.B2B), db.KosongJadiNil(u.BusinessCode), db.KosongJadiNil(u.PercentRNM))
		if err != nil {
			return err
		}
		if err := db.PastikanSatuBaris(hasil, "penulisan riwayat produksi"); err != nil {
			return err
		}
	}
	return nil
}

// bacaUsulan mengisi SuggestList halaman dari riwayat produksi kasus itu
// (riwayat catatan layar `Section/ListSuggest`), berurut NOURUT.
func (g *Gudang) bacaUsulan(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t, err := g.nama(tabelRiwayatProduksi)
	if err != nil {
		return err
	}
	q := sqlBacaUsulan(t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	rows, err := g.pembaca(tx).QueryContext(ctx, q, models.KunciInstans(id))
	if err != nil {
		return fmt.Errorf("repository: membaca riwayat produksi: %w", err)
	}
	defer rows.Close()
	var catatan []models.Baris
	for rows.Next() {
		var v [6]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5]); err != nil {
			return fmt.Errorf("repository: membaca riwayat produksi: %w", err)
		}
		no, _ := strconv.Atoi(v[0].String)
		catatan = append(catatan, models.BarisCatatan(models.UsulanProduksi{NoUrut: no, PIC: v[1].String,
			TglInp: v[2].String, Approval: v[3].String, Keterangan: v[4].String, AksesLogin: v[5].String}))
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("repository: membaca riwayat produksi: %w", err)
	}
	h.SetelDaftar(models.DaftarUsulan, catatan)
	return nil
}
