package services

// Pencarian diagnosa — kelompok Medis.
//
// Untuk apa berkas ini: tombol `Find Disease` (`ClaimLifeDetailGCNM.xml`
// b5061) membuka `Diagnose_Harness`, dan pencariannya berjalan atas
// `DISEASE_LIFE` yang **97.586 baris**.
//
// ⛔ Perannya MEDIS. `Find Disease` hidup di jalur Medical Check, dan
// ADR-U-0002 menyebut `ReasLifeMedicalAdvisor` sebagai penelaah medis.
// Gerbangnya ada DI SINI dan bukan di handler: layar bukan penjaga.
//
// Dibaca sesudah: dokumen.go.

import (
	"context"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// Diagnosa melayani pencarian penyakit.
type Diagnosa struct{ svc *Service }

// Penyakit menyusun layanannya.
func (s *Service) Penyakit() *Diagnosa { return &Diagnosa{svc: s} }

// Cari mengembalikan diagnosa yang cocok, SELALU berbatas.
//
// ⚠️ Kriteria KOSONG diterima, dan itu ditiru apa adanya: di Pega
// `Contains ""` cocok dengan semua baris, dan yang menahannya adalah
// `pyMaxRecords` 500. Melarang pencarian kosong menutup jalan yang di sistem
// lama terbuka - menelusuri daftar tanpa tahu kata kuncinya.
func (d *Diagnosa) Cari(ctx context.Context, pelaku Pelaku, kodeICD, nama string,
	batas int) ([]models.Penyakit, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if d == nil || d.svc == nil || !d.svc.PunyaDatabase() {
		return nil, repository.ErrTanpaOracle
	}
	k := models.NormalkanKriteriaPenyakit(kodeICD, nama)
	return repository.NewPenyakit(d.svc.db).Cari(ctx, k, batas)
}
