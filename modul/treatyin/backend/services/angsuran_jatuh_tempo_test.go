package services

import (
	"reflect"
	"testing"
	"time"
)

// Due Date Installment dari Commencement, dibagi rata sampai Termination -
// keputusan pemakai 9 Oktober 2026 (menyimpang dari `@CurrentDate` ekspor).
func TestJadwalJatuhTempoDibagiRataDariCommencement(t *testing.T) {
	kini := time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)
	kasus := []struct {
		nama, mulai, akhir string
		n                  int
		ingin              []string
	}{
		{"setahun 4x = per kuartal", "20260701", "20270630", 4,
			[]string{"20260701", "20261001", "20270101", "20270401"}},
		{"setahun 12x = bulanan", "20260101", "20261231", 12,
			[]string{"20260101", "20260201", "20260301", "20260401", "20260501", "20260601",
				"20260701", "20260801", "20260901", "20261001", "20261101", "20261201"}},
		{"setahun 5x dibulatkan", "20260101", "20261231", 5,
			[]string{"20260101", "20260301", "20260601", "20260801", "20261101"}},
		{"hari dipotong ke akhir bulan", "20260131", "20270130", 4,
			[]string{"20260131", "20260430", "20260731", "20261031"}},
		{"bentuk kabel DD-MM-YYYY", "01-07-2026", "30-06-2027", 2,
			[]string{"20260701", "20270101"}},
		{"Termination kosong = semua Commencement", "20260701", "", 3,
			[]string{"20260701", "20260701", "20260701"}},
		{"Commencement kosong = hari ini (ekspor)", "", "20270630", 2,
			[]string{"20261009", "20261009"}},
	}
	for _, k := range kasus {
		if got := jadwalJatuhTempo(k.mulai, k.akhir, k.n, kini); !reflect.DeepEqual(got, k.ingin) {
			t.Errorf("%s: %v, ingin %v", k.nama, got, k.ingin)
		}
	}
}

// Lewat aksi `nilai` (TreatyInSetValueInstallment) - baris memakai jadwal itu.
func TestNilaiAngsuranMemakaiCommencement(t *testing.T) {
	h, err := HitungAngsuran(MasukanAngsuran{
		Aksi: AksiAngsuranNilai, InstallmentNo: "4",
		NetPremium:   []NilaiMataUang{{Currency: "IDR", Value: "1000"}},
		Commencement: "20260701", Termination: "20270630",
	}, time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	got := kolomAngsuran(h.Angsuran[0], func(r BarisAngsuran) string { return r.DueDate })
	if !reflect.DeepEqual(got, []string{"20260701", "20261001", "20270101", "20270401"}) {
		t.Errorf("DueDate = %v", got)
	}
}
