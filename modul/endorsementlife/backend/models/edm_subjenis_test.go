package models

import (
	"errors"
	"testing"
)

// TestPeriksaSubJenisEdm - kombinasi EDM Type × sub-jenis (OQ-EDM-005, keputusan work owner 07-10-2026).
func TestPeriksaSubJenisEdm(t *testing.T) {
	sah := []struct{ tipe, ubah, batal string }{
		{EdmTypeBatal, "", EdmTypeBatalInternal},
		{EdmTypeBatal, "", "Batal Eksternal"},
	}
	for _, o := range OpsiEdmTypePerubahanData {
		sah = append(sah, struct{ tipe, ubah, batal string }{EdmTypePerubahanData, o, ""})
	}
	for _, k := range sah {
		if err := PeriksaSubJenisEdm(k.tipe, k.ubah, k.batal); err != nil {
			t.Errorf("%+v ditolak: %v", k, err)
		}
	}
	tolak := []struct{ tipe, ubah, batal string }{
		{EdmTypePerubahanData, "", ""},                                            // wajib
		{EdmTypePerubahanData, "UJI-BUKAN-OPSI", ""},                              // di luar daftar
		{EdmTypePerubahanData, "perubahan data (dob)", ""},                        // persis, bukan tanpa huruf besar
		{EdmTypePerubahanData, OpsiEdmTypePerubahanData[0], EdmTypeBatalInternal}, // Batal harus kosong
		{EdmTypeBatal, "", ""},                                                    // wajib
		{EdmTypeBatal, "", "Batal"},                                               // di luar daftar
		{EdmTypeBatal, OpsiEdmTypePerubahanData[3], EdmTypeBatalInternal},         // Perubahan Data harus kosong
	}
	for _, k := range tolak {
		if err := PeriksaSubJenisEdm(k.tipe, k.ubah, k.batal); !errors.Is(err, ErrSubJenisEdmTidakSah) {
			t.Errorf("%+v diterima: %v", k, err)
		}
	}
	// Daftar opsi = layar Pega (diberikan work owner); bawaan Batal = yang pertama.
	if len(OpsiEdmTypePerubahanData) != 6 || OpsiEdmTypePerubahanData[5] != "Perubahan data (UP & Premi)" {
		t.Errorf("opsi Perubahan Data = %q", OpsiEdmTypePerubahanData)
	}
	if len(OpsiEdmTypeBatal) != 2 || OpsiEdmTypeBatal[0] != "Batal Internal" || OpsiEdmTypeBatal[1] != "Batal Eksternal" {
		t.Errorf("opsi Batal = %q", OpsiEdmTypeBatal)
	}
	for _, o := range append(append([]string{}, OpsiEdmTypePerubahanData...), OpsiEdmTypeBatal...) {
		if len(o) > 100 {
			t.Errorf("opsi %q melebihi VARCHAR2(100) migrasi 483", o)
		}
	}
}
