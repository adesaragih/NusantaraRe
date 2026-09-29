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
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
)

// ErrBarisTidakSah menandai baris adjustment yang tidak dapat dibentuk.
var ErrBarisTidakSah = errors.New("services: baris adjustment tidak sah")

// ErrPesertaSudahDiputus - peserta sudah diaksep atau ditolak; isian layar
// Detail-nya beku (tujuh gerbang `.STS_REJECT`, lihat BarisPertama).
var ErrPesertaSudahDiputus = errors.New(
	"services: peserta sudah diputus, isian layar Detail-nya tidak dapat diubah")

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
const PeranSimpanOutstanding = PeranSPV

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

// BarisPertama membentuk baris adjustment PERTAMA seorang peserta - butir bo
// (GILIRAN-13, `[DIPUTUSKAN; veto work owner]`).
//
// `[terverifikasi]` `Section/ClaimLifeDetailGCNM.xml`: grid `.AdjustmentList`
// b17126, tombol `Add` b17937 -> `addRow` b17947 + `SetIndexAdjustmentList`
// b17991. Pada grid KOSONG, langkah 3 activity itu (`.AdjustmentList(<LAST>).X
// = .AdjustmentList(1).X`, b570-744) menyalin baris ke DIRINYA SENDIRI - jadi
// baris pertama lahir KOSONG seluruhnya, dan itulah yang dikembalikan.
//
// ⛔ Bukan `TambahBaris`. Cabang baris-pertama di sana mengisi mata uang dari
// peserta; butir bo memutuskan "seperti XML", dan XML tidak mengisinya.
// Penanda dipilih (`SetIndexAdjustmentList` langkah 1 b328, `.IsCheck = true`)
// dipasang layanannya di transaksi yang sama - ia milik PESERTA, bukan baris.
//
// ⛔ Gerbang `.STS_REJECT` peserta: `pyDisabledWhen` `.STS_REJECT=='1' ||
// .STS_REJECT=='2'` tujuh kali di layar yang sama (b2628, b4682, b5059,
// b5870, b6152, b7335, b15234) - peserta yang sudah diputus membekukan isian
// layar Detail. `Add` b17937 sendiri TIDAK membawanya (hanya syarat tampil
// b18160); butir bo memberlakukannya pada baris pertama. Pada jalur putaran
// gerbang ini TIDAK berlaku: peserta yang barisnya ditolak berkode "2", dan
// justru dialah yang dibuka putaran berikutnya.
//
// ⛔ Pembacaan seluruh activity Claim Life yang menyebut `AdjustmentList`,
// sebagai pohon (`pyStepsBlockName` dicetak; `//` = mati, dan tidak satu pun
// langkah di bawah bertanda itu): `SaveInsuredClaim_Act` 2.2 b1341 menulis
// `TempDetail.pxResults(<LAST>).AdjustmentList(<LAST>)` - halaman sementara
// unggahan; `SavePesertaClaim` 7.8 b3671 (hidup, WHEN b3919) menulis baris
// pertama saat PENDAFTARAN; `SaveOutStandingLife_Act` 22.1 b10548 MELEWATI
// peserta tanpa baris; `SpreadingClaimLife_Act` 7.1/8.2 dan `SaveAdjustment_Act`
// 1.5/1.6, `DeletePesertaClaimLife` 1.1, serta `serviceInsertArasapasClaimLife_act`
// 1.1 b311 hanya menyunting atau membaca baris yang sudah ada. Yang
// melahirkan baris hanya dua: `Add` dan pendaftaran 7.8.
func BarisPertama(p models.Peserta) (models.BarisAdjustment, error) {
	if len(p.Baris) > 0 {
		return models.BarisAdjustment{}, fmt.Errorf(
			"%w: peserta %q sudah berbaris %d - baris berikutnya lahir lewat putaran",
			ErrBarisTidakSah, p.ID, len(p.Baris))
	}
	if models.DiagnosaTerkunci(p.KodeStatus) {
		return models.BarisAdjustment{}, fmt.Errorf("%w: peserta %q berkode %q",
			ErrPesertaSudahDiputus, p.ID, strings.TrimSpace(p.KodeStatus))
	}
	return models.BarisAdjustment{}, nil
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
