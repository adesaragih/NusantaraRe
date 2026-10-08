// Package backend merakit modul Komite Claim Prop (`komiteclaimprop`) - komite klaim treaty inward proporsional,
// kasus `ASM-FW-GCNMFW-Work-KomiteTreaty` (`Flow/KomiteTreaty_Flow.xml`: assignment "KomiteRouter" -> Decision
// `KomiteLoop` -> Resolved-Completed). Kasusnya dilahirkan Claim Prop (opsi B 07-10-2026); modul ini membaca dan
// menulis keputusan penyetuju.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_komiteclaimprop_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/handlers"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 680-719, slot menu 986-987 - MODUL.md), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "komiteclaimprop"

// berkasKasir - kode tetap muatan Kasir (`HitServiceToKasirKMT_Act` S14.1.3-S14.1.4): konfigurasi berdokumen
// (`konfigurasi/kasir.json`, MODUL.md), bukan literal kode (`[penyimpangan sadar]` CLAUDE.md §10). Env tidak dipakai:
// env hanya dibaca `inti/backend/config` (ADR-U-0013).
//
//go:embed konfigurasi/kasir.json
var berkasKasir []byte

// KonfigurasiKasir membaca kode tetap muatan Kasir.
func KonfigurasiKasir() (models.KonfigurasiKasir, error) {
	var c struct {
		CompanyName  string `json:"companyName"`
		LjtdID       string `json:"ljtdId"`
		LdcID        string `json:"ldcId"`
		LdcIDSyariah string `json:"ldcIdSyariah"`
	}
	if err := json.Unmarshal(berkasKasir, &c); err != nil {
		return models.KonfigurasiKasir{}, fmt.Errorf("komiteclaimprop: konfigurasi/kasir.json: %w", err)
	}
	return models.KonfigurasiKasir{CompanyName: c.CompanyName, LjtdID: c.LjtdID, LdcID: c.LdcID,
		LdcIDSyariah: c.LdcIDSyariah}, nil
}

// Pendaftaran menyerahkan modul ini kepada perakit. MEMBUTUHKAN `kontrak.KlaimTreatyKomite` (disediakan `claimprop`,
// keputusan work owner 08-10-2026): klaim induk dibaca dan ditulis kembali HANYA lewat kontrak itu.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimTreatyKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			kasir, err := KonfigurasiKasir()
			if err != nil {
				return nil, err
			}
			svc := services.DariDasar(p.Dasar(), inti.Ambil[kontrak.KlaimTreatyKomite](p)).DenganKasir(kasir)
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Komite Claim Prop untuk `cmd/api` - `inti.Modul`.
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

// JalankanPekerja - pekerja outbox tidak dijalankan modul ini (pola Claim Prop / Komite Claim Life): efek keluar
// terantre di T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaKomiteClaimProp`) berhenti terang
// sampai panggilan nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
