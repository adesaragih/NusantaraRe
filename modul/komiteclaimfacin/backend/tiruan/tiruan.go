// Package tiruan adalah gudang, acuan, dan KONTRAK PALSU Claim Fac In dalam memori untuk uji Komite Claim Fac In (seam
// HTTP handlers -> services, dan server tiruan uji manual). Ia meniru apa yang Oracle simpan: tangga dan kepala kasus,
// tabel warisan, outbox, serta kasus klaim induk di balik `kontrak.KlaimFacInKomite` (daftar putih ditegakkan seperti
// penyedia aslinya). Transaksi gagal memulihkan seluruh keadaan, termasuk klaim induk. Fixture berawalan `UJI-`.
package tiruan

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/repository"
)

// Efek - satu baris outbox tiruan.
type Efek struct {
	Jenis, Rujukan, Muatan string
}

// Gudang - penyimpanan tiruan.
type Gudang struct {
	mu     sync.Mutex
	Kasus  map[string]models.Kasus
	Lini   map[string]string
	Tangga map[string][]models.Anggota
	// Posisi - T_WORK_CLAIM.POSITION kasus komite (KomiteID tingkat berjalan; kosong = selesai).
	Posisi map[string]string
	// Tabel warisan dan outbox.
	OS         []models.BarisOS
	JSONKlaim  map[string][2]string
	Log        []models.LogLayanan
	Riwayat    []models.RiwayatAkseptasi
	SubProgres map[string]string
	KlaimTolak []models.KlaimDitolak
	Efek       []Efek
	urutTangga int
	// Urut - penghitung nomor akseptasi (jenis -> urut terakhir); KodeProduksi - KODE_PRODUKSI NONLIFE.
	Urut         map[string]int
	KodeProduksi string
	// Klaim - kontrak palsu (keadaannya ikut dipulihkan saat transaksi gagal).
	Klaim *KlaimPalsu
}

// Baru membuat gudang kosong beserta kontrak palsunya.
func Baru() *Gudang {
	return &Gudang{Kasus: map[string]models.Kasus{}, Lini: map[string]string{}, Tangga: map[string][]models.Anggota{},
		Posisi: map[string]string{}, JSONKlaim: map[string][2]string{}, SubProgres: map[string]string{},
		Urut: map[string]int{}, KodeProduksi: "UJI-", Klaim: KlaimBaru()}
}

type cadangan struct {
	kasus      map[string]models.Kasus
	tangga     map[string][]models.Anggota
	posisi     map[string]string
	os         []models.BarisOS
	json       map[string][2]string
	log        []models.LogLayanan
	riw        []models.RiwayatAkseptasi
	sub        map[string]string
	tolak      []models.KlaimDitolak
	efek       []Efek
	urut       map[string]int
	urutTangga int
	klaim      map[string]*klaimTiruan
}

func salinPeta[V any](m map[string]V) map[string]V {
	out := make(map[string]V, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func (g *Gudang) salin() cadangan {
	c := cadangan{kasus: salinPeta(g.Kasus), tangga: map[string][]models.Anggota{}, posisi: salinPeta(g.Posisi),
		os: slices.Clone(g.OS), json: salinPeta(g.JSONKlaim), log: slices.Clone(g.Log), riw: slices.Clone(g.Riwayat),
		sub: salinPeta(g.SubProgres), tolak: slices.Clone(g.KlaimTolak), efek: slices.Clone(g.Efek),
		urut: salinPeta(g.Urut), urutTangga: g.urutTangga}
	for k, v := range g.Tangga {
		c.tangga[k] = slices.Clone(v)
	}
	if g.Klaim != nil {
		c.klaim = g.Klaim.salin()
	}
	return c
}

func (g *Gudang) pulihkan(c cadangan) {
	g.Kasus, g.Tangga, g.Posisi, g.OS, g.JSONKlaim, g.Log, g.Riwayat = c.kasus, c.tangga, c.posisi, c.os, c.json, c.log,
		c.riw
	g.SubProgres, g.KlaimTolak, g.Efek, g.Urut, g.urutTangga = c.sub, c.tolak, c.efek, c.urut, c.urutTangga
	if g.Klaim != nil {
		g.Klaim.mu.Lock()
		g.Klaim.klaim = c.klaim
		g.Klaim.mu.Unlock()
	}
}

// Transaksi - tiruan: keadaan dipulihkan bila fn gagal (rollback). `tx` selalu nil.
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	c := g.salin()
	g.mu.Unlock()
	if err := fn(nil); err != nil {
		g.mu.Lock()
		g.pulihkan(c)
		g.mu.Unlock()
		return err
	}
	return nil
}

// Lahirkan - kasus komite KMT- baru seperti `BuatKasusKomite` Claim Fac In (LINI FACIN, COUNT 1, tangga menunggu,
// SUBPROGRESSCLAIM "Waiting Committee").
func (g *Gudang) Lahirkan(id, klaimID, adjID, transfer, pembuat, namaPembuat string, anggota []models.Anggota,
	saat time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()
	t := make([]models.Anggota, 0, len(anggota))
	for i, a := range anggota {
		a.Urut = i + 1
		if a.ID == "" {
			a.ID = fmt.Sprintf("%s-L%d", id, i+1)
		}
		a.Keputusan = models.KeputusanMenunggu
		t = append(t, a)
	}
	g.Kasus[id] = models.Kasus{ID: id, KlaimID: klaimID, AdjustmentID: adjID, TransferType: transfer, Loop: len(t),
		Count: 1, UsulTutup: models.UsulTidak, UsulCadang: models.UsulTidak, Tahap: models.TahapKomite, PembuatID: pembuat,
		PembuatNama: namaPembuat, TglCreate: saat, TglUpdate: saat}
	g.Lini[id] = models.LiniFacIn
	g.Tangga[id] = t
	if len(t) > 0 {
		g.Posisi[id] = t[0].OperatorID
	}
	g.SubProgres[id] = "Waiting Committee"
}

// BacaKasus - lihat `repository.Gudang.BacaKasus`.
func (g *Gudang) BacaKasus(_ context.Context, _ *db.Tx, id string, _ bool) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada || g.Lini[id] != models.LiniFacIn || !strings.HasPrefix(id, models.AwalanKomite) {
		return models.Kasus{}, fmt.Errorf("%w: %q", repository.ErrKasusTidakAda, id)
	}
	return k, nil
}

// BacaTangga - tangga urut jenjang.
func (g *Gudang) BacaTangga(_ context.Context, _ *db.Tx, id string) ([]models.Anggota, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	t := slices.Clone(g.Tangga[id])
	sort.SliceStable(t, func(i, j int) bool { return t[i].Urut < t[j].Urut })
	return t, nil
}

// TambahAnggota - perluasan tangga (ApprovalKomite_Act S7.1).
func (g *Gudang) TambahAnggota(_ context.Context, _ *db.Tx, id string, baru []models.Anggota) ([]models.Anggota, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]models.Anggota, 0, len(baru))
	for _, a := range baru {
		g.urutTangga++
		a.ID = fmt.Sprintf("%s-P%d", id, g.urutTangga)
		a.Keputusan = models.KeputusanMenunggu
		g.Tangga[id] = append(g.Tangga[id], a)
		out = append(out, a)
	}
	return out, nil
}

// TulisAnggota - keputusan satu baris (harus masih menunggu).
func (g *Gudang) TulisAnggota(_ context.Context, _ *db.Tx, id string, u models.UbahAnggota) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, a := range g.Tangga[id] {
		if a.ID != u.ID {
			continue
		}
		if a.Keputusan != models.KeputusanMenunggu {
			return fmt.Errorf("%w: baris tangga %s", repository.ErrKeputusanBersamaan, u.ID)
		}
		a.Keputusan, a.Tanggal = u.Keputusan, models.FormatWaktu(u.Tanggal)
		if u.IsiKomentar {
			a.Komentar = u.Komentar
			if u.Pemutus != "" {
				a.OperatorID = u.Pemutus
			}
		}
		g.Tangga[id][i] = a
		return nil
	}
	return fmt.Errorf("%w: baris tangga %s", repository.ErrKeputusanBersamaan, u.ID)
}

// SimpanKepala - bersyarat KOMITE_COUNT lama.
func (g *Gudang) SimpanKepala(_ context.Context, _ *db.Tx, id string, countLama int, kp models.Kepala) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	if k.Count != countLama {
		return fmt.Errorf("%w: kepala kasus %s", repository.ErrKeputusanBersamaan, id)
	}
	k.Count, k.Loop, k.AcceptStatus, k.UsulTutup, k.UsulCadang = kp.Count, kp.Loop, kp.AcceptStatus, kp.UsulTutup,
		kp.UsulCadang
	g.Kasus[id] = k
	return nil
}

// TutupKasus - selesai (Resolved-Completed, POSITION kosong) atau POSITION tingkat berikut.
func (g *Gudang) TutupKasus(_ context.Context, _ *db.Tx, id string, selesai bool, posisi string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := g.Kasus[id]
	if k.Tertutup() {
		return fmt.Errorf("%w: work object %s", repository.ErrKeputusanBersamaan, id)
	}
	k.TglUpdate = saat
	if selesai {
		k.StatusWork, posisi = models.StatusSelesai, ""
	}
	g.Kasus[id] = k
	g.Posisi[id] = posisi
	return nil
}

// DaftarKerja - baris tingkat berjalan yang menunggu, akun atau workbasket pelaku (sudah PeranKerja).
func (g *Gudang) DaftarKerja(_ context.Context, akun string, peran []string) ([]models.BarisKerja, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.BarisKerja{}
	for id, k := range g.Kasus {
		if k.Tertutup() || g.Lini[id] != models.LiniFacIn {
			continue
		}
		for _, a := range g.Tangga[id] {
			if a.Urut != k.Count || a.Keputusan != models.KeputusanMenunggu {
				continue
			}
			if a.OperatorID != akun && !slices.Contains(peran, a.OperatorID) {
				continue
			}
			out = append(out, models.BarisKerja{KasusID: id, KlaimID: k.KlaimID, Jenis: models.JudulTransfer[k.TransferType],
				Tingkat: a.Urut, Count: k.Count, Loop: k.Loop, Jabatan: a.Jabatan, TglUpdate: k.TglUpdate})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].KasusID > out[j].KasusID })
	return out, nil
}

// UrutNomorAkseptasi - penghitung tiruan (periode tetap 10.2026).
func (g *Gudang) UrutNomorAkseptasi(_ context.Context, _ *db.Tx, _ time.Time) (models.BahanNomor, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	jenis := g.KodeProduksi + models.HurufAkseptasi
	g.Urut[jenis]++
	return models.BahanNomor{Jenis: jenis, MMYYYY: "10.2026", Urut: g.Urut[jenis]}, nil
}

// SisipOS - baris OS_AKSEPTASI_KLAIM.
func (g *Gudang) SisipOS(_ context.Context, _ *db.Tx, b models.BarisOS, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !json.Valid([]byte(b.DataJSON)) { // CHECK `DATA_JSON IS JSON`
		return fmt.Errorf("tiruan: DATA_JSON bukan JSON: %q", b.DataJSON)
	}
	g.OS = append(g.OS, b)
	return nil
}

// SalinJSONKlaim - INSERT bila IDPEGA belum ada.
func (g *Gudang) SalinJSONKlaim(_ context.Context, _ *db.Tx, idPega, noKlaim, noPolis string, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.JSONKlaim[idPega]; !ada {
		g.JSONKlaim[idPega] = [2]string{noKlaim, noPolis}
	}
	return nil
}

// CatatLogLayanan - MONITORING_KLAIM_LOG.
func (g *Gudang) CatatLogLayanan(_ context.Context, _ *db.Tx, l models.LogLayanan, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Log = append(g.Log, l)
	return nil
}

// CatatRiwayatAkseptasi - HISTORYAKSEPTASIPEGA.
func (g *Gudang) CatatRiwayatAkseptasi(_ context.Context, _ *db.Tx, r models.RiwayatAkseptasi, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Riwayat = append(g.Riwayat, r)
	return nil
}

// UbahSubProgres - SUBPROGRESSCLAIM.POSITION2 kasus komite.
func (g *Gudang) UbahSubProgres(_ context.Context, _ *db.Tx, komiteID, posisi string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.SubProgres[komiteID]; ada {
		g.SubProgres[komiteID] = posisi
	}
	return nil
}

// SisipKlaimDitolak - CLAIMREJECTED.
func (g *Gudang) SisipKlaimDitolak(_ context.Context, _ *db.Tx, k models.KlaimDitolak) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.KlaimTolak = append(g.KlaimTolak, k)
	return nil
}

// AntreEfek - outbox tiruan.
func (g *Gudang) AntreEfek(_ context.Context, _ *db.Tx, jenis, rujukan, muatan string, _ time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Efek = append(g.Efek, Efek{Jenis: jenis, Rujukan: rujukan, Muatan: muatan})
	return fmt.Sprintf("UJI-EFEK-%d", len(g.Efek)), nil
}

// ---------------------------------------------------------------- acuan

// Acuan - acuan tiruan.
type Acuan struct {
	Roster    []models.AnggotaRoster
	Kurs      map[string]string
	JenisReas map[string]string
	Bank      map[string]string
	Nama      map[string]string
	Email     map[string]string
	Anggota   map[string][]string
	Ceding    map[string]string
	// Konversi - status konversi per nomor akseptasi tanpa titik (getStatusKonversi_Act, "1" = sudah di TRLOSS_DETAIL_T).
	Konversi map[string]string
}

// AcuanBaru - acuan kosong.
func AcuanBaru() *Acuan {
	return &Acuan{Kurs: map[string]string{}, JenisReas: map[string]string{}, Bank: map[string]string{},
		Nama: map[string]string{}, Email: map[string]string{}, Anggota: map[string][]string{}, Ceding: map[string]string{},
		Konversi: map[string]string{}}
}

// StatusKonversi - lihat `repository.Acuan.StatusKonversi` (tiruan: peta Konversi).
func (a *Acuan) StatusKonversi(_ context.Context, noAksep string) (string, error) {
	return a.Konversi[noAksep], nil
}

// RosterKomite - roster FACIN aktif.
func (a *Acuan) RosterKomite(context.Context) ([]models.AnggotaRoster, error) { return a.Roster, nil }

// KursStandar - kurs tiruan.
func (a *Acuan) KursStandar(_ context.Context, cur string) (string, error) { return a.Kurs[cur], nil }

// NamaJenisReas - REINSURANCETYPE tiruan.
func (a *Acuan) NamaJenisReas(context.Context) (map[string]string, error) { return a.JenisReas, nil }

// IDBankRekening - kunci "bank|cabang|akun".
func (a *Acuan) IDBankRekening(_ context.Context, bank, cabang, akun string) (string, error) {
	return a.Bank[bank+"|"+cabang+"|"+akun], nil
}

// EmailCeding - email ceding tiruan.
func (a *Acuan) EmailCeding(_ context.Context, ceding string) (string, error) {
	return a.Ceding[ceding], nil
}

// NamaPelaku - nama akun; tanpa nama = ID akun.
func (a *Acuan) NamaPelaku(_ context.Context, akun string) (string, error) {
	if n := a.Nama[akun]; n != "" {
		return n, nil
	}
	return akun, nil
}

// EmailPelaku - email akun.
func (a *Acuan) EmailPelaku(_ context.Context, akun string) (string, error) {
	return a.Email[akun], nil
}

// EmailAnggotaWorkbasket - email anggota workbasket.
func (a *Acuan) EmailAnggotaWorkbasket(_ context.Context, wb string) ([]string, error) {
	return a.Anggota[wb], nil
}

// ---------------------------------------------------------------- kontrak palsu Claim Fac In

type klaimTiruan struct {
	nilai   map[string]string
	daftar  map[string][]map[string]string
	status  string
	riwayat []kontrak.RiwayatKlaimFacIn
}

func (k *klaimTiruan) salin() *klaimTiruan {
	c := &klaimTiruan{nilai: salinPeta(k.nilai), daftar: map[string][]map[string]string{}, status: k.status,
		riwayat: slices.Clone(k.riwayat)}
	for j, rows := range k.daftar {
		s := make([]map[string]string, 0, len(rows))
		for _, b := range rows {
			s = append(s, salinPeta(b))
		}
		c.daftar[j] = s
	}
	return c
}

// KlaimPalsu memenuhi `kontrak.KlaimFacInKomite` di atas kasus klaim dalam memori. Baris adjustment dikenali lewat
// properti "ID" di pohon `ClaimData.ObjectList(o).ObjectItemList(i).Adjustment`.
type KlaimPalsu struct {
	mu    sync.Mutex
	klaim map[string]*klaimTiruan
	// Dikunci - klaim yang pernah dikunci.
	Dikunci []string
}

var _ kontrak.KlaimFacInKomite = (*KlaimPalsu)(nil)

// KlaimBaru - kontrak palsu kosong.
func KlaimBaru() *KlaimPalsu { return &KlaimPalsu{klaim: map[string]*klaimTiruan{}} }

func (p *KlaimPalsu) salin() map[string]*klaimTiruan {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := map[string]*klaimTiruan{}
	for id, k := range p.klaim {
		out[id] = k.salin()
	}
	return out
}

// Setel - fixture kasus klaim induk.
func (p *KlaimPalsu) Setel(id string, nilai map[string]string, daftar map[string][]map[string]string) {
	k := (&klaimTiruan{nilai: nilai, daftar: daftar}).salin()
	k.nilai["pyID"] = id
	p.mu.Lock()
	defer p.mu.Unlock()
	p.klaim[id] = k
}

// Nilai / Daftar / Status / Riwayat - keadaan klaim induk (asersi uji).
func (p *KlaimPalsu) Nilai(id string) map[string]string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return salinPeta(p.klaim[id].nilai)
}

func (p *KlaimPalsu) Daftar(id, jalur string) []map[string]string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.klaim[id].salin().daftar[jalur]
}

func (p *KlaimPalsu) Status(id string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.klaim[id].status
}

func (p *KlaimPalsu) Riwayat(id string) []kontrak.RiwayatKlaimFacIn {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Clone(p.klaim[id].riwayat)
}

func posisi(k *klaimTiruan, adjID string) (o, i, a int) {
	if adjID == "" {
		return 0, 0, 0
	}
	for on := range k.daftar[models.DaftarObjek] {
		for in := range k.daftar[models.DaftarItem(on+1)] {
			for an, b := range k.daftar[models.DaftarAdj(on+1, in+1)] {
				if b["ID"] == adjID {
					return on + 1, in + 1, an + 1
				}
			}
		}
	}
	return 0, 0, 0
}

// BacaKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) BacaKlaimFacIn(_ context.Context, _ *db.Tx, klaimID, adjID string) (kontrak.KlaimFacIn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ada := p.klaim[klaimID]
	if !ada {
		return kontrak.KlaimFacIn{}, kontrak.ErrKlaimFacInTidakAda
	}
	o, i, a := posisi(k, adjID)
	if adjID != "" && a == 0 {
		return kontrak.KlaimFacIn{}, kontrak.ErrAdjustmentFacInTidakAda
	}
	c := k.salin()
	return kontrak.KlaimFacIn{Nilai: c.nilai, Daftar: c.daftar, Objek: o, Item: i, Adjustment: a, Tertutup: k.status != ""}, nil
}

// KunciKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) KunciKlaimFacIn(_ context.Context, _ *db.Tx, klaimID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ada := p.klaim[klaimID]
	if !ada {
		return kontrak.ErrKlaimFacInTidakAda
	}
	if k.status != "" {
		return kontrak.ErrKlaimFacInTertutup
	}
	p.Dikunci = append(p.Dikunci, klaimID)
	return nil
}

func putih(nama string, ubah, daftar map[string]string) error {
	for j := range ubah {
		if _, ok := daftar[j]; !ok {
			return fmt.Errorf("%w: %s %q", kontrak.ErrUbahanKlaimFacInTidakSah, nama, j)
		}
	}
	return nil
}

func setel(b map[string]string, u map[string]string) {
	for k, v := range u {
		if v == "" {
			delete(b, k)
		} else {
			b[k] = v
		}
	}
}

// TulisBalikKlaimFacIn - lihat `kontrak.KlaimFacInKomite` (daftar putih ditegakkan seperti penyedia aslinya).
func (p *KlaimPalsu) TulisBalikKlaimFacIn(_ context.Context, _ *db.Tx, klaimID, adjID string,
	u kontrak.UbahanKlaimFacIn) error {
	for _, x := range []struct {
		n    string
		u, d map[string]string
	}{{"header", u.Header, kontrak.JalurHeaderKomiteFacIn}, {"adjustment", u.Adjustment, kontrak.PropAdjustmentKomiteFacIn},
		{"objek", u.Objek, kontrak.PropObjekKomiteFacIn}, {"item", u.Item, kontrak.PropItemKomiteFacIn}} {
		if err := putih(x.n, x.u, x.d); err != nil {
			return err
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ada := p.klaim[klaimID]
	if !ada {
		return kontrak.ErrKlaimFacInTidakAda
	}
	if k.status != "" {
		return kontrak.ErrKlaimFacInTertutup
	}
	if len(u.Adjustment) > 0 || len(u.Objek) > 0 || len(u.Item) > 0 {
		o, i, a := posisi(k, adjID)
		if a == 0 {
			return kontrak.ErrAdjustmentFacInTidakAda
		}
		setel(k.daftar[models.DaftarObjek][o-1], u.Objek)
		setel(k.daftar[models.DaftarItem(o)][i-1], u.Item)
		setel(k.daftar[models.DaftarAdj(o, i)][a-1], u.Adjustment)
	}
	setel(k.nilai, u.Header)
	k.riwayat = append(k.riwayat, u.Riwayat...)
	return nil
}

// TutupKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (p *KlaimPalsu) TutupKlaimFacIn(_ context.Context, _ *db.Tx, klaimID, _, status string, _ time.Time) error {
	if status != kontrak.StatusKlaimDitolak && status != kontrak.StatusKlaimSelesai {
		return fmt.Errorf("%w: status %q", kontrak.ErrUbahanKlaimFacInTidakSah, status)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ada := p.klaim[klaimID]
	if !ada {
		return kontrak.ErrKlaimFacInTidakAda
	}
	if k.status != "" {
		return kontrak.ErrKlaimFacInTertutup
	}
	k.status = status
	return nil
}
