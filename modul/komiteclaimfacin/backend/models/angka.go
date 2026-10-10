package models

// Untuk apa berkas ini: ANGKA UANG dan NOMOR AKSEPTASI (`KomitePost_Adjustment` S7.2.1.2.4-S7.2.1.2.7). Uang lewat
// `apd.Decimal`, nol float, nol pembulatan di tengah hitungan; tampilan 4 desimal urusan layar.

import (
	"errors"
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

// hitung - akumulator galat aritmetika (galat pertama menang).
type hitung struct{ err error }

func (h *hitung) dari(jalur, s string) *apd.Decimal {
	if h.err != nil {
		return apd.New(0, 0)
	}
	d, err := desimal(jalur, s)
	if err != nil {
		h.err = err
		return apd.New(0, 0)
	}
	return d
}

func (h *hitung) tambah(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	if h.err == nil {
		_, h.err = konteksUang.Add(r, a, b)
	}
	return r
}

func (h *hitung) kali(a, b *apd.Decimal) *apd.Decimal {
	r := new(apd.Decimal)
	if h.err == nil {
		_, h.err = konteksUang.Mul(r, a, b)
	}
	return r
}

// errBagiNol - pembagian dengan nol (`@divide` Pega menolak).
var errBagiNol = errors.New("models: pembagian dengan nol")

// bagiBulat = `@divide(a, b, n)` Pega: a / b dibulatkan n desimal (HALF_UP).
func (h *hitung) bagiBulat(a, b *apd.Decimal, n int32) *apd.Decimal {
	r := new(apd.Decimal)
	if h.err != nil {
		return r
	}
	if b.IsZero() {
		h.err = errBagiNol
		return r
	}
	c := *konteksUang
	if _, h.err = c.Quo(r, a, b); h.err != nil {
		return r
	}
	c.Rounding = apd.RoundHalfUp
	_, h.err = c.Quantize(r, r, -n)
	return r
}

// lebih - a > b.
func lebih(a, b *apd.Decimal) bool { return a.Cmp(b) > 0 }

// lpad5 = LPAD(v_seq, 5, '0') (`PROC_GENERATE_SEQUENCE_NUMBER` keluaran `p_seq_number`).
func lpad5(n int) string { return fmt.Sprintf("%05d", n) }

// HurufAkseptasi - S7.2.1.2.5 `ParamSeq.CARI2 := ParamSeq.HASIL3 + "A"`.
const HurufAkseptasi = "A"

// RakitNomorAkseptasi = S7.2.1.2.7 `ParamSeq.CARI2 + BusinessOldId + "." + HASIL1 + "." + HASIL2`: jenis (kode produksi
// NONLIFE + "A"), kode lama bisnis polis, periode MM.YYYY, urut 5 digit. TANPA ".TP" (Komite Prop) / ".TX" (Komite Non
// Prop) - panjang 21 / 22 = gerbang kasir jalur CLM (`HitServiceToKasirKMT_Act` S14.2).
func RakitNomorAkseptasi(jenis, oldID, mmyyyy string, urut int) string {
	return jenis + oldID + "." + mmyyyy + "." + lpad5(urut)
}

// PanjangNoAksepCLM - HitServiceToKasirKMT_Act S14.2 `@length(.AcceptedNo) = "21" || ... = "22"` (F -> 6 keluar).
func PanjangNoAksepCLM(no string) bool { return len(no) == 21 || len(no) == 22 }
