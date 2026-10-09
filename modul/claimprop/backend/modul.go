// Package backend merakit modul Claim Prop (`claimprop`) - klaim treaty inward proporsional, kasus
// `ASM-FW-GCNMFW-Work-ClaimTreaty` (`Flow/Flow_TreatyIn.xml`: Outstanding Claim -> Input Acceptation -> Resolved).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_claimprop_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 520-559, slot menu 980-981 - MODUL.md), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "claimprop"

// Pendaftaran menyerahkan modul ini kepada perakit. MENYEDIAKAN `kontrak.KlaimTreatyKomite` (keputusan work owner
// 08-10-2026): Komite Claim Prop membaca kasus klaim induk dan menulis kembali hasil keputusannya HANYA lewat kontrak
// itu (`services.KlaimUntukKomite`). Nol kontrak dipakai: tabel bersama (T_WORK_CLAIM, T_GENERAL_CLAIM,
// T_VIEW_SUGGEST, T_GENERAL_KOMITE, T_KOMITE_KOMITELIST) dibaca dan ditulis lewat SQL modul ini sendiri.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimTreatyKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			svc := services.DariDasar(p.Dasar())
			if d := p.Dasar(); d != nil && d.PunyaDatabase() { // GCNMSaveAttachments -> InsertDocument_Act (pola Bordereaux)
				svc = svc.DenganPenyimpanan(penyimpanan.Oracle(d, p.Config().StorageTokenSalt))
			}
			inti.Sediakan[kontrak.KlaimTreatyKomite](p, svc.KlaimUntukKomite())
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Claim Prop untuk `cmd/api` - `inti.Modul`.
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

// JalankanPekerja - pekerja outbox tidak dijalankan modul ini (pola Komite Claim Life): efek keluar terantre di
// T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaClaimProp`) berhenti terang sampai panggilan
// nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
