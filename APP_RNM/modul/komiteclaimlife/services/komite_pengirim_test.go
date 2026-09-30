package services

// Pelaksana outbox Komite - tiket 07. TANPA Oracle.

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
)

type riwayatUji struct {
	sudah   bool
	ditanya int
}

func (r *riwayatUji) SudahSelesai(context.Context, *db.Tx, string, string, string, string) (bool, error) {
	r.ditanya++
	return r.sudah, nil
}

type resolverUjiKomite struct {
	alamat  string
	diminta []layanan.KunciLayanan
}

func (r *resolverUjiKomite) Resolve(_ context.Context, k layanan.KunciLayanan) (string, error) {
	r.diminta = append(r.diminta, k)
	return r.alamat, nil
}

func barisKomite(t *testing.T, jenis string, kunci layanan.KunciLayanan) outbox.BarisEfekKeluar {
	t.Helper()
	b, err := json.Marshal(muatanOutboxKomite{KasusID: "KMTLF-UJI", KlaimID: "UJI-K",
		Kategori1: kunci.Kategori1, Kategori2: kunci.Kategori2})
	if err != nil {
		t.Fatal(err)
	}
	return outbox.BarisEfekKeluar{ID: "9", Modul: ModulKomiteLife, Jenis: jenis,
		Rujukan: "KMTLF-UJI", Muatan: string(b)}
}

// TestKasirTidakDikirimDuaKali - AC 22 spec.
func TestKasirTidakDikirimDuaKali(t *testing.T) {
	r := &riwayatUji{sudah: true}
	res := &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	p := PelaksanaKomite{Lingkungan: inti.Produksi, Resolver: res, Riwayat: r}
	if err := p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite)); err != nil {
		t.Errorf("kasir yang sudah selesai: %v, mau nil (tuntas tanpa kirim)", err)
	}
	if len(res.diminta) != 0 {
		t.Error("kasir yang sudah selesai masih me-resolve alamat (akan dikirim ulang)")
	}
}

// TestKunciHilangPermanenJaringanBukan - AC "kunci tidak ditemukan ≠ jaringan gagal".
func TestKunciHilangPermanenJaringanBukan(t *testing.T) {
	p := PelaksanaKomite{Lingkungan: inti.Produksi, Resolver: &resolverUjiKomite{alamat: ""}, Riwayat: &riwayatUji{}}
	err := p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite))
	if !errors.Is(err, layanan.ErrEndpointTidakDitemukan) || outbox.LayakDicobaUlang(err) {
		t.Errorf("kunci hilang: %v (layak ulang %v), mau permanen", err, outbox.LayakDicobaUlang(err))
	}
	p.Resolver = &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	err = p.Laksanakan(context.Background(), nil, barisKomite(t, JenisEfekKomiteKasir, KunciKasirKomite))
	if !errors.Is(err, outbox.ErrKasirBelumDisetujui) || outbox.LayakDicobaUlang(err) {
		t.Errorf("kasir belum disetujui: %v, mau permanen", err)
	}
	if !outbox.LayakDicobaUlang(errors.New("UJI jaringan putus")) {
		t.Error("galat jaringan tidak layak dicoba ulang")
	}
}

// TestNonProduksiStubTanpaPanggilan - ADR-0005, km4.
func TestNonProduksiStubTanpaPanggilan(t *testing.T) {
	res := &resolverUjiKomite{alamat: "UJI-ALAMAT"}
	p := PelaksanaKomite{Lingkungan: inti.BukanProduksi, Resolver: res, Riwayat: &riwayatUji{}}
	for _, j := range []string{JenisEfekKomiteArasapas, JenisEfekKomiteEmail, JenisEfekKomiteKasir} {
		err := p.Laksanakan(context.Background(), nil, barisKomite(t, j, KunciKasirKomite))
		// ⛔ Stub TIDAK BOLEH mengaku terkirim (temuan /code-review).
		if !errors.Is(err, outbox.ErrPengirimStubNonProduksi) || outbox.LayakDicobaUlang(err) {
			t.Errorf("%s di non-produksi: %v, mau stub permanen", j, err)
		}
	}
	if len(res.diminta) != 0 {
		t.Error("non-produksi me-resolve alamat layanan")
	}
}

// TestBarisAsingDitolak - pelaksana Komite hanya untuk modulnya.
func TestBarisAsingDitolak(t *testing.T) {
	p := PelaksanaKomite{Lingkungan: inti.Produksi, Riwayat: &riwayatUji{}}
	b := barisKomite(t, JenisEfekKomiteEmail, layanan.KunciLayanan{})
	// Modul LAIN - nilai kolom MODUL milik Claim Life. Refactor bentuk B:
	// literal, sebab uji modul ini tidak boleh mengimpor Claim Life.
	b.Modul = "CLAIMLIFE"
	if err := p.Laksanakan(context.Background(), nil, b); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("baris modul lain: %v", err)
	}
	b = barisKomite(t, "asing", layanan.KunciLayanan{})
	if err := p.Laksanakan(context.Background(), nil, b); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("jenis asing: %v", err)
	}
}
