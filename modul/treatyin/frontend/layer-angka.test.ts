import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

// Permintaan pemakai 10 Oktober 2026: kotak Layer dan Part (tab Limits
// Non-Prop) harus BOLEH diisi huruf. Penyaring angka GitHub (`saringAngka`)
// menolaknya; keduanya kini `teks` — dikecualikan dari penyaring. Medan angka
// lain (Adjustment Rate, MDP %, …) tetap disaring.
describe('tab Limits Non-Prop — Layer dan Part boleh huruf', () => {
  const src = readFileSync(join(__dirname, 'components/TabLimitsNonProp.tsx'), 'utf8')

  it('Layer dan Part ber-`teks` (tanpa penyaring angka)', () => {
    expect(src).toMatch(/nilai=\{teksDari\(l, 'Layer'\)\}[^\n]*onUbah=\{set\('Layer'\)\} teks/)
    expect(src).toMatch(/nilai=\{teksDari\(l, 'LayerPart'\)\}[^\n]*onUbah=\{set\('LayerPart'\)\} teks/)
  })

  it('medan lain tanpa `teks` tetap disaring saringAngka', () => {
    expect(src).toContain('onUbah(teks ? v : saringAngka(v))')
    for (const k of ['AdjRate', 'MDPMinPct', 'MDPPct']) {
      const baris = src.split('\n').filter((b) => b.includes(`'${k}'`) && b.includes('MedanNP'))
      for (const b of baris) expect(b).not.toContain(' teks')
    }
  })
})

// ⛔ Laporan pemakai 10 Oktober 2026: ROL % tampil 33.390.000.000 (Pega 33,39).
// `saringAngka` menganggap titik desimal, jadi limit `2.000.000.000` menjadi
// `2.000000000` = 2. Kotak uang di matriks limit dan sel Reinstatement WAJIB
// `FieldAngka` (titik = ribuan), bukan input mentah + `saringAngka`.
describe('tab Limits Non-Prop — kotak uang tidak merusak titik ribuan', () => {
  const src = readFileSync(join(__dirname, 'components/TabLimitsNonProp.tsx'), 'utf8')

  it('matriks limit dan sel Reinstatement memakai FieldAngka', () => {
    expect(src).not.toContain('onUbah({ ...l, [kunci]: saringAngka(e.target.value) })')
    expect(src).not.toContain('ubahSel(saringAngka(e.target.value))')
    expect((src.match(/<FieldAngka/g) ?? []).length).toBeGreaterThanOrEqual(3)
  })
})
