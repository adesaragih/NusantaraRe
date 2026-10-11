// Disalin dari `modul/claimnonprop/frontend/gaya.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Gaya modul Claim Fac In terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar `.claimfacin`, kelas modul
// berawalan `claimfacin__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'claimfacin.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
// blok @keyframes (pemilih from / to bukan pemilih elemen) - namanya global, diuji berawalan claimfacin- tersendiri
const KEYFRAMES = /@keyframes\s+([\w-]+)\s*\{(?:[^{}]*\{[^{}]*\})*[^{}]*\}/g
const ATURAN_ELEMEN = ATURAN.replace(KEYFRAMES, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Claim Fac In', () => {
  it('satu-satunya berkas CSS adalah claimfacin.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['claimfacin.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './claimfacin.css'")
    expect(rute).toContain('<div className="claimfacin">')
    expect(CSS).toMatch(/^\.claimfacin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .claimfacin - nol pemilih global', () => {
    const pemilih = [...ATURAN_ELEMEN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.claimfacin(\s|$)/)
    const nama = [...ATURAN.matchAll(KEYFRAMES)].map((m) => m[1]!)
    for (const n of nama) expect(n, n).toMatch(/^claimfacin-/)
    expect(ATURAN.match(/@keyframes/g)?.length ?? 0).toBe(nama.length)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan claimfacin__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(claimfacin__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(claimfacin__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('kulit Kelola User: akar .claimfacin__akar bertoken --cfi-* terang dan gelap; warna hanya di blok token', () => {
    expect(CSS).toMatch(/^\.claimfacin \.claimfacin__akar \{[^}]*--cfi-latar: #eef1f6;/m)
    expect(CSS).toMatch(/^:root\[data-theme="dark"\] \.claimfacin \.claimfacin__akar \{[^}]*--cfi-latar: #1b2130;/m)
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.claimfacin__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    for (const [berkas, kelas] of [
      [['pages', 'ClaimFacIn.tsx'], 'className="inbox claimfacin__akar"'],
      [['components', 'LayarKasus.tsx'], 'className="inbox claimfacin__akar"'],
      [['components', 'TataView.tsx'], 'className="claimfacin__tabel"'],
      [['components', 'Popup.tsx'], 'className="claimfacin__tabel"'],
    ] as const) {
      expect(readFileSync(join(AKAR, ...berkas), 'utf8'), kelas).toContain(kelas)
    }
  })

  it('halaman awal memakai kelas inti Kelola User: toolbar, field__input, inbox__tabel, btn--primary', () => {
    const awal = readFileSync(join(AKAR, 'pages', 'ClaimFacIn.tsx'), 'utf8')
    for (const k of [
      'className="toolbar"',
      'field__input',
      'className="inbox__tabel claimfacin__tabel"',
      'btn btn--primary',
    ]) {
      expect(awal, k).toContain(k)
    }
  })

  it('layar kasus = layout lama Pega dirapikan: label kiri, dua kolom, tab, kartu panel Kelola User', () => {
    const tata = readFileSync(join(AKAR, 'components', 'TataView.tsx'), 'utf8')
    for (const k of [
      'className="panel"',
      'panel__title',
      'claimfacin__baris',
      'claimfacin__label-medan',
      'claimfacin__dua',
      '<TabGeser',
      'field__error',
    ]) {
      expect(tata, k).toContain(k)
    }
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__baris \{[^}]*grid-template-columns: 150px minmax\(0, 1fr\)/)
    // tata letak tanpa flex kolom (basis flex medan pernah menjadi tinggi - 12.204px, 07-10-2026)
    expect(ATURAN).not.toMatch(/flex-direction:\s*column/)
  })

  it('huruf modul diperkecil (work owner 08-10-2026 "font juga kecilin"): akar 13px, isian dan tombol 32px', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \{\s*font-size: 13px;\s*\}/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \.field__input \{[^}]*height: 32px;[^}]*font-size: 13px;/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \.btn \{[^}]*height: 32px;[^}]*font-size: 13px;/)
  })

  it('Input Register berlabel di atas (duaAtas - screenshot Pega work owner 11-10-2026 "UBAH LAYOUTNYA BEGINI")', () => {
    expect(readFileSync(join(AKAR, 'components', 'TataView.tsx'), 'utf8')).toContain("u.t.letak === 'duaAtas'")
    expect(CSS).toMatch(
      /\.claimfacin \.claimfacin__dua\.claimfacin__dua-atas \.claimfacin__baris \{\s*grid-template-columns: minmax\(0, 1fr\);/,
    )
    expect(CSS).toMatch(/\.claimfacin__dua-atas \.claimfacin__isi-medan > \.field__error \{\s*flex: 1 0 100%;/)
  })

  it('panel Adjustment = layout lama Pega: Inline grid triple berlabel atas, sel kolom / kiri (work owner 11-10-2026)', () => {
    const tata = readFileSync(join(AKAR, 'components', 'TataView.tsx'), 'utf8')
    expect(tata).toContain(`if (u.t.letak === 'tiga') return <Grup key={i} t={u.t} k={k} kelas="claimfacin__tiga" />`)
    expect(tata).toContain(`if (u.t.letak === 'kolom') return <Grup key={i} t={u.t} k={k} kelas="claimfacin__kolom" />`)
    expect(tata).toContain('className="claimfacin__kiri"')
    // kontainer tersembunyi seluruhnya tidak memakan sel
    expect(tata).toContain("(t.anak ?? []).filter((c) => c.jenis !== 'bagian' || (c.anak ?? []).length > 0)")
    // paling banyak tiga sel per baris, turun sendiri di panel sempit (tanpa container query: containment = Modal terkurung)
    expect(CSS).toMatch(
      /\.claimfacin \.claimfacin__tiga \{[^}]*grid-template-columns: repeat\(auto-fill, minmax\(max\(280px, calc\(\(100% - 2 \* 28px\) \/ 3\)\), 1fr\)\);/,
    )
    expect(ATURAN).not.toMatch(/container-type|container:/)
    expect(CSS).toMatch(
      /\.claimfacin \.claimfacin__tiga \.claimfacin__baris,\s*\.claimfacin \.claimfacin__kolom \.claimfacin__baris \{\s*grid-template-columns: minmax\(0, 1fr\);/,
    )
    expect(CSS).toMatch(
      /\.claimfacin \.claimfacin__tiga \.claimfacin__kiri \.claimfacin__baris \{\s*grid-template-columns: 104px minmax\(0, 1fr\);/,
    )
    expect(CSS).toMatch(
      /\.claimfacin__baris:has\(> \.claimfacin__isi-medan > \.claimfacin__centang\)\s*> \.claimfacin__label-medan \{\s*order: 1;/,
    )
    // nilai hanya-baca kosong ditandai, bukan label menggantung (Detail Chronology to Comitee, 11-10-2026 "RAPIHKAN!")
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__tiga \.claimfacin__nilai:empty::after \{\s*content: "—";/)
  })

  it('panel baris dan kaki grid selebar area terlihat; kaki = bilah ringkasan (work owner 11-10-2026 "tolong rapihin")', () => {
    const tata = readFileSync(join(AKAR, 'components', 'TataView.tsx'), 'utf8')
    expect(tata).toContain("'--cfi-lebar-grid'")
    expect(tata).toContain('{CFI.gridKosong}')
    expect(CSS).toMatch(/tr\.claimfacin__baris-rinci > td > \.claimfacin__rinci \{[^}]*position: sticky;[^}]*left: 0;/)
    expect(CSS).toMatch(/\.claimfacin__grid > \.claimfacin__grid-kaki \{[^}]*position: sticky;/)
    expect(CSS).toMatch(/\.claimfacin__grid-kaki \.claimfacin__nilai:empty::after \{\s*content: "—";/)
  })

  it('nol properti yang memerangkap Modal tanpa portal', () => {
    expect(ATURAN).not.toMatch(
      /(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m,
    )
  })

  it('tanpa kotak masuk Beranda (bukan bagian Pega)', () => {
    expect(readFileSync(join(AKAR, 'menu.ts'), 'utf8')).not.toMatch(/antreanBeranda\s*:|daftarBeranda\s*:/)
  })
})

describe('gaya panel bersarang Claim Fac In', () => {
  it('panel baris di dalam panel baris tetap bergaris aksen dan berlatar tenang', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__rinci \.claimfacin__rinci \{/)
  })

  it('SEMUA tab modul = TabGeser: kapsul penanda meluncur, isi tab dan layar memudar (work owner 11-10-2026 "SEMUA TAB")', () => {
    const geser = readFileSync(join(AKAR, 'components', 'TabGeser.tsx'), 'utf8')
    expect(geser).toContain('className="tabs claimfacin__tabs"')
    expect(geser).toContain('className="claimfacin__tabs-penanda"')
    // nol strip tab lain: StripTab inti (garis bawah) dan tablist buatan sendiri
    for (const p of berkas(AKAR, '.tsx')) {
      const isi = readFileSync(p, 'utf8')
      expect(isi, p).not.toMatch(/\bStripTab\b/)
      if (!p.endsWith('TabGeser.tsx')) expect(isi, p).not.toContain('role="tablist"')
    }
    for (const b of [
      ['components', 'TataView.tsx'],
      ['components', 'PolisNB.tsx'],
      ['components', 'PolisNBTab.tsx'],
      ['pages', 'ClaimFacIn.tsx'],
    ]) {
      expect(readFileSync(join(AKAR, ...b), 'utf8'), b.join('/')).toContain('<TabGeser')
    }
    // isi tab berkunci tab aktif (dipasang ulang -> animasi muncul berjalan lagi); layar kasus berkunci tahap
    for (const [b, kunci] of [
      [['components', 'TataView.tsx'], '<div key={i} role="tabpanel" className="claimfacin__tab-isi">'],
      [['components', 'PolisNB.tsx'], '<div key={tab} role="tabpanel" className="claimfacin__tab-isi">'],
      [['components', 'PolisNBTab.tsx'], '<div key={aktif} role="tabpanel" className="claimfacin__tab-isi">'],
      [['components', 'LayarKasus.tsx'], 'className="claimfacin__halaman"'],
    ] as const) {
      expect(readFileSync(join(AKAR, ...b), 'utf8'), b.join('/')).toContain(kunci)
    }
    expect(CSS).toMatch(
      /\.claimfacin \.claimfacin__tabs-penanda \{[^}]*position: absolute;[^}]*transition:[^}]*left [^}]*width /,
    )
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__tab-isi \{\s*animation: claimfacin-muncul [^;]*backwards;/)
    expect(CSS).toMatch(/@keyframes claimfacin-muncul \{\s*from \{\s*opacity: 0;/)
    expect(CSS).toMatch(
      /@media \(prefers-reduced-motion: reduce\) \{[^@]*\.claimfacin__tabs-penanda[^@]*transition: none;[^@]*animation: none;/,
    )
    // tebal huruf tab aktif = tidak aktif (lebar tombol tidak melompat saat kapsul meluncur)
    expect(CSS).not.toMatch(/\.tabs__item--aktif[^{]*\{[^}]*font-weight/)
  })

  it('grid lebar menggulir di kartunya; daftar autocomplete sel yang terbuka diberi ruang, tidak terpotong', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__grid \{[^}]*overflow-x: auto;/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__grid:has\(.*\.pilih-saring__daftar\) \{\s*padding-bottom: 270px;/)
  })
})
