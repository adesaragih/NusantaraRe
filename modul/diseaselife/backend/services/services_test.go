package services_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/repository"
	"nusantarare/modul/diseaselife/backend/services"
	"nusantarare/modul/diseaselife/backend/tiruan"
)

var (
	ctx   = context.Background()
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan)
}

// Grid 10 baris, bawaan ID MENURUN (b5012); urut ICD Code (b5036); saring ICD Code / Disease sampai ke gudang
// (server), dirapikan; halaman < 1 = 1; urut tak dikenal = ID.
func TestDaftar(t *testing.T) {
	g := tiruan.Contoh()
	for i := 100006; i <= 100012; i++ {
		id := strconv.Itoa(i)
		g.Baris[id] = models.Penyakit{ID: id, ICDCode: "UJI" + id[4:], Disease: "UJI PENYAKIT " + id}
	}
	l := layanan(g)
	h, err := l.Daftar(ctx, models.Saringan{Halaman: 0, Urut: "disease"})
	if err != nil || h.Total != 12 || h.Ukuran != 10 || h.Halaman != 1 || len(h.Daftar) != 10 || h.Daftar[0].ID != "100012" {
		t.Fatalf("%+v %v", h, err)
	}
	if g.Diminta.Urut != models.UrutID {
		t.Errorf("urut tak dikenal diteruskan %q", g.Diminta.Urut)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Halaman: 2}); len(h.Daftar) != 2 || h.Daftar[1].ID != "100001" {
		t.Errorf("halaman 2 %+v", h.Daftar)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Urut: models.UrutICD, Naik: true, Halaman: 1}); h.Daftar[0].ICDCode != "UJI01" {
		t.Errorf("urut icd %+v", h.Daftar[0])
	}
	h, _ = l.Daftar(ctx, models.Saringan{ICDCode: "  uji0 ", Disease: "demam", Halaman: 1})
	if h.Total != 2 || g.Diminta.ICDCode != "uji0" || g.Diminta.Disease != "demam" {
		t.Errorf("saring %+v %+v", h, g.Diminta)
	}
}

// Add: ID = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL) (D1.1); ICD Code / Disease dipangkas dan HURUF BESAR (SetUpperCase_DT).
func TestSimpanAdd(t *testing.T) {
	g := tiruan.Contoh()
	p, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{ICDCode: " uji99 ", Disease: "  uji demam berdarah "})
	if err != nil || p != (models.Penyakit{ID: "100006", ICDCode: "UJI99", Disease: "UJI DEMAM BERDARAH"}) || g.Seq != 100007 {
		t.Fatalf("add %+v %v %d", p, err, g.Seq)
	}
}

// D3: wajib; ICD Code tidak kembar tanpa beda huruf (di luar XML) - Edit boleh menyimpan ICD-nya sendiri; ID tak ada
// 404; ditolak = nol perubahan.
func TestSimpanWajibDanUnik(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{ICDCode: " ", Disease: "UJI X"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "ICD Code is required") {
		t.Errorf("icd wajib %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{ICDCode: "UJI98", Disease: ""}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Disease is required") {
		t.Errorf("disease wajib %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{ICDCode: " uji02", Disease: "UJI X"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "ICD Code UJI02 is already used by ID 100002") {
		t.Errorf("kembar %v", err)
	}
	if len(g.Baris) != 5 || g.Seq != 100006 {
		t.Errorf("ditolak tetapi berubah %d %d", len(g.Baris), g.Seq)
	}
	if p, err := l.Simpan(ctx, penuh, "100002", models.Isian{ICDCode: "uji02", Disease: "uji demam tifoid berat"}); err != nil ||
		p != (models.Penyakit{ID: "100002", ICDCode: "UJI02", Disease: "UJI DEMAM TIFOID BERAT"}) {
		t.Errorf("edit icd sendiri %+v %v", p, err)
	}
	if _, err := l.Simpan(ctx, penuh, "100003", models.Isian{ICDCode: "UJI04", Disease: "UJI Y"}); err == nil || !strings.Contains(err.Error(), "already used by ID 100004") {
		t.Errorf("edit jadi kembar %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "999999", models.Isian{ICDCode: "UJI97", Disease: "UJI Z"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada %v", err)
	}
}

// D1.1: NEXTVAL -> ID yang sudah ada (mis. ditulis prosedur Pega) = ErrIDTerpakai berkalimat; balapan PK juga; nol
// baris berubah.
func TestSimpanIDDariSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 100002
	isi := models.Isian{ICDCode: "UJI96", Disease: "UJI BARU"}
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "new ID 100002 from SEQ_DISEASE_LIFE already exists") {
		t.Fatalf("%v", err)
	}
	if len(g.Baris) != 5 {
		t.Errorf("baris berubah")
	}
	b := services.BaruLayanan(&balapan{tiruan.Contoh()}, tiruan.Transaksi{G: tiruan.Contoh()}.Jalankan)
	if _, err := b.Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("balapan %v", err)
	}
}

type balapan struct{ *tiruan.Gudang }

func (balapan) AdaID(context.Context, *db.Tx, string) (bool, error) { return false, nil }
func (balapan) Sisip(context.Context, *db.Tx, models.Penyakit) error {
	return repository.Bungkus(errors.New("ORA-00001: unique constraint (UJI.PK_DISEASE_LIFE) violated"), "menyimpan")
}

// View only / tanpa login: 403 sebelum apa pun.
func TestSimpanHak(t *testing.T) {
	g := tiruan.Contoh()
	for _, a := range []services.Aktor{{AkunID: "UJI-LIHAT"}, {Penuh: true}} {
		if _, err := layanan(g).Simpan(ctx, a, "", models.Isian{ICDCode: "UJI95", Disease: "UJI"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v %v", a, err)
		}
	}
	if g.Seq != 100006 {
		t.Errorf("sequence terpakai %d", g.Seq)
	}
}
