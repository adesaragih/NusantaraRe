package services

// Jalur buat dan baca kontrak - tiket 14.
//
// ⭐ DI SINILAH INV-53 DAN INV-29 DITEGAKKAN. Migrasi 401 menuliskannya
// sebagai pernyataan keputusan dengan baris "DITAGIH: tiket lapisan aplikasi,
// yang menulis jalur simpannya" - dan ini jalur simpan itu. ADR-0056 (K-4)
// menahan aturan bisnis di lapisan ini, bukan di basis data.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// PolaTanggal - bentuk tanggal di seluruh permukaan modul ini.
const PolaTanggal = "2006-01-02"

var (
	// ErrMasukanTidakSah - JSON-nya sah, isinya ditolak gerbang.
	ErrMasukanTidakSah = errors.New("masukan tidak sah")
	// ErrKontrakTidakAda - pengenal yang diminta tidak menunjuk kontrak mana pun.
	ErrKontrakTidakAda = errors.New("kontrak tidak ada")
	// ErrNomorUrutVersiGanda - INV-04 ditolak basis data; pesannya menyebut nomornya.
	ErrNomorUrutVersiGanda = errors.New("nomor urut versi sudah dipakai")
)

// MasukanKontrak adalah kontrak baru beserta KEPALA versi pertamanya.
//
// ⛔ Angka uang dan persen TEKS sepanjang jalan - tidak pernah float
// (ADR-0003). Kosong berarti kosong, bukan nol.
type MasukanKontrak struct {
	NomorKontrakWarisan string `json:"nomorKontrakWarisan"`
	IDCedant            int64  `json:"idCedant"`
	IDAsalBisnis        int64  `json:"idAsalBisnis"`
	SifatProporsi       string `json:"sifatProporsi"`
	TanggalMulai        string `json:"tanggalMulai"`
	TanggalBerakhir     string `json:"tanggalBerakhir"`

	NamaKontrak         string `json:"namaKontrak"`
	KodeMataUangKontrak int64  `json:"kodeMataUangKontrak"`
	PersenBagianNure    string `json:"persenBagianNure"`
	BagianNureSeragam   string `json:"bagianNureSeragam"`
	MemakaiBordereaux   string `json:"memakaiBordereaux"`
	CaraPembukuan       string `json:"caraPembukuan"`
	MemakaiProrata      string `json:"memakaiProrata"`
	RetroBerganda       string `json:"retroBerganda"`
	KeadaanSiklusHidup  string `json:"keadaanSiklusHidup"`
}

// HasilBuatKontrak adalah pengenal yang SISTEM berikan - tiket 14 menuntut
// kontraknya dapat ditemukan kembali dengannya.
type HasilBuatKontrak struct {
	IDKontrak int64 `json:"idKontrak"`
	IDVersi   int64 `json:"idVersi"`
	// Peringatan tiket 16 - penyimpanan BERHASIL walau ia terisi.
	Peringatan []Peringatan `json:"peringatan,omitempty"`
}

// periksaMasukan menjalankan seluruh gerbang, dan mengumpulkan SELURUH
// pelanggaran sebelum menolak.
//
// ⛔ Mengumpulkan, bukan berhenti di yang pertama: pengisi yang memperbaiki
// satu ruas lalu ditolak lagi oleh ruas berikutnya akan menebak berapa kali
// lagi ia harus mencoba.
func periksaMasukan(m MasukanKontrak) (models.Kontrak, models.VersiKontrak, error) {
	var salah []string
	wajibTeks := func(nama, nilai string) string {
		nilai = strings.TrimSpace(nilai)
		if nilai == "" {
			salah = append(salah, nama+" wajib diisi")
		}
		return nilai
	}
	wajibAngka := func(nama string, nilai int64) int64 {
		if nilai <= 0 {
			salah = append(salah, nama+" wajib diisi dengan pengenal yang sah")
		}
		return nilai
	}

	// INV-29 - dua nilai, disebut satu per satu.
	sifat := strings.TrimSpace(m.SifatProporsi)
	if sifat != models.SifatProporsional && sifat != models.SifatNonProporsional {
		salah = append(salah, fmt.Sprintf("sifatProporsi %q bukan %s maupun %s",
			m.SifatProporsi, models.SifatProporsional, models.SifatNonProporsional))
	}

	// INV-53 - batas INKLUSIF keduanya (ADR-0022).
	mulai, errMulai := time.Parse(PolaTanggal, strings.TrimSpace(m.TanggalMulai))
	if errMulai != nil {
		salah = append(salah, "tanggalMulai bukan tanggal "+PolaTanggal)
	}
	akhir, errAkhir := time.Parse(PolaTanggal, strings.TrimSpace(m.TanggalBerakhir))
	if errAkhir != nil {
		salah = append(salah, "tanggalBerakhir bukan tanggal "+PolaTanggal)
	}
	if errMulai == nil && errAkhir == nil && akhir.Before(mulai) {
		salah = append(salah, fmt.Sprintf(
			"tanggalBerakhir %s lebih awal daripada tanggalMulai %s (INV-53, batas inklusif)",
			akhir.Format(PolaTanggal), mulai.Format(PolaTanggal)))
	}

	nama := wajibTeks("namaKontrak", m.NamaKontrak)
	keadaan := wajibTeks("keadaanSiklusHidup", m.KeadaanSiklusHidup)
	seragam := wajibTeks("bagianNureSeragam", m.BagianNureSeragam)
	bordereaux := wajibTeks("memakaiBordereaux", m.MemakaiBordereaux)
	pembukuan := wajibTeks("caraPembukuan", m.CaraPembukuan)
	prorata := wajibTeks("memakaiProrata", m.MemakaiProrata)
	retro := wajibTeks("retroBerganda", m.RetroBerganda)
	cedant := wajibAngka("idCedant", m.IDCedant)
	asal := wajibAngka("idAsalBisnis", m.IDAsalBisnis)
	mataUang := wajibAngka("kodeMataUangKontrak", m.KodeMataUangKontrak)

	persen, _, errPersen := apd.NewFromString(strings.TrimSpace(m.PersenBagianNure))
	if errPersen != nil || strings.TrimSpace(m.PersenBagianNure) == "" {
		salah = append(salah, "persenBagianNure wajib diisi angka desimal")
	}

	if len(salah) > 0 {
		return models.Kontrak{}, models.VersiKontrak{}, fmt.Errorf("%w: %s",
			ErrMasukanTidakSah, strings.Join(salah, "; "))
	}

	satu := int64(1)
	return models.Kontrak{
			NomorKontrakWarisan: strings.TrimSpace(m.NomorKontrakWarisan),
			IDCedant:            cedant,
			IDAsalBisnis:        asal,
			SifatProporsi:       sifat,
			TanggalMulai:        mulai.Format(PolaTanggal),
			TanggalBerakhir:     akhir.Format(PolaTanggal),
		}, models.VersiKontrak{
			// Versi PERTAMA selalu bernomor 1. Penomoran versi berikutnya
			// milik tiket 01 papan Adjustment, yang menyimpan rujukan dasarnya.
			NomorUrutVersi:     &satu,
			KeadaanSiklusHidup: keadaan,
			NamaKontrak:        nama,
			KodeMataUangKontak: mataUang,
			PersenBagianNure:   persen,
			BagianNureSeragam:  seragam,
			MemakaiBordereaux:  bordereaux,
			CaraPembukuan:      pembukuan,
			MemakaiProrata:     prorata,
			RetroBerganda:      retro,
		}, nil
}

// BuatKontrak membuat kontrak beserta versi pertamanya, dan mengembalikan
// pengenal yang SISTEM berikan (INV-02: dari sequence, tidak pernah dari teks).
func (l *Layanan) BuatKontrak(ctx context.Context, p inti.Pelaku, m MasukanKontrak) (HasilBuatKontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilBuatKontrak{}, err
	}
	k, v, err := periksaMasukan(m)
	if err != nil {
		return HasilBuatKontrak{}, err
	}
	// Tiket 16: kunci alami kembar MEMPERINGATKAN, tidak menghalangi. Dicari
	// SEBELUM menulis supaya pembandingnya tidak memuat kontrak ini sendiri.
	serupa, err := l.gudang.CariKontrakSerupa(ctx, k)
	if err != nil {
		return HasilBuatKontrak{}, err
	}
	idK, idV, err := l.gudang.BuatKontrakDenganVersiPertama(ctx, k, v)
	if err != nil {
		return HasilBuatKontrak{}, err
	}
	hasil := HasilBuatKontrak{IDKontrak: idK, IDVersi: idV}
	if p := peringatkanKunciAlamiGanda(serupa); p != nil {
		hasil.Peringatan = append(hasil.Peringatan, *p)
	}
	return hasil, nil
}

// BacaKontrak menemukan kembali kontrak dengan pengenalnya, beserta seluruh
// versinya - lapisan beku SEKALI, tidak disalin ke tiap versi (ADR-0040).
func (l *Layanan) BacaKontrak(ctx context.Context, p inti.Pelaku, id int64) (models.KontrakDenganVersi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.KontrakDenganVersi{}, err
	}
	if id <= 0 {
		return models.KontrakDenganVersi{}, fmt.Errorf("%w: pengenal kontrak %d", ErrMasukanTidakSah, id)
	}
	return l.gudang.BacaKontrak(ctx, id)
}
