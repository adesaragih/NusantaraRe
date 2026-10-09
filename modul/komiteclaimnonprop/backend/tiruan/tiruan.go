// Package tiruan adalah gudang, acuan, dan KONTRAK PALSU Claim Non Prop dalam memori untuk uji Komite Claim Non Prop
// (seam HTTP handlers -> services, dan server tiruan uji manual). Ia meniru apa yang Oracle simpan: tangga dan
// kepala kasus, tabel warisan, outbox, serta kasus klaim induk di balik `kontrak.KlaimTreatyNonPropKomite` (daftar putih
// ditegakkan seperti penyedia aslinya). Transaksi gagal memulihkan seluruh keadaan, termasuk klaim induk. Fixture
// berawalan `UJI-`.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
	"nusantarare/modul/komiteclaimnonprop/backend/repository"
)

// Gudang - penyimpanan tiruan.
type Gudang struct {
	mu     sync.Mutex
	Kasus  map[string]models.Kasus
	Lini   map[string]string
	Tangga map[string][]models.Anggota
	// Posisi - T_WORK_CLAIM.POSITION kasus komite (KomiteID tingkat berjalan; kosong = selesai).
	Posisi map[string]string
	// Tabel warisan dan outbox.
	OS        []models.BarisOSAkseptasi
	XOL2      []models.BarisXOL2
	OSSubj    map[string]models.BarisOSSubjectivity
	JSONKlaim map[string][2]string
	Riwayat   []models.RiwayatAkseptasi
	Efek      []string
	// Urut - penghitung nomor akseptasi (jenis -> urut terakhir); KodeProduksi - KODE_PRODUKSI NONLIFE.
	Urut         map[string]int
	KodeProduksi string
	// Klaim - kontrak palsu (keadaannya ikut dipulihkan saat transaksi gagal).
	Klaim *KlaimPalsu
}

// Baru membuat gudang kosong beserta kontrak palsunya.
func Baru() *Gudang {
	return &Gudang{Kasus: map[string]models.Kasus{}, Lini: map[string]string{}, Tangga: map[string][]models.Anggota{},
		Posisi: map[string]string{}, OSSubj: map[string]models.BarisOSSubjectivity{}, JSONKlaim: map[string][2]string{},
		Urut: map[string]int{}, KodeProduksi: "UJI", Klaim: KlaimBaru()}
}

type cadangan struct {
	kasus  map[string]models.Kasus
	tangga map[string][]models.Anggota
	posisi map[string]string
	os     []models.BarisOSAkseptasi
	xol2   []models.BarisXOL2
	osSubj map[string]models.BarisOSSubjectivity
	json   map[string][2]string
	riw    []models.RiwayatAkseptasi
	efek   []string
	urut   map[string]int
	klaim  map[string]*klaimTiruan
}

func (g *Gudang) salin() cadangan {
	c := cadangan{kasus: map[string]models.Kasus{}, tangga: map[string][]models.Anggota{}, posisi: map[string]string{},
		os: append([]models.BarisOSAkseptasi{}, g.OS...), xol2: append([]models.BarisXOL2{}, g.XOL2...),
		osSubj: map[string]models.BarisOSSubjectivity{}, json: map[string][2]string{},
		riw: append([]models.RiwayatAkseptasi{}, g.Riwayat...), efek: append([]string{}, g.Efek...), urut: map[string]int{}}
	for k, v := range g.OSSubj {
		c.osSubj[k] = v
	}
	for k, v := range g.Kasus {
		c.kasus[k] = v
	}
	for k, v := range g.Tangga {
		c.tangga[k] = append([]models.Anggota{}, v...)
	}
	for k, v := range g.JSONKlaim {
		c.json[k] = v
	}
	for k, v := range g.Posisi {
		c.posisi[k] = v
	}
	for k, v := range g.Urut {
		c.urut[k] = v
	}
	if g.Klaim != nil {
		c.klaim = g.Klaim.salin()
	}
	return c
}

func (g *Gudang) pulihkan(c cadangan) {
	g.Kasus, g.Tangga, g.OS, g.JSONKlaim, g.Riwayat, g.Efek, g.Urut = c.kasus, c.tangga, c.os, c.json, c.riw, c.efek,
		c.urut
	g.XOL2, g.OSSubj, g.Posisi = c.xol2, c.osSubj, c.posisi
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

// Lahirkan - kasus komite KMTNP- baru seperti `BuatKasusKomite` Claim Non Prop (LINI NONPROP, COUNT 1, tangga menunggu).
func (g *Gudang) Lahirkan(id, klaimID, adjID, pembuat, namaPembuat string, anggota []models.Anggota, saat time.Time) {
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
	g.Kasus[id] = models.Kasus{ID: id, KlaimID: klaimID, AdjustmentID: adjID, Loop: len(t), Count: 1,
		UsulTutup: models.UsulTidak, UsulCadang: models.UsulTidak, Tahap: models.TahapKomite, PembuatID: pembuat,
		PembuatNama: namaPembuat, TglCreate: saat, TglUpdate: saat}
	g.Lini[id] = models.LiniNonProp
	g.Tangga[id] = t
}

// BacaKasus - lihat `repository.Gudang.BacaKasus`.
func (g *Gudang) BacaKasus(_ context.Context, _ *db.Tx, id string, _ bool) (models.Kasus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada || g.Lini[id] != models.LiniNonProp || !strings.HasPrefix(id, models.AwalanKomite) {
		return models.Kasus{}, fmt.Errorf("%w: %q", repository.ErrKasusTidakAda, id)
	}
	k.Tangga = nil
	return k, nil
}

// BacaTangga - lihat `repository.Gudang.BacaTangga`.
func (g *Gudang) BacaTangga(_ context.Context, _ *db.Tx, id string) ([]models.Anggota, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]models.Anggota{}, g.Tangga[id]...), nil
}

// TulisAnggota - lihat `repository.Gudang.TulisAnggota` (baris harus masih menunggu).
func (g *Gudang) TulisAnggota(_ context.Context, _ *db.Tx, id string, u models.UbahAnggota) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	t := g.Tangga[id]
	for i := range t {
		if t[i].ID != u.ID {
			continue
		}
		if t[i].Keputusan != models.KeputusanMenunggu {
			return fmt.Errorf("%w: baris tangga %s", repository.ErrKeputusanBersamaan, u.ID)
		}
		t[i].Keputusan = u.Keputusan
		t[i].Tanggal = models.FormatWaktu(u.Tanggal)
		if u.IsiKomentar {
			t[i].Komentar = u.Komentar
			if u.Pemutus != "" {
				t[i].OperatorID = u.Pemutus
			}
		}
		return nil
	}
	return fmt.Errorf("%w: baris tangga %s", repository.ErrKeputusanBersamaan, u.ID)
}

// SimpanKepala - lihat `repository.Gudang.SimpanKepala`.
func (g *Gudang) SimpanKepala(_ context.Context, _ *db.Tx, id string, countLama int, kp models.Kepala) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada || k.Count != countLama {
		return fmt.Errorf("%w: kepala kasus %s", repository.ErrKeputusanBersamaan, id)
	}
	k.Count, k.AcceptStatus, k.UsulTutup, k.UsulCadang = kp.Count, kp.AcceptStatus, kp.UsulTutup, kp.UsulCadang
	k.Subjectivity, k.SubjectivityNote = kp.Subjectivity, kp.SubjectivityNote
	g.Kasus[id] = k
	return nil
}

// KomentarAwal - lihat `repository.Gudang.KomentarAwal`: komentar anggota pertama kasus komite lain (paling awal)
// untuk klaim dan akseptasi yang sama.
func (g *Gudang) KomentarAwal(_ context.Context, klaimID, adjID, id string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var calon []models.Kasus
	for kid, k := range g.Kasus {
		if kid != id && k.KlaimID == klaimID && k.AdjustmentID == adjID {
			calon = append(calon, k)
		}
	}
	sort.Slice(calon, func(i, j int) bool { return calon[i].TglCreate.Before(calon[j].TglCreate) })
	if len(calon) == 0 || len(g.Tangga[calon[0].ID]) == 0 {
		return "", nil
	}
	return g.Tangga[calon[0].ID][0].Komentar, nil
}

// TutupKasus - lihat `repository.Gudang.TutupKasus`.
func (g *Gudang) TutupKasus(_ context.Context, _ *db.Tx, id string, selesai bool, posisi string, saat time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	k, ada := g.Kasus[id]
	if !ada || k.StatusWork != "" {
		return fmt.Errorf("%w: work object %s", repository.ErrKeputusanBersamaan, id)
	}
	if selesai {
		k.StatusWork, posisi = models.StatusSelesai, ""
	}
	k.TglUpdate = saat
	g.Kasus[id] = k
	if g.Posisi == nil {
		g.Posisi = map[string]string{}
	}
	g.Posisi[id] = posisi
	return nil
}

// DaftarKerja - lihat `repository.Gudang.DaftarKerja` (KomiteRouter S6.1; akun atau workbasket `peran`).
func (g *Gudang) DaftarKerja(_ context.Context, akun string, peran []string) ([]models.BarisKerja, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.BarisKerja{}
	for id, k := range g.Kasus {
		if g.Lini[id] != models.LiniNonProp || k.StatusWork != "" || !strings.HasPrefix(id, models.AwalanKomite) {
			continue
		}
		k.Tangga = g.Tangga[id]
		a, ada := k.Giliran()
		if !ada || !k.Pemegang(akun, peran) {
			continue
		}
		b := models.BarisKerja{KasusID: id, KlaimID: k.KlaimID, Tingkat: a.Urut, Count: k.Count, Loop: k.Loop,
			Jabatan: a.Jabatan, TglUpdate: k.TglUpdate}
		if g.Klaim != nil {
			if kl, ok := g.Klaim.lihat(k.KlaimID); ok {
				b.NoKlaim = kl.nilai["ClaimData.NoClaim"]
				for _, r := range kl.daftar["ClaimData.AdjustmentList"] {
					if r["ID"] == k.AdjustmentID {
						b.StatusBaris = r["AcceptanceStatus"]
					}
				}
			}
		}
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].KasusID > out[j].KasusID })
	return out, nil
}

// UrutNomorAkseptasi - penghitung bersama tiruan (periode MM.YYYY bulan berjalan Jakarta).
func (g *Gudang) UrutNomorAkseptasi(_ context.Context, _ *db.Tx, saat time.Time) (models.BahanNomor, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	jenis := g.KodeProduksi + models.HurufAkseptasi
	g.Urut[jenis]++
	return models.BahanNomor{Jenis: jenis, MMYYYY: saat.In(models.Jakarta).Format("01.2006"), Urut: g.Urut[jenis]}, nil
}

// SisipOS - lihat `repository.Gudang.SisipOS`.
func (g *Gudang) SisipOS(_ context.Context, _ *db.Tx, b models.BarisOSAkseptasi, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.OS = append(g.OS, b)
	return nil
}

// SisipXOL2 - lihat `repository.Gudang.SisipXOL2`.
func (g *Gudang) SisipXOL2(_ context.Context, _ *db.Tx, b models.BarisXOL2, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.XOL2 = append(g.XOL2, b)
	return nil
}

// SimpanOSSubjectivity - lihat `repository.Gudang.SimpanOSSubjectivity` (CASEID ada = ubah, selainnya sisip).
func (g *Gudang) SimpanOSSubjectivity(_ context.Context, _ *db.Tx, b models.BarisOSSubjectivity, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if lama, ada := g.OSSubj[b.CaseID]; ada {
		lama.StsSubjectivity, lama.DataJSON = b.StsSubjectivity, b.DataJSON
		g.OSSubj[b.CaseID] = lama
		return nil
	}
	g.OSSubj[b.CaseID] = b
	return nil
}

// SalinJSONKlaim - lihat `repository.Gudang.SalinJSONKlaim`.
func (g *Gudang) SalinJSONKlaim(_ context.Context, _ *db.Tx, idPega, noKlaim, noPolis string, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.JSONKlaim[idPega]; !ada {
		g.JSONKlaim[idPega] = [2]string{noKlaim, noPolis}
	}
	return nil
}

// CatatRiwayatAkseptasi - lihat `repository.Gudang.CatatRiwayatAkseptasi`.
func (g *Gudang) CatatRiwayatAkseptasi(_ context.Context, _ *db.Tx, r models.RiwayatAkseptasi, _ time.Time) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Riwayat = append(g.Riwayat, r)
	return nil
}

// AntreEfek - outbox tiruan: "jenis:rujukan:muatan".
func (g *Gudang) AntreEfek(_ context.Context, _ *db.Tx, jenis, rujukan, muatan string, _ time.Time) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Efek = append(g.Efek, jenis+":"+rujukan+":"+muatan)
	return fmt.Sprintf("UJI-EFEK-%d", len(g.Efek)), nil
}

// Acuan - acuan tiruan.
type Acuan struct {
	// Bank - "bank|cabang|rekening" -> IDOFBANK; Nama - akun -> M_LOGIN_GO.NAME.
	Bank map[string]string
	Nama map[string]string
	// Konversi - nomor akseptasi tanpa titik -> status konversi; Email - ceding -> email.
	Konversi map[string]string
	Email    map[string]string
	// Surel - akun -> M_LOGIN_GO.EMAIL.
	Surel map[string]string
	// AnggotaWB - workbasket -> email anggota aktifnya.
	AnggotaWB map[string][]string
}

// AcuanBaru membuat acuan kosong.
func AcuanBaru() *Acuan {
	return &Acuan{Bank: map[string]string{}, Nama: map[string]string{}, Konversi: map[string]string{},
		Email: map[string]string{}, Surel: map[string]string{}}
}

// IDBankRekening - lihat `repository.Acuan.IDBankRekening`.
func (a *Acuan) IDBankRekening(_ context.Context, bank, cabang, akun string) (string, error) {
	return a.Bank[bank+"|"+cabang+"|"+akun], nil
}

// StatusKonversi - lihat `repository.Acuan.StatusKonversi` (tiruan: peta Konversi).
func (a *Acuan) StatusKonversi(_ context.Context, noAksep string) (string, error) {
	return a.Konversi[noAksep], nil
}

// EmailCeding - lihat `repository.Acuan.EmailCeding`.
func (a *Acuan) EmailCeding(_ context.Context, ceding string) (string, error) {
	return a.Email[ceding], nil
}

// EmailPelaku - lihat `repository.Acuan.EmailPelaku`.
func (a *Acuan) EmailPelaku(_ context.Context, akun string) (string, error) {
	return a.Surel[akun], nil
}

// EmailAnggotaWorkbasket - lihat `repository.Acuan.EmailAnggotaWorkbasket`.
func (a *Acuan) EmailAnggotaWorkbasket(_ context.Context, workbasket string) ([]string, error) {
	return a.AnggotaWB[workbasket], nil
}

// NamaPelaku - lihat `repository.Acuan.NamaPelaku`.
func (a *Acuan) NamaPelaku(_ context.Context, akun string) (string, error) {
	if n, ok := a.Nama[akun]; ok {
		return n, nil
	}
	return akun, nil
}

// ---------------------------------------------------------------- kontrak palsu Claim Non Prop

type klaimTiruan struct {
	nilai    map[string]string
	daftar   map[string][]map[string]string
	tertutup bool
}

func (k *klaimTiruan) salin() *klaimTiruan {
	c := &klaimTiruan{nilai: map[string]string{}, daftar: map[string][]map[string]string{}, tertutup: k.tertutup}
	for j, v := range k.nilai {
		c.nilai[j] = v
	}
	for j, rows := range k.daftar {
		s := make([]map[string]string, 0, len(rows))
		for _, b := range rows {
			m := map[string]string{}
			for p, v := range b {
				m[p] = v
			}
			s = append(s, m)
		}
		c.daftar[j] = s
	}
	return c
}

// KlaimPalsu memenuhi `kontrak.KlaimTreatyNonPropKomite` di atas kasus klaim dalam memori. Baris akseptasi dikenali
// lewat properti `ID` (sama dengan penyedia asli).
type KlaimPalsu struct {
	mu    sync.Mutex
	klaim map[string]*klaimTiruan
	// Dikunci - klaim yang pernah dikunci (S5).
	Dikunci []string
}

var _ kontrak.KlaimTreatyNonPropKomite = (*KlaimPalsu)(nil)

// KlaimBaru membuat kontrak palsu kosong.
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

func (p *KlaimPalsu) lihat(id string) (*klaimTiruan, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ok := p.klaim[id]
	if !ok {
		return nil, false
	}
	return k.salin(), true
}

// Setel menyimpan kasus klaim `id` (nilai halaman + daftar).
func (p *KlaimPalsu) Setel(id string, nilai map[string]string, daftar map[string][]map[string]string) {
	k := (&klaimTiruan{nilai: nilai, daftar: daftar}).salin()
	p.mu.Lock()
	p.klaim[id] = k
	p.mu.Unlock()
}

// Tutup menandai kasus klaim `id` tertutup.
func (p *KlaimPalsu) Tutup(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if k, ok := p.klaim[id]; ok {
		k.tertutup = true
	}
}

// Nilai / Daftar - keadaan kasus klaim (salinan) untuk asersi uji.
func (p *KlaimPalsu) Nilai(id string) map[string]string {
	if k, ok := p.lihat(id); ok {
		return k.nilai
	}
	return nil
}

// Daftar - lihat `Nilai`.
func (p *KlaimPalsu) Daftar(id, jalur string) []map[string]string {
	if k, ok := p.lihat(id); ok {
		return k.daftar[jalur]
	}
	return nil
}

func posisi(k *klaimTiruan, adjID string) int {
	for i, b := range k.daftar["ClaimData.AdjustmentList"] {
		if adjID != "" && b["ID"] == adjID {
			return i + 1
		}
	}
	return 0
}

// BacaKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite`.
func (p *KlaimPalsu) BacaKlaimTreaty(_ context.Context, _ *db.Tx, klaimID, adjID string) (kontrak.KlaimTreaty, error) {
	k, ok := p.lihat(klaimID)
	if !ok {
		return kontrak.KlaimTreaty{}, kontrak.ErrKlaimTreatyTidakAda
	}
	n := posisi(k, adjID)
	if n == 0 {
		return kontrak.KlaimTreaty{}, kontrak.ErrAdjustmentTreatyTidakAda
	}
	return kontrak.KlaimTreaty{Nilai: k.nilai, Daftar: k.daftar, Adjustment: n, Tertutup: k.tertutup}, nil
}

// KunciKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite`.
func (p *KlaimPalsu) KunciKlaimTreaty(_ context.Context, _ *db.Tx, klaimID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ok := p.klaim[klaimID]
	if !ok {
		return kontrak.ErrKlaimTreatyTidakAda
	}
	if k.tertutup {
		return kontrak.ErrKlaimTreatyTertutup
	}
	p.Dikunci = append(p.Dikunci, klaimID)
	return nil
}

// TulisBalikKlaimTreaty - lihat `kontrak.KlaimTreatyNonPropKomite` (daftar putih ditegakkan).
func (p *KlaimPalsu) TulisBalikKlaimTreaty(_ context.Context, _ *db.Tx, klaimID, adjID string,
	u kontrak.UbahanKlaimTreaty) error {
	for j := range u.Header {
		if _, ok := kontrak.JalurHeaderKomiteNonProp[j]; !ok {
			return fmt.Errorf("%w: header %q", kontrak.ErrUbahanKlaimTreatyTidakSah, j)
		}
	}
	for j := range u.Adjustment {
		if _, ok := kontrak.PropAdjustmentKomiteNonProp[j]; !ok {
			return fmt.Errorf("%w: adjustment %q", kontrak.ErrUbahanKlaimTreatyTidakSah, j)
		}
	}
	if len(u.FacRetro) > 0 { // XML Non Prop tanpa Fac Retro (penyedia asli menolak)
		return fmt.Errorf("%w: FacRetro", kontrak.ErrUbahanKlaimTreatyTidakSah)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	k, ok := p.klaim[klaimID]
	if !ok {
		return kontrak.ErrKlaimTreatyTidakAda
	}
	if k.tertutup {
		return kontrak.ErrKlaimTreatyTertutup
	}
	n := posisi(k, adjID)
	if n == 0 {
		return kontrak.ErrAdjustmentTreatyTidakAda
	}
	for j, v := range u.Header {
		k.nilai[j] = v
	}
	b := k.daftar["ClaimData.AdjustmentList"][n-1]
	for j, v := range u.Adjustment {
		b[j] = v
	}
	for _, r := range u.Riwayat {
		k.daftar["ClaimData.SuggestList"] = append(k.daftar["ClaimData.SuggestList"], map[string]string{
			"CommentSuggest": r.Teks, "PICSuggest": r.Pelaku, "IsCedingConfirm": r.Tingkat,
			"DateSuggest": models.FormatWaktu(r.Saat)})
	}
	return nil
}
