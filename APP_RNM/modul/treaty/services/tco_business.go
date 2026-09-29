package services

// Business pada kombinasi kontrak - tiket 07 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan panel `ViewDetailTreatyBusinessGrid`
// (disertakan `InputTreatyContractReinsType.xml` b14064), dibuka tombol baris
// kontrak `Business List` b10842 (`BrowseTreatyBusinessList_Act` +
// `SetTreatyBusinessList_Act`).
//
// ⛔ Tiga perbaikan sadar atas perilaku warisan, masing-masing dari AC:
//   - baris nonaktif TETAP tampil (AC 22) - grid Pega menyaring `IsActive = 1`;
//   - pembaruan menulis SELURUH medan (AC 23) - prosedur hanya lima kolom;
//   - hapus menyentuh SATU tabel (AC 63/64) - rule Pega memegang dua.
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/utils"
	"nusantarare/modul/treaty/models"
	"nusantarare/modul/treaty/repository"
)

var (
	// ErrGudangBusinessBelumDisuntik - gudang bisnis belum dipasang.
	ErrGudangBusinessBelumDisuntik = errors.New("services: gudang bisnis belum disuntik")
	// ErrBusinessTidakAda - baris bisnis bukan milik kombinasi kontrak itu (404).
	ErrBusinessTidakAda = repository.ErrBusinessTidakAda
	// ErrBusinessDiLuarMaster - kode bisnis tidak ada di master BUSINESS (422).
	ErrBusinessDiLuarMaster = repository.ErrBusinessMasterTidakAda
	// ErrBusinessDobel - kode bisnis sudah ada pada kombinasi (409).
	ErrBusinessDobel = errors.New("services: bisnis dobel")
)

// GalatBusinessDobel menyebut baris mana yang sudah memegang kode itu.
//
// Pesan VERBATIM `SaveTreatyBusinessDetail_Act.xml` langkah 6 (`ERRMSG4`),
// yang layar tampilkan saat prosedur menjawab galat (`ViewDetailTreatyBusinessGrid`
// b10095 + b10215).
type GalatBusinessDobel struct{ IDLain, BizCode string }

func (g GalatBusinessDobel) Error() string {
	return fmt.Sprintf("%s: kode bisnis %s sudah ada pada baris %s kombinasi ini", PesanKontrakDobelTCO, g.BizCode, g.IDLain)
}

// Is membuat `errors.Is(err, ErrBusinessDobel)` benar.
func (GalatBusinessDobel) Is(target error) bool { return target == ErrBusinessDobel }

// GudangBusinessTCO membaca, menulis, dan menghapus baris bisnis.
type GudangBusinessTCO interface {
	Daftar(ctx context.Context, k models.KombinasiTCO) ([]models.BusinessTreaty, error)
	Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.BusinessTreaty, error)
	CariDobel(ctx context.Context, tx *db.Tx, k models.KombinasiTCO, bizCode, kecualiID string) (string, error)
	Sisip(ctx context.Context, tx *db.Tx, b models.BusinessTreaty) (string, error)
	Perbarui(ctx context.Context, tx *db.Tx, b models.BusinessTreaty) error
	Hapus(ctx context.Context, tx *db.Tx, k models.KombinasiTCO, id string) error
}

// PembacaBusinessMasterTCO membaca master `BUSINESS` (dibaca saja).
type PembacaBusinessMasterTCO interface {
	Daftar(ctx context.Context) ([]repository.BusinessMasterTCO, error)
	Ambil(ctx context.Context, id string) (repository.BusinessMasterTCO, error)
}

type businessBelumDisuntik struct{}

func (businessBelumDisuntik) Daftar(context.Context, models.KombinasiTCO) ([]models.BusinessTreaty, error) {
	return nil, ErrGudangBusinessBelumDisuntik
}
func (businessBelumDisuntik) Ambil(context.Context, models.KombinasiTCO, string) (models.BusinessTreaty, error) {
	return models.BusinessTreaty{}, ErrGudangBusinessBelumDisuntik
}
func (businessBelumDisuntik) CariDobel(context.Context, *db.Tx, models.KombinasiTCO, string, string) (string, error) {
	return "", ErrGudangBusinessBelumDisuntik
}
func (businessBelumDisuntik) Sisip(context.Context, *db.Tx, models.BusinessTreaty) (string, error) {
	return "", ErrGudangBusinessBelumDisuntik
}
func (businessBelumDisuntik) Perbarui(context.Context, *db.Tx, models.BusinessTreaty) error {
	return ErrGudangBusinessBelumDisuntik
}
func (businessBelumDisuntik) Hapus(context.Context, *db.Tx, models.KombinasiTCO, string) error {
	return ErrGudangBusinessBelumDisuntik
}

type masterBusinessBelumDisuntik struct{}

func (masterBusinessBelumDisuntik) Daftar(context.Context) ([]repository.BusinessMasterTCO, error) {
	return nil, ErrGudangBusinessBelumDisuntik
}
func (masterBusinessBelumDisuntik) Ambil(context.Context, string) (repository.BusinessMasterTCO, error) {
	return repository.BusinessMasterTCO{}, ErrGudangBusinessBelumDisuntik
}

type gudangBusinessOracle struct {
	*repository.MasterBusinessKombinasiTCO
	db *db.DB
}

// GudangBusinessOracle menyusun gudang bisnis di atas Oracle.
func GudangBusinessOracle(svc *Service) GudangBusinessTCO {
	return gudangBusinessOracle{MasterBusinessKombinasiTCO: repository.NewMasterBusinessKombinasiTCO(svc.DB()), db: svc.DB()}
}

// MasterBusinessOracle menyusun pembaca master `BUSINESS`.
func MasterBusinessOracle(svc *Service) PembacaBusinessMasterTCO {
	return repository.NewMasterBusiness(svc.DB())
}

// BusinessMasuk adalah badan simpan - medan yang TAMPIL di form (`Business Name`
// b6241, `Active` b6499). `BizName` DIABAIKAN: namanya dari master.
type BusinessMasuk struct {
	ID       string `json:"id"`
	BizCode  string `json:"bizCode"`
	BizName  string `json:"bizName"`
	IsActive string `json:"isActive"`
}

// BusinessTampil adalah satu baris bisnis seperti dikirim ke layar.
type BusinessTampil struct {
	ID              string `json:"id"`
	IsActive        string `json:"isActive"`
	Aktif           bool   `json:"aktif"`
	TreatyYear      string `json:"treatyYear"`
	TreatyYearID    string `json:"treatyYearId"`
	TreatyGroupID   string `json:"treatyGroupId"`
	TreatyGroupName string `json:"treatyGroupName"`
	ReinsTypeID     string `json:"reinsTypeId"`
	ReinsTypeName   string `json:"reinsTypeName"`
	BizCode         string `json:"bizCode"`
	BizName         string `json:"bizName"`
	UserID          string `json:"userId"`
	TglUpdate       string `json:"tglUpdate"`
}

// TampilBusiness menerjemahkan satu baris.
func TampilBusiness(b models.BusinessTreaty) BusinessTampil {
	return BusinessTampil{ID: b.ID, IsActive: b.IsActive, Aktif: models.BusinessAktifTCO(b), TreatyYear: b.TreatyYear,
		TreatyYearID: b.TreatyYearID, TreatyGroupID: b.TreatyGroupID, TreatyGroupName: b.TreatyGroupName,
		ReinsTypeID: b.ReinsTypeID, ReinsTypeName: b.ReinsTypeName, BizCode: b.BizCode, BizName: b.BizName,
		UserID: b.UserID, TglUpdate: utils.FormatTanggalWaktu(b.TglUpdate)}
}

// DaftarBusinessTampil adalah grid bisnis satu kombinasi.
type DaftarBusinessTampil struct {
	Daftar    []BusinessTampil `json:"daftar"`
	Total     int              `json:"total"`
	Kombinasi KombinasiTampil  `json:"kombinasi"`
}

// BusinessMasterTampil adalah satu pilihan bisnis.
type BusinessMasterTampil struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

// BusinessTCO melayani panel bisnis.
type BusinessTCO struct {
	svc       *Service
	gudang    GudangBusinessTCO
	kontrak   PemegangKontrakTCO
	tahun     PemeriksaTahunTCO
	master    PembacaBusinessMasterTCO
	catat     func(string) // log aplikasi: pelaku yang tidak ditulis ke kolom warisan (OQ-TCO-25)
	transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error
}

// BusinessTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) BusinessTCO() *BusinessTCO {
	return &BusinessTCO{svc: s, gudang: businessBelumDisuntik{}, kontrak: pemegangKontrakBelumDisuntik{},
		tahun: gudangTahunTreatyBelumDisuntik{}, master: masterBusinessBelumDisuntik{},
		catat: func(p string) { log.Print(p) }, transaksi: s.DalamTransaksi}
}

func (l *BusinessTCO) salin() *BusinessTCO { s := *l; return &s }

// DenganCatat mengganti tujuan log aplikasi - dipakai uji.
func (l *BusinessTCO) DenganCatat(c func(string)) *BusinessTCO { s := l.salin(); s.catat = c; return s }

// DenganGudang memasang gudang bisnis.
func (l *BusinessTCO) DenganGudang(g GudangBusinessTCO) *BusinessTCO {
	s := l.salin()
	s.gudang = g
	return s
}

// DenganKontrak memasang pemegang kontrak pembuka kombinasi.
func (l *BusinessTCO) DenganKontrak(k PemegangKontrakTCO) *BusinessTCO {
	s := l.salin()
	s.kontrak = k
	return s
}

// DenganTahun memasang pemeriksa tahun treaty.
func (l *BusinessTCO) DenganTahun(t PemeriksaTahunTCO) *BusinessTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganMaster memasang pembaca master bisnis.
func (l *BusinessTCO) DenganMaster(m PembacaBusinessMasterTCO) *BusinessTCO {
	s := l.salin()
	s.master = m
	return s
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *BusinessTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *BusinessTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

func (l *BusinessTCO) kombinasi(ctx context.Context, tahunID, kontrakID string) (models.KombinasiTCO, error) {
	t, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return models.KombinasiTCO{}, err
	}
	k, err := l.kontrak.Ambil(ctx, tahunID, kontrakID)
	if err != nil {
		return models.KombinasiTCO{}, err
	}
	return models.KombinasiDari(t, k), nil
}

// Master membaca pilihan `Business Name`.
func (l *BusinessTCO) Master(ctx context.Context, pelaku inti.Pelaku) ([]BusinessMasterTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	d, err := l.master.Daftar(ctx)
	if err != nil {
		return nil, err
	}
	hasil := make([]BusinessMasterTampil, 0, len(d))
	for _, b := range d {
		hasil = append(hasil, BusinessMasterTampil{ID: b.ID, Note: b.Note})
	}
	return hasil, nil
}

// Daftar membaca seluruh baris bisnis kombinasi, aktif maupun nonaktif.
func (l *BusinessTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID string) (DaftarBusinessTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return DaftarBusinessTampil{}, err
	}
	k, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return DaftarBusinessTampil{}, err
	}
	baris, err := l.gudang.Daftar(ctx, k)
	if err != nil {
		return DaftarBusinessTampil{}, err
	}
	hasil := make([]BusinessTampil, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, TampilBusiness(b))
	}
	return DaftarBusinessTampil{Daftar: hasil, Total: len(hasil), Kombinasi: TampilKombinasi(k)}, nil
}

// Simpan menulis baris bisnis baru atau memperbarui SELURUH medannya - `Save`
// b6966 (`SaveTreatyBusinessDetail_Act`).
func (l *BusinessTCO) Simpan(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID string, m BusinessMasuk) (BusinessTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return BusinessTampil{}, err
	}
	k, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return BusinessTampil{}, err
	}
	// `SaveTreatyBusinessDetail_Act` langkah 1: TreatyYearID = IDTreatyYear kontrak.
	// OQ-TCO-25 (lanjutan 4): USERID/TGLUPDATE tidak diisi, seperti Pega - tidak
	// ada langkah yang mengisi `InputTreatyBusiness.UserID/TglUpdate`, dan
	// `SaveMasterTreatyBusiness_SQL` b85 meneruskannya kosong; data DEV 2/4.621.
	b := models.BusinessTreaty{ID: strings.TrimSpace(m.ID), IsActive: strings.TrimSpace(m.IsActive),
		TreatyYear: k.TreatyYear, TreatyYearID: tahunID, TreatyGroupID: k.TreatyGroupID,
		TreatyGroupName: k.TreatyGroupName, ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName,
		BizCode: strings.TrimSpace(m.BizCode)}
	if err := models.PeriksaBusinessTCO(b); err != nil {
		return BusinessTampil{}, err
	}
	master, err := l.master.Ambil(ctx, b.BizCode)
	if err != nil {
		if errors.Is(err, ErrBusinessDiLuarMaster) {
			return BusinessTampil{}, fmt.Errorf("%w: BizCode %q", err, b.BizCode)
		}
		return BusinessTampil{}, err
	}
	b.BizName = master.Note
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		lain, err := l.gudang.CariDobel(ctx, tx, k, b.BizCode, b.ID)
		if err != nil {
			return err
		}
		if lain != "" {
			return GalatBusinessDobel{IDLain: lain, BizCode: b.BizCode}
		}
		if b.ID == "" {
			if b.ID, err = l.gudang.Sisip(ctx, tx, b); err != nil {
				return err
			}
			return nil
		}
		// tco4: UPDATE procedure hanya lima kolom - kolom lain baris lama yang
		// tersimpan, jadi itu pula yang dijawab (temuan /code-review).
		lama, err := l.gudang.Ambil(ctx, k, b.ID)
		if err != nil {
			return err
		}
		b.TreatyYearID, b.TreatyGroupName, b.ReinsTypeName = lama.TreatyYearID, lama.TreatyGroupName, lama.ReinsTypeName
		return l.gudang.Perbarui(ctx, tx, b)
	})
	if err != nil {
		return BusinessTampil{}, err
	}
	l.catat(fmt.Sprintf("treaty contract out: business %s disimpan oleh akun %s", b.ID, pelaku.AkunID))
	return TampilBusiness(b), nil
}

// Hapus membuang SATU baris bisnis - `Delete` b4826 (`DeleteRowBusiness`).
//
// Mengembalikan pesan VERBATIM `DeleteRowBusiness.xml` langkah 3:
// `"Data Dengan ID" + " " + ID + " " + "Berhasil di Hapus"`.
func (l *BusinessTCO) Hapus(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID, id string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	k, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return "", err
	}
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.gudang.Hapus(ctx, tx, k, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return "Data Dengan ID " + id + " Berhasil di Hapus", nil
}
