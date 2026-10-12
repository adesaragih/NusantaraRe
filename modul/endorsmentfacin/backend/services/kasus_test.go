package services

// Uji Layanan.Buat (irisan I3) di atas Gudang tiruan bertransaksi. Data UJI buatan.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
	"nusantarare/modul/endorsmentfacin/backend/tiruan"
)

const akunUji = "uji.akun"

func gudangKasus() *tiruan.Gudang {
	return &tiruan.Gudang{JenisBisnisNilai: "Fire",
		Periode: &repository.PeriodeTeks{Mulai: "20250101120000", Akhir: "20260101120000"}}
}

func masukanKasus() models.MasukanKasus {
	return models.MasukanKasus{
		MasukanPeriksa: models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "1",
			EndorsementDate: "01-06-2025", EndorsementInternalRetro: "0"},
		Note: "Catatan uji", EndorsementSource: "Ceding", EndorsementSourceNote: "Surat uji", DetailType: "2",
		EndorsementStatus: "0",
	}
}

func layananKasus(g *tiruan.Gudang) *Layanan { return BaruLayanan(g).DenganTransaksi(g.Transaksi) }

func TestBuatKasusLahir(t *testing.T) {
	g := gudangKasus()
	h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "EDM-00001" || h.Periksa != nil {
		t.Fatalf("hasil = %+v", h)
	}
	if g.WorkPolis["EDM-00001"] != akunUji {
		t.Fatalf("T_WORK_POLIS = %v", g.WorkPolis)
	}
	want := models.KepalaKasus{ID: "EDM-00001", OldPolicyNo: polisUji, EdmType: "4", EdmTypeNew: "1", DetailType: "2",
		EndorsementReason: "Catatan uji", EndorsementSource: "Ceding", EndorsementSourceNote: "Surat uji",
		// AdminCharge / Survey kosong → NULL (A50).
		EndorsementStatus: "0", AdminCharge: "", Survey: "", InternalRetro: "0", NBStatus: "NEW EDM",
		IsCedingConfirm: "Policy"}
	got := g.Kasus["EDM-00001"]
	tanggal := got.EndorsementDate
	got.EndorsementDate = time.Time{}
	if got != want {
		t.Fatalf("kepala =\n%+v\ningin\n%+v", got, want)
	}
	if tanggal.Format("2006-01-02") != "2025-06-01" {
		t.Fatalf("tanggal = %v", tanggal)
	}
	// Kasus kedua: nomor berikutnya.
	if h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus()); err != nil || h.CaseID != "EDM-00002" {
		t.Fatalf("kasus kedua = %+v, %v", h, err)
	}
}

func TestBuatKasusTanpaAkun(t *testing.T) {
	g := gudangKasus()
	_, err := layananKasus(g).Buat(context.Background(), " ", masukanKasus())
	if !errors.Is(err, ErrDilarang) || !strings.Contains(err.Error(), "Log in to create an endorsement") {
		t.Fatalf("galat = %v", err)
	}
	if len(g.Panggilan) != 0 {
		t.Fatalf("query dijalankan tanpa akun: %v", g.Panggilan)
	}
}

func TestBuatKasusTanpaTransaksi(t *testing.T) {
	if _, err := BaruLayanan(gudangKasus()).Buat(context.Background(), akunUji, masukanKasus()); !errors.Is(err, ErrTanpaTransaksi) {
		t.Fatalf("galat = %v", err)
	}
}

// Wajib bersyarat bukan-Life (pyRequiredWhen CARI2!='Life'): kosong lolos ValidasiKasus, ditagih sesudah pemeriksaan.
func TestValidasiSesudahPeriksa(t *testing.T) {
	for _, k := range []struct {
		ubah  func(*models.MasukanKasus)
		jenis string
		teks  string // "" = lolos
	}{
		{func(m *models.MasukanKasus) { m.EndorsementStatus = " " }, "Fire", "Endorsement Status is required"},
		{func(m *models.MasukanKasus) { m.EndorsementInternalRetro = "" }, "Fire", "Endorsement Internal Retro is required"},
		// Urut layar: Detail EDM Type (Section EdmType1) lalu Status (L8906) lalu Retro (L9968).
		{func(m *models.MasukanKasus) { m.EndorsementStatus, m.EndorsementInternalRetro = "", "" }, "Fire", "Endorsement Status is required"},
		{func(m *models.MasukanKasus) { m.DetailType, m.EndorsementStatus = "", "" }, "Fire", "Detail EDM Type is required"},
		// Life: ketiganya tersembunyi, dikirim "".
		{func(m *models.MasukanKasus) {
			m.DetailType, m.EndorsementStatus, m.EndorsementInternalRetro = "", "", ""
		}, "Life", ""},
		{func(*models.MasukanKasus) {}, "Fire", ""},
	} {
		m := masukanKasus()
		k.ubah(&m)
		if err := ValidasiKasus(m); err != nil {
			t.Fatalf("%q: ValidasiKasus menolak yang kosong: %v", k.teks, err)
		}
		err := validasiSesudahPeriksa(m, k.jenis)
		if k.teks == "" {
			if err != nil {
				t.Errorf("%s: galat = %v", k.jenis, err)
			}
			continue
		}
		if !errors.Is(err, ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.teks) {
			t.Errorf("ingin %q, galat = %v", k.teks, err)
		}
	}
}

// Polis bukan Life tanpa Status / Retro: 422 sesudah pemeriksaan, nol tulisan.
func TestBuatKasusBukanLifeWajibStatusRetro(t *testing.T) {
	for _, k := range []struct {
		status, retro, teks string
	}{
		{"", "0", "Endorsement Status is required"},
		{"0", "", "Endorsement Internal Retro is required"},
	} {
		g := gudangKasus()
		m := masukanKasus()
		m.EndorsementStatus, m.EndorsementInternalRetro = k.status, k.retro
		_, err := layananKasus(g).Buat(context.Background(), akunUji, m)
		if !errors.Is(err, ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.teks) {
			t.Errorf("ingin %q, galat = %v", k.teks, err)
		}
		tanpaTulisan(t, g)
	}
}

// Polis Life dengan Status/Retro "" (layar menyembunyikannya): jawabannya penolakan bisnis, bukan 422.
func TestBuatKasusLifeTanpaStatusRetro(t *testing.T) {
	g := gudangKasus()
	m := masukanKasus()
	m.PolicyNo, m.DetailType, m.EndorsementStatus, m.EndorsementInternalRetro = polisUjiLife, "", "", ""
	h, err := layananKasus(g).Buat(context.Background(), akunUji, m)
	if err != nil {
		t.Fatalf("galat = %v", err)
	}
	if h.CaseID != "" || h.Periksa == nil || h.Periksa.Bisnis.JenisBisnis != "Life" || h.Periksa.Bisnis.Hasil != T {
		t.Fatalf("hasil = %+v", h.Periksa)
	}
	tanpaTulisan(t, g)
}

func TestValidasiKasus(t *testing.T) {
	for _, k := range []struct {
		ubah func(*models.MasukanKasus)
		teks string
	}{
		{func(m *models.MasukanKasus) { m.PolicyNo = " " }, "Policy Number is required"},
		{func(m *models.MasukanKasus) { m.EndorsementSourceNote = "" }, "Endorsement Source is required"},
		{func(m *models.MasukanKasus) { m.Note = "" }, "Endorsement Note is required"},
		{func(m *models.MasukanKasus) { m.EdmType = "" }, "Endorsement Type is required"},
		{func(m *models.MasukanKasus) { m.EdmTypeNew = "" }, "EDM Type is required"},
		{func(m *models.MasukanKasus) { m.EndorsementDate = "" }, "Endorsement Date is required"},
		// Prompt List (work owner 11-10-2026), diperiksa bila terisi.
		{func(m *models.MasukanKasus) { m.EndorsementStatus = "4" }, `Endorsement Status "4" is not an option`},
		{func(m *models.MasukanKasus) { m.EndorsementStatus = "Internal" }, `Endorsement Status "Internal" is not an option`},
		{func(m *models.MasukanKasus) { m.EndorsementInternalRetro = "2" }, `Endorsement Internal Retro "2" is not an option`},
		{func(m *models.MasukanKasus) { m.EndorsementInternalRetro = "Yes" }, `Endorsement Internal Retro "Yes" is not an option`},
		// Medan pertama yang gagal menang.
		{func(m *models.MasukanKasus) { m.Note, m.EndorsementDate = "", "" }, "Endorsement Note is required"},
		{func(m *models.MasukanKasus) { m.Note = strings.Repeat("a", 4001) }, "Endorsement Note is too long (4001 bytes, maximum 4000)"},
		{func(m *models.MasukanKasus) { m.EndorsementSourceNote = strings.Repeat("é", 501) }, "Endorsement Source is too long (1002 bytes, maximum 1000)"},
		{func(m *models.MasukanKasus) { m.PolicyNo = strings.Repeat("9", 101) }, "Policy Number is too long"},
	} {
		m := masukanKasus()
		k.ubah(&m)
		err := ValidasiKasus(m)
		if !errors.Is(err, ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.teks) {
			t.Errorf("ingin %q, galat = %v", k.teks, err)
		}
	}
	// Tepat di batas: lolos. EdmTypeNew hanya wajib untuk EdmType 4.
	m := masukanKasus()
	m.Note, m.EndorsementSourceNote = strings.Repeat("a", 4000), strings.Repeat("é", 500)
	if err := ValidasiKasus(m); err != nil {
		t.Fatalf("batas: %v", err)
	}
	m = masukanKasus()
	m.EdmType, m.EdmTypeNew = "1", ""
	if err := ValidasiKasus(m); err != nil {
		t.Fatalf("EdmType 1 tanpa EdmTypeNew: %v", err)
	}
	// Setiap pilihan Prompt List diterima.
	for _, s := range []string{"0", "1", "2", "3"} {
		for _, r := range []string{"0", "1"} {
			m = masukanKasus()
			m.EndorsementStatus, m.EndorsementInternalRetro = s, r
			if err := ValidasiKasus(m); err != nil {
				t.Fatalf("status %s retro %s: %v", s, r, err)
			}
		}
	}
}

func TestValidasiDetailType(t *testing.T) {
	for _, k := range []struct {
		edm, baru, detail, jenis string
		lolos                    bool
	}{
		{"4", "1", "12", "Fire", true},
		{"4", "1", "9", "Fire", false},  // 9 hanya pilihan perubahan
		{"4", "2", "12", "Fire", false}, // K-046: SetEdmType menaruh 12 di TypePenambahan
		{"4", "2", "11", "Fire", true},
		{"4", "3", "0", "Fire", true},
		{"4", "3", "", "Fire", false}, // wajib
		{"4", "4", "", "Fire", true},  // RI slip: tanpa Detail Type
		{"4", "4", "2", "Fire", false},
		{"1", "", "2", "Fire", false}, // bukan EdmType 4: wajib kosong
		{"1", "", "", "Fire", true},
		{"4", "1", "", "Life", true}, // pyRequiredWhen CARI2!='Life'
	} {
		m := masukanKasus()
		m.EdmType, m.EdmTypeNew, m.DetailType = k.edm, k.baru, k.detail
		err := validasiDetailType(m, k.jenis)
		if (err == nil) != k.lolos || (err != nil && !errors.Is(err, ErrMasukanTidakSah)) {
			t.Errorf("%+v: galat = %v", k, err)
		}
	}
}

// tanpaTulisan - penolakan: nol nomor diambil, nol baris.
func tanpaTulisan(t *testing.T, g *tiruan.Gudang) {
	t.Helper()
	if g.Dipanggil("NomorKasusBerikut") || len(g.WorkPolis) != 0 || len(g.Kasus) != 0 {
		t.Fatalf("ada tulisan: panggilan %v", g.Panggilan)
	}
}

func TestBuatKasusDitolakPeriksa(t *testing.T) {
	g := gudangKasus()
	g.Renewal = true
	h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "" || h.Periksa == nil || !h.Periksa.Ditolak || h.Periksa.Gerbang[4].Hasil != T {
		t.Fatalf("hasil = %+v", h)
	}
	tanpaTulisan(t, g)
}

func TestBuatKasusTidakLengkap(t *testing.T) {
	g := gudangKasus()
	g.GalatPembayaran = fmt.Errorf("%w: ORA-01031", repository.ErrSkemaLainTakTerbaca)
	m := masukanKasus()
	m.EdmType, m.EdmTypeNew, m.DetailType = "1", "", ""
	h, err := layananKasus(g).Buat(context.Background(), akunUji, m)
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "" || h.Periksa == nil || h.Periksa.Lengkap || h.Periksa.Ditolak {
		t.Fatalf("hasil = %+v", h.Periksa)
	}
	tanpaTulisan(t, g)
}

func TestBuatKasusDetailTypeSesudahPeriksa(t *testing.T) {
	g := gudangKasus()
	m := masukanKasus()
	m.DetailType = "9"
	if _, err := layananKasus(g).Buat(context.Background(), akunUji, m); !errors.Is(err, ErrMasukanTidakSah) {
		t.Fatalf("galat = %v", err)
	}
	tanpaTulisan(t, g)
}

func TestBuatKasusOpenCase(t *testing.T) {
	// EdmType polis lama 1, klep Batal membuka gerbang 1 → OpenCase langkah 5-6 tetap menolak.
	g := gudangKasus()
	g.StatusEDMNilai, g.Proteksi = "1", map[int]bool{1: true}
	h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "" || h.Periksa == nil || !h.Periksa.Ditolak {
		t.Fatalf("hasil = %+v", h)
	}
	if got := teksPesan(*h.Periksa); !samaIrisan(got, []string{"PolicyNo|Sudah Di endorsement Batal"}) {
		t.Fatalf("pesan = %v", got)
	}
	tanpaTulisan(t, g)

	// Polis di daftar BATAL (tabel 261): kasus lahir, pesan langkah 4 tetap dipasang.
	g = gudangKasus()
	g.StatusEDMNilai, g.Proteksi = "1", map[int]bool{1: true}
	g.PengecualianPer = map[string][]string{JenisPengecualianBatal: {polisUji}}
	h, err = layananKasus(g).Buat(context.Background(), akunUji, masukanKasus())
	if err != nil {
		t.Fatal(err)
	}
	if h.CaseID != "EDM-00001" || len(h.Pesan) != 1 || h.Pesan[0].Teks != "Sudah Di endorsement Batal" {
		t.Fatalf("hasil = %+v", h)
	}
	// Tabel 261 belum ada: daftar kosong.
	g = gudangKasus()
	g.StatusEDMNilai, g.Proteksi = "1", map[int]bool{1: true}
	g.GalatPengecualian = fmt.Errorf("%w: ORA-00942", repository.ErrBelumAda)
	if h, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus()); err != nil || h.CaseID != "" {
		t.Fatalf("tabel 261 belum ada: %+v %v", h, err)
	}
}

// RI slip melompati G1 di /periksa, tetapi OpenCase (SetValueToEDMWork langkah 2) tetap membaca GetEDMStatus_SQL.
func TestBuatKasusRISlipMembacaStatusEDM(t *testing.T) {
	g := gudangKasus()
	g.StatusEDMNilai = "1"
	m := masukanKasus()
	m.EdmTypeNew, m.DetailType = "4", ""
	h, err := layananKasus(g).Buat(context.Background(), akunUji, m)
	if err != nil {
		t.Fatal(err)
	}
	if !g.Dipanggil("StatusEDM") || h.CaseID != "" || h.Periksa == nil {
		t.Fatalf("hasil = %+v panggilan %v", h, g.Panggilan)
	}
}

func TestBuatKasusRollback(t *testing.T) {
	boom := errors.New("ORA-12899")
	for _, set := range []func(*tiruan.Gudang){
		func(g *tiruan.Gudang) { g.GalatSisipKasus = boom },
		func(g *tiruan.Gudang) { g.GalatSisipWorkPolis = boom },
		func(g *tiruan.Gudang) { g.GalatNomor = boom },
	} {
		g := gudangKasus()
		set(g)
		if _, err := layananKasus(g).Buat(context.Background(), akunUji, masukanKasus()); !errors.Is(err, boom) {
			t.Fatalf("galat = %v", err)
		}
		if len(g.WorkPolis) != 0 || len(g.Kasus) != 0 {
			t.Fatalf("tulisan tersisa sesudah rollback: %v %v", g.WorkPolis, g.Kasus)
		}
	}
}

func TestPeriksaGerbangDuaKasusSistemBaru(t *testing.T) {
	g := gudangBersih()
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
	if err != nil {
		t.Fatal(err)
	}
	if h.Gerbang[1].Hasil != L || h.Gerbang[1].Catatan != catatanEDMPega || !h.Lengkap {
		t.Fatalf("tanpa kasus terbuka: %+v lengkap %v", h.Gerbang[1], h.Lengkap)
	}
	g = gudangBersih()
	g.EDMTerbuka = "EDM-00007"
	h, err = BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "4"})
	if err != nil {
		t.Fatal(err)
	}
	// Gerbang 2 juga berjalan untuk RI slip, tanpa klep.
	if h.Gerbang[1].Hasil != T || !samaIrisan(teksPesan(h), []string{"PolicyNo|There's EDM with this policy no that haven't finish yet! EDM-00007"}) {
		t.Fatalf("kasus terbuka: %+v %v", h.Gerbang[1], teksPesan(h))
	}
	g = gudangBersih()
	g.GalatEDMBelumSelesai = fmt.Errorf("%w: ORA-00942", repository.ErrBelumAda)
	h, err = BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
	if err != nil {
		t.Fatal(err)
	}
	if h.Gerbang[1].Hasil != BD || h.Gerbang[1].Catatan != catatanKasusBelumAda || h.Lengkap {
		t.Fatalf("tabel 262 belum ada: %+v", h.Gerbang[1])
	}
	g = gudangBersih()
	g.GalatEDMBelumSelesai = errors.New("ORA-03113")
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err == nil {
		t.Fatal("galat gerbang 2 ditelan")
	}
}
