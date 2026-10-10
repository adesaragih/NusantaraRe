package services

// Tab Object keluarga Aneka (tiket 73 P1): baca / Save daftar lokasi -> okupasi Aneka / Golf -> objek. Proteksi objek
// (`ProtectionObject<Lini>_Act`, 422) di proteksianeka.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// DenganObjekAneka memasang penyimpan tab Object Aneka (tiket 73 P1).
func (s *Service) DenganObjekAneka(o repository.PenyimpanObjekAneka) *Service {
	s.objekAneka = o
	return s
}

// suretyTypeAneka - `[terverifikasi]` `NB FacIn\Activity\InputDtlObjectAneka_PreAct.xml` (riset form tiket 73 bagian B): loop
// LocationList -> AnekaList, `.MemorandumAccept.SuretyType` di-SET menurut `pyWorkPage.Quotation.BusinessCode` 5 digit; kode
// lain tidak men-SET apa pun (dan SuretyType tidak ditulis di tempat lain di korpus -> kosong). Arti kode = belum terverifikasi.
var suretyTypeAneka = map[string]string{
	"10076": "SB-BB", "10007": "SB-PB", "10077": "SB-APB", "10083": "SB-MB", "10011": "SB-PayB", "10075": "CB", "10126": "KBG",
}

// lebarTeksAneka - lebar kolom VARCHAR2 (BYTE, migrasi 214) tiap medan teks objek; jalur = nama kabel.
var lebarTeksAneka = []struct {
	jalur string
	ambil func(models.ObjekAneka) string
	n     int
}{
	{"objectName", func(o models.ObjekAneka) string { return o.ObjectName }, 4000},
	{"section", func(o models.ObjekAneka) string { return o.Section }, 50},
	{"noKontrak", func(o models.ObjekAneka) string { return o.NoKontrak }, 50},
	{"noKtp", func(o models.ObjekAneka) string { return o.NoKtp }, 50},
	{"noNpwp", func(o models.ObjekAneka) string { return o.NoNpwp }, 50},
	{"currency", func(o models.ObjekAneka) string { return o.Currency }, 10},
	{"draftWordingId", func(o models.ObjekAneka) string { return o.DraftWordingID }, 50},
	{"others", func(o models.ObjekAneka) string { return o.Others }, 500},
	{"year", func(o models.ObjekAneka) string { return o.Year }, 50},
	{"quantity", func(o models.ObjekAneka) string { return o.Quantity }, 50},
	{"unitQuantity", func(o models.ObjekAneka) string { return o.UnitQuantity }, 50},
	{"serialNo", func(o models.ObjekAneka) string { return o.SerialNo }, 50},
	{"builtIn", func(o models.ObjekAneka) string { return o.BuiltIn }, 50},
	{"flag", func(o models.ObjekAneka) string { return o.Flag }, 1000},
	{"geographicalLimits", func(o models.ObjekAneka) string { return o.GeographicalLimits }, 50},
	{"invoiceNo", func(o models.ObjekAneka) string { return o.InvoiceNo }, 50},
	{"almRiskId", func(o models.ObjekAneka) string { return o.AlmRiskID }, 100},
	{"briguna.loanAccountNumber", func(o models.ObjekAneka) string { return o.Briguna.LoanAccountNumber }, 50},
	{"briguna.loanName", func(o models.ObjekAneka) string { return o.Briguna.LoanName }, 500},
	{"briguna.job", func(o models.ObjekAneka) string { return o.Briguna.Job }, 500},
	{"briguna.depositAccountNo", func(o models.ObjekAneka) string { return o.Briguna.DepositAccountNo }, 50},
	{"briguna.idCardNo", func(o models.ObjekAneka) string { return o.Briguna.IDCardNo }, 50},
	{"briguna.institution", func(o models.ObjekAneka) string { return o.Briguna.Institution }, 500},
	{"briguna.loanType", func(o models.ObjekAneka) string { return o.Briguna.LoanType }, 50},
	{"marineHull.typeOfVessel", func(o models.ObjekAneka) string { return o.MarineHull.TypeOfVessel }, 50},
	{"marineHull.nrtrt", func(o models.ObjekAneka) string { return o.MarineHull.NRTRT }, 1000},
	{"marineHull.lengthm", func(o models.ObjekAneka) string { return o.MarineHull.LengthM }, 50},
	{"marineHull.classification", func(o models.ObjekAneka) string { return o.MarineHull.Classification }, 1000},
	{"marineHull.grtrt", func(o models.ObjekAneka) string { return o.MarineHull.GRTRT }, 1000},
	{"marineHull.breadthm", func(o models.ObjekAneka) string { return o.MarineHull.BreadthM }, 50},
	{"marineHull.construction", func(o models.ObjekAneka) string { return o.MarineHull.Construction }, 1000},
	{"marineHull.dwtrt", func(o models.ObjekAneka) string { return o.MarineHull.DWTRT }, 1000},
	{"marineHull.depthm", func(o models.ObjekAneka) string { return o.MarineHull.DepthM }, 50},
	{"aviationHull.registrationMarks", func(o models.ObjekAneka) string { return o.AviationHull.RegistrationMarks }, 50},
	{"aviationHull.maxPassengers", func(o models.ObjekAneka) string { return o.AviationHull.MaxPassengers }, 50},
	{"aviationHull.pilots", func(o models.ObjekAneka) string { return o.AviationHull.Pilots }, 50},
	{"aviationHull.standardUses", func(o models.ObjekAneka) string { return o.AviationHull.StandardUses }, 50},
	{"aviationHull.specialUses", func(o models.ObjekAneka) string { return o.AviationHull.SpecialUses }, 50},
	{"aviationHull.specialRentalUses", func(o models.ObjekAneka) string { return o.AviationHull.SpecialRentalUses }, 50},
	{"vehicleHe.objectNameHe", func(o models.ObjekAneka) string { return o.VehicleHE.ObjectNameHE }, 500},
	{"vehicleHe.objectIdHe", func(o models.ObjekAneka) string { return o.VehicleHE.ObjectIDHE }, 50},
	{"vehicleHe.brandName", func(o models.ObjekAneka) string { return o.VehicleHE.BrandName }, 500},
	{"vehicleHe.brand", func(o models.ObjekAneka) string { return o.VehicleHE.Brand }, 50},
	{"vehicleHe.typeName", func(o models.ObjekAneka) string { return o.VehicleHE.TypeName }, 500},
	{"vehicleHe.type", func(o models.ObjekAneka) string { return o.VehicleHE.Type }, 50},
	{"vehicleHe.chassisNumber", func(o models.ObjekAneka) string { return o.VehicleHE.ChassisNumber }, 50},
	{"vehicleHe.engineNumber", func(o models.ObjekAneka) string { return o.VehicleHE.EngineNumber }, 50},
	{"memorandum.applicationNumber", func(o models.ObjekAneka) string { return o.Memorandum.ApplicationNumber }, 500},
	{"memorandum.accumulation", func(o models.ObjekAneka) string { return o.Memorandum.Accumulation }, 500},
	{"memorandum.collateral", func(o models.ObjekAneka) string { return o.Memorandum.Collateral }, 500},
	{"memorandum.insuredSiup", func(o models.ObjekAneka) string { return o.Memorandum.InsuredSiup }, 500},
	{"memorandum.keyPerson", func(o models.ObjekAneka) string { return o.Memorandum.KeyPerson }, 500},
	{"memorandum.originalRate", func(o models.ObjekAneka) string { return o.Memorandum.OriginalRate }, 500},
	{"memorandum.principalsOutGo", func(o models.ObjekAneka) string { return o.Memorandum.PrincipalsOutGo }, 500},
	{"memorandum.typeOfWording", func(o models.ObjekAneka) string { return o.Memorandum.TypeOfWording }, 500},
	{"memorandum.buildingStatus", func(o models.ObjekAneka) string { return o.Memorandum.BuildingStatus }, 4000},
	{"memorandum.supportingDocument", func(o models.ObjekAneka) string { return o.Memorandum.SupportingDocument }, 4000},
	{"memorandum.principalsData", func(o models.ObjekAneka) string { return o.Memorandum.PrincipalsData }, 4000},
	{"memorandum.conclusionsProposals", func(o models.ObjekAneka) string { return o.Memorandum.ConclusionsProposals }, 4000},
	{"memorandum.principalsBackground", func(o models.ObjekAneka) string { return o.Memorandum.PrincipalsBackground }, 4000},
	{"memorandum.summaryFinancial", func(o models.ObjekAneka) string { return o.Memorandum.SummaryFinancial }, 4000},
	{"memorandum.principalsGuarantee", func(o models.ObjekAneka) string { return o.Memorandum.PrincipalsGuarantee }, 4000},
}

// lebarLokasiAneka - medan alamat lokasi = kolom tab Object FIRE (186 / 192).
var lebarLokasiAneka = []struct {
	jalur string
	ambil func(models.LokasiAneka) string
	n     int
}{
	{"objectNo", func(l models.LokasiAneka) string { return l.ObjectNo }, 50},
	{"roadType", func(l models.LokasiAneka) string { return l.RoadType }, lebarAlamatRingkas},
	{"roadName", func(l models.LokasiAneka) string { return l.RoadName }, lebarRiskAddress},
	{"buildingNo", func(l models.LokasiAneka) string { return l.BuildingNo }, 50},
	{"zipCode", func(l models.LokasiAneka) string { return l.ZipCode }, lebarAlamatRingkas},
	{"country", func(l models.LokasiAneka) string { return l.Country }, lebarAlamatRingkas},
	{"riskLocation", func(l models.LokasiAneka) string { return l.RiskLocation }, lebarRiskAddress},
	{"territory", func(l models.LokasiAneka) string { return l.Territory }, lebarAlamatRingkas},
	{"city", func(l models.LokasiAneka) string { return l.City }, lebarAlamatRingkas},
	{"district", func(l models.LokasiAneka) string { return l.District }, lebarAlamatRingkas},
	{"province", func(l models.LokasiAneka) string { return l.Province }, lebarAlamatRingkas},
	{"riskAddressId", func(l models.LokasiAneka) string { return l.RiskAddressID }, lebarAlamatRingkas},
}

// tanggalAneka - tiga medan tanggal objek (kabel DD-MM-YYYY, disimpan teks Pega DateTime pola Date of Loss).
var tanggalAneka = []struct {
	jalur string
	ke    func(*models.ObjekAneka) *string
}{
	{"startDate", func(o *models.ObjekAneka) *string { return &o.StartDate }},
	{"endDate", func(o *models.ObjekAneka) *string { return &o.EndDate }},
	{"briguna.dateOfBirth", func(o *models.ObjekAneka) *string { return &o.Briguna.DateOfBirth }},
}

// telusuriObjekAneka - setiap objek (okupasi lalu Golf) dengan jalur kabelnya.
func telusuriObjekAneka(lokasi []models.LokasiAneka, f func(jalur string, o *models.ObjekAneka)) {
	for i := range lokasi {
		for j := range lokasi[i].Okupasi {
			for k := range lokasi[i].Okupasi[j].Objek {
				f(fmt.Sprintf("lokasi[%d].okupasi[%d].objek[%d]", i, j, k), &lokasi[i].Okupasi[j].Objek[k])
			}
		}
		for k := range lokasi[i].ObjekGolf {
			f(fmt.Sprintf("lokasi[%d].objekGolf[%d]", i, k), &lokasi[i].ObjekGolf[k])
		}
	}
}

// periksaAneka - batas baris (SEQ_NO NUMBER(5)), lebar kolom, tanggal DD-MM-YYYY. Nilai pilihan (Section, Currency di luar
// CURRENCY, Draft Wording) tidak dicocokkan di sini kecuali mata uang (periksaMataUangAneka).
func periksaAneka(lokasi []models.LokasiAneka) error {
	var masalah []string
	if len(lokasi) > batasBarisObjek {
		return fmt.Errorf("%w: paling banyak %d lokasi", ErrMasukanObjek, batasBarisObjek)
	}
	for i, l := range lokasi {
		for _, x := range lebarLokasiAneka {
			if len(x.ambil(l)) > x.n {
				masalah = append(masalah, fmt.Sprintf("lokasi[%d].%s paling banyak %d byte", i, x.jalur, x.n))
			}
		}
		if len(l.Okupasi) > batasBarisObjek || len(l.ObjekGolf) > batasBarisObjek {
			masalah = append(masalah, fmt.Sprintf("lokasi[%d]: paling banyak %d okupasi / objek", i, batasBarisObjek))
		}
		for j, o := range l.Okupasi {
			if len(o.OccupationID) > 1000 || len(o.OccupationName) > 1000 {
				masalah = append(masalah, fmt.Sprintf("lokasi[%d].okupasi[%d] occupationId / occupationName paling banyak 1000 byte", i, j))
			}
			if len(o.Objek) > batasBarisObjek {
				masalah = append(masalah, fmt.Sprintf("lokasi[%d].okupasi[%d]: paling banyak %d objek", i, j, batasBarisObjek))
			}
		}
	}
	telusuriObjekAneka(lokasi, func(jalur string, o *models.ObjekAneka) {
		for _, x := range lebarTeksAneka {
			if len(x.ambil(*o)) > x.n {
				masalah = append(masalah, fmt.Sprintf("%s.%s paling banyak %d byte", jalur, x.jalur, x.n))
			}
		}
		for _, d := range tanggalAneka {
			if v := *d.ke(o); v != "" {
				if _, err := time.ParseInLocation(BentukTanggalKabel, v, WIB); err != nil {
					masalah = append(masalah, fmt.Sprintf("%s.%s bukan tanggal DD-MM-YYYY", jalur, d.jalur))
				}
			}
		}
		masalah = append(masalah, periksaCoverageAneka(jalur, o.Coverages)...)
	})
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}

// periksaMataUangAneka - mata uang objek yang diisi ada di CURRENCY (pilihan BrowseCurrency_RD, pola tiket 39). Kosong boleh
// (Aviation Hull C7 tanpa Currency). Basis data dibaca hanya bila ada yang diisi.
func (s *Service) periksaMataUangAneka(ctx context.Context, lokasi []models.LokasiAneka) error {
	var di []mataUangDi
	telusuriObjekAneka(lokasi, func(jalur string, o *models.ObjekAneka) {
		if o.Currency != "" {
			di = append(di, mataUangDi{jalur, o.Currency})
		}
		// tiket 73 P2: mata uang coverage dan deductible yang diisi.
		for k, c := range o.Coverages {
			jc := fmt.Sprintf("%s.coverages[%d]", jalur, k)
			if c.Currency != "" {
				di = append(di, mataUangDi{jc, c.Currency})
			}
			for d, de := range c.Deductibles {
				if de.Currency != "" {
					di = append(di, mataUangDi{fmt.Sprintf("%s.deductibles[%d]", jc, d), de.Currency})
				}
			}
		}
	})
	if len(di) == 0 {
		return nil
	}
	if s.mataUang == nil {
		return ErrPilihanItemTanpaDatabase
	}
	daftar, err := s.mataUang.DaftarMataUang(ctx)
	if err != nil {
		return err
	}
	sah := make(map[string]bool, len(daftar))
	for _, c := range daftar {
		sah[c] = true
	}
	var masalah []string
	for _, mu := range di {
		if !sah[mu.kode] {
			masalah = append(masalah, fmt.Sprintf("%s.currency %q tidak ada di daftar mata uang", mu.jalur, mu.kode))
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}

// salinAneka - salinan larik (lokasi, okupasi, objek) supaya pengubahan medan objek tidak menyentuh larik pemanggil.
func salinAneka(lokasi []models.LokasiAneka) []models.LokasiAneka {
	hasil := make([]models.LokasiAneka, len(lokasi))
	for i, l := range lokasi {
		ok := make([]models.OkupasiAneka, len(l.Okupasi))
		for j, o := range l.Okupasi {
			o.Objek = salinObjekAneka(o.Objek)
			ok[j] = o
		}
		l.Okupasi, l.ObjekGolf = ok, salinObjekAneka(l.ObjekGolf)
		hasil[i] = l
	}
	return hasil
}

// salinObjekAneka - salinan larik objek beserta larik coverage-nya (nil tetap nil: arti "coverage tidak dikirim").
func salinObjekAneka(objek []models.ObjekAneka) []models.ObjekAneka {
	hasil := append([]models.ObjekAneka{}, objek...)
	for i := range hasil {
		if hasil[i].Coverages != nil {
			hasil[i].Coverages = append([]models.CoverageAneka{}, hasil[i].Coverages...)
		}
	}
	return hasil
}

// lebarCoverageAneka - lebar kolom VARCHAR2 T_COVERAGELIST (193 / 195 / 214) medan teks coverage Aneka.
var lebarCoverageAneka = []struct {
	jalur string
	ambil func(models.CoverageAneka) string
	n     int
}{
	{"coverage", func(c models.CoverageAneka) string { return c.Coverage }, 50},
	{"oldId", func(c models.CoverageAneka) string { return c.OldID }, 50},
	{"coverageNote", func(c models.CoverageAneka) string { return c.CoverageNote }, 500},
	{"calculateMethodFacIn", func(c models.CoverageAneka) string { return c.CalculateMethodFacIn }, 50},
	{"calculateMethod", func(c models.CoverageAneka) string { return c.CalculateMethod }, 50},
	{"day", func(c models.CoverageAneka) string { return c.Day }, 50},
	{"currency", func(c models.CoverageAneka) string { return c.Currency }, 10},
	{"indemnity", func(c models.CoverageAneka) string { return c.Indemnity }, 50},
	{"unit", func(c models.CoverageAneka) string { return c.Unit }, 50},
	{"subLimitNote", func(c models.CoverageAneka) string { return c.SubLimitNote }, 500},
}

// lebarLostLimit - T_COVERAGELIST.LOST_LIMIT VARCHAR2(50) rancangan (teks desimal).
const lebarLostLimit = 50

// periksaCoverageAneka - coverage satu objek (nil = tidak dikirim, tidak diperiksa): batas baris, lebar kolom, deductible.
func periksaCoverageAneka(jalurObjek string, cov []models.CoverageAneka) []string {
	if len(cov) > batasItem {
		return []string{fmt.Sprintf("%s.coverages paling banyak %d", jalurObjek, batasItem)}
	}
	var masalah []string
	for k, c := range cov {
		jc := fmt.Sprintf("%s.coverages[%d]", jalurObjek, k)
		for _, l := range lebarCoverageAneka {
			if len(l.ambil(c)) > l.n {
				masalah = append(masalah, fmt.Sprintf("%s.%s paling banyak %d byte", jc, l.jalur, l.n))
			}
		}
		if c.LostLimit != nil && len(utils.FormatDecimal(c.LostLimit)) > lebarLostLimit {
			masalah = append(masalah, fmt.Sprintf("%s.lostLimit paling banyak %d karakter", jc, lebarLostLimit))
		}
		masalah = append(masalah, periksaDeductible(jc+".", c.Deductibles)...)
	}
	return masalah
}

// siapkanAneka - salinan untuk disimpan: SuretyType dari BusinessCode (isian klien diabaikan), tanggal kabel -> teks Pega.
func siapkanAneka(lokasi []models.LokasiAneka, kodeBisnis string) []models.LokasiAneka {
	hasil := salinAneka(lokasi)
	telusuriObjekAneka(hasil, func(_ string, o *models.ObjekAneka) {
		o.Memorandum.SuretyType = suretyTypeAneka[kodeBisnis]
		for _, d := range tanggalAneka {
			*d.ke(o) = tanggalKerugianKePega(*d.ke(o))
		}
	})
	return hasil
}

// BacaObjekAneka - GET tab Object Aneka: lokasi urut SEQ_NO (selalu larik), tanggal teks Pega -> DD-MM-YYYY.
func (s *Service) BacaObjekAneka(ctx context.Context, id string) ([]models.LokasiAneka, error) {
	if s.objekAneka == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	lokasi = salinAneka(lokasi)
	telusuriObjekAneka(lokasi, func(_ string, o *models.ObjekAneka) {
		for _, d := range tanggalAneka {
			*d.ke(o) = tanggalKerugianKeKabel(*d.ke(o))
		}
	})
	// tiket 73 P2: total per mata uang per lokasi (gerbang menurut lini case; case tak terbaca -> kosong).
	kb, _ := s.kasusBisnisDari(ctx, id)
	for i := range lokasi {
		lokasi[i].TotalPerCurrency = totalAnekaLokasi(id, kb, lokasi[i])
	}
	return lokasi, nil
}

// pertahankanSpreadingAneka - tiket 73 P3 (pola K48-7 FIRE): coverage yang DIKIRIM ulang kehilangan baris lamanya (dihapus
// repository), jadi `.SpreadingList` coverage tersimpan di posisi (jalur objek, indeks coverage) yang oldId-nya sama dipasang ke
// coverage baru. Coverage yang tidak dikirim (nil) dipindah utuh oleh repository bersama spreading-nya.
func (s *Service) pertahankanSpreadingAneka(ctx context.Context, id string, simpan []models.LokasiAneka) error {
	dikirim := false
	telusuriObjekAneka(simpan, func(_ string, o *models.ObjekAneka) { dikirim = dikirim || o.Coverages != nil })
	if !dikirim {
		return nil
	}
	lama, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return ErrKasusTidakAda
	}
	if err != nil {
		return err
	}
	tersimpan := map[string][]models.CoverageAneka{}
	telusuriObjekAneka(lama, func(jalur string, o *models.ObjekAneka) { tersimpan[jalur] = o.Coverages })
	telusuriObjekAneka(simpan, func(jalur string, o *models.ObjekAneka) {
		for m := range o.Coverages {
			c := &o.Coverages[m]
			c.Spreading = []models.BarisSpreading{}
			if l := tersimpan[jalur]; m < len(l) && l[m].OldID == c.OldID {
				c.Spreading = salinBaris(l[m].Spreading)
			}
		}
	})
	return nil
}

// GantiObjekAneka - Save tab Object Aneka: identitas (401) -> hak ubah tahap (409, sama dengan tab Object FIRE) -> isian (400)
// -> basis data (503) -> mata uang (400) -> proteksi objek (422) -> case (404); seluruh lokasi diganti di satu transaksi, lalu
// dibaca ulang.
func (s *Service) GantiObjekAneka(ctx context.Context, pelaku inti.Pelaku, id string, lokasi []models.LokasiAneka,
	proRate *apd.Decimal) ([]models.LokasiAneka, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return nil, err
	}
	if err := periksaAneka(lokasi); err != nil {
		return nil, err
	}
	if s.objekAneka == nil || s.transaksi == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	if err := s.periksaMataUangAneka(ctx, lokasi); err != nil {
		return nil, err
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	if p := proteksiObjekAneka(id, kb, lokasi); len(p) > 0 {
		return nil, GalatProteksiAneka{Pelanggaran: p}
	}
	simpan := siapkanAneka(lokasi, kb.businessCode)
	if err := s.pertahankanSpreadingAneka(ctx, id, simpan); err != nil {
		return nil, err
	}
	share, err := s.objekAneka.PercentShareKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	medanServerAneka(simpan, share)
	// tiket 73 P3 (pola K48-7): spreading coverage yang dikirim dihitung ulang dari % Share RNM tersimpan seperti setiap Save Pega
	// (DT percent + ClaimEstimation cabang 2.4); coverage yang tidak dikirim tidak berubah nilainya.
	if share != nil {
		l := liniSpreadingDari(id, kb)
		terapkanPersenAneka(simpan, l)
		totalLokasiAneka(simpan, share, l)
		bulatkanSimpanAneka(simpan)
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.objekAneka.GantiObjekAneka(ctx, tx, id, simpan, simpan8(proRate)) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaObjekAneka(ctx, id)
}
