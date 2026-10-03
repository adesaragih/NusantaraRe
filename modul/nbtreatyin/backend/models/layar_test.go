package models

// Uji aturan layar, bentuk simpan, dan pilih bisnis - nilai harapan dari
// section/aktivitas Pega (INVENTARIS bab 5-7) dan spec-penyimpanan.

import (
	"errors"
	"testing"
	"time"
)

func TestBentukProporsionalMenolakXOLDanRincian(t *testing.T) { // spec-penyimpanan AC 31, 33
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	if err := PeriksaBentukSimpan(h); err != nil {
		t.Fatal(err)
	}
	h.SetelDaftar(TabelXOL.Daftar, []Baris{{"Currency": "UJI"}})
	if err := PeriksaBentukSimpan(h); !errors.Is(err, ErrBentukTidakSah) {
		t.Fatalf("XOL pada proporsional: %v", err)
	}
	h.SetelDaftar(TabelXOL.Daftar, nil)
	h.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1"}})
	h.SetelDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList"), []Baris{{"InstallmentNo": "1"}})
	if err := PeriksaBentukSimpan(h); !errors.Is(err, ErrBentukTidakSah) {
		t.Fatalf("rincian angsuran pada proporsional: %v", err)
	}
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisNonProporsional)
	if err := PeriksaBentukSimpan(h); err != nil {
		t.Fatalf("non-proporsional boleh bersarang: %v", err)
	}
}

func TestSpreadingBagiRataPresisiSepuluh(t *testing.T) { // spec-penyimpanan AC 36
	h := HalamanBaru()
	h.SetelDaftar(DaftarSpreading, []Baris{{}, {}, {}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	if got := h.AmbilDaftar(DaftarSpreading)[0]["SharePercentage"]; got != "33.3333333333" {
		t.Fatalf("100/3 presisi 10: %q", got)
	}
}

func TestMedanWajibAtasanTanpaEnamMedanAdmin(t *testing.T) { // AC 46, 47
	jalur := map[string]bool{}
	for _, m := range DaftarMedanWajib(PosisiDeptHead) {
		jalur[m.Jalur] = true
	}
	for _, m := range []string{"ClaimPaymentType", "ClaimType", "IDCurrency", "Quartal", "TypeTax", "YearOfQuartal"} {
		if jalur["PolicyTreatyIn."+m] {
			t.Errorf("%s tidak wajib di layar Dept Head", m)
		}
	}
	if !jalur["PolicyTreatyIn.ResultOnp1"] {
		t.Error("ResultOnp1 wajib di layar Dept Head")
	}
	// gabungan dua layar: 26 medan berbeda (+ DateofSurvey layar survei = 27, AC 45)
	semua := map[string]bool{}
	for _, p := range PosisiTangga {
		for _, m := range DaftarMedanWajib(p) {
			semua[m.Jalur] = true
		}
	}
	if len(semua) != 25 {
		t.Fatalf("medan wajib berbeda di layar realisasi: %d", len(semua))
	}
}

func TestMedanWajibBersyarat(t *testing.T) {
	h := HalamanBaru()
	h.Setel("Quotation.ProportionalType", JenisNonProporsional)
	kosong := MedanWajibKosong(h, PosisiAdmin)
	for _, l := range kosong {
		if l == "Premi Ogp" || l == "Type Tax" || l == "Claim Type" {
			t.Errorf("%s tidak wajib untuk NonProportional / FlagPPH kosong / Claim kosong", l)
		}
	}
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	h.Setel("PolicyTreatyIn.Claim", "5")
	kosong = MedanWajibKosong(h, PosisiAdmin)
	harap := map[string]bool{"Type Tax": false, "Claim Type": false, "Payment Type": false}
	for _, l := range kosong {
		if _, ada := harap[l]; ada {
			harap[l] = true
		}
	}
	for l, ada := range harap {
		if !ada {
			t.Errorf("%s wajib bila syaratnya terpenuhi", l)
		}
	}
}

func TestGabungMasukanAtasanHanyaMedanTerbuka(t *testing.T) { // AC 49-52
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.PremiOgp", "1000")
	h.Setel("PolicyTreatyIn.DueTo", "1")
	m := HalamanBaru()
	m.Setel("PolicyTreatyIn.PremiOgp", "1")
	m.Setel("PolicyTreatyIn.Suggest", "UJI")
	GabungMasukanLayar(h, m, PosisiSecHead)
	if h.Ambil("PolicyTreatyIn.PremiOgp") != "1000" || h.Ambil("PolicyTreatyIn.Suggest") != "UJI" {
		t.Fatal("atasan: medan terkunci diabaikan, Suggest diterima")
	}
	if h.Ambil("PolicyTreatyIn.DueTo") != "1" {
		t.Fatal("medan yang tidak dikirim tidak dikosongkan")
	}
	GabungMasukanLayar(h, m, PosisiAdmin)
	if h.Ambil("PolicyTreatyIn.PremiOgp") != "1" {
		t.Fatal("admin boleh mengubah PremiOgp")
	}
	m.Setel("PolicyTreatyIn.StatementDate", "2000-01-01 00:00:00")
	GabungMasukanLayar(h, m, PosisiAdmin)
	if h.Ambil("PolicyTreatyIn.StatementDate") != "" {
		t.Fatal("StatementDate terkunci ALWAYS di layar admin")
	}
}

func TestPascaAdminTolakNBStatus(t *testing.T) { // AC 44, 71
	h := HalamanBaru()
	h.Setel("PositionNote", PosisiAdmin)
	h.Setel("PolicyTreatyIn.IsApproved", "0")
	h.Setel("PolicyTreatyIn.Suggest", "UJI-alasan")
	h.Setel("PolicyTreatyIn.SuggestDate", "2026-10-03 09:00:00")
	h.Setel("PolicyTreatyIn.OperatorName", "Uji Nama")
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "10218"}})
	PascaAdmin(h, "UJI-AKUN", "Uji Nama")
	if h.Ambil("NBStatus") != "NB WAS DECLINED BY  UJI NAMA" || h.Ambil("PolicyTreatyIn.HasFacOut") != "1" {
		t.Fatalf("NBStatus %q HasFacOut %q", h.Ambil("NBStatus"), h.Ambil("PolicyTreatyIn.HasFacOut"))
	}
	d := h.AmbilDaftar(DaftarUsulan)
	if len(d) != 1 || d[0]["IsApproved"] != "0" || d[0]["OperatorID"] != "UJI-AKUN" || d[0]["Date"] == "" {
		t.Fatalf("catatan %+v", d)
	}
	// Approved kosong: IsApproved baris tidak diisi (WHEN Param.Approved != "")
	h.Setel("PolicyTreatyIn.IsApproved", "")
	TambahCatatan(h, "UJI-AKUN")
	if _, ada := h.AmbilDaftar(DaftarUsulan)[1]["IsApproved"]; ada {
		t.Fatal("IsApproved kosong tidak ditulis ke baris catatan")
	}
}

func TestTombolPerPosisi(t *testing.T) {
	h := HalamanBaru()
	for _, tt := range []struct {
		posisi, ia string
		harap      TombolKirim
	}{
		{PosisiAdmin, "1", TombolKirimLangsung},
		{PosisiAdmin, "0", TombolKonfirmasiTolak},
		{PosisiAdmin, "", TombolTidakAda},
		{PosisiSecHead, "1", TombolKirimLangsung},
		{PosisiSecHead, "0", TombolKirimLangsung},
		{PosisiDeptHead, "1", TombolNomorPolis},
		{PosisiDeptHead, "0", TombolKirimLangsung},
	} {
		h.Setel("PolicyTreatyIn.IsApproved", tt.ia)
		if got := TombolUntuk(h, tt.posisi); got != tt.harap {
			t.Errorf("%s IsApproved %q: %q, harap %q", tt.posisi, tt.ia, got, tt.harap)
		}
	}
}

func TestKunciCariBisnis(t *testing.T) {
	for _, tt := range []struct {
		nama  string
		pra   bool
		harap string
	}{
		{"UJI MBU CAR", false, "MOTOR VEHICLE"},
		{"PERFORMANCE BONDS", false, "PERFORMANCE BOND"},
		{"BID OR TENDER BONDS", false, "BID OR TENDER BOND"},
		{"BID OR TENDER BONDS", true, "BID BOND"},
		{"CUSTOMS BOND", true, "OTHERS CUSTOMS BOND"},
		{"GOLF INSURANCE", false, "HOLE IN ONE"},
		{"UJI LAIN", false, "UJI LAIN"},
	} {
		if got := KunciCariBisnis(tt.nama, tt.pra); got != tt.harap {
			t.Errorf("%q pra=%v: %q, harap %q", tt.nama, tt.pra, got, tt.harap)
		}
	}
}

func TestTanggalProduksiDariHariClosing(t *testing.T) {
	h := HalamanBaru()
	sekarang := time.Date(2026, 10, 27, 9, 0, 0, 0, time.UTC)
	PraprosesTanggal(h, sekarang, 25)
	// Pre_Act langkah 9: `@substring(ProductionDate,8)` membawa jam StatementDate
	// (RALAT putaran 2: semula 00:00:00).
	if h.Ambil("PolicyTreatyIn.StatementDate") != "2026-10-27 09:00:00" ||
		h.Ambil("PolicyTreatyIn.ProductionDate") != "2026-11-01 09:00:00" {
		t.Fatalf("hari 27 > closing 25: %q / %q", h.Ambil("PolicyTreatyIn.StatementDate"), h.Ambil("PolicyTreatyIn.ProductionDate"))
	}
	PraprosesTanggal(h, sekarang, 28)
	if h.Ambil("PolicyTreatyIn.ProductionDate") != "2026-10-27 09:00:00" {
		t.Fatal("hari 27 <= closing 28: tetap")
	}
}

func TestEnableDisableSelaluNol(t *testing.T) {
	for _, awal := range []string{"", "0", "1"} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", awal)
		TreatyEnableDisableInput(h)
		if h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "0" {
			t.Errorf("awal %q -> %q; urutan WHEN DT berakhir di 0", awal, h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp"))
		}
	}
}
