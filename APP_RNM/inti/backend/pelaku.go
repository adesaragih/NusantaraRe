package backend

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrTanpaWewenang dikembalikan bila peran tidak mencukupi.
	ErrTanpaWewenang = errors.New("services: wewenang tidak mencukupi")
)

// Pelaku adalah identitas yang meminta sebuah tindakan.
//
// ⛔ Nol nama orang di kode (ADR-U-0030, CLAUDE.md bab 4 butir 10). Yang
// dibawa hanya pengenal akun dan peran; nama tidak pernah menjadi dasar
// keputusan dan tidak pernah ditulis ke kode atau fixture.
type Pelaku struct {
	AkunID string
	Peran  []string
}

// PunyaPeran memeriksa satu peran.
//
// Sumber peran adalah SATU tabel (ADR-U-0030). Tabel itu belum ada di Fase 0;
// pembacanya lahir bersama tiket yang memerlukannya, dan pemeriksaan tetap
// dilakukan di lapisan ini - tidak pernah di handlers, tidak pernah di SQL.
func (p Pelaku) PunyaPeran(peran string) bool {
	for _, x := range p.Peran {
		if x == peran {
			return true
		}
	}
	return false
}

// ErrTanpaIdentitas menandai permintaan tanpa pengenal akun.
//
// ⛔ DIPISAH dari ErrTanpaWewenang 26-09-2026. Satu galat yang berarti dua hal
// - "aku tidak tahu kamu siapa" dan "aku tahu kamu siapa, tetapi kamu tidak
// boleh" - memaksa pemanggil menebak, dan dua handler memang menerjemahkannya
// ke dua kode HTTP yang berbeda untuk galat yang sama. 401 dan 403 menjawab
// pertanyaan yang berbeda.
var ErrTanpaIdentitas = errors.New("services: permintaan tanpa identitas pelaku")

// WajibIdentitas menolak permintaan tanpa pengenal akun.
func WajibIdentitas(p Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaIdentitas
	}
	return nil
}

// WajibPeran mengembalikan galat bila peran tidak dimiliki.
//
// Jalur yang DITOLAK wajib punya uji tersendiri, bukan hanya jalur yang
// berhasil (brief bab 5).
func WajibPeran(p Pelaku, peran string) error {
	if !p.PunyaPeran(peran) {
		return fmt.Errorf("%w: perlu peran %q", ErrTanpaWewenang, peran)
	}
	return nil
}

// Peran yang muncul di gerbang XML dan di Flow.
//
// `[terverifikasi]` `Flow/Register_Flow.xml`: `Assignment1` (Outstanding
// Claim) dan `Assignment2` (Input Register) dipegang `ReasLifeAdmin`;
// `Assignment3` (Medical Check) dan `Decision3` dipegang
// `ReasLifeMedicalAdvisor`; `Decision1` dipegang `ReasLifeSPV`.
//
// ⛔ Teks perannya hidup HANYA di sini. Tiket 03 dan 05 menamai IZIN-nya
// lebih dulu (`PeranSimpanOutstanding`, `PeranRejectOutstanding`), dan
// keduanya kini bernilai dari konstanta di bawah - bukan menaruh teks
// `"ReasLifeAdmin"` untuk kedua kalinya.
const (
	PeranAdmin          = "ReasLifeAdmin"
	PeranSPV            = "ReasLifeSPV"
	PeranMedicalAdvisor = "ReasLifeMedicalAdvisor"
)
