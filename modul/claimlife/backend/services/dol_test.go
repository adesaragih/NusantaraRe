package services_test

// Validasi Date of Loss dan penurunan jenis klaim - TANPA Oracle.
//
// Pemilik: tiket 06. Dibaca sesudah: dol.go.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// Jendela valuasi uji. Gross dan retro sengaja BERBEDA, supaya cabang yang
// salah tidak lolos hanya karena kedua jendela kebetulan sama.
const (
	grossMulai   = "2025-03-01 00:00:00"
	grossSelesai = "2025-09-01 00:00:00"
	retroMulai   = "2025-04-01 00:00:00"
	retroSelesai = "2025-10-01 00:00:00"
)

func pesertaDOL(t *testing.T) models.Peserta {
	t.Helper()
	return models.Peserta{
		ID:                  "UJI-P-DOL",
		ValuasiGrossMulai:   grossMulai,
		ValuasiGrossSelesai: grossSelesai,
		ValuasiRetroMulai:   retroMulai,
		ValuasiRetroSelesai: retroSelesai,
	}
}

func saat(t *testing.T, s string) time.Time {
	t.Helper()
	w, err := utils.ParseTanggal(s)
	if err != nil {
		t.Fatalf("saat %q: %v", s, err)
	}
	return w
}

// TestBatasJendelaTiapType menguji KEEMPAT batas untuk tiap Type.
//
// Keputusan work owner 11-10-2026 (menyimpang sadar dari Pega): jendela [BEGIN, EXPIRED] INKLUSIF di kedua ujung,
// tanpa pergeseran retro b698, per tanggal kalender - aturan yang SAMA dengan Save to RNM 11.7-11.8 (DOLDalamJendela).
func TestBatasJendelaTiapType(t *testing.T) {
	p := pesertaDOL(t)
	kasus := []struct {
		tipe string
		dol  string
		sah  bool
		apa  string
	}{
		{"QR", "2025-02-28 00:00:00", false, "BEGIN-1"},
		{"QR", grossMulai, true, "tepat BEGIN: termasuk"},
		{"QR", grossSelesai, true, "tepat EXPIRED: termasuk"},
		{"QR", "2025-09-01 23:00:00", true, "hari EXPIRED, jam diabaikan"},
		{"QR", "2025-09-02 00:00:00", false, "EXPIRED+1"},
		{"QP", "2025-02-28 00:00:00", false, "BEGIN-1"},
		{"QP", grossMulai, true, "tepat BEGIN"},
		{"QP", grossSelesai, true, "tepat EXPIRED"},
		{"QP", "2025-09-02 00:00:00", false, "EXPIRED+1"},
		{"TR", "2025-03-31 00:00:00", false, "BEGIN-1 (dulu lolos karena geser +1 hari)"},
		{"TR", retroMulai, true, "tepat BEGIN"},
		{"TR", retroSelesai, true, "tepat EXPIRED (dulu ditolak karena geser)"},
		{"TR", "2025-10-02 00:00:00", false, "EXPIRED+1"},
		{"TP", "2025-03-31 00:00:00", false, "BEGIN-1"},
		{"TP", retroMulai, true, "tepat BEGIN"},
		{"TP", retroSelesai, true, "tepat EXPIRED"},
		{"TP", "2025-10-02 00:00:00", false, "EXPIRED+1"},
		// Cabang terpisah: sah di satu jendela, tidak di jendela lain.
		{"QR", "2025-09-15 00:00:00", false, "sesudah gross berakhir, di dalam retro"},
		{"TR", "2025-09-15 00:00:00", true, "sesudah gross berakhir, di dalam retro"},
		{"QR", "2025-03-15 00:00:00", true, "di dalam gross, sebelum retro mulai"},
		{"TR", "2025-03-15 00:00:00", false, "di dalam gross, sebelum retro mulai"},
	}
	for _, k := range kasus {
		err := services.ValidasiDOL(k.tipe, saat(t, k.dol), p)
		if k.sah && err != nil {
			t.Errorf("%s %s (%s): ditolak - %v", k.tipe, k.dol, k.apa, err)
		}
		if !k.sah && !errors.Is(err, services.ErrDOLTidakSah) {
			t.Errorf("%s %s (%s): galat = %v, mau ErrDOLTidakSah", k.tipe, k.dol, k.apa, err)
		}
	}
}

// TestPesanDOLPersisSepertiXML - `local.errmsg = "Invalid DOL"`, tanpa apa pun
// yang ditambahkan. Peserta mana yang bermasalah diketahui pemanggilnya, sebab
// ia menyerahkan pesertanya satu per satu.
func TestPesanDOLPersisSepertiXML(t *testing.T) {
	err := services.ValidasiDOL("QR", saat(t, "2025-02-28 00:00:00"), pesertaDOL(t))
	if err == nil {
		t.Fatal("tanggal sebelum BEGIN lolos")
	}
	const mau = "Invalid DOL"
	if err.Error() != mau {
		t.Errorf("pesan = %q, mau persis %q", err.Error(), mau)
	}
}

// TestTypeTidakDikenalDitolak - nol cabang di XML berarti pertanyaan yang belum
// dijawab, bukan izin lewat.
func TestTypeTidakDikenalDitolak(t *testing.T) {
	for _, tipe := range []string{"", "XX", "Q", "qr "} {
		err := services.ValidasiDOL(tipe, saat(t, "2025-06-01 00:00:00"), pesertaDOL(t))
		if !errors.Is(err, services.ErrTypeTidakDikenal) {
			t.Errorf("Type %q: galat = %v, mau ErrTypeTidakDikenal", tipe, err)
		}
	}
}

// TestValuasiKosongDitolakBukanDianggapNol - teks kosong pada kolom tanggal
// menjadi KOSONG, bukan tanggal nol (ADR-U-0022 Akibat 2). Tanggal nol akan
// membuat SELURUH DOL tampak sesudah tanggal mulai.
func TestValuasiKosongDitolakBukanDianggapNol(t *testing.T) {
	p := pesertaDOL(t)
	p.ValuasiGrossMulai = ""
	err := services.ValidasiDOL("QR", saat(t, "2025-06-01 00:00:00"), p)
	if !errors.Is(err, models.ErrValuasiKosong) {
		t.Fatalf("galat = %v, mau ErrValuasiKosong", err)
	}
	if !strings.Contains(err.Error(), "UJI-P-DOL") {
		t.Errorf("pesan tidak menyebut peserta: %v", err)
	}
	// Cabang retro tidak ikut terpengaruh kolom gross yang kosong.
	if err := services.ValidasiDOL("TR", saat(t, "2025-06-01 00:00:00"), p); err != nil {
		t.Errorf("cabang retro ikut gagal karena kolom gross kosong: %v", err)
	}
}

// TestValuasiTakTeruraiDilaporkan - teks yang bukan tanggal dilaporkan dengan
// nama kolomnya, tidak dibulatkan menjadi "tidak sah".
func TestValuasiTakTeruraiDilaporkan(t *testing.T) {
	p := pesertaDOL(t)
	p.ValuasiGrossSelesai = "bukan-tanggal"
	err := services.ValidasiDOL("QR", saat(t, "2025-06-01 00:00:00"), p)
	if err == nil || errors.Is(err, services.ErrDOLTidakSah) {
		t.Fatalf("galat = %v; teks tak terurai bukan DOL tidak sah", err)
	}
	if !strings.Contains(err.Error(), "GROSS_VALUATION_EXPIRED_DATE") {
		t.Errorf("pesan tidak menyebut NAMA KOLOM SEBENARNYA: %v", err)
	}
}

// TestContentNoteSeluruhDuaPuluhSatuKode mengunci tabel data acuan.
//
// Angkanya dua puluh satu, dan itu dicacah DUA CARA: daftar di bawah ditulis
// tangan dari CONTEXT.md, dan cacah kunci tabelnya dibandingkan dengannya.
func TestContentNoteSeluruhDuaPuluhSatuKode(t *testing.T) {
	mau := map[string]string{
		"L1": "DEATH", "L2": "DEATH", "L3": "DEATH", "L4": "DEATH",
		"L5": "DEATH", "L6": "DEATH", "L7": "DEATH", "L8": "DEATH",
		"L9": "DEATH", "L10": "DEATH", "L11": "DEATH",
		"L12": "HEALTH", "L13": "HEALTH", "L14": "HEALTH", "L15": "HEALTH",
		"L16": "CI", "L17": "TPD", "L18": "HEALTH",
		"L19": "CI", "L20": "TPD", "L21": "TI",
	}
	if n := len(models.BusinessCodeContentNote()); n != len(mau) {
		t.Fatalf("cacah kode = %d, mau %d", n, len(mau))
	}
	for kode, note := range mau {
		got, err := services.ContentNoteDari(kode)
		if err != nil {
			t.Errorf("%s: %v", kode, err)
			continue
		}
		if got != note {
			t.Errorf("%s = %q, mau %q", kode, got, note)
		}
	}
}

// TestContentNoteKodeAsingDitolak - kode di luar daftar adalah produk yang
// belum dikenal, bukan alasan mengarang jenis klaim.
func TestContentNoteKodeAsingDitolak(t *testing.T) {
	for _, kode := range []string{"", "L0", "L22", "l1", "X1", "1"} {
		if _, err := services.ContentNoteDari(kode); !errors.Is(
			err, services.ErrBusinessCodeTidakDikenal) {
			t.Errorf("kode %q: galat = %v, mau ErrBusinessCodeTidakDikenal", kode, err)
		}
	}
}

// TestDaftarContentNoteSalinanBukanAslinya - pemanggil tidak boleh dapat
// mengubah data acuan seluruh aplikasi.
func TestDaftarContentNoteSalinanBukanAslinya(t *testing.T) {
	d := models.BusinessCodeContentNote()
	d["L1"] = "DIUBAH"
	if got, _ := services.ContentNoteDari("L1"); got != "DEATH" {
		t.Errorf("daftar aslinya ikut berubah: L1 = %q", got)
	}
}

// TestDOLKosongDibedakanDariDOLDiLuarJendela - waktu nol pasti di luar jendela
// mana pun, jadi tanpa penjaga tersendiri pengguna mendapat "Invalid DOL"
// untuk tanggal yang tidak pernah ia isi. Pesan yang menyesatkan ke arah yang
// salah lebih buruk daripada tidak ada pesan.
func TestDOLKosongDibedakanDariDOLDiLuarJendela(t *testing.T) {
	err := services.ValidasiDOL("QR", time.Time{}, pesertaDOL(t))
	if !errors.Is(err, services.ErrDOLKosong) {
		t.Fatalf("galat = %v, mau ErrDOLKosong", err)
	}
	if errors.Is(err, services.ErrDOLTidakSah) {
		t.Error("DOL kosong dilaporkan sebagai DOL tidak sah")
	}
}

// TestGalatDOLMembawaPesertaDiMedanSendiri - kalimatnya persis XML, dan
// pengenal pesertanya tetap dapat diambil pemanggil.
func TestGalatDOLMembawaPesertaDiMedanSendiri(t *testing.T) {
	err := services.ValidasiDOL("QR", saat(t, "2025-02-28 00:00:00"), pesertaDOL(t))
	var g services.GalatDOL
	if !errors.As(err, &g) {
		t.Fatalf("galat = %v, mau services.GalatDOL", err)
	}
	if g.PesertaID != "UJI-P-DOL" {
		t.Errorf("PesertaID = %q, mau UJI-P-DOL", g.PesertaID)
	}
	if g.Error() != "Invalid DOL" {
		t.Errorf("kalimatnya berubah: %q", g.Error())
	}
}

// TestSetTanggalKejadianTanpaPelakuDitolak - fail-closed atas pelaku anonim,
// alasan yang sama dengan pendaftaran: perubahan tanggal kejadian mengubah
// hasil validasi klaim, dan perubahan tanpa identitas tidak dapat ditelusuri.
func TestSetTanggalKejadianTanpaPelakuDitolak(t *testing.T) {
	svc := services.New(nil)
	err := svc.TanggalKejadian().Set(context.Background(), inti.Pelaku{},
		"CLM-000001", "UJI-P-1", saat(t, "2025-06-01 00:00:00"))
	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("galat = %v, mau ErrTanpaIdentitas", err)
	}
}

// TestSetTanggalKejadianTanpaPengenalDitolak - pengenal yang kosong tidak
// pernah diteruskan ke SQL sebagai teks kosong.
func TestSetTanggalKejadianTanpaPengenalDitolak(t *testing.T) {
	svc := services.New(nil)
	pelaku := inti.Pelaku{AkunID: "UJI-AKUN"}
	for _, k := range []struct{ klaim, peserta string }{
		{"", "UJI-P-1"}, {"CLM-000001", ""}, {"   ", "   "},
	} {
		err := svc.TanggalKejadian().Set(context.Background(), pelaku,
			k.klaim, k.peserta, saat(t, "2025-06-01 00:00:00"))
		if !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("klaim %q peserta %q: galat = %v, mau ErrPermintaanTidakSah",
				k.klaim, k.peserta, err)
		}
	}
}

// TestSetTanggalKejadianTanpaOracleGagal - tanpa basis data ia gagal terang,
// bukan diam-diam mengaku berhasil.
func TestSetTanggalKejadianTanpaOracleGagal(t *testing.T) {
	svc := services.New(nil)
	err := svc.TanggalKejadian().Set(context.Background(),
		pelakuBerperan(inti.PeranAdmin), "CLM-000001", "UJI-P-1",
		saat(t, "2025-06-01 00:00:00"))
	if !errors.Is(err, db.ErrTanpaOracle) {
		t.Fatalf("galat = %v, mau ErrTanpaOracle", err)
	}
}

// TestSetTanggalKejadianHanyaAdmin - butir bj (28-09-2026): gerbang DOL SAMA
// dengan tiga tanggal lainnya (b1000 `pyPosition!='ReasLifeAdmin'` baca-saja).
func TestSetTanggalKejadianHanyaAdmin(t *testing.T) {
	svc := services.New(nil)
	for _, pelaku := range []inti.Pelaku{
		{AkunID: "UJI-AKUN"},
		pelakuBerperan(inti.PeranMedicalAdvisor),
		pelakuBerperan(inti.PeranSPV),
	} {
		err := svc.TanggalKejadian().Set(context.Background(), pelaku, "CLM-000001", "UJI-P-1",
			saat(t, "2025-06-01 00:00:00"))
		if !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%v: galat = %v, mau ErrTanpaWewenang", pelaku.Peran, err)
		}
	}
}

// TestPergeseranDOLSatuHariBukanSatuJam - ralat A0 atas `[dugaan]` tiket 06.
//
// ⛔ Tanda tangan `@addCalendar` DITENTUKAN oleh korpus, bukan dicocokkan dari
// cacah argumen. Jangkarnya `Claim Life/Activity/LoadDataPeserta_Act.xml`
// pecahan baris 1097, 1143, 1163, 1183, 1203, 1223, 1243, dan 1263 - kedelapan
// tanggal peserta dimuat lewat `@addCalendar(<tanggal>,0,0,0,0,7,0,0)`.
//
// Angka **7** di posisi kelima tidak mungkin menit maupun minggu: tidak ada
// yang menggeser tanggal LAHIR tujuh menit, dan tujuh minggu akan memindahkan
// orang ke bulan lain. Tujuh JAM adalah WIB (UTC+7) - normalisasi zona yang
// wajar untuk data Indonesia. Maka posisi kelima = JAM, dan posisi keempat =
// HARI, sesuai tanda tangan baku
// `(tanggal, tahun, bulan, minggu, hari, jam, menit, detik)`.
//
// Akibatnya `ValidasiDOL_Act` baris 698 `(.DATE_OF_LOSS,0,0,0,1,0,0,0)`
// menggeser DOL **satu HARI**, bukan satu jam seperti dugaan ronde pertama -
// yang menaruh jam di posisi keempat atas dasar kecocokan cacah argumen saja.
func TestPergeseranDOLSatuHariBukanSatuJam(t *testing.T) {
	// Keputusan work owner 11-10-2026: pergeseran retro b698 DIBUANG - DOL = RETRO BEGIN lolos, sehari sebelumnya ditolak.
	if err := services.ValidasiDOL("TR", saat(t, retroMulai), pesertaDOL(t)); err != nil {
		t.Errorf("RETRO BEGIN ditolak: %v", err)
	}
	if err := services.ValidasiDOL("TR", saat(t, "2025-03-31 00:00:00"), pesertaDOL(t)); !errors.Is(err, services.ErrDOLTidakSah) {
		t.Errorf("RETRO BEGIN-1 lolos: %v (pergeseran masih ada?)", err)
	}
}

// --- Tiga tanggal klaim dialog Edit Date (sensus §3.1, 28-09-2026) ---

func tanggalKlaimUji(t *testing.T) models.TanggalKlaim {
	t.Helper()
	w := saat(t, "2026-03-02 00:00:00")
	return models.TanggalKlaim{TerimaKlaim: &w}
}

// TestSetTanggalKlaimMenjagaPagarnya - identitas, pengenal, lalu peran, SEBELUM
// basis data. `pyPosition=='ReasLifeAdmin'` (b1000/b1313/b1550/b1788) berarti
// hanya Admin yang dapat membuka isiannya.
func TestSetTanggalKlaimMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	tgl := tanggalKlaimUji(t)
	admin := pelakuBerperan(inti.PeranAdmin)

	if err := svc.TanggalKejadian().SetTanggalKlaim(ctx, inti.Pelaku{},
		"CLM-1", "P-1", tgl); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v, mau ErrTanpaIdentitas", err)
	}
	for _, k := range []struct{ klaim, peserta string }{{"", "P-1"}, {"CLM-1", " "}} {
		if err := svc.TanggalKejadian().SetTanggalKlaim(ctx, admin,
			k.klaim, k.peserta, tgl); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%q/%q: %v, mau ErrPermintaanTidakSah", k.klaim, k.peserta, err)
		}
	}
	for _, peran := range []string{inti.PeranMedicalAdvisor, inti.PeranSPV} {
		if err := svc.TanggalKejadian().SetTanggalKlaim(ctx, pelakuBerperan(peran),
			"CLM-1", "P-1", tgl); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s: %v, mau ErrTanpaWewenang", peran, err)
		}
	}
	if err := svc.TanggalKejadian().SetTanggalKlaim(ctx, admin,
		"CLM-1", "P-1", tgl); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("admin tanpa Oracle: %v, mau ErrTanpaOracle", err)
	}
}

// OQ-M1 DITUTUP (GILIRAN-17): gerbang bersama kedua rute dialog Edit Date
// memakai penanda turunan "sudah Save to RNM" - satu tempat, bukan dua.
func TestGerbangDialogTanggalMemakaiPenandaSaveRNM(t *testing.T) {
	isi, err := os.ReadFile("dol.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	i := strings.Index(s, "func gerbangTahapDialogTanggal(")
	if i < 0 {
		t.Fatal("gerbangTahapDialogTanggal tidak ada")
	}
	badan := s[i:]
	if j := strings.Index(badan[1:], "\nfunc "); j >= 0 {
		badan = badan[:j+1]
	}
	for _, mau := range []string{".SudahSaveRNM(ctx, klaimID)", "models.BolehUbahTanggalKlaim(", "ErrTanggalTerkunciSesudahSaveRNM"} {
		if !strings.Contains(badan, mau) {
			t.Errorf("gerbang tanpa %q", mau)
		}
	}
	if strings.Count(s, "gerbangTahapDialogTanggal(ctx, baca, klaimID)") != 2 {
		t.Error("kedua rute (DOL dan tiga tanggal) harus melewati gerbang yang sama")
	}
}

// Keputusan work owner 11-10-2026: DOL di luar jendela ditampilkan dengan kalimat jelas berisi rentang valuasi jendela
// yang dipakai (gross QR/QP, retro TR/TP) dan nomor peserta; "Invalid DOL" tetap kode (Error()).
func TestKalimatDOLMenyebutRentangDanPeserta(t *testing.T) {
	for _, u := range []struct{ tipe, dol, rentang string }{
		{"QR", "2025-09-02 00:00:00", "01/03/2025 – 01/09/2025"},
		{"TP", "2025-03-31 00:00:00", "01/04/2025 – 01/10/2025"},
	} {
		err := services.BeriNomorDOL(services.ValidasiDOL(u.tipe, saat(t, u.dol), pesertaDOL(t)), 2)
		var g services.GalatDOL
		if !errors.As(err, &g) {
			t.Fatalf("%s: %v", u.tipe, err)
		}
		mau := "Insured number 2: Date of Loss must be within the valuation period (" + u.rentang + ")."
		if g.Kalimat() != mau || g.Error() != "Invalid DOL" || !errors.Is(err, services.ErrDOLTidakSah) {
			t.Errorf("%s: kalimat %q kode %q", u.tipe, g.Kalimat(), g.Error())
		}
	}
	if services.KalimatDOLDiLuar("18/03/2025", "18/03/2026") != "Date of Loss must be within the valuation period (18/03/2025 – 18/03/2026)." {
		t.Error("kalimat dasar berubah")
	}
	// Galat lain tidak diubah BeriNomorDOL.
	if err := services.BeriNomorDOL(services.ErrDOLKosong, 1); !errors.Is(err, services.ErrDOLKosong) {
		t.Errorf("galat lain: %v", err)
	}
}
