// Audit 8 Oktober 2026 — "jangan semisal isi data table kosong tp bisa
// ditambah". Penyesuaian yang isi tabnya TIDAK tersedia (`terdarat = false`)
// dibuka dalam mode LIHAT di panel New: grid kosong tidak dapat ditambah
// lalu disimpan sebagai data separuh. Di Pega penyesuaian selalu membawa
// salinan dokumen utuh, jadi grid tidak pernah kosong karena datanya hilang.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PENYESUAIAN } from './labelsPenyesuaian'
import { simpanTampil } from './pages/PenyesuaianKontrak'

const HALAMAN = readFileSync(join(__dirname, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')

describe('penyesuaian tanpa isi tab — panel New terkunci', () => {
  it('terdarat=false memaksa mode LIHAT untuk panel New dan deret tombol; panel Old tidak berubah', () => {
    expect(HALAMAN).toContain("const modeBaru: ModeLayar = terkunciPendaratan ? '1' : mode")
    expect(HALAMAN).toMatch(/bacaSaja=\{false\}\s*mode=\{modeBaru\}/)
    expect(HALAMAN).toMatch(/<DeretTombol\s*mode=\{modeBaru\}/)
    expect(HALAMAN).toMatch(/bacaSaja\s*mode=\{mode\}/)
  })

  it('mode lihat = Save tidak tampil (syarat tombol Save)', () => {
    expect(simpanTampil('1', {})).toBe(false)
    expect(simpanTampil('0', {})).toBe(true)
  })

  it('pemakai diberi tahu sebabnya — tanpa istilah teknis', () => {
    expect(HALAMAN).toContain('{PENYESUAIAN.belumTerdarat}')
    expect(PENYESUAIAN.belumTerdarat).toContain('Treaty In')
  })
})
