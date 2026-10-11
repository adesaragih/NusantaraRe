package services

// Nilai baris datar OS_AKSEPTASI_KLAIM_LIFE yang ditulis Save to RNM - MURNI (tanpa Oracle).
//
// Keputusan work owner 11-10-2026: "sistem melakukan generate no claim, insert data outstanding ke table
// os_akseptasi_klaim_life (kolom NO_ACCEPTATION diisi dengan tanggal outstanding-nya, policy holder ambil dari policy
// holder milik polisnya, sex peserta, age peserta, plan dari plan polis, begin date peserta, lapse date dari date of
// loss peserta, expired date peserta, wpc milik polis, CEDINGCONAME milik polis, sob, sobname, business id dan business
// name dari polis, dan semua nilai dari peserta)". Urutan / sumber kolom lain mengikuti InsertJsonKlaimLife_sql dan
// SaveOutStandingLife_Act b8772-b10080 (TempInputDetail.CARI1-49).
//
// Dibaca sesudah: repository.KolomCerminSaveRNM (daftar kolom yang ditulis), simpanrnm.go (pemanggil).

import (
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// FormatTanggalOutstanding - NO_ACCEPTATION (VARCHAR2(100) di katalog) = tanggal outstanding, 'DD/MM/YYYY' - bentuk tanggal
// teks Pega (@FormatDateTime "dd/MM/YYYY") dan layar. Work owner dapat meminta bentuk lain; ubah di sini saja.
const FormatTanggalOutstanding = "02/01/2006"

// zonaOutstanding - tanggal outstanding dihitung di WIB (Asia/Jakarta, bentuk @FormatDateTime Pega), bukan UTC server.
var zonaOutstanding = func() *time.Location {
	if z, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return z
	}
	return time.FixedZone("WIB", 7*3600)
}()

// TanggalOutstandingTeks - teks NO_ACCEPTATION dari waktu Save to RNM (waktu yang sama dengan jejak OutStanding Date).
func TanggalOutstandingTeks(saat time.Time) string {
	return saat.In(zonaOutstanding).Format(FormatTanggalOutstanding)
}

// SusunCerminSaveRNM melengkapi baris datar dasar (repository.BarisLamaDari: kunci, nomor klaim, PL, sertifikat,
// Claim RNM, Claim Retro header, Type, pembuat, bank) dengan nilai polis, peserta, dan baris adjustment.
//
//	kolom                  sumber
//	NO_ACCEPTATION         tanggal outstanding (saat Save to RNM), DD/MM/YYYY
//	ACCEPTATION_DATE       tanggal outstanding, DATE (waktu WIB)
//	POLICY_NO              klaim (header); kosong -> polis.NomorPolis
//	POLICY_HOLDER          polis.PolicyHolderName (T_PREMIUM_LIST.POLICY_HOLDER_NAME)
//	SEX, AGE               peserta (baris sumber M_LIFE_PREMIUM_DETAIL)
//	PLAN                   polis.BusinessName (aturan PremiumList: PLAN peserta = T_PREMIUM_LIST.BUSINESS_NAME)
//	BEGIN_DATE/EXPIRED_DATE peserta; LAPSE_DATE = DATE OF LOSS peserta
//	EM_PERCENT             peserta
//	SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE
//	                       baris adjustment (halaman langkah Pega = baris adjustment); kosong -> peserta
//	CLAIM_AMOUNT           Claim RNM baris (dari dasar); CLAIM_RETRO header (dari dasar)
//	WPC, CEDINGCONAME, SOBNAME, BUSINESSID (= BusinessCode), BUSINESSNAME, PRODUCTNAMEID, PRODUCTNAME  polis
//	SOB                    T_PREMIUM_LIST.SOB (dibaca terpisah - PolisRingkas inti tidak membawanya)
//	RETROID/RETRONAME/SECURITYREINSURERID/SECURITYREINSURER  polis; DITUKAR bila Type TP/TR dan Security Reinsurer
//	                       terisi ("Tukar SecurityReinsurer dengan RetroName" b9671, sama dengan ArasapasDilewatiRetro)
//	TYPECEDING             polis.TypeCeding 1/2/3/4 -> QS / SURPLUS / QS + SURPLUS / XOL (CARI40)
//	CONFIRMATION_DATE, CLAIM_RECEIVED_DATE, COMPLETE_DATE  peserta
func SusunCerminSaveRNM(dasar repository.BarisLama, p models.Peserta, b models.BarisAdjustment,
	polis kontrak.PolisRingkas, sob string, saat time.Time) repository.BarisLama {

	d := dasar
	d.NO_ACCEPTATION = TanggalOutstandingTeks(saat)
	// ACCEPTATION_DATE (DATE) = tanggal outstanding juga, waktu WIB (keputusan work owner 11-10-2026).
	d.ACCEPTATION_DATE = saat.In(zonaOutstanding).Format("2006-01-02 15:04:05")
	d.POLICY_NO = pertamaTerisi(d.POLICY_NO, polis.NomorPolis)
	d.POLICY_HOLDER = strings.TrimSpace(polis.PolicyHolderName)
	d.SEX = strings.TrimSpace(p.JenisKelamin)
	d.AGE = strings.TrimSpace(p.Umur)
	d.PLAN = strings.TrimSpace(polis.BusinessName)
	d.BEGIN_DATE = strings.TrimSpace(p.TanggalMulai)
	d.LAPSE_DATE = strings.TrimSpace(p.TanggalKejadian)
	d.EXPIRED_DATE = strings.TrimSpace(p.TanggalExpired)
	d.EM_PERCENT = utils.FormatDecimal(p.EMPercent.Value)
	d.SUM_INSURED = uangAtau(b.SumInsured, p.SumInsured)
	d.CEDING_RETENTION = uangAtau(b.CedingRetention, p.CedingRetention)
	d.SUM_REASURED = uangAtau(b.SumReasured, p.SumReasured)
	d.SHARE_NUSANTARA_RE = uangAtau(b.ShareNusantaraRe, p.ShareNusantaraRe)
	d.SHARE_RETRO = uangAtau(b.ShareRetro, p.ShareRetro)
	d.RETROCEDED_SHARE = uangAtau(b.RetrocededShare, p.RetrocededShare)
	if polis.WPC != nil {
		d.WPC = strings.TrimSpace(*polis.WPC)
	} else {
		d.WPC = ""
	}
	d.CEDINGCONAME = strings.TrimSpace(polis.CedingCoName)
	d.SOB = strings.TrimSpace(sob)
	d.SOBNAME = strings.TrimSpace(polis.SobName)
	d.BUSINESSID = strings.TrimSpace(polis.BusinessCode)
	d.BUSINESSNAME = pertamaTerisi(polis.BusinessName, d.BUSINESSNAME)
	d.PRODUCTNAMEID = strings.TrimSpace(polis.ProductNameID)
	d.PRODUCTNAME = strings.TrimSpace(polis.ProductName)
	d.RETROID, d.RETRONAME = strings.TrimSpace(polis.RetroID), strings.TrimSpace(polis.RetroName)
	d.SECURITYREINSURERID, d.SECURITYREINSURER = strings.TrimSpace(polis.SecurityReinsurerID), strings.TrimSpace(polis.SecurityReinsurer)
	tipe := strings.ToUpper(strings.TrimSpace(polis.Type))
	if (tipe == TypeTP || tipe == TypeTR) && d.SECURITYREINSURERID != "" && d.SECURITYREINSURER != "" {
		d.RETROID, d.SECURITYREINSURERID = d.SECURITYREINSURERID, d.RETROID
		d.RETRONAME, d.SECURITYREINSURER = d.SECURITYREINSURER, d.RETRONAME
	}
	d.TYPECEDING = namaTypeCeding(polis.TypeCeding)
	d.CONFIRMATION_DATE = strings.TrimSpace(p.TanggalKonfirmasi)
	d.CLAIM_RECEIVED_DATE = strings.TrimSpace(p.TanggalTerimaKlaim)
	d.COMPLETE_DATE = strings.TrimSpace(p.TanggalDokumenLengkap)
	return d
}

// namaTypeCeding - CARI40 b9xxx: @if(TypeCeding="1","QS",...,"2","SURPLUS","3","QS + SURPLUS","4","XOL","").
func namaTypeCeding(kode string) string {
	switch strings.TrimSpace(kode) {
	case "1":
		return "QS"
	case "2":
		return "SURPLUS"
	case "3":
		return "QS + SURPLUS"
	case "4":
		return "XOL"
	}
	return ""
}

func pertamaTerisi(nilai ...string) string {
	for _, n := range nilai {
		if t := strings.TrimSpace(n); t != "" {
			return t
		}
	}
	return ""
}

// uangAtau - jumlah teks desimal baris; kosong -> nilai cadangan (peserta); keduanya kosong -> "" (NULL).
func uangAtau(utama, cadangan uang.Money) string {
	if !utama.Kosong() {
		return utils.FormatDecimal(utama.Amount)
	}
	if !cadangan.Kosong() {
		return utils.FormatDecimal(cadangan.Amount)
	}
	return ""
}
