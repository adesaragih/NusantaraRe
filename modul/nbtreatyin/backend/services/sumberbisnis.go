package services

// Untuk apa berkas ini: PEMILIH SOURCE OF BUSINESS (XOL Retro) layar admin -
// tombol `Select Source Of Business` `Section/DetailPolicyTreatyIn` (pyVisible
// `.ClaimType = 'XOL Retro'`) -> showHarness `SOB` -> TreeGrid
// `Section/SourceHierarki`. Aturan murninya di `models/sumberbisnis.go`.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// DaftarSumberBisnis - isi TreeGrid `Section/SourceHierarki`
// (`AgentSourceBizTreatyIn_Act` -> RD `BrowseAgentHierarkiList_RD`, keadaan NB:
// daftar akar - lihat `repository/agen.go`).
func (l *Layanan) DaftarSumberBisnis(ctx context.Context, p inti.Pelaku) ([]models.BarisAgen, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	return l.g.DaftarAgenHierarki(ctx)
}

// PilihSumberBisnis = klik satu baris TreeGrid popup `SOB`:
//
//	tombol   `InputQuotation_PreAct(Acton=SOB)` langkah 1 -> `btnSOB_DT`
//	         (`Quotation.btnQuotation = "SOB"`; langkah 2 `SearchSOB.CARI1 = ""`
//	         halaman di luar pyWorkPage, hanya untuk autocomplete - tidak dibangun)
//	klik     pyRowEditing masterDetail -> flow action `AgentSourceBizDetails`
//	         -> pra-proses `SearchHierarkiSourceBizAgent_PostDT`
//
// Hanya layar admin dan hanya bila `.ClaimType = 'XOL Retro'` (pyVisible
// tombol). Isian layar ikut dikirim (showHarness `pySubmitData=Yes`) dan
// digabung menurut daftar izin admin (`kerjakan`). Baris yang diklik DIBACA
// ULANG dari RD - nilai ID/ClientName/ChildCount layar tidak dipercaya.
//
// ⚠️ `[penyesuaian sadar]` Pega memegang hasil PostDT di clipboard sampai Save/
// Submit layar utama (PostDT dan aktivitas Choose tanpa Obj-Save). Layar baru
// tidak dapat mengirim `Quotation.*` (medan milik server, `models.medanAdmin`),
// jadi hasilnya DISIMPAN saat klik - di bawah kunci kasus, satu transaksi,
// pola yang sama dengan `PilihBisnis`.
//
// ⛔ Tombol `Choose` (`SearchHierarkiSourceBizAgentTreatyIn_Act`) tidak
// dibangun di server: ia hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.*`
// (dibaca nol rule NB) lalu `window.close` + `opener.location.reload` - layar
// menutup popup.
func (l *Layanan) PilihSumberBisnis(ctx context.Context, p inti.Pelaku, id, idAgen string, masuk *models.Halaman) (Layar, error) {
	if strings.TrimSpace(idAgen) == "" {
		return Layar{}, fmt.Errorf("%w: ID sumber bisnis kosong", ErrPermintaanTidakSah)
	}
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return Layar{}, err
	}
	if k.PositionNote != models.PosisiAdmin {
		return Layar{}, ErrTindakanTakAdaDiPosisi
	}
	if !models.TampilPilihSumberBisnis(h) {
		return Layar{}, fmt.Errorf("%w: tombol Select Source Of Business tampil hanya bila ClaimType '%s'",
			ErrTindakanTakAdaDiPosisi, models.KlaimXOLRetro)
	}
	b, ada, err := l.g.AgenHierarki(ctx, idAgen)
	if err != nil {
		return Layar{}, err
	}
	if !ada {
		return Layar{}, fmt.Errorf("%w: sumber bisnis %q tidak ada di daftar BrowseAgentHierarkiList_RD", ErrPermintaanTidakSah, idAgen)
	}
	models.TombolSOB(h)
	if err := models.TerapkanSumberBisnis(h, b); err != nil {
		return Layar{}, err
	}
	models.SalinKeQuotationData(h, models.MedanSumberBisnis...)
	if err := l.tulis(ctx, k, func(tx *db.Tx) error { return l.g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		return Layar{}, err
	}
	return l.layar(ctx, p, k, h, true)
}
