package services

// Tombol `Copy` daftar kontrak — kontrak BARU yang diisi dari kontrak tuntas.
//
// ---------------------------------------------------------------------
// Rantai ekspor (`Section/InputTreatyInOffer.xml` cell 993, label `Copy`)
// ---------------------------------------------------------------------
//
//	syarat tampil  `(OperatorID.pyWorkBasketList(2).pyWorkBasketName =
//	               'ReasTreatyInAdmin' && .Position = '' &&
//	               .StatusAkseptasi = 'Resolve Complete')` — SAMA dengan
//	               cell 994 `Revision`
//	aksi (click)   `refresh thisSection` + `SetTreatyIn_Act(ID=.ID)` — tanpa
//	               `viewstate`/`revisionstate`, jadi langkah [6]–[11]
//	               SetTreatyIn_Act (yang menyimpan) TIDAK berjalan
//	               lalu `refresh thisSection` + `TreatyInCopy`
//
// `Activity/TreatyInCopy.xml` (empat langkah, semuanya `Property-Set` atau
// `RDB-List` baca — ⛔ NOL `SaveTreatyIn`):
//
//	[1]  TreatyIn.CommentList = ""            (riwayat dikosongkan)
//	     TreatyIn.RevisionState = ""
//	[2]  RDB-List `GetCurrentDate` → halaman `time`
//	[3]  CommentList(<APPEND>).Date         = @toDateTime(time.pxResults(1).CARI1)
//	     CommentList(<LAST>).OperatorName   = OperatorID.pxInsName
//	     CommentList(<LAST>).Suggest        = "Copied from ID " + TreatyIn.ID
//	     TreatyIn.OLDID                     = TreatyIn.ID   (ID LAMA — [4] belum jalan)
//	[4]  TreatyIn.ID = "UnknownId"          (penanda kontrak BARU, sama dengan
//	                                         `TreatyInInputVis` Add=1)
//	     TreatyIn.Position = "ReasTreatyInAdmin"
//	     TreatyIn.StatusAkseptasi = ""
//	     TreatyIn.ViewState = 0
//	     TreatyIn.IsEditData = 0             (`pyMemo` "tambah iseditdata = 0")
//
// ⭐ JADI COPY TIDAK MENULIS. Ia membuka form berisi salinan clipboard
// kontrak sumber dengan `ID = "UnknownId"`; basis data baru disentuh ketika
// pemakai menekan Save (atau Submit/Decline) — dan karena ID-nya
// `UnknownId`, prosedur `SaveTreatyIn` melahirkan pengenal BARU, persis
// tombol `Add`. Padanannya di sini:
//
//	GET  /kontrak-warisan/{id}/salin  BacaDraftSalinan — draf tampil, NOL tulisan
//	POST /kontrak/salin       SimpanSalinan    — Save/Submit/Decline draf
//
// ⛔ Mengapa Save draf TIDAK lewat `/kontrak/simpan` biasa: untuk kontrak
// baru `dokumenKerja` hanya memegang kiriman layar, padahal layar hanya
// mengirim tab yang pernah DIBUKA. Clipboard Pega memegang SELURUH halaman
// sumber; tanpa dokumen sumber, tab yang tidak dibuka hilang dari salinan.
// Dan `CommentList` (komentar "Copied from ID …") milik server — kiriman
// layar untuknya diabaikan `dokumenKerja`.
//
// ⚠️ Yang TIDAK ikut tersalin, sebab tinggal di tabel berkunci ID kontrak
// dan bukan di halaman `TreatyIn`: lampiran (`M_ATTACHMENTTREATY_2`, dibaca
// `GetAllAttachment2_Sql` menurut `TreatyIn.ID` — `UnknownId` memberi nol),
// log achievement, dan polis produksi. Kurs `TREATYEXCHANGEYEARLY` berkunci
// TAHUN, bukan kontrak — tidak perlu disalin.

import (
	"context"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// SuggestSalin - `TreatyInCopy` [3], apa adanya: `"Copied from ID " + TreatyIn.ID`.
const SuggestSalin = "Copied from ID "

// MasukanSalin - draf Copy yang tombol Save/Submit/Decline kirim.
type MasukanSalin struct {
	// IDSumber - kontrak yang tombol Copy salin (`TreatyIn.OLDID` salinan).
	IDSumber string `json:"idSumber"`
	// Dokumen, Kurs - isi layar, sama dengan `MasukanSimpan`.
	Dokumen map[string]any     `json:"dokumen"`
	Kurs    *models.KursSimpan `json:"kurs"`
	// Aksi - kosong = Save; `submit` atau `decline` = tombol tab
	// `Information & Submit`. `akseptasi` ditolak: tombol Actions hanya
	// tampil bagi penyetuju, dan draf salinan selalu di posisi Admin.
	Aksi string `json:"aksi"`
}

// sumberSalinan - clipboard kontrak sumber (`SetTreatyIn_Act(ID=.ID)` tanpa
// viewstate/revisionstate), setelah syarat tampil tombol Copy ditegakkan.
func (l *Layanan) sumberSalinan(ctx context.Context, p inti.Pelaku, idSumber string) (map[string]any, string, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, "", err
	}
	id := strings.TrimSpace(idSumber)
	if id == "" {
		return nil, "", ditolak("Pilih kontrak yang akan disalin.")
	}
	// Syarat tampil cell 993 — workbasket `ReasTreatyInAdmin`.
	if !punyaPeran(p, models.PosisiAdmin) {
		return nil, "", galatTombol{ErrBukanPemegangPosisi,
			fmt.Sprintf("Tombol Copy hanya untuk workbasket %s, dan akun %s tidak memegangnya.", models.PosisiAdmin, p.AkunID)}
	}
	doc, err := l.dokumenKerja(ctx, MasukanSimpan{IDKontrak: id})
	if err != nil {
		return nil, "", err
	}
	// Syarat tampil cell 993 — `.Position = '' && .StatusAkseptasi =
	// 'Resolve Complete'`, ditegakkan juga di server.
	if teksDok(doc, "Position") != models.PosisiKosong || teksDok(doc, "StatusAkseptasi") != models.StatusTuntas {
		return nil, "", ditolak("Copy hanya untuk kontrak Resolve Complete yang tidak sedang menunggu di posisi mana pun.")
	}
	return doc, id, nil
}

// TerapkanSalin - `TreatyInCopy` [1]–[4] atas clipboard sumber.
//
// ⛔ `ID` DIHAPUS, bukan diisi `"UnknownId"`: padanan `UnknownId` di sini
// adalah `IDKontrak` kosong, yang membuat `idKontrakBaru` melahirkan
// pengenal baru — dan `simpanDalam` menimpa `doc["ID"]` dengannya.
//
// ⛔ `IsEditData` TIDAK ditaruh di dokumen: ia keadaan LAYAR (form terbuka
// dapat disunting), bukan kolom — sama dengan catatan `MulaiRevisi`.
func TerapkanSalin(doc map[string]any, idSumber, operator, stempel string) {
	// [1]
	doc["RevisionState"] = ""
	// [2]–[3] — riwayat lama dibuang; satu baris komentar salinan, tanpa
	// `IsApproved` (langkah [3] tidak menyetelnya).
	doc["CommentList"] = []any{map[string]any{
		"Date":         stempel,
		"OperatorName": operator,
		"Suggest":      SuggestSalin + idSumber,
	}}
	doc["OLDID"] = idSumber
	// [4]
	delete(doc, "ID")
	doc["Position"] = models.PosisiAdmin
	doc["StatusAkseptasi"] = ""
	doc["ViewState"] = "0"
}

// BacaDraftSalinan - tombol `Copy`: kontrak sumber, siap tampil sebagai
// kontrak BARU. ⛔ NOL tulisan — `TreatyInCopy` tidak menyimpan.
func (l *Layanan) BacaDraftSalinan(ctx context.Context, p inti.Pelaku, idSumber string) (models.KontrakWarisan, error) {
	_, id, err := l.sumberSalinan(ctx, p, idSumber)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	k, err := l.BacaKontrakWarisan(ctx, p, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	stempel := stempelPega(time.Now())
	// [4] — `UnknownId`: form memperlakukan pengenal kosong sebagai kontrak baru.
	k.ID = ""
	k.Posisi = models.PosisiAdmin
	k.StatusAkseptasi = ""
	// [1] dan [4] — `RevisionState`/`ViewState` hidup di penampung halaman.
	if k.Penampung == nil {
		k.Penampung = map[string]string{}
	}
	k.Penampung["RevisionState"] = ""
	k.Penampung["ViewState"] = "0"
	// [1]–[3] — panel History hanya memuat komentar salinan.
	k.Catatan = []models.BarisCatatanWarisan{{
		Tanggal:  TanggalTampil(stempel),
		Operator: p.AkunID,
		Catatan:  SuggestSalin + id,
	}}
	// ⚠️ Lampiran berkunci `TreatyIn.ID`; panel Attachment kontrak
	// `UnknownId` kosong. Kategorinya tetap tampil, cacahnya nol.
	k.Lampiran = []models.BarisLampiranWarisan{}
	for i := range k.KategoriLampiran {
		k.KategoriLampiran[i].Cacah = 0
	}
	// Panel `Existing Policy for Master ID` — `FetchTreatyExistingProduction`
	// (pemuat tunda layout, berjalan lagi tiap `refresh thisSection`):
	// `CARI1 := @if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)`.
	// Sesudah Copy, `ID = UnknownId` (nol baris) dan `OLDID` = sumber —
	// yang terakhir itu persis pengenal yang `BacaKontrakWarisan` sudah
	// pakai, jadi daftarnya dibiarkan.
	if strings.TrimSpace(k.EDMState) == "" {
		k.PolisProduksi = []models.BarisPolisProduksi{}
	}
	return k, nil
}

// SimpanSalinan - Save (atau Submit/Decline) draf Copy: clipboard sumber →
// `TreatyInCopy` → kiriman layar → `SaveTreatyIn_Act` dengan ID baru.
func (l *Layanan) SimpanSalinan(ctx context.Context, p inti.Pelaku, m MasukanSalin) (HasilSimpan, error) {
	switch m.Aksi {
	case "", AksiSubmit, AksiDecline:
	default:
		return HasilSimpan{}, ditolak(fmt.Sprintf("Aksi %q tidak berlaku untuk salinan yang belum disimpan.", m.Aksi))
	}
	doc, id, err := l.sumberSalinan(ctx, p, m.IDSumber)
	if err != nil {
		return HasilSimpan{}, err
	}
	TerapkanSalin(doc, id, p.AkunID, stempelPega(time.Now()))
	// Isian pemakai di atas clipboard salinan — kunci milik server dan
	// halaman sesi diabaikan, aturan yang SAMA dengan `dokumenKerja`.
	for k, v := range m.Dokumen {
		if kunciMilikServer[k] || kunciSesi(k) {
			continue
		}
		doc[k] = v
	}
	// ⛔ `OLDID` milik tombol Copy, bukan isian layar.
	doc["OLDID"] = id
	baru := MasukanSimpan{Dokumen: m.Dokumen, Kurs: m.Kurs}
	if m.Aksi == AksiSubmit {
		// `TreatyInCheckError` DIPERIKSA SEBELUM salinan ditulis — Submit
		// yang ditolak di Pega tidak menyimpan apa pun.
		if msg := periksaGalatSubmit(doc); msg != "" {
			return HasilSimpan{}, ditolak(msg)
		}
	}
	// Save draf: `Position` sudah `ReasTreatyInAdmin` ([4]), jadi DT
	// `TreatyInAddNew` [1] tidak mengubah apa pun; status kosong, jadi
	// penolakan `Resolve Complete` milik Save tidak berlaku.
	h, err := l.tulis(ctx, p, baru, doc)
	if err != nil || m.Aksi == "" {
		return h, err
	}
	// Submit/Decline draf: Pega menjalankan tangganya atas clipboard yang
	// SAMA lalu satu `SaveTreatyIn`. Di sini kontrak barunya sudah lahir,
	// dan tombol yang sama dijalankan atas pengenal itu — dokumennya kini
	// salinan tersimpan, jadi hasil akhirnya sama.
	//
	// ⚠️ DUA transaksi, bukan satu: bila langkah kedua gagal, salinan
	// tetap tersimpan sebagai draf Admin berpengenal `h.ID`.
	k, err := l.KirimKontrak(ctx, p, MasukanKirim{
		MasukanSimpan: MasukanSimpan{IDKontrak: h.ID, Dokumen: m.Dokumen},
		Aksi:          m.Aksi,
	})
	if err != nil {
		return HasilSimpan{}, fmt.Errorf("salinan tersimpan dengan ID %s, tetapi %s gagal: %w", h.ID, m.Aksi, err)
	}
	return k, nil
}
