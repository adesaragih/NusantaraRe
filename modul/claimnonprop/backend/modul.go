// Package backend merakit modul Claim Non Prop (`claimnonprop`) - klaim treaty inward non proporsional / XoL, kasus
// `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` (`Flow/Flow_TreatyIn.xml`: Outstanding Claim -> Input Acceptation ->
// Resolved). Tahap 1 dari 2: kasus komite KMTNP- dilahirkan modul ini; keputusan komite = `komiteclaimnonprop` (tahap 2).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_claimnonprop_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimnonprop/backend/handlers"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 600-639, slot menu 982-983 - MODUL.md), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "claimnonprop"

// berkasKasir - kode tetap muatan Kasir (`HitServiceToKasir_Act` 9.3-9.4): konfigurasi berdokumen
// (`konfigurasi/kasir.json`, MODUL.md), bukan literal kode. Env tidak dipakai: env hanya dibaca `inti/backend/config`
// (ADR-U-0013). Kata sandi Edit XOL Allocation TIDAK di sini (OQ-CNP-04: tombol nonaktif).
//
//go:embed konfigurasi/kasir.json
var berkasKasir []byte

// KonfigurasiKasir membaca kode tetap muatan Kasir.
func KonfigurasiKasir() (models.KonfigKasir, error) {
	var c struct {
		CompanyName  string `json:"companyName"`
		LjtdID       string `json:"ljtdId"`
		LdcID        string `json:"ldcId"`
		LdcIDSyariah string `json:"ldcIdSyariah"`
		StsAp        string `json:"stsAp"`
	}
	if err := json.Unmarshal(berkasKasir, &c); err != nil {
		return models.KonfigKasir{}, fmt.Errorf("claimnonprop: konfigurasi/kasir.json: %w", err)
	}
	return models.KonfigKasir{CompanyName: c.CompanyName, LjtdID: c.LjtdID, LdcID: c.LdcID,
		LdcIDSyariah: c.LdcIDSyariah, StsAp: c.StsAp}, nil
}

// Pendaftaran menyerahkan modul ini kepada perakit. Nol kontrak disediakan / dipakai: tabel bersama (T_WORK_CLAIM,
// T_GENERAL_CLAIM, T_CLAIM_*, T_VIEW_SUGGEST, T_GENERAL_KOMITE, T_KOMITE_KOMITELIST) dibaca dan ditulis lewat SQL modul
// ini sendiri; kontrak klaim <-> komite Non Prop dibuat di tahap 2 (prompt §2).
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			kasir, err := KonfigurasiKasir()
			if err != nil {
				return nil, err
			}
			return Baru(services.DariDasar(p.Dasar()).DenganKasir(kasir), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Claim Non Prop untuk `cmd/api` - `inti.Modul`.
type Modul struct {
	svc        *services.Layanan
	stubPelaku bool
}

// Baru merakit modul di atas layanan yang sudah disambung.
func Baru(svc *services.Layanan, stubPelaku bool) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - pekerja outbox tidak dijalankan modul ini (pola Claim Prop): efek keluar terantre di
// T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaClaimNonProp`) berhenti terang sampai panggilan
// nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
