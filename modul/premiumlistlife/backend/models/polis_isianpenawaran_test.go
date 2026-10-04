package models

// Uji isian layar Input Offer - tiket 01 bagian 3. MURNI.

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func isianLengkap() IsianPenawaran {
	return IsianPenawaran{
		CedingCo: "UJI-L01", CedingCoName: "UJI-CEDING",
		PolicyHolder: "UJI-C1", PolicyHolderName: "UJI-PEMEGANG",
		TypeCeding: "2", BusinessCode: "L3", Status: "Pending", Description: "UJI-KOMENTAR",
		DateReceived: tanggalTerimaUji(),
	}
}

func tanggalTerimaUji() *time.Time {
	t := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	return &t
}

// Nama turunan dari kode - VERBATIM InputOfferLife_ACT langkah 3 dan SetCoBName_Act.
func TestNamaTurunanPenawaran(t *testing.T) {
	for kode, mau := range map[string]string{"1": "QS", "2": "SURPLUS", "3": "QS + SURPLUS", "4": "XOL", "5": "", "": ""} {
		if got := NamaTypeCeding(kode); got != mau {
			t.Errorf("NamaTypeCeding(%q) = %q, mau %q", kode, got, mau)
		}
	}
	if got := NamaBusiness("L16"); got != "INDIVIDU CRITICAL ILNESS" {
		t.Errorf("ejaan warisan tidak dipertahankan: %q", got)
	}
	if got := NamaBusiness("L99"); got != NamaBusinessTakDikenal {
		t.Errorf("cabang OTHERWISE: %q", got)
	}
	if len(PilihanClassOfBusiness) != 21 {
		t.Errorf("SetCoBName_Act punya 21 cabang WHEN, daftar = %d", len(PilihanClassOfBusiness))
	}
}

// SetReinsuranceType: hanya "4" (XOL) yang Non Proportional - kosong pun
// Proportional; dan SusunPenawaran menurunkannya, bukan menerimanya.
func TestJenisAsuransiTurunanSystemReinsurance(t *testing.T) {
	for kode, mau := range map[string]string{"4": JenisAsuransiNonProporsional, "1": JenisAsuransiProporsional,
		"3": JenisAsuransiProporsional, "": JenisAsuransiProporsional} {
		if got := JenisAsuransi(kode); got != mau {
			t.Errorf("JenisAsuransi(%q) = %q, mau %q", kode, got, mau)
		}
	}
	isi := isianLengkap()
	isi.TypeCeding = "4"
	if got, err := SusunPenawaran(FlagPolisPenawaran, isi); err != nil || got.JenisAsuransi != JenisAsuransiNonProporsional {
		t.Errorf("XOL = %q, %v - mau Non Proportional", got.JenisAsuransi, err)
	}
}

// Wajib-isi mengikuti pyRequired sel, dan SELURUH yang kosong dilaporkan sekaligus.
func TestSusunPenawaranMelaporkanSeluruhKekurangan(t *testing.T) {
	_, err := SusunPenawaran(FlagPolisPenawaran, IsianPenawaran{Description: "  "})
	if !errors.Is(err, ErrIsianPenawaranBelumLengkap) {
		t.Fatalf("galat = %v", err)
	}
	for _, pesan := range []string{PesanCedingKosong, PesanPolicyHolderKosong, PesanTypeCedingKosong,
		PesanBusinessCodeKosong, PesanDateReceivedPenawaranKosong, PesanStatusKosong, PesanCommentKosong} {
		if !strings.Contains(err.Error(), pesan) {
			t.Errorf("pesan %q tidak dilaporkan: %v", pesan, err)
		}
	}
	// Gerbang Confirm tidak memeriksa Status - kolomnya tidak ada di header.
	k := KekuranganPenawaran(WajibPenawaran{CedingCoName: "UJI-C", PolicyHolderName: "UJI-P",
		TypeCeding: "1", BusinessCode: "L1", Description: "UJI", DateReceived: tanggalTerimaUji()})
	if len(k) != 0 {
		t.Errorf("Confirm tanpa Status ditolak: %v", k)
	}
	if k := KekuranganPenawaran(WajibPenawaran{PeriksaStatus: true}); len(k) != 7 || k[0] != PesanCedingKosong || k[4] != PesanDateReceivedPenawaranKosong {
		t.Errorf("urutan pesan bukan urutan layar: %v", k)
	}
}

func TestSusunPenawaranMenolakKodeDiLuarPilihan(t *testing.T) {
	for nama, ubah := range map[string]func(*IsianPenawaran){
		"type ceding":    func(i *IsianPenawaran) { i.TypeCeding = "9" },
		"cob":            func(i *IsianPenawaran) { i.BusinessCode = "L99" },
		"status offer":   func(i *IsianPenawaran) { i.Status = "1" },
		"ceding separuh": func(i *IsianPenawaran) { i.CedingCoName = "" },
		"pemegang separuh": func(i *IsianPenawaran) {
			i.PolicyHolder = ""
		},
	} {
		isi := isianLengkap()
		ubah(&isi)
		if _, err := SusunPenawaran(FlagPolisPenawaran, isi); !errors.Is(err, ErrIsianPenawaranTidakSah) {
			t.Errorf("%s: galat = %v, mau ErrIsianPenawaranTidakSah", nama, err)
		}
	}
	// Radio yang tampil bergantung bendera: "1" sah untuk Input Premium.
	isi := isianLengkap()
	isi.Status = "1"
	if _, err := SusunPenawaran(FlagPolisPremium, isi); err != nil {
		t.Errorf("EmailTypePL 1 pada bendera premium: %v", err)
	}
}

func TestSusunPenawaranMenurunkanNamaDanMerapikan(t *testing.T) {
	isi := isianLengkap()
	isi.TypeCeding = " 4 "
	got, err := SusunPenawaran(FlagPolisPenawaran, isi)
	if err != nil {
		t.Fatal(err)
	}
	if got.TypeCeding != "4" || got.TypeCedingName != "XOL" || got.BusinessName != "INDIVIDUAL WHOLE LIFE" {
		t.Errorf("hasil = %+v", got)
	}
}

// AddHistorySuggest langkah 1 (Offer) dan 2 (Bind, termasuk cabang else).
func TestSuggestBaruMengikutiBendera(t *testing.T) {
	w := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	offer := SuggestBaru(FlagPolisPenawaran, "Pending", " UJI ", "UJI-AKUN", w)
	if offer.InitialSuggest != InitialSuggestOffer || offer.IsCedingConfirm != "Pending" ||
		offer.CommentSuggest != "UJI" || offer.PICSuggest != "UJI-AKUN" || !offer.DateSuggest.Equal(w) {
		t.Errorf("offer = %+v", offer)
	}
	for status, mau := range map[string]string{"1": "Accept", "2": "Reject", "3": "Decline", "": "Decline"} {
		b := SuggestBaru(FlagPolisPremium, status, "", "UJI-AKUN", w)
		if b.InitialSuggest != InitialSuggestBind || b.IsCedingConfirm != mau {
			t.Errorf("bind status %q = %+v, mau %s", status, b, mau)
		}
	}
}

// SetMaxTBCLife_Act: Confirmation Date + TBC HARI; nil bila salah satunya kosong.
func TestMaxTBCMenambahHari(t *testing.T) {
	k := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	tbc := 25
	got := MaxTBC(&k, &tbc)
	if got == nil || got.Format("2006-01-02") != "2026-10-24" {
		t.Errorf("MaxTBC = %v, mau 2026-10-24 (layar Pega: 29/09/2026 + 25 = 24/10/2026)", got)
	}
	if MaxTBC(nil, &tbc) != nil || MaxTBC(&k, nil) != nil {
		t.Error("Max TBC dihitung dari isian kosong")
	}
	isi := isianLengkap()
	isi.TanggalKonfirmasi, isi.TBC = &k, &tbc
	siap, err := SusunPenawaran(FlagPolisPenawaran, isi)
	if err != nil || siap.TanggalTBC == nil || !siap.TanggalTBC.Equal(*got) {
		t.Errorf("SusunPenawaran tidak menurunkan Max TBC: %+v, %v", siap.TanggalTBC, err)
	}
	minus := -1
	isi.TBC = &minus
	if _, err := SusunPenawaran(FlagPolisPenawaran, isi); !errors.Is(err, ErrIsianPenawaranTidakSah) {
		t.Errorf("TBC negatif: %v", err)
	}
}
