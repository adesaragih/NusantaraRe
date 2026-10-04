// Package services memuat layanan modul Master Data: daftar, tambah, ubah, dan aktif / nonaktif baris master
// (rencana `docs/issues/01-rencana-master-data.md`, keputusan MD-1 ... MD-9 di `docs/issues/02-api-master-data.md`).
// Tidak ada hapus - hanya nonaktif (M-3).
package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterdata/backend/models"
	"nusantarare/modul/masterdata/backend/repository"
)

var (
	// ErrMasterTidakDikenal - `{tabel}` bukan master yang dikenal. 404.
	ErrMasterTidakDikenal = errors.New("services: master tidak dikenal")
	// ErrBarisTidakAda - ID tidak ada. 404.
	ErrBarisTidakAda = errors.New("services: baris master tidak ada")
	// ErrMasukanMaster - isian tidak sah. 400.
	ErrMasukanMaster = errors.New("services: isian master tidak sah")
	// ErrSudahAda - ID (atau akumulasi Note + zip) sudah ada. 409.
	ErrSudahAda = errors.New("services: data master sudah ada")
	// ErrMasterTanpaDatabase - basis data tidak dikonfigurasi. 503.
	ErrMasterTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, master tidak terbaca")
)

const (
	// UkuranHalaman - baris per halaman daftar (MD-6).
	UkuranHalaman = 50
	// lebarKata - batas `q` (pola A73).
	lebarKata = 255
	// lebarNegara - kode negara ID akumulasi (contoh korpus `INA-61151-041589`: tiga huruf; batas 10 = NATION.ID).
	lebarNegara = 10
	// KunciNegara - medan masukan tambahan akumulasi baru (bukan kolom; MD-4).
	KunciNegara = "negara"
	// lebarPelaku - CREATE_OP / UPDATE_OP VARCHAR2(64) = M_LOGIN_GO.LOGIN_ID (migrasi 882).
	lebarPelaku = 64
)

// pelakuJejak - akun pelaku untuk jejak ubah (MD-7); wajib ada dan muat kolomnya.
func pelakuJejak(p inti.Pelaku) (string, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return "", err
	}
	akun := strings.TrimSpace(p.AkunID)
	if len(akun) > lebarPelaku {
		return "", fmt.Errorf("%w: akun pelaku paling banyak %d byte", ErrMasukanMaster, lebarPelaku)
	}
	return akun, nil
}

// Transaksi - pembuka transaksi (inti.Dasar.DalamTransaksi).
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Service - layanan Master Data.
type Service struct {
	simpan    repository.Penyimpan
	transaksi Transaksi
}

// Baru merakit layanan; penyimpan / transaksi nil = tanpa basis data (503).
func Baru(p repository.Penyimpan, tx Transaksi) *Service { return &Service{simpan: p, transaksi: tx} }

// DariDasar merakit layanan dari dasar aplikasi.
func DariDasar(d *inti.Dasar) *Service {
	if !d.PunyaDatabase() {
		return Baru(nil, nil)
	}
	return Baru(repository.NewMasterOracle(d.DB()), d.DalamTransaksi)
}

func (s *Service) master(kunci string) (models.TabelMaster, error) {
	t, ada := models.CariMaster(kunci)
	if !ada {
		return t, fmt.Errorf("%w: %q", ErrMasterTidakDikenal, kunci)
	}
	if s.simpan == nil || s.transaksi == nil {
		return t, ErrMasterTanpaDatabase
	}
	return t, nil
}

// Daftar - satu halaman master `kunci`. status: "", "aktif", "nonaktif"; nomor >= 1.
func (s *Service) Daftar(ctx context.Context, kunci, kata, status string, nomor int) (models.Halaman, error) {
	if _, ada := models.CariMaster(kunci); !ada {
		return models.Halaman{}, fmt.Errorf("%w: %q", ErrMasterTidakDikenal, kunci)
	}
	kata = strings.TrimSpace(kata)
	switch {
	case utf8.RuneCountInString(kata) > lebarKata:
		return models.Halaman{}, fmt.Errorf("%w: q paling banyak %d karakter", ErrMasukanMaster, lebarKata)
	case status != "" && status != "aktif" && status != "nonaktif":
		return models.Halaman{}, fmt.Errorf("%w: status harus aktif atau nonaktif", ErrMasukanMaster)
	case nomor < 1:
		return models.Halaman{}, fmt.Errorf("%w: halaman mulai dari 1", ErrMasukanMaster)
	}
	t, err := s.master(kunci)
	if err != nil {
		return models.Halaman{}, err
	}
	return s.simpan.Daftar(ctx, t, kata, status, nomor, UkuranHalaman)
}

// periksa - pangkas, buang medan turunan, tolak kunci asing; wajib dan lebar. `baru` = tambah (ID wajib bila tidak
// otomatis).
func periksa(t models.TabelMaster, masukan map[string]string, baru bool) (models.Baris, error) {
	sah := map[string]models.Kolom{}
	// Kolom jejak ubah (MD-7) dikenal tetapi turunan: dikirim balik layar pun diabaikan, tidak ditolak.
	for _, k := range t.SeluruhKolom() {
		sah[k.JSON] = k
	}
	var masalah []string
	b := models.Baris{}
	for kunci, nilai := range masukan {
		k, ada := sah[kunci]
		if !ada && !(kunci == KunciNegara && t.IDOtomatis && baru) {
			masalah = append(masalah, kunci+" bukan medan master "+t.Kunci)
			continue
		}
		if k.Turunan {
			continue // diisi dari rujukan (M-5)
		}
		b[kunci] = strings.TrimSpace(nilai)
	}
	for _, k := range t.Kolom {
		v := b[k.JSON]
		if k.HurufBesar {
			v = strings.ToUpper(v)
			b[k.JSON] = v
		}
		lewati := k.Nama == "ID" && (t.IDOtomatis || !baru)
		if k.Wajib && v == "" && !lewati {
			masalah = append(masalah, k.JSON+" wajib diisi")
		}
		if len(v) > k.Lebar {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d byte", k.JSON, k.Lebar))
		}
	}
	if t.IDOtomatis && baru {
		n := b[KunciNegara]
		if n == "" {
			masalah = append(masalah, KunciNegara+" wajib diisi (awalan ID akumulasi)")
		} else if len(n) > lebarNegara {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d byte", KunciNegara, lebarNegara))
		}
	}
	if len(masalah) > 0 {
		sort.Strings(masalah)
		return nil, fmt.Errorf("%w: %s", ErrMasukanMaster, strings.Join(masalah, "; "))
	}
	return b, nil
}

// isiRujukan - setiap rujukan yang diisi wajib ada; kolom turunannya diisi (M-5). Rujukan kosong -> turunan kosong.
func (s *Service) isiRujukan(ctx context.Context, t models.TabelMaster, b models.Baris) error {
	jsonDari := map[string]string{}
	for _, k := range t.Kolom {
		jsonDari[k.Nama] = k.JSON
	}
	var masalah []string
	for _, r := range t.Rujukan {
		nilai := b[jsonDari[r.Sumber]]
		turunan := ""
		if nilai != "" {
			isi, ada, err := s.simpan.NilaiRujukan(ctx, r, nilai)
			if err != nil {
				return err
			}
			if !ada {
				masalah = append(masalah, fmt.Sprintf("%s %q tidak ada di %s", jsonDari[r.Sumber], nilai, r.Tabel))
				continue
			}
			turunan = isi
		}
		if r.Turunan != "" {
			b[jsonDari[r.Turunan]] = turunan
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanMaster, strings.Join(masalah, "; "))
	}
	return nil
}

// Tambah - baris baru berstatus aktif; mengembalikan ID-nya.
func (s *Service) Tambah(ctx context.Context, pelaku inti.Pelaku, kunci string, masukan map[string]string) (string, error) {
	akun, err := pelakuJejak(pelaku)
	if err != nil {
		return "", err
	}
	t, ada := models.CariMaster(kunci)
	if !ada {
		return "", fmt.Errorf("%w: %q", ErrMasterTidakDikenal, kunci)
	}
	b, err := periksa(t, masukan, true)
	if err != nil {
		return "", err
	}
	if t, err = s.master(kunci); err != nil {
		return "", err
	}
	if err := s.isiRujukan(ctx, t, b); err != nil {
		return "", err
	}
	if t.IDOtomatis {
		ada, err := s.simpan.AdaCatatanAkumulasi(ctx, b["note"], b["zipCode"])
		if err != nil {
			return "", err
		}
		if ada {
			// Pesan prosedur RDBMASTERACCUMULATION: "Error master item Akumulasi Sudah Ada".
			return "", fmt.Errorf("%w: master item Akumulasi sudah ada (note dan zipCode sama)", ErrSudahAda)
		}
	} else {
		ada, err := s.simpan.Ada(ctx, t, b["id"])
		if err != nil {
			return "", err
		}
		if ada {
			return "", fmt.Errorf("%w: id %q", ErrSudahAda, b["id"])
		}
	}
	negara := b[KunciNegara]
	delete(b, KunciNegara)
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		if t.IDOtomatis {
			id, err := s.simpan.IDAkumulasi(ctx, tx, negara, b["zipCode"])
			if err != nil {
				return err
			}
			b["id"] = id
		}
		return s.simpan.Sisip(ctx, tx, t, b, akun)
	})
	if err != nil {
		return "", err
	}
	return b["id"], nil
}

// Ubah - kolom data selain ID (ID dari rute; `id` di badan diabaikan).
func (s *Service) Ubah(ctx context.Context, pelaku inti.Pelaku, kunci, id string, masukan map[string]string) error {
	akun, err := pelakuJejak(pelaku)
	if err != nil {
		return err
	}
	t, ada := models.CariMaster(kunci)
	if !ada {
		return fmt.Errorf("%w: %q", ErrMasterTidakDikenal, kunci)
	}
	b, err := periksa(t, masukan, false)
	if err != nil {
		return err
	}
	if t, err = s.master(kunci); err != nil {
		return err
	}
	if err := s.isiRujukan(ctx, t, b); err != nil {
		return err
	}
	b["id"] = strings.TrimSpace(id)
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.simpan.Ubah(ctx, tx, t, b, akun) })
	if errors.Is(err, repository.ErrBarisTidakAda) {
		return fmt.Errorf("%w: %s %q", ErrBarisTidakAda, kunci, id)
	}
	return err
}

// UbahStatus - aktif / nonaktif (tanpa hapus).
func (s *Service) UbahStatus(ctx context.Context, pelaku inti.Pelaku, kunci, id string, aktif bool) error {
	akun, err := pelakuJejak(pelaku)
	if err != nil {
		return err
	}
	t, err := s.master(kunci)
	if err != nil {
		return err
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.simpan.UbahStatus(ctx, tx, t, strings.TrimSpace(id), aktif, akun) })
	if errors.Is(err, repository.ErrBarisTidakAda) {
		return fmt.Errorf("%w: %s %q", ErrBarisTidakAda, kunci, id)
	}
	return err
}
