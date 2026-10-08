// Tombol `Add Revision` / `Add Adjustment Premium` dan picker-nya —
// `Section/InputTreatyInAdjustment.xml` @500554 / @515429,
// `Section/PickerTreatyInMasterRevisi.xml`, `Section/PickerTreatyInMaster.xml`.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Penyesuaian } from './api'
import PilihMaster, { judulPicker, selMaster } from './komponen/PilihMaster'
import { KOLOM_PICKER, LEBAR_PICKER_PREMI, LEBAR_PICKER_REVISI, PENYESUAIAN } from './labelsPenyesuaian'
import { modeDraf } from './pages/PenyesuaianKontrak'

const AKAR = __dirname
const HALAMAN = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
const PICKER = readFileSync(join(AKAR, 'komponen', 'PilihMaster.tsx'), 'utf8')

const tanpaAksi = () => undefined

describe('tombol Add di layar daftar', () => {
  it('⭐ KEDUANYA HIDUP dan membuka picker-nya masing-masing', () => {
    const i = HALAMAN.indexOf('function Daftar')
    const blok = HALAMAN.slice(HALAMAN.indexOf('tria__aksi', i), HALAMAN.indexOf('<PilihMaster', i))
    expect(blok).toContain('PENYESUAIAN.tambahRevisi')
    expect(blok).toContain('PENYESUAIAN.tambahPremi')
    expect(blok).not.toContain('disabled')
    expect(blok).toContain("setPicker('revisi')")
    expect(blok).toContain("setPicker('premi')")
  })

  it('draf dari Choose membuka layar detail TANPA dibaca ulang dari server', () => {
    expect(HALAMAN).toContain('setBuka({ id: p.id, mode: modeDraf(p), draf: p })')
    expect(HALAMAN).toMatch(/if \(draf !== undefined\) \{\s*setP\(draf\)\s*return/)
    expect(HALAMAN).toContain('PENYESUAIAN.drafBelumTersimpan')
  })

  it('mode draf = ViewState TreatyInSetEdit — 0 Edit, 1 bila RevisionState = 1', () => {
    const draf = (v: string): Penyesuaian => ({
      id: 'x',
      idAsal: 'y',
      baru: { medan: { ViewState: v }, larik: {} },
      lama: { medan: {}, larik: {} },
    })
    expect(modeDraf(draf('0'))).toBe('0')
    expect(modeDraf(draf('1'))).toBe('1')
  })
})

describe('picker Add Revision', () => {
  const html = renderToStaticMarkup(<PilihMaster jenis="revisi" onTutup={tanpaAksi} onDraf={tanpaAksi} />)

  it('judul mengikuti EDMState — TreatyCreateEDM menyetel "1" lebih dulu', () => {
    expect(html).toContain(PENYESUAIAN.pickerRevisiInternal)
    expect(judulPicker('revisi', '2')).toBe('Choose Master to Create External Revision')
    expect(judulPicker('premi', '3')).toBe('Choose Master to create Premium Adjustment')
  })

  it('radio Internal / External (Internal tercentang) dan Material Type (kosong)', () => {
    expect(html).toContain('Internal / External')
    // Urutan atribut milik React — dicocokkan tanpa bergantung padanya.
    expect(html).toMatch(/<input(?=[^>]*name="tria-picker-edmstate")(?=[^>]*value="1")(?=[^>]*checked="")[^>]*>/)
    expect(html).not.toMatch(/<input(?=[^>]*name="tria-picker-edmstate")(?=[^>]*value="2")(?=[^>]*checked="")[^>]*>/)
    expect(html).not.toMatch(/<input(?=[^>]*name="tria-picker-material")(?=[^>]*checked="")[^>]*>/)
    expect(html).toContain('Material Type')
  })

  it('Choose mengirim parameter TreatyInEDMSetValue persis Section-nya', () => {
    expect(PICKER).toContain("internalType: jenis === 'premi' ? '3' : edmState")
    expect(PICKER).toContain("materialType: jenis === 'premi' ? '1' : material")
  })
})

describe('picker Add Adjustment Premium', () => {
  const html = renderToStaticMarkup(<PilihMaster jenis="premi" onTutup={tanpaAksi} onDraf={tanpaAksi} />)

  it('tanpa radio — hanya judul dan grid', () => {
    expect(html).toContain(PENYESUAIAN.pickerPremi)
    expect(html).not.toContain('tria-picker-edmstate')
    expect(html).not.toContain('tria-picker-material')
  })
})

describe('grid picker', () => {
  it('tombol Choose di kolom PERTAMA, lalu tujuh kolom CARI1…CARI7', () => {
    expect(KOLOM_PICKER).toEqual(['ID', 'Contract Name', 'Reinsurance Type', 'Source of Business', 'Ceding', 'Commencement', 'Termination'])
    expect(LEBAR_PICKER_REVISI).toHaveLength(8)
    expect(LEBAR_PICKER_PREMI).toHaveLength(8)
    expect(LEBAR_PICKER_REVISI[0]).toBe(109)
    expect(PICKER.indexOf('PENYESUAIAN.pilihMaster}\n')).toBeLessThan(PICKER.indexOf('selMaster(b).map'))
  })

  it('Commencement / Termination diformat sebagai tanggal (pxDateTime)', () => {
    const sel = selMaster({
      id: '1000506',
      namaKontrak: 'N',
      sifatProporsi: 'NonProportional',
      asalBisnis: 'S',
      cedant: 'C',
      tanggalMulai: '20250101',
      tanggalBerakhir: '20251231',
    })
    expect(sel[0]).toBe('1000506')
    expect(sel[5]).not.toBe('20250101')
    expect(sel[6]).not.toBe('20251231')
  })
})
