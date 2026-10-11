package services

// Tombol `Submit` sub-tab Achievement — `InsertToLogAchievement`
// (`Section/DetailLimits.xml`, sel 3 deret tombol Achievement; tampil bila
// `FlagExcel.CARI1=='1'`). Identik di korpus Treaty In dan Adjustment.
//
//	[1]   InputParam.CARI1 = TreatyIn.ID · CARI2 = OperatorID.pxUpdateOperator
//	[2]   untuk tiap .AchievementLists:
//	[2.1]   (Quarter == "" → lewati langkah ini) TempInput.* = baris itu
//	[2.2]   RDB InsertToLogAchievement_SQL — INSERT LOG_ACHIEVEMENT
//
// ⭐ Keputusan pemakai 8 Oktober 2026: sasarannya `LOG_ACHIEVEMENT` apa
// adanya; setiap klik menyisipkan (Pega tanpa pencegah duplikat), KECUALI
// baris ber-Quarter kosong — di Pega [2.1] terlewat sementara [2.2] tetap
// berjalan, sehingga baris itu tersisip dengan nilai baris SEBELUMNYA. Di
// sini baris itu dilewati utuh, dan cacahnya dilaporkan.
//
// ⚠️ Operator: Pega memakai `OperatorID.pxUpdateOperator` (pengubah TERAKHIR
// rekaman operator, bukan yang menekan). Di sini akun yang menekan — satu-
// satunya identitas yang aplikasi punya.
//
// ⚠️ Pega tidak menampilkan pesan apa pun sesudahnya; jawaban rute ini
// membawa cacah yang disisipkan dan dilewati untuk layar.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// MasukanLogAchievement - isi tombol Submit: pengenal `TreatyIn.ID`
// (kontrak atau penyesuaian `…/Rnn`) dan baris `AchievementLists` rincian itu.
type MasukanLogAchievement struct {
	IDKontrak string                       `json:"idKontrak"`
	Baris     []models.BarisLogAchievement `json:"baris"`
}

// HasilLogAchievement - cacah baris log yang tersisip dan yang dilewati.
type HasilLogAchievement struct {
	Disisipkan int `json:"disisipkan"`
	Dilewati   int `json:"dilewati"`
}

// CatatLogAchievement - tombol Submit sub-tab Achievement.
func (l *Layanan) CatatLogAchievement(ctx context.Context, p inti.Pelaku, m MasukanLogAchievement) (HasilLogAchievement, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilLogAchievement{}, err
	}
	id := strings.TrimSpace(m.IDKontrak)
	if id == "" {
		return HasilLogAchievement{}, ditolak("Pengenal kontrak kosong — Achievement tidak dapat dicatat.")
	}
	siap, dilewati, err := uraiLogAchievement(m.Baris)
	if err != nil {
		return HasilLogAchievement{}, err
	}
	if err := l.gudang.CatatLogAchievement(ctx, id, p.AkunID, siap); err != nil {
		return HasilLogAchievement{}, err
	}
	return HasilLogAchievement{Disisipkan: len(siap), Dilewati: dilewati}, nil
}

// uraiLogAchievement - baris layar → baris log. Murni.
//
// ⛔ Angka yang TIDAK terbaca ditolak dengan namanya, bukan dicatat 0: log
// yang diam-diam berisi nol lebih buruk daripada tombol yang menolak.
func uraiLogAchievement(baris []models.BarisLogAchievement) ([]repository.BarisLogSiap, int, error) {
	out := []repository.BarisLogSiap{}
	dilewati := 0
	for i, b := range baris {
		if strings.TrimSpace(b.Quarter) == "" {
			dilewati++
			continue
		}
		s := repository.BarisLogSiap{Quarter: b.Quarter, QuarterYear: b.QUARTERYEAR, CurrencyID: b.CurrencyID, Currency: b.Currency}
		for j, x := range []struct{ nama, v string }{
			{"Premium", b.PREMIUM}, {"R/I COMM", b.RICOMM}, {"Brokerage", b.BROKERAGE},
			{"Net Premium Before Claim", b.NETPREMIUM}, {"Paid Claim", b.PaidClaim}, {"Cash Call Claim", b.CASHCALL},
			{"Outstanding Claim", b.OutstandingClaim}, {"Incured Claim", b.IncuredClaim}, {"Total", b.Total},
			{"Net Loss Ratio", b.LossRatio},
		} {
			v := strings.TrimSpace(x.v)
			if v == "" {
				continue
			}
			d, _, err := apd.NewFromString(v)
			if err != nil {
				return nil, 0, ditolak(fmt.Sprintf("Baris Achievement ke-%d: %s %q bukan angka.", i+1, x.nama, x.v))
			}
			s.Angka[j] = d
		}
		out = append(out, s)
	}
	return out, dilewati, nil
}
