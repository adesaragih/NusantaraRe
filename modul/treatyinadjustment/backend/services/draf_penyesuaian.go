package services

// Tombol `Choose` picker Add Revision / Add Adjustment Premium —
// `TreatyInEDMSetValue` langkah [1]–[6], DI LAYAR.
//
// ---------------------------------------------------------------------
// ⛔ DRAF, BUKAN SIMPANAN — keputusan pemilik proses 7 Oktober 2026
// ---------------------------------------------------------------------
// Di Pega `Choose` langsung menyimpan: [8] `TreatyRevisionCopyAttachment`
// dan [9] `SaveTreatyIn_EDM_Act`. Aturan aplikasi ini: isian tidak masuk
// basis data sebelum Save/Submit. Jadi langkah [1]–[6] dijalankan di sini
// TANPA menulis apa pun, dan [8]–[9] menunggu jalur Save. [7]
// `TreatyInEdmCheckDuplicate` ber-`//` — mati di ekspor.
//
// Urutan langkahnya (semua dari ekspor):
//
//	[1]/[2] SetTreatyInEDM_Act / SetTreatyIn_Act — muat dokumen `Param.ID`
//	        (panjang 7 = master Treaty In, selainnya revisi)
//	[3]     Property-Set: OLDID=ID · CommentList="" · EDMState=InternalType ·
//	        EDMMaterialType=MaterialType · RevisionDate=@CurrentDateTime()
//	[4]     DT TreatyInSetEditPre — satu baris CommentList per EDMState;
//	        EDMState 1: EDMEffective=Commencement; EDMState 3:
//	        ActualValue.EGNPI=EGNPI, AddendumPremi="1"
//	[5]     DT TreatyInSetEdit — ViewState 0 (RevisionState dikosongkan),
//	        Position="ReasTreatyInAdmin", IsEditData 0, PositionUsername,
//	        StatusAkseptasi="", Comment=""
//	[6]     TreatyInRevisi_post — [1] salin TreatyIn → OLDDATA, lalu
//	        [2]–[5] pengenal revisi baru
//
// ⚠️ `SetTreatyIn_Act` juga menjalankan `TreatyInInputVis`,
// `ConvertHistoryDate`, `TreatySetReinstatement`, `CheckDuplicateOffer`,
// dan penyegaran Achievement: visibilitas dan turunan yang layar ini hitung
// sendiri, bukan nilai yang draf bawa.

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// DaftarMasterPilihan - grid picker. `jenis` = `revisi` (Add Revision) atau
// `premi` (Add Adjustment Premium, hanya NonProportional).
func (l *Layanan) DaftarMasterPilihan(ctx context.Context, p inti.Pelaku, jenis string) ([]models.BarisMasterPilihan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	switch jenis {
	case "revisi":
		return l.gudang.DaftarMasterPilihan(ctx, false)
	case "premi":
		return l.gudang.DaftarMasterPilihan(ctx, true)
	default:
		return nil, fmt.Errorf("%w: jenis picker harus revisi atau premi, bukan %q", ErrIDTidakSah, jenis)
	}
}

// DrafPenyesuaian - tombol `Choose`: penyesuaian baru yang BELUM tersimpan.
func (l *Layanan) DrafPenyesuaian(ctx context.Context, p inti.Pelaku, m models.MasukanDraf) (models.Penyesuaian, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Penyesuaian{}, err
	}
	m.ID = strings.TrimSpace(m.ID)
	if m.ID == "" {
		return models.Penyesuaian{}, fmt.Errorf("%w: pengenal master kosong", ErrIDTidakSah)
	}
	switch m.InternalType {
	case "1", "2", "3":
	default:
		return models.Penyesuaian{}, fmt.Errorf("%w: jenis penyesuaian harus 1, 2, atau 3, bukan %q", ErrIDTidakSah, m.InternalType)
	}
	dok, ada, err := l.gudang.BacaDokumenMaster(ctx, m.ID)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	if !ada {
		return models.Penyesuaian{}, fmt.Errorf("%w: %s", ErrKontrakTidakAda, m.ID)
	}
	adaRevisi, err := l.gudang.AdaRevisi(ctx, m.ID)
	if err != nil {
		return models.Penyesuaian{}, err
	}
	return SusunDraf(dok, m, adaRevisi, p.AkunID, time.Now()), nil
}

// Komentar `TreatyInSetEditPre` per EDMState — `OperatorID.pxInsName + "  " + …`.
var komentarDraf = map[string]string{
	"1": "Had Created Internal Edit",
	"2": "Had Created External Addendum",
	"3": "Had Created Addendum Premium",
}

// SusunDraf - langkah [3]–[6] atas dokumen yang sudah dimuat. Murni.
//
// `operator` = `OperatorID.pxInsName`; `kini` = `@CurrentDateTime()`.
func SusunDraf(dok models.SisiPenyesuaian, m models.MasukanDraf, adaRevisi bool, operator string, kini time.Time) models.Penyesuaian {
	s := salinSisi(dok)
	stempel := StempelPega(kini)

	// [1]/[2] — `TreatyIn.ID = Param.ID`.
	s.Medan["ID"] = m.ID

	// [3] Property-Set "Set EDM Properties".
	// ⚠️ `EDMMaterialType`: ekspor memuat `Param.MaterialType` dan juga
	// ekspresi `@if(Param.InternalType=3, "1", Param.MaterialType)` (memo
	// rule: "param material type removed the complicated logic"). Keduanya
	// SAMA hasilnya — picker Premi mengirim MaterialType '1'.
	s.Medan["OLDID"] = s.Medan["ID"]
	s.Larik["CommentList"] = []map[string]string{}
	s.Medan["EDMState"] = m.InternalType
	s.Medan["EDMMaterialType"] = m.MaterialType
	s.Medan["RevisionDate"] = stempel

	// [4] DT `TreatyInSetEditPre`.
	if catatan, ada := komentarDraf[m.InternalType]; ada {
		s.Larik["CommentList"] = append(s.Larik["CommentList"], map[string]string{
			"IsApproved":   "",
			"OperatorName": operator,
			"Date":         stempel,
			"Suggest":      operator + "  " + catatan,
		})
	}
	switch m.InternalType {
	case "1":
		s.Medan["EDMEffective"] = s.Medan["Commencement"]
	case "3":
		if egnpi, ada := s.Larik["EGNPI"]; ada {
			s.Larik["ActualValue.EGNPI"] = salinLarik(egnpi)
		}
		s.Medan["AddendumPremi"] = "1"
	}

	// [5] DT `TreatyInSetEdit`.
	//
	// ⛔ `RevisionState` TERSIMPAN TIDAK MENGUNCI LAYAR INI, dan itu ralat
	// atas salah tafsir yang dilaporkan pemilik proses 8 Oktober 2026
	// (*"kenapa setelah saya coba tidak bisa edit juga"*).
	//
	// `TreatyInSetEdit[2]` memang berbunyi `WHEN TreatyIn.RevisionState==1 →
	// ViewState = 1`, dan dahulu itu kami baca sebagai kolom. Ia bukan:
	// `SetTreatyIn_Act[7]` — satu-satunya yang menyetel properti itu —
	// berprasyarat `param.revisionstate==1`, yakni PARAMETER TOMBOL.
	// `Section/InputTreatyInOffer` memasangkannya berpasangan:
	//
	//	Edit  → viewstate=<kosong>  revisionstate=<kosong>
	//	View  → viewstate=1         revisionstate=1
	//
	// Jadi `RevisionState` pada halaman berarti *"halaman ini dibuka lewat
	// View sebuah revisi"*, bukan *"kontrak ini pernah direvisi"*. Kolom
	// `REVISIONSTATE` (migrasi 448) menyimpan yang KEDUA — ia dipakai tangga
	// akseptasi `modul/treatyin` — dan memakainya untuk mengunci layar ini
	// adalah pemakaian ulang yang keliru.
	//
	// ⚠️ Akibat salah tafsir itu nyata: kontrak yang pernah lewat tombol
	// Revision di layar Treaty In tersimpan ber-`REVISIONSTATE = 1`, maka
	// SETIAP penyesuaian turunannya lahir terkunci dan Add Revision menjadi
	// tombol yang tidak menghasilkan apa-apa.
	s.Medan["ViewState"] = "0"

	// ⭐ Dokumen BARU, jadi `RevisionState` dikosongkan — pola `TreatyInCopy[1]`
	// (`TreatyIn.RevisionState = ""`), satu-satunya tempat di korpus yang
	// melahirkan dokumen dari dokumen lain. Tanpa ini nilai warisan master
	// ikut terbawa ke revisi dan menguncinya di pembukaan berikutnya.
	s.Medan["RevisionState"] = ""
	s.Medan["Position"] = "ReasTreatyInAdmin"
	s.Medan["IsEditData"] = "0"
	s.Medan["PositionUsername"] = operator
	s.Medan["StatusAkseptasi"] = ""
	s.Medan["Comment"] = ""

	// [6] `TreatyInRevisi_post` [1] — `TreatyInSetAddendumToHistory`:
	// salinan SEBELUM pengenalnya berganti.
	lama := salinSisi(s)

	// [2]–[5] pengenal revisi baru.
	s.Medan["ID"] = IDRevisiBaru(m.ID, adaRevisi)
	s.Medan["OLDID"] = m.ID

	// ⭐ Draf dari master yang isinya TIDAK ada di pendaratan — layar
	// menguncinya (`Penyesuaian.Terdarat`).
	return models.Penyesuaian{ID: s.Medan["ID"], IDAsal: m.ID, Baru: s, Lama: lama, Terdarat: dok.Terdarat}
}

// IDRevisiBaru - `TreatyInRevisi_post` [4]/[5].
//
//	[4] belum ada revisi → `ID + "/R01"`
//	[5] sudah ada        → `@If(n < 10, dasar+"/R0"+(n+1), dasar+"/R"+(n+1))`
//
// `n` = dua digit akhiran pengenal yang dipilih (`SUBSTR(ID,10,2)` di SQL
// `GetTreatyRevisionID`). Diadu dengan `TREATY_IN_EDM`: 263 `R01` ber-OLDID
// master, 16 `R02` dan 1 `R03` ber-OLDID revisi sebelumnya — `X/R01` → `X/R02`.
//
// ⚠️ Disalin APA ADANYA, termasuk dua kejanggalannya: syaratnya `n < 10`,
// bukan `n+1 < 10` (`R09` → `R010`), dan master 7 aksara yang SUDAH punya
// revisi mendapat `n = 0` → `R01` lagi. Pemeriksa duplikat ([7]) ber-`//`.
// Pengenal draf ini belum tersimpan; jalur Save yang memutuskan nasibnya.
func IDRevisiBaru(id string, adaRevisi bool) string {
	if !adaRevisi {
		return id + "/R01"
	}
	dasar := id
	if len(dasar) > 7 {
		dasar = dasar[:7]
	}
	n := 0
	if len(id) >= 11 {
		n, _ = strconv.Atoi(id[9:11])
	}
	if n < 10 {
		return dasar + "/R0" + strconv.Itoa(n+1)
	}
	return dasar + "/R" + strconv.Itoa(n+1)
}

// StempelPega - `@CurrentDateTime()` dalam bentuk tersimpan Pega
// (`20201031T170000.000 GMT`), bentuk yang `formatDate` layar kenal.
func StempelPega(t time.Time) string {
	return t.UTC().Format("20060102T150405.000") + " GMT"
}

func salinLarik(xs []map[string]string) []map[string]string {
	out := make([]map[string]string, len(xs))
	for i, b := range xs {
		c := make(map[string]string, len(b))
		for k, v := range b {
			c[k] = v
		}
		out[i] = c
	}
	return out
}

func salinSisi(s models.SisiPenyesuaian) models.SisiPenyesuaian {
	out := models.SisiPenyesuaian{
		Medan: make(map[string]string, len(s.Medan)),
		Larik: make(map[string][]map[string]string, len(s.Larik)),
		Pohon: make(map[string][]map[string]any, len(s.Pohon)),
	}
	for k, v := range s.Medan {
		out.Medan[k] = v
	}
	for k, v := range s.Larik {
		out.Larik[k] = salinLarik(v)
	}
	for k, v := range s.Pohon {
		out.Pohon[k] = salinPohon(v)
	}
	return out
}

func salinPohon(xs []map[string]any) []map[string]any {
	out := make([]map[string]any, len(xs))
	for i, b := range xs {
		c := make(map[string]any, len(b))
		for k, v := range b {
			if anak, ok := v.([]map[string]any); ok {
				c[k] = salinPohon(anak)
			} else {
				c[k] = v
			}
		}
		out[i] = c
	}
	return out
}
