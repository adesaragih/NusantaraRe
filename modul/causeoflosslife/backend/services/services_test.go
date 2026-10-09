package services_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/causeoflosslife/backend/models"
	"nusantarare/modul/causeoflosslife/backend/repository"
	"nusantarare/modul/causeoflosslife/backend/services"
	"nusantarare/modul/causeoflosslife/backend/tiruan"
)

var (
	ctx   = context.Background()
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan)
}

// Grid 10 baris, ID menaik (b3923); halaman < 1 = 1; baris 100001 kosong tetap tampil (tidak diisi, tidak dihapus).
func TestDaftar(t *testing.T) {
	g := tiruan.Contoh()
	for i := 5; i <= 12; i++ {
		id, _ := models.BentukID(strconv.Itoa(i))
		g.Baris[id] = models.CauseOfLoss{ID: id, CauseOfLoss: "UJI " + id}
	}
	l := layanan(g)
	h, err := l.Daftar(ctx, models.Saringan{Halaman: 0})
	if err != nil || h.Total != 12 || h.Ukuran != 10 || h.Halaman != 1 || len(h.Daftar) != 10 || h.Daftar[0] != (models.CauseOfLoss{ID: "100001"}) {
		t.Fatalf("%+v %v", h, err)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Halaman: 2}); len(h.Daftar) != 2 || h.Daftar[1].ID != "100012" {
		t.Errorf("halaman 2 %+v", h.Daftar)
	}
}

// Add: ID '1' || LPAD(5, 5) = 100005 (sequence DEV last_number 5); dipangkas, huruf TIDAK diubah (K4).
func TestSimpanAdd(t *testing.T) {
	g := tiruan.Contoh()
	c, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{CauseOfLoss: "  uji Bencana Alam "})
	if err != nil || c != (models.CauseOfLoss{ID: "100005", CauseOfLoss: "uji Bencana Alam"}) || g.Seq != 6 {
		t.Fatalf("add %+v %v %d", c, err, g.Seq)
	}
}

// K4: wajib (XML); tidak kembar tanpa beda huruf (di luar XML) - Edit boleh menyimpan namanya sendiri; ID tak ada 404;
// baris kosong 100001 dapat diisi lewat Edit.
func TestSimpanWajibDanUnik(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{CauseOfLoss: "  "}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Cause of Loss is required") {
		t.Errorf("wajib %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{CauseOfLoss: " uji sakit"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Cause of Loss uji sakit is already used by ID 100002") {
		t.Errorf("kembar %v", err)
	}
	if len(g.Baris) != 4 || g.Seq != 5 {
		t.Errorf("ditolak tetapi berubah %d %d", len(g.Baris), g.Seq)
	}
	if c, err := l.Simpan(ctx, penuh, "100002", models.Isian{CauseOfLoss: "UJI Sakit"}); err != nil || c.CauseOfLoss != "UJI Sakit" {
		t.Errorf("edit nama sendiri %+v %v", c, err)
	}
	if _, err := l.Simpan(ctx, penuh, "100003", models.Isian{CauseOfLoss: "uji semua sebab"}); err == nil || !strings.Contains(err.Error(), "already used by ID 100004") {
		t.Errorf("edit jadi kembar %v", err)
	}
	if c, err := l.Simpan(ctx, penuh, "100001", models.Isian{CauseOfLoss: "UJI LAIN"}); err != nil || c.ID != "100001" || g.Seq != 5 {
		t.Errorf("edit baris kosong %+v %v", c, err)
	}
	if _, err := l.Simpan(ctx, penuh, "999999", models.Isian{CauseOfLoss: "UJI Z"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada %v", err)
	}
}

// K3: NEXTVAL -> ID yang sudah ada = ErrIDTerpakai berkalimat; nomor 6 angka juga; balapan PK juga; nol baris berubah.
func TestSimpanIDDariSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 2 // -> 100002 sudah ada
	isi := models.Isian{CauseOfLoss: "UJI BARU"}
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "new ID 100002 from M_CAUSEOFLOSS_LIFE_SEQ already exists") {
		t.Fatalf("%v", err)
	}
	g.Seq = 100000
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "does not fit the ID formula") {
		t.Errorf("6 angka %v", err)
	}
	if len(g.Baris) != 4 {
		t.Errorf("baris berubah")
	}
	b := services.BaruLayanan(&balapan{tiruan.Contoh()}, tiruan.Transaksi{G: tiruan.Contoh()}.Jalankan)
	if _, err := b.Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("balapan %v", err)
	}
}

type balapan struct{ *tiruan.Gudang }

func (balapan) AdaID(context.Context, *db.Tx, string) (bool, error) { return false, nil }
func (balapan) Sisip(context.Context, *db.Tx, models.CauseOfLoss) error {
	return repository.Bungkus(errors.New("ORA-00001: unique constraint (UJI.SYS_C008825) violated"), "menyimpan")
}

// View only / tanpa login: 403 sebelum apa pun.
func TestSimpanHak(t *testing.T) {
	g := tiruan.Contoh()
	for _, a := range []services.Aktor{{AkunID: "UJI-LIHAT"}, {Penuh: true}} {
		if _, err := layanan(g).Simpan(ctx, a, "", models.Isian{CauseOfLoss: "UJI"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v %v", a, err)
		}
	}
	if g.Seq != 5 {
		t.Errorf("sequence terpakai %d", g.Seq)
	}
}
