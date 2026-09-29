package services

// Reinsurer pada kombinasi (tahun, grup, jenis) - tiket 05 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan panel `ViewDetailTreatyReinsurerGrid1`
// (disertakan `InputTreatyContractReinsType.xml` b13311), dibuka tombol baris
// kontrak `Reinsurer List` b11308 (`BrowseTreatyReinsurerList_Act` +
// `SetTreatyReinsurerList_Act`).
//
// ⛔ Kombinasinya dibuka KONTRAK: jalur HTTP menyebut tahun dan kontraknya, dan
// kombinasi (TREATYYEAR tahun, TREATYGROUPID tahun, REINSTYPEID kontrak)
// diturunkan server. Rekam reinsurer sendiri tidak menyimpan ID kontrak.
//
// ⛔ Kontrak DIKUNCI selama transaksi penulis: total share dijumlah lalu
// ditulis, dan dua penulis serentak tidak boleh sama-sama lolos batas 100.
//
// Dibaca sesudah: tco_kontrak.go, models/tco_reinsurer.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/pkg/utils"
)

var (
	// ErrGudangReinsurerBelumDisuntik - gudang reinsurer belum dipasang.
	ErrGudangReinsurerBelumDisuntik = errors.New("services: gudang reinsurer belum disuntik")
	// ErrReinsurerTidakAda - reinsurer bukan milik kombinasi kontrak itu (404).
	ErrReinsurerTidakAda = repository.ErrReinsurerTidakAda
	// ErrReinsurerDiLuarMaster - ReinsurerID tidak ada di master AGENT aktif (422).
	ErrReinsurerDiLuarMaster = repository.ErrReinsurerMasterTidakAda
)

// GudangReinsurerTCO membaca dan menulis reinsurer + jejaknya.
type GudangReinsurerTCO interface {
	Daftar(ctx context.Context, k models.KombinasiTCO) ([]models.ReinsurerTreaty, error)
	Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.ReinsurerTreaty, error)
	ShareLain(ctx context.Context, tx *repository.Tx, k models.KombinasiTCO, kecualiID string) ([]*apd.Decimal, error)
	Sisip(ctx context.Context, tx *repository.Tx, r models.ReinsurerTreaty) (string, error)
	Perbarui(ctx context.Context, tx *repository.Tx, r models.ReinsurerTreaty) error
}

// PemegangKontrakTCO membaca dan mengunci kontrak pembuka kombinasi.
type PemegangKontrakTCO interface {
	Ambil(ctx context.Context, tahunID, id string) (models.KontrakTreaty, error)
	Kunci(ctx context.Context, tx *repository.Tx, tahunID, id string) error
}

// PembacaReinsurerMasterTCO membaca master `AGENT` (dibaca saja).
type PembacaReinsurerMasterTCO interface {
	Cari(ctx context.Context, teks string) ([]repository.ReinsurerMasterTCO, error)
	Ambil(ctx context.Context, id string) (repository.ReinsurerMasterTCO, error)
}

type reinsurerBelumDisuntik struct{}

func (reinsurerBelumDisuntik) Daftar(context.Context, models.KombinasiTCO) ([]models.ReinsurerTreaty, error) {
	return nil, ErrGudangReinsurerBelumDisuntik
}
func (reinsurerBelumDisuntik) Ambil(context.Context, models.KombinasiTCO, string) (models.ReinsurerTreaty, error) {
	return models.ReinsurerTreaty{}, ErrGudangReinsurerBelumDisuntik
}
func (reinsurerBelumDisuntik) ShareLain(context.Context, *repository.Tx, models.KombinasiTCO, string) ([]*apd.Decimal, error) {
	return nil, ErrGudangReinsurerBelumDisuntik
}
func (reinsurerBelumDisuntik) Sisip(context.Context, *repository.Tx, models.ReinsurerTreaty) (string, error) {
	return "", ErrGudangReinsurerBelumDisuntik
}
func (reinsurerBelumDisuntik) Perbarui(context.Context, *repository.Tx, models.ReinsurerTreaty) error {
	return ErrGudangReinsurerBelumDisuntik
}

type pemegangKontrakBelumDisuntik struct{}

func (pemegangKontrakBelumDisuntik) Ambil(context.Context, string, string) (models.KontrakTreaty, error) {
	return models.KontrakTreaty{}, ErrGudangReinsurerBelumDisuntik
}
func (pemegangKontrakBelumDisuntik) Kunci(context.Context, *repository.Tx, string, string) error {
	return ErrGudangReinsurerBelumDisuntik
}

type masterReinsurerBelumDisuntik struct{}

func (masterReinsurerBelumDisuntik) Cari(context.Context, string) ([]repository.ReinsurerMasterTCO, error) {
	return nil, ErrGudangReinsurerBelumDisuntik
}
func (masterReinsurerBelumDisuntik) Ambil(context.Context, string) (repository.ReinsurerMasterTCO, error) {
	return repository.ReinsurerMasterTCO{}, ErrGudangReinsurerBelumDisuntik
}

type gudangReinsurerOracle struct {
	m  *repository.MasterReinsurerTCO
	db *repository.DB
}

func (g gudangReinsurerOracle) Daftar(ctx context.Context, k models.KombinasiTCO) ([]models.ReinsurerTreaty, error) {
	return g.m.Daftar(ctx, k)
}
func (g gudangReinsurerOracle) Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.ReinsurerTreaty, error) {
	return g.m.Ambil(ctx, k, id)
}
func (g gudangReinsurerOracle) ShareLain(ctx context.Context, tx *repository.Tx, k models.KombinasiTCO, kecualiID string) ([]*apd.Decimal, error) {
	return g.m.ShareLain(ctx, tx, k, kecualiID)
}
func (g gudangReinsurerOracle) Sisip(ctx context.Context, tx *repository.Tx, r models.ReinsurerTreaty) (string, error) {
	return g.m.Sisip(ctx, tx, r)
}
func (g gudangReinsurerOracle) Perbarui(ctx context.Context, tx *repository.Tx, r models.ReinsurerTreaty) error {
	return g.m.Perbarui(ctx, tx, r)
}

// GudangReinsurerOracle menyusun gudang reinsurer di atas Oracle.
func GudangReinsurerOracle(svc *Service) GudangReinsurerTCO {
	return gudangReinsurerOracle{m: repository.NewMasterReinsurerTCO(svc.db), db: svc.db}
}

// PemegangKontrakOracle menyusun pembaca + pengunci kontrak.
func PemegangKontrakOracle(svc *Service) PemegangKontrakTCO {
	return repository.NewMasterKontrakTCO(svc.db)
}

// MasterReinsurerOracle menyusun pembaca master `AGENT`.
func MasterReinsurerOracle(svc *Service) PembacaReinsurerMasterTCO {
	return repository.NewMasterReinsurerAgent(svc.db)
}

// ReinsurerMasuk adalah badan simpan - medan yang TAMPIL di form
// (`Reins.ID`/`Reinsurer` b8042/b8226, `%Share` b8522, `%Comm` b8800,
// `Rating` b9076). Delapan medan lain form itu tersembunyi permanen
// (`pyCondition 1=2`) dan tidak diterima dari klien.
type ReinsurerMasuk struct {
	ID          string `json:"id"`
	ReinsurerID string `json:"reinsurerId"`
	// PctShare dan Ricomm TEKS - koma atau titik desimal; dinormalkan sekali
	// di `models.UraiPersenMasukTCO`.
	PctShare  string `json:"pctShare"`
	Ricomm    string `json:"ricomm"`
	StdRating string `json:"stdRating"`
}

// ReinsurerTampil adalah satu reinsurer seperti dikirim ke layar.
//
// ⛔ Uang dan persen dikirim sebagai TEKS (ADR-0003): angka JavaScript adalah
// float64.
type ReinsurerTampil struct {
	ID              string `json:"id"`
	TreatyYear      string `json:"treatyYear"`
	TreatyGroupID   string `json:"treatyGroupId"`
	TreatyGroupName string `json:"treatyGroupName"`
	ReinsTypeID     string `json:"reinsTypeId"`
	ReinsTypeName   string `json:"reinsTypeName"`
	ReinsurerID     string `json:"reinsurerId"`
	ClientID        string `json:"clientId"`
	Name            string `json:"name"`
	Ricomm          string `json:"ricomm"`
	PctShare        string `json:"pctShare"`
	IUDate          string `json:"iuDate"`
	UserID          string `json:"userId"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	StatusOn        string `json:"statusOn"`
	StdRating       string `json:"stdRating"`
	OperatorName    string `json:"operatorName"`
	TglUpdate       string `json:"tglUpdate"`
}

// TampilReinsurer menerjemahkan satu baris.
func TampilReinsurer(r models.ReinsurerTreaty) ReinsurerTampil {
	return ReinsurerTampil{ID: r.ID, TreatyYear: r.TreatyYear, TreatyGroupID: r.TreatyGroupID,
		TreatyGroupName: r.TreatyGroupName, ReinsTypeID: r.ReinsTypeID, ReinsTypeName: r.ReinsTypeName,
		ReinsurerID: r.ReinsurerID, ClientID: r.ClientID, Name: r.Name, Ricomm: utils.FormatDecimal(r.Ricomm),
		PctShare: utils.FormatDecimal(r.PctShare), IUDate: r.IUDate, UserID: r.UserID,
		StartDate: utils.FormatTanggal(r.StartDate), EndDate: utils.FormatTanggal(r.EndDate), StatusOn: r.StatusOn,
		StdRating: r.StdRating, OperatorName: r.OperatorName, TglUpdate: utils.FormatTanggalWaktu(r.TglUpdate)}
}

// KombinasiTampil adalah kunci gabungan seperti dikirim ke layar
// (`OutputData.HASIL1/2/3` b2284-b2296).
type KombinasiTampil struct {
	TreatyYear      string `json:"treatyYear"`
	TreatyGroupID   string `json:"treatyGroupId"`
	TreatyGroupName string `json:"treatyGroupName"`
	ReinsTypeID     string `json:"reinsTypeId"`
	ReinsTypeName   string `json:"reinsTypeName"`
}

// TampilKombinasi menerjemahkan kombinasi.
func TampilKombinasi(k models.KombinasiTCO) KombinasiTampil {
	return KombinasiTampil{TreatyYear: k.TreatyYear, TreatyGroupID: k.TreatyGroupID,
		TreatyGroupName: k.TreatyGroupName, ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName}
}

// DaftarReinsurerTampil adalah grid reinsurer + baris `Total Share -->>` b6186.
type DaftarReinsurerTampil struct {
	Daftar     []ReinsurerTampil `json:"daftar"`
	Total      int               `json:"total"`
	TotalShare string            `json:"totalShare"`
	Kombinasi  KombinasiTampil   `json:"kombinasi"`
}

// HasilReinsurerTampil adalah jawaban simpan.
type HasilReinsurerTampil struct {
	Reinsurer  ReinsurerTampil `json:"reinsurer"`
	TotalShare string          `json:"totalShare"`
}

// ReinsurerMasterTampil adalah satu pilihan reinsurer.
type ReinsurerMasterTampil struct {
	ID         string `json:"id"`
	ClientName string `json:"clientName"`
	ClientID   string `json:"clientId"`
}

// ReinsurerTCO melayani panel reinsurer.
type ReinsurerTCO struct {
	svc       *Service
	gudang    GudangReinsurerTCO
	kontrak   PemegangKontrakTCO
	tahun     PemeriksaTahunTCO
	master    PembacaReinsurerMasterTCO
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *repository.Tx) error) error
}

// ReinsurerTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) ReinsurerTCO() *ReinsurerTCO {
	return &ReinsurerTCO{svc: s, gudang: reinsurerBelumDisuntik{}, kontrak: pemegangKontrakBelumDisuntik{},
		tahun: gudangTahunTreatyBelumDisuntik{}, master: masterReinsurerBelumDisuntik{},
		jam: time.Now, transaksi: s.DalamTransaksi}
}

func (l *ReinsurerTCO) salin() *ReinsurerTCO { s := *l; return &s }

// DenganGudang memasang gudang reinsurer.
func (l *ReinsurerTCO) DenganGudang(g GudangReinsurerTCO) *ReinsurerTCO {
	s := l.salin()
	s.gudang = g
	return s
}

// DenganKontrak memasang pemegang kontrak pembuka kombinasi.
func (l *ReinsurerTCO) DenganKontrak(k PemegangKontrakTCO) *ReinsurerTCO {
	s := l.salin()
	s.kontrak = k
	return s
}

// DenganTahun memasang pemeriksa tahun treaty.
func (l *ReinsurerTCO) DenganTahun(t PemeriksaTahunTCO) *ReinsurerTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganMaster memasang pembaca master reinsurer.
func (l *ReinsurerTCO) DenganMaster(m PembacaReinsurerMasterTCO) *ReinsurerTCO {
	s := l.salin()
	s.master = m
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *ReinsurerTCO) DenganJam(j func() time.Time) *ReinsurerTCO {
	s := l.salin()
	s.jam = j
	return s
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *ReinsurerTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *repository.Tx) error) error) *ReinsurerTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

// kombinasi menurunkan kunci gabungan dari tahun dan kontraknya.
func (l *ReinsurerTCO) kombinasi(ctx context.Context, tahunID, kontrakID string) (models.KombinasiTCO, error) {
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

// Daftar membaca reinsurer kombinasi + total share-nya.
func (l *ReinsurerTCO) Daftar(ctx context.Context, pelaku Pelaku, tahunID, kontrakID string) (DaftarReinsurerTampil, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return DaftarReinsurerTampil{}, err
	}
	k, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return DaftarReinsurerTampil{}, err
	}
	baris, err := l.gudang.Daftar(ctx, k)
	if err != nil {
		return DaftarReinsurerTampil{}, err
	}
	shares := make([]*apd.Decimal, 0, len(baris))
	hasil := make([]ReinsurerTampil, 0, len(baris))
	for _, b := range baris {
		shares = append(shares, b.PctShare)
		hasil = append(hasil, TampilReinsurer(b))
	}
	total, err := models.TotalShareTCO(shares)
	if err != nil {
		return DaftarReinsurerTampil{}, err
	}
	return DaftarReinsurerTampil{Daftar: hasil, Total: len(hasil), TotalShare: utils.FormatDecimal(total),
		Kombinasi: TampilKombinasi(k)}, nil
}

// Simpan menulis reinsurer baru atau memperbarui yang ada - `Save` b11405
// (`SaveTreatyReinsurerDetail1_Act`).
func (l *ReinsurerTCO) Simpan(ctx context.Context, pelaku Pelaku, tahunID, kontrakID string, m ReinsurerMasuk) (
	HasilReinsurerTampil, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return HasilReinsurerTampil{}, err
	}
	k, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return HasilReinsurerTampil{}, err
	}
	r := models.ReinsurerTreaty{ID: strings.TrimSpace(m.ID), TreatyYear: k.TreatyYear, TreatyGroupID: k.TreatyGroupID,
		TreatyGroupName: k.TreatyGroupName, ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName,
		ReinsurerID: strings.TrimSpace(m.ReinsurerID), StdRating: strings.TrimSpace(m.StdRating),
		OperatorName: pelaku.AkunID, TglUpdate: l.jam()}
	if err := models.PeriksaReinsurerTCO(r); err != nil {
		return HasilReinsurerTampil{}, err
	}
	if r.PctShare, err = models.UraiPersenMasukTCO("PctShare", m.PctShare); err != nil {
		return HasilReinsurerTampil{}, err
	}
	if r.Ricomm, err = models.UraiPersenMasukTCO("Ricomm", m.Ricomm); err != nil {
		return HasilReinsurerTampil{}, err
	}
	master, err := l.master.Ambil(ctx, r.ReinsurerID)
	if err != nil {
		if errors.Is(err, ErrReinsurerDiLuarMaster) {
			return HasilReinsurerTampil{}, fmt.Errorf("%w: ReinsurerID %q", err, r.ReinsurerID)
		}
		return HasilReinsurerTampil{}, err
	}
	// ⛔ NAME dan CLIENTID dari MASTER (pemilih b8285: `.ClientName` -> NAME,
	// `.ClientID` -> CLIENTID) - tidak pernah dari klien.
	r.Name, r.ClientID = master.ClientName, master.ClientID

	var total *apd.Decimal
	err = l.transaksi(ctx, func(tx *repository.Tx) error {
		if err := l.kontrak.Kunci(ctx, tx, tahunID, kontrakID); err != nil {
			return err
		}
		if r.ID == "" {
			r.UserID = pelaku.AkunID
		} else {
			lama, err := l.gudang.Ambil(ctx, k, r.ID)
			if err != nil {
				return err
			}
			// Medan tersembunyi form dipertahankan dari barisnya.
			r.IUDate, r.StartDate, r.EndDate, r.StatusOn = lama.IUDate, lama.StartDate, lama.EndDate, lama.StatusOn
			r.UserID = lama.UserID
			if strings.TrimSpace(r.UserID) == "" {
				r.UserID = pelaku.AkunID
			}
		}
		lain, err := l.gudang.ShareLain(ctx, tx, k, r.ID)
		if err != nil {
			return err
		}
		if total, err = models.PeriksaTotalShareTCO(lain, r.PctShare); err != nil {
			return err
		}
		if r.ID == "" {
			if r.ID, err = l.gudang.Sisip(ctx, tx, r); err != nil {
				return err
			}
		} else if err := l.gudang.Perbarui(ctx, tx, r); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return HasilReinsurerTampil{}, err
	}
	return HasilReinsurerTampil{Reinsurer: TampilReinsurer(r), TotalShare: utils.FormatDecimal(total)}, nil
}

// CariMaster membaca pilihan reinsurer - `BrowseAgentReinsSOA_RD`.
func (l *ReinsurerTCO) CariMaster(ctx context.Context, pelaku Pelaku, teks string) ([]ReinsurerMasterTampil, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	d, err := l.master.Cari(ctx, teks)
	if err != nil {
		return nil, err
	}
	hasil := make([]ReinsurerMasterTampil, 0, len(d))
	for _, r := range d {
		hasil = append(hasil, ReinsurerMasterTampil{ID: r.ID, ClientName: r.ClientName, ClientID: r.ClientID})
	}
	return hasil, nil
}
