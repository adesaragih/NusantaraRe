package services

// Transactional outbox keputusan Komite - tiket 06 Komite Claim Life.
//
// Untuk apa berkas ini: efek keluar keputusan komite DIANTRE di transaksi yang
// SAMA dengan keputusannya (ADR-0015 - menyimpang dari ADR-0008 Claim Life:
// Kasir memindahkan uang, jadi efeknya wajib sampai, bukan boleh hilang).
// Pengirimannya tiket 07.
//
// `[terverifikasi]` `KomitePostAdjustment` sesudah langkah 7 `Obj-Save`:
//
//	8  `InsertJsonClaimLife_Act`              JSON - DIBUANG (keputusan 2026-09-16)
//	10 `serviceInsertArasapasClaimLife_act`   AcceptStatus = 1 && Count == Loop (b8648), IsPEGAPROD
//	11 `SendEmailKlaimLife`                   IsPEGAPROD SAJA - setiap keputusan
//	12 `HitServiceToKasirKMTLife_Act`         AcceptStatus = 1 && Count == Loop (b8887),
//	                                          Type TP||TR, IsKPR=="KPR"  (lihat catatan)
//
// ⚠️ RALAT tiket 06: AC "tingkat bukan-terakhir tidak menghasilkan entri
// outbox" dibantah langkah 11 - email tidak bergerbang tingkat, jadi ia
// diantre pada SETIAP keputusan.
//
// ⚠️ `[terbuka — pemilik ekspor]` langkah 12 punya TIGA `When` (tingkat akhir,
// `Type TP||TR`, `IsKPR=="KPR"`) dengan transisi yang ekspornya tidak
// jelaskan (lanjut/lewati). Membacanya "ketiganya AND" berarti Kasir hanya
// untuk polis treaty ber-KPR - tafsiran yang terlalu sempit untuk ditebak.
// Yang dibangun: gerbang tingkat akhir Setuju saja; kedua syarat lain dicatat.
//
// ⛔ `IsPEGAPROD` TIDAK menggerbangi pengantrean - ia menggerbangi PENGIRIMAN
// (tiket 07). Niat efek yang tidak tercatat di non-produksi tidak dapat diuji
// di mana pun.

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// ModulKomiteLife mengisi kolom `MODUL` outbox untuk baris Komite.
const ModulKomiteLife = "KOMITELIFE"

// Jenis efek Komite - kolom `JENIS_EFEK`.
const (
	JenisEfekKomiteArasapas = "arasapas-komite"
	JenisEfekKomiteEmail    = "email-komite"
	JenisEfekKomiteKasir    = "kasir-komite"
)

// KunciKasirKomite adalah kunci endpoint Kasir.
//
// `[terverifikasi]` `Activity/HitServiceToKasirKMTLife_Act.xml` `GetLinkService`
// `Kategori_1 = "Kasir"`, `Kategori_2 = "insertAllPaymentKasir"`. ⚠️ Langkah
// `Connect-REST`-nya sendiri ber-`//` dengan deskripsi "kalau diserver dev
// jangan dijalanin" - panggilan nyatanya menunggu persetujuan manusia.
var KunciKasirKomite = KunciLayanan{Kategori1: "Kasir", Kategori2: "insertAllPaymentKasir"}

// muatanOutboxKomite adalah JSON kolom `MUATAN` baris Komite.
//
// ⛔ Pengenal, waktu, dan KUNCI KATEGORI - tidak pernah URL (ADR-0013); nol
// nama, nol email.
type muatanOutboxKomite struct {
	KasusID      string `json:"kasus_id"`
	KlaimID      string `json:"klaim_id"`
	AdjustmentID string `json:"adjustment_id"`
	AkunID       string `json:"akun_id"`
	Waktu        string `json:"waktu"`
	Kategori1    string `json:"kategori_1,omitempty"`
	Kategori2    string `json:"kategori_2,omitempty"`
}

// efekKomite adalah satu efek yang diantre.
type efekKomite struct {
	Jenis string
	Kunci KunciLayanan
}

// EfekKeputusanKomite menyusun daftar efek satu keputusan - MURNI.
func EfekKeputusanKomite(a models.AkibatKeputusanKomite) []efekKomite {
	var e []efekKomite
	if a.AkseptasiAkhir {
		e = append(e, efekKomite{JenisEfekKomiteArasapas, KunciArasapasLife})
	}
	// Langkah 11 - setiap keputusan. SMTP langsung, bukan M_LINK_SERVICE
	// (lihat `EfekEmail` Claim Life), jadi tanpa kunci kategori.
	e = append(e, efekKomite{Jenis: JenisEfekKomiteEmail})
	if a.AkseptasiAkhir {
		e = append(e, efekKomite{JenisEfekKomiteKasir, KunciKasirKomite})
	}
	return e
}

// RujukanEfekKomite menyusun kolom `RUJUKAN` - kunci anti-dobel tiket 07.
//
// ⛔ Email diantre pada SETIAP keputusan (langkah 11), jadi rujukannya membawa
// TINGKAT: `KMTLF-…#T2`. Tanpa itu email tingkat 1 yang sudah `selesai`
// membuat email tingkat 2..n dianggap "sudah terkirim" dan tidak pernah
// dikirim (temuan /code-review giliran 10). Arasapas dan Kasir hanya sekali
// per kasus - rujukannya kasus itu sendiri.
func RujukanEfekKomite(kasusID, jenis string, tingkat int) string {
	if jenis == JenisEfekKomiteEmail {
		return kasusID + "#T" + strconv.Itoa(tingkat)
	}
	return kasusID
}

// antreEfekKomite menulis efek keputusan ke outbox, DI DALAM transaksinya.
//
// Mengembalikan jenis efek yang diantre - "tersimpan, belum tuntas" (ADR-0015).
func antreEfekKomite(ctx context.Context, svc *Service, tx *repository.Tx,
	kasus repository.KasusKomite, a models.AkibatKeputusanKomite,
	pelaku Pelaku, saat time.Time) ([]string, error) {

	pohon := repository.NewPohonKlaim(svc.db)
	diantre := []string{}
	for _, e := range EfekKeputusanKomite(a) {
		muatan, err := json.Marshal(muatanOutboxKomite{
			KasusID: kasus.Baris.KasusID, KlaimID: kasus.Baris.KlaimID,
			AdjustmentID: kasus.AdjID, AkunID: pelaku.AkunID,
			Waktu: saat.Format(time.RFC3339), Kategori1: e.Kunci.Kategori1,
			Kategori2: e.Kunci.Kategori2,
		})
		if err != nil {
			return nil, fmt.Errorf("services: merakit muatan outbox komite: %w", err)
		}
		// ⛔ ID baris outbox = `SEQ_LOG_SERVICE_RNM` - pengenal idempoten unik
		// sejak lahir (AC 21 spec); tiket 07 memakainya sebagai kunci anti-dobel.
		if _, err := pohon.AntreEfek(ctx, tx, LiniLife, ModulKomiteLife, e.Jenis,
			RujukanEfekKomite(kasus.Baris.KasusID, e.Jenis, a.TingkatDiputus),
			string(muatan), saat); err != nil {
			return nil, err
		}
		diantre = append(diantre, e.Jenis)
	}
	return diantre, nil
}
