// Popup Choose Accumulation tahap 1 (tiket 46) - data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { POPUP_AKUMULASI as P, TEKS_AKUMULASI } from '../labels'
import PopupAkumulasi, { zipDariId } from './PopupAkumulasi'

const SUMBER = readFileSync(join(__dirname, 'PopupAkumulasi.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('zipDariId (SetDataAccum_Act langkah 2-3)', () => {
  it('format NEGARA-ZIP-NOMOR: substring(4,9) tanpa "-"', () => {
    expect(zipDariId('UJI-12345-000001')).toBe('12345')
  })
  it('karakter ke-3 bukan "-": substring(3,8)', () => {
    expect(zipDariId('AB012345X')).toBe('12345')
  })
  it('ID terlalu pendek -> kosong', () => {
    expect(zipDariId('UJI-1')).toBe('')
  })
})

describe('PopupAkumulasi', () => {
  it('awal: judul jendela, Zip Code = zip lokasi risiko, tombol Filter / Clear Column', () => {
    const html = renderToStaticMarkup(<PopupAkumulasi zipRisiko="12345" onTutup={() => {}} onPilih={() => {}} />)
    for (const t of [P.judul, P.cari.label, P.accumulationCode.label, P.road.label, P.zipCode.label, P.czone.label, P.filter.label, P.bersih.label]) {
      expect(html).toContain(t)
    }
    expect(html).toContain('value="12345"')
  })

  it('Choose: zip beda -> pesan verbatim, popup tetap; sama -> onPilih(ID, Note)', () => {
    expect(TEKS_AKUMULASI.zipBeda('12345', '54321')).toBe('ZIpCode Harus Sama dengan ZIpCode  Yang di Object Item >>>> 12345 != 54321')
    expect(SUMBER).toMatch(/if \(dariId !== zipRisiko\) \{\s*setPesan\(TEKS_AKUMULASI\.zipBeda\(zipRisiko, dariId\)\)\s*return\s*\}\s*onPilih\(b\.id, b\.note\)/)
  })
})
