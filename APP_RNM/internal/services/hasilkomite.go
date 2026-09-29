package services

// Membaca hasil keputusan Komite, dan memulai putaran berikutnya - tiket 11.
//
// Untuk apa berkas ini: MEMBACA hasil yang sudah ditulis Komite, dan membuka
// putaran baru bila barisnya ditolak.
//
// ⛔ Berkas ini TIDAK menulis status baris. `[terverifikasi]` penulisnya
// `Komite Claim Life/Activity/KomitePostAdjustment.xml`, dan
// `[keputusan work owner 2026-09-15]` menetapkan kepemilikannya:
// **Komite menulis, Claim Life membaca.** Dua tiket yang sama-sama mengklaim
// penulisan yang sama berarti dua agent dapat mengerjakannya berdua.
//
// Dibaca sesudah: komite.go, adjustment.go.
//
// Istilah:
//   - putaran : satu siklus baris adjustment dari Outstanding sampai keputusan.
//   - tingkat : anak tangga komite; hasil hanya berlaku di tingkat TERAKHIR.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// ErrBukanPenolakan - putaran baru lahir dari penolakan, bukan dari apa pun.
var ErrBukanPenolakan = errors.New(
	"services: baris lanjutan hanya lahir sesudah baris terakhir ditolak")

// ErrTahapTanpaAddAdjustment - `Add` grid adjustment tampil hanya di Claim
// Analis: `pyWorkPage.pyPosition =='ReasLifeSPV'` (ClaimLifeDetailGCNM b18160).
var ErrTahapTanpaAddAdjustment = errors.New(
	"services: Add baris adjustment hanya tersedia di tahap Claim Analis")

// ⛔ JALUR BACA DIBUANG - ia kode mati, dan saya sudah pernah berjanji tidak
// mengulanginya.
//
// Ronde pertama melahirkan `HasilKomite`, `SumberHasilKomite`,
// `SumberHasilBelumDiputuskan`, `DenganSumber`, `BacaHasil`,
// `PeriksaHasilFinal`, dan `KeputusanKomite` - **nol** pemanggil produksi,
// nol handler, dan `Tambah` tidak pernah menyentuh satu pun. Itu persis
// jebakan `WajibWewenangKomite` tiket 07: antarmuka lahir tanpa pemanggil,
// lalu AC-nya tercentang atas bentuknya.
//
// Yang AC minta sudah berjalan tanpa semua itu: keputusan Komite ditulis ke
// `STS_REJECT` baris, dan baris dibaca `AmbilBaris` lalu ditampilkan sebagai
// KATA oleh layar - jalur yang sudah ada sejak tiket 01 dan 03. Rantai join
// tiga lompatan hanya menambah RINCIAN (komentar, penyetuju, tingkat), dan
// tidak satu pun AC memintanya.
//
// Pengetahuan XML-nya tidak hilang: ia hidup di bab pembacaan XML tiket 11,
// lengkap dengan nomor baris. Kode mati yang menyimpan pengetahuan tetap kode
// mati - kesimpulan yang sama dengan `DariTingkatProduksi` di tiket 12.
//
// ⛔ Dan penegakan "tidak diterapkan lebih awal" kini STRUKTURAL, bukan
// sekadar pemeriksaan: konteks ini tidak dapat menerapkan keputusan pada
// tingkat mana pun, sebab `repository.PeriksaBarisBaru` menolak baris baru
// yang berkeputusan dan berkas ini tidak memanggil satu pun pengubah status
// baris lama. Struktur yang tidak memungkinkan lebih kuat daripada gerbang
// yang memeriksa.

// BarisLanjutan membentuk baris putaran berikutnya bagi seorang peserta.
//
// ⛔ Inilah yang membuat klaim TIDAK TERMINAL (ADR-U-0011): penolakan
// menghasilkan putaran berikutnya, bukan akhir. Yang terminal adalah BARIS,
// bukan klaimnya.
//
// Pewarisan delapan kolom dipinjam dari `TambahBaris`, yang meniru
// `SetIndexAdjustmentList` - satu aturan, satu tempat.
func BarisLanjutan(p models.Peserta) (models.BarisAdjustment, error) {
	terakhir, ada := barisTerakhirPeserta(p)
	if !ada || terakhir.KodeStatus != models.KodeDitolak {
		punya := "tidak ada baris"
		if ada {
			punya = fmt.Sprintf("baris terakhir berkode %q", terakhir.KodeStatus)
		}
		return models.BarisAdjustment{}, fmt.Errorf(
			"%w: peserta %q %s", ErrBukanPenolakan, p.ID, punya)
	}

	salin := p
	if err := TambahBaris(&salin, models.BarisAdjustment{}); err != nil {
		return models.BarisAdjustment{}, err
	}
	baru := salin.Baris[len(salin.Baris)-1]
	// ⛔ Status Outstanding ditulis DI SINI, dan ia bukan warisan: baris baru
	// memulai putaran, ia tidak mewarisi keputusan putaran sebelumnya.
	baru.KodeStatus = models.KodeOutstanding
	// ⛔ Ketiganya milik putaran LAMA. Baris baru yang membawa `KOMITE_ID`
	// barisnya sendiri akan tampak sudah diserahkan, dan gerbang serah-ganda
	// tiket 10 akan menolaknya.
	//
	// ⚠️ TERUS TERANG: ketiga baris di bawah **no-op hari ini**. `TambahBaris`
	// berangkat dari baris kosong dan `WarisiKolom` hanya menyalin kedelapan
	// kolom warisan, jadi ketiganya memang sudah kosong. Menghapusnya tidak
	// membuat satu pun test merah - saya mencobanya. Yang dijaga bukan baris
	// ini melainkan PROPERTI-nya: ketika `WarisiKolom` dicoba menyalin
	// `KomiteID`, penjaga tiket 11 langsung merah. Keduanya dipertahankan
	// sebagai pagar kedua, dan sifat no-op-nya dinyatakan supaya tidak ada
	// yang mengira ia sudah terbukti.
	baru.ID = ""
	baru.NomorAkseptasi = ""
	baru.KomiteID = ""
	baru.TanggalAkseptasi = time.Time{}
	return baru, nil
}

// pesertaMilikKlaim memastikan peserta itu benar-benar milik klaimnya.
//
// ⛔ Pengenal peserta datang dari JALUR URL, dan jalur URL datang dari siapa
// saja.
func pesertaMilikKlaim(ctx context.Context, baca *repository.KlaimLife,
	klaimID, pesertaID string) error {

	semua, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	for _, p := range semua {
		if p.ID == pesertaID {
			return nil
		}
	}
	return fmt.Errorf("%w: peserta %q bukan milik klaim %q",
		ErrPermintaanTidakSah, pesertaID, klaimID)
}

// barisTerakhirPeserta mengembalikan SALINAN baris paling akhir.
//
// ⛔ Salinan, bukan penunjuk. Ronde pertama mengembalikan
// `&p.Baris[len(p.Baris)-1]` - dan meski `p` diterima sebagai NILAI, slice
// berbagi array yang sama dengan pemanggilnya, sehingga penunjuk itu mengarah
// ke baris ASLI. Menulis lewatnya menimpa keputusan Komite dari konteks yang
// hanya berhak membacanya, dan test yang mengandalkan "parameter nilai"
// tidak akan menangkapnya.
func barisTerakhirPeserta(p models.Peserta) (models.BarisAdjustment, bool) {
	if len(p.Baris) == 0 {
		return models.BarisAdjustment{}, false
	}
	return p.Baris[len(p.Baris)-1], true
}

// Putaran membuka putaran adjustment berikutnya.
type Putaran struct {
	svc   *Service
	jejak Jejak
}

// Putaran menyusun layanan itu dengan ketergantungan yang gagal terang.
func (s *Service) Putaran() *Putaran {
	return &Putaran{svc: s, jejak: JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekam jejaknya.
func (pt *Putaran) DenganJejak(j Jejak) *Putaran {
	salin := *pt
	salin.jejak = j
	return &salin
}

// Tambah menambahkan satu baris adjustment putaran berikutnya - tombol `Add`
// b17937 (`ClaimLifeDetailGCNM.xml`, tampil bila `pyPosition=='ReasLifeSPV'`
// b18160).
//
// ⛔ RALAT 29-09-2026 (GILIRAN-14 butir bp, meralat bo). GILIRAN-13 membuat
// rute ini juga melahirkan baris PERTAMA pada grid kosong. Keliru jalan
// lahirnya: `SavePesertaClaim` 7.8 (b3671, hidup) melahirkan baris pertama
// saat Submit Register, sehingga grid peserta terpilih tidak pernah kosong.
// Cabang itu dibuang; kini pendaftaran yang melahirkannya
// (`LahirkanBarisPendaftaran`).
//
// ⛔ Ia TIDAK menulis keputusan. Yang ditulisnya baris BARU berstatus
// Outstanding; keputusan atas baris lama tetap milik Komite dan tidak
// disentuh. Penjaga statik tiket 11 memastikan berkas ini tidak pernah
// memanggil penulis status.
func (pt *Putaran) Tambah(ctx context.Context, pelaku Pelaku,
	klaimID, pesertaID string, saat time.Time) error {

	if err := WajibIdentitas(pelaku); err != nil {
		return err
	}
	// Menambah baris lanjutan adalah menyimpan ke Outstanding - perannya sama
	// dengan penyimpanan Outstanding, bukan peran penolak.
	if err := WajibPeran(pelaku, PeranSimpanOutstanding); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi",
			ErrPermintaanTidakSah)
	}
	if !pt.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}

	// ⛔ BUTIR bb: kasus yang sudah ditutup tidak dapat diubah lagi.
	// Satu pintu untuk seluruh rute pengubah - lihat
	// services.PastikanKasusTerbuka, yang pemanggilannya ditagih penjaga
	// statik. Diperiksa SESUDAH wewenang: pemanggil yang tidak berhak tidak
	// berhak pula tahu keadaan kasusnya.
	if err := pt.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	baca := repository.NewKlaimLife(pt.svc.db)
	// ⛔ GERBANG TAHAP - GILIRAN-14 (tinjauan): "Add tetap = putaran
	// bergerbang ReasLifeSPV b18160". `pyPosition=='ReasLifeSPV'` berarti kasus
	// dipegang SPV - tahap Claim Analis. Sebelumnya gerbang ini hanya di layar;
	// layar bukan pagar.
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if tahap != models.TahapClaimAnalis {
		return fmt.Errorf("%w: tahap %s", ErrTahapTanpaAddAdjustment, tahap)
	}
	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return err
	}
	daftar, ada := perBaris[pesertaID]
	if !ada {
		// ⛔ `AmbilBaris` menggabung ke baris adjustment, jadi peserta TANPA
		// baris tidak muncul di sini. Sejak butir bp peserta seperti itu
		// hanya ada pada klaim LAMA (A4 migrasi data). Kepemilikannya
		// diperiksa lewat `AmbilPeserta`, supaya jawabannya jujur: 409 "tidak
		// ada baris" dari `BarisLanjutan`, bukan 400 "bukan milik klaim".
		if err := pesertaMilikKlaim(ctx, baca, klaimID, pesertaID); err != nil {
			return err
		}
	}

	baru, err := BarisLanjutan(models.Peserta{ID: pesertaID, Baris: daftar})
	if err != nil {
		return err
	}

	pohon := repository.NewPohonKlaim(pt.svc.db)
	return pt.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		id, err := pohon.SisipkanBaris(ctx, tx, pesertaID, baru)
		if err != nil {
			return err
		}
		// ⛔ PENANDA DIPILIH DIPULIHKAN, di transaksi yang sama.
		//
		// `Tolak` mencabutnya (`IS_CHECK='false'`) supaya peserta dapat
		// dipilih ulang - dan memilih ulang itu PERSIS yang terjadi di sini.
		// `[terverifikasi]` `Komite Claim Life/Activity/
		// KomitePostAdjustment.xml` pecahan baris 1900, 2117, 4896, 5341,
		// 7622, dan 8067: keenam prasyaratnya menuntut `.IsCheck = true`.
		//
		// ⚠️ SENSUS, jendelanya DINAMAI. Di berkas pecahan `KomitePostAdjustment.xml`
		// kata `IsCheck` muncul di **13** baris; **4** di antaranya tag
		// `<pyStepsPreCondParamsWhen>` langsung (4896, 5341, 7622, 8067), dan **2**
		// lagi tersimpan sebagai `rowdata REPEATINGINDEX="pyStepsPreCondParamsWhen"`
		// (1900, 2117) - masing-masing berpasangan dengan kembaran `<pyExpression>`
		// (1897, 2114) yang berisi kondisi yang sama.
		//
		// Jendela yang dipakai: **LANGKAH PRASYARAT YANG BERBEDA** → **6**. Bila
		// yang dihitung ELEMEN XML pembawa kondisi, angkanya **8**; bila seluruh
		// penyebutan, **13**. Ketiganya benar untuk pertanyaan yang berbeda, dan
		// angka 6 tidak berarti apa-apa tanpa kalimat ini.
		//
		// Audit: `py` + pemecah `><` → `>\n<`, lalu cacah baris ber-`IsCheck` yang
		// juga ber-`pyStepsPreCondParamsWhen` (4) dan yang ber-`rowdata` (2).
		//
		// Tanpa pemulihan ini putaran berikutnya LAHIR TETAPI TIDAK DAPAT
		// DIAMBIL SIAPA PUN - dan justru itu yang tiket ini janjikan.
		if err := baca.PasangPenandaDipilih(ctx, tx, pesertaID); err != nil {
			return err
		}
		// Header mencerminkan baris TERAKHIR - kini baris yang baru lahir,
		// yang Outstanding dan belum bernomor akseptasi.
		//
		// ⚠️ `[terbuka]` Pada klaim BERPESERTA BANYAK, "baris terakhir" belum
		// ditetapkan artinya - butir yang sudah terbuka sejak tiket 04. Yang
		// dilakukan di sini mengikuti kontrak tiket 04 apa adanya; bila kelak
		// artinya ditetapkan lain, tempat ini ikut berubah.
		if err := baca.CerminkanHeader(ctx, tx, klaimID,
			models.KodeOutstanding, ""); err != nil {
			return err
		}
		return pt.jejak.Rekam(ctx, tx, CatatanJejak{
			AdjustmentID: id,
			KlaimID:      klaimID,
			Dari:         "",
			Ke:           models.KodeOutstanding,
			AkunID:       pelaku.AkunID,
			Waktu:        saat,
		})
	})
}
