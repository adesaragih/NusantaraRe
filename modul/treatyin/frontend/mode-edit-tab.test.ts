// Mode Edit harus membuat SELURUH tab dapat disunting.
//
// ---------------------------------------------------------------------
// ⛔ KEPUTUSAN PEMILIK PROSES
// ---------------------------------------------------------------------
// "pastikan di aplikasi semua button nya rumusnya bisa digunakan semua nya".
//
// Pencocokan dengan ekspor 6 Oktober 2026 menemukan tiga tab yang TETAP
// baca-saja walau tombol Edit ditekan — Accumulation, Exclusions, dan
// Special Conditions. Ketiganya di Pega dapat diisi.
//
// ⚠️ Cacatnya DIAM: tab yang menolak suntingan terlihat persis seperti tab
// yang memang kosong datanya, dan yang membukanya akan menagih pemuatan data
// padahal yang kurang kemampuan menyuntingnya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const FORM = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
const TEKS = readFileSync(join(AKAR, 'components', 'TabTeksPanjang.tsx'), 'utf8')
const GRID = readFileSync(join(AKAR, 'components', 'TabGridWarisan.tsx'), 'utf8')

/** Potongan sumber satu cabang tab, dari penandanya sampai `/>` penutup. */
function cabangTab(nama: string): string {
  const i = FORM.indexOf(`tabTampil === '${nama}'`)
  expect(i).toBeGreaterThan(0)
  const j = FORM.indexOf('/>', i)
  expect(j).toBeGreaterThan(i)
  return FORM.slice(i, j)
}

describe('setiap tab menerima mode form', () => {
  // ⛔ Ketiganya SEBELUMNYA tidak menerimanya, dan itu sebab mereka
  // baca-saja: komponennya mendukung `mode`, pemanggilnya yang lupa.
  it('⛔ Accumulation, Exclusions, Special Conditions menerima mode', () => {
    for (const tab of ['Accumulation', 'Exclusions', 'Special Conditions']) {
      expect(cabangTab(tab)).toContain('mode={mode}')
    }
  })

  it('EGNPI juga menerimanya — grid yang sama, kelalaian yang sama', () => {
    expect(cabangTab('EGNPI')).toContain('mode={mode}')
  })
})

describe('tab teks panjang dapat diketik di mode Ubah', () => {
  it('⭐ merender textarea, bukan teks baca-saja', () => {
    expect(TEKS).toContain('<textarea')
    expect(TEKS).toContain("const bisaUbah = mode === 'ubah'")
  })

  // ⛔ Tab yang menolak isian PERTAMA tidak akan pernah terisi. Kosong di
  // mode Ubah harus tetap dapat diketik, bukan menampilkan `Kosong`.
  it('⛔ teks KOSONG pun dapat diketik di mode Ubah', () => {
    const i = TEKS.indexOf('{bisaUbah ? (')
    const j = TEKS.indexOf("isi === '' ? (")
    expect(i).toBeGreaterThan(0)
    expect(j).toBeGreaterThan(i)
  })

  it('huruf isiannya diwarisi — rupa tidak berubah saat Edit ditekan', () => {
    const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
    expect(CSS).toContain('.trin__teks-isian')
    expect(CSS).toMatch(/\.trin__teks-isian \{[^}]*font-family: inherit/s)
  })
})

describe('grid warisan menghormati mode', () => {
  it('tombol Add/Delete hanya dirender di mode ubah', () => {
    expect(GRID).toContain("const bisaUbah = mode === 'ubah'")
    expect(GRID).toContain('const tombol = bisaUbah && bisaTambah')
  })
})
