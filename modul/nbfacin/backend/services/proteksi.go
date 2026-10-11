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

// flagProteksi - properti halaman Protect* yang di-SET Pega bersama sebuah pesan proteksi (tiket 57b). Dibaca syarat tombol
// Send Email Banding; pesan tanpa flag (Scoring Risk) = "".
type flagProteksi string

const (
	flagSpreading1  flagProteksi = "ProtectSpreading.CARI1"
	flagSpreading23 flagProteksi = "ProtectSpreading.CARI23"
	flagSpreading28 flagProteksi = "ProtectSpreading.CARI28"
	flagPayment1    flagProteksi = "ProtectPayment.CARI1"
	flagMarine2     flagProteksi = "ProtectMarine.CARI2"
	flagNetRate1    flagProteksi = "ProtectNeRate.CARI1"
	flagZipCode1    flagProteksi = "ProtectZipCode.CARI1"
	flagOkupasi1    flagProteksi = "ProtectOccupation.CARI1"
	flagMarketing1  flagProteksi = "ProtectMarketing.CARI1"
	flagDeductible1 flagProteksi = "ProtectDeductible.CARI1"
	flagDeductible2 flagProteksi = "ProtectDeductible.CARI2"
	flagDeductible4 flagProteksi = "ProtectDeductible.CARI4"
)

// Flag yang MENAHAN Banding = flag Protect* di syarat Visible tombol Send Email Banding (tiket 57b, keputusan work owner
// 08-10-2026 "disamakan"):
//   - flagBandingMarketing - `Section\EmailSection.xml` (tahap marketing): `ProtectSpreading.CARI1 == 1|| ProtectSpreading.CARI23
//     == 1 || … || ProtectPayment.CARI1==1 || ProtectMarine.CARI2==1|| ProtectNeRate.CARI1==…`.
//   - flagBandingBinding - `Section\EmailSectionCeding.xml` (tahap Binding): hanya `ProtectSpreading.CARI1 == 1||
//     ProtectSpreading.CARI23 == 1` - BERBEDA dari EmailSection.
var (
	flagBandingMarketing = []flagProteksi{flagSpreading1, flagSpreading23, flagPayment1, flagMarine2, flagNetRate1}
	flagBandingBinding   = []flagProteksi{flagSpreading1, flagSpreading23}
)

// Penahan tahap Binding lain (tiket 64b, `Section\EmailSectionCeding.xml` pyDisabledWhen `[terverifikasi]` sesi 0f):
//   - flagAcceptBinding - tombol "Send Email" Accept (SetAkseptasiProposal): `… || ProtectSpreading.CARI23 = 1
//     ||ProtectSpreading.CARI1 == 1||ProtectZipCode.CARI1==1` (+ `.OfferFacIn.ConfirmBinding != 1` = pesanBelumBinding);
//   - flagTombolBinding - tombol "Binding": `… || ProtectSpreading.CARI23 = 1 || ProtectDeductible.CARI2=='1' ||
//     ProtectDeductible.CARI4=='1'||ProtectSpreading.CARI1 == 1` (+ ConfirmBinding = 1 = ErrSudahBinding).
var (
	flagAcceptBinding = []flagProteksi{flagSpreading1, flagSpreading23, flagZipCode1}
	flagTombolBinding = []flagProteksi{flagSpreading1, flagSpreading23, flagDeductible2, flagDeductible4}
)

// pengumpulPesan - Page-Set-Messages: urutan munculnya, kembar sekali; tiap pesan membawa flag Protect* yang di-SET
// langkah yang sama (tiket 57b).
type pengumpulPesan struct {
	pesan []string
	flag  [][]flagProteksi
	lihat map[string]int
}

func (p *pengumpulPesan) tambah(s string) { p.tambahFlag(s, "") }

func (p *pengumpulPesan) tambahFlag(s string, f flagProteksi) {
	if s == "" {
		return
	}
	if i, ada := p.lihat[s]; ada {
		if f != "" && !dalamFlag(p.flag[i], f) {
			p.flag[i] = append(p.flag[i], f)
		}
		return
	}
	if p.lihat == nil {
		p.lihat = map[string]int{}
	}
	p.lihat[s] = len(p.pesan)
	p.pesan = append(p.pesan, s)
	var fs []flagProteksi
	if f != "" {
		fs = []flagProteksi{f}
	}
	p.flag = append(p.flag, fs)
}

// penahan - pesan yang flag-nya ada di himpunan tahan, urutan asli; [] bila tidak ada.
func (p *pengumpulPesan) penahan(tahan []flagProteksi) []string {
	out := []string{}
	for i, s := range p.pesan {
		for _, f := range p.flag[i] {
			if dalamFlag(tahan, f) {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

func dalamFlag(fs []flagProteksi, f flagProteksi) bool {
	for _, x := range fs {
		if x == f {
			return true
		}
	}
	return false
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

// isianDeductibleTL - isian Team Leader untuk Protection_Act 9-12.
type isianDeductibleTL struct {
	terima   bool
	komentar string
}

// deductibleBedaMataUang - padanan Pega `IsBedaCurr` (`Activity\SaveDeductible.xml` 2-3: deductible ber-Currency terisi yang
// beda dari Currency item induknya). Pega menyetelnya per Save satu deductible (nilai = deductible TERAKHIR disimpan); di
// sini dihitung atas SELURUH deductible tersimpan `[keputusan agent]`. Currency kosong / 'UNKNOWN' (DEFAULT kolom, A169) =
// tidak diisi.
func deductibleBedaMataUang(objek []models.ObjekFire) bool {
	for _, o := range objek {
		for _, it := range o.Items {
			for _, c := range it.Coverages {
				for _, d := range c.Deductibles {
					if cur := strings.TrimSpace(d.Currency); cur != "" && cur != "UNKNOWN" && cur != it.Currency {
						return true
					}
				}
			}
		}
	}
	return false
}

// ProteksiAccept - pesan proteksi Accept; kosong = lolos. tl = isian Team Leader (nil = tahap marketing).
func (s *Service) ProteksiAccept(ctx context.Context, id string, tl *isianDeductibleTL) ([]string, error) {
	p, err := s.proteksiBerflag(ctx, id, tl, false)
	if err != nil {
		return nil, err
	}
	if p.pesan == nil {
		return []string{}, nil
	}
	return p.pesan, nil
}

// proteksiBanding - Send Email Banding (tiket 57b, satu fungsi tahap marketing dan Binding): port Protection_Act yang sama
// dengan Accept dijalankan, tetapi yang MENAHAN hanya pesan ber-flag di himpunan tahan (flagBandingMarketing /
// flagBandingBinding); pesan lain tidak ditampilkan dan tidak menahan.
func (s *Service) proteksiBanding(ctx context.Context, id string, tahan []flagProteksi) ([]string, error) {
	p, err := s.proteksiBerflag(ctx, id, nil, false)
	if err != nil {
		return nil, err
	}
	return p.penahan(tahan), nil
}

// proteksiBerflag - badan Protection_Act; pesan + flag Protect* per pesan.
// deductibleBinding - ProtectFIREMBUPA_Act 9.2.30.6.6 / 9.2.30.6.13.2 berjalan (FlagOnGoingPolicy 2 DAN EmailTypeBinding
// bukan 4): tahap Binding selain Banding (tiket 64b).
func (s *Service) proteksiBerflag(ctx context.Context, id string, tl *isianDeductibleTL, deductibleBinding bool) (*pengumpulPesan, error) {
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
	p := &pengumpulPesan{}
	// Langkah 2: keluar tanpa proteksi.
	if insuredDikecualikan(kasus.Opportunity.InsuredID) || len(objek) > batasLokasiProteksi {
		return p, nil
	}
	// Tiket 73 P4 (K73-28): case keluarga Aneka (IsAneka - termasuk Bonding / Engineering - atau isGolfInsurance) tidak
	// menjalankan proteksi IsFire (langkah 18 Scoring Risk L3873, 21 ProtectCurrencyTSI L1301, 24 ProtectFIREMBUPA L17817);
	// langkah 24 memakai cabang ProtectCoverage Aneka / Golf. Lini tak terbaca -> perilaku FIRE (tiket 57) tidak berubah.
	aneka, err := s.dataProteksiAneka(ctx, id, kasus, deductibleBinding)
	if err != nil {
		return nil, err
	}
	// Tiket 74 P5: case MARINE CARGO - langkah IsFire tidak berjalan; langkah 24 cabang MARINE, langkah 27 MARINE CARGO.
	cargo, adaCargo, err := s.dataProteksiCargo(ctx, id)
	if err != nil {
		return nil, err
	}
	// Langkah 9-12 (PositionNote ReasFacInTeamLeader): IsBedaCurr && IsDeductibleAcceptance "false" -> "Deductible Acceptance
	// cannot be empty"; IsDeductibleAcceptance "true" && CommentDeductible kosong -> "Comment Deductible cannot be empty".
	// Centang tidak dicentang = "false" `[dugaan]` (checkbox Pega mengirim "false").
	if tl != nil {
		beda := deductibleBedaMataUang(objek)
		if aneka != nil {
			beda = deductibleBedaMataUangAneka(aneka.lokasi) // K73-30 `[dugaan]`
		}
		if adaCargo {
			beda = false // addDeductibleCargo_ACT L248 IsBedaCurr = 0
		}
		if beda && !tl.terima {
			p.tambahFlag("Deductible Acceptance cannot be empty", flagDeductible1)
		}
		if tl.terima && tl.komentar == "" {
			p.tambahFlag("Comment Deductible cannot be empty", flagDeductible2)
		}
	}
	// Langkah 16-17 MOID kosong; langkah 18 Scoring Risk (bukan EDM, bukan B2B ASM, IsFire).
	if strings.TrimSpace(kasus.General.MarketingID) == "" {
		p.tambahFlag(pesanMOIDKosong, flagMarketing1) // langkah 16
	}
	if aneka == nil && !adaCargo && kasus.RiskScoring.FinalScore == nil {
		p.tambah(pesanScoringKosong)
	}
	// Langkah 21 ProtectCurrencyTSI_Act (IsFire).
	if aneka == nil && !adaCargo {
		if err := s.proteksiMataUang(ctx, objek, p); err != nil {
			return nil, err
		}
	}
	// Langkah 22 CheckSpreadingProtect_ACT: langkah 2 keluar untuk PositionNote ReasFacInMarketing kecuali Group - tidak ada
	// pesan di tahap ini.
	// Langkah 23 ProtectPremiPolicy_Act.
	if err := s.proteksiPayment(ctx, id, kodePolicyType(kasus.General.PolicyType), p); err != nil {
		return nil, err
	}
	// Langkah 24 ProtectCoverage_Act -> cabang ANEKA / GOLF (tiket 73 P4) atau ProtectFIREMBUPA_Act (IsFire).
	switch {
	case aneka != nil:
		if err := s.proteksiCoverageAneka(ctx, aneka.lokasi, aneka.isi, p); err != nil {
			return nil, err
		}
	case adaCargo:
		if err := s.proteksiCoverageCargo(ctx, cargo, p); err != nil {
			return nil, err
		}
	default:
		if err := s.proteksiLokasi(ctx, kasus.CaseID, objek, p, deductibleBinding); err != nil {
			return nil, err
		}
	}
	// Langkah 27 MARINE CARGO (tiket 74 P5): ProtectShipData_Act + "Multi Location Must Have Toprisk".
	if adaCargo {
		proteksiKapalCargo(cargo, p)
	}
	// Langkah 28 ProtectShareCedant_Act (keluar bila bukan Group), 29 ProtectRenewal_Act (bukan renewal), 30
	// ProtectCedingCo (IsB2B ASM), 31-34 Survey report / Ceding Co (FlagOnGoingPolicy 2 / Group / data survey) - lihat
	// komentar berkas.
	return p, nil
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
	p.tambahFlag(terakhir, flagSpreading1) // 1.2.1.3
	return nil
}

// proteksiDeductible - ProtectFIREMBUPA_Act 9.2.30.6.6 (`@SizeOfPropertyList(.DeductibleList)=0` -> CARI4, "should be added!")
// dan 9.2.30.6.13.2 (`.TypeDeductible == "" && .Condition == "" && .MinMax == "" && .Amount == "" && .TimeExcess == ""` ->
// CARI2, "can't be empty!" - dua spasi sebelum DeductibleList VERBATIM).
func proteksiDeductible(p *pengumpulPesan, n, o, c int, ded []models.Deductible) {
	if len(ded) == 0 {
		p.tambahFlag(fmt.Sprintf(" Deductible on location %d object item %d Coverage %d should be added!", n, o, c), flagDeductible4)
	}
	for d, x := range ded {
		if x.TypeDeductible == "" && x.Condition == "" && x.MinMax == "" && x.Amount == nil && x.TimeExcess == nil {
			p.tambahFlag(fmt.Sprintf(" Deductible on location %d object item %d Coverage %d  DeductibleList %d can't be empty!",
				n, o, c, d+1), flagDeductible2)
		}
	}
}

// kosongAtauNol4 - `== "" || == "0.0000"` atas @divide(x,1,4).
func kosongAtauNol4(d *apd.Decimal) bool { return d == nil || bulat4(d).IsZero() }

// beda / dekat - `a <> b` dan |a - b| < 0.001 atas nilai 4 desimal (bulat4 / selisihKecil tab Payment).
func beda(a, b *apd.Decimal) bool { return bulat4(nolBila(a)).Cmp(bulat4(nolBila(b))) != 0 }

func dekat(a, b *apd.Decimal) bool { return selisihKecil(bulat4(nolBila(a)), bulat4(nolBila(b))) }

// proteksiPayment - ProtectPremiPolicy_Act langkah 3-11 per `.OfferFacIn.CurrencyList` (mata uang tab Payment). Bukan Life,
// bukan MOP (IsMOP tidak ada di sistem baru), bukan EDM. Langkah 7 / 11 ("Premi on paymnet can't be 0") dilewati bila
// `PolicyType==1` (Master Policy) `[terverifikasi]` - kode standard value (policytype.go, ralat butir 78.2). Langkah 8-9 (Master
// Policy + MOP, premi harus 0) tidak berlaku: IsMOP tidak ada.
func (s *Service) proteksiPayment(ctx context.Context, id, policyType string, p *pengumpulPesan) error {
	t, err := s.TampilanPayment(ctx, id)
	if err != nil {
		return err
	}
	if len(t.MataUang) == 0 {
		p.tambahFlag("List Payment can't be null!", flagPayment1) // langkah 3
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
			p.tambahFlag("List Installment can't be null!", flagSpreading1) // 5.8 (syarat = 5.7)
		}
	}
	p.tambahFlag(errIns, flagSpreading28) // 5.5
	if premiKosong && totalPremiAll.Sign() <= 0 && policyType != policyTypeMaster {
		p.tambahFlag("Premi on paymnet can't be 0", flagSpreading23) // langkah 7 / 11
	}
	return nil
}

// proteksiLokasi - ProtectFIREMBUPA_Act langkah 9 (IsFire), per `.OfferFacIn.LocationList(n)`. TERITORY (9.2.23) tidak
// diperiksa: kolom TERITORYNAME tidak ada di DDL RW. Double Input Case (9.2.32-35) tidak diperiksa: BusinessCode tidak ada.
// langkahPreferred - satu okupasi 9.2.28 (= SetFlagOccupation_ACT, dipakai juga jenjang Accept tahap uw tiket 61): kategori
// okupasi menjadi IsPreferredRisk bila status top risk lokasi sama dengan kasus, kecuali IsPreferredRisk sudah
// "Non-Preferred Risk", atau sudah "Preferred Risk" dan kategori bukan "Non-Preferred Risk".
func langkahPreferred(preferred, kat string, cocokTopRisk bool) string {
	if cocokTopRisk && !(preferred == "Non-Preferred Risk" || (preferred == "Preferred Risk" && kat != "Non-Preferred Risk")) {
		return kat
	}
	return preferred
}

func (s *Service) proteksiLokasi(ctx context.Context, caseID string, objek []models.ObjekFire, p *pengumpulPesan,
	deductibleBinding bool) error {
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
			p.tambahFlag(fmt.Sprintf("Province on Location %dcan't be null!", n), flagSpreading1) // 9.2.3
		}
		if o.City == "" {
			p.tambahFlag(fmt.Sprintf("City on Location %dcan't be null!", n), flagSpreading1) // 9.2.5
		}
		if o.ZipCode == "" {
			// 9.2.7: satu Property-Set (ZipCode kosong ATAU OccupationList kosong ATAU loss record kosong) men-SET CARI1;
			// When kode 5 = Skip Whens (KEPUTUSAN-30-09-2026: sisa When dilewati, langkah dijalankan `[dugaan kuat]`) -> ATAU.
			p.tambahFlag(fmt.Sprintf("ZipCode on Location %dcan't be null!", n), flagSpreading1)
		}
		if len(o.Occupations) == 0 {
			p.tambahFlag(fmt.Sprintf("Occupation on location %d can't be null!", n), flagSpreading1)
		}
		if len(o.LossRecords) == 0 && caseID != kasusTanpaLossRecord {
			p.tambahFlag(fmt.Sprintf("Loss record on location %d can't be empty !", n), flagSpreading1)
		}
		rw, err := s.submit.CekRW(ctx, o.ZipCode, o.City, o.District, o.Province)
		if err != nil {
			return err
		}
		if !rw.Zip {
			p.tambahFlag(fmt.Sprintf("ZipCode on Location %d doesn’t exist in Master!", n), flagZipCode1) // 9.2.21
		}
		// 9.2.24-26: Obj-Browse dilewati bila kode pos tidak di RW -> halaman Cek* kosong -> pesan tetap muncul `[dugaan]`.
		if !rw.City {
			p.tambahFlag(fmt.Sprintf("CITYNAME %s on Location %d doesn’t exist in Master!", o.City, n), flagZipCode1) // 9.2.18
		}
		if !rw.District {
			p.tambahFlag(fmt.Sprintf("DISTRICT %s on Location %d doesn’t exist in Master!", o.District, n), flagZipCode1) // 9.2.19
		}
		if !rw.Province {
			p.tambahFlag(fmt.Sprintf("PROVINCE %s on Location %d doesn’t exist in Master!", o.Province, n), flagZipCode1) // 9.2.20
		}
		// 9.2.27 / 9.2.28 okupasi: OccupationId kosong; kategori -> IsPreferredRisk.
		for _, ok := range o.Occupations {
			if ok.OccupationID == "" {
				p.tambahFlag(fmt.Sprintf("Occupation ID on location%dcan't be null!", n), flagSpreading1) // 9.2.27.3 / 9.2.28.3
			}
			preferred = langkahPreferred(preferred, kategori[ok.OccupationID], topRiskKasus == topRiskLokal)
			if preferred == "" {
				pesanKategori = fmt.Sprintf("Kategori Okupasi ID %s lokasi%d tidak boleh kosong", ok.OccupationID, n)
			}
		}
		for m, it := range o.Items {
			proteksiItem(p, caseID, n, m+1, it, deductibleBinding)
		}
	}
	if preferred == "" {
		p.tambahFlag(pesanKategori, flagOkupasi1) // 9.2.29
	}
	return nil
}

// proteksiItem - ProtectFIREMBUPA_Act 9.2.30 satu `.Property.PropertyItemList(o)`. IndexProperty / IndexPropertyItem pada
// pesan coverage = nomor lokasi / item `[dugaan]` (pengisinya tidak ditelusur). Deductible (9.2.30.6.6-7, 13) hanya bila
// deductibleBinding (FlagOnGoingPolicy 2, EmailTypeBinding bukan 4; tiket 64b): per coverage SESUDAH pesan coverage /
// akumulasinya - DeductibleList kosong (CARI4) lalu tiap baris deductible kosong seluruhnya (CARI2). IdxObject = nomor item,
// IdxCoverage / IdxDeductible = .pxListSubscript (mulai 1). Pega tidak mereset CARI4 / CARI2 di dalam loop (Page-Set-Messages
// mengulang pesan lama) - pesan kembar ditampilkan sekali, hasilnya sama.
func proteksiItem(p *pengumpulPesan, caseID string, n, o int, it models.ItemObjek, deductibleBinding bool) {
	if it.Currency == "" {
		// 9.2.30.3: Currency ATAU ItemTypeID kosong (When kode 5 = Skip Whens `[dugaan kuat]`) men-SET CARI1.
		p.tambahFlag(fmt.Sprintf("Currency on location %d object Item %d must be choose!", n, o), flagSpreading1)
	}
	if it.ItemTypeID == "" {
		p.tambahFlag(fmt.Sprintf("Choose Object Item on location %d !", n), flagSpreading1)
	}
	netRate := it.TotalNetRate != nil && it.TotalNetRate.Sign() > 0
	ada := map[string]bool{}
	for c, cv := range it.Coverages {
		if cv.Coverage == "" || cv.CoverageNote == "" {
			p.tambahFlag(fmt.Sprintf(" Coverage %d on location %d object item %d can't be empty!", c+1, n, o), flagMarine2) // 9.2.30.6.2
		}
		if (cv.AccumulationCode == "" || cv.AccumulationCode == "UNKNOWNID") && caseID != kasusTanpaLossRecord {
			p.tambahFlag(fmt.Sprintf(" Accumulation on location %d object item %d can't be empty!", n, o), flagMarine2)
		}
		if deductibleBinding {
			proteksiDeductible(p, n, o, c+1, cv.Deductibles)
		}
		ada[cv.OldID] = true
	}
	if netRate {
		for _, r := range rateCoverageNetRate {
			if !ada[r] {
				p.tambahFlag(fmt.Sprintf("coverage on location %d object item %d can't be Net Rate!", n, o), flagNetRate1) // 9.2.30.7
				break
			}
		}
	}
}
