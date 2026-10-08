package models

// Untuk apa berkas ini: ISI EMAIL dan DOKUMEN AKSEPTASI - dua stream HTML yang diekspor work owner 08-10-2026 ke korpus
// `Komite Claim Prop`:
//
//	EmailKlaim_HTML_KMT  SendEmailKlaim_KMT S5-S16 (Temp.CARI*, InputData.pxResults) -> badan email `Param.Message`
//	FILEAcceptanceNote   PrintFileAcceptance_TKMT S5-S10 (halaman TempAcceptedNo) -> markup yang S11 ubah jadi PDF
//
// Markup VERBATIM di `templat/`; nilai dirakit di sini (murni). Kedua isi DIRAKIT SAAT EFEK DIKIRIM dari pengenal di
// MUATAN outbox (layanan `SusunEmailKomite` / `SusunDokumenAkseptasi`): MUATAN T_LOG_SERVICE_RNM hanya memuat pengenal
// dan angka, tanpa nama atau alamat (claimlife/015).
//
// `[penyimpangan sadar]` tanggal: `pega:reference` properti Date (dan `format="date"`) memakai aturan tampilan Pega
// yang tidak diekspor - di sini "dd/MM/yyyy" (pola yang dideklarasikan stream FILEAcceptanceNote baris 3).

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

//go:embed templat/*.html
var berkasTemplat embed.FS

var (
	templatEmail   = template.Must(template.ParseFS(berkasTemplat, "templat/email_klaim_kmt.html"))
	templatDokumen = template.Must(template.ParseFS(berkasTemplat, "templat/file_acceptance_note.html"))
)

// AngkaPega = `DecimalFormat("#,###.####")` locale id_ID (kedua stream): ribuan ".", desimal ",", paling banyak empat
// angka desimal dibulatkan HALF_EVEN (bawaan DecimalFormat), tanpa nol ekor. Pola tanpa digit '0': bagian bulat nol
// tidak dicetak (0,5 -> ",5"); nol -> "0". Kosong = 0.
func AngkaPega(jalur, s string) (string, error) {
	d, err := desimal(jalur, s)
	if err != nil {
		return "", err
	}
	c := *konteksUang
	c.Rounding = apd.RoundHalfEven
	r := new(apd.Decimal)
	if _, err := c.Quantize(r, d, -4); err != nil {
		return "", fmt.Errorf("models: %s: %w", jalur, err)
	}
	r.Reduce(r)
	t := r.Text('f')
	tanda := ""
	if strings.HasPrefix(t, "-") {
		tanda, t = "-", t[1:]
	}
	bulat, pecahan, _ := strings.Cut(t, ".")
	if bulat == "0" && pecahan != "" {
		bulat = ""
	}
	var b strings.Builder
	for i, c := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	if pecahan != "" {
		b.WriteString("," + pecahan)
	}
	return tanda + b.String(), nil
}

// bagiPega = `@divide(x, 1, skala)` - pembulatan HALF_UP (konvensi repo, `claimlife` adjustment.go 7.8).
func bagiPega(jalur, s string, skala int32) (string, error) {
	d, err := desimal(jalur, s)
	if err != nil {
		return "", err
	}
	c := *konteksUang
	c.Rounding = apd.RoundHalfUp
	r := new(apd.Decimal)
	if _, err := c.Quantize(r, d, -skala); err != nil {
		return "", fmt.Errorf("models: %s: %w", jalur, err)
	}
	return r.Text('f'), nil
}

// TanggalSurat - nilai tanggal halaman ("yyyyMMdd..." / "yyyy-MM-dd...") sebagai "dd/MM/yyyy"; kosong tetap kosong.
func TanggalSurat(v string) string {
	d := YMD(v)
	if len(d) < 8 {
		return ""
	}
	return d[6:8] + "/" + d[4:6] + "/" + d[0:4]
}

// ---------------------------------------------------------------- email

// KonfigurasiEmail - SendEmailKlaim_KMT S4 (`Local.EmailCC`, hanya IsPEGAPROD) dan S17-S18 (`pyNotifyAccountName`):
// konfigurasi berdokumen (`konfigurasi/email.json`, MODUL.md), bukan literal kode. BCC pribadi S3 tidak disalin.
type KonfigurasiEmail struct {
	Akun, AkunSyariah, CC string
}

// AkunUntuk - S17 akun bawaan; S18 `@contains(Local.Emailto,"syariah")` -> akun syariah.
func (c KonfigurasiEmail) AkunUntuk(kepada string) string {
	if strings.Contains(kepada, "syariah") {
		return c.AkunSyariah
	}
	return c.Akun
}

// SurelKomite - satu email SendEmailKlaim_KMT S17-S19 (`SendEmailWithAttachments`, HTMLmessage).
type SurelKomite struct {
	Akun, Kepada, CC, Subjek, HTML string
}

// BarisEmailSpreading - satu `InputData.pxResults` S11.1-S11.3.
type BarisEmailSpreading struct {
	No, Currency, TreatyName, Share, ClaimSpreaded string
}

// DataEmailKomite - nilai stream EmailKlaim_HTML_KMT.
type DataEmailKomite struct {
	Sapaan              string // Temp.CARI14: S12 IDKomite penyetuju berikut; S14-S15 pyUserName pembuat
	KlaimID, KomiteNo   string // pyWorkCover.pyID / Temp.CARI10
	DateOfLoss          string // pyWorkCover.ClaimData.DateOfLoss
	PolicyNo            string // Temp.CARI24
	Ceding, InsuredName string // Temp.CARI25 / CARI26
	SOB                 string // Temp.CARI27
	Baris               []BarisEmailSpreading
	TotalShare          string // Temp.CARI11 = @divide(.TotalSharePersen,1,0)
	TotalClaim          string // Temp.CARI12 = .TotalSpreadAdjustment
	Status              string // Temp.CARI30
	AcceptedNo          string // Temp.CARI16 (S14-S15)
	Remarks             string // Temp.CARI31 = pyWorkPage.Comment
	Pengirim            string // Temp.CARI13 = OperatorID.pyUserName
}

// BahanEmail - nilai yang tidak tinggal di halaman klaim: keputusan, komentar, dan nama dari tangga / M_LOGIN_GO.
type BahanEmail struct {
	Jenis      string // EmailPenyetujuBerikut / EmailPembuatSetuju / EmailPembuatTolak
	KomiteID   string
	Keputusan  string // pyWorkPage.AcceptStatus
	Komentar   string
	Sapaan     string
	Pengirim   string
	AcceptedNo string
}

// SusunDataEmail = SendEmailKlaim_KMT S5-S15 jalur CLMP (S7). Ceding / SOB S7 = `ClaimData.QuotationData.CedingCoName`
// / `SobName`, yang SetValueToClaim_Act Claim Prop isi dari master treaty yang SAMA dengan `TreatyInMaster.Ceding` /
// `LeadingReinsSource` (registrasi.go) dan tidak disimpan - dibaca dari `TreatyInMaster`.
func SusunDataEmail(kl kontrak.KlaimTreaty, b BahanEmail) (DataEmailKomite, error) {
	v := kl.Nilai
	d := DataEmailKomite{Sapaan: b.Sapaan, KlaimID: v["pyID"], KomiteNo: NoKomite(kl, b.KomiteID),
		DateOfLoss: TanggalSurat(v["ClaimData.DateOfLoss"]), PolicyNo: v["ClaimData.PolicyData.PolicyNo"],
		Ceding: v["TreatyInMaster.Ceding"], InsuredName: v["ClaimData.InsuredName"],
		SOB: v["TreatyInMaster.LeadingReinsSource"], Status: "Rejected", Remarks: b.Komentar, Pengirim: b.Pengirim}
	if b.Keputusan == KeputusanSetuju { // S9 @if(.AcceptStatus==1,"Accepted","Rejected")
		d.Status = "Accepted"
	}
	if b.Jenis != EmailPenyetujuBerikut { // S14-S15 Temp.CARI16 := .AcceptedNo
		d.AcceptedNo = b.AcceptedNo
	}
	// CountSpreadingADJ_Act S6.3 / S10: TotalSharePersen dan TotalSpreadAdjustment = jumlah `.SpreadingAdjustment`
	// baris adjustment (Claim Prop tidak menyimpannya).
	share, klaim := apd.New(0, 0), apd.New(0, 0)
	for i, s := range kl.Daftar[jalurAdj(kl.Adjustment, "SpreadingAdjustment")] {
		sp, err := bagiPega("SharePercentage", s["SharePercentage"], 0)
		if err != nil {
			return DataEmailKomite{}, err
		}
		cs, err := AngkaPega("ClaimSpreaded", s["ClaimSpreaded"])
		if err != nil {
			return DataEmailKomite{}, err
		}
		d.Baris = append(d.Baris, BarisEmailSpreading{No: strconv.Itoa(i + 1), Currency: s["Currency"],
			TreatyName: s["TreatyName"], Share: sp, ClaimSpreaded: cs})
		for _, x := range []struct {
			tot  *apd.Decimal
			prop string
		}{{share, "SharePercentage"}, {klaim, "ClaimSpreaded"}} {
			n, err := desimal(x.prop, s[x.prop])
			if err != nil {
				return DataEmailKomite{}, err
			}
			if _, err := konteksUang.Add(x.tot, x.tot, n); err != nil {
				return DataEmailKomite{}, err
			}
		}
	}
	var err error
	if d.TotalShare, err = bagiPega("TotalSharePersen", TeksAngka(share), 0); err != nil {
		return DataEmailKomite{}, err
	}
	if d.TotalClaim, err = AngkaPega("TotalSpreadAdjustment", TeksAngka(klaim)); err != nil {
		return DataEmailKomite{}, err
	}
	return d, nil
}

// RenderEmailKomite - S16 Property-Set-HTML `EmailKlaim_HTML_KMT`.
func RenderEmailKomite(d DataEmailKomite) (string, error) {
	var b bytes.Buffer
	if err := templatEmail.Execute(&b, d); err != nil {
		return "", fmt.Errorf("models: EmailKlaim_HTML_KMT: %w", err)
	}
	return b.String(), nil
}

// ---------------------------------------------------------------- dokumen akseptasi

// Konstanta PrintFileAcceptance_TKMT S9 dan InsertDocument_Act S4.
const (
	KategoriDokumenAkseptasi = "AcceptanceNote"
	FolderDokumenKlaim       = "Claim"
)

// NamaBerkasAkseptasi - S9 `param.PDFName` (VERBATIM, termasuk spasinya).
func NamaBerkasAkseptasi(acceptedNo string) string {
	return "Persetujuan Klaim" + "  " + " AcceptNo " + acceptedNo + ".pdf"
}

// BarisTSI - satu `ClaimData.InterestList` baris "TSI RNM".
type BarisTSI struct{ ObjectName, Currency, Nilai string }

// BarisSpread - satu baris tabel "Spreading Adjustment" / "BreakDown Spreading (QS)".
type BarisSpread struct{ TreatyName, Share, ClaimSpread string }

// DataAcceptanceNote - halaman `TempAcceptedNo` stream FILEAcceptanceNote.
type DataAcceptanceNote struct {
	AcceptedNo, PolicyNo                                      string
	TampilLini                                                bool // pega:when TempAcceptedNo.IsTreatyIn!='1'
	LineOfBusiness                                            string
	NoClaim, ClaimID, SOB, Ceding, InsuredName                string
	InterestInsured                                           string
	TSI                                                       []BarisTSI
	PeriodeAwal, PeriodeAkhir, DateOfLoss, CauseOfLoss        string
	Location                                                  string
	Currency, AcceptedClaim                                   string
	PayableTo, NameOfBank, SwiftCode, BranchOfBank, NoAccount string
	PIC                                                       string
	Spreading                                                 []BarisSpread
	TampilQS                                                  bool // NameSpread == "10196"
	QS                                                        []BarisSpread
	TanggalCetak, Pembuat                                     string
}

// TreatyTypeQS - TreatyType yang memunculkan tabel "BreakDown Spreading (QS)" (stream baris 225 / 242).
const TreatyTypeQS = "10196"

// SusunAcceptanceNote = PrintFileAcceptance_TKMT S5-S8 + isi stream. `adj` = baris adjustment sesudah keputusan (berisi
// AcceptedNo); `saat` = S4 `NameInput.CARI24` (@CurrentDateTime); `pembuat` = `OperatorID.pyUserName`.
func SusunAcceptanceNote(kl kontrak.KlaimTreaty, adj map[string]string, saat time.Time, pembuat string) (
	DataAcceptanceNote, error) {
	v := kl.Nilai
	d := DataAcceptanceNote{AcceptedNo: adj["AcceptedNo"], PolicyNo: v["ClaimData.PolicyData.PolicyNo"],
		TampilLini: v["IsTreatyIn"] != "1", LineOfBusiness: v["OfferFacIn.QuotationData.BusinessName"],
		NoClaim: v["ClaimData.NoClaim"], ClaimID: v["pyID"], SOB: v["TreatyInMaster.LeadingReinsSource"],
		Ceding: v["TreatyInMaster.Ceding"], InsuredName: v["ClaimData.InsuredName"],
		PeriodeAwal:  TanggalSurat(v["ClaimData.PolicyData.StartDateTime"]),
		PeriodeAkhir: TanggalSurat(v["ClaimData.PolicyData.EndDateTime"]),
		DateOfLoss:   TanggalSurat(v["ClaimData.DateOfLoss"]), CauseOfLoss: v["ClaimData.CauseOfLoss"],
		// S8: `.Location := TempMainPage.ClaimData.Location` lalu DITIMPA `.Location := ...CauseOfLoss` (XML apa adanya).
		Location: v["ClaimData.CauseOfLoss"],
		Currency: adj["Currency"], PayableTo: adj["PayableTo"], NameOfBank: adj["NameOfBank"],
		SwiftCode: adj["SwiftCode"], BranchOfBank: adj["BranchOfBank"], NoAccount: adj["NoAccount"],
		PIC:          adj["pxCreateOpName"],
		TanggalCetak: saat.In(Jakarta).Format("02 January 2006"), // dd MMMM yyyy, Asia/Jakarta, en_US
		Pembuat:      pembuat}
	var err error
	if d.AcceptedClaim, err = AngkaPega("AdjustmentValue", adj["AdjustmentValue"]); err != nil {
		return DataAcceptanceNote{}, err
	}
	for i, b := range kl.Daftar["ClaimData.InterestList"] {
		if i == 0 {
			d.InterestInsured = b["ObjectName"]
		}
		n, err := AngkaPega("TSIPerObject", b["TSIPerObject"])
		if err != nil {
			return DataAcceptanceNote{}, err
		}
		d.TSI = append(d.TSI, BarisTSI{ObjectName: b["ObjectName"], Currency: b["Currency"], Nilai: n})
	}
	spread := func(jalur string) ([]BarisSpread, error) {
		var out []BarisSpread
		for _, s := range kl.Daftar[jalurAdj(kl.Adjustment, jalur)] {
			sh, err := AngkaPega("SharePercentage", s["SharePercentage"])
			if err != nil {
				return nil, err
			}
			cs, err := AngkaPega("ClaimSpreaded", s["ClaimSpreaded"])
			if err != nil {
				return nil, err
			}
			if jalur == "SpreadingAdjustment" && s["TreatyType"] == TreatyTypeQS {
				d.TampilQS = true
			}
			out = append(out, BarisSpread{TreatyName: s["TreatyName"], Share: sh, ClaimSpread: cs})
		}
		return out, nil
	}
	if d.Spreading, err = spread("SpreadingAdjustment"); err != nil {
		return DataAcceptanceNote{}, err
	}
	if d.TampilQS {
		if d.QS, err = spread("SpreadingQuotaShare"); err != nil {
			return DataAcceptanceNote{}, err
		}
	}
	return d, nil
}

// RenderAcceptanceNote - S10 Property-Set-HTML `FILEAcceptanceNote` (markup sebelum S11 HTMLToPDF).
func RenderAcceptanceNote(d DataAcceptanceNote) (string, error) {
	var b bytes.Buffer
	if err := templatDokumen.Execute(&b, d); err != nil {
		return "", fmt.Errorf("models: FILEAcceptanceNote: %w", err)
	}
	return b.String(), nil
}
