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

// HasilSumberBisnis - jawaban klik satu baris: nilai yang ditulis PostDT per
// jalur halaman (`Quotation.SourceOfBusiness`, `.SobName`, `.SobLeader0`,
// `.SobLeader1`). Layar memegangnya dan mengirimnya kembali pada Save/Submit.
type HasilSumberBisnis struct {
	Nilai map[string]string `json:"nilai"`
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
// layar memegang `HasilSumberBisnis` dan mengirimnya bersama Save/Submit, dan
// server menerimanya hanya bila cocok dengan RD yang dijalankan ulang
// (`terimaSumberBisnis`).
//
// ⛔ Tombol `Choose` (`SearchHierarkiSourceBizAgentTreatyIn_Act`) tidak
// dibangun di server: ia hanya menulis `pyWorkPage.OfferTreatyIn.QuotationData.*`
// (dibaca nol rule NB) lalu `opener.location.reload` + `window.close` - layar
// menutup popup. XML tidak menjalankan refresh berhitung apa pun sesudah klik
// atau Choose (tombol `Select Source Of Business` hanya `showHarness`, tanpa
// `refresh` - berbeda dengan `Choose Business`).
func (l *Layanan) PilihSumberBisnis(ctx context.Context, p inti.Pelaku, id, idAgen string, masuk *models.Halaman) (HasilSumberBisnis, error) {
	if strings.TrimSpace(idAgen) == "" {
		return HasilSumberBisnis{}, fmt.Errorf("%w: ID sumber bisnis kosong", ErrPermintaanTidakSah)
	}
	k, h, err := l.kerjakan(ctx, p, id, masuk)
	if err != nil {
		return HasilSumberBisnis{}, err
	}
	if k.PositionNote != models.PosisiAdmin {
		return HasilSumberBisnis{}, ErrTindakanTakAdaDiPosisi
	}
	if !models.TampilPilihSumberBisnis(h) {
		return HasilSumberBisnis{}, fmt.Errorf("%w: tombol Select Source Of Business tampil hanya bila ClaimType '%s'",
			ErrTindakanTakAdaDiPosisi, models.KlaimXOLRetro)
	}
	b, ada, err := l.g.AgenHierarki(ctx, idAgen)
	if err != nil {
		return HasilSumberBisnis{}, err
	}
	if !ada {
		return HasilSumberBisnis{}, fmt.Errorf("%w: sumber bisnis %q tidak ada di daftar BrowseAgentHierarkiList_RD", ErrPermintaanTidakSah, idAgen)
	}
	nilai, err := models.HasilPostDT(b)
	if err != nil {
		return HasilSumberBisnis{}, err
	}
	return HasilSumberBisnis{Nilai: nilai}, nil
}

// terimaSumberBisnis = F4: pilihan Source Of Business yang dipegang layar admin
// (hasil klik `PilihSumberBisnis`) ikut kiriman Save/Submit/refresh dan
// diterima di sini - dipanggil `kerjakan` sesudah isian layar digabung dan
// SEBELUM medan turunan dihitung (`turunkan` -> `SetPPNPPH`), sehingga status
// PKP memakai nilai yang dipegang layar saat itu.
//
//   - tidak ada pilihan baru (`Quotation.SourceOfBusiness` tidak dikirim atau
//     sama dengan tersimpan) -> tidak berbuat apa pun;
//   - ClaimType (isian layar) bukan 'XOL Retro' -> tombolnya tidak tampil, medan
//     itu TERKUNCI: kiriman diabaikan, nilai tersimpan dipakai (pola AC 49-51
//     `GabungMasukanLayar`);
//   - selain itu RD `BrowseAgentHierarkiList_RD` DIJALANKAN ULANG dan pilihan
//     diterima hanya bila sama persis dengan hasil PostDT salah satu barisnya
//     (`models.CocokHasilPostDT`); tidak cocok -> 422 berpesan jelas.
func (l *Layanan) terimaSumberBisnis(ctx context.Context, h, masuk *models.Halaman) error {
	pilihan, berubah := models.SumberBisnisKiriman(h, masuk)
	if !berubah || !models.TampilPilihSumberBisnis(h) {
		return nil
	}
	daftar, err := l.g.DaftarAgenHierarki(ctx)
	if err != nil {
		return err
	}
	if !models.CocokHasilPostDT(daftar, pilihan) {
		id := pilihan[models.HalamanQuotation+".SourceOfBusiness"]
		return &GalatValidasi{Pesan: []string{models.PesanSumberBisnisTidakCocok(id)}}
	}
	models.TerapkanPilihanSumberBisnis(h, pilihan)
	return nil
}
