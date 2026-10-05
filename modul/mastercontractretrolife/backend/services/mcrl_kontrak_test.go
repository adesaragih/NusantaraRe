package services_test

// Kontrak dan batas proteksi (paket 3, tiket 02/04): `SaveTreatyLimit_Act`
// tanpa procedure - wajib-isi langkah 3 b609 VERBATIM b313, tanggal dari tahun
// (R5), selisih dihitung lalu ditulis (K5), jenis dari master life (tiket 04),
// salinan jenis ke anak (K4), gerbang tahun TIDAK ditegakkan (K7).

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func gudangKontrak() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Tahun["UJI-T1"] = models.TahunTreaty{ID: "UJI-T1", TreatyYear: "2026",
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)}
	g.Jenis = []models.JenisReasuransi{{ID: "10197", Note: "2ND QS"}, {ID: "10196", Note: "QS"}}
	return g
}

func kontrakLengkap() services.KontrakMasuk {
	return services.KontrakMasuk{ReinsTypeID: "10196", BIDR: "0", IDR: "1500000000.10", BUSD: "0", USD: "100000.5"}
}

func TestKontrakWajibEnamMedanVerbatim(t *testing.T) {
	kosong := map[string]func(*services.KontrakMasuk){
		"REINS TYPE":          func(m *services.KontrakMasuk) { m.ReinsTypeID = "" },
		"MINIMUM LIMIT (IDR)": func(m *services.KontrakMasuk) { m.BIDR = "" },
		"MAXIMUM LIMIT (IDR)": func(m *services.KontrakMasuk) { m.IDR = " " },
		"MINIMUM LIMIT (USD)": func(m *services.KontrakMasuk) { m.BUSD = "" },
	}
	for nama, ubah := range kosong {
		g := gudangKontrak()
		m := kontrakLengkap()
		ubah(&m)
		_, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
		if !errors.Is(err, services.ErrWajibIsi) || services.Pesan(err) != "All value cannot be empty." {
			t.Errorf("%s kosong: %v", nama, err)
		}
		if len(g.Kontrak) != 0 {
			t.Errorf("%s kosong: baris lahir", nama)
		}
	}
	// Tanggal kontrak berasal dari tahun (R5): tahun tanpa tanggal = wajib-isi gagal.
	g := gudangKontrak()
	th := g.Tahun["UJI-T1"]
	th.StartDate = time.Time{}
	g.Tahun["UJI-T1"] = th
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", kontrakLengkap()); !errors.Is(err, services.ErrWajibIsi) {
		t.Errorf("TREATY START kosong (tahun tanpa tanggal): %v", err)
	}
}

// MAXIMUM LIMIT (USD) WAJIB sejak 04-10-2026 (keputusan work owner) - dulu
// boleh kosong (NULL).
func TestKontrakUSDKosongDitolak(t *testing.T) {
	g := gudangKontrak()
	m := kontrakLengkap()
	m.USD = " "
	_, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
	if !errors.Is(err, services.ErrWajibIsi) {
		t.Errorf("USD kosong: galat %v, mau ErrWajibIsi", err)
	}
}

func TestKontrakSelisihUSDDanNilaiIdentik(t *testing.T) {
	g := gudangKontrak()
	m := kontrakLengkap()
	m.BUSD, m.USD = "0.000000001", "100000.123456789012345678"
	k, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
	if err != nil {
		t.Fatal(err)
	}
	if k.USD.Text('f') != "100000.123456789012345678" {
		t.Errorf("USD tidak identik (ADR-0003): %s", k.USD.Text('f'))
	}
	if k.USDSelisih.Text('f') != "100000.123456788012345678" {
		t.Errorf("USD_SELISIH = %s", k.USDSelisih.Text('f'))
	}
}

func TestKontrakBatasBawahMelebihiAtasDitolakMenyebutMataUang(t *testing.T) {
	g := gudangKontrak()
	m := kontrakLengkap()
	m.BUSD, m.USD = "200", "100"
	_, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !containsAll(services.Pesan(err), "USD") {
		t.Errorf("B_USD > USD: %v (OQ-MCRL-09)", err)
	}
	// Layer menaik: batas atas satu layer = batas bawah layer berikut - sah.
	m = kontrakLengkap()
	m.BIDR, m.IDR = "1500000000.10", "1500000000.10"
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m); err != nil {
		t.Errorf("B_IDR = IDR harus sah: %v", err)
	}
}

func TestKontrakTanggalDanJenisDariMasterBukanKlien(t *testing.T) {
	g := gudangKontrak()
	m := kontrakLengkap()
	m.ReinsTypeName = "UJI-NAMA-KARANGAN"
	k, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != "1000044" || k.IDTreatyYear != "UJI-T1" || k.UserID != "UJI-PELAKU" {
		t.Errorf("identitas/induk/pelaku: %+v", k)
	}
	if k.ReinsTypeName != "QS" {
		t.Errorf("REINSTYPENAME %q, mau .Note master (TreatyLimit_TypeProtect 4.1)", k.ReinsTypeName)
	}
	if !k.TreatyStartDate.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) ||
		!k.TreatyEndDate.Equal(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("tanggal kontrak bukan salinan tahun (R5): %v %v", k.TreatyStartDate, k.TreatyEndDate)
	}
	m.ReinsTypeID = "10004"
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("jenis di luar master life: %v", err)
	}
}

// ⛔ K7: gerbang tahun MATI di Pega - ikut XML sampai OQ-MCRL-01 dijawab.
func TestKontrakTahunBedaTetapTersimpan(t *testing.T) {
	g := gudangKontrak()
	th := g.Tahun["UJI-T1"]
	th.TreatyYear = "2027" // TRANSACTION YEAR berbeda dari tahun START DATE 2026
	g.Tahun["UJI-T1"] = th
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", kontrakLengkap()); err != nil {
		t.Errorf("gerbang tahun ditegakkan padahal mati di Pega (K7): %v", err)
	}
}

func TestKontrakUbahMenyalinJenisKeAnak(t *testing.T) {
	g := gudangKontrak()
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "UJI-T1", ReinsTypeID: "10196", ReinsTypeName: "QS"}
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		ReinsTypeID: "10196", ReinsTypeName: "QS"}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		ReinsTypeID: "10196", ReinsTypeName: "QS"}
	m := kontrakLengkap()
	m.ID, m.ReinsTypeID = "UJI-K1", "10197"
	k, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m)
	if err != nil {
		t.Fatal(err)
	}
	if k.ReinsTypeName != "2ND QS" || len(g.Kontrak) != 1 {
		t.Errorf("upsert kontrak: %+v (%d baris)", k, len(g.Kontrak))
	}
	for _, r := range g.Reinsurer {
		if r.ReinsTypeID != "10197" || r.ReinsTypeName != "2ND QS" {
			t.Errorf("salinan jenis reinsurer tidak ikut kontrak (K4): %+v", r)
		}
	}
	for _, b := range g.Business {
		if b.ReinsTypeID != "10197" || b.ReinsTypeName != "2ND QS" {
			t.Errorf("salinan jenis business tidak ikut kontrak (K4): %+v", b)
		}
	}
}

func TestKontrakUbahTahunLainAtauTidakAda(t *testing.T) {
	g := gudangKontrak()
	g.Tahun["UJI-T2"] = models.TahunTreaty{ID: "UJI-T2"}
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T2"}
	m := kontrakLengkap()
	m.ID = "UJI-K2"
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-T1", m); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak tahun lain lewat tahun UJI-T1: %v", err)
	}
	if _, err := layananUji(g).SimpanKontrak(context.Background(), pelaku, "UJI-TAK-ADA", kontrakLengkap()); !errors.Is(err, services.ErrTahunTidakAda) {
		t.Errorf("tahun tak ada: %v", err)
	}
}
