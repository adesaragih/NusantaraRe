package services_test

// Uji layanan kaskade hapus - TANPA Oracle (tiket 10).

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

type kaskadeUji struct {
	dampak, terhapus repository.DampakHapusTCO
	dihapus          int
	kombinasi        models.KombinasiTCO
}

func (k *kaskadeUji) DampakKontrak(_ context.Context, kom models.KombinasiTCO, _ string) (repository.DampakHapusTCO, error) {
	k.kombinasi = kom
	return k.dampak, nil
}
func (k *kaskadeUji) HapusKontrak(context.Context, *repository.Tx, models.KombinasiTCO, string, string) (repository.DampakHapusTCO, error) {
	k.dihapus++
	return k.terhapus, nil
}
func (k *kaskadeUji) DampakReinsurer(context.Context, string) (repository.DampakHapusTCO, error) {
	return k.dampak, nil
}
func (k *kaskadeUji) HapusReinsurer(context.Context, *repository.Tx, models.KombinasiTCO, string) (repository.DampakHapusTCO, error) {
	k.dihapus++
	return k.terhapus, nil
}

type jejakKaskadeUji struct{ baris []string }

func (j *jejakKaskadeUji) rekam(_ context.Context, _ *repository.Tx, akun, tabel, baris, aksi, ket string, _ time.Time) error {
	j.baris = append(j.baris, strings.Join([]string{akun, tabel, baris, aksi, ket}, "|"))
	return nil
}

func layananKaskade(k *kaskadeUji, j *jejakKaskadeUji, dikunci *int) *services.KaskadeTCO {
	return services.New(nil).KaskadeTCO().DenganKaskade(k).DenganKontrak(kontrakPemegangUji{dikunci: dikunci}).
		DenganTahun(tahunReinsurerUji{}).DenganReinsurer(reinsurerIndukSec()).DenganJejak(j.rekam).
		DenganTransaksi(transaksiUji).DenganJam(jamUji)
}

func dampakUji() repository.DampakHapusTCO {
	return repository.DampakHapusTCO{Kontrak: 1, Reinsurer: 2, Security: 3, Business: 1, KlausulTetap: 5}
}

// AC 43/44: popup menyebut jumlah tiap jenis DAN klausul yang tetap hidup.
func TestDampakHapusKontrak(t *testing.T) {
	k := &kaskadeUji{dampak: dampakUji()}
	d, err := layananKaskade(k, &jejakKaskadeUji{}, new(int)).DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003")
	if err != nil || d.Reinsurer != 2 || d.Security != 3 || d.Business != 1 || d.KlausulTetap != 5 {
		t.Errorf("dampak: %+v %v", d, err)
	}
	if k.kombinasi.TreatyYear != "2026" || k.kombinasi.TreatyGroupID != "10001" || k.kombinasi.ReinsTypeID != "10003" {
		t.Errorf("kombinasi: %+v", k.kombinasi)
	}
	if _, err := layananKaskade(k, &jejakKaskadeUji{}, new(int)).DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000099"); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak asing: %v", err)
	}
}

// AC 42, 45, jejak: kaskade dalam satu transaksi, kontrak dikunci, jejak
// menyebut jumlah tiap jenis dan klausul yang tidak disentuh; pesan VERBATIM.
func TestHapusKontrakKaskade(t *testing.T) {
	k, j, dikunci := &kaskadeUji{dampak: dampakUji(), terhapus: dampakUji()}, &jejakKaskadeUji{}, 0
	pesan, err := layananKaskade(k, j, &dikunci).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1})
	if err != nil || pesan != "Data Berhasil di Hapus" || k.dihapus != 1 || dikunci != 1 {
		t.Fatalf("hapus: %q %v dihapus %d dikunci %d", pesan, err, k.dihapus, dikunci)
	}
	if len(j.baris) != 1 || !strings.Contains(j.baris[0], "T_TREATYCONTRACT|1000003|hapus|") ||
		!strings.Contains(j.baris[0], "2 reinsurer, 3 security, 1 business; 5 klausul tidak disentuh") {
		t.Errorf("jejak: %v", j.baris)
	}
}

// Batal di popup = tidak ada permintaan hapus. Angka yang berubah sejak popup
// ditolak SEBELUM menghapus; angka terhapus yang berbeda membatalkan transaksi.
func TestHapusKontrakAngkaHarusSamaDenganPopup(t *testing.T) {
	k, j := &kaskadeUji{dampak: dampakUji(), terhapus: dampakUji()}, &jejakKaskadeUji{}
	_, err := layananKaskade(k, j, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 2, Business: 1})
	if !errors.Is(err, services.ErrDampakBerubah) || k.dihapus != 0 || len(j.baris) != 0 {
		t.Errorf("berubah sebelum hapus: %v dihapus %d", err, k.dihapus)
	}
	lain := dampakUji()
	lain.Business = 2
	k2 := &kaskadeUji{dampak: dampakUji(), terhapus: lain}
	if _, err := layananKaskade(k2, j, new(int)).HapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.KonfirmasiHapus{Reinsurer: 2, Security: 3, Business: 1}); !errors.Is(err, services.ErrDampakBerubah) || len(j.baris) != 0 {
		t.Errorf("terhapus berbeda dari hitungan: %v", err)
	}
	if _, err := layananKaskade(k, j, new(int)).HapusKontrak(context.Background(), services.Pelaku{}, "1000001", "1000003",
		services.KonfirmasiHapus{}); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("identitas: %v", err)
	}
}

func TestHapusReinsurerKaskade(t *testing.T) {
	d := repository.DampakHapusTCO{Reinsurer: 1, Security: 2}
	k, j := &kaskadeUji{dampak: d, terhapus: d}, &jejakKaskadeUji{}
	l := layananKaskade(k, j, new(int))
	dt, err := l.DampakHapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007")
	if err != nil || dt.Security != 2 {
		t.Fatalf("dampak reinsurer: %+v %v", dt, err)
	}
	if _, err := l.DampakHapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000099"); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("reinsurer asing: %v", err)
	}
	pesan, err := l.HapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", services.KonfirmasiHapus{Reinsurer: 1, Security: 2})
	if err != nil || !strings.Contains(pesan, "1000007") || !strings.Contains(j.baris[0], "T_TREATYREINSURER|1000007|hapus|reinsurer dihapus beserta 2 security") {
		t.Errorf("hapus reinsurer: %q %v %v", pesan, err, j.baris)
	}
	if _, err := l.HapusReinsurer(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", services.KonfirmasiHapus{Reinsurer: 1, Security: 1}); !errors.Is(err, services.ErrDampakBerubah) {
		t.Errorf("security berubah: %v", err)
	}
	if _, err := services.New(nil).KaskadeTCO().DenganTahun(tahunReinsurerUji{}).DenganKontrak(kontrakPemegangUji{dikunci: new(int)}).
		DampakHapusKontrak(context.Background(), pelakuUjiTCO, "1000001", "1000003"); !errors.Is(err, services.ErrGudangKaskadeBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
}
