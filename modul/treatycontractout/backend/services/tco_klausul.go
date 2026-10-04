package services

// Klausul - satu tabel, 25 jenis, validasi per jenis - tiket 08 Treaty Contract Out.
//
// Untuk apa berkas ini: layanan layar `InboxTreatyContractDescription` (harness
// b359, dibuka tombol `List Description` baris tahun treaty
// `InputTreatyContract.xml` b22196): daftar jenis dari master `TREATYDESC`
// (grid `For Non XOL` / `For XOL`, `BrowseTreatyDesc_RD`), lalu panel per jenis
// (`SetKirimIDDesc` + `BrowseDescriptionLimit` + `PanggilID`).
//
// ⛔ Aturan wajib-isi ada di KODE (`models.AturanKlausulTCO`, AC 34). Jenis di
// master yang belum punya aturan di kode TIDAK dapat disimpan - tampil dengan
// alasannya, bukan diam-diam tanpa validasi (AC 36 berlaku serupa).
//
// ⛔ Rp/Usd baris ANAK dihitung SERVER dari induknya (`HitungRpUsd`); klien hanya
// mengirim Pct.
//
// Dibaca sesudah: models/tco_klausul.go, tco_reinsurer.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/repository"
)

var (
	// ErrGudangKlausulBelumDisuntik - gudang klausul belum dipasang.
	ErrGudangKlausulBelumDisuntik = errors.New("services: clause store is not injected")
	// ErrKlausulTidakAda - klausul bukan milik tahun treaty itu (404).
	ErrKlausulTidakAda = repository.ErrKlausulTidakAda
	// ErrKlausulDobel - "Data sudah pernah di Input" (409, AC 30).
	ErrKlausulDobel = errors.New("services: duplicate clause")
	// ErrJenisKlausulDiLuarMaster - DescID tidak ada di master TREATYDESC (422).
	ErrJenisKlausulDiLuarMaster = errors.New("services: clause type not in the TREATYDESC master")
	// ErrKlausulJenisBerubah - pembaruan memindah baris ke jenis/induk lain (400).
	ErrKlausulJenisBerubah = errors.New("services: an update must not move a clause to another type or parent")
	// ErrKlausulIndukBeranak - jenis reasuransi induk diubah padahal ia beranak (409).
	ErrKlausulIndukBeranak = errors.New("services: the parent reinsurance type cannot be changed while it has child rows")
	// ErrPilihanDiLuarMaster - ID occupation/clause tidak ada di master FIRE (422).
	ErrPilihanDiLuarMaster = repository.ErrPilihanMasterTidakAda
)

// GalatKlausulDobel menyebut baris yang sudah memegang isian yang sama.
type GalatKlausulDobel struct{ IDLain, Jenis string }

func (g GalatKlausulDobel) Error() string {
	return fmt.Sprintf("%s: the same %s clause already exists (row %s)", PesanKontrakDobelTCO, g.Jenis, g.IDLain)
}

// Is membuat `errors.Is(err, ErrKlausulDobel)` benar.
func (GalatKlausulDobel) Is(target error) bool { return target == ErrKlausulDobel }

// ErrKlausulSatuBaris - jenis satu baris (`models.AturanKlausul.SatuBaris`)
// sudah berisi; ubah lewat `Edit` (409).
var ErrKlausulSatuBaris = errors.New("services: this clause type holds a single row per treaty year")

// GalatKlausulSatuBaris menyebut baris yang sudah ada.
type GalatKlausulSatuBaris struct{ IDAda, Jenis string }

func (g GalatKlausulSatuBaris) Error() string {
	return fmt.Sprintf("%s only holds one row per treaty year; row %s already exists - use Edit", g.Jenis, g.IDAda)
}

// Is membuat `errors.Is(err, ErrKlausulSatuBaris)` benar.
func (GalatKlausulSatuBaris) Is(target error) bool { return target == ErrKlausulSatuBaris }

// GudangKlausulTCO membaca dan menulis klausul.
type GudangKlausulTCO interface {
	Daftar(ctx context.Context, tahunID, descID, parentReinsTypeID string) ([]models.KlausulTreaty, error)
	Ambil(ctx context.Context, tahunID, id string) (models.KlausulTreaty, error)
	Induk(ctx context.Context, tahunID, descID, reinsTypeID string) (models.KlausulTreaty, error)
	PctAnakLain(ctx context.Context, tx *db.Tx, tahunID, descID, parentReinsTypeID, kecualiID string) ([]*apd.Decimal, error)
	CariDobel(ctx context.Context, tx *db.Tx, k models.KlausulTreaty, kunci []string) (string, error)
	Sisip(ctx context.Context, tx *db.Tx, k models.KlausulTreaty) (string, error)
	Perbarui(ctx context.Context, tx *db.Tx, k models.KlausulTreaty) error
}

// PembacaMasterKlausulTCO membaca master TREATYDESC / OCCUPATION / CLAUSE.
type PembacaMasterKlausulTCO interface {
	JenisKlausul(ctx context.Context, isXOL string) ([]repository.JenisKlausulMasterTCO, error)
	CariPilihan(ctx context.Context, master, teks string) ([]repository.PilihanMasterTCO, error)
	AmbilPilihan(ctx context.Context, master, id string) (repository.PilihanMasterTCO, error)
}

// PengunciTahunTCO membaca dan mengunci tahun treaty induk klausul.
type PengunciTahunTCO interface {
	Ambil(ctx context.Context, id string) (models.TahunTreaty, error)
	Kunci(ctx context.Context, tx *db.Tx, id string) error
}

type klausulBelumDisuntik struct{}

func (klausulBelumDisuntik) Daftar(context.Context, string, string, string) ([]models.KlausulTreaty, error) {
	return nil, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) Ambil(context.Context, string, string) (models.KlausulTreaty, error) {
	return models.KlausulTreaty{}, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) Induk(context.Context, string, string, string) (models.KlausulTreaty, error) {
	return models.KlausulTreaty{}, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) PctAnakLain(context.Context, *db.Tx, string, string, string, string) ([]*apd.Decimal, error) {
	return nil, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) CariDobel(context.Context, *db.Tx, models.KlausulTreaty, []string) (string, error) {
	return "", ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) Sisip(context.Context, *db.Tx, models.KlausulTreaty) (string, error) {
	return "", ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) Perbarui(context.Context, *db.Tx, models.KlausulTreaty) error {
	return ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) JenisKlausul(context.Context, string) ([]repository.JenisKlausulMasterTCO, error) {
	return nil, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) CariPilihan(context.Context, string, string) ([]repository.PilihanMasterTCO, error) {
	return nil, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) AmbilPilihan(context.Context, string, string) (repository.PilihanMasterTCO, error) {
	return repository.PilihanMasterTCO{}, ErrGudangKlausulBelumDisuntik
}
func (klausulBelumDisuntik) Kunci(context.Context, *db.Tx, string) error {
	return ErrGudangKlausulBelumDisuntik
}

type pengunciTahunBelumDisuntik struct{ klausulBelumDisuntik }

func (pengunciTahunBelumDisuntik) Ambil(context.Context, string) (models.TahunTreaty, error) {
	return models.TahunTreaty{}, ErrGudangKlausulBelumDisuntik
}

type gudangKlausulOracle struct {
	*repository.MasterKlausulTCO
	db *db.DB
}

// GudangKlausulOracle menyusun gudang klausul di atas Oracle.
func GudangKlausulOracle(svc *Service) GudangKlausulTCO {
	return gudangKlausulOracle{MasterKlausulTCO: repository.NewMasterKlausulTCO(svc.DB()), db: svc.DB()}
}

// MasterKlausulOracle menyusun pembaca master klausul.
func MasterKlausulOracle(svc *Service) PembacaMasterKlausulTCO {
	return repository.NewMasterKlausulPilihan(svc.DB())
}

// PengunciTahunOracle menyusun pembaca + pengunci tahun treaty.
func PengunciTahunOracle(svc *Service) PengunciTahunTCO {
	return repository.NewMasterTahunTreaty(svc.DB())
}

// AturanTampil adalah satu aturan jenis seperti dikirim ke layar - supaya form
// per jenis dirakit dari SATU tabel kebenaran, bukan ditulis ulang di layar.
type AturanTampil struct {
	Jenis    string   `json:"jenis"`
	Anak     bool     `json:"anak"`
	Subjenis string   `json:"subjenis"`
	Medan    []string `json:"medan"`
	Wajib    []string `json:"wajib"`
	Turunan  []string `json:"turunan"`
	Ditahan  string   `json:"ditahan"`
	// Berkurs / Konversi - tiket 11.
	Berkurs  bool   `json:"berkurs"`
	Konversi string `json:"konversi"`
	Sumber   string `json:"sumber"`
	// PilihanReins - `models.AturanKlausul.PilihanReins`: layar memilih
	// pemilih ReinsType dari sini, bukan dari nama jenis.
	PilihanReins string `json:"pilihanReins"`
	// SatuBaris - `Add` hilang begitu jenis ini berisi satu baris.
	SatuBaris bool `json:"satuBaris"`
}

// JenisKlausulTampil adalah satu baris grid jenis + aturannya.
type JenisKlausulTampil struct {
	ID          string         `json:"id"`
	DescName    string         `json:"descName"`
	IsXOL       string         `json:"isXol"`
	StatusAktif string         `json:"statusAktif"`
	Aturan      []AturanTampil `json:"aturan"`
	// Catatan - jenis master tanpa aturan di kode.
	Catatan string `json:"catatan"`
}

// KlausulMasuk adalah badan simpan satu klausul.
type KlausulMasuk struct {
	ID       string `json:"id"`
	DescID   string `json:"descId"`
	Anak     bool   `json:"anak"`
	Subjenis string `json:"subjenis"`
	// ParentReinsTypeID - jenis reasuransi induk; hanya untuk baris anak.
	ParentReinsTypeID string            `json:"parentReinsTypeId"`
	Medan             map[string]string `json:"medan"`
}

// KlausulTampil adalah satu klausul seperti dikirim ke layar.
//
// ⛔ Seluruh uang dan persen TEKS (ADR-0003).
type KlausulTampil struct {
	ID                string            `json:"id"`
	TreatyYear        string            `json:"treatyYear"`
	TreatyYearID      string            `json:"treatyYearId"`
	TreatyGroupID     string            `json:"treatyGroupId"`
	TreatyDescID      string            `json:"treatyDescId"`
	TreatyDescName    string            `json:"treatyDescName"`
	ReinsTypeID       string            `json:"reinsTypeId"`
	ReinsTypeName     string            `json:"reinsTypeName"`
	ParentReinsTypeID string            `json:"parentReinsTypeId"`
	Subjenis          string            `json:"subjenis"`
	Medan             map[string]string `json:"medan"`
	Kurs              string            `json:"kurs"`
	UserID            string            `json:"userId"`
	TglUpdate         string            `json:"tglUpdate"`
}

// TampilKlausul menerjemahkan satu baris - medan jenisnya saja, VERBATIM.
func TampilKlausul(k models.KlausulTreaty) KlausulTampil {
	medan := map[string]string{}
	for _, m := range []string{models.MedanReinsTypeID, models.MedanLine, models.MedanRp, models.MedanUsd,
		models.MedanPct, models.MedanPctMe, models.MedanYdcf, models.MedanMethod, models.MedanTerritorialLimit,
		models.MedanCoInsMin, models.MedanCoInsMax, models.MedanTreatyLimit, models.MedanIDOccupation,
		models.MedanOccupation, models.MedanIDClause, models.MedanClause, models.MedanLayer, models.MedanMoreRp,
		models.MedanMoreUsd} {
		if v := models.NilaiMedanKlausul(k, m); v != "" {
			medan[m] = v
		}
	}
	sub := models.SubjenisKlausulTCO(k)
	return KlausulTampil{ID: k.ID, TreatyYear: k.TreatyYear, TreatyYearID: k.TreatyYearID,
		TreatyGroupID: k.TreatyGroupID, TreatyDescID: k.TreatyDescID, TreatyDescName: k.TreatyDescName,
		ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName, ParentReinsTypeID: k.ParentReinsTypeID,
		Subjenis: sub, Medan: medan, Kurs: utils.FormatDecimal(k.Kurs), UserID: k.UserID,
		TglUpdate: utils.FormatTanggalWaktu(k.TglUpdate)}
}

// DaftarKlausulTampil adalah grid satu jenis (induk, atau anak satu induk).
type DaftarKlausulTampil struct {
	Daftar []KlausulTampil `json:"daftar"`
	Total  int             `json:"total"`
	// TotalPct dan Peringatan - hanya untuk daftar anak.
	TotalPct   string `json:"totalPct"`
	Peringatan string `json:"peringatan"`
}

// HasilKlausulTampil adalah jawaban simpan.
type HasilKlausulTampil struct {
	Klausul    KlausulTampil `json:"klausul"`
	TotalPct   string        `json:"totalPct"`
	Peringatan string        `json:"peringatan"`
}

// PilihanTampil adalah satu pilihan occupation/clause.
type PilihanTampil struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// KlausulTCO melayani layar klausul.
type KlausulTCO struct {
	svc       *Service
	gudang    GudangKlausulTCO
	master    PembacaMasterKlausulTCO
	tahun     PengunciTahunTCO
	jenis     PembacaJenisReasuransiTCO
	kurs      PembacaKursTCO
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error
}

// KlausulTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) KlausulTCO() *KlausulTCO {
	return &KlausulTCO{svc: s, gudang: klausulBelumDisuntik{}, master: klausulBelumDisuntik{},
		tahun: pengunciTahunBelumDisuntik{}, jenis: pembacaJenisReasuransiBelumDisuntik{}, kurs: kursBelumDisuntik{},
		jam: time.Now, transaksi: s.DalamTransaksi}
}

func (l *KlausulTCO) salin() *KlausulTCO { s := *l; return &s }

// DenganGudang memasang gudang klausul.
func (l *KlausulTCO) DenganGudang(g GudangKlausulTCO) *KlausulTCO {
	s := l.salin()
	s.gudang = g
	return s
}

// DenganMaster memasang pembaca master klausul.
func (l *KlausulTCO) DenganMaster(m PembacaMasterKlausulTCO) *KlausulTCO {
	s := l.salin()
	s.master = m
	return s
}

// DenganTahun memasang pembaca + pengunci tahun treaty.
func (l *KlausulTCO) DenganTahun(t PengunciTahunTCO) *KlausulTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganJenis memasang daftar jenis reasuransi tersaring (tiket 02).
func (l *KlausulTCO) DenganJenis(j PembacaJenisReasuransiTCO) *KlausulTCO {
	s := l.salin()
	s.jenis = j
	return s
}

// DenganKurs memasang pembaca kurs berlaku (tiket 11).
func (l *KlausulTCO) DenganKurs(k PembacaKursTCO) *KlausulTCO {
	s := l.salin()
	s.kurs = k
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *KlausulTCO) DenganJam(j func() time.Time) *KlausulTCO { s := l.salin(); s.jam = j; return s }

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *KlausulTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *KlausulTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

// larik - nil menjadi larik kosong: layar membaca `medan`/`wajib` sebagai larik
// (`aturan.wajib.includes`), dan jenis tanpa wajib-isi (LimitMB) mengirim null.
func larik(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func tampilAturan(a models.AturanKlausul) AturanTampil {
	return AturanTampil{Jenis: a.Jenis, Anak: a.Anak, Subjenis: a.Subjenis, Medan: larik(a.Medan), Wajib: larik(a.Wajib),
		Turunan: a.Turunan, Ditahan: a.Ditahan, Berkurs: a.Berkurs, Konversi: a.Konversi, Sumber: a.Sumber,
		PilihanReins: a.PilihanReins, SatuBaris: a.SatuBaris}
}

// JenisKlausul membaca grid jenis (`BrowseTreatyDesc_RD`) + aturan tiap jenis.
func (l *KlausulTCO) JenisKlausul(ctx context.Context, pelaku inti.Pelaku, isXOL string) ([]JenisKlausulTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	d, err := l.master.JenisKlausul(ctx, strings.TrimSpace(isXOL))
	if err != nil {
		return nil, err
	}
	hasil := make([]JenisKlausulTampil, 0, len(d))
	for _, j := range d {
		t := JenisKlausulTampil{ID: j.ID, DescName: j.DescName, IsXOL: j.IsXOL, StatusAktif: j.StatusAktif,
			Aturan: []AturanTampil{}}
		for _, a := range models.AturanKlausulTCO {
			if a.DescID == j.ID {
				t.Aturan = append(t.Aturan, tampilAturan(a))
			}
		}
		if len(t.Aturan) == 0 {
			// Jenis master tanpa aturan di kode - kalimat untuk pemakai (30-09-2026).
			t.Catatan = "this type cannot be saved yet: its required-field rules are not available yet"
		}
		hasil = append(hasil, t)
	}
	return hasil, nil
}

// Daftar membaca klausul satu jenis: induk (`parent` kosong/"00") atau anak.
func (l *KlausulTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, tahunID, descID, parent string) (DaftarKlausulTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return DaftarKlausulTampil{}, err
	}
	if strings.TrimSpace(parent) == "" {
		parent = models.ParentReinsTypeTanpaInduk
	}
	if _, err := l.tahun.Ambil(ctx, tahunID); err != nil {
		return DaftarKlausulTampil{}, err
	}
	baris, err := l.gudang.Daftar(ctx, tahunID, descID, parent)
	if err != nil {
		return DaftarKlausulTampil{}, err
	}
	if strings.TrimSpace(descID) == models.DescCoinsPanel {
		models.UrutCoInsScaleTCO(baris)
	}
	hasil := DaftarKlausulTampil{Daftar: make([]KlausulTampil, 0, len(baris)), Total: len(baris)}
	var pct []*apd.Decimal
	for _, b := range baris {
		hasil.Daftar = append(hasil.Daftar, TampilKlausul(b))
		pct = append(pct, b.Pct)
	}
	if parent != models.ParentReinsTypeTanpaInduk {
		total, err := models.TotalShareTCO(pct)
		if err != nil {
			return DaftarKlausulTampil{}, err
		}
		hasil.TotalPct = total.Text('f')
		if a, ok := models.CariAturanKlausul(descID, true, ""); ok && a.PeringatanSpreading && len(baris) > 0 {
			hasil.Peringatan = models.PeringatanSpreadingTCO(total)
		}
	}
	return hasil, nil
}

// namaReinsType memeriksa ID jenis reasuransi di daftar pilihan jenisnya:
// induk - daftar tersaring tiket 02 [keputusan work owner 29-09-2026]
// (OQ-TCO-15, ditutup); SETIAP anak - `TreatyContractSetReinsTypeList` atas
// nama ReinsType baris induknya (`namaInduk`), nama tersimpan = `.CARI2`
// [keputusan work owner 02-10-2026].
func (l *KlausulTCO) namaReinsType(ctx context.Context, a models.AturanKlausul, id, namaInduk string) (string, error) {
	if a.PilihanReins == models.PilihanReinsAnakTreatyLimit {
		for _, p := range models.PilihanReinsAnakDari(namaInduk) {
			if p.ID == id {
				return p.Nama, nil
			}
		}
		return "", fmt.Errorf("%w: ReinsTypeID %q is not a child option of parent %q (TreatyContractSetReinsTypeList)",
			ErrJenisReasuransiDiLuarDaftar, id, namaInduk)
	}
	daftar, err := l.jenis.DaftarNonLife(ctx)
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

// namaDesc mencari `DescName` di master jenis.
func (l *KlausulTCO) namaDesc(ctx context.Context, descID string) (string, error) {
	d, err := l.master.JenisKlausul(ctx, "")
	if err != nil {
		return "", err
	}
	for _, j := range d {
		if j.ID == descID {
			return j.DescName, nil
		}
	}
	return "", fmt.Errorf("%w: TreatyDescID %q", ErrJenisKlausulDiLuarMaster, descID)
}

// Simpan menulis satu klausul - `Save` tiap form jenis (`SaveTreatyArr*`).
func (l *KlausulTCO) Simpan(ctx context.Context, pelaku inti.Pelaku, tahunID string, m KlausulMasuk) (HasilKlausulTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilKlausulTampil{}, err
	}
	tahun, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return HasilKlausulTampil{}, err
	}
	a, ok := models.CariAturanKlausul(m.DescID, m.Anak, m.Subjenis)
	if !ok {
		return HasilKlausulTampil{}, fmt.Errorf("%w: TreatyDescID %q child %v subtype %q", models.ErrKlausulJenisTakDikenal,
			m.DescID, m.Anak, m.Subjenis)
	}
	if a.Ditahan != "" {
		return HasilKlausulTampil{}, models.PeriksaKlausulTCO(a, models.KlausulTreaty{})
	}
	descName, err := l.namaDesc(ctx, a.DescID)
	if err != nil {
		return HasilKlausulTampil{}, err
	}
	parent := models.ParentReinsTypeTanpaInduk
	if a.Anak {
		parent = strings.TrimSpace(m.ParentReinsTypeID)
		if parent == "" || parent == models.ParentReinsTypeTanpaInduk {
			return HasilKlausulTampil{}, fmt.Errorf("%w: a child row must name its parent ParentReinsTypeID", galat.ErrPermintaanTidakSah)
		}
	}
	k := models.KlausulTreaty{}
	id := strings.TrimSpace(m.ID)
	var lama models.KlausulTreaty
	if id != "" {
		if lama, err = l.gudang.Ambil(ctx, tahunID, id); err != nil {
			return HasilKlausulTampil{}, err
		}
		if lama.TreatyDescID != a.DescID || lama.ParentReinsTypeID != parent {
			return HasilKlausulTampil{}, fmt.Errorf("%w: row %s", ErrKlausulJenisBerubah, id)
		}
		// Co-Ins Scale: baris tidak pindah grid (`SetTreatyArrExclustionCoins_Act`
		// membaca `.SpreadingOrder` barisnya sendiri, b721/b1049).
		if a.DescID == models.DescCoinsPanel && models.SubjenisKlausulTCO(lama) != a.Subjenis {
			return HasilKlausulTampil{}, fmt.Errorf("%w: row %s belongs to grid %q", ErrKlausulJenisBerubah, id, lama.SpreadingOrder)
		}
		// Medan di luar form (Kurs, Layer*, SpreadingOrder, ...) dipertahankan.
		k = lama
	}
	if a.DescID == models.DescCoinsPanel {
		// `NewTreatyArrCoins` b405/b406: `.SpreadingOrder = Param.Type` grid tombol Add-nya.
		k.SpreadingOrder = a.Subjenis
	}
	k.ID, k.TreatyYear, k.TreatyYearID = id, tahun.TreatyYear, tahunID
	k.TreatyGroupID, k.TreatyGroupName = tahun.TreatyGroupID, tahun.TreatyGroupName
	k.TreatyDescID, k.TreatyDescName, k.ParentReinsTypeID = a.DescID, descName, parent
	k.UserID, k.TglUpdate = pelaku.AkunID, l.jam()
	for _, turunan := range a.Turunan {
		if _, dikirim := m.Medan[turunan]; dikirim {
			return HasilKlausulTampil{}, fmt.Errorf("%w: %s is computed by the server, not sent by the client",
				models.ErrMedanBukanMilikJenis, turunan)
		}
	}
	if err := models.IsiMedanKlausulTCO(a, &k, m.Medan); err != nil {
		return HasilKlausulTampil{}, err
	}
	// Tiket 11: form berkurs menuntut kurs berlaku (`NewTreatyArr*`); induk
	// ber-Rp/Usd menurunkan `Usd = Rp / Kurs` (`HitungRpUsd_depan`) dan
	// menyimpan kurs yang dipakai di `KURS` [keputusan work owner 29-09-2026] (OQ-TCO-18).
	if a.Berkurs {
		kurs, err := l.kurs.Berlaku(ctx, tahun)
		if err != nil {
			return HasilKlausulTampil{}, err
		}
		if a.Konversi == models.KonversiRpKeUsd {
			k.Kurs = kurs.ToIDR
			if k.Usd, err = models.UsdDariRpTCO(k.Rp, kurs.ToIDR, models.SkalaUsdDariRpTCO); err != nil {
				return HasilKlausulTampil{}, err
			}
		}
	}
	// Baris induk lebih dulu: pilihan ReinsType anak dihitung dari NAMANYA.
	var induk models.KlausulTreaty
	if a.Anak {
		if induk, err = l.gudang.Induk(ctx, tahunID, a.DescID, parent); err != nil {
			return HasilKlausulTampil{}, err
		}
	}
	if err := l.lengkapiDariMaster(ctx, a, &k, lama, induk.ReinsTypeName); err != nil {
		return HasilKlausulTampil{}, err
	}
	if a.Anak {
		if k.Rp, k.Usd, err = models.RpUsdAnakTCO(k.Pct, induk.Rp, induk.Usd); err != nil {
			return HasilKlausulTampil{}, err
		}
	}
	if err := models.PeriksaKlausulTCO(a, k); err != nil {
		return HasilKlausulTampil{}, err
	}
	var total *apd.Decimal
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.tahun.Kunci(ctx, tx, tahunID); err != nil {
			return err
		}
		lain, err := l.gudang.CariDobel(ctx, tx, k, a.KunciDobel)
		if err != nil {
			return err
		}
		if lain != "" {
			return GalatKlausulDobel{IDLain: lain, Jenis: a.Jenis}
		}
		if a.SatuBaris && k.ID == "" {
			ada, err := l.gudang.Daftar(repository.DenganBacaTxTCO(ctx, tx), tahunID, a.DescID, parent)
			if err != nil {
				return err
			}
			if len(ada) > 0 {
				return GalatKlausulSatuBaris{IDAda: ada[0].ID, Jenis: a.Jenis}
			}
		}
		if a.BatasTotalAnak {
			pctLain, err := l.gudang.PctAnakLain(ctx, tx, tahunID, a.DescID, parent, k.ID)
			if err != nil {
				return err
			}
			if total, err = models.PeriksaTotalAnakTCO(pctLain, k.Pct); err != nil {
				return err
			}
		}
		if !a.Anak && k.ID != "" && lama.ReinsTypeID != "" && lama.ReinsTypeID != k.ReinsTypeID {
			anak, err := l.gudang.Daftar(ctx, tahunID, a.DescID, lama.ReinsTypeID)
			if err != nil {
				return err
			}
			if len(anak) > 0 {
				return fmt.Errorf("%w: parent %s has %d child rows", ErrKlausulIndukBeranak, k.ID, len(anak))
			}
		}
		if k.ID == "" {
			if k.ID, err = l.gudang.Sisip(ctx, tx, k); err != nil {
				return err
			}
		} else if err := l.gudang.Perbarui(ctx, tx, k); err != nil {
			return err
		} else if err := l.hitungUlangAnak(repository.DenganBacaTxTCO(ctx, tx), tx, pelaku, a, lama, k); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return HasilKlausulTampil{}, err
	}
	h := HasilKlausulTampil{Klausul: TampilKlausul(k)}
	if total != nil {
		h.TotalPct = total.Text('f')
		if a.PeringatanSpreading {
			h.Peringatan = models.PeringatanSpreadingTCO(total)
		}
	}
	return h, nil
}

// lengkapiDariMaster mengisi nama dari master: jenis reasuransi (tiket 02),
// occupation, clause - klien hanya memilih ID.
//
// `lama` - baris sebelum Edit (kosong untuk baris baru). ReinsType baris anak
// lama yang TIDAK diganti tetap boleh walau di luar pilihan anak sekarang
// (data lama ber-ReinsType induk, mis. ORS di bawah ORS) - Edit Pct-nya tidak
// terkunci; menggantinya wajib memakai pilihan anak `TreatyContractSetReinsTypeList`.
func (l *KlausulTCO) lengkapiDariMaster(ctx context.Context, a models.AturanKlausul, k *models.KlausulTreaty,
	lama models.KlausulTreaty, namaInduk string) error {
	punya := func(m string) bool {
		for _, x := range a.Medan {
			if x == m {
				return true
			}
		}
		return false
	}
	var err error
	if punya(models.MedanReinsTypeID) && k.ReinsTypeID != "" {
		if a.Anak && lama.ID != "" && lama.ReinsTypeID == k.ReinsTypeID {
			k.ReinsTypeName = lama.ReinsTypeName
		} else if k.ReinsTypeName, err = l.namaReinsType(ctx, a, k.ReinsTypeID, namaInduk); err != nil {
			return err
		}
	}
	if a.DescID == models.DescLimitMB {
		// MB Capacity: nama dari `SetOccupationLimitMB` (empat ID tetap), bukan master FIRE.
		k.Occupation = ""
		if k.IDOccupation != "" {
			nama, ok := models.NamaOccupationLimitMB(k.IDOccupation)
			if !ok {
				return fmt.Errorf("%w: ID_Occupation %q is not an MB Capacity occupation (SetOccupationLimitMB)",
					ErrPilihanDiLuarMaster, k.IDOccupation)
			}
			k.Occupation = nama
		}
	} else if punya(models.MedanIDOccupation) && k.IDOccupation != "" {
		p, err := l.master.AmbilPilihan(ctx, repository.MasterOccupationTCO, k.IDOccupation)
		if err != nil {
			return err
		}
		k.Occupation = p.Nama
	}
	if punya(models.MedanIDClause) && k.IDClause != "" {
		p, err := l.master.AmbilPilihan(ctx, repository.MasterClauseTCO, k.IDClause)
		if err != nil {
			return err
		}
		k.Clause = p.Nama
	}
	return nil
}

// Pilihan membaca pemilih ExclutionTreaty - `occupation` atau `clause`.
func (l *KlausulTCO) Pilihan(ctx context.Context, pelaku inti.Pelaku, master, cari string) ([]PilihanTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	var tabel string
	switch master {
	case "occupation-limitmb":
		// MB Capacity - empat pilihan `SetOccupationLimitMB`, tanpa baca master.
		hasil := make([]PilihanTampil, 0, len(models.PilihanOccupationLimitMB))
		for _, p := range models.PilihanOccupationLimitMB {
			hasil = append(hasil, PilihanTampil{ID: p.ID, Nama: p.Nama})
		}
		return hasil, nil
	case "occupation":
		tabel = repository.MasterOccupationTCO
	case "clause":
		tabel = repository.MasterClauseTCO
	default:
		return nil, fmt.Errorf("%w: option master %q", galat.ErrPermintaanTidakSah, master)
	}
	d, err := l.master.CariPilihan(ctx, tabel, cari)
	if err != nil {
		return nil, err
	}
	hasil := make([]PilihanTampil, 0, len(d))
	for _, p := range d {
		hasil = append(hasil, PilihanTampil{ID: p.ID, Nama: p.Nama})
	}
	return hasil, nil
}

// hitungUlangAnak - temuan /code-review: Rp/Usd anak adalah TURUNAN induknya
// (`HitungRpUsd`), yang tersimpan. Induk yang Rp/Usd-nya berubah menghitung
// ulang seluruh anaknya di transaksi yang sama, supaya nilai tersimpan tetap
// Pct x induk / 100 (Pega hanya menghitung saat anak disimpan).
func (l *KlausulTCO) hitungUlangAnak(ctx context.Context, tx *db.Tx, pelaku inti.Pelaku, a models.AturanKlausul,
	lama, induk models.KlausulTreaty) error {

	if a.Anak || induk.ReinsTypeID == "" {
		return nil
	}
	_, ada := models.CariAturanKlausul(a.DescID, true, "")
	if !ada || (desimalSama(lama.Rp, induk.Rp) && desimalSama(lama.Usd, induk.Usd)) {
		return nil
	}
	anak, err := l.gudang.Daftar(ctx, induk.TreatyYearID, a.DescID, induk.ReinsTypeID)
	if err != nil {
		return err
	}
	for _, c := range anak {
		rp, usd, err := models.RpUsdAnakTCO(c.Pct, induk.Rp, induk.Usd)
		if err != nil {
			return err
		}
		if desimalSama(c.Rp, rp) && desimalSama(c.Usd, usd) {
			continue
		}
		c.Rp, c.Usd, c.UserID, c.TglUpdate = rp, usd, pelaku.AkunID, induk.TglUpdate
		if err := l.gudang.Perbarui(ctx, tx, c); err != nil {
			return err
		}
	}
	return nil
}

// desimalSama - dua desimal (boleh nil) bernilai sama.
func desimalSama(a, b *apd.Decimal) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Cmp(b) == 0
}
