package services

// Pemeriksaan nama sebelum Create (perintah work owner 04-10-2026; title di nama kini DITOLAK isiOrg - "ga boleh
// disimpan dong"): "sebelum create jika namanya mengandung title nya
// kasih warning; lalu jika namanya sudah pernah terdaftar atau memiliki kemiripan kasih popup ... lalu ada pilihan
// lanjut save / tidak". Layanan ini hanya MENJAWAB temuannya; keputusan lanjut atau batal ada di layar - Create tidak
// pernah ditolak karena nama mirip.

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"nusantarare/modul/companydetail/backend/models"
)

// Ambang dan batas pemeriksaan nama.
const (
	// AmbangMirip - kemiripan huruf (1 - jarak edit / panjang terpanjang) minimum untuk "mirip".
	AmbangMirip = 0.8
	// MinHurufMemuat - nama yang lebih pendek minimal sepanjang ini (tanpa spasi) untuk dianggap "memuat".
	MinHurufMemuat = 5
	// BatasMirip - temuan terbanyak yang dikirim ke layar.
	BatasMirip = 10
)

// Jenis temuan nama.
const (
	NamaSama  = "sama"
	NamaMirip = "mirip"
)

// NamaSerupa adalah satu organisasi yang namanya sama atau mirip.
type NamaSerupa struct {
	models.BarisDaftar
	// Jenis - `sama` (sama sesudah title dan tanda baca dibuang) atau `mirip`.
	Jenis string `json:"jenis"`
	// Skor - kemiripan 0..1.
	Skor float64 `json:"skor"`
}

// PeriksaNama adalah jawaban pemeriksaan nama.
type PeriksaNama struct {
	// TitleDalamNama - title (LABEL, mis. `PT.`) yang ditemukan sebagai kata di nama; kosong = tidak ada.
	TitleDalamNama string       `json:"titleDalamNama"`
	Serupa         []NamaSerupa `json:"serupa"`
}

// kataTitle - kata title AKTIF tanpa titik, huruf besar (`PT.` -> `PT`): PT., CV., PD., UD. Title nonaktif TN., NY.,
// NN. adalah sapaan orang (dan `NY` wajar di nama perusahaan) - tidak dihitung.
func kataTitle(r referensi) map[string]string {
	hasil := map[string]string{}
	for _, p := range r.enum[models.JenisTitle] {
		if !p.Aktif {
			continue
		}
		k := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(p.Label), ".", ""))
		if k != "" {
			hasil[k] = p.Label
		}
	}
	return hasil
}

// kataNama memecah nama menjadi kata huruf besar: titik dibuang (`P.T.` -> `PT`), tanda baca lain pemisah kata.
func kataNama(nama string) []string {
	nama = strings.ToUpper(strings.ReplaceAll(nama, ".", ""))
	return strings.FieldsFunc(nama, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}

// NormalNama - nama untuk dibandingkan: kata huruf besar tanpa kata title, dipisah satu spasi.
func NormalNama(nama string, title map[string]string) string {
	var inti []string
	for _, k := range kataNama(nama) {
		if _, adalahTitle := title[k]; !adalahTitle {
			inti = append(inti, k)
		}
	}
	return strings.Join(inti, " ")
}

// TitleDiNama - LABEL title pertama yang muncul sebagai KATA di nama (`PT ABC`, `P.T. ABC`, `ABC, PT`).
func TitleDiNama(nama string, title map[string]string) string {
	for _, k := range kataNama(nama) {
		if label, ada := title[k]; ada {
			return label
		}
	}
	return ""
}

// jarakEdit - jarak Levenshtein dua teks (per huruf).
func jarakEdit(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	sebelum := make([]int, len(b)+1)
	kini := make([]int, len(b)+1)
	for j := range sebelum {
		sebelum[j] = j
	}
	for i := 1; i <= len(a); i++ {
		kini[0] = i
		for j := 1; j <= len(b); j++ {
			biaya := 1
			if a[i-1] == b[j-1] {
				biaya = 0
			}
			kini[j] = min(sebelum[j]+1, kini[j-1]+1, sebelum[j-1]+biaya)
		}
		sebelum, kini = kini, sebelum
	}
	return sebelum[len(b)]
}

// KemiripanNama - jenis dan skor dua nama yang SUDAH dinormalkan; jenis kosong = tidak serupa.
func KemiripanNama(a, b string) (string, float64) {
	if a == "" || b == "" {
		return "", 0
	}
	if a == b {
		return NamaSama, 1
	}
	ra, rb := []rune(a), []rune(b)
	panjang := max(len(ra), len(rb))
	skor := 1 - float64(jarakEdit(ra, rb))/float64(panjang)
	if skor >= AmbangMirip {
		return NamaMirip, skor
	}
	// Satu nama memuat nama lain sebagai rangkaian kata utuh (`ABC` di `ABC INDONESIA`).
	pendek, pajang := a, b
	if len(pendek) > len(pajang) {
		pendek, pajang = pajang, pendek
	}
	if len(strings.ReplaceAll(pendek, " ", "")) >= MinHurufMemuat && strings.Contains(" "+pajang+" ", " "+pendek+" ") {
		return NamaMirip, skor
	}
	return "", 0
}

// Periksa memeriksa nama organisasi sebelum disimpan: title di dalam nama, dan organisasi lain bernama sama atau
// mirip (paling banyak BatasMirip, yang sama lebih dulu). `kecuali` = ID organisasi yang sedang diubah.
func (l *Layanan) Periksa(ctx context.Context, nama, kecuali string) (PeriksaNama, error) {
	hasil := PeriksaNama{Serupa: []NamaSerupa{}}
	nama = strings.TrimSpace(nama)
	if nama == "" {
		return hasil, nil
	}
	r, _, _, err := l.muatReferensi(ctx)
	if err != nil {
		return hasil, err
	}
	title := kataTitle(r)
	hasil.TitleDalamNama = TitleDiNama(nama, title)
	normal := NormalNama(nama, title)
	if normal == "" {
		return hasil, nil
	}
	semua, err := l.gudang.DaftarNamaOrg(ctx)
	if err != nil {
		return hasil, err
	}
	for _, o := range semua {
		if o.ID == kecuali {
			continue
		}
		if jenis, skor := KemiripanNama(normal, NormalNama(o.Nama, title)); jenis != "" {
			hasil.Serupa = append(hasil.Serupa, NamaSerupa{BarisDaftar: o, Jenis: jenis, Skor: skor})
		}
	}
	sort.SliceStable(hasil.Serupa, func(i, j int) bool {
		a, b := hasil.Serupa[i], hasil.Serupa[j]
		if (a.Jenis == NamaSama) != (b.Jenis == NamaSama) {
			return a.Jenis == NamaSama
		}
		if a.Skor != b.Skor {
			return a.Skor > b.Skor
		}
		return a.IDView < b.IDView
	})
	if len(hasil.Serupa) > BatasMirip {
		hasil.Serupa = hasil.Serupa[:BatasMirip]
	}
	return hasil, nil
}
