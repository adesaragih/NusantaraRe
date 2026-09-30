package services_test

// Uji layanan kaskade hapus - TANPA Oracle (tiket 10).

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
)

type kaskadeUji struct {
	dampak, terhapus repository.DampakHapusTCO
	dihapus          int
	kombinasi        models.KombinasiTCO
}

func (k *kaskadeUji) DampakKontrak(_ context.Context, kom models.KombinasiTCO, _, _ string) (repository.DampakHapusTCO, error) {
	k.kombinasi = kom
	return k.dampak, nil
}
func (k *kaskadeUji) HapusKontrak(context.Context, *db.Tx, models.KombinasiTCO, string, string) (repository.DampakHapusTCO, error) {
	k.dihapus++
	return k.terhapus, nil
}
func (k *kaskadeUji) DampakReinsurer(context.Context, string) (repository.DampakHapusTCO, error) {
	return k.dampak, nil
}
func (k *kaskadeUji) HapusReinsurer(context.Context, *db.Tx, models.KombinasiTCO, string) (repository.DampakHapusTCO, error) {
	k.dihapus++
	return k.terhapus, nil
}

func layananKaskade(k *kaskadeUji, dikunci *int) *services.KaskadeTCO {
	return services.New(nil).KaskadeTCO().DenganKaskade(k).DenganKontrak(kontrakPemegangUji{dikunci: dikunci}).
		DenganTahun(tahunReinsurerUji{}).DenganReinsurer(reinsurerIndukSec()).
		DenganTransaksi(transaksiUji).DenganJam(jamUji)
}

func dampakUji() repository.DampakHapusTCO {
	return repository.DampakHapusTCO{Kontrak: 1, Reinsurer: 2, Security: 3, Business: 1, KlausulTetap: 5}
}

// AC 43/44: popup menyebut jumlah tiap jenis DAN klausul yang tetap hidup.
func TestDampakHapusKontrak(t *testing.T) {
	k := &kaskadeUji{dampak: dampakUji()}
	d, err := layananKaskade(k, new(int)).DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003")
	if err != nil || d.Reinsurer != 2 || d.Security != 3 || d.Business != 1 || d.KlausulTetap != 5 {
		t.Errorf("dampak: %+v %v", d, err)
	}
	if k.kombinasi.TreatyYear != "2026" || k.kombinasi.TreatyGroupID != "10001" || k.kombinasi.ReinsTypeID != "10003" {
		t.Errorf("kombinasi: %+v", k.kombinasi)
	}
	if _, err := layananKaskade(k, new(int)).DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000099"); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak asing: %v", err)
	}
}

// AC 42, 45: kaskade dalam satu transaksi, kontrak dikunci; pesan VERBATIM.
// tco4: nol jejak modul (Pega tidak mencatatnya).
func TestHapusKontrakKaskade(t *testing.T) {
	k, dikunci := &kaskadeUji{dampak: dampakUji(), terhapus: dampakUji()}, 0
	pesan, err := layananKaskade(k, &dikunci).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1})
	if err != nil || pesan != "Data successfully deleted" || k.dihapus != 1 || dikunci != 1 {
		t.Fatalf("hapus: %q %v dihapus %d dikunci %d", pesan, err, k.dihapus, dikunci)
	}
}

// Batal di popup = tidak ada permintaan hapus. Angka yang berubah sejak popup
// ditolak SEBELUM menghapus; angka terhapus yang berbeda membatalkan transaksi.
func TestHapusKontrakAngkaHarusSamaDenganPopup(t *testing.T) {
	k := &kaskadeUji{dampak: dampakUji(), terhapus: dampakUji()}
	_, err := layananKaskade(k, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 2, Business: 1})
	if !errors.Is(err, services.ErrDampakBerubah) || k.dihapus != 0 {
		t.Errorf("berubah sebelum hapus: %v dihapus %d", err, k.dihapus)
	}
	lain := dampakUji()
	lain.Business = 2
	k2 := &kaskadeUji{dampak: dampakUji(), terhapus: lain}
	if _, err := layananKaskade(k2, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1}); !errors.Is(err, services.ErrDampakBerubah) {
		t.Errorf("terhapus berbeda dari hitungan: %v", err)
	}
	if _, err := layananKaskade(k, new(int)).HapusKontrak(context.Background(), inti.Pelaku{}, "1000001", "1000003",
		services.KonfirmasiHapus{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("identitas: %v", err)
	}
}

func TestHapusReinsurerKaskade(t *testing.T) {
	d := repository.DampakHapusTCO{Reinsurer: 1, Security: 2}
	k := &kaskadeUji{dampak: d, terhapus: d}
	l := layananKaskade(k, new(int))
	dt, err := l.DampakHapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007")
	if err != nil || dt.Security != 2 {
		t.Fatalf("dampak reinsurer: %+v %v", dt, err)
	}
	if _, err := l.DampakHapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000099"); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("reinsurer asing: %v", err)
	}
	pesan, err := l.HapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", services.KonfirmasiHapus{Reinsurer: 1, Security: 2})
	if err != nil || !strings.Contains(pesan, "1000007") {
		t.Errorf("hapus reinsurer: %q %v", pesan, err)
	}
	if _, err := l.HapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", services.KonfirmasiHapus{Reinsurer: 1, Security: 1}); !errors.Is(err, services.ErrDampakBerubah) {
		t.Errorf("security berubah: %v", err)
	}
	if _, err := services.New(nil).KaskadeTCO().DenganTahun(tahunReinsurerUji{}).DenganKontrak(kontrakPemegangUji{dikunci: new(int)}).
		DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003"); !errors.Is(err, services.ErrGudangKaskadeBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
}

// OQ-TCO-21: cacah kontrak lain terdampak ikut dikonfirmasi.
func TestHapusKontrakBersamaDikonfirmasi(t *testing.T) {
	d := dampakUji()
	d.Bersama = 2
	k := &kaskadeUji{dampak: d, terhapus: d}
	if _, err := layananKaskade(k, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1}); !errors.Is(err, services.ErrDampakBerubah) || k.dihapus != 0 {
		t.Errorf("kontrak lain tidak dikonfirmasi: %v (dihapus %d)", err, k.dihapus)
	}
	if _, err := layananKaskade(k, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1, Bersama: 2}); err != nil || k.dihapus != 1 {
		t.Fatalf("hapus bersama: %v", err)
	}
}
