// Package models memuat bentuk data dan aturan murni modul Endorsement Life.
//
// Aturan murni = tanpa Oracle, tanpa HTTP: nilai yang sah, pesan VERBATIM
// korpus, daftar kolom yang dibalik tandanya, dan perakitan nomor
// endorsement. Repository dan services memakainya; uji membuktikannya dari
// korpus (`D:\XML\RNM_BRD\Endorsement Life\`, dilewati bila tak terjangkau).
//
// Bukti `bNNN` mengikuti kepala `docs/PARITAS-LAYAR-DAN-AKSI.md`.
package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Maksud endorsement (`EdmType`) - spec §5 `[keputusan work owner]`: hanya
// `1` dan `3`; `2` dan `4` tidak dipakai. Nilai `3` terbukti dari kode:
// `SetPremi_EDM` 2.3 b1583 `pyWorkPage.EdmType==3`.
const (
	EdmTypePerubahanData = "1"
	EdmTypeBatal         = "3"
)

// LabelEdmType - label opsi `EDM Type` (spec §5; opsi properti tidak
// diekspor, `pyListSource associated` `EndorsmentLife_Section.xml` b1347).
var LabelEdmType = map[string]string{
	EdmTypePerubahanData: "Perubahan Data",
	EdmTypeBatal:         "Batal",
}

// ErrEdmTypeTidakSah - nilai di luar `1`/`3` (AC 11, AC 58).
var ErrEdmTypeTidakSah = errors.New("models: EDM Type only accepts 1 (Perubahan Data) or 3 (Batal)")

// PeriksaEdmType menolak nilai di luar `1`/`3`, tanpa merapikan spasi: nilai
// datang dari pilihan, bukan ketikan.
func PeriksaEdmType(v string) error {
	if v != EdmTypePerubahanData && v != EdmTypeBatal {
		return fmt.Errorf("%w: %q", ErrEdmTypeTidakSah, v)
	}
	return nil
}

// Status per peserta (properti `EDMStatus` → kolom `EDM_STATUS`, dan kolom status
// EDM tabel warisan peserta) - spec §5.
const (
	// StatusOld - warisan versi sebelumnya (`MappingEDMLife` 11.1 b2609).
	StatusOld = "Old"
	// StatusNew - peserta tambahan CSV, hanya pada `EdmType=1`
	// (`SaveCSVEDMLife` 4.3 b2715).
	StatusNew = "New"
	// StatusDelete - peserta dikeluarkan (`SetPremi_EDM` 2.2 b1385).
	StatusDelete = "Delete"
	// StatusBatal - seluruh peserta polis batal (`SetPremi_EDM` 2.3 b1507).
	StatusBatal = "Batal"
)

// StatusHidup menjawab apakah peserta berstatus EDM ini masih ditanggung -
// penyaring hilir Claim Life (spec §14): kosong/NULL, `Old`, `New` hidup;
// `Delete`, `Batal` mati. ⛔ Kosong HIDUP: baris new business tidak pernah
// mengisi kolom ini (AC 48a).
func StatusHidup(s string) bool { return s != StatusDelete && s != StatusBatal }

// Status kerja kasus (`pyStatusWork` → `T_PREMIUM_LIST.STATUSS`, RALAT R20).
// Kosong = terbuka; kedua nilai akhir VERBATIM `Flow/InputEDMLife.xml`.
const (
	StatusKasusTerbuka  = ""
	StatusKasusSelesai  = "Resolved-Completed" // END52 b686
	StatusKasusDitolak  = "Resolved-Rejected"  // End1 b642
	TahapInputEDMLife   = "InputEDMLife"       // assignment b731
	AwalanKasus         = "EDMLF-"             // spec §16 AC 72
	panjangMaksPengenal = 32                   // `T_PREMIUM_LIST.ID VARCHAR2(32)`
)

// AksiSimpan - komentar jejak `Save` (label tombol b37202): tahap tidak berpindah,
// tetapi jurnal balik dicatat siapa dan kapan (`T_PREMIUM_LIST` tanpa kolom pengubah).
const AksiSimpan = "Save"

// RakitPengenalKasus menyusun `EDMLF-<n>` dari angka sequence, tanpa padding
// - sejajar `NBLF-<n>` PremiumList.
func RakitPengenalKasus(urut string) (string, error) {
	u := strings.TrimSpace(urut)
	if _, err := strconv.ParseUint(u, 10, 64); err != nil || u == "" {
		return "", fmt.Errorf("models: sequence number %q is not a positive integer", urut)
	}
	id := AwalanKasus + u
	if len(id) > panjangMaksPengenal {
		return "", fmt.Errorf("models: case id %q exceeds %d bytes", id, panjangMaksPengenal)
	}
	return id, nil
}

// KasusEDM menjawab apakah sebuah pengenal adalah kasus modul ini.
func KasusEDM(id string) bool {
	return strings.HasPrefix(id, AwalanKasus) && len(id) > len(AwalanKasus)
}

// Jenis transaksi (`.Type`) - QR/QP memakai kolom gross, TP/TR retro.
const (
	TypeQR = "QR"
	TypeQP = "QP"
	TypeTP = "TP"
	TypeTR = "TR"
)

// StatusJenis adalah `CARI48` `InsertJsonPolisLife_Act` 11.2 b4568 (`QR`/`QP`
// → `0`) dan 11.3 b4706 (`TR`/`TP` → `1`), ditulis ke kolom `STATUS`.
//
// ⚠️ Tipe lain: kedua prakondisi gagal, `CARI48` tidak disetel - nilai
// sebelumnya di halaman dipakai ulang. Di sini kosong (NULL).
func StatusJenis(tipe string) string {
	switch tipe {
	case TypeQR, TypeQP:
		return "0"
	case TypeTR, TypeTP:
		return "1"
	}
	return ""
}

// StatusLama adalah `CARI47`: `1` bila peserta `Old` (11.4 b4844), selain itu
// `0` (11.5 b4982), ditulis ke `STATUSOLD`/`STATUS_OLD`.
func StatusLama(edmStatus string) string {
	if edmStatus == StatusOld {
		return "1"
	}
	return "0"
}

// KolomJurnalBalik adalah 32 kolom uang yang `SetPremi_EDM` 2.1 kalikan `-1`,
// VERBATIM urutannya (b551, prakondisi b1273). ⛔ Deskripsi langkahnya
// berbunyi "Set 0 jika EDM Batal" - kodenya yang berlaku (spec §6).
//
// ⚠️ Kolom uang lain (`COMM`, `OVR_COMM`, `TAX`, `PROF_COMM`, `CLAIM`,
// `FLEET_DISCOUNT`, `*_REFUND` komisi/tax, `OVR_COMM_*RETRO`) TIDAK dibalik
// di korpus, meski spec user story 21 menyebut "komisi, tax" (RALAT R28).
var KolomJurnalBalik = []string{
	"SUM_INSURED", "CEDING_RETENTION", "SUM_REASURED", "SHARE_NUSANTARA_RE_GROSS",
	"SUM_AT_RISK_GROSS", "GROSS_PREMIUM", "DEDUCTION", "NET_PREMIUM", "FACTOR",
	"CLAIM_AMOUNT", "GROSS_PREMIUM_REFUND", "DEDUCTION_REFUND", "NET_PREMIUM_REFUND",
	"SHARE_NUSANTARA_RE", "SUM_AT_RISK_RETRO", "RETROCEDED_SHARE", "SHARE_RETRO", "RATE",
	"GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO",
	"RI_ADMIN_FEE_REFUND_RETRO", "BROKERAGE_FEE_RETRO", "NET_PREMIUM_REFUND_RETRO",
	"GROSS_PREMIUM_RETRO", "RI_ADMIN_FEE_RETRO", "NET_PREMIUM_RETRO", "RI_ADMIN_FEE",
	"RI_ADMIN_FEE_REFUND", "BROKERAGE_FEE", "BROKERAGE_FEE_REFUND", "BROKERAGE_FEE_REFUND_RETRO",
}

// NomorEndorsement merakit `PL_NUMBER_EDM` = `NOENDORS` dari nomor polis dan
// `PRODKE` versi berjalan - `GenerateNoEDM_Life` 4 b872 (`CARI4 = Prodke+1`
// b946, `CARI14` pad nol bila satu digit b967) lalu `Generate_NoEndorsmentLife`
// b84 `NOPOLIS||'/'||{InputData.CARI14}`.
//
// ⛔ BUKAN dari `PROC_GENERATE_SEQUENCE_NUMBER` (ADR-0006 tidak berlaku).
// Mengembalikan nomor dan `PRODKE` versi baru.
func NomorEndorsement(nomorPolis string, prodKeBerjalan int) (string, int, error) {
	np := strings.TrimSpace(nomorPolis)
	if np == "" {
		return "", 0, errors.New("models: policy number is empty")
	}
	if prodKeBerjalan < 1 {
		return "", 0, fmt.Errorf("models: current production number %d is not a version", prodKeBerjalan)
	}
	baru := prodKeBerjalan + 1
	teks := strconv.Itoa(baru)
	if len(teks) == 1 {
		teks = "0" + teks
	}
	return np + "/" + teks, baru, nil
}

// NomorInvoiceArasapas - `SetErrorBatalEndorsement_Act` 1 b464
// `@replaceAll(TempWork.PolicyNo,".","")`: nomor polis tanpa titik.
func NomorInvoiceArasapas(nomorPolis string) string {
	return strings.ReplaceAll(nomorPolis, ".", "")
}

// Keputusan `Status` (`EmailTypePL`) di `ConfirmSection` b496.
const (
	// KeputusanAccept - `IsLifeAccepted` b284 `1` → `Confirm` b317.
	KeputusanAccept = "1"
	// KeputusanDecline - nilai kedua tombol `Submit` b38109 (`2`).
	KeputusanDecline = "2"
	// KeputusanTujuh - nilai ketiga tombol yang sama (`7`), tanpa label
	// (OQ-EDM-006): diterima, jalurnya Decline.
	KeputusanTujuh = "7"
)

// ErrKeputusanTidakSah - nilai `Status` di luar ketiga nilai korpus.
var ErrKeputusanTidakSah = errors.New("models: decision status only accepts 1, 2, or 7")

// PeriksaKeputusan menolak nilai `Status` di luar `1`, `2`, `7`.
func PeriksaKeputusan(v string) error {
	switch v {
	case KeputusanAccept, KeputusanDecline, KeputusanTujuh:
		return nil
	}
	return fmt.Errorf("%w: %q", ErrKeputusanTidakSah, v)
}

// Diterima - `IsLifeAccepted`: `1` → Confirm, selain itu (bawaan b86) Decline.
func Diterima(keputusan string) bool { return keputusan == KeputusanAccept }

// LabelKeputusanRiwayat - `AddHistorySuggest` b353
// `@if(.EmailTypePL=1,"Accept","Decline")` → `IS_CEDING_CONFIRM`.
func LabelKeputusanRiwayat(keputusan string) string {
	if Diterima(keputusan) {
		return "Accept"
	}
	return "Decline"
}
