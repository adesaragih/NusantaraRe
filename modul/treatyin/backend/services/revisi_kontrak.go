package services

// Tombol `Revision` daftar kontrak — membuka kontrak tuntas untuk direvisi.
//
// ---------------------------------------------------------------------
// Rantai ekspor (`Section/InputTreatyInOffer.xml` cell 994)
// ---------------------------------------------------------------------
//
//	syarat tampil  WB ReasTreatyInAdmin && .Position = ''
//	               && .StatusAkseptasi = 'Resolve Complete'
//	aksi           SetTreatyIn_Act(ID, viewstate='1', revisionstate='1')
//
// `SetTreatyIn_Act`:
//
//	[6]  ViewState = 1, IsEditData = 1
//	[7]  RevisionState = 1, StatusAkseptasi = "", PositionUsername = operator
//	     (⛔ Position TIDAK diubah — tetap "")
//	[8]  GetCurrentDate
//	[9]  CommentList(<APPEND>): Date, OperatorName, Suggest "Create Revision"
//	     (tanpa IsApproved)
//	[10–11] SaveTreatyIn — ditulis SEKETIKA, bukan menunggu Save
//
// Sesudahnya tangga memakai cabang revisi `Akseptasi_DT` (`LangkahBerikut`
// dengan `revisi`), yang tuntas di Sec Head.
//
// ⛔ `IsEditData` keadaan LAYAR (form terkunci), bukan kolom: form dibuka
// dalam mode lihat, dan `RevisionState` yang menghidupkan Comment serta
// Submit revisi (`TreatyInfoSubmit` cell 9 dan 21).

import (
	"context"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// SuggestBuatRevisi - `SetTreatyIn_Act` [9], apa adanya.
const SuggestBuatRevisi = "Create Revision"

// MasukanRevisi - kontrak yang tombol Revision buka.
type MasukanRevisi struct {
	IDKontrak string `json:"idKontrak"`
}

// MulaiRevisi - tombol `Revision`.
func (l *Layanan) MulaiRevisi(ctx context.Context, p inti.Pelaku, m MasukanRevisi) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	id := strings.TrimSpace(m.IDKontrak)
	if id == "" {
		return HasilSimpan{}, ditolak("Pilih kontrak yang akan direvisi.")
	}
	if !punyaPeran(p, models.PosisiAdmin) {
		return HasilSimpan{}, galatTombol{ErrBukanPemegangPosisi,
			fmt.Sprintf("Tombol Revision hanya untuk workbasket %s, dan akun %s tidak memegangnya.", models.PosisiAdmin, p.AkunID)}
	}
	// ⛔ Tanpa kolom `REVISIONSTATE`/`VIEWSTATE` (migrasi `448`) keadaan revisi
	// HILANG begitu disimpan, dan tangganya jatuh ke cabang biasa — lebih
	// baik ditolak dengan alasannya daripada berhasil setengah.
	belum, err := l.gudang.KunciBelumTerpasang(ctx, map[string]any{"RevisionState": "1", "ViewState": "1"})
	if err != nil {
		return HasilSimpan{}, err
	}
	if len(belum) > 0 {
		return HasilSimpan{}, ditolak("Kolom " + strings.Join(belum, ", ") +
			" belum terpasang di T_TREATY_REVISION (migrasi 448) — revisi belum dapat disimpan.")
	}
	doc, err := l.dokumenKerja(ctx, MasukanSimpan{IDKontrak: id})
	if err != nil {
		return HasilSimpan{}, err
	}
	// Syarat tampil tombolnya — ditegakkan juga di server.
	if teksDok(doc, "Position") != models.PosisiKosong || teksDok(doc, "StatusAkseptasi") != models.StatusTuntas {
		return HasilSimpan{}, ditolak("Revision hanya untuk kontrak Resolve Complete yang tidak sedang menunggu di posisi mana pun.")
	}
	operator := p.AkunID
	// [6]–[7]
	doc["ViewState"] = "1"
	doc["RevisionState"] = "1"
	doc["StatusAkseptasi"] = ""
	doc["PositionUsername"] = operator
	// [8]–[9] — tanpa `IsApproved`.
	l2, _ := doc["CommentList"].([]any)
	doc["CommentList"] = append(l2, map[string]any{
		"Date":         stempelPega(time.Now()),
		"OperatorName": operator,
		"Suggest":      SuggestBuatRevisi,
	})
	// [10–11] — `tulis` menolak memasang Position sendiri; Save biasa yang
	// memasang Admin (`TreatyInAddNew`), tombol ini tidak.
	return l.tulis(ctx, p, MasukanSimpan{IDKontrak: id}, doc)
}
