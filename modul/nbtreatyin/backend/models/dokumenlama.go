package models

// Untuk apa berkas ini: PEMECAH DOKUMEN LAMA - tiket 22 (spec-penyimpanan
// ID-3, ID-19, ID-27; AC 21-22, 52-59; KEPUTUSAN-RONDE-12 butir 5; jawaban
// work owner K15, K17 PROMPT-NB-TREATY-IN-PUTARAN-2.md bab 2).
//
// Dokumen lama = satu baris `POOLDATA.JSON_POLIS`. `[terverifikasi]`
// `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 (halaman langkah
// `pyWorkPage.PolicyTreatyIn`, kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`):
// `InputData.CARI3 = @ASM.GetPageJSONString()` - DATA_JSON adalah HALAMAN
// `PolicyTreatyIn` apa adanya. `RDBList\SavePolisTreatyIn_SQL.xml`:
// `PEGA_JSON_POLIS_TREATYIN({pyWorkPage.pzInsKey}, {PolicyTreatyIn.PolicyNo},
// NULL, '0', {InputData.CARI21}, {OperatorID.pyUserIdentifier}, {CARI3})` -
// IDPEGA, NOPOLIS, NOENDORS, PRODKE, TGL_PROD, USERNAME, DATA_JSON.
//
// Pemecah ini MURNI (seam 3 spec.md §6.2): tidak menyentuh basis data,
// berkas, maupun jam. Ia menghasilkan halaman kerja yang lalu disimpan lewat
// antarmuka penyimpanan yang SAMA dengan jalur biasa (`SimpanHalaman`, ID-3).
// Pemetaan jalur dokumen -> kolom DIGERAKKAN KATALOG (`katalog.go`), bukan
// daftar kolom salinan.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// KelasDokumenTreatyIn - `pxObjClass` akar DATA_JSON polis treaty inward
// (`SaveJsonPolisTreatyIn_Act` langkah 6, `pyStepsClassName`). POOLDATA.JSON_POLIS
// dipakai bersama lini lain (contoh pertama P29 milik Fac In, rancangan 4q.4):
// dokumen berkelas lain BUKAN lingkup pemuat ini dan hanya dihitung.
const KelasDokumenTreatyIn = "ASM-FW-GISFW-Data-PolicyTreatyIn"

// BarisJSONPolis - satu baris POOLDATA.JSON_POLIS: tujuh kolom datar (ID-21)
// dan dokumennya. Tanggal datar berbentuk `utils.TanggalWaktu` (dibaca
// repository lewat TO_CHAR); kosong = NULL.
type BarisJSONPolis struct {
	IDPega, NoPolis, NoEndors, ProdKe string
	TglInput, TglProd, Username       string
	DataJSON                          []byte
}

// KolomDatarLama - kolom datar json_polis yang ditulis apa adanya ke
// T_GENERAL_POLIS di luar katalog (ID-21). NOPOLIS ditulis lewat
// `SetelNomorPolis`, PRODKE selalu 0 (`SisipKasus`), TGL_PROD lewat katalog
// (`ProductionDate`).
type KolomDatarLama struct {
	IDPega, NoEndors, TglInput, Username string
}

// Medan - satu medan daun dokumen. Jalur memakai notasi properti Pega
// relatif `pyWorkPage` (`PolicyTreatyIn.ListInstallment(2).Premium`); Pola
// adalah jalur tanpa nomor baris (`PolicyTreatyIn.ListInstallment().Premium`).
type Medan struct {
	Jalur, Pola, Nilai string
	// letak di halaman: daftar == "" berarti nilai halaman `Jalur`.
	daftar string
	baris  int
	nama   string
}

// GalatDokumen - satu sebab dokumen tidak dimuat (AC 58).
type GalatDokumen struct {
	Jalur, Nilai string
	Err          error
}

// HasilPecah - keluaran pemecah untuk satu dokumen.
type HasilPecah struct {
	// ID - kunci T_WORK_POLIS / T_GENERAL_POLIS = pyID dari IDPEGA.
	ID      string
	NoPolis string
	Halaman *Halaman
	Datar   KolomDatarLama
	// TakDikenal - medan tanpa kolom dan tanpa keputusan tertulis; masuk
	// berkas laporan CSV, TIDAK dibuang (K17, AC 57).
	TakDikenal []Medan
	// Diabaikan - alasan tertulis -> cacah medan yang sengaja tidak disimpan.
	Diabaikan map[string]int
	// Galat - dokumen TIDAK dimuat bila terisi (K15, AC 58).
	Galat []GalatDokumen
}

// Galat struktural dokumen lama.
var (
	// ErrDokumenRusak - DATA_JSON bukan objek JSON yang dapat diurai.
	ErrDokumenRusak = errors.New("models: DATA_JSON tidak dapat diurai sebagai objek JSON")
	// ErrBukanTreatyIn - dokumen lini lain (bukan galat; dihitung saja).
	ErrBukanTreatyIn = errors.New("models: dokumen bukan polis Treaty In")
	// ErrGenerasiEndorsemen - PRODKE > 0: generasi endorsemen, milik pemuat
	// EDM (edmtreatyin tiket 10), bukan galat.
	ErrGenerasiEndorsemen = errors.New("models: generasi endorsemen (PRODKE > 0) - milik pemuat EDM tiket 10")
	// ErrProdKe - PRODKE kosong atau bukan bilangan.
	ErrProdKe = errors.New("models: PRODKE kosong atau bukan bilangan")
	// ErrIDPega - IDPEGA bukan `<kelas> <pyID>`.
	ErrIDPega = errors.New("models: IDPEGA tidak berbentuk \"<kelas> <awalan>-<nomor>\"")
	// ErrNilaiKolom - nilai tidak dapat ditulis ke kolomnya (aturan sama
	// dengan konversi repository `nilaiTulis`).
	ErrNilaiKolom = errors.New("models: nilai tidak sesuai tipe kolomnya")
	// ErrNoPolisKosong - JSON_POLIS.NOPOLIS kosong. Dokumen hanya ditulis
	// bila nomor polis terisi (Flow InputRealizationTreatyIn: Decision8
	// "Nopolis not empty" -> Utility1 `SaveJsonPolisTreatyIn_Act`).
	ErrNoPolisKosong = errors.New("models: NOPOLIS kosong")
	// ErrNoPolisBeda - PolicyNo dokumen berbeda dari JSON_POLIS.NOPOLIS
	// (keduanya dari `PolicyTreatyIn.PolicyNo`, SavePolisTreatyIn_SQL).
	ErrNoPolisBeda = errors.New("models: PolicyNo dokumen berbeda dari JSON_POLIS.NOPOLIS")
)

// Galat pembacaan tanggal lama.
var (
	// ErrTanggalAmbigu - susunan hari/bulan tidak dapat ditentukan dari
	// nilainya (`05/06/2017`). `[keputusan work owner]` K15: TIDAK ditebak;
	// dokumennya masuk laporan galat dan jumlahnya dilaporkan.
	ErrTanggalAmbigu = errors.New("models: tanggal ambigu (susunan hari/bulan tidak dapat ditentukan) - tidak ditebak (K15)")
	// ErrFormatTanggal - bukan salah satu dari dua format dokumen lama
	// (ID-19), atau tanggalnya tidak ada di kalender.
	ErrFormatTanggal = errors.New("models: format tanggal di luar YYYYMMDD / cap waktu Pega ber-GMT")
)

var (
	// polaYYYYMMDD - tanggal Pega (`"20171130"`, P29 sifat 2).
	polaYYYYMMDD = regexp.MustCompile(`^\d{8}$`)
	// polaCapWaktuGMT - cap waktu Pega (`"20170930T170000.000 GMT"`, P29 sifat 2).
	polaCapWaktuGMT = regexp.MustCompile(`^(\d{8}T\d{6})(\.\d{1,3})? GMT$`)
	// polaGarisMiring - susunan `dd/MM/yyyy` atau `MM/dd/yyyy …` (P32:
	// `InputPolicyTreatyIn_preDT` memakai keduanya).
	polaGarisMiring = regexp.MustCompile(`^(\d{1,2})/(\d{1,2})/(\d{4})(\s.*)?$`)
)

// BacaTanggalLama mengubah nilai tanggal dokumen lama menjadi bentuk
// pertukaran repository (`utils.TanggalSaja` / `utils.TanggalWaktu`).
//
//   - `YYYYMMDD` (AC 21) -> `YYYY-MM-DD`; jam tidak dikarang, juga untuk
//     kolom tanggal-waktu.
//   - `YYYYMMDDTHHMMSS.mmm GMT` (AC 22) -> jam dinding Asia/Jakarta. Pega
//     menyimpan cap waktu dalam GMT dan rule-nya sendiri membaca harinya di
//     Asia/Jakarta (`GeneratePolicyNoTreaty_Act` langkah 5.3, `TanggalProduksiNomor`):
//     `20170930T170000.000 GMT` adalah 1 Oktober 2017 pukul 00.00 WIB.
//     Kolom bertanggal saja menerima tanggal kalender Jakarta itu.
//   - `05/06/2017` -> ErrTanggalAmbigu (K15), tidak ditebak.
//   - selain itu -> ErrFormatTanggal, juga garis miring yang tidak ambigu:
//     cara menentukan susunan per baris belum diputuskan (P32 butir 1).
func BacaTanggalLama(teks string, g Golongan) (string, error) {
	teks = strings.TrimSpace(teks)
	switch {
	case teks == "":
		return "", nil
	case polaYYYYMMDD.MatchString(teks):
		t, err := time.Parse("20060102", teks)
		if err != nil {
			return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
		}
		return t.Format("2006-01-02"), nil
	}
	if m := polaCapWaktuGMT.FindStringSubmatch(teks); m != nil {
		t, err := time.ParseInLocation("20060102T150405", m[1], time.UTC)
		if err != nil {
			return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
		}
		lokal := t.In(zonaJakarta())
		if g == GolTanggal {
			return lokal.Format("2006-01-02"), nil
		}
		return lokal.Format("2006-01-02 15:04:05"), nil
	}
	if m := polaGarisMiring.FindStringSubmatch(teks); m != nil {
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		if a >= 1 && a <= 12 && b >= 1 && b <= 12 && a != b {
			return "", fmt.Errorf("%w: %q", ErrTanggalAmbigu, teks)
		}
	}
	return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
}

// ---------------------------------------------------------------- penggolong

// Alasan medan dokumen yang SENGAJA tidak disimpan - masing-masing keputusan
// tertulis. Dihitung per alasan di ringkasan pemuat, tidak hilang diam-diam.
const (
	AlasanInternalPega  = "pxObjClass - internal Pega, tidak dimigrasi (P29 sifat 4; rancangan §4.1)"
	AlasanNourut        = "pxListSubscript - nomor baris, diwakili NOURUT (spec-penyimpanan ID-11)"
	AlasanKeadaanLayar  = "Show/ViewState/pxResults/FillPaymentInstallmentEDMT - keadaan layar (rancangan §4.1, 4q5.7)"
	AlasanTurunan       = "Total* tingkat polis - turunan baris spreading, dihitung saat dibaca (katalog.go; RALAT AC 38)"
	AlasanPantulanLayer = "Layer* tingkat polis - pantulan baris layer pertama (ID-22, rancangan 4q.10)"
	AlasanBreakdown     = "BreakDownSpreadList - tidak dimigrasi (KEPUTUSAN-RONDE-12 butir 3/3b)"
	AlasanPersetujuanDH = "isApprovedtoDeptHead - tidak dibangun (spec AC 64)"
	AlasanOldData       = "OldData - diganti penunjuk OLD_POLIS_ID, kosong di NB (rancangan 4ter.1, ID-9)"
	AlasanSelisih       = "TreatyDifference/TreatyXOLDifferenceList - selisih endorsemen, nol baris di NB (rancangan 4ter.2)"
)

// IndukDaftarBersarang - tabel cucu -> jalur daftar induknya di dokumen
// (`ListInstallment(n).InstallmentList`, `TreatyXOLList(n).ValueList`;
// pasangan `cucuAngsuran` / `cucuLayer` repository). Tabel katalog yang
// `Daftar`-nya relatif WAJIB terdaftar di sini (uji `TestPetaKatalogDokumen`).
//
// Diturunkan dari `keturunanKatalog` (katalog.go) - SATU sumber silsilah
// tabel yang juga dipakai `ProyeksiKatalog`.
var IndukDaftarBersarang = func() map[string]string {
	m := map[string]string{}
	for _, kt := range keturunanKatalog {
		if kt.cucu != nil {
			m[kt.cucu.Nama] = kt.anak.Daftar
		}
	}
	return m
}()

// jalurNomorPolis - `PolicyNo` ditulis ke NOPOLIS (SetelNomorPolis), bukan
// lewat katalog.
const jalurNomorPolis = HalamanPolis + ".PolicyNo"

// petaKatalogDokumen menurunkan pola jalur dokumen -> kolom dari KATALOG.
// Tabel 1:1 memakai jalur properti; `T_POLIS_QUOTATION` dibaca repository
// dari `PolicyTreatyIn.QuotationData.<medan>` (`nilaiQuotation`); tabel anak
// dari `<Daftar>().<medan>`; tabel cucu dari `<induk>().<Daftar>().<medan>`.
// Medan halaman kerja (PositionNote, NBStatus, TreatyIn.ID) bukan isi dokumen.
func petaKatalogDokumen() (map[string]Kolom, error) {
	peta := map[string]Kolom{}
	for _, t := range SemuaTabel {
		var awalan string
		switch {
		case t.Daftar == "" && t.Nama == TabelQuotation.Nama:
			awalan = HalamanPolis + ".QuotationData."
		case t.Daftar == "":
			awalan = ""
		case strings.HasPrefix(t.Daftar, HalamanPolis+"."):
			awalan = t.Daftar + "()."
		default:
			induk, ada := IndukDaftarBersarang[t.Nama]
			if !ada {
				return nil, fmt.Errorf("models: daftar bersarang %s (%s) tanpa induk di IndukDaftarBersarang", t.Daftar, t.Nama)
			}
			awalan = induk + "()." + t.Daftar + "()."
		}
		for _, k := range t.Kolom {
			pola := awalan + k.Properti
			if !strings.HasPrefix(pola, HalamanPolis+".") {
				continue
			}
			peta[pola] = k
		}
	}
	return peta, nil
}

// alasanDiabaikan - alasan tertulis bila pola medan sengaja tidak disimpan.
//
// Simpul utuh (OldData, TreatyDifference, ...) diperiksa lebih dulu: seluruh
// isinya - termasuk pxObjClass-nya - tercatat di bawah alasan simpul itu.
func alasanDiabaikan(pola string) string {
	akar := strings.TrimPrefix(pola, HalamanPolis+".")
	puncak := akar
	if i := strings.IndexAny(akar, ".("); i >= 0 {
		puncak = akar[:i]
	}
	switch puncak {
	case "Show", "ViewState", "pxResults", "FillPaymentInstallmentEDMT":
		return AlasanKeadaanLayar
	case "BreakDownSpreadList":
		return AlasanBreakdown
	case "OldData":
		return AlasanOldData
	case "TreatyDifference", "TreatyXOLDifferenceList":
		return AlasanSelisih
	}
	ruas := pola[strings.LastIndex(pola, ".")+1:]
	switch {
	case ruas == "pxObjClass":
		return AlasanInternalPega
	case ruas == "pxListSubscript" && strings.Contains(pola, "()"):
		return AlasanNourut
	case puncak != akar:
		return ""
	}
	switch puncak {
	case "TotalPremium", "TotalClaim", "TotalSharePercentagePremium", "TotalSharePercentageClaim":
		return AlasanTurunan
	case "Layer", "LayerType", "LayerPart", "LayerPartType":
		return AlasanPantulanLayer
	case "isApprovedtoDeptHead":
		return AlasanPersetujuanDH
	}
	return ""
}

// ---------------------------------------------------------------- pemecah

// polaPyID - pyID kasus: awalan huruf, tanda hubung, nomor (`NB-77`).
var polaPyID = regexp.MustCompile(`^[A-Z][A-Z0-9]*-\d+$`)

// IDKasusDariIDPega mengambil pyID dari `pyWorkPage.pzInsKey`
// (`<kelas> <pyID>`, `KunciInstans`). pyID menjadi ID T_WORK_POLIS - diagram
// grilling: T_WORK_POLIS "diambil dari pyWorkPage.pzInsKey".
func IDKasusDariIDPega(idpega string) (string, error) {
	s := strings.TrimSpace(idpega)
	i := strings.LastIndex(s, " ")
	if i <= 0 {
		return "", fmt.Errorf("%w: %q", ErrIDPega, idpega)
	}
	id := s[i+1:]
	if !polaPyID.MatchString(id) || len(id) > 32 {
		return "", fmt.Errorf("%w: %q", ErrIDPega, idpega)
	}
	return id, nil
}

// periksaNilai menolak nilai yang akan ditolak konversi repository
// (`repository.nilaiTulis`) - supaya uji-kering melaporkan galat yang sama
// tanpa menyentuh basis data.
func periksaNilai(k Kolom, teks string) error {
	s := strings.TrimSpace(teks)
	switch {
	case k.Golongan.Desimal():
		if s == "" {
			return nil
		}
		if _, err := utils.ParseDecimal(s); err != nil {
			return fmt.Errorf("%w: %s bukan angka desimal", ErrNilaiKolom, k.Kolom)
		}
	case k.Golongan == GolCacah:
		if s != "" && !utils.AngkaSaja(s) {
			return fmt.Errorf("%w: %s bukan bilangan bulat", ErrNilaiKolom, k.Kolom)
		}
	case !k.Golongan.Tanggal():
		if k.Panjang > 0 && len([]rune(teks)) > k.Panjang {
			return fmt.Errorf("%w: %s melebihi %d karakter", ErrNilaiKolom, k.Kolom, k.Panjang)
		}
	}
	return nil
}

// pejalan menelusuri pohon JSON menjadi halaman kerja dan daftar medan daun.
type pejalan struct {
	h     *Halaman
	medan []Medan
}

func kunciUrut(m map[string]any) []string {
	k := make([]string, 0, len(m))
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	return k
}

// teksSkalar - seluruh nilai dokumen lama bertipe teks (P29 sifat 1); angka
// JSON dibawa sebagai teks literalnya (json.Number), tidak pernah float.
func teksSkalar(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case bool:
		return strconv.FormatBool(x)
	}
	return fmt.Sprint(v)
}

// objek - simpul objek di tingkat halaman (akar, QuotationData, OldData, ...).
func (p *pejalan) objek(jalur, pola string, m map[string]any) {
	for _, n := range kunciUrut(m) {
		j, q := jalur+"."+n, pola+"."+n
		switch v := m[n].(type) {
		case map[string]any:
			p.objek(j, q, v)
		case []any:
			p.daftar(j, q, v)
		default:
			s := teksSkalar(v)
			p.h.Setel(j, s)
			p.medan = append(p.medan, Medan{Jalur: j, Pola: q, Nilai: s})
		}
	}
}

// daftar - PageList: satu Baris per unsur; daftar bersarang di jalur
// `JalurAnak(daftar, n, nama)`.
func (p *pejalan) daftar(jalur, pola string, xs []any) {
	baris := make([]Baris, len(xs))
	p.h.SetelDaftar(jalur, baris)
	for i, x := range xs {
		baris[i] = Baris{}
		j := fmt.Sprintf("%s(%d)", jalur, i+1)
		switch v := x.(type) {
		case map[string]any:
			p.baris(baris[i], jalur, i, "", j, pola+"()", v)
		case []any:
			p.daftar(j, pola+"()", v)
		default:
			p.medan = append(p.medan, Medan{Jalur: j, Pola: pola + "()", Nilai: teksSkalar(v)})
		}
	}
}

// baris - anggota satu unsur PageList. Objek di dalam baris diratakan
// (`anggota.sub`) - tidak ada di katalog, jadi berakhir sebagai medan tak dikenal.
func (p *pejalan) baris(b Baris, daftar string, i int, awalan, jalur, pola string, m map[string]any) {
	for _, n := range kunciUrut(m) {
		nama := awalan + n
		switch v := m[n].(type) {
		case map[string]any:
			p.baris(b, daftar, i, nama+".", jalur, pola, v)
		case []any:
			p.daftar(jalur+"."+nama, pola+"."+nama, v)
		default:
			s := teksSkalar(v)
			b[nama] = s
			p.medan = append(p.medan, Medan{Jalur: jalur + "." + nama, Pola: pola + "." + nama, Nilai: s,
				daftar: daftar, baris: i, nama: nama})
		}
	}
}

// setel menulis ulang nilai satu medan di halaman (hasil konversi tanggal).
func (p *pejalan) setel(m Medan, nilai string) {
	if m.daftar == "" {
		p.h.Setel(m.Jalur, nilai)
		return
	}
	p.h.AmbilDaftar(m.daftar)[m.baris][m.nama] = nilai
}

// periksaProdKe - "0" = generasi NB (SavePolisTreatyIn_SQL mengirim '0').
func periksaProdKe(prodke string) error {
	s := strings.TrimSpace(prodke)
	n, err := strconv.Atoi(s)
	switch {
	case s == "" || err != nil || n < 0:
		return fmt.Errorf("%w: %q", ErrProdKe, prodke)
	case n > 0:
		return ErrGenerasiEndorsemen
	}
	return nil
}

// PecahDokumenLama memecah satu baris JSON_POLIS menjadi halaman kerja yang
// siap disimpan lewat `SimpanHalaman`.
//
// Galat yang dikembalikan menghentikan dokumen seluruhnya (dokumen rusak,
// IDPEGA/PRODKE tak terbaca) atau menandai dokumen di luar lingkup
// (ErrBukanTreatyIn, ErrGenerasiEndorsemen). Galat per medan - tanggal ambigu
// (K15), format tanggal, nilai kolom, nomor polis - terkumpul di
// `HasilPecah.Galat`; dokumen bergalat TIDAK dimuat sebagian (P31).
//
// ⛔ Nilai lama tidak dihitung ulang dan tidak dibulatkan (AC 55, P29):
// teks desimal dibawa apa adanya; pembulatan hanya terjadi di Oracle pada
// desimal kesembilan (NUMBER(38,8), AC 19, 20b). Satu-satunya pengisian:
// EndDate kosong = StartDate (`[keputusan work owner]` spec AC 69, §5.8).
func PecahDokumenLama(b BarisJSONPolis) (HasilPecah, error) {
	dek := json.NewDecoder(bytes.NewReader(b.DataJSON))
	dek.UseNumber()
	var akar any
	if err := dek.Decode(&akar); err != nil {
		return HasilPecah{}, fmt.Errorf("%w: %v", ErrDokumenRusak, err)
	}
	m, ok := akar.(map[string]any)
	if !ok {
		return HasilPecah{}, ErrDokumenRusak
	}
	if kelas, _ := m["pxObjClass"].(string); kelas != KelasDokumenTreatyIn {
		return HasilPecah{}, fmt.Errorf("%w: pxObjClass %q", ErrBukanTreatyIn, kelas)
	}
	if err := periksaProdKe(b.ProdKe); err != nil {
		return HasilPecah{}, err
	}
	id, err := IDKasusDariIDPega(b.IDPega)
	if err != nil {
		return HasilPecah{}, err
	}
	peta, err := petaKatalogDokumen()
	if err != nil {
		return HasilPecah{}, err
	}

	p := &pejalan{h: HalamanBaru()}
	p.objek(HalamanPolis, HalamanPolis, m)
	hasil := HasilPecah{
		ID: id, NoPolis: strings.TrimSpace(b.NoPolis), Halaman: p.h, Diabaikan: map[string]int{},
		Datar: KolomDatarLama{IDPega: b.IDPega, NoEndors: b.NoEndors, TglInput: b.TglInput, Username: b.Username},
	}
	galat := func(m Medan, err error) {
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: m.Jalur, Nilai: m.Nilai, Err: err})
	}
	for _, md := range p.medan {
		if k, dikenal := peta[md.Pola]; dikenal {
			if !k.Golongan.Tanggal() {
				if err := periksaNilai(k, md.Nilai); err != nil {
					galat(md, err)
				}
				continue
			}
			v, err := BacaTanggalLama(md.Nilai, k.Golongan)
			if err != nil {
				galat(md, err)
				continue
			}
			p.setel(md, v)
			continue
		}
		if md.Pola == jalurNomorPolis {
			if md.Nilai != "" && strings.TrimSpace(md.Nilai) != hasil.NoPolis {
				galat(md, ErrNoPolisBeda)
			}
			continue
		}
		if alasan := alasanDiabaikan(md.Pola); alasan != "" {
			hasil.Diabaikan[alasan]++
			continue
		}
		hasil.TakDikenal = append(hasil.TakDikenal, md)
	}
	if hasil.NoPolis == "" {
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: "NOPOLIS", Err: ErrNoPolisKosong})
	}
	h := p.h
	if h.Ambil(HalamanPolis+".EndDate") == "" && h.Ambil(HalamanPolis+".StartDate") != "" {
		h.Setel(HalamanPolis+".EndDate", h.Ambil(HalamanPolis+".StartDate"))
	}
	// TGL_PROD datar json_polis (ID-21) bila ProductionDate dokumen kosong.
	if h.Ambil(HalamanPolis+".ProductionDate") == "" && b.TglProd != "" {
		h.Setel(HalamanPolis+".ProductionDate", b.TglProd)
	}
	return hasil, nil
}
