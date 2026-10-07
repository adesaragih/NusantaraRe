package services

// Proteksi Accept tahap marketing (tiket 57) - port `Activity\Protection_Act.xml` untuk kasus FIRE NB di tahap Offer
// (PositionNote ReasFacInMarketing, FlagOnGoingPolicy 0, EmailType 1). Pesan VERBATIM, urutan = urutan langkah Pega,
// pesan kembar ditampilkan sekali. Langkah ber-label `//` (nonaktif) tidak diport. Tabel lengkap - pesan, rule, terpenuhi /
// belum, alasan - di docs/issues/57.
//
// Yang BELUM dapat diperiksa karena datanya tidak ada di sistem baru (dilaporkan, bukan dilewati diam-diam - keputusan
// work owner 07-10-2026 butir 3): IsGroup, IsB2B, IsMOP, StatusBusiness (EDM), IsSurveyReport / SurveyReportList,
// BusinessCode (Double Input Case), RW.TERITORYNAME (kolom tidak ada di DDL RW). Teks `ASMMessageQuotation` (field value
// tidak ada di korpus) = teks keputusan work owner. Kasus diperlakukan: bukan Group, bukan B2B ASM, bukan MOP, bukan EDM.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// pesanMOIDKosong - Protection_Act langkah 17: Page-Set-Messages `ASMMessageQuotation` (Rule-Obj-FieldValue
	// pyMessageLabel; teksnya tidak ada di korpus). Teks = keputusan work owner 07-10-2026 (BUKAN teks Pega).
	pesanMOIDKosong = "Harus memilih dari daftar nama marketing yang tersedia"
	// pesanScoringKosong - Protection_Act langkah 4 `Local.errScoringRisk`, ditampilkan langkah 18.
	pesanScoringKosong = "Scoring Risk Cannot Empty !!"
	// batasLokasiProteksi - langkah 2: LocationList > 400 -> keluar tanpa proteksi.
	batasLokasiProteksi = 400
	// kasusTanpaLossRecord - ProtectFIREMBUPA_Act 9.2.7 / 9.2.10 / 9.2.30.6.4: pengecualian kasus lama apa adanya.
	kasusTanpaLossRecord = "NB-83801"
)

// insuredTanpaProteksi - langkah 2: InsuredID "ASM-SFAGIS-WORK-ORG ORG-7889" / "ORG-2344" -> keluar tanpa proteksi.
var insuredTanpaProteksi = []string{"ORG-7889", "ORG-2344"}

// rateCoverageNetRate - ProtectFIREMBUPA_Act 9.2.30.6.8-12: item ber-TotalNetRate > 0 wajib memuat kelima OLDID ini.
var rateCoverageNetRate = []string{"FLEXAS", "4.1A CC", "4.3", "4.2 PRGBI", "OTHERS"}

// pengumpulPesan - Page-Set-Messages: urutan munculnya, kembar sekali.
type pengumpulPesan struct {
	pesan []string
	lihat map[string]bool
}

func (p *pengumpulPesan) tambah(s string) {
	if s == "" || p.lihat[s] {
		return
	}
	if p.lihat == nil {
		p.lihat = map[string]bool{}
	}
	p.lihat[s] = true
	p.pesan = append(p.pesan, s)
}

// insuredDikecualikan - InsuredID utuh Pega atau kode ORG-nya saja (bentuk INSURED_ID sistem baru belum terverifikasi).
func insuredDikecualikan(id string) bool {
	id = strings.TrimSpace(id)
	for _, x := range insuredTanpaProteksi {
		if id == x || strings.HasSuffix(id, " "+x) {
			return true
		}
	}
	return false
}

// ProteksiAccept - pesan proteksi Accept; kosong = lolos.
func (s *Service) ProteksiAccept(ctx context.Context, id string) ([]string, error) {
	if s.kasus == nil || s.submit == nil {
		return nil, ErrSubmitTanpaDatabase
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	objek, err := s.BacaObjek(ctx, id)
	if err != nil {
		return nil, err
	}
	// Langkah 2: keluar tanpa proteksi.
	if insuredDikecualikan(kasus.Opportunity.InsuredID) || len(objek) > batasLokasiProteksi {
		return []string{}, nil
	}
	var p pengumpulPesan
	// Langkah 16-17 MOID kosong; langkah 18 Scoring Risk (bukan EDM, bukan B2B ASM, IsFire).
	if strings.TrimSpace(kasus.General.MarketingID) == "" {
		p.tambah(pesanMOIDKosong)
	}
	if kasus.RiskScoring.FinalScore == nil {
		p.tambah(pesanScoringKosong)
	}
	// Langkah 21 ProtectCurrencyTSI_Act.
	if err := s.proteksiMataUang(ctx, objek, &p); err != nil {
		return nil, err
	}
	// Langkah 22 CheckSpreadingProtect_ACT: langkah 2 keluar untuk PositionNote ReasFacInMarketing kecuali Group - tidak ada
	// pesan di tahap ini.
	// Langkah 23 ProtectPremiPolicy_Act.
	if err := s.proteksiPayment(ctx, id, &p); err != nil {
		return nil, err
	}
	// Langkah 24 ProtectCoverage_Act -> ProtectFIREMBUPA_Act (IsFire).
	if err := s.proteksiLokasi(ctx, kasus.CaseID, objek, &p); err != nil {
		return nil, err
	}
	// Langkah 28 ProtectShareCedant_Act (keluar bila bukan Group), 29 ProtectRenewal_Act (bukan renewal), 30
	// ProtectCedingCo (IsB2B ASM), 31-34 Survey report / Ceding Co (FlagOnGoingPolicy 2 / Group / data survey) - lihat
	// komentar berkas.
	if p.pesan == nil {
		return []string{}, nil
	}
	return p.pesan, nil
}

// proteksiMataUang - ProtectCurrencyTSI_Act: `CURRENCY.CURRENCY` ber-id `.CurrencyID` wajib = `.Currency`. Sistem baru hanya
// menyimpan nama mata uang item, jadi padanannya: nama itu ada di CURRENCY `[dugaan]`. Hanya pesan item gagal TERAKHIR
// yang tampil (Local.Err ditimpa di dalam loop, ditampilkan sekali sesudahnya).
func (s *Service) proteksiMataUang(ctx context.Context, objek []models.ObjekFire, p *pengumpulPesan) error {
	ada := false
	for _, o := range objek {
		ada = ada || len(o.Items) > 0
	}
	if !ada {
		return nil
	}
	if s.mataUang == nil {
		return ErrSubmitTanpaDatabase
	}
	daftar, err := s.mataUang.DaftarMataUang(ctx)
	if err != nil {
		return err
	}
	sah := map[string]bool{}
	for _, c := range daftar {
		sah[c] = true
	}
	terakhir := ""
	for n, o := range objek {
		for m, it := range o.Items {
			// Currency kosong: CURRENCY(id kosong) tak berbaris = "" = .Currency -> Pega diam `[dugaan]`; kosong dilaporkan
			// 9.2.30.4 "must be choose!".
			if it.Currency != "" && !sah[it.Currency] {
				terakhir = fmt.Sprintf("Currency Pada Lokasi %d Object Item %d Bermasalah, Coba Pilih Ulang Currency", n+1, m+1)
			}
		}
	}
	p.tambah(terakhir)
	return nil
}

// kosongAtauNol4 - `== "" || == "0.0000"` atas @divide(x,1,4).
func kosongAtauNol4(d *apd.Decimal) bool { return d == nil || bulat4(d).IsZero() }

// beda / dekat - `a <> b` dan |a - b| < 0.001 atas nilai 4 desimal (bulat4 / selisihKecil tab Payment).
func beda(a, b *apd.Decimal) bool { return bulat4(nolBila(a)).Cmp(bulat4(nolBila(b))) != 0 }

func dekat(a, b *apd.Decimal) bool { return selisihKecil(bulat4(nolBila(a)), bulat4(nolBila(b))) }

// proteksiPayment - ProtectPremiPolicy_Act langkah 3-11 per `.OfferFacIn.CurrencyList` (mata uang tab Payment). Bukan Life,
// bukan MOP, bukan EDM, PolicyType bukan 1 `[dugaan]` (PolicyType tersimpan sebagai label layar, butir 78.2).
func (s *Service) proteksiPayment(ctx context.Context, id string, p *pengumpulPesan) error {
	t, err := s.TampilanPayment(ctx, id)
	if err != nil {
		return err
	}
	if len(t.MataUang) == 0 {
		p.tambah("List Payment can't be null!")
	}
	premiKosong := false
	totalPremiAll := nol()
	errIns := ""
	for _, mu := range t.MataUang {
		nilaiPremi := mu.Baris.Premium
		if kosongAtauNol4(nilaiPremi) {
			premiKosong = true
		}
		var premiIns, commIns, totalIns *apd.Decimal
		for _, c := range mu.Cicilan {
			premiIns, commIns, totalIns = tambah(premiIns, c.Premium), tambah(commIns, c.RICommision), tambah(totalIns, c.PaymentTotal)
			if kosongAtauNol4(c.Premium) {
				premiKosong = true
			}
		}
		totalPremiAll = tambah(totalPremiAll, bulat4(nilaiPremi))
		comm, net := mu.Baris.Commision, mu.NetPremium
		tidakCocok := beda(premiIns, nilaiPremi) || beda(commIns, comm) || beda(totalIns, net)
		if tidakCocok && !(dekat(premiIns, nilaiPremi) || dekat(commIns, comm) || dekat(totalIns, net)) {
			errIns = "Total premium RNM on installment didn't match with premium RNM!"
		}
		if len(mu.Cicilan) == 0 {
			p.tambah("List Installment can't be null!")
		}
	}
	p.tambah(errIns)
	if premiKosong && totalPremiAll.Sign() <= 0 {
		p.tambah("Premi on paymnet can't be 0")
	}
	return nil
}

// proteksiLokasi - ProtectFIREMBUPA_Act langkah 9 (IsFire), per `.OfferFacIn.LocationList(n)`. TERITORY (9.2.23) tidak
// diperiksa: kolom TERITORYNAME tidak ada di DDL RW. Double Input Case (9.2.32-35) tidak diperiksa: BusinessCode tidak ada.
func (s *Service) proteksiLokasi(ctx context.Context, caseID string, objek []models.ObjekFire, p *pengumpulPesan) error {
	var kode []string
	lihat := map[string]bool{}
	topRiskKasus := false
	for _, o := range objek {
		topRiskKasus = topRiskKasus || o.IsTopRisk
		for _, ok := range o.Occupations {
			if id := ok.OccupationID; id != "" && !lihat[id] {
				lihat[id] = true
				kode = append(kode, id)
			}
		}
	}
	kategori, err := s.submit.KategoriOkupasi(ctx, kode)
	if err != nil {
		return err
	}
	preferred, pesanKategori, topRiskLokal := "", "", false
	for i, o := range objek {
		n := i + 1
		if o.IsTopRisk {
			topRiskLokal = true // Local.Toprisk 9.2.2: tidak pernah dikosongkan lagi di dalam loop
		}
		if o.Province == "" {
			p.tambah(fmt.Sprintf("Province on Location %dcan't be null!", n))
		}
		if o.City == "" {
			p.tambah(fmt.Sprintf("City on Location %dcan't be null!", n))
		}
		if o.ZipCode == "" {
			p.tambah(fmt.Sprintf("ZipCode on Location %dcan't be null!", n))
		}
		if len(o.Occupations) == 0 {
			p.tambah(fmt.Sprintf("Occupation on location %d can't be null!", n))
		}
		if len(o.LossRecords) == 0 && caseID != kasusTanpaLossRecord {
			p.tambah(fmt.Sprintf("Loss record on location %d can't be empty !", n))
		}
		rw, err := s.submit.CekRW(ctx, o.ZipCode, o.City, o.District, o.Province)
		if err != nil {
			return err
		}
		if !rw.Zip {
			p.tambah(fmt.Sprintf("ZipCode on Location %d doesn’t exist in Master!", n))
		}
		// 9.2.24-26: Obj-Browse dilewati bila kode pos tidak di RW -> halaman Cek* kosong -> pesan tetap muncul `[dugaan]`.
		if !rw.City {
			p.tambah(fmt.Sprintf("CITYNAME %s on Location %d doesn’t exist in Master!", o.City, n))
		}
		if !rw.District {
			p.tambah(fmt.Sprintf("DISTRICT %s on Location %d doesn’t exist in Master!", o.District, n))
		}
		if !rw.Province {
			p.tambah(fmt.Sprintf("PROVINCE %s on Location %d doesn’t exist in Master!", o.Province, n))
		}
		// 9.2.27 / 9.2.28 okupasi: OccupationId kosong; kategori -> IsPreferredRisk.
		for _, ok := range o.Occupations {
			if ok.OccupationID == "" {
				p.tambah(fmt.Sprintf("Occupation ID on location%dcan't be null!", n))
			}
			kat := kategori[ok.OccupationID]
			if topRiskKasus == topRiskLokal && !(preferred == "Non-Preferred Risk" || (preferred == "Preferred Risk" && kat != "Non-Preferred Risk")) {
				preferred = kat
			}
			if preferred == "" {
				pesanKategori = fmt.Sprintf("Kategori Okupasi ID %s lokasi%d tidak boleh kosong", ok.OccupationID, n)
			}
		}
		for m, it := range o.Items {
			proteksiItem(p, caseID, n, m+1, it)
		}
	}
	if preferred == "" {
		p.tambah(pesanKategori)
	}
	return nil
}

// proteksiItem - ProtectFIREMBUPA_Act 9.2.30 satu `.Property.PropertyItemList(o)`. IndexProperty / IndexPropertyItem pada
// pesan coverage = nomor lokasi / item `[dugaan]` (pengisinya tidak ditelusur). Deductible (9.2.30.6.6-7, 13) hanya untuk
// FlagOnGoingPolicy 2 - tidak di tahap ini.
func proteksiItem(p *pengumpulPesan, caseID string, n, o int, it models.ItemObjek) {
	if it.Currency == "" {
		p.tambah(fmt.Sprintf("Currency on location %d object Item %d must be choose!", n, o))
	}
	if it.ItemTypeID == "" {
		p.tambah(fmt.Sprintf("Choose Object Item on location %d !", n))
	}
	netRate := it.TotalNetRate != nil && it.TotalNetRate.Sign() > 0
	ada := map[string]bool{}
	for c, cv := range it.Coverages {
		if cv.Coverage == "" || cv.CoverageNote == "" {
			p.tambah(fmt.Sprintf(" Coverage %d on location %d object item %d can't be empty!", c+1, n, o))
		}
		if (cv.AccumulationCode == "" || cv.AccumulationCode == "UNKNOWNID") && caseID != kasusTanpaLossRecord {
			p.tambah(fmt.Sprintf(" Accumulation on location %d object item %d can't be empty!", n, o))
		}
		ada[cv.OldID] = true
	}
	if netRate {
		for _, r := range rateCoverageNetRate {
			if !ada[r] {
				p.tambah(fmt.Sprintf("coverage on location %d object item %d can't be Net Rate!", n, o))
				break
			}
		}
	}
}
