package models

// Untuk apa berkas ini: MUATAN KASIR - `HitServiceToKasirKMT_Act` jalur CLM (S2, S7, S12-S14.2; dipanggil
// KomitePost_Adjustment S25.2.1.1 bila `.DirectToKasir == "true" && .StatusKasir = ""`). Bentuk muatan disalin dari
// Claim Fac In `models.SusunMuatanKasir` (HitServiceToKasir_Act 13.2, rule REST yang sama `SendAcceptationToKasir`),
// bukan impor. Kode tetap Kasir dari `konfigurasi/kasir.json` (pola komite lain), bukan literal kode.
//
// PERBAIKAN prompt §5 butir 2: S25.2.1.1 `.StatusKasir = ""` (tanda `=` tunggal) dibaca pembandingan. S7
// `.NoAccount := angka saja` mengubah baris adjustment - di sini hanya muatan (NoAccount di luar daftar putih kontrak).
// S3 `getStatusKonversi_Act` menggerbang lewat transisi PASCA-langkah `.StatusKonversi=="1"` (T=2 F=6) - ditegakkan
// layanan (`services.kasir`, `Acuan.StatusKonversi`); RALAT 10-10-2026 atas OQ-KCFI-04 (alat dump hanya mencetak prasyarat).

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// KonfigurasiKasir - kode tetap muatan Kasir (S14.2.2.1.1.3-S14.2.2.1.1.4).
type KonfigurasiKasir struct {
	CompanyName, LjtdID, LdcID, LdcIDSyariah string
}

// MuatanKasir - satu muatan `SendAcceptationToKasir` (bentuk sama dengan Claim Fac In).
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuId, NoPolis, AcceptType, Kepada, AccountNo, TglAksep string
	Nett, Deductible, KaliDeduct                                              string
	StsSyariah, CompanyName, LjtdId, LdcId, StsAp, LkuId, LbgID               string
	TglBolehBayar, Email, UserInput                                           string
}

var bukanAngka = regexp.MustCompile(`[^0-9]`)

// BolehKasir - S2 `.DirectToKasir == "true" && .StatusKasir = ""`.
func BolehKasir(b map[string]string) bool {
	return b["DirectToKasir"] == "true" && b["StatusKasir"] == ""
}

// KunciCedingKasir - S14.2.2.1.1.1 `@substring(pyWorkCover.Quotation.CedingCo, 0, 8)` (`Quotation` = salinan
// `OfferFacIn.QuotationData`, pola Claim Fac In).
func KunciCedingKasir(kl kontrak.KlaimFacIn) string {
	c := kl.Nilai[OQ+"CedingCo"]
	if len(c) > 8 {
		return c[:8]
	}
	return c
}

// ymd - tanggal teks halaman ("2006-01-02 ..." atau "yyyyMMdd...") -> "yyyyMMdd".
func ymd(s string) string {
	d := bukanAngka.ReplaceAllString(s, "")
	if len(d) < 8 {
		return ""
	}
	return d[:8]
}

// tglAksepTeks = `@substring(d,6,8)+"-"+@substring(d,4,6)+"-"+@substring(d,0,4)` atas "yyyyMMdd".
func tglAksepTeks(t string) string {
	if len(t) < 8 {
		return ""
	}
	return t[6:8] + "-" + t[4:6] + "-" + t[0:4]
}

// TglBolehBayar = S14.2.2.1.1.5: hari > 25 -> tanggal 01 dua bulan sesudahnya, selainnya hari yang sama bulan berikutnya;
// bulan 13 -> "1"; tahun naik bila bulan hasil 01 DAN bulan berjalan Desember.
func TglBolehBayar(t string, bulanKini int) string {
	if len(t) < 8 {
		return ""
	}
	hari, _ := strconv.Atoi(t[6:8])
	bulan, _ := strconv.Atoi(t[4:6])
	tahun, _ := strconv.Atoi(t[0:4])
	bulan++
	if hari > 25 {
		bulan++
	}
	if bulan >= 13 {
		bulan = 1
	}
	hh := t[6:8]
	if hari > 25 {
		hh = "01"
	}
	if bulan == 1 && bulanKini == 12 {
		tahun++
	}
	return hh + "-" + strconv.Itoa(100 + bulan)[1:] + "-" + strconv.Itoa(tahun)
}

// SusunMuatanKasir = S7 + S14.2.2 (Local.IndexObject = objek kasus komite): seluruh adjustment item objek ber-AcceptedNo
// sama dijumlah - Nett menurut Payment Type (1/2/5 AdjustmentValue, 4/6 AdjusterFeeValue, 3 SalvageValue), Deductible
// Σ IndividualRiskRNM, KaliDeduct Σ IndividualRiskPercentage; medan lain dari baris cocok TERAKHIR. `adj` = baris adjustment
// sesudah keputusan (AcceptedNo, IDOfBank S12-S13). `pengirim` = akun pemutus (`pxRequestor.pxUserIdentifier`).
func SusunMuatanKasir(kl kontrak.KlaimFacIn, adj map[string]string, email, pengirim string, c KonfigurasiKasir,
	bulanKini int) (MuatanKasir, error) {
	var h hitung
	nett, ded, kali := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	akhir := adj
	for in := range kl.Daftar[DaftarItem(kl.Objek)] { // S14.2.2
		for an, x := range kl.Daftar[DaftarAdj(kl.Objek, in+1)] {
			if in+1 == kl.Item && an+1 == kl.Adjustment {
				x = adj // baris kasus komite sesudah keputusan
			}
			if x["AcceptedNo"] != adj["AcceptedNo"] { // S14.2.2.1.1
				continue
			}
			akhir = x
			ded = h.tambah(ded, h.dari("IndividualRiskRNM", x["IndividualRiskRNM"]))
			kali = h.tambah(kali, h.dari("IndividualRiskPercentage", x["IndividualRiskPercentage"]))
			switch x["PaymentType"] { // S14.2.2.1.1.6-8
			case "1", "2", "5":
				nett = h.tambah(nett, h.dari("AdjustmentValue", x["AdjustmentValue"]))
			case "4", "6":
				nett = h.tambah(nett, h.dari("AdjusterFeeValue", x["AdjusterFeeValue"]))
			case "3":
				nett = h.tambah(nett, h.dari("SalvageValue", x["SalvageValue"]))
			}
		}
	}
	if h.err != nil {
		return MuatanKasir{}, h.err
	}
	tgl := ymd(akhir["AcceptedDate"])
	user := akhir["pxCreateOperator"]
	if user == "" {
		user = pengirim
	}
	return MuatanKasir{
		NoTrans: akhir["AcceptedNo"], NoKlaim: kl.Nilai["ClaimData.NoClaim"], LbuId: kl.Nilai[OQ+"BusinessOldId"],
		NoPolis: kl.Nilai["OfferFacIn.PolicyData.PolicyNo"], AcceptType: akhir["PaymentType"], Kepada: akhir["PayableTo"],
		AccountNo: bukanAngka.ReplaceAllString(akhir["NoAccount"], ""), TglAksep: tglAksepTeks(tgl), Nett: TeksAngka(nett),
		Deductible: TeksAngka(ded), KaliDeduct: TeksAngka(kali), StsSyariah: "0", CompanyName: c.CompanyName,
		LjtdId: c.LjtdID, LdcId: c.LdcID, StsAp: "0", LkuId: akhir["CurrencyID"], LbgID: akhir["IDOfBank"],
		TglBolehBayar: TglBolehBayar(tgl, bulanKini), Email: strings.TrimSpace(email), UserInput: user,
	}, nil
}
