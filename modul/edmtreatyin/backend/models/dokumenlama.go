package models

// Untuk apa berkas ini: PEMECAH DOKUMEN LAMA ENDORSEMEN - tiket EDM 10 (spec-penyimpanan-relasional.md ID-3, ID-8,
// ID-25, ID-33, AC 39, AC 44; KEPUTUSAN 23-09-2026 butir 4: setiap polis, setiap generasi). Asal: adaptasi
// `modul/nbtreatyin/backend/models/dokumenlama.go` (tiket NB 22, `PecahDokumenLama`) - NB menolak generasi
// `PRODKE > 0` dengan `ErrGenerasiEndorsemen` "milik pemuat EDM tiket 10"; berkas ini pemecah generasi itu.
//
// Dokumen lama = satu baris `POOLDATA.JSON_POLIS` ber-PRODKE > 0. `[terverifikasi]`
// `Activity/SaveJsonPolisTreatyInEDM_Act.xml` (Utility1 kasus `ASM-FW-GISFW-Work-EndorsementTreaty`):
//
//	langkah "Copy pzInsKey"   InputData.CARI1 = pyWorkPage.pzInsKey; CARI21 = ProductionDate (dd/MM/yyyy hh:mm:ss);
//	                          CARI2 = PolicyTreatyIn.PolicyNo; CARI4 = PolicyTreatyIn.EDMNo
//	langkah "Copy datajson"   halaman langkah pyWorkPage.PolicyTreatyIn (kelas Data-PolicyTreatyIn):
//	                          InputData.CARI3 = @ASM.GetPageJSONString() - DATA_JSON = HALAMAN PolicyTreatyIn
//	RDB TreatyInSearchProdKe  select count(*) as CARI1 from pooldata.json_polis where Nopolis = {CARI2} -> CARI5
//	RDB SavePolisTreatyInEDM_SQL  PEGA_JSON_POLIS_TREATYIN({pzInsKey}, {PolicyNo}, {EDMNo}, {CARI5}, {CARI21},
//	                          {OperatorID.pyUserIdentifier}, {CARI3}) - IDPEGA, NOPOLIS, NOENDORS, PRODKE, TGL_PROD,
//	                          USERNAME, DATA_JSON (urutan argumen = procedure NB, PERTANYAAN-untuk-DBA P1 §2)
//
// ⇒ PRODKE = cacah baris json_polis bernomor polis itu SEBELUM generasi ini ditulis (NB '0' + endorsemen
// sebelumnya); NOENDORS = EDMNo; NOPOLIS = nomor polis INDUK (SetEDMTNoPolis langkah 3: PolicyNo = OldData.PolicyNo).
//
// Isi dokumen dipetakan MENURUT KATALOG yang sama dengan jalur biasa - katalog polis (`katalog.go`) dan katalog
// proyeksi selisih (`katalog_selisih.go`): `TreatyDifference.*` dan `TreatyXOLDifferenceList().ValueList().*`
// masuk proyeksi SUMBER 'PEGA' APA ADANYA (ID-33, AC 39: tidak dihitung ulang, digit galat utuh). `OldData` TIDAK
// dipetakan (ID-8): generasi sebelumnya dibaca lewat OLD_POLIS_ID; hanya `OldData.EDMNo` dibaca untuk memeriksa
// pemilih varian rumus berlapis (`penanda_migrasi.go`).
//
// Pemecah ini MURNI: tidak menyentuh basis data, berkas, maupun jam.

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

// KelasDokumenTreatyIn - `pxObjClass` akar DATA_JSON (halaman langkah "Copy datajson",
// `pyStepsClassName` ASM-FW-GISFW-Data-PolicyTreatyIn). JSON_POLIS dipakai bersama lini lain: kelas lain di
// luar lingkup pemuat ini dan hanya dihitung.
const KelasDokumenTreatyIn = "ASM-FW-GISFW-Data-PolicyTreatyIn"

// KelasKerjaEDM - kelas work object endorsemen treaty (`Activity/CreateEDMT` langkah 18 CreateWorkPage
// `ASM-FW-GISFW-Work-EndorsementTreaty`). IDPEGA = `pyWorkPage.pzInsKey` = "<kelas> <pyID>"; huruf kelas di kunci
// instans Pega ditulis kapital, jadi dibandingkan tanpa beda huruf.
const KelasKerjaEDM = "ASM-FW-GISFW-Work-EndorsementTreaty"

// BarisJSONPolis - satu baris POOLDATA.JSON_POLIS: kolom datar dan dokumennya. Tanggal datar berbentuk
// `utils.TanggalWaktu` (dibaca repository lewat TO_CHAR); kosong = NULL.
type BarisJSONPolis struct {
	IDPega, NoPolis, NoEndors, ProdKe string
	TglInput, TglProd, Username       string
	DataJSON                          []byte
}

// KunciJSONPolis - kunci baca satu baris JSON_POLIS endorsemen beserta kunci urutnya (NOPOLIS, PRODKE): generasi
// dimuat menurut nomornya (tiket 10 "urutan pemuatan penting").
type KunciJSONPolis struct {
	Kunci, NoPolis, ProdKe string
}

// UrutKunciGenerasi mengurutkan kunci menurut NOPOLIS lalu PRODKE BILANGAN naik (bukan teks: json_polis.PRODKE
// VARCHAR2 di DEV - "9" > "10", catatan `repository/generasi.go`); PRODKE tak terbaca di belakang nomor polisnya
// (dilaporkan pemecah), lalu kunci baca supaya urutan pasti.
func UrutKunciGenerasi(k []KunciJSONPolis) {
	sort.SliceStable(k, func(i, j int) bool {
		a, b := k[i], k[j]
		if a.NoPolis != b.NoPolis {
			return a.NoPolis < b.NoPolis
		}
		na, ea := strconv.Atoi(strings.TrimSpace(a.ProdKe))
		nb, eb := strconv.Atoi(strings.TrimSpace(b.ProdKe))
		switch {
		case ea == nil && eb == nil && na != nb:
			return na < nb
		case (ea == nil) != (eb == nil):
			return ea == nil
		}
		return a.Kunci < b.Kunci
	})
}

// KolomDatarLama - kolom datar json_polis yang ditulis apa adanya ke T_GENERAL_POLIS_TREATY di luar katalog
// (pola NB `SetelKolomDatarLama`). NOPOLIS lewat `SetelNomorPolisSelesai`, PRODKE / NOENDORS lewat `SisipKasus`,
// TGL_PROD lewat katalog (`ProductionDate`).
type KolomDatarLama struct {
	TglInput, Username string
}

// KunciGenerasi - kunci generasi yang sudah tersimpan (pembeda "sudah dimuat" dari "ID bentrok").
type KunciGenerasi struct {
	NoPolis  string
	ProdKe   int
	NoEndors string
}

// MedanDokumen - satu medan daun dokumen. Jalur memakai notasi properti Pega relatif `pyWorkPage`
// (`PolicyTreatyIn.ListInstallment(2).Premium`); Pola = jalur tanpa nomor baris (`...ListInstallment().Premium`).
type MedanDokumen struct {
	Jalur, Pola, Nilai string
	// letak di halaman: daftar == "" berarti nilai halaman `Jalur`.
	daftar string
	baris  int
	nama   string
}

// MedanArsip - satu medan dokumen yang tidak ditulis ke kolom: Kunci = kunci alasan keputusan tertulis
// (`medan_abaikan_lama.json`), "" = belum diputuskan.
type MedanArsip struct {
	MedanDokumen
	Kunci string
}

// GalatDokumen - satu sebab dokumen tidak dimuat.
type GalatDokumen struct {
	Jalur, Nilai string
	Err          error
}

// HasilPecahEDM - keluaran pemecah untuk satu dokumen endorsemen.
type HasilPecahEDM struct {
	// ID - kunci T_WORK_POLIS / T_GENERAL_POLIS_TREATY = pyID `EDMT-<n>` dari IDPEGA.
	ID string
	// NoPolis - JSON_POLIS.NOPOLIS = nomor polis induk.
	NoPolis string
	// ProdKe - JSON_POLIS.PRODKE (>= 1) - nomor generasi.
	ProdKe int
	// EDMNo - NOENDORS (kolom datar json_polis; sama dengan `PolicyTreatyIn.EDMNo` dokumen).
	EDMNo string
	// EDMType - `PolicyTreatyIn.EDMType` -> EDM_TYPE (SisipKasus dan katalog).
	EDMType string
	// OldDataEDMNo - `PolicyTreatyIn.OldData.EDMNo` dokumen: pemilih varian rumus Pega
	// (`EDMTCalculateTreatyDifference` langkah 1 `.OldData.EDMNo==""`). Hanya diperiksa, tidak disimpan.
	OldDataEDMNo string
	// Halaman - halaman generasi siap `SimpanHalaman` + `SimpanSelisih` (tanpa OldData).
	Halaman *Halaman
	Datar   KolomDatarLama
	// BelumDiputuskan - medan tanpa kolom dan tanpa keputusan tertulis (F3: wajib nol sebelum selesai).
	BelumDiputuskan []MedanDokumen
	// Diabaikan - teks alasan tertulis -> cacah medan yang sengaja tidak disimpan.
	Diabaikan map[string]int
	// Arsip - SETIAP medan daun yang tidak ditulis ke kolom maupun ke riwayat produksi, urut dokumen, beserta
	// kunci keputusannya ("" = belum diputuskan).
	Arsip []MedanArsip
	// Galat - dokumen TIDAK dimuat bila terisi.
	Galat []GalatDokumen
	// Usulan - baris SuggestList yang disalin ke POOLDATA.HISTORYAKSEPTASIPRODUCTION (`usulanlama.go`).
	Usulan []UsulanProduksi
}

// Galat dokumen lama endorsemen. Yang "lewat" (bukan galat): ErrBukanTreatyIn, ErrBukanGenerasiEndorsemen.
var (
	// ErrDokumenRusak - DATA_JSON bukan objek JSON yang dapat diurai.
	ErrDokumenRusak = errors.New("models: DATA_JSON tidak dapat diurai sebagai objek JSON")
	// ErrBukanTreatyIn - dokumen lini lain (dihitung saja).
	ErrBukanTreatyIn = errors.New("models: dokumen bukan polis Treaty In")
	// ErrBarisAplikasiBaru - baris json_polis yang ditulis Utility1 APLIKASI BARU (`repository/produksi.go`: IDPEGA =
	// ID kasus polos `EDMT-<n>`, DATA_JSON tidak ditulis). Bukan dokumen Pega lama - generasinya sudah di tabel flat
	// (dihitung saja; ditemukan popup Copy Old DEV 07-10-2026).
	ErrBarisAplikasiBaru = errors.New("models: baris json_polis tulisan aplikasi baru (tanpa DATA_JSON), bukan dokumen Pega")
	// ErrBukanGenerasiEndorsemen - PRODKE 0: generasi NB, milik pemuat NB tiket 22 (dihitung saja).
	ErrBukanGenerasiEndorsemen = errors.New("models: generasi NB (PRODKE 0) - milik pemuat NB tiket 22")
	// ErrProdKe - PRODKE kosong atau bukan bilangan.
	ErrProdKe = errors.New("models: PRODKE kosong atau bukan bilangan")
	// ErrIDPega - IDPEGA bukan "<ASM-FW-GISFW-Work-EndorsementTreaty> EDMT-<n>": dokumen Treaty In ber-PRODKE > 0
	// yang bukan kasus endorsemen (korpus: hanya SavePolisTreatyInEDM_SQL yang menulisnya) - tidak ditebak.
	ErrIDPega = errors.New("models: IDPEGA tidak berbentuk \"ASM-FW-GISFW-Work-EndorsementTreaty EDMT-<nomor>\"")
	// ErrNilaiKolom - nilai tidak dapat ditulis ke kolomnya (aturan sama dengan konversi repository).
	ErrNilaiKolom = errors.New("models: nilai tidak sesuai tipe kolomnya")
	// ErrNoPolisKosong - JSON_POLIS.NOPOLIS kosong.
	ErrNoPolisKosong = errors.New("models: NOPOLIS kosong")
	// ErrNoPolisBeda - PolicyNo dokumen berbeda dari JSON_POLIS.NOPOLIS (keduanya `PolicyTreatyIn.PolicyNo`).
	ErrNoPolisBeda = errors.New("models: PolicyNo dokumen berbeda dari JSON_POLIS.NOPOLIS")
	// ErrProdKeBeda - ProdKe dokumen (SetEDMTNoPolis langkah 3) berbeda dari JSON_POLIS.PRODKE (cacah
	// TreatyInSearchProdKe) - nomor generasi mana yang benar tidak ditebak.
	ErrProdKeBeda = errors.New("models: ProdKe dokumen berbeda dari JSON_POLIS.PRODKE")
	// ErrEDMNoKosong - NOENDORS dan EDMNo dokumen sama-sama kosong.
	ErrEDMNoKosong = errors.New("models: nomor endorsemen kosong (NOENDORS dan EDMNo dokumen)")
	// ErrEDMNoBeda - EDMNo dokumen berbeda dari JSON_POLIS.NOENDORS (keduanya `PolicyTreatyIn.EDMNo`).
	ErrEDMNoBeda = errors.New("models: EDMNo dokumen berbeda dari JSON_POLIS.NOENDORS")
	// ErrDokumenGanda - dua baris JSON_POLIS untuk kasus yang sama dalam satu jalankan; hanya yang pertama dimuat.
	ErrDokumenGanda = errors.New("models: dokumen ganda untuk kasus yang sama")
	// ErrGenerasiSebelumnyaTidakAda - generasi bernomor polis sama ber-PRODKE - 1 belum dimuat (NB: pemuat NB tiket
	// 22; endorsemen: dokumen sebelumnya gagal). OLD_POLIS_ID tidak ditebak.
	ErrGenerasiSebelumnyaTidakAda = errors.New("models: generasi sebelumnya (PRODKE - 1) belum dimuat")
	// ErrPercabangan - generasi sebelumnya sudah punya penerus: UNIQUE OLD_POLIS_ID (spec-penyimpanan ID-10, AC 3).
	ErrPercabangan = errors.New("models: percabangan ditolak - generasi sebelumnya sudah punya penerus (UNIQUE OLD_POLIS_ID)")
	// ErrKeutuhan - generasi baru kehilangan baris spreading generasi sebelumnya (ID-15, AC 8).
	ErrKeutuhan = errors.New("models: keutuhan ditolak - baris spreading generasi sebelumnya hilang (ID-15)")
	// ErrOldDataTakSesuai - OldData.EDMNo dokumen (pemilih varian rumus Pega) tidak sesuai generasi OLD_POLIS_ID:
	// penanda rumus berlapis tidak dapat ditentukan pasti (ID-36).
	ErrOldDataTakSesuai = errors.New("models: OldData.EDMNo dokumen tidak sesuai nomor endorsemen generasi sebelumnya")
	// ErrIDKasusBentrok - ID kasus sudah dipakai generasi lain (nomor polis / PRODKE / NOENDORS berbeda); dokumen
	// tidak dimuat, tidak menimpa.
	ErrIDKasusBentrok = errors.New("models: ID kasus sudah dipakai generasi lain")
)

// ---------------------------------------------------------------- penggolong

// pmBerkasMedanAbaikan - daftar medan dokumen yang SENGAJA tidak disimpan beserta alasan tertulisnya. Berkas data,
// bukan literal kode: nama properti nomor baris Pega tidak boleh muncul di luar komentar berkas .go (penjaga lintas
// modul claimlife `polaIndeksPosisi`). go:embed terisi saat kompilasi - pemecah tetap murni.
//
//go:embed medan_abaikan_lama.json
var pmBerkasMedanAbaikan []byte

// pmPenggolong - isi `medan_abaikan_lama.json`.
type pmPenggolong struct {
	alasan      map[string]string // kunci -> teks alasan
	simpul      map[string]string // puncak -> kunci; seluruh isi simpul
	skalarPolis map[string]string // puncak skalar tingkat polis -> kunci
	ruas        map[string]pmRuas
	pola        map[string]pmPola
}

// pmRuas - nama properti daun di mana pun; hanyaDiBarisDaftar = hanya bila medannya berada di baris PageList.
type pmRuas struct {
	alasan             string
	hanyaDiBarisDaftar bool
}

// pmPola - keputusan per medan; hanyaBilaKosong = keputusan bersandar pada "medan ini selalu kosong".
type pmPola struct {
	alasan          string
	hanyaBilaKosong bool
}

// pmMuatPenggolong mengurai berkas data penggolong. Rujukan ke alasan tak terdefinisi, keputusan per medan tanpa
// bukti, ruas F3 tanpa bukti, dan medan JSON tak dikenal ditolak.
func pmMuatPenggolong(isi []byte) (pmPenggolong, error) {
	var mentah struct {
		Catatan     []string          `json:"catatan"`
		Alasan      map[string]string `json:"alasan"`
		Simpul      map[string]string `json:"simpul"`
		SkalarPolis map[string]string `json:"skalar_polis"`
		Ruas        map[string]struct {
			Alasan             string `json:"alasan"`
			HanyaDiBarisDaftar bool   `json:"hanya_di_baris_daftar"`
		} `json:"ruas"`
		BuktiRuas map[string]string `json:"bukti_ruas"`
		Pola      map[string]struct {
			Alasan          string `json:"alasan"`
			Bukti           string `json:"bukti"`
			HanyaBilaKosong bool   `json:"hanya_bila_kosong"`
		} `json:"pola"`
	}
	dek := json.NewDecoder(bytes.NewReader(isi))
	dek.DisallowUnknownFields()
	if err := dek.Decode(&mentah); err != nil {
		return pmPenggolong{}, fmt.Errorf("models: medan_abaikan_lama.json: %w", err)
	}
	var galat error
	kunci := func(bagian, nama, k string) string {
		if mentah.Alasan[k] == "" && galat == nil {
			galat = fmt.Errorf("models: medan_abaikan_lama.json %s %q merujuk alasan tak terdefinisi %q", bagian, nama, k)
		}
		return k
	}
	p := pmPenggolong{alasan: mentah.Alasan, simpul: map[string]string{}, skalarPolis: map[string]string{},
		ruas: map[string]pmRuas{}, pola: map[string]pmPola{}}
	for nama, k := range mentah.Simpul {
		p.simpul[nama] = kunci("simpul", nama, k)
	}
	for nama, k := range mentah.SkalarPolis {
		p.skalarPolis[nama] = kunci("skalar_polis", nama, k)
	}
	for nama, r := range mentah.Ruas {
		p.ruas[nama] = pmRuas{alasan: kunci("ruas", nama, r.Alasan), hanyaDiBarisDaftar: r.HanyaDiBarisDaftar}
		if strings.HasPrefix(r.Alasan, "f3_") && strings.TrimSpace(mentah.BuktiRuas[nama]) == "" && galat == nil {
			galat = fmt.Errorf("models: medan_abaikan_lama.json ruas F3 %q tanpa bukti XML", nama)
		}
	}
	for nama, x := range mentah.Pola {
		p.pola[nama] = pmPola{alasan: kunci("pola", nama, x.Alasan), hanyaBilaKosong: x.HanyaBilaKosong}
		if strings.TrimSpace(x.Bukti) == "" && galat == nil {
			galat = fmt.Errorf("models: medan_abaikan_lama.json pola %q tanpa bukti XML (F3)", nama)
		}
	}
	if galat != nil {
		return pmPenggolong{}, galat
	}
	return p, nil
}

// pmMedanAbaikan - penggolong dari berkas tertanam. Gagal urai = berkas tertanam cacat (galat pemrogram).
var pmMedanAbaikan = func() pmPenggolong {
	p, err := pmMuatPenggolong(pmBerkasMedanAbaikan)
	if err != nil {
		panic(err)
	}
	return p
}()

// pmTeksAlasan - teks alasan berkunci `kunci` di berkas tertanam.
func pmTeksAlasan(kunci string) string {
	s, ada := pmMedanAbaikan.alasan[kunci]
	if !ada {
		panic("models: medan_abaikan_lama.json tanpa alasan " + kunci)
	}
	return s
}

// pmKunciPropHasilAntara - alasan daftar hasil antara generasi Proportional (dipilih di Go menurut
// ProportionalType dokumen, bukan oleh berkas data).
const pmKunciPropHasilAntara = "prop_hasil_antara"

var _ = pmTeksAlasan(pmKunciPropHasilAntara) // berkas data wajib memuatnya

// pmKunciDibuang - KUNCI alasan bila medan sengaja tidak disimpan; "" = belum diputuskan. Urutan: simpul utuh,
// keputusan per medan (pola persis), ruas daun, skalar tingkat polis.
func pmKunciDibuang(pola, nilai string) string {
	p := pmMedanAbaikan
	akar := strings.TrimPrefix(pola, HalamanPolis+".")
	puncak := akar
	if i := strings.IndexAny(akar, ".("); i >= 0 {
		puncak = akar[:i]
	}
	if k, ada := p.simpul[puncak]; ada {
		return k
	}
	if x, ada := p.pola[akar]; ada {
		if x.hanyaBilaKosong && strings.TrimSpace(nilai) != "" {
			return ""
		}
		return x.alasan
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

// pmDaftarHasilAntaraProp - daftar hasil antara generasi Proportional yang dibuang `RapikanBentukSimpan`.
var pmDaftarHasilAntaraProp = []string{
	TabelXOL.Daftar + "()",
	DaftarSelisihXOL + "()",
	DaftarAngsuran + "()." + TabelAngsuranRinci.Daftar + "()",
}

// pmHasilAntaraProp - pola medan berada di daftar hasil antara generasi Proportional.
func pmHasilAntaraProp(pola string) bool {
	for _, d := range pmDaftarHasilAntaraProp {
		if strings.HasPrefix(pola, d+".") {
			return true
		}
	}
	return false
}

// pmPetaKatalogDokumen menurunkan pola jalur dokumen -> kolom dari KATALOG polis DAN katalog selisih. Tabel 1:1
// memakai jalur properti; T_POLIS_QUOTATION dari `PolicyTreatyIn.QuotationData.<medan>`; tabel anak dari
// `<Daftar>().<medan>`; tabel cucu dari `<induk>().<Daftar>().<medan>`; lapisan selisih XOL dari
// `TreatyXOLDifferenceList().ValueList().<medan>`. Medan halaman kerja (PositionNote, NBStatus, TreatyIn.ID) bukan
// isi dokumen.
func pmPetaKatalogDokumen() map[string]Kolom {
	induk := map[string]string{}
	for _, kt := range keturunanKatalog {
		if kt.cucu != nil {
			induk[kt.cucu.Nama] = kt.anak.Daftar
		}
	}
	induk[TabelSelisihLapisan.Nama] = DaftarSelisihXOL
	peta := map[string]Kolom{}
	semua := append(append([]Tabel{}, SemuaTabel...), SemuaTabelSelisih...)
	for _, t := range semua {
		var awalan string
		switch {
		case t.Daftar == "" && t.Nama == TabelQuotation.Nama:
			awalan = HalamanPolis + ".QuotationData."
		case t.Daftar == "":
			awalan = ""
		case strings.HasPrefix(t.Daftar, HalamanPolis+"."):
			awalan = t.Daftar + "()."
		default:
			i, ada := induk[t.Nama]
			if !ada {
				panic("models: daftar bersarang " + t.Daftar + " (" + t.Nama + ") tanpa induk")
			}
			awalan = i + "()." + t.Daftar + "()."
		}
		for _, k := range t.Kolom {
			pola := awalan + k.Properti
			if strings.HasPrefix(pola, HalamanPolis+".") {
				peta[pola] = k
			}
		}
	}
	return peta
}

// pmPetaKatalog - hasil `pmPetaKatalogDokumen` (katalog tetap selama proses).
var pmPetaKatalog = pmPetaKatalogDokumen()

// Jalur kunci generasi di dokumen (bukan kolom katalog).
const (
	pmJalurNomorPolis = HalamanPolis + ".PolicyNo"
	pmJalurEDMNo      = HalamanPolis + ".EDMNo"
	pmJalurProdKe     = HalamanPolis + ".ProdKe"
	pmJalurOldEDMNo   = HalamanPolis + ".OldData.EDMNo"
)

// ---------------------------------------------------------------- pemecah

// pmPolaPyIDEDM - pyID kasus endorsemen: `EDMT-<nomor>` (`pyWorkIDPrefix=="EDMT-"`, EDMChooseBusiness_Act 11-12).
var pmPolaPyIDEDM = regexp.MustCompile(`^` + regexp.QuoteMeta(AwalanKasus) + `\d+$`)

// BarisAplikasiBaru - baris json_polis tulisan Utility1 aplikasi baru: IDPEGA tanpa kelas Pega (tanpa spasi; Pega selalu
// `<kelas> <pyID>`) DAN DATA_JSON kosong. Dokumen Pega ber-kelas tanpa JSON tetap galat (`ErrDokumenRusak`).
func BarisAplikasiBaru(b BarisJSONPolis) bool {
	id := strings.TrimSpace(b.IDPega)
	return id != "" && !strings.Contains(id, " ") && len(bytes.TrimSpace(b.DataJSON)) == 0
}

// KelasGrupKerjaPega - kelas di depan pyID pada `pzInsKey` kasus Pega. DEV 07-10-2026 (DATAPEGA.PC_ASM_FW_GISFW_WORK):
// 4 kasus EDMT dan 264 kasus NB berkunci `ASM-FW-GISFW-WORK <pyID>` - kelas GRUP, bukan kelas kerjanya
// (`KelasKerjaEDM` tetap PXOBJCLASS-nya).
const KelasGrupKerjaPega = "ASM-FW-GISFW-WORK"

// PanjangIDKasus - lebar kolom ID kasus (T_WORK_POLIS.ID, T_GENERAL_POLIS_TREATY.ID / OLD_POLIS_ID, T_POLIS_*.POLIS_ID:
// VARCHAR2(32) di DEV 07-10-2026). pzInsKey kasus Pega terpanjang di DEV: 27 karakter.
const PanjangIDKasus = 32

// PyIDKasus - pyID sebuah ID kasus: ID salinan dokumen Pega = pzInsKey utuh (`<kelas> <pyID>`), ID aplikasi baru = pyID.
func PyIDKasus(id string) string {
	s := strings.TrimSpace(id)
	return s[strings.LastIndex(s, " ")+1:]
}

// IDKasusDariIDPegaEDM - ID kasus salinan dokumen Pega = IDPEGA (`pyWorkPage.pzInsKey`) UTUH, tidak dipotong.
// ⛔ Perintah work owner 07-10-2026: "IDPEGA BAWAAN PEGA JANGAN DI POTONG, BERLAKU UNTUK SEMUA NB TREATY DAN EDM
// TREATY" - riwayat HISTORYAKSEPTASIPEGA / HISTORYAKSEPTASIPRODUCTION berkas Pega berkunci pzInsKey yang sama
// (`KunciInstans`). Diperiksa, tidak ditebak (ErrIDPega): kelas = `KelasGrupKerjaPega` (tanpa beda huruf; RALAT
// 07-10-2026 - dulu `KelasKerjaEDM`, yang tidak pernah dipakai pzInsKey), pyID `EDMT-<nomor>`, panjang <=
// `PanjangIDKasus`.
func IDKasusDariIDPegaEDM(idpega string) (string, error) {
	s := strings.TrimSpace(idpega)
	i := strings.LastIndex(s, " ")
	if i <= 0 {
		return "", fmt.Errorf("%w: %q", ErrIDPega, idpega)
	}
	kelas, py := strings.TrimSpace(s[:i]), s[i+1:]
	if !strings.EqualFold(kelas, KelasGrupKerjaPega) || !pmPolaPyIDEDM.MatchString(py) || len(s) > PanjangIDKasus {
		return "", fmt.Errorf("%w: %q", ErrIDPega, idpega)
	}
	return s, nil
}

// pmPeriksaNilai menolak nilai yang akan ditolak konversi repository (`repository.nilaiTulis`) - uji-kering
// melaporkan galat yang sama tanpa basis data.
func pmPeriksaNilai(k Kolom, teks string) error {
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

// pmPejalan menelusuri pohon JSON menjadi halaman kerja dan daftar medan daun.
type pmPejalan struct {
	h     *Halaman
	medan []MedanDokumen
}

func pmKunciUrut(m map[string]any) []string {
	k := make([]string, 0, len(m))
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	return k
}

// objek - simpul objek di tingkat halaman (akar, QuotationData, OldData, TreatyDifference, ...).
func (p *pmPejalan) objek(jalur, pola string, m map[string]any) {
	for _, n := range pmKunciUrut(m) {
		j, q := jalur+"."+n, pola+"."+n
		switch v := m[n].(type) {
		case map[string]any:
			p.objek(j, q, v)
		case []any:
			p.daftar(j, q, v)
		default:
			s, _ := TeksSkalarJSON(v)
			p.h.Setel(j, s)
			p.medan = append(p.medan, MedanDokumen{Jalur: j, Pola: q, Nilai: s})
		}
	}
}

// daftar - PageList: satu Baris per unsur; daftar bersarang di jalur `JalurAnak(daftar, n, nama)`.
func (p *pmPejalan) daftar(jalur, pola string, xs []any) {
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
			p.medan = append(p.medan, MedanDokumen{Jalur: j, Pola: pola + "()", Nilai: s})
		}
	}
}

// baris - anggota satu unsur PageList. Objek di dalam baris diratakan (`anggota.sub`).
func (p *pmPejalan) baris(b Baris, daftar string, i int, awalan, jalur, pola string, m map[string]any) {
	for _, n := range pmKunciUrut(m) {
		nama := awalan + n
		switch v := m[n].(type) {
		case map[string]any:
			p.baris(b, daftar, i, nama+".", jalur, pola, v)
		case []any:
			p.daftar(jalur+"."+nama, pola+"."+nama, v)
		default:
			s, _ := TeksSkalarJSON(v)
			b[nama] = s
			p.medan = append(p.medan, MedanDokumen{Jalur: jalur + "." + nama, Pola: pola + "." + nama, Nilai: s,
				daftar: daftar, baris: i, nama: nama})
		}
	}
}

// setel menulis ulang nilai satu medan di halaman (hasil konversi tanggal).
func (p *pmPejalan) setel(m MedanDokumen, nilai string) {
	if m.daftar == "" {
		p.h.Setel(m.Jalur, nilai)
		return
	}
	p.h.AmbilDaftar(m.daftar)[m.baris][m.nama] = nilai
}

// pmUrai mengurai DATA_JSON menjadi objek akar berkelas PolicyTreatyIn.
func pmUrai(data []byte) (map[string]any, error) {
	dek := json.NewDecoder(bytes.NewReader(data))
	dek.UseNumber()
	var akar any
	if err := dek.Decode(&akar); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDokumenRusak, err)
	}
	m, ok := akar.(map[string]any)
	if !ok {
		return nil, ErrDokumenRusak
	}
	if kelas, _ := m["pxObjClass"].(string); kelas != KelasDokumenTreatyIn {
		return nil, fmt.Errorf("%w: pxObjClass %q", ErrBukanTreatyIn, kelas)
	}
	return m, nil
}

// pmProdKe - PRODKE bilangan >= 0.
func pmProdKe(prodke string) (int, error) {
	s := strings.TrimSpace(prodke)
	n, err := strconv.Atoi(s)
	if s == "" || err != nil || n < 0 {
		return 0, fmt.Errorf("%w: %q", ErrProdKe, prodke)
	}
	return n, nil
}

// pmProporsional - jenis proporsi halaman = Proportional (QuotationData, lalu Quotation; `PeriksaBentukSimpan`).
func pmProporsional(h *Halaman) bool {
	jenis := h.Ambil(HalamanPolis + ".QuotationData.ProportionalType")
	if jenis == "" {
		jenis = h.Ambil(HalamanQuotation + ".ProportionalType")
	}
	return jenis == JenisProporsional
}

// pmBuangOldData menghapus `PolicyTreatyIn.OldData.*` dari halaman: OldData bukan isi generasi ini (ID-8).
func pmBuangOldData(h *Halaman) {
	awalan := HalamanPolis + ".OldData."
	for j := range h.Nilai {
		if strings.HasPrefix(j, awalan) {
			delete(h.Nilai, j)
		}
	}
	for j := range h.Daftar {
		if strings.HasPrefix(j, awalan) {
			delete(h.Daftar, j)
		}
	}
}

// PecahDokumenEDM memecah satu baris JSON_POLIS generasi endorsemen menjadi halaman generasi yang siap disimpan
// lewat antarmuka penyimpanan yang SAMA dengan jalur biasa (`SimpanHalaman`, `SimpanSelisih`).
//
// Galat yang dikembalikan menghentikan dokumen seluruhnya (dokumen rusak, PRODKE / IDPEGA tak terbaca) atau
// menandai dokumen di luar lingkup (ErrBukanTreatyIn, ErrBukanGenerasiEndorsemen). Galat per medan terkumpul di
// `HasilPecahEDM.Galat`; dokumen bergalat TIDAK dimuat sebagian.
//
// ⛔ Nilai lama tidak dihitung ulang dan tidak dibulatkan (ID-33, AC 39) - termasuk `TreatyDifference.*` yang
// lahir dari varian rumus berlapis Pega: teks desimal dibawa apa adanya (pembulatan hanya di Oracle pada desimal
// kesebelas, NUMBER(38,10)). Tidak ada pengisian nilai: EndDate / ProductionDate kosong tetap kosong.
func PecahDokumenEDM(b BarisJSONPolis) (HasilPecahEDM, error) {
	if BarisAplikasiBaru(b) {
		return HasilPecahEDM{}, ErrBarisAplikasiBaru
	}
	m, err := pmUrai(b.DataJSON)
	if err != nil {
		return HasilPecahEDM{}, err
	}
	prodke, err := pmProdKe(b.ProdKe)
	if err != nil {
		return HasilPecahEDM{}, err
	}
	if prodke == 0 {
		return HasilPecahEDM{}, ErrBukanGenerasiEndorsemen
	}
	id, err := IDKasusDariIDPegaEDM(b.IDPega)
	if err != nil {
		return HasilPecahEDM{}, err
	}
	p := &pmPejalan{h: HalamanBaru()}
	p.objek(HalamanPolis, HalamanPolis, m)
	h := p.h
	hasil := HasilPecahEDM{
		ID: id, NoPolis: strings.TrimSpace(b.NoPolis), ProdKe: prodke, Halaman: h, Diabaikan: map[string]int{},
		Datar:        KolomDatarLama{TglInput: b.TglInput, Username: b.Username},
		EDMType:      h.Ambil(HalamanPolis + ".EDMType"),
		OldDataEDMNo: h.Ambil(pmJalurOldEDMNo),
	}
	galat := func(md MedanDokumen, err error) {
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: md.Jalur, Nilai: md.Nilai, Err: err})
	}
	arsip := func(md MedanDokumen, kunci string) {
		hasil.Arsip = append(hasil.Arsip, MedanArsip{MedanDokumen: md, Kunci: kunci})
		if kunci == "" {
			hasil.BelumDiputuskan = append(hasil.BelumDiputuskan, md)
			return
		}
		hasil.Diabaikan[pmTeksAlasan(kunci)]++
	}
	prop := pmProporsional(h)
	var edmDok string
	for _, md := range p.medan {
		if prop && pmHasilAntaraProp(md.Pola) {
			arsip(md, pmKunciPropHasilAntara)
			continue
		}
		if k, dikenal := pmPetaKatalog[md.Pola]; dikenal {
			if !k.Golongan.Tanggal() {
				if err := pmPeriksaNilai(k, md.Nilai); err != nil {
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
		switch md.Pola {
		case pmJalurNomorPolis:
			if md.Nilai != "" && strings.TrimSpace(md.Nilai) != hasil.NoPolis {
				galat(md, ErrNoPolisBeda)
			}
			continue
		case pmJalurEDMNo:
			edmDok = strings.TrimSpace(md.Nilai)
			continue
		case pmJalurProdKe:
			if s := strings.TrimSpace(md.Nilai); s != "" {
				if n, err := strconv.Atoi(s); err != nil || n != prodke {
					galat(md, fmt.Errorf("%w: dokumen %q, PRODKE %d", ErrProdKeBeda, md.Nilai, prodke))
				}
			}
			continue
		}
		if nama, ok := pmAnggotaUsulanLama(md.Pola); ok {
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
		arsip(md, pmKunciDibuang(md.Pola, md.Nilai))
	}
	if hasil.NoPolis == "" {
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: "NOPOLIS", Err: ErrNoPolisKosong})
	}
	switch kolom := strings.TrimSpace(b.NoEndors); {
	case kolom != "" && edmDok != "" && kolom != edmDok:
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: pmJalurEDMNo, Nilai: edmDok,
			Err: fmt.Errorf("%w: NOENDORS %q", ErrEDMNoBeda, kolom)})
	case kolom != "":
		hasil.EDMNo = kolom
	case edmDok != "":
		hasil.EDMNo = edmDok
	default:
		hasil.Galat = append(hasil.Galat, GalatDokumen{Jalur: "NOENDORS", Err: ErrEDMNoKosong})
	}
	pmBuangOldData(h)
	if prop {
		RapikanBentukSimpan(h)
	}
	// Kunci generasi di halaman sama dengan yang dibaca `repository.BacaGenerasi` (NOENDORS, PRODKE) - bukan kolom
	// katalog; pembanding generasi berikutnya saat uji-kering.
	if hasil.EDMNo != "" {
		h.Setel(pmJalurEDMNo, hasil.EDMNo)
	}
	h.Setel(pmJalurProdKe, strconv.Itoa(prodke))
	hasil.Usulan = UsulanDokumenLamaEDM(PyIDKasus(id), h)
	return hasil, nil
}

// HalamanPembanding menyusun halaman GENERASI SEBELUMNYA dari dokumen json_polis-nya (generasi NB ber-PRODKE 0
// maupun endorsemen) untuk uji-kering - padanan `BacaGenerasi(OLD_POLIS_ID)` tanpa basis data: hanya daftar dan
// nilai halaman (tanpa konversi, tanpa arsip), `EDMNo` = NOENDORS kolom datar (NB: NULL di SavePolisTreatyIn_SQL).
// Generasi Proportional dirapikan sama dengan jalur simpan (`RapikanBentukSimpan`).
func HalamanPembanding(b BarisJSONPolis) (*Halaman, error) {
	m, err := pmUrai(b.DataJSON)
	if err != nil {
		return nil, err
	}
	p := &pmPejalan{h: HalamanBaru()}
	p.objek(HalamanPolis, HalamanPolis, m)
	h := p.h
	pmBuangOldData(h)
	if pmProporsional(h) {
		RapikanBentukSimpan(h)
	}
	h.Setel(pmJalurEDMNo, strings.TrimSpace(b.NoEndors))
	return h, nil
}
