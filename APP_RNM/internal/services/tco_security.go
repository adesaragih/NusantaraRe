package services

// Security reinsurer - tiket 06 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan grid security (`InputTreatyContractReinsType.xml`
// b14885, tampil bila `HASILD21 == 1`) yang dibuka tombol baris reinsurer
// `Security Reinsurer` (`ViewDetailTreatyReinsurerGrid1.xml` b5277 ->
// `SetSecurityReinsurer`: `THN_TREATY = .TreatyYear`, `REAS_ID = .ID`).
//
// ⛔ Tiga cacat warisan yang TIDAK dibawa (penyimpangan sadar 5):
//   - UPDATE mencocokkan `trim(REAS_SECURITY) = trim(CARI10)` dengan CARI10 =
//     nama BARU (`SaveSecurityReinsurer_Act.xml` b420/b953) - mengganti nama
//     security tidak mengenai baris lamanya. Di sini baris dicocokkan `ID`.
//   - DELETE berkunci nama - menghapus SEMUA baris bernama sama. Di sini satu ID.
//   - `Local.IsUpdate` dihitung (langkah 3, b760) lalu tidak dipakai: sisip vs
//     perbarui diputus `HASILD3` (b1225/b1409), sehingga security yang sama
//     dapat tersisip dua kali. Di sini dobel ditolak 409 - PENYIMPANGAN SADAR
//     dari Pega [keputusan work owner 29-09-2026] (OQ-TCO-17, ditutup); begitu pula
//     `%Share` wajib 0..100 (`models.UraiShareSecurityTCO`).
//
// ⛔ Reinsurer diturunkan server dari jalur (tahun -> kontrak -> kombinasi ->
// reinsurer); `THN_TREATY` = `TreatyYear` reinsurer (b5301), tidak dari klien.
//
// Dibaca sesudah: tco_reinsurer.go, models/tco_security.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/pkg/utils"
)

var (
	// ErrGudangSecurityBelumDisuntik - gudang security belum dipasang.
	ErrGudangSecurityBelumDisuntik = errors.New("services: gudang security belum disuntik")
	// ErrSecurityTidakAda - security bukan milik reinsurer itu (404).
	ErrSecurityTidakAda = repository.ErrSecurityTidakAda
	// ErrSecurityDobel - security yang sama sudah tercatat pada reinsurer itu (409).
	ErrSecurityDobel = errors.New("services: security dobel")
)

// GalatSecurityDobel menyebut baris mana yang sudah memegang security itu.
type GalatSecurityDobel struct{ IDLain, ReasSecurity string }

func (g GalatSecurityDobel) Error() string {
	return fmt.Sprintf("security %s sudah tercatat pada baris %s reinsurer ini", g.ReasSecurity, g.IDLain)
}

// Is membuat `errors.Is(err, ErrSecurityDobel)` benar.
func (GalatSecurityDobel) Is(target error) bool { return target == ErrSecurityDobel }

// GudangSecurityTCO membaca, menulis, dan menghapus security + jejaknya.
type GudangSecurityTCO interface {
	Daftar(ctx context.Context, reasID, thnTreaty string) ([]repository.SecurityTCO, error)
	Ambil(ctx context.Context, reasID, id string) (repository.SecurityTCO, error)
	CariDobel(ctx context.Context, tx *repository.Tx, reasID, reasSecurity, kecualiID string) (string, error)
	Sisip(ctx context.Context, tx *repository.Tx, s models.SecurityReinsurer) (string, error)
	Perbarui(ctx context.Context, tx *repository.Tx, s models.SecurityReinsurer) error
	Hapus(ctx context.Context, tx *repository.Tx, reasID, id string) error
	Jejak(ctx context.Context, tx *repository.Tx, akunID, barisID, aksi, keterangan string, waktu time.Time) error
}

// PembacaReinsurerTCO membaca satu reinsurer milik kombinasi - induk security.
type PembacaReinsurerTCO interface {
	Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.ReinsurerTreaty, error)
}

type securityBelumDisuntik struct{}

func (securityBelumDisuntik) Daftar(context.Context, string, string) ([]repository.SecurityTCO, error) {
	return nil, ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Ambil(context.Context, string, string) (repository.SecurityTCO, error) {
	return repository.SecurityTCO{}, ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) CariDobel(context.Context, *repository.Tx, string, string, string) (string, error) {
	return "", ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Sisip(context.Context, *repository.Tx, models.SecurityReinsurer) (string, error) {
	return "", ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Perbarui(context.Context, *repository.Tx, models.SecurityReinsurer) error {
	return ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Hapus(context.Context, *repository.Tx, string, string) error {
	return ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Jejak(context.Context, *repository.Tx, string, string, string, string, time.Time) error {
	return ErrGudangSecurityBelumDisuntik
}

type gudangSecurityOracle struct {
	*repository.MasterSecurityTCO
	db *repository.DB
}

func (g gudangSecurityOracle) Jejak(ctx context.Context, tx *repository.Tx, akunID, barisID, aksi,
	keterangan string, waktu time.Time) error {
	return g.db.SisipJejakTCO(ctx, tx, akunID, repository.TabelSecurityTCO, barisID, aksi, keterangan, waktu)
}

// GudangSecurityOracle menyusun gudang security di atas Oracle.
func GudangSecurityOracle(svc *Service) GudangSecurityTCO {
	return gudangSecurityOracle{MasterSecurityTCO: repository.NewMasterSecurityTCO(svc.db), db: svc.db}
}

// SecurityMasuk adalah badan simpan - medan form `Security Name` b19648
// (`REAS_SECURITY` dari pemilih) dan `%Share` b19888. Nama tampil DIABAIKAN:
// ia dari master.
type SecurityMasuk struct {
	ID           string `json:"id"`
	ReasSecurity string `json:"reasSecurity"`
	ClientName   string `json:"clientName"`
	PctShare     string `json:"pctShare"`
}

// SecurityTampil adalah satu security seperti dikirim ke layar.
type SecurityTampil struct {
	ID           string `json:"id"`
	ThnTreaty    string `json:"thnTreaty"`
	ReasID       string `json:"reasId"`
	ReasSecurity string `json:"reasSecurity"`
	ClientName   string `json:"clientName"`
	PctShare     string `json:"pctShare"`
	TopID        string `json:"topId"`
	TpTreaty     string `json:"tpTreaty"`
	UserID       string `json:"userId"`
}

// TampilSecurity menerjemahkan satu baris.
func TampilSecurity(s repository.SecurityTCO) SecurityTampil {
	return SecurityTampil{ID: s.ID, ThnTreaty: s.ThnTreaty, ReasID: s.ReasID, ReasSecurity: s.ReasSecurity,
		ClientName: s.ClientName, PctShare: utils.FormatDecimal(s.PctShare), TopID: s.TopID, TpTreaty: s.TpTreaty,
		UserID: s.UserID}
}

// DaftarSecurityTampil adalah grid security seorang reinsurer.
type DaftarSecurityTampil struct {
	Daftar    []SecurityTampil `json:"daftar"`
	Total     int              `json:"total"`
	Reinsurer ReinsurerTampil  `json:"reinsurer"`
}

// SecurityTCO melayani grid security.
type SecurityTCO struct {
	svc       *Service
	gudang    GudangSecurityTCO
	reinsurer PembacaReinsurerTCO
	kontrak   PemegangKontrakTCO
	tahun     PemeriksaTahunTCO
	master    PembacaReinsurerMasterTCO
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *repository.Tx) error) error
}

// SecurityTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) SecurityTCO() *SecurityTCO {
	return &SecurityTCO{svc: s, gudang: securityBelumDisuntik{}, reinsurer: reinsurerBelumDisuntik{},
		kontrak: pemegangKontrakBelumDisuntik{}, tahun: gudangTahunTreatyBelumDisuntik{},
		master: masterReinsurerBelumDisuntik{}, jam: time.Now, transaksi: s.DalamTransaksi}
}

func (l *SecurityTCO) salin() *SecurityTCO { s := *l; return &s }

// DenganGudang memasang gudang security.
func (l *SecurityTCO) DenganGudang(g GudangSecurityTCO) *SecurityTCO {
	s := l.salin()
	s.gudang = g
	return s
}

// DenganReinsurer memasang pembaca reinsurer induk.
func (l *SecurityTCO) DenganReinsurer(r PembacaReinsurerTCO) *SecurityTCO {
	s := l.salin()
	s.reinsurer = r
	return s
}

// DenganKontrak memasang pemegang kontrak pembuka kombinasi.
func (l *SecurityTCO) DenganKontrak(k PemegangKontrakTCO) *SecurityTCO {
	s := l.salin()
	s.kontrak = k
	return s
}

// DenganTahun memasang pemeriksa tahun treaty.
func (l *SecurityTCO) DenganTahun(t PemeriksaTahunTCO) *SecurityTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganMaster memasang pembaca master `AGENT` (pemilih b19711).
func (l *SecurityTCO) DenganMaster(m PembacaReinsurerMasterTCO) *SecurityTCO {
	s := l.salin()
	s.master = m
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *SecurityTCO) DenganJam(j func() time.Time) *SecurityTCO { s := l.salin(); s.jam = j; return s }

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *SecurityTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *repository.Tx) error) error) *SecurityTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

// induk menurunkan reinsurer induk dari tahun, kontrak, dan ID reinsurer.
func (l *SecurityTCO) induk(ctx context.Context, tahunID, kontrakID, reinsurerID string) (models.ReinsurerTreaty, error) {
	t, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return models.ReinsurerTreaty{}, err
	}
	k, err := l.kontrak.Ambil(ctx, tahunID, kontrakID)
	if err != nil {
		return models.ReinsurerTreaty{}, err
	}
	return l.reinsurer.Ambil(ctx, models.KombinasiDari(t, k), reinsurerID)
}

// Daftar membaca security seorang reinsurer - `SelectSecurityReinsurer`.
func (l *SecurityTCO) Daftar(ctx context.Context, pelaku Pelaku, tahunID, kontrakID, reinsurerID string) (
	DaftarSecurityTampil, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return DaftarSecurityTampil{}, err
	}
	r, err := l.induk(ctx, tahunID, kontrakID, reinsurerID)
	if err != nil {
		return DaftarSecurityTampil{}, err
	}
	baris, err := l.gudang.Daftar(ctx, r.ID, r.TreatyYear)
	if err != nil {
		return DaftarSecurityTampil{}, err
	}
	hasil := make([]SecurityTampil, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, TampilSecurity(b))
	}
	return DaftarSecurityTampil{Daftar: hasil, Total: len(hasil), Reinsurer: TampilReinsurer(r)}, nil
}

// Simpan menulis security baru atau memperbarui yang ada - `Save` b20246
// (`SaveSecurityReinsurer_Act`).
func (l *SecurityTCO) Simpan(ctx context.Context, pelaku Pelaku, tahunID, kontrakID, reinsurerID string,
	m SecurityMasuk) (SecurityTampil, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return SecurityTampil{}, err
	}
	r, err := l.induk(ctx, tahunID, kontrakID, reinsurerID)
	if err != nil {
		return SecurityTampil{}, err
	}
	s := models.SecurityReinsurer{ID: strings.TrimSpace(m.ID), ThnTreaty: r.TreatyYear, ReasID: r.ID,
		ReasSecurity: strings.TrimSpace(m.ReasSecurity)}
	if err := models.PeriksaSecurityTCO(s); err != nil {
		return SecurityTampil{}, err
	}
	if s.PctShare, err = models.UraiShareSecurityTCO(m.PctShare); err != nil {
		return SecurityTampil{}, err
	}
	master, err := l.master.Ambil(ctx, s.ReasSecurity)
	if err != nil {
		if errors.Is(err, ErrReinsurerDiLuarMaster) {
			return SecurityTampil{}, fmt.Errorf("%w: %s %q", err, models.LabelSecurityNameTCO, s.ReasSecurity)
		}
		return SecurityTampil{}, err
	}
	tampil := repository.SecurityTCO{SecurityReinsurer: s, ClientName: master.ClientName}
	waktu := l.jam()
	err = l.transaksi(ctx, func(tx *repository.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		ket := "security baru"
		if s.ID != "" {
			ket = "security diperbarui"
			lama, err := l.gudang.Ambil(ctx, r.ID, s.ID)
			if err != nil {
				return err
			}
			// Kolom yang UPDATE warisan tidak sentuh dipertahankan dari barisnya.
			s.TopID, s.TpTreaty, s.UserID = lama.TopID, lama.TpTreaty, lama.UserID
			if lama.ReasSecurity != s.ReasSecurity {
				ket = fmt.Sprintf("security diperbarui (dari %s)", lama.ReasSecurity)
			}
		}
		lain, err := l.gudang.CariDobel(ctx, tx, r.ID, s.ReasSecurity, s.ID)
		if err != nil {
			return err
		}
		if lain != "" {
			return GalatSecurityDobel{IDLain: lain, ReasSecurity: s.ReasSecurity}
		}
		if s.ID == "" {
			if s.ID, err = l.gudang.Sisip(ctx, tx, s); err != nil {
				return err
			}
		} else if err := l.gudang.Perbarui(ctx, tx, s); err != nil {
			return err
		}
		return l.gudang.Jejak(ctx, tx, pelaku.AkunID, s.ID, repository.AksiJejakSimpan,
			// Reinsurer disebut dengan kode agennya - identitas bisnis yang stabil.
			fmt.Sprintf("%s %s share %s reinsurer %s", ket, s.ReasSecurity, s.PctShare.Text('f'), r.ReinsurerID), waktu)
	})
	if err != nil {
		return SecurityTampil{}, err
	}
	tampil.SecurityReinsurer = s
	return TampilSecurity(tampil), nil
}

// Hapus membuang SATU security - `Delete` b17559 (`DeleteSecurityReinsurer`).
// Reinsurer induknya tidak disentuh.
func (l *SecurityTCO) Hapus(ctx context.Context, pelaku Pelaku, tahunID, kontrakID, reinsurerID, id string) (string, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	r, err := l.induk(ctx, tahunID, kontrakID, reinsurerID)
	if err != nil {
		return "", err
	}
	err = l.transaksi(ctx, func(tx *repository.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		if err := l.gudang.Hapus(ctx, tx, r.ID, id); err != nil {
			return err
		}
		return l.gudang.Jejak(ctx, tx, pelaku.AkunID, id, repository.AksiJejakHapus,
			"security dihapus reinsurer "+r.ID, l.jam())
	})
	if err != nil {
		return "", err
	}
	// `DeleteSecurityReinsurer` tidak menampilkan pesan - `[tidak ada di korpus]`.
	return "Security dengan ID " + id + " dihapus", nil
}
