package services

// Security reinsurer - tiket 06 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan grid security (`InputTreatyContractReinsType.xml`
// b14885, tampil bila `HASILD21 == 1`) yang dibuka tombol baris reinsurer
// `Security Reinsurer` (`ViewDetailTreatyReinsurerGrid1.xml` b5277 ->
// `SetSecurityReinsurer`: `THN_TREATY = .TreatyYear`, `REAS_ID = .ID`).
//
// ⛔ tco4: tabel warisan `MTREATYSECURITY` tanpa identitas - SQL-nya PERSIS
// `UpdateMTreatySecurity`/`DeleteSecurityReinsurer` (kunci `REAS_ID` +
// `TRIM(REAS_SECURITY)`, semua baris senama). Cacat warisan yang TIDAK dibawa:
//   - UPDATE Pega mengisi kuncinya dengan CARI10 = nama BARU
//     (`SaveSecurityReinsurer_Act.xml` b420/b953) - mengganti nama security
//     tidak mengenai baris lamanya. Di sini kuncinya nama LAMA dari rute (`ID`).
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

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/utils"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
)

var (
	// ErrGudangSecurityBelumDisuntik - gudang security belum dipasang.
	ErrGudangSecurityBelumDisuntik = errors.New("services: security store is not injected")
	// ErrSecurityTidakAda - security bukan milik reinsurer itu (404).
	ErrSecurityTidakAda = repository.ErrSecurityTidakAda
	// ErrSecurityDobel - security yang sama sudah tercatat pada reinsurer itu (409).
	ErrSecurityDobel = errors.New("services: duplicate security")
)

// GalatSecurityDobel menyebut baris mana yang sudah memegang security itu.
type GalatSecurityDobel struct{ IDLain, ReasSecurity string }

func (g GalatSecurityDobel) Error() string {
	return fmt.Sprintf("security %s is already recorded in row %s of this reinsurer", g.ReasSecurity, g.IDLain)
}

// Is membuat `errors.Is(err, ErrSecurityDobel)` benar.
func (GalatSecurityDobel) Is(target error) bool { return target == ErrSecurityDobel }

// GudangSecurityTCO membaca, menulis, dan menghapus security.
type GudangSecurityTCO interface {
	Daftar(ctx context.Context, reasID, thnTreaty string) ([]repository.SecurityTCO, error)
	Ambil(ctx context.Context, reasID, id string) (repository.SecurityTCO, error)
	CariDobel(ctx context.Context, tx *db.Tx, reasID, reasSecurity, kecualiID string) (string, error)
	Sisip(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) (string, error)
	Perbarui(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) error
	Hapus(ctx context.Context, tx *db.Tx, reasID, id string) error
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
func (securityBelumDisuntik) CariDobel(context.Context, *db.Tx, string, string, string) (string, error) {
	return "", ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Sisip(context.Context, *db.Tx, models.SecurityReinsurer) (string, error) {
	return "", ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Perbarui(context.Context, *db.Tx, models.SecurityReinsurer) error {
	return ErrGudangSecurityBelumDisuntik
}
func (securityBelumDisuntik) Hapus(context.Context, *db.Tx, string, string) error {
	return ErrGudangSecurityBelumDisuntik
}

type gudangSecurityOracle struct {
	*repository.MasterSecurityTCO
	db *db.DB
}

// GudangSecurityOracle menyusun gudang security di atas Oracle.
func GudangSecurityOracle(svc *Service) GudangSecurityTCO {
	return gudangSecurityOracle{MasterSecurityTCO: repository.NewMasterSecurityTCO(svc.DB()), db: svc.DB()}
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
	transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error
}

// SecurityTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) SecurityTCO() *SecurityTCO {
	return &SecurityTCO{svc: s, gudang: securityBelumDisuntik{}, reinsurer: reinsurerBelumDisuntik{},
		kontrak: pemegangKontrakBelumDisuntik{}, tahun: gudangTahunTreatyBelumDisuntik{},
		master: masterReinsurerBelumDisuntik{}, transaksi: s.DalamTransaksi}
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

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *SecurityTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *SecurityTCO {
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
func (l *SecurityTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID, reinsurerID string) (
	DaftarSecurityTampil, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
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
func (l *SecurityTCO) Simpan(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID, reinsurerID string,
	m SecurityMasuk) (SecurityTampil, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
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
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		if s.ID != "" {
			lama, err := l.gudang.Ambil(ctx, r.ID, s.ID)
			if err != nil {
				return err
			}
			// Kolom yang UPDATE warisan tidak sentuh dipertahankan dari barisnya.
			s.TopID, s.TpTreaty, s.UserID = lama.TopID, lama.TpTreaty, lama.UserID
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
		// tco4: MTREATYSECURITY tanpa identitas - "ID" security = REAS_SECURITY
		// terpangkas; mengganti nama security mengganti kuncinya.
		s.ID = strings.TrimSpace(s.ReasSecurity)
		return nil
	})
	if err != nil {
		return SecurityTampil{}, err
	}
	tampil.SecurityReinsurer = s
	return TampilSecurity(tampil), nil
}

// Hapus membuang SATU security - `Delete` b17559 (`DeleteSecurityReinsurer`).
// Reinsurer induknya tidak disentuh.
func (l *SecurityTCO) Hapus(ctx context.Context, pelaku inti.Pelaku, tahunID, kontrakID, reinsurerID, id string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	r, err := l.induk(ctx, tahunID, kontrakID, reinsurerID)
	if err != nil {
		return "", err
	}
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		if err := l.gudang.Hapus(ctx, tx, r.ID, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	// `DeleteSecurityReinsurer` tidak menampilkan pesan - `[tidak ada di korpus]`.
	return "Security with ID " + id + " deleted", nil
}
