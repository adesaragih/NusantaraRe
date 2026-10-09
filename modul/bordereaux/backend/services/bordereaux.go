package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/repository"
)

// Batas.
const (
	UkuranHalaman  = 50
	UkuranMaksimum = 200
	BatasTreaty    = 500
	BatasCedant    = 20
	BatasKomentar  = 2000
	BatasReff      = 1000
	cobaIDMaksimum = 1000
)

// lokasiJakarta - Asia/Jakarta (WIB, tanpa musim panas); FixedZone supaya tidak bergantung tzdata mesin.
var lokasiJakarta = time.FixedZone("WIB", 7*3600)

func sekarangJakarta() time.Time { return time.Now().In(lokasiJakarta) }

func milidetik(t time.Time) string { return fmt.Sprintf("%03d", t.Nanosecond()/int(time.Millisecond)) }

// IDBerkas - `"BDX-" + FormatDateTime(now, "yyyy.MM.dd.ssSSS", "Asia/Jakarta")` (`InputNew` langkah 3, `BdxSave_Act`
// langkah 2) - format XML diikuti apa adanya (keputusan work owner); bentrok ditangani pemanggil.
func IDBerkas(t time.Time) string {
	t = t.In(lokasiJakarta)
	return "BDX-" + t.Format("2006.01.02.05") + milidetik(t)
}

// IDDetail - `"BDX_DTL-" + CurrentDate("yyyyMMddhhmmssSSS", "Asia/Jakarta")` (`SaveDetailData`; `hh` 12 jam seperti XML).
func IDDetail(t time.Time) string {
	t = t.In(lokasiJakarta)
	return "BDX_DTL-" + t.Format("20060102030405") + milidetik(t)
}

// BarisDaftar - satu baris layar daftar beserta hak aktor atasnya.
type BarisDaftar struct {
	models.Header
	Hak Hak `json:"hak"`
}

// Halaman - satu halaman daftar.
type Halaman struct {
	Daftar  []BarisDaftar `json:"daftar"`
	Total   int           `json:"total"`
	Halaman int           `json:"halaman"`
	Ukuran  int           `json:"ukuran"`
}

func tanggalSah(s string) bool {
	if s == "" {
		return true
	}
	_, err := time.Parse("02-01-2006", s)
	return err == nil && len(s) == 10
}

// Daftar - satu halaman daftar (`InboxBordereaux_RD`).
func (l *Layanan) Daftar(ctx context.Context, a Aktor, f models.Filter, halaman, ukuran int) (Halaman, error) {
	f.ReportStart, f.ReportEnd = strings.TrimSpace(f.ReportStart), strings.TrimSpace(f.ReportEnd)
	if !tanggalSah(f.ReportStart) || !tanggalSah(f.ReportEnd) {
		return Halaman{}, tolak("Bordereaux Report Start/End must be DD-MM-YYYY")
	}
	if ukuran <= 0 {
		ukuran = UkuranHalaman
	}
	if ukuran > UkuranMaksimum {
		ukuran = UkuranMaksimum
	}
	if halaman <= 0 {
		halaman = 1
	}
	hs, total, err := l.gudang.Daftar(ctx, f, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return Halaman{}, err
	}
	d := make([]BarisDaftar, 0, len(hs))
	for _, h := range hs {
		d = append(d, BarisDaftar{Header: h, Hak: HakAtas(a, h)})
	}
	return Halaman{Daftar: d, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}

// Rincian - isi form satu berkas.
type Rincian struct {
	Header    models.Header      `json:"header"`
	Baris     []models.Baris     `json:"baris"`
	Ringkasan []models.Ringkasan `json:"ringkasan"`
	Riwayat   []models.Riwayat   `json:"riwayat"`
	Hak       Hak                `json:"hak"`
}

// Buka - header, detail, ringkasan (tab Summary), riwayat, dan hak aktor (`OpenDataBordereaux_Act`, tanpa JSON).
func (l *Layanan) Buka(ctx context.Context, a Aktor, id string) (Rincian, error) {
	h, err := l.gudang.AmbilHeader(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return Rincian{}, ErrTidakAda
	}
	if err != nil {
		return Rincian{}, err
	}
	r := Rincian{Header: h, Baris: []models.Baris{}, Ringkasan: []models.Ringkasan{}, Hak: HakAtas(a, h)}
	r.Hak.Lampiran = BolehLampiran(a, h)
	if k, e := models.CariKombinasi(h.Type, h.TypeBusiness); e == nil {
		if r.Baris, err = l.gudang.Detail(ctx, nil, k, h.BdxID); err != nil {
			return Rincian{}, err
		}
		r.Ringkasan = HitungRingkasan(k, r.Baris)
	}
	if r.Riwayat, err = l.gudang.Riwayat(ctx, h.BdxID); err != nil {
		return Rincian{}, err
	}
	return r, nil
}

// kolomRingkasan - kolom Reinsurer dan RNM yang dijumlah tab Summary per Type (`CountSummaryBdx` langkah 7).
func kolomRingkasan(k models.KombinasiBdx) (reins, rnm string) {
	ada := map[string]bool{}
	for _, c := range k.Kolom {
		ada[c.Kolom] = true
	}
	calon := map[string][][2]string{
		models.TypePremium:     {{"PREMIUM_REINSURER", "PREMIUM_RNM"}},
		models.TypeClaim:       {{"PAID_CLAIMS_REINSURER", "PAID_CLAIMS_RNM"}, {"PAID_CLAIM_REINSURER", "PAID_CLAIM_RNM"}},
		models.TypeSubrogation: {{"PAID_CLAIMS_SUBROGATION_REINSURER", "PAID_CLAIMS_SUBROGATION_RNM"}},
	}[k.Type]
	for _, c := range calon {
		if ada[c[0]] && ada[c[1]] {
			return c[0], c[1]
		}
	}
	return "", ""
}

// HitungRingkasan - total Reinsurer dan RNM per mata uang, eksak (apd), urut kemunculan mata uang.
func HitungRingkasan(k models.KombinasiBdx, baris []models.Baris) []models.Ringkasan {
	reins, rnm := kolomRingkasan(k)
	hasil := []models.Ringkasan{}
	if reins == "" {
		return hasil
	}
	type jumlah struct{ reins, rnm apd.Decimal }
	urut := []string{}
	per := map[string]*jumlah{}
	ctx := utils.DecimalContext()
	tambah := func(tujuan *apd.Decimal, v string) {
		if v == "" {
			return
		}
		d, err := utils.ParseDecimal(v)
		if err != nil {
			return
		}
		_, _ = ctx.Add(tujuan, tujuan, d)
	}
	for _, b := range baris {
		cur := strings.ToUpper(strings.TrimSpace(b["CURRENCY"]))
		j, ok := per[cur]
		if !ok {
			j = &jumlah{}
			per[cur] = j
			urut = append(urut, cur)
		}
		tambah(&j.reins, b[reins])
		tambah(&j.rnm, b[rnm])
	}
	for _, cur := range urut {
		j := per[cur]
		hasil = append(hasil, models.Ringkasan{Currency: cur, Reinsurer: utils.FormatDecimal(&j.reins), RNM: utils.FormatDecimal(&j.rnm)})
	}
	return hasil
}

// Pratinjau - hasil Upload CSV untuk Type dan Business terpilih (`UploadCSVBordereaux_Act` + 29 `Mapping*Bdx*`).
type Pratinjau struct {
	Baris     []models.Baris     `json:"baris"`
	Ringkasan []models.Ringkasan `json:"ringkasan"`
}

// UnggahCSV membaca berkas CSV.
func (l *Layanan) UnggahCSV(_ context.Context, tipe, bisnis, teks string) (Pratinjau, error) {
	tipe, bisnis = normalType(tipe, bisnis)
	k, err := models.CariKombinasi(tipe, bisnis)
	if err != nil {
		return Pratinjau{}, tolak("%s", err.Error())
	}
	b, err := PetakanCSV(k, teks)
	if err != nil {
		return Pratinjau{}, err
	}
	return Pratinjau{Baris: b, Ringkasan: HitungRingkasan(k, b)}, nil
}

// normalType - huruf besar; SUBROGATION selalu BONDING (`SetSubrogation_ACT` langkah 1).
func normalType(tipe, bisnis string) (string, string) {
	tipe, bisnis = strings.ToUpper(strings.TrimSpace(tipe)), strings.ToUpper(strings.TrimSpace(bisnis))
	if tipe == models.TypeSubrogation {
		bisnis = models.BusinessBonding
	}
	return tipe, bisnis
}

// PermintaanSimpan - isian form Save (`BdxSave_Act`). BdxID kosong = berkas baru.
type PermintaanSimpan struct {
	BdxID       string         `json:"bdxId"`
	Type        string         `json:"type"`
	Business    string         `json:"business"`
	MasterID    string         `json:"masterId"`
	ReportStart string         `json:"reportStart"`
	ReportEnd   string         `json:"reportEnd"`
	ReffNoSOA   string         `json:"reffNoSoa"`
	ReffNoBDX   string         `json:"reffNoBdx"`
	Baris       []models.Baris `json:"baris"`
}

// periksaHeader menilai isian header; seluruh kesalahan dikumpulkan.
func (l *Layanan) periksaHeader(ctx context.Context, p *PermintaanSimpan) (models.KombinasiBdx, models.MasterTreaty, error) {
	var (
		pesan []string
		k     models.KombinasiBdx
		t     models.MasterTreaty
	)
	p.Type, p.Business = normalType(p.Type, p.Business)
	p.ReportStart, p.ReportEnd = strings.TrimSpace(p.ReportStart), strings.TrimSpace(p.ReportEnd)
	p.ReffNoSOA, p.ReffNoBDX, p.MasterID = strings.TrimSpace(p.ReffNoSOA), strings.TrimSpace(p.ReffNoBDX), strings.TrimSpace(p.MasterID)
	switch {
	case p.Type == "":
		pesan = append(pesan, "Type is required")
	case p.Business == "":
		pesan = append(pesan, "Type Business is required")
	default:
		var err error
		if k, err = models.CariKombinasi(p.Type, p.Business); err != nil {
			pesan = append(pesan, err.Error())
		}
	}
	if p.MasterID == "" {
		pesan = append(pesan, "Choose a Master Treaty")
	} else {
		var err error
		t, err = l.gudang.AmbilTreaty(ctx, p.MasterID)
		if errors.Is(err, repository.ErrTidakAda) {
			pesan = append(pesan, fmt.Sprintf("Master Treaty %s is not found", p.MasterID))
		} else if err != nil {
			return k, t, err
		}
	}
	switch {
	case p.ReportStart == "" || p.ReportEnd == "":
		pesan = append(pesan, "Bordereaux Report Start and End are required")
	case !tanggalSah(p.ReportStart) || !tanggalSah(p.ReportEnd):
		pesan = append(pesan, "Bordereaux Report Start/End must be DD-MM-YYYY")
	case waktuKanonik(p.ReportEnd).Before(waktuKanonik(p.ReportStart)):
		pesan = append(pesan, "Bordereaux Report End is before Bordereaux Report Start")
	}
	if utf8.RuneCountInString(p.ReffNoSOA) > BatasReff || utf8.RuneCountInString(p.ReffNoBDX) > BatasReff {
		pesan = append(pesan, fmt.Sprintf("Reff No of SOA / Bordereaux is longer than %d characters", BatasReff))
	}
	if len(p.Baris) > MaksBarisCSV {
		pesan = append(pesan, fmt.Sprintf("More than %d detail rows", MaksBarisCSV))
	}
	if len(pesan) == 0 && k.Tabel != "" {
		if err := PeriksaBaris(k, p.Baris); err != nil {
			return k, t, err
		}
	}
	if len(pesan) > 0 {
		return k, t, GalatCSV{pesan}
	}
	return k, t, nil
}

// Simpan - Save: header ke BORDEREAUX, detail ke tabel kombinasinya - SETIAP Save, tidak menunggu Resolve-Complete
// dan tanpa JSON (keputusan work owner). Menjawab BDX_ID.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, p PermintaanSimpan) (string, error) {
	if a.AkunID == "" {
		return "", dilarang("Log in to save")
	}
	k, t, err := l.periksaHeader(ctx, &p)
	if err != nil {
		return "", err
	}
	h := models.Header{BdxID: strings.TrimSpace(p.BdxID), Type: p.Type, TypeBusiness: p.Business, MasterID: t.ID,
		CedingID: t.CedingID, CedingName: t.CedingName, SobID: t.SobID, SobName: t.SobName, TreatyName: t.ContractName,
		ReportStart: p.ReportStart, ReportEnd: p.ReportEnd, ReffNoSOA: p.ReffNoSOA, ReffNoBDX: p.ReffNoBDX}
	baru := h.BdxID == ""
	if baru {
		if !a.BolehBuat() {
			return "", dilarang("Your access to the Bordereaux menu is View only")
		}
		h.UserInput, h.Position = a.AkunID, a.AkunID
	} else {
		lama, err := l.gudang.AmbilHeader(ctx, nil, h.BdxID)
		if errors.Is(err, repository.ErrTidakAda) {
			return "", ErrTidakAda
		}
		if err != nil {
			return "", err
		}
		if !HakAtas(a, lama).Ubah {
			return "", dilarang("This bordereaux cannot be edited by you at its current position")
		}
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if baru {
			jam := l.sekarang()
			for i := 0; ; i++ {
				if i >= cobaIDMaksimum {
					return fmt.Errorf("services: tidak menemukan BDX_ID kosong")
				}
				h.BdxID = IDBerkas(jam.Add(time.Duration(i) * time.Millisecond))
				ada, err := l.gudang.AdaBerkas(ctx, tx, h.BdxID)
				if err != nil {
					return err
				}
				if !ada {
					break
				}
			}
			if err := l.gudang.SisipHeader(ctx, tx, h); err != nil {
				return err
			}
		} else if err := l.gudang.UbahHeader(ctx, tx, h); err != nil {
			return err
		}
		// Detail lama di SEMUA tabel kombinasi (Type/Business mungkin berubah), lalu baris baru. Pega menghapus
		// tabel yang salah untuk klaim Marine Cargo/Hull dan memakai BDX_ID di tabel ber-ID_BDX - diperbaiki.
		for _, kk := range models.Kombinasi {
			if _, err := l.gudang.HapusDetail(ctx, tx, kk, h.BdxID); err != nil {
				return err
			}
		}
		jam, geser := l.sekarang(), 0
		for _, b := range p.Baris {
			for {
				if geser >= cobaIDMaksimum+len(p.Baris) {
					return fmt.Errorf("services: tidak menemukan ID detail kosong")
				}
				id := IDDetail(jam.Add(time.Duration(geser) * time.Millisecond))
				geser++
				err := l.gudang.SisipDetail(ctx, tx, k, h.BdxID, id, b)
				if errors.Is(err, repository.ErrIDTerpakai) {
					continue
				}
				if err != nil {
					return err
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return h.BdxID, nil
}

// Submit - tab Submit (`ActionSubmit` + `AkseptasiBdx_DT`):
//   - pembuat (berkas di tangannya): Accept, ke Checker - minimal satu baris detail;
//   - Checker: setuju = Accept, ke Supervisor; tolak = Rejected, kembali ke pembuat;
//   - Supervisor: setuju = Resolve-Complete, POSITION kosong; tolak = Rejected, kembali ke pembuat.
//
// Tolak wajib berkomentar. Setiap langkah dicatat di BORDEREAUX_HISTORY. Berkas yang sudah berpindah sejak dibaca
// ditolak (ErrSudahDiproses).
func (l *Layanan) Submit(ctx context.Context, a Aktor, id string, setuju bool, komentar string) error {
	komentar = strings.TrimSpace(komentar)
	if len(komentar) > BatasKomentar {
		return tolak("Comment is longer than %d characters", BatasKomentar)
	}
	h, err := l.gudang.AmbilHeader(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	if err != nil {
		return err
	}
	hak := HakAtas(a, h)
	var posisi, status string
	switch {
	case hak.Submit:
		k, err := models.CariKombinasi(h.Type, h.TypeBusiness)
		if err != nil {
			return tolak("%s", err.Error())
		}
		b, err := l.gudang.Detail(ctx, nil, k, h.BdxID)
		if err != nil {
			return err
		}
		if len(b) == 0 {
			return tolak("Upload and save at least one detail row before Submit")
		}
		setuju, posisi, status = true, models.PosisiChecker, models.StatusAccept
	case hak.Putuskan:
		if !setuju && komentar == "" {
			return tolak("Comment is required to reject")
		}
		switch {
		case !setuju:
			posisi, status = h.UserInput, models.StatusRejected
		case h.Position == models.PosisiChecker:
			posisi, status = models.PosisiSupervisor, models.StatusAccept
		default:
			posisi, status = "", models.StatusResolveComplete
		}
	default:
		return dilarang("This bordereaux is not waiting for you")
	}
	return l.tx(ctx, func(tx *dbTx) error {
		ok, err := l.gudang.UbahStatus(ctx, tx, h.BdxID, h.Position, h.StatusAksep, posisi, status)
		if err != nil {
			return err
		}
		if !ok {
			return ErrSudahDiproses
		}
		return l.gudang.SisipRiwayat(ctx, tx, h.BdxID, a.AkunID, setuju, komentar)
	})
}

// Hapus - Delete (`ConfrimDelete`): seluruh isi berkas, hanya selama Edit diizinkan.
func (l *Layanan) Hapus(ctx context.Context, a Aktor, id string) error {
	h, err := l.gudang.AmbilHeader(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	if err != nil {
		return err
	}
	if !HakAtas(a, h).Hapus {
		return dilarang("This bordereaux cannot be deleted by you at its current position")
	}
	return l.tx(ctx, func(tx *dbTx) error { return l.gudang.HapusBerkas(ctx, tx, h.BdxID) })
}

// Chart - isi chart daftar: jumlah berkas per Business x Type x Ceding (dirangkum di layar).
func (l *Layanan) Chart(ctx context.Context) ([]models.IrisanChart, error) {
	d, err := l.gudang.Chart(ctx)
	if d == nil {
		d = []models.IrisanChart{}
	}
	return d, err
}

// CariCedant - autocomplete Cedant popup Choose Master Treaty.
func (l *Layanan) CariCedant(ctx context.Context, kata string) ([]models.Cedant, error) {
	return l.gudang.CariCedant(ctx, kata, BatasCedant)
}

// CariTreaty - grid Master Treaty satu cedant; cedant kosong = nol baris.
func (l *Layanan) CariTreaty(ctx context.Context, cedingID string) ([]models.MasterTreaty, error) {
	if strings.TrimSpace(cedingID) == "" {
		return []models.MasterTreaty{}, nil
	}
	return l.gudang.CariTreaty(ctx, strings.TrimSpace(cedingID), BatasTreaty)
}

// Pilihan - nilai radio Type dan Business (pengganti pilihan properti Pega yang tidak ikut diekspor), dan status.
type Pilihan struct {
	Type     []string            `json:"type"`
	Business map[string][]string `json:"business"`
	Status   []string            `json:"status"`
	Kolom    map[string][]Kolom  `json:"kolom"`
	// Templat - per `TYPE|BUSINESS`: kode slot Template Manager dan nama berkas unduhan.
	Templat map[string]Templat `json:"templat"`
	// BolehBuat - tombol Input Data.
	BolehBuat bool `json:"bolehBuat"`
	// CopyOld - tombol Copy Old Data (superadmin).
	CopyOld bool `json:"copyOld"`
}

// Templat - slot Template Manager satu kombinasi.
type Templat struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// AwalanTemplat - awalan kode slot Template Manager (sama dengan `backend.AwalanTemplat`).
const AwalanTemplat = "bordereaux."

// Kolom - satu kolom grid detail (urutan CSV).
type Kolom struct {
	Kolom string `json:"kolom"`
	// Judul - judul kepala templat CSV (pesan validasi).
	Judul string `json:"judul"`
	Jenis string `json:"jenis"`
	// Kepala, Grup, Sembunyi - tampilan kepala grid menurut format Excel bordereaux 2025 (`models.KepalaKolom`,
	// perintah work owner 05-10-2026); Grup kosong = kolom tanpa grup.
	Kepala   string `json:"kepala"`
	Grup     string `json:"grup,omitempty"`
	Sembunyi bool   `json:"sembunyi,omitempty"`
}

// DaftarPilihan - Type, Business per Type, status, kolom grid per kombinasi (`TYPE|BUSINESS`), dan hak Input Data.
func DaftarPilihan(a Aktor) Pilihan {
	p := Pilihan{Type: models.DaftarType, Business: map[string][]string{}, Kolom: map[string][]Kolom{}, Templat: map[string]Templat{},
		Status: []string{models.StatusAccept, models.StatusRejected, models.StatusResolveComplete}, BolehBuat: a.BolehBuat(),
		CopyOld: wajibSuperadmin(a) == nil}
	jenis := map[models.Jenis]string{models.Teks: "teks", models.Angka: "angka", models.Tanggal: "tanggal"}
	for _, t := range models.DaftarType {
		p.Business[t] = models.DaftarBusiness(t)
	}
	for _, k := range models.Kombinasi {
		var ks []Kolom
		for _, c := range k.Kolom {
			kg := models.KepalaKolom(k, c)
			ks = append(ks, Kolom{Kolom: c.Kolom, Judul: c.Judul, Jenis: jenis[c.Jenis], Kepala: kg.Judul, Grup: kg.Grup, Sembunyi: kg.Sembunyi})
		}
		p.Kolom[k.Type+"|"+k.Business] = ks
		p.Templat[k.Type+"|"+k.Business] = Templat{Kode: AwalanTemplat + k.Kode, Nama: k.BerkasTemplat}
	}
	return p
}
