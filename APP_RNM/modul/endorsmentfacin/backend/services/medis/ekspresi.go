package medis

// Penafsir subset ekspresi Pega yang dipakai `CalculateScorLife_Act` - tiket E20.
//
// Untuk apa berkas ini: ekspresi skoring medis dijalankan APA ADANYA dari
// teks korpus (aturan_gen.go), bukan diterjemahkan tangan satu per satu.
// Menerjemahkan ratusan `@if` bersarang dengan tangan adalah jalan paling
// pasti untuk salah salin ambang klinis (keputusan A04).
//
// Dibaca sesudah: aturan_gen.go. Dibaca sebelum: medis.go.
//
// [terverifikasi] Subset yang dipakai - cacah atas seluruh 728 ekspresi
// (prakondisi + penugasan) activity itu. Perintah audit (Git Bash, di
// `Endorsment Fac In/Activity`):
//
//	grep -c '<pyStepsPreCondParamsWhen>[^<]' CalculateScorLife_Act.xml   # 390 prakondisi
//	grep -c '<pyStepsPreCondParamsWhenTrue>[0-9]' CalculateScorLife_Act.xml  # 390 - cara kedua (170 "5" + 220 "2")
//	grep -c '<PropertiesName>' CalculateScorLife_Act.xml                 # 331 penugasan bernama
//	grep -c '<PropertiesName/>' CalculateScorLife_Act.xml                # 7 penugasan kosong
//
// (390 + 338 = 728.) Isinya:
//   - fungsi: `@if` · `@toDecimal` · `@replaceAll` · `@divide` (hanya empat);
//     `CalculatePhysicalExam` menambah `@divide` DUA argumen dan
//     `@StrategyUtils.pow`;
//   - operator: `==` `=` `!=` `<` `<=` `>` `>=` `&&` `||`, kurung, koma;
//     TANPA aritmetika;
//   - rujukan: halaman langkah (`.X`), `param.X`, `ParamLab.X`;
//   - literal teks berkutip dan angka desimal.
// Yang di luar subset ini ditolak sebagai galat sintaks, bukan ditebak.
//
// ⚠️ Semantik Pega yang belum terverifikasi, dan sikapnya di sini:
//   - nilai kosong di konteks angka (`@toDecimal("")`, `.X>=6` dengan X
//     kosong) → `ErrNilaiBukanAngka`, tidak ditebak nol;
//   - mode pembulatan `@divide(a,b,n)` → ekspresi dijalankan dengan
//     beberapa mode; hasil yang berbeda antar mode → `ErrBergantungModePembulatan`;
//   - `@if` dievaluasi MALAS (cabang yang tidak dipilih tidak dijalankan)
//     `[dugaan]`; `&&`/`||` dipotong pendek seperti Java `[dugaan]`;
//   - `=` dan `==` sama-sama kesetaraan. Tipe properti tidak ada di korpus
//     (nol rule Property), jadi perbandingan mengikuti sikap registry NB
//     (keputusan work owner 01-10-2026 butir 20): dua nilai BERTIPE angka
//     (hasil `@toDecimal`/`@divide`/`@StrategyUtils.pow` atau literal angka)
//     dibandingkan sebagai angka; bila salah satunya teks yang terbaca angka,
//     tafsir teks dan tafsir angka dihitung dua-duanya - berbeda →
//     `ErrTafsirBerbeda`, tidak ditebak;
//   - teks kosong pada `=`/`!=` dibandingkan sebagai teks (tidak ambigu);
//   - argumen ketiga `@divide` dibaca sebagai jumlah desimal `[dugaan]`;
//     hasil bagi dihitung dulu pada presisi 38 lalu dibulatkan ke skala itu -
//     pembulatan ganda hanya berbeda bila digit 39 dst. tepat di batas,
//     diabaikan `[dugaan]`.
//
// ⛔ PRIVASI (E20): teks galat TIDAK PERNAH memuat ekspresi (yang berisi
// ambang klinis) maupun nilai (yang berisi hasil lab) - hanya jenis galat,
// nama fungsi, nomor argumen, dan posisi. Pemanggil yang mencatat galat
// tidak dapat membocorkan data medis lewat paket ini.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

var (
	// ErrSintaks - ekspresi di luar subset yang dikenali.
	ErrSintaks = errors.New("medis: ekspresi tidak dapat diurai")
	// ErrNilaiBukanAngka - konteks angka menerima teks kosong atau bukan angka.
	ErrNilaiBukanAngka = errors.New("medis: nilai bukan angka (perilaku Pega atas nilai kosong/bukan angka belum terverifikasi)")
	// ErrBagiNol - `@divide` dengan penyebut nol.
	ErrBagiNol = errors.New("medis: @divide dengan penyebut nol")
	// ErrBergantungModePembulatan - hasil berubah menurut mode pembulatan
	// `@divide`, yang belum terverifikasi.
	ErrBergantungModePembulatan = errors.New("medis: hasil bergantung mode pembulatan @divide yang belum terverifikasi")
	// ErrSkalaBawaanBelumTerverifikasi - `@divide(a,b)` tanpa skala memberi
	// hasil tak eksak; skala dan pembulatan bawaan Pega belum terverifikasi.
	ErrSkalaBawaanBelumTerverifikasi = errors.New("medis: @divide tanpa skala tidak eksak; skala bawaan Pega belum terverifikasi")
	// ErrPangkat - `@StrategyUtils.pow` dengan pangkat bukan bilangan bulat
	// 0..100, atau hasil yang melampaui presisi.
	ErrPangkat = errors.New("medis: @StrategyUtils.pow hanya untuk pangkat bulat 0..100 yang eksak")
	// ErrTafsirBerbeda - perbandingan yang hasilnya berbeda menurut tafsir
	// teks dan tafsir angka, padahal tipe properti belum terverifikasi.
	ErrTafsirBerbeda = errors.New("medis: tafsir teks dan angka berbeda, tipe properti belum terverifikasi")
	// ErrBukanBoolean - syarat `@if` atau operand `&&`/`||` bukan benar/salah.
	ErrBukanBoolean = errors.New("medis: syarat bukan benar/salah")
)

// halaman - sumber nilai rujukan. Galat = properti yang penugasannya sendiri
// gagal; galat itu menjalar ke ekspresi yang membacanya.
type halaman interface {
	ambil(nama string) (string, error)
}

// halamanBertipe - halaman yang tahu properti mana berisi nilai bertipe
// angka (diisi hasil fungsi angka). Properti lain dibaca sebagai teks.
type halamanBertipe interface {
	bertipeAngka(nama string) bool
}

type jenisNilai int

const (
	nilaiTeks jenisNilai = iota
	nilaiAngka
	nilaiBool
)

type nilai struct {
	jenis jenisNilai
	teks  string
	angka *apd.Decimal
	benar bool
}

func (n nilai) String() string {
	switch n.jenis {
	case nilaiAngka:
		return utils.FormatDecimal(n.angka)
	case nilaiBool:
		if n.benar {
			return "true"
		}
		return "false"
	}
	return n.teks
}

// sebagaiAngka - nilai angka, atau teks yang terbaca desimal.
func (n nilai) sebagaiAngka() (*apd.Decimal, bool) {
	switch n.jenis {
	case nilaiAngka:
		return n.angka, true
	case nilaiTeks:
		if n.teks == "" {
			return nil, false
		}
		d, err := utils.ParseDecimal(n.teks)
		return d, err == nil
	}
	return nil, false
}

func (n nilai) sebagaiBool() (bool, error) {
	switch {
	case n.jenis == nilaiBool:
		return n.benar, nil
	case n.jenis == nilaiTeks && n.teks == "true":
		return true, nil
	case n.jenis == nilaiTeks && n.teks == "false":
		return false, nil
	}
	return false, ErrBukanBoolean
}

// modeUji - mode pembulatan yang dicoba bila `@divide` membuang digit.
var modeUji = []apd.Rounder{
	apd.RoundHalfUp, apd.RoundHalfEven, apd.RoundHalfDown,
	apd.RoundDown, apd.RoundUp, apd.RoundCeiling, apd.RoundFloor,
}

// konteks - satu evaluasi: sumber nilai, mode pembulatan, dan penanda apakah
// ada pembulatan yang membuang digit.
type konteks struct {
	h          halaman
	mode       apd.Rounder
	dibulatkan bool
}

// evaluasiTeks - jalankan ekspresi, hasil sebagai teks. Bila ada `@divide`
// yang membuang digit, ekspresi diulang dengan setiap mode di `modeUji`;
// hasilnya harus sama.
func evaluasiTeks(ekspresi string, h halaman) (string, error) {
	s, _, err := evaluasiNilai(ekspresi, h)
	return s, err
}

// evaluasiNilai - seperti evaluasiTeks, ditambah apakah hasilnya bertipe angka.
func evaluasiNilai(ekspresi string, h halaman) (teks string, angka bool, err error) {
	n, err := urai(ekspresi)
	if err != nil {
		return "", false, err
	}
	k := &konteks{h: h, mode: modeUji[0]}
	v, err := n.eval(k)
	if err != nil {
		return "", false, err
	}
	if !k.dibulatkan {
		return v.String(), v.jenis == nilaiAngka, nil
	}
	pertama := v.String()
	for _, m := range modeUji[1:] {
		v, err := n.eval(&konteks{h: h, mode: m})
		if err != nil || v.String() != pertama {
			return "", false, ErrBergantungModePembulatan
		}
	}
	return pertama, v.jenis == nilaiAngka, nil
}

// evaluasiBool - ekspresi prakondisi.
func evaluasiBool(ekspresi string, h halaman) (bool, error) {
	s, err := evaluasiTeks(ekspresi, h)
	if err != nil {
		return false, err
	}
	return nilai{jenis: nilaiTeks, teks: s}.sebagaiBool()
}

// ---- pohon ekspresi ----

type simpul interface {
	eval(k *konteks) (nilai, error)
}

type literalTeks string

func (l literalTeks) eval(*konteks) (nilai, error) {
	return nilai{jenis: nilaiTeks, teks: string(l)}, nil
}

type literalAngka struct{ d *apd.Decimal }

func (l literalAngka) eval(*konteks) (nilai, error) { return nilai{jenis: nilaiAngka, angka: l.d}, nil }

type rujukan string

func (r rujukan) eval(k *konteks) (nilai, error) {
	s, err := k.h.ambil(string(r))
	if err != nil {
		return nilai{}, fmt.Errorf("%s: %w", r, err)
	}
	if hb, ok := k.h.(halamanBertipe); ok && hb.bertipeAngka(string(r)) {
		if d, err := utils.ParseDecimal(s); err == nil {
			return nilai{jenis: nilaiAngka, angka: d}, nil
		}
	}
	return nilai{jenis: nilaiTeks, teks: s}, nil
}

type logika struct {
	dan         bool
	kiri, kanan simpul
}

func (l logika) eval(k *konteks) (nilai, error) {
	a, err := l.kiri.eval(k)
	if err != nil {
		return nilai{}, err
	}
	ab, err := a.sebagaiBool()
	if err != nil {
		return nilai{}, err
	}
	if l.dan && !ab {
		return nilai{jenis: nilaiBool}, nil
	}
	if !l.dan && ab {
		return nilai{jenis: nilaiBool, benar: true}, nil
	}
	b, err := l.kanan.eval(k)
	if err != nil {
		return nilai{}, err
	}
	bb, err := b.sebagaiBool()
	if err != nil {
		return nilai{}, err
	}
	return nilai{jenis: nilaiBool, benar: bb}, nil
}

type banding struct {
	op          string
	kiri, kanan simpul
}

func (b banding) eval(k *konteks) (nilai, error) {
	x, err := b.kiri.eval(k)
	if err != nil {
		return nilai{}, err
	}
	y, err := b.kanan.eval(k)
	if err != nil {
		return nilai{}, err
	}
	dx, okx := x.sebagaiAngka()
	dy, oky := y.sebagaiAngka()
	kesetaraan := b.op == "==" || b.op == "=" || b.op == "!="
	var hasil bool
	switch {
	case okx && oky && x.jenis == nilaiAngka && y.jenis == nilaiAngka:
		hasil = b.menurut(dx.Cmp(dy))
	case okx && oky:
		// Salah satu sisi teks yang terbaca angka: tipe properti menentukan
		// hasilnya, dan tipe itu tidak ada di korpus.
		angka, teks := b.menurut(dx.Cmp(dy)), b.menurut(strings.Compare(x.String(), y.String()))
		if angka != teks {
			return nilai{}, fmt.Errorf("%w: operator %s", ErrTafsirBerbeda, b.op)
		}
		hasil = angka
	case kesetaraan && x.jenis != nilaiAngka && y.jenis != nilaiAngka:
		hasil = b.menurut(strings.Compare(x.String(), y.String()))
	default:
		// Urutan atas bukan-angka, atau angka bertipe lawan teks bukan angka
		// (mis. kosong).
		return nilai{}, fmt.Errorf("%w: operator %s", ErrNilaiBukanAngka, b.op)
	}
	return nilai{jenis: nilaiBool, benar: hasil}, nil
}

// menurut - hasil operator atas hasil perbandingan c (<0, 0, >0).
func (b banding) menurut(c int) bool {
	switch b.op {
	case "==", "=":
		return c == 0
	case "!=":
		return c != 0
	case "<":
		return c < 0
	case "<=":
		return c <= 0
	case ">":
		return c > 0
	}
	return c >= 0
}

type panggilan struct {
	fungsi string
	arg    []simpul
}

func (p panggilan) eval(k *konteks) (nilai, error) {
	switch p.fungsi {
	case "@if":
		s, err := p.arg[0].eval(k)
		if err != nil {
			return nilai{}, err
		}
		benar, err := s.sebagaiBool()
		if err != nil {
			return nilai{}, err
		}
		if benar {
			return p.arg[1].eval(k)
		}
		return p.arg[2].eval(k)
	case "@toDecimal":
		v, err := p.arg[0].eval(k)
		if err != nil {
			return nilai{}, err
		}
		d, ok := v.sebagaiAngka()
		if !ok {
			return nilai{}, fmt.Errorf("%w: @toDecimal", ErrNilaiBukanAngka)
		}
		return nilai{jenis: nilaiAngka, angka: d}, nil
	case "@replaceAll":
		var s [3]string
		for i := range s {
			v, err := p.arg[i].eval(k)
			if err != nil {
				return nilai{}, err
			}
			s[i] = v.String()
		}
		// Java String.replaceAll memakai regex; pola korpus hanya "," -
		// literal yang sama artinya dalam regex.
		return nilai{jenis: nilaiTeks, teks: strings.ReplaceAll(s[0], s[1], s[2])}, nil
	case "@divide":
		return p.bagi(k)
	case "@StrategyUtils.pow":
		return p.pangkat(k)
	}
	return nilai{}, fmt.Errorf("%w: fungsi %s", ErrSintaks, p.fungsi)
}

func (p panggilan) bagi(k *konteks) (nilai, error) {
	d := make([]*apd.Decimal, len(p.arg))
	for i := range d {
		v, err := p.arg[i].eval(k)
		if err != nil {
			return nilai{}, err
		}
		x, ok := v.sebagaiAngka()
		if !ok {
			return nilai{}, fmt.Errorf("%w: @divide argumen %d", ErrNilaiBukanAngka, i+1)
		}
		d[i] = x
	}
	if d[1].IsZero() {
		return nilai{}, ErrBagiNol
	}
	if len(p.arg) == 2 {
		// Tanpa skala: hanya hasil eksak yang diterima.
		hasil := new(apd.Decimal)
		kondisi, err := utils.DecimalContext().Quo(hasil, d[0], d[1])
		if err != nil {
			return nilai{}, err
		}
		if kondisi.Inexact() {
			return nilai{}, ErrSkalaBawaanBelumTerverifikasi
		}
		// Nol ekor dari presisi 38 dibuang ("25", bukan "25.000…").
		// [dugaan] Bentuk teks yang disimpan Pega belum terverifikasi; nilai
		// angkanya sama.
		if _, _, err := utils.DecimalContext().Reduce(hasil, hasil); err != nil {
			return nilai{}, err
		}
		return nilai{jenis: nilaiAngka, angka: hasil}, nil
	}
	skala, err := d[2].Int64()
	if err != nil {
		return nilai{}, fmt.Errorf("%w: skala @divide bukan bilangan bulat", ErrSintaks)
	}
	ctx := utils.DecimalContext()
	hasil := new(apd.Decimal)
	kondisi, err := ctx.Quo(hasil, d[0], d[1])
	if err != nil {
		return nilai{}, err
	}
	ctx.Rounding = k.mode
	kondisiQ, err := ctx.Quantize(hasil, hasil, -int32(skala))
	if err != nil {
		return nilai{}, err
	}
	if kondisi.Inexact() || kondisiQ.Inexact() {
		k.dibulatkan = true
	}
	return nilai{jenis: nilaiAngka, angka: hasil}, nil
}

// pangkat - `@StrategyUtils.pow(x, n)`, n bulat 0..100, dihitung eksak.
func (p panggilan) pangkat(k *konteks) (nilai, error) {
	var d [2]*apd.Decimal
	for i := range d {
		v, err := p.arg[i].eval(k)
		if err != nil {
			return nilai{}, err
		}
		x, ok := v.sebagaiAngka()
		if !ok {
			return nilai{}, fmt.Errorf("%w: @StrategyUtils.pow argumen %d", ErrNilaiBukanAngka, i+1)
		}
		d[i] = x
	}
	n, err := d[1].Int64()
	if err != nil || n < 0 || n > 100 {
		return nilai{}, ErrPangkat
	}
	ctx := utils.DecimalContext()
	hasil := apd.New(1, 0)
	for i := int64(0); i < n; i++ {
		kondisi, err := ctx.Mul(hasil, hasil, d[0])
		if err != nil {
			return nilai{}, err
		}
		if kondisi.Inexact() {
			return nilai{}, ErrPangkat
		}
	}
	return nilai{jenis: nilaiAngka, angka: hasil}, nil
}

// ---- pengurai ----

type pengurai struct {
	s   string
	pos int
}

func urai(s string) (simpul, error) {
	p := &pengurai{s: s}
	n, err := p.atau()
	if err != nil {
		return nil, err
	}
	p.spasi()
	if p.pos != len(p.s) {
		return nil, fmt.Errorf("%w: sisa teks di posisi %d", ErrSintaks, p.pos)
	}
	return n, nil
}

func (p *pengurai) spasi() {
	for p.pos < len(p.s) && (p.s[p.pos] == ' ' || p.s[p.pos] == '\t' || p.s[p.pos] == '\n' || p.s[p.pos] == '\r') {
		p.pos++
	}
}

func (p *pengurai) awalan(t string) bool {
	p.spasi()
	return strings.HasPrefix(p.s[p.pos:], t)
}

func (p *pengurai) makan(t string) bool {
	if p.awalan(t) {
		p.pos += len(t)
		return true
	}
	return false
}

func (p *pengurai) atau() (simpul, error) {
	kiri, err := p.dan()
	if err != nil {
		return nil, err
	}
	for p.makan("||") {
		kanan, err := p.dan()
		if err != nil {
			return nil, err
		}
		kiri = logika{dan: false, kiri: kiri, kanan: kanan}
	}
	return kiri, nil
}

func (p *pengurai) dan() (simpul, error) {
	kiri, err := p.bandingan()
	if err != nil {
		return nil, err
	}
	for p.makan("&&") {
		kanan, err := p.bandingan()
		if err != nil {
			return nil, err
		}
		kiri = logika{dan: true, kiri: kiri, kanan: kanan}
	}
	return kiri, nil
}

func (p *pengurai) bandingan() (simpul, error) {
	kiri, err := p.utama()
	if err != nil {
		return nil, err
	}
	// Urutan penting: operator dua karakter lebih dulu.
	for _, op := range []string{"==", "!=", "<=", ">=", "<", ">", "="} {
		if p.makan(op) {
			kanan, err := p.utama()
			if err != nil {
				return nil, err
			}
			return banding{op: op, kiri: kiri, kanan: kanan}, nil
		}
	}
	return kiri, nil
}

func (p *pengurai) utama() (simpul, error) {
	p.spasi()
	if p.pos >= len(p.s) {
		return nil, fmt.Errorf("%w: ekspresi terputus", ErrSintaks)
	}
	c := p.s[p.pos]
	switch {
	case c == '(':
		p.pos++
		n, err := p.atau()
		if err != nil {
			return nil, err
		}
		if !p.makan(")") {
			return nil, fmt.Errorf("%w: kurung tak tertutup", ErrSintaks)
		}
		return n, nil
	case c == '"':
		akhir := strings.IndexByte(p.s[p.pos+1:], '"')
		if akhir < 0 {
			return nil, fmt.Errorf("%w: teks tak tertutup di posisi %d", ErrSintaks, p.pos)
		}
		t := p.s[p.pos+1 : p.pos+1+akhir]
		p.pos += akhir + 2
		return literalTeks(t), nil
	case c == '@':
		return p.fungsi()
	case c >= '0' && c <= '9':
		mulai := p.pos
		for p.pos < len(p.s) && (p.s[p.pos] >= '0' && p.s[p.pos] <= '9' || p.s[p.pos] == '.') {
			p.pos++
		}
		d, err := utils.ParseDecimal(p.s[mulai:p.pos])
		if err != nil {
			return nil, fmt.Errorf("%w: angka di posisi %d", ErrSintaks, mulai)
		}
		return literalAngka{d: d}, nil
	case c == '.' || c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		mulai := p.pos
		for p.pos < len(p.s) && isNamaChar(p.s[p.pos]) {
			p.pos++
		}
		return rujukan(p.s[mulai:p.pos]), nil
	}
	return nil, fmt.Errorf("%w: karakter tak dikenal di posisi %d", ErrSintaks, p.pos)
}

func isNamaChar(c byte) bool {
	return c == '.' || c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9'
}

// cacahArgumen - jumlah argumen yang sah per fungsi yang dikenali.
var cacahArgumen = map[string][]int{
	"@if": {3}, "@toDecimal": {1}, "@replaceAll": {3}, "@divide": {2, 3}, "@StrategyUtils.pow": {2},
}

func (p *pengurai) fungsi() (simpul, error) {
	mulai := p.pos
	p.pos++
	for p.pos < len(p.s) && isNamaChar(p.s[p.pos]) {
		p.pos++
	}
	nama := p.s[mulai:p.pos]
	sah, dikenal := cacahArgumen[nama]
	if !dikenal {
		return nil, fmt.Errorf("%w: fungsi %s", ErrSintaks, nama)
	}
	if !p.makan("(") {
		return nil, fmt.Errorf("%w: %s tanpa kurung", ErrSintaks, nama)
	}
	var arg []simpul
	for {
		a, err := p.atau()
		if err != nil {
			return nil, err
		}
		arg = append(arg, a)
		if p.makan(",") {
			continue
		}
		if p.makan(")") {
			break
		}
		return nil, fmt.Errorf("%w: argumen %s tak tertutup", ErrSintaks, nama)
	}
	for _, n := range sah {
		if len(arg) == n {
			return panggilan{fungsi: nama, arg: arg}, nil
		}
	}
	return nil, fmt.Errorf("%w: %s dengan %d argumen", ErrSintaks, nama, len(arg))
}
