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
func TestAkhirKontrakBawaanTCO(t *testing.T) {
	kasus := []struct {
		mulai, tahun, mau string
	}{
		{"2026-01-01", "2026", "2027-01-01"}, // 365
		{"2024-01-01", "2024", "2024-12-31"}, // kabisat, Januari: TETAP 365 (anomali OQ-TCO-10)
		{"2024-10-01", "2024", "2025-10-02"}, // kabisat, Oktober: 366
		{"2024-12-01", "2024", "2025-12-02"}, // kabisat, Desember: 366
		{"2024-10-01", "2025", "2025-10-01"}, // tahun mulai beda: 365
		{"2100-10-01", "2100", "2101-10-01"}, // 2100 bukan kabisat: 365
		{"2000-11-01", "2000", "2001-11-02"}, // 2000 kabisat: 366
		{"2024-10-01", "dua ribu", "2025-10-01"},
	}
	for _, k := range kasus {
		if dapat := AkhirKontrakBawaanTCO(tglTeks(k.mulai), k.tahun).Format("2006-01-02"); dapat != k.mau {
			t.Errorf("AkhirKontrakBawaanTCO(%s, %q) = %s, mau %s", k.mulai, k.tahun, dapat, k.mau)
		}
	}
}
