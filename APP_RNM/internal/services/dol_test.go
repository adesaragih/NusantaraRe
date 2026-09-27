package services_test

// Validasi Date of Loss dan penurunan jenis klaim - TANPA Oracle.
//
// Pemilik: tiket 06. Dibaca sesudah: dol.go.

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
	"nusantarare/pkg/utils"
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
// Dua asumsi pustaka Pega menentukan seluruh tabel ini, dan keduanya `[dugaan]`
// yang ditulis sebagai konstanta bernama di dol.go:
//
//	(i)  @CompareDates(a,b) benar bila a SESUDAH b (ketat)
//	(ii) argumen keempat @addCalendar adalah JAM
//
// Akibatnya jendelanya ASIMETRIS - QP/QR terbuka di awal dan tertutup di
// akhir, TP/TR sebaliknya. Daftar LENGKAP kasus yang berbalik bila salah satu
// asumsi diputuskan lain ada di bab "Pembacaan ulang XML" tiket 06; ia ditulis
// di SATU tempat supaya tiga salinan tidak berjalan sendiri-sendiri.
func TestBatasJendelaTiapType(t *testing.T) {
	p := pesertaDOL(t)
	kasus := []struct {
		tipe string
		dol  string
		sah  bool
		apa  string
	}{
		// QP/QR - jendela gross, tanpa pergeseran. Keempat batas, kedua Type.
		{"QR", grossMulai, false, "tepat di tanggal mulai: DOL tidak SESUDAH mulai"},
		{"QR", "2025-03-01 01:00:00", true, "satu jam sesudah mulai"},
		{"QR", "2025-08-31 23:00:00", true, "satu jam sebelum berakhir"},
		{"QR", grossSelesai, true, "tepat di tanggal berakhir: belum melewatinya"},
		{"QR", "2025-09-01 01:00:00", false, "satu jam sesudah berakhir"},
		{"QP", grossMulai, false, "tepat di tanggal mulai"},
		{"QP", "2025-03-01 01:00:00", true, "satu jam sesudah mulai"},
		{"QP", "2025-08-31 23:00:00", true, "satu jam sebelum berakhir"},
		{"QP", grossSelesai, true, "tepat di tanggal berakhir"},
		{"QP", "2025-09-01 01:00:00", false, "satu jam sesudah berakhir"},

		// TP/TR - jendela retro, DOL digeser SATU HARI sebelum dibandingkan.
		//
		// ⛔ RALAT A0: kasus-kasus ini semula memakai pergeseran satu JAM, dan
		// batas-batasnya ikut berubah ketika bacaan korpus diralat. Inilah
		// "daftar kasus yang berbalik" yang komentar tiket 06 janjikan.
		//
		// Dengan geser +1 hari, DOL d dibandingkan sebagai d+1hari:
		//   mulai retro 2025-04-01, berakhir 2025-10-01.
		{"TR", "2025-03-30 00:00:00", false, "dua hari sebelum mulai: d+1 masih sebelum mulai"},
		{"TR", "2025-03-31 12:00:00", true, "sehari sebelum mulai: d+1 melewati mulai"},
		{"TR", retroMulai, true, "tepat di tanggal mulai: pergeseran membuatnya masuk"},
		{"TR", "2025-09-29 00:00:00", true, "dua hari sebelum berakhir: d+1 masih di dalam"},
		{"TR", "2025-09-30 12:00:00", false, "sehari sebelum berakhir: d+1 melewati akhir"},
		{"TR", retroSelesai, false, "tepat di tanggal berakhir: pergeseran melewatinya"},
		{"TP", "2025-03-30 00:00:00", false, "dua hari sebelum mulai"},
		{"TP", retroMulai, true, "tepat di tanggal mulai"},
		{"TP", "2025-09-29 00:00:00", true, "dua hari sebelum berakhir"},
		{"TP", retroSelesai, false, "tepat di tanggal berakhir"},

		// Cabang yang benar-benar terpisah: tanggal yang sah di satu jendela
		// dan TIDAK sah di jendela yang lain.
		{"QR", "2025-09-15 00:00:00", false, "sesudah gross berakhir, tetapi masih di dalam retro"},
		{"TR", "2025-09-15 00:00:00", true, "sesudah gross berakhir, tetapi masih di dalam retro"},
		{"QR", "2025-03-15 00:00:00", true, "di dalam gross, tetapi sebelum retro mulai"},
		{"TR", "2025-03-15 00:00:00", false, "di dalam gross, tetapi sebelum retro mulai"},
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
	err := services.ValidasiDOL("QR", saat(t, grossMulai), pesertaDOL(t))
	if err == nil {
		t.Fatal("tanggal di batas mulai lolos")
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
	err := services.ValidasiDOL("QR", saat(t, grossMulai), pesertaDOL(t))
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
	err := svc.TanggalKejadian().Set(context.Background(), services.Pelaku{},
		"CLM-000001", "UJI-P-1", saat(t, "2025-06-01 00:00:00"))
	if !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Fatalf("galat = %v, mau ErrTanpaIdentitas", err)
	}
}

// TestSetTanggalKejadianTanpaPengenalDitolak - pengenal yang kosong tidak
// pernah diteruskan ke SQL sebagai teks kosong.
func TestSetTanggalKejadianTanpaPengenalDitolak(t *testing.T) {
	svc := services.New(nil)
	pelaku := services.Pelaku{AkunID: "UJI-AKUN"}
	for _, k := range []struct{ klaim, peserta string }{
		{"", "UJI-P-1"}, {"CLM-000001", ""}, {"   ", "   "},
	} {
		err := svc.TanggalKejadian().Set(context.Background(), pelaku,
			k.klaim, k.peserta, saat(t, "2025-06-01 00:00:00"))
		if !errors.Is(err, services.ErrPermintaanTidakSah) {
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
		services.Pelaku{AkunID: "UJI-AKUN"}, "CLM-000001", "UJI-P-1",
		saat(t, "2025-06-01 00:00:00"))
	if !errors.Is(err, repository.ErrTanpaOracle) {
		t.Fatalf("galat = %v, mau ErrTanpaOracle", err)
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
	const seharusnya = 24 * time.Hour
	if services.PergeseranDOLRetro() != seharusnya {
		t.Errorf("pergeseran = %v, mau %v (satu hari). Jangkarnya angka 7 = WIB "+
			"di LoadDataPeserta_Act, yang menetapkan posisi kelima sebagai jam "+
			"dan posisi keempat sebagai hari",
			services.PergeseranDOLRetro(), seharusnya)
	}
}
