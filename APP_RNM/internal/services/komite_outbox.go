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
// ⛔ RALAT (temuan /code-review Spec, 28-09-2026) - langkah 12 TERBACA dari
// transisi tiap `When`-nya. Di seluruh `KomitePostAdjustment` pola normal
// adalah `WhenTrue 2` (lanjut) / `WhenFalse 3` (lewati langkah) - 48 lawan 32
// kemunculan. Baris `Type=="TP"||"TR"` langkah 12 TERBALIK: `WhenTrue 3`,
// `WhenFalse 2`. Jadi Kasir berjalan HANYA bila:
//
//	Setuju di tingkat akhir  DAN  Type BUKAN TP/TR  DAN  IsKPR == "KPR"  (DAN IsPEGAPROD)
//
// Bacaan sebelumnya ("gerbang tingkat akhir saja") mengantre Kasir untuk
// polis treaty dan klaim non-KPR - terlalu luas untuk efek yang memindahkan uang.
//
// ⛔ `IsPEGAPROD` TIDAK menggerbangi pengantrean - ia menggerbangi PENGIRIMAN
// (tiket 07). Niat efek yang tidak tercatat di non-produksi tidak dapat diuji
// di mana pun.

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// KasirBerlaku adalah gerbang langkah 12, VERBATIM - MURNI.
func KasirBerlaku(a models.AkibatKeputusanKomite, tipe, isKPR string) bool {
	t := strings.TrimSpace(tipe)
	return a.AkseptasiAkhir && t != models.TipePLTreatyProposal && t != models.TipePLTreatyRealisasi &&
		strings.TrimSpace(isKPR) == "KPR"
}

// EfekKeputusanKomite menyusun daftar efek satu keputusan - MURNI.
func EfekKeputusanKomite(a models.AkibatKeputusanKomite, tipe, isKPR string) []efekKomite {
	var e []efekKomite
	if a.AkseptasiAkhir {
		e = append(e, efekKomite{JenisEfekKomiteArasapas, KunciArasapasLife})
	}
	// Langkah 11 - setiap keputusan. SMTP langsung, bukan M_LINK_SERVICE
	// (lihat `EfekEmail` Claim Life), jadi tanpa kunci kategori.
	e = append(e, efekKomite{Jenis: JenisEfekKomiteEmail})
	if KasirBerlaku(a, tipe, isKPR) {
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
	// `TempOpenPage.PolicyDataLife.Type` - klaim induk.
	tipe, err := repository.NewKlaimLife(svc.db).TypeKlaim(ctx, kasus.Baris.KlaimID)
	if err != nil {
		return nil, err
	}
	for _, e := range EfekKeputusanKomite(a, tipe, kasus.IsKPR) {
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
