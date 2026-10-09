// Package repository adalah SATU-SATUNYA lapisan modul Claim Prop yang berbicara ke Oracle. Setiap nama tabel lewat
// `db.Qualify` (skema eksplisit), setiap teks SQL lewat `db.PeriksaSQL` (nol COMMIT), nol stored procedure dipanggil
// (prompt implementasi §3: procedure Pega ditulis ulang sebagai SQL langsung di transaksi aplikasi).
//
// Untuk apa berkas ini: KONVERSI KOLOM yang digerakkan katalog (`models/katalog.go`). Konversi terjadi SEKALI, di
// sini:
//
//	uang, persen   teks desimal  <-> NUMBER(38,10) lewat koefisien / 10^skala
//	tanggal        "2006-01-02"  <-> DATE
//	tanggal-waktu  "2006-01-02 15:04:05" <-> DATE
//	kode, penanda, teks           <-> VARCHAR2 apa adanya ("006" tetap "006")
//
// Pola DISALIN dari modul NB Treaty In (`repository/kolom.go`), bukan diimpor.
//
// ⛔ Teks kosong pada kolom angka atau tanggal ditulis NULL, bukan 0. Nol float di seluruh paket.
package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimprop/backend/models"
)

// fmtTanggal - SATU format pertukaran tanggal dengan Oracle, pasangan `utils.TanggalWaktu`.
const fmtTanggal = "YYYY-MM-DD HH24:MI:SS"

// pecahAngka memecah desimal menjadi koefisien bulat dan skala (`TO_NUMBER(:k) / POWER(10, :s)` - eksak, tidak
// bergantung pemisah desimal sesi).
func pecahAngka(d *apd.Decimal) (string, int64) {
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return tanda + koef, int64(len(pecahan))
}

// ekspresiTulis - ekspresi SQL penampung satu kolom bernomor `n` (dan `n+1` untuk angka).
func ekspresiTulis(k models.Kolom, n int) (string, int) {
	switch {
	case k.Golongan.Desimal():
		return fmt.Sprintf("(TO_NUMBER(:%d) / POWER(10, :%d))", n, n+1), 2
	case k.Golongan.Tanggal():
		return fmt.Sprintf("TO_DATE(:%d, '%s')", n, fmtTanggal), 1
	}
	return fmt.Sprintf(":%d", n), 1
}

// nilaiTulis mengubah teks halaman menjadi argumen bind satu kolom. Galat masukan = `galat.ErrPermintaanTidakSah`
// (400): yang salah isian layar, bukan basis data.
func nilaiTulis(k models.Kolom, teks string) ([]any, error) {
	s := strings.TrimSpace(teks)
	switch {
	case k.Golongan.Desimal():
		if s == "" {
			return []any{nil, nil}, nil
		}
		d, err := utils.ParseDecimal(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %w: %s = %q", galat.ErrPermintaanTidakSah, models.ErrBukanAngka, k.Properti, teks)
		}
		koef, skala := pecahAngka(d)
		return []any{koef, skala}, nil
	case k.Golongan.Tanggal():
		if s == "" {
			return []any{nil}, nil
		}
		t, ok := models.UraiTanggal(s)
		if !ok {
			return nil, fmt.Errorf("%w: %s = %q bukan tanggal", galat.ErrPermintaanTidakSah, k.Properti, teks)
		}
		return []any{t.Format(utils.TanggalWaktu)}, nil
	}
	if teks == "" {
		return []any{nil}, nil
	}
	if k.Panjang > 0 && len([]rune(teks)) > k.Panjang {
		return nil, fmt.Errorf("%w: %s melebihi %d karakter", galat.ErrPermintaanTidakSah, k.Properti, k.Panjang)
	}
	return []any{teks}, nil
}

// ekspresiBaca - ekspresi SELECT satu kolom: angka dan tanggal dibaca sebagai TEKS.
func ekspresiBaca(k models.Kolom, alias string) string {
	kol := alias + k.Kolom
	switch {
	case k.Golongan.Desimal():
		return fmt.Sprintf(db.FmtDesimal, kol)
	case k.Golongan.Tanggal():
		return fmt.Sprintf("TO_CHAR(%s, '%s')", kol, fmtTanggal)
	}
	return kol
}

// nilaiBaca mengubah hasil kolom menjadi teks halaman.
func nilaiBaca(k models.Kolom, v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	s := v.String
	switch {
	case k.Golongan.Desimal():
		s = rapikanDesimal(s)
	case k.Golongan == models.GolTanggal:
		if len(s) >= 10 {
			s = s[:10]
		}
	}
	return s
}

// rapikanDesimal - TO_CHAR TM9 menulis ".5" untuk 0,5; halaman memakai "0.5".
func rapikanDesimal(s string) string {
	switch {
	case strings.HasPrefix(s, "."):
		return "0" + s
	case strings.HasPrefix(s, "-."):
		return "-0" + s[1:]
	}
	return s
}

// daftarBaca merangkai ekspresi SELECT.
func daftarBaca(ks []models.Kolom, alias string) string {
	b := make([]string, len(ks))
	for i, k := range ks {
		b[i] = ekspresiBaca(k, alias)
	}
	return strings.Join(b, ", ")
}

// teksAtauNil - VARCHAR2 kosong = NULL.
func teksAtauNil(s string) any { return db.KosongJadiNil(s) }

// angkaAtauNil - bind angka eksak (dua penampung) dari teks; kosong = NULL.
func angkaAtauNil(jalur, s string) ([]any, error) {
	return nilaiTulis(models.Kolom{Properti: jalur, Golongan: models.GolUang}, s)
}
