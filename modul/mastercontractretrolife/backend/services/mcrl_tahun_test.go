package services_test

// Tahun treaty (paket 2, tiket 01): `SaveTreatyYearLife_Act` tanpa procedure
// (keputusan o) - wajib-isi hidup langkah 3 b609 berpesan VERBATIM b313,
// ID baru dari sequence (K3), USERID pelaku (K6), salinan anak (K4).

import (
	"context"
	"errors"
	"strings"
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

// TestTahunUbahTidakMenyentuhBusinessYatim - perbaikan 01-10-2026. Penghapus
// kontrak Pega (`DeleteTreatyLimit_SQL`) datar, jadi DEV menyimpan baris
// business yang TREATYCONTRACTID-nya sudah tidak ada (mis. kontrak yang
// dihapus di Pega). Baris itu tidak tampil di layar mana pun (daftar business
// dikunci TREATYYEARID + TREATYCONTRACTID, `BrowseTreatyBusiness_Life_RD`),
// maka salinan TREATYYEAR (K4) tidak boleh menimpanya - termasuk USERID dan
// TGLUPDATE-nya.
func TestTahunUbahTidakMenyentuhBusinessYatim(t *testing.T) {
	g := tiruan.Baru()
	g.Tahun["1000044"] = models.TahunTreaty{ID: "1000044", TreatyYear: "2025", UnderwritingYear: "2025",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		UserID: "UJI-LAMA"}
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "1000044",
		TreatyStartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), TreatyEndDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyYearID: "1000044", TreatyContractID: "UJI-K1"}
	g.Business["UJI-YATIM"] = models.Business{ID: "UJI-YATIM", TreatyYearID: "1000044",
		TreatyContractID: "UJI-K-DIHAPUS-PEGA", UserID: "UJI-LAMA"}

	m := tahunLengkap()
	m.ID = "1000044"
	if _, err := layananUji(g).SimpanTahun(context.Background(), pelaku, m, false); err != nil {
		t.Fatal(err)
	}
	if b := g.Business["UJI-B1"]; b.TreatyYear != "2026" {
		t.Errorf("business berinduk tidak ikut tahun: %+v", b)
	}
	if b := g.Business["UJI-YATIM"]; b.TreatyYear != "" || b.UserID != "UJI-LAMA" {
		t.Errorf("business yatim ikut disentuh salinan tahun: %+v", b)
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

func TestTerjemahRelasiHanyaGalatRelasi(t *testing.T) {
	if err := services.TerjemahRelasi(errors.New("ORA-01722: invalid number")); errors.Is(err, services.ErrRelasiDitolak) {
		t.Errorf("galat bukan relasi ikut diterjemahkan: %v", err)
	}
	asal := errors.New("ORA-02292: child record found")
	err := services.TerjemahRelasi(asal)
	if !errors.Is(err, services.ErrRelasiDitolak) || !errors.Is(err, asal) {
		t.Errorf("ORA-02292 tidak menjadi ErrRelasiDitolak yang membawa asalnya: %v", err)
	}
	if p := services.Pesan(err); strings.Contains(p, "ORA-") || p == "" {
		t.Errorf("pesan layar memuat teks Oracle atau kosong: %q", p)
	}
	if services.TerjemahRelasi(nil) != nil {
		t.Error("nil harus tetap nil")
	}
}

// TRANSACTION YEAR dan UNDERWRITING YEAR tidak boleh sama dengan tahun treaty lain
// (keputusan work owner 05-10-2026); Edit baris sendiri tetap boleh.
func TestTahunTreatyTidakBolehGanda(t *testing.T) {
	g := tiruan.Baru()
	l := layananUji(g)
	ada, err := l.SimpanTahun(context.Background(), pelaku, tahunLengkap(), true) // 2026 / 2026
	if err != nil {
		t.Fatal(err)
	}
	for nama, k := range map[string]struct{ ty, uy, pesan string }{
		"transaksi": {"2026", "2027", "TRANSACTION YEAR 2026 is already used (ID " + ada.ID + ")."},
		"uw":        {"2027", " 2026 ", "UNDERWRITING YEAR 2026 is already used (ID " + ada.ID + ")."},
	} {
		m := tahunLengkap()
		m.TreatyYear, m.UnderwritingYear = k.ty, k.uy
		_, err := l.SimpanTahun(context.Background(), pelaku, m, true)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), k.pesan) {
			t.Errorf("%s: %v, mau %q", nama, err, k.pesan)
		}
	}
	if len(g.Tahun) != 1 {
		t.Errorf("%d baris tersimpan, mau 1", len(g.Tahun))
	}
	// Edit baris sendiri dengan tahun yang sama - boleh.
	m := tahunLengkap()
	m.ID = ada.ID
	if _, err := l.SimpanTahun(context.Background(), pelaku, m, false); err != nil {
		t.Errorf("edit baris sendiri ditolak: %v", err)
	}
	// Tahun lain - boleh.
	m = tahunLengkap()
	m.TreatyYear, m.UnderwritingYear = "2027", "2027"
	// periode juga tidak boleh bertabrakan (06-10-2026)
	m.StartDate, m.EndDate = "2027-01-01", "2027-12-31"
	if _, err := l.SimpanTahun(context.Background(), pelaku, m, true); err != nil {
		t.Errorf("tahun berbeda ditolak: %v", err)
	}
}

// Periode treaty year tidak boleh bertabrakan (keputusan work owner 06-10-2026):
// sebagian, di dalam, sama persis ditolak; bersambung dan edit diri sendiri boleh.
func TestTahunTreatyPeriodeTidakBertabrakan(t *testing.T) {
	g := tiruan.Baru()
	l := layananUji(g)
	m := services.TahunMasuk{TreatyYear: "2025", UnderwritingYear: "2025", StartDate: "2025-01-01", EndDate: "2036-12-31"}
	ada, err := l.SimpanTahun(context.Background(), pelaku, m, true)
	if err != nil {
		t.Fatal(err)
	}
	for nama, k := range map[string][2]string{
		"sebagian":    {"2020-06-01", "2025-01-01"},
		"di dalam":    {"2026-01-01", "2026-12-31"},
		"sama persis": {"2025-01-01", "2036-12-31"},
		"mencakup":    {"2024-01-01", "2040-12-31"},
	} {
		b := services.TahunMasuk{TreatyYear: "2099", UnderwritingYear: "2099", StartDate: k[0], EndDate: k[1]}
		_, err := l.SimpanTahun(context.Background(), pelaku, b, true)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "overlaps treaty year "+ada.ID+" (01/01/2025 - 31/12/2036)") {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if got := services.Pesan(func() error {
		_, err := l.SimpanTahun(context.Background(), pelaku, services.TahunMasuk{TreatyYear: "2098", UnderwritingYear: "2098",
			StartDate: "2026-01-01", EndDate: "2026-12-31"}, true)
		return err
	}()); !strings.Contains(got, "Period 01/01/2026 - 31/12/2026 overlaps treaty year "+ada.ID+" (01/01/2025 - 31/12/2036)") {
		t.Errorf("pesan = %q", got)
	}
	// Bersambung: berakhir 31/12/2024, mulai 01/01/2037 - boleh.
	for _, b := range []services.TahunMasuk{
		{TreatyYear: "2018", UnderwritingYear: "2018", StartDate: "2018-01-01", EndDate: "2024-12-31"},
		{TreatyYear: "2037", UnderwritingYear: "2037", StartDate: "2037-01-01", EndDate: "2037-12-31"},
	} {
		if _, err := l.SimpanTahun(context.Background(), pelaku, b, true); err != nil {
			t.Errorf("bersambung %s ditolak: %v", b.StartDate, err)
		}
	}
	// Edit diri sendiri dengan periode yang sama - boleh.
	m.ID = ada.ID
	if _, err := l.SimpanTahun(context.Background(), pelaku, m, false); err != nil {
		t.Errorf("edit diri sendiri ditolak: %v", err)
	}
	// Pemeriksaan di dalam tx sesudah kunci tabel.
	if len(g.Kunci) == 0 || g.Kunci[0] != "tabel|TREATYYEAR_LIFE" {
		t.Errorf("kunci = %v", g.Kunci)
	}
}
