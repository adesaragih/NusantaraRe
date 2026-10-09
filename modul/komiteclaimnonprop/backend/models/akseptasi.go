package models

// Untuk apa berkas ini: LANGKAH TINGKAT AKHIR yang membentuk data (korpus `Komite Claim Non Prop`) - baris OS akseptasi
// (`InsertOSKlaimCNP` -> `SaveOSClaim_SQL` / `PEGA_JSON_OS_AKSEP_KLAIM`), baris CLAIMXOL2 (`InsertXOLKlaimCNP` ->
// `XOL2_AKSEP_KLAIM`), baris OS subjectivity (`InsertOSSubjectivityCNP` -> `PEGA_JSON_OS_AKSEP_SUBJECTIVITY`), muatan
// Kasir (`HitServiceToKasirKMT_Act` cabang IsCLMNP), dan pembantu email (`SendEmailKlaim_KMT` cabang IsCLMNP). Murni;
// layanan yang membaca acuan dan menulis. Procedure tidak dipanggil (isinya dibaca dari ALL_SOURCE DEV 09-10-2026).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// KunciInstans - padanan `pzInsKey` kasus sistem baru di tabel warisan (OS_AKSEPTASI_KLAIM.CASEID, JSON_KLAIM.IDPEGA,
// HISTORYAKSEPTASIPEGA.ID_PEGA / ID_KOMITE, CLAIMXOL2.CASEID, OS_AKSEPTASI_SUBJECTIVITY.CASEID): ID T_WORK_CLAIM apa
// adanya - sama dengan Claim Non Prop `models.KunciInstans`.
func KunciInstans(id string) string { return id }

// DaftarAdjustment - jalur daftar baris akseptasi halaman klaim induk (`ClaimData.AdjustmentList`).
const DaftarAdjustment = "ClaimData.AdjustmentList"

// jalurAdj - jalur properti / anak baris akseptasi ke-n halaman klaim induk.
func jalurAdj(n int, p string) string { return fmt.Sprintf("ClaimData.AdjustmentList(%d).%s", n, p) }

// AdjustmentKlaim - baris akseptasi klaim induk yang sedang diputus (kosong bila tidak ada).
func AdjustmentKlaim(kl kontrak.KlaimTreaty) map[string]string {
	rows := kl.Daftar[DaftarAdjustment]
	if kl.Adjustment < 1 || kl.Adjustment > len(rows) {
		return map[string]string{}
	}
	return rows[kl.Adjustment-1]
}

// anakAdj - daftar anak baris akseptasi yang diputus (`pyWorkPage.Adjustment.<anak>` = salinan baris akseptasi).
func anakAdj(kl kontrak.KlaimTreaty, anak string) []map[string]string {
	return kl.Daftar[jalurAdj(kl.Adjustment, anak)]
}

// Anak baris akseptasi (katalog Claim Non Prop).
const (
	AnakXOL        = "SpreadingRisk"
	AnakXOLDibayar = "AlokasiXOLPaid"
	AnakSpreadIn   = "SpreadingAdjustment"
	AnakSpreadOut  = "SpreadingQuotaShare"
	AnakClaimAcc   = "ListClaimAcceptation"
	AnakLossAlloc  = "CNPSpreadLoss"
	// TreatyUR - layer retensi (`.TreatyName == "UR"`), tidak masuk CNPLayerList / CLAIMXOL2.
	TreatyUR = "UR"
)

// Kelas halaman JSON (DATA_JSON OS DEV, pola Claim Non Prop).
const (
	KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"
	KelasAdjustment  = "ASM-FW-GCNMFW-Data-Adjustment"
)

// Kode tetap efek tingkat akhir.
const (
	// StsRejectKonversi - KomitePostAdjustment S14.22 `KonversiKlaim_Act` parameter STSREJECT (tertulis mati).
	StsRejectKonversi = "1"
	// StatusKonversiSudah - HitServiceToKasirKMT_Act S3 transisi `.StatusKonversi` sudah dikonversi.
	StatusKonversiSudah = "1"
)

// StsOSAkseptasi = InsertOSKlaimCNP S3 `InputData.CARI10 := @if(Adjustment.PaymentType == "1", "4", "1")`.
func StsOSAkseptasi(paymentType string) string {
	if paymentType == "1" {
		return "4"
	}
	return "1"
}

// HalamanJSON - satu halaman untuk `@ASM.GetPageJSONString()`: nilai teks + daftar halaman (PageList).
type HalamanJSON struct {
	Nilai  map[string]string
	Daftar []DaftarJSON
}

// DaftarJSON - satu PageList bernama.
type DaftarJSON struct {
	Nama string
	Isi  []HalamanJSON
}

// JSONHalamanPega = `@ASM.GetPageJSONString()` (disalin dari Claim Non Prop `JSONHalamanPega`): "{" LF, pasangan
// dipisah LF ",", ditutup LF "}", nilai teks urut kunci tanpa membedakan huruf, properti kosong tidak ditulis, PageList
// sesudahnya; halaman puncak diakhiri LF.
func JSONHalamanPega(p HalamanJSON) string { return tulisHalamanJSON(p) + "\n" }

func tulisHalamanJSON(p HalamanJSON) string {
	kunci := make([]string, 0, len(p.Nilai))
	for k, v := range p.Nilai {
		if v != "" {
			kunci = append(kunci, k)
		}
	}
	sort.Slice(kunci, func(i, j int) bool {
		a, b := strings.ToLower(kunci[i]), strings.ToLower(kunci[j])
		if a != b {
			return a < b
		}
		return kunci[i] < kunci[j]
	})
	var bagian []string
	for _, k := range kunci {
		bagian = append(bagian, teksJSON(k)+":"+teksJSON(p.Nilai[k]))
	}
	for _, d := range p.Daftar {
		if len(d.Isi) == 0 {
			continue
		}
		var isi []string
		for _, a := range d.Isi {
			isi = append(isi, tulisHalamanJSON(a))
		}
		bagian = append(bagian, teksJSON(d.Nama)+":[ \n"+strings.Join(isi, "\n,")+"\n] ")
	}
	return "{\n" + strings.Join(bagian, "\n,") + "\n}"
}

func teksJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(buf.String(), "\n")
}

func k100() *apd.Decimal { return apd.New(100, 0) }

// LayerListAkseptasi = CreateChildKomiteCNP_Act 20-23 (`ChildWorkPage.Adjustment.CNPLayerList`, disalin dari Claim Non
// Prop `LayerListAkseptasi`): satu layer per nama XoL non-UR, satu `CNPCurrencyList` per baris mata uang XOL Allocation
// akseptasi; Previously Calculated dikurangkan.
func LayerListAkseptasi(kl kontrak.KlaimTreaty) ([]HalamanJSON, error) {
	var kal Kalkulator
	rnm := kal.BagiPega(kal.Teks("TreatyInMaster.RNMShare", kl.Nilai["TreatyInMaster.RNMShare"]), k100())
	type layer struct {
		xol, xolID string
		mu         []map[string]string
	}
	var daftar []*layer
	for _, s := range anakAdj(kl, AnakXOL) {
		if s["TreatyName"] == TreatyUR {
			continue
		}
		pct := kal.BagiPega(kal.B(s, "ClaimPercentage"), k100())
		var l *layer
		for _, d := range daftar {
			if d.xol == s["TreatyName"] {
				l = d
			}
		}
		mu := map[string]string{"pxObjClass": KelasAdjustment, "AdjusterFeeRNM": s["AdjusterFee"],
			"CNPLimit": s["CNPLimit"], "CNPMindep": s["CNPMDP"], "CNPPctReinstate": s["CNPPctReinstate"],
			"CNPReinstatement": s["CNPReinstatement"], "CNPReinstatementRNM": s["CNPReinstatementRNM"],
			"Currency": s["Currency"], "CurrencyID": s["CurrencyID"], "GrossAdjustment": s["ClaimEstimation"],
			"GrossValue": s["ClaimSpreaded"], "KursIDR": s["Kurs"], "PersenRNM": s["ClaimPercentage"],
			"SalvageRNM": s["Salvage"], "TotalXOL": s["TotalClaim"], "TotalXOLGross": s["TotalClaim"],
			"CNPOthersFeeRNM": s["CNPOthersFee"]}
		total := kal.Kali(kal.B(s, "TotalClaim"), pct)
		mu["TotalXOLLayerRNM"], mu["TotalXOLRNM"] = Teks(total), Teks(total)
		if l == nil { // layer baru - pembagi RNM Share master
			mu["AdjusterFeeValue"] = Teks(kal.BagiPega(kal.B(s, "AdjusterFee"), rnm))
			mu["SalvageValue"] = Teks(kal.BagiPega(kal.B(s, "Salvage"), rnm))
			mu["CNPOthersFee"] = Teks(kal.BagiPega(kal.B(s, "CNPOthersFee"), rnm))
			daftar = append(daftar, &layer{xol: s["TreatyName"], xolID: s["TreatyType"], mu: []map[string]string{mu}})
			continue
		}
		mu["AdjusterFeeValue"] = Teks(kal.BagiPega(kal.B(s, "AdjusterFee"), pct))
		mu["SalvageValue"] = Teks(kal.BagiPega(kal.B(s, "Salvage"), pct))
		mu["CNPOthersFee"] = Teks(kal.BagiPega(kal.B(s, "CNPOthersFee"), pct))
		l.mu = append(l.mu, mu)
	}
	if dibayar := anakAdj(kl, AnakXOLDibayar); len(dibayar) > 0 {
		for _, l := range daftar {
			for _, mu := range l.mu {
				var p map[string]string
				for _, d := range dibayar {
					if d["Currency"] == mu["Currency"] && d["TreatyName"] == l.xol {
						p = d
					}
				}
				if p == nil {
					continue
				}
				kurang := func(k string, v *apd.Decimal) { mu[k] = Teks(kal.Kurang(kal.Teks(k, mu[k]), v)) }
				kurang("AdjusterFeeRNM", kal.B(p, "AdjusterFee"))
				kurang("AdjusterFeeValue", kal.BagiPega(kal.B(p, "AdjusterFee"), rnm))
				kurang("CNPOthersFee", kal.BagiPega(kal.B(p, "CNPOthersFee"), rnm))
				kurang("CNPOthersFeeRNM", kal.B(p, "CNPOthersFee"))
				kurang("CNPReinstatement", kal.B(p, "CNPReinstatement"))
				kurang("CNPReinstatementRNM", kal.B(p, "CNPReinstatementRNM"))
				kurang("GrossAdjustment", kal.B(p, "ClaimAmountAdjust"))
				kurang("GrossValue", kal.B(p, "ClaimSpreaded"))
				kurang("SalvageRNM", kal.B(p, "Salvage"))
				kurang("SalvageValue", kal.BagiPega(kal.B(p, "Salvage"), rnm))
				g := kal.Tambah(kal.Kurang(kal.Tambah(kal.Teks("g", mu["GrossAdjustment"]), kal.Teks("a", mu["AdjusterFeeValue"])),
					kal.Teks("s", mu["SalvageValue"])), kal.Teks("o", mu["CNPOthersFee"]))
				mu["TotalXOL"], mu["TotalXOLGross"] = Teks(g), Teks(g)
				r := kal.Kurang(kal.Tambah(kal.Tambah(kal.Teks("g", mu["GrossValue"]), kal.Teks("a", mu["AdjusterFeeRNM"])),
					kal.Teks("o", mu["CNPOthersFeeRNM"])), kal.Teks("s", mu["SalvageRNM"]))
				mu["TotalXOLLayerRNM"], mu["TotalXOLRNM"] = Teks(r), Teks(r)
			}
		}
	}
	var out []HalamanJSON
	for _, l := range daftar {
		var mu []HalamanJSON
		for _, m := range l.mu {
			mu = append(mu, HalamanJSON{Nilai: m})
		}
		out = append(out, HalamanJSON{Nilai: map[string]string{"pxObjClass": KelasAdjustment, "XOL": l.xol, "XOLID": l.xolID},
			Daftar: []DaftarJSON{{Nama: "CNPCurrencyList", Isi: mu}}})
	}
	return out, kal.Galat()
}

// BarisOSAkseptasi - satu baris OS_AKSEPTASI_KLAIM (`PEGA_JSON_OS_AKSEP_KLAIM` baris 16: CASEID, NOCLAIM, DATA_JSON,
// TANGGAL hari ini, NOPOLIS, STS_REJECT, STS_KONVERSI, STS_DLA, MASTERID, CLAIMOLD). STS_KONVERSI / STS_DLA (CARI16 /
// CARI17) tidak diisi InsertOSKlaimCNP (NULL); CLAIMOLD (CARI20) hanya untuk satu akun yang tertulis mati (OQ-CNP-03,
// dibuang).
type BarisOSAkseptasi struct {
	CaseID, NoClaim, NoPolis, StsReject, MasterID string
	DataJSON                                      string
}

// BarisOSSubjectivity - satu baris OS_AKSEPTASI_SUBJECTIVITY (`PEGA_JSON_OS_AKSEP_SUBJECTIVITY`: UPDATE STS_SUBJECTIVITY
// / DATA_JSON / TANGGAL bila CASEID ada, selainnya INSERT CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT,
// STS_KONVERSI, STS_DLA, STS_SUBJECTIVITY). STS_KONVERSI / STS_DLA tidak diisi (NULL).
type BarisOSSubjectivity struct {
	CaseID, NoClaim, NoPolis, StsReject, StsSubjectivity string
	DataJSON                                             string
}

// HalamanOSAkseptasi = InsertOSKlaimCNP S2 / InsertOSSubjectivityCNP S1: halaman TempOSAkseptasi / TempOSSubjectivity
// (CauseOfLoss, CauseOfLossID, EstimationDate `@getCurrentDateStamp()`, NoClaim, PersenRNM, AcceptedNo, Type 1,
// CNPLayerList salinan kasus komite).
func HalamanOSAkseptasi(kl kontrak.KlaimTreaty, acceptedNo string, saat time.Time) (string, error) {
	layer, err := LayerListAkseptasi(kl)
	if err != nil {
		return "", err
	}
	v := kl.Nilai
	p := HalamanJSON{Nilai: map[string]string{"pxObjClass": KelasOSAkseptasi, "CauseOfLoss": v["ClaimData.CauseOfLoss"],
		"CauseOfLossID": v["ClaimData.CauseOfLossID"], "EstimationDate": saat.In(Jakarta).Format("20060102"),
		"NoClaim": v["ClaimData.NoClaim"], "PersenRNM": v["TreatyInMaster.RNMShare"], "AcceptedNo": acceptedNo, "Type": "1"},
		Daftar: []DaftarJSON{{Nama: "CNPLayerList", Isi: layer}}}
	return JSONHalamanPega(p), nil
}

// SusunOSAkseptasi = InsertOSKlaimCNP S2-S6.
func SusunOSAkseptasi(kl kontrak.KlaimTreaty, klaimID, acceptedNo string, saat time.Time) (BarisOSAkseptasi, error) {
	dj, err := HalamanOSAkseptasi(kl, acceptedNo, saat)
	if err != nil {
		return BarisOSAkseptasi{}, err
	}
	v := kl.Nilai
	return BarisOSAkseptasi{CaseID: KunciInstans(klaimID), NoClaim: v["ClaimData.NoClaim"],
		NoPolis: v["ClaimData.PolicyData.PolicyNo"], StsReject: StsOSAkseptasi(AdjustmentKlaim(kl)["PaymentType"]),
		MasterID: v["ClaimData.IDMaster"], DataJSON: dj}, nil
}

// SusunOSSubjectivity = InsertOSSubjectivityCNP S1-S4 (STS_SUBJECTIVITY = 1 bila subjectivity, selainnya 0).
func SusunOSSubjectivity(kl kontrak.KlaimTreaty, klaimID, acceptedNo string, subjectivity bool, saat time.Time) (
	BarisOSSubjectivity, error) {
	dj, err := HalamanOSAkseptasi(kl, acceptedNo, saat)
	if err != nil {
		return BarisOSSubjectivity{}, err
	}
	sts := "0"
	if subjectivity {
		sts = "1"
	}
	v := kl.Nilai
	return BarisOSSubjectivity{CaseID: KunciInstans(klaimID), NoClaim: v["ClaimData.NoClaim"],
		NoPolis: v["ClaimData.PolicyData.PolicyNo"], StsReject: StsOSAkseptasi(AdjustmentKlaim(kl)["PaymentType"]),
		StsSubjectivity: sts, DataJSON: dj}, nil
}

// BarisXOL2 - satu baris CLAIMXOL2 (`XOL2_AKSEP_KLAIM`: CASEID, "GrossAdjustment", "CNPReinstatement", "Currency",
// "KursIDR", XOL, TANGGAL). Kolom teks (VARCHAR2) - nilai apa adanya.
type BarisXOL2 struct {
	CaseID, GrossAdjustment, CNPReinstatement, Currency, KursIDR, XOL string
}

// SusunXOL2 = InsertXOLKlaimCNP S3-S5: setiap baris XOL Allocation akseptasi non-UR (TotalClaim, CNPReinstatement,
// Currency, KursIDR, TreatyName). Subjectivity -> nol baris (S1).
func SusunXOL2(kl kontrak.KlaimTreaty, klaimID string, subjectivity bool) []BarisXOL2 {
	if subjectivity {
		return nil
	}
	var out []BarisXOL2
	for _, s := range anakAdj(kl, AnakXOL) {
		if s["TreatyName"] == TreatyUR {
			continue
		}
		out = append(out, BarisXOL2{CaseID: KunciInstans(klaimID), GrossAdjustment: s["TotalClaim"],
			CNPReinstatement: s["CNPReinstatement"], Currency: s["Currency"], KursIDR: s["KursIDR"], XOL: s["TreatyName"]})
	}
	return out
}

var bukanAngka = regexp.MustCompile(`[^0-9]`)

// YMD - "yyyyMMdd" (8 digit pertama) dari teks tanggal halaman ("2006-01-02[ 15:04:05]" atau "20060102...").
func YMD(s string) string {
	d := bukanAngka.ReplaceAllString(s, "")
	if len(d) > 8 {
		return d[:8]
	}
	return d
}

// KonfigurasiKasir - kode tetap muatan `SendAcceptationToKasir` (HitServiceToKasirKMT_Act 10.3: CompanyName, LjtdId,
// LdcId, StsAp; 10.4 IsPEGASyariah LdcId) dari konfigurasi berdokumen (MODUL.md), bukan literal kode.
type KonfigurasiKasir struct {
	CompanyName, LjtdID, LdcID, LdcIDSyariah string
}

// MuatanKasir - badan `TAllPaymentData` (HitServiceToKasirKMT_Act 10.6 Java).
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuId, NoPolis, AcceptType, Kepada, AccountNo, TglAksep string
	Nett, Deductible, KaliDeduct                                              string
	StsSyariah, CompanyName, LjtdId, LdcId, StsAp, LkuId, LbgID               string
	TglBolehBayar, Email, UserInput                                           string
}

// PanjangNoAksepCNP - langkah 9 / 10 `@length(Primary.AcceptedNo) = 23 || 24`; selainnya keluar.
func PanjangNoAksepCNP(no string) bool { n := len([]rune(no)); return n == 23 || n == 24 }

// tglAksepTeks = `@substring(d,6,8)+"-"+@substring(d,4,6)+"-"+@substring(d,0,4)` atas "yyyyMMdd".
func tglAksepTeks(ymd string) string {
	if len(ymd) < 8 {
		return ""
	}
	return ymd[6:8] + "-" + ymd[4:6] + "-" + ymd[0:4]
}

// TanggalBolehBayar = HitServiceToKasirKMT_Act 10.5: tanggal akseptasi + 1 bulan, +1 lagi bila tanggal > 25 (tanggal
// 01). PERBAIKAN yang sama dengan Claim Non Prop OQ-CNP-05 butir 5 (`[penyimpangan sadar]`, PARITAS): XML menaikkan
// tahun hanya bila bulan hasil 01 DAN bulan SEKARANG Desember; di sini tahun bergulir bersama bulan.
func TanggalBolehBayar(ymd string) string {
	if len(ymd) < 8 {
		return ""
	}
	t, err := time.Parse("20060102", ymd)
	if err != nil {
		return ""
	}
	hari, tambah := t.Day(), 1
	if hari > 25 {
		hari, tambah = 1, 2
	}
	awal := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, tambah, 0)
	return fmt.Sprintf("%02d-%02d-%d", hari, int(awal.Month()), awal.Year())
}

// SusunMuatanKasir = HitServiceToKasirKMT_Act 9-10 (cabang IsCLMNP): baris Spreading In akseptasi dijumlah per CurrencyID
// (TotalClaim, PremiumSpreaded; NoAccount / IDOfBank baris pertama mata uang itu), satu muatan per mata uang: Nett =
// TotalClaim - PremiumSpreaded, Deductible / KaliDeduct 0. `b` = baris akseptasi induk, `email` = email ceding (10.2),
// `syariah` = IsPEGASyariah (OQ-CNP-40, false).
func SusunMuatanKasir(kl kontrak.KlaimTreaty, b map[string]string, email string, syariah bool, cfg KonfigurasiKasir,
	pelaku string) ([]MuatanKasir, error) {
	ldc := cfg.LdcID
	if syariah {
		ldc = cfg.LdcIDSyariah
	}
	type jumlah struct {
		cur, akun, bank   string
		total, premiSebar *apd.Decimal
	}
	var kal Kalkulator
	var urut []*jumlah
	for _, s := range anakAdj(kl, AnakSpreadIn) {
		var j *jumlah
		for _, x := range urut {
			if x.cur == s["CurrencyID"] {
				j = x
			}
		}
		if j == nil {
			urut = append(urut, &jumlah{cur: s["CurrencyID"], akun: s["NoAccount"], bank: s["IDOfBank"],
				total: kal.B(s, "TotalClaim"), premiSebar: kal.B(s, "PremiumSpreaded")})
			continue
		}
		j.total = kal.Tambah(j.total, kal.B(s, "TotalClaim"))
		j.premiSebar = kal.Tambah(j.premiSebar, kal.B(s, "PremiumSpreaded"))
	}
	tgl := YMD(b["AcceptedDate"])
	user := b["pxCreateOperator"]
	if user == "" {
		user = pelaku // `pxRequestor.pxUserIdentifier`
	}
	v := kl.Nilai
	var out []MuatanKasir
	for _, j := range urut {
		akun := RekeningAngka(j.akun)
		if akun == "" { // 10.3 `@replaceAll(Primary.NoAccount,"-","")` atas NoAccount yang sudah angka saja (S7)
			akun = RekeningAngka(b["NoAccount"])
		}
		out = append(out, MuatanKasir{NoTrans: b["AcceptedNo"], NoKlaim: v["ClaimData.NoClaim"],
			LbuId: v["OfferFacIn.QuotationData.BusinessOldId"], NoPolis: v["ClaimData.PolicyData.PolicyNo"],
			AcceptType: b["PaymentType"], Kepada: b["PayableTo"], AccountNo: akun, TglAksep: tglAksepTeks(tgl),
			Nett: Teks(kal.Kurang(j.total, j.premiSebar)), Deductible: "0", KaliDeduct: "0", StsSyariah: "0",
			CompanyName: cfg.CompanyName, LjtdId: cfg.LjtdID, LdcId: ldc, StsAp: "0", LkuId: j.cur, LbgID: j.bank,
			TglBolehBayar: TanggalBolehBayar(tgl), Email: email, UserInput: user})
	}
	return out, kal.Galat()
}

// RekeningAngka = HitServiceToKasirKMT_Act S7 `.NoAccount := @pxReplaceAllViaRegex(.NoAccount,"[^0-9]","")`.
func RekeningAngka(no string) string { return bukanAngka.ReplaceAllString(no, "") }

// bulanEmail - SendEmailKlaim_KMT S5 (VERBATIM, termasuk ejaannya).
var bulanEmail = map[string]string{"01": "Januari", "02": "Febuari", "03": "Maret", "04": "April", "05": "Mei",
	"06": "Juni", "07": "July", "08": "Agustus", "09": "September", "10": "Oktober", "11": "November", "12": "Desember"}

// TanggalEmail = S5 `Temp.CARI23`: "dd <Bulan> yyyy" dari DateOfLoss.
func TanggalEmail(dol string) string {
	d := YMD(dol)
	if len(d) < 8 {
		return ""
	}
	return d[6:8] + " " + bulanEmail[d[4:6]] + " " + d[0:4]
}

// SubjekEmail = SendEmailKlaim_KMT S12 / S14: "<awalan>Pengajuan Akseptasi : <ID klaim>/<KomiteNo> <InsuredName> DOL
// <tanggal>". Awalan "" (penyetuju berikut), "(Approval) " (tingkat akhir).
func SubjekEmail(jenis string, kl kontrak.KlaimTreaty, komiteNo string) string {
	awal := ""
	if jenis == EmailPembuatSetuju {
		awal = "(Approval) "
	}
	v := kl.Nilai
	return awal + "Pengajuan Akseptasi : " + v["pyID"] + "/" + komiteNo + " " + v["ClaimData.InsuredName"] + " DOL " +
		TanggalEmail(v["ClaimData.DateOfLoss"])
}

// AdaSpreadingAdjustment - SendEmailKlaim_KMT S1: `.SpreadingAdjustment` kosong -> keluar (tanpa email).
func AdaSpreadingAdjustment(kl kontrak.KlaimTreaty) bool { return len(anakAdj(kl, AnakSpreadIn)) > 0 }

// NoKomite - `.KomiteNo` akseptasi (nomor kasus komite); kosong -> ID kasus komite.
func NoKomite(kl kontrak.KlaimTreaty, komiteID string) string {
	if no := strings.TrimSpace(AdjustmentKlaim(kl)["KomiteNo"]); no != "" {
		return no
	}
	return komiteID
}
