package services

// Keputusan satu tingkat Komite - tiket 02. TANPA Oracle.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/repository"
)

// TestGiliranHanyaAnggotaBerjalan - ADR-0014, AC 1 spec.
func TestGiliranHanyaAnggotaBerjalan(t *testing.T) {
	k := kasusUji() // count berjalan 2 di kasusUji? - disetel di sini
	k.Baris.KomiteCount, k.Baris.KomiteLoop = 2, 3
	if err := periksaGiliran(k, "UJI-B"); err != nil {
		t.Errorf("anggota berjalan ditolak: %v", err)
	}
	for _, akun := range []string{"UJI-A", "UJI-C", "UJI-LUAR", ""} {
		if err := periksaGiliran(k, akun); !errors.Is(err, ErrTanpaWewenang) {
			t.Errorf("%q: %v, mau ErrTanpaWewenang", akun, err)
		}
	}
}

// TestTanggaBerhentiTidakMenerimaKeputusan - Tolak menghentikan (AC 6).
func TestTanggaBerhentiTidakMenerimaKeputusan(t *testing.T) {
	k := kasusUji()
	k.Baris.KomiteCount, k.Baris.KomiteLoop, k.Baris.AcceptStatus = 2, 3, "2"
	if err := periksaGiliran(k, "UJI-B"); !errors.Is(err, ErrTanggaKomiteBerhenti) {
		t.Errorf("tangga berhenti: %v", err)
	}
	k.Baris.AcceptStatus, k.Baris.StatusWork = "1", "Resolved-Completed"
	if err := periksaGiliran(k, "UJI-B"); !errors.Is(err, ErrTanggaKomiteBerhenti) {
		t.Errorf("kasus tertutup: %v", err)
	}
}

// TestTingkatSudahDiputuskanTidakDiulang - dua klik, satu keputusan.
func TestTingkatSudahDiputuskanTidakDiulang(t *testing.T) {
	k := kasusUji()
	k.Baris.KomiteCount, k.Baris.KomiteLoop = 1, 3
	if err := periksaGiliran(k, "UJI-A"); !errors.Is(err, ErrKeputusanKomiteBersamaan) {
		t.Errorf("tingkat sudah berkode 1 masih diterima: %v", err)
	}
}

// TestKeputusanAsingDitolakSebelumBasisData - AC 35 spec.
func TestKeputusanAsingDitolakSebelumBasisData(t *testing.T) {
	kk := New(nil).KeputusanKomite()
	p := Pelaku{AkunID: "UJI"}
	if _, err := kk.Putuskan(context.Background(), p, "K", "3", "", time.Now()); !errors.Is(err, ErrKeputusanKomiteTidakDikenal) {
		t.Errorf("keputusan 3: %v", err)
	}
	if _, err := kk.Putuskan(context.Background(), p, "K", "1", "", time.Now()); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("keputusan sah tanpa Oracle: %v", err)
	}
}

// TestTingkatAkhirBawaanGagalTerang - langkah 4/5 belum dibangun.
func TestTingkatAkhirBawaanGagalTerang(t *testing.T) {
	var p PenyelesaiAkhirKomite = PenyelesaiAkhirBelumAda{}
	if _, err := p.Akseptasi(context.Background(), nil, repository.KasusKomite{}, Pelaku{}, time.Now()); !errors.Is(err, ErrPenyelesaianAkhirBelumAda) {
		t.Errorf("akseptasi bawaan: %v", err)
	}
	if err := p.Tolak(context.Background(), nil, repository.KasusKomite{}, Pelaku{}, time.Now()); !errors.Is(err, ErrPenyelesaianAkhirBelumAda) {
		t.Errorf("tolak bawaan: %v", err)
	}
}

// TestUrutanKeputusanDalamSatuTransaksi - langkah 3/13, lalu 4/5, lalu jejak.
func TestUrutanKeputusanDalamSatuTransaksi(t *testing.T) {
	isi, err := os.ReadFile("komite_keputusan.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := badanFungsi(t, string(isi), "func (k *KeputusanKomite) Putuskan(")
	urut := []string{"PastikanKasusTerbuka(ctx, klaimID)", "periksaGiliran(", "DalamTransaksi(ctx",
		"baca.CatatKeputusan(", "k.akhir.Akseptasi(", "k.akhir.Tolak(", "antreEfekKomite(", "k.jejak.Rekam("}
	lalu := -1
	for _, s := range urut {
		i := strings.Index(badan, s)
		if i < 0 || i < lalu {
			t.Errorf("urutan Putuskan salah di %q (urutan: %v)", s, urut)
		}
		lalu = i
	}
	if n := strings.Count(badan, "DalamTransaksi("); n != 1 {
		t.Errorf("Putuskan membuka %d transaksi, mau 1", n)
	}
}

// TestEskalasiHanyaAdmin - AC 11 spec Komite.
func TestEskalasiHanyaAdmin(t *testing.T) {
	kk := New(nil).KeputusanKomite()
	_, err := kk.Eskalasi(context.Background(), Pelaku{AkunID: "UJI-B"}, "K", time.Now())
	if !errors.Is(err, ErrTanpaWewenang) {
		t.Errorf("eskalasi bukan admin: %v", err)
	}
	_, err = kk.Eskalasi(context.Background(), Pelaku{AkunID: "UJI-ADM", Peran: []string{PeranAdmin}}, "K", time.Now())
	if !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("eskalasi admin tanpa Oracle: %v", err)
	}
}

// TestEskalasiBukanPintuBelakangKeputusan - admin tetap bukan anggota berjalan.
func TestEskalasiBukanPintuBelakangKeputusan(t *testing.T) {
	k := kasusUji()
	k.Baris.KomiteCount, k.Baris.KomiteLoop = 2, 3
	if err := periksaGiliran(k, "UJI-ADM"); !errors.Is(err, ErrTanpaWewenang) {
		t.Errorf("admin dapat memutuskan atas nama tingkat berjalan: %v", err)
	}
}
