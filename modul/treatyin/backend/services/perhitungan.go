package services

// Validasi berskema dan sakelar pemindahan — tiket 21, 35, 36, 43.
//
//	21  perhitungan yang gagal menghasilkan KETERANGAN, bukan nol
//	35  tepat satu dari persen quota share atau jumlah lines surplus
//	36  baris surplus tanpa baris quota share pada versi yang sama
//	43  sakelar penegakan selama pemindahan, dan keadaannya terlihat

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
)

// Dua nilai sah JENIS_TREATY (INV-30) — himpunan TERTUTUP.
const (
	JenisQuotaShare = "QUOTA_SHARE"
	JenisSurplus    = "SURPLUS"
)

var (
	// ErrHitungGagal — tiket 21. Perhitungan yang tidak dapat diselesaikan.
	ErrHitungGagal = errors.New("perhitungan tidak dapat diselesaikan")
	// ErrKetentuanProporsional — tiket 35.
	ErrKetentuanProporsional = errors.New("ketentuan proporsional tidak sah")
	// ErrSurplusTanpaQuotaShare — tiket 36.
	ErrSurplusTanpaQuotaShare = errors.New("baris surplus menuntut baris quota share pada versi yang sama")
)

// KonversiMataUang mengubah nilai memakai kurs, dan MENOLAK ketika tidak dapat.
//
// ⛔ Yang dilarang di sini bukan kegagalannya, melainkan MENYAMARKANNYA.
// `ADR-0035`: sistem lama menjawab kurs yang tidak ditemukan dengan **satu**,
// dan pembagi nol dengan **nol** — keduanya angka yang terlihat sah di laporan,
// dan tidak ada yang tahu ia karangan. Yang dikembalikan di sini keterangan
// APA yang gagal, bukan angka.
//
// Nilai nol yang SAH tetap nol: yang ditolak nol yang lahir dari kegagalan.
func (l *Layanan) KonversiMataUang(_ context.Context, p inti.Pelaku, nilai, kurs *apd.Decimal, kodeMataUang string) (*apd.Decimal, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	kode := strings.TrimSpace(kodeMataUang)
	switch {
	case kode == "":
		return nil, fmt.Errorf("%w: kode mata uang kosong, tidak ada kurs yang dapat dicari", ErrHitungGagal)
	case nilai == nil:
		return nil, fmt.Errorf("%w: nilai sumber kosong untuk mata uang %s", ErrHitungGagal, kode)
	case kurs == nil:
		// ⛔ TIDAK memakai satu sebagai pengganti. Itu persis yang ADR-0035
		// larang: hasilnya terlihat benar dan tidak ada yang tahu kursnya
		// tidak pernah ditemukan.
		return nil, fmt.Errorf("%w: kurs untuk mata uang %s tidak ditemukan", ErrHitungGagal, kode)
	case kurs.IsZero():
		return nil, fmt.Errorf("%w: kurs untuk mata uang %s bernilai nol", ErrHitungGagal, kode)
	}
	hasil := new(apd.Decimal)
	if _, err := apd.BaseContext.Mul(hasil, nilai, kurs); err != nil {
		return nil, fmt.Errorf("%w: mengalikan nilai dengan kurs %s: %v", ErrHitungGagal, kode, err)
	}
	return hasil, nil
}

// MasukanKetentuanProporsional adalah satu baris `DETAIL_PROPORSIONAL`.
type MasukanKetentuanProporsional struct {
	IDLayer            int64  `json:"idLayer"`
	IDKelompokTreaty   int64  `json:"idKelompokTreaty"`
	JenisTreaty        string `json:"jenisTreaty"`
	PersenQuotaShare   string `json:"persenQuotaShare"`
	JumlahLinesSurplus *int64 `json:"jumlahLinesSurplus"`
}

// periksaKetentuanProporsional menegakkan tiket 35.
//
// ⛔ Kedua kolom MENGUKUR HAL BERBEDA DENGAN SATUAN BERBEDA: persen bagian
// versus kelipatan retensi. Baris yang mengisi keduanya tidak punya arti
// tunggal; baris yang mengosongkan keduanya tidak punya arti sama sekali.
func periksaKetentuanProporsional(m MasukanKetentuanProporsional) error {
	jenis := strings.TrimSpace(m.JenisTreaty)
	adaQS := strings.TrimSpace(m.PersenQuotaShare) != ""
	adaSurplus := m.JumlahLinesSurplus != nil

	var salah []string
	if jenis != JenisQuotaShare && jenis != JenisSurplus {
		salah = append(salah, fmt.Sprintf("jenisTreaty %q bukan %s maupun %s (INV-30)",
			m.JenisTreaty, JenisQuotaShare, JenisSurplus))
	}
	switch {
	case adaQS && adaSurplus:
		salah = append(salah, "persenQuotaShare dan jumlahLinesSurplus terisi KEDUANYA; tepat satu yang boleh")
	case !adaQS && !adaSurplus:
		salah = append(salah, "persenQuotaShare dan jumlahLinesSurplus KOSONG keduanya; tepat satu wajib terisi")
	case jenis == JenisQuotaShare && !adaQS:
		salah = append(salah, "jenisTreaty QUOTA_SHARE menuntut persenQuotaShare, yang terisi justru jumlahLinesSurplus")
	case jenis == JenisSurplus && !adaSurplus:
		salah = append(salah, "jenisTreaty SURPLUS menuntut jumlahLinesSurplus, yang terisi justru persenQuotaShare")
	}
	if adaSurplus && *m.JumlahLinesSurplus <= 0 {
		salah = append(salah, "jumlahLinesSurplus wajib bilangan positif")
	}
	if len(salah) > 0 {
		return fmt.Errorf("%w: %s", ErrKetentuanProporsional, strings.Join(salah, "; "))
	}
	return nil
}

// CatatKetentuanProporsional menegakkan tiket 35 dan 36.
//
// Tiket 36: baris `SURPLUS` menuntut adanya baris `QUOTA_SHARE` pada versi yang
// sama. Kapasitas surplus dihitung retensi × jumlah lines, dan retensinya datang
// dari baris quota share — tanpa baris itu hasilnya bukan nol melainkan **tidak
// terdefinisi**, dan menyajikannya sebagai nol adalah ADR-0035 sekali lagi.
func (l *Layanan) CatatKetentuanProporsional(ctx context.Context, p inti.Pelaku, idVersi int64, m MasukanKetentuanProporsional) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	if idVersi <= 0 || m.IDLayer <= 0 || m.IDKelompokTreaty <= 0 {
		return fmt.Errorf("%w: pengenal versi, layer, dan kelompok treaty wajib positif", ErrMasukanTidakSah)
	}
	if err := periksaKetentuanProporsional(m); err != nil {
		return err
	}
	if strings.TrimSpace(m.JenisTreaty) == JenisSurplus {
		ada, err := l.gudang.AdaQuotaSharePadaVersi(ctx, idVersi)
		if err != nil {
			return err
		}
		if !ada {
			// Pesannya menyebut APA yang kurang, bukan sekadar "ditolak".
			return fmt.Errorf("%w: versi %d belum punya baris QUOTA_SHARE, sehingga retensi "+
				"yang menjadi pengali kapasitas surplus tidak ada", ErrSurplusTanpaQuotaShare, idVersi)
		}
	}
	return l.gudang.CatatKetentuanProporsional(ctx, idVersi, m.IDLayer, m.IDKelompokTreaty,
		strings.TrimSpace(m.JenisTreaty), strings.TrimSpace(m.PersenQuotaShare), m.JumlahLinesSurplus)
}

// SakelarPemindahan mematikan penegakan selama pemindahan data warisan — tiket 43.
//
// ⛔ ADR-0042 memerintahkan sejarah pindah APA ADANYA, dan sebagian data lama
// melanggar invarian yang kita tegakkan sekarang. Tanpa sakelar, pemindahan
// berhenti di baris pertama yang menyimpang; dengan sakelar yang tidak terlihat,
// ia bisa tertinggal mati tanpa ada yang tahu.
//
// Dua tuntutan tiket 43 yang mudah terlewat:
//   - keadaannya terlihat TANPA membuka basis data
//   - menyalakan kembali MEMERIKSA ULANG baris yang masuk selagi ia mati
type SakelarPemindahan struct {
	mu     sync.RWMutex
	mati   bool
	alasan string
	// masuk mencatat baris yang diterima selagi penegakan mati — tanpa daftar
	// ini, "memeriksa ulang" tidak punya yang diperiksa.
	masuk []int64
}

// KeadaanSakelar adalah potret yang dapat dibaca tanpa menyentuh Oracle.
type KeadaanSakelar struct {
	PenegakanMati bool   `json:"penegakanMati"`
	Alasan        string `json:"alasan,omitempty"`
	BarisMasuk    int    `json:"barisMasuk"`
}

// Matikan mematikan penegakan; alasan WAJIB.
//
// Sakelar yang dapat dimatikan tanpa alasan akan ditemukan mati tanpa ada yang
// ingat kenapa.
func (s *SakelarPemindahan) Matikan(alasan string) error {
	if strings.TrimSpace(alasan) == "" {
		return fmt.Errorf("%w: alasan mematikan penegakan wajib diisi", ErrMasukanTidakSah)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mati, s.alasan = true, strings.TrimSpace(alasan)
	return nil
}

// Keadaan mengembalikan potretnya — terlihat tanpa membuka basis data.
func (s *SakelarPemindahan) Keadaan() KeadaanSakelar {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return KeadaanSakelar{PenegakanMati: s.mati, Alasan: s.alasan, BarisMasuk: len(s.masuk)}
}

// Catat menandai satu baris masuk selagi penegakan mati.
func (s *SakelarPemindahan) Catat(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mati {
		s.masuk = append(s.masuk, id)
	}
}

// Nyalakan menyalakan penegakan kembali DAN memeriksa ulang baris yang masuk
// selagi ia mati; pengenal yang gagal dikembalikan.
//
// ⛔ Menyalakan tanpa memeriksa ulang membuat sakelar ini lebih buruk daripada
// tidak ada: ia memberi keyakinan bahwa penegakan berlaku, atas data yang tidak
// pernah dinilai.
func (s *SakelarPemindahan) Nyalakan(periksa func(id int64) error) ([]int64, error) {
	if periksa == nil {
		return nil, fmt.Errorf("%w: menyalakan penegakan menuntut pemeriksa ulang", ErrMasukanTidakSah)
	}
	s.mu.Lock()
	masuk := s.masuk
	s.masuk = nil
	s.mati, s.alasan = false, ""
	s.mu.Unlock()

	var gagal []int64
	for _, id := range masuk {
		if periksa(id) != nil {
			gagal = append(gagal, id)
		}
	}
	return gagal, nil
}
