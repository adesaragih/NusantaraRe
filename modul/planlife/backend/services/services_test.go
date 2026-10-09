package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/planlife/backend/models"
	"nusantarare/modul/planlife/backend/repository"
	"nusantarare/modul/planlife/backend/services"
	"nusantarare/modul/planlife/backend/tiruan"
)

var (
	ctx   = context.Background()
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan)
}

// Grid 10 baris, bawaan ID menaik; urut Plan Name / Benefit; kolom lain = bawaan.
func TestDaftar(t *testing.T) {
	l := layanan(tiruan.Contoh())
	h, err := l.Daftar(ctx, models.Saringan{})
	if err != nil || h.Total != 2 || h.Ukuran != 10 || h.Daftar[0].ID != "100001" {
		t.Fatalf("%+v %v", h, err)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Urut: "BENEFIT", Turun: true}); h.Daftar[0].ID != "100002" {
		t.Errorf("urut benefit %+v", h.Daftar)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Urut: "business", Turun: true}); h.Daftar[0].ID != "100001" {
		t.Errorf("Business tidak dapat diurutkan (b4269) -> bawaan %+v", h.Daftar)
	}
}

// K4: pilihan Business = grup 009 saja; Benefit = seluruh BENEFIT_LIFE.
func TestPilihan(t *testing.T) {
	l := layanan(tiruan.Contoh())
	b, _ := l.PilihanBusiness(ctx)
	if len(b) != 2 || b[0] != (models.PilihanBusiness{ID: "9001", OldID: "L1", Note: "UJI KREDIT"}) {
		t.Errorf("business %+v", b)
	}
	if e, _ := l.PilihanBenefit(ctx); len(e) != 2 {
		t.Errorf("benefit %+v", e)
	}
}

// Add: ID '1' || LPAD(44, 5) = 100044; K4 teks dicocokkan TANPA beda huruf -> nama + ID dari master.
func TestSimpanAddCocokMaster(t *testing.T) {
	g := tiruan.Contoh()
	p, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{CoverName: " UJI Plan Baru ", Business: " uji jiwa KUMPULAN", Benefit: "uji rawat inap"})
	if err != nil || p != (models.Plan{ID: "100044", CoverName: "UJI Plan Baru", Business: "UJI Jiwa Kumpulan", BusinessID: "9002",
		Benefit: "UJI RAWAT INAP", BenefitID: "100001"}) {
		t.Fatalf("add %+v %v", p, err)
	}
}

// K4: teks yang tidak ada di master (atau business grup lain) = 422 berkalimat, nol baris, nol sequence.
func TestSimpanTolakBukanMaster(t *testing.T) {
	g := tiruan.Contoh()
	_, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{CoverName: "UJI X", Business: "UJI KEBAKARAN", Benefit: "UJI LAIN"})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), `Business "UJI KEBAKARAN" is not in the Business list`) ||
		!strings.Contains(err.Error(), `Benefit "UJI LAIN" is not in the Benefit list`) {
		t.Errorf("%v", err)
	}
	if len(g.Baris) != 2 || g.Seq != 44 {
		t.Errorf("berubah %d %d", len(g.Baris), g.Seq)
	}
	g.Benefit = append(g.Benefit, models.PilihanBenefit{ID: "100009", Benefit: "uji rawat inap"})
	if _, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{CoverName: "UJI X", Business: "UJI KREDIT", Benefit: "UJI RAWAT INAP"}); err == nil ||
		!strings.Contains(err.Error(), "matches more than one Benefit") {
		t.Errorf("ganda %v", err)
	}
}

// K5: ketiga medan wajib; Plan Name unik tanpa beda huruf (Edit boleh menyimpan namanya sendiri).
func TestSimpanWajibDanUnik(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Plan Name is required; Business is required; Benefit is required") {
		t.Errorf("wajib %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{CoverName: "uji plan a", Business: "UJI KREDIT", Benefit: "UJI MENINGGAL"}); err == nil ||
		!strings.Contains(err.Error(), "Plan Name uji plan a is already used by ID 100001") {
		t.Errorf("kembar %v", err)
	}
	p, err := l.Simpan(ctx, penuh, "100001", models.Isian{CoverName: "UJI PLAN A", Business: "UJI KREDIT", Benefit: "UJI RAWAT INAP"})
	if err != nil || p.BenefitID != "100001" || g.Seq != 44 {
		t.Errorf("edit %+v %v", p, err)
	}
	if _, err := l.Simpan(ctx, penuh, "999999", models.Isian{CoverName: "UJI Z", Business: "UJI KREDIT", Benefit: "UJI RAWAT INAP"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada %v", err)
	}
}

// K3: NEXTVAL -> ID yang sudah ada = ErrIDTerpakai berkalimat; nomor 6 angka juga; balapan PK juga.
func TestSimpanIDDariSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 2 // -> 100002 sudah ada
	isi := models.Isian{CoverName: "UJI BARU", Business: "UJI KREDIT", Benefit: "UJI MENINGGAL"}
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "new ID 100002 from M_PRODUCT_TYPE_LIFE_SEQ already exists") {
		t.Fatalf("%v", err)
	}
	g.Seq = 100000
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("6 angka %v", err)
	}
	if len(g.Baris) != 2 {
		t.Errorf("baris berubah")
	}
	b := services.BaruLayanan(&balapan{tiruan.Contoh()}, tiruan.Transaksi{G: tiruan.Contoh()}.Jalankan)
	if _, err := b.Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("balapan %v", err)
	}
}

type balapan struct{ *tiruan.Gudang }

func (balapan) AdaID(context.Context, *db.Tx, string) (bool, error) { return false, nil }
func (balapan) Sisip(context.Context, *db.Tx, models.Plan) error {
	return repository.Bungkus(errors.New("ORA-00001: unique constraint (UJI.PK_PRODUCT_TYPE_LIFE) violated"), "menyimpan")
}

// View only / tanpa login: 403 sebelum apa pun.
func TestSimpanHak(t *testing.T) {
	g := tiruan.Contoh()
	for _, a := range []services.Aktor{{AkunID: "UJI-LIHAT"}, {Penuh: true}} {
		if _, err := layanan(g).Simpan(ctx, a, "", models.Isian{CoverName: "UJI", Business: "UJI KREDIT", Benefit: "UJI MENINGGAL"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v %v", a, err)
		}
	}
}
