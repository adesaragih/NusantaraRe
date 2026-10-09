package models

// Untuk apa berkas ini: ISI EMAIL komite - stream `EmailKlaim_HTML_KMT` (aturan kelas Data-Adjustment yang sama dipakai
// Komite Claim Prop dan Komite Claim Non Prop; markup VERBATIM disalin dari ekspor korpus `Komite Claim Prop`) yang diisi
// `SendEmailKlaim_KMT` cabang IsCLMNP (korpus `Komite Claim Non Prop`). Isi DIRAKIT SAAT EFEK DIKIRIM dari pengenal di
// MUATAN outbox: MUATAN T_LOG_SERVICE_RNM hanya memuat pengenal, tanpa nama atau alamat (claimlife/015).
//
// Dokumen PDF persetujuan (`GenerateAccCNP_act`, stream `AccClaimKomite_HTML`) TIDAK dibuat: stream-nya tidak diekspor
// (OQ-CNP-22).
//
// `[penyimpangan sadar]` tanggal: `pega:reference` properti Date memakai aturan tampilan Pega yang tidak diekspor - di
// sini "dd/MM/yyyy" (pola Komite Claim Prop).

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

//go:embed templat/*.html
var berkasTemplat embed.FS

var templatEmail = template.Must(template.ParseFS(berkasTemplat, "templat/email_klaim_kmt.html"))

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
	Jenis      string // EmailPenyetujuBerikut / EmailPembuatSetuju
	KomiteID   string
	Keputusan  string // pyWorkPage.AcceptStatus
	Komentar   string
	Sapaan     string
	Pengirim   string
	AcceptedNo string
}

// SusunDataEmail = SendEmailKlaim_KMT S5-S15 jalur CLMNP (S8): Ceding / SOB = `TreatyInMaster.Ceding` /
// `LeadingReinsSource`; nilai baris spreading (S11.2 IsCLMNP) = `.TotalClaim` Spreading In.
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
	// TotalSharePersen dan TotalSpreadAdjustment (S10) tanpa penulis di korpus Non Prop: jumlah `.SpreadingAdjustment`
	// baris akseptasi seperti Claim Prop `[inferensi]` (PARITAS) - Share dan Total Claim.
	share, klaim := apd.New(0, 0), apd.New(0, 0)
	for i, s := range anakAdj(kl, AnakSpreadIn) {
		sp, err := bagiPega("SharePercentage", s["SharePercentage"], 0)
		if err != nil {
			return DataEmailKomite{}, err
		}
		cs, err := AngkaPega("TotalClaim", s["TotalClaim"])
		if err != nil {
			return DataEmailKomite{}, err
		}
		d.Baris = append(d.Baris, BarisEmailSpreading{No: strconv.Itoa(i + 1), Currency: s["Currency"],
			TreatyName: s["TreatyName"], Share: sp, ClaimSpreaded: cs})
		for _, x := range []struct {
			tot  *apd.Decimal
			prop string
		}{{share, "SharePercentage"}, {klaim, "TotalClaim"}} {
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
