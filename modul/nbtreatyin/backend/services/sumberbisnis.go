package services

// Untuk apa berkas ini: PEMILIH SOURCE OF BUSINESS (XOL Retro) layar admin -
// tombol `Select Source Of Business` `Section/DetailPolicyTreatyIn` (pyVisible
// `.ClaimType = 'XOL Retro'`) -> showHarness `SOB` -> TreeGrid
// `Section/SourceHierarki`. Aturan murninya di `models/sumberbisnis.go`.
//
// F4 (keputusan WO 04-10-2026: IKUTI XML). Klik baris menjalankan PostDT dan
// hasilnya DIPEGANG layar; ia disimpan bersama Save/Submit layar admin
// (`terimaSumberBisnis`, dipanggil `kerjakan`), bukan saat klik.

import (
	"context"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
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

// PilihSumberBisnis = klik satu baris TreeGrid popup `SOB` - PENCARIAN TANPA
// SIMPAN:
//
//	tombol   `InputQuotation_PreAct(Acton=SOB)` langkah 1 -> `btnSOB_DT`
//	         (`Quotation.btnQuotation = "SOB"`; langkah 2 `SearchSOB.CARI1 = ""`
//	         halaman di luar pyWorkPage, hanya untuk autocomplete - tidak dibangun)
//	klik     pyRowEditing masterDetail -> flow action `AgentSourceBizDetails`
//	         -> pra-proses `SearchHierarkiSourceBizAgent_PostDT`
//
// Hanya layar admin dan hanya bila `.ClaimType = 'XOL Retro'` (pyVisible
// tombol; ClaimType dari isian layar - showHarness `pySubmitData=Yes`). Baris
// yang diklik DIBACA ULANG dari RD - nilai ID/ClientName/ChildCount layar tidak
// dipercaya.
//
// ⭐ F4: PostDT dan aktivitas `Choose` tidak ber-Obj-Save - hasilnya hanya ada
// di clipboard sampai Save/Submit layar utama. Maka di sini NOL penyimpanan:
// layar memegang `models.SumberBisnisPostDT` dan mengirimnya bersama Save/Submit, dan
// server menerimanya hanya bila cocok dengan RD yang dijalankan ulang
// (`terimaSumberBisnis`).
//
// ⛔ Tombol `Choose` (`SearchHierarkiSourceBizAgentTreatyIn_Act`) tidak
// dibangun di server: ia hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.*`
// (dibaca nol rule NB) lalu `opener.location.reload` + `window.close` - layar
// menutup popup. XML tidak menjalankan refresh berhitung apa pun sesudah klik
// atau Choose (tombol `Select Source Of Business` hanya `showHarness`, tanpa
// `refresh` - berbeda dengan `Choose Business`).
func (l *Layanan) PilihSumberBisnis(ctx context.Context, p inti.Pelaku, id, idAgen string, masuk *models.Halaman) (models.SumberBisnisPostDT, error) {
	if strings.TrimSpace(idAgen) == "" {
		return models.SumberBisnisPostDT{}, fmt.Errorf("%w: ID sumber bisnis kosong", ErrPermintaanTidakSah)
	}
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return models.SumberBisnisPostDT{}, err
	}
	if k.PositionNote != models.PosisiAdmin {
		return models.SumberBisnisPostDT{}, ErrTindakanTakAdaDiPosisi
	}
	if !models.TampilPilihSumberBisnis(h) {
		return models.SumberBisnisPostDT{}, fmt.Errorf("%w: tombol Select Source Of Business tampil hanya bila ClaimType '%s'",
			ErrTindakanTakAdaDiPosisi, models.KlaimXOLRetro)
	}
	b, ada, err := l.g.AgenHierarki(ctx, idAgen)
	if err != nil {
		return models.SumberBisnisPostDT{}, err
	}
	if !ada {
		return models.SumberBisnisPostDT{}, fmt.Errorf("%w: sumber bisnis %q tidak ada di daftar BrowseAgentHierarkiList_RD", ErrPermintaanTidakSah, idAgen)
	}
	return models.HasilPostDT(b)
}

// terimaSumberBisnis = F4: pilihan Source Of Business yang dipegang layar admin
// (hasil klik `PilihSumberBisnis`) ikut kiriman Save/Submit/refresh dan
// diterima dengan pola kiriman terkunci (`models.KirimanSumberBisnis` +
// `models.TerimaKirimanTerkunci`; RD `BrowseAgentHierarkiList_RD` DIJALANKAN
// ULANG hanya bila ada pilihan baru) - dipanggil `kerjakan` sesudah isian layar
// digabung dan SEBELUM medan turunan dihitung (`turunkan` -> `SetPPNPPH`),
// sehingga status PKP memakai nilai yang dipegang layar saat itu. Tidak cocok
// -> `models.GalatKiriman` (422, `jawabKiriman`).
func (l *Layanan) terimaSumberBisnis(ctx context.Context, h, masuk *models.Halaman) error {
	_, err := models.TerimaKirimanTerkunci(h, masuk, models.KirimanSumberBisnis(h, func() ([]models.BarisAgen, error) {
		return l.g.DaftarAgenHierarki(ctx)
	}))
	return err
}
