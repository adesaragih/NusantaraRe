// Package backend merakit modul Komite Claim Fac In (`komiteclaimfacin`) - komite klaim fakultatif inward, kasus
// `ASM-FW-GCNMFW-Work-Komite` (`Flow/Komite_Flow.xml` korpus `Komite Claim FacIn`: assignment "KomiteRouter" ->
// `ViewTransferDtl` (`SetValueKomite` / `ShowTransfer` / `KomitePostAct`) -> decision IsKomiteLoop). Kasusnya
// dilahirkan Claim Fac In (`CreateKMTNo_Act` TT2, `SendRejectClaimToKomite2` TT3, `SendCloseClaimToKomite` TT4); modul
// ini membaca dan menulis keputusan penyetuju. Pola Komite Claim Prop / Non Prop (prompt work owner tahap 2
// 10-10-2026): tanpa menu sendiri, dibuka dari inbox Claim Fac In.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_komiteclaimfacin_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/komiteclaimfacin/backend/handlers"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/services"
)

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "komiteclaimfacin"

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 640-679, slot menu `—` - MODUL.md), ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// berkasKasir - kode tetap muatan Kasir (`HitServiceToKasirKMT_Act` S14.2.2.1.1.3-S14.2.2.1.1.4): konfigurasi
// berdokumen (`konfigurasi/kasir.json`, MODUL.md), bukan literal kode (`[penyimpangan sadar]` CLAUDE.md §10). Env tidak
// dipakai: env hanya dibaca `inti/backend/config` (ADR-U-0013).
//
//go:embed konfigurasi/kasir.json
var berkasKasir []byte

// berkasEmail - akun notifikasi dan CC email komite (`SendEmailKlaim_KMT`): konfigurasi berdokumen
// (`konfigurasi/email.json`, MODUL.md). BCC pribadi tidak disalin.
//
//go:embed konfigurasi/email.json
var berkasEmail []byte

// KonfigurasiEmail membaca akun notifikasi dan CC email komite.
func KonfigurasiEmail() (models.KonfigurasiEmail, error) {
	var c struct {
		Akun        string `json:"akun"`
		AkunSyariah string `json:"akunSyariah"`
		CC          string `json:"cc"`
	}
	if err := json.Unmarshal(berkasEmail, &c); err != nil {
		return models.KonfigurasiEmail{}, fmt.Errorf("komiteclaimfacin: konfigurasi/email.json: %w", err)
	}
	return models.KonfigurasiEmail{Akun: c.Akun, AkunSyariah: c.AkunSyariah, CC: c.CC}, nil
}

// KonfigurasiKasir membaca kode tetap muatan Kasir.
func KonfigurasiKasir() (models.KonfigurasiKasir, error) {
	var c struct {
		CompanyName  string `json:"companyName"`
		LjtdID       string `json:"ljtdId"`
		LdcID        string `json:"ldcId"`
		LdcIDSyariah string `json:"ldcIdSyariah"`
	}
	if err := json.Unmarshal(berkasKasir, &c); err != nil {
		return models.KonfigurasiKasir{}, fmt.Errorf("komiteclaimfacin: konfigurasi/kasir.json: %w", err)
	}
	return models.KonfigurasiKasir{CompanyName: c.CompanyName, LjtdID: c.LjtdID, LdcID: c.LdcID,
		LdcIDSyariah: c.LdcIDSyariah}, nil
}

// Pendaftaran menyerahkan modul ini kepada perakit. MEMBUTUHKAN `kontrak.KlaimFacInKomite` (disediakan `claimfacin`,
// prompt tahap 2 §2 butir 1): klaim induk dibaca dan ditulis kembali HANYA lewat kontrak itu.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimFacInKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			kasir, err := KonfigurasiKasir()
			if err != nil {
				return nil, err
			}
			surel, err := KonfigurasiEmail()
			if err != nil {
				return nil, err
			}
			svc := services.DariDasar(p.Dasar(), inti.Ambil[kontrak.KlaimFacInKomite](p)).DenganKasir(kasir).
				DenganEmail(surel)
			if d := p.Dasar(); d != nil && d.PunyaDatabase() { // PrintPDFAccep_MultiAksep_KMT S27 (pola Komite Claim Prop)
				svc = svc.DenganPenyimpanan(penyimpanan.Oracle(d, p.Config().StorageTokenSalt))
			}
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Komite Claim Fac In untuk `cmd/api` - `inti.Modul`.
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

// JalankanPekerja - pekerja outbox tidak dijalankan modul ini (pola Komite Claim Prop): efek keluar terantre di
// T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaKomiteClaimFacIn`) berhenti terang sampai
// panggilan nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
