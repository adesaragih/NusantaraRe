package services

// Kontrak treaty di dalam tahun treaty - tiket 04 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan editor kontrak `InputTreatyContractReinsType`
// (harness `InboxTreatyContractReinsType` b151, dibuka sebagai popup dari
// tombol `ReinsType` baris tahun treaty `InputTreatyContract.xml` b20778).
//
// ⛔ Gerbangnya dua lapis: yang HIDUP di Pega (`SetTanggalTreatyContract`
// langkah 4: tahun mulai = tahun treaty) dan yang dituntut AC tiket 04 walau
// langkahnya DIKOMENTARI di `SaveTreatyContract_Act` (jenis reasuransi wajib
// dari daftar tersaring, tanggal wajib, satu jenis satu kontrak per tahun).
// Ralatnya bertanggal di tiket.
//
// Dibaca sesudah: tco_tahun.go, tco_jenisreasuransi.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/repository"
)

var (
	// ErrKontrakBeranak - jenis reasuransi tidak dapat diganti selama kombinasinya beranak (409).
	ErrKontrakBeranak = errors.New("services: the contract reinsurance type cannot be changed while reinsurer/business rows still exist")
	// ErrGudangKontrakBelumDisuntik - gudang kontrak belum dipasang.
	ErrGudangKontrakBelumDisuntik = errors.New("services: treaty contract store is not injected")
	// ErrKontrakTidakAda - kontrak bukan milik tahun treaty itu (404).
	ErrKontrakTidakAda = repository.ErrKontrakTidakAda
	// ErrKontrakDobel - jenis reasuransi itu sudah punya kontrak di tahun ini (409).
	ErrKontrakDobel = errors.New("services: duplicate contract")
	// ErrJenisReasuransiDiLuarDaftar - ID jenis reasuransi tidak ada di daftar
	// tersaring tiket 02 (AC 10, 11): tidak diketik bebas (422).
	ErrJenisReasuransiDiLuarDaftar = errors.New("services: reinsurance type outside the filtered REINSURANCETYPE list")
)

// PesanKontrakDobelTCO - `SaveTreatyContract_Act.xml` langkah 12 (b2270
// `OutputParam.ERRMSG3` "Data sudah pernah di Input"), langkah yang dikomentari
// di Pega; diterjemahkan [keputusan work owner 30-09-2026: bahasa Inggris].
const PesanKontrakDobelTCO = "Data has already been entered"

// GalatKontrakDobel menyebut kontrak mana yang sudah memegang jenis itu.
type GalatKontrakDobel struct {
	IDLain      string
	ReinsTypeID string
}

func (g GalatKontrakDobel) Error() string {
	return fmt.Sprintf("%s: reinsurance type %s is already used by contract %s in this treaty year",
		PesanKontrakDobelTCO, g.ReinsTypeID, g.IDLain)
}

// Is membuat `errors.Is(err, ErrKontrakDobel)` benar.
func (GalatKontrakDobel) Is(target error) bool { return target == ErrKontrakDobel }

// GudangKontrakTCO membaca dan menulis kontrak.
type GudangKontrakTCO interface {
	Daftar(ctx context.Context, tahunID string) ([]models.KontrakTreaty, error)
	Ambil(ctx context.Context, tahunID, id string) (models.KontrakTreaty, error)
	Sisip(ctx context.Context, tx *db.Tx, k models.KontrakTreaty) (string, error)
	Perbarui(ctx context.Context, tx *db.Tx, k models.KontrakTreaty) error
	CariDobel(ctx context.Context, tx *db.Tx, tahunID, reinsTypeID, kecualiID string) (string, error)
	JumlahAnakKombinasi(ctx context.Context, tx *db.Tx, kom models.KombinasiTCO, tahunID string) (int64, error)
}

type gudangKontrakBelumDisuntik struct{}

func (gudangKontrakBelumDisuntik) JumlahAnakKombinasi(context.Context, *db.Tx, models.KombinasiTCO, string) (int64, error) {
	return 0, ErrGudangKontrakBelumDisuntik
}

func (gudangKontrakBelumDisuntik) Daftar(context.Context, string) ([]models.KontrakTreaty, error) {
	return nil, ErrGudangKontrakBelumDisuntik
}
func (gudangKontrakBelumDisuntik) Ambil(context.Context, string, string) (models.KontrakTreaty, error) {
	return models.KontrakTreaty{}, ErrGudangKontrakBelumDisuntik
}
func (gudangKontrakBelumDisuntik) Sisip(context.Context, *db.Tx, models.KontrakTreaty) (string, error) {
	return "", ErrGudangKontrakBelumDisuntik
}
func (gudangKontrakBelumDisuntik) Perbarui(context.Context, *db.Tx, models.KontrakTreaty) error {
	return ErrGudangKontrakBelumDisuntik
}
func (gudangKontrakBelumDisuntik) CariDobel(context.Context, *db.Tx, string, string, string) (string, error) {
	return "", ErrGudangKontrakBelumDisuntik
}

type gudangKontrakOracle struct {
	m  *repository.MasterKontrakTCO
	db *db.DB
}

func (g gudangKontrakOracle) Daftar(ctx context.Context, tahunID string) ([]models.KontrakTreaty, error) {
	return g.m.Daftar(ctx, tahunID)
}
func (g gudangKontrakOracle) Ambil(ctx context.Context, tahunID, id string) (models.KontrakTreaty, error) {
	return g.m.Ambil(ctx, tahunID, id)
}
func (g gudangKontrakOracle) Sisip(ctx context.Context, tx *db.Tx, k models.KontrakTreaty) (string, error) {
	return g.m.Sisip(ctx, tx, k)
}
func (g gudangKontrakOracle) Perbarui(ctx context.Context, tx *db.Tx, k models.KontrakTreaty) error {
	return g.m.Perbarui(ctx, tx, k)
}
func (g gudangKontrakOracle) CariDobel(ctx context.Context, tx *db.Tx, tahunID, reinsTypeID, kecualiID string) (string, error) {
	return g.m.CariDobel(ctx, tx, tahunID, reinsTypeID, kecualiID)
}
func (g gudangKontrakOracle) JumlahAnakKombinasi(ctx context.Context, tx *db.Tx, kom models.KombinasiTCO,
	tahunID string) (int64, error) {
	return g.m.JumlahAnakKombinasi(ctx, tx, kom, tahunID)
}

// GudangKontrakOracle menyusun gudang kontrak di atas Oracle.
func GudangKontrakOracle(svc *Service) GudangKontrakTCO {
	return gudangKontrakOracle{m: repository.NewMasterKontrakTCO(svc.DB()), db: svc.DB()}
}

// KontrakMasuk adalah badan simpan kontrak.
//
// ⛔ `ReinsTypeName` DIABAIKAN: namanya diambil dari master tersaring
// berdasarkan `ReinsTypeID` (AC 10, 11 - tidak diketik bebas). Medannya ada
// hanya supaya klien yang mengirimnya tidak ditolak.
type KontrakMasuk struct {
	ID            string `json:"id"`
	ReinsTypeID   string `json:"reinsTypeId"`
	ReinsTypeName string `json:"reinsTypeName"`
	// TreatyStartDate, TreatyEndDate - DIABAIKAN server sejak 30-09-2026:
	// tanggal kontrak = tanggal tahun treaty induknya (`Simpan`).
	TreatyStartDate string `json:"treatyStartDate"`
	TreatyEndDate   string `json:"treatyEndDate"`
}

// KontrakTampil adalah satu kontrak seperti dikirim ke layar.
type KontrakTampil struct {
	ID              string `json:"id"`
	IDTreatyYear    string `json:"idTreatyYear"`
	ReinsTypeID     string `json:"reinsTypeId"`
	ReinsTypeName   string `json:"reinsTypeName"`
	TreatyStartDate string `json:"treatyStartDate"`
	TreatyEndDate   string `json:"treatyEndDate"`
	UserID          string `json:"userId"`
	TglUpdate       string `json:"tglUpdate"`
}

// TampilKontrak menerjemahkan satu baris.
func TampilKontrak(k models.KontrakTreaty) KontrakTampil {
	return KontrakTampil{ID: k.ID, IDTreatyYear: k.IDTreatyYear, ReinsTypeID: k.ReinsTypeID,
		ReinsTypeName: k.ReinsTypeName, TreatyStartDate: utils.FormatTanggal(k.TreatyStartDate),
		TreatyEndDate: utils.FormatTanggal(k.TreatyEndDate), UserID: k.UserID,
		TglUpdate: utils.FormatTanggalWaktu(k.TglUpdate)}
}

// KontrakTreatyTCO melayani editor kontrak di dalam tahun treaty.
type KontrakTreatyTCO struct {
	svc       *Service
	gudang    GudangKontrakTCO
	tahun     PemeriksaTahunTCO
	jenis     PembacaJenisReasuransiTCO
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error
}

// KontrakTreatyTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) KontrakTreatyTCO() *KontrakTreatyTCO {
	return &KontrakTreatyTCO{svc: s, gudang: gudangKontrakBelumDisuntik{},
		tahun: gudangTahunTreatyBelumDisuntik{}, jenis: pembacaJenisReasuransiBelumDisuntik{},
		jam: time.Now, transaksi: s.DalamTransaksi}
}

func (k *KontrakTreatyTCO) salin() *KontrakTreatyTCO { s := *k; return &s }

// DenganGudang memasang gudang kontrak.
func (k *KontrakTreatyTCO) DenganGudang(g GudangKontrakTCO) *KontrakTreatyTCO {
	s := k.salin()
	s.gudang = g
	return s
}

// DenganTahun memasang pemeriksa tahun treaty induk.
func (k *KontrakTreatyTCO) DenganTahun(t PemeriksaTahunTCO) *KontrakTreatyTCO {
	s := k.salin()
	s.tahun = t
	return s
}

// DenganJenis memasang pembaca daftar jenis reasuransi tersaring (tiket 02).
func (k *KontrakTreatyTCO) DenganJenis(j PembacaJenisReasuransiTCO) *KontrakTreatyTCO {
	s := k.salin()
	s.jenis = j
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (k *KontrakTreatyTCO) DenganJam(j func() time.Time) *KontrakTreatyTCO {
	s := k.salin()
	s.jam = j
	return s
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (k *KontrakTreatyTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *KontrakTreatyTCO {
	s := k.salin()
	s.transaksi = f
	return s
}

// Daftar membaca kontrak satu tahun treaty - grid `BrowseTreatyContract_RD`.
func (k *KontrakTreatyTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, tahunID string) ([]KontrakTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	baris, err := k.gudang.Daftar(ctx, tahunID)
	if err != nil {
		return nil, err
	}
	if _, err := k.tahun.Ambil(ctx, tahunID); err != nil {
		return nil, err
	}
	hasil := make([]KontrakTampil, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, TampilKontrak(b))
	}
	return hasil, nil
}

func uraiTanggalKontrak(nama, teks string) (time.Time, error) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return time.Time{}, nil
	}
	v, err := utils.ParseTanggal(t)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s %q is not a YYYY-MM-DD date", galat.ErrPermintaanTidakSah, nama, t)
	}
	return v, nil
}

// namaJenisDariMaster mencari `.Note` ID pilihan di daftar tersaring.
func (k *KontrakTreatyTCO) namaJenisDariMaster(ctx context.Context, id string) (string, error) {
	daftar, err := k.jenis.DaftarNonLife(ctx)
	if err != nil {
		return "", err
	}
	if len(daftar) == 0 {
		return "", ErrMasterJenisReasuransiKosong
	}
	for _, j := range daftar {
		if j.ID == id {
			return j.Note, nil
		}
	}
	return "", fmt.Errorf("%w: ReinsTypeID %q", ErrJenisReasuransiDiLuarDaftar, id)
}

// Simpan menulis kontrak baru atau memperbarui yang ada - `Save` b3618
// (`SaveTreatyContract_Act`).
//
// Urutannya: identitas -> tahun induk ada -> tanggal dari tahun induk -> gerbang
// murni -> jenis dari master -> satu transaksi {dobel, sisip/perbarui} (tco4: tanpa jejak).
//
// ⛔ Tanggal kontrak = Start/End Date TAHUN treaty induknya [keputusan work owner 30-09-2026: "start date dan end date pada ReinsType read only, datanya diambil dari depan"] -
// "dari depan sampai belakang tanggalnya sama". Tanggal kiriman klien
// DIABAIKAN; gerbang periode dan "tahun mulai = tahun treaty" tetap berjalan
// atas tanggal tahun itu.
func (k *KontrakTreatyTCO) Simpan(ctx context.Context, pelaku inti.Pelaku, tahunID string, m KontrakMasuk) (KontrakTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KontrakTampil{}, err
	}
	tahun, err := k.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return KontrakTampil{}, err
	}
	kontrak := models.KontrakTreaty{ID: strings.TrimSpace(m.ID), IDTreatyYear: tahunID,
		ReinsTypeID: strings.TrimSpace(m.ReinsTypeID), TreatyStartDate: tahun.StartDate, TreatyEndDate: tahun.EndDate,
		UserID: pelaku.AkunID, TglUpdate: k.jam()}
	if err := models.PeriksaKontrakTreaty(kontrak, tahun.TreatyYear); err != nil {
		return KontrakTampil{}, err
	}
	if kontrak.ReinsTypeName, err = k.namaJenisDariMaster(ctx, kontrak.ReinsTypeID); err != nil {
		return KontrakTampil{}, err
	}
	err = k.transaksi(ctx, func(tx *db.Tx) error {
		lain, err := k.gudang.CariDobel(ctx, tx, tahunID, kontrak.ReinsTypeID, kontrak.ID)
		if err != nil {
			return err
		}
		// Satu jenis reasuransi satu kontrak per tahun [keputusan work owner 29-09-2026]
		// (OQ-TCO-11, ditutup).
		if lain != "" {
			return GalatKontrakDobel{IDLain: lain, ReinsTypeID: kontrak.ReinsTypeID}
		}
		// Temuan /code-review: mengganti jenis reasuransi mengganti KOMBINASI
		// tempat reinsurer, security, dan business menggantung - ditolak selama
		// anak-anaknya masih ada (seperti induk klausul beranak).
		if kontrak.ID != "" {
			lama, err := k.gudang.Ambil(repository.DenganBacaTxTCO(ctx, tx), tahunID, kontrak.ID)
			if err != nil {
				return err
			}
			if lama.ReinsTypeID != kontrak.ReinsTypeID {
				n, err := k.gudang.JumlahAnakKombinasi(ctx, tx, models.KombinasiDari(tahun, lama), tahunID)
				if err != nil {
					return err
				}
				if n > 0 {
					return fmt.Errorf("%w: contract %s still has %d reinsurer/business rows on type %s", ErrKontrakBeranak,
						kontrak.ID, n, lama.ReinsTypeID)
				}
			}
		}
		if kontrak.ID == "" {
			if kontrak.ID, err = k.gudang.Sisip(ctx, tx, kontrak); err != nil {
				return err
			}
		} else if err := k.gudang.Perbarui(ctx, tx, kontrak); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return KontrakTampil{}, err
	}
	return TampilKontrak(kontrak), nil
}

// AkhirBawaan menghitung tanggal akhir yang layar isikan saat tanggal mulai
// dipilih: mulai + 1 tahun kalender (models.AkhirKontrakBawaanTCO, OQ-TCO-10).
func (k *KontrakTreatyTCO) AkhirBawaan(ctx context.Context, pelaku inti.Pelaku, tahunID, mulai string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	if _, err := k.tahun.Ambil(ctx, tahunID); err != nil {
		return "", err
	}
	return AkhirTahunBawaan(pelaku, mulai)
}

// AkhirTahunBawaan - tanggal akhir bawaan dari tanggal mulai, MULAI + 1 TAHUN
// KALENDER (`models.AkhirKontrakBawaanTCO`, OQ-TCO-10). Dipakai kontrak DAN
// tahun treaty baru yang belum ber-ID [keputusan work owner 30-09-2026] - satu
// aturan, satu tempat.
func AkhirTahunBawaan(pelaku inti.Pelaku, mulai string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	t, err := uraiTanggalKontrak("start date", mulai)
	if err != nil {
		return "", err
	}
	if t.IsZero() {
		return "", fmt.Errorf("%w: start date is required", galat.ErrPermintaanTidakSah)
	}
	return utils.FormatTanggal(models.AkhirKontrakBawaanTCO(t)), nil
}
