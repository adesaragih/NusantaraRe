package services

// Layar Inward Facultative tahap 2 - tiket 31, butir 78 (keputusan work owner
// 03-10-2026). Baca case, simpan blok General (tombol "Save for later"), pilihan
// Marketing Name. Submit TIDAK di sini (tahap terpisah).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// DenganKasus memasang pembaca/penulis case NB (tiket 31); transaksinya lewat DenganTransaksi.
func (s *Service) DenganKasus(k repository.PenyimpanKasus) *Service {
	s.kasus = k
	return s
}

// DenganMarketingOfficer memasang pembaca pilihan Marketing Name (tiket 31).
func (s *Service) DenganMarketingOfficer(m repository.PembacaMarketingOfficer) *Service {
	s.marketing = m
	return s
}

// DenganJam mengganti jam (uji); nil = time.Now.
func (s *Service) DenganJam(jam func() time.Time) *Service {
	s.jam = jam
	return s
}

// ErrKasusTidakAda - case NB tidak ada (atau bukan LINI Fac In). 404.
var ErrKasusTidakAda = errors.New("services: case NB tidak ditemukan")

// ErrKasusTanpaDatabase - layanan tanpa basis data. 503.
var ErrKasusTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, case NB tidak terbaca")

// ErrMasukanGeneral - isian General tidak sah (tanggal / panjang kolom). 400.
var ErrMasukanGeneral = errors.New("services: isian General tidak sah")

// ErrMarketingTanpaDatabase - tabel MARKETINGOFFICER tidak terbaca. 503.
var ErrMarketingTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel MARKETINGOFFICER tidak terbaca")

// WIB - zona tanggal layar (butir 78.1): UTC+7 tetap, tanpa musim panas.
var WIB = time.FixedZone("WIB", 7*60*60)

// jamTulisWaktu - Begin/End date ditulis pukul 12:00 WIB = 05:00 GMT (butir 78.1 diralat
// 03-10-2026): bentuk mayoritas nilai TERKINI fixture, dan delapan digit tanggal di teks =
// tanggal pilihan (00:00 WIB = 17:00 GMT hari SEBELUMNYA).
const jamTulisWaktu = 12 * time.Hour

// panjangIDKasus - lebar T_WORK_POLIS.ID (VARCHAR2(32), premiumlistlife 050): ID lebih
// panjang tidak mungkin ada = 404.
const panjangIDKasus = 32

func idKasusSah(id string) bool { return strings.TrimSpace(id) != "" && len(id) <= panjangIDKasus }

// Bentuk teks Pega yang tersimpan (butir 78.1) - `[terverifikasi]` di 5 fixture
// (`services/premium/testdata/kasus`), wadah PolicyData terkini: OfferingDate 8 digit
// (5), StartDateTime 'YYYYMMDDTHHMMSS.mmm GMT' jam 050000 x4 / 170000 x1, EndDateTime
// 050000 x3 / 170000 x2 (salinan OldData seluruhnya 170000).
const (
	BentukTanggalPega = "20060102"
	BentukWaktuPega   = "20060102T150405.000 GMT"
)

// GeneralKabel - blok General di kabel (tiket 31). Tanggal DD-MM-YYYY, kosong = "".
type GeneralKabel struct {
	ReffNumber, QQName, BeginDate, OfferingDate, EndDate, PolicyType, MarketingID, Day, TypeFacultative string
	SourceOfBusiness, CedingCoName, GroupName, OldPolicyNumber                                          string
}

// KasusKabel - jawaban GET/PUT case.
type KasusKabel struct {
	CaseID, Position, StatusWork string
	Opportunity                  IsianOpportunity
	InsuredName                  string
	General                      GeneralKabel
}

// IsianGeneral - badan PUT .../general: medan yang dapat diubah (tanpa medan tampil-saja).
type IsianGeneral struct {
	ReffNumber, QQName, BeginDate, OfferingDate, EndDate, PolicyType, MarketingID, Day, TypeFacultative string
}

func (s *Service) sekarang() time.Time {
	if s.jam != nil {
		return s.jam()
	}
	return time.Now()
}

// tanggalKeKabel - 'YYYYMMDD' -> DD-MM-YYYY. Teks lain dikembalikan APA ADANYA (A89):
// tidak dibuang diam-diam, dan simpan berikutnya menolaknya (400) alih-alih menimpanya.
func tanggalKeKabel(s string) string {
	t, err := time.Parse(BentukTanggalPega, s)
	if err != nil || t.Format(BentukTanggalPega) != s {
		return s
	}
	return t.Format(BentukTanggalKabel)
}

// waktuKeKabel - 'YYYYMMDDTHHMMSS.mmm GMT' -> tanggal WIB DD-MM-YYYY; teks lain apa adanya.
func waktuKeKabel(s string) string {
	t, err := time.Parse(BentukWaktuPega, s)
	if err != nil {
		return s
	}
	return t.In(WIB).Format(BentukTanggalKabel)
}

// uraiKabel - DD-MM-YYYY ketat (31-02 ditolak); "" = kosong.
func uraiKabel(medan, s string, masalah *[]string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation(BentukTanggalKabel, s, WIB)
	if err != nil || t.Format(BentukTanggalKabel) != s {
		*masalah = append(*masalah, medan+" bukan tanggal DD-MM-YYYY yang sah")
		return time.Time{}, false
	}
	return t, true
}

// batasGeneral - lebar kolom rancangan (BYTE) medan yang ditulis (migrasi 183).
var batasGeneral = []struct {
	nama  string
	nilai func(IsianGeneral) string
	n     int
}{
	{"reffNumber", func(i IsianGeneral) string { return i.ReffNumber }, 50},
	{"qqName", func(i IsianGeneral) string { return i.QQName }, 500},
	{"policyType", func(i IsianGeneral) string { return i.PolicyType }, 50},
	{"marketingId", func(i IsianGeneral) string { return i.MarketingID }, 50},
	{"day", func(i IsianGeneral) string { return i.Day }, 50},
	{"typeFacultative", func(i IsianGeneral) string { return i.TypeFacultative }, 50},
}

// periksaGeneral - A90: Save for later menyimpan isian SEBAGIAN - tidak ada medan wajib
// (wajib-isi Pega ditegakkan saat Submit, `[dugaan]` dari brief sesi 0f). Yang ditolak
// hanya tanggal tak sah dan isian melebihi lebar kolom. Teks lain apa adanya (butir
// 78.2/78.3: Policy Type dan Type facultative disimpan seperti dikirim layar).
// ⚠️ Kesembilan medan SELALU ditulis: medan yang tidak dikirim / "" mengosongkan kolomnya.
func periksaGeneral(i IsianGeneral) (models.General, error) {
	var masalah []string
	g := models.General{ReffNumber: i.ReffNumber, QQName: i.QQName, PolicyType: i.PolicyType,
		MarketingID: i.MarketingID, Day: i.Day, TypeFacultative: i.TypeFacultative}
	if t, ada := uraiKabel("beginDate", i.BeginDate, &masalah); ada {
		g.StartDateTime = t.Add(jamTulisWaktu).UTC().Format(BentukWaktuPega)
	}
	if t, ada := uraiKabel("offeringDate", i.OfferingDate, &masalah); ada {
		g.OfferingDate = t.Format(BentukTanggalPega)
	}
	if t, ada := uraiKabel("endDate", i.EndDate, &masalah); ada {
		g.EndDateTime = t.Add(jamTulisWaktu).UTC().Format(BentukWaktuPega)
	}
	for _, b := range batasGeneral {
		if len(b.nilai(i)) > b.n {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d byte", b.nama, b.n))
		}
	}
	if len(masalah) > 0 {
		return models.General{}, fmt.Errorf("%w: %s", ErrMasukanGeneral, strings.Join(masalah, "; "))
	}
	return g, nil
}

// keKabel - case tersimpan -> kabel. OfferingDate kosong diisi HARI INI (WIB) -
// `InwardFacultative_PreDT` langkah 5-5.1 `[terverifikasi]` (aktif): WHEN
// `pyWorkPage.OfferFacIn.PolicyData.OfferingDate==""` SET `@getCurrentDateStamp()`.
// Langkah 2-2.1 (`@DateTime.CurrentDate("dd/MM/yyyy","")`) pyDisabled true - tidak berlaku.
// A91: diterapkan saat baca, tidak ditulis sampai disimpan. Nilai bawaan sel 26 (`0`) dan
// sel 43 (`FacultativeIn`) TIDAK diterapkan - keduanya kode Pega, sedangkan butir
// 78.2/78.3 menyimpan label layar.
func (s *Service) keKabel(k models.Kasus) KasusKabel {
	o := k.Opportunity
	tutup := ""
	if !o.EstimatedClosingDate.IsZero() {
		tutup = o.EstimatedClosingDate.Format(BentukTanggalKabel)
	}
	g := k.General
	tawar := tanggalKeKabel(g.OfferingDate)
	if g.OfferingDate == "" {
		tawar = s.sekarang().In(WIB).Format(BentukTanggalKabel)
	}
	return KasusKabel{CaseID: k.CaseID, Position: k.Position, StatusWork: k.StatusWork, InsuredName: k.InsuredName,
		Opportunity: IsianOpportunity{EstimatedClosingDate: tutup, BusinessProspectName: o.BusinessProspectName,
			AccountID: o.AccountID, InsuredID: o.InsuredID, GroupBusinessID: o.GroupBusinessID,
			GroupBusiness: o.GroupBusiness, ClassOfBusiness: o.ClassOfBusiness, TypeOfInward: o.TypeOfInward,
			TypeOfFacultative: o.TypeOfFacultative, Phase: o.Phase, Stage: o.Stage,
			OpportunitySource: o.OpportunitySource, BusinessStatus: o.BusinessStatus, Description: o.Description},
		General: GeneralKabel{ReffNumber: g.ReffNumber, QQName: g.QQName, BeginDate: waktuKeKabel(g.StartDateTime),
			OfferingDate: tawar, EndDate: waktuKeKabel(g.EndDateTime), PolicyType: g.PolicyType,
			MarketingID: g.MarketingID, Day: g.Day, TypeFacultative: g.TypeFacultative,
			SourceOfBusiness: g.SourceOfBusiness, CedingCoName: g.CedingCoName, GroupName: g.GroupName,
			OldPolicyNumber: g.OldPolicyNumber}}
}

// BacaKasus - case NB `id` untuk layar Inward Facultative.
func (s *Service) BacaKasus(ctx context.Context, id string) (KasusKabel, error) {
	if s.kasus == nil {
		return KasusKabel{}, ErrKasusTanpaDatabase
	}
	if !idKasusSah(id) {
		return KasusKabel{}, ErrKasusTidakAda
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return KasusKabel{}, ErrKasusTidakAda
	}
	if err != nil {
		return KasusKabel{}, err
	}
	return s.keKabel(k), nil
}

// SimpanGeneral - tombol "Save for later": simpan blok General lalu kembalikan case.
// Urutan: identitas (401) -> isian (400) -> basis data (503) -> case (404).
func (s *Service) SimpanGeneral(ctx context.Context, pelaku inti.Pelaku, id string, isian IsianGeneral) (KasusKabel, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KasusKabel{}, err
	}
	g, err := periksaGeneral(isian)
	if err != nil {
		return KasusKabel{}, err
	}
	if s.kasus == nil || s.transaksi == nil {
		return KasusKabel{}, ErrKasusTanpaDatabase
	}
	if !idKasusSah(id) {
		return KasusKabel{}, ErrKasusTidakAda
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.kasus.SimpanGeneral(ctx, tx, id, g) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return KasusKabel{}, ErrKasusTidakAda
	}
	if err != nil {
		return KasusKabel{}, err
	}
	return s.BacaKasus(ctx, id)
}

// DaftarMarketingOfficer - pilihan Marketing Name (RD BrowseMarketingOfficer_RD).
func (s *Service) DaftarMarketingOfficer(ctx context.Context) ([]models.MarketingOfficer, error) {
	if s.marketing == nil {
		return nil, ErrMarketingTanpaDatabase
	}
	return s.marketing.DaftarMarketingOfficer(ctx)
}
