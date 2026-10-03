// Paritas popup Choose Risk Address (tiket 36) - harness Pega `ChooseRiskAddress` + tangkapan layar work owner
// 03-10-2026. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { POPUP_RISK as R } from '../labels'
import PopupRiskAddress, { adaSaring } from './PopupRiskAddress'

const KOSONG = { address: '', zipCode: '', country: '', province: '', city: '', district: '', territory: '' }
const HTML = renderToStaticMarkup(<PopupRiskAddress awal={{ ...KOSONG, zipCode: '99999' }} onTutup={() => {}} onPilih={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'PopupRiskAddress.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('PopupRiskAddress', () => {
  it('judul Choose Risk Location; tujuh saringan berurutan, diisi awal dari objek', () => {
    expect(HTML).toContain(`>${R.judul}</h3>`)
    const label = [...HTML.matchAll(/<label class="field__label"[^>]*>([^<]*)/g)].map((m) => m[1])
    expect(label).toEqual(Object.values(R.saring).map((s) => s.label))
    expect(HTML).toContain('value="99999"')
  })

  it('kolom grid Type … Zip Code lalu kolom Pilih; grid kosong sampai Search (H-2)', () => {
    const kolom = [...HTML.matchAll(/<th[^>]*>([^<]*)<\/th>|<th[^>]*\/>/g)].map((m) => m[1] ?? '')
    expect(kolom).toEqual([...R.kolom.map((k) => k.label), ''])
    expect(HTML).not.toContain('<tbody')
  })

  it('kaki: Search = kirim form (Enter juga), Add nonaktif (H-4)', () => {
    expect(HTML).toMatch(new RegExp(`<button type="submit" class="btn btn--primary">${R.cari.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--ghost" disabled="">${R.tambah.label}</button>`))
    expect(SUMBER).toContain('onKirim={() => void muat(saring, 1)}')
  })

  it('Search tanpa saringan tidak memanggil server, menampilkan pesan; paging memakai saringan Search terakhir', () => {
    expect(adaSaring(KOSONG)).toBe(false)
    expect(adaSaring({ ...KOSONG, territory: ' x ' })).toBe(true)
    expect(adaSaring({ ...KOSONG, city: '   ' })).toBe(false)
    expect(SUMBER).toMatch(/if \(!adaSaring\(s\)\) \{\s*setKurang\(true\)\s*return/)
    expect(SUMBER).toContain('onPindah={(h) => void muat(saringCari.current, h)}')
  })

  it('baris: delapan kolom berurutan, Pilih memanggil onPilih', () => {
    expect(SUMBER).toMatch(
      /<td>\{b\.title\}<\/td>\s*<td>\{b\.address\}<\/td>\s*<td>\{b\.nationName\}<\/td>\s*<td>\{b\.provinceName\}<\/td>\s*<td>\{b\.cityName\}<\/td>\s*<td>\{b\.districtName\}<\/td>\s*<td>\{b\.territoryName\}<\/td>\s*<td>\{b\.postalCode\}<\/td>/,
    )
    expect(SUMBER).toContain('onClick={() => onPilih(b)}')
  })
})
