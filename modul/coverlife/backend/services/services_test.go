package services_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/coverlife/backend/models"
	"nusantarare/modul/coverlife/backend/repository"
	"nusantarare/modul/coverlife/backend/services"
	"nusantarare/modul/coverlife/backend/tiruan"
)

var (
	ctx   = context.Background()
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan)
}

// Grid 50 baris (b4101), ID menaik (b3967); halaman < 1 = 1.
func TestDaftar(t *testing.T) {
	g := tiruan.Contoh()
	for i := 5; i <= 55; i++ {
		id, _ := models.BentukID(strconv.Itoa(i))
		g.Baris[id] = models.Cover{ID: id, Cover: "UJI " + id}
	}
	l := layanan(g)
	h, err := l.Daftar(ctx, models.Saringan{Halaman: 0})
	if err != nil || h.Total != 55 || h.Ukuran != 50 || h.Halaman != 1 || len(h.Daftar) != 50 || h.Daftar[0].ID != "100001" {
		t.Fatalf("%+v %v", h.Total, err)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Halaman: 2}); len(h.Daftar) != 5 || h.Daftar[4].ID != "100055" {
		t.Errorf("halaman 2 %+v", h.Daftar)
	}
}

// Add: ID '1' || LPAD(5, 5) = 100005 (sequence DEV last_number 5); Cover dipangkas, huruf TIDAK diubah; Note apa
// adanya, hanya spasi = kosong.
func TestSimpanAdd(t *testing.T) {
	g := tiruan.Contoh()
	c, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{Cover: "  uji Penyakit Kritis ", Note: "uji catatan"})
	if err != nil || c != (models.Cover{ID: "100005", Cover: "uji Penyakit Kritis", Note: "uji catatan"}) || g.Seq != 6 {
		t.Fatalf("add %+v %v %d", c, err, g.Seq)
	}
	c, err = layanan(g).Simpan(ctx, penuh, "", models.Isian{Cover: "UJI CACAT", Note: "   "})
	if err != nil || c.Note != "" || c.ID != "100006" {
		t.Errorf("note spasi %+v %v", c, err)
	}
}

// C3: Cover wajib (XML); tidak kembar tanpa beda huruf (di luar XML) - Edit boleh menyimpan namanya sendiri; ID tak
// ada 404; ditolak = nol perubahan.
func TestSimpanWajibDanUnik(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{Cover: "  ", Note: "UJI"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Cover is required") {
		t.Errorf("wajib %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{Cover: " uji jiwa"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "Cover uji jiwa is already used by ID 100002") {
		t.Errorf("kembar %v", err)
	}
	if len(g.Baris) != 4 || g.Seq != 5 {
		t.Errorf("ditolak tetapi berubah %d %d", len(g.Baris), g.Seq)
	}
	if c, err := l.Simpan(ctx, penuh, "100002", models.Isian{Cover: "UJI Jiwa", Note: "uji catatan jiwa"}); err != nil ||
		c != (models.Cover{ID: "100002", Cover: "UJI Jiwa", Note: "uji catatan jiwa"}) {
		t.Errorf("edit nama sendiri %+v %v", c, err)
	}
	if _, err := l.Simpan(ctx, penuh, "100003", models.Isian{Cover: "uji rider"}); err == nil || !strings.Contains(err.Error(), "already used by ID 100004") {
		t.Errorf("edit jadi kembar %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, "999999", models.Isian{Cover: "UJI Z"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada %v", err)
	}
}

// C2: NEXTVAL -> ID yang sudah ada = ErrIDTerpakai berkalimat; nomor 6 angka juga; balapan PK juga; nol baris berubah.
func TestSimpanIDDariSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 2 // -> 100002 sudah ada
	isi := models.Isian{Cover: "UJI BARU"}
	if _, err := layanan(g).Simpan(ctx, penuh, "", isi); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "new ID 100002 from M_COVER_LIFE_SEQ already exists") {
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
func (balapan) Sisip(context.Context, *db.Tx, models.Cover) error {
	return repository.Bungkus(errors.New("ORA-00001: unique constraint (UJI.SYS_C009203) violated"), "menyimpan")
}

// View only / tanpa login: 403 sebelum apa pun.
func TestSimpanHak(t *testing.T) {
	g := tiruan.Contoh()
	for _, a := range []services.Aktor{{AkunID: "UJI-LIHAT"}, {Penuh: true}} {
		if _, err := layanan(g).Simpan(ctx, a, "", models.Isian{Cover: "UJI"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v %v", a, err)
		}
	}
	if g.Seq != 5 {
		t.Errorf("sequence terpakai %d", g.Seq)
	}
}
