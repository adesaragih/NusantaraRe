// Isi tab New Data DIRENDER (react-dom/server): `PropNewData2` (admin, dapat diisi) dan `PropNewData` (atasan,
// hanya-baca). Harapan dari XML (dibaca 06-10-2026). Fixture UJI-.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Baris, Halaman } from '../api'
import { tataTab, type VarianTab } from '../medan'
import TabData from './TabData'

const h = (nilai: Record<string, string>, daftar: Record<string, Baris[]> = {}): Halaman => ({ nilai, daftar })
const kosong = { mataUang: [], mo: [], jenisEdm: [] }

const render = (varian: VarianTab, halaman: Halaman, boleh = true, wajib: string[] = []) =>
  renderToStaticMarkup(
    <TabData
      tata={tataTab(varian)}
      halaman={halaman}
      wajib={new Set(wajib)}
      boleh={boleh}
      opsi={kosong}
      opsiJenisReas={[{ nilai: '7', label: 'UJI-QS' }]}
      onUbah={() => {}}
      onSelesai={() => {}}
      onUbahBaris={() => {}}
      onRefresh={() => {}}
    />,
  )

describe('PropNewData2 (admin)', () => {
  const halaman = h(
    { 'PolicyTreatyIn.PremiOgp': '1000', 'PolicyTreatyIn.Installment': '2' },
    { 'PolicyTreatyIn.SpreadingRiskList': [{ TreatyType: '7', SharePercentage: '100', PremiumSpreaded: '1000' }] },
  )

  it('medan uang = isian angka rata kanan; wajib dari backend bertanda bintang', () => {
    const html = render('baruAdmin', halaman, true, ['PolicyTreatyIn.PremiOgp'])
    expect(html).toMatch(/<input[^>]*class="field__input edmt__angka"[^>]*aria-label="Premi Ogp"/)
    expect(html).toMatch(/Premi Ogp<span class="field__req">\*<\/span>/)
    expect(html).toContain('(%) Deduction In A (OGP)')
    expect(html).toContain('Total Premium Before Claim')
    expect(html).toContain('Outstanding Claim')
  })

  // keputusan work owner 07-10-2026 ("INI READ ONLY JUGA", sama dengan NB): tanpa Add / Delete, Type Treaty teks
  it('grid spreading hanya-baca: TANPA Add / Delete, Type Treaty teks; hanya %Share tersunting', () => {
    const html = render('baruAdmin', halaman)
    expect(html).not.toContain('>Add</button>')
    expect(html).not.toContain('Delete')
    expect(html).not.toMatch(/<select[^>]*aria-label="Type Treaty"/)
    expect(html).toContain('<td>UJI-QS</td>')
    expect(html).toMatch(/<input[^>]*aria-label="% Share"/)
  })

  // keputusan work owner 07-10-2026 ("HAPUS AJA"): Save di dalam tab dibuang - Save kaki halaman (S19) tetap
  // WO 08-10-2026 ("COBA CEK TOMBOLITU, APAKAH MASIH DIPERLUKAN? KALAU SUDAH TIDAK DIHAPUS AJA!"): tombol S24 dibuang
  it('tanpa Save di dalam tab; tanpa Calculate Value Difference; isian Installment', () => {
    const html = render('baruAdmin', halaman)
    expect(html).not.toContain('>Save</button>')
    expect(html).not.toContain('Calculate Value Difference')
    expect(html).toMatch(/<input[^>]*aria-label="Installment"[^>]*value="2"/)
  })

  it('tanpa bolehKerja: seluruhnya hanya-baca, tanpa tombol', () => {
    const html = render('baruAdmin', halaman, false)
    expect(html).not.toMatch(/<input/)
    expect(html).not.toContain('<select')
    expect(html).not.toContain('<button')
    expect(html).toContain('data-jalur="PolicyTreatyIn.PremiOgp">1.000,0000<')
  })
})

describe('PropNewData (atasan)', () => {
  it('hanya-baca, label "% Deduction", Net Premium, total berlabel 2 desimal, tanpa Installment', () => {
    const html = render('baru', h({ 'PolicyTreatyIn.TotalPremium': '1234.5', 'PolicyTreatyIn.Installment': '3' }))
    expect(html).not.toMatch(/<input/)
    expect(html).toContain('% Deduction In A (OGP)')
    expect(html).toContain('>Net Premium<')
    expect(html).toContain('data-jalur="PolicyTreatyIn.TotalPremium">1.234,50<')
    expect(html).not.toContain('>Installment<')
    expect(html).not.toContain('Outstanding Claim')
  })
})
