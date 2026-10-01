package services_test

// Tahun treaty (paket 2, tiket 01): `SaveTreatyYearLife_Act` tanpa procedure
// (keputusan o) - wajib-isi hidup langkah 3 b609 berpesan VERBATIM b313,
// ID baru dari sequence (K3), USERID pelaku (K6), salinan anak (K4).

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func tahunLengkap() services.TahunMasuk {
	return services.TahunMasuk{TreatyYear: "2026", UnderwritingYear: "2026", StartDate: "2026-01-01", EndDate: "2026-12-31"}
}

func TestTahunBaruWajibEmpatMedanBerpesanVerbatim(t *testing.T) {
	kosongkan := map[string]func(*services.TahunMasuk){
		"UNDERWRITING YEAR": func(m *services.TahunMasuk) { m.UnderwritingYear = " " },
		"TRANSACTION YEAR":  func(m *services.TahunMasuk) { m.TreatyYear = "" },
		"START DATE":        func(m *services.TahunMasuk) { m.StartDate = "" },
		"END DATE":          func(m *services.TahunMasuk) { m.EndDate = "" },
	}
	for nama, ubah := range kosongkan {
		g := tiruan.Baru()
		m := tahunLengkap()
		ubah(&m)
		_, err := layananUji(g).SimpanTahun(context.Background(), pelaku, m, true)
		if !errors.Is(err, services.ErrWajibIsi) {
			t.Errorf("%s kosong: galat %v, mau ErrWajibIsi", nama, err)
			continue
		}
		if got := services.Pesan(err); got != "Value cannot be empty." {
			t.Errorf("%s kosong: pesan %q, mau VERBATIM SaveTreatyYearLife_Act b313", nama, got)
		}
		if len(g.Tahun) != 0 {
			t.Errorf("%s kosong: %d baris tersimpan, mau 0", nama, len(g.Tahun))
		}
	}
}

func TestTahunBaruIDDariSequenceDanPelakuDicatat(t *testing.T) {
	g := tiruan.Baru()
	hasil, err := layananUji(g).SimpanTahun(context.Background(), pelaku, tahunLengkap(), true)
	if err != nil {
		t.Fatal(err)
	}
	if hasil.ID != "1000044" {
		t.Errorf("ID baru %q, mau '1' || LPAD(44, 6, '0') dari TREATYYEAR_LIFE_SEQ", hasil.ID)
	}
	if hasil.UserID != "UJI-PELAKU" || hasil.TglUpdate.IsZero() {
		t.Errorf("USERID %q TGLUPDATE %v - K6", hasil.UserID, hasil.TglUpdate)
	}
	if !hasil.StartDate.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("STARTDATE %v", hasil.StartDate)
	}
}

func TestTahunBaruMenolakIDDariKlien(t *testing.T) {
	g := tiruan.Baru()
	m := tahunLengkap()
	m.ID = "UJI-KARANGAN"
	if _, err := layananUji(g).SimpanTahun(context.Background(), pelaku, m, true); !errors.Is(err, services.ErrIDDariKlien) {
		t.Errorf("ID klien pada baris baru: %v, mau ErrIDDariKlien (ADR-0006)", err)
	}
	if len(g.Tahun) != 0 {
		t.Error("baris lahir walau ID datang dari klien")
	}
}

func TestTahunUbahMemperbaruiSalinanAnak(t *testing.T) {
	g := tiruan.Baru()
	g.Tahun["1000044"] = models.TahunTreaty{ID: "1000044", TreatyYear: "2025", UnderwritingYear: "2025",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		UserID: "UJI-LAMA"}
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "1000044",
		TreatyStartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), TreatyEndDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "1000044", TreatyContractID: "UJI-K1", TreatyYear: "2025"}
	g.Business["UJI-B2"] = models.Business{ID: "UJI-B2", TreatyYearID: "UJI-TAHUN-LAIN", TreatyYear: "2025"}

	m := tahunLengkap()
	m.ID = "1000044"
	hasil, err := layananUji(g).SimpanTahun(context.Background(), pelaku, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if hasil.TreatyYear != "2026" || hasil.UserID != "UJI-PELAKU" || len(g.Tahun) != 1 {
		t.Errorf("upsert dikunci ID: %+v, %d baris", hasil, len(g.Tahun))
	}
	if b := g.Business["UJI-B1"]; b.TreatyYear != "2026" || b.UserID != "UJI-PELAKU" {
		t.Errorf("salinan TREATYYEAR business tidak ikut induk (K4): %+v", b)
	}
	if b := g.Business["UJI-B2"]; b.TreatyYear != "2025" {
		t.Errorf("business tahun LAIN ikut berubah: %+v", b)
	}
	k := g.Kontrak["UJI-K1"]
	if !k.TreatyStartDate.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) ||
		!k.TreatyEndDate.Equal(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("tanggal kontrak tidak ikut tahun induk (R5, OQ-MCRL-10): %v %v", k.TreatyStartDate, k.TreatyEndDate)
	}
}

func TestTahunUbahYangTidakAda404(t *testing.T) {
	m := tahunLengkap()
	m.ID = "UJI-TIDAK-ADA"
	if _, err := layananUji(tiruan.Baru()).SimpanTahun(context.Background(), pelaku, m, false); !errors.Is(err, services.ErrTahunTidakAda) {
		t.Errorf("ubah tahun tak ada: %v", err)
	}
}

func TestTahunTanggalTidakSahDitolakMenyebutMedan(t *testing.T) {
	m := tahunLengkap()
	m.EndDate = "31/12/2026"
	_, err := layananUji(tiruan.Baru()).SimpanTahun(context.Background(), pelaku, m, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("tanggal bukan YYYY-MM-DD: %v", err)
	}
	if got := services.Pesan(err); got == "" || !containsAll(got, "END DATE", "31/12/2026") {
		t.Errorf("pesan harus menyebut medan dan nilainya: %q", got)
	}
}

func TestTahunTanpaIdentitasDitolak(t *testing.T) {
	g := tiruan.Baru()
	if _, err := layananUji(g).SimpanTahun(context.Background(), noPelaku, tahunLengkap(), true); err == nil || len(g.Tahun) != 0 {
		t.Errorf("simpan tanpa identitas: %v, %d baris", err, len(g.Tahun))
	}
}
