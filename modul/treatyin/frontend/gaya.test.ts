// Penjaga gaya modul Treaty In (02-10-2026): CSS modul TERISOLASI dan TIDAK menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `treatyin.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.treatyin` (`display: contents`, jadi tata
//     letak kerangka tidak berubah);
//   - SETIAP pemilih di `treatyin.css` diawali kelas akar (atau `:where(.treatyin ...)` yang menjaga kekhususan
//     nol) - nol pemilih global;
//   - kelas khusus modul ini tidak didefinisikan di `inti/frontend/styles.css` (work owner 02-10-2026:
//     "untuk css style per modul tidak ada lagi menumpang per inti");
//   - tanpa properti yang mengurung popup `position: fixed`.
//
// ⚠️ Modul ini TIDAK pernah menumpang di inti - ia memang belum punya gaya sama sekali sampai hari ini.
// Penjaganya tetap sama bentuknya dengan tujuh modul yang dipindah, supaya yang menambah aturan besok
// tidak perlu tahu modul mana yang lahir dari pemindahan dan mana yang lahir baru.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini: wajib hidup di berkas modul, dan nol jejak di inti. */
const KELAS_MODUL = [
  'trin__aturan',
  'trin__lencana',
  // Ronde layar 1 — kelas kedua layar baru.
  'trin__redup',
  'trin__catatan',
  'trin__galat',
  'trin__kepala',
  'trin__id',
  'trin__radio',
  'trin__pilih-luar',
  'trin__centang',
  'trin__aksi',
  'trin__aksi--kaki',
  // Ronde layar 2 — tata letak.
  'trin__dwikolom',
  'trin__kolom',
  'trin__panel-kepala',
  'trin__kaki',
  'trin__tabel',
  'trin__tabel-kurs',
  'trin__kol-mata-uang',
  'trin__kol-nilai',
  'trin__kol-tanggal',
  'trin__kepala-kolom',
  'trin__kepala-kanan',
  'trin__ikon-saring',
  'trin__baris-saring',
  'trin__belum',
  'trin__belum-judul',
  'trin__belum-petunjuk',
  'trin__teks',
  'trin__teks-asal',
  'trin__teks-lain',
  'trin__spanduk',
] as const

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

/** Pisah daftar pemilih pada koma tingkat atas (koma di dalam `:is(...)`/`:where(...)` tidak memisah). */
function pisahKoma(p: string): string[] {
  const hasil: string[] = []
  let dalam = 0
  let awal = 0
  for (let i = 0; i < p.length; i++) {
    if (p[i] === '(') dalam++
    else if (p[i] === ')') dalam--
    else if (p[i] === ',' && dalam === 0) {
      hasil.push(p.slice(awal, i))
      awal = i + 1
    }
  }
  hasil.push(p.slice(awal))
  return hasil.map((s) => s.trim())
}

/** Pemilih CSS tingkat atas (tanpa komentar; isi @media ikut, kepala @media tidak). */
function pemilihCSS(css: string): string[] {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(/([^{};]+)\{/g)]
    .map((m) => (m[1] ?? '').trim())
    .filter((p) => p !== '' && !p.startsWith('@'))
    .flatMap(pisahKoma)
}

function terisolasi(p: string): boolean {
  // ⭐ 8 Oktober 2026 — tema Treaty Exchange Yearly: varian GELAP ditulis
  // `:root[data-theme="dark"] .treatyin …`, pengecualian SEMPIT yang sama
  // dengan penjaga `treatyexchangeyearly/frontend/gaya.test.ts`. Ia tetap
  // berakar `.treatyin`: hanya atribut tema `:root` yang mendahuluinya.
  const p2 = p.startsWith(':root[data-theme="dark"] ') ? p.slice(':root[data-theme="dark"] '.length) : p
  return p2 === '.treatyin' || p2.startsWith('.treatyin ') || p2.startsWith(':where(.treatyin ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul Treaty In', () => {
  it('satu-satunya berkas CSS modul adalah treatyin.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['treatyin.css'])
    expect(RUTE).toContain("import './treatyin.css'")
    expect(RUTE).toContain('<div className="treatyin">')
    expect(CSS).toMatch(/^\.treatyin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di treatyin.css diawali kelas akar .treatyin - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(1)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS('table { x: 1 } .treatyin .a, .b { y: 2 } .treatyin :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }')
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
    // Tema gelap berakar modul diterima; tema gelap GLOBAL tetap ditolak.
    const gelap = pemilihCSS(':root[data-theme="dark"] .treatyin > .inbox { a: 1 } :root[data-theme="dark"] .x { b: 2 }')
    expect(gelap.filter((p) => !terisolasi(p))).toEqual([':root[data-theme="dark"] .x'])
  })

  it('kelas khusus modul ini tidak menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .trin__aturan:hover { a: 1 } .trin__aturan-lain { b: 2 }', ['trin__aturan'])).toEqual(['trin__aturan'])
    expect(masihDi('.trin__aturan-lain { b: 2 }', ['trin__aturan'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
  })
})

// ===========================================================================
// Kerapatan badan layar — keputusan pemilik proses 6 Oktober 2026
// ===========================================================================
//
// "agar tidak terlalu besar … ambil rentang sedang nya saja … tema nya bisa
//  mencontoh menu kelola user"
describe('kerapatan badan layar Treaty In', () => {
  // ⛔ UJI PALING PENTING DI BLOK INI.
  //
  // Jaraknya dahulu dihitung DUA KALI — `gap: 14px` kolomnya ditambah
  // `margin-bottom: 16px` yang tema inti berikan ke tiap `.field`, menjadi
  // 30px. Yang terlihat "terlalu besar" justru penjumlahan itu, bukan
  // ukuran kontrolnya.
  it('⛔ jarak antar medan punya SATU pemilik, bukan dua', () => {
    expect(CSS).toMatch(/\.treatyin \.trin__kolom \.field \{\s*margin-bottom: 0;/)
  })

  // ⛔ BENTUK PEGA: LABEL DI KIRI MEDAN.
  //
  // Tangkapan layar Pega 6 Oktober 2026 menjawab pertanyaan yang empat ronde
  // sebelumnya salah tebak: yang membuat layar terasa salah BUKAN lebarnya,
  // melainkan letak labelnya. Label di kiri membuat satu baris setinggi satu
  // kotak, bukan dua — dan dua kolom selebar layar tetap terbaca rapat.
  it('⭐ label duduk di KIRI medan, bukan di atasnya', () => {
    expect(CSS).toMatch(/grid-template-columns: 132px minmax\(0, 1fr\);/)
    expect(CSS).toMatch(/\.treatyin \.trin__kolom \.field \.field__label \{\s*margin: 0;/)
  })

  // ⛔ Keempat bentuk yang ditolak salah pada hal yang SAMA: mengatur lebar,
  // padahal yang keliru tata letak label. Nol di antaranya boleh kembali.
  it('⛔ wadah form TIDAK lagi dibatasi lebar maupun ditengahkan', () => {
    expect(CSS).not.toMatch(/\.treatyin \.trin__dwikolom[^}]*max-width: 1[01]\d0px/)
    expect(CSS).not.toMatch(/\.treatyin \.form-grid[^}]*margin-inline: auto/)
  })

  it('aturannya SATU untuk kepala dan untuk grid tab', () => {
    expect(CSS).toMatch(/\.treatyin \.trin__dwikolom,\s*\.treatyin \.form-grid \{/)
  })

  // ⚠️ Di bawah 900px label kembali ke ATAS medan: kolom label 132px pada
  // lebar telepon menyisakan kotak isian yang terlalu ramping untuk diisi.
  it('⚠️ bentuk label-kiri hanya dari 900px ke atas', () => {
    const i = CSS.indexOf('grid-template-columns: 132px')
    expect(i).toBeGreaterThan(0)
    expect(CSS.lastIndexOf('@media (min-width: 900px)', i)).toBeGreaterThan(0)
  })

  // ⚠️ TABEL DATA tidak ikut dibatasi. Pemilik proyek mencabut `max-width`
  // dari `.page` justru karena tabel yang berhenti di tengah layar membuang
  // ruang yang paling dibutuhkannya.
  //
  // ⭐ SATU pengecualian bernama: ketiga grid total tab Share proporsional
  // (`.trin__share-total`). Ia bukan tabel data melainkan pasangan
  // judul↔`Value` dua kolom; dibentangkan selebar layar, kolom `Value`
  // berdiri ~1500px dari judulnya dan berhenti terbaca sebagai satu baris.
  //
  // ⛔ Daftar pengecualiannya TERTUTUP — tabel kedua yang dibatasi lebarnya
  // membuat uji ini gagal, dan itu memang maksudnya.
  it('⚠️ tabel data TIDAK dibatasi lebarnya, kecuali grid total Share', () => {
    const dibatasi = [...CSS.matchAll(/([^}]*?)\{[^}]*max-width:\s*\d+px[^}]*\}/g)]
      .map((m) => m[1] ?? '')
      .filter((sel) => /\.trin__tabel|\.table-wrap/.test(sel))
      .map((sel) => sel.trim())
    for (const sel of dibatasi) {
      expect(sel).toContain('trin__share-total')
    }
  })

  // ⛔ `min-height` inti dicabut untuk panel ini: `.table-wrap` memesan 160px
  // bahkan ketika isinya `No items`, sehingga panel KOSONG pun memakan ruang
  // satu tabel penuh.
  it('Existing Policy menampilkan empat baris lalu menggulir', () => {
    expect(CSS).toMatch(
      /\.treatyin \.trin__polis \.table-wrap \{\s*min-height: 0;\s*max-height: 170px;/,
    )
  })

  // ⭐ RALAT 7 Oktober 2026 — SATU pengecualian bernama, atas permintaan
  // pemilik proses: *"untuk design treaty in dikecilin lagi krn terlalu
  // besar"*.
  //
  // ⛔ Larangannya TIDAK dicabut, ia dipersempit. Yang boleh hanya SATU
  // aturan: `.treatyin .field__input` setinggi 34px, di blok kerapatan yang
  // menyebut sebabnya. Ukuran tombol tetap NOL boleh ditimpa, dan tinggi
  // lain pada `.field__input` tetap ditolak — pengecualian tanpa batas
  // sama saja dengan mencabut penjaganya.
  //
  // ⚠️ Bawaan `inti` (42px/15px) tidak disentuh: ia ukuran seluruh
  // aplikasi, dan permintaan ini hanya tentang Treaty In.
  it('⛔ ukuran kontrol hanya boleh ditimpa oleh blok kerapatan bernama', () => {
    // ⛔ Tinggi tombol: `.btn` telanjang NOL pengecualian.
    expect(CSS).not.toMatch(/\.treatyin[^{]*\.btn\s*\{[^}]*height:/)

    // ⚠️ PAGAR DIRAPATKAN 8 Oktober 2026. Bentuk sebelumnya hanya mencari
    // `.btn {` — varian seperti `.btn--sm {` LOLOS tanpa suara, dan sebuah
    // penimpaan tinggi tombol sempat masuk tanpa terdaftar di sini.
    //
    // ⭐ SATU pengecualian, dan ia harus menyebut dirinya: tombol DI DALAM
    // grid mengikuti tinggi kotak isian (34px), bukan 36px bawaan `.btn--sm`
    // — selisih dua piksel membuat baris grid tampak tidak rata.
    const tinggiTombol = [...CSS.matchAll(/\.treatyin[^{]*\.btn--[a-z]+\s*\{([^}]*)\}/g)]
      .map((m) => /height:\s*([^;]+);/.exec(m[1] ?? '')?.[1]?.trim())
      .filter((v): v is string => v !== undefined)
    expect(tinggiTombol).toEqual(['28px', '34px'])
    expect(CSS).toContain('Tombol baris mengikuti tinggi kotak isian')

    // `.field__input`: tepat SATU aturan bertinggi, dan tingginya 34px.
    const tinggi = [...CSS.matchAll(/\.treatyin[^{]*\.field__input\s*\{([^}]*)\}/g)]
      .map((m) => /height:\s*([^;]+);/.exec(m[1] ?? '')?.[1]?.trim())
      .filter((v): v is string => v !== undefined)
    // ⚠️ DUA nilai, dan keduanya perlu:
    //   `28px` kotak satu baris — permintaan pengecilan 8 Oktober 2026
    //          (*"design inputan nya diubah seperti pega jangan kegedean"*);
    //          sebelumnya 34px, dan sebelum itu 42px bawaan `inti`;
    //   `auto`  AREA TEKS, yang justru harus LEPAS dari tinggi itu. Area
    //           teks setinggi 34px memperlihatkan satu baris kalimat dan
    //           memaksa pembacanya menggulir di dalam kotak.
    //   `34px` panel LAMPIRAN — blok kerapatannya sendiri sejak 8 Oktober
    //          2026 (*"dynamic layout antara table attachment dengan inputan
    //          yg di tab dipisah jangan di samakan"*). Isian tab adalah
    //          FORMULIR DATA (kerapatan menang); panel lampiran PENGELOLA
    //          BERKAS (ketepatan tekan menang). Keduanya tidak lagi bergerak
    //          bersama.
    expect(tinggi).toEqual(['28px', 'auto', '34px'])

    // ⛔ Dan blok itu WAJIB menyebut sebabnya — aturan tanpa alasan adalah
    // aturan yang akan disalin ke modul berikutnya tanpa dipikirkan.
    expect(CSS).toContain('KERAPATAN — layar Treaty In dikecilkan')
    expect(CSS).toContain('BAWAAN `inti` TIDAK DISENTUH')

    // ⛔ Blok lampiran pun WAJIB menyebut sebabnya — dan sebabnya BUKAN
    // selera, melainkan dua urusan yang berbeda.
    expect(CSS).toContain('KERAPATAN PANEL LAMPIRAN — TERPISAH dari kerapatan isian tab')
    expect(CSS).toContain('PENGELOLA BERKAS')
  })

  // ⭐ 7 Oktober 2026 — deret kaki Save · Close · Actions kini berjarak bawah
  // (History menempel ke tombol), tetapi lewat pengubah `--kaki`. `.trin__aksi`
  // sendiri tetap tidak disentuh: ia dipakai juga deret tombol DI DALAM kartu.
  it('`.trin__aksi` tidak disentuh — jarak kaki lewat pengubah `--kaki`', () => {
    expect(CSS).toMatch(/\.treatyin \.trin__aksi \{\s*display: flex;\s*flex-wrap: wrap;\s*gap: 8px;\s*margin-top: 12px;\s*\}/)
    expect(CSS).toMatch(/\.treatyin \.trin__aksi--kaki \{\s*margin-bottom: 18px;\s*\}/)
  })
})

// ⭐ 8 Oktober 2026 — "ubah tema … menjadi seperti tema Treaty Exchange Yearly".
describe('tema Treaty Exchange Yearly', () => {
  const css = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const akar = css.slice(css.indexOf('.treatyin > .inbox {'))

  it('token terang di AKAR MODUL (juga layar tanpa `.inbox`), kartu akar halaman: `--trin-*` bernilai `--tey-*`, gradasi lembut, sudut 26px', () => {
    expect(css).toMatch(/^\.treatyin \{[^}]*--trin-latar: #eef1f6;/m)
    expect(akar.slice(0, 3000)).toContain('border-radius: 26px;')
    expect(akar.slice(0, 3000)).toContain('background: linear-gradient(180deg, var(--trin-latar-atas) 0%, var(--trin-latar-bawah) 100%);')
  })

  it('varian gelap ditulis berakar modul', () => {
    expect(css).toMatch(/^:root\[data-theme="dark"\] \.treatyin \{[^}]*--trin-isi: #1f2638;/m)
  })

  it('tombol kapsul: utama merah bergradasi, kartu panel putih timbul, isian cekung', () => {
    expect(css).toMatch(/\.treatyin \.btn \{\s*border-radius: 999px;/)
    expect(css).toMatch(/\.treatyin \.btn--primary,\s*\.treatyin \.tl-tambah \{\s*background: linear-gradient\(180deg, var\(--trin-tombol\), var\(--trin-tombol-ujung\)\);/)
    expect(css).toMatch(/\.treatyin \.panel \{\s*background: var\(--trin-isi\);[^}]*border-radius: 18px;[^}]*box-shadow: var\(--trin-timbul\);/)
    expect(css).toMatch(/\.treatyin \.field__input \{[^}]*box-shadow: var\(--trin-cekung-kecil\);/)
  })

  it('strip tab = kontrol segmen kapsul; kepala tabel navy muda dan TIDAK tembus pandang', () => {
    expect(css).toMatch(/\.treatyin \.tabs__item--aktif,\s*\.treatyin \.tabs__item--aktif:hover \{\s*background: var\(--trin-isi\);/)
    expect(css).toMatch(/\.treatyin \.table-wrap \.trin__tabel thead th \{\s*background: var\(--trin-kepala-tabel\);/)
  })

  it('lapisan tema TIDAK menimpa ukuran kontrol — kerapatan tetap milik bloknya', () => {
    const tema = css.slice(css.indexOf('.treatyin > .inbox {'))
    expect(tema).not.toMatch(/\.(btn|field__input)[^{]*\{[^}]*height:/)
  })
})

// ⛔ KOLOM LABEL SATU UKURAN DI SELURUH LAYAR — permintaan pemilik proses
// 8 Oktober 2026 atas tab Account Reporting Period: *"sebisa mungkin
// dirapikan letaknya … designnya harus seperti yg di pega"*.
//
// ⚠️ Yang terbaca "berantakan" BUKAN jarak antarbarisnya, melainkan kolom
// labelnya: bentuk lama `minmax(110px, 38%)` berukuran PERSEN, jadi tiap
// blok punya lebar label sendiri dan nol baris segaris dengan yang lain.
describe('kolom label blok tata Pega', () => {
  it('⛔ lebarnya TETAP, bukan persen', () => {
    const i = CSS.indexOf('.treatyin .trin__tata--kiri > .field {')
    expect(i).toBeGreaterThan(0)
    const aturan = CSS.slice(i, CSS.indexOf('}', i))
    expect(aturan).toContain('grid-template-columns: 132px minmax(0, 1fr)')
    expect(aturan).not.toContain('38%')
  })

  // ⭐ Dan angkanya SAMA dengan kolom label form — satu ukuran untuk seluruh
  // layar, bukan dua yang kebetulan mirip.
  it('⭐ angkanya sama dengan kolom label form', () => {
    expect([...CSS.matchAll(/grid-template-columns: 132px minmax\(0, 1fr\)/g)].length).toBeGreaterThanOrEqual(2)
  })

  it('⚠️ batas lebar ada pada KOTAKNYA, bukan pada selnya', () => {
    const i = CSS.indexOf('.treatyin .trin__tata--kiri > .field > .field__input,')
    expect(i).toBeGreaterThan(0)
    expect(CSS.slice(i, CSS.indexOf('}', i))).toContain('max-width: 320px')
    // Bentuk yang dilarang: membatasi SEL-nya, yang menggeser kolom berikutnya.
    expect(CSS).not.toMatch(/\.treatyin \.trin__tata--kiri > \.field \{[^}]*max-width/)
  })
})

// ⛔ RUPA GRID SERAGAM — permintaan pemilik proses 8 Oktober 2026:
// *"tampilan di view saya mau di semuanya seperti ini, tidak hanya view"*.
describe('rupa grid seragam', () => {
  const TRIA = readFileSync(
    join(AKAR, '..', '..', 'treatyinadjustment', 'frontend', 'treatyinadjustment.css'),
    'utf8',
  )

  /** Bantalan sel yang BENAR-BENAR berlaku — aturan TERAKHIR yang menyebutnya. */
  const bantalanTerakhir = (css: string, pemilih: string): string => {
    const i = css.lastIndexOf(pemilih)
    expect(i).toBeGreaterThan(0)
    const m = /padding:\s*([^;]+);/.exec(css.slice(i, css.indexOf('}', i)))
    return m?.[1]?.trim() ?? ''
  }

  // ⚠️ Sengaja memeriksa aturan TERAKHIR, bukan yang pertama. Bantalan grid
  // kurs pernah ditulis dua kali di berkas yang sama, dan yang di bawah
  // menimpa yang di atas — memeriksa yang pertama akan hijau sementara layar
  // memakai nilai yang lain.
  it('⛔ kedua modul memakai bantalan sel yang SAMA', () => {
    expect(bantalanTerakhir(CSS, '.treatyin .trin__tabel td,')).toBe('7px 12px')
    expect(bantalanTerakhir(TRIA, '.treatyinadjustment .tria__tabel th,')).toBe('7px 12px')
  })

  // ⛔ BELANG `inti` TETAP DINETRALKAN — ia berkekhususan NOL
  // (`:where(tbody tr…)`) dan memakai `--alt-row` (= `--surface-2`), cukup
  // pekat untuk terbaca sebagai dua warna berselang-seling. Itu yang ditolak
  // pemilik proses, dan penetralnya harus tetap ada.
  it('⛔ belang `--alt-row` inti dinetralkan di kedua modul', () => {
    for (const [css, pemilih] of [
      [CSS, '.treatyin .trin__tabel tbody tr:nth-child(even)'],
      [TRIA, '.treatyinadjustment .tria__tabel tbody tr:nth-child(even)'],
    ] as const) {
      const i = css.indexOf(pemilih)
      expect(i).toBeGreaterThan(0)
      expect(css.slice(i, css.indexOf('}', i))).toContain('background: transparent')
    }
  })

  // ⭐ DAN BELANG MODUL YANG HAMPIR TIDAK TERLIHAT DIPASANG — permintaan
  // pemilik proses 8 Oktober 2026: *"buat design zebra nya tapi jangan gelap
  // terang, tapi setidaknya bisa tahu perbedaan antar baris"*.
  //
  // ⚠️ Yang ditolak dahulu KEKUATANNYA, bukan belangnya. Uji ini karena itu
  // menjaga DUA hal sekaligus: belangnya ada, dan ia tetap tipis.
  it('⭐ belang modul ada, dan nilainya tipis — bukan gelap terang', () => {
    for (const [css, akar, pre] of [
      [CSS, 'treatyin', 'trin'],
      [TRIA, 'treatyinadjustment', 'tria'],
    ] as const) {
      const i = css.indexOf(`--${pre}-belang: `)
      expect(i, akar).toBeGreaterThan(0)
      const nilai = css.slice(i, css.indexOf(';', i))
      // Alfa maksimum 6% — di atas itu ia terbaca sebagai dua warna.
      const alfa = Number(/,\s*([0-9.]+)\)/.exec(nilai)?.[1] ?? '1')
      expect(alfa, `${akar} ${nilai}`).toBeGreaterThan(0)
      expect(alfa, `${akar} ${nilai}`).toBeLessThanOrEqual(0.06)
      // ⛔ BUKAN `--alt-row`: itu token yang justru dicabut.
      expect(nilai).not.toContain('--alt-row')
      // Dipakai, bukan sekadar dideklarasikan.
      expect(css).toContain(`background: var(--${pre}-belang)`)
    }
  })

  // ⛔ Baris RINCIAN adalah `<tr>` juga; mewarnai `<tr>` membuat seluruh
  // panel yang terbuka ikut berbelang. Dan `:hover` harus MENANG atas
  // belang, kalau tidak baris genap yang disorot terbaca lebih pucat.
  it('⛔ belang melewati baris rincian dan tidak menimpa sorot', () => {
    for (const [css, pre] of [
      [CSS, 'trin'],
      [TRIA, 'tria'],
    ] as const) {
      const i = css.indexOf(`background: var(--${pre}-belang)`)
      // ⚠️ DIPERIKSA PER BARIS PEMILIH, bukan gabungannya. Aturannya punya
      // dua pemilih berkoma, dan memeriksa gabungan membuat satu pemilih
      // yang benar menutupi satu pemilih yang bocor.
      const pemilih = css
        .slice(css.lastIndexOf('}', i) + 1, i)
        .split(',')
        .map((x) => x.trim())
        .filter((x) => x !== '' && !x.startsWith('{'))
      expect(pemilih.length).toBeGreaterThan(0)
      for (const satu of pemilih) {
        expect(satu, satu).toContain(`:not(.${pre}__rincian)`)
        expect(satu, satu).toContain(':not(:hover)')
        // Diwarnai pada selnya, bukan pada barisnya.
        expect(satu, satu).toContain('> td')
      }
    }
  })

  // ⭐ Sorot baris TETAP: bersama belang, sorotlah penanda baris mana yang
  // sedang ditunjuk mata.
  it('⭐ sorot baris dipertahankan', () => {
    expect(CSS).toContain('.treatyin .trin__tabel tbody tr:hover')
    expect(TRIA).toContain('.treatyinadjustment .tria__tabel tbody tr:hover')
  })
})

// ⛔ PANEL LAMPIRAN TIDAK IKUT BERGERAK — permintaan pemilik proses
// 8 Oktober 2026: *"dynamic layout antara table attachment dengan inputan
// yg di tab dipisah jangan di samakan"*.
//
// ⚠️ Yang dijaga BUKAN angkanya, melainkan PEMISAHANNYA: setiap kali kotak
// isian tab dikecilkan (42 → 34 → 28px), kontrol panel lampiran dulu ikut
// mengecil tanpa satu pun permintaan menyebutnya.
describe('kerapatan panel lampiran terpisah', () => {
  const LAMP = readFileSync(join(AKAR, 'components', 'PanelLampiran.tsx'), 'utf8')

  it('⛔ panel punya pengaitnya sendiri — tanpa itu pemisahan mustahil', () => {
    expect(LAMP).toContain('<div className="trin__lampiran">')
    expect(CSS).toContain('.treatyin .trin__lampiran .field__input')
  })

  it('⛔ tingginya BERBEDA dari kotak isian tab', () => {
    const tinggi = (pemilih: string): string => {
      const i = CSS.indexOf(pemilih)
      expect(i).toBeGreaterThan(0)
      return /height:\s*([^;]+);/.exec(CSS.slice(i, CSS.indexOf('}', i)))?.[1]?.trim() ?? ''
    }
    expect(tinggi('.treatyin .field__input {')).toBe('28px')
    expect(tinggi('.treatyin .trin__lampiran .field__input {')).toBe('34px')
  })

  it('⭐ grid lampiran juga punya bantalannya sendiri', () => {
    const i = CSS.indexOf('.treatyin .trin__lampiran .trin__tabel td,')
    expect(i).toBeGreaterThan(0)
    expect(CSS.slice(i, CSS.indexOf('}', i))).toContain('padding: 9px 12px')
  })
})
