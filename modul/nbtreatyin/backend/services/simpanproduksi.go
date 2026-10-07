package services

// Untuk apa berkas ini: alat `simpanproduksi` - baris Utility1 (json_polis tanpa DATA_JSON, ACHIEVEMENT,
// TREATYINPRODUCTION; models/produksi.go) untuk kasus yang SUDAH selesai sebelum penulisan produksi ada
// (permintaan work owner 06-10-2026: "bagaimana cara insert ulang ke json_polis dan treatyinproduction").
// Dijalankan manusia dari baris perintah, tidak pernah oleh aplikasi; seam services tidak diuji (spec.md §6.2),
// sama dengan `Pemuat`.
//
// Bedanya dengan submit: ProductionDate yang tersimpan dipakai apa adanya (tanggalnya sudah tetap saat nomor
// polis terbit), halaman kasus TIDAK disimpan ulang (hanya salinan yang dipra-proses), dan ACHIEVEMENT dilewati
// bila IDPEGA sudah punya baris - alat yang dijalankan dua kali tidak menggandakannya. json_polis dan
// TREATYINPRODUCTION sudah dilewati bila ada (`repository.SimpanPolisProduksi`).

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// PenulisProduksi - penulis baris Utility1 atas Oracle.
type PenulisProduksi struct{ g penyimpanOracle }

// PenulisProduksiDariDasar menyusun penulis di atas basis data bersama.
func PenulisProduksiDariDasar(d *inti.Dasar) (*PenulisProduksi, error) {
	if d == nil || !d.PunyaDatabase() {
		return nil, ErrTanpaOracle
	}
	return &PenulisProduksi{g: penyimpanOracle{Gudang: repository.Baru(d.DB()), dasar: d}}, nil
}

// HasilProduksi - baris yang disusun, jumlah baris yang SUDAH ada sebelum alat jalan, dan apakah ditulis.
type HasilProduksi struct {
	Simpanan models.SimpananPolis
	Ada      models.CacahProduksi
	Pengguna string
	Ditulis  bool
}

// ErrBelumSelesaiBernomor - alat hanya untuk kasus Resolved-Completed yang bernomor polis.
var ErrBelumSelesaiBernomor = errors.New("services: kasus belum selesai atau belum bernomor polis")

// Jalankan menyusun (dan bila `tulis`, menulis) baris Utility1 satu kasus. `pengguna` kosong = OPERATORID
// riwayat terakhir kasus (yang menekan Submit penutup) - USERNAME json_polis (`OperatorID.pyUserIdentifier`).
func (p *PenulisProduksi) Jalankan(ctx context.Context, id, pengguna string, tulis bool) (HasilProduksi, error) {
	k, err := p.g.Keadaan(ctx, nil, id)
	if err != nil {
		return HasilProduksi{}, err
	}
	if k.StatusWork != models.StatusSelesai || k.NoPolis == "" {
		return HasilProduksi{}, fmt.Errorf("%w: %s status %q, nomor polis %q", ErrBelumSelesaiBernomor, id, k.StatusWork, k.NoPolis)
	}
	h, err := p.g.BacaHalaman(ctx, nil, id)
	if err != nil {
		return HasilProduksi{}, err
	}
	if pengguna == "" {
		r, err := p.g.DaftarRiwayat(ctx, models.KunciInstans(id))
		if err != nil {
			return HasilProduksi{}, err
		}
		if len(r) > 0 {
			pengguna = r[len(r)-1].OperatorID
		}
	}
	if pengguna == "" {
		return HasilProduksi{}, fmt.Errorf("services: %s tanpa riwayat - isi -pengguna", id)
	}
	models.PrasimpanMedan(h, id)
	s, err := models.SusunSimpananPolis(h, id, pengguna)
	if err != nil {
		return HasilProduksi{}, err
	}
	hasil := HasilProduksi{Simpanan: s, Pengguna: pengguna}
	if !tulis {
		hasil.Ada, err = p.g.CacahProduksi(ctx, nil, s.IDPega)
		return hasil, err
	}
	err = p.g.Transaksi(ctx, func(tx *db.Tx) error {
		ada, err := p.g.CacahProduksi(ctx, tx, s.IDPega)
		if err != nil {
			return err
		}
		hasil.Ada = ada
		if ada.Capaian > 0 {
			s.Capaian = nil
		}
		return p.g.SimpanPolisProduksi(ctx, tx, s)
	})
	if err != nil {
		return HasilProduksi{}, err
	}
	hasil.Simpanan, hasil.Ditulis = s, true
	return hasil, nil
}
