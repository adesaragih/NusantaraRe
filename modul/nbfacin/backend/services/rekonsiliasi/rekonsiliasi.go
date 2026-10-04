// Package rekonsiliasi adalah kerangka rekonsiliasi eksak TAHAP 1 (tiket NB-16,
// ADR-F-0001): kasus terekam dijalankan lewat perhitungan sistem baru dan
// dibandingkan sampai digit terakhir dengan nilai sistem lama.
//
// ⛔ Tanpa toleransi: premi coverage dibandingkan sebagai teks desimal
// (`utils.FormatDecimal`), agregat per mata uang sebagai nilai desimal eksak (A31).
// Tidak ada epsilon di paket ini.
// ⛔ Hanya fixture ter-de-identifikasi: BerkasKasus menolak berkas yang belum
// lolos `deidentifikasi.SudahBersih`; fungsi menerima ISI berkas, bukan jalur,
// jadi paket ini tidak pernah membuka berkas sendiri.
//
// Tahap 1 = perhitungan murni atas masukan terekam. Prasyarat tahap berikutnya:
//   - jalur `Param.PremiStatus == "amount"` FIRE dan pro-rata dari tanggal (tiket 18
//     "Di luar"): coverage seperti itu terhitung dengan rumus percent dan muncul
//     sebagai selisih, bukan belum tercakup - `PremiStatus` tidak terekam di fixture;
//   - pemetaan coverage berkas kasus → premium.Input untuk PA dan MBU (belum ada
//     fixture P-5 PA/MBU untuk mengujinya);
//   - nilai tangga sistem lama (DataSearch.CARID2, LetterNo, PositionNote) terekam
//     per kasus - fixture P-5 tidak memuatnya;
//   - jalur tulis produksi (per-modul lalu end-to-end).
package rekonsiliasi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/alat/deidentifikasi"
	"nusantarare/modul/nbfacin/backend/services/pembayaran"
	"nusantarare/modul/nbfacin/backend/services/premium"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// Status - hasil satu perbandingan.
type Status string

const (
	Cocok         Status = "cocok"
	Selisih       Status = "selisih"
	BelumTercakup Status = "belum tercakup" // bukan lulus: belum ada yang dibandingkan
	Galat         Status = "galat"
)

// Hasil - satu baris laporan.
type Hasil struct {
	Kasus   string
	Lini    string
	Langkah string // rule + langkah sistem lama tempat nilai dibentuk
	Status  Status
	// SistemLama, SistemBaru - teks desimal apa adanya; kosong bila tidak dibandingkan.
	SistemLama string
	SistemBaru string
	Catatan    string
}

// Langkah untuk baris laporan.
const (
	langkahPremi      = "premi coverage"
	langkahNilaiDasar = "nilai dasar akseptasi (DataSearch.CARID2)"
	langkahAgregatMBU = "FillPremiMBU_FacIn langkah 2.6 CurrencyList.Premium"
	langkahPembayaran = "PremiPaymentMarine langkah 1.1 Policy.Payment.Premium"
	langkahTotal      = "SumTotalTSIPremiGross_Act langkah 1"
)

// DaftarIzin membandingkan fixture daftar-izin premi (`premium/testdata/daftarizin`):
// larik {kasus, input premium.Input, premi}. Urutan hasil = urutan masukan.
func DaftarIzin(isi []byte) ([]Hasil, error) {
	var kasus []struct {
		Kasus string
		Input premium.Input
		Premi string
	}
	if err := json.Unmarshal(isi, &kasus); err != nil {
		return nil, fmt.Errorf("rekonsiliasi: daftar-izin: %w", err)
	}
	hasil := make([]Hasil, 0, len(kasus))
	for _, k := range kasus {
		hasil = append(hasil, bandingkanPremi(k.Kasus, k.Input, k.Premi))
	}
	return hasil, nil
}

func bandingkanPremi(kasus string, in premium.Input, lama string) (h Hasil) {
	h = Hasil{Kasus: kasus, Lini: string(in.LiniBisnis), Langkah: premium.AsalRumus(in), SistemLama: lama}
	defer func() {
		// Lini di luar peta K-018 → panic resolver; satu kasus, bukan seluruh laporan.
		if r := recover(); r != nil {
			h.Status, h.Catatan = Galat, fmt.Sprint(r)
		}
	}()
	baru, err := premium.Calculate(in)
	switch {
	case errors.Is(err, premium.ErrBentukBelumDiport):
		h.Status, h.Catatan = BelumTercakup, err.Error()
	case err != nil:
		h.Status, h.Catatan = Galat, err.Error()
	default:
		h.SistemBaru = utils.FormatDecimal(baru.Amount)
		h.Status = Selisih
		if h.SistemBaru == lama {
			h.Status = Cocok
		}
	}
	return h
}

// AgregatMBU membandingkan premi per mata uang MBU (tiket 07,
// `premium/testdata/daftarizin/mbu_mata_uang.json`): satu baris per mata uang per kasus.
// Kesamaan NILAI desimal eksak, bukan teks (A31): CurrencyList.Premium bertipe Decimal
// (butir 52) dan tersimpan tanpa nol di belakang koma. Mata uang yang hanya ada di
// salah satu sisi = selisih.
func AgregatMBU(isi []byte) ([]Hasil, error) {
	var kasus []struct {
		Kasus                string
		UrutanMasterMataUang []string
		Coverage             []struct{ TSI, Rate, Loading, ProRatePercentCoverage, MataUang, FlagDelete string }
		CurrencyList         []struct{ Name, Premium string }
	}
	if err := json.Unmarshal(isi, &kasus); err != nil {
		return nil, fmt.Errorf("rekonsiliasi: agregat MBU: %w", err)
	}
	var hasil []Hasil
	for _, k := range kasus {
		var cov []premium.CoverageMBU
		for _, c := range k.Coverage {
			cov = append(cov, premium.CoverageMBU{FlagDelete: c.FlagDelete, Input: premium.Input{LiniBisnis: premium.LiniMBU,
				MataUang: c.MataUang, TSI: c.TSI, Rate: c.Rate, Loading: c.Loading, ProRatePercentCoverage: c.ProRatePercentCoverage}})
		}
		baru, err := premium.PremiPerMataUang(cov, k.UrutanMasterMataUang)
		if err != nil {
			hasil = append(hasil, Hasil{Kasus: k.Kasus, Lini: string(premium.LiniMBU), Langkah: langkahAgregatMBU, Status: Galat, Catatan: err.Error()})
			continue
		}
		nilaiBaru := map[string]string{}
		var urut []string
		for _, b := range baru {
			nilaiBaru[b.MataUang] = utils.FormatDecimal(b.Premi.Amount)
			urut = append(urut, b.MataUang)
		}
		nilaiLama := map[string]string{}
		for _, c := range k.CurrencyList {
			if _, ada := nilaiBaru[c.Name]; !ada {
				urut = append(urut, c.Name)
			}
			nilaiLama[c.Name] = c.Premium
		}
		for _, mu := range urut {
			h := Hasil{Kasus: k.Kasus, Lini: string(premium.LiniMBU), Langkah: langkahAgregatMBU + " " + mu,
				SistemLama: nilaiLama[mu], SistemBaru: nilaiBaru[mu], Status: Selisih}
			if sama, err := samaNilai(h.SistemLama, h.SistemBaru); err == nil && sama {
				h.Status = Cocok
			}
			hasil = append(hasil, h)
		}
	}
	return hasil, nil
}

// samaNilai - kesamaan nilai desimal eksak; teks kosong atau tak terbaca = tidak sama.
func samaNilai(a, b string) (bool, error) {
	da, err := utils.ParseDecimal(a)
	if err != nil {
		return false, err
	}
	db, err := utils.ParseDecimal(b)
	if err != nil {
		return false, err
	}
	return da.Cmp(db) == 0, nil
}

// BerkasKasus melaporkan satu fixture kasus ter-de-identifikasi (`premium/testdata/kasus`):
// lini bisnisnya (lewat premium.LiniDariPredikat) dan apa yang dapat dibandingkan.
func BerkasKasus(nama string, isi []byte) ([]Hasil, error) {
	if err := deidentifikasi.SudahBersih(isi); err != nil {
		return nil, fmt.Errorf("rekonsiliasi: %s: %w", nama, err)
	}
	var akar map[string]any
	if err := json.Unmarshal(isi, &akar); err != nil {
		return nil, fmt.Errorf("rekonsiliasi: %s: %w", nama, err)
	}
	hasil, err := barisPremi(nama, kasusJSON{akar})
	if err != nil {
		return nil, err
	}
	hasil = append(hasil, barisPembayaran(nama, kasusJSON{akar})...)
	hasil = append(hasil, barisTotal(nama, kasusJSON{akar})...)
	return append(hasil, Hasil{Kasus: nama, Langkah: langkahNilaiDasar, Status: BelumTercakup,
		Catatan: "nilai sistem lama (DataSearch.CARID2, LetterNo) tidak terekam di fixture"}), nil
}

// barisPembayaran - tiket 19: `pyWorkPage.Policy.Payment` kasus MARINE CARGO dihitung
// lewat pembayaran.Marine, tetapi nilai lamanya tidak terekam. `[terverifikasi]`
// PremiPaymentMarine berkelas ASM-FW-GISFW-Work (L66), jadi `.Policy` dan
// `.CargoList()` relatif terhadap pyWorkPage; akar fixture adalah halaman OfferFacIn
// (`CurrencyList(n).Policy.Payment` halaman lain). `[dugaan]` `pyWorkPage.CargoList`
// berisi sama dengan `OfferFacIn.CargoList` yang dijumlahkan di sini - belum
// terverifikasi, jadi barisnya belum tercakup, bukan pembanding. Bukan MARINE = tanpa
// baris.
func barisPembayaran(nama string, k kasusJSON) []Hasil {
	var kargo [][]pembayaran.Coverage
	for _, c := range telusuri(k.akar, "", []string{"CargoList"}) {
		var cov []pembayaran.Coverage
		for _, v := range telusuri(c.m, "", []string{"CoverageList"}) {
			cov = append(cov, pembayaran.Coverage{MataUang: teksDi(v.m, "Currency", "Name"),
				Premium: teksDi(v.m, "Premium"), Discount: teksDi(v.m, "Discount")})
		}
		kargo = append(kargo, cov)
	}
	h, err := pembayaran.Marine(k, kargo)
	switch {
	case err != nil:
		return []Hasil{{Kasus: nama, Lini: string(premium.LiniMarineCargo), Langkah: langkahPembayaran, Status: Galat, Catatan: err.Error()}}
	case !h.Berlaku:
		return nil
	}
	return []Hasil{{Kasus: nama, Lini: string(premium.LiniMarineCargo), Langkah: langkahPembayaran, Status: BelumTercakup,
		SistemBaru: utils.FormatDecimal(h.Premium.Amount),
		Catatan:    "pyWorkPage.Policy.Payment sistem lama tidak terekam di fixture (akar = OfferFacIn)"}}
}

// barisTotal - tiket 07: total premi per item properti (`.TotalGrossPremi`) dan per
// mata uang lokasi (`.Property.TotalTSIPremiGrossList`: TSI, Premium, Rate) kasus
// FIRE, lewat premium.TotalFireLokasi (SumTotalTSIPremiGross_Act langkah 1) dari premi
// coverage TERSIMPAN, dibandingkan sebagai teks persis (TSI per mata uang: nilai, A47).
// [terverifikasi] Kasus EDM pun memakai rule ini: nilai tersimpan `edm-fire-1` terulang
// oleh model langkah 1 (dengan double, A48) - jadi tidak diistimewakan. Entri
// LewatDouble (cabang 1.3.4) → belum tercakup, bukan dibandingkan. Bukan FIRE = tanpa
// baris.
func barisTotal(nama string, k kasusJSON) []Hasil {
	if fire, err := rules.Eval("IsFire", k); err != nil || !fire {
		return nil
	}
	lini := string(premium.LiniFire)
	var hasil []Hasil
	baris := func(langkah, lama, baru string) {
		h := Hasil{Kasus: nama, Lini: lini, Langkah: langkahTotal + " " + langkah, SistemLama: lama, SistemBaru: baru, Status: Selisih}
		if lama == baru {
			h.Status = Cocok
		}
		hasil = append(hasil, h)
	}
	for _, l := range telusuri(k.akar, "", []string{"LocationList", "Property"}) {
		var items []premium.ItemProperti
		var lamaItem []string
		for _, it := range telusuri(l.m, l.jalur, []string{"PropertyItemList"}) {
			// `Local.Currency = .Currency` (1.3.3) - teks; [terverifikasi] 11/11 item FIRE
			// fixture ber-`.Currency` teks.
			ip := premium.ItemProperti{MataUang: teksDi(it.m, "Currency"), TSIObjectItem: teksDi(it.m, "TSIObjectItem")}
			for _, c := range telusuri(it.m, "", []string{"CoverageList"}) {
				ip.PremiCoverage = append(ip.PremiCoverage, teksDi(c.m, "Premium"))
			}
			items = append(items, ip)
			lamaItem = append(lamaItem, it.jalur+".TotalGrossPremi="+teksDi(it.m, "TotalGrossPremi"))
		}
		tot, err := premium.TotalFireLokasi(items)
		if err != nil {
			hasil = append(hasil, Hasil{Kasus: nama, Lini: lini, Langkah: langkahTotal + " | " + l.jalur, Status: Galat, Catatan: err.Error()})
			continue
		}
		for i, p := range tot.PremiItem {
			jalur, lama, _ := strings.Cut(lamaItem[i], "=")
			baris("1.3.3 | "+jalur, lama, utils.FormatDecimal(p.Amount))
		}
		lamaList := telusuri(l.m, l.jalur, []string{"TotalTSIPremiGrossList"})
		if len(lamaList) != len(tot.PerMataUang) {
			hasil = append(hasil, Hasil{Kasus: nama, Lini: lini, Langkah: langkahTotal + " 1.3.4-1.3.5 | " + l.jalur, Status: Selisih,
				Catatan: fmt.Sprintf("%d mata uang tersimpan, %d dihitung", len(lamaList), len(tot.PerMataUang))})
			continue
		}
		for i, e := range tot.PerMataUang {
			g := lamaList[i]
			if e.LewatDouble {
				hasil = append(hasil, Hasil{Kasus: nama, Lini: lini, Langkah: langkahTotal + " 1.3.4 | " + g.jalur, Status: BelumTercakup,
					Catatan: "≥2 item bermata uang sama: sistem lama membulatkan item ke-2 dst. ke double (Local.Premi/Local.TSI, L293-301) lalu mengakumulasi Decimal; port tetap eksak (A48, butir 63)"})
				continue
			}
			// TSI dibandingkan sebagai NILAI (A47): tersimpan `8400000.0` untuk
			// TSIObjectItem `8400000`. `[dugaan]` perbedaan format angka; sebabnya belum
			// terverifikasi (tipe properti TSI kelas OfferFacIn-Currency tidak ada di korpus).
			// Name, Premium, Rate: teks persis.
			tsiLama := teksDi(g.m, "TSI")
			tsiBaru := utils.FormatDecimal(e.TSI.Amount)
			if sama, err := samaNilai(tsiLama, tsiBaru); err == nil && sama {
				tsiBaru = tsiLama
			}
			lama := strings.Join([]string{teksDi(g.m, "Name"), tsiLama, teksDi(g.m, "Premium"), teksDi(g.m, "Rate")}, " ")
			baru := strings.Join([]string{e.TSI.Currency, tsiBaru, utils.FormatDecimal(e.Premium.Amount), utils.FormatDecimal(e.Rate)}, " ")
			baris("1.3.4-1.3.5 | "+g.jalur+" (Name TSI Premium Rate; TSI dibanding nilai, A47)", lama, baru)
		}
	}
	return hasil
}

func barisPremi(nama string, k kasusJSON) (hasil []Hasil, err error) {
	defer func() {
		// Tidak satu gerbang lini terbuka → panic jembatan; dilaporkan per kasus.
		if r := recover(); r != nil {
			hasil, err = []Hasil{{Kasus: nama, Langkah: langkahPremi, Status: Galat, Catatan: fmt.Sprint(r)}}, nil
		}
	}()
	lini, err := premium.LiniDariPredikat(k)
	if err != nil {
		return nil, fmt.Errorf("rekonsiliasi: %s: %w", nama, err)
	}
	catatanKasus := strings.Join(lini.Peringatan, "; ")
	edm := teksDi(k.akar, "QuotationData", "StatusBusiness") == "3"
	if edm {
		// CountGrossPremi_Act langkah 5 bergerbang StatusBusiness != 3; langkah 6
		// memanggil CountGrossPremiEDM_Act, yang hanya berjalan bila == 3 (A39, A44).
		catatanKasus = strings.TrimPrefix(catatanKasus+"; StatusBusiness 3: jalur EDM tidak diport, rumus coverage NB", "; ")
	}
	for _, l := range lini.Lini {
		if edm && l.Lini != premium.LiniFire {
			hasil = append(hasil, Hasil{Kasus: nama, Lini: string(l.Lini), Langkah: langkahPremi, Status: BelumTercakup,
				Catatan: "StatusBusiness 3: CountGrossPremiEDM_Act cabang lini ini belum dibandingkan (A44)"})
			continue
		}
		jalur, ada := jalurCoverage[l.Lini]
		if !ada {
			catatan := "rumus premi ada, tetapi pemetaan coverage berkas kasus → premium.Input belum dibuat"
			if _, err := premium.Calculate(premium.Input{LiniBisnis: l.Lini}); errors.Is(err, premium.ErrBentukBelumDiport) {
				catatan = "rumus premi lini ini belum diport"
			}
			hasil = append(hasil, Hasil{Kasus: nama, Lini: string(l.Lini), Langkah: langkahPremi, Status: BelumTercakup,
				Catatan: strings.Trim(catatan+" | "+catatanKasus, " |")})
			continue
		}
		mbd := false
		if l.Lini == premium.LiniAneka || l.Lini == premium.LiniGolf {
			if mbd, err = rules.Eval("IsMBD", k); err != nil {
				hasil = append(hasil, Hasil{Kasus: nama, Lini: string(l.Lini), Langkah: langkahPremi, Status: Galat, Catatan: "IsMBD: " + err.Error()})
				continue
			}
		}
		for _, c := range telusuriCoverage(k.akar, jalur) {
			if edm {
				if alasan := bedaEDM(c); alasan != "" {
					hasil = append(hasil, Hasil{Kasus: nama, Lini: string(l.Lini), Langkah: langkahPremi + " | " + c.jalur,
						Status: Galat, Catatan: alasan})
					continue
				}
			}
			in := masukanCoverage(k.akar, c.m, l.Lini, mbd)
			h := bandingkanPremi(nama, in, teksDi(c.m, "Premium"))
			h.Langkah += " | " + c.jalur
			if edm {
				h.Catatan = strings.Trim(h.Catatan+" | "+catatanTSIItem(c), " |")
			}
			h.Catatan = strings.Trim(h.Catatan+" | "+catatanKasus, " |")
			hasil = append(hasil, h)
		}
	}
	return hasil, nil
}

// jalurCoverage - letak CoverageList tiap lini di fixture kasus (akar = halaman
// OfferFacIn). `[terverifikasi]` sesuai halaman langkah pemilih rumus:
// CountGrossPremi_Act 5.1.2 / 5.2.1 / 5.3.1, CountGPWMarinePAMbu_Act 1.1.1.
var jalurCoverage = map[premium.LiniBisnis][]string{
	premium.LiniFire:        {"LocationList", "Property", "PropertyItemList", "CoverageList"},
	premium.LiniAneka:       {"LocationList", "Property", "RiskLocation", "OccupationList", "AnekaList", "CoverageList"},
	premium.LiniGolf:        {"LocationList", "Property", "RiskLocation", "AnekaList", "CoverageList"},
	premium.LiniMarineCargo: {"CargoList", "CoverageList"},
}

// masukanCoverage - premium.Input satu coverage. Pro-rata dari
// `pyWorkPage.OfferFacIn.ProRatePercent` (CountPremi_ACT langkah 10 /
// CountPremiCoverageAneka langkah 6 menyalin Local.Prorate ke sana); master policy
// MARINE dari QuotationData.
func masukanCoverage(akar, c map[string]any, lini premium.LiniBisnis, mbd bool) premium.Input {
	return premium.Input{LiniBisnis: lini, MataUang: teksDi(c, "Currency", "Name"),
		CalculateMethod: teksDi(c, "CalculateMethod_FacIn"), TSI: teksDi(c, "TSI"), Rate: teksDi(c, "Rate"),
		ProRatePercent: teksDi(akar, "ProRatePercent"), PctShortPeriod: teksDi(c, "PctShortPeriod"),
		Loading: teksDi(c, "Loading"), IndemnityPercentage: teksDi(c, "IndemnityPercentage"),
		LossLimit: teksDi(c, "LostLimit"), PctAdjustment: teksDi(c, "PctAdjustment"),
		CoverageBasis: teksDi(c, "CoverageBasis"), NetRate: teksDi(c, "NetRate"), FirstScale: teksDi(c, "FirstScale"),
		MBD:          mbd,
		MasterPolicy: teksDi(akar, "QuotationData", "PolicyType") == "1" && teksDi(akar, "QuotationData", "IsMOP") == "MOP"}
}

// simpul - satu halaman di fixture beserta jalurnya, mis. `CargoList[3].CoverageList[0]`,
// dan halaman induk langsung coverage (item properti FIRE; dipakai bedaEDM).
type simpul struct {
	jalur string
	m     map[string]any
	induk map[string]any
}

// telusuriCoverage - seluruh coverage di ujung `jalur` (berakhir `CoverageList`),
// urut seperti fixture, masing-masing dengan halaman induknya.
func telusuriCoverage(akar map[string]any, jalur []string) []simpul {
	var hasil []simpul
	for _, ind := range telusuri(akar, "", jalur[:len(jalur)-1]) {
		for _, c := range telusuri(ind.m, ind.jalur, jalur[len(jalur)-1:]) {
			c.induk = ind.m
			hasil = append(hasil, c)
		}
	}
	return hasil
}

// bedaEDM - alasan coverage FIRE kasus EDM tidak dapat dibandingkan dengan rumus NB.
// `[terverifikasi]` CountGrossPremiEDM_Act (`NB FacIn\Activity\CountGrossPremiEDM_Act.xml`,
// isi identik varian NB, RNW, dan EDM setelah awalan `RH_n` disamakan) cabang 1.1 FIRE
// memanggil CountPremi_ACT yang sama (1.1.2.1.4.2, L1468), lalu 1.1.2.1.4.3 menegasikan
// `.TSI` (L1623), `.TSILiability` (L1677), `.Premium` (L1697) bila `.FlagDelete=="1" &&
// Local.deleteobjitem!="1"` - belum diport (A44, milik endorsmentfacin) → galat.
// Kosong = tidak ada beda.
func bedaEDM(c simpul) string {
	if teksDi(c.m, "FlagDelete") == "1" && teksDi(c.induk, "FlagDelete") != "1" {
		return "EDM 1.1.2.1.4.3: premi coverage FlagDelete dinegasikan (belum diport, A44)"
	}
	return ""
}

// catatanTSIItem - `[dugaan]`, BUKAN galat: pasangan `.TSI = Local.TSIObjectItem`
// (L1273) hanya ada di pyParamArray langkah 1.1.2.1.4 yang bermetode kosong (loop
// EMBEDDED atas .CoverageList), jadi kemungkinan tidak dieksekusi. Bila TSI coverage
// berbeda dari TSIObjectItem item, perbandingan tetap jalan dan diberi tanda.
func catatanTSIItem(c simpul) string {
	// Dibandingkan sebagai NILAI: `8400000` lawan `8400000.0` bukan perbedaan.
	if sama, err := samaNilai(teksDi(c.m, "TSI"), teksDi(c.induk, "TSIObjectItem")); err != nil || !sama {
		return "[dugaan] TSIObjectItem item berbeda dari TSI coverage; EDM 1.1.2.1.4 mungkin menimpa .TSI (A44)"
	}
	return ""
}

// telusuri - semua halaman di ujung `kunci`, menembus daftar; urut seperti fixture.
func telusuri(m map[string]any, awalan string, kunci []string) []simpul {
	if len(kunci) == 0 {
		return []simpul{{jalur: strings.TrimPrefix(awalan, "."), m: m}}
	}
	var hasil []simpul
	switch v := m[kunci[0]].(type) {
	case map[string]any:
		hasil = append(hasil, telusuri(v, awalan+"."+kunci[0], kunci[1:])...)
	case []any:
		for i, x := range v {
			if anak, ok := x.(map[string]any); ok {
				hasil = append(hasil, telusuri(anak, fmt.Sprintf("%s.%s[%d]", awalan, kunci[0], i), kunci[1:])...)
			}
		}
	}
	return hasil
}

// teksDi - teks di jalur halaman; tidak ada = kosong.
func teksDi(m map[string]any, kunci ...string) string {
	var v any = m
	for _, k := range kunci {
		mm, ok := v.(map[string]any)
		if !ok {
			return ""
		}
		v = mm[k]
	}
	s, _ := v.(string)
	return s
}

// kasusJSON - rules.Kasus atas fixture yang akarnya halaman OfferFacIn.
//
//	pyWorkPage.OfferFacIn.<jalur> → <jalur> di akar
//	pyWorkPage.Quotation.<jalur>  → QuotationData.<jalur> di akar (A29, butir 49: fixture hanya
//	                                menyimpan jalur ini; pemuat mengisi KEDUA jalur dan
//	                                LiniDariPredikat tetap mencatat bila berbeda)
//
// Kedua jalur berisi sama (keputusan work owner 01-10-2026, butir 49). Jalur di luar
// kedua awalan = tidak ada. Dipakai gerbang lini (LiniDariPredikat) dan, sejak
// 01-10-2026, juga `IsMBD`, `IsFire`, `IsNB`, `IsEDM`, `IsMarineCargo` (barisPremi,
// barisTotal, barisPembayaran). ⚠️ Batas bukti: `IsNB` (`Quotation.StatusBusiness`)
// dan `IsEDM` (`OfferFacIn.QuotationData.StatusBusiness`) membaca kolom yang SAMA di
// fixture, jadi keadaan "keduanya benar" tidak teruji dari data nyata (ralat kedua).
type kasusJSON struct{ akar map[string]any }

func (k kasusJSON) Nilai(jalur string) (string, bool) {
	var rel string
	switch {
	case strings.HasPrefix(jalur, "pyWorkPage.OfferFacIn."):
		rel = strings.TrimPrefix(jalur, "pyWorkPage.OfferFacIn.")
	case strings.HasPrefix(jalur, "pyWorkPage.Quotation."):
		rel = "QuotationData." + strings.TrimPrefix(jalur, "pyWorkPage.Quotation.")
	default:
		return "", false
	}
	var v any = k.akar
	for _, seg := range strings.Split(rel, ".") {
		m, ok := v.(map[string]any)
		if !ok {
			return "", false
		}
		if v, ok = m[seg]; !ok {
			return "", false
		}
	}
	s, ok := v.(string)
	return s, ok
}

// Ringkas - jumlah hasil per status.
func Ringkas(hasil []Hasil) map[Status]int {
	r := map[Status]int{}
	for _, h := range hasil {
		r[h.Status]++
	}
	return r
}
