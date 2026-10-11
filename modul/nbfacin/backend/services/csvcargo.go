package services

// Unggah CSV MARINE CARGO (tiket 74 P6, keputusan work owner 11-10-2026 K74-14..17) di atas infrastruktur tiket 75 (berkas
// mengalir -> penampung, validasi latar, simpan latar per kelompok, Ulangi, kunci case). `[terverifikasi]` NB FacIn:
//   - ChooseActivityMappingCSV_Act L539 / L599: IsMarineCargo -> InsertUploadMarineCargo_act (jenis pekerjaan "cargo").
//   - Judul dibaca per NAMA (pxUploadCSVResults), urutan berkas bebas: models.KolomCsvCargo (18).
//   - InsertUploadMarineCargo_act: 1 Page-Remove, 2 Property-Remove CargoList / CurrencyList (GANTI semua cargo case); 3 per baris
//     ListB2BHostUpload: 3.1 DEPARTURE / DESTINATION @toUpperCase; TSI berkoma -> "." dibuang lalu "," -> "." (L1136-1383);
//     CURRENCY "RMB" -> "CNY"; 3.2 Call ProtectUploadMarineCargo_act lalu TRANSISI `TempError.ERRMSG==""` salah -> kode 6 Exit
//     Activity (satu baris bermasalah = seluruh unggahan tidak tersimpan); 3.3-3.8 ber-label `//` (MATI: TradingList / GoodsList /
//     ConveyanceList / Ship case / InvoiceNumber case TIDAK diisi); 3.9 PRAKONDISI `TRADING, GOODS, CONVEYENCE, PACKING baris ini ==
//     baris sebelumnya` T=3 lewati / F=2 jalan: cargo BARU hanya bila salah satu berbeda - Trading/Goods/Conveyance/Packing ID + Note,
//     InvoiceNumber, SailDate (dari ETD), BLNumber, IsTopRisk (TOP_RISK = "TRUE"); 3.10 / 3.11 / 3.12 TANPA prakondisi pada cargo
//     TERAKHIR: Ship + AdditionalShip(1) = Ship (BrowseShip2_RD pxResults(1)), FromRute / ToRute, SurveyAgent (BrowseLLOYDAGENT_RD),
//     CoverageList = coverage polis master (tanpa master = kosong) lalu CoverageList(1) TSI / Currency Name / ID / OldID
//     (GetCurrencyIDByName). Baris berturutan berkunci sama MENIMPA kapal / rute / agen / coverage cargo yang sama (TSI baris
//     sebelumnya hilang) - perilaku Pega dipertahankan (K74-19, keputusan work owner 11-10-2026 "ikut Pega").
//   - ProtectUploadMarineCargo_act per baris (urutan dan transisi): 1 Trading (BrowseTrading_RD `.Description =`) tidak ada ->
//     "* Trading X Not Available in Master", Exit Activity; 2 Goods (BrowseGoodsType) -> pesan, Exit; 3 PACKING (BrowsePacking_RD) ->
//     pesan, LANJUT (transisinya memeriksa GoodsTemp); 4 Conveyence (BrowseConveyance_RD) -> pesan, Exit; 5 Agent bila AGENT terisi
//     (BrowseLLOYDAGENT_RD `.ID =`; pesan 5.4 MATI - agen tak ada = kosong); 6 BrowseShip_RD `.ID = UPPER(DATASHIP)` ADA ->
//     ERRMSG DIGANTI "Jangan gunakan ID kapal dikolom DATASHIP pada CSV" (DATASHIP kosong = filter diabaikan = semua kapal = pesan
//     itu); 7 BrowseShip2_RD `.Name = UPPER(DATASHIP)` tidak ada -> daftar kapal (RDBMASTERSHIP, repository/kapalmaster.go).
//   - Langkah 5 membangun CurrencyList dari coverage: TIDAK - K66 D5 (CurrencyList + cicilan dihapus tanpa dibangun ulang, sama
//     dengan CSV FIRE) - K74-21, keputusan work owner 11-10-2026.
// K74-18 (keputusan work owner 11-10-2026, template Pega): judul wajib = 17 kolom template tanpa INVOICEDATE; COVERAGE_NAME /
// INVOICEDATE opsional. AGENT terisi selagi master agen (LLOYDAGENT, DDL menunggu) belum ada -> validasi berhenti di baris itu
// dengan galat jelas (gagal_validasi) `[keputusan agent]`, bukan ditebak kosong.
// Keputusan work owner 11-10-2026 K74-20 / K74-21 (penyimpangan sadar / ikut CSV FIRE): TRADING / GOODS / PACKING / CONVEYENCE kosong = tidak
// ada di master (Pega: filter kosong diabaikan -> baris master pertama sembarang); mata uang tak ada di master CURRENCY = galat (K66
// D6 "Ketat"); TSI / ETD wajib terurai (K66 D7); AGENT kosong = tanpa agen (Pega memakai hasil agen baris SEBELUMNYA karena
// AgentTemp tidak dihapus per baris); kapal baru didaftarkan saat SIMPAN (Pega: di Protect, juga untuk baris yang gagal);
// SailDate disimpan bentuk tanggal kita (tanggalKerugianKePega), bukan "yyyyMMddT170000.000 GMT" Pega.

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

const (
	// pesanKapalIDCsv - ProtectUploadMarineCargo_act 6.4 verbatim.
	pesanKapalIDCsv = "Jangan gunakan ID kapal dikolom DATASHIP pada CSV"
	// pesanMasterCsv - "\n * <Jenis> <nilai> Not Available in Master" ProtectUploadMarineCargo_act 1.4 / 2.4 / 3.4 / 4.4 (tanpa "\n").
	pesanMasterCsv = "* %s %s Not Available in Master"
)

// petaJudulCsv - pemetaan judul berkas: idx[i] = indeks kolom penampung untuk sel ke-i (-1 = diabaikan); kolom = judul dikenali
// urut berkas (nil = urutan template); n = jumlah kolom penampung.
type petaJudulCsv struct {
	idx   []int
	kolom []string
	n     int
}

// petaJudulFire - template FIRE (judulCsvSah, posisi tetap).
func petaJudulFire(baris string) (petaJudulCsv, error) {
	if !judulCsvSah(baris) {
		return petaJudulCsv{}, fmt.Errorf("%w: %s", ErrMasukanObjek, pesanJudulCsv)
	}
	p := petaJudulCsv{idx: make([]int, len(models.KolomCsvObjek)), n: len(models.KolomCsvObjek)}
	for i := range p.idx {
		p.idx[i] = i
	}
	return p, nil
}

// kolomCsvCargoOpsional - K74-18 (keputusan work owner 11-10-2026, template Pega `templatecsvmarinecargo.csv`): judul wajib = 17
// kolom template tanpa INVOICEDATE; COVERAGE_NAME (dibaca Insert, tidak ada di template) opsional = kosong. INVOICEDATE (di template,
// TIDAK dibaca InsertUploadMarineCargo_act / ProtectUploadMarineCargo_act) diterima dan diabaikan seperti kolom tak dikenal lain -
// tidak disimpan di penampung, tidak tampil di View Upload.
var kolomCsvCargoOpsional = []string{"COVERAGE_NAME"}

// petaJudulCargo - judul cargo per NAMA (tak peka huruf / spasi tepi), urutan bebas; nama tak dikenal diabaikan; ganda / wajib
// kurang -> ErrMasukanObjek.
func petaJudulCargo(baris string) (petaJudulCsv, error) {
	sel := pecahBarisCsv(baris, ';')
	p := petaJudulCsv{idx: make([]int, len(sel)), n: len(models.KolomCsvCargo)}
	ada := map[string]bool{}
	for i, s := range sel {
		k := strings.ToUpper(strings.TrimSpace(s))
		p.idx[i] = slices.Index(models.KolomCsvCargo, k)
		if p.idx[i] < 0 {
			continue
		}
		if ada[k] {
			return petaJudulCsv{}, fmt.Errorf("%w: judul berkas CSV Marine Cargo memuat kolom %s dua kali", ErrMasukanObjek, k)
		}
		ada[k] = true
		p.kolom = append(p.kolom, k)
	}
	var kurang []string
	for _, k := range models.KolomCsvCargo {
		if !ada[k] && !slices.Contains(kolomCsvCargoOpsional, k) {
			kurang = append(kurang, k)
		}
	}
	if len(kurang) > 0 {
		return petaJudulCsv{}, fmt.Errorf("%w: judul berkas CSV tidak sesuai template Marine Cargo - kolom %s tidak ada (judul dibaca per "+
			"nama, pemisah \";\")", ErrMasukanObjek, strings.Join(kurang, ", "))
	}
	return p, nil
}

// petaJudulJenis - pemeta judul importer `jenis`.
func petaJudulJenis(jenis string) func(string) (petaJudulCsv, error) {
	switch jenis {
	case models.JenisCsvCargo:
		return petaJudulCargo
	case models.JenisCsvAneka: // tiket 73 P6
		return petaJudulAneka
	}
	return petaJudulFire
}

// selDataPeta - sel baris data dipangkas ke kolom penampung menurut peta; sel apa pun yang berisi (termasuk yang diabaikan)
// menandai baris berisi.
func selDataPeta(baris string, p petaJudulCsv) ([]string, bool) {
	sel := pecahBarisCsv(baris, ';')
	isi := make([]string, p.n)
	ada := false
	for i, s := range sel {
		v := strings.TrimSpace(s)
		ada = ada || v != ""
		if i < len(p.idx) && p.idx[i] >= 0 {
			isi[p.idx[i]] = v
		}
	}
	return isi, ada
}

// DenganUnggahCsvCargo memasang penulis cargo hasil unggah dan master kapal (tiket 74 P6).
func (s *Service) DenganUnggahCsvCargo(p repository.PenulisCargoCsv, k repository.PenyimpanKapalMaster) *Service {
	s.penulisCargoCsv, s.kapalMaster = p, k
	return s
}

// ---- satu baris ----------------------------------------------------------------------------------------------------

// barisCargoCsv - satu baris CSV cargo terurai (langkah 3.1) dan hasil master (Protect).
type barisCargoCsv struct {
	no                                  int
	trading, goods, packing, conveyence string // teks CSV (kunci 3.9 + Note)
	departure, destination              string // UPPER
	tsi                                 *apd.Decimal
	mataUang                            string // CURRENCY sesudah RMB -> CNY
	agen, kapal                         string // AGENT apa adanya, UPPER(DATASHIP)
	invoice, sailDate, blNumber         string
	topRisk                             bool
	// hasil master (periksaBarisCargoCsv)
	tradingID, goodID, packingID, conveyanceID string
	mu                                         models.MataUangCsv
	agenSurvei                                 models.AgenSurvei
}

// kunci - prakondisi 3.9: TRADING, GOODS, CONVEYENCE, PACKING.
func (b barisCargoCsv) kunci() [4]string {
	return [4]string{b.trading, b.goods, b.conveyence, b.packing}
}

// polaEtdCsv - ETD: @substring(.ETD,"6","10") + @substring(.ETD,"3","5") + @substring(.ETD,"0","02") -> dd?MM?yyyy.
var polaEtdCsv = regexp.MustCompile(`^(\d{2})\D(\d{2})\D(\d{4})$`)

// lebarCargoCsv - kolom CSV yang ditulis ke kolom bertepi (migrasi 256: T_CARGOLIST / T_COVERAGELIST.CURRENCY_CODE).
var lebarCargoCsv = []struct {
	kolom string
	n     int
}{{"TRADING", 500}, {"GOODS", 500}, {"PACKING", 500}, {"CONVEYENCE", 500}, {"DEPARTURE", 50}, {"DESTINATION", 50},
	{"INVOICE", 50}, {"BLNUMBER", 4000}, {"CURRENCY", 10}, {"DATASHIP", 1000}, {"AGENT", 100}}

// uraiBarisCargoCsv - langkah 3.1 + lebar + angka / tanggal (K66 D7); galat dikembalikan, tidak dihentikan.
func uraiBarisCargoCsv(no int, isi []string) (barisCargoCsv, []string) {
	v := func(k string) string {
		if i := slices.Index(models.KolomCsvCargo, k); i < len(isi) {
			return isi[i]
		}
		return ""
	}
	var masalah []string
	for _, l := range lebarCargoCsv {
		if len(v(l.kolom)) > l.n {
			masalah = append(masalah, fmt.Sprintf("baris %d kolom %s: paling banyak %d byte", no, l.kolom, l.n))
		}
	}
	b := barisCargoCsv{no: no, trading: v("TRADING"), goods: v("GOODS"), packing: v("PACKING"), conveyence: v("CONVEYENCE"),
		departure: strings.ToUpper(v("DEPARTURE")), destination: strings.ToUpper(v("DESTINATION")), mataUang: v("CURRENCY"),
		agen: v("AGENT"), kapal: strings.ToUpper(v("DATASHIP")), invoice: v("INVOICE"), blNumber: v("BLNUMBER"),
		topRisk: v("TOP_RISK") == "TRUE"}
	if b.mataUang == "RMB" {
		b.mataUang = "CNY"
	}
	if s := v("TSI"); s != "" {
		if strings.Contains(s, ",") {
			s = normalAngkaCsv(s)
		}
		d, err := utils.ParseDecimal(s)
		switch {
		case err != nil:
			masalah = append(masalah, fmt.Sprintf("baris %d kolom TSI: angka tidak sah %q", no, v("TSI")))
		case !repository.PolaAngkaMasuk.MatchString(utils.FormatDecimal(simpan8(d))):
			masalah = append(masalah, fmt.Sprintf("baris %d kolom TSI: melebihi 30 digit bulat", no))
		default:
			b.tsi = simpan8(d)
		}
	}
	if s := v("ETD"); s != "" {
		m := polaEtdCsv.FindStringSubmatch(s)
		var t time.Time
		ok := m != nil
		if ok {
			var err error
			t, err = time.ParseInLocation("02-01-2006", m[1]+"-"+m[2]+"-"+m[3], WIB)
			ok = err == nil
		}
		if !ok {
			masalah = append(masalah, fmt.Sprintf("baris %d kolom ETD: tanggal harus dd/MM/yyyy: %q", no, s))
		} else {
			b.sailDate = tanggalKerugianKePega(t.Format(BentukTanggalKabel))
		}
	}
	return b, masalah
}

// cacheCargoCsv - hasil master per pekerjaan (validasi / simpan).
type cacheCargoCsv struct {
	master  map[[2]string]string // [jenis, deskripsi] -> ID ("" = tidak ada)
	agen    map[string]*models.AgenSurvei
	kapalID map[string]bool
	kapal   map[string]models.KapalCargo
	mu      map[string]*models.MataUangCsv
}

func cacheCargoCsvBaru() *cacheCargoCsv {
	return &cacheCargoCsv{master: map[[2]string]string{}, agen: map[string]*models.AgenSurvei{}, kapalID: map[string]bool{},
		kapal: map[string]models.KapalCargo{}, mu: map[string]*models.MataUangCsv{}}
}

// periksaBarisCargoCsv - master satu baris: mata uang (K66 D6) lalu ProtectUploadMarineCargo_act berurutan (lihat kepala berkas).
// Mengisi ID master / mata uang / agen pada b. Galat basis data = error.
func (s *Service) periksaBarisCargoCsv(ctx context.Context, b *barisCargoCsv, c *cacheCargoCsv) ([]string, error) {
	var masalah []string
	mu, ada := c.mu[b.mataUang]
	if !ada {
		if b.mataUang != "" {
			m, ok, err := s.objekCargo.MataUangCargo(ctx, b.mataUang)
			if err != nil {
				return nil, err
			}
			if ok {
				mu = &m
			}
		}
		c.mu[b.mataUang] = mu
	}
	if mu == nil {
		masalah = append(masalah, fmt.Sprintf("baris %d kolom CURRENCY: mata uang %q tidak ada di master CURRENCY", b.no, b.mataUang))
	} else {
		b.mu = *mu
	}
	protect, err := s.protectBarisCargoCsv(ctx, b, c)
	if err != nil {
		return nil, err
	}
	masalah = append(masalah, protect...)
	return append(masalah, lebarMasterCargoCsv(*b)...), nil
}

// lebarMasterCargoCsv - nilai MASTER yang ditulis cargo CSV melebihi kolom tujuan (migrasi 256, VARCHAR2 BYTE = len UTF-8) -> galat
// jelas per baris (keputusan work owner 11-10-2026, audit lebar sesudah bug ORA-12899; tanpa migrasi; batas = validasi layar
// PUT objek-cargo). Master TRADING / GOODSTYPE / PACKING / CONVEYANCE / LLOYDAGENT = view atas JSON (lebar tidak tertulis di
// DDL). Warranty Extension tidak ditulis CSV cargo (K74-15: coverage hanya TSI + Currency) - tidak ada nilai master untuk diperiksa.
func lebarMasterCargoCsv(b barisCargoCsv) []string {
	var masalah []string
	for _, l := range []struct {
		medan, nilai, kolom string
		n                   int
	}{{"ID Trading", b.tradingID, "T_CARGOLIST.TRADING_ID", 50}, {"ID Goods", b.goodID, "T_CARGOLIST.GOOD_ID", 50},
		{"ID Packing", b.packingID, "T_CARGOLIST.PACKING_ID", 50}, {"ID Conveyance", b.conveyanceID, "T_CARGOLIST.CONVEYANCE_ID", 50},
		{"Address agen", b.agenSurvei.Address, "T_SURVEYAGENT.ADDRESS", 500}, {"City agen", b.agenSurvei.City, "T_SURVEYAGENT.CITY", 500},
		{"Country agen", b.agenSurvei.Country, "T_SURVEYAGENT.COUNTRY", 500}} {
		if len(l.nilai) > l.n {
			masalah = append(masalah, fmt.Sprintf("baris %d: master %s %d byte melebihi kolom %s (paling banyak %d byte)", b.no, l.medan,
				len(l.nilai), l.kolom, l.n))
		}
	}
	return masalah
}

// protectBarisCargoCsv - ProtectUploadMarineCargo_act langkah 1-6 (tanpa pendaftaran kapal - itu saat simpan).
func (s *Service) protectBarisCargoCsv(ctx context.Context, b *barisCargoCsv, c *cacheCargoCsv) ([]string, error) {
	cari := func(jenis, nilai string) (string, error) {
		if nilai == "" {
			return "", nil
		}
		k := [2]string{jenis, nilai}
		if id, ada := c.master[k]; ada {
			return id, nil
		}
		id, ok, err := s.masterCargo.IDMaster(ctx, jenis, nilai)
		if err != nil {
			return "", err
		}
		if !ok {
			id = ""
		}
		c.master[k] = id
		return id, nil
	}
	var protect []string
	for _, l := range []struct {
		jenis, label, nilai string
		ke                  *string
		keluar              bool
	}{{MasterTrading, "Trading", b.trading, &b.tradingID, true}, {MasterGoods, "Goods", b.goods, &b.goodID, true},
		{MasterPacking, "PACKING", b.packing, &b.packingID, false}, {MasterConveyance, "Conveyence", b.conveyence, &b.conveyanceID, true}} {
		id, err := cari(l.jenis, l.nilai)
		if err != nil {
			return nil, err
		}
		if id == "" {
			protect = append(protect, fmt.Sprintf(pesanMasterCsv, l.label, l.nilai))
			if l.keluar {
				return protect, nil
			}
			continue
		}
		*l.ke = id
	}
	if b.agen != "" {
		a, ada := c.agen[b.agen]
		if !ada {
			x, ok, err := s.masterCargo.Agen(ctx, b.agen)
			if errors.Is(err, repository.ErrMasterBelumAda) {
				return nil, fmt.Errorf("%w: baris %d kolom AGENT %q: master agen (LLOYDAGENT) belum tersedia di sistem baru - "+
					"kosongkan kolom AGENT atau tunggu master agen dipasang", errAgenBelumAda, b.no, b.agen)
			}
			if err != nil {
				return nil, err
			}
			if ok {
				a = &x
			}
			c.agen[b.agen] = a
		}
		if a != nil {
			b.agenSurvei = *a
		}
	}
	kapalID := b.kapal == ""
	if !kapalID {
		ada, sudah := c.kapalID[b.kapal]
		if !sudah {
			var err error
			if ada, err = s.kapalMaster.AdaKapalID(ctx, b.kapal); err != nil {
				return nil, err
			}
			c.kapalID[b.kapal] = ada
		}
		kapalID = ada
	}
	if kapalID {
		protect = []string{pesanKapalIDCsv} // 6.4: ERRMSG diganti, bukan ditambah
	}
	return protect, nil
}

// kapalCsv - kapal baris (BrowseShip2_RD pxResults(1)); tidak ada -> didaftarkan (RDBMASTERSHIP, transaksi sendiri - Pega
// men-COMMIT di prosedur).
func (s *Service) kapalCsv(ctx context.Context, pelaku inti.Pelaku, nama string, c *cacheCargoCsv) (models.KapalCargo, error) {
	if k, ada := c.kapal[nama]; ada {
		return k, nil
	}
	k, ok, err := s.kapalMaster.KapalDariNama(ctx, nama)
	if err != nil {
		return models.KapalCargo{}, err
	}
	if !ok {
		err = s.transaksi(ctx, func(tx *db.Tx) error {
			var err error
			k, err = s.kapalMaster.DaftarKapal(ctx, tx, nama, pelaku.AkunID)
			return err
		})
		if err != nil {
			return models.KapalCargo{}, err
		}
	}
	c.kapal[nama] = k
	return k, nil
}

// cargoBaruCsv - 3.9 Set Object.
func cargoBaruCsv(b barisCargoCsv) models.ObjekCargo {
	return models.ObjekCargo{TradingID: b.tradingID, TradingNote: b.trading, GoodID: b.goodID, GoodNote: b.goods,
		ConveyanceID: b.conveyanceID, ConveyanceNote: b.conveyence, PackingID: b.packingID, PackingNote: b.packing,
		InvoiceNumber: b.invoice, SailDate: b.sailDate, BLNumber: b.blNumber, IsTopRisk: b.topRisk}
}

// terapkanBarisCargoCsv - 3.10 kapal + AdditionalShip(1) + rute, 3.11 agen, 3.12.6 coverage (tanpa polis master) atas cargo
// terakhir.
func terapkanBarisCargoCsv(c *models.ObjekCargo, b barisCargoCsv, k models.KapalCargo) {
	k.AdditionalShip = []models.KapalTambahan{{ID: k.ID, Register: k.Register, NmShip: k.NmShip, ShipTypeID: k.ShipTypeID,
		FlagName: k.FlagName, NRT: k.NRT, DWT: k.DWT, GRT: k.GRT, YMake1: k.YMake1, Construction: k.Construction}}
	c.Ship = k
	c.FromRute, c.ToRute = b.departure, b.destination
	c.SurveyAgent = b.agenSurvei
	c.Coverages = []models.CoverageCargo{{TSI: b.tsi, Currency: b.mataUang, CurrencyID: b.mu.ID, CurrencyOldID: b.mu.OldID}}
}

// ---- validasi & simpan latar ------------------------------------------------------------------------------------------

// errAgenBelumAda - AGENT terisi, master agen belum terpasang: validasi berhenti (galat baris, bukan galat proses).
var errAgenBelumAda = errors.New("master agen belum tersedia")

// validasiKelompokCargo - satu kelompok baris: galat bernomor baris; berhenti = errAgenBelumAda di baris terakhir galat. Galat basis
// data = error.
func (s *Service) validasiKelompokCargo(ctx context.Context, baris []models.BarisUnggahCsv, c *cacheCargoCsv) ([]models.PesanUnggahCsv, bool, error) {
	var galat []models.PesanUnggahCsv
	for _, x := range baris {
		b, masalah := uraiBarisCargoCsv(x.No, x.Isi)
		if len(masalah) == 0 {
			var err error
			masalah, err = s.periksaBarisCargoCsv(ctx, &b, c)
			if errors.Is(err, errAgenBelumAda) {
				return append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: strings.TrimPrefix(err.Error(), errAgenBelumAda.Error()+": ")}),
					true, nil
			}
			if err != nil {
				return nil, false, err
			}
		}
		for _, m := range masalah {
			galat = append(galat, models.PesanUnggahCsv{NoBaris: x.No, Teks: m})
		}
	}
	return galat, false, nil
}

// validasiUnggahCargo - pekerjaan latar fase validasi jenis cargo: semua kelompok; satu galat = gagal_validasi (Pega Exit Activity).
func (s *Service) validasiUnggahCargo(ctx context.Context, u models.UnggahCsv) {
	c := cacheCargoCsvBaru()
	for dari := 1; dari <= u.TotalBaris; dari += UkuranKelompokCsv {
		baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, dari, min(dari+UkuranKelompokCsv-1, u.TotalBaris))
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		galat, berhenti, err := s.validasiKelompokCargo(ctx, baris, c)
		if err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}

		urut := u.JumlahGalat + 1
		u.BarisDiproses, u.JumlahGalat = dari+len(baris)-1, u.JumlahGalat+len(galat)
		if berhenti {
			u.BarisDiproses = galat[len(galat)-1].NoBaris
		}
		if err := s.tulisUnggah(ctx, u, func(tx *db.Tx) error {
			return s.unggahCsv.SisipPesanUnggah(ctx, tx, u.ID, jenisGalatUnggah, urut, galat)
		}); err != nil {
			s.gagalUnggah(ctx, u, dari, err)
			return
		}
		if berhenti {
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

// simpanUnggahCargo - pekerjaan latar fase simpan jenis cargo dari baris BARIS_TERSIMPAN+1. Satu kelompok = paling sedikit
// UkuranKelompokCsv baris DAN berakhir di batas cargo (baris berikutnya membuka cargo baru, prakondisi 3.9) - kelompok tidak
// pernah memotong cargo, jadi Ulangi selalu mulai di cargo baru. Kelompok pertama: GANTI cargo case + hapus CurrencyList
// (langkah 2); kemajuan ditulis di transaksi cargo kelompoknya.
func (s *Service) simpanUnggahCargo(ctx context.Context, pelaku inti.Pelaku, u models.UnggahCsv) {
	c := cacheCargoCsvBaru()
	for dari := u.BarisTersimpan + 1; dari <= u.TotalBaris; {
		var cargo []models.ObjekCargo
		var kunciAkhir [4]string
		akhir, jumlah := dari-1, 0
	baca:
		for p := dari; p <= u.TotalBaris; p += UkuranKelompokCsv {
			baris, err := s.unggahCsv.BacaBarisUnggah(ctx, u.ID, u.Jenis, p, min(p+UkuranKelompokCsv-1, u.TotalBaris))
			if err != nil {
				s.gagalUnggah(ctx, u, dari, err)
				return
			}
			for _, x := range baris {
				b, masalah := uraiBarisCargoCsv(x.No, x.Isi)
				if len(masalah) == 0 {
					if masalah, err = s.periksaBarisCargoCsv(ctx, &b, c); err != nil {
						s.gagalUnggah(ctx, u, x.No, err)
						return
					}
				}
				if len(masalah) > 0 {
					s.gagalUnggah(ctx, u, x.No, fmt.Errorf("%w: data master berubah sejak validasi: %s", ErrMasukanObjek, masalah[0]))
					return
				}
				baru := len(cargo) == 0 || b.kunci() != kunciAkhir
				if baru && jumlah >= UkuranKelompokCsv {
					break baca
				}
				k, err := s.kapalCsv(ctx, pelaku, b.kapal, c)
				if err != nil {
					s.gagalUnggah(ctx, u, x.No, err)
					return
				}
				if baru {
					cargo = append(cargo, cargoBaruCsv(b))
				}
				terapkanBarisCargoCsv(&cargo[len(cargo)-1], b, k)
				kunciAkhir, akhir, jumlah = b.kunci(), x.No, jumlah+1
			}
		}
		maju := u
		maju.KelompokTerakhir, maju.BarisTersimpan = u.KelompokTerakhir+1, akhir
		ganti := u.BarisTersimpan == 0
		err := s.transaksi(ctx, func(tx *db.Tx) error {
			if err := s.penulisCargoCsv.TulisCargoCsv(ctx, tx, u.CaseID, cargo, ganti); err != nil {
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
