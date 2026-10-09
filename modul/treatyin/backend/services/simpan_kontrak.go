package services

// Tombol tulis form Treaty In — `Save`, `Submit`, `Actions`, `Decline offer`.
//
// ---------------------------------------------------------------------
// ⭐ KEPUTUSAN PEMILIK PROSES (6–7 Oktober 2026)
// ---------------------------------------------------------------------
//   - Save/Submit menyimpan ke tabel masing-masing (`T_TREATY_*`, kepala
//     `TREATY_IN`, kurs `TREATYEXCHANGEYEARLY`), tidak lewat Pega.
//   - Isian TIDAK masuk basis data sebelum tombol ditekan.
//   - Save ≠ Submit: Save hanya menyimpan; Submit menjalankan tangga
//     ReasTreatyInAdmin > SecHead > DeptHead > Director.
//   - `PositionUsername` jalur naik = SATU pemegang workbasket posisi
//     berikutnya (yang pertama secara urut) — tampilan saja, akses dari peran
//     di Kelola User.
//
// ---------------------------------------------------------------------
// Rantai ekspor tiap tombol (Section `TreatyInActionButtons`,
// `TreatyInfoSubmit`, `TreatyInAction`, `TreatyInDeclineConfirmation`)
// ---------------------------------------------------------------------
//
//	Save           DT TreatyInAddNew (Position = ReasTreatyInAdmin) →
//	               SaveTreatyIn_Act
//	Submit         TreatyInSubmit: TreatyInCheckError → Akseptasi_DT
//	               (Accept) → AddCommentList_Act → SaveTreatyIn_Act
//	Actions        TreatyInAkseptasi_Act: Akseptasi_DT (pilihan radio) →
//	               AddCommentList_Act → SaveTreatyIn_Act
//	Decline offer  TreatyInDeclineConfirmation_postact: AddCommentList_Act →
//	               IsApproved/StatusAkseptasi = "Decline" → SaveTreatyIn_Act
//
// ⭐ "Clipboard" Pega di sini = dokumen TERSIMPAN (`T_TREATY_*`, di atas
// kepala `TREATY_IN`) ditimpa kunci yang layar kirim. Tab yang tidak dibuka
// tidak dikirim layar, dan tidak terhapus.

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// ErrTombolDitolak - penekanan tombol ditolak aturan ekspor (pesan galat
// Activity). Handler menjawab 422 dengan pesan APA ADANYA.
var ErrTombolDitolak = errors.New("tombol ditolak")

// ErrBukanPemegangPosisi - yang menekan tidak memegang workbasket posisi
// berkas. Handler menjawab 403 dengan pesannya.
var ErrBukanPemegangPosisi = errors.New("bukan pemegang posisi")

// galatTombol membawa pesan yang ditampilkan apa adanya.
type galatTombol struct {
	jenis error
	pesan string
}

func (g galatTombol) Error() string        { return g.pesan }
func (g galatTombol) Is(target error) bool { return target == g.jenis }

func ditolak(pesan string) error { return galatTombol{ErrTombolDitolak, pesan} }

// Aksi tombol tulis selain Save.
const (
	AksiSubmit    = "submit"
	AksiAkseptasi = "akseptasi"
	AksiDecline   = "decline"
)

// MasukanSimpan - isi layar yang tombol Save kirim.
type MasukanSimpan struct {
	// IDKontrak - kosong = kontrak BARU (`Add`).
	IDKontrak string `json:"idKontrak"`
	// Dokumen - properti `TreatyIn` yang layar pegang (kepala, penampung
	// halaman, Limits/Share Non-Prop), ejaan Pega.
	Dokumen map[string]any `json:"dokumen"`
	// Kurs - grid Rate of Exchange; `nil` = tidak disimpan.
	Kurs *models.KursSimpan `json:"kurs"`
}

// MasukanKirim - Submit, Actions, atau Decline offer.
type MasukanKirim struct {
	MasukanSimpan
	// Aksi - `submit`, `akseptasi` (Actions), atau `decline` (Decline offer).
	Aksi string `json:"aksi"`
	// Pilihan - `ChooseStatusAkseptasi` radio Actions: Accept/Reject/Decline.
	Pilihan string `json:"pilihan"`
}

// HasilSimpan - jawaban tombol tulis.
type HasilSimpan struct {
	ID string `json:"id"`
	// Pesan - `ErrMsg` prosedur Pega: "Data Sudah Disimpan Dengan ID : …".
	Pesan string `json:"pesan"`
	// Posisi/Status/PemegangPosisi SESUDAH tombol — untuk layar.
	Posisi         string `json:"posisi"`
	Status         string `json:"status"`
	PemegangPosisi string `json:"pemegangPosisi"`
	// KunciTakTersimpan - properti yang layar kirim tetapi BELUM punya kolom
	// di tabel pendaratan. ⛔ Dilaporkan, tidak ditelan.
	KunciTakTersimpan []string `json:"kunciTakTersimpan"`
	// SalinanLampiran - draf penyesuaian yang PERTAMA kali tersimpan:
	// nasib salinan lampiran master (`TreatyRevisionCopyAttachment`,
	// `salin_lampiran_penyesuaian.go`). nil = tidak ada yang disalin.
	SalinanLampiran *HasilSalinLampiran `json:"salinanLampiran,omitempty"`
}

// kunciMilikServer - properti yang HANYA tombol yang mengubah; nilai kiriman
// layar untuknya diabaikan.
var kunciMilikServer = map[string]bool{
	"ID": true, "StatusAkseptasi": true, "Position": true, "PositionUsername": true,
	"ChooseStatusAkseptasi": true, "CommentList": true,
	// ⭐ Jalur revisi — hanya tombol Revision dan tangga yang mengubahnya.
	"RevisionState": true, "ViewState": true,
}

// kunciSesi - halaman SESI Pega (`SearchData`, `TempQuarter*`, `FlagExcel`)
// yang ikut penampung layar tetapi bukan properti `TreatyIn`.
func kunciSesi(k string) bool {
	return strings.HasPrefix(k, "SearchData.") || strings.HasPrefix(k, "FlagExcel.") || strings.HasPrefix(k, "TempQuarter")
}

// SimpanKontrak - tombol Save.
func (l *Layanan) SimpanKontrak(ctx context.Context, p inti.Pelaku, m MasukanSimpan) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	doc, err := l.dokumenKerja(ctx, m)
	if err != nil {
		return HasilSimpan{}, err
	}
	// Syarat tampil Save: `StatusAkseptasi != 'Resolve Complete'`.
	// ⛔ Kecuali Force Edit divisi IT (keputusan pemakai 9 Oktober 2026) —
	// menyimpang dari Pega, yang Save-nya tetap tersembunyi di sana.
	tuntas := teksDok(doc, "StatusAkseptasi") == models.StatusTuntas
	if tuntas {
		it, err := l.adalahIT(ctx, p)
		if err != nil {
			return HasilSimpan{}, err
		}
		if !it {
			return HasilSimpan{}, ditolak("Kontrak sudah Resolve Complete — tombol Save tidak berlaku.")
		}
	}
	if err := l.bolehUbah(ctx, p, doc); err != nil {
		return HasilSimpan{}, err
	}
	// DT `TreatyInAddNew` [1] (pra-DT tombol Save) — DIBERI SYARAT.
	//
	// ---------------------------------------------------------------------
	// ⛔ CACAT YANG DITUTUP DI SINI (8 Oktober 2026)
	// ---------------------------------------------------------------------
	// Laporan pemakai atas kontrak uji `1002305`: *"di submit terus di
	// actions sudah di accept terus tetapi belum pernah resolve complete"*.
	//
	// Keadaan terukur di basis data: `POSITION = ReasTreatyInDeptHead`
	// dengan TIGA komentar akseptasi. Tanpa Save, tiga akseptasi berakhir di
	// `ReasTreatyInDirector`; satu-satunya urutan yang menghasilkan Dept Head
	// adalah Submit → SAVE → Accept → Accept.
	//
	// Sebabnya baris ini: dulu ia menyetel `Admin` TANPA SYARAT, sehingga
	// setiap Save menjatuhkan berkas kembali ke pengajunya. Tombol Save
	// tampil di SEMUA anak tangga (`IsEditData != '1' && StatusAkseptasi !=
	// 'Resolve Complete'`), jadi penyetuju yang menyunting apa pun lalu
	// menekan Save menghapus kemajuannya sendiri — dan `Resolve Complete`
	// tidak pernah tercapai.
	//
	// ---------------------------------------------------------------------
	// ⭐ SYARATNYA DIAMBIL DARI DT YANG SAMA, BUKAN DIKARANG
	// ---------------------------------------------------------------------
	// `DataTransform/TreatyInAddNew.xml` berisi empat langkah:
	//
	//	[1]      SET  TreatyIn.Position = "ReasTreatyInAdmin"
	//	[2]      WHEN Param.Status == 1
	//	[2.1]    WHEN TreatyIn.Position == "" || Position == "ReasTreatyInAdmin"
	//	[2.1.1]  SET  TreatyIn.PositionUsername = OperatorID.pxInsName
	//
	// ⚠️ Syarat `[2.1]` MUSTAHIL SALAH bila `[1]` baru saja memaksa Admin —
	// ia kode mati. Keberadaannya pertanda penulisnya mengira `[1]` tidak
	// selalu berlaku. Syarat itulah yang dipakai di sini.
	//
	// ⛔ PERBEDAAN DARI PEGA DINYATAKAN, bukan disembunyikan: Pega
	// menjalankan `[1]` tanpa syarat (tombol Save ber-`pyPreDataTransform`
	// `TreatyInAddNew` sebelum `SaveTreatyIn_Act`). Keputusan pemilik proses
	// 8 Oktober 2026 memilih syarat DT-nya sendiri daripada perilaku itu.
	// Kontrak tuntas (Force Edit IT) TIDAK disentuh posisinya: `Position`
	// kosong di sana berarti "selesai", bukan "draf Admin".
	if posisi := teksDok(doc, "Position"); !tuntas && (posisi == models.PosisiKosong || posisi == models.PosisiAdmin) {
		doc["Position"] = models.PosisiAdmin
	}
	return l.tulis(ctx, p, m, doc)
}

// DivisiIT - divisi akun (`M_LOGIN_GO.DIVISION_CODE`) pemegang Force Edit.
const DivisiIT = "IT"

// bolehUbah - syarat tombol `Edit` Pega (`Section/InputTreatyInOffer.xml`
// @268974), ditegakkan di Save:
//
//	pyWorkBasketList(2) = 'ReasTreatyInAdmin' && (Position = 'ReasTreatyInAdmin'
//	|| Position = '') && StatusAkseptasi != 'Decline' && != 'Resolve Complete'
//
// SecHead / DeptHead / Director hanya punya `View` (form hanya-baca) dan
// bertindak lewat `Actions`; Admin pun tidak dapat mengubah berkas yang sudah
// naik sampai di-Reject kembali. Satu-satunya pengecualian: Force Edit akun
// divisi IT — padanan `Force Edit (dev)` Pega (keputusan pemakai 9 Oktober
// 2026). Force Edit TIDAK memindah posisi: `TreatyInAddNew` [1] di
// `SimpanKontrak` hanya menyetel Admin di pangkal tangga.
func (l *Layanan) bolehUbah(ctx context.Context, p inti.Pelaku, doc map[string]any) error {
	posisi := teksDok(doc, "Position")
	status := teksDok(doc, "StatusAkseptasi")
	pangkal := posisi == models.PosisiKosong || posisi == models.PosisiAdmin
	if punyaPeran(p, models.PosisiAdmin) && pangkal && status != models.StatusDecline {
		return nil
	}
	it, err := l.adalahIT(ctx, p)
	if err != nil {
		return err
	}
	if it {
		return nil
	}
	switch {
	case status == models.StatusDecline:
		return ditolak("Kontrak berstatus Decline — tidak dapat diubah.")
	case !pangkal:
		return galatTombol{ErrBukanPemegangPosisi, fmt.Sprintf(
			"Berkas ini di posisi %s — hanya Admin yang dapat mengubahnya, setelah berkas dikembalikan (Reject).", posisi)}
	default:
		return galatTombol{ErrBukanPemegangPosisi, fmt.Sprintf(
			"Akun %s tidak memegang workbasket %s — kontrak hanya dapat diubah Admin.", p.AkunID, models.PosisiAdmin)}
	}
}

// adalahIT - akun berdivisi IT (`M_LOGIN_GO.DIVISION_CODE`), pemegang Force Edit.
func (l *Layanan) adalahIT(ctx context.Context, p inti.Pelaku) (bool, error) {
	divisi, err := l.gudang.DivisiAkun(ctx, p.AkunID)
	if err != nil {
		return false, err
	}
	return divisi == DivisiIT, nil
}

// KirimKontrak - Submit, Actions, atau Decline offer.
func (l *Layanan) KirimKontrak(ctx context.Context, p inti.Pelaku, m MasukanKirim) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	doc, err := l.dokumenKerja(ctx, m.MasukanSimpan)
	if err != nil {
		return HasilSimpan{}, err
	}
	status := teksDok(doc, "StatusAkseptasi")
	posisi := teksDok(doc, "Position")
	if posisi == models.PosisiKosong {
		posisi = models.PosisiAdmin
	}
	if !punyaPeran(p, posisi) {
		return HasilSimpan{}, galatTombol{ErrBukanPemegangPosisi,
			fmt.Sprintf("Berkas ini menunggu di posisi %s, dan akun %s tidak memegang workbasket itu.", posisi, p.AkunID)}
	}
	operator := p.AkunID
	stempel := stempelPega(time.Now())

	switch m.Aksi {
	case AksiDecline:
		// TreatyInDeclineConfirmation_postact [2]–[3].
		tambahKomentar(doc, stempel, operator, status, teksDok(doc, "Comment"))
		setelKomentarTerakhir(doc, "IsApproved", models.StatusDecline)
		doc["StatusAkseptasi"] = models.StatusDecline
		return l.tulis(ctx, p, m.MasukanSimpan, doc)
	case AksiSubmit, AksiAkseptasi:
	default:
		return HasilSimpan{}, ditolak(fmt.Sprintf("Aksi %q tidak dikenal.", m.Aksi))
	}

	pilihan := models.PilihAccept
	if m.Aksi == AksiSubmit {
		// TreatyInCheckError — langkah berurutan; yang TERAKHIR menang.
		if msg := periksaGalatSubmit(doc); msg != "" {
			return HasilSimpan{}, ditolak(msg)
		}
	} else {
		pilihan = m.Pilihan
	}
	langkah, err := models.LangkahBerikut(posisi, pilihan, status, teksDok(doc, "RevisionState") == "1")
	switch {
	case errors.Is(err, models.ErrSudahTuntas):
		return HasilSimpan{}, ditolak("Kontrak sudah Resolve Complete — tidak dapat dikirim ulang.")
	case errors.Is(err, models.ErrPilihanTakBerlaku), errors.Is(err, models.ErrPosisiTakDikenal):
		return HasilSimpan{}, ditolak(fmt.Sprintf("Pilihan %q tidak berlaku dari posisi %s.", pilihan, posisi))
	case err != nil:
		return HasilSimpan{}, err
	}
	pemegang := []string{}
	if langkah.AsalNama == models.AsalNamaKosong && langkah.Posisi != models.PosisiKosong {
		if pemegang, err = l.gudang.PemegangPosisi(ctx, langkah.Posisi); err != nil {
			return HasilSimpan{}, err
		}
	}
	// Akseptasi_DT, lalu AddCommentList_Act — riwayat MEMBACA status baru.
	TerapkanLangkah(doc, langkah, pilihan, pemegang)
	tambahKomentar(doc, stempel, operator, teksDok(doc, "StatusAkseptasi"), teksDok(doc, "Comment"))
	return l.tulis(ctx, p, m.MasukanSimpan, doc)
}

// periksaGalatSubmit - `TreatyInCheckError`: Ceding lalu Source of Business;
// pesan langkah TERAKHIR yang menang (`OutputParam.ERRMSG` ditimpa).
func periksaGalatSubmit(doc map[string]any) string {
	msg := ""
	if strings.TrimSpace(teksDok(doc, "Ceding")) == "" {
		msg = "Please input Ceding"
	}
	if strings.TrimSpace(teksDok(doc, "LeadingReinsSource")) == "" {
		msg = "Please input Source of Business (SoB)"
	}
	return msg
}

// TerapkanLangkah - `Akseptasi_DT` atas dokumen: Position, StatusAkseptasi,
// ChooseStatusAkseptasi, PositionUsername (dan RevisionState/ViewState bila
// langkah menyentuhnya).
//
//	komentar-pertama   `CommentList(1).OperatorName`   (Reject jalur biasa)
//	komentar-terakhir  `CommentList(<LAST>).OperatorName` (Reject revisi)
//	naik               pemegang workbasket posisi baru, dipisah ", "
//	tuntas / Decline   kosong
func TerapkanLangkah(doc map[string]any, l models.Langkah, pilihan string, pemegang []string) {
	doc["Position"] = l.Posisi
	doc["StatusAkseptasi"] = l.Status
	doc["ChooseStatusAkseptasi"] = pilihan
	komentar := barisKomentar(doc)
	switch l.AsalNama {
	case models.AsalNamaKomentarPertama:
		doc["PositionUsername"] = ""
		if len(komentar) > 0 {
			doc["PositionUsername"] = teksDok(komentar[0], "OperatorName")
		}
	case models.AsalNamaKomentarTerakhir:
		doc["PositionUsername"] = ""
		if len(komentar) > 0 {
			doc["PositionUsername"] = teksDok(komentar[len(komentar)-1], "OperatorName")
		}
	default:
		// ⛔ SATU NAMA, bukan daftar — permintaan pemilik proses 8 Oktober
		// 2026: *"dibuat salah satu nya di tampilan tp bisa diakses semua
		// yg dapat Workbasket itu"*.
		//
		// Sebelumnya seluruh pemegang digabung `", "`, dan kolom
		// `Position To` berbunyi `Dastin, JEFRIHARI, SUPERADMIN` — tiga
		// akun yang sama memegang KEEMPAT anak tangga, sehingga daftarnya
		// sama di rung mana pun dan nol memberi tahu apa pun.
		//
		// ⭐ AKSES NOL BERUBAH, dan itu bagian kedua permintaannya. Kolom ini
		// MURNI TAMPILAN: pagar akseptasi `punyaPeran(p, posisi)` membaca
		// PERAN pelaku, dan tombol `Actions` di layar membaca daftar peran
		// (`bolehActions`). Nol jalur yang membandingkan nama pemakai dengan
		// `PositionUsername` — `TestSemuaPemegangTetapBolehBertindak`
		// memakukannya.
		//
		// ⚠️ Yang dipilih yang PERTAMA SECARA URUT, bukan yang pertama dari
		// basis data: urutan baris SQL tanpa `ORDER BY` tidak dijamin, dan
		// nama yang berganti-ganti antar penyimpanan akan terbaca sebagai
		// berkas yang berpindah tangan.
		urut := append([]string{}, pemegang...)
		sort.Strings(urut)
		doc["PositionUsername"] = ""
		if len(urut) > 0 {
			doc["PositionUsername"] = urut[0]
		}
		if l.Posisi == models.PosisiKosong {
			doc["PositionUsername"] = ""
		}
	}
	if l.RevisionState != nil {
		doc["RevisionState"] = *l.RevisionState
	}
	if l.ViewState != nil {
		doc["ViewState"] = *l.ViewState
	}
}

// tambahKomentar - `AddCommentList_Act` [3]: baris riwayat baru.
func tambahKomentar(doc map[string]any, stempel, operator, status, komentar string) {
	l, _ := doc["CommentList"].([]any)
	doc["CommentList"] = append(l, map[string]any{
		"Date":         stempel,
		"OperatorName": operator,
		"IsApproved":   status,
		"Suggest":      komentar,
	})
}

func setelKomentarTerakhir(doc map[string]any, kunci, nilai string) {
	l, _ := doc["CommentList"].([]any)
	if len(l) == 0 {
		return
	}
	if el, ok := l[len(l)-1].(map[string]any); ok {
		el[kunci] = nilai
	}
}

func barisKomentar(doc map[string]any) []map[string]any {
	l, _ := doc["CommentList"].([]any)
	out := make([]map[string]any, 0, len(l))
	for _, e := range l {
		if el, ok := e.(map[string]any); ok {
			out = append(out, el)
		}
	}
	return out
}

// dokumenKerja - "clipboard": kepala `TREATY_IN` ⊕ dokumen pendaratan ⊕
// kiriman layar (kecuali kunci milik server dan halaman sesi).
func (l *Layanan) dokumenKerja(ctx context.Context, m MasukanSimpan) (map[string]any, error) {
	doc := map[string]any{}
	id := strings.TrimSpace(m.IDKontrak)
	if id != "" {
		kepala, ada, err := l.gudang.BacaKepalaTreatyIn(ctx, id)
		if err != nil {
			return nil, err
		}
		tersimpan, err := l.gudang.BacaDokumenPendaratan(ctx, id)
		if err != nil {
			return nil, err
		}
		if !ada && len(tersimpan) == 0 {
			return nil, fmt.Errorf("%w: %s", ErrKontrakTidakAda, id)
		}
		for k, v := range kepala {
			doc[k] = v
		}
		for k, v := range tersimpan {
			doc[k] = v
		}
	}
	for k, v := range m.Dokumen {
		if kunciMilikServer[k] || kunciSesi(k) {
			continue
		}
		doc[k] = v
	}
	return doc, nil
}

// tulis - `SaveTreatyIn_Act`: langkah [1] lalu tulisan, dan pesan prosedur.
func (l *Layanan) tulis(ctx context.Context, p inti.Pelaku, m MasukanSimpan, doc map[string]any) (HasilSimpan, error) {
	// SaveTreatyIn_Act [1] — hanya bila sumbernya ada di dokumen: properti
	// layar `FacShare*` tidak punya kolom, dan menyalin kekosongan akan
	// menghapus nilai tersimpan.
	if v, ada := doc["FacShare"]; ada {
		doc["FacultativeShare"] = v
	}
	if v, ada := doc["FacShareBrokerage"]; ada {
		doc["FacultativeShareBrokerage"] = v
	}
	var kurs *models.KursSimpan
	if m.Kurs != nil {
		k := *m.Kurs
		if strings.TrimSpace(k.Tahun) == "" {
			k.Tahun = teksDok(doc, "TreatyYear")
		}
		// ⛔ Kotak tanggal layar mengirim bentuk TERSIMPAN (`YYYYMMDD`) di
		// `…Asli`; `TREATYEXCHANGEYEARLY.STARTDATE` menyimpan STEMPEL Pega
		// (`20180101T000000.000 GMT`). Dikembalikan ke bentuk itu di sini,
		// bukan di repository: `stempelPega` tinggal di lapis ini, dan dua
		// pengubah bentuk tanggal adalah cara termudah keduanya berbeda.
		//
		// ⚠️ Yang KOSONG tetap kosong — pemakai yang mengosongkan kotaknya
		// memang bermaksud mengosongkannya.
		k.Baris = append([]models.BarisKursWarisan{}, k.Baris...)
		for i := range k.Baris {
			b := &k.Baris[i]
			b.BerlakuDari = stempelHariPega(b.BerlakuDariAsli)
			b.BerlakuSampai = stempelHariPega(b.BerlakuSampaiAsli)
		}
		kurs = &k
	}
	// Kepala `TREATY_IN` mendarat di tabelnya sendiri — bukan kunci asing.
	diKepala := map[string]bool{}
	for _, k := range repository.KolomKepalaTreatyIn {
		diKepala[k[1]] = true
	}
	// ⭐ Laporan atas properti BERISI saja — larik kosong dan nilai kosong
	// tidak membawa apa pun yang dapat hilang (8 Oktober 2026: pembaca pohon
	// menyisipkan larik kosong, dan laporan menyebutnya "TIDAK tersimpan").
	laporan := repository.TanpaKosong(doc)
	asing := []string{}
	for _, ks := range repository.KunciTakTerpetakan(laporan) {
		for _, k := range ks {
			if !diKepala[k] {
				asing = append(asing, k)
			}
		}
	}
	// ⭐ Dan yang peta KENAL tetapi kolomnya menunggu migrasi dipasang.
	belum, err := l.gudang.KunciBelumTerpasang(ctx, laporan)
	if err != nil {
		return HasilSimpan{}, err
	}
	asing = unik(append(asing, belum...))
	sort.Strings(asing)
	id, err := l.gudang.SimpanKontrak(ctx, models.RencanaSimpan{
		ID: strings.TrimSpace(m.IDKontrak), Dokumen: doc, Kurs: kurs,
		Operator: p.AkunID, Stempel: stempelPega(time.Now()),
	})
	if err != nil {
		return HasilSimpan{}, err
	}
	return HasilSimpan{
		ID: id, Pesan: "Data Sudah Disimpan Dengan ID : " + id,
		Posisi: teksDok(doc, "Position"), Status: teksDok(doc, "StatusAkseptasi"),
		PemegangPosisi: teksDok(doc, "PositionUsername"), KunciTakTersimpan: asing,
	}, nil
}

func punyaPeran(p inti.Pelaku, peran string) bool {
	for _, x := range p.Peran {
		if x == peran {
			return true
		}
	}
	return false
}

func teksDok(m map[string]any, k string) string {
	t, _ := repository.NilaiTeks(m[k])
	return t
}

// stempelPega - `@CurrentDateTime()` dalam bentuk tersimpan Pega.
func stempelPega(t time.Time) string {
	return t.UTC().Format("20060102T150405.000") + " GMT"
}

// stempelHariPega - `YYYYMMDD` menjadi stempel Pega tengah malam.
//
// ⭐ Jamnya NOL, bukan jam sekarang: yang disimpan tanggal berlakunya kurs,
// dan menambahkan jam membuat dua baris bertanggal sama terbaca berbeda.
//
// ⚠️ Yang BUKAN delapan digit dikembalikan apa adanya — termasuk teks
// kosong. Memaksanya menjadi tanggal akan mengarang nilai yang tidak
// seorang pun ketikkan.
func stempelHariPega(yyyymmdd string) string {
	s := strings.TrimSpace(yyyymmdd)
	if len(s) != 8 {
		return s
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return s
		}
	}
	return s + "T000000.000 GMT"
}
