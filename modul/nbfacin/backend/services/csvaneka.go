package services

// Unggah CSV keluarga ANEKA, mulai Marine Hull (tiket 73 P6, keputusan work owner 11-10-2026 K73-52..62) di atas infrastruktur
// tiket 75 / 74 P6 (berkas mengalir -> penampung, validasi latar, simpan latar per kelompok, Ulangi, kunci case). `[terverifikasi]`
// `D:\migrasi\RNM\NB FacIn\Activity\`:
//   - ChooseActivityMappingCSV_Act langkah 5: IsAneka -> InsertUploadAneka_act (jenis pekerjaan "aneka"). Template menurut
//     SetDownloadMasterCSV (LiniKasus.JenisCsv): hanya "marinehull" yang templatenya ada (K73-60); "aneka" / "cpm" -> 409
//     ErrCsvTemplateAnekaBelumAda sampai template dikirim.
//   - Judul dibaca per NAMA: models.KolomCsvAneka (35); judul wajib = 31 judul template Marine Hull (K73-60), 4 lainnya opsional.
//   - InsertUploadAneka_act: 1 Property-Remove LocationList + CurrencyList (GANTI); 4 per baris ListB2BHostUpload:
//     4.1 NO kosong -> baris dilewati (K73-54: dibandingkan dengan baris TERPROSES sebelumnya, bukan indeks Pega yang bergeser);
//     4.2 ProtectUploadAneka_Act, galat PERTAMA menghentikan impor (TRANS F=6, K73-62);
//     4.6 / 4.7 lokasi baru bila NO berubah: ObjectNo = NO, RoadType / RoadName / BuildingNo / Country, RiskLocation = TYPE_ADDRESS
//     + " " + ADDRESS, ZipCode = POSTAL_CODE, kota / provinsi / district / RW dari master RW (K73-58: kode pos + KELURAHAN bila
//     terisi, selain itu kode pos saja); NO yang muncul lagi tidak berurutan = galat (K73-53);
//     4.9 OccupationList(1) saja, ditimpa bila NO / OCCUPATION_NAME berubah;
//     4.10-4.11 mata uang GetCurrencyIDByName (tidak ada = galat, K73-55);
//     4.12 objek baru bila NO / OBJECT_DESCRIPTION / SECTION berubah (medan dari baris PERTAMA objek): ObjectName, Currency,
//     Section (SECTION I 1 / SECTION II 2 / WARRANTY 3 / lainnya 4), TSI, SerialNo, VehicleHE Chassis / Engine, Year, MarineHull
//     NRTRT / GRTRT / DWTRT / Construction / TypeOfVessel (OBJECT_TYPE) / Classification, Flag (FLAGNAME), BuiltIn
//     (PLACE_OF_MANUFACTURE);
//     4.13 coverage baru bila objek / COVERAGE_NAME berubah: Coverage / CoverageNote = baris pertama BrowseCoverageFacIn_RD (Contains,
//     K73-59), CalculateMethod_FacIn (PRO RATA 1 / SHORT PERIOD 2 / lainnya 3), Day (kosong 365), Loading = PCT_LOADING, LostLimit
//     = PCT_SCALE, Rate = PCT_RATE, Premium / DiscountPercentage / Discount / PctLoL / LimitofLiability / TSI dari CSV APA ADANYA
//     tanpa hitung ulang (K73-57), Currency = mata uang objek;
//     angka TSI / PREMIUM / DISCOUNT / LOL: "." dibuang lalu "," -> "."; PCT_*: "," -> "." (wajib terurai, preseden K66 D7);
//     5 / 6 / 7 CurrencyList TIDAK dibangun ulang - dihapus beserta cicilan (K73-56 = K66 D5; Java langkah 6 cacat);
//     8 Work-.Save (medan server coverage tidak dihitung - CountPremiAndTSIRNM* hanya di Save layar).
//   - ProtectUploadAneka_Act (pesan verbatim): "Occupation Cannot NULL" lalu "Coverage Cannot NULL" (Property-Set MENIMPA);
//     "* Occupation <nama> Not Available in Master" (BrowseOccupationANEKA_RD, Contains, baris pertama, K73-59) lalu KELUAR
//     (TRANS 5.4 menguji TradingTemp - halaman cargo yang tidak ada di alur Aneka `[dugaan]`); "* Coverage <nama> Not Available in
//     Master".
//   - Bonding / CustomBonds (IsAneka) diimpor ke AnekaList generik seperti Pega (K73-61).

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrCsvTemplateAnekaBelumAda - case Aneka yang templatenya (TemplateCSVAneka / TemplateCSVAneka_CPM) belum dikirim (K73-60). 409.
var ErrCsvTemplateAnekaBelumAda = errors.New("Template CSV untuk Class of Business ini (Pega TemplateCSVAneka / TemplateCSVAneka_CPM) " +
	"belum tersedia di sistem baru - unggah CSV baru tersedia untuk Marine Hull.")

const (
	// pesanOkupasiKosongCsv / pesanCoverageKosongCsv - ProtectUploadAneka_Act langkah 2 / 3 verbatim.
	pesanOkupasiKosongCsv  = "Occupation Cannot NULL"
	pesanCoverageKosongCsv = "Coverage Cannot NULL"
)

// DenganUnggahCsvAneka memasang penulis objek Aneka hasil unggah (tiket 73 P6).
func (s *Service) DenganUnggahCsvAneka(p repository.PenulisAnekaCsv) *Service {
	s.penulisAnekaCsv = p
	return s
}

// petaJudulAneka - judul Aneka per NAMA (tak peka huruf / spasi tepi), urutan bebas; nama tak dikenal diabaikan; ganda / judul wajib
// template Marine Hull kurang -> ErrMasukanObjek.
func petaJudulAneka(baris string) (petaJudulCsv, error) {
	sel := pecahBarisCsv(baris, ';')
	p := petaJudulCsv{idx: make([]int, len(sel)), n: len(models.KolomCsvAneka)}
	ada := map[string]bool{}
	for i, s := range sel {
		k := strings.ToUpper(strings.TrimSpace(s))
		p.idx[i] = slices.Index(models.KolomCsvAneka, k)
		if p.idx[i] < 0 {
			continue
		}
		if ada[k] {
			return petaJudulCsv{}, fmt.Errorf("%w: judul berkas CSV Marine Hull memuat kolom %s dua kali", ErrMasukanObjek, k)
		}
		ada[k] = true
		p.kolom = append(p.kolom, k)
	}
	var kurang []string
	for _, k := range models.KolomCsvAneka[:models.JumlahJudulMarineHull] {
		if !ada[k] {
			kurang = append(kurang, k)
		}
	}
	if len(kurang) > 0 {
		return petaJudulCsv{}, fmt.Errorf("%w: judul berkas CSV tidak sesuai template Marine Hull - kolom %s tidak ada (judul dibaca per "+
			"nama, pemisah \";\")", ErrMasukanObjek, strings.Join(kurang, ", "))
	}
	return p, nil
}

// ---- satu baris ----------------------------------------------------------------------------------------------------

// barisAnekaCsv - satu baris CSV Aneka terurai (angka / pemetaan Insert 4.12-4.13) dan hasil master (Protect / 4.10).
type barisAnekaCsv struct {
	no                                       int
	teks                                     map[string]string
	section, calculateMethod, day            string
	tsi, premium, discount, lol              *apd.Decimal
	pctLoading, pctScale, pctRate, pctDiscnt *apd.Decimal
	pctLol                                   *apd.Decimal
	// hasil master
	okupasi  models.PilihanOkupasiAneka
	coverage models.BarisCoverage
	mu       models.MataUangCsv
}

func (b barisAnekaCsv) t(k string) string { return b.teks[k] }

// lebarAnekaCsv - kolom CSV yang ditulis ke kolom bertepi (lebar tujuan tab Object Aneka: lebarLokasiAneka / lebarTeksAneka /
// lebarCoverageAneka, migrasi 186 / 192 / 214).
var lebarAnekaCsv = []struct {
	kolom string
	n     int
}{{"NO", 50}, {"TYPE_ADDRESS", lebarAlamatRingkas}, {"ADDRESS", lebarRiskAddress}, {"BUILDING_NO", 50},
	{"POSTAL_CODE", lebarAlamatRingkas}, {"COUNTRY", lebarAlamatRingkas}, {"OBJECT_DESCRIPTION", 4000}, {"CURRENCY", 10},
	{"YEAR", 50}, {"PLACE_OF_MANUFACTURE", 50}, {"OBJECT_TYPE", 50}, {"NRT", 1000}, {"GRT", 1000}, {"DWT", 1000},
	{"FLAGNAME", 1000}, {"CONSTRUCTION", 1000}, {"CLASSIFICATION", 1000}, {"DAY", 50}, {"SERIALNO", 50},
	{"CHASSISNUMBER", 50}, {"ENGINENUMBER", 50}}

// angkaAnekaCsv - kolom angka: ribuan = "." dibuang lalu "," -> "." (TSI / PREMIUM / DISCOUNT / LOL); persen = "," -> "." saja.
var angkaAnekaCsv = []struct {
	kolom  string
	ribuan bool
	ke     func(*barisAnekaCsv) **apd.Decimal
}{
	{"TSI", true, func(b *barisAnekaCsv) **apd.Decimal { return &b.tsi }},
	{"PREMIUM", true, func(b *barisAnekaCsv) **apd.Decimal { return &b.premium }},
	{"DISCOUNT", true, func(b *barisAnekaCsv) **apd.Decimal { return &b.discount }},
	{"LOL", true, func(b *barisAnekaCsv) **apd.Decimal { return &b.lol }},
	{"PCT_LOADING", false, func(b *barisAnekaCsv) **apd.Decimal { return &b.pctLoading }},
	{"PCT_SCALE", false, func(b *barisAnekaCsv) **apd.Decimal { return &b.pctScale }},
	{"PCT_RATE", false, func(b *barisAnekaCsv) **apd.Decimal { return &b.pctRate }},
	{"PCT_DISCOUNT", false, func(b *barisAnekaCsv) **apd.Decimal { return &b.pctDiscnt }},
	{"PCT_LOL", false, func(b *barisAnekaCsv) **apd.Decimal { return &b.pctLol }},
}

// sectionAnekaCsv - Insert 4.12 `@if(.SECTION="SECTION I","1",@if(.SECTION="SECTION II","2",@if(.SECTION="WARRANTY","3","4")))`.
func sectionAnekaCsv(s string) string {
	switch s {
	case "SECTION I":
		return "1"
	case "SECTION II":
		return "2"
	case "WARRANTY":
		return "3"
	}
	return "4"
}

// metodeAnekaCsv - Insert 4.13 `@if(.CALCULATE_METHOD=="PRO RATA","1",@if(.CALCULATE_METHOD=="SHORT PERIOD","2","3"))`.
func metodeAnekaCsv(s string) string {
	switch s {
	case "PRO RATA":
		return "1"
	case "SHORT PERIOD":
		return "2"
	}
	return "3"
}

// uraiBarisAnekaCsv - lebar + angka (K66 D7) + pemetaan; galat dikembalikan, tidak dihentikan.
func uraiBarisAnekaCsv(no int, isi []string) (barisAnekaCsv, []string) {
	b := barisAnekaCsv{no: no, teks: make(map[string]string, len(models.KolomCsvAneka))}
	for i, k := range models.KolomCsvAneka {
		if i < len(isi) {
			b.teks[k] = isi[i]
		}
	}
	var masalah []string
	for _, l := range lebarAnekaCsv {
		if len(b.t(l.kolom)) > l.n {
			masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: paling banyak %d byte", no, l.kolom, l.n))
		}
	}
	if n := len(b.t("TYPE_ADDRESS")) + 1 + len(b.t("ADDRESS")); n > lebarRiskAddress {
		masalah = append(masalah, fmt.Sprintf("baris %d kolom TYPE_ADDRESS + ADDRESS: Risk Location paling banyak %d byte", no,
			lebarRiskAddress))
	}
	for _, a := range angkaAnekaCsv {
		s := b.t(a.kolom)
		if s == "" {
			continue
		}
		if a.ribuan {
			s = normalAngkaCsv(s)
		} else {
			s = strings.ReplaceAll(s, ",", ".")
		}
		d, err := utils.ParseDecimal(s)
		switch {
		case err != nil:
			masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: angka tidak sah %q", no, a.kolom, b.t(a.kolom)))
		case !repository.PolaAngkaMasuk.MatchString(utils.FormatDecimal(simpan8(d))):
			masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: melebihi 30 digit bulat", no, a.kolom))
		default:
			*a.ke(&b) = simpan8(d)
		}
	}
	// K73-63 (keputusan work owner 11-10-2026, PENYIMPANGAN SADAR dari Pega - Pega tanpa proteksi premi > TSI): PREMIUM baris >
	// TSI baris yang sama (sesudah urai) = galat; sama = sah. Hanya bila keduanya terisi `[keputusan agent]`.
	if b.premium != nil && b.tsi != nil && lebihBesar(b.premium, b.tsi) {
		masalah = append(masalah, fmt.Sprintf("baris %d: PREMIUM %s melebihi TSI %s", no, utils.FormatDecimal(b.premium),
			utils.FormatDecimal(b.tsi)))
	}
	if b.pctScale != nil && len(utils.FormatDecimal(b.pctScale)) > lebarLostLimit {
		masalah = append(masalah, fmt.Sprintf("baris %d kolom PCT_SCALE: paling banyak %d karakter", no, lebarLostLimit))
	}
	b.section, b.calculateMethod, b.day = sectionAnekaCsv(b.t("SECTION")), metodeAnekaCsv(b.t("CALCULATE_METHOD")), b.t("DAY")
	if b.day == "" {
		b.day = "365"
	}
	return b, masalah
}

// cacheAnekaCsv - hasil master per pekerjaan (validasi / simpan) dan konteks case.
type cacheAnekaCsv struct {
	kelas, kodeBisnis string // Class of Business (BusinessName okupasi) / BUSINESS.ID (BizCode coverage)
	okupasi           map[string]*models.PilihanOkupasiAneka
	coverage          map[string]*models.BarisCoverage
	rw                map[string]*models.RWCsv
}

// cacheAnekaCsv - BusinessName = Class of Business case (pola popup tiket 73 P1), BizCode = BUSINESS.ID (tiket 73 P2).
func (s *Service) cacheAnekaCsv(ctx context.Context, id string) (*cacheAnekaCsv, error) {
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	return &cacheAnekaCsv{kelas: strings.TrimSpace(k.Opportunity.ClassOfBusiness), kodeBisnis: kb.businessCode,
		okupasi: map[string]*models.PilihanOkupasiAneka{}, coverage: map[string]*models.BarisCoverage{},
		rw: map[string]*models.RWCsv{}}, nil
}

// protectBarisAnekaCsv - ProtectUploadAneka_Act (lihat kepala berkas). Mengisi okupasi / coverage pada b.
func (s *Service) protectBarisAnekaCsv(ctx context.Context, b *barisAnekaCsv, c *cacheAnekaCsv) ([]string, error) {
	var protect []string
	if b.t("OCCUPATION_NAME") == "" {
		protect = []string{pesanOkupasiKosongCsv}
	}
	if b.t("COVERAGE_NAME") == "" {
		protect = []string{pesanCoverageKosongCsv} // Property-Set ERRMSG MENIMPA
	}
	if nama := b.t("OCCUPATION_NAME"); nama != "" {
		o, ada := c.okupasi[nama]
		if !ada {
			baris, err := s.lookupAneka.CariOkupasiAneka(ctx, c.kelas, nama)
			if err != nil {
				return nil, err
			}
			if len(baris) > 0 {
				o = &baris[0]
			}
			c.okupasi[nama] = o
		}
		if o == nil {
			return append(protect, fmt.Sprintf(pesanMasterCsv, "Occupation", nama)), nil
		}
		b.okupasi = *o
	}
	if nama := b.t("COVERAGE_NAME"); nama != "" {
		v, ada := c.coverage[nama]
		if !ada {
			baris, err := s.lookupAneka.CariCoverageAneka(ctx, c.kodeBisnis, nama)
			if err != nil {
				return nil, err
			}
			if len(baris) > 0 {
				v = &baris[0]
			}
			c.coverage[nama] = v
		}
		if v == nil {
			return append(protect, fmt.Sprintf(pesanMasterCsv, "Coverage", nama)), nil
		}
		b.coverage = *v
	}
	return protect, nil
}

// periksaBarisAnekaCsv - Protect lalu mata uang (K73-55). Galat basis data = error.
func (s *Service) periksaBarisAnekaCsv(ctx context.Context, b *barisAnekaCsv, c *cacheAnekaCsv) ([]string, error) {
	masalah, err := s.protectBarisAnekaCsv(ctx, b, c)
	if err != nil {
		return nil, err
	}
	var m models.MataUangCsv
	ok := false
	if b.t("CURRENCY") != "" {
		if m, ok, err = s.csv.MataUang(ctx, b.t("CURRENCY")); err != nil {
			return nil, err
		}
	}
	if !ok {
		return append(masalah, fmt.Sprintf("baris %d kolom CURRENCY: mata uang %q tidak ada di master CURRENCY", b.no, b.t("CURRENCY"))), nil
	}
	b.mu = m
	return masalah, nil
}

// rwAnekaCsv - K73-58: KELURAHAN terisi -> RW ber-kode pos + kelurahan (BrowseRW2_SQL, baris pertama); kosong -> RW pertama
// ber-kode pos. Kode pos kosong / tidak ada -> nil (kota / provinsi / district / RW kosong, Pega OutRW tanpa baris).
func (s *Service) rwAnekaCsv(ctx context.Context, b barisAnekaCsv, c *cacheAnekaCsv) (*models.RWCsv, error) {
	zip, kel := b.t("POSTAL_CODE"), b.t("KELURAHAN")
	if zip == "" {
		return nil, nil
	}
	k := zip + "|" + kel
	if r, ada := c.rw[k]; ada {
		return r, nil
	}
	var r *models.RWCsv
	if kel != "" {
		baris, err := s.csv.RWKelurahan(ctx, zip, kel)
		if err != nil {
			return nil, err
		}
		if len(baris) > 0 {
			r = &baris[0]
		}
	} else {
		v, ada, err := s.csv.RWKodePos(ctx, zip)
		if err != nil {
			return nil, err
		}
		if ada {
			r = &v
		}
	}
	c.rw[k] = r
	return r, nil
}

// lebarRWAnekaCsv - medan master RW yang ditulis lokasi CSV Aneka dan kolom tujuannya (migrasi 192, VARCHAR2(100) BYTE); master
// RW VARCHAR2(4000) (`DDL\RW.txt`). Keputusan work owner 11-10-2026 (audit lebar sesudah bug ORA-12899): tanpa migrasi - melebihi
// = galat jelas per baris (lebar BYTE = len UTF-8).
var lebarRWAnekaCsv = []struct {
	medan, kolom string
	ambil        func(models.RWCsv) string
}{
	{"NOTE", "T_RISKLOCATION.ASMRW", func(r models.RWCsv) string { return r.Note }},
	{"DISTRICTNAME", "T_RISKLOCATION.ASM_DISTRICT", func(r models.RWCsv) string { return r.District }},
	{"CITYNAME", "T_RISKLOCATION.ASM_CITY", func(r models.RWCsv) string { return r.City }},
	{"PROVINCENAME", "T_PROPERTY.PROVINCE", func(r models.RWCsv) string { return r.Province }},
}

// periksaRWAnekaCsv - lihat lebarRWAnekaCsv (rw nil = tidak ada RW, sah).
func periksaRWAnekaCsv(no int, zip string, rw *models.RWCsv) []string {
	if rw == nil {
		return nil
	}
	var masalah []string
	for _, l := range lebarRWAnekaCsv {
		if n := len(l.ambil(*rw)); n > lebarAlamatRingkas {
			masalah = append(masalah, fmt.Sprintf("baris %d: master RW kode pos %q - %s %d byte melebihi kolom %s (paling banyak %d byte)",
				no, zip, l.medan, n, l.kolom, lebarAlamatRingkas))
		}
	}
	return masalah
}

// ---- susun pohon --------------------------------------------------------------------------------------------------------

// susunAnekaCsv - pohon lokasi Aneka dibangun baris demi baris (Insert 4.7 / 4.9 / 4.12 / 4.13); sebelum = baris TERPROSES
// sebelumnya (K73-54).
type susunAnekaCsv struct {
	lokasi  []models.LokasiAneka
	sebelum *barisAnekaCsv
}

// lokasiBaru - Insert 4.7 `NO != NO baris sebelumnya` (baris pertama selalu).
func (z *susunAnekaCsv) lokasiBaru(b barisAnekaCsv) bool {
	return z.sebelum == nil || b.t("NO") != z.sebelum.t("NO")
}

// terapkan - satu baris ber-NO; rw = master RW lokasi (dipakai hanya bila lokasi baru).
func (z *susunAnekaCsv) terapkan(b barisAnekaCsv, rw *models.RWCsv) {
	p := z.sebelum
	lokasiBaru := z.lokasiBaru(b)
	if lokasiBaru {
		l := models.LokasiAneka{ObjectNo: b.t("NO"), RoadType: b.t("TYPE_ADDRESS"), RoadName: b.t("ADDRESS"),
			BuildingNo: b.t("BUILDING_NO"), RiskLocation: b.t("TYPE_ADDRESS") + " " + b.t("ADDRESS"), Country: b.t("COUNTRY"),
			ZipCode: b.t("POSTAL_CODE")}
		if rw != nil {
			l.City, l.Province, l.District, l.Territory = rw.City, rw.Province, rw.District, rw.Note
		}
		l.Okupasi = []models.OkupasiAneka{{}}
		z.lokasi = append(z.lokasi, l)
	}
	l := &z.lokasi[len(z.lokasi)-1]
	if lokasiBaru || b.t("OCCUPATION_NAME") != p.t("OCCUPATION_NAME") {
		l.Okupasi[0].OccupationID, l.Okupasi[0].OccupationName = b.okupasi.ID, b.okupasi.Name
	}
	ok := &l.Okupasi[0]
	objekBaru := lokasiBaru || b.t("OBJECT_DESCRIPTION") != p.t("OBJECT_DESCRIPTION") || b.t("SECTION") != p.t("SECTION")
	if objekBaru {
		ok.Objek = append(ok.Objek, models.ObjekAneka{ObjectName: b.t("OBJECT_DESCRIPTION"), Section: b.section,
			Currency: b.t("CURRENCY"), TSI: b.tsi, SerialNo: b.t("SERIALNO"), Year: b.t("YEAR"), BuiltIn: b.t("PLACE_OF_MANUFACTURE"),
			Flag: b.t("FLAGNAME"),
			MarineHull: models.MarineHull{TypeOfVessel: b.t("OBJECT_TYPE"), NRTRT: b.t("NRT"), GRTRT: b.t("GRT"), DWTRT: b.t("DWT"),
				Construction: b.t("CONSTRUCTION"), Classification: b.t("CLASSIFICATION")},
			VehicleHE: models.VehicleHE{ChassisNumber: b.t("CHASSISNUMBER"), EngineNumber: b.t("ENGINENUMBER")},
			Coverages: []models.CoverageAneka{}})
	}
	o := &ok.Objek[len(ok.Objek)-1]
	if objekBaru || b.t("COVERAGE_NAME") != p.t("COVERAGE_NAME") {
		o.Coverages = append(o.Coverages, models.CoverageAneka{Coverage: b.coverage.ID, CoverageNote: b.coverage.Nama,
			CalculateMethodFacIn: b.calculateMethod, Day: b.day, Loading: b.pctLoading, LostLimit: b.pctScale, Rate: b.pctRate,
			Premium: b.premium, DiscountPercentage: b.pctDiscnt, Discount: b.discount, PctLoL: b.pctLol, LimitOfLiability: b.lol,
			TSI: b.tsi, Currency: o.Currency})
	}
	z.sebelum = &b
}

// ---- validasi & simpan latar ------------------------------------------------------------------------------------------

// validasiUnggahAneka - pekerjaan latar fase validasi jenis aneka: berhenti di baris PERTAMA yang bermasalah (K73-62, Pega
// TRANS F=6); NO berulang tidak berurutan = galat (K73-53).
func (s *Service) validasiUnggahAneka(ctx context.Context, u models.UnggahCsv) {
	c, err := s.cacheAnekaCsv(ctx, u.CaseID)
	if err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	lihat, noSebelum := map[string]bool{}, ""
	for dari := 1; dari <= u.TotalBaris; dari += UkuranKelompokCsv {
		baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, dari, min(dari+UkuranKelompokCsv-1, u.TotalBaris))
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		var galat []models.PesanUnggahCsv
		for _, x := range baris {
			b, masalah := uraiBarisAnekaCsv(x.No, x.Isi)
			if b.t("NO") == "" {
				continue // 4.1 baris dilewati (K73-54)
			}
			lokasiBaru := b.t("NO") != noSebelum
			if no := b.t("NO"); lokasiBaru {
				if lihat[no] {
					masalah = append(masalah, fmt.Sprintf("baris %d kolom NO: lokasi NO %q sudah muncul sebelumnya - baris satu lokasi "+
						"wajib berurutan", x.No, no))
				}
				lihat[no], noSebelum = true, no
			}
			if len(masalah) == 0 {
				if masalah, err = s.periksaBarisAnekaCsv(ctx, &b, c); err != nil {
					s.gagalUnggah(ctx, u, x.No, err)
					return
				}
			}
			if len(masalah) == 0 && lokasiBaru { // RW lokasi baru (audit lebar 11-10-2026)
				rw, err := s.rwAnekaCsv(ctx, b, c)
				if err != nil {
					s.gagalUnggah(ctx, u, x.No, err)
					return
				}
				masalah = periksaRWAnekaCsv(x.No, b.t("POSTAL_CODE"), rw)
			}
			if len(masalah) > 0 {
				for _, m := range masalah {
					galat = append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: m})
				}
				break
			}
		}
		urut := u.JumlahGalat + 1
		u.BarisDiproses, u.JumlahGalat = dari+len(baris)-1, u.JumlahGalat+len(galat)
		if len(galat) > 0 {
			u.BarisDiproses = galat[len(galat)-1].NoBaris
		}
		if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error {
			return s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisGalatUnggah, urut, galat)
		}); err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		if len(galat) > 0 {
			break
		}
	}
	u.Status = StatusUnggahValid
	if u.JumlahGalat > 0 {
		u.Status = StatusUnggahGagalValidasi
	}
	if err := s.tulisUnggah(ctx, u, nil); err != nil {
		s.gagalUnggah(ctx, u, 0, err)
	}
}

// simpanUnggahAneka - pekerjaan latar fase simpan jenis aneka dari baris BARIS_TERSIMPAN+1. Satu kelompok = paling sedikit
// UkuranKelompokCsv baris DAN berakhir di batas lokasi (baris berikutnya membuka lokasi baru) - kelompok tidak pernah memotong
// lokasi, jadi Ulangi selalu mulai di lokasi baru. Kelompok pertama: GANTI objek case + hapus CurrencyList (langkah 1); kemajuan
// ditulis di transaksi objek kelompoknya.
func (s *Service) simpanUnggahAneka(ctx context.Context, u models.UnggahCsv) {
	c, err := s.cacheAnekaCsv(ctx, u.CaseID)
	if err != nil {
		s.gagalUnggah(ctx, u, 0, err)
		return
	}
	for dari := u.BarisTersimpan + 1; dari <= u.TotalBaris; {
		z := &susunAnekaCsv{}
		akhir, jumlah := dari-1, 0
	baca:
		for p := dari; p <= u.TotalBaris; p += UkuranKelompokCsv {
			baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, p, min(p+UkuranKelompokCsv-1, u.TotalBaris))
			if err != nil {
				s.gagalUnggah(ctx, u, dari, err)
				return
			}
			for _, x := range baris {
				b, masalah := uraiBarisAnekaCsv(x.No, x.Isi)
				if b.t("NO") == "" {
					akhir = x.No // 4.1 dilewati; ikut kelompok yang sedang berjalan
					continue
				}
				if len(masalah) == 0 {
					if masalah, err = s.periksaBarisAnekaCsv(ctx, &b, c); err != nil {
						s.gagalUnggah(ctx, u, x.No, err)
						return
					}
				}
				if len(masalah) > 0 {
					s.gagalUnggah(ctx, u, x.No, fmt.Errorf("%w: data master berubah sejak validasi: %s", ErrMasukanObjek, masalah[0]))
					return
				}
				baru := z.lokasiBaru(b)
				if baru && jumlah >= UkuranKelompokCsv {
					break baca
				}
				var rw *models.RWCsv
				if baru {
					if rw, err = s.rwAnekaCsv(ctx, b, c); err != nil {
						s.gagalUnggah(ctx, u, x.No, err)
						return
					}
					if m := periksaRWAnekaCsv(x.No, b.t("POSTAL_CODE"), rw); len(m) > 0 {
						s.gagalUnggah(ctx, u, x.No, fmt.Errorf("%w: data master berubah sejak validasi: %s", ErrMasukanObjek, m[0]))
						return
					}
				}
				z.terapkan(b, rw)
				akhir, jumlah = x.No, jumlah+1
			}
		}
		maju := u
		maju.KelompokTerakhir, maju.BarisTersimpan = u.KelompokTerakhir+1, akhir
		ganti := u.BarisTersimpan == 0
		err := s.transaksi(ctx, func(tx *db.Tx) error {
			if err := s.penulisAnekaCsv.TulisAnekaCsv(ctx, tx, u.CaseID, z.lokasi, ganti); err != nil {
				return err
			}
			if ganti {
				if err := s.csv.HapusCurrencyList(ctx, tx, u.CaseID); err != nil {
					return err
				}
			}
			return s.unggahCsv.UbahUnggah(ctx, tx, maju)
		})
		if errors.Is(err, repository.ErrKasusTidakAda) {
			err = ErrKasusTidakAda
		}
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		u, dari = maju, akhir+1
	}
	u.Status = StatusUnggahSelesai
	if err := s.tulisUnggah(ctx, u, nil); err != nil {
		s.gagalUnggah(ctx, u, 0, err)
	}
}
