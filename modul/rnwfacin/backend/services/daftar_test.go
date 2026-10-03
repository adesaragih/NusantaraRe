package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/rnwfacin/backend/models"
)

// sumberUji - repository.SumberDaftarRenewal dari irisan baris. Seluruh nilai
// sintetis: tidak ada nomor kasus, nama, atau login nyata.
type sumberUji []models.BarisDaftarRenewal

func (s sumberUji) KasusRenewal(context.Context) ([]models.BarisDaftarRenewal, error) { return s, nil }

type sumberGagal struct{}

func (sumberGagal) KasusRenewal(context.Context) ([]models.BarisDaftarRenewal, error) {
	return nil, errors.New("sumber mati")
}

var (
	t0          = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	saringanUji = Saringan{PembuatID: "UJI-op-a", TeamGroup: "1"}
)

func baris(id string, menit int, ubah func(*models.BarisDaftarRenewal)) models.BarisDaftarRenewal {
	b := models.BarisDaftarRenewal{IDKasus: id, PembuatID: "UJI-op-a", TeamGroup: "1", StatusKerja: "Pending-Approval",
		DibuatPada: t0.Add(time.Duration(menit) * time.Minute)}
	if ubah != nil {
		ubah(&b)
	}
	return b
}

func idKasus(daftar []models.BarisDaftarRenewal) []string {
	var hasil []string
	for _, b := range daftar {
		hasil = append(hasil, b.IDKasus)
	}
	return hasil
}

// TestDaftarRenewalSaringan - RenewalList_RD filter `B AND A AND C AND D`:
// pembuat = Param.UserNameID, TeamGroup = Param.TeamGroup, status kerja bukan
// Resolved-Completed dan bukan Resolved-Rejected.
func TestDaftarRenewalSaringan(t *testing.T) {
	sumber := sumberUji{
		baris("UJI-R-1", 1, nil),
		baris("UJI-R-2", 2, func(b *models.BarisDaftarRenewal) { b.PembuatID = "UJI-op-b" }),
		baris("UJI-R-3", 3, func(b *models.BarisDaftarRenewal) { b.TeamGroup = "2" }),
		baris("UJI-R-4", 4, func(b *models.BarisDaftarRenewal) { b.StatusKerja = models.StatusKerjaSelesai }),
		baris("UJI-R-5", 5, func(b *models.BarisDaftarRenewal) { b.StatusKerja = models.StatusKerjaDitolak }),
		// Pembanding persis: beda huruf besar-kecil bukan pembuat yang sama.
		baris("UJI-R-7", 7, func(b *models.BarisDaftarRenewal) { b.PembuatID = "UJI-OP-A" }),
	}
	got, err := NewServiceDaftar(sumber).Daftar(context.Background(), saringanUji)
	if err != nil {
		t.Fatal(err)
	}
	if mau := []string{"UJI-R-1"}; !reflect.DeepEqual(idKasus(got.Baris), mau) || got.Terpotong {
		t.Errorf("dapat %v (terpotong %v), mau %v", idKasus(got.Baris), got.Terpotong, mau)
	}
}

// TestDaftarRenewalUrutan - pxCreateDateTime DESC (urutan 1), lalu pyID DESC
// (urutan 2), dibandingkan sebagai teks.
func TestDaftarRenewalUrutan(t *testing.T) {
	sumber := sumberUji{baris("UJI-R-10", 1, nil), baris("UJI-R-9", 1, nil), baris("UJI-R-11", 3, nil), baris("UJI-R-2", 2, nil)}
	got, err := NewServiceDaftar(sumber).Daftar(context.Background(), saringanUji)
	if err != nil {
		t.Fatal(err)
	}
	// Menit 1 seri: "UJI-R-9" > "UJI-R-10" sebagai teks.
	if mau := []string{"UJI-R-11", "UJI-R-2", "UJI-R-9", "UJI-R-10"}; !reflect.DeepEqual(idKasus(got.Baris), mau) {
		t.Errorf("dapat %v, mau %v", idKasus(got.Baris), mau)
	}
}

// TestDaftarRenewalBatas - pyMaxRecords 500: sesudah diurutkan, 500 teratas.
func TestDaftarRenewalBatas(t *testing.T) {
	var sumber sumberUji
	for i := 0; i < 501; i++ {
		sumber = append(sumber, baris(fmt.Sprintf("UJI-R-%04d", i), i, nil))
	}
	got, err := NewServiceDaftar(sumber).Daftar(context.Background(), saringanUji)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Baris) != 500 || !got.Terpotong || got.Baris[0].IDKasus != "UJI-R-0500" || got.Baris[499].IDKasus != "UJI-R-0001" {
		t.Errorf("%d baris, terpotong %v, pertama %s", len(got.Baris), got.Terpotong, got.Baris[0].IDKasus)
	}
	got, _ = NewServiceDaftar(sumber[:500]).Daftar(context.Background(), saringanUji)
	if len(got.Baris) != 500 || got.Terpotong {
		t.Errorf("tepat 500: %d baris, terpotong %v", len(got.Baris), got.Terpotong)
	}
}

// TestDaftarRenewalDitolak - parameter kosong (perilaku Pega belum terverifikasi)
// dan galat sumber diteruskan.
func TestDaftarRenewalDitolak(t *testing.T) {
	s := NewServiceDaftar(sumberUji{baris("UJI-R-1", 1, nil)})
	for _, sr := range []Saringan{{TeamGroup: "1"}, {PembuatID: "UJI-op-a"}} {
		if _, err := s.Daftar(context.Background(), sr); !errors.Is(err, ErrSaringanKosong) {
			t.Errorf("%+v: galat %v", sr, err)
		}
	}
	// Status kerja kosong: SQL `<>` Oracle akan membuang baris itu (kosong = NULL), sedangkan
	// pembanding Go akan meloloskannya - belum terverifikasi, jadi ditolak.
	kosong := sumberUji{baris("UJI-R-1", 1, func(b *models.BarisDaftarRenewal) { b.StatusKerja = "" })}
	if _, err := NewServiceDaftar(kosong).Daftar(context.Background(), saringanUji); !errors.Is(err, ErrStatusKerjaKosong) {
		t.Errorf("status kosong: galat %v", err)
	}
	// Pesan galat tidak memuat login pembuat (CLAUDE.md §4.10).
	if _, err := s.Daftar(context.Background(), Saringan{PembuatID: "UJI-op-a"}); err == nil || strings.Contains(err.Error(), "UJI-op-a") {
		t.Errorf("pesan galat: %v", err)
	}
	if _, err := NewServiceDaftar(sumberGagal{}).Daftar(context.Background(), saringanUji); err == nil {
		t.Error("galat sumber hilang")
	}
}
