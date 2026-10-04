package models

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func tglTeks(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func kontrakWajar() KontrakTreaty {
	return KontrakTreaty{IDTreatyYear: "1000001", ReinsTypeID: "10003", ReinsTypeName: "UJI QS",
		TreatyStartDate: tglTeks("2026-01-01"), TreatyEndDate: tglTeks("2027-01-01")}
}

func TestPeriksaKontrakTreaty(t *testing.T) {
	if err := PeriksaKontrakTreaty(kontrakWajar(), "2026"); err != nil {
		t.Fatalf("kontrak wajar ditolak: %v", err)
	}
	kasus := []struct {
		ubah func(*KontrakTreaty)
		mau  error
		teks []string
	}{
		{func(k *KontrakTreaty) { k.ReinsTypeID = " " }, ErrKontrakJenisReasuransiKosong, []string{"ReinsTypeID"}},
		{func(k *KontrakTreaty) { k.TreatyStartDate = time.Time{} }, ErrKontrakMulaiKosong, []string{"TreatyStartDate"}},
		{func(k *KontrakTreaty) { k.TreatyEndDate = time.Time{} }, ErrKontrakAkhirKosong, []string{"TreatyEndDate"}},
		// AC 9 sisi kontrak: pesan menyebut kedua medan.
		{func(k *KontrakTreaty) { k.TreatyEndDate = tglTeks("2025-12-31") }, ErrPeriodeTerbalik,
			[]string{"TreatyStartDate", "TreatyEndDate", "2025-12-31"}},
		// SetTanggalTreatyContract b847: tahun tanggal mulai = tahun treaty.
		{func(k *KontrakTreaty) {
			k.TreatyStartDate, k.TreatyEndDate = tglTeks("2027-01-01"), tglTeks("2028-01-01")
		}, ErrKontrakTahunMulaiBeda, []string{"TreatyStartDate", "2027", "2026"}},
	}
	for _, k := range kasus {
		m := kontrakWajar()
		k.ubah(&m)
		err := PeriksaKontrakTreaty(m, "2026")
		if !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
			continue
		}
		for _, s := range k.teks {
			if !strings.Contains(err.Error(), s) {
				t.Errorf("pesan %q tanpa %q", err, s)
			}
		}
	}
}

// SetTanggalTreatyContract langkah 3-6 DITIRU APA ADANYA, termasuk anomalinya:
// 366 hari hanya bila (a) tahun mulai = tahun treaty, (b) dua digit awal
// cap waktu Pega (abad) di 1..31, (c) digit puluhan bulan di 1..2 - yaitu
// Oktober-Desember, sebab cap waktunya `yyyyMMdd...` - dan (d) tahun treaty
// kabisat. Selain itu 365 hari.
// OQ-TCO-10 [keputusan work owner 29-09-2026]: akhir = mulai + 1 TAHUN
// KALENDER, semantik `ADD_MONTHS(mulai, 12)` - menggantikan anomali 366 hari.
func TestAkhirKontrakBawaanTCO(t *testing.T) {
	kasus := []struct {
		mulai, mau, sebab string
	}{
		{"2026-01-01", "2027-01-01", "tahun biasa"},
		{"2026-03-15", "2027-03-15", "tengah bulan"},
		{"2024-01-01", "2025-01-01", "tahun kabisat, Januari (dulu 2024-12-31)"},
		{"2024-10-01", "2025-10-01", "tahun kabisat, Oktober (dulu 366 hari)"},
		{"2023-03-01", "2024-03-01", "melewati 29 Februari"},
		{"2028-02-29", "2029-02-28", "29 Februari dijepit ke akhir bulan"},
		{"2027-02-28", "2028-02-29", "akhir bulan tetap akhir bulan (ADD_MONTHS)"},
		{"2026-06-30", "2027-06-30", "akhir bulan 30 hari"},
		{"2099-12-31", "2100-12-31", "akhir tahun"},
	}
	for _, k := range kasus {
		if dapat := AkhirKontrakBawaanTCO(tglTeks(k.mulai)).Format("2006-01-02"); dapat != k.mau {
			t.Errorf("%s: AkhirKontrakBawaanTCO(%s) = %s, mau %s", k.sebab, k.mulai, dapat, k.mau)
		}
	}
}
