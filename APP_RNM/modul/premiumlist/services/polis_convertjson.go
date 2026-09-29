package services

// `convertJsonNusareToProduction` lewat outbox - OQ-PL-14 DITUTUP 29-09-2026
// (GILIRAN-17) `[keputusan work owner]`: ditiru lewat outbox, pelaksana STUB.
//
// `[terverifikasi]` `Activity/serviceInsertArasapasLife_act.xml` langkah 5
// (b915, `Connect-REST` `ServiceName=convertJsonNusareToProduction` b963, POST
// b966-967; `pyStepsBlockName` kosong - hanya langkah 9 yang ter-remark) - satu-
// satunya panggilan keluar activity itu, dipanggil `InsertJsonPolisLife_Act`
// langkah 15 (b5168) bergerbang `IsPEGAPROD` (b5236). Parameternya
// `ConnectREST/ConvertJsonNusareToProduction.xml`:
//
//	b200 noPolis  <- .OfferFacIn.PolicyData.PolicyNo   (b199)
//	b209 caseId   <- .pzInsKey                          (b207)
//	b221 tglInput <- .pxCreateDateTime                  (b218)
//
// Efek segera sesudah commit (`EfekArasapasPolis`, `NamaEfekArasapasPolis`)
// adalah padanan langkah itu; yang gagal mendarat di outbox bersama
// (`T_LOG_SERVICE_RNM`, `MODUL = PREMIUMLISTLIFE`). Berkas ini menambah
// PELAKSANA baris outbox tersebut - pola `PelaksanaKomite`.
//
// ⛔ STUB (OQ-PL-11 `[terbuka - pemilik Arasapas]`): parameternya DIRAKIT dari
// tabel, lalu berhenti terang - nol panggilan keluar. Layanan hilir tampaknya
// membaca `JSON_POLIS` yang tidak lagi ditulis (pl1).
//
// ⚠️ `tglInput`: `.pxCreateDateTime` (saat kasus lahir) tidak punya kolom -
// `T_WORK_POLIS` tanpa waktu lahir. Yang dirakit waktu dari muatan outbox
// (saat keputusan), selisih yang sudah tercatat di OQ-PL-11 (tiket 05b).
// Nomor polis TIDAK disimpan di muatan outbox (pengenal + waktu saja); ia
// dibaca saat pelaksanaan.
//
// Dibaca sesudah: polis_efekkeluar.go, komite_pengirim.go (pola pelaksana).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
	"nusantarare/modul/premiumlist/repository"
)

// ParameterConvertJsonPolis - ketiga parameter REST, nama VERBATIM (b200,
// b209, b221).
type ParameterConvertJsonPolis struct {
	NoPolis  string `json:"noPolis"`
	CaseID   string `json:"caseId"`
	TglInput string `json:"tglInput"`
}

// RakitParameterConvertJson merakit parameter - MURNI. `tglInput` berbentuk
// stempel Pega `YYYYMMDDTHHMMSS.mmm GMT`, bentuk `.pxCreateDateTime`.
func RakitParameterConvertJson(caseID, noPolis string, saat time.Time) (ParameterConvertJsonPolis, error) {
	if strings.TrimSpace(caseID) == "" || strings.TrimSpace(noPolis) == "" {
		return ParameterConvertJsonPolis{}, fmt.Errorf(
			"%w: caseId dan noPolis wajib ada untuk convertJsonNusareToProduction", galat.ErrPermintaanTidakSah)
	}
	return ParameterConvertJsonPolis{
		NoPolis:  strings.TrimSpace(noPolis),
		CaseID:   strings.TrimSpace(caseID),
		TglInput: saat.UTC().Format("20060102T150405.000") + " GMT",
	}, nil
}

// PembacaNomorPolis membaca `NO_POLIS` polis dari pengenal work-nya - padanan
// `.OfferFacIn.PolicyData.PolicyNo`.
type PembacaNomorPolis interface {
	NomorPolisDariID(ctx context.Context, polisID string) (string, error)
}

// PelaksanaPremiumList menjalankan baris outbox PremiumList Life - STUB.
type PelaksanaPremiumList struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
	Pembaca    PembacaNomorPolis
}

// Laksanakan memenuhi `PelaksanaEfek`.
func (p PelaksanaPremiumList) Laksanakan(ctx context.Context, _ *db.Tx,
	b outbox.BarisEfekKeluar) error {

	if b.Modul != ModulPremiumListLife {
		return fmt.Errorf("%w: pelaksana PremiumList menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m outbox.MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox PremiumList tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case NamaEfekArasapasPolis:
	case NamaEfekAlarmPolis:
		return EfekAlarmEmailPolis{}.Jalankan(ctx, outbox.MuatanEfek{KlaimID: m.KlaimID})
	default:
		return fmt.Errorf("%w: jenis efek PremiumList %q", galat.ErrPermintaanTidakSah, b.Jenis)
	}
	if !p.Lingkungan.AdalahProduksi() {
		// Seperti pengirim Komite: nol panggilan keluar di luar produksi, dan
		// catatannya jujur - gagal bersebab, bukan `selesai`.
		return outbox.ErrPengirimStubNonProduksi
	}
	if p.Pembaca == nil {
		return errors.New("services: pelaksana PremiumList tanpa pembaca nomor polis")
	}
	noPolis, err := p.Pembaca.NomorPolisDariID(ctx, m.KlaimID)
	if err != nil {
		return err
	}
	saat, err := time.Parse(time.RFC3339, m.Waktu)
	if err != nil {
		return fmt.Errorf("%w: waktu muatan outbox %q", galat.ErrPermintaanTidakSah, m.Waktu)
	}
	if _, err := RakitParameterConvertJson(m.KlaimID, noPolis, saat); err != nil {
		return err
	}
	if p.Resolver == nil {
		return layanan.ErrResolverBelumDiputuskan
	}
	if _, err := layanan.AlamatLayanan(ctx, p.Resolver, KunciArasapasPremiumList); err != nil {
		return err
	}
	return outbox.ErrArasapasBelumDisetujui
}

// PelaksanaPremiumListOracle menyusun pelaksana dengan resolver dan pembaca
// Oracle. ⚠️ Belum ada pekerja yang memanggilnya (sama dengan modul lain:
// tidak ada penjadwal outbox di `cmd/api`).
func PelaksanaPremiumListOracle(svc *Service) PelaksanaPremiumList {
	return PelaksanaPremiumList{Lingkungan: svc.Lingkungan(), Resolver: layanan.ResolverLinkServiceOracle(svc),
		Pembaca: repository.NewRingkasPolisLife(svc.DB())}
}
