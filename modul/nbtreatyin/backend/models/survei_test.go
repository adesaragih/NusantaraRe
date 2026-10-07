package models

// Keputusan work owner 06-10-2026 (membatalkan K7): daftar survei historis
// `PolicyTreatyIn.QuotationData.SurveyReportList` disimpan di T_POLIS_SURVEY - empat medan sel grid
// `Section/HistoricalSurveyReportDtl` (DateofSurvey, SurveyedBy, LossPrevention, Remarks); medan lain hilang.

import "testing"

func TestProyeksiKatalogMenyimpanDaftarSurvei(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarSurvei, []Baris{
		{"DateofSurvey": "2026-09-01", "SurveyedBy": "UJI-SURVEYOR", "LossPrevention": "12.5", "Remarks": "UJI-R", "pxObjClass": "UJI"},
		{"DateofSurvey": "2026-09-15", "SurveyedBy": "UJI-SURVEYOR-2"},
	})
	a := ProyeksiKatalog(h).AmbilDaftar(DaftarSurvei)
	if len(a) != 2 {
		t.Fatalf("baris survei %d, harap 2: %+v", len(a), a)
	}
	for k, harap := range map[string]string{"DateofSurvey": "2026-09-01", "SurveyedBy": "UJI-SURVEYOR", "LossPrevention": "12.5", "Remarks": "UJI-R"} {
		if a[0][k] != harap {
			t.Errorf("baris 1 %s = %q, harap %q", k, a[0][k], harap)
		}
	}
	if _, ada := a[0]["pxObjClass"]; ada {
		t.Error("medan di luar sel grid tidak punya kolom")
	}
	if DaftarSurvei != HalamanPolis+".QuotationData.SurveyReportList" {
		t.Errorf("DaftarSurvei = %q", DaftarSurvei)
	}
}

// Popup survei hanya dapat disunting di layar ADMIN (`Section/HistoricalSurveyReportDtl`; layar atasan
// `HistoricalSurveyReportDtlUW` hanya-baca) dan hanya bila tombol Survey Report tampil DAN aktif
// (pyVisible ProportionalType != NonProportional; pyDisabledWhen IsSurveyReport 'No' / ”). Selain itu
// kiriman daftar diabaikan - nilai server dipakai (pola AC 49-51). Hanya empat medan sel yang diterima.
func TestGabungMasukanLayarDaftarSurvei(t *testing.T) {
	kiriman := func(survei string) *Halaman {
		m := HalamanBaru()
		m.Setel(HalamanPolis+".QuotationData.IsSurveyReport", survei)
		m.SetelDaftar(DaftarSurvei, []Baris{
			{"DateofSurvey": "2026-09-01", "SurveyedBy": "UJI-S", "LossPrevention": "1", "Remarks": "UJI-R", "Lain": "x"},
			{},
		})
		return m
	}
	server := func() *Halaman {
		h := HalamanBaru()
		h.SetelDaftar(DaftarSurvei, []Baris{{"DateofSurvey": "2026-01-01"}})
		return h
	}

	h := server()
	if _, err := GabungMasukanLayar(h, kiriman("Yes"), PosisiAdmin, nil); err != nil {
		t.Fatal(err)
	}
	b := h.AmbilDaftar(DaftarSurvei)
	if len(b) != 1 || b[0]["SurveyedBy"] != "UJI-S" || b[0]["Remarks"] != "UJI-R" {
		t.Fatalf("admin + Yes: daftar = %+v, harap satu baris kiriman (baris kosong dibuang)", b)
	}
	if _, ada := b[0]["Lain"]; ada {
		t.Error("medan di luar sel grid tidak diterima")
	}

	for _, c := range []struct {
		nama, survei, posisi, jenis string
	}{
		{"Survey Report No", "No", PosisiAdmin, ""},
		{"Survey Report kosong", "", PosisiAdmin, ""},
		{"layar atasan", "Yes", PosisiSecHead, ""},
		{"NonProportional", "Yes", PosisiAdmin, JenisNonProporsional},
	} {
		h := server()
		h.Setel(HalamanPolis+".QuotationData.IsSurveyReport", c.survei)
		h.Setel(HalamanQuotation+".ProportionalType", c.jenis)
		if _, err := GabungMasukanLayar(h, kiriman(c.survei), c.posisi, nil); err != nil {
			t.Fatal(err)
		}
		if b := h.AmbilDaftar(DaftarSurvei); len(b) != 1 || b[0]["DateofSurvey"] != "2026-01-01" {
			t.Errorf("%s: kiriman harus diabaikan, daftar = %+v", c.nama, b)
		}
	}
}
