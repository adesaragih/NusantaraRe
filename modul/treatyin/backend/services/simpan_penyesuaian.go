package services

// Tombol tulis layar ADJUSTMENT (EDM) — `Save`, `Submit`, `Actions`,
// `Decline offer` modul Treaty In Adjustment. Rutenya di modul ini sebab
// penulis pendaratannya di sini (`repository/simpan_penyesuaian.go`).
//
// ---------------------------------------------------------------------
// Rantai ekspor (korpus `Treaty In Adjustment`, blok `TreatyMasterInEDM`)
// ---------------------------------------------------------------------
//
//	Save     pra-DT TreatyInAddNew(Status='1'): Position = Admin,
//	         PositionUsername = operator → SaveTreatyIn_EDM_Act
//	Submit   TreatyInSubmitEDM: [2]–[4] (CheckID, CheckError) MATI — nol
//	         validasi; [5] Akseptasi_DT, [6] AddCommentList_Act,
//	         [7] TreatyEDMCalculateDifference bila EDMState 1/2,
//	         [8] SaveTreatyIn_EDM_Act
//	Actions  TreatyInAkseptasiEDM_Act: Akseptasi_DT (pilihan radio) →
//	         AddCommentList_Act → SaveTreatyIn_EDM_Act
//	Decline  TreatyInDeclineConfirmation_postactEDM: [1]–[4] MATI (nol
//	offer    komentar, nol status); [6]–[7] HAPUS FISIK baris EDM
//
// Tangga akseptasinya `Akseptasi_DT` yang SAMA dengan Treaty In (diadu:
// identik di kedua korpus) — `models.LangkahBerikut`, keputusan pemilik
// proses tentang jalur naik dan `PositionUsername` ikut berlaku.
//
// ⭐ Halaman `ActualValue` (8 Oktober 2026, keputusan pemakai) mendarat di
// MASTERID `ID + AkhiranSisiAktual` — termasuk salinan `TreatyIn` →
// `ActualValue` `SaveTreatyIn_EDM_Act` [2]–[4] untuk EDM 1/2 (`halamanAktual`).
//
// ⚠️ Yang SENGAJA tidak dibangun:
//   - `SaveTreatyInDetailEdm_Act` saat Resolve Complete — pemilik proses:
//     data ditarik dari tabel tiap tab, status dari kepalanya.
//   - (8 Oktober 2026: `TreatyRevisionCopyAttachment` KINI dibangun —
//     `salin_lampiran_penyesuaian.go`, dikaitkan sesudah tulisan draf.)
//   - grid Rate of Exchange: prosedur EDM tidak menulis
//     `TREATYEXCHANGEYEARLY`; yang dikirim layar DILAPORKAN tak tersimpan.

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

// SisiKiriman - satu halaman layar Adjustment apa adanya: `medan` (skalar,
// kunci boleh bertitik) dan `larik` (baris, boleh bersarang — pohon yang
// layar gabungkan).
type SisiKiriman struct {
	Medan map[string]any `json:"medan"`
	Larik map[string]any `json:"larik"`
}

// MasukanPenyesuaian - isi layar Adjustment yang tombol Save kirim.
type MasukanPenyesuaian struct {
	ID string `json:"id"`
	// Draf - penyesuaian dari `Choose` yang belum pernah disimpan.
	Draf bool `json:"draf"`
	// Baru - panel New (`TreatyIn`).
	Baru SisiKiriman `json:"baru"`
	// Lama - panel Old (`OLDDATA`); hanya dipakai untuk draf.
	Lama *SisiKiriman `json:"lama"`
}

// MasukanKirimPenyesuaian - Submit atau Actions layar Adjustment.
type MasukanKirimPenyesuaian struct {
	MasukanPenyesuaian
	// Aksi - `submit` atau `akseptasi`.
	Aksi    string `json:"aksi"`
	Pilihan string `json:"pilihan"`
}

// MasukanHapusPenyesuaian - Decline offer layar Adjustment.
type MasukanHapusPenyesuaian struct {
	ID string `json:"id"`
}

// larikKurs - grid Rate of Exchange layar Adjustment: cermin
// `TREATYEXCHANGEYEARLY`, bukan larik dokumen.
const larikKurs = "CurrencyList"

// SimpanPenyesuaian - tombol Save layar Adjustment.
func (l *Layanan) SimpanPenyesuaian(ctx context.Context, p inti.Pelaku, m MasukanPenyesuaian) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	doc, lama, err := l.dokumenPenyesuaian(ctx, m)
	if err != nil {
		return HasilSimpan{}, err
	}
	if teksDok(doc, "StatusAkseptasi") == models.StatusTuntas {
		return HasilSimpan{}, ditolak("Penyesuaian sudah Resolve Complete — tombol Save tidak berlaku.")
	}
	// DT `TreatyInAddNew` [1]–[2]: Position Admin, lalu (Status=1, Position
	// kosong/Admin) PositionUsername = operator.
	doc["Position"] = models.PosisiAdmin
	doc["PositionUsername"] = p.AkunID
	h, err := l.tulisPenyesuaian(ctx, m, doc, lama)
	// TreatyInEDMSetValue [8] — draf saja, sesudah kepala tersimpan.
	return l.salinLampiranDraf(ctx, p, m, doc, h, err)
}

// KirimPenyesuaian - Submit atau Actions layar Adjustment.
func (l *Layanan) KirimPenyesuaian(ctx context.Context, p inti.Pelaku, m MasukanKirimPenyesuaian) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	doc, lama, err := l.dokumenPenyesuaian(ctx, m.MasukanPenyesuaian)
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
			fmt.Sprintf("Penyesuaian ini menunggu di posisi %s, dan akun %s tidak memegang workbasket itu.", posisi, p.AkunID)}
	}
	pilihan := models.PilihAccept
	switch m.Aksi {
	case AksiSubmit:
		// TreatyInSubmitEDM [2]–[4] ber-`//`: nol validasi.
	case AksiAkseptasi:
		pilihan = m.Pilihan
	default:
		return HasilSimpan{}, ditolak(fmt.Sprintf("Aksi %q tidak dikenal.", m.Aksi))
	}
	langkah, err := models.LangkahBerikut(posisi, pilihan, status, teksDok(doc, "RevisionState") == "1")
	switch {
	case errors.Is(err, models.ErrSudahTuntas):
		return HasilSimpan{}, ditolak("Penyesuaian sudah Resolve Complete — tidak dapat dikirim ulang.")
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
	TerapkanLangkah(doc, langkah, pilihan, pemegang)
	tambahKomentar(doc, stempelPega(time.Now()), p.AkunID, teksDok(doc, "StatusAkseptasi"), teksDok(doc, "Comment"))
	// TreatyInSubmitEDM [7] — Value Difference, EDM 1/2.
	if m.Aksi == AksiSubmit {
		if s := teksDok(doc, "EDMState"); s == "1" || s == "2" {
			if err := l.terapkanSelisih(ctx, m.MasukanPenyesuaian, doc, lama); err != nil {
				return HasilSimpan{}, err
			}
		}
	}
	h, err := l.tulisPenyesuaian(ctx, m.MasukanPenyesuaian, doc, lama)
	// TreatyInEDMSetValue [8] — draf saja, sesudah kepala tersimpan.
	return l.salinLampiranDraf(ctx, p, m.MasukanPenyesuaian, doc, h, err)
}

// HapusPenyesuaian - Decline offer layar Adjustment: baris EDM dihapus
// fisik (`TreatyInDeclineConfirmation_postactEDM` [6]–[7]).
//
// ⚠️ Satu pagar yang TIDAK ada di ekspor: penyesuaian `Resolve Complete`
// ditolak. Ekspor hanya bersyarat `ViewState != '1'`, sehingga penyesuaian
// tuntas yang dibuka lewat Edit dapat dihapus — penghapusan fisik yang tidak
// dapat dibatalkan, jadi pagarnya dipasang dan dinyatakan di sini.
func (l *Layanan) HapusPenyesuaian(ctx context.Context, p inti.Pelaku, m MasukanHapusPenyesuaian) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	id := strings.TrimSpace(m.ID)
	kepala, ada, err := l.gudang.BacaKepalaPenyesuaian(ctx, id)
	if err != nil {
		return HasilSimpan{}, err
	}
	if !ada {
		return HasilSimpan{}, fmt.Errorf("%w: %s", ErrKontrakTidakAda, id)
	}
	if teksDok(kepala, "StatusAkseptasi") == models.StatusTuntas {
		return HasilSimpan{}, ditolak("Penyesuaian sudah Resolve Complete — tidak dapat di-Decline.")
	}
	if err := l.gudang.HapusPenyesuaian(ctx, id); err != nil {
		if errors.Is(err, repository.ErrBukanPengenalPenyesuaian) {
			return HasilSimpan{}, ditolak(fmt.Sprintf("%q bukan pengenal penyesuaian.", id))
		}
		return HasilSimpan{}, err
	}
	return HasilSimpan{ID: id, KunciTakTersimpan: []string{}}, nil
}

// dokumenPenyesuaian - "clipboard" Adjustment: kepala `TREATY_IN_EDM` ⊕
// pendaratan sisi New ⊕ kiriman layar. Draf: kiriman layar saja (belum ada
// yang tersimpan), beserta sisi Old-nya.
func (l *Layanan) dokumenPenyesuaian(ctx context.Context, m MasukanPenyesuaian) (doc, lama map[string]any, err error) {
	id := strings.TrimSpace(m.ID)
	if id == "" {
		return nil, nil, fmt.Errorf("%w: pengenal penyesuaian kosong", ErrMasukanTidakSah)
	}
	doc = map[string]any{}
	kiriman := dokDariSisi(m.Baru)
	if m.Draf {
		// Draf: properti milik server SUDAH disusun `Choose` (SusunDraf) —
		// diterima apa adanya, lalu DT tombolnya menimpa yang perlu.
		for k, v := range kiriman {
			doc[k] = v
		}
		lama = map[string]any{}
		if m.Lama != nil {
			lama = dokDariSisi(*m.Lama)
		}
	} else {
		kepala, ada, err := l.gudang.BacaKepalaPenyesuaian(ctx, id)
		if err != nil {
			return nil, nil, err
		}
		if !ada {
			return nil, nil, fmt.Errorf("%w: %s", ErrKontrakTidakAda, id)
		}
		tersimpan, err := l.gudang.BacaDokumenPendaratan(ctx, id)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range tersimpan {
			doc[k] = v
		}
		for k, v := range kepala {
			doc[k] = v
		}
		for k, v := range kiriman {
			if kunciMilikServer[k] || kunciSesi(k) {
				continue
			}
			doc[k] = v
		}
	}
	doc["ID"] = id
	return doc, lama, nil
}

// dokDariSisi - halaman layar → dokumen Pega: kunci bertitik kembali ke
// halaman tertanamnya.
func dokDariSisi(s SisiKiriman) map[string]any {
	doc := map[string]any{}
	for k, v := range s.Medan {
		setelJalurDok(doc, k, v)
	}
	for k, v := range s.Larik {
		setelJalurDok(doc, k, v)
	}
	return doc
}

func setelJalurDok(m map[string]any, kunci string, v any) {
	bagian := strings.SplitN(kunci, ".", 2)
	if len(bagian) == 1 {
		m[kunci] = v
		return
	}
	anak, ok := m[bagian[0]].(map[string]any)
	if !ok {
		anak = map[string]any{}
		m[bagian[0]] = anak
	}
	setelJalurDok(anak, bagian[1], v)
}

// terapkanSelisih - `TreatyEDMCalculateDifference` (rumus yang SAMA dengan
// tombol `Update Value`, `hitung_selisih.go`) atas dokumen yang dikirim.
func (l *Layanan) terapkanSelisih(ctx context.Context, m MasukanPenyesuaian, doc, lama map[string]any) error {
	sisiLama := lama
	if !m.Draf {
		var err error
		if sisiLama, err = l.gudang.BacaDokumenPendaratan(ctx, strings.TrimSpace(m.ID)+repository.AkhiranSisiLama); err != nil {
			return err
		}
	}
	aktual, _ := doc["ActualValue"].(map[string]any)
	h := selisihEDM(MasukanSelisih{
		Aksi: AksiSelisihEDM, Akar: halamanDariDok(doc), Lama: halamanDariDok(sisiLama), Actual: halamanDariDok(aktual),
	})
	doc["ValueDifference"] = dokDariHalaman(h.Selisih)
	if h.SebelumProrata != nil {
		doc["ValueBeforeProrate"] = dokDariHalaman(*h.SebelumProrata)
	}
	if h.FacultativeShareList != nil {
		doc["FacultativeShareList"] = larikDok(h.FacultativeShareList)
	}
	return nil
}

// halamanDariDok - dokumen → halaman rumus: skalar teks ke `Medan`, larik
// baris ke `Larik`; halaman tertanam tidak dibawa.
func halamanDariDok(d map[string]any) HalamanPohon {
	h := HalamanPohon{Medan: map[string]string{}, Larik: map[string][]map[string]any{}}
	for k, v := range d {
		switch x := v.(type) {
		case string:
			h.Medan[k] = x
		case []any:
			baris := make([]map[string]any, 0, len(x))
			for _, e := range x {
				if b, ok := e.(map[string]any); ok {
					baris = append(baris, b)
				}
			}
			h.Larik[k] = baris
		case []map[string]any:
			h.Larik[k] = x
		}
	}
	return h
}

func dokDariHalaman(h HalamanPohon) map[string]any {
	d := map[string]any{}
	for k, v := range h.Medan {
		d[k] = v
	}
	for k, v := range h.Larik {
		d[k] = larikDok(v)
	}
	return d
}

func larikDok(xs []map[string]any) []any {
	out := make([]any, len(xs))
	for i, b := range xs {
		out[i] = b
	}
	return out
}

// tulisPenyesuaian - `SaveTreatyIn_EDM_Act` [5]–[8]: kepala dan pendaratan,
// lalu pesan prosedur.
func (l *Layanan) tulisPenyesuaian(ctx context.Context, m MasukanPenyesuaian, doc, lama map[string]any) (HasilSimpan, error) {
	id := strings.TrimSpace(m.ID)
	takTersimpan := []string{}
	// ⛔ Grid Rate of Exchange bukan larik dokumen; prosedur EDM tidak
	// menulis `TREATYEXCHANGEYEARLY`. Layar mengirimnya HANYA bila berubah —
	// dan perubahan itu dilaporkan, tidak ditelan.
	if _, ada := doc[larikKurs]; ada {
		delete(doc, larikKurs)
		takTersimpan = append(takTersimpan, larikKurs)
	}
	delete(lama, larikKurs)
	// ⭐ Halaman `ActualValue` — MASTERID sendiri (`AkhiranSisiAktual`,
	// keputusan pemakai 8 Oktober 2026), bukan bagian dokumen sisi New.
	aktual, adaAktual := halamanAktual(doc)
	delete(doc, kunciAktual)
	// ⭐ Halaman `ValueBeforeProrate` — MASTERID sendiri, pola yang SAMA.
	//
	// ⛔ Sampai 8 Oktober 2026 ia tertinggal di dokumen sisi New, nol punya
	// kolom, dan karena itu dilaporkan "TIDAK tersimpan" pada setiap Save
	// penyesuaian ber-pro-rate — lalu hilang. Tab
	// `TreatyInTabsNPValueDifference_NoProRate` yang membacanya kosong
	// sesudah kontrak dimuat ulang.
	sebelumProrata, adaSebelumProrata := doc[kunciSebelumProrata].(map[string]any)
	delete(doc, kunciSebelumProrata)
	asing, err := l.kunciTakTersimpan(ctx, doc, repository.KolomKepalaPenyesuaian)
	if err != nil {
		return HasilSimpan{}, err
	}
	takTersimpan = append(takTersimpan, asing...)
	// Yang TIDAK tersimpan dari halaman Actual dilaporkan HANYA bila isinya
	// isian pemakai (EDMState 3). Salinan [2]–[4] identik dengan sisi New,
	// jadi kuncinya sudah dilaporkan sekali di atas.
	if adaAktual && teksDok(doc, "EDMState") == "3" {
		asingAktual, err := l.kunciTakTersimpan(ctx, aktual, nil)
		if err != nil {
			return HasilSimpan{}, err
		}
		for _, k := range asingAktual {
			takTersimpan = append(takTersimpan, kunciAktual+"."+k)
		}
	}
	sort.Strings(takTersimpan)
	r := models.RencanaPenyesuaian{ID: id, Draf: m.Draf, Baru: doc}
	if m.Draf {
		r.Lama = lama
	}
	if adaAktual {
		r.Aktual = aktual
	}
	if adaSebelumProrata {
		r.SebelumProrata = sebelumProrata
	}
	if err := l.gudang.SimpanPenyesuaian(ctx, r); err != nil {
		if errors.Is(err, repository.ErrPenyesuaianSudahAda) {
			return HasilSimpan{}, ditolak(fmt.Sprintf("Pengenal %s sudah dipakai penyesuaian lain — pilih revisi terakhirnya di picker.", id))
		}
		if errors.Is(err, repository.ErrBukanPengenalPenyesuaian) {
			return HasilSimpan{}, ditolak(fmt.Sprintf("%q bukan pengenal penyesuaian.", id))
		}
		return HasilSimpan{}, err
	}
	return HasilSimpan{
		ID: id, Pesan: "Data Sudah Disimpan Dengan ID : " + id,
		Posisi: teksDok(doc, "Position"), Status: teksDok(doc, "StatusAkseptasi"),
		PemegangPosisi: teksDok(doc, "PositionUsername"), KunciTakTersimpan: takTersimpan,
	}, nil
}

// kunciTakTersimpan - properti yang peta tidak kenal (kecuali kolom kepala)
// ditambah yang kolomnya menunggu migrasi.
func (l *Layanan) kunciTakTersimpan(ctx context.Context, doc map[string]any, kepala [][2]string) ([]string, error) {
	diKepala := map[string]bool{}
	for _, k := range kepala {
		diKepala[k[1]] = true
	}
	out := []string{}
	// Properti BERISI saja — lihat `tulis`.
	doc = repository.TanpaKosong(doc)
	for _, ks := range repository.KunciTakTerpetakan(doc) {
		for _, k := range ks {
			if !diKepala[k] {
				out = append(out, k)
			}
		}
	}
	belum, err := l.gudang.KunciBelumTerpasang(ctx, doc)
	if err != nil {
		return nil, err
	}
	return unik(append(out, belum...)), nil
}

// kunciAktual - halaman `TreatyIn.ActualValue`.
const kunciAktual = "ActualValue"

// kunciSebelumProrata - halaman `TreatyIn.ValueBeforeProrate`, salinan
// `ValueDifference` SEBELUM pro-rate (`TreatyEDMProRateCalculation`).
const kunciSebelumProrata = "ValueBeforeProrate"

// halamanAktual - `SaveTreatyIn_EDM_Act` [2]–[4] lalu halaman Actual yang
// didaratkan.
//
//	[2]  EDMState == "3" → lompat ke `jmp`: `ActualValue` apa adanya
//	     (isian cabang Adjust Premium)
//	[2]–[4] selainnya: Page-Copy TreatyIn → TreatyInTemp, buang `.OLDDATA`
//	     dan `.ValueDifference`, Page-Copy TreatyInTemp → `ActualValue`
//
// ⚠️ `ActualValue` di DALAM salinan (halaman Actual lama yang ikut tersalin
// Page-Copy) dibuang: tingkat kedua tidak mendarat — aturan yang sama dengan
// `OLDDATA` di dalam `OLDDATA`.
//
// `ada` palsu = tidak ada halaman Actual sama sekali (EDMState 3 tanpa
// isian) — sisi `#AKTUAL` tidak disentuh.
func halamanAktual(doc map[string]any) (map[string]any, bool) {
	if teksDok(doc, "EDMState") == "3" {
		a, ok := doc[kunciAktual].(map[string]any)
		return a, ok
	}
	salin := make(map[string]any, len(doc))
	for k, v := range doc {
		switch k {
		case "OLDDATA", "ValueDifference", kunciAktual:
			continue
		}
		salin[k] = v
	}
	return salin, true
}
