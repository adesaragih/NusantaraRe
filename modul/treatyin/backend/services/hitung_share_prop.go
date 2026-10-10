package services

// RUMUS TAB SHARE (cabang PROPORSIONAL) - dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ KETERGANTUNGAN ANTARTAB: Share MEMBACA DAN MENULIS `TreatyIn.Limits`
// ---------------------------------------------------------------------
// Grid `Kind of Treaty` sub-tab RNM Share (`Section/TreatyInShareProp.xml`)
// adalah page list `TreatyIn.Limits` YANG SAMA dengan tab Limits; rinciannya
// (`TotalLimits` → `DetailShare`) adalah `.Detail[]` yang sama. Rumus di
// bawah membaca CessionList/IOOLimitList yang diisi di tab Limits dan
// menulis RNMShare/RNMShareList/spreading ke Detail yang sama. Layar
// mengirim `Limits` dari penampung halaman (`frontend/halaman.tsx`) dan
// menaruh jawabannya kembali di sana - tab Limits dan Achievement melihat
// perubahan yang sama.
//
// ---------------------------------------------------------------------
// Pemicu (ekspor)
// ---------------------------------------------------------------------
//
//	share      tombol Refresh; isian % RNM Share / Option / Share to Other
//	           Retro (perilaku change)        → Activity/TreatyInPropshare.xml
//	detail     isian % RNM Share rincian Detail → Activity/TreatyInPropshareDetail.xml
//	spreading  isian Spreading Type rincian Detail
//	           → Activity/FetchQSfromMaster.xml (ParentReinsTypeID = .SpreadingTypeID)
//	sebar-nama sel Reins Type / Pct Share spreading manual
//	           → Activity/SetSpreadName.xml
//
// ---------------------------------------------------------------------
// TreatyInPropshare
// ---------------------------------------------------------------------
//
//	1   Page-Clear-Messages
//	3   Property-Remove TotalShareRnmLimit, TotalShareRnmProp,
//	    TotalSpreadedRnmProp, TotalSpreadedRnmRIProp
//	4   per Limits.Detail: .Currency = .CurrencyIOOLimit;
//	    .NusaReLimit = .Cession × @divide(RNMShareP,100,4); hapus .RNMShareList
//	5   per Limits.Detail: .RNMShare = RNMShareP; .Brokerage = BrokeragePercentP;
//	    CalculateShareList(type=calculate, TreatyType = .TreatyType Detail);
//	    SpreadingTypeID != "" → FetchQSfromMaster (TANPA parameter)
//	    SpreadingTypeID == "" → SetSpreadName
//	6   total per mata uang atas SELURUH Detail: RNMShareList → TotalShareRnmProp,
//	    RNMSpreadedList(RI) → TotalSpreadedRnm(RI)Prop; nilai kosong/0 di
//	    sebaran (kecuali SpreadingTypeID "10007") → ErrorFlag
//	7   Param.ErrorSpreading = "Spreading is incomplete at : " + TrtKind +
//	    " ; " + TrtGroup (Kind/Group TERAKHIR yang ditapaki)
//	8   ⛔ MATI (`pyStepsBlockName == "//"`) — Property-Set-Messages atas teks
//	    langkah 7. Nol yang membaca `Param.ErrorSpreading` sesudah ia mati,
//	    jadi pesan itu TIDAK PERNAH tampil di Pega. TIDAK dibangun.
//
// CalculateShareList (type=calculate):
//
//	2  ⛔ MATI (`pyStepsBlockName == "//"`) — RNMShareAcrossTheBoard == true:
//	   CessionList × @divide(RNMShareP,100,20). TIDAK dibangun; selama ia
//	   dibangun, `Total Share RNM Limit` tampil dua kali lipat.
//	3  OptionLimit == 1: CessionList × @divide(.RNMShare,100,20);
//	   ShareNote = QUOTA SHARE ? " of " + CessionPct + "% of 100%" : " of 100%"
//	4  OptionLimit == 2: IOOLimitList × @divide(.RNMShare,100,20); ShareNote " of 100%"
//	(ShareNote ditulis DI DALAM perulangan - daftar kosong tidak mengubahnya)
//
// ⚠️ DISALIN APA ADANYA WALAU JANGGAL - rumus, bukan alamat:
//   - `SetSpreadName` langkah 5 MENJUMLAH total seluruh Detail ke total akar
//     setiap kali dipanggil (sekali per Detail tanpa Spreading Type), lalu
//     `TreatyInPropshare` langkah 6 menjumlah lagi. Untuk Detail tanpa
//     Spreading Type total menjadi berlipat.
//
//     ⛔ DIGANTIKAN 8 Oktober 2026: layar Pega PRODUKSI memperlihatkan total
//     = Σ Detail SATU kali (lihat `jumlahkanTotal`). Total kini selalu
//     dihitung ulang dari nol; catatan di bawah dipertahankan sebagai
//     riwayat bacaan ekspor yang ternyata tidak cocok dengan layar Pega.
//
//     (Riwayat) KEPUTUSAN PEMILIK PROSES 7 Oktober 2026: IKUTI PEGA, BIARKAN
//     BERLIPAT. Ditanyakan dengan angka di tangan dan dijawab dengan sadar -
//     ini BUKAN cacat yang belum ketahuan, dan BUKAN pula sesuatu yang
//     menunggu diperbaiki.
//
//     Yang memicunya, terukur dari laporan pemakai hari itu:
//
//     	Cession to R/I     1.757.675.000.000,00
//     	% RNM Share                   1,28 %
//     	RNMShareList          22.498.240.000,00   <- per Detail, benar
//     	TotalShareRnmProp     44.996.480.000,00   <- tepat 2x
//
//     Mengapa ia berlipat, dari ekspor: `SetSpreadName` langkah [5]
//     ber-`pyStepsObjectName = TreatyIn.Limits` - ia menapaki SELURUH pohon,
//     bukan Detail yang sedang dikerjakan - dan `TreatyInPropshare` langkah
//     [6] menapaki pohon yang sama dengan rumus yang persis sama.
//     Pengosongan totalnya hanya SEKALI, di langkah [3], sebelum kedua
//     gelung.
//
//     (Riwayat) ⛔ JANGAN "MEMPERBAIKINYA". Cabang ber-Spreading Type memanggil
//     `FetchQSfromMaster`, bukan `SetSpreadName`, sehingga ia TIDAK berlipat
//     (Gambar Pega 17, dan terukur ulang 7 Oktober 2026: 2.000.000.000 lawan
//     2.000.000.000). Selisih kedua cabang itu NYATA di Pega, dan layar kita
//     memang harus memperlihatkannya.
//     `TestSharePropTotalSatuKaliSepertiPega` kini memaku perilaku baru.
//   - `SetSpreadName` langkah 6 memasang "Total share must equal RNM
//     share.!!" bila Total Spreading Pct ≠ % RNM Share - juga pada Detail
//     tanpa baris spreading (total 0).
//   - Langkah 4 menimpa `.Currency` Detail dengan `.CurrencyIOOLimit`.
//
// ⚠️ TIDAK BERBUKTI, DIPUTUSKAN: `Page-Clear-Messages` pada `TreatyIn`
// dianggap membuang SEMUA pesan sebelumnya (termasuk pesan Detail) - pesan
// yang dikembalikan hanya yang lahir sesudah pembersihan terakhir.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// Aksi tab Share Prop.
const (
	AksiSharePropShare     = "share"
	AksiSharePropDetail    = "detail"
	AksiSharePropSpreading = "spreading"
	AksiSharePropSebarNama = "sebar-nama"
)

// Pesan Activity, apa adanya.
const (
	pesanSpreadingBelumLengkap = "Spreading is incomplete at : "
	pesanTotalShareSpreading   = "Total share must equal RNM share.!!"
	namaQSRI                   = "QS (R/I)"
	idSpreadingTanpaGalat      = "10007"
	jenisFacOut                = "FAC-OUT"
	jenisQuotaShare            = "QUOTA SHARE"
)

// MasukanShareProp - satu aksi tab Share Prop. `Limits` = halaman
// `TreatyIn.Limits` utuh (pohon, ejaan Pega).
type MasukanShareProp struct {
	Aksi                   string           `json:"aksi"`
	Limits                 []map[string]any `json:"Limits"`
	RNMShareP              string           `json:"RNMShareP"`
	BrokeragePercentP      string           `json:"BrokeragePercentP"`
	OptionLimit            string           `json:"OptionLimit"`
	RNMShareAcrossTheBoard string           `json:"RNMShareAcrossTheBoard"`
	// Commencement - `TreatyIn.Commencement` (YYYYMMDD), saringan RD induk.
	Commencement string `json:"Commencement"`
	// Total akar saat ini - `spreading`/`sebar-nama` menjumlah ke sini.
	TotalShareRnmProp      []NilaiMataUang `json:"TotalShareRnmProp"`
	TotalSpreadedRnmProp   []NilaiMataUang `json:"TotalSpreadedRnmProp"`
	TotalSpreadedRnmRIProp []NilaiMataUang `json:"TotalSpreadedRnmRIProp"`
	// ParentReinsTypeID - `Param.ParentReinsTypeID` aksi `spreading`. `nil` =
	// `.SpreadingTypeID` Detail itu (isian Spreading Type); teks kosong =
	// dipanggil TANPA induk (rantai Treaty Group rincian Limits:
	// `FetchQSfromMaster(ParentReinsTypeID="")`).
	ParentReinsTypeID *string `json:"ParentReinsTypeID"`
	// Detail sasaran (mulai 0) untuk aksi selain `share`.
	IndeksLimit  int `json:"indeksLimit"`
	IndeksDetail int `json:"indeksDetail"`
}

// HasilShareProp - halaman sesudah aksi.
type HasilShareProp struct {
	Limits                 []map[string]any `json:"Limits"`
	TotalShareRnmProp      []NilaiMataUang  `json:"TotalShareRnmProp"`
	TotalSpreadedRnmProp   []NilaiMataUang  `json:"TotalSpreadedRnmProp"`
	TotalSpreadedRnmRIProp []NilaiMataUang  `json:"TotalSpreadedRnmRIProp"`
	Pesan                  []string         `json:"pesan"`
}

// SumberSpreadingProp - kedua RD spreading seperti cabang Prop memanggilnya.
type SumberSpreadingProp interface {
	// Induk - `BrowseTreatyArrangement_ParentReinsMasterTrt`.
	Induk(grup, mulai string) []models.SusunanSpreading
	// AnakProp - `BrowseTreatyArrangement_Limit_MstTrt_RD`, lima filter.
	AnakProp(treatyYear, grup, desc, induk, treatyYearID string) []models.SusunanSpreading
}

// sumberGudangProp - sumber dari gudang; galat pertama disimpan.
type sumberGudangProp struct {
	ctx context.Context
	g   Gudang
	err error
}

func (s *sumberGudangProp) Induk(grup, mulai string) []models.SusunanSpreading {
	if s.err != nil {
		return nil
	}
	// `FetchQSfromMaster[2]` mengirim `TreatyGroupID` DAN
	// `TreatyDescID = "10001"`; hanya parameter H yang dikosongkan. Berbeda
	// dari dropdown `Section/DetailShare.xml`, yang tidak mengirim grup.
	v, err := s.g.BacaIndukSpreading(s.ctx, grup, DescSpreading, mulai, "")
	if err != nil {
		s.err = err
	}
	return v
}

func (s *sumberGudangProp) AnakProp(treatyYear, grup, desc, induk, treatyYearID string) []models.SusunanSpreading {
	if s.err != nil {
		return nil
	}
	v, err := s.g.BacaAnakSpreadingProp(s.ctx, treatyYear, grup, desc, induk, treatyYearID)
	if err != nil {
		s.err = err
	}
	return v
}

// HitungShareProp - bentuk ber-pelaku untuk handler (membaca RD master).
func (l *Layanan) HitungShareProp(ctx context.Context, p inti.Pelaku, m MasukanShareProp) (HasilShareProp, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilShareProp{}, err
	}
	src := &sumberGudangProp{ctx: ctx, g: l.gudang}
	h, err := HitungShareProp(m, src)
	if err != nil {
		return HasilShareProp{}, err
	}
	if src.err != nil {
		return HasilShareProp{}, src.err
	}
	return h, nil
}

// HitungShareProp menjalankan satu aksi atas salinan halaman.
func HitungShareProp(m MasukanShareProp, src SumberSpreadingProp) (HasilShareProp, error) {
	k := &konteksProp{
		m:      m,
		src:    src,
		limits: salinPohon(m.Limits),
		share:  salinNilai(m.TotalShareRnmProp),
		or:     salinNilai(m.TotalSpreadedRnmProp),
		ri:     salinNilai(m.TotalSpreadedRnmRIProp),
		pesan:  []string{},
	}
	switch m.Aksi {
	case AksiSharePropShare:
		k.propshare()
	case AksiSharePropDetail, AksiSharePropSpreading, AksiSharePropSebarNama:
		d := k.detail(m.IndeksLimit, m.IndeksDetail)
		if d == nil {
			return HasilShareProp{}, fmt.Errorf("%w: Detail %d/%d tidak ada", ErrMasukanTidakSah, m.IndeksLimit, m.IndeksDetail)
		}
		switch m.Aksi {
		case AksiSharePropDetail:
			k.propshareDetail(d)
		case AksiSharePropSpreading:
			// Isian Spreading Type: `FetchQSfromMaster(ParentReinsTypeID = .SpreadingTypeID)`.
			induk := teksSimpul(d, "SpreadingTypeID")
			if m.ParentReinsTypeID != nil {
				induk = *m.ParentReinsTypeID
			}
			k.fetchQS(d, induk)
		default:
			k.setSpreadName(d)
		}
	default:
		return HasilShareProp{}, fmt.Errorf("%w: aksi share prop %q", ErrMasukanTidakSah, m.Aksi)
	}
	return HasilShareProp{
		Limits:                 k.limits,
		TotalShareRnmProp:      k.share,
		TotalSpreadedRnmProp:   k.or,
		TotalSpreadedRnmRIProp: k.ri,
		Pesan:                  k.pesan,
	}, nil
}

type konteksProp struct {
	m             MasukanShareProp
	src           SumberSpreadingProp
	limits        []map[string]any
	share, or, ri []NilaiMataUang
	pesan         []string
}

func (k *konteksProp) detail(i, j int) map[string]any {
	if i < 0 || i >= len(k.limits) {
		return nil
	}
	ds := larikSimpul(k.limits[i], "Detail")
	if j < 0 || j >= len(ds) {
		return nil
	}
	return ds[j]
}

// propshare = `TreatyInPropshare`.
func (k *konteksProp) propshare() {
	k.pesan = []string{}                                                          // [1]
	k.share, k.or, k.ri = []NilaiMataUang{}, []NilaiMataUang{}, []NilaiMataUang{} // [3]
	persen4 := bagiBulat(angka(k.m.RNMShareP), 100, 4)
	for _, l := range k.limits { // [4]
		for _, d := range larikSimpul(l, "Detail") {
			d["Currency"] = teksSimpul(d, "CurrencyIOOLimit")
			d["NusaReLimit"] = teks(kali(angka(teksSimpul(d, "Cession")), persen4))
			delete(d, "RNMShareList")
		}
	}
	for _, l := range k.limits { // [5]
		for _, d := range larikSimpul(l, "Detail") {
			d["RNMShare"] = k.m.RNMShareP
			d["Brokerage"] = k.m.BrokeragePercentP
			k.calculateShareList(d)
			if teksSimpul(d, "SpreadingTypeID") != "" {
				k.fetchQS(d, "")
			} else {
				k.setSpreadName(d)
			}
		}
	}
	// [6]
	//
	// ⛔ LANGKAH [8] DICABUT 7 Oktober 2026 — ia MATI di ekspor.
	//
	// Dulu di sini berdiri:
	//
	//	galat, jenis, grup := k.jumlahkanTotal(true)
	//	if galat { // [7-8]
	//		k.pesan = append(k.pesan, pesanSpreadingBelumLengkap+jenis+" ; "+grup)
	//	}
	//
	// ⚠️ Langkah [8] (`Property-Set-Messages`, pra-syarat
	// `Local.ErrorFlag=="1"`) ber-`pyStepsBlockName == "//"`: dikomentari,
	// bukan sekadar pra-syarat nonaktif. Langkah [7] yang MENYUSUN teksnya
	// tetap hidup — ia hanya mengisi `Param.ErrorSpreading`, dan nol yang
	// membacanya sesudah [8] dimatikan. Pesan ini karena itu TIDAK PERNAH
	// tampil di Pega.
	//
	// ⭐ `jumlahkanTotal` TETAP DIPANGGIL: ia yang mengisi ketiga total
	// akar. Hanya PESANNYA yang tidak lahir.
	k.jumlahkanTotal(true)
}

// propshareDetail = `TreatyInPropshareDetail` atas satu Detail.
func (k *konteksProp) propshareDetail(d map[string]any) {
	if k.m.RNMShareAcrossTheBoard == "true" { // [1] keluar
		return
	}
	k.share, k.or, k.ri = []NilaiMataUang{}, []NilaiMataUang{}, []NilaiMataUang{} // [2]
	delete(d, "RNMShareList")                                                     // [3]
	delete(d, "RNMSpreadedList")
	delete(d, "RNMSpreadedListRI")
	k.calculateShareList(d) // [4-5]
	if teksSimpul(d, "SpreadingTypeID") != "" {
		k.fetchQS(d, "") // [6]
	} else {
		k.setSpreadName(d) // [7]
	}
	k.jumlahkanTotal(false) // [8] - ErrorFlag-nya tidak dipakai Activity ini
}

// calculateShareList = `CalculateShareList` (type=calculate).
func (k *konteksProp) calculateShareList(d map[string]any) {
	jenis := teksSimpul(d, "TreatyType")
	daftar := larikSimpul(d, "RNMShareList")
	tambahBaris := func(sumber []map[string]any, persen string) {
		f := bagiBulat(angka(persen), 100, 20)
		for _, c := range sumber {
			daftar = append(daftar, map[string]any{
				"Currency": teksSimpul(c, "Currency"),
				"Value":    teks(kali(angka(teksSimpul(c, "Value")), f)),
			})
		}
	}
	// ⛔ LANGKAH [2] DICABUT 7 Oktober 2026 — ia MATI di ekspor.
	//
	// Dulu di sini berdiri:
	//
	//	if k.m.RNMShareAcrossTheBoard == "true" { // [2]
	//		tambahBaris(larikSimpul(d, "CessionList"), k.m.RNMShareP)
	//	}
	//
	// ⚠️ Langkah [2] ber-`pyStepsBlockName == "//"` di `Treaty In` MAUPUN
	// `Treaty In Adjustment` — artinya DIKOMENTARI, bukan sekadar
	// pra-syaratnya nonaktif. Ketiga penanda "mati" itu berbeda akibatnya,
	// dan yang ini satu-satunya yang benar-benar menghentikan langkahnya
	// (lihat `_migration-docs/alat-baca-ekspor/README.md`, Aturan 2).
	//
	// ⛔ AKIBAT NYATA SELAMA IA HIDUP: `CessionList` ditambahkan DUA KALI
	// — sekali di sini ber-`RNMShareP`, sekali lagi di langkah [3]
	// ber-`.RNMShare` — dan keduanya bernilai sama sebab langkah 5
	// `TreatyInPropshare` menyalin `.RNMShare = RNMShareP`. `Total Share RNM
	// Limit` karena itu tampil TEPAT DUA KALI LIPAT.
	//
	// Terukur pada laporan pemakai 7 Oktober 2026:
	//
	//	Cession to R/I   1.757.675.000.000,00
	//	% RNM Share               1,28 %
	//	seharusnya          22.498.240.000,00
	//	yang tampil         44.996.480.000,00   <- tepat 2x
	//
	// ⚠️ `RNMShareAcrossTheBoard` bawaannya `true` ketika kolomnya kosong
	// (`share_np_muat.go`), jadi cabang ini praktis selalu menyala — itu
	// sebabnya ia menggandakan hampir setiap kontrak, bukan sebagian.
	//
	// ⭐ Propertinya SENDIRI tetap dipakai di tempat lain yang sah:
	// `propshareDetail` langkah [1] (keluar lebih awal) dan cabang Non-Prop.
	// Yang dicabut hanya langkah yang ekspornya komentari.
	opsi := strings.TrimSpace(k.m.OptionLimit)
	if opsi == "1" { // [3]
		cession := larikSimpul(d, "CessionList")
		tambahBaris(cession, teksSimpul(d, "RNMShare"))
		if len(cession) > 0 {
			if jenis == jenisQuotaShare {
				d["ShareNote"] = " of " + teksSimpul(d, "CessionPct") + "% of 100%"
			} else {
				d["ShareNote"] = " of 100%"
			}
		}
	}
	if opsi == "2" { // [4]
		ioo := larikSimpul(d, "IOOLimitList")
		tambahBaris(ioo, teksSimpul(d, "RNMShare"))
		if len(ioo) > 0 {
			d["ShareNote"] = " of 100%"
		}
	}
	d["RNMShareList"] = daftar
}

// fetchQS = `FetchQSfromMaster` atas satu Detail; `induk` =
// `Param.ParentReinsTypeID` (kosong bila dipanggil `TreatyInPropshare`).
func (k *konteksProp) fetchQS(d map[string]any, induk string) {
	// [2]
	d["SpreadingList"] = []map[string]any{}
	d["RNMSpreadedList"] = []map[string]any{}
	d["RNMSpreadedListRI"] = []map[string]any{}
	tipe := teksSimpul(d, "SpreadingTypeID")
	tahun, tahunID := "", ""
	// ⛔ Grup DIBUANG — lihat catatan di `fetchQS` (`hitung_share_np.go`).
	if tipe != "" { // [3] RD induk; [6.1] pasangan Param.ParentReinsTypeID
		for _, p := range k.src.Induk("", k.m.Commencement) {
			if induk == p.ReinsTypeID {
				d["SpreadingType"] = p.ReinsTypeName
				tahunID = p.TreatyYearID
				tahun = p.TreatyYear
			}
		}
	}
	// [7]-[9] RD anak.
	daftar := []map[string]any{}
	if tipe != "" {
		pertama := apd.New(0, 0)
		if rs := larikSimpul(d, "RNMShareList"); len(rs) > 0 {
			pertama = angka(teksSimpul(rs[0], "Value"))
		}
		// ⚠️ Grup dibuang DI SINI JUGA: induk di luar Treaty Group baris
		// punya anak yang juga di luar grup itu, jadi menyaringnya akan
		// mengosongkan pecahan yang baru saja dibuat dapat dipilih.
		for _, a := range k.src.AnakProp(tahun, "", DescSpreading, tipe, tahunID) {
			daftar = append(daftar, map[string]any{
				"ReinsTypeName":     a.ReinsTypeName,
				"ReinsTypeID":       a.ReinsTypeID,
				"ParentReinsTypeID": a.ParentReinsTypeID,
				"Value":             teks(kali(bagiPolos(angka(a.Pct), seratus), pertama)),
				"Pct":               a.Pct,
				"Rp":                a.Rp,
				"Usd":               a.Usd,
			})
		}
	}
	d["SpreadingList"] = daftar
	// [10]
	total := apd.New(0, 0)
	for _, r := range daftar {
		total = tambah(total, angka(teksSimpul(r, "Pct")))
	}
	d["SpreadingTotalPct"] = teks(total)
	// [11]-[12]
	jenisSebar := teksSimpul(d, "SpreadingType")
	qs := apd.New(0, 0)
	for _, r := range daftar {
		nama := teksSimpul(r, "ReinsTypeName")
		if nama == namaQSOR || nama == namaORS {
			qs = angka(teksSimpul(r, "Pct"))
		}
		if jenisSebar == namaORS {
			r["ReinsTypeName"] = namaORS
		}
	}
	// [13]
	or, ri := []map[string]any{}, []map[string]any{}
	fOR, fRI := bagiPolos(qs, seratus), bagiPolos(kurang(seratus, qs), seratus)
	for _, c := range larikSimpul(d, "RNMShareList") {
		v := angka(teksSimpul(c, "Value"))
		or = append(or, map[string]any{"Value": teks(kali(fOR, v)), "Currency": teksSimpul(c, "Currency")})
		ri = append(ri, map[string]any{"Value": teks(kali(fRI, v)), "Currency": teksSimpul(c, "Currency")})
	}
	d["RNMSpreadedList"], d["RNMSpreadedListRI"] = or, ri
	if teksSimpul(d, "TreatyType") == jenisFacOut { // [14]
		d["RNMSpreadedListRI"] = []map[string]any{}
		d["RNMSpreadedList"] = []map[string]any{}
		d["SpreadingType"] = ""
		d["SpreadingTypeID"] = ""
	}
}

// setSpreadName = `SetSpreadName` atas satu Detail (spreading MANUAL).
func (k *konteksProp) setSpreadName(d map[string]any) {
	k.pesan = []string{} // [1] Page-Clear-Messages TreatyIn
	d["RNMSpreadedList"] = []map[string]any{}
	d["RNMSpreadedListRI"] = []map[string]any{}
	or, ri := []NilaiMataUang{}, []NilaiMataUang{}
	tahunID := "" // Param.TreatyYearID bertahan antarbaris (satu halaman param).
	rnm := angka(teksSimpul(d, "RNMShare"))
	for _, r := range larikSimpul(d, "SpreadingList") { // [3]
		delete(r, "BreakDownSprdList")
		var pecahan []map[string]any
		pecahan, tahunID = pecahanSebar(k.src, k.m.Commencement, d, r, tahunID)
		r["BreakDownSprdList"] = pecahan
		for _, b := range pecahan { // [3.3]
			switch teksSimpul(b, "ReinsName") {
			case namaQSOR, namaORS:
				or = tambahPerMataUang(or, teksSimpul(b, "Currency"), "", angka(teksSimpul(b, "Amount")))
			case namaQSRI:
				ri = tambahPerMataUang(ri, teksSimpul(b, "Currency"), "", angka(teksSimpul(b, "Amount")))
			}
		}
	}
	d["RNMSpreadedList"], d["RNMSpreadedListRI"] = keSimpul(or), keSimpul(ri)
	// [4] DataTransform `CountTotalPctSpead`: Σ .Pct.
	total := apd.New(0, 0)
	for _, r := range larikSimpul(d, "SpreadingList") {
		total = tambah(total, angka(teksSimpul(r, "Pct")))
	}
	d["SpreadingTotalPct"] = teks(total)
	// [5] total seluruh Detail - DITAMBAHKAN ke total akar (lihat kepala berkas).
	k.jumlahkanTotal(false)
	// [6]
	if total.Cmp(rnm) != 0 {
		k.pesan = append(k.pesan, pesanTotalShareSpreading)
	}
}

// pecahanSebar = `SetSpreadName` [3.2] atas SATU baris spreading manual `r`
// Detail `d`: nama & `TreatyYearID` induknya (RD `ParentReinsMasterTrt`),
// lalu anak RD `Limit_MstTrt` × `RNMShareList`. `tahunID` = `Param.
// TreatyYearID` yang bertahan antarbaris; dikembalikan untuk baris berikut.
//
// ⛔ 9 Oktober 2026 — laporan pemakai: rincian satu spread `ORS` berisi 24
// baris `ORS` identik. Dropdown Reins Type grid manual menawarkan induk
// SEMUA Treaty Group, sedangkan [3.2.2] mencarinya dengan Treaty Group
// Detail ini — tidak ketemu, `TreatyYearID` kosong, filter E RD anak
// dilewati, dan anak `ORS` (10007) dari SELURUH tahun ikut (24 baris,
// 2017–2025). Dua penjaga:
//  1. tidak ketemu di grup Detail → dicari tanpa grup (periode
//     Commencement tetap berlaku) — induk yang sama dengan isi dropdown;
//  2. `TreatyYearID` tetap kosong → NOL pecahan, tidak pernah "semua tahun".
func pecahanSebar(src SumberSpreadingProp, mulai string, d, r map[string]any, tahunID string) ([]map[string]any, string) {
	pecahan := []map[string]any{}
	id := teksSimpul(r, "ReinsTypeID")
	if id == "" {
		return pecahan, tahunID
	}
	// [3.2.2] TreatyGroupID = TempSprd.TreatyGroupID = `.TreatyGroupID`
	// Detail ini (`AddDelSpreadingTreatyin` [1]).
	for _, grup := range []string{teksSimpul(d, "TreatyGroupID"), ""} {
		ketemu := false
		for _, p := range src.Induk(grup, mulai) {
			if p.ReinsTypeID == id {
				r["ReinsTypeName"] = p.ReinsTypeName
				tahunID = p.TreatyYearID
				ketemu = true
			}
		}
		if ketemu || grup == "" {
			break
		}
	}
	if tahunID == "" {
		return pecahan, tahunID
	}
	belah := bagiBulatDes(angka(teksSimpul(r, "Pct")), angka(teksSimpul(d, "RNMShare")), 20)
	// ⛔ SATU baris per tipe anak per mata uang — pemakai 9 Oktober 2026:
	// *"maksimal 1 ORS, kalau pun ada yang 2 itu karena OR atau RI"*. Satu
	// TreatyYearID di Pega memberi `ORS` → 1 anak (`ORS` 100%), `… TRT` → 2
	// (`QS (OR)` + `QS (R/I)`); baris kembar dari master dibuang di sini.
	sudah := map[string]bool{}
	// [3.2.5] TreatyYear/TreatyGroupID/TreatyDescID dikosongkan.
	for _, a := range src.AnakProp("", "", "", id, tahunID) {
		for _, c := range larikSimpul(d, "RNMShareList") {
			kunci := a.ReinsTypeID + "|" + a.ReinsTypeName + "|" + teksSimpul(c, "Currency")
			if sudah[kunci] {
				continue
			}
			sudah[kunci] = true
			jumlah := kali(kali(angka(teksSimpul(c, "Value")), belah), angka(a.Pct))
			pecahan = append(pecahan, map[string]any{
				"ReinsName": a.ReinsTypeName,
				"ReinsID":   a.ReinsTypeID,
				"Currency":  teksSimpul(c, "Currency"),
				"Amount":    teks(bagiBulat(jumlah, 100, 20)),
				"SharePct":  a.Pct,
			})
		}
	}
	return pecahan, tahunID
}

// LengkapiPecahanSpreading — `BreakDownSprdList` baris spreading MANUAL yang
// KOSONG diisi saat kontrak dibuka, dengan rumus `pecahanSebar` (=
// `SetSpreadName` [3.2]). Laporan pemakai 9 Oktober 2026: sesudah dimuat
// ulang, ▾ baris `2025 QS 181M TRT` berbunyi "No items" padahal Value
// Spreading OR/R/I terisi — kontrak tersimpan sebelum pecahan punya tabel
// (migrasi 454). Mode View tidak punya tombol untuk memunculkannya.
//
// ⛔ HANYA pecahan yang kosong; nol medan lain berubah (Value Spreading,
// total, pesan). Spreading LAMA (Spreading Type terisi) tidak disentuh.
func LengkapiPecahanSpreading(src SumberSpreadingProp, mulai string, limits []map[string]any) {
	for _, l := range limits {
		for _, d := range larikSimpul(l, "Detail") {
			if teksSimpul(d, "SpreadingTypeID") != "" {
				continue
			}
			tahunID := ""
			for _, r := range larikSimpul(d, "SpreadingList") {
				if len(larikSimpul(r, "BreakDownSprdList")) > 0 {
					continue
				}
				var pecahan []map[string]any
				pecahan, tahunID = pecahanSebar(src, mulai, d, r, tahunID)
				if len(pecahan) > 0 {
					r["BreakDownSprdList"] = pecahan
				}
			}
		}
	}
}

// lengkapiPecahanKontrak — `LengkapiPecahanSpreading` atas kontrak yang
// dibuka. Galat RD master TIDAK menggagalkan pembukaan kontrak: pecahan
// tetap kosong, seperti sebelumnya.
func (l *Layanan) lengkapiPecahanKontrak(ctx context.Context, k *models.KontrakWarisan) {
	src := &sumberGudangProp{ctx: ctx, g: l.gudang}
	LengkapiPecahanSpreading(src, k.TanggalMulaiAsli, k.LimitsPohon)
}

// jumlahkanTotal - pola total `Local.test` atas SELURUH `Limits.Detail`
// (`TreatyInPropshare` [6], `PropshareDetail` [8], `SetSpreadName` [5]).
// `periksa` = deteksi ErrorFlag `TreatyInPropshare` [6.1.2.2]/[6.1.3.2].
//
// ⛔ RALAT 8 Oktober 2026 — total DIMULAI DARI NOL setiap kali, BUKAN
// ditambahkan ke total yang sudah ada. Bukti: layar Pega PRODUKSI pemakai
// (Treaty Group HOSPITAL, spreading manual `2025 QS 181M TRT` 25%):
// RNMShareList 56.250.000 → Total Share RNM Limit 56.250.000, Value
// Spreading OR/R/I 22.500.000 / 33.750.000 → Total OR/R/I sama persis.
// Penjumlahan-ke-yang-ada (bacaan ekspor `SetSpreadName` [5]) memberi
// 112.500.000 / 45.000.000 / 67.500.000 di aplikasi — dua kali lipat, dan
// makin berlipat tiap kali Reins Type / Pct Share diubah. Pemakai: *"perbaiki
// yang di aplikasi seharusnya seperti yang dipega"*. Keputusan 7 Oktober
// ("biarkan berlipat") DIGANTIKAN oleh hasil layar Pega yang nyata.
func (k *konteksProp) jumlahkanTotal(periksa bool) (galat bool, jenis, grup string) {
	k.share, k.or, k.ri = []NilaiMataUang{}, []NilaiMataUang{}, []NilaiMataUang{}
	for _, l := range k.limits {
		jenis = teksSimpul(l, "TreatyType")
		for _, d := range larikSimpul(l, "Detail") {
			grup = teksSimpul(d, "TreatyGroup")
			sebar := teksSimpul(d, "SpreadingTypeID")
			for _, c := range larikSimpul(d, "RNMShareList") {
				k.share = tambahPerMataUang(k.share, teksSimpul(c, "Currency"), "", angka(teksSimpul(c, "Value")))
			}
			for _, x := range []struct {
				larik  string
				tujuan *[]NilaiMataUang
			}{{"RNMSpreadedList", &k.or}, {"RNMSpreadedListRI", &k.ri}} {
				for _, c := range larikSimpul(d, x.larik) {
					v := teksSimpul(c, "Value")
					if periksa && sebar != idSpreadingTanpaGalat && (v == "" || v == "0") {
						galat = true
					}
					*x.tujuan = tambahPerMataUang(*x.tujuan, teksSimpul(c, "Currency"), "", angka(v))
				}
			}
		}
	}
	return galat, jenis, grup
}

// --- pohon generik --------------------------------------------------------

// teksSimpul - nilai teks satu kunci simpul; yang bukan teks → "".
func teksSimpul(s map[string]any, k string) string {
	v, _ := s[k].(string)
	return v
}

// larikSimpul - larik anak satu kunci, sebagai SIMPUL YANG SAMA (bukan
// salinan): penulisan ke elemennya mengubah pohon. Bentuk JSON (`[]any`)
// dinormalkan sekali ke `[]map[string]any` dan dipasang balik.
func larikSimpul(s map[string]any, k string) []map[string]any {
	switch v := s[k].(type) {
	case []map[string]any:
		return v
	case []any:
		out := make([]map[string]any, 0, len(v))
		for _, x := range v {
			if m, ok := x.(map[string]any); ok {
				out = append(out, m)
			}
		}
		s[k] = out
		return out
	}
	return nil
}

// salinPohon - salinan dalam; masukan tidak berubah.
func salinPohon(xs []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(xs))
	for _, x := range xs {
		out = append(out, salinSimpul(x))
	}
	return out
}

func salinSimpul(s map[string]any) map[string]any {
	out := make(map[string]any, len(s))
	for k, v := range s {
		switch t := v.(type) {
		case []any:
			l := make([]map[string]any, 0, len(t))
			for _, x := range t {
				if m, ok := x.(map[string]any); ok {
					l = append(l, salinSimpul(m))
				}
			}
			out[k] = l
		case []map[string]any:
			out[k] = salinPohon(t)
		default:
			out[k] = v
		}
	}
	return out
}

func keSimpul(xs []NilaiMataUang) []map[string]any {
	out := make([]map[string]any, 0, len(xs))
	for _, x := range xs {
		out = append(out, map[string]any{"Currency": x.Currency, "Value": x.Value})
	}
	return out
}
