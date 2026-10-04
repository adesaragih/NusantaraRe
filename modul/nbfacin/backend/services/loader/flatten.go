// Package loader meratakan SATU dokumen penawaran Fac In (isi CLOB
// POOLDATA.JSON_POLIS.DATA_JSON) menjadi himpunan baris untuk tabel flat (78 rancangan +
// amandemen work owner, amandemen.go) - seam
// `loader.Flatten` spec 11 (`D:\migrasi\RNM\OUTPUT\04-spec\11-spec-pemuatan-data-lama.md`),
// tiket 22.
//
// ⛔ Murni: tidak menyentuh basis data, jam, maupun berkas. Identitas surrogate (ID)
// dan PARENT_ID diberikan repository (tiket 24) dari Baris.Kunci / Baris.Induk.
//
// Mesin generik membaca skema_gen.go (bangkitan workbook rancangan, ./bangkit) yang
// digabung dengan amandemen rancangan keputusan work owner (amandemen.go); aturan
// khusus hanya yang tertulis sebagai keputusan V-xx / butir register (aturan.go,
// amandemen.go), dan tafsiran yang belum pasti ditandai `[dugaan]` di tempatnya. Medan yang
// tidak terpetakan TIDAK dibuang diam-diam: nilainya masuk Hasil.Penampung
// (ADR-0023) dan hitungannya - tanpa nilai, karena nilai dapat berupa data pribadi -
// masuk Diagnostik menurut jalur dan nama medan.
package loader

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Keadaan berhenti keras (spec 11 "Yang wajib berhenti keras", ditambah keadaan
// mesin yang tidak boleh diselesaikan dengan nilai bawaan).
var (
	ErrDokumen       = errors.New("loader: dokumen penawaran tidak dapat diurai")
	ErrIDPega        = errors.New("loader: IDPEGA tidak berbentuk kelas + spasi + NB/RNW/EDM-nomor")
	ErrLiniBisnis    = errors.New("loader: kelompok lini bisnis (V-16) tidak dapat ditetapkan")
	ErrIndukTakTetap = errors.New("loader: tabel induk baris tidak dapat ditetapkan")
	ErrTerlaluDalam  = errors.New("loader: kedalaman wadah berulang melebihi 8 tingkat")
	ErrBukanAngka    = errors.New("loader: nilai kolom NUMBER bukan angka berformat dikenal")
	ErrKolomGanda    = errors.New("loader: satu kolom terisi dua kali di satu baris")
	ErrBentuk        = errors.New("loader: bentuk simpul tidak sesuai jalurnya")
)

// Masukan - satu baris POOLDATA.JSON_POLIS: kunci kerja Pega dan isi dokumennya.
type Masukan struct {
	IDPega   string
	DataJSON []byte
}

// Nilai - isi satu kolom. Kolom NUMBER: Angka (eksak, tidak dibulatkan - ADR-0005);
// kolom lain, termasuk delapan kolom NUMBER yang disimpangkan ke teks
// (kolomTeksMenyimpang, butir 68.1): Teks apa adanya.
type Nilai struct {
	Teks  string
	Angka *apd.Decimal
}

// Baris - satu baris tabel flat. Kunci = identitas surrogate lokal 1..N menurut
// urutan lahir (induk selalu lahir sebelum anaknya - urutan muat spec 11 butir 15);
// Induk = Kunci baris induk (0 = T_WORK_POLIS). Kolom tidak memuat ID dan PARENT_ID.
type Baris struct {
	Kunci, Induk int
	Kolom        map[string]Nilai
}

// Diagnostik - hitungan saja; tidak ada nilai data di dalamnya.
//
// Setiap medan daun TERISI yang bukan metadata px/pz/py berakhir di tepat satu
// ember: Terpetakan, Dibuang (selain dua alasan metadata), TakTerpetakan, atau
// PenunjukBelumDikonversi - ditagih TestSetiapMedanTerhitung.
type Diagnostik struct {
	// Terpetakan - medan daun yang tertulis ke kolom (satu medan dihitung sekali
	// walau mengisi dua kolom, mis. CurrencyList.Name).
	Terpetakan    int
	BarisPerTabel map[string]int
	// MataUangTakDiketahui - "TABEL.KOLOM" -> baris yang diisi UNKNOWN (K-069).
	MataUangTakDiketahui map[string]int
	// KomaDesimal - "TABEL.KOLOM" -> nilai berkoma desimal yang dinormalkan (BAHAN §7).
	KomaDesimal map[string]int
	// LebihDari38Digit - "TABEL.KOLOM" -> angka dengan > 38 digit BERMAKNA (nol di
	// ujung tidak dihitung; batas presisi NUMBER Oracle). Tidak dibulatkan di sini.
	LebihDari38Digit map[string]int
	// Dibuang - alasan -> medan daun terisi yang dibuang menurut keputusan tertulis.
	Dibuang map[string]int
	// CabangKosong - jalur -> unsur bertabel yang nol medan terisinya (V-27).
	CabangKosong map[string]int
	// TakTerpetakan - "jalur :: medan" -> medan daun terisi tanpa kolom di rancangan
	// ("jalur :: *" = seluruh subpohon yang jalurnya tidak ada di Jalur Sumber).
	TakTerpetakan map[string]int
	// PenunjukBelumDikonversi - "jalur :: medan" Idx*/Index* di luar 48 kolom penunjuk
	// teks mentah (butir 72); di fixture dan korpus 02-10-2026: nol.
	PenunjukBelumDikonversi map[string]int
}

// MedanTakDikenal - satu medan daun terisi yang tidak punya kolom, DISIMPAN
// beserta nilainya (ADR-0023 "penampung medan tak dikenal": masuk penampung, tidak
// dibuang). Kunci = Baris.Kunci baris terdekat yang memuatnya; Jalur = jalur
// rancangan + ruas di bawahnya. ⛔ Nilai dapat berupa data pribadi: penampung
// bagian dari Hasil, tidak pernah dicatat ke log (Diagnostik hanya hitungan).
type MedanTakDikenal struct {
	Kunci        int
	Jalur, Medan string
	Nilai        string
	// Penunjuk - medan Idx*/Index* di luar 48 kolom penunjuk teks mentah (butir 72).
	Penunjuk bool
}

// Hasil - keluaran seam: nama tabel -> deretan baris (urutan dipertahankan).
// Penampung wajib KOSONG sebelum pemuatan dinyatakan selesai (ADR-0023 akibat 2).
type Hasil struct {
	Baris      map[string][]Baris
	Penampung  []MedanTakDikenal
	Diagnostik Diagnostik
}

// barisAktif - baris yang sedang diisi beserta rantai leluhurnya.
type barisAktif struct {
	tabel string
	b     *Baris
	induk *barisAktif
}

type pejalan struct {
	idpega, cob string
	jalurKe     map[string]jalurSkema
	peta        map[string]map[string][]kolomSkema // tabel -> FIELD ASLI -> kolom (bisa lebih dari satu)
	kolomAda    map[string]map[string]kolomSkema   // tabel -> nama kolom -> kolom
	hasil       *Hasil
	urut        []*barisAktif
	seqFaktor   map[*Baris]int
	bukanAngka  map[string]int
}

// Flatten - satu dokumen penawaran -> himpunan baris 78 tabel. Galat = berhenti
// keras: tidak ada Hasil sebagian.
func Flatten(m Masukan) (*Hasil, error) { return flatten(m, nil) }

// flatten - kumpulBukanAngka non-nil HANYA dari uji pengukuran: nilai NUMBER yang
// tak terurai dihitung per "TABEL.KOLOM" alih-alih menghentikan, supaya seluruh
// konflik tipe rancangan-vs-data terukur dalam satu jalan. Flatten selalu nil.
func flatten(m Masukan, kumpulBukanAngka map[string]int) (*Hasil, error) {
	noWork, jenis, err := uraiIDPega(m.IDPega)
	if err != nil {
		return nil, err
	}
	akar, err := uraiJSON(m.DataJSON)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDokumen, sebabDokumen(err))
	}
	cob, err := kelompokLini(akar)
	if err != nil {
		return nil, err
	}
	p := &pejalan{idpega: m.IDPega, cob: cob, jalurKe: map[string]jalurSkema{}, peta: map[string]map[string][]kolomSkema{},
		kolomAda: map[string]map[string]kolomSkema{}, seqFaktor: map[*Baris]int{}, bukanAngka: kumpulBukanAngka,
		hasil: &Hasil{Baris: map[string][]Baris{}, Diagnostik: Diagnostik{BarisPerTabel: map[string]int{},
			MataUangTakDiketahui: map[string]int{}, KomaDesimal: map[string]int{}, LebihDari38Digit: map[string]int{},
			Dibuang: map[string]int{}, CabangKosong: map[string]int{}, TakTerpetakan: map[string]int{},
			PenunjukBelumDikonversi: map[string]int{}}}}
	for _, j := range jalurSumber {
		p.jalurKe[j.jalur] = j
	}
	for t, ks := range skemaTabel {
		p.peta[t], p.kolomAda[t] = map[string][]kolomSkema{}, map[string]kolomSkema{}
		for _, k := range ks {
			p.kolomAda[t][k.nama] = k
			if k.medan != "" && (!k.turunan || k.kolomPay()) {
				p.peta[t][k.medan] = append(p.peta[t][k.medan], k)
			}
		}
	}

	kerja := p.lahir("T_WORK_POLIS", nil, "", "", 0)
	if err := p.isi(kerja, kolomJenisWork, jenis); err != nil {
		return nil, err
	}
	if err := p.isi(kerja, kolomNoWork, noWork); err != nil {
		return nil, err
	}
	umum := p.lahir("T_GENERAL_POLIS", kerja, "", "", 0)
	if err := p.objek(akar, "", "", umum, "", 0); err != nil {
		return nil, err
	}
	if err := p.mataUang(); err != nil {
		return nil, err
	}
	return p.hasil, nil
}

// sebabDokumen - sebab galat urai TANPA potongan isi dokumen: pesan encoding/json
// dapat mengutip karakter atau teks masukan.
func sebabDokumen(err error) string {
	var sintaks *json.SyntaxError
	if errors.As(err, &sintaks) {
		return fmt.Sprintf("sintaks JSON tidak sah di bita %d", sintaks.Offset)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return "dokumen terpotong"
	}
	var bentuk galatBentukDokumen
	if errors.As(err, &bentuk) {
		return string(bentuk)
	}
	return "dokumen tidak dapat diurai"
}

// uraiIDPega - "kelas spasi pyID", pyID = JENIS-nomor (lembar Kolom T_WORK_POLIS).
func uraiIDPega(id string) (noWork, jenis string, err error) {
	i := strings.LastIndexByte(id, ' ')
	if i <= 0 || i == len(id)-1 {
		return "", "", ErrIDPega
	}
	noWork = id[i+1:]
	j, nomor, ok := strings.Cut(noWork, "-")
	if !ok || !jenisWork[j] || nomor == "" {
		return "", "", ErrIDPega
	}
	return noWork, j, nil
}

// kelompokLini - V-16, berjenjang, berhenti di kecocokan pertama. "ada" = larik
// dengan sedikitnya satu unsur `[dugaan]`: V-16 diukur di XML, yang tidak memuat
// larik kosong; di JSON larik kosong (mis. "VehicleList":[]) lazim.
func kelompokLini(akar *simpul) (string, error) {
	lokasi := akar.medan("LocationList")
	switch {
	case adaDiBawah(lokasi, func(s *simpul) bool {
		return berisiLarik(s.medan("Property").medan("PropertyItemList"))
	}):
		return "FIRE", nil
	case adaKunciDiBawah(lokasi, "AnekaList"):
		return "Aneka", nil
	case berisiLarik(akar.medan("VehicleList")):
		return "MBUCar", nil
	case berisiLarik(akar.medan("CargoList")):
		return "MarineCargo", nil
	}
	// Langkah 5 - satu-satunya yang membaca NILAI (V-16a). Yang mengikat: medan
	// QuotationData.BusinessType di akar; kemunculan lain di fixture seluruhnya di
	// cabang OldData yang dibuang V-19.
	bt := akar.medan("QuotationData").medan("BusinessType")
	if bt == nil || bt.jenis != simpulNilai || !businessTypeDikenal[bt.teks] {
		return "", fmt.Errorf("%w: BusinessType kosong atau di luar 18 nilai yang terbaca", ErrLiniBisnis)
	}
	if bt.teks != "Life" && bt.teks != "PA" {
		return "", fmt.Errorf("%w: bentuk dokumen jatuh ke langkah 5 tetapi BusinessType bukan Life/PA", ErrLiniBisnis)
	}
	return bt.teks, nil
}

func berisiLarik(s *simpul) bool { return s != nil && s.jenis == simpulLarik && len(s.anak) > 0 }

// adaDiBawah - uji f pada tiap unsur larik s.
func adaDiBawah(s *simpul, f func(*simpul) bool) bool {
	if s == nil || s.jenis != simpulLarik {
		return false
	}
	for _, u := range s.anak {
		if f(u) {
			return true
		}
	}
	return false
}

// adaKunciDiBawah - ada larik berisi bernama k di kedalaman mana pun di bawah s,
// MELEWATI cabang yang dibuang (OldData, Total*, ...): V-16 menuntut LocationList
// "bisnis", dan aturan lama justru salah karena ikut membaca cabang non-bisnis.
func adaKunciDiBawah(s *simpul, k string) bool {
	if s == nil {
		return false
	}
	for i, a := range s.anak {
		if s.jenis == simpulObjek {
			if _, dibuang := cabangDibuang[s.kunci[i]]; dibuang {
				continue
			}
			if s.kunci[i] == k && berisiLarik(a) {
				return true
			}
		}
		if adaKunciDiBawah(a, k) {
			return true
		}
	}
	return false
}

// lahir - baris baru tabel t, anak induk. jalur = jalur rancangan, posisi = jalur
// mentah berposisi (bahan ROW_UID), seq = posisi di lariknya (0 = bukan larik).
func (p *pejalan) lahir(t string, induk *barisAktif, jalur, posisi string, seq int) *barisAktif {
	b := &Baris{Kunci: len(p.urut) + 1, Kolom: map[string]Nilai{}}
	if induk != nil {
		b.Induk = induk.b.Kunci
	}
	ba := &barisAktif{tabel: t, b: b, induk: induk}
	p.urut = append(p.urut, ba)
	p.hasil.Diagnostik.BarisPerTabel[t]++
	kol := p.kolomAda[t]
	b.Kolom[kolomIDPega] = Nilai{Teks: p.idpega}
	if _, ada := kol[kolomCOB]; ada {
		b.Kolom[kolomCOB] = Nilai{Teks: p.cob}
	}
	if _, ada := kol[kolomTabelInduk]; ada && induk != nil {
		b.Kolom[kolomTabelInduk] = Nilai{Teks: induk.tabel}
	}
	if _, ada := kol[kolomJalurSumber]; ada {
		b.Kolom[kolomJalurSumber] = Nilai{Teks: jalur}
	}
	if _, ada := kol[kolomSeq]; ada {
		if seq == 0 {
			seq = 1
		}
		b.Kolom[kolomSeq] = Nilai{Angka: apd.New(int64(seq), 0)}
	}
	if _, ada := kol[kolomRowUID]; ada {
		b.Kolom[kolomRowUID] = Nilai{Teks: uidBaris(p.idpega, t, posisi)}
	}
	return ba
}

// uidBaris - ROW_UID fase 1 (K-073: sementara, tidak boleh dijadikan sandaran di
// luar tabel flat). Dibangkitkan DETERMINISTIK dari kunci kerja + tabel + posisi
// mentah - pilihan agent, supaya Flatten tetap murni dan muat ulang memberi isi
// yang sama (spec 11, uji repository butir 4). Bentuk UUID versi 5 (SHA-1).
func uidBaris(idpega, tabel, posisi string) string {
	h := sha1.Sum([]byte("loader.Flatten\x00" + idpega + "\x00" + tabel + "\x00" + posisi))
	h[6] = h[6]&0x0f | 0x50
	h[8] = h[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", h[0:4], h[4:6], h[6:8], h[8:10], h[10:16])
}

// objek - jalani medan objek o ke baris ba. jalur = jalur rancangan baris ba
// (lipatan tidak mengubahnya), awalan = awalan medan lipatan (V-22).
func (p *pejalan) objek(o *simpul, jalur, posisi string, ba *barisAktif, awalan string, dalam int) error {
	for i, k := range o.kunci {
		v := o.anak[i]
		posisiK := posisi + "/" + k
		if v.jenis == simpulNilai {
			if err := p.daun(ba, jalur, awalan, k, v.teks); err != nil {
				return err
			}
			continue
		}
		if alasan, ada := cabangDibuang[k]; ada {
			selamat := p.selamatkan(ba, jalurAnakDari(jalur, k), v, medanDiselamatkan[k])
			if n := hitungDaun(v) - selamat; n > 0 {
				p.hasil.Diagnostik.Dibuang[alasan] += n
			}
			continue
		}
		if awalan == "" {
			if aturan, ada := lipat[lipatan{ba.tabel, k}]; ada {
				if v.jenis != simpulObjek {
					return fmt.Errorf("%w: %s/%s dilipat tetapi bukan objek", ErrBentuk, jalur, k)
				}
				if err := p.lipatan(v, jalur, k, posisiK, ba, aturan, dalam); err != nil {
					return err
				}
				continue
			}
		}
		if ba.tabel == tabelSkoring && v.jenis == simpulObjek {
			ok, err := p.skoring(v, k, jalur, posisiK, ba)
			if err != nil {
				return err
			}
			if ok {
				continue
			}
		}
		if err := p.anak(v, k, jalur, posisiK, ba, dalam); err != nil {
			return err
		}
	}
	return nil
}

// lipatan - halaman dilipat ke ba: medannya ke baris yang sama (bila
// aturan.medanSendiri); halaman di dalamnya yang ada di aturan.dalam ikut dilipat
// dengan awalannya; sisanya anak biasa dengan jalur rancangan baris ba (V-22a
// ListInstallment, V-24b Ship/LC).
func (p *pejalan) lipatan(o *simpul, jalur, halaman, posisi string, ba *barisAktif, aturan aturanLipat, dalam int) error {
	for i, k := range o.kunci {
		v := o.anak[i]
		if awal, ada := aturan.dalam[k]; ada && v.jenis == simpulObjek {
			bersarang := &simpul{jenis: simpulObjek}
			for j, kk := range v.kunci {
				if v.anak[j].jenis == simpulNilai {
					if err := p.daun(ba, jalur, awal, kk, v.anak[j].teks); err != nil {
						return err
					}
					continue
				}
				bersarang.kunci, bersarang.anak = append(bersarang.kunci, kk), append(bersarang.anak, v.anak[j])
			}
			if err := p.objek(bersarang, jalur, posisi+"/"+k, ba, "", dalam); err != nil {
				return err
			}
			continue
		}
		if kol, ada := aturan.khusus[k]; ada && v.jenis == simpulNilai {
			if v.teks != "" {
				if err := p.isiMedan(ba, kol, v.teks); err != nil {
					return err
				}
			}
			continue
		}
		if v.jenis == simpulNilai {
			if aturan.medanSendiri {
				if err := p.daun(ba, jalur, "", k, v.teks); err != nil {
					return err
				}
			} else if v.teks != "" {
				p.catatTak(ba, jalur+"/"+halaman, k, v.teks)
			}
			continue
		}
		if err := p.objek(&simpul{jenis: simpulObjek, kunci: []string{k}, anak: []*simpul{v}}, jalur, posisi, ba, "", dalam); err != nil {
			return err
		}
	}
	return nil
}

// anak - halaman / larik k di bawah baris ba: tabel baru menurut Jalur Sumber.
func (p *pejalan) anak(v *simpul, k, jalur, posisi string, ba *barisAktif, dalam int) error {
	ruas := k
	if g, ada := gantiNama[[2]string{ruasTerakhir(jalur), k}]; ada {
		ruas = g
	}
	jalurAnak := jalurAnakDari(jalur, ruas)
	j, ada := p.jalurKe[jalurAnak]
	if !ada {
		p.catatSubpohon(ba, jalurAnak, v)
		return nil
	}
	var unsur []*simpul
	larik := v.jenis == simpulLarik
	if larik {
		unsur = v.anak
		dalam++
		if dalam > kedalamanMaks {
			return fmt.Errorf("%w: %s", ErrTerlaluDalam, jalurAnak)
		}
	} else {
		unsur = []*simpul{v}
	}
	for i, u := range unsur {
		if u.jenis != simpulObjek {
			return fmt.Errorf("%w: unsur %s bukan objek", ErrBentuk, jalurAnak)
		}
		if j.induk != ba.tabel {
			return fmt.Errorf("%w: %s di bawah %s, Jalur Sumber menetapkan %s", ErrIndukTakTetap, jalurAnak, ba.tabel, j.induk)
		}
		// SEQ_NO = posisi di larik sumber, bukan di antara baris yang jadi: unsur
		// kosong yang dibatalkan V-27 meninggalkan celah, supaya posisi tetap sejajar
		// lintas versi (V-40 menjodohkan menurut posisi; V-43 baris tidak pernah dihapus).
		seq, pos := 0, posisi
		if larik {
			seq, pos = i+1, posisiLarik(posisi, i)
		}
		if err := p.mungkinLahir(j.tabel, ba, jalurAnak, pos, seq, func(baru *barisAktif) error {
			return p.objek(u, jalurAnak, pos, baru, "", dalam)
		}); err != nil {
			return err
		}
	}
	return nil
}

// mungkinLahir - V-27: "cabang yang punya baris tetapi NOL kolom terisi di seluruh
// keturunannya tidak dijadikan tabel". Subpohon dijalani lebih dulu; bila tidak
// satu medan pun tertulis ke kolom (medan yang dibuang, tak terpetakan, atau
// kosong tidak mengisi kolom), baris itu beserta seluruh turunannya dibatalkan.
// Diagnostik subpohon tetap tercatat - hanya barisnya yang batal.
//
// Entri penampung milik baris yang batal DIPINDAH ke baris induk - nomor Kunci baris
// batal dipakai ulang baris berikutnya, jadi entri yang tetap memegangnya akan
// menempel ke baris yang salah.
func (p *pejalan) mungkinLahir(t string, induk *barisAktif, jalur, posisi string, seq int, isi func(*barisAktif) error) error {
	awalUrut, awalTerpetakan, awalPenampung := len(p.urut), p.hasil.Diagnostik.Terpetakan, len(p.hasil.Penampung)
	baru := p.lahir(t, induk, jalur, posisi, seq)
	if err := isi(baru); err != nil {
		return err
	}
	if p.hasil.Diagnostik.Terpetakan > awalTerpetakan {
		return nil
	}
	for _, ba := range p.urut[awalUrut:] {
		p.hasil.Diagnostik.BarisPerTabel[ba.tabel]--
		if p.hasil.Diagnostik.BarisPerTabel[ba.tabel] == 0 {
			delete(p.hasil.Diagnostik.BarisPerTabel, ba.tabel)
		}
	}
	p.urut = p.urut[:awalUrut]
	for i := awalPenampung; i < len(p.hasil.Penampung); i++ {
		p.hasil.Penampung[i].Kunci = induk.b.Kunci
	}
	p.hasil.Diagnostik.CabangKosong[jalur]++
	return nil
}

// digitBermakna - jumlah digit koefisien tanpa nol di ujungnya.
func digitBermakna(d *apd.Decimal) int64 {
	var r apd.Decimal
	r.Reduce(d)
	return r.NumDigits()
}

func jalurAnakDari(jalur, k string) string {
	if jalur == "" {
		return k
	}
	return jalur + "/" + k
}

// selamatkan - medan bernama di `medan` di dalam cabang dibuang v (langsung, atau di
// unsur lariknya) ke penampung; mengembalikan jumlahnya, yang pemanggil KURANGKAN dari
// hitungan Dibuang cabang itu. Pohon tidak diubah. Tanpa ini K-069 (7b) hilang
// bersama ViewSuggest.
func (p *pejalan) selamatkan(ba *barisAktif, jalur string, v *simpul, medan map[string]bool) int {
	n := 0
	ambil := func(o *simpul, j string) {
		for i, k := range o.kunci {
			if a := o.anak[i]; medan[k] && a.jenis == simpulNilai && a.teks != "" {
				p.catatTak(ba, j, k, a.teks)
				n++
			}
		}
	}
	switch v.jenis {
	case simpulObjek:
		ambil(v, jalur)
	case simpulLarik:
		for i, u := range v.anak {
			if u.jenis == simpulObjek {
				ambil(u, posisiLarik(jalur, i))
			}
		}
	}
	return n
}

// posisiLarik - jalur unsur ke-i (dari 0) sebuah larik: "jalur[i+1]".
func posisiLarik(jalur string, i int) string { return jalur + "[" + strconv.Itoa(i+1) + "]" }

func ruasTerakhir(jalur string) string {
	if i := strings.LastIndexByte(jalur, '/'); i >= 0 {
		return jalur[i+1:]
	}
	return jalur
}

// hitungDaun - medan daun terisi (bukan kosong, bukan px/pz/py) di subpohon s.
// Unsur larik yang berupa nilai ikut terhitung (tidak bernama, jadi bukan metadata).
func hitungDaun(s *simpul) int {
	if s.jenis == simpulNilai {
		return 0
	}
	n := 0
	for i, a := range s.anak {
		if a.jenis != simpulNilai {
			n += hitungDaun(a)
			continue
		}
		if a.teks == "" {
			continue
		}
		if s.jenis == simpulLarik || !metadata(s.kunci[i]) {
			n++
		}
	}
	return n
}

var polaPenunjuk = regexp.MustCompile(`^(Idx|Index)[A-Z]`)

// daun - satu medan daun ke kolom baris ba.
func (p *pejalan) daun(ba *barisAktif, jalur, awalan, medan, teks string) error {
	if teks == "" {
		return nil
	}
	if al := alasanDaun(medan); al != "" {
		p.hasil.Diagnostik.Dibuang[al]++
		return nil
	}
	nama := awalan + medan
	sasaran, b := ba.tabel, ba.b
	ks := p.peta[sasaran][nama]
	if awalan == "" && medan == "ID" {
		if kol, v49 := medanID[sasaran]; v49 {
			ks = []kolomSkema{p.kolomAda[sasaran][kol]}
		}
	}
	if awalan == "" {
		if kol, ganda := medanGanda[sasaran][medan]; ganda {
			ks = append(ks, p.kolomAda[sasaran][kol])
		}
	}
	if len(ks) == 0 && sasaran == "T_GENERAL_POLIS" {
		// V-48: enam medan akar menerangkan jalannya pekerjaan -> T_WORK_POLIS.
		if kk := p.peta["T_WORK_POLIS"][nama]; len(kk) > 0 {
			ks, sasaran, b = kk, "T_WORK_POLIS", ba.induk.b
		}
	}
	if len(ks) == 0 {
		p.catatTak(ba, jalur, nama, teks)
		return nil
	}
	semua := true
	for _, k := range ks {
		ok, err := p.tulis(sasaran, b, k, teks)
		if err != nil {
			return err
		}
		semua = semua && ok
	}
	if semua {
		p.hasil.Diagnostik.Terpetakan++
	}
	return nil
}

// catatSubpohon - subpohon tanpa tabel: medan terisinya ke TakTerpetakan "jalur :: *"
// dan ke penampung satu per satu (jalur lengkapnya), milik baris ba.
func (p *pejalan) catatSubpohon(ba *barisAktif, jalur string, v *simpul) {
	if n := hitungDaun(v); n > 0 {
		p.hasil.Diagnostik.TakTerpetakan[jalur+" :: *"] += n
		p.tampungSubpohon(ba.b.Kunci, jalur, v)
	}
}

// tampungSubpohon - setiap medan terisi subpohon v (bukan px/pz/py) ke penampung.
// Unsur larik bernilai bernama "[n]".
func (p *pejalan) tampungSubpohon(kunci int, jalur string, v *simpul) {
	for i, a := range v.anak {
		nama := "[" + strconv.Itoa(i+1) + "]"
		if v.jenis == simpulObjek {
			nama = v.kunci[i]
		}
		if a.jenis != simpulNilai {
			sub := jalur + "/" + nama
			if v.jenis == simpulLarik {
				sub = jalur + nama
			}
			p.tampungSubpohon(kunci, sub, a)
			continue
		}
		if a.teks == "" || v.jenis == simpulObjek && metadata(nama) {
			continue
		}
		p.hasil.Penampung = append(p.hasil.Penampung, MedanTakDikenal{Kunci: kunci, Jalur: jalur, Medan: nama, Nilai: a.teks})
	}
}

// catatTak - medan daun terisi tanpa kolom: Idx*/Index* (di luar 48 kolom penunjuk) ke
// PenunjukBelumDikonversi, selainnya ke TakTerpetakan; keduanya juga ke penampung beserta nilainya
// (ADR-0023). Kunci Diagnostik = jalur + nama medan, tanpa nilai.
func (p *pejalan) catatTak(ba *barisAktif, jalur, medan, teks string) {
	if al := alasanDaun(medan); al != "" {
		p.hasil.Diagnostik.Dibuang[al]++
		return
	}
	kunci := jalurAtauAkar(jalur) + " :: " + medan
	penunjuk := polaPenunjuk.MatchString(medan)
	if penunjuk {
		p.hasil.Diagnostik.PenunjukBelumDikonversi[kunci]++
	} else {
		p.hasil.Diagnostik.TakTerpetakan[kunci]++
	}
	p.hasil.Penampung = append(p.hasil.Penampung, MedanTakDikenal{Kunci: ba.b.Kunci, Jalur: jalurAtauAkar(jalur),
		Medan: medan, Nilai: teks, Penunjuk: penunjuk})
}

func jalurAtauAkar(j string) string {
	if j == "" {
		return "(akar)"
	}
	return j
}

// isi - kolom sistem/turunan baris ba diisi teks (atau angka bila NUMBER).
func (p *pejalan) isi(ba *barisAktif, kolom, teks string) error {
	_, err := p.tulis(ba.tabel, ba.b, p.kolomAda[ba.tabel][kolom], teks)
	return err
}

// isiMedan - isi untuk nilai yang berasal dari SATU medan daun (V-30 ChechBoxN /
// ScoreN): ikut terhitung Terpetakan.
func (p *pejalan) isiMedan(ba *barisAktif, kolom, teks string) error {
	ok, err := p.tulis(ba.tabel, ba.b, p.kolomAda[ba.tabel][kolom], teks)
	if ok {
		p.hasil.Diagnostik.Terpetakan++
	}
	return err
}

var (
	polaBulat = regexp.MustCompile(`^-?[0-9]+$`)
	polaTitik = regexp.MustCompile(`^-?[0-9]+\.[0-9]+$`)
	polaKoma  = regexp.MustCompile(`^-?[0-9]+,[0-9]+$`)
)

// tulis - satu nilai ke kolom k. NUMBER diurai eksak dari tiga format yang terukur
// (BAHAN §7: bulat, titik desimal, koma desimal sebagai pengecualian); format lain
// berhenti keras - tanpa menyebut nilainya. ok = nilai tertulis.
func (p *pejalan) tulis(tabel string, b *Baris, k kolomSkema, teks string) (bool, error) {
	if k.nama == "" {
		// Aturan menyebut kolom yang tidak ada di skema - bug program, dijaga
		// TestAturanMenunjukKolomYangAda; bukan keadaan dokumen.
		panic("loader: aturan menyebut kolom yang tidak ada di skema tabel " + tabel)
	}
	if _, sudah := b.Kolom[k.nama]; sudah {
		return false, fmt.Errorf("%w: %s.%s", ErrKolomGanda, tabel, k.nama)
	}
	if !k.angka() || kolomTeksMenyimpang[tabel+"."+k.nama] != "" {
		b.Kolom[k.nama] = Nilai{Teks: teks}
		return true, nil
	}
	s := teks
	switch {
	case polaBulat.MatchString(s), polaTitik.MatchString(s):
	case polaKoma.MatchString(s):
		s = strings.Replace(s, ",", ".", 1)
		p.hasil.Diagnostik.KomaDesimal[tabel+"."+k.nama]++
	default:
		if p.bukanAngka != nil {
			p.bukanAngka[tabel+"."+k.nama]++
			return false, nil
		}
		return false, fmt.Errorf("%w: %s.%s", ErrBukanAngka, tabel, k.nama)
	}
	d, _, err := apd.NewFromString(s)
	if err != nil {
		return false, fmt.Errorf("%w: %s.%s", ErrBukanAngka, tabel, k.nama)
	}
	if digitBermakna(d) > 38 {
		p.hasil.Diagnostik.LebihDari38Digit[tabel+"."+k.nama]++
	}
	b.Kolom[k.nama] = Nilai{Angka: d}
	return true, nil
}

// skoring - V-30 di bawah baris T_DATASCORINGRISKLIST: halaman k adalah faktor
// (punya Checked atau ChechBoxN), atau kelompok yang SELURUH halaman anaknya faktor.
// false = bukan bentuk skoring; jalani sebagai anak biasa.
func (p *pejalan) skoring(v *simpul, k, jalur, posisi string, ba *barisAktif) (bool, error) {
	if faktor(v) {
		return true, p.faktor(v, "", k, jalur, posisi, ba)
	}
	if len(v.kunci) == 0 {
		return false, nil
	}
	ada := false
	for i, a := range v.anak {
		switch {
		case a.jenis == simpulObjek && faktor(a):
			ada = true
		case a.jenis == simpulNilai && alasanDaun(v.kunci[i]) == alasanMeta:
		default:
			return false, nil
		}
	}
	if !ada {
		return false, nil
	}
	for i, a := range v.anak {
		if a.jenis == simpulObjek {
			if err := p.faktor(a, k, v.kunci[i], jalur, posisi+"/"+v.kunci[i], ba); err != nil {
				return true, err
			}
		}
	}
	return true, nil
}

func faktor(s *simpul) bool {
	if s.jenis != simpulObjek {
		return false
	}
	for _, k := range s.kunci {
		if k == medanTandaCek || polaOpsi.MatchString(k) && strings.HasPrefix(k, "ChechBox") {
			return true
		}
	}
	return false
}

func (p *pejalan) faktor(o *simpul, grup, nama, jalur, posisi string, ba *barisAktif) error {
	jalurF := jalur + "/" + ruasFaktor
	p.seqFaktor[ba.b]++
	return p.mungkinLahir(tabelFaktor, ba, jalurF, posisi, p.seqFaktor[ba.b], func(f *barisAktif) error {
		return p.isiFaktor(o, grup, nama, jalurF, posisi, f)
	})
}

// isiFaktor - medan satu faktor V-30 ke baris f dan baris opsinya.
func (p *pejalan) isiFaktor(o *simpul, grup, nama, jalurF, posisi string, f *barisAktif) error {
	if grup != "" {
		if err := p.isi(f, "FACTOR_GROUP", grup); err != nil {
			return err
		}
	}
	if err := p.isi(f, "FACTOR_NAME", nama); err != nil {
		return err
	}
	type opsi struct{ cek, skor string }
	semua := map[int]*opsi{}
	for i, k := range o.kunci {
		a := o.anak[i]
		m := polaOpsi.FindStringSubmatch(k)
		if m == nil || a.jenis != simpulNilai {
			if a.jenis != simpulNilai {
				p.catatSubpohon(f, jalurF+"/"+k, a)
				continue
			}
			if err := p.daun(f, jalurF, "", k, a.teks); err != nil {
				return err
			}
			continue
		}
		n, _ := strconv.Atoi(m[2])
		if semua[n] == nil {
			semua[n] = &opsi{}
		}
		if m[1] == "ChechBox" {
			semua[n].cek = a.teks
		} else {
			semua[n].skor = a.teks
		}
	}
	nomor := make([]int, 0, len(semua))
	for n := range semua {
		nomor = append(nomor, n)
	}
	sort.Ints(nomor)
	for _, n := range nomor {
		op := semua[n]
		if op.cek == "" && op.skor == "" {
			continue
		}
		// SEQ_NO = N (posisi pasangan di sumber), sejalan A50: opsi kosong meninggalkan celah.
		bo := p.lahir(tabelOpsi, f, jalurF+"/"+ruasOpsi, posisi+"#"+strconv.Itoa(n), n)
		if err := p.isi(bo, "OPTION_NO", strconv.Itoa(n)); err != nil {
			return err
		}
		if op.cek != "" {
			if err := p.isiMedan(bo, "CHECK_BOX", op.cek); err != nil {
				return err
			}
		}
		if op.skor != "" {
			if err := p.isiMedan(bo, "SCORE", op.skor); err != nil {
				return err
			}
		}
	}
	return nil
}

// mataUang - K-063 (c) pewarisan lalu K-069 sentinel UNKNOWN, menurut urutan lahir
// (leluhur selalu lebih dulu, jadi nilai warisannya sudah final).
func (p *pejalan) mataUang() error {
	if err := p.kodeDariHalamanCurrency(); err != nil {
		return err
	}
	for _, ba := range p.urut {
		for _, k := range skemaTabel[ba.tabel] {
			if !k.bawaanUnknown() {
				continue
			}
			if _, ada := ba.b.Kolom[k.nama]; ada {
				continue
			}
			if sumber, waris := warisMataUang[ba.tabel]; waris && k.nama == kolomKodeUang {
				for l := ba.induk; l != nil; l = l.induk {
					if l.tabel == sumber {
						if v, ada := l.b.Kolom[kolomKodeUang]; ada {
							ba.b.Kolom[k.nama] = v
						}
						break
					}
				}
				if v, ada := ba.b.Kolom[k.nama]; ada {
					if v.Teks == kodeTakDiketahui {
						p.hasil.Diagnostik.MataUangTakDiketahui[ba.tabel+"."+k.nama]++
					}
					continue
				}
			}
			ba.b.Kolom[k.nama] = Nilai{Teks: kodeTakDiketahui}
			p.hasil.Diagnostik.MataUangTakDiketahui[ba.tabel+"."+k.nama]++
		}
		p.hasil.Baris[ba.tabel] = append(p.hasil.Baris[ba.tabel], *ba.b)
	}
	return nil
}

// kodeDariHalamanCurrency - butir 68.3: CURRENCY_CODE baris kodeDariCurrency diisi
// NAME baris T_CURRENCY anaknya (halaman Currency/Name), SEBELUM pewarisan K-063 (c)
// berjalan, supaya T_SPREADINGLIST dkk. mewarisi kode yang sudah final. Bila baris
// itu sendiri sudah membawa kode (FIELD ASLI Name) dan berbeda dari anaknya -
// ErrKolomGanda: mana yang benar tidak ditebak.
func (p *pejalan) kodeDariHalamanCurrency() error {
	anakCurrency := map[*Baris]string{}
	for _, ba := range p.urut {
		if ba.induk != nil && kodeDariCurrency[ba.induk.tabel] == ba.tabel {
			if v, ada := ba.b.Kolom[kolomNamaCurrency]; ada && v.Teks != "" {
				if lama, sudah := anakCurrency[ba.induk.b]; sudah && lama != v.Teks {
					return fmt.Errorf("%w: %s punya dua halaman Currency berkode berbeda", ErrKolomGanda, ba.induk.tabel)
				}
				anakCurrency[ba.induk.b] = v.Teks
			}
		}
	}
	for _, ba := range p.urut {
		kode, ada := anakCurrency[ba.b]
		if !ada {
			continue
		}
		if v, sudah := ba.b.Kolom[kolomKodeUang]; sudah {
			if v.Teks != kode {
				return fmt.Errorf("%w: %s.%s berbeda dari Currency/Name anaknya", ErrKolomGanda, ba.tabel, kolomKodeUang)
			}
			continue
		}
		ba.b.Kolom[kolomKodeUang] = Nilai{Teks: kode}
	}
	return nil
}
