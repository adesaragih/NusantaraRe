// Program bangkit menulis `registry_gen.go` paket rules dari ekspor rule `When`
// korpus Pega. Ia dijalankan tangan, BUKAN lewat `go generate`: korpus tidak ada
// di repositori, dan jalurnya konfigurasi, bukan literal.
//
//	go run ./modul/nbfacin/backend/services/rules/bangkit \
//	    -akar "D:\migrasi\RNM" \
//	    -keluar modul/nbfacin/backend/services/rules/registry_gen.go
//
// `-akar` adalah akar korpus. Yang dibaca: folder `folderUtama`, ditambah berkas
// `pinjaman` (daftar eksplisit, masing-masing berdasar keputusan work owner).
// Penjaganya `TestRegistrySamaDenganKeluaranBangkit`: dengan env
// `KORPUS_RNM_AKAR` terisi, ia membangkitkan ulang di memori dan menolak
// `registry_gen.go` yang berbeda.
//
// Folder siklus lain (opsi (b) registry predikat EDM, keputusan work owner
// 01-10-2026, butir 56 `docs/KEPUTUSAN-30-09-2026.md`): `-folder` dan `-paket`.
// Folder selain folderUtama dibaca APA ADANYA - tanpa pinjaman, catatan sikap, dan
// sikap khusus NB (CLAUDE.md §4.6) - mis.
//
//	go run ./modul/nbfacin/backend/services/rules/bangkit \
//	    -akar "D:\migrasi\RNM" -folder "Endorsment Fac In\When" -paket edm \
//	    -keluar <berkas di paket tujuan>
//
// ⚠️ Keluarannya merujuk tipe paket rules (`predikat`, `kondisi`, `operand`, …):
// paket tujuan wajib memuat mesin yang sama. Mesin TIDAK dipindah ke inti: NB dan
// EDM modul terpisah (butir 57) - modul tujuan memegang mesin evaluatornya sendiri.
// [terverifikasi] keluaran folder EDM terkompilasi terhadap mesin paket ini
// (`go build -overlay`, 01-10-2026).
//
// Yang dibaca per berkas HANYA: `<pxInsName>` (identitas), `<pyLogic>` dan
// `<pyCondition>/<rowdata>` anak langsung akar (`pyConditionLabel`,
// `pyConditionValue1` - ekspresi TERSIMPAN). `<pyConditionString>` (teks
// tampilan, 55 berkas NB berisi placeholder) sengaja tidak dibaca (NB-08).
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"nusantarare/modul/nbfacin/backend/services/rules/logika"
)

// folderUtama - folder `When` siklus NB, relatif terhadap akar korpus.
const folderUtama = `NB FacIn\When`

// pinjaman - berkas `When` dari folder siklus lain yang dimuat karena keputusan
// work owner. Hanya boleh untuk nama yang TIDAK ada di folderUtama (CLAUDE.md
// §4.6: varian modul lain tidak dipinjam tanpa keputusan).
var pinjaman = []struct{ jalur, dasar string }{
	// K-003: dirujuk SpreadingAdditionalProtection.xml sebagai pra-kondisi,
	// tidak ada di NB FacIn\When, terbaca penuh di folder Endorsment.
	{`Endorsment Fac In\When\IsSpreadingDepan.xml`, "K-003"},
}

// catatanSikap - catatan sikap khusus per predikat (tiket 09).
var catatanSikap = map[string]string{
	"ISOFFERFACIN": "K-002: ekspresi tersimpan yang berlaku. Teks tampilan <pyConditionString> menyebut " +
		`Kode Bisnis = "02"/"58"/"SB"/"SG" - kandidat perbaikan, TIDAK dieksekusi; tersangka pertama ` +
		"bila paralel run berselisih pada gerbang masuk NB.",
	"ISFACOUT": "K-019: nama dipertahankan demi ketertelusuran, tetapi isinya menguji " +
		"ProposalAcceptStatus = 4 (Banding), BUKAN fac out. Fitur usang secara bisnis, kode tetap diport.",
}

// sikapKhusus - predikat yang sikapnya ditetapkan spec, bukan kondisinya.
var sikapKhusus = map[string]string{
	// spec Modul 3: kondisi terbaca tetapi label belum ter-resolve, tampilan
	// placeholder, pyTempText true - belum terbukti sebagai yang dieksekusi.
	"ISPKSASM": "sikap panic (spec Modul 3, tiket 09): kondisi belum terbukti dieksekusi",
}

// profil - apa yang dibaca dan ke paket mana ditulis.
type profil struct {
	folder, paket string
	pinjaman      []struct{ jalur, dasar string }
	catatan       map[string]string // catatanSikap
	sikap         map[string]string // sikapKhusus
}

var polaPaket = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// pilihProfil - folderUtama → profil NB lengkap (hanya ke paket rules); folder
// lain → tanpa pinjaman, catatan, dan sikap NB.
func pilihProfil(folder, paket string) (profil, error) {
	// `NB FacIn/When` dan `NB FacIn\When\` = folder yang sama.
	folder = strings.ReplaceAll(strings.TrimRight(folder, `/\`), "/", `\`)
	switch {
	case folder == "":
		return profil{}, fmt.Errorf("bangkit: -folder kosong")
	case !polaPaket.MatchString(paket):
		return profil{}, fmt.Errorf("bangkit: -paket %q bukan nama paket Go huruf kecil", paket)
	case folder == folderUtama && paket != "rules":
		return profil{}, fmt.Errorf("bangkit: registry %s hanya untuk paket rules nbfacin", folderUtama)
	case folder == folderUtama:
		return profil{folder: folder, paket: paket, pinjaman: pinjaman, catatan: catatanSikap, sikap: sikapKhusus}, nil
	}
	return profil{folder: folder, paket: paket, catatan: map[string]string{}, sikap: map[string]string{}}, nil
}

type berkasWhen struct {
	InsName string `xml:"pxInsName"`
	Versi   string `xml:"pyRuleSetVersion"`
	Logika  string `xml:"pyLogic"`
	Baris   []struct {
		Label  string `xml:"pyConditionLabel"`
		Nilai1 string `xml:"pyConditionValue1"`
	} `xml:"pyCondition>rowdata"`
}

var (
	polaRujuk   = regexp.MustCompile(`^@\(Pega-RULES:ExpressionEvaluators\)\.evaluateWhen\("([^"]+)"\)$`)
	awalBanding = "@(Pega-RULES:ExpressionEvaluators).compareTwoValues("
	polaJalur   = regexp.MustCompile(`^\.?[A-Za-z_]\w*(\(\d+\))?(\.[A-Za-z_]\w*(\(\d+\))?)*$`)
	polaAngka   = regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`)
	// polaKata - satu kata tanpa titik/kurung; yang bertitik tetap jalur properti.
	polaKata = regexp.MustCompile(`^[A-Za-z_]\w*$`)
	// medanLogin - segmen akhir sisi kiri yang berisi login operator. Nilainya
	// nama orang, jadi tidak disalin (CLAUDE.md §4 butir 10) - termasuk ke pesan panik.
	medanLogin = map[string]bool{"pyUserIdentifier": true, "pxCreateOperator": true, "pxUpdateOperator": true}
	// medanIdentitas - segmen akhir sisi kiri yang literal TIDAK-KOSONG-nya tidak
	// disalin (literal kosong - cek isian, mis. NOPOLISEMPTY - tetap diport):
	//   - nomor polis dan nama marketing (nama orang): data identitas, CLAUDE.md §4
	//     butir 10 dan §10 - dua nomor polis (ISERRORSPREADING) dan tiga nama
	//     (ISTBONDING) sempat tersalin, tidak pernah masuk riwayat git (dicek 01-10-2026);
	//   - `OperatorID.pyTelephone` (ISTREATY1, ISSPVTREATY1): BUKAN telepon pribadi -
	//     kode peran yang disimpan di medan telepon operator (dibahas di discovery
	//     sebagai RBAC). Seperti medanLogin, diport lewat model peran, bukan literal.
	medanIdentitas = map[string]bool{"OldPolicyNo": true, "PolicyNo": true, "MarketingName": true, "pyTelephone": true}
)

// polaMedanRahasia - ekspresi yang MENYEBUT medan login/identitas (huruf besar-kecil
// diabaikan), dalam bentuk apa pun: pesan paniknya tidak menyalin ekspresi itu.
var polaMedanRahasia = regexp.MustCompile(`(?i)\.(oldpolicyno|policyno|marketingname|pytelephone|pyuseridentifier|pxcreateoperator|pxupdateoperator)\b`)

// cocokMedan - segmen akhir jalur ada di `m`, tanpa peduli huruf besar-kecil.
func cocokMedan(m map[string]bool, seg string) bool {
	for k := range m {
		if strings.EqualFold(k, seg) {
			return true
		}
	}
	return false
}

const (
	alasanLogin     = "membandingkan login operator dengan literal; nilainya tidak disalin, diport lewat model peran"
	alasanIdentitas = "membandingkan nomor polis, nama, atau atribut peran operator dengan literal; nilainya tidak disalin"
)

// kondisiGen - satu baris kondisi yang sudah dikenali.
type kondisiGen struct {
	literal string // literal Go `kondisi{...}`
	rujukan string // nama predikat yang dirujuk, huruf besar; kosong bila banding
}

// entri - satu predikat siap ditulis.
type entri struct {
	nama    string
	asal    []string
	logika  string
	kondisi map[string]kondisiGen
	dibuang []string // label yang tidak dirujuk logika
	panik   string
	catatan string
	// sidikIsi - logika + baris kondisi yang dirujuk; pembanding identitas
	// yang tersimpan di lebih dari satu berkas.
	sidikIsi string
}

func main() {
	akar := flag.String("akar", "", "akar korpus Pega (berisi folder siklus)")
	keluar := flag.String("keluar", "", "berkas Go yang ditulis")
	folder := flag.String("folder", folderUtama, "folder When relatif terhadap -akar")
	paket := flag.String("paket", "rules", "nama paket Go keluaran")
	flag.Parse()
	if *akar == "" || *keluar == "" {
		log.Fatal("bangkit: -akar dan -keluar wajib")
	}
	p, err := pilihProfil(*folder, *paket)
	if err != nil {
		log.Fatal(err)
	}
	src, ringkas, err := bangkitkan(*akar, p)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*keluar, src, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bangkit: %s → %s\n", ringkas, *keluar)
}

// bangkitkan membaca seluruh berkas `When` folder profil di bawah `akar`, ditambah
// pinjamannya, dan mengembalikan sumber Go registry beserta ringkasannya.
func bangkitkan(akar string, p profil) ([]byte, string, error) {
	korpus, label := filepath.Join(akar, p.folder), p.folder
	nama, err := os.ReadDir(korpus)
	if err != nil {
		return nil, "", err
	}
	semua := map[string]*entri{}
	berkas := 0
	var urutan []*entri
	for _, d := range nama {
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".xml") {
			continue
		}
		berkas++
		e, err := baca(filepath.Join(korpus, d.Name()), label+`\`+d.Name(), "", p.sikap)
		if err != nil {
			return nil, "", fmt.Errorf("bangkit: %s: %w", d.Name(), err)
		}
		if lama, ada := semua[e.nama]; ada {
			if lama.sidikIsi != e.sidikIsi {
				return nil, "", fmt.Errorf("bangkit: identitas %s di %v dan %v BERBEDA isi - tidak disatukan (CLAUDE.md §4.6)", e.nama, lama.asal, e.asal)
			}
			lama.asal = append(lama.asal, e.asal...)
			continue
		}
		semua[e.nama] = e
		urutan = append(urutan, e)
	}
	for _, pj := range p.pinjaman {
		e, err := baca(filepath.Join(akar, pj.jalur), pj.jalur, pj.dasar, p.sikap)
		if err != nil {
			return nil, "", fmt.Errorf("bangkit: pinjaman %s: %w", pj.jalur, err)
		}
		if _, ada := semua[e.nama]; ada {
			return nil, "", fmt.Errorf("bangkit: pinjaman %s bernama %s, padahal nama itu ada di %s", pj.jalur, e.nama, p.folder)
		}
		semua[e.nama] = e
		urutan = append(urutan, e)
	}
	for _, e := range urutan {
		e.catatan = p.catatan[e.nama]
	}
	// Rujukan ke predikat yang tidak ada di korpus → panic saat dievaluasi.
	for _, e := range urutan {
		if e.panik != "" {
			continue
		}
		for _, l := range urutLabel(e.kondisi) {
			if r := e.kondisi[l].rujukan; r != "" {
				if _, ada := semua[r]; !ada {
					e.panik = fmt.Sprintf("merujuk predikat %s yang tidak ada di korpus", r)
					break
				}
			}
		}
	}
	sort.Slice(urutan, func(i, j int) bool { return urutan[i].nama < urutan[j].nama })
	src, err := tulis(urutan, berkas, label, p.paket, len(p.pinjaman))
	if err != nil {
		return nil, "", err
	}
	return src, fmt.Sprintf("%d berkas %s + %d pinjaman, %d predikat", berkas, label, len(p.pinjaman), len(urutan)), nil
}

// baca membaca satu berkas `When`. `dasarPinjam` terisi bila berkas ini pinjaman
// dari folder lain; asalnya lalu menyebut keputusan dan versi ruleset-nya.
func baca(jalur, asal, dasarPinjam string, sikap map[string]string) (*entri, error) {
	isi, err := os.ReadFile(jalur)
	if err != nil {
		return nil, err
	}
	var b berkasWhen
	if err := xml.Unmarshal(isi, &b); err != nil {
		return nil, err
	}
	kelasNama := strings.SplitN(b.InsName, "!", 2)
	if len(kelasNama) != 2 {
		return nil, fmt.Errorf("pxInsName %q", b.InsName)
	}
	asal += " (" + b.InsName + ")"
	if dasarPinjam != "" {
		asal += fmt.Sprintf(" - dipinjam %s, pyRuleSetVersion %s", dasarPinjam, strings.TrimSpace(b.Versi))
	}
	e := &entri{nama: strings.ToUpper(kelasNama[1]), asal: []string{asal},
		logika: strings.TrimSpace(b.Logika), kondisi: map[string]kondisiGen{}}
	if s, ada := sikap[e.nama]; ada {
		e.panik = s
	}
	dirujuk := map[string]bool{}
	if pohon, err := logika.Urai(e.logika); err != nil {
		if e.panik == "" {
			e.panik = fmt.Sprintf("logika %q tidak terurai: %v", e.logika, err)
		}
	} else {
		for _, l := range pohon.Label() {
			dirujuk[l] = true
		}
	}
	var sidik []string
	for _, br := range b.Baris {
		label, nilai := strings.TrimSpace(br.Label), strings.TrimSpace(br.Nilai1)
		if !dirujuk[label] {
			e.dibuang = append(e.dibuang, label)
			continue
		}
		sidik = append(sidik, label+"="+nilai)
		kd, alasan := kenaliKondisi(nilai)
		if alasan != "" && e.panik == "" {
			e.panik = fmt.Sprintf("baris %s belum diport (%s): %s", label, alasan, nilai)
			// Ekspresi yang menyebut medan login/identitas tidak ditulis - apa pun alasannya
			// (mis. "bentuk ekspresi lain" berisi @String.equals atas nomor polis).
			if alasan == alasanLogin || alasan == alasanIdentitas || polaMedanRahasia.MatchString(nilai) {
				e.panik = fmt.Sprintf("baris %s belum diport (%s)", label, alasan)
			}
		}
		e.kondisi[label] = kd
		delete(dirujuk, label)
	}
	if len(dirujuk) > 0 && e.panik == "" {
		var hilang []string
		for l := range dirujuk {
			hilang = append(hilang, l)
		}
		sort.Strings(hilang)
		e.panik = fmt.Sprintf("logika merujuk label tanpa baris kondisi: %v", hilang)
	}
	sort.Strings(sidik)
	e.sidikIsi = e.logika + "|" + strings.Join(sidik, "|")
	return e, nil
}

// kenaliKondisi - baris kondisi yang diport, atau alasan belum diport.
func kenaliKondisi(nilai string) (kondisiGen, string) {
	if m := polaRujuk.FindStringSubmatch(nilai); m != nil {
		r := strings.ToUpper(m[1])
		return kondisiGen{literal: fmt.Sprintf("{jenis: rujukWhen, rujukan: %s}", strconv.Quote(r)), rujukan: r}, ""
	}
	if !strings.HasPrefix(nilai, awalBanding) || !strings.HasSuffix(nilai, ")") {
		return kondisiGen{}, "bentuk ekspresi lain"
	}
	arg := pecahArgumen(nilai[len(awalBanding) : len(nilai)-1])
	if len(arg) != 3 {
		return kondisiGen{}, "jumlah argumen compareTwoValues"
	}
	if !polaJalur.MatchString(arg[0]) {
		return kondisiGen{}, "sisi kiri bukan jalur properti"
	}
	seg := strings.Split(arg[0], ".")
	if cocokMedan(medanLogin, seg[len(seg)-1]) {
		return kondisiGen{}, alasanLogin
	}
	if cocokMedan(medanIdentitas, seg[len(seg)-1]) && arg[2] != `""` {
		return kondisiGen{}, alasanIdentitas
	}
	op, err := strconv.Unquote(arg[1])
	if err != nil || (op != "=" && op != "!=") {
		return kondisiGen{}, "operator " + arg[1]
	}
	var kanan string
	switch k := arg[2]; {
	case strings.HasPrefix(k, `"`):
		s, err := strconv.Unquote(k)
		if err != nil {
			return kondisiGen{}, "literal teks"
		}
		kanan = fmt.Sprintf("operand{teksLiteral, %s}", strconv.Quote(s))
	case polaAngka.MatchString(k):
		kanan = fmt.Sprintf("operand{angkaLiteral, %s}", strconv.Quote(k))
	case k == "true" || k == "false":
		kanan = fmt.Sprintf("operand{booleanLiteral, %s}", strconv.Quote(k))
	case strings.EqualFold(k, "true") || strings.EqualFold(k, "false"):
		// `True`: boolean atau teks "True" belum terverifikasi; booleanLiteral
		// membandingkan teks huruf kecil, jadi keduanya tidak setara.
		return kondisiGen{}, "operand kanan boolean berhuruf besar " + k
	case polaKata.MatchString(k):
		// Kata telanjang (`ReasFacInDirector`, `A1`, `Offer`) dibaca sebagai teks
		// (keputusan work owner 01-10-2026, butir 28).
		kanan = fmt.Sprintf("operand{teksLiteral, %s}", strconv.Quote(k))
	default:
		return kondisiGen{}, "operand kanan bukan literal " + k
	}
	return kondisiGen{literal: fmt.Sprintf("{jenis: banding, kiri: %s, tidakSama: %v, kanan: %s}", strconv.Quote(arg[0]), op == "!=", kanan)}, ""
}

// pecahArgumen memecah argumen di koma tingkat atas, menghormati kutip dan kurung.
func pecahArgumen(s string) []string {
	var hasil []string
	var buf strings.Builder
	kutip, kurung := false, 0
	for _, ch := range s {
		switch {
		case ch == '"':
			kutip = !kutip
		case !kutip && ch == '(':
			kurung++
		case !kutip && ch == ')':
			kurung--
		case !kutip && kurung == 0 && ch == ',':
			hasil = append(hasil, strings.TrimSpace(buf.String()))
			buf.Reset()
			continue
		}
		buf.WriteRune(ch)
	}
	return append(hasil, strings.TrimSpace(buf.String()))
}

func urutLabel(k map[string]kondisiGen) []string {
	var label []string
	for l := range k {
		label = append(label, l)
	}
	sort.Strings(label)
	return label
}

func tulis(urutan []*entri, berkas int, label, paket string, nPinjam int) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by ./bangkit DO NOT EDIT.\n\n")
	fmt.Fprintf(&b, "// Sumber: %s - %d berkas + %d pinjaman, %d predikat (identitas `<pxInsName>`, nama huruf besar).\n", label, berkas, nPinjam, len(urutan))
	fmt.Fprintf(&b, "// Ekspresi dibaca dari `<pyLogic>` + `<pyConditionValue1>`, bukan `<pyConditionString>`.\n\n")
	fmt.Fprintf(&b, "package %s\n\nvar registry = map[string]predikat{\n", paket)
	for _, e := range urutan {
		for _, a := range e.asal {
			fmt.Fprintf(&b, "\t// Asal: %s\n", a)
		}
		if len(e.dibuang) > 0 {
			fmt.Fprintf(&b, "\t// Baris tidak dirujuk pyLogic, tidak dimuat: %v\n", e.dibuang)
		}
		fmt.Fprintf(&b, "\t%s: {\n\t\tasal: %#v,\n\t\tlogika: %s,\n", strconv.Quote(e.nama), e.asal, strconv.Quote(e.logika))
		if e.catatan != "" {
			fmt.Fprintf(&b, "\t\tcatatan: %s,\n", strconv.Quote(e.catatan))
		}
		if e.panik != "" {
			fmt.Fprintf(&b, "\t\tpanik: %s,\n", strconv.Quote(e.panik))
		} else if len(e.kondisi) > 0 {
			fmt.Fprintf(&b, "\t\tkondisi: map[string]kondisi{\n")
			for _, l := range urutLabel(e.kondisi) {
				fmt.Fprintf(&b, "\t\t\t%s: %s,\n", strconv.Quote(l), e.kondisi[l].literal)
			}
			fmt.Fprintf(&b, "\t\t},\n")
		}
		fmt.Fprintf(&b, "\t},\n")
	}
	fmt.Fprintf(&b, "}\n")
	return format.Source(b.Bytes())
}
