package services

// Baris AdjustmentList - tiket 03.
//
// Untuk apa berkas ini: satu baris AdjustmentList adalah UNIT KEPUTUSAN mesin
// status (ADR-U-0011) - bukan klaimnya, bukan pesertanya. Berkas ini membentuk
// baris itu; perhitungan pecahan retronya ada di spreading.go.
//
// Dibaca sesudah: pendaftaran.go.
//
// Istilah:
//   - spreading : pembagian satu baris klaim ke beberapa reinsurer retro.
//   - retrosesi : meneruskan sebagian risiko ke reinsurer lain.

import (
	"errors"
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/inti"
	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
)

// ErrBarisTidakSah menandai baris adjustment yang tidak dapat dibentuk.
var ErrBarisTidakSah = errors.New("services: baris adjustment tidak sah")

// kolomDiwarisi adalah kedelapan kolom yang diwarisi baris kedua dan seterusnya.
//
// `[terverifikasi]` dari `SetIndexAdjustmentList` langkah 3 (baris 570-744 di
// berkas pecahan): `.AdjustmentList(<LAST>).X = .AdjustmentList(1).X` untuk
// tepat kedelapan nama ini.
//
// ⛔ STS_REJECT SENGAJA TIDAK ADA di daftar ini, dan itu inti AC 5: status
// adalah keputusan per baris. Baris baru yang mewarisi status orang lain
// berarti keputusan yang tidak pernah diambil siapa pun.
var kolomDiwarisi = []string{
	"SHARE_NUSANTARA_RE", "CEDING_RETENTION", "SUM_REASURED", "SUM_INSURED",
	"SHARE_RETRO", "RETROCEDED_SHARE", "CURRENCYID", "CURRENCY",
}

// KolomDiwarisi mengeluarkan salinan daftar itu untuk test dan dokumentasi.
func KolomDiwarisi() []string {
	out := make([]string, len(kolomDiwarisi))
	copy(out, kolomDiwarisi)
	return out
}

// WarisiKolom menyalin kedelapan kolom dari baris PERTAMA ke baris baru.
//
// Meniru `SetIndexAdjustmentList` langkah 3 (baris 570-744):
// `.AdjustmentList(<LAST>).X = .AdjustmentList(1).X` untuk tepat kedelapan
// nama di `kolomDiwarisi` - tidak lebih, tidak kurang.
//
// ⛔ Status TIDAK diwarisi, dan tanggal akseptasi TIDAK distempel: keduanya
// milik aksi, bukan milik baris sebelumnya (AC 5, AC 43). `CLAIM_AMOUNT` juga
// tidak: tiap baris punya jumlah klaimnya sendiri.
func WarisiKolom(pertama models.BarisAdjustment, baru *models.BarisAdjustment) {
	if baru == nil {
		return
	}
	baru.ShareNusantaraRe = pertama.ShareNusantaraRe
	baru.CedingRetention = pertama.CedingRetention
	baru.SumReasured = pertama.SumReasured
	baru.SumInsured = pertama.SumInsured
	baru.ShareRetro = pertama.ShareRetro
	baru.RetrocededShare = pertama.RetrocededShare
	baru.CurrencyID = pertama.CurrencyID
	baru.JumlahKlaim.Currency = pertama.JumlahKlaim.Currency

	baru.KodeStatus = ""
	baru.TanggalAkseptasi = time.Time{}
}

// nilaiAtauNol memperlakukan desimal kosong sebagai nol untuk perkalian.
//
// Kosong berarti "tidak diisi" (ADR-U-0027), dan dalam perkalian share itu
// sama artinya dengan nol bagian - bukan galat.
func nilaiAtauNol(d *apd.Decimal) *apd.Decimal {
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}

// PeranSimpanOutstanding adalah peran yang boleh menyimpan ke Outstanding.
//
// `[terverifikasi]` `Section/AdjustmentDetail_Section.xml`: kedua tombol
// `SaveOutstandingLife_Act` (baris 16249 dan 16396) bergerbang
// `pyWorkPage.pyPosition =='ReasLifeSPV'` (baris 16468).
//
// ⚠️ Tiket 03 semula menulis `ReasLifeAdmin`. XML menang atas tiket (aturan
// work owner 26-09-2026), dan tiket sudah diralat. Yang bergerbang
// `ReasLifeAdmin` adalah REJECT Outstanding - tiket 05, kondisi baris 15399.
const PeranSimpanOutstanding = inti.PeranSPV

// TambahBaris membentuk satu baris adjustment baru bagi seorang peserta.
//
// Meniru `SetIndexAdjustmentList`: baris kedua dan seterusnya mewarisi
// delapan kolom dari baris PERTAMA, peserta ditandai dipilih-untuk-diklaim,
// dan status TIDAK diwarisi.
//
// ⛔ Fungsi ini MURNI dan mengubah pesertanya di tempat: penulisan ke basis
// data milik repository, dan aturannya milik sini. Memisahkan keduanya membuat
// aturan dapat diuji tanpa Oracle - itulah seam tiket ini.
func TambahBaris(p *models.Peserta, baru models.BarisAdjustment) error {
	if p == nil {
		return fmt.Errorf("%w: peserta kosong", ErrBarisTidakSah)
	}
	// `SetIndexAdjustmentList` langkah 1 baris 328: `.IsCheck = true`. Penanda
	// dipilih-untuk-diklaim lahir saat baris dibuat, bukan di layar.
	p.IsCheck = "true"

	if len(p.Baris) > 0 {
		WarisiKolom(p.Baris[0], &baru)
	} else {
		// Baris pertama pun tidak membawa status maupun tanggal akseptasi.
		baru.KodeStatus = ""
		baru.TanggalAkseptasi = time.Time{}
		if baru.JumlahKlaim.Currency == "" {
			baru.JumlahKlaim.Currency = p.MataUang
		}
	}
	p.Baris = append(p.Baris, baru)
	return nil
}

// angkaDesimalPendaftaran adalah skala `@divide(…,1,4)` langkah 7.8.
//
// `[terverifikasi]` `SavePesertaClaim.xml` b3697-b3849: pembagian dengan SATU
// pada empat angka - pembulatan yang ditulis sebagai pembagian. KEANEHAN
// WARISAN, disalin apa adanya; bentuk yang sama ditiru rekap PremiumList
// (`models.AngkaDesimalRekap`), dengan konteks desimal yang sama
// (`utils.DecimalContext`, setengah ke atas).
const angkaDesimalPendaftaran = 4

// BarisPendaftaran membentuk baris adjustment yang lahir saat Submit Register
// - GILIRAN-14 butir bp (`[DIPUTUSKAN; veto work owner]`).
//
// `[terverifikasi]` `Activity/SavePesertaClaim.xml` (tombol `Submit`
// `Section/InputRegisterClaimLife.xml` b27369 -> b27393), langkah 7.8 b3671
// "Set property AdjustmentList" - HIDUP (`pyStepsBlockName` kosong),
// `pyStepsPreCondition=true` b3691, WHEN `.IsCheck=="true"` b3919 (True=2
// lanjut, False=3 lewati). Ia menulis `PremiumListDetail(<LAST>).
// AdjustmentList(<LAST>)` - baris PERTAMA peserta yang baru saja 7.7
// tambahkan - dengan delapan medan dari baris sumber yang sama:
//
//	CEDING_RETENTION b3696   SHARE_RETRO      b3808
//	SHARE_NUSANTARA_RE b3742 CLAIM_AMOUNT     b3828
//	SUM_INSURED b3768        RETROCEDED_SHARE b3848
//	SUM_REASURED b3788       CURRENCY         b3868
//
// Tujuh yang pertama `@divide(@toDecimal(@replaceAll(.X,",",".")),1,4)`.
// (Catatan pembaca Go, BUKAN bunyi korpus: `@replaceAll` tidak berbuat
// apa-apa di sini, sebab `repository.kolomSalin` membaca sumbernya lewat
// `TO_CHAR(…,'TM9')` bertitik desimal.) `SHARE_NUSANTARA_RE` b3743 memakai
// `@if` yang SAMA dengan peserta 7.7 b2845, jadi nilai peserta - yang sudah
// memilihnya lewat `repository.ShareNusantaraReTeks` - dipakai apa adanya.
//
// ⛔ Status TIDAK ditulis: 7.8 tidak menyentuh `STS_REJECT`, dan `Save to RNM`
// 22.1.3.2 yang kemudian menulis "0" bagi baris tanpa status.
//
// ⛔ Kolom sumber yang KOSONG dibaca NOL - `kosongJadiNol78`, keputusan work
// owner 29-09-2026, HANYA di langkah ini (lihat komentar fungsi itu).
func BarisPendaftaran(p models.Peserta) (models.BarisAdjustment, bool, error) {
	// `.IsCheck=="true"` - teks, persis (b3919).
	if p.IsCheck != models.PenandaDipilih {
		return models.BarisAdjustment{}, false, nil
	}
	var b models.BarisAdjustment
	for _, m := range []struct {
		nama string
		dari uang.Money
		ke   *uang.Money
	}{
		{"CEDING_RETENTION", p.CedingRetention, &b.CedingRetention},
		{"SHARE_NUSANTARA_RE", p.ShareNusantaraRe, &b.ShareNusantaraRe},
		{"SUM_INSURED", p.SumInsured, &b.SumInsured},
		{"SUM_REASURED", p.SumReasured, &b.SumReasured},
		{"SHARE_RETRO", p.ShareRetro, &b.ShareRetro},
		{"CLAIM_AMOUNT", p.JumlahKlaim, &b.JumlahKlaim},
		{"RETROCEDED_SHARE", p.RetrocededShare, &b.RetrocededShare},
	} {
		*m.ke = uang.Money{Currency: p.MataUang}
		empat, err := bulatEmpatPendaftaran(kosongJadiNol78(m.dari), p.ID, m.nama)
		if err != nil {
			return models.BarisAdjustment{}, false, err
		}
		m.ke.Amount = empat
	}
	return b, true, nil
}

// ErrPembulatanPendaftaran - satu nilai pendaftaran tidak dapat dibulatkan
// ke empat angka (7.7 peserta atau 7.8 baris).
var ErrPembulatanPendaftaran = errors.New("services: nilai pendaftaran tidak dapat dibulatkan")

// bulatEmpatPendaftaran membulatkan satu nilai `@divide(…,1,4)` dengan `bulat`
// bersama (spreading.go), dan menamai medannya bila gagal.
func bulatEmpatPendaftaran(d *apd.Decimal, pesertaID, medan string) (*apd.Decimal, error) {
	empat, err := bulat(utils.DecimalContext(), d, angkaDesimalPendaftaran)
	if err != nil {
		return nil, fmt.Errorf("%w: peserta %q %s: %w", ErrPembulatanPendaftaran, pesertaID, medan, err)
	}
	return empat, nil
}

// kosongJadiNol78 membaca satu medan uang sumber langkah 7.8: KOSONG -> 0.
//
// ⛔ PENYIMPANGAN BERTANGGAL TERHADAP ADR-U-0027 - keputusan work owner
// 29-09-2026 (GILIRAN-15, "ikuti rekomendasi" - "ikut Pega hanya di langkah
// 7.8"): ketujuh medan `@divide(@toDecimal(…),1,4)` langkah 7.8
// `SavePesertaClaim` (baris adjustment yang lahir saat Submit Register)
// membaca sumber kosong sebagai nol. `[dugaan]` itulah yang `@toDecimal("")`
// Pega hasilkan - perilaku fungsi Pega itu tidak ada di korpus; yang pasti
// adalah KEPUTUSANNYA. Baris pertama inilah yang diwarisi setiap putaran
// berikutnya. ADR-U-0027 ("kosong bukan nol") tetap berlaku di SELURUH tempat
// lain - termasuk pembulatan peserta 7.7 (`BulatkanPesertaPendaftaran`).
//
// ⛔ Satu fungsi konversi, SATU pemanggil (`BarisPendaftaran`), dan namanya
// terikat langkahnya supaya tidak tampak sebagai peniru `@toDecimal` umum.
// Penjaga `TestKosongJadiNolHanyaSatuPemanggil` menagih satu pemanggil itu.
func kosongJadiNol78(m uang.Money) *apd.Decimal {
	if m.Kosong() {
		return apd.New(0, 0)
	}
	return m.Amount
}

// BulatkanPesertaPendaftaran meniru pembulatan peserta langkah 7.7 - OQ-N10
// ditutup (keputusan work owner 29-09-2026).
//
// `[terverifikasi]` `SavePesertaClaim.xml` 7.7 b2744-b3631 (hidup, WHEN
// `.IsCheck=="true"` b3631; rentang "b3600-b3671" brief GILIRAN-15 keliru -
// b3671 adalah awal 7.8): sepuluh medan `@divide(@toDecimal(@replaceAll(.X,",",".")),1,4)` -
// GROSS_PREMIUM b2770, NET_PREMIUM b2824, SHARE_NUSANTARA_RE b2845 (`@if`
// pemilih GROSS sudah diterapkan `repository.ShareNusantaraReTeks`),
// SUM_INSURED b2872, CEDING_RETENTION b2893, SUM_REASURED b3107, EM_PERCENT
// b3134, CLAIM_AMOUNT b3281, SHARE_RETRO b3401, RETROCEDED_SHARE b3561.
//
// ⛔ KOSONG TETAP KOSONG di sini (ADR-U-0027): penyimpangan "kosong = nol"
// diputuskan hanya untuk 7.8.
func BulatkanPesertaPendaftaran(peserta []models.Peserta) error {
	for i := range peserta {
		p := &peserta[i]
		for _, m := range []struct {
			nama string
			ke   *uang.Money
		}{
			{"GROSS_PREMIUM", &p.GrossPremium}, {"NET_PREMIUM", &p.NetPremium},
			{"SHARE_NUSANTARA_RE", &p.ShareNusantaraRe}, {"SUM_INSURED", &p.SumInsured},
			{"CEDING_RETENTION", &p.CedingRetention}, {"SUM_REASURED", &p.SumReasured},
			{"CLAIM_AMOUNT", &p.JumlahKlaim}, {"SHARE_RETRO", &p.ShareRetro},
			{"RETROCEDED_SHARE", &p.RetrocededShare},
		} {
			if m.ke.Kosong() {
				continue
			}
			empat, err := bulatEmpatPendaftaran(m.ke.Amount, p.ID, m.nama)
			if err != nil {
				return err
			}
			m.ke.Amount = empat
		}
		if p.EMPercent.Value != nil {
			empat, err := bulatEmpatPendaftaran(p.EMPercent.Value, p.ID, "EM_PERCENT")
			if err != nil {
				return err
			}
			p.EMPercent.Value = empat
		}
	}
	return nil
}

// LahirkanBarisPendaftaran memasang baris 7.8 pada setiap peserta terpilih.
//
// ⛔ Dipanggil pendaftaran di DALAM transaksinya, sebelum pohon disimpan: satu
// transaksi untuk peserta dan barisnya, seperti satu `Obj-Save` Pega. Baris
// ini `.AdjustmentList(1)` yang diwarisi putaran berikutnya
// (`SetIndexAdjustmentList` langkah 3, `WarisiKolom`).
func LahirkanBarisPendaftaran(peserta []models.Peserta) error {
	for i := range peserta {
		b, lahir, err := BarisPendaftaran(peserta[i])
		if err != nil {
			return err
		}
		if lahir {
			peserta[i].Baris = append(peserta[i].Baris, b)
		}
	}
	return nil
}

// TandaiOutstanding menuliskan STS_REJECT = 0 pada baris yang BELUM pernah
// disimpan, dan mengembalikan cacah baris yang tersentuh.
//
// `[terverifikasi]` `SaveOutStandingLife_Act` langkah 22.1.3.1 dan 22.1.3.2:
// insert-nya bergerbang `.PrintFaceClaim==""` dan baru sesudah itu
// `.PrintFaceClaim=1` dan `.STS_REJECT=0` ditulis. Jadi nol adalah nilai
// MENURUT AKSI, dan baris yang sudah punya status tidak disentuh ulang.
//
// ⛔ ACCEPTATION_DATE TIDAK distempel di sini (AC 43). Menyimpan ke Outstanding
// bukan mengaksep; tanggal akseptasi milik aksi akseptasi yang sebenarnya.
func TandaiOutstanding(pohon *models.PohonKlaim) int {
	if pohon == nil {
		return 0
	}
	// Aturannya hidup di TandaiOutstandingKlaim (tiket 04) - satu tempat.
	return TandaiOutstandingKlaim(&pohon.Klaim)
}
