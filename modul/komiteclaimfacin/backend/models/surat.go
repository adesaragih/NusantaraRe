package models

// Untuk apa berkas ini: EMAIL komite - subjek, penerima, dan akun pengirim `SendEmailKlaim_KMT` (TT2, KomitePost_Adjustment
// 7.2.1.15) dan KomitePost_Reject / KomitePost_CloseClaim S13.10 / S11.10 (TT3 / TT4). Isi DIRAKIT SAAT EFEK DIKIRIM dari
// pengenal di MUATAN outbox (claimlife/015: MUATAN tanpa nama / alamat).
//
// Badan email (stream `EmailKlaim_HTML_KMT`, `CommitteReject_CC`, `CommitteCloseClaim`) TIDAK ada di korpus Komite
// Claim FacIn - isinya tidak dikarang (OQ-KCFI-01). BCC pribadi (S3 / S13.10.7) tidak disalin. Ejaan bulan VERBATIM
// ("Febuari", "July" - dibaca penerima, prompt §5).

import (
	"strings"
)

// KonfigurasiEmail - akun notifikasi dan CC kotak surat klaim (`konfigurasi/email.json`).
type KonfigurasiEmail struct {
	Akun, AkunSyariah, CC string
}

// namaBulan - SendEmailKlaim_KMT S5 / KomitePost_Reject S13.10.2 (VERBATIM).
var namaBulan = map[string]string{"01": "Januari", "02": "Febuari", "03": "Maret", "04": "April", "05": "Mei",
	"06": "Juni", "07": "July", "08": "Agustus", "09": "September", "10": "Oktober", "11": "November", "12": "Desember"}

// TanggalSurat - `@substring(DOL,6,8)+" "+<bulan>+" "+@substring(DOL,0,4)`.
func TanggalSurat(s string) string {
	t := ymd(s)
	if t == "" {
		return ""
	}
	return t[6:8] + " " + namaBulan[t[4:6]] + " " + t[0:4]
}

// SurelKomite - satu email komite yang siap dikirim (badan menunggu stream, OQ-KCFI-01).
type SurelKomite struct {
	Akun, Subjek, CC string
	Kepada           []string
}

// SubjekSurel - subjek email. TT2: S12 "Pengajuan Akseptasi : " (tingkat berikut), S14 "(Approval) ...", S15
// "(Reject) ..."; TT3 / TT4 S13.10.5-S13.10.6 "(Approval) / (Reject) Pengajuan Reject / Close Klaim : ".
// `CLM/KMT nama DOL tanggal` - Temp.CARI21 / CARI10 / CARI22 / CARI23.
func SubjekSurel(transfer, jenis, klaimID, kmt, nama, dol string) string {
	ekor := klaimID + "/" + kmt + " " + nama + " DOL " + TanggalSurat(dol)
	if transfer == TransferReject || transfer == TransferClose {
		sub := "Close"
		if transfer == TransferReject { // S13.10.1 `@if(TransferType = "3", "Reject", "Close")`
			sub = "Reject"
		}
		awal := "(Reject) "
		if jenis == EmailPembuatSetuju {
			awal = "(Approval) "
		}
		return awal + "Pengajuan " + sub + " Klaim : " + ekor
	}
	switch jenis {
	case EmailPembuatSetuju:
		return "(Approval) Pengajuan Akseptasi : " + ekor
	case EmailPembuatTolak:
		return "(Reject) Pengajuan Akseptasi : " + ekor
	}
	return "Pengajuan Akseptasi : " + ekor
}

// NamaTertanggungSurel - Temp.CARI22: TT2 IsCLM `OfferFacIn.QuotationData.InsuredName` (S6); TT3 / TT4
// `TempOpenPage.ClaimData.InsuredName` (S13.10.2).
func NamaTertanggungSurel(transfer string, nilai map[string]string) string {
	if transfer == TransferReject || transfer == TransferClose {
		return nilai["ClaimData.InsuredName"]
	}
	return nilai[OQ+"InsuredName"]
}

// AkunSurel - S17 "NUSARE"; S18 (TT2) penerima memuat "syariah" -> "NUSARESYARIAH". TT3 / TT4 S13.10.7 selalu NUSARE.
func AkunSurel(transfer string, kepada []string, c KonfigurasiEmail) string {
	if transfer == TransferAdjustment {
		for _, e := range kepada {
			if strings.Contains(e, "syariah") {
				return c.AkunSyariah
			}
		}
	}
	return c.Akun
}

// CCSurel - TT2 S4 `claim@...` (hanya IsPEGAPROD); TT3 / TT4 S13.10.1 CC kosong.
func CCSurel(transfer string, c KonfigurasiEmail) string {
	if transfer == TransferAdjustment {
		return c.CC
	}
	return ""
}
