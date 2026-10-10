// Package backend merakit modul Claim Fac In (`claimfacin`) - klaim fakultatif inward, kasus `ASM-FW-GCNMFW-Work-PNC`
// (`Flow/Register_Flow.xml`: Input Register -> Input Estimasi -> Choose Surveyor / Input Adjustment -> Resolved).
// Kasus komite KMT- (TT2 / TT3 / TT4) dilahirkan modul ini; keputusan komite = `komiteclaimfacin` (tahap 2), yang
// membaca dan menulis kembali klaim induk lewat `kontrak.KlaimFacInKomite`.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_claimfacin_gen.go`) mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 560-599, slot menu 978-979 - MODUL.md), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "claimfacin"

// Pendaftaran menyerahkan modul ini kepada perakit. MENYEDIAKAN `kontrak.KlaimFacInKomite` untuk `komiteclaimfacin`
// (prompt work owner tahap 2 10-10-2026 §2 butir 1, pola Claim Non Prop); tabel bersama (T_WORK_CLAIM, T_GENERAL_CLAIM,
// T_CLAIM_*, T_VIEW_SUGGEST, T_GENERAL_KOMITE, T_KOMITE_KOMITELIST, OS_AKSEPTASI_KLAIM, JSON_KLAIM) dibaca dan ditulis
// lewat SQL modul ini sendiri (pola Claim Prop / Non Prop, disalin bukan diimpor).
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimFacInKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			svc := services.DariDasar(p.Dasar())
			if d := p.Dasar(); d != nil && d.PunyaDatabase() { // lampiran klaim pola Claim Prop (penyimpanan bersama inti)
				svc = svc.DenganPenyimpanan(penyimpanan.Oracle(d, p.Config().StorageTokenSalt))
			}
			inti.Sediakan[kontrak.KlaimFacInKomite](p, svc.KlaimUntukKomite())
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Claim Fac In untuk `cmd/api` - `inti.Modul`.
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

// JalankanPekerja - pekerja outbox tidak dijalankan modul ini (pola Claim Prop / Non Prop): efek keluar terantre di
// T_LOG_SERVICE_RNM hanya di produksi, pelaksananya (`services.PelaksanaClaimFacIn`) berhenti terang sampai panggilan
// nyata disetujui manusia.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
