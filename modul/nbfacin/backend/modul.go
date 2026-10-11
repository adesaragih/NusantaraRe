// Package backend merakit modul NB FacIn (`nbfacin`) - tiket 20.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/nbfacin/backend/handlers"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
	"nusantarare/modul/nbfacin/backend/services/kontrakfacin"
)

// berkasMigrasi adalah folder `migrations/` modul ini: 180 T_NB_OPPORTUNITY dan 181
// SEQ_WORK_POLIS_NB (tiket 29, butir 76 - angka awal 181 diisi work owner/DBA, penanda
// {NB_MULAI}), 182 T_GENERAL_POLIS dan 183 T_QUOTATIONDATA sebagian (tiket 31), slot menu
// 962 (butir 59).
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "nbfacin"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`).
//
// Kontrak: MENYEDIAKAN `kontrak.PenilaiPredikatFacIn` (registry `When` NB) dan
// `kontrak.MesinPremiFacIn` (premi coverage NB) - keduanya mesin murni, tanpa data.
// `kontrak.TanggaAkseptasiFacIn` BELUM disediakan lewat perakit: tangga butuh tabel
// limit Oracle per permintaan, dan pemakai pertamanya (`rnwfacin`) dibekukan. Modul
// lain belum dapat memakainya - mengimpor `modul/nbfacin` terlarang (CLAUDE.md §5);
// sisa A32 menunggu penyedia berbasis repository.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Menyediakan: []inti.Kontrak{
			inti.KontrakDari[kontrak.PenilaiPredikatFacIn](),
			inti.KontrakDari[kontrak.MesinPremiFacIn](),
		},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			inti.Sediakan[kontrak.PenilaiPredikatFacIn](p, kontrakfacin.Predikat{})
			inti.Sediakan[kontrak.MesinPremiFacIn](p, kontrakfacin.Premi{})
			cfg := p.Config()
			// Tiket 56: SMTP korespondensi dari NBFACIN_SMTP_* (persetujuan work owner 06-10-2026); kosong = Send 503.
			// ⛔ Sandi tidak pernah dicetak (config.SMTP.String menyembunyikannya).
			m := cfg.SMTPNBFacIn
			svc := services.DariDasar(p.Dasar()).DenganSMTP(repository.KonfigSMTP{Host: m.Host, Port: m.Port, User: m.User,
				Sandi: m.Password, Dari: m.Dari, TanpaTLS: m.TanpaTLS})
			p.Catat("nbfacin: korespondensi " + m.String())
			// Tiket 63: penanda group (NBFACIN_GROUP_*, nilai tidak dicetak) dan IS_PEGA_PROD (email "Waiting for Bind").
			svc.DenganGrup(services.KonfigGrup{Login: cfg.GrupNBFacIn.Login, MarketingCode: cfg.GrupNBFacIn.MarketingCode,
				BondingMarketingName: cfg.GrupNBFacIn.BondingMarketingName}, cfg.IsPegaProd)
			p.Catat("nbfacin: group " + cfg.GrupNBFacIn.String())
			// Tiket 67: lampiran DOCUMENT_POLIS / T_STORAGE_IMAGE; ⛔ garam STORAGE_TOKEN_SALT tidak pernah dicetak.
			// Merge GitHub 10-10-2026 (keputusan work owner "Ikut GitHub, sesuaikan nbfacin"): PELAKSANA_STORAGE dibuang dari config
			// inti - penyimpanan SELALU nyata (M_LINK_SERVICE + STORAGE_TOKEN_SALT), sama dengan `inti/backend/penyimpanan`; stub
			// lokal hanya untuk uji (DenganLampiranUji).
			svc.DenganLampiran(p.Dasar(), true, cfg.StorageTokenSalt)
			return Modul{svc: svc, stubPelaku: cfg.AuthStub, pelaksana: pelaksanaNyata,
				interval: services.IntervalPekerjaLampiran, catat: p.Catat}, nil
		},
	}
}

// pelaksanaNyata - teks log pelaksana penyimpanan lampiran (dulu nilai PELAKSANA_STORAGE).
const pelaksanaNyata = "nyata"

// Modul memenuhi inti.Modul.
type Modul struct {
	svc *services.Service
	// stubPelaku - config AuthStub: penunda identitas X-Pelaku (POST opportunity, tiket 29).
	stubPelaku bool
	// pelaksana, interval, catat - pekerja latar antrean lampiran (tiket 67); interval nol = mati.
	pelaksana string
	interval  time.Duration
	catat     func(string)
}

func (Modul) Nama() string { return Nama }
func (m Modul) DaftarkanRute(mux *http.ServeMux) {
	handlers.DaftarkanRute(mux, m.svc, m.stubPelaku)
}

// JalankanPekerja - pekerja latar antrean lampiran (tiket 67, pola Treaty Contract Out): mati tanpa Oracle atau interval
// nol; berhenti bersama ctx proses, `Selesai` tertutup saat benar-benar berhenti.
func (m Modul) JalankanPekerja(ctx context.Context) inti.Pekerja {
	if m.catat == nil {
		m.catat = func(string) {}
	}
	m.catat("nbfacin lampiran: pelaksana penyimpanan " + m.pelaksana)
	selesai := make(chan struct{})
	pekerja := inti.Pekerja{Selesai: selesai,
		PesanTerlambat: "nbfacin lampiran: pekerja latar belum berhenti saat batas waktu penutupan"}
	if m.interval <= 0 || m.svc == nil || !m.svc.PunyaLampiran() {
		m.catat("nbfacin lampiran: pekerja latar mati (tanpa Oracle / interval kosong)")
		close(selesai)
		return pekerja
	}
	m.catat("nbfacin lampiran: pekerja latar setiap " + m.interval.String())
	go func() {
		defer close(selesai)
		m.svc.JalankanPekerjaLampiran(ctx, m.interval, m.catat)
	}()
	return pekerja
}
