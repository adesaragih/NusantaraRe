package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// antreanTiruan mencatat setiap kegagalan yang diantre.
type antreanTiruan struct{ catatan []outbox.CatatanEfekGagal }

func (a *antreanTiruan) Antre(_ context.Context, c outbox.CatatanEfekGagal) error {
	a.catatan = append(a.catatan, c)
	return nil
}

// efekTiruan - efek keluar yang dihitung panggilannya.
type efekTiruan struct {
	nama    string
	galat   error
	dipakai int
}

func (e *efekTiruan) Nama() string { return e.nama }
func (e *efekTiruan) Jalankan(context.Context, outbox.MuatanEfek) error {
	e.dipakai++
	return e.galat
}

// resolverTiruan menjawab satu teks untuk kunci Arasapas Endorsement.
type resolverTiruan struct {
	kunci layanan.KunciLayanan
	isi   string
}

func (r *resolverTiruan) Resolve(_ context.Context, k layanan.KunciLayanan) (string, error) {
	r.kunci = k
	return r.isi, nil
}

func TestConfirmBukanProduksiMelewatiEfekKeluar(t *testing.T) {
	_, l, _ := gudangPutusan(t)
	h, err := l.Putuskan(context.Background(), pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"})
	if err != nil || h.EfekKeluar == nil || !h.EfekKeluar.Dilewati || len(h.EfekKeluar.Gagal) != 0 || h.Status != models.StatusKasusSelesai {
		t.Fatalf("%+v %v", h, err)
	}
}

func TestAlarmHanyaBilaKirimanGagal(t *testing.T) {
	ctx := context.Background()
	// Kiriman berhasil: alarm tidak menyala (alarm bukan notifikasi bisnis).
	_, l, _ := gudangPutusan(t)
	kirim, alarm, a := &efekTiruan{nama: "arasapas"}, &efekTiruan{nama: "alarm"}, &antreanTiruan{}
	l.DenganPenyalur(services.NewPenyalurEDM(inti.Produksi, a, kirim, alarm))
	h, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"})
	if err != nil || kirim.dipakai != 1 || alarm.dipakai != 0 || len(h.EfekKeluar.Gagal) != 0 || h.EfekKeluar.Dilewati {
		t.Fatalf("berhasil: %+v kirim %d alarm %d %v", h.EfekKeluar, kirim.dipakai, alarm.dipakai, err)
	}

	// Kiriman gagal: alarm menyala, keduanya diantre, versi tetap resmi, alamat tidak tampil.
	g, l2, _ := gudangPutusan(t)
	var log []string
	l2 = services.BaruLayanan(g, tiruan.Transaksi, func(s string) { log = append(log, s) })
	kirim = &efekTiruan{nama: "arasapas", galat: errors.New("UJI gagal di uji-alamat-rahasia")}
	alarm, a = &efekTiruan{nama: "alarm", galat: outbox.ErrEmailBelumDisetujui}, &antreanTiruan{}
	l2.DenganJejak(&jejakTiruan{}).DenganPenyalur(services.NewPenyalurEDM(inti.Produksi, a, kirim, alarm))
	h, err = l2.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"})
	if err != nil || h.Status != models.StatusKasusSelesai || g.Polis["EDMLF-1"].Status != models.StatusKasusSelesai {
		t.Fatalf("kegagalan efek membatalkan Confirm: %+v %v", h, err)
	}
	if alarm.dipakai != 1 || len(a.catatan) != 2 || len(h.EfekKeluar.Gagal) != 2 {
		t.Fatalf("alarm %d diantre %d ringkas %+v", alarm.dipakai, len(a.catatan), h.EfekKeluar)
	}
	semua := strings.Join(append(append([]string{}, h.EfekKeluar.Gagal...), log...), "\n")
	if strings.Contains(semua, "uji-alamat-rahasia") || !strings.Contains(semua, "recorded in the outbox") || len(log) != 2 {
		t.Errorf("alarm layar/log: %q", semua)
	}
}

func TestDeclineTanpaEfekKeluar(t *testing.T) {
	_, l, _ := gudangPutusan(t)
	kirim := &efekTiruan{nama: "arasapas"}
	l.DenganPenyalur(services.NewPenyalurEDM(inti.Produksi, &antreanTiruan{}, kirim, &efekTiruan{nama: "alarm"}))
	h, err := l.Putuskan(context.Background(), pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "2"})
	if err != nil || h.EfekKeluar != nil || kirim.dipakai != 0 {
		t.Fatalf("%+v %v kirim %d", h, err, kirim.dipakai)
	}
}

func TestEfekArasapasMeresolveKunciEndorsement(t *testing.T) {
	r := &resolverTiruan{isi: "UJI-ALAMAT"}
	err := services.EfekArasapasEDM{Resolver: r}.Jalankan(context.Background(), outbox.MuatanEfek{})
	if !errors.Is(err, outbox.ErrArasapasBelumDisetujui) || r.kunci != services.KunciArasapasEndorsement {
		t.Fatalf("%v %+v", err, r.kunci)
	}
	if services.KunciArasapasEndorsement.Kategori1 != "Production" || services.KunciArasapasEndorsement.Kategori2 != "convertJsonNusareToProduction" {
		t.Error("kunci M_LINK_SERVICE bukan b793/b795")
	}
	if err := (services.EfekArasapasEDM{Resolver: &resolverTiruan{}}).Jalankan(context.Background(), outbox.MuatanEfek{}); !errors.Is(err, layanan.ErrEndpointTidakDitemukan) {
		t.Errorf("alamat kosong: %v", err)
	}
	if err := (services.EfekArasapasEDM{}).Jalankan(context.Background(), outbox.MuatanEfek{}); !errors.Is(err, layanan.ErrResolverBelumDiputuskan) {
		t.Errorf("tanpa resolver: %v", err)
	}
	if err := (services.EfekAlarmEmailEDM{}).Jalankan(context.Background(), outbox.MuatanEfek{}); !errors.Is(err, outbox.ErrEmailBelumDisetujui) {
		t.Errorf("email: %v", err)
	}
}
