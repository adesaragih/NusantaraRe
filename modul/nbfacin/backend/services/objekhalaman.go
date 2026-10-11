package services

// Tab Object FIRE untuk case besar (tiket 75 F3, K75-3): dibaca per halaman dari server, bukan utuh. Layar memilih mode dari
// `jumlahLokasi` GET kasus (usulan ambang AmbangHalamanObjek); di bawah ambang GET / PUT …/objek utuh tetap dipakai apa adanya.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// AmbangHalamanObjek - tab Object per halaman bila jumlah lokasi > ambang; WAJIB sama dengan frontend api.ts
// AMBANG_HALAMAN_OBJEK (dijaga TestAmbangHalamanObjekSamaDenganApiTs).
const AmbangHalamanObjek = 500

// DenganHalamanObjek memasang penyimpan halaman tab Object (tiket 75 F3).
func (s *Service) DenganHalamanObjek(h repository.PenyimpanHalamanObjek) *Service {
	s.halamanObjek = h
	return s
}

// JumlahLokasi - jumlah lokasi objek FIRE case; nil bila tidak terbaca (tanpa basis data / case tidak ada / galat - dicatat).
func (s *Service) JumlahLokasi(ctx context.Context, id string) *int {
	if s.halamanObjek == nil || !idKasusSah(id) {
		return nil
	}
	n, err := s.halamanObjek.JumlahLokasi(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrKasusTidakAda) {
			log.Printf("nbfacin: jumlah lokasi kasus %s: %v", id, err)
		}
		return nil
	}
	return &n
}

// BacaObjekHalaman - GET …/objek/halaman: halaman >= 1, ukuran 1..BatasUkuranHalamanObjek (400); lokasi diolah seperti BacaObjek.
func (s *Service) BacaObjekHalaman(ctx context.Context, id string, halaman, ukuran int, cari string) (models.HalamanObjek, error) {
	if halaman < 1 || ukuran < 1 || ukuran > repository.BatasUkuranHalamanObjek {
		return models.HalamanObjek{}, fmt.Errorf("%w: halaman >= 1, ukuran 1..%d", ErrMasukanObjek, repository.BatasUkuranHalamanObjek)
	}
	if s.halamanObjek == nil {
		return models.HalamanObjek{}, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return models.HalamanObjek{}, ErrKasusTidakAda
	}
	h, err := s.halamanObjek.BacaObjekHalaman(ctx, id, halaman, ukuran, cari)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return models.HalamanObjek{}, ErrKasusTidakAda
	}
	if err != nil {
		return models.HalamanObjek{}, err
	}
	siapkanObjekBaca(h.Lokasi)
	return h, nil
}

var (
	// ErrLokasiObjekBerubah - ID lokasi PATCH tidak ada lagi di case. 409.
	ErrLokasiObjekBerubah = errors.New("Lokasi objek sudah berubah sejak halaman dimuat - muat ulang halaman.")
	// ErrLokasiBerSpreading - lokasi yang diubah sudah ber-spreading (tiket 75 F3, spreading case besar = F4). 409.
	ErrLokasiBerSpreading = errors.New("Lokasi yang diubah sudah memiliki spreading; perubahan per halaman untuk lokasi ber-spreading " +
		"belum tersedia (tiket 75 F4).")
)

// polaIDLokasi - T_LOCATIONLIST.ID NUMBER(19).
var polaIDLokasi = regexp.MustCompile(`^[0-9]{1,19}$`)

// PerubahanObjek - badan PATCH …/objek: lokasi diubah (ber-ID), ditambah, dan ID lokasi dihapus.
type PerubahanObjek struct {
	Ubah   []models.LokasiUbah
	Tambah []models.ObjekFire
	Hapus  []string
}

// periksaPerubahan - isian PATCH (400): ada perubahan, batas jumlah, ID sah dan tidak ganda, aturan objek tab Object.
func periksaPerubahan(p PerubahanObjek) error {
	if len(p.Ubah)+len(p.Tambah)+len(p.Hapus) == 0 {
		return fmt.Errorf("%w: tidak ada perubahan", ErrMasukanObjek)
	}
	if n := len(p.Ubah) + len(p.Hapus); n > repository.BatasLokasiPerubahan {
		return fmt.Errorf("%w: paling banyak %d lokasi diubah / dihapus sekali simpan", ErrMasukanObjek, repository.BatasLokasiPerubahan)
	}
	if len(p.Tambah) > repository.BatasLokasiPerubahan {
		return fmt.Errorf("%w: paling banyak %d lokasi ditambah sekali simpan", ErrMasukanObjek, repository.BatasLokasiPerubahan)
	}
	lihat := map[string]bool{}
	cek := func(x string) error {
		if !polaIDLokasi.MatchString(x) {
			return fmt.Errorf("%w: id lokasi %q tidak sah", ErrMasukanObjek, x)
		}
		if lihat[x] {
			return fmt.Errorf("%w: id lokasi %s ganda", ErrMasukanObjek, x)
		}
		lihat[x] = true
		return nil
	}
	ubah := make([]models.ObjekFire, len(p.Ubah))
	for i, u := range p.Ubah {
		if err := cek(u.ID); err != nil {
			return err
		}
		ubah[i] = u.Objek
	}
	for _, x := range p.Hapus {
		if err := cek(x); err != nil {
			return err
		}
	}
	if err := periksaObjek(ubah); err != nil {
		return fmt.Errorf("ubah: %w", err)
	}
	if err := periksaObjek(p.Tambah); err != nil {
		return fmt.Errorf("tambah: %w", err)
	}
	return nil
}

// PerbaruiObjek - PATCH …/objek (tiket 75 F3): urutan galat sama dengan GantiObjek (identitas 401 -> hak ubah 409 -> isian 400
// -> basis data 503 -> kunci unggah 409 -> mata uang 400 -> loss ratio / coverage), lalu satu transaksi (repository
// PerbaruiObjek). Mengembalikan jumlah lokasi sesudahnya; halaman dibaca ulang pemanggil.
func (s *Service) PerbaruiObjek(ctx context.Context, pelaku inti.Pelaku, id string, p PerubahanObjek) (int, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return 0, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return 0, err
	}
	if err := periksaPerubahan(p); err != nil {
		return 0, err
	}
	if s.halamanObjek == nil || s.transaksi == nil {
		return 0, ErrObjekTanpaDatabase
	}
	if err := s.cekKunciUnggah(ctx, id); err != nil {
		return 0, err
	}
	if !idKasusSah(id) {
		return 0, ErrKasusTidakAda
	}
	baris := make([]models.ObjekFire, 0, len(p.Ubah)+len(p.Tambah))
	for _, u := range p.Ubah {
		baris = append(baris, u.Objek)
	}
	baris = append(baris, p.Tambah...)
	if err := s.periksaMataUang(ctx, baris); err != nil {
		return 0, err
	}
	simpan, err := s.siapkanKerugian(ctx, id, baris)
	if err != nil {
		return 0, err
	}
	if err := s.siapkanCoverageDengan(ctx, id, simpan, nil); err != nil {
		return 0, err
	}
	ubah := make([]models.LokasiUbah, len(p.Ubah))
	for i, u := range p.Ubah {
		ubah[i] = models.LokasiUbah{ID: u.ID, Objek: simpan[i]}
	}
	var jumlah int
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		var err error
		jumlah, err = s.halamanObjek.PerbaruiObjek(ctx, tx, id, ubah, simpan[len(p.Ubah):], p.Hapus)
		return err
	})
	switch {
	case errors.Is(err, repository.ErrKasusTidakAda):
		return 0, ErrKasusTidakAda
	case errors.Is(err, repository.ErrLokasiBukanKasus):
		return 0, ErrLokasiObjekBerubah
	case errors.Is(err, repository.ErrLokasiBerSpreading):
		return 0, ErrLokasiBerSpreading
	}
	return jumlah, err
}
