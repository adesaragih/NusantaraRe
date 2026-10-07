package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/premium"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

type kasus map[string]string

func (k kasus) Nilai(j string) (string, bool) { v, ada := k[j]; return v, ada }

type limitTiruan struct {
	a   []models.BarisLimitA
	b   []models.BarisLimitB
	err error
}

func (l limitTiruan) MuatLimit(context.Context) ([]models.BarisLimitA, []models.BarisLimitB, error) {
	return l.a, l.b, l.err
}

func d(s string) *apd.Decimal {
	v, err := utils.ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return v
}

// limitPropertyTG1 - potongan M_LIMIT_PROPERTYY team group 1 (fixture tiket 11, tanpa NAMA/LOGIN).
var limitPropertyTG1 = []models.BarisLimitA{
	{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "SENIORUW", TeamGroup: "1", LimitBottom: d("178500000001"), LimitBottom2: d("178500000001")},
	{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "KADIVTEKNIK", TeamGroup: "1", LimitBottom: d("255000000001"), LimitBottom2: d("289000000001")},
	{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "MANAGERTEKNIK", TeamGroup: "1", LimitBottom: d("255000000001"), LimitBottom2: d("255000000001")},
	{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "KADIVFACULTATIVE", TeamGroup: "1", LimitBottom: d("255000000001"), LimitBottom2: d("255000000001")},
	{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "DIREKTURMARKETING", TeamGroup: "1", LimitBottom: d("331500000001"), LimitBottom2: d("331500000001")},
}

// fireTG1 - kasus FIRE Preferred Risk TG 1, TSI 300 M (sama dengan tes acceptance).
var fireTG1 = kasus{
	"pyWorkPage.OfferFacIn.TotalTSINusaRe": "300000000000", "pyWorkPage.OfferFacIn.TotalTSITopRisk": "0",
	"pyWorkPage.IsAdaTopRisk": "false", "pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "1",
	"pyWorkPage.OfferFacIn.QuotationData.TeamGroup": "1", "pyWorkPage.OfferFacIn.QuotationData.BusinessType": "Fire",
	"pyWorkPage.OfferFacIn.QuotationData.BusinessOldId": "22", "pyWorkPage.OfferFacIn.IsPreferredRisk": "Preferred Risk",
	"pyWorkPage.OfferFacIn.IsBanding": "false", "pyWorkPage.OfferFacIn.IsFlagReject": "false",
	"pyWorkPage.PositionNote": "ReasFacInSeniorUnderwriting",
}

// TestLangkahAkseptasiDariRepository - tabel limit datang dari repository; hasilnya
// sama dengan acceptance.Next atas fixture tiket 11.
func TestLangkahAkseptasiDariRepository(t *testing.T) {
	s := Baru(limitTiruan{a: limitPropertyTG1})
	got, err := s.LangkahAkseptasi(context.Background(), fireTG1, kontrak.PenggunaFacIn{Jabatan: "SENIORUW"})
	mau := kontrak.TransisiFacIn{JabatanTujuan: "KADIVTEKNIK", Antrean: "ReasFacInGroupLeader", PositionNoteDitulis: true}
	if err != nil || got != mau {
		t.Fatalf("dapat %+v (%v), mau %+v", got, err, mau)
	}
}

// TestLangkahAkseptasiDitolak - tabel tak dikenal, tabel kosong, galat repository,
// galat tangga - galat masukan bertanda ErrTidakDapatDiproses, galat repository tidak.
func TestLangkahAkseptasiDitolak(t *testing.T) {
	ctx := context.Background()
	p := kontrak.PenggunaFacIn{Jabatan: "SENIORUW"}
	asing := []models.BarisLimitA{{Tabel: "M_LIMIT_LAIN", Jabatan: "X", TeamGroup: "1", LimitBottom: d("1")}}
	if _, err := Baru(limitTiruan{a: asing}).LangkahAkseptasi(ctx, fireTG1, p); err == nil || errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("tabel asing: galat %v (bukan galat masukan)", err)
	}
	if _, err := Baru(limitTiruan{}).LangkahAkseptasi(ctx, fireTG1, p); !errors.Is(err, ErrTabelLimitTakTersedia) || errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("tabel kosong: galat %v (mau ErrTabelLimitTakTersedia)", err)
	}
	mati := errors.New("oracle mati")
	if _, err := Baru(limitTiruan{err: mati}).LangkahAkseptasi(ctx, fireTG1, p); !errors.Is(err, mati) || errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("repository: galat %v", err)
	}
	rusak := kasus{"pyWorkPage.OfferFacIn.TotalTSINusaRe": ""}
	if _, err := Baru(limitTiruan{a: limitPropertyTG1}).LangkahAkseptasi(ctx, rusak, p); !errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("kasus rusak: galat %v", err)
	}
}

// TestHitungPremi - premi satu coverage lewat mesin NB; lini tak dikenal dan masukan
// tak sah menjadi ErrTidakDapatDiproses, bukan panic.
func TestHitungPremi(t *testing.T) {
	s := Baru(limitTiruan{})
	got, err := s.HitungPremi(kontrak.MasukanPremiFacIn{LiniBisnis: "MARINE CARGO", MataUang: "IDR", TSI: "1000000", Rate: "0.125"})
	if err != nil || utils.FormatDecimal(got.Premi.Amount) != "1250.0000" || got.Premi.Currency != "IDR" ||
		got.AsalRumus != "CountGPWMarinePAMbu_Act langkah 1.1.1.1.1" {
		t.Fatalf("dapat %+v (%v)", got, err)
	}
	for nama, in := range map[string]kontrak.MasukanPremiFacIn{
		"lini tak dikenal": {LiniBisnis: "UJI-LINI", TSI: "1", Rate: "1"},
		"lini kosong":      {TSI: "1", Rate: "1"},
		"rate tak terbaca": {LiniBisnis: "MARINE CARGO", TSI: "1", Rate: "1,5"},
		"belum diport":     {LiniBisnis: "BONDING", TSI: "1", Rate: "1"},
	} {
		if _, err := s.HitungPremi(in); !errors.Is(err, ErrTidakDapatDiproses) {
			t.Errorf("%s: galat %v", nama, err)
		}
	}
}

// TestGalatDariPanic - hanya panic sikap predikat (rules.PanikSikap: predikat belum
// diport, termasuk pembanding data identitas) yang menjadi galat masukan (422); panic
// lain - bug program: registry rusak, rujukan melingkar - tetap galat server (500).
func TestGalatDariPanic(t *testing.T) {
	if err := galatDariPanic(rules.PanikSikap{Predikat: "UJI", Alasan: "belum diport"}); !errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("panic sikap: %v", err)
	}
	// Lini dari luar yang tidak ada di peta skala = masukan, bukan bug.
	if err := galatDariPanic(premium.PanikLini{Lini: "UJI-LINI"}); !errors.Is(err, ErrTidakDapatDiproses) {
		t.Errorf("panic lini: %v", err)
	}
	for _, r := range []any{"rules: rujukan melingkar lewat UJI", errors.New("lain")} {
		if err := galatDariPanic(r); err == nil || errors.Is(err, ErrTidakDapatDiproses) {
			t.Errorf("%v: galat %v, mau galat server", r, err)
		}
	}
}

// TestTabelRepositoryDikenalTangga - kelima tabel bentuk A repository dikenal
// susunTangga (dua daftar di dua lapisan; tanpa uji ini selisihnya baru terlihat saat
// jalan).
func TestTabelRepositoryDikenalTangga(t *testing.T) {
	for _, nama := range repository.TabelBentukA {
		if _, ada := tabelDikenal[nama]; !ada {
			t.Errorf("tabel %s dibaca repository tetapi tidak dikenal tangga", nama)
		}
	}
	if len(tabelDikenal) != len(repository.TabelBentukA) {
		t.Errorf("tangga mengenal %d tabel, repository membaca %d", len(tabelDikenal), len(repository.TabelBentukA))
	}
}

// akunTiruan - PembacaAkun tiruan: mencatat argumen, mengembalikan baris sintetis UJI-.
type akunTiruan struct {
	baris          []models.Akun
	total          int
	err            error
	cari           *string
	offset, ukuran *int
}

func (a akunTiruan) CariAkun(_ context.Context, cari string, offset, ukuran int) ([]models.Akun, int, error) {
	if a.cari != nil {
		*a.cari, *a.offset, *a.ukuran = cari, offset, ukuran
	}
	return a.baris, a.total, a.err
}

// TestCariAkun - tiket 27: offset dari halaman (ukuran 15, keputusan work owner), teks cari diteruskan
// apa adanya (pola LIKE dibentuk repository - TestPolaCari), total dari repository
// (bukan panjang halaman), 400 untuk halaman/cari tak sah, 503 tanpa basis data.
func TestCariAkun(t *testing.T) {
	var cari string
	var offset, ukuran int
	svc := Baru(nil).DenganAkun(akunTiruan{baris: []models.Akun{{ID: "UJI-1", InsuredID: "UJI-INS"}}, total: 41,
		cari: &cari, offset: &offset, ukuran: &ukuran})
	h, err := svc.CariAkun(context.Background(), "uji_50%", 3)
	if err != nil || h.Total != 41 || h.Halaman != 3 || h.Ukuran != 15 || len(h.Baris) != 1 {
		t.Fatalf("hasil %+v (%v)", h, err)
	}
	if cari != "uji_50%" || offset != 30 || ukuran != 15 {
		t.Errorf("ke repository: cari %q offset %d ukuran %d", cari, offset, ukuran)
	}
	if _, err := svc.CariAkun(context.Background(), "", 1); err != nil || cari != "" || offset != 0 {
		t.Errorf("cari kosong: cari %q offset %d (%v)", cari, offset, err)
	}
	panjang := make([]rune, 256)
	for i := range panjang {
		panjang[i] = 'é'
	}
	for _, u := range []struct {
		cari    string
		halaman int
	}{{"", 0}, {"", -1}, {string(panjang), 1}, {"", halamanMaks + 1}} {
		if _, err := svc.CariAkun(context.Background(), u.cari, u.halaman); !errors.Is(err, ErrMasukanAkun) {
			t.Errorf("halaman %d, cari %d karakter: %v", u.halaman, len([]rune(u.cari)), err)
		}
	}
	if _, err := svc.CariAkun(context.Background(), string(panjang[:255]), 1); err != nil {
		t.Errorf("255 karakter (multibita) ditolak: %v", err)
	}
	if _, err := Baru(nil).CariAkun(context.Background(), "", 1); !errors.Is(err, ErrAkunTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	var _ repository.PembacaAkun = akunTiruan{}
}

// kelasBisnisTiruan - PembacaKelasBisnis tiruan: mencatat argumen, baris sintetis UJI-.
type kelasBisnisTiruan struct {
	baris []models.KelasBisnis
	err   error
	group *string
}

func (k kelasBisnisTiruan) KelasBisnis(_ context.Context, groupBusinessID string) ([]models.KelasBisnis, error) {
	if k.group != nil {
		*k.group = groupBusinessID
	}
	return k.baris, k.err
}

// TestKelasBisnis - tiket 28: groupBusinessId diteruskan APA ADANYA (tanpa pangkas),
// baris dan galat repository diteruskan, 400 untuk kosong/spasi/terlalu panjang (diperiksa
// SEBELUM basis data), 503 tanpa basis data.
func TestKelasBisnis(t *testing.T) {
	var group string
	baris := []models.KelasBisnis{{ID: "UJI-B1", Note: "UJI NOTE A"}}
	svc := Baru(nil).DenganKelasBisnis(kelasBisnisTiruan{baris: baris, group: &group})
	h, err := svc.KelasBisnis(context.Background(), " UJI-GB ")
	if err != nil || len(h) != 1 || h[0] != baris[0] || group != " UJI-GB " {
		t.Fatalf("hasil %+v (%v), ke repository %q", h, err, group)
	}
	galatUji := errors.New("UJI galat repository")
	if _, err := Baru(nil).DenganKelasBisnis(kelasBisnisTiruan{err: galatUji}).KelasBisnis(context.Background(), "UJI-GB"); !errors.Is(err, galatUji) {
		t.Errorf("galat repository tidak diteruskan: %v", err)
	}
	for _, g := range []string{"", "   ", strings.Repeat("U", 4001)} {
		if _, err := svc.KelasBisnis(context.Background(), g); !errors.Is(err, ErrMasukanKelasBisnis) {
			t.Errorf("groupBusinessId panjang %d: %v, mau ErrMasukanKelasBisnis", len(g), err)
		}
		if _, err := Baru(nil).KelasBisnis(context.Background(), g); !errors.Is(err, ErrMasukanKelasBisnis) {
			t.Errorf("tanpa DB, groupBusinessId panjang %d: %v, mau 400 dulu", len(g), err)
		}
	}
	if _, err := svc.KelasBisnis(context.Background(), strings.Repeat("U", 4000)); err != nil {
		t.Errorf("tepat 4000 byte: %v", err)
	}
	if _, err := Baru(nil).KelasBisnis(context.Background(), "UJI-GB"); !errors.Is(err, ErrKelasBisnisTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}
