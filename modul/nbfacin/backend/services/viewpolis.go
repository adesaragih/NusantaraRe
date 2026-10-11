package services

// View Policy tab Facultative (List Offer) / Treaty / Life - tiket 71 perluasan (permintaan work owner 11-10-2026 "menu view policy ini
// buatin untuk life dan treaty juga"; kontrak sesi 9d). Bukti `[terverifikasi]` `D:\migrasi\RNM\DDL\GetListViewPolicy_act.xml`
// (bukan korpus) langkah 6 (kutip dibuang -> filterPolis), 7, 9-12; kode prakondisi dibaca per kode (2 lanjut, 3 lewati,
// 5 lewati when / jalankan). Langkah berurutan: yang belakangan MENIMPA CARI5 yang lebih dulu.
//   - `TRANS 1==1 -> SKP` (10.1, 10.3, 11.1-11.5): label SKP TIDAK ADA di activity (nol `pyStepsLabel`) - `[dugaan]` (keputusan sesi
//     9d) lompatan tak berefek, langkah berikutnya tetap dievaluasi. SATU penyimpangan sadar dari bacaan itu (keputusan work owner
//     11-10-2026): filter TEPAT 5 karakter -> 11.2 (5 digit terakhir) MENANG atas 11.6 (bacaan SKP akan membuat 11.6 menimpanya).
//     Prakondisi 11.2 `@length(filter)==5 T=(kosong) F=3`: benar -> lanjut (jalankan), salah -> lewati; tanpa syarat isi.
//   - langkah 7: filter kosong -> keluar hanya untuk POLIFAC / TREATYPOL / OFFERFAC; POLISLIFE / OFFERLIFE di Pega lanjut dengan
//     `LIKE '%%'` (semua baris). Di sini SEMUA tab: filter kosong -> daftar kosong (kontrak 9d; penyimpangan dicatat di tiket 71).
//   - OFFERFAC: Pega blok 9 mendeteksi isi (RNM -> NOOFFER, NB-/RNW- -> IDPEGA, selainnya SOB OR RISLIPCEDING); layar baru memakai
//     View Option seperti polis Facultative (kontrak 9d): 1 Case ID (NB-/RNW- seperti 9.2), 2 No Offer, 4 Reff, 5 Insured awalan,
//     6 SOB awalan; 3 / 7 / lainnya 400 `[keputusan agent]` (No Slip ditunda K70-3; QQ offer selalu '').
// Akses (K71-3): semua pengguna login.

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrBerkasTreatyTidakAda - nomor polis Treaty tanpa berkas NB / EDM Treaty In di sistem baru (polis Pega lama). 404.
var ErrBerkasTreatyTidakAda = errors.New("Polis ini belum punya berkas NB / EDM Treaty In di sistem baru.")

// DenganViewPolis - memasang pembaca tab Offer Facultative / Treaty / Life (tiket 71 perluasan).
func (s *Service) DenganViewPolis(p repository.PembacaViewPolis) *Service { s.viewPolis = p; return s }

// HalamanView - satu halaman tab View Policy (20 baris, seperti polis Facultative).
type HalamanView[T any] struct {
	Baris                  []T
	Total, Halaman, Ukuran int
}

// jenisCariOfferFac - View Option tab Offer (lihat kepala berkas) -> (jenis, nilai bind); jenis "" = tanpa kueri.
func jenisCariOfferFac(opsi, filter string) (string, string, error) {
	switch opsi {
	case "1", "2", "4", "5", "6":
	default:
		return "", "", errors.Join(ErrMasukanFollowing, errors.New("View Option Offer wajib 1, 2, 4, 5 atau 6"))
	}
	if filter == "" {
		return "", "", nil
	}
	switch opsi {
	case "1":
		if strings.Contains(filter, "NB-") || strings.Contains(filter, "RNW-") {
			return repository.OfferFacIDPega, filter, nil
		}
		return "", "", nil
	case "2":
		return repository.OfferFacNoOffer, filter, nil
	case "4":
		return repository.OfferFacReff, "%" + filter + "%", nil
	case "5":
		return repository.OfferFacInsured, filter + "%", nil
	}
	return repository.OfferFacSOB, filter + "%", nil
}

// jenisCariTreaty - blok 10 (FilterForViewPolicy 1-5) berurutan; @length = jumlah karakter.
func jenisCariTreaty(opsi, filter string) (string, string, error) {
	switch opsi {
	case "1", "2", "3", "4", "5":
	default:
		return "", "", errors.Join(ErrMasukanFollowing, errors.New("View Option Treaty wajib 1, 2, 3, 4 atau 5"))
	}
	if filter == "" {
		return "", "", nil // langkah 7
	}
	var jenis, nilai string
	if opsi == "2" { // 10.1
		jenis, nilai = repository.TreatyNoPolis, filter
	}
	if opsi == "2" && utf8.RuneCountInString(filter) == 28 { // 10.2
		jenis, nilai = repository.TreatyNoEndors, filter
	}
	if opsi == "1" && strings.Contains(filter, "NB-") { // 10.3
		jenis, nilai = repository.TreatyIDPega, filter
	}
	if opsi == "1" && strings.Contains(filter, "EDMT-") { // 10.4
		jenis, nilai = repository.TreatyEDMT, filter
	}
	switch opsi {
	case "3": // 10.5
		jenis, nilai = repository.TreatyMaster, filter
	case "4": // 10.6
		jenis, nilai = repository.TreatyInsured, "%"+filter+"%"
	case "5": // 10.7
		jenis, nilai = repository.TreatySOB, "%"+filter+"%"
	}
	return jenis, nilai, nil
}

// jenisCariLife - blok 11 (deteksi isi teks, tanpa View Option) berurutan. 11.2 / 11.4 / 11.5 `T=` kosong = lanjut (jalankan).
func jenisCariLife(filter string) (string, string) {
	if filter == "" {
		return "", ""
	}
	n := utf8.RuneCountInString(filter)
	rnm, nblf, edml := strings.Contains(filter, "RNM"), strings.Contains(filter, "NBLF-"), strings.Contains(filter, "EDML-")
	var jenis, nilai string
	if rnm && n <= 22 { // 11.1
		jenis, nilai = repository.LifeNoPolis, filter
	}
	if n == 5 { // 11.2
		jenis, nilai = repository.LifeDigitAkhir, "%"+filter
	}
	if nblf { // 11.3
		jenis, nilai = repository.LifeIDPega, filter
	}
	if n > 22 && rnm { // 11.4
		jenis, nilai = repository.LifeNoEndors, filter
	}
	if edml { // 11.5
		jenis, nilai = repository.LifeEDML, filter
	}
	if !rnm && !nblf && !edml && n != 5 { // 11.6; panjang 5 = 11.2 menang (keputusan work owner 11-10-2026)
		jenis, nilai = repository.LifeSOB, "%"+filter+"%"
	}
	return jenis, nilai
}

// jenisCariOfferLife - blok 12 (deteksi isi teks) berurutan; 12.2 nilai = teks utuh (SQL memotong SUBSTR(:1, 5)).
func jenisCariOfferLife(filter string) (string, string) {
	if filter == "" {
		return "", ""
	}
	nblf, ofr := strings.Contains(filter, "NBLF-"), strings.Contains(filter, "OFR-")
	var jenis, nilai string
	if nblf { // 12.1
		jenis, nilai = repository.OfferLifeIDPega, filter
	}
	if ofr { // 12.2
		jenis, nilai = repository.OfferLifeOFR, filter
	}
	if !nblf && !ofr { // 12.3
		jenis, nilai = repository.OfferLifeCeding, "%"+filter+"%"
	}
	return jenis, nilai
}

// cariView - identitas, halaman, jenis; jenis "" -> halaman kosong tanpa membaca basis data.
func cariView[T any](ctx context.Context, s *Service, pelaku inti.Pelaku, halaman int, jenis, nilai string, err error,
	baca func(r repository.PembacaViewPolis, offset, batas int) ([]T, int, error)) (HalamanView[T], error) {
	if e := inti.WajibIdentitas(pelaku); e != nil {
		return HalamanView[T]{}, e
	}
	if halaman < 1 {
		return HalamanView[T]{}, errors.Join(ErrMasukanFollowing, errors.New("halaman wajib >= 1"))
	}
	if err != nil {
		return HalamanView[T]{}, err
	}
	h := HalamanView[T]{Baris: []T{}, Halaman: halaman, Ukuran: ukuranHalamanFollowing}
	if jenis == "" {
		return h, nil
	}
	if s.viewPolis == nil {
		return HalamanView[T]{}, ErrPolisTanpaDatabase
	}
	h.Baris, h.Total, err = baca(s.viewPolis, (halaman-1)*ukuranHalamanFollowing, ukuranHalamanFollowing)
	return h, err
}

// CariOfferFac - GET /api/nbfacin/polis/offer (SEARCHON OFFERFAC, BrowseDataViewOfferF).
func (s *Service) CariOfferFac(ctx context.Context, pelaku inti.Pelaku, opsi, filter string, halaman int) (HalamanView[models.BarisOfferFac], error) {
	jenis, nilai, err := jenisCariOfferFac(opsi, filterPolis(filter))
	return cariView(ctx, s, pelaku, halaman, jenis, nilai, err, func(r repository.PembacaViewPolis, o, b int) ([]models.BarisOfferFac, int, error) {
		return r.CariOfferFac(ctx, jenis, nilai, o, b)
	})
}

// CariPolisTreaty - GET /api/nbfacin/polis/treaty (SEARCHON TREATYPOL, BrowsePolisTreaty).
func (s *Service) CariPolisTreaty(ctx context.Context, pelaku inti.Pelaku, opsi, filter string, halaman int) (HalamanView[models.BarisPolisTreaty], error) {
	jenis, nilai, err := jenisCariTreaty(opsi, filterPolis(filter))
	return cariView(ctx, s, pelaku, halaman, jenis, nilai, err, func(r repository.PembacaViewPolis, o, b int) ([]models.BarisPolisTreaty, int, error) {
		return r.CariPolisTreaty(ctx, jenis, nilai, o, b)
	})
}

// CariPolisLife - GET /api/nbfacin/polis/life (SEARCHON POLISLIFE, BrowsePolisLife).
func (s *Service) CariPolisLife(ctx context.Context, pelaku inti.Pelaku, filter string, halaman int) (HalamanView[models.BarisPolisLife], error) {
	jenis, nilai := jenisCariLife(filterPolis(filter))
	return cariView(ctx, s, pelaku, halaman, jenis, nilai, nil, func(r repository.PembacaViewPolis, o, b int) ([]models.BarisPolisLife, int, error) {
		return r.CariPolisLife(ctx, jenis, nilai, o, b)
	})
}

// CariOfferLife - GET /api/nbfacin/polis/life-offer (SEARCHON OFFERLIFE, BrowseOfferLife).
func (s *Service) CariOfferLife(ctx context.Context, pelaku inti.Pelaku, filter string, halaman int) (HalamanView[models.BarisOfferLife], error) {
	jenis, nilai := jenisCariOfferLife(filterPolis(filter))
	return cariView(ctx, s, pelaku, halaman, jenis, nilai, nil, func(r repository.PembacaViewPolis, o, b int) ([]models.BarisOfferLife, int, error) {
		return r.CariOfferLife(ctx, jenis, nilai, o, b)
	})
}

// BerkasPolisTreaty - GET /api/nbfacin/polis/treaty/berkas (tombol View tab Treaty, work owner 11-10-2026; pola Claim Non Prop
// OQ-CNP-13): nomor polis kosong 400, hanya login, tanpa berkas 404, tanpa pembaca 503.
func (s *Service) BerkasPolisTreaty(ctx context.Context, pelaku inti.Pelaku, noPolis string) (models.BerkasPolisTreaty, error) {
	if err := s.siapViewPolis(pelaku); err != nil {
		return models.BerkasPolisTreaty{}, err
	}
	noPolis = strings.TrimSpace(noPolis)
	if noPolis == "" {
		return models.BerkasPolisTreaty{}, errors.Join(ErrMasukanPolis, errors.New("nopolis wajib diisi"))
	}
	if s.viewPolis == nil {
		return models.BerkasPolisTreaty{}, ErrPolisTanpaDatabase
	}
	b, ada, err := s.viewPolis.BerkasPolisTreaty(ctx, noPolis)
	if err != nil {
		return models.BerkasPolisTreaty{}, err
	}
	if !ada {
		return models.BerkasPolisTreaty{}, ErrBerkasTreatyTidakAda
	}
	return b, nil
}

// siapViewPolis - identitas wajib (K71-3: semua pengguna login).
func (s *Service) siapViewPolis(pelaku inti.Pelaku) error { return inti.WajibIdentitas(pelaku) }
