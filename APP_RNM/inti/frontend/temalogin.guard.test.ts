// Penjaga tema halaman login di styles.css inti (permintaan work owner 02-10-2026):
//
//   - latar sidebar TETAP latar lama (`var(--surface)`) dan latar topbar tetap `var(--topbar-bg)`; blok
//     soft UI sidebar dan topbar hanya membentuk butir, tombol, kotak cari, lencana, dan pembatas golongan,
//     dengan kontras teks WCAG AA, dan tidak menambah properti penampung `position: fixed`;
//   - gaya soft UI Kelola User terisolasi di kelas akar `.kelola-user`, dan teksnya tetap terbaca (kontras WCAG AA);
//   - tidak satu pun aturan `.kelola-user` memakai properti yang menjadikan elemen blok penampung bagi
//     `position: fixed`. Popup inti (`.modal__backdrop`, `position: fixed; inset: 0`) dirender di dalam
//     `.kelola-user`, bukan lewat portal; properti itu membuat popup menempel ke kotak dan tidak lagi di
//     tengah layar (bug 02-10-2026 pada panel modul bertema `backdrop-filter`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const CSS = readFileSync(join(__dirname, 'styles.css'), 'utf8')

/** Aturan terdalam `pemilih { isi }` (tanpa komentar; aturan di dalam @media ikut). */
function aturan(css: string): Array<{ pemilih: string[]; isi: string }> {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(/([^{};]+)\{([^{}]*)\}/g)].map((m) => ({
    pemilih: (m[1] ?? '').split(',').map((p) => p.trim()),
    isi: m[2] ?? '',
  }))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

function penampungFixed(isi: string): string[] {
  return [...isi.matchAll(PENAMPUNG_FIXED)].map((m) => m[1] ?? '')
}

/** Nilai `background` aturan yang pemilihnya PERSIS salah satu `pemilih`. */
function latar(css: string, pemilih: string[]): string[] {
  return aturan(css)
    .filter((a) => a.pemilih.some((p) => pemilih.includes(p)))
    .flatMap((a) => [...a.isi.matchAll(/(?:^|[;\s])background(?:-color|-image)?\s*:\s*([^;]+)/g)].map((m) => (m[1] ?? '').trim()))
}

/** Rasio kontras WCAG 2 dua warna hex `#rrggbb`. */
function kontras(a: string, b: string): number {
  const terang = (hex: string): number => {
    const [r, g, bl] = [1, 3, 5].map((i) => {
      const v = parseInt(hex.slice(i, i + 2), 16) / 255
      return v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4
    })
    return 0.2126 * (r ?? 0) + 0.7152 * (g ?? 0) + 0.0722 * (bl ?? 0)
  }
  const [t, g] = [terang(a), terang(b)].sort((x, y) => y - x)
  return ((t ?? 0) + 0.05) / ((g ?? 0) + 0.05)
}

/**
 * Token bernilai hex `#rrggbb` dari SEMUA aturan yang pemilihnya PERSIS `pemilih`, digabung berurutan
 * seperti kaskade (aturan belakangan menang).
 */
function token(daftar: Array<{ pemilih: string[]; isi: string }>, pemilih: string): Record<string, string> {
  return Object.fromEntries(
    daftar
      .filter((x) => x.pemilih.length === 1 && x.pemilih[0] === pemilih)
      .flatMap((a) => [...a.isi.matchAll(/(--[a-z0-9-]+)\s*:\s*(#[0-9a-f]{6})\s*;/gi)].map((m) => [m[1] ?? '', m[2] ?? ''])),
  )
}

/** Pasangan [teks, latar] yang kontrasnya di bawah 4,5:1, atau yang tokennya hilang. */
function kontrasKurang(t: Record<string, string>, pasangan: ReadonlyArray<readonly [string, string]>): string[] {
  return pasangan.flatMap(([depan, belakang]) => {
    const a = t[depan]
    const b = t[belakang]
    if (a === undefined || b === undefined) return [`${depan}/${belakang} token hilang`]
    const r = kontras(a, b)
    return r < 4.5 ? [`${depan}/${belakang} ${r.toFixed(2)}`] : []
  })
}

/** Aturan di antara dua kepala blok komentar (dipotong dari pembuka komentar, supaya komentar terbuang utuh). */
function blok(dari: string, sampai: string): Array<{ pemilih: string[]; isi: string }> {
  const awal = CSS.lastIndexOf('/*', CSS.indexOf(dari))
  const akhir = CSS.lastIndexOf('/*', CSS.indexOf(sampai))
  return aturan(CSS.slice(awal, akhir))
}

const AKAR_SIDEBAR = ['.shell__sidebar', ':root[data-theme="dark"] .shell__sidebar']
const AKAR_TOPBAR = ['.shell__topbar', ':root[data-theme="dark"] .shell__topbar']

describe('sidebar dan topbar bergaya soft UI', () => {
  const shell = blok('SIDEBAR DAN TOPBAR BERGAYA SOFT UI', 'KELOLA USER BERGAYA SOFT UI')

  it('latar sidebar tetap latar lama var(--surface), tidak ditimpa blok tema', () => {
    expect(latar(CSS, AKAR_SIDEBAR)).toEqual(['var(--surface)'])
  })

  it('latar topbar tetap var(--topbar-bg), tidak ditimpa blok tema', () => {
    expect(latar(CSS, AKAR_TOPBAR)).toEqual(['var(--topbar-bg)'])
  })

  it('aturan latar menggigit', () => {
    expect(latar('.shell__sidebar { color: red; background: var(--sb-latar) } .shell__butir { background: x }', AKAR_SIDEBAR)).toEqual([
      'var(--sb-latar)',
    ])
  })

  it('setiap golongan GROUPMENU diberi garis pembatas', () => {
    const golongan = aturan(CSS).filter((a) => a.pemilih.includes('.shell__sidebar .shell__golongan'))
    expect(golongan.some((a) => /border-top\s*:\s*1px solid var\(--sb-pembatas\)/.test(a.isi))).toBe(true)
  })

  it('blok tidak menambah properti yang mengurung popup position: fixed', () => {
    expect(shell.length).toBeGreaterThan(20)
    expect(shell.flatMap((a) => penampungFixed(a.isi).map((x) => `${a.pemilih.join(', ')}: ${x}`))).toEqual([])
  })

  it('teks terbaca: kontras token teks sidebar dan topbar terhadap latarnya minimal 4,5:1, terang dan gelap', () => {
    const semua = aturan(CSS)
    const akar = token(semua, ':root')
    const akarGelap = { ...akar, ...token(semua, ':root[data-theme="dark"]') }
    const sb = token(shell, '.shell__sidebar')
    const tb = token(shell, '.shell__topbar')
    const terang = { ...akar, ...sb, ...tb }
    const gelap = {
      ...akarGelap,
      ...sb,
      ...tb,
      ...token(shell, ':root[data-theme="dark"] .shell__sidebar'),
      ...token(shell, ':root[data-theme="dark"] .shell__topbar'),
    }
    // `--sb-teks-redup` sengaja tidak diuji: hanya dipakai butir "belum dimigrasi" yang nonaktif
    // (WCAG 1.4.3 mengecualikan komponen nonaktif).
    const pasangan = [
      ['--sb-teks', '--surface'],
      ['--sb-judul', '--surface'],
      ['--sb-aktif-teks', '--sb-aktif'],
      ['--sb-aktif-teks', '--sb-hover'],
      ['--sb-lencana-teks', '--sb-lencana'],
      ['--sb-lencana-aktif-teks', '--sb-lencana-aktif'],
      ['--tb-teks', '--tb-tombol'],
      ['--tb-teks-redup', '--tb-tombol'],
      ['--tb-teks', '--tb-cari'],
      ['--tb-teks-redup', '--tb-cari'],
      ['--tb-avatar-teks', '--tb-avatar'],
    ] as const
    expect([...kontrasKurang(terang, pasangan).map((x) => `terang ${x}`), ...kontrasKurang(gelap, pasangan).map((x) => `gelap ${x}`)]).toEqual([])
  })
})

describe('Kelola User bergaya soft UI', () => {
  const tema = blok('KELOLA USER BERGAYA SOFT UI', 'BERANDA BERGAYA SOFT UI')

  it('teks tetap terbaca: kontras token teks terhadap latar, kartu, dan kepala tabel minimal 4,5:1, terang dan gelap', () => {
    const terang = token(tema, '.kelola-user')
    const gelap = { ...terang, ...token(tema, ':root[data-theme="dark"] .kelola-user') }
    const pasangan = [
      ['--ku-teks', '--ku-latar'],
      ['--ku-teks', '--ku-kartu'],
      ['--ku-teks-redup', '--ku-latar'],
      ['--ku-teks-redup', '--ku-kartu'],
      ['--ku-teks-redup', '--ku-kepala-tabel'],
      ['--ku-aksen-teks', '--ku-kartu'],
      ['--ku-judul-golongan', '--ku-latar'],
      ['--ku-aksen-teks', '--ku-aksen-lembut'],
      ['--ku-sukses-teks', '--ku-sukses-latar'],
      ['--ku-netral-teks', '--ku-netral-latar'],
      ['--ku-bahaya-teks', '--ku-bahaya-latar'],
      ['--ku-waspada-teks', '--ku-waspada-latar'],
    ] as const
    const gagal = [
      ['terang', terang],
      ['gelap', gelap],
    ].flatMap(([nama, t]) =>
      pasangan
        .map(([depan, belakang]) => {
          const tk = t as Record<string, string>
          const r = kontras(tk[depan] ?? '#ffffff', tk[belakang] ?? '#ffffff')
          return { label: `${String(nama)} ${depan}/${belakang} ${r.toFixed(2)}`, r }
        })
        .filter((x) => x.r < 4.5)
        .map((x) => x.label),
    )
    expect(Object.keys(terang).length).toBeGreaterThan(8)
    expect(gagal).toEqual([])
  })

  it('rumus kontras menggigit: hitam/putih 21:1, #777777/putih di bawah 4,5:1', () => {
    expect(kontras('#000000', '#ffffff')).toBeCloseTo(21, 5)
    expect(kontras('#777777', '#ffffff')).toBeLessThan(4.5)
    expect(kontras('#767676', '#ffffff')).toBeGreaterThanOrEqual(4.5)
  })

  it('blok tema ada dan setiap pemilihnya di bawah kelas akar .kelola-user', () => {
    expect(tema.length).toBeGreaterThan(10)
    const lepas = tema
      .flatMap((a) => a.pemilih)
      .filter((p) => {
        const tanpaTema = p.replace(/^:root\[data-theme="dark"\]\s+/, '')
        return tanpaTema !== '.kelola-user' && !tanpaTema.startsWith('.kelola-user ')
      })
    expect(lepas).toEqual([])
  })

  it('tidak satu pun aturan .kelola-user mengurung popup position: fixed', () => {
    const kena = aturan(CSS)
      .filter((a) => a.pemilih.some((p) => p.includes('kelola-user')))
      .flatMap((a) => penampungFixed(a.isi).map((x) => `${a.pemilih.join(', ')}: ${x}`))
    expect(kena).toEqual([])
  })

  it('aturan penampung menggigit', () => {
    expect(penampungFixed(' -webkit-backdrop-filter: blur(4px); backdrop-filter: blur(4px); transform: none; container-type: inline-size')).toEqual([
      '-webkit-backdrop-filter',
      'backdrop-filter',
      'transform',
    ])
  })
})

describe('Beranda bergaya soft UI', () => {
  // Blok Beranda adalah blok terakhir styles.css.
  const tema = aturan(CSS.slice(CSS.lastIndexOf('/*', CSS.indexOf('BERANDA BERGAYA SOFT UI'))))

  it('blok ada dan setiap pemilihnya di bawah kelas akar .beranda', () => {
    expect(tema.length).toBeGreaterThan(15)
    const lepas = tema
      .flatMap((a) => a.pemilih)
      .filter((p) => {
        const tanpaTema = p.replace(/^:root\[data-theme="dark"\]\s+/, '')
        return tanpaTema !== '.beranda' && !tanpaTema.startsWith('.beranda ')
      })
    expect(lepas).toEqual([])
  })

  it('blok tidak memakai properti yang mengurung popup position: fixed', () => {
    expect(tema.flatMap((a) => penampungFixed(a.isi).map((x) => `${a.pemilih.join(', ')}: ${x}`))).toEqual([])
  })

  it('teks terbaca: kontras token teks Beranda minimal 4,5:1, terang dan gelap', () => {
    const terang = token(tema, '.beranda')
    const gelap = { ...terang, ...token(tema, ':root[data-theme="dark"] .beranda') }
    const pasangan = [
      ['--br-teks', '--br-latar'],
      ['--br-teks', '--br-kartu'],
      ['--br-teks', '--br-baris-hover'],
      ['--br-teks-redup', '--br-latar'],
      ['--br-teks-redup', '--br-kartu'],
      ['--br-teks-redup', '--br-kepala-tabel'],
      ['--br-aksen-teks', '--br-kartu'],
      ['--br-aksen-teks', '--br-ikon'],
      ['--br-sukses-teks', '--br-sukses-latar'],
      ['--br-netral-teks', '--br-netral-latar'],
    ] as const
    expect([...kontrasKurang(terang, pasangan).map((x) => `terang ${x}`), ...kontrasKurang(gelap, pasangan).map((x) => `gelap ${x}`)]).toEqual([])
  })
})

describe('area kerja dilebarkan (permintaan work owner 02-10-2026)', () => {
  /** Nilai `prop` dari setiap aturan yang pemilihnya PERSIS `.shell__isi` (aturan di @media ikut). */
  const nilai = (css: string, prop: string): string[] =>
    aturan(css)
      .filter((a) => a.pemilih.includes('.shell__isi'))
      .flatMap((a) => [...a.isi.matchAll(new RegExp(`(?:^|[;\\s])${prop}\\s*:\\s*([^;]+)`, 'g'))].map((m) => (m[1] ?? '').trim()))

  it('.shell__isi berbatas lebar 1600px, padding samping paling lebar 20px', () => {
    expect(nilai(CSS, 'max-width')).toEqual(['1600px'])
    expect(nilai(CSS, 'padding-inline')).toEqual(['20px'])
  })

  it('aturan nilai menggigit', () => {
    expect(nilai('.shell__isi { max-width: 1400px; } @media (min-width: 1px) { .shell__isi { padding-inline: 32px } }', 'max-width')).toEqual(['1400px'])
  })
})
