import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

// Permintaan pemakai 11-10-2026 (grid Share / Limits / total NB NonProp): "kalau nilai
// pada kolom tidak ada atau kosong, isi 0,00 jangan kosong".
describe('DetailNonProp - sel angka kosong tampil 0,00', () => {
  const src = readFileSync(join(__dirname, 'components/DetailNonProp.tsx'), 'utf8')

  it('sel angka kosong menjadi 0,00 sebelum format apa pun (termasuk kolom mentah)', () => {
    expect(src).toContain("const NOL_TAMPIL = '0,00'")
    const i = src.indexOf('function teksSel(')
    const badan = src.slice(i, src.indexOf('\n}\n', i))
    expect(badan.indexOf('if (angkaKosong(b, c)) return NOL_TAMPIL')).toBeGreaterThan(-1)
    expect(badan.indexOf('if (angkaKosong(b, c)) return NOL_TAMPIL')).toBeLessThan(badan.indexOf("c.format === 'mentah'"))
  })

  it('Note, Currency, dan tanggal tidak diisi 0,00', () => {
    expect(src).toContain("const kolomAngka = (c: Kolom) => c.m !== '' && !TEKS.has(c.m) && c.m !== 'DueDate'")
    expect(src).toContain("const TEKS = new Set(['Note', 'Currency'])")
  })

  it('sel 0,00 ikut rata kanan seperti angka lain', () => {
    expect(src).toContain('className={kelasAngka(nilaiKelas(b, c))}')
  })
})
