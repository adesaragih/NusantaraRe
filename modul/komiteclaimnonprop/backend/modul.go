// Package backend merakit modul Komite Claim Non Prop (`komiteclaimnonprop`) - komite klaim treaty inward non
// proporsional (XoL), kasus `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp` (`Flow/KomiteTreaty_Flow.xml` korpus `Komite Claim
// Non Prop`: assignment "KomiteRouter" -> `KomitePostAdjustment` -> Resolved-Completed). Kasusnya dilahirkan Claim Non
// Prop (`CreateChildKomiteCNP_Act`); modul ini membaca dan menulis keputusan penyetuju. Pola disalin dari Komite Claim
// Prop (perintah work owner 09-10-2026: tanpa menu sendiri, dibuka dari inbox Claim Non Prop).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_komiteclaimnonprop_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimnonprop/backend/handlers"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
	"nusantarare/modul/komiteclaimnonprop/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 720-759, slot menu 988-989 - MODUL.md), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "komiteclaimnonprop"

// berkasKasir - kode tetap muatan Kasir (`HitServiceToKasirKMT_Act` S10.3-S10.4): konfigurasi berdokumen
// (`konfigurasi/kasir.json`, MODUL.md), bukan literal kode (`[penyimpangan sadar]` CLAUDE.md §10). Env tidak dipakai:
// env hanya dibaca `inti/backend/config` (ADR-U-0013).
//
//go:embed konfigurasi/kasir.json
var berkasKasir []byte

// berkasEmail - akun notifikasi dan CC email komite (`SendEmailKlaim_KMT`): konfigurasi berdokumen
// (`konfigurasi/email.json`, MODUL.md). BCC pribadi S3 tidak disalin.
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
		return models.KonfigurasiEmail{}, fmt.Errorf("komiteclaimnonprop: konfigurasi/email.json: %w", err)
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
		return models.KonfigurasiKasir{}, fmt.Errorf("komiteclaimnonprop: konfigurasi/kasir.json: %w", err)
	}
	return models.KonfigurasiKasir{CompanyName: c.CompanyName, LjtdID: c.LjtdID, LdcID: c.LdcID,
		LdcIDSyariah: c.LdcIDSyariah}, nil
}

// Pendaftaran menyerahkan modul ini kepada perakit. MEMBUTUHKAN `kontrak.KlaimTreatyNonPropKomite` (disediakan
// `claimnonprop`, perintah work owner 09-10-2026): klaim induk dibaca dan ditulis kembali HANYA lewat kontrak itu.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimTreatyNonPropKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			kasir, err := KonfigurasiKasir()
			if err != nil {
				return nil, err
			}
			surel, err := KonfigurasiEmail()
			if err != nil {
				return nil, err
			}
			svc := services.DariDasar(p.Dasar(), inti.Ambil[kontrak.KlaimTreatyNonPropKomite](p)).DenganKasir(kasir).
				DenganEmail(surel)
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Komite Claim Non Prop untuk `cmd/api` - `inti.Modul`.
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
// T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaKomiteClaimNonProp`) berhenti terang
// sampai panggilan nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
