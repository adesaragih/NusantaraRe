package services

// Save Data bilah CSV objek General (tiket 66). `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\`:
// ChooseActivityMappingCSV_Act langkah 3 (IsFire) -> InsertUploadFire_act (ASM-FW-GISFW-DATA-OFFERFACIN!INSERTUPLOADFIRE_ACT):
//   - 2 "Please Upload CSV" bila tanpa baris; 3 Property-Remove LocationList + CurrencyList (MENGGANTI, K66-4);
//   - 5.2 normalisasi angka `@replaceAll(@replaceAll(x, ".", ""), ",", ".")` dan huruf besar TYPE_ADDRESS / ADDRESS / COUNTRY /
//     OBJECT_TYPE / KELURAHAN; ACCUMULATION_NAME = ADDRESS;
//   - 5.3 header polis ditimpa per baris; 5.4 SATU lokasi per baris; 5.6-5.8 RW (BrowseRW2_SQL) + pesan KodePos (tidak
//     menghentikan); 5.9 risk address cari-atau-tambah; 5.10 akumulasi cari-atau-tambah; 5.11-5.13 okupasi; 5.14 ListCauseOfLoss
//     "-" (5.13 / 5.14 / 5.10.2.3 ber-pyStepsPreCondition=false = selalu jalan); 5.15 MappingJsonH2HNusare2_act (item ITEM_TYPE,
//     _1.._6 + MappingCoverage); 6 "Success Insert Data" HANYA bila ERRMSG kosong; 7 AddPaymentCurency_ACT; 8 Obj-Save.
// Keputusan work owner / agent 08-10-2026 (issues/66): D1 (b) header hanya dari kolom terisi; D2 (b) Object Type dicocokkan tak
// peka huruf; D3 (a) master lebih dulu, objek + header + currency satu transaksi; D4 akumulasi baru hanya bila provinsi / CZone /
// negara ditemukan; D5 coverage dihitung siapkanCoverage, spreading dibuang, T_CURRENCYLIST + cicilan dihapus; D6 / D7 ketat
// (galat dikumpulkan dalam satu 400); D8 ID coverage ter-hardcode wajib ada di master.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// KolomCsvObjek - 61 kolom template FIRE (`D:\migrasi\RNM\DDL\templatecsvfire.csv` baris judul, K66-3), urutan berkas; sama
// dengan KOLOM_CSV_OBJEK frontend. Kolom 1-55 = parser Pega; 56-61 ditambah work owner 08-10-2026 (bukan kolom Pega):
// KODE_POS / KECAMATAN / KABUPATEN_KOTA / PROVINSI membantu mencari RW (cariRWCsv), SUMBER_KODE_POS / TINGKAT_KEYAKINAN hanya
// ditampilkan (tidak memengaruhi Save Data).
var KolomCsvObjek = models.KolomCsvObjek

const (
	// BatasBarisCsvObjek - paling banyak 3000 baris per unggahan (RALAT D7 08-10-2026, work owner: "jangan dibatasin 500. 3000
	// aja buat batas nya"; sebelumnya 500).
	BatasBarisCsvObjek = 3000
	// pesanUnggahCsv / pesanSuksesCsv - verbatim InsertUploadFire_act langkah 2 (b.492) dan 6 (b.7577).
	pesanUnggahCsv = "Please Upload CSV"
	pesanSuksesCsv = "Success Insert Data"
	// typeFacultativeCsv - langkah 5.3 `TypeFacultative = "FacultativeIn"`; disimpan sebagai label layar "Facultative In"
	// (butir 78.3: data FacultativeIn = layar Facultative In, disimpan "as sent by the screen").
	typeFacultativeCsv = "Facultative In"
	// Akumulasi baru langkah 5.10.2.2 (literal korpus; bukan nama orang).
	akumulasiBaruCsv, tipeAkumulasiCsv, scopeAkumulasiCsv = "10000001", "FLEXAS", "AREA"
)

// ItemCsvObjek - akhiran slot item MappingJsonH2HNusare2_act langkah 5.1-5.7 (ITEM_TYPE, ITEM_TYPE_1 … _6).
var ItemCsvObjek = []string{"", "_1", "_2", "_3", "_4", "_5", "_6"}

// perilCsv - MappingCoverage langkah 2.1-2.5 (literal OLDID / ID master / catatan; D8): coverage dibuat bila RATE_<kolom> > 0.
var perilCsv = []struct{ kolom, oldID, coverage, catatan string }{
	{"FLEXAS", "FLEXAS", "100815", "FLEXAS (PSAKI)"},
	{"RSMDCC", "4.1A CC", "100828", "RSMD 4.1A CIVIL COMMOTION"},
	{"FWTWD", "4.3", "100825", "FLOOD"},
	{"EQ", "4.2 PRGBI", "100829", "Earthquake Shock, Fire Explosion Following EQ, Volcanic Eruption, Tsunami"},
	{"OTHERS", "OTHERS", "100840", "BURGLARY; VEHICLE IMPACT; ACCIDENTAL DAMAGE ETC"},
}

// jenisObjekCsv - 9 nilai baku langkah 5.4 dan Keyword akumulasi baru langkah 5.10.2.2 (selain itu "Others" / "ZONA LAIN-LAIN").
var jenisObjekCsv = []struct{ nama, keyword string }{
	{"Dwelling House", "KOMPLEKS/KAV.PERUMAHAN/ VILLA/PEMUKIMAN"}, {"Shop Houses", "RUKO/RUKAN"},
	{"Factory", "KAWASAN BERIKAT/INDUSTRI"}, {"Office", "BUSINESS CENTRE / OFFICE BUILDING"},
	{"Shop", "MAL/SHOPPING CENTRE/TRADE CENTRE"}, {"Apartment", "APARTEMEN/RUMAH SUSUN"}, {"Restaurant", "RESTORAN"},
	{"Private Warehouse", "KAWASAN PERGUDANGAN"}, {"Public Warehouse", "KAWASAN PERGUDANGAN"},
}

const (
	jenisObjekLain   = "Others"
	keywordObjekLain = "ZONA LAIN-LAIN"
)

// jenisObjekBaku - D2 (b): tak peka huruf ke 9 nilai baku -> ejaan baku; selain itu "Others". Pega membandingkan teks yang
// sudah di-uppercase dengan ejaan campuran (`[dugaan]` selalu "Others") - penyimpangan sadar keputusan work owner.
func jenisObjekBaku(teks string) string {
	for _, j := range jenisObjekCsv {
		if strings.EqualFold(strings.TrimSpace(teks), j.nama) {
			return j.nama
		}
	}
	return jenisObjekLain
}

// keywordAkumulasi - Keyword akumulasi baru menurut Object Type (langkah 5.10.2.2, 9 cabang).
func keywordAkumulasi(jenis string) string {
	for _, j := range jenisObjekCsv {
		if j.nama == jenis {
			return j.keyword
		}
	}
	return keywordObjekLain
}

// kolomAngkaCsv - kolom template yang dinormalisasi langkah 5.2 (TSI_OBJECT_ITEM_6 ikut: Pega lupa, D7) + ADJUSTABLE_PCT
// (tidak dinormalisasi Pega; diurai apa adanya).
var kolomAngkaCsv = func() []string {
	k := []string{"LOL"}
	for _, s := range ItemCsvObjek {
		k = append(k, "TSI_OBJECT_ITEM"+s)
	}
	for _, p := range perilCsv {
		k = append(k, "RATE_"+p.kolom, "NET_RATE_"+p.kolom, "DISCOUNT_RATE_"+p.kolom)
	}
	return k
}()

// barisCsv - satu baris sesudah langkah 5.2: teks (dipangkas; kolom huruf besar sudah di-uppercase) dan angka terurai.
type barisCsv struct {
	n            int // nomor baris CSV (pxListSubscript, mulai 1)
	teks         map[string]string
	angka        map[string]*apd.Decimal
	mulai, akhir time.Time // BEGIN_DATE / END_DATE (zero = kosong)
}

func (b barisCsv) t(k string) string { return b.teks[k] }

// normalAngkaCsv - langkah 5.2: "." dibuang, "," -> "." (format ribuan titik, desimal koma).
func normalAngkaCsv(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
}

// tanggalCsv - BEGIN_DATE / END_DATE dd/MM/yyyy (`[dugaan]` bentuk baris contoh template 99/99/9999) di WIB.
func tanggalCsv(s string) (time.Time, bool) {
	for _, b := range []string{"02/01/2006", "2/1/2006"} {
		if t, err := time.ParseInLocation(b, s, WIB); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// uraiCsvObjek - kunci (55 kolom), jumlah baris, angka dan tanggal; semua galat dikumpulkan dalam satu 400 (D7).
func uraiCsvObjek(baris []map[string]string) ([]barisCsv, error) {
	if len(baris) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrMasukanObjek, pesanUnggahCsv)
	}
	if len(baris) > BatasBarisCsvObjek {
		return nil, fmt.Errorf("%w: paling banyak %d baris CSV", ErrMasukanObjek, BatasBarisCsvObjek)
	}
	var masalah []string
	hasil := make([]barisCsv, len(baris))
	for i, r := range baris {
		b, m := uraiBarisCsv(i+1, r)
		masalah = append(masalah, m...)
		hasil[i] = b
	}
	if len(masalah) > 0 {
		slices.Sort(masalah)
		return nil, fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return hasil, nil
}

// uraiBarisCsv - satu baris CSV nomor n (langkah 5.2): kolom tak dikenal, normalisasi teks / angka, tanggal; galat baris
// dikembalikan (tidak dihentikan). Dipakai jalur JSON (uraiCsvObjek) dan jalur unggah mengalir (tiket 75) - aturan SATU.
func uraiBarisCsv(n int, r map[string]string) (barisCsv, []string) {
	var masalah []string
	b := barisCsv{n: n, teks: map[string]string{}, angka: map[string]*apd.Decimal{}}
	for k := range r {
		if !slices.Contains(KolomCsvObjek, k) {
			masalah = append(masalah, fmt.Sprintf("baris %d: kolom %s tidak dikenal", b.n, k))
		}
	}
	for _, k := range KolomCsvObjek {
		v := strings.TrimSpace(r[k])
		switch k {
		case "TYPE_ADDRESS", "ADDRESS", "OBJECT_TYPE", "KELURAHAN":
			v = strings.ToUpper(v)
		}
		b.teks[k] = v
	}
	b.teks["ACCUMULATION_NAME"] = b.teks["ADDRESS"]
	for _, k := range append(slices.Clone(kolomAngkaCsv), "ADJUSTABLE_PCT") {
		v := b.teks[k]
		if v == "" {
			continue
		}
		if k != "ADJUSTABLE_PCT" {
			v = normalAngkaCsv(v)
		}
		d, err := utils.ParseDecimal(v)
		if err != nil {
			masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: angka tidak sah %q", b.n, k, b.teks[k]))
			continue
		}
		b.angka[k] = d
	}
	for k, ke := range map[string]*time.Time{"BEGIN_DATE": &b.mulai, "END_DATE": &b.akhir} {
		if v := b.teks[k]; v != "" {
			t, ok := tanggalCsv(v)
			if !ok {
				masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: tanggal harus dd/MM/yyyy: %q", b.n, k, v))
			}
			*ke = t
		}
	}
	return b, masalah
}

// headerCsv - langkah 5.3 dengan D1 (b): kolom terisi saja, baris terakhir menang (beda nilai -> log).
type headerCsv struct {
	noOfferSlip  string
	mulai, akhir time.Time
}

func headerDariCsv(id string, baris []barisCsv) headerCsv {
	var h headerCsv
	beda := map[string]bool{}
	ambil := func(nama, v string, ke *string) {
		if v == "" {
			return
		}
		if *ke != "" && *ke != v {
			beda[nama] = true
		}
		*ke = v
	}
	var mulai, akhir string
	for _, b := range baris {
		ambil("NO_OFFER_SLIP", b.t("NO_OFFER_SLIP"), &h.noOfferSlip)
		ambil("BEGIN_DATE", b.t("BEGIN_DATE"), &mulai)
		ambil("END_DATE", b.t("END_DATE"), &akhir)
		if !b.mulai.IsZero() {
			h.mulai = b.mulai
		}
		if !b.akhir.IsZero() {
			h.akhir = b.akhir
		}
	}
	for k := range beda {
		log.Printf("nbfacin: CSV objek kasus %s: %s berbeda antar baris - baris terakhir menang", id, k)
	}
	return h
}

// rencanaCsv - satu lokasi + hasil pencarian master untuk tulisan master sesudahnya.
type rencanaCsv struct {
	objek       models.ObjekFire
	rw          *models.RWCsv
	riskID      string // "" = belum ada di master
	akumulasiID string // "" = belum ada di master
}

// rencanakanBarisCsv - satu baris rencanakanCsv: rencana, pesan KodePos (0 / 1), galat isian (mata uang, D6); ID coverage master
// yang dipakai dicatat ke pakaiCoverage (dicek sekali oleh pemanggil, D8). Galat basis data dikembalikan sebagai error.
func (s *Service) rencanakanBarisCsv(ctx context.Context, b barisCsv, pakaiCoverage map[string]string) (rencanaCsv, []string, []string, error) {
	var pesan, masalah []string
	jenis := jenisObjekBaku(b.t("OBJECT_TYPE"))
	o := models.ObjekFire{ObjectNo: strconv.Itoa(b.n), RoadType: b.t("TYPE_ADDRESS"), RoadName: b.t("ADDRESS"),
		BuildingNo: b.t("BUILDING_NO"), RiskLocation: b.t("TYPE_ADDRESS") + " " + b.t("ADDRESS"), ZipCode: b.t("POSTAL_CODE"),
		IsTopRisk: b.t("TOP_RISK") == "Y", ObjectName: b.t("OBJECT_TYPE"), ObjectType: jenis,
		LossRecords: []models.CatatanKerugian{{Detail: "-"}}}
	r := rencanaCsv{}
	// 5.6-5.8: BrowseRW2_SQL (cariRWCsv); tanpa hasil -> pesan (juga bila kode pos / kelurahan kosong).
	rw, zip, err := s.cariRWCsv(ctx, b)
	if err != nil {
		return rencanaCsv{}, nil, nil, err
	}
	if rw != nil {
		r.rw, o.ZipCode = rw, zip
		o.City, o.District, o.Territory, o.Province, o.Country = rw.City, rw.District, rw.Note, rw.Province, rw.Nation
	}
	if r.rw == nil {
		pesan = append(pesan, "KodePos "+b.t("POSTAL_CODE")+" and kelurahan "+b.t("KELURAHAN")+" pada lokasi ke - "+strconv.Itoa(b.n)+
			" tidak ditemukan!")
	}
	// 5.9 risk address persis ADDRESS + POSTALCODE (kode pos yang menemukan RW).
	id, err := s.csv.IDRiskAddress(ctx, b.t("ADDRESS"), o.ZipCode)
	if err != nil {
		return rencanaCsv{}, nil, nil, err
	}
	r.riskID, o.RiskAddressID = id, id
	// 5.10 akumulasi persis NOTE = ACCUMULATION_NAME (= ADDRESS); tidak ada -> kode CSV dipakai sampai master dibuat.
	if r.akumulasiID, err = s.csv.IDAkumulasi(ctx, b.t("ACCUMULATION_NAME")); err != nil {
		return rencanaCsv{}, nil, nil, err
	}
	// 5.11-5.13 okupasi (selalu satu per lokasi); nama kosong -> master OCCUPATION FIRE ber-OLDID kode. Kode kosong tidak
	// dicari (RD Pega membuang filter kosong dan memungut okupasi pertama - tidak ditiru, keputusan agent).
	nama := b.t("OCCUPATION_NAME")
	if nama == "" && b.t("OCCUPATION_CODE") != "" {
		if nama, err = s.csv.NamaOkupasiFire(ctx, b.t("OCCUPATION_CODE")); err != nil {
			return rencanaCsv{}, nil, nil, err
		}
	}
	o.Occupations = []models.OkupasiObjek{{OccupationID: b.t("OCCUPATION_CODE"), OccupationName: nama,
		ConstructionClass: b.t("CLASS_OF_CONSTRUCTION")}}
	// MappingJsonH2HNusare2_act 5.1-5.7: item per slot terisi; mata uang wajib ada di master (D6).
	var mu models.MataUangCsv
	adaItem := false
	for _, sl := range ItemCsvObjek {
		adaItem = adaItem || b.t("ITEM_TYPE"+sl) != ""
	}
	if adaItem {
		m, ok, err := s.csv.MataUang(ctx, b.t("CURRENCY"))
		if err != nil {
			return rencanaCsv{}, nil, nil, err
		}
		if !ok {
			masalah = append(masalah, fmt.Sprintf("baris %d kolom CURRENCY: mata uang %q tidak ada di master CURRENCY", b.n, b.t("CURRENCY")))
		}
		mu = m
	}
	adj := b.t("ADJUSTABLE_PCT")
	for _, sl := range ItemCsvObjek {
		jenisItem := b.t("ITEM_TYPE" + sl)
		if jenisItem == "" {
			continue
		}
		idItem, err := s.csv.IDJenisItemFire(ctx, jenisItem)
		if err != nil {
			return rencanaCsv{}, nil, nil, err
		}
		it := models.ItemObjek{ItemTypeID: idItem, ItemType: jenisItem, TSI: b.angka["TSI_OBJECT_ITEM"+sl], Currency: mu.Currency,
			Note: b.t("INTEREST"), IsAdjustable: adj != "" && adj != "100", PctAdjustOther: b.angka["ADJUSTABLE_PCT"],
			Coverages: []models.CoverageObjek{}}
		for _, p := range perilCsv {
			rate := b.angka["RATE_"+p.kolom]
			if rate == nil || rate.Sign() <= 0 {
				continue
			}
			pakaiCoverage[p.coverage] = p.oldID
			it.Coverages = append(it.Coverages, models.CoverageObjek{OldID: p.oldID, Coverage: p.coverage, CoverageNote: p.catatan,
				CoverageBasis: "1", Rate: rate, NetRate: b.angka["NET_RATE_"+p.kolom], DiscountPercentage: b.angka["DISCOUNT_RATE_"+p.kolom],
				IndemnityPercentage: apd.New(100, 0), LimitOfLiability: b.angka["LOL"], AccumulationCode: b.t("ACCUMULATION_CODE"),
				AccumulationDescription: b.t("ACCUMULATION_NAME")})
		}
		o.Items = append(o.Items, it)
	}
	r.objek = o
	return r, pesan, masalah, nil
}

// rencanakanCsv - langkah 5.4-5.14 + MappingJsonH2HNusare2_act / MappingCoverage tanpa menulis: objek per baris, pesan KodePos,
// galat mata uang (400, D6) dan ID coverage master (400, D8) dikumpulkan.
func (s *Service) rencanakanCsv(ctx context.Context, baris []barisCsv) ([]rencanaCsv, []string, error) {
	var pesan, masalah []string
	var rencana []rencanaCsv
	pakaiCoverage := map[string]string{}
	for _, b := range baris {
		r, p, m, err := s.rencanakanBarisCsv(ctx, b, pakaiCoverage)
		if err != nil {
			return nil, nil, err
		}
		pesan, masalah = append(pesan, p...), append(masalah, m...)
		rencana = append(rencana, r)
	}
	ids := make([]string, 0, len(pakaiCoverage))
	for id := range pakaiCoverage {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	ada, err := s.csv.CoverageAda(ctx, ids)
	if err != nil {
		return nil, nil, err
	}
	for _, id := range ids {
		if !ada[id] {
			masalah = append(masalah, fmt.Sprintf("coverage %s (%s) tidak ada di master COVERAGE", id, pakaiCoverage[id]))
		}
	}
	if len(masalah) > 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return rencana, pesan, nil
}

// cariRWCsv - BrowseRW2_SQL (zip + kelurahan) dengan aturan work owner 08-10-2026 (tiket 66 lanjutan): POSTAL_CODE dulu; tanpa
// hasil dan KODE_POS terisi (berbeda) -> coba KODE_POS. Kode pos yang menemukan RW dipakai lokasi, risk address dan akumulasi.
// Kelurahan kosong -> tidak dicari (Pega: langkah 5.6 dilewati). Hasil > 1 baris -> pilihRWCsv.
func (s *Service) cariRWCsv(ctx context.Context, b barisCsv) (*models.RWCsv, string, error) {
	pos, kode, kel := b.t("POSTAL_CODE"), b.t("KODE_POS"), b.t("KELURAHAN")
	if kel == "" {
		return nil, pos, nil
	}
	var calon []string
	if pos != "" {
		calon = append(calon, pos)
	}
	if kode != "" && kode != pos {
		calon = append(calon, kode)
	}
	for _, z := range calon {
		baris, err := s.csv.RWKelurahan(ctx, z, kel)
		if err != nil {
			return nil, pos, err
		}
		if len(baris) > 0 {
			rw := pilihRWCsv(b, baris)
			return &rw, z, nil
		}
	}
	return nil, pos, nil
}

// pilihRWCsv - "Bantu cari RW" (keputusan work owner; aturan pilih `[keputusan agent]`): > 1 baris -> baris pertama yang cocok
// (tak peka huruf, dipangkas) dengan setiap KECAMATAN / KABUPATEN_KOTA / PROVINSI yang terisi (master RW DISTRICTNAME /
// CITYNAME / PROVINCENAME); tidak ada yang cocok -> baris pertama (perilaku Pega) + log. Satu baris -> dipakai (beda wilayah ->
// log saja).
func pilihRWCsv(b barisCsv, baris []models.RWCsv) models.RWCsv {
	cocok := func(r models.RWCsv) bool {
		for _, w := range []struct{ csv, rw string }{{b.t("KECAMATAN"), r.District}, {b.t("KABUPATEN_KOTA"), r.City},
			{b.t("PROVINSI"), r.Province}} {
			if c := strings.TrimSpace(w.csv); c != "" && !strings.EqualFold(c, strings.TrimSpace(w.rw)) {
				return false
			}
		}
		return true
	}
	if len(baris) == 1 {
		if !cocok(baris[0]) {
			log.Printf("nbfacin: CSV objek baris %d: RW satu-satunya berbeda wilayah dengan KECAMATAN / KABUPATEN_KOTA / PROVINSI", b.n)
		}
		return baris[0]
	}
	for _, r := range baris {
		if cocok(r) {
			return r
		}
	}
	log.Printf("nbfacin: CSV objek baris %d: %d RW tanpa yang cocok wilayah - baris pertama dipakai", b.n, len(baris))
	return baris[0]
}

// tulisMasterCsv - D3 (a): risk address / akumulasi baru (transaksi masing-masing, seperti prosedur Pega yang men-COMMIT
// sendiri) SEBELUM transaksi objek; galat sesudahnya tidak membatalkan master baru. Kode hasil dipasang ke objek / coverage.
func (s *Service) tulisMasterCsv(ctx context.Context, pelaku inti.Pelaku, id string, rencana []rencanaCsv) error {
	for i := range rencana {
		r := &rencana[i]
		o := &r.objek
		if r.rw == nil {
			continue // 5.9.2 / 5.10.2: hanya bila RW ditemukan
		}
		if r.riskID == "" {
			// Cari ulang: baris sebelumnya mungkin sudah menambah alamat yang sama (Pega: Obj-Browse tiap baris sesudah
			// prosedur baris sebelumnya men-COMMIT).
			ada, err := s.csv.IDRiskAddress(ctx, o.RoadName, o.ZipCode)
			if err != nil {
				return err
			}
			r.riskID = ada
		}
		if r.riskID == "" {
			a := models.AlamatBaru{NationName: r.rw.Nation, ProvinceName: r.rw.Province, DistrictName: r.rw.District,
				CityName: r.rw.City, TerritoryName: r.rw.Note, Title: o.RoadType, Address: o.RoadName, PostalCode: o.ZipCode}
			err := s.transaksi(ctx, func(tx *db.Tx) error {
				baru, err := s.rw.SisipRiskAddress(ctx, tx, a)
				r.riskID = baru
				return err
			})
			if err != nil {
				return err
			}
			catatCacheCsv(s.csv, kunciRisk(o.RoadName, o.ZipCode), r.riskID)
		}
		o.RiskAddressID = r.riskID
		nama := o.RoadName // ACCUMULATION_NAME = ADDRESS (langkah 5.2)
		if r.akumulasiID == "" {
			ada, err := s.csv.IDAkumulasi(ctx, nama)
			if err != nil {
				return err
			}
			r.akumulasiID = ada
		}
		if r.akumulasiID == "" {
			kode, err := s.akumulasiBaruCsv(ctx, pelaku, id, *o)
			if err != nil {
				return err
			}
			r.akumulasiID = kode
			if kode != "" {
				catatCacheCsv(s.csv, kunciAkumulasi(nama), kode)
			}
		}
		if r.akumulasiID == "" {
			continue // D4: tidak dibuat -> ACCUMULATION_CODE tetap nilai CSV
		}
		for m := range o.Items {
			for k := range o.Items[m].Coverages {
				o.Items[m].Coverages[k].AccumulationCode = r.akumulasiID
			}
		}
	}
	return nil
}

// akumulasiBaruCsv - langkah 5.10.2 lewat port tiket 46 (TambahAkumulasi): provinsi -> ID, CZone menurut kode pos, negara ->
// NATIONINITIAL. Salah satunya tidak ditemukan, atau master menolak (ganda / isian) -> "" + log, tanpa pesan baru (D4).
func (s *Service) akumulasiBaruCsv(ctx context.Context, pelaku inti.Pelaku, id string, o models.ObjekFire) (string, error) {
	if s.addAkumulasi == nil {
		return "", ErrAkumulasiTanpaDatabase
	}
	prov, err := s.csv.IDProvinsi(ctx, o.Province)
	if err != nil {
		return "", err
	}
	cz, ada, err := s.addAkumulasi.CZoneZip(ctx, o.ZipCode)
	if err != nil {
		return "", err
	}
	negara, err := s.csv.InisialNegara(ctx, o.Country)
	if err != nil {
		return "", err
	}
	if prov == "" || !ada || cz.Code == "" || negara == "" || o.ZipCode == "" {
		log.Printf("nbfacin: CSV objek kasus %s lokasi %s: akumulasi baru tidak dibuat (provinsi %t, CZone %t, negara %t)", id,
			o.ObjectNo, prov != "", ada && cz.Code != "", negara != "")
		return "", nil
	}
	kode, _, err := s.TambahAkumulasi(ctx, pelaku, IsianAkumulasi{Accumulation: akumulasiBaruCsv, AccumulationType: tipeAkumulasiCsv,
		Note: o.RoadName, Keyword: keywordAkumulasi(o.ObjectType), ScopeArea: scopeAkumulasiCsv, CZone: cz.Code, CZoneID: cz.ID,
		ProvinceID: prov, ZipCode: o.ZipCode, Negara: negara})
	switch {
	case errors.Is(err, ErrAkumulasiSudahAda), errors.Is(err, ErrWajibAkumulasi), errors.Is(err, ErrMasukanAkumulasi):
		log.Printf("nbfacin: CSV objek kasus %s lokasi %s: akumulasi baru ditolak master: %v", id, o.ObjectNo, err)
		return "", nil
	case err != nil:
		return "", err
	}
	return kode, nil
}

// HasilCsvObjek - jawaban POST …/objek/csv: daftar objek sesudah diganti (= GET …/objek) dan pesan Pega.
type HasilCsvObjek struct {
	Baris []models.ObjekFire
	Pesan []string
}

// DenganCsvObjek memasang pembaca master Save Data CSV (tiket 66).
func (s *Service) DenganCsvObjek(c repository.PembacaCsvObjek) *Service { s.csv = c; return s }

// SimpanCsvObjek - POST …/objek/csv. Urutan: 401 -> penjaga PUT …/objek (403 / 409) -> 503 -> 404 -> 400 isian CSV (semua
// galat) -> 400 mata uang / coverage master -> 400 objek -> master baru -> SATU transaksi (header General, hapus CurrencyList +
// cicilan, objek diganti, spreading dibuang) -> baca ulang. Pesan = pesan KodePos, atau "Success Insert Data" bila tanpa pesan.
func (s *Service) SimpanCsvObjek(ctx context.Context, pelaku inti.Pelaku, id string, baris []map[string]string) (HasilCsvObjek, error) {
	return s.simpanCsvObjek(ctx, pelaku, id, baris, true)
}

// simpanCsvObjek - pakaiCache false hanya untuk uji kesamaan hasil dengan / tanpa cache.
func (s *Service) simpanCsvObjek(ctx context.Context, pelaku inti.Pelaku, id string, baris []map[string]string, pakaiCache bool) (HasilCsvObjek, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilCsvObjek{}, err
	}
	if err := s.wajibBolehUbah(ctx, pelaku, id); err != nil {
		return HasilCsvObjek{}, err
	}
	if s.objek == nil || s.transaksi == nil || s.csv == nil || s.kasus == nil || s.rw == nil {
		return HasilCsvObjek{}, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return HasilCsvObjek{}, ErrKasusTidakAda
	}
	if err := s.cekKunciUnggah(ctx, id); err != nil { // tiket 75 K75-2
		return HasilCsvObjek{}, err
	}
	csv, err := uraiCsvObjek(baris)
	if err != nil {
		return HasilCsvObjek{}, err
	}
	if pakaiCache {
		s = s.denganCacheCsv()
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return HasilCsvObjek{}, ErrKasusTidakAda
	}
	if err != nil {
		return HasilCsvObjek{}, err
	}
	if err := s.wajibCsvFire(ctx, id); err != nil {
		return HasilCsvObjek{}, err
	}
	g, err := s.generalCsv(kasus, headerDariCsv(id, csv))
	if err != nil {
		return HasilCsvObjek{}, err
	}
	rencana, pesan, err := s.rencanakanCsv(ctx, csv)
	if err != nil {
		return HasilCsvObjek{}, err
	}
	objek := make([]models.ObjekFire, len(rencana))
	for i, r := range rencana {
		objek[i] = r.objek
	}
	if err := periksaObjek(objek); err != nil {
		return HasilCsvObjek{}, err
	}
	if err := s.tulisMasterCsv(ctx, pelaku, id, rencana); err != nil {
		return HasilCsvObjek{}, err
	}
	for i, r := range rencana {
		objek[i] = r.objek
	}
	err = s.gantiObjek(ctx, id, objek, opsiGantiObjek{
		periode:        &Periode{Mulai: g.StartDateTime, Akhir: g.EndDateTime, EDMDay: g.Day},
		buangSpreading: true,
		sebelum: func(tx *db.Tx) error {
			if err := s.kasus.SimpanGeneral(ctx, tx, id, g); err != nil {
				return err
			}
			return s.csv.HapusCurrencyList(ctx, tx, id)
		},
	})
	if err != nil {
		return HasilCsvObjek{}, err
	}
	if len(pesan) == 0 {
		pesan = []string{pesanSuksesCsv}
	}
	hasil, err := s.BacaObjek(ctx, id)
	if err != nil {
		return HasilCsvObjek{}, err
	}
	log.Printf("nbfacin: CSV objek kasus %s oleh %s: %d lokasi, %d pesan", id, pelaku.AkunID, len(objek), len(pesan))
	return HasilCsvObjek{Baris: hasil, Pesan: pesan}, nil
}

// generalCsv - General tersimpan (lewat kabel, seperti Save for later) dengan langkah 5.3 versi D1 (b): No Offer Slip / Begin /
// End dari kolom terisi, OfferingDate = hari ini WIB (@CurrentDateTime), TypeFacultative = typeFacultativeCsv; % Share, RI Comm,
// IsProRate tidak disentuh. Ceding Co dipertahankan. Divalidasi periksaGeneral (lebar / tanggal, 400).
func (s *Service) generalCsv(k models.Kasus, h headerCsv) (models.General, error) {
	kab := s.keKabel(k).General
	isi := IsianGeneral{ReffNumber: kab.ReffNumber, QQName: kab.QQName, BeginDate: kab.BeginDate, EndDate: kab.EndDate,
		PolicyType: kab.PolicyType, MarketingID: kab.MarketingID, Day: kab.Day, SourceOfBusinessID: kab.SourceOfBusinessID,
		OfferingDate: s.sekarang().In(WIB).Format(BentukTanggalKabel), TypeFacultative: typeFacultativeCsv}
	for _, c := range kab.CedingList {
		isi.CedingIDs = append(isi.CedingIDs, c.ID)
	}
	if h.noOfferSlip != "" {
		isi.ReffNumber = h.noOfferSlip
	}
	if !h.mulai.IsZero() {
		isi.BeginDate = h.mulai.Format(BentukTanggalKabel)
	}
	if !h.akhir.IsZero() {
		isi.EndDate = h.akhir.Format(BentukTanggalKabel)
	}
	return periksaGeneral(isi)
}

// ---- cache per permintaan (tiket 66 lanjutan: 3000 baris) -----------------------------------------------------------
//
// Pencarian master berulang (RW per kode pos + kelurahan, risk address / akumulasi per kunci, negara, provinsi, mata uang,
// ItemTypeID, okupasi, CZone) dijawab sekali per permintaan. Hasil SAMA dengan tanpa cache (diuji): galat tidak di-cache, dan
// master yang baru ditambah dicatat ke cache (catatCacheCsv) sehingga baris berikutnya memakainya - seperti pencarian ulang ke
// basis data. Cache hidup selama satu permintaan saja (salinan Service), tidak dibagi antarpermintaan.

func kunciRisk(alamat, zip string) string { return "risk|" + alamat + "|" + zip }
func kunciAkumulasi(nama string) string   { return "akumulasi|" + nama }

type mataUangCache struct {
	m  models.MataUangCsv
	ok bool
}

// cacheCsv - PembacaCsvObjek ber-cache; CoverageAda / HapusCurrencyList diteruskan apa adanya.
type cacheCsv struct {
	repository.PembacaCsvObjek
	rw   map[string][]models.RWCsv
	teks map[string]string
	mu   map[string]mataUangCache
}

func (c *cacheCsv) RWKelurahan(ctx context.Context, zip, kel string) ([]models.RWCsv, error) {
	k := zip + "|" + kel
	if v, ok := c.rw[k]; ok {
		return v, nil
	}
	v, err := c.PembacaCsvObjek.RWKelurahan(ctx, zip, kel)
	if err == nil {
		c.rw[k] = v
	}
	return v, err
}

// RWKodePos - tiket 73 P6; ber-cache per kode pos (tidak ada juga di-cache).
func (c *cacheCsv) RWKodePos(ctx context.Context, zip string) (models.RWCsv, bool, error) {
	k := "kodepos|" + zip
	if v, ok := c.rw[k]; ok {
		if len(v) == 0 {
			return models.RWCsv{}, false, nil
		}
		return v[0], true, nil
	}
	v, ada, err := c.PembacaCsvObjek.RWKodePos(ctx, zip)
	if err != nil {
		return models.RWCsv{}, false, err
	}
	c.rw[k] = nil
	if ada {
		c.rw[k] = []models.RWCsv{v}
	}
	return v, ada, nil
}

func (c *cacheCsv) satu(kunci string, ambil func() (string, error)) (string, error) {
	if v, ok := c.teks[kunci]; ok {
		return v, nil
	}
	v, err := ambil()
	if err == nil {
		c.teks[kunci] = v
	}
	return v, err
}

func (c *cacheCsv) IDRiskAddress(ctx context.Context, alamat, zip string) (string, error) {
	return c.satu(kunciRisk(alamat, zip), func() (string, error) { return c.PembacaCsvObjek.IDRiskAddress(ctx, alamat, zip) })
}

func (c *cacheCsv) IDAkumulasi(ctx context.Context, nama string) (string, error) {
	return c.satu(kunciAkumulasi(nama), func() (string, error) { return c.PembacaCsvObjek.IDAkumulasi(ctx, nama) })
}

func (c *cacheCsv) InisialNegara(ctx context.Context, nama string) (string, error) {
	return c.satu("negara|"+nama, func() (string, error) { return c.PembacaCsvObjek.InisialNegara(ctx, nama) })
}

func (c *cacheCsv) IDProvinsi(ctx context.Context, nama string) (string, error) {
	return c.satu("provinsi|"+nama, func() (string, error) { return c.PembacaCsvObjek.IDProvinsi(ctx, nama) })
}

func (c *cacheCsv) IDJenisItemFire(ctx context.Context, nama string) (string, error) {
	return c.satu("item|"+nama, func() (string, error) { return c.PembacaCsvObjek.IDJenisItemFire(ctx, nama) })
}

func (c *cacheCsv) NamaOkupasiFire(ctx context.Context, kode string) (string, error) {
	return c.satu("okupasi|"+kode, func() (string, error) { return c.PembacaCsvObjek.NamaOkupasiFire(ctx, kode) })
}

func (c *cacheCsv) MataUang(ctx context.Context, nama string) (models.MataUangCsv, bool, error) {
	if v, ok := c.mu[nama]; ok {
		return v.m, v.ok, nil
	}
	m, ok, err := c.PembacaCsvObjek.MataUang(ctx, nama)
	if err == nil {
		c.mu[nama] = mataUangCache{m, ok}
	}
	return m, ok, err
}

// catatCacheCsv - master yang baru ditambah (risk address / akumulasi) masuk cache; tanpa cache tidak berbuat apa-apa
// (pencarian ulang ke basis data menemukannya karena tambahan master sudah di-COMMIT, D3).
func catatCacheCsv(p repository.PembacaCsvObjek, kunci, nilai string) {
	if c, ok := p.(*cacheCsv); ok {
		c.teks[kunci] = nilai
	}
}

type czoneCache struct {
	c   models.CZoneZip
	ada bool
}

// cacheAddAkumulasi - CZoneZip ber-cache per kode pos; CariZip diteruskan.
type cacheAddAkumulasi struct {
	repository.PembacaAddAkumulasi
	cz map[string]czoneCache
}

func (c *cacheAddAkumulasi) CZoneZip(ctx context.Context, zip string) (models.CZoneZip, bool, error) {
	if v, ok := c.cz[zip]; ok {
		return v.c, v.ada, nil
	}
	v, ada, err := c.PembacaAddAkumulasi.CZoneZip(ctx, zip)
	if err == nil {
		c.cz[zip] = czoneCache{v, ada}
	}
	return v, ada, err
}

// denganCacheCsv - salinan Service untuk SATU permintaan Save Data CSV dengan pencarian master ber-cache.
func (s *Service) denganCacheCsv() *Service {
	salin := *s
	salin.csv = &cacheCsv{PembacaCsvObjek: s.csv, rw: map[string][]models.RWCsv{}, teks: map[string]string{},
		mu: map[string]mataUangCache{}}
	if s.addAkumulasi != nil {
		salin.addAkumulasi = &cacheAddAkumulasi{PembacaAddAkumulasi: s.addAkumulasi, cz: map[string]czoneCache{}}
	}
	return &salin
}
