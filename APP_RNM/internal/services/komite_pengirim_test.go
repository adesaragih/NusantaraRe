package services

// Pelaksana outbox Komite - tiket 07. TANPA Oracle.

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"nusantarare/internal/repository"
)

type riwayatUji struct {
	sudah   bool
	ditanya int
}

func (r *riwayatUji) SudahSelesai(context.Context, *repository.Tx, string, string, string, string) (bool, error) {
	r.ditanya++
	return r.sudah, nil
}

type resolverUjiKomite struct {
	alamat  string
	diminta []KunciLayanan
}

func (r *resolverUjiKomite) Resolve(_ context.Context, k KunciLayanan) (string, error) {
	r.diminta = append(r.diminta, k)
	return r.alamat, nil
}

func barisKomite(t *testing.T, jenis string, kunci KunciLayanan) repository.BarisEfekKeluar {
	t.Helper()
	b, err := json.Marshal(muatanOutboxKomite{KasusID: "KMTLF-UJI", KlaimID: "UJI-K",
		Kategori1: kunci.Kategori1, Kategori2: kunci.Kategori2})
	if err != nil {
		t.Fatal(err)
	}
	return repository.BarisEfekKeluar{ID: "9", Modul: ModulKomiteLife, Jenis: jenis,
		Rujukan: "KMTLF-UJI", Muatan: string(b)}
}

// TestKasirTidakDikirimDuaKali - AC 22 spec.
func TestKasirTidakDikirimDuaKali(t *testing.T) {
	r := &riwayatUji{sudah: true}
	res := &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	p := PelaksanaKomite{Lingkungan: Produksi, Resolver: res, Riwayat: r}
	if err := p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite)); err != nil {
		t.Errorf("kasir yang sudah selesai: %v, mau nil (tuntas tanpa kirim)", err)
	}
	if len(res.diminta) != 0 {
		t.Error("kasir yang sudah selesai masih me-resolve alamat (akan dikirim ulang)")
	}
}

// TestKunciHilangPermanenJaringanBukan - AC "kunci tidak ditemukan ≠ jaringan gagal".
func TestKunciHilangPermanenJaringanBukan(t *testing.T) {
	p := PelaksanaKomite{Lingkungan: Produksi, Resolver: &resolverUjiKomite{alamat: ""}, Riwayat: &riwayatUji{}}
	err := p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite))
	if !errors.Is(err, ErrEndpointTidakDitemukan) || LayakDicobaUlang(err) {
		t.Errorf("kunci hilang: %v (layak ulang %v), mau permanen", err, LayakDicobaUlang(err))
	}
	p.Resolver = &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	err = p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite))
	if !errors.Is(err, ErrKasirBelumDisetujui) || LayakDicobaUlang(err) {
		t.Errorf("kasir belum disetujui: %v, mau permanen", err)
	}
	if !LayakDicobaUlang(errors.New("UJI jaringan putus")) {
		t.Error("galat jaringan tidak layak dicoba ulang")
	}
}

// TestNonProduksiStubTanpaPanggilan - ADR-0005, km4.
func TestNonProduksiStubTanpaPanggilan(t *testing.T) {
	res := &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	p := PelaksanaKomite{Lingkungan: BukanProduksi, Resolver: res, Riwayat: &riwayatUji{}}
	for _, j := range []string{JenisEfekKomiteArasapas, JenisEfekKomiteEmail, JenisEfekKomiteKasir} {
		if err := p.Laksanakan(context.Background(), nil, barisKomite(t, j, KunciKasirKomite)); err != nil {
			t.Errorf("%s di non-produksi: %v", j, err)
		}
	}
	if len(res.diminta) != 0 {
		t.Error("non-produksi me-resolve alamat layanan")
	}
}

// TestBarisAsingDitolak - pelaksana Komite hanya untuk modulnya.
func TestBarisAsingDitolak(t *testing.T) {
	p := PelaksanaKomite{Lingkungan: Produksi, Riwayat: &riwayatUji{}}
	b := barisKomite(t, JenisEfekKomiteEmail, KunciLayanan{})
	b.Modul = ModulClaimLife
	if err := p.Laksanakan(context.Background(), nil, b); !errors.Is(err, ErrPermintaanTidakSah) {
		t.Errorf("baris modul lain: %v", err)
	}
	b = barisKomite(t, "asing", KunciLayanan{})
	if err := p.Laksanakan(context.Background(), nil, b); !errors.Is(err, ErrPermintaanTidakSah) {
		t.Errorf("jenis asing: %v", err)
	}
}
