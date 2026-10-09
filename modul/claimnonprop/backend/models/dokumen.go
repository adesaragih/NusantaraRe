package models

// Untuk apa berkas ini: DOKUMEN dan POPUP DATA Non Prop - Print CFS (`GenerateCFS_act`), Print PLA (`GeneratePlaCNP_Act`),
// Generate Claim Analysis (`GenerateCACNP_Act`), popup Hitung_Test (`GetSelisihActual_Act`), dan Claim History Master ID
// (`GetHistoryMasterID_NP`).
//
// Stream HTML dokumen (CFSClaimXOL, PLACNP_Html, ClaimAnalysisHTML) TIDAK diekspor: nomor dan data dibangun penuh, isi
// berkas PDF tidak dikarang (OQ-CNP-22, prompt §6 butir 8). Tombol tetap menandai IsCFS / menerbitkan nomor PLA seperti
// XML supaya Submit dapat dicapai.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// PesanEstimasiKosong - GenerateCFS_act langkah 1 (VERBATIM).
const PesanEstimasiKosong = "Estimation Is Null!"

// PesanBerkasOQ - pemberitahuan tombol dokumen: berkas PDF menunggu OQ-CNP-22.
const PesanBerkasOQ = "OQ-CNP-22: stream HTML dokumen tidak diekspor - nomor dan data tersimpan, berkas PDF belum dibuat"

// DataCFS - data laporan CFS (`TempDataOutStanding` GenerateCFS_act): Claim Amount per mata uang dan Loss Allocation per
// treaty.
type DataCFS struct {
	KlaimMataUang  []Baris
	LossAllocation []Baris
}

// SusunDataCFS = GenerateCFS_act langkah 3-7.
//
// PERBAIKAN OQ-CNP-05 butir 1 (PARITAS `[penyimpangan sadar]`):
//   - langkah 4 (pembuang daftar salinan `TempDataOutStanding`) ber-REMARK, sehingga agregasi langkah 6-7 menjumlah ke
//     daftar yang SUDAH berisi salinan pyWorkPage (jumlah ganda). Di sini agregasi dimulai dari daftar kosong.
//   - langkah 7.2.1 menulis `.CNPOthersFee += Local.ClaimAmount` dan `.ClaimAmountAdjust += Local.OthersFee` (Fee <->
//     Claim Amount tertukar). Di sini Fee ke CNPOthersFee, Claim Amount ke ClaimAmountAdjust.
func SusunDataCFS(h *Halaman) (DataCFS, error) {
	var kal Kalkulator
	var out DataCFS
	for _, r := range h.AmbilDaftar(DaftarClaimAmount) { // 6
		klaim := kal.Tambah(kal.Kurang(kal.B(r, "ClaimAmountCedant"), kal.B(r, "TPL")), kal.B(r, "CNPDeductible"))
		var t Baris
		for _, o := range out.KlaimMataUang {
			if o["Currency"] == r["Currency"] {
				t = o
			}
		}
		if t == nil {
			out.KlaimMataUang = append(out.KlaimMataUang, Baris{"Currency": r["Currency"], "AdjusterFee": r["AdjusterFee"],
				"CNPOthersFee": r["CNPOthersFee"], "Salvage": r["Salvage"], "ClaimAmountCedant": Teks(klaim),
				"PctProrateClaim": r["PctProrateClaim"], "CNPDeductible": r["CNPDeductible"], "TPL": r["TPL"],
				"Note": "( ROE 1 " + r["Currency"] + " = IDR  " + r["AltValue"] + ")"})
			continue
		}
		t["ClaimAmountCedant"] = Teks(kal.Tambah(kal.B(t, "ClaimAmountCedant"), klaim))
		for _, p := range []string{"AdjusterFee", "CNPOthersFee", "Salvage", "CNPDeductible", "TPL"} {
			t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(r, p)))
		}
	}
	for _, l := range h.AmbilDaftar(DaftarLossAlloc) { // 7
		var t Baris
		for _, o := range out.LossAllocation {
			if o["TreatyName"] == l["TreatyName"] {
				t = o
			}
		}
		if t == nil {
			out.LossAllocation = append(out.LossAllocation, Baris{"TreatyName": l["TreatyName"], "Currency": l["Currency"],
				"AdjusterFee": l["AdjusterFee"], "CNPOthersFee": l["CNPOthersFee"], "Salvage": l["Salvage"],
				"ClaimAmountAdjust": l["ClaimAmountAdjust"]})
			continue
		}
		for _, p := range []string{"AdjusterFee", "CNPOthersFee", "Salvage", "ClaimAmountAdjust"} {
			t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(l, p)))
		}
	}
	return out, kal.Galat()
}

// CetakCFS = GenerateCFS_act langkah 1-2 dan 5: tanpa ListClaimAmount -> pesan (keluar); selain itu Insured Name huruf
// besar dan `IsCFS = 1`. Berkas PDF (langkah 11-16) = OQ-CNP-22.
func CetakCFS(h *Halaman) (DataCFS, bool, error) {
	if len(h.AmbilDaftar(DaftarClaimAmount)) == 0 {
		h.TambahPesan("", PesanEstimasiKosong)
		return DataCFS{}, false, nil
	}
	h.Setel(CD+"InsuredName", strings.ToUpper(h.Ambil(CD+"InsuredName")))
	h.Setel("IsCFS", "1")
	d, err := SusunDataCFS(h)
	return d, err == nil, err
}

// RakitNomorPLA = RDB `GenerateNoPLATNP`: 'RNM-M' || BusinessOldId || '.' || MM || '.' || yyyy || '.TX' ||
// lpad(PLATNP_SEQ.nextval, 5, '0') (tanggal sistem).
func RakitNomorPLA(oldID string, saat time.Time, urut int) string {
	t := saat.In(Jakarta)
	return fmt.Sprintf("RNM-M%s.%02d.%d.TX%05d", oldID, int(t.Month()), t.Year(), urut)
}

// RevisiNomorPLA = GeneratePlaCNP_Act langkah 7: nomor PLA sudah ada dan nilai klaim berubah (FlagPrintPla.CARI28 = 1)
// -> `NoPla + "/" + (substring(NoPla, 24) + 1)`. Pada nomor pertama `substring(…,24)` kosong (= 0) sehingga revisi
// pertama "/1"; pada revisi berikut `substring` memuat "/n" lama - Pega `@toDecimal` atas teks itu = 0 -> "/1" lagi
// (ditiru: hanya digit setelah posisi 24 dibaca).
func RevisiNomorPLA(no string) string {
	r := []rune(no)
	sisa := ""
	if len(r) > 24 {
		sisa = string(r[24:])
	}
	n, err := strconv.Atoi(strings.TrimSpace(sisa))
	if err != nil {
		n = 0
	}
	return no + "/" + strconv.Itoa(n+1)
}

// CetakCA = GenerateCACNP_Act langkah awal: Insured Name huruf besar (berkas PDF = OQ-CNP-22).
func CetakCA(h *Halaman) {
	h.Setel(CD+"InsuredName", strings.ToUpper(h.Ambil(CD+"InsuredName")))
}

// SelisihAktual - isi popup Hitung_Test.
type SelisihAktual struct {
	EstimasiTotal  string  `json:"estimasiTotal"`  // TempEstimation.CARI2
	EstimasiSpread string  `json:"estimasiSpread"` // TempEstimation.CARI1
	Akseptasi      []Baris `json:"akseptasi"`      // TempAksep.pxResults (CARI2 Total Claim, CARI1 Claim Spreaded)
	SelisihTotal   string  `json:"selisihTotal"`   // InputParamOs.GrossValue
	SelisihSpread  string  `json:"selisihSpread"`  // InputParamOs.Value
	AktualTotal    string  `json:"aktualTotal"`    // TempActual.CARI2
	AktualSpread   string  `json:"aktualSpread"`   // TempActual.CARI1
}

// HitungSelisihAktual = `GetSelisihActual_Act` (tombol View Waiting For Actual Premium). `os` = Σ baris OS STS 0 per
// layer x mata uang (urut sama dengan `KelompokLayerOS`). Nilai layer TERAKHIR yang tampil (popup menimpa per iterasi).
//
// PERBAIKAN OQ-CNP-05 butir 2 (PARITAS `[penyimpangan sadar]`): langkah 6.3 membaca `OutOSAcc.pxResults(1).CARI1 /
// CARI2`, padahal SQL `GetDataOS` mengembalikan alias Value / GrossValue / Adjusterfee / CNPOthersFee / Salvage - nilai
// estimasi selalu kosong. Di sini Value dan GrossValue yang dibaca.
func HitungSelisihAktual(h *Halaman, os []NilaiOS) (SelisihAktual, error) {
	var kal Kalkulator
	var out SelisihAktual
	aktual := h.Ambil("FlagActualPremium") == "true"
	for i, l := range KelompokLayerOS(h) {
		var lama NilaiOS
		if i < len(os) {
			lama = os[i]
		}
		value, gross := kal.Teks("Value", lama.Value), kal.Teks("GrossValue", lama.GrossValue)
		out.EstimasiSpread, out.EstimasiTotal = Teks(value), Teks(gross)
		if aktual { // 6.4
			for j := range h.AmbilDaftar(DaftarAdjustment) {
				for _, s := range h.AmbilDaftar(JalurAdj(j+1, AnakXOL)) {
					if s["TreatyName"] == TreatyUR {
						continue
					}
					value = kal.Kurang(value, kal.B(s, "ClaimSpreaded"))
					gross = kal.Kurang(gross, kal.B(s, "TotalClaim"))
					out.Akseptasi = append(out.Akseptasi, Baris{"CARI1": s["ClaimSpreaded"], "CARI2": s["TotalClaim"]})
				}
			}
		}
		out.AktualSpread, out.AktualTotal = l["ClaimSpreaded"], l["TotalClaim"] // 6.5
		out.SelisihSpread = Teks(kal.Kurang(kal.B(l, "ClaimSpreaded"), value))  // 6.6
		out.SelisihTotal = Teks(kal.Kurang(kal.B(l, "TotalClaim"), gross))
	}
	return out, kal.Galat()
}

// IDMasterRiwayat = GetHistoryMasterID_NP langkah 2: 7 karakter pertama IDMaster, + "/R01", + "/R02" (hardcode master
// 1001130 langkah 4-5 dibuang, OQ-CNP-03); judul " CLAIM HISTORY MASTER ID " + 7 karakter.
func IDMasterRiwayat(idMaster string) ([3]string, string) {
	a := AwalanMaster(idMaster)
	return [3]string{a, a + "/R01", a + "/R02"}, " CLAIM HISTORY MASTER ID " + a
}

// BarisXOL2 - satu baris CLAIMXOL2 (`GetHistoryMasterID`).
type BarisXOL2 struct {
	CaseID           string `json:"caseId"`
	XOL              string `json:"xol"`
	Currency         string `json:"currency"`
	GrossAdjustment  string `json:"grossAdjustment"`
	CNPReinstatement string `json:"cnpReinstatement"`
	KursIDR          string `json:"kursIdr"`
}

// TotalRiwayatMaster = GetHistoryMasterID_NP langkah 7-9: total per layer x mata uang (Java langkah 8 membuang ganda
// dari belakang - kemunculan PERTAMA dipertahankan).
func TotalRiwayatMaster(rows []BarisXOL2) ([]Baris, error) {
	var kal Kalkulator
	var out []Baris
	for _, r := range rows {
		var t Baris
		for _, o := range out {
			if o["XOL"] == r.XOL && o["Currency"] == r.Currency {
				t = o
			}
		}
		if t == nil {
			t = Baris{"XOL": r.XOL, "Currency": r.Currency, "GrossAdjustment": "0", "CNPReinstatement": "0"}
			out = append(out, t)
		}
		t["GrossAdjustment"] = Teks(kal.Tambah(kal.B(t, "GrossAdjustment"), kal.Teks("g", r.GrossAdjustment)))
		t["CNPReinstatement"] = Teks(kal.Tambah(kal.B(t, "CNPReinstatement"), kal.Teks("r", r.CNPReinstatement)))
	}
	return out, kal.Galat()
}

// nolDesimal - 0.
func nolDesimal() *apd.Decimal { return apd.New(0, 0) }
