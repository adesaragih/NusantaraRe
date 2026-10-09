package services_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/benefitlife/backend/models"
	"nusantarare/modul/benefitlife/backend/repository"
	"nusantarare/modul/benefitlife/backend/services"
	"nusantarare/modul/benefitlife/backend/tiruan"
)

var (
	ctx   = context.Background()
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan)
}

// Grid: 10 per halaman (b4526), ID MENURUN bawaan (b4391), saring ID / Benefit.
func TestDaftar(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	h, err := l.Daftar(ctx, models.Saringan{})
	if err != nil || h.Total != 3 || h.Ukuran != 10 || h.Halaman != 1 || h.Daftar[0].ID != "100004" || h.Daftar[2].ID != "100001" {
		t.Fatalf("daftar %+v %v", h, err)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Naik: true}); h.Daftar[0].ID != "100001" {
		t.Errorf("naik %+v", h.Daftar)
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Benefit: " rawat "}); h.Total != 1 || h.Daftar[0].ID != "100001" {
		t.Errorf("saring benefit %+v", h)
	}
	for i := 0; i < 12; i++ {
		id := strconv.Itoa(200000 + i)
		g.Baris[id] = models.Benefit{ID: id, Benefit: "UJI"}
	}
	if h, _ := l.Daftar(ctx, models.Saringan{Halaman: 2}); len(h.Daftar) != 5 || h.Total != 15 {
		t.Errorf("halaman 2 %d dari %d", len(h.Daftar), h.Total)
	}
}

// Save Add: ID = '1' || LPAD(12, 5) = 100012; Benefit dipangkas dan huruf besar (SetUpperCase_DT).
func TestSimpanAdd(t *testing.T) {
	g := tiruan.Contoh()
	b, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{Benefit: "  uji kritis  "})
	if err != nil || b != (models.Benefit{ID: "100012", Benefit: "UJI KRITIS"}) {
		t.Fatalf("add %+v %v", b, err)
	}
	if g.Seq != 13 {
		t.Errorf("sequence dipakai %d kali", g.Seq-12)
	}
}

// Save Edit (EditList_DT): ubah Benefit baris yang ada; ID tak ada = ErrTidakAda.
func TestSimpanEdit(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	b, err := l.Simpan(ctx, penuh, " 100002 ", models.Isian{Benefit: "uji meninggal"})
	if err != nil || b != (models.Benefit{ID: "100002", Benefit: "UJI MENINGGAL"}) || g.Seq != 12 {
		t.Fatalf("edit %+v %v (seq %d)", b, err, g.Seq)
	}
	if _, err := l.Simpan(ctx, penuh, "999999", models.Isian{Benefit: "UJI"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("ID tak ada %v", err)
	}
}

// K6: Benefit wajib (b1137) dan muat 200 byte; tidak ada penolakan kembar (bukan aturan XML).
func TestSimpanValidasi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for _, isi := range []string{"", "   ", strings.Repeat("A", 201)} {
		if _, err := l.Simpan(ctx, penuh, "", models.Isian{Benefit: isi}); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("%q: %v", isi, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{Benefit: "uji rawat inap"}); err != nil {
		t.Errorf("Benefit sama dengan baris lain ditolak - XML tidak memuat aturan kembar: %v", err)
	}
	if g.Seq != 13 {
		t.Errorf("isian yang ditolak memakai sequence: %d", g.Seq)
	}
}

// K3: NEXTVAL menghasilkan ID yang SUDAH ada -> galat ErrIDTerpakai yang jelas (bukan 500, tidak dilompati diam-diam),
// tanpa baris baru; nomor yang melewati LPAD 5 juga.
func TestSimpanIDDariSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 4 // -> 100004, sudah ada
	l := layanan(g)
	_, err := l.Simpan(ctx, penuh, "", models.Isian{Benefit: "UJI BARU"})
	if !errors.Is(err, services.ErrIDTerpakai) || !strings.Contains(err.Error(), "new ID 100004 from M_BENEFIT_LIFE_SEQ already exists") {
		t.Fatalf("mau ErrIDTerpakai yang menyebut ID, dapat %v", err)
	}
	if len(g.Baris) != 3 || g.Baris["100004"].Benefit != "UJI CACAT TETAP" {
		t.Errorf("baris berubah: %+v", g.Baris)
	}
	g.Seq = 100000
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{Benefit: "UJI BARU"}); !errors.Is(err, services.ErrIDTerpakai) ||
		!strings.Contains(err.Error(), "100000") {
		t.Errorf("nomor 6 angka: %v", err)
	}
}

// K3 (balapan): AdaID kosong tetapi INSERT ditolak PK (ORA-00001) -> galat yang sama, transaksi batal.
func TestSimpanBalapanPK(t *testing.T) {
	g := tiruan.Contoh()
	l := services.BaruLayanan(&gudangBalapan{Gudang: g}, tiruan.Transaksi{G: g}.Jalankan)
	if _, err := l.Simpan(ctx, penuh, "", models.Isian{Benefit: "UJI"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("mau ErrIDTerpakai, dapat %v", err)
	}
}

// gudangBalapan - AdaID selalu "kosong", Sisip selalu ORA-00001 (penulis lain mendahului).
type gudangBalapan struct{ *tiruan.Gudang }

func (gudangBalapan) AdaID(context.Context, *db.Tx, string) (bool, error) { return false, nil }
func (gudangBalapan) Sisip(context.Context, *db.Tx, models.Benefit) error {
	return repository.Bungkus(errors.New("ORA-00001: unique constraint (UJI.SYS_C009031) violated"), "menyimpan")
}

// View only / tanpa login: 403, nol sequence terpakai.
func TestSimpanHak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for _, a := range []services.Aktor{{AkunID: "UJI-LIHAT"}, {Penuh: true}} {
		if _, err := l.Simpan(ctx, a, "", models.Isian{Benefit: "UJI"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v: %v", a, err)
		}
	}
	if g.Seq != 12 {
		t.Errorf("sequence dipakai tanpa hak")
	}
}
