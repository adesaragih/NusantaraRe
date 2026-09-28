package services

// Efek keluar PremiumList Life - tiket 06. TANPA Oracle.
//
// ⛔ Yang diperiksa EFEKNYA (AC tiket 06): panggilan terjadi atau tidak,
// kegagalan terantre atau tidak, alarm terpicu atau tidak.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

type efekPolisUji struct {
	nama      string
	galat     error
	dipanggil int
	muatan    MuatanEfek
}

func (e *efekPolisUji) Nama() string { return e.nama }
func (e *efekPolisUji) Jalankan(_ context.Context, m MuatanEfek) error {
	e.dipanggil++
	e.muatan = m
	return e.galat
}

type antreanPolisUji struct{ isi []CatatanEfekGagal }

func (a *antreanPolisUji) Antre(_ context.Context, c CatatanEfekGagal) error {
	a.isi = append(a.isi, c)
	return nil
}

var muatanPolisUji = MuatanEfek{KlaimID: "UJI-POLIS-1", AkunID: "UJI-AKUN",
	Waktu: time.Date(2026, 9, 28, 10, 0, 0, 0, time.UTC)}

// TestAlarmHanyaBilaArasapasGagal - AC 27 spec, dan AC relasional tiket 06.
func TestAlarmHanyaBilaArasapasGagal(t *testing.T) {
	ara := &efekPolisUji{nama: "arasapas"}
	alarm := &efekPolisUji{nama: "alarm"}
	antre := &antreanPolisUji{}
	h := NewPenyalurPolis(Produksi, antre, ara, alarm).Salurkan(context.Background(), muatanPolisUji)
	if ara.dipanggil != 1 || alarm.dipanggil != 0 || len(h.Gagal) != 0 || len(antre.isi) != 0 {
		t.Errorf("berhasil: arasapas %d, alarm %d, gagal %d, antre %d; mau 1/0/0/0",
			ara.dipanggil, alarm.dipanggil, len(h.Gagal), len(antre.isi))
	}

	ara = &efekPolisUji{nama: "arasapas", galat: errors.New("UJI gagal")}
	alarm = &efekPolisUji{nama: "alarm"}
	antre = &antreanPolisUji{}
	h = NewPenyalurPolis(Produksi, antre, ara, alarm).Salurkan(context.Background(), muatanPolisUji)
	if alarm.dipanggil != 1 {
		t.Errorf("Arasapas gagal tetapi alarm dipanggil %d kali, mau 1", alarm.dipanggil)
	}
	if len(h.Gagal) != 1 || h.Gagal[0].Nama != "arasapas" || len(antre.isi) != 1 {
		t.Errorf("kegagalan Arasapas tidak tercatat dan terantre: %+v / %+v", h.Gagal, antre.isi)
	}
	if alarm.muatan.KlaimID != "UJI-POLIS-1" {
		t.Errorf("alarm tidak membawa pengenal polisnya: %+v", alarm.muatan)
	}
}

// TestAlarmYangGagalIkutTerantre - alarm yang hilang diam-diam bukan alarm.
func TestAlarmYangGagalIkutTerantre(t *testing.T) {
	ara := &efekPolisUji{nama: "arasapas", galat: errors.New("UJI gagal")}
	alarm := &efekPolisUji{nama: "alarm", galat: errors.New("UJI email gagal")}
	antre := &antreanPolisUji{}
	h := NewPenyalurPolis(Produksi, antre, ara, alarm).Salurkan(context.Background(), muatanPolisUji)
	if len(h.Gagal) != 2 || len(antre.isi) != 2 {
		t.Errorf("gagal %d, antre %d; mau keduanya 2", len(h.Gagal), len(antre.isi))
	}
}

// TestBukanProduksiNolPanggilan - AC 25 spec, ADR-U-0005.
func TestBukanProduksiNolPanggilan(t *testing.T) {
	ara := &efekPolisUji{nama: "arasapas", galat: errors.New("UJI")}
	alarm := &efekPolisUji{nama: "alarm"}
	antre := &antreanPolisUji{}
	h := NewPenyalurPolis(BukanProduksi, antre, ara, alarm).Salurkan(context.Background(), muatanPolisUji)
	if !h.Dilewati || ara.dipanggil+alarm.dipanggil != 0 || len(antre.isi) != 0 {
		t.Errorf("bukan produksi: dilewati %v, panggilan %d, antre %d",
			h.Dilewati, ara.dipanggil+alarm.dipanggil, len(antre.isi))
	}
}

// TestKunciArasapasPolisVERBATIM - dibaca dari activity, bukan ditebak.
func TestKunciArasapasPolisVERBATIM(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\PremiumList Life\Activity\serviceInsertArasapasLife_act.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v)", err)
	}
	teks := string(isi)
	for _, s := range []string{
		`<Kategori_1>"` + KunciArasapasPremiumList.Kategori1 + `"</Kategori_1>`,
		`<Kategori_2>"` + KunciArasapasPremiumList.Kategori2 + `"</Kategori_2>`,
	} {
		if !strings.Contains(teks, s) {
			t.Errorf("serviceInsertArasapasLife_act tidak memuat %s", s)
		}
	}
	if KunciArasapasPremiumList == KunciArasapasLife {
		t.Error("kunci PremiumList sama dengan kunci Claim Life; keduanya activity berbeda")
	}
}

// TestEfekPolisBawaanGagalTerang - stub, bukan panggilan nyata.
func TestEfekPolisBawaanGagalTerang(t *testing.T) {
	ara := EfekArasapasPolis{Resolver: ResolverBelumDiputuskan{}}
	if err := ara.Jalankan(context.Background(), muatanPolisUji); !errors.Is(err, ErrResolverBelumDiputuskan) {
		t.Errorf("Arasapas tanpa resolver: %v", err)
	}
	ara = EfekArasapasPolis{Resolver: resolverUjiPolis{}}
	if err := ara.Jalankan(context.Background(), muatanPolisUji); !errors.Is(err, ErrArasapasBelumDisetujui) {
		t.Errorf("Arasapas ber-resolver: %v, mau ErrArasapasBelumDisetujui", err)
	}
	if err := (EfekAlarmEmailPolis{}).Jalankan(context.Background(), muatanPolisUji); !errors.Is(err, ErrEmailBelumDisetujui) {
		t.Errorf("alarm email: %v", err)
	}
}

type resolverUjiPolis struct{}

func (resolverUjiPolis) Resolve(context.Context, KunciLayanan) (string, error) {
	return "UJI-ALAMAT", nil
}

// TestEfekBerjalanSesudahCommitBukanDiDalamnya - ADR-U-0008.
//
// ⛔ Kegagalan efek keluar tidak boleh membatalkan premium list tersimpan:
// `Salurkan` harus dipanggil SESUDAH `DalamTransaksi` selesai, di luar
// fungsinya.
func TestEfekBerjalanSesudahCommitBukanDiDalamnya(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := badanFungsi(t, string(isi), "func (p *Penawaran) terapkan(")
	iTx := strings.Index(badan, "DalamTransaksi(ctx")
	iTutupTx := strings.Index(badan, "\n\t})")
	iSalur := strings.Index(badan, "p.penyalur.Salurkan(")
	if iTx < 0 || iTutupTx < 0 || iSalur < 0 {
		t.Fatal("terapkan kehilangan transaksi atau penyalurnya")
	}
	if iSalur < iTutupTx {
		t.Error("efek keluar dipanggil DI DALAM transaksi; kegagalannya akan membatalkan simpan")
	}
	if !strings.Contains(badan[iTutupTx:], "akibat.SimpanPolis") {
		t.Error("efek keluar tidak dikurung SimpanPolis; Decline/Reject akan mengirim ke Arasapas")
	}
}
