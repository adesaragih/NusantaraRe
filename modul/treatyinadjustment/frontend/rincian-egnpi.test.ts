// ⭐ Rincian EGNPI layar Adjustment = tangkapan Pega pemakai (9 Oktober 2026):
//
//   Treaty Group   MARINE CARGO
//   As At          01/01/2025
//   Amount         [IDR v]  [700.000.000,00 ..........]
//   Amount in IDR  [IDR  ]  [700.000.000,00 ..........]
//   Proportion %   3,7037037037
//   Note           ——
//
// Dulu: baris `Amount in IDR` hilang (sel `.pyTemplateRichTextEditor` dibuang
// pembangkit), labelnya `Text Area` (nama kontrol designer), dan kotak nilai
// bergeser karena kolom label kosong.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KERANGKA_RINCIAN } from './ekspor/kerangka.gen'
import { kelasBlok } from './komponen/KerangkaTab'

const CSS = readFileSync(join(__dirname, 'treatyinadjustment.css'), 'utf8')
const KT = readFileSync(join(__dirname, 'komponen', 'KerangkaTab.tsx'), 'utf8')

describe('rincian EGNPI & Max Retention Adjustment — seperti Pega', () => {
  for (const nama of ['DetailEGNPI', 'DetailEGNPIOldData']) {
    it(`${nama}: kotak satuan Amount in IDR [IDR] dan label Note`, () => {
      const teks = JSON.stringify(KERANGKA_RINCIAN[nama])
      expect(teks).toContain('"t":"satuan","at":')
      expect(teks).toContain('"label":"Amount in IDR","teks":"IDR"')
      expect(teks).toContain('"label":"Note"')
      expect(teks).not.toContain('"label":"Text Area"')
    })
  }

  it('satuan dirender kotak baca-saja; pasangan 30/70 tanpa kolom label kosong', () => {
    expect(KT).toContain("case 'satuan':")
    expect(KT).toContain('<Field key={b.at} label={b.label} value={b.teks} readOnly onChange={tanpaAksi} />')
    expect(CSS).toContain('.tria__blok.tria__tata--t3070:has(> .tria__tanpa-label:last-child)')
    expect(CSS).toContain('.tria__blok.tria__tata--kiri > .tria__pemicu > .field {')
    expect(CSS).toContain('.tria__blok.tria__tata--kiri > .field > .pilih-saring,')
    // ⛔ `:has()` TIDAK boleh bersarang — aturan seperti itu dibuang browser.
    expect(CSS).not.toMatch(/:has\([^)]*:has\(/)
  })

  it('kelasBlok menandai blok yang semua medannya tanpa label (sel 70% Amount)', () => {
    const blok = (anak: unknown[]) => ({ t: 'blok', at: 1, judul: '', syarat: [], tata: 'kiri', anak }) as never
    const medan = (label: string) => ({ t: 'medan', at: 2, label, dari: 'sisi', kunci: 'Amount', format: 'pxTextInput', desimal: null, syarat: [] })
    expect(kelasBlok(blok([medan('')]))).toBe('tria__blok tria__tata--kiri tria__tanpa-label')
    expect(kelasBlok(blok([medan('Amount')]))).toBe('tria__blok tria__tata--kiri')
    expect(kelasBlok(blok([]))).toBe('tria__blok tria__tata--kiri')
  })
})
