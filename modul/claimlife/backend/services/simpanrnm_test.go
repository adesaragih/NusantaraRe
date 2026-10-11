package services_test

// Aturan murni `Save to RNM` - SaveOutStandingLife_Act. TANPA Oracle.
//
// Pemilik: tiket 03 (giliran 11 paket 1). Dibaca sesudah: simpanrnm.go.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	intiuang "nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/modul/claimlife/backend/services"
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
		ID:              id,
		IsCheck:         "true",
		IsAccept:        "true",
		NomorPolis:      "UJI-POLIS",
		NomorSertifikat: "UJI-SERT-" + id,
		TanggalKejadian: "2025-06-01 00:00:00",
		// 11-10-2026: Claim Received Date dan Confirmation Date wajib sebelum Save to RNM.
		TanggalKlaimTeks:    models.TanggalKlaimTeks{TanggalTerimaKlaim: "2025-06-02 00:00:00", TanggalKonfirmasi: "2025-06-03 00:00:00"},
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
	// Is KPR terisi: wajib sebelum Save to RNM (11-10-2026).
	return services.MasukanRNM{Tipe: tipe, ContentNote: "DEATH", IsKPR: "Non-KPR", Peserta: p}
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

// Langkah 11.7/11.8/11.10: jendela [BEGIN, EXPIRED] INKLUSIF tanpa geser (keputusan work owner 11-10-2026) - dan Save
// grid (ValidasiDOL) memberi hasil yang SAMA untuk tanggal yang sama.
func TestSimpanRNMDOLJendelaTanpaGeserRetro(t *testing.T) {
	const pesan = "DOL cannot be blank or outside the valuation period No 1"
	for _, u := range []struct {
		tipe, dol string
		tolak     bool
	}{
		{"QR", grossMulai, false},   // BEGIN termasuk
		{"QR", grossSelesai, false}, // EXPIRED termasuk
		{"QR", "2025-02-28 00:00:00", true},
		{"QP", "2025-09-02 00:00:00", true},
		{"TR", retroMulai, false},
		{"TR", "2025-03-31 00:00:00", true},
		{"TP", retroSelesai, false},
		{"TP", "2025-10-02 00:00:00", true},
		{"QR", "", true},
	} {
		p := pesertaRNM(t, "1")
		p.TanggalKejadian = u.dol
		err := services.PeriksaSimpanRNM(masukanRNM(t, u.tipe, p))
		if u.tolak {
			mau := pesan
			if u.dol != "" {
				// DOL terisi di luar jendela = kalimat layar dengan rentang valuasi (gross / retro).
				mau = "Insured number 1: Date of Loss must be within the valuation period (" + jendelaUji(u.tipe) + ")."
			}
			if got := pesanRNM(t, err); got != mau {
				t.Errorf("%s %q: %q", u.tipe, u.dol, got)
			}
		} else if err != nil {
			t.Errorf("%s %q ditolak: %v", u.tipe, u.dol, err)
		}
		if u.dol == "" {
			continue
		}
		// Save grid (ValidasiDOL) sepakat dengan Save to RNM untuk tanggal yang sama.
		grid := services.ValidasiDOL(u.tipe, saat(t, u.dol), p.Peserta)
		if (grid != nil) != u.tolak {
			t.Errorf("%s %q: Save grid %v, Save to RNM tolak=%v - harus sama", u.tipe, u.dol, grid, u.tolak)
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
		{"GROSS", func(p *services.PesertaRNM) { p.Baris[0].JumlahKlaim.Amount = nil }, "Claim RNM is required for insured number 2"},
	} {
		p2 := pesertaRNM(t, "2")
		u.ubah(&p2)
		got := pesanRNM(t, services.PeriksaSimpanRNM(masukanRNM(t, "QR", pesertaRNM(t, "1"), p2)))
		if got != u.pesan {
			t.Errorf("%s: %q, mau %q", u.nama, got, u.pesan)
		}
	}
	// DOB kosong DAN BEGIN kosong: KEDUANYA dilaporkan (11-10-2026), tetap dalam urutan XML.
	p := pesertaRNM(t, "1")
	p.DOBKosong, p.TanggalMulai = true, ""
	if got := pesanRNM(t, services.PeriksaSimpanRNM(masukanRNM(t, "QR", p))); got != "DOB cannnot be blank No 1\nBegin Date cannnot be blank No 1" {
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

// Keputusan work owner 11-10-2026: Claim Received Date dan Confirmation Date WAJIB sebelum Save to RNM (per peserta,
// semua peserta seperti gerbang 11.x); DOL kosong tetap ditolak langkah 11.10 dengan kalimat XML.
func TestSimpanRNMTanggalWajib(t *testing.T) {
	for _, u := range []struct {
		nama string
		ubah func(*services.PesertaRNM)
		mau  string
	}{
		{"claim received", func(p *services.PesertaRNM) { p.TanggalTerimaKlaim = "" }, "Claim Received Date is required for insured number 1"},
		{"confirmation", func(p *services.PesertaRNM) { p.TanggalKonfirmasi = "" }, "Confirmation Date is required for insured number 1"},
		{"dol", func(p *services.PesertaRNM) { p.TanggalKejadian = "" }, "DOL cannot be blank or outside the valuation period No 1"},
	} {
		p := pesertaRNM(t, "1")
		u.ubah(&p)
		err := services.PeriksaSimpanRNM(services.MasukanRNM{Tipe: "QP", ContentNote: "HEALTH", IsKPR: "KPR", Peserta: []services.PesertaRNM{p}})
		var tolak *services.PelanggaranRNM
		if !errors.As(err, &tolak) || tolak.Pesan != u.mau {
			t.Errorf("%s: galat %v, mau %q", u.nama, err, u.mau)
		}
	}
}

// Keputusan work owner 11-10-2026: Save to RNM mengumpulkan SELURUH kekurangan kolom wajib semua peserta dalam SATU
// penolakan (bukan berhenti di peserta/kolom pertama seperti Pega).
func TestSimpanRNMKekuranganDikumpulkanSemuaPeserta(t *testing.T) {
	p1, p2 := pesertaRNM(t, "1"), pesertaRNM(t, "2")
	p1.NamaTertanggung = "UJI-NAMA-1"
	p1.TanggalKejadian, p1.TanggalTerimaKlaim = "", ""
	p1.Baris[0].JumlahKlaim.Amount = nil
	p2.TanggalKonfirmasi = ""
	p2.Dokumen = nil
	m := masukanRNM(t, "QR", p1, p2)
	m.IsKPR = " "
	var tolak *services.PelanggaranRNM
	if !errors.As(services.PeriksaSimpanRNM(m), &tolak) {
		t.Fatal("mau PelanggaranRNM")
	}
	type k struct {
		no          int
		id, sert    string
		kolom, nama string
	}
	mau := []k{
		{0, "", "", "Is KPR", ""},
		{1, "1", "UJI-SERT-1", "Date of Loss", "UJI-NAMA-1"},
		{1, "1", "UJI-SERT-1", "Claim Received Date", "UJI-NAMA-1"},
		{1, "1", "UJI-SERT-1", "Claim RNM", "UJI-NAMA-1"},
		{2, "2", "UJI-SERT-2", "Document", ""},
		{2, "2", "UJI-SERT-2", "Confirmation Date", ""},
	}
	if len(tolak.Kekurangan) != len(mau) {
		t.Fatalf("kekurangan %d: %+v", len(tolak.Kekurangan), tolak.Kekurangan)
	}
	for i, x := range tolak.Kekurangan {
		if got := (k{x.NomorPeserta, x.PesertaID, x.Sertifikat, x.Kolom, x.Nama}); got != mau[i] {
			t.Errorf("#%d %+v, mau %+v", i, got, mau[i])
		}
	}
	// DOL kosong tetap kalimat XML 11.10; Pesan = gabungan baris tiap kekurangan.
	if tolak.Kekurangan[1].Pesan != "DOL cannot be blank or outside the valuation period No 1" {
		t.Errorf("DOL: %q", tolak.Kekurangan[1].Pesan)
	}
	if !strings.Contains(tolak.Pesan, "Is KPR is required.\nDOL cannot be blank") || strings.Count(tolak.Pesan, "\n") != len(mau)-1 {
		t.Errorf("pesan: %q", tolak.Pesan)
	}
	if tolak.Langkah != "wajib" {
		t.Errorf("langkah %q", tolak.Langkah)
	}
}

// Claim RNM wajib > 0: kosong dan 0 ditolak, nilai positif lolos (KlaimRNMMinimal).
func TestSimpanRNMClaimRNMWajibLebihDariNol(t *testing.T) {
	for _, u := range []struct {
		nilai string
		tolak bool
	}{{"", true}, {"0", true}, {"0.00", true}, {"0.01", false}} {
		p := pesertaRNM(t, "1")
		if u.nilai == "" {
			p.Baris[0].JumlahKlaim.Amount = nil
		} else {
			p.Baris[0].JumlahKlaim = uangRNM(t, u.nilai)
		}
		err := services.PeriksaSimpanRNM(masukanRNM(t, "QR", p))
		if !u.tolak {
			if err != nil {
				t.Errorf("%q ditolak: %v", u.nilai, err)
			}
			continue
		}
		if got := pesanRNM(t, err); got != "Claim RNM is required for insured number 1" {
			t.Errorf("%q: %q", u.nilai, got)
		}
	}
}

// Klaim ganda (bukan kolom kosong) tetap berhenti dini - kekurangan kolom TIDAK ikut dikumpulkan.
func TestSimpanRNMKlaimGandaBerhentiDini(t *testing.T) {
	p1, p2 := pesertaRNM(t, "1"), pesertaRNM(t, "2")
	p1.TanggalTerimaKlaim = ""
	p2.StatusWarisanTerakhir = "1"
	var tolak *services.PelanggaranRNM
	if !errors.As(services.PeriksaSimpanRNM(masukanRNM(t, "QR", p1, p2)), &tolak) {
		t.Fatal("mau PelanggaranRNM")
	}
	if tolak.Pesan != "Person number 2 has already been accepted." || len(tolak.Kekurangan) != 0 {
		t.Errorf("%q %+v", tolak.Pesan, tolak.Kekurangan)
	}
}

// Nomor klaim TIDAK diambil sebelum seluruh gerbang lolos: PeriksaSimpanRNM mendahului NomorBerikut di Simpan.
func TestSimpanRNMGerbangSebelumPenomoran(t *testing.T) {
	isi, err := os.ReadFile("simpanrnm.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	gerbang := strings.Index(s, "if err := PeriksaSimpanRNM(m); err != nil")
	nomor := strings.Index(s, "x.penomor.NomorBerikut(")
	if gerbang < 0 || nomor < 0 || gerbang > nomor {
		t.Errorf("gerbang %d, penomoran %d: gerbang harus lebih dulu", gerbang, nomor)
	}
}

// jendelaUji - rentang valuasi fixture (grossMulai..grossSelesai / retroMulai..retroSelesai) dalam bentuk layar.
func jendelaUji(tipe string) string {
	if tipe == "QR" || tipe == "QP" {
		return "01/03/2025 – 01/09/2025"
	}
	return "01/04/2025 – 01/10/2025"
}

// 11-10-2026: Claim RNM > Share Nusantara Re dipindah dari Save grid ke Save to RNM - dikumpulkan bersama kekurangan lain,
// kalimat sama dengan layar (angka bergaya Indonesia).
func TestSimpanRNMMenolakClaimRNMMelebihiShare(t *testing.T) {
	p1, p2 := pesertaRNM(t, "1"), pesertaRNM(t, "2")
	p1.Baris[0].ShareNusantaraRe = uangRNM(t, "5077429.5")
	p1.Baris[0].JumlahKlaim = uangRNM(t, "6000000")
	p2.TanggalKejadian = "2026-10-11 00:00:00" // di luar jendela gross
	var tolak *services.PelanggaranRNM
	if !errors.As(services.PeriksaSimpanRNM(masukanRNM(t, "QR", p1, p2)), &tolak) {
		t.Fatal("mau PelanggaranRNM")
	}
	mau := []string{
		"Insured number 1: Claim RNM cannot be greater than Share Nusantara Re (5.077.429,5).",
		"Insured number 2: Date of Loss must be within the valuation period (01/03/2025 – 01/09/2025).",
	}
	if len(tolak.Kekurangan) != 2 || tolak.Kekurangan[0].Pesan != mau[0] || tolak.Kekurangan[1].Pesan != mau[1] {
		t.Errorf("kekurangan %+v", tolak.Kekurangan)
	}
	// Sama dengan share: lolos.
	p1.Baris[0].JumlahKlaim = uangRNM(t, "5077429.50")
	if err := services.PeriksaSimpanRNM(masukanRNM(t, "QR", p1)); err != nil {
		t.Errorf("Claim RNM = Share ditolak: %v", err)
	}
}

// Keputusan work owner 11-10-2026: baris OS_AKSEPTASI_KLAIM_LIFE Save to RNM membawa nilai polis, peserta, dan baris.
func TestSusunCerminSaveRNMDariPolisPesertaBaris(t *testing.T) {
	wpc := "2025-01-15"
	polis := kontrak.PolisRingkas{NomorPolis: "UJI-POL", Type: "QR", PolicyHolderName: "UJI PEMEGANG", CedingCoName: "UJI CEDING",
		BusinessName: "UJI PLAN BISNIS", BusinessCode: "L01", SobName: "UJI SOB", ProductNameID: "UJI-PRD", ProductName: "UJI PRODUK",
		TypeCeding: "3", WPC: &wpc, RetroID: "R1", RetroName: "UJI RETRO", SecurityReinsurerID: "S1", SecurityReinsurer: "UJI SEC"}
	p := pesertaRNM(t, "1").Peserta
	p.JenisKelamin, p.Umur = "M", "45"
	p.TanggalKejadian = "2025-06-01 00:00:00"
	b := p.Baris[0]
	b.SumInsured, b.ShareNusantaraRe = uangRNM(t, "200000000"), uangRNM(t, "100000")
	dasar := repository.BarisLama{ID: "A-1", CASEID: "K-1", NO_CLAIM: "UJI-CLM", NO_ACCEPTATION: "LAMA", CLAIM_AMOUNT: "100"}
	saat := time.Date(2026, 10, 11, 18, 30, 0, 0, time.UTC) // 01:30 WIB tanggal 12
	d := services.SusunCerminSaveRNM(dasar, p, b, polis, "UJI-SOB-ID", saat)
	cek := map[string][2]string{
		"NO_ACCEPTATION": {d.NO_ACCEPTATION, "12/10/2026"},
		// 11-10-2026: ACCEPTATION_DATE = tanggal outstanding juga (DATE, waktu WIB).
		"ACCEPTATION_DATE": {d.ACCEPTATION_DATE, "2026-10-12 01:30:00"},
		// Kolom Medical Check tetap kosong di Save to RNM.
		"STATUS":        {d.STATUS + d.DISEASE + d.ICD_CODE + d.NOTES + d.KETERANGAN, ""},
		"POLICY_NO":     {d.POLICY_NO, "UJI-POL"},
		"POLICY_HOLDER": {d.POLICY_HOLDER, "UJI PEMEGANG"},
		"SEX":           {d.SEX, "M"}, "AGE": {d.AGE, "45"},
		"PLAN":         {d.PLAN, "UJI PLAN BISNIS"},
		"BEGIN_DATE":   {d.BEGIN_DATE, p.TanggalMulai},
		"LAPSE_DATE":   {d.LAPSE_DATE, "2025-06-01 00:00:00"},
		"EXPIRED_DATE": {d.EXPIRED_DATE, p.TanggalExpired},
		"WPC":          {d.WPC, "2025-01-15"},
		"CEDINGCONAME": {d.CEDINGCONAME, "UJI CEDING"},
		"SOB":          {d.SOB, "UJI-SOB-ID"}, "SOBNAME": {d.SOBNAME, "UJI SOB"},
		"BUSINESSID": {d.BUSINESSID, "L01"}, "BUSINESSNAME": {d.BUSINESSNAME, "UJI PLAN BISNIS"},
		"SUM_INSURED": {d.SUM_INSURED, "200000000"}, "SHARE_NUSANTARA_RE": {d.SHARE_NUSANTARA_RE, "100000"},
		"CLAIM_AMOUNT": {d.CLAIM_AMOUNT, "100"},
		"TYPECEDING":   {d.TYPECEDING, "QS + SURPLUS"},
		"RETROID":      {d.RETROID, "R1"}, "SECURITYREINSURER": {d.SECURITYREINSURER, "UJI SEC"},
		"PRODUCTNAMEID":       {d.PRODUCTNAMEID, "UJI-PRD"},
		"CLAIM_RECEIVED_DATE": {d.CLAIM_RECEIVED_DATE, p.TanggalTerimaKlaim},
		"CONFIRMATION_DATE":   {d.CONFIRMATION_DATE, p.TanggalKonfirmasi},
	}
	for k, v := range cek {
		if v[0] != v[1] {
			t.Errorf("%s = %q, mau %q", k, v[0], v[1])
		}
	}
	// TP/TR dengan Security Reinsurer: ditukar (b9671).
	polis.Type = "TR"
	d = services.SusunCerminSaveRNM(dasar, p, b, polis, "", saat)
	if d.RETROID != "S1" || d.RETRONAME != "UJI SEC" || d.SECURITYREINSURERID != "R1" || d.SECURITYREINSURER != "UJI RETRO" {
		t.Errorf("tukar retro: %+v", d)
	}
	if err := repository.PeriksaNilaiWarisan(d); err != nil {
		t.Error(err)
	}
	if err := repository.PeriksaLebarWarisan(d); err != nil {
		t.Error(err)
	}
}

// Save to RNM menyusun baris datar dengan SusunCerminSaveRNM dan SOB polis berkunci IDPEGA yang sama dengan header polis.
func TestSimpanRNMMemakaiSusunCermin(t *testing.T) {
	isi, err := os.ReadFile("simpanrnm.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	for _, mau := range []string{"polis, kunciPolis, err := RingkasPolisIdPega(", "sob, err := baca.SOBPolis(ctx, kunciPolis)",
		"datar[b.ID] = SusunCerminSaveRNM(d, p, b, polis, sob, saat)"} {
		if !strings.Contains(s, mau) {
			t.Errorf("simpanrnm.go tanpa %q", mau)
		}
	}
	// Ditulis di transaksi UTAMA, sesudah transaksi penomoran.
	if strings.Index(s, "x.penomor.NomorBerikut(") > strings.Index(s, "pohon.TulisCerminSaveRNM(") {
		t.Error("baris datar ditulis sebelum penomoran")
	}
}

// 11-10-2026 ("iya itu hanya di claim analis saja"): penyerahan ke Komite hanya dari tahap Claim Analis (409 di tahap lain).
func TestSerahkanKomiteHanyaClaimAnalis(t *testing.T) {
	isi, err := os.ReadFile("komite.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	fn := s[strings.Index(s, "func (p *Penyerahan) Serahkan("):]
	gerbang := strings.Index(fn, "if tahap != models.TahapClaimAnalis {")
	tulis := strings.Index(fn, "DalamTransaksi(")
	if gerbang < 0 || (tulis >= 0 && gerbang > tulis) {
		t.Errorf("gerbang tahap Claim Analis %d harus mendahului penulisan %d", gerbang, tulis)
	}
	if !errors.Is(fmt.Errorf("x: %w", services.ErrSerahkanBukanClaimAnalis), services.ErrSerahkanBukanClaimAnalis) {
		t.Error("galat tidak dapat dikenali")
	}
}
