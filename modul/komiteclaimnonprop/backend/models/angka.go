package models

// Untuk apa berkas ini: ANGKA UANG - pembaca angka halaman dan nomor akseptasi (`KomitePostAdjustment` S14.8-S14.11). Uang lewat `apd.Decimal`, nol float, nol pembulatan di tengah hitungan (spec bab 7
// "ketelitian angka di sistem baru"); tampilan 4 desimal urusan layar.

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// konteksUang - presisi lebar: perkalian dan penjumlahan tidak dibulatkan di tengah.
var konteksUang = apd.BaseContext.WithPrecision(60)

// desimal - teks angka halaman (kosong = 0). Bukan angka -> galat (bukan ditebak).
func desimal(jalur, s string) (*apd.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return apd.New(0, 0), nil
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return nil, fmt.Errorf("models: %s bukan angka (%q): %w", jalur, s, err)
	}
	return d, nil
}

// TeksAngka - angka sebagai teks halaman (tanpa nol ekor).
func TeksAngka(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	r := new(apd.Decimal)
	r.Reduce(d)
	return utils.FormatDecimal(r)
}

// lpad5 = LPAD(v_seq, 5, '0') (`PROC_GENERATE_SEQUENCE_NUMBER` keluaran `p_seq_number`).
func lpad5(n int) string { return fmt.Sprintf("%05d", n) }

// HurufAkseptasi - S14.9 `ParamSeq.CARI2 := ParamSeq.HASIL3 + "A"`.
const HurufAkseptasi = "A"

// RakitNomorAkseptasi = S14.11 `CARI2 + BusinessOldId + "." + HASIL1 + ".TX" + HASIL2`: jenis (kode produksi NONLIFE +
// "A"), kode lama bisnis, periode MM.YYYY, urut 5 digit.
func RakitNomorAkseptasi(jenis, oldID, mmyyyy string, urut int) string {
	return jenis + oldID + "." + mmyyyy + ".TX" + lpad5(urut)
}
