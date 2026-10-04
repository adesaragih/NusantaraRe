package services

// Pencarian diagnosa — kelompok Medis.
//
// Untuk apa berkas ini: tombol `Find Disease` (`ClaimLifeDetailGCNM.xml`
// b5061) membuka `Diagnose_Harness`, dan pencariannya berjalan atas
// `DISEASE_LIFE` yang **97.586 baris**.
//
// ⛔ PERANNYA TIDAK DISEMPITKAN, dan itu RALAT atas komentar berkas ini
// sendiri (27-09-2026). Ia dulu berbunyi "Perannya MEDIS ... gerbangnya ada
// DI SINI dan bukan di handler" sementara `Cari` di bawah hanya menuntut
// IDENTITAS - komentar yang menjanjikan penjaga yang tidak pernah ada.
// Bentuk cacat yang sama dengan `totalpeserta.go` sebelum diralat: yang
// berbahaya bukan penjaga yang hilang, melainkan kalimat yang membuat
// pembaca berhenti mencarinya.
//
// Yang benar menurut pohonnya, dan kenapa kode di bawah justru BENAR:
// `Find Disease` b5061 berdiri di `ClaimLifeDetailGCNM`, dan section itu
// dibuka `ViewClaimDetailLifeGCNM` - aksi sunting baris grid peserta di
// TIGA section: `InputOSClaimLife` b18252, `MedicalCheckClaimLife` b17416,
// `InputAkseptasiClaimLife` b17387, ketiganya atas grid yang sama
// (b15924, b15764, b15735). Ketiga pemuatnya ber-`pyWhenName` KOSONG dan
// ber-`pyPrivilegeName` KOSONG. Menyempitkannya ke Medical Advisor akan
// menutup pencarian bagi Admin dan SPV yang di sistem lama membukanya.
//
// Gerbang yang SUNGGUH ada atas diagnosa - tahap, pemegangnya, dan
// `STS_REJECT` peserta - berlaku atas rute yang MENGUBAH, dan tinggal di
// `diagnosa.go`. Membaca katalog tidak mengubah apa pun.
//
// Dibaca sesudah: dokumen.go.

import (
	"context"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// PencarianPenyakit melayani pencarian KATALOG penyakit.
//
// ⚠️ Namanya berubah 27-09-2026. Ia lahir sebagai `Diagnosa`, dan nama
// itu kini dipakai hal yang berbeda: diagnosa yang benar-benar DILEKATKAN
// pada seorang peserta (`diagnosa.go`). Satu kata untuk katalog dan untuk
// isinya adalah satu kata yang harus dibaca dua kali setiap kali.
type PencarianPenyakit struct{ svc *Service }

// Penyakit menyusun layanannya.
func (s *Service) Penyakit() *PencarianPenyakit { return &PencarianPenyakit{svc: s} }

// Cari mengembalikan diagnosa yang cocok, SELALU berbatas.
//
// ⚠️ Kriteria KOSONG diterima, dan itu ditiru apa adanya: di Pega
// `Contains ""` cocok dengan semua baris, dan yang menahannya adalah
// `pyMaxRecords` 500. Melarang pencarian kosong menutup jalan yang di sistem
// lama terbuka - menelusuri daftar tanpa tahu kata kuncinya.
func (d *PencarianPenyakit) Cari(ctx context.Context, pelaku inti.Pelaku, kodeICD, nama string,
	batas int) ([]models.Penyakit, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if d == nil || d.svc == nil || !d.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	k := models.NormalkanKriteriaPenyakit(kodeICD, nama)
	return repository.NewPenyakit(d.svc.DB()).Cari(ctx, k, batas)
}
