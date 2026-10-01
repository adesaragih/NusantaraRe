// Penjaga tema halaman login di styles.css inti (permintaan work owner 02-10-2026):
//
//   - latar sidebar TETAP latar lama (`var(--surface)`); blok tema sidebar hanya mewarnai butir,
//     pembatas golongan, dan aksen, tidak menimpa latar;
//   - tema Kelola User terisolasi di kelas akar `.kelola-user`;
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

const AKAR_SIDEBAR = ['.shell__sidebar', ':root[data-theme="dark"] .shell__sidebar']

describe('sidebar bertema login', () => {
  it('latar sidebar tetap latar lama var(--surface), tidak ditimpa blok tema', () => {
    expect(latar(CSS, AKAR_SIDEBAR)).toEqual(['var(--surface)'])
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
})

describe('Kelola User bertema login', () => {
  // Potong dari pembuka komentar kepala blok, supaya komentar itu ikut terbuang utuh.
  const tema = aturan(CSS.slice(CSS.lastIndexOf('/*', CSS.indexOf('KELOLA USER BERTEMA HALAMAN LOGIN'))))

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
