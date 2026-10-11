package models

// Untuk apa berkas ini: ANGKA UANG - `SetKomiteList_Act` (pra-proses flow action ViewTransferDtl) dan nomor akseptasi
// (`KomitePostAdjustment` S16.8). Uang lewat `apd.Decimal`, nol float, nol pembulatan di tengah hitungan (spec bab 7
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

// TotalMataUang - satu baris grid "Total Adjustment" (`TempTotalAdj.pxResults`).
type TotalMataUang struct {
	Currency string `json:"currency"`
	// AdjustmentGross - "Total Gross (100%)"; AdjustmentValue - "Total Adjustment RNM".
	AdjustmentGross string `json:"adjustmentGross"`
	AdjustmentValue string `json:"adjustmentValue"`
}

// LabelTotalIDR - baris penutup `SetKomiteList_Act` S8 (VERBATIM).
const LabelTotalIDR = "Total in IDR"

// TotalPenyesuaian = `SetKomiteList_Act` (TT 2; selainnya S1 keluar):
//
//	S4-S5   mata uang setiap baris `ClaimData.AdjustmentList` (salinan TempClaimData)
//	S6      Java: mata uang kembar dibuang, dipindai dari BELAKANG - yang tersisa kemunculan TERAKHIR, urutan asal
//	S7      per mata uang: jumlah GrossAdjustment / AdjustmentValue baris ber-`.AcceptanceStatus != 2` bermata uang itu
//	        (S7.1.1), sekaligus jumlah (nilai x KursIDR) semua baris itu dalam IDR
//	S8      baris "Total in IDR"
func TotalPenyesuaian(rows []map[string]string) ([]TotalMataUang, error) {
	kunci := make([]string, 0, len(rows))
	for _, b := range rows {
		kunci = append(kunci, b["Currency"])
	}
	sudah := map[string]bool{}
	var unik []string
	for i := len(kunci) - 1; i >= 0; i-- { // S6: HashSet dari belakang
		if sudah[kunci[i]] {
			continue
		}
		sudah[kunci[i]] = true
		unik = append([]string{kunci[i]}, unik...)
	}
	grossIDR, valueIDR := apd.New(0, 0), apd.New(0, 0)
	var out []TotalMataUang
	for _, mu := range unik {
		gross, value := apd.New(0, 0), apd.New(0, 0)
		for _, b := range rows {
			if b["AcceptanceStatus"] == KeputusanTolak || b["Currency"] != mu { // S7.1.1
				continue
			}
			g, err := desimal("GrossAdjustment", b["GrossAdjustment"])
			if err != nil {
				return nil, err
			}
			v, err := desimal("AdjustmentValue", b["AdjustmentValue"])
			if err != nil {
				return nil, err
			}
			kurs, err := desimal("KursIDR", b["KursIDR"])
			if err != nil {
				return nil, err
			}
			if _, err := konteksUang.Add(gross, gross, g); err != nil {
				return nil, err
			}
			if _, err := konteksUang.Add(value, value, v); err != nil {
				return nil, err
			}
			gi, vi := new(apd.Decimal), new(apd.Decimal)
			if _, err := konteksUang.Mul(gi, g, kurs); err != nil {
				return nil, err
			}
			if _, err := konteksUang.Mul(vi, v, kurs); err != nil {
				return nil, err
			}
			if _, err := konteksUang.Add(grossIDR, grossIDR, gi); err != nil {
				return nil, err
			}
			if _, err := konteksUang.Add(valueIDR, valueIDR, vi); err != nil {
				return nil, err
			}
		}
		out = append(out, TotalMataUang{Currency: mu, AdjustmentGross: TeksAngka(gross), AdjustmentValue: TeksAngka(value)})
	}
	out = append(out, TotalMataUang{Currency: LabelTotalIDR, AdjustmentGross: TeksAngka(grossIDR),
		AdjustmentValue: TeksAngka(valueIDR)})
	return out, nil
}

// lpad5 = LPAD(v_seq, 5, '0') (`PROC_GENERATE_SEQUENCE_NUMBER` keluaran `p_seq_number`).
func lpad5(n int) string { return fmt.Sprintf("%05d", n) }

// HurufAkseptasi - S16.6 `ParamSeq.CARI2 := ParamSeq.HASIL3 + "A"`.
const HurufAkseptasi = "A"

// RakitNomorAkseptasi = S16.8 `CARI2 + BusinessOldId + "." + HASIL1 + ".TP" + HASIL2`: jenis (kode produksi NONLIFE +
// "A"), kode lama bisnis, periode MM.YYYY, urut 5 digit.
func RakitNomorAkseptasi(jenis, oldID, mmyyyy string, urut int) string {
	return jenis + oldID + "." + mmyyyy + ".TP" + lpad5(urut)
}
