// Package deidentifikasi mengubah berkas kasus JSON nyata menjadi fixture yang
// aman disimpan: identitas hilang, setiap angka utuh (tiket NB-15).
//
// ⛔ Penyaringan memakai DAFTAR KUNCI EKSPLISIT (`DaftarBuang`), tidak pernah pola,
// substring, atau regex kata kunci: filter pola akan membuang `TotalSumInsured`
// ("insured") dan `CedingRetention` / `ShareCeding` ("ceding") - angka yang justru
// diuji. Medan identitas lain ditambahkan lewat TINJAUAN MANUAL; alat ini hanya
// melaporkan kandidatnya (`KandidatIdentitas`) dan kebocorannya (`Periksa`).
package deidentifikasi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// DaftarBuang - kunci yang NILAINYA dikosongkan. Kunci tetap ada, sehingga bentuk
// dokumen tidak berubah. Baris 1-2: dua belas kunci tiket 15; baris 3-4: sembilan
// medan tempat nilai itu muncul lagi di berkas P-5 (keputusan work owner
// 01-10-2026, butir 32); baris 5: dua medan bocor sisa (butir 36); sisanya nomor
// polis/kasus (butir 38). Nilainya identitas, jadi sekaligus SUMBER deteksi
// kebocoran di `Periksa`.
var DaftarBuang = []string{
	"InsuredName", "InsuredID", "ASMAddress", "SelectedLocationAddress", "RoadName", "ASMZipCode",
	"RiskZipCode", "Email", "pxCreateOpName", "MarketingName", "CedingCoName", "CoinsName",
	"pxCreateOperator", "AccumulationDescription", "AccumulationCode", "PIC", "PICSuggest",
	"CommentSuggest", "PropertiItemNote", "SobName", "TopRiskLocation",
	"Comment", "OperatorID",
	"PolicyNo", "OldPolicyNo", "EndorsementNo", "InvoiceNumber", "NoOfferSlip", "PolicyMasterNumber",
	"PolicyMasterIDPega", "IDNewBisnis", "IDFollowingNB", "Following",
}

// DaftarKosongkan - teks bebas dan bagian alamat (butir 38): nilainya dikosongkan
// karena MUNGKIN memuat identitas, tetapi bukan sumber deteksi kebocoran - isinya
// menyalin medan biasa (angka TSI, nama okupasi), sehingga memburunya di medan
// lain hanya memberi positif palsu.
var DaftarKosongkan = []string{
	"Remarks", "REMARK", "EdmNote", "EdmSourceNote", "GoodNote", "OccupationNote", "CoverageNote",
	"TradingNote", "ConveyanceNote", "PackingNote", "EndorsmentReason", "Description", "NoteFinalScore",
	"Detail", "NM_SHIP", "BranchName",
	"ASMCity", "ASMDistrict", "ASMRW",
}

// JalurBuang - JALUR yang nilainya dikosongkan, dicocokkan persis (bukan pola):
// `ID` tingkat dokumen adalah identitas kasus, sedangkan `ID` di jalur lain kode
// (`Currency.ID`, `Ship.ID`) - jadi kuncinya tidak bisa dibuang utuh (butir 38).
var JalurBuang = []string{"$.ID", "$.OldData.ID", "$.OldData.OldData.ID"}

// sumberRahasia - nilai di `jalur` identitas: dikosongkan DAN diburu di medan lain.
func sumberRahasia(jalur, kunci string) bool {
	return dibuang(kunci) || ada(JalurBuang, jalur)
}

// buang - nilai di `jalur` (kunci terakhirnya `kunci`) dikosongkan.
func buang(jalur, kunci string) bool {
	return sumberRahasia(jalur, kunci) || ada(DaftarKosongkan, kunci)
}

func ada(daftar []string, s string) bool {
	for _, x := range daftar {
		if x == s {
			return true
		}
	}
	return false
}

func dibuang(kunci string) bool { return ada(DaftarBuang, kunci) }

// Laporan - berapa nilai yang dikosongkan per kunci DaftarBuang, atau per jalur
// JalurBuang.
type Laporan struct {
	Dibuang map[string]int
}

// Bersihkan menyalin JSON `masuk` token demi token - urutan kunci dan teks angka
// apa adanya - dan mengganti nilai kunci DaftarBuang dengan teks kosong.
func Bersihkan(masuk []byte) ([]byte, Laporan, error) {
	dec := json.NewDecoder(bytes.NewReader(masuk))
	dec.UseNumber()
	var out bytes.Buffer
	lap := Laporan{Dibuang: map[string]int{}}
	if err := salinNilai(dec, &out, &lap, "$", false); err != nil {
		return nil, Laporan{}, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, Laporan{}, errors.New("deidentifikasi: ada data sesudah nilai JSON utama")
	}
	return out.Bytes(), lap, nil
}

// salinNilai menyalin satu nilai JSON di `jalur` (bentuk jalur sama dengan
// `daunSemua`). `kosongkan` = nilai ini milik kunci atau jalur BUANG.
func salinNilai(dec *json.Decoder, out *bytes.Buffer, lap *Laporan, jalur string, kosongkan bool) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	switch v := tok.(type) {
	case json.Delim:
		if kosongkan {
			return fmt.Errorf("deidentifikasi: nilai kunci BUANG berupa %v, bukan teks - tinjau manual", v)
		}
		switch v {
		case '{':
			out.WriteByte('{')
			for i := 0; dec.More(); i++ {
				if i > 0 {
					out.WriteByte(',')
				}
				kt, err := dec.Token()
				if err != nil {
					return err
				}
				kunci := kt.(string)
				tulisTeks(out, kunci)
				out.WriteByte(':')
				anak := jalur + "." + kunci
				kosong := buang(anak, kunci)
				if err := salinNilai(dec, out, lap, anak, kosong); err != nil {
					return err
				}
				if kosong && ada(JalurBuang, anak) {
					lap.Dibuang[anak]++
				} else if kosong {
					lap.Dibuang[kunci]++
				}
			}
			if _, err := dec.Token(); err != nil { // '}'
				return err
			}
			out.WriteByte('}')
		case '[':
			out.WriteByte('[')
			for i := 0; dec.More(); i++ {
				if i > 0 {
					out.WriteByte(',')
				}
				if err := salinNilai(dec, out, lap, fmt.Sprintf("%s[%d]", jalur, i), false); err != nil {
					return err
				}
			}
			if _, err := dec.Token(); err != nil { // ']'
				return err
			}
			out.WriteByte(']')
		}
	case string:
		if kosongkan {
			v = ""
		}
		tulisTeks(out, v)
	case json.Number:
		if kosongkan {
			return errors.New("deidentifikasi: nilai kunci BUANG berupa angka - tinjau manual")
		}
		out.WriteString(v.String())
	case bool:
		if kosongkan {
			return errors.New("deidentifikasi: nilai kunci BUANG berupa boolean - tinjau manual")
		}
		fmt.Fprintf(out, "%v", v)
	case nil:
		out.WriteString("null")
	}
	return nil
}

func tulisTeks(out *bytes.Buffer, s string) {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	out.Truncate(out.Len() - 1) // Encode menambah baris baru
}

// daun - satu nilai skalar beserta jalurnya.
type daun struct {
	jalur, kunci, teks string
	jenis              byte // 's' teks, 'n' angka, 'b' boolean, '0' null
}

func daunSemua(masuk []byte) ([]daun, error) {
	dec := json.NewDecoder(bytes.NewReader(masuk))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	var hasil []daun
	var jalan func(x any, jalur, kunci string)
	jalan = func(x any, jalur, kunci string) {
		switch t := x.(type) {
		case map[string]any:
			nama := make([]string, 0, len(t))
			for k := range t {
				nama = append(nama, k)
			}
			sort.Strings(nama)
			for _, k := range nama {
				jalan(t[k], jalur+"."+k, k)
			}
		case []any:
			for i, e := range t {
				jalan(e, fmt.Sprintf("%s[%d]", jalur, i), kunci)
			}
		case string:
			hasil = append(hasil, daun{jalur, kunci, t, 's'})
		case json.Number:
			hasil = append(hasil, daun{jalur, kunci, t.String(), 'n'})
		case bool:
			hasil = append(hasil, daun{jalur, kunci, fmt.Sprint(t), 'b'})
		case nil:
			hasil = append(hasil, daun{jalur, kunci, "", '0'})
		}
	}
	jalan(v, "$", "")
	return hasil, nil
}

// Periksa membandingkan `asli` dengan `bersih` secara otomatis:
//   - setiap daun di luar kunci BUANG IDENTIK (teks dan jenisnya), termasuk
//     setiap angka; dan setiap daun kunci BUANG kosong - selain itu galat;
//   - setiap nilai identitas asli (DaftarBuang / JalurBuang, ≥ 3 huruf) yang muncul lagi sebagai KATA UTUH di
//     daun teks lain dilaporkan sebagai kebocoran (jalur daunnya), bukan galat:
//     itu tanda medan identitas di luar daftar, untuk ditinjau manusia.
func Periksa(asli, bersih []byte) ([]string, error) {
	a, err := daunSemua(asli)
	if err != nil {
		return nil, err
	}
	b, err := daunSemua(bersih)
	if err != nil {
		return nil, err
	}
	if len(a) != len(b) {
		return nil, fmt.Errorf("deidentifikasi: bentuk berubah, %d daun asli vs %d daun bersih", len(a), len(b))
	}
	var rahasia []string
	for i := range a {
		if a[i].jalur != b[i].jalur {
			return nil, fmt.Errorf("deidentifikasi: jalur berubah %s vs %s", a[i].jalur, b[i].jalur)
		}
		if buang(a[i].jalur, a[i].kunci) {
			if b[i].teks != "" {
				return nil, fmt.Errorf("deidentifikasi: %s masih berisi", b[i].jalur)
			}
			if sumberRahasia(a[i].jalur, a[i].kunci) && len([]rune(strings.TrimSpace(a[i].teks))) >= 3 {
				rahasia = append(rahasia, strings.TrimSpace(a[i].teks))
			}
			continue
		}
		if a[i].teks != b[i].teks || a[i].jenis != b[i].jenis {
			return nil, fmt.Errorf("deidentifikasi: %s berubah", a[i].jalur)
		}
	}
	var bocor []string
	for _, r := range uniq(rahasia) {
		pola := regexp.MustCompile(`(?i)(^|[^\pL\pN])` + regexp.QuoteMeta(r) + `($|[^\pL\pN])`)
		for _, d := range b {
			if d.jenis == 's' && !buang(d.jalur, d.kunci) && pola.MatchString(d.teks) {
				bocor = append(bocor, d.jalur)
			}
		}
	}
	return uniq(bocor), nil
}

// KandidatIdentitas - NAMA kunci (bukan nilainya) yang bernilai teks bukan-angka
// dan tidak dibuang (DaftarBuang / JalurBuang), urut abjad, untuk ditinjau manusia. Tidak ada
// yang dibuang otomatis berdasarkan daftar ini.
func KandidatIdentitas(masuk []byte) ([]string, error) {
	d, err := daunSemua(masuk)
	if err != nil {
		return nil, err
	}
	angka := regexp.MustCompile(`^\s*-?[0-9.,]+\s*$`)
	var k []string
	for _, x := range d {
		if x.jenis == 's' && strings.TrimSpace(x.teks) != "" && !angka.MatchString(x.teks) && !buang(x.jalur, x.kunci) {
			k = append(k, x.kunci)
		}
	}
	return uniq(k), nil
}

func uniq(s []string) []string {
	sort.Strings(s)
	var h []string
	for i, x := range s {
		if i == 0 || x != s[i-1] {
			h = append(h, x)
		}
	}
	return h
}

// ErrBelumBersih - berkas masih berisi nilai di kunci atau jalur BUANG: bukan
// keluaran alat ini, mungkin berkas kasus mentah.
var ErrBelumBersih = errors.New("deidentifikasi: berkas belum dibersihkan")

// SudahBersih memeriksa bahwa setiap daun kunci/jalur BUANG kosong - syarat bagi
// pembaca yang hanya boleh menerima fixture ter-de-identifikasi (tiket 16).
func SudahBersih(isi []byte) error {
	d, err := daunSemua(isi)
	if err != nil {
		return err
	}
	for _, x := range d {
		if buang(x.jalur, x.kunci) && x.teks != "" {
			return fmt.Errorf("%w: %s berisi", ErrBelumBersih, x.jalur)
		}
	}
	return nil
}
