package services

// Efek kasus dan laporan harian Komite - tiket 08. TANPA Oracle.

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimlife/backend/models"
	"nusantarare/modul/komiteclaimlife/backend/repository"
)

// TestKodeEfekSamaDiTigaLapis - km5: kode di models, repository, dan outbox sama.
func TestKodeEfekSamaDiTigaLapis(t *testing.T) {
	for m, r := range map[string]string{
		models.KodeEfekAntre: outbox.StatusEfekAntre, models.KodeEfekJalan: outbox.StatusEfekJalan,
		models.KodeEfekSelesai: outbox.StatusEfekSelesai, models.KodeEfekGagalPermanen: outbox.StatusEfekGagalPermanen,
	} {
		if m != r {
			t.Errorf("kode efek models %q ≠ repository %q", m, r)
		}
	}
	if repository.ModulOutboxKomite != ModulKomiteLife {
		t.Errorf("modul outbox %q ≠ %q", repository.ModulOutboxKomite, ModulKomiteLife)
	}
}

// TestLaporanHarianKosongDinyatakan - AC "tetap terkirim meskipun kosong".
func TestLaporanHarianKosongDinyatakan(t *testing.T) {
	saat := time.Date(2026, 9, 28, 6, 0, 0, 0, time.UTC)
	l := susunLaporan(saat, nil)
	if !l.Kosong || l.Baris == nil || l.Tanggal != "2026-09-28" {
		t.Errorf("laporan kosong %+v", l)
	}
	l = susunLaporan(saat, []repository.EfekKasusKomite{{KasusID: "KMTLF-UJI", Jenis: JenisEfekKomiteKasir,
		Status: "gagal-permanen", Percobaan: 8, Dibuat: saat.Add(-time.Hour),
		Diperbarui: sql.NullTime{Time: saat, Valid: true}}})
	if l.Kosong || len(l.Baris) != 1 || l.Baris[0].Keadaan != models.KataEfekPerluIntervensi ||
		l.Baris[0].Sejak != "2026-09-28 06:00:00" {
		t.Errorf("laporan berisi %+v", l)
	}
}

// TestLaporanHarianHanyaAdmin - ADR-0014 [asumsi OQ-007/021].
func TestLaporanHarianHanyaAdmin(t *testing.T) {
	i := New(nil).InboxKomite()
	if _, err := i.LaporanHarian(context.Background(), inti.Pelaku{AkunID: "UJI"}, time.Now()); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("bukan admin: %v", err)
	}
}

// TestEfekMenempelPadaKasus - ringkasan efek di layar kasus.
func TestEfekMenempelPadaKasus(t *testing.T) {
	r := ringkasEfekKasus([]repository.EfekKasusKomite{
		{Jenis: JenisEfekKomiteArasapas, Status: "selesai", Dibuat: time.Now()},
		{Jenis: JenisEfekKomiteKasir, Status: "antre", Dibuat: time.Now()},
	})
	if r.Keadaan != models.KataKeputusanTersimpan || len(r.Efek) != 2 || r.Efek[1].Keadaan != models.KataEfekTertunda {
		t.Errorf("ringkasan efek %+v", r)
	}
}
