package services

// Tahun treaty - tiket 03 Treaty Contract Out.
//
// Untuk apa berkas ini: orkestrasi `SaveTreatyYear_Act` (b280 UserID dari
// operator, b327 TglUpdate sekarang, prasyarat b388/b411/b434, RDB-List
// `SaveMasterTreatyYear_SQL` -> upsert ber-ID) TANPA memanggil procedure
// (keputusan o), ditambah dua gerbang sadar: periode terbalik ditolak (AC 9)
// dan anti-dobel (StartDate, EndDate, TreatyGroupID) ditolak (AC 73).
//
// ⛔ Identitas TIDAK PERNAH datang dari pemanggil untuk baris baru (AC 5):
// ID kosong = baru, dari sequence; ID terisi = pembaruan seluruh medan (AC 8).
//
// ⛔ Anti-dobel diperiksa DI DALAM transaksi penyimpanan - ia pertanyaan
// keunikan di basis data, dan jawabannya menyebut tahun treaty mana yang
// sudah memakai kombinasi itu.
//
// ⛔ tco4: nol jejak modul (Pega tidak mencatatnya); gudang dan pelaksana
// transaksi DISUNTIK supaya gerbangnya teruji tanpa Oracle.
//
// ⛔ Fitur salin tahun treaty TIDAK ADA di sini maupun di mana pun (AC 72).

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
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
)

var (
	// ErrGrupTreatyDiLuarMaster - TreatyGroupID tidak ada di master TREATYGROUP (422).
	ErrGrupTreatyDiLuarMaster = errors.New("services: treaty group not in master")
	// ErrTahunBeranak - tahun/grup tidak dapat diganti selama tahun treaty beranak (409).
	ErrTahunBeranak = errors.New("services: the treaty year and its group cannot be changed while contracts/clauses still exist")
	// ErrGudangTahunTreatyBelumDisuntik - handler lupa memasang gudang.
	ErrGudangTahunTreatyBelumDisuntik = errors.New("services: treaty year store is not injected")
	// ErrTahunTreatyDobel - AC 73.
	ErrTahunTreatyDobel = errors.New("services: a treaty year with the same period and group already exists")
	// ErrTahunTreatyTidakAda dirujuk ulang supaya handler tidak mengimpor repository.
	ErrTahunTreatyTidakAda = repository.ErrTahunTreatyTidakAda
)

// GalatTahunTreatyDobel menyebut tahun treaty yang sudah memakai kombinasi itu.
type GalatTahunTreatyDobel struct {
	IDLain        string
	TreatyGroupID string
	StartDate     time.Time
	EndDate       time.Time
}

func (g GalatTahunTreatyDobel) Error() string {
	return fmt.Sprintf("treaty year %s already uses StartDate %s, EndDate %s, Treaty Group %s",
		g.IDLain, utils.FormatTanggal(g.StartDate), utils.FormatTanggal(g.EndDate), g.TreatyGroupID)
}

// Is membuat errors.Is(err, ErrTahunTreatyDobel) benar.
func (GalatTahunTreatyDobel) Is(target error) bool { return target == ErrTahunTreatyDobel }

// GudangTahunTreatyTCO adalah seluruh sentuhan basis data tahun treaty.
type GudangTahunTreatyTCO interface {
	Daftar(ctx context.Context, halaman, ukuran int) (repository.HalamanTahunTreaty, error)
	Ambil(ctx context.Context, id string) (models.TahunTreaty, error)
	Sisip(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (string, error)
	Perbarui(ctx context.Context, tx *db.Tx, t models.TahunTreaty) error
	CariDobel(ctx context.Context, tx *db.Tx, grupID string, mulai, akhir time.Time,
		kecualiID string) (string, error)
	JumlahAnak(ctx context.Context, tx *db.Tx, tahunID string) (int64, error)
}

type gudangTahunTreatyBelumDisuntik struct{}

func (gudangTahunTreatyBelumDisuntik) Daftar(context.Context, int, int) (repository.HalamanTahunTreaty, error) {
	return repository.HalamanTahunTreaty{}, ErrGudangTahunTreatyBelumDisuntik
}
func (gudangTahunTreatyBelumDisuntik) Ambil(context.Context, string) (models.TahunTreaty, error) {
	return models.TahunTreaty{}, ErrGudangTahunTreatyBelumDisuntik
}
func (gudangTahunTreatyBelumDisuntik) Sisip(context.Context, *db.Tx, models.TahunTreaty) (string, error) {
	return "", ErrGudangTahunTreatyBelumDisuntik
}
func (gudangTahunTreatyBelumDisuntik) Perbarui(context.Context, *db.Tx, models.TahunTreaty) error {
	return ErrGudangTahunTreatyBelumDisuntik
}
func (gudangTahunTreatyBelumDisuntik) CariDobel(context.Context, *db.Tx, string, time.Time, time.Time, string) (string, error) {
	return "", ErrGudangTahunTreatyBelumDisuntik
}
func (gudangTahunTreatyBelumDisuntik) JumlahAnak(context.Context, *db.Tx, string) (int64, error) {
	return 0, ErrGudangTahunTreatyBelumDisuntik
}

type gudangTahunTreatyOracle struct {
	svc  *Service
	baca *repository.MasterTahunTreaty
}

func (g gudangTahunTreatyOracle) Daftar(ctx context.Context, halaman, ukuran int) (repository.HalamanTahunTreaty, error) {
	return g.baca.Daftar(ctx, halaman, ukuran)
}
func (g gudangTahunTreatyOracle) Ambil(ctx context.Context, id string) (models.TahunTreaty, error) {
	return g.baca.Ambil(ctx, id)
}
func (g gudangTahunTreatyOracle) Sisip(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (string, error) {
	return g.baca.Sisip(ctx, tx, t)
}
func (g gudangTahunTreatyOracle) Perbarui(ctx context.Context, tx *db.Tx, t models.TahunTreaty) error {
	return g.baca.Perbarui(ctx, tx, t)
}
func (g gudangTahunTreatyOracle) CariDobel(ctx context.Context, tx *db.Tx, grupID string,
	mulai, akhir time.Time, kecualiID string) (string, error) {
	return g.baca.CariDobel(ctx, tx, grupID, mulai, akhir, kecualiID)
}
func (g gudangTahunTreatyOracle) JumlahAnak(ctx context.Context, tx *db.Tx, tahunID string) (int64, error) {
	return g.baca.JumlahAnak(ctx, tx, tahunID)
}

// GudangTahunTreatyOracle adalah gudang sungguhan, dipasang handler.
func GudangTahunTreatyOracle(svc *Service) GudangTahunTreatyTCO {
	return gudangTahunTreatyOracle{svc: svc, baca: repository.NewMasterTahunTreaty(svc.DB())}
}

// TahunTreatyMasuk adalah badan permintaan simpan.
//
// Tanggal TEKS `YYYY-MM-DD` (utils.ParseTanggal); kosong berarti kosong.
type TahunTreatyMasuk struct {
	ID               string `json:"id"`
	TreatyYear       string `json:"treatyYear"`
	UnderwritingYear string `json:"underwritingYear"`
	TreatyGroupID    string `json:"treatyGroupId"`
	TreatyGroupName  string `json:"treatyGroupName"`
	Proportion       string `json:"proportion"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
}

// TahunTreatyTampil adalah satu tahun treaty untuk layar.
//
// Nama kunci mengikuti properti Pega `InputTreatyYear.*`.
type TahunTreatyTampil struct {
	ID               string `json:"id"`
	TreatyYear       string `json:"treatyYear"`
	UnderwritingYear string `json:"underwritingYear"`
	TreatyGroupID    string `json:"treatyGroupId"`
	TreatyGroupName  string `json:"treatyGroupName"`
	Proportion       string `json:"proportion"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	UserID           string `json:"userId"`
	TglUpdate        string `json:"tglUpdate"`
}

// HalamanTahunTreatyTampil adalah satu halaman daftar.
type HalamanTahunTreatyTampil struct {
	Baris   []TahunTreatyTampil `json:"baris"`
	Total   int                 `json:"total"`
	Halaman int                 `json:"halaman"`
	Ukuran  int                 `json:"ukuran"`
}

// TampilTahunTreaty mengubah model menjadi bentuk layar.
func TampilTahunTreaty(t models.TahunTreaty) TahunTreatyTampil {
	return TahunTreatyTampil{
		ID: t.ID, TreatyYear: t.TreatyYear, UnderwritingYear: t.UnderwritingYear,
		TreatyGroupID: t.TreatyGroupID, TreatyGroupName: t.TreatyGroupName, Proportion: t.Proportion,
		StartDate: utils.FormatTanggal(t.StartDate), EndDate: utils.FormatTanggal(t.EndDate),
		UserID: t.UserID, TglUpdate: utils.FormatTanggalWaktu(t.TglUpdate),
	}
}

// ukuranHalamanTahunMaks meniru `pyMaxRecords` 500 (BrowseTreatyYear_RD b757).
const ukuranHalamanTahunMaks = 500

// TahunTreatyTCO melayani daftar dan penyimpanan tahun treaty.
//
// ⚠️ Berakhiran TCO sebab `services.TahunTreaty` sudah dipakai Claim Life
// (spreading.go) untuk TREATYYEAR_LIFE - dua hal berbeda, dua nama.
type TahunTreatyTCO struct {
	svc    *Service
	gudang GudangTahunTreatyTCO
	// grup - master TREATYGROUP: ID diperiksa, nama DARI master (temuan /code-review).
	grup      PembacaGrupTreatyTCO
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error
}

// TahunTreatyTCO menyusun layanannya dengan gudang yang GAGAL TERANG.
func (s *Service) TahunTreatyTCO() *TahunTreatyTCO {
	return &TahunTreatyTCO{grup: pembacaGrupTreatyBelumDisuntik{}, svc: s, gudang: gudangTahunTreatyBelumDisuntik{}, jam: time.Now,
		transaksi: s.DalamTransaksi}
}

// DenganGrup memasang pembaca master grup treaty.
func (t *TahunTreatyTCO) DenganGrup(g PembacaGrupTreatyTCO) *TahunTreatyTCO {
	salin := *t
	salin.grup = g
	return &salin
}

// DenganGudang memasang gudang.
func (t *TahunTreatyTCO) DenganGudang(g GudangTahunTreatyTCO) *TahunTreatyTCO {
	salin := *t
	salin.gudang = g
	return &salin
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (t *TahunTreatyTCO) DenganJam(j func() time.Time) *TahunTreatyTCO {
	salin := *t
	salin.jam = j
	return &salin
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji tanpa Oracle.
func (t *TahunTreatyTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *TahunTreatyTCO {
	salin := *t
	salin.transaksi = f
	return &salin
}

// Daftar membaca satu halaman, terbaru dahulu.
func (t *TahunTreatyTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, halaman, ukuran int) (HalamanTahunTreatyTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanTahunTreatyTampil{}, err
	}
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 {
		ukuran = 20
	}
	if ukuran > ukuranHalamanTahunMaks {
		ukuran = ukuranHalamanTahunMaks
	}
	hal, err := t.gudang.Daftar(ctx, halaman, ukuran)
	if err != nil {
		return HalamanTahunTreatyTampil{}, err
	}
	out := HalamanTahunTreatyTampil{Baris: make([]TahunTreatyTampil, 0, len(hal.Baris)),
		Total: hal.Total, Halaman: hal.Halaman, Ukuran: hal.Ukuran}
	for _, b := range hal.Baris {
		out.Baris = append(out.Baris, TampilTahunTreaty(b))
	}
	return out, nil
}

// Ambil membaca satu tahun treaty.
func (t *TahunTreatyTCO) Ambil(ctx context.Context, pelaku inti.Pelaku, id string) (TahunTreatyTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TahunTreatyTampil{}, err
	}
	b, err := t.gudang.Ambil(ctx, id)
	if err != nil {
		return TahunTreatyTampil{}, err
	}
	return TampilTahunTreaty(b), nil
}

// uraiTanggalMasuk membaca tanggal badan permintaan; kosong = kosong.
func uraiTanggalMasuk(nama, teks string) (time.Time, error) {
	if strings.TrimSpace(teks) == "" {
		return time.Time{}, nil
	}
	d, err := utils.ParseTanggal(strings.TrimSpace(teks))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s is not a recognised date (%q)", galat.ErrPermintaanTidakSah, nama, teks)
	}
	return d, nil
}

// Simpan menyisipkan (ID kosong) atau memperbarui (ID terisi) tahun treaty.
func (t *TahunTreatyTCO) Simpan(ctx context.Context, pelaku inti.Pelaku, masuk TahunTreatyMasuk) (TahunTreatyTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TahunTreatyTampil{}, err
	}
	mulai, err := uraiTanggalMasuk("startDate", masuk.StartDate)
	if err != nil {
		return TahunTreatyTampil{}, err
	}
	akhir, err := uraiTanggalMasuk("endDate", masuk.EndDate)
	if err != nil {
		return TahunTreatyTampil{}, err
	}
	saat := t.jam()
	model := models.TahunTreaty{
		ID:               strings.TrimSpace(masuk.ID),
		TreatyYear:       strings.TrimSpace(masuk.TreatyYear),
		UnderwritingYear: strings.TrimSpace(masuk.UnderwritingYear),
		TreatyGroupID:    strings.TrimSpace(masuk.TreatyGroupID),
		TreatyGroupName:  strings.TrimSpace(masuk.TreatyGroupName),
		Proportion:       strings.TrimSpace(masuk.Proportion),
		StartDate:        mulai,
		EndDate:          akhir,
		// b280 `UserID <- OperatorID.pyUserName`; b327 `TglUpdate <- sekarang`.
		UserID:    pelaku.AkunID,
		TglUpdate: saat,
	}
	if err := models.PeriksaTahunTreaty(model); err != nil {
		return TahunTreatyTampil{}, err
	}
	// Temuan /code-review: grup diperiksa ke master dan NAMANYA dari master -
	// nama itu disalin ke kontrak/reinsurer/business lewat kombinasi dan dibaca
	// hilir; teks bebas dari klien tidak boleh sampai ke sana.
	if model.TreatyGroupName, err = t.namaGrupDariMaster(ctx, model.TreatyGroupID); err != nil {
		return TahunTreatyTampil{}, err
	}

	err = t.transaksi(ctx, func(tx *db.Tx) error {
		idLain, err := t.gudang.CariDobel(ctx, tx, model.TreatyGroupID, model.StartDate, model.EndDate, model.ID)
		if err != nil {
			return err
		}
		if idLain != "" {
			return GalatTahunTreatyDobel{IDLain: idLain, TreatyGroupID: model.TreatyGroupID,
				StartDate: model.StartDate, EndDate: model.EndDate}
		}
		// Temuan /code-review: TREATYYEAR teks dan TREATYGROUPID ikut ditulis ke
		// kontrak (kombinasi) dan klausul - tidak dapat diganti selama tahun ini
		// sudah beranak.
		if model.ID != "" {
			lama, err := t.gudang.Ambil(repository.DenganBacaTxTCO(ctx, tx), model.ID)
			if err != nil {
				return err
			}
			if lama.TreatyYear != model.TreatyYear || lama.TreatyGroupID != model.TreatyGroupID {
				n, err := t.gudang.JumlahAnak(ctx, tx, model.ID)
				if err != nil {
					return err
				}
				if n > 0 {
					return fmt.Errorf("%w: treaty year %s has %d contracts/clauses", ErrTahunBeranak, model.ID, n)
				}
			}
		}
		if model.ID == "" {
			id, err := t.gudang.Sisip(ctx, tx, model)
			if err != nil {
				return err
			}
			model.ID = id
		} else if err := t.gudang.Perbarui(ctx, tx, model); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TahunTreatyTampil{}, err
	}
	return TampilTahunTreaty(model), nil
}

// namaGrupDariMaster mencari nama grup treaty di master `TREATYGROUP`.
func (t *TahunTreatyTCO) namaGrupDariMaster(ctx context.Context, id string) (string, error) {
	daftar, err := t.grup.Daftar(ctx)
	if err != nil {
		return "", err
	}
	if len(daftar) == 0 {
		return "", ErrMasterGrupTreatyKosong
	}
	for _, g := range daftar {
		if g.ID == id {
			return g.TreatyGroupName, nil
		}
	}
	return "", fmt.Errorf("%w: TreatyGroupID %q", ErrGrupTreatyDiLuarMaster, id)
}
