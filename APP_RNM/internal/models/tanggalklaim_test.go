package models

import (
	"encoding/json"
	"testing"
)

// TestTanggalKlaimHanyaDiOutstanding - gerbang baca-saja keempat isian
// `EditDateClaimLife_Section` (b1000, b1313, b1550, b1788), disilangkan dengan
// tahap yang membuka grid peserta (`ClaimLifeDetailGCNM`, tempat tombol
// `Edit Date` b14115 berdiri).
//
// Input Register memegang ReasLifeAdmin tetapi tidak membuka grid peserta;
// Medical Check dan Claim Analis membuka grid tetapi tidak dipegang Admin.
func TestTanggalKlaimHanyaDiOutstanding(t *testing.T) {
	for _, u := range []struct {
		tahap Tahap
		mau   bool
	}{
		{TahapOutstanding, true},
		{TahapInputRegister, false},
		{TahapMedicalCheck, false},
		{TahapClaimAnalis, false},
		{TahapTidakDikenal, false},
	} {
		if got := TahapBolehUbahTanggalKlaim(u.tahap); got != u.mau {
			t.Errorf("%s: %v, mau %v", u.tahap, got, u.mau)
		}
	}
}

// TestTigaTanggalKlaimMenyeberang - dialog Edit Date menampilkan tanggal yang
// SEDANG berlaku (`UpdateDateClaimLife_Act` b252-b384 membacanya dari peserta).
func TestTigaTanggalKlaimMenyeberang(t *testing.T) {
	b, err := json.Marshal(Peserta{ID: "P-1", TanggalTerimaKlaim: "2026-03-02",
		TanggalDokumenLengkap: "2026-03-03", TanggalKonfirmasi: "2026-03-04"})
	if err != nil {
		t.Fatal(err)
	}
	var isi map[string]any
	if err := json.Unmarshal(b, &isi); err != nil {
		t.Fatal(err)
	}
	for kunci, mau := range map[string]string{
		"tanggalTerimaKlaim": "2026-03-02", "tanggalDokumenLengkap": "2026-03-03",
		"tanggalKonfirmasi": "2026-03-04",
	} {
		if isi[kunci] != mau {
			t.Errorf("%s = %v, mau %s", kunci, isi[kunci], mau)
		}
	}
}
