package services

// Menu View Policy (Facultative) - tiket 71 (kontrak sesi 9d; K71-1..3, A71-1..5 di docs/issues/71-view-policy.md). Bukti
// `[terverifikasi]` `D:\migrasi\RNM\DDL\GetListViewPolicy_act.xml`:
//   - langkah 6: `.FilterTermForPolicy` dibuang tanda kutip ganda / tunggal (`@pxReplaceAllViaRegex(…,"[\"']","")`); langkah 7:
//     filter kosong -> keluar;
//   - blok 8 "POLIS FACIN" = GetListViewFollowing_act blok 4 (tiket 70) KECUALI 8.7 `INSUREDNAME LIKE 'x%'`, 8.8 `SOB LIKE 'x%'`,
//     8.9 `UPPER(QQNAME) LIKE 'x%'` (awalan; A71-2: filter ikut di-UPPER). 8.3 opsi 3 ARASAPAS ditunda (K70-3).
// Akses (K71-3): semua pengguna login - tanpa penjaga case.

import (
	"context"
	"errors"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrPolisTanpaDatabase - pembaca View Policy belum dipasang (503).
	ErrPolisTanpaDatabase = errors.New("services: View Policy tidak tersedia tanpa basis data")
	// ErrMasukanPolis - nomor polis / versi kosong (400).
	ErrMasukanPolis = errors.New("isian View Policy tidak sah")
	// ErrPolisTidakAda - polis / versi tidak ditemukan (404).
	ErrPolisTidakAda = errors.New("Polis tidak ditemukan.")
)

// DenganPolis - memasang pembaca menu View Policy (tiket 71).
func (s *Service) DenganPolis(p repository.PembacaPolis) *Service { s.polis = p; return s }

// filterPolis - langkah 6 (kutip dibuang) + pangkas spasi.
func filterPolis(f string) string {
	return strings.TrimSpace(strings.NewReplacer(`"`, "", "'", "").Replace(f))
}

// jenisCariPolis - jenisCariFollowing lalu beda blok 8: opsi 5 / 6 awalan, opsi 7 UPPER awalan.
func jenisCariPolis(opsi, filter string) (string, string, error) {
	jenis, nilai, err := jenisCariFollowing(opsi, filter)
	if err != nil || jenis == "" {
		return jenis, nilai, err
	}
	switch opsi {
	case "5", "6":
		nilai = filter + "%"
	case "7":
		jenis, nilai = repository.CariPolisQQ, filter+"%"
	}
	return jenis, nilai, nil
}

// CariPolis - GET /api/nbfacin/polis (GetListViewPolicy_act SEARCHON "POLIFAC"); baris = BarisFollowing (tiket 70).
func (s *Service) CariPolis(ctx context.Context, pelaku inti.Pelaku, opsi, filter string, halaman int) (HalamanFollowing, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanFollowing{}, err
	}
	if halaman < 1 {
		return HalamanFollowing{}, errors.Join(ErrMasukanFollowing, errors.New("halaman wajib >= 1"))
	}
	jenis, nilai, err := jenisCariPolis(opsi, filterPolis(filter))
	if err != nil {
		return HalamanFollowing{}, err
	}
	h := HalamanFollowing{Baris: []models.BarisFollowing{}, Halaman: halaman, Ukuran: ukuranHalamanFollowing}
	if jenis == "" {
		return h, nil
	}
	if s.following == nil {
		return HalamanFollowing{}, ErrPolisTanpaDatabase
	}
	h.Baris, h.Total, err = s.following.CariFollowing(ctx, jenis, nilai, (halaman-1)*ukuranHalamanFollowing, ukuranHalamanFollowing)
	return h, err
}

// siapPolis - identitas, pembaca, nomor polis.
func (s *Service) siapPolis(pelaku inti.Pelaku, noPolis string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(noPolis) == "" {
		return errors.Join(ErrMasukanPolis, errors.New("nopolis wajib diisi"))
	}
	if s.polis == nil {
		return ErrPolisTanpaDatabase
	}
	return nil
}

// galatPolis - repository.ErrPolisTidakAda -> ErrPolisTidakAda (404).
func galatPolis(err error) error {
	if errors.Is(err, repository.ErrPolisTidakAda) {
		return ErrPolisTidakAda
	}
	return err
}

// VersiPolis - GET /api/nbfacin/polis/versi (K71-2: dropdown hanya bila > 1 - keputusan layar).
func (s *Service) VersiPolis(ctx context.Context, pelaku inti.Pelaku, noPolis string) ([]models.VersiPolis, error) {
	if err := s.siapPolis(pelaku, noPolis); err != nil {
		return nil, err
	}
	v, err := s.polis.VersiPolis(ctx, strings.TrimSpace(noPolis))
	return v, galatPolis(err)
}

// RingkasanPolis - GET /api/nbfacin/polis/ringkasan (A71-3, polis lama). Policy Type = kode (kodePolicyType); Rate kolom VARCHAR2
// dinormalkan ke desimal bertitik bila terurai (TO_NUMBER implisit Pega dapat menulis ".15"), selain itu apa adanya.
func (s *Service) RingkasanPolis(ctx context.Context, pelaku inti.Pelaku, noPolis, prodKe string) (models.RingkasanPolis, error) {
	if err := s.siapPolis(pelaku, noPolis); err != nil {
		return models.RingkasanPolis{}, err
	}
	if strings.TrimSpace(prodKe) == "" {
		return models.RingkasanPolis{}, errors.Join(ErrMasukanPolis, errors.New("prodke wajib diisi"))
	}
	r, err := s.polis.RingkasanPolis(ctx, strings.TrimSpace(noPolis), strings.TrimSpace(prodKe))
	if err != nil {
		return models.RingkasanPolis{}, galatPolis(err)
	}
	r.Kepala.PolicyType = kodePolicyType(r.Kepala.PolicyType)
	for i := range r.Baris {
		r.Baris[i].Rate = rateTeks(r.Baris[i].Rate)
	}
	return r, nil
}

// rateTeks - teks RATE FACINPRODUCTION -> desimal bertitik bila terurai.
func rateTeks(s string) string {
	t := strings.TrimSpace(s)
	if t == "" {
		return ""
	}
	d, _, err := apd.NewFromString(t)
	if err != nil {
		return s
	}
	return utils.FormatDecimal(d)
}
