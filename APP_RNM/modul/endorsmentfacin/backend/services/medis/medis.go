// Package medis memuat skoring pemeriksaan medis tertanggung jiwa - SEAM 6
// (K-052), tiket E20.
//
// Untuk apa berkas ini: hasil laboratorium dan usia tertanggung diubah menjadi
// skor per pemeriksaan, persis seperti `Endorsment Fac In/Activity/
// CalculateScorLife_Act.xml` (`DATA-PARTY-PERSON!CALCULATESCORLIFE_ACT`).
// Dua pintu pendamping memakai mesin yang sama: `NilaiRujukanLab`
// (`SetParamLab_Act`) dan `PeriksaFisik` (`CalculatePhysicalExam`).
//
// Dibaca sesudah: ekspresi.go, aturan_gen.go.
//
// ⛔ BUKAN perhitungan premi: `[terverifikasi]` activity ini tidak menyentuh
// rate, TSI, premi, maupun pembagian. Keluarannya TEKS skor per pemeriksaan
// (mis. "Standard", "Decline", "Postpone for 12 Months", "Invalid Age",
// "150%", "175") - diport apa adanya, TIDAK dipetakan ke kategori lain: arti
// tiap teks, termasuk persen dan angka, adalah `[pertanyaan terbuka]` milik
// Underwriting. Penggabungan menjadi satu keputusan akhir TIDAK ada di activity
// ini.
//
// ⛔ Domain sensitif: hasil lab, usia, dan skor TIDAK PERNAH ditulis ke log.
// Paket ini tidak mencatat apa pun.
//
// Skoring RISIKO yang bersama NB (`ScoringRisk` dkk) adalah sistem lain dan
// bukan isi paket ini (spec Life §1.3).
package medis

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Kode transisi prakondisi langkah activity.
//
// [terverifikasi] Hanya tiga pola di activity ini: satu baris (2,3) ×30 · dua
// baris (2,3)(2,3) ×10 · dua baris (5,2)(2,3) ×170. Dua cara, sepakat:
//
//	urai pohon: TestTransisiKorpusTigaPola (medis_test.go) atas aturanSkor
//	grep -c '<pyStepsPreCondParamsWhenTrue>5<'  CalculateScorLife_Act.xml   # 170
//	grep -c '<pyStepsPreCondParamsWhenTrue>2<'  CalculateScorLife_Act.xml   # 220 = 30 + 10×2 + 170
//	grep -c '<pyStepsPreCondParamsWhenFalse>3<' CalculateScorLife_Act.xml   # 220
//	grep -c '<pyStepsPreCondParamsWhenFalse>2<' CalculateScorLife_Act.xml   # 170
//
// ⚠️ Kesepuluh pola (2,3)(2,3) ada di langkah berlabel `//` (di-remark) - tidak
// pernah dijalankan.
//
// `pyStepsRepeatDefHasRepeat=REPEAT` di 50 langkah blok `[terverifikasi]`
// bernilai Start=1 · Limit=1 · Iteration=1 di seluruh 56 definisinya - satu kali
// jalan; blok dijalankan sekali:
//
//		grep -c '<pyStepsRepeatDefHasRepeat>REPEAT<' CalculateScorLife_Act.xml   # 50
//		grep -c '<pyStepsRepeatDefLimit>1<'          CalculateScorLife_Act.xml   # 56
//		grep -o '<pyStepsRepeatDefLimit>[^<]*'       CalculateScorLife_Act.xml | sort -u   # hanya 1
//
//	  - "2" lanjut ke baris berikutnya; sesudah baris terakhir, langkah dijalankan;
//	  - "3" langkah dilewati;
//	  - "6" keluar activity (tidak muncul di sini; arti dari `SetOldData`);
//	  - "5" `[dugaan]` langkah dijalankan TANPA memeriksa baris berikutnya.
//	    Tidak satu pun dari 170 kemunculannya membawa `WhenTruePrms` (jadi
//	    bukan "lompat ke langkah"), dan hanya bacaan ini yang membuat cabang
//	    `param.JN=="All"` bermakna: pola (5,2)(2,3) = "A ATAU B".
const (
	transisiLanjut    = "2"
	transisiLewati    = "3"
	transisiJalankan  = "5"
	transisiKeluar    = "6"
	transisiTakTerisi = ""
)

// aksiLangkah - keputusan mesin atas satu langkah. Sengaja tipe tersendiri,
// terpisah dari kode transisi korpus: "jalan tanpa syarat" (prakondisi
// nonaktif) tidak sama dengan kode "5" yang artinya masih `[dugaan]`.
type aksiLangkah uint8

const (
	aksiJalankan aksiLangkah = iota + 1
	aksiLewati
	aksiKeluar
)

// ErrTransisiTakDikenal - kode transisi di luar yang dimengerti.
var ErrTransisiTakDikenal = errors.New("medis: kode transisi prakondisi tidak dikenal")

type barisPrakondisi struct {
	ekspresi             string
	bilaBenar, bilaSalah string
}

type penugasan struct {
	nama, ekspresi string
}

type langkah struct {
	nomor, metode, deskripsi string
	berprakondisi            bool
	// diRemark - label `//` (`pyStepsBlockName`): langkah beserta seluruh
	// sub-langkahnya TIDAK dijalankan (keputusan work owner 01-10-2026 butir 43,
	// nbfacin `docs/KEPUTUSAN-30-09-2026.md`).
	diRemark bool
	// halamanDihapus - halaman sasaran langkah `Page-Remove`.
	halamanDihapus string
	prakondisi     []barisPrakondisi
	penugasan      []penugasan
	anak           []langkah
}

// aktivitas - satu activity Pega yang dibangkitkan.
type aktivitas struct {
	asal                                      string
	cacahLangkah, cacahPenugasan, cacahRemark int
	langkah                                   []langkah
}

// Masukan - satu tertanggung dan satu permintaan skoring.
type Masukan struct {
	// JenisKelamin - `.Medical.Sex` apa adanya (korpus menguji `==1` dan `==2`).
	JenisKelamin string
	// Lab - `.Medical.Lab` apa adanya (korpus: "Prodia", "Pramita", "Biotest").
	Lab string
	// Usia - `.Age` apa adanya.
	Usia string
	// JN, JNSub - parameter `param.JN` / `param.JN_Sub`: kelompok dan jenis
	// pemeriksaan yang diminta, atau "All".
	JN, JNSub string
	// Pemeriksaan - `.Medical.Medical_Examination.<nama>` → nilai apa adanya
	// (boleh berkoma desimal; activity sendiri mengganti "," dengan ".").
	Pemeriksaan map[string]string
	// ParamLab - `ParamLab.Medical_Examination.<nama>`, halaman lain yang
	// dibaca dua prakondisi (keberadaan nilai AFP).
	ParamLab map[string]string
}

// Galat - satu langkah atau satu penugasan yang tidak dapat dijalankan tanpa
// menebak. Pemeriksaan lain tetap dinilai.
type Galat struct {
	// Langkah - nomor langkah activity, mis. "1.1.2.1".
	Langkah string
	// Properti - properti yang gagal diisi; kosong bila prakondisinya yang gagal.
	Properti string
	Err      error
}

func (g Galat) Error() string {
	if g.Properti == "" {
		return fmt.Sprintf("langkah %s prakondisi: %v", g.Langkah, g.Err)
	}
	return fmt.Sprintf("langkah %s %s: %v", g.Langkah, g.Properti, g.Err)
}

func (g Galat) Unwrap() error { return g.Err }

// Hasil - keluaran skoring.
type Hasil struct {
	// Skor - `.Medical.Medical_Examination.Score.<nama>` → teks skor apa adanya.
	Skor map[string]string
	// Pemeriksaan - nilai pemeriksaan sesudah dinormalkan activity
	// (`@toDecimal(@replaceAll(…,",","."))`).
	Pemeriksaan map[string]string
	// Galat - langkah/penugasan yang tidak dinilai, berurutan.
	Galat []Galat
}

const (
	awalanPemeriksaan = ".Medical.Medical_Examination."
	awalanSkor        = ".Medical.Medical_Examination.Score."
	awalanParamLab    = "ParamLab.Medical_Examination."
)

// halamanKerja - clipboard satu skoring: nilai teks per rujukan, ditambah
// penanda properti yang penugasannya gagal.
type halamanKerja struct {
	nilai map[string]string
	gagal map[string]error
	// angka - properti yang diisi hasil bertipe angka.
	angka map[string]bool
}

func (h *halamanKerja) bertipeAngka(nama string) bool { return h.angka[nama] }

func (h *halamanKerja) ambil(nama string) (string, error) {
	if err, ada := h.gagal[nama]; ada {
		return "", err
	}
	return h.nilai[nama], nil
}

// ScoreMedical - SEAM 6. Menjalankan pohon langkah `CalculateScorLife_Act`
// atas satu tertanggung.
//
// Galat yang dikembalikan hanya galat STRUKTUR (kode transisi tak dikenal);
// galat per pemeriksaan ada di `Hasil.Galat`, dan skor yang gagal tidak
// muncul di `Hasil.Skor` - tidak pernah diisi tebakan.
func ScoreMedical(m Masukan) (Hasil, error) {
	h := &halamanKerja{nilai: map[string]string{
		".Medical.Sex": m.JenisKelamin,
		".Medical.Lab": m.Lab,
		".Age":         m.Usia,
		"param.JN":     m.JN,
		"param.JN_Sub": m.JNSub,
	}, gagal: map[string]error{}}
	for k, v := range m.Pemeriksaan {
		h.nilai[awalanPemeriksaan+k] = v
	}
	for k, v := range m.ParamLab {
		h.nilai[awalanParamLab+k] = v
	}
	galat, err := jalankanAktivitas(aturanSkor, h)
	if err != nil {
		return Hasil{}, err
	}
	periksa := ambilAwalan(h, awalanPemeriksaan)
	for k := range periksa {
		// `Score.` tinggal di bawah halaman yang sama; ia keluaran, bukan nilai pemeriksaan.
		if strings.HasPrefix(k, "Score.") {
			delete(periksa, k)
		}
	}
	return Hasil{Skor: ambilAwalan(h, awalanSkor), Pemeriksaan: periksa, Galat: galat}, nil
}

// NilaiRujukanLab - `SetParamLab_Act`: nilai rujukan pemeriksaan per jenis
// kelamin × lab, sebagai teks apa adanya. Hasilnya menjadi `Masukan.ParamLab`.
//
// [terverifikasi] Langkah 1 `Page-Remove ParamLab` bila `.Medical.Lab==""`;
// langkah 2-3 per jenis kelamin, anak per lab (Prodia, Pramita, Biotest),
// 438 penugasan seluruhnya (`grep -c '<PropertiesName>' SetParamLab_Act.xml`).
func NilaiRujukanLab(jenisKelamin, lab string) (map[string]string, []Galat, error) {
	h := &halamanKerja{nilai: map[string]string{".Medical.Sex": jenisKelamin, ".Medical.Lab": lab}, gagal: map[string]error{}}
	galat, err := jalankanAktivitas(aturanParamLab, h)
	if err != nil {
		return nil, nil, err
	}
	return ambilAwalan(h, awalanParamLab), galat, nil
}

// MasukanFisik - pemeriksaan fisik satu tertanggung, nilai apa adanya.
type MasukanFisik struct {
	JenisKelamin string // .Medical.Sex
	// Pemeriksaan - `.Medical.Physcycal_Examination.<nama>` (ejaan korpus):
	// `Weight`, `Height`, `Girth_of_Abdomen`.
	Pemeriksaan map[string]string
}

// HasilFisik - `.Medical.Physcycal_Examination.*` sesudah activity, termasuk
// `BM_Ratio`, `BMI_Ratio_Note`, `EM`, `Girth_of_Abdomen_Note`,
// `Score.Girth_of_Abdomen`.
type HasilFisik struct {
	Nilai map[string]string
	Galat []Galat
}

const awalanFisik = ".Medical.Physcycal_Examination."

// PeriksaFisik - `CalculatePhysicalExam`: rasio BMI dan lingkar perut.
//
// ⚠️ `@divide(berat, tinggi²)` di korpus TANPA argumen skala; skala dan
// pembulatan bawaannya belum terverifikasi, sehingga BMI yang tak eksak
// menjadi `ErrSkalaBawaanBelumTerverifikasi` - catatan BMI dan `EM`
// sesudahnya ikut tidak dinilai.
func PeriksaFisik(m MasukanFisik) (HasilFisik, error) {
	h := &halamanKerja{nilai: map[string]string{".Medical.Sex": m.JenisKelamin}, gagal: map[string]error{}}
	for k, v := range m.Pemeriksaan {
		h.nilai[awalanFisik+k] = v
	}
	galat, err := jalankanAktivitas(aturanFisik, h)
	if err != nil {
		return HasilFisik{}, err
	}
	return HasilFisik{Nilai: ambilAwalan(h, awalanFisik), Galat: galat}, nil
}

// ambilAwalan - nilai yang sah (bukan gagal) berawalan tertentu, tanpa awalannya.
func ambilAwalan(h *halamanKerja, awalan string) map[string]string {
	hasil := map[string]string{}
	for k, v := range h.nilai {
		if _, gagal := h.gagal[k]; gagal {
			continue
		}
		if strings.HasPrefix(k, awalan) {
			hasil[strings.TrimPrefix(k, awalan)] = v
		}
	}
	return hasil
}

// jalankanAktivitas - satu activity utuh atas halaman kerja.
func jalankanAktivitas(a aktivitas, h *halamanKerja) ([]Galat, error) {
	var galat []Galat
	if _, err := jalankan(a.langkah, h, &galat); err != nil {
		return nil, fmt.Errorf("%s: %w", a.asal, err)
	}
	return galat, nil
}

// jalankan - langkah berurutan. `keluar` = transisi "6" (keluar activity).
func jalankan(daftar []langkah, h *halamanKerja, galat *[]Galat) (keluar bool, err error) {
	for _, l := range daftar {
		if l.diRemark {
			continue
		}
		aksi, err := putuskan(l, h)
		if err != nil {
			if errors.Is(err, ErrTransisiTakDikenal) {
				return false, fmt.Errorf("langkah %s: %w", l.nomor, err)
			}
			*galat = append(*galat, Galat{Langkah: l.nomor, Err: err})
			continue
		}
		switch aksi {
		case aksiKeluar:
			return true, nil
		case aksiLewati:
			continue
		}
		if l.metode == "Page-Remove" {
			hapusHalaman(h, l.halamanDihapus)
			continue
		}
		for _, p := range l.penugasan {
			if p.nama == "" && p.ekspresi == "" {
				continue // `SET  = ` kosong di korpus: tanpa efek
			}
			v, bertipeAngka, err := evaluasiNilai(p.ekspresi, h)
			if err != nil {
				h.gagal[p.nama] = err
				*galat = append(*galat, Galat{Langkah: l.nomor, Properti: p.nama, Err: err})
				continue
			}
			delete(h.gagal, p.nama)
			h.nilai[p.nama] = v
			if h.angka == nil {
				h.angka = map[string]bool{}
			}
			h.angka[p.nama] = bertipeAngka
		}
		keluar, err := jalankan(l.anak, h, galat)
		if err != nil || keluar {
			return keluar, err
		}
	}
	return false, nil
}

// hapusHalaman - `Page-Remove`: seluruh properti halaman itu hilang.
func hapusHalaman(h *halamanKerja, halaman string) {
	awalan := halaman + "."
	for k := range h.nilai {
		if strings.HasPrefix(k, awalan) {
			delete(h.nilai, k)
		}
	}
	for k := range h.gagal {
		if strings.HasPrefix(k, awalan) {
			delete(h.gagal, k)
		}
	}
}

// putuskan - jalankan, lewati, atau keluar activity.
//
// `pyStepsPreCondition=false` → tetap jalan tanpa syarat (P-11, tertutup).
func putuskan(l langkah, h *halamanKerja) (aksiLangkah, error) {
	if !l.berprakondisi || len(l.prakondisi) == 0 {
		return aksiJalankan, nil
	}
	for _, b := range l.prakondisi {
		benar, err := evaluasiBool(b.ekspresi, h)
		if err != nil {
			return 0, err
		}
		kode := b.bilaSalah
		if benar {
			kode = b.bilaBenar
		}
		switch kode {
		case transisiLanjut, transisiTakTerisi:
			continue
		case transisiLewati:
			return aksiLewati, nil
		case transisiJalankan:
			return aksiJalankan, nil
		case transisiKeluar:
			return aksiKeluar, nil
		default:
			return 0, fmt.Errorf("%w: %q", ErrTransisiTakDikenal, kode)
		}
	}
	return aksiJalankan, nil
}

// namaSkor - seluruh properti skor yang dapat diisi aturan, terurut. Dipakai
// test untuk menjaga bahwa tabel bangkitan memuat apa yang diharapkan.
func namaSkor() []string {
	set := map[string]bool{}
	var kumpul func([]langkah)
	kumpul = func(ls []langkah) {
		for _, l := range ls {
			for _, p := range l.penugasan {
				if strings.HasPrefix(p.nama, awalanSkor) {
					set[strings.TrimPrefix(p.nama, awalanSkor)] = true
				}
			}
			kumpul(l.anak)
		}
	}
	kumpul(aturanSkor.langkah)
	var h []string
	for k := range set {
		h = append(h, k)
	}
	sort.Strings(h)
	return h
}
