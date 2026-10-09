package services_test

// Jalur baca (paket 1): identitas wajib, induk wajib ada, dan grid anak
// dibaca dengan KEDUA kunci induk dari baris induk yang sungguh ada - persis
// parameter RD Pega (PARITAS §3–§6).

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

var (
	pelaku   = inti.Pelaku{AkunID: "UJI-PELAKU"}
	noPelaku = inti.Pelaku{}
)

// containsAll - seluruh potongan ada di s.
func containsAll(s string, potongan ...string) bool {
	for _, p := range potongan {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}

func layananUji(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, g.Transaksi, nil)
}

func gudangBerisi() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Tahun["UJI-T1"] = models.TahunTreaty{ID: "UJI-T1", TreatyYear: "2026"}
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "UJI-T1", ReinsTypeID: "10196"}
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	// Baris yang salinan tahunnya menyimpang TIDAK tampil - saringan dua kunci RD.
	g.Reinsurer["UJI-R2"] = models.Reinsurer{ID: "UJI-R2", TreatyYearID: "UJI-TLAIN", TreatyContractID: "UJI-K1"}
	g.Security["UJI-S1"] = models.SecurityReinsurer{ID: "UJI-S1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		TreatyReinsurerID: "UJI-R1"}
	return g
}

func TestBacaTanpaIdentitasDitolak(t *testing.T) {
	l := layananUji(gudangBerisi())
	ctx := context.Background()
	kosong := inti.Pelaku{}
	cek := func(nama string, err error) {
		t.Helper()
		if !errors.Is(err, inti.ErrTanpaIdentitas) {
			t.Errorf("%s tanpa identitas: galat %v, mau ErrTanpaIdentitas", nama, err)
		}
	}
	_, err := l.DaftarTahun(ctx, kosong)
	cek("DaftarTahun", err)
	_, err = l.DaftarKontrak(ctx, kosong, "UJI-T1")
	cek("DaftarKontrak", err)
	_, err = l.DaftarReinsurer(ctx, kosong, "UJI-K1")
	cek("DaftarReinsurer", err)
	_, err = l.DaftarSecurity(ctx, kosong, "UJI-R1")
	cek("DaftarSecurity", err)
	_, err = l.DaftarBusiness(ctx, kosong, "UJI-K1")
	cek("DaftarBusiness", err)
	_, err = l.JenisReasuransi(ctx, kosong)
	cek("JenisReasuransi", err)
	_, err = l.DaftarRate(ctx, kosong, "1")
	cek("DaftarRate", err)
}

func TestIndukTidakAdaAdalah404BukanDaftarKosong(t *testing.T) {
	l := layananUji(gudangBerisi())
	ctx := context.Background()
	if _, err := l.DaftarKontrak(ctx, pelaku, "UJI-TIDAK-ADA"); !errors.Is(err, services.ErrTahunTidakAda) {
		t.Errorf("kontrak tahun tak ada: %v", err)
	}
	if _, err := l.DaftarReinsurer(ctx, pelaku, "UJI-TIDAK-ADA"); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("reinsurer kontrak tak ada: %v", err)
	}
	if _, err := l.DaftarBusiness(ctx, pelaku, "UJI-TIDAK-ADA"); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("business kontrak tak ada: %v", err)
	}
	if _, err := l.DaftarSecurity(ctx, pelaku, "UJI-TIDAK-ADA"); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("security reinsurer tak ada: %v", err)
	}
}

func TestGridAnakMemakaiKunciIndukDariBarisInduk(t *testing.T) {
	g := gudangBerisi()
	l := layananUji(g)
	ctx := context.Background()
	jr, err := l.DaftarReinsurer(ctx, pelaku, "UJI-K1")
	if err != nil {
		t.Fatal(err)
	}
	if len(jr.Daftar) != 1 || jr.Daftar[0].ID != "UJI-R1" {
		t.Errorf("reinsurer: %+v, mau hanya UJI-R1 (UJI-R2 bersalinan tahun lain)", jr.Daftar)
	}
	if jr.Kontrak.ID != "UJI-K1" {
		t.Errorf("kontrak kepala panel: %q", jr.Kontrak.ID)
	}
	js, err := l.DaftarSecurity(ctx, pelaku, "UJI-R1")
	if err != nil {
		t.Fatal(err)
	}
	if len(js.Daftar) != 1 || js.Induk.ID != "UJI-R1" {
		t.Errorf("security: %+v induk %q", js.Daftar, js.Induk.ID)
	}
	if _, err := l.DaftarBusiness(ctx, pelaku, "UJI-K1"); err != nil {
		t.Fatal(err)
	}
	mau := []string{"reinsurer|UJI-T1|UJI-K1", "security|UJI-T1|UJI-K1|UJI-R1", "business|UJI-T1|UJI-K1"}
	if !reflect.DeepEqual(g.Panggilan, mau) {
		t.Errorf("kunci induk ke gudang: %q, mau %q", g.Panggilan, mau)
	}
}

func TestDaftarKosongBukanNil(t *testing.T) {
	g := tiruan.Baru()
	g.Tahun["UJI-T1"] = models.TahunTreaty{ID: "UJI-T1"}
	j, err := layananUji(g).DaftarKontrak(context.Background(), pelaku, "UJI-T1")
	if err != nil {
		t.Fatal(err)
	}
	if j.Daftar == nil {
		t.Error("daftar kontrak kosong harus [] (JSON), bukan nil")
	}
}

func TestMasterJenisKosongAdalahKeadaanServer(t *testing.T) {
	_, err := layananUji(tiruan.Baru()).JenisReasuransi(context.Background(), pelaku)
	if !errors.Is(err, services.ErrMasterTidakTerbaca) {
		t.Fatalf("master kosong: %v, mau ErrMasterTidakTerbaca", err)
	}
	if got := services.Pesan(err); got == "" || !strings.Contains(got, "REINSURANCETYPE") {
		t.Errorf("pesan harus menyebut objeknya: %q", got)
	}
}

// K1 keputusan work owner 01-10-2026 (OQ-MCRL-13): kedua view rate dibaca. Autocomplete = `Contains`
// atas USEDBY urut `ID ASC`; Rate List satu IDUSEDBY urut `ID DESC, RATE ASC`, RATE teks apa adanya.
func TestRateDibacaDariViewRate(t *testing.T) {
	g := tiruan.Baru()
	g.RingkasanRate = []models.RingkasanRate{{ID: "UJI-3", UsedBy: "UJI RATE TIGA"}, {ID: "UJI-1", UsedBy: "uji rate satu"},
		{ID: "UJI-2", UsedBy: "LAIN"}}
	g.Rate["UJI-1"] = []models.BarisRate{{ID: "UJI-A", Rate: "0,5"}, {ID: "UJI-B", Rate: "1.25"}, {ID: "UJI-B", Rate: "0,75"}}
	l := layananUji(g)
	ctx := context.Background()
	d, err := l.CariRingkasanRate(ctx, pelaku, " rate ")
	if err != nil || len(d) != 2 || d[0].ID != "UJI-1" || d[1].ID != "UJI-3" {
		t.Errorf("CariRingkasanRate: %v %+v", err, d)
	}
	r, err := l.DaftarRate(ctx, pelaku, "UJI-1")
	urut := []string{}
	for _, b := range r.Daftar {
		urut = append(urut, b.ID+"="+b.Rate)
	}
	if err != nil || strings.Join(urut, " ") != "UJI-B=0,75 UJI-B=1.25 UJI-A=0,5" || r.Total != 3 || r.Terpotong {
		t.Errorf("DaftarRate: %v %v %+v", err, urut, r)
	}
	// IDUSEDBY tanpa baris = daftar kosong 200, bukan galat.
	if r, err := l.DaftarRate(ctx, pelaku, "UJI-9"); err != nil || r.Daftar == nil || len(r.Daftar) != 0 {
		t.Errorf("rate kosong: %v %+v", err, r)
	}
}

// `BrowseRateLife_RD` b729 `pyMaxRecords` 500: baris ke-501 tidak tampil, dan itu DINYATAKAN.
func TestRateDipotong500Dinyatakan(t *testing.T) {
	g := tiruan.Baru()
	for i := 0; i < 501; i++ {
		g.Rate["UJI-1"] = append(g.Rate["UJI-1"], models.BarisRate{ID: fmt.Sprintf("UJI-%04d", i), Rate: "1"})
	}
	r, err := layananUji(g).DaftarRate(context.Background(), pelaku, "UJI-1")
	if err != nil || len(r.Daftar) != 500 || r.Total != 500 || !r.Terpotong || r.Daftar[0].ID != "UJI-0500" {
		t.Errorf("501 baris: %v len %d terpotong %v", err, len(r.Daftar), r.Terpotong)
	}
}

// View rate tak terbaca = ErrMasterTidakTerbaca yang menyebut objeknya (503), bukan daftar kosong.
func TestRateTakTerbacaMenyebutView(t *testing.T) {
	g := tiruan.Baru()
	g.GalatMaster = errors.New("ORA-00942")
	l := layananUji(g)
	if _, err := l.DaftarRate(context.Background(), pelaku, "UJI-1"); !errors.Is(err, services.ErrMasterTidakTerbaca) ||
		!strings.Contains(services.Pesan(err), "M_RATE_LIFE") {
		t.Errorf("DaftarRate: %v", err)
	}
	if _, err := l.CariRingkasanRate(context.Background(), pelaku, ""); !errors.Is(err, services.ErrMasterTidakTerbaca) ||
		!strings.Contains(services.Pesan(err), "M_RATE_LIFE_SUMMARY") {
		t.Errorf("CariRingkasanRate: %v", err)
	}
}

func TestRateTanpaIdUsedByDitolak(t *testing.T) {
	_, err := layananUji(tiruan.Baru()).DaftarRate(context.Background(), pelaku, "  ")
	if !errors.Is(err, services.ErrParameterWajib) {
		t.Fatalf("idusedby kosong: %v", err)
	}
}

// Pembersih galat di SATU tempat (tiket 10): tag dibuang BERSAMA isinya.
func TestPesanMembuangMarkupDanAwalanLapisan(t *testing.T) {
	err := errors.New(`services: <span style="color:red">Data gagal</span> disimpan`)
	if got := services.Pesan(err); got != "Data gagal disimpan" {
		t.Errorf("Pesan = %q", got)
	}
}
