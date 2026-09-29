package services_test

// Aturan murni `Save to RNM` - SaveOutStandingLife_Act. TANPA Oracle.
//
// Pemilik: tiket 03 (giliran 11 paket 1). Dibaca sesudah: simpanrnm.go.

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	intiuang "nusantarare/inti/uang"
)

func uangRNM(t *testing.T, s string) intiuang.Money {
	t.Helper()
	m, err := intiuang.NewMoney(s, "IDR")
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// pesertaRNM adalah peserta yang LOLOS seluruh gerbang langkah 11.
func pesertaRNM(t *testing.T, id string) services.PesertaRNM {
	t.Helper()
	return services.PesertaRNM{Peserta: models.Peserta{
		ID:                  id,
		IsCheck:             "true",
		NomorPolis:          "UJI-POLIS",
		NomorSertifikat:     "UJI-SERT-" + id,
		TanggalKejadian:     "2025-06-01 00:00:00",
		TanggalMulai:        "2025-01-01 00:00:00",
		TanggalExpired:      "2026-01-01 00:00:00",
		ValuasiGrossMulai:   grossMulai,
		ValuasiGrossSelesai: grossSelesai,
		ValuasiRetroMulai:   retroMulai,
		ValuasiRetroSelesai: retroSelesai,
		Dokumen:             []models.Dokumen{{Kategori2: "KTP"}, {Kategori2: "SURAT"}},
		Baris:               []models.BarisAdjustment{{ID: "A-" + id, JumlahKlaim: uangRNM(t, "100")}},
	}}
}

func masukanRNM(t *testing.T, tipe string, p ...services.PesertaRNM) services.MasukanRNM {
	t.Helper()
	return services.MasukanRNM{Tipe: tipe, ContentNote: "DEATH", Peserta: p}
}

func pesanRNM(t *testing.T, err error) string {
	t.Helper()
	var p *services.PelanggaranRNM
	if !errors.As(err, &p) {
		t.Fatalf("galat = %v, mau PelanggaranRNM", err)
	}
	if !errors.Is(err, services.ErrSimpanRNMDitolak) {
		t.Error("PelanggaranRNM harus dapat dikenali sebagai ErrSimpanRNMDitolak")
	}
	return p.Pesan
}

func TestSimpanRNMLolosSeluruhGerbang(t *testing.T) {
	if err := services.PeriksaSimpanRNM(masukanRNM(t, "QR", pesertaRNM(t, "1"), pesertaRNM(t, "2"))); err != nil {
		t.Fatalf("masukan sah ditolak: %v", err)
	}
}

// Langkah 3-4: pesan BERTUMPUK untuk seluruh peserta, lalu keluar.
func TestSimpanRNMDokumenBelumDiunggahMenyebutSetiapNomor(t *testing.T) {
	p1, p2, p3 := pesertaRNM(t, "1"), pesertaRNM(t, "2"), pesertaRNM(t, "3")
	p1.Dokumen, p3.Dokumen = nil, nil
	got := pesanRNM(t, services.PeriksaSimpanRNM(masukanRNM(t, "QR", p1, p2, p3)))
	mau := "The document hasn’t been uploaded person number 1\n" +
		"The document hasn’t been uploaded person number 3"
	if got != mau {
		t.Errorf("pesan = %q\nmau    %q", got, mau)
	}
}

// Langkah 3 `Type=="TP"||"TR"` WhenTrue 3: treaty tidak dituntut berdokumen -
// dan langkah 12 (dokumen lengkap) ter-remark, jadi TIDAK ada gerbang kedua.
//
// ⛔ RALAT GILIRAN-11: uji ini dulu menuntut pesan langkah 12.
func TestSimpanRNMTreatyTanpaDokumenLolos(t *testing.T) {
	for _, tipe := range []string{"TP", "TR"} {
		p := pesertaRNM(t, "1")
		p.Dokumen = nil
		if err := services.PeriksaSimpanRNM(masukanRNM(t, tipe, p)); err != nil {
			t.Errorf("%s tanpa dokumen ditolak: %v", tipe, err)
		}
	}
}

// Langkah 11.4: STS_REJECT baris warisan terbaru 0 atau 1 menolak; "" dan 2 tidak.
func TestSimpanRNMKlaimGandaDeath(t *testing.T) {
	for _, u := range []struct {
		status string
		tolak  bool
	}{{"0", true}, {"1", true}, {"2", false}, {"", false}, {"4", false}} {
		p := pesertaRNM(t, "1")
		p.StatusWarisanTerakhir = u.status
		err := services.PeriksaSimpanRNM(masukanRNM(t, "QR", p))
		if u.tolak {
			if got := pesanRNM(t, err); got != "Person number 1 has already been accepted." {
				t.Errorf("status %q: %q", u.status, got)
			}
		} else if err != nil {
			t.Errorf("status %q ditolak: %v", u.status, err)
		}
	}
}

// Langkah 11.6: bukan DEATH - baris warisan bertertanggung dan DOL sama menolak.
func TestSimpanRNMKlaimGandaHealthHanyaBilaBukanDeath(t *testing.T) {
	p := pesertaRNM(t, "1")
	p.AdaKlaimSehatSamaDOL = true
	m := masukanRNM(t, "QR", p)
	if err := services.PeriksaSimpanRNM(m); err != nil {
		t.Fatalf("DEATH memakai gerbang health: %v", err)
	}
	m.ContentNote = "HEALTH"
	if got := pesanRNM(t, services.PeriksaSimpanRNM(m)); got != "Person number 1 has already been accepted." {
		t.Errorf("health ganda: %q", got)
	}
}

// Langkah 11.7/11.8/11.10: jendela (BEGIN, EXPIRED] - dan TANPA pergeseran retro.
func TestSimpanRNMDOLJendelaTanpaGeserRetro(t *testing.T) {
	const pesan = "DOL cannot be blank or outside the valuation period No 1"
	for _, u := range []struct {
		tipe, dol string
		tolak     bool
	}{
		{"QR", grossMulai, true},    // sama dengan awal: @CompareDates ketat
		{"QR", grossSelesai, false}, // sama dengan akhir: masih di dalam
		{"QP", "2025-09-02 00:00:00", true},
		// ⛔ Inti uji ini: ValidasiDOL menggeser DOL retro SATU HARI (b698),
		// langkah 11.8 TIDAK (b4464 `@addCalendar(.DATE_OF_LOSS,0,...,0)`).
		{"TR", retroMulai, true},
		{"TP", retroSelesai, false},
		{"QR", "", true},
	} {
		p := pesertaRNM(t, "1")
		p.TanggalKejadian = u.dol
		err := services.PeriksaSimpanRNM(masukanRNM(t, u.tipe, p))
		if u.tolak {
			if got := pesanRNM(t, err); got != pesan {
				t.Errorf("%s %q: %q", u.tipe, u.dol, got)
			}
		} else if err != nil {
			t.Errorf("%s %q ditolak: %v", u.tipe, u.dol, err)
		}
	}
}

// Langkah 11.12-11.17, urut, dengan kalimat VERBATIM (termasuk salah ejanya).
func TestSimpanRNMMedanKosongUrutXML(t *testing.T) {
	for _, u := range []struct {
		nama  string
		ubah  func(p *services.PesertaRNM)
		pesan string
	}{
		{"DOB", func(p *services.PesertaRNM) { p.DOBKosong = true }, "DOB cannnot be blank No 2"},
		{"BEGIN", func(p *services.PesertaRNM) { p.TanggalMulai = "" }, "Begin Date cannnot be blank No 2"},
		{"EXPIRED", func(p *services.PesertaRNM) { p.TanggalExpired = "" }, "Expired Date cannnot be blank No 2"},
		{"POLICY", func(p *services.PesertaRNM) { p.NomorPolis = "" }, "Policy No cannnot be blank No 2, please contact IT"},
		{"CERT", func(p *services.PesertaRNM) { p.NomorSertifikat = "" }, "Certificate No cannnot be blank No 2, please contact IT"},
		{"GROSS", func(p *services.PesertaRNM) { p.Baris[0].JumlahKlaim.Amount = nil }, "Claim Gross No 2 can't null"},
	} {
		p2 := pesertaRNM(t, "2")
		u.ubah(&p2)
		got := pesanRNM(t, services.PeriksaSimpanRNM(masukanRNM(t, "QR", pesertaRNM(t, "1"), p2)))
		if got != u.pesan {
			t.Errorf("%s: %q, mau %q", u.nama, got, u.pesan)
		}
	}
	// DOB kosong DAN BEGIN kosong: yang PERTAMA di urutan XML yang menang.
	p := pesertaRNM(t, "1")
	p.DOBKosong, p.TanggalMulai = true, ""
	if got := pesanRNM(t, services.PeriksaSimpanRNM(masukanRNM(t, "QR", p))); !strings.HasPrefix(got, "DOB") {
		t.Errorf("urutan: %q", got)
	}
}

// Langkah 12 ter-remark (b6178): dokumen yang ADA tetapi tidak lengkap lolos.
func TestSimpanRNMDokumenTidakLengkapLolos(t *testing.T) {
	p := pesertaRNM(t, "1")
	p.Dokumen = []models.Dokumen{{Kategori2: "KTP"}, {Kategori2: "KTP"}}
	if err := services.PeriksaSimpanRNM(masukanRNM(t, "QR", p)); err != nil {
		t.Errorf("dokumen kembar ditolak: %v", err)
	}
}

func TestSimpanRNMTypeAsingGagalTerang(t *testing.T) {
	if err := services.PeriksaSimpanRNM(masukanRNM(t, "XX", pesertaRNM(t, "1"))); !errors.Is(err, services.ErrTypeTidakDikenal) {
		t.Errorf("Type asing: %v", err)
	}
}

// Langkah 27: tiga kode retro/security reinsurer KELUAR sebelum Arasapas.
func TestArasapasDilewatiUntukTigaKodeRetro(t *testing.T) {
	for _, u := range []struct {
		retro, sec string
		lewat      bool
	}{
		{"L0000141", "", true}, {"", "L0000134", true}, {"1000013", "", true},
		{"L0000134", "", false}, {"", "L0000141", false}, {"", "", false},
	} {
		got := services.ArasapasDilewatiRetro(services.PolisRetro{
			Tipe: "QR", RetroID: u.retro, SecurityReinsurerID: u.sec})
		if got != u.lewat {
			t.Errorf("retro %q sec %q: %v; mau %v", u.retro, u.sec, got, u.lewat)
		}
	}
}

// Langkah 27 membaca SALINAN yang ditukar `InsertJsonClaimLife_Act` langkah 2.
// OQ-N5 DITUTUP (GILIRAN-17): keputusan Komite "cutover 7 Feb 2025 tidak
// dipakai" berlaku juga di sini - tukar cukup DUA WHEN (b1268 TP/TR, b1291
// SecurityReinsurer terisi); b1314 `ProdDateTime` tidak dipakai.
func TestArasapasRetroTukarDuaSyarat(t *testing.T) {
	polis := func(tipe, retro, sec, secNama string) services.PolisRetro {
		return services.PolisRetro{Tipe: tipe, RetroID: retro, SecurityReinsurerID: sec,
			SecurityReinsurer: secNama}
	}
	for _, u := range []struct {
		p     services.PolisRetro
		lewat bool
	}{
		// Dua WHEN benar -> DITUKAR, dan gerbang membaca nilai sesudah tukar.
		{polis("TP", "L0000141", "S-1", "UJI-SEC"), false}, // retro sesudah tukar S-1
		{polis("TR", "R-1", "L0000141", "UJI-SEC"), true},  // retro sesudah tukar L0000141
		{polis("TP", "R-1", "L0000134", "UJI-SEC"), false}, // security sesudah tukar R-1
		{polis("TR", "L0000134", "R-1", "UJI-SEC"), true},  // security sesudah tukar L0000134
		{polis("TP", "R-1", "S-1", "UJI-SEC"), false},
		{polis("TR", "L0000141", "L0000141", "UJI-SEC"), true},
		// b1291: nama security reinsurer kosong -> tidak ditukar.
		{polis("TP", "L0000141", "S-1", ""), true},
		// b1268: bukan TP/TR -> tidak ditukar.
		{polis("QR", "R-1", "L0000141", "UJI-SEC"), false},
	} {
		if got := services.ArasapasDilewatiRetro(u.p); got != u.lewat {
			t.Errorf("%+v: %v; mau %v", u.p, got, u.lewat)
		}
	}
}

// TestSimpanRNMMenjagaPagarnya - identitas, pengenal, dan basis data SEBELUM
// apa pun dibaca.
func TestSimpanRNMMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	w := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC)
	if _, err := svc.SimpanRNM().Simpan(ctx, inti.Pelaku{}, "K-1", w); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v", err)
	}
	admin := pelakuBerperan(inti.PeranAdmin)
	if _, err := svc.SimpanRNM().Simpan(ctx, admin, " ", w); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("tanpa pengenal: %v", err)
	}
	if _, err := svc.SimpanRNM().Simpan(ctx, admin, "K-1", w); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
}
