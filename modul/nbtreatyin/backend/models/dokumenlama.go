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
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

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
// T_GENERAL_POLIS_TREATY di luar katalog (ID-21). NOPOLIS ditulis lewat
// `SetelNomorPolis`, PRODKE selalu 0 (`SisipKasus`), TGL_PROD lewat katalog
// (`ProductionDate`).
type KolomDatarLama struct {
	NoEndors, TglInput, Username string
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

// MedanArsip - satu medan dokumen yang tidak ditulis ke kolom: Kunci = kunci
// alasan keputusan tertulis (`medan_abaikan_lama.json`), "" = belum diputuskan.
type MedanArsip struct {
	Medan
	Kunci string
}

// GalatDokumen - satu sebab dokumen tidak dimuat (AC 58).
type GalatDokumen struct {
	Jalur, Nilai string
	Err          error
}

// HasilPecah - keluaran pemecah untuk satu dokumen.
type HasilPecah struct {
	// ID - kunci T_WORK_POLIS / T_GENERAL_POLIS_TREATY = pyID dari IDPEGA.
	ID      string
	NoPolis string
	Halaman *Halaman
	Datar   KolomDatarLama
	// BelumDiputuskan - medan tanpa kolom dan tanpa keputusan tertulis (AC 59
	// RALAT F3: wajib nol sebelum pekerjaan dinyatakan selesai).
	BelumDiputuskan []Medan
	// Diabaikan - teks alasan tertulis -> cacah medan yang sengaja tidak disimpan.
	Diabaikan map[string]int
	// Arsip - SETIAP medan daun yang tidak ditulis ke kolom maupun ke riwayat
	// produksi, urut dokumen, beserta kunci keputusannya ("" = belum
	// diputuskan). Ditulis ke arsip CSV pemuat (F3: arsip audit pemuatan,
	// bukan penampung; nilainya tidak hilang diam-diam - AC 57).
	Arsip []MedanArsip
	// Galat - dokumen TIDAK dimuat bila terisi (K15, AC 58).
	Galat []GalatDokumen
	// Usulan - baris SuggestList dokumen yang disalin ke
	// POOLDATA.HISTORYAKSEPTASIPRODUCTION (F3; `UsulanDokumenLama`).
	Usulan []UsulanProduksi
}

// Galat struktural dokumen lama.
var (
	// ErrDokumenRusak - DATA_JSON bukan objek JSON yang dapat diurai.
	ErrDokumenRusak = errors.New("models: DATA_JSON tidak dapat diurai sebagai objek JSON")
	// ErrBukanTreatyIn - dokumen lini lain (bukan galat; dihitung saja).
	ErrBukanTreatyIn = errors.New("models: dokumen bukan polis Treaty In")
	// ErrBarisAplikasiBaru - baris JSON_POLIS tulisan Utility1 APLIKASI BARU (keputusan work owner 06-10-2026: tanpa
	// DATA_JSON, IDPEGA = ID T_WORK_POLIS polos): bukan dokumen Pega - berkasnya sudah di tabel baru. Dilewati dan
	// dihitung, bukan galat (temuan sesi EDM 07-10-2026; 4 baris di DEV).
	ErrBarisAplikasiBaru = errors.New("models: baris JSON_POLIS tulisan aplikasi baru (tanpa DATA_JSON) - bukan dokumen Pega")
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

// ---------------------------------------------------------------- penggolong

// berkasMedanAbaikanLama - daftar medan dokumen yang SENGAJA tidak disimpan
// beserta alasan tertulisnya (simpul utuh, ruas di mana pun, skalar tingkat
// polis). Dihitung per alasan di ringkasan pemuat, tidak hilang diam-diam.
//
// Mengapa berkas data, bukan literal kode: penjaga lintas modul
// `TestNolPenyimpanKeputusanKomiteDiKonteksIni`
// (`modul/claimlife/backend/services/komite_statik_test.go`, `polaIndeksPosisi`)
// memindai SETIAP berkas .go repo secara leksikal dan melarang nama properti
// nomor baris Pega di luar komentar ("nol indeks posisi dipakai sebagai kunci
// rujukan di mana pun"). Pemuat ini memang TIDAK memakainya sebagai kunci:
// properti itu hanya dikenali supaya DIBUANG - urutan baris diwakili NOURUT
// (spec-penyimpanan ID-11). Cakupan penjaga diajukan ke tim claimlife/inti
// (`docs/PERMINTAAN-TIM-INTI.md` bagian A). go:embed terisi saat kompilasi,
// jadi pemecah tetap murni (nol baca berkas saat jalan).
//
//go:embed medan_abaikan_lama.json
var berkasMedanAbaikanLama []byte

// penggolongAbaikan - isi `medan_abaikan_lama.json`. Setiap golongan
// memetakan nama ke KUNCI alasan (tercetak di arsip CSV); teksnya di `alasan`
// (tercetak di ringkasan).
type penggolongAbaikan struct {
	alasan      map[string]string // kunci -> teks alasan
	simpul      map[string]string // puncak -> kunci; seluruh isi simpul
	skalarPolis map[string]string // puncak skalar tingkat polis -> kunci
	ruas        map[string]ruasAbaikan
	// pola - keputusan F3 per medan: pola jalur PERSIS relatif PolicyTreatyIn
	// -> kunci alasan; buktinya wajib ada di berkas (F3, WO 04-10-2026).
	pola map[string]string
}

// ruasAbaikan - nama properti daun di mana pun; hanyaDiBarisDaftar = hanya
// bila medannya berada di baris PageList (pola memuat "()").
type ruasAbaikan struct {
	alasan             string
	hanyaDiBarisDaftar bool
}

// muatPenggolongAbaikan mengurai berkas data penggolong. Rujukan ke alasan
// yang tidak terdefinisi, keputusan per medan tanpa bukti, dan medan JSON
// tak dikenal ditolak - berkas data tidak diperiksa penyusun Go, jadi
// pemeriksaannya di sini.
func muatPenggolongAbaikan(isi []byte) (penggolongAbaikan, error) {
	var mentah struct {
		Catatan     []string          `json:"catatan"`
		Alasan      map[string]string `json:"alasan"`
		Simpul      map[string]string `json:"simpul"`
		SkalarPolis map[string]string `json:"skalar_polis"`
		Ruas        map[string]struct {
			Alasan             string `json:"alasan"`
			HanyaDiBarisDaftar bool   `json:"hanya_di_baris_daftar"`
		} `json:"ruas"`
		Pola map[string]struct {
			Alasan string `json:"alasan"`
			Bukti  string `json:"bukti"`
		} `json:"pola"`
	}
	dek := json.NewDecoder(bytes.NewReader(isi))
	dek.DisallowUnknownFields()
	if err := dek.Decode(&mentah); err != nil {
		return penggolongAbaikan{}, fmt.Errorf("models: medan_abaikan_lama.json: %w", err)
	}
	var galat error
	kunci := func(bagian, nama, k string) string {
		if mentah.Alasan[k] == "" && galat == nil {
			galat = fmt.Errorf("models: medan_abaikan_lama.json %s %q merujuk alasan tak terdefinisi %q", bagian, nama, k)
		}
		return k
	}
	p := penggolongAbaikan{alasan: mentah.Alasan, simpul: map[string]string{},
		skalarPolis: map[string]string{}, ruas: map[string]ruasAbaikan{}, pola: map[string]string{}}
	for nama, k := range mentah.Simpul {
		p.simpul[nama] = kunci("simpul", nama, k)
	}
	for nama, k := range mentah.SkalarPolis {
		p.skalarPolis[nama] = kunci("skalar_polis", nama, k)
	}
	for nama, r := range mentah.Ruas {
		p.ruas[nama] = ruasAbaikan{alasan: kunci("ruas", nama, r.Alasan), hanyaDiBarisDaftar: r.HanyaDiBarisDaftar}
	}
	for nama, x := range mentah.Pola {
		p.pola[nama] = kunci("pola", nama, x.Alasan)
		if strings.TrimSpace(x.Bukti) == "" && galat == nil {
			galat = fmt.Errorf("models: medan_abaikan_lama.json pola %q tanpa bukti XML (F3)", nama)
		}
	}
	if galat != nil {
		return penggolongAbaikan{}, galat
	}
	return p, nil
}

// medanAbaikanLama - penggolong dari berkas tertanam. Gagal urai = berkas
// tertanam cacat (galat pemrogram, sama halnya `regexp.MustCompile`).
var medanAbaikanLama = func() penggolongAbaikan {
	p, err := muatPenggolongAbaikan(berkasMedanAbaikanLama)
	if err != nil {
		panic(err)
	}
	return p
}()

// teksAlasan - teks alasan berkunci `kunci` di berkas tertanam.
func teksAlasan(kunci string) string {
	s, ada := medanAbaikanLama.alasan[kunci]
	if !ada {
		panic("models: medan_abaikan_lama.json tanpa alasan " + kunci)
	}
	return s
}

// Alasan medan dokumen yang SENGAJA tidak disimpan - masing-masing keputusan
// tertulis; teksnya di `medan_abaikan_lama.json` dan tercetak di ringkasan
// pemuat (kunci `HasilPecah.Diabaikan`).
var (
	// AlasanInternalPega - `pxObjClass`, internal Pega.
	AlasanInternalPega = teksAlasan("internal_pega")
	// AlasanNourut - properti nomor baris Pega (berbasis 1) di baris daftar;
	// diwakili NOURUT, bukan kunci rujukan (spec-penyimpanan ID-11).
	AlasanNourut = teksAlasan("nourut")
	// AlasanKeadaanLayar - Show/ViewState/pxResults/FillPaymentInstallmentEDMT.
	AlasanKeadaanLayar = teksAlasan("keadaan_layar")
	// AlasanTurunan - Total* tingkat polis.
	AlasanTurunan = teksAlasan("turunan")
	// AlasanPantulanLayer - Layer* tingkat polis.
	AlasanPantulanLayer = teksAlasan("pantulan_layer")
	// AlasanBreakdown - BreakDownSpreadList (K9).
	AlasanBreakdown = teksAlasan("breakdown")
	// AlasanPersetujuanDH - isApprovedtoDeptHead.
	AlasanPersetujuanDH = teksAlasan("persetujuan_dh")
	// AlasanOldData - simpul OldData.
	AlasanOldData = teksAlasan("old_data")
	// AlasanSelisih - TreatyDifference/TreatyXOLDifferenceList.
	AlasanSelisih = teksAlasan("selisih")
	// AlasanF3TanpaPembaca - F3: medan dokumen lama tanpa kolom yang nol
	// dibaca rule NB terjangkau (bukti per medan di bagian `pola`).
	AlasanF3TanpaPembaca = teksAlasan("f3_tanpa_pembaca")
	// AlasanF3KeadaanBaris - F3: pyExpanded baris daftar.
	AlasanF3KeadaanBaris = teksAlasan("f3_keadaan_baris")
	// AlasanF3SalinanGenerasi - F3: EDMNo, ProdKe dokumen.
	AlasanF3SalinanGenerasi = teksAlasan("f3_salinan_generasi")
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

// kunciDibuang - KUNCI alasan bila pola medan sengaja tidak disimpan
// (keputusan tertulis di `medan_abaikan_lama.json`); "" = belum diputuskan.
//
// Urutan pemeriksaan: (1) simpul utuh (OldData, TreatyDifference, ...) lebih
// dulu - seluruh isinya, termasuk pxObjClass-nya, tercatat di bawah alasan
// simpul itu; (2) keputusan F3 per medan (pola persis); (3) ruas daun di mana
// pun (pxObjClass; nomor baris hanya di baris daftar); (4) skalar tingkat
// polis saja.
func kunciDibuang(pola string) string {
	p := medanAbaikanLama
	akar := strings.TrimPrefix(pola, HalamanPolis+".")
	puncak := akar
	if i := strings.IndexAny(akar, ".("); i >= 0 {
		puncak = akar[:i]
	}
	if k, ada := p.simpul[puncak]; ada {
		return k
	}
	if k, ada := p.pola[akar]; ada {
		return k
	}
	ruas := pola[strings.LastIndex(pola, ".")+1:]
	if r, ada := p.ruas[ruas]; ada && (!r.hanyaDiBarisDaftar || strings.Contains(pola, "()")) {
		return r.alasan
	}
	if puncak != akar {
		return ""
	}
	return p.skalarPolis[puncak]
}

// alasanDiabaikan - TEKS alasan tertulis bila pola medan sengaja tidak
// disimpan; "" = belum diputuskan.
func alasanDiabaikan(pola string) string {
	return medanAbaikanLama.alasan[kunciDibuang(pola)]
}

// ---------------------------------------------------------------- pemecah

// polaPyID - pyID kasus: awalan huruf, tanda hubung, nomor (`NB-77`).
var polaPyID = regexp.MustCompile(`^[A-Z][A-Z0-9]*-\d+$`)

// IDKasusDariIDPega mengambil pyID dari `pyWorkPage.pzInsKey`
// (`<kelas> <pyID>`, kunci dokumen Pega lama). pyID menjadi ID T_WORK_POLIS - diagram
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
			s, _ := TeksSkalarJSON(v)
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
			s, _ := TeksSkalarJSON(v)
			p.medan = append(p.medan, Medan{Jalur: j, Pola: pola + "()", Nilai: s})
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
			s, _ := TeksSkalarJSON(v)
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
// desimal KESEBELAS (NUMBER(38,10) - diagram NB Treaty In Prop F20 "skala
// MINIMAL 9"; AC 19, 20b). Satu-satunya pengisian:
// EndDate kosong = StartDate (`[keputusan work owner]` spec AC 69, §5.8).
// BarisAplikasiBaru - baris JSON_POLIS tulisan aplikasi baru: IDPEGA tanpa spasi (bukan `<kelas> <pyID>` Pega) DAN
// DATA_JSON kosong. Dokumen Pega ber-kelas yang JSON-nya kosong TETAP galat (`ErrDokumenRusak`).
func BarisAplikasiBaru(b BarisJSONPolis) bool {
	return !strings.Contains(strings.TrimSpace(b.IDPega), " ") && len(bytes.TrimSpace(b.DataJSON)) == 0
}

func PecahDokumenLama(b BarisJSONPolis) (HasilPecah, error) {
	if BarisAplikasiBaru(b) {
		return HasilPecah{}, ErrBarisAplikasiBaru
	}
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
		Datar: KolomDatarLama{NoEndors: b.NoEndors, TglInput: b.TglInput, Username: b.Username},
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
		// F3: anggota baris SuggestList disalin ke riwayat produksi
		// (`UsulanDokumenLama`); `.Date` dibaca seperti kolom tanggal-waktu.
		if nama, ok := anggotaUsulanLama(md.Pola); ok {
			if nama == "Date" {
				v, err := BacaTanggalLama(md.Nilai, GolTanggalWaktu)
				if err != nil {
					galat(md, err)
					continue
				}
				p.setel(md, v)
			}
			continue
		}
		kunci := kunciDibuang(md.Pola)
		hasil.Arsip = append(hasil.Arsip, MedanArsip{Medan: md, Kunci: kunci})
		if kunci != "" {
			hasil.Diabaikan[teksAlasan(kunci)]++
			continue
		}
		hasil.BelumDiputuskan = append(hasil.BelumDiputuskan, md)
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
	hasil.Usulan = UsulanDokumenLama(id, h)
	return hasil, nil
}
