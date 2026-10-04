// Paritas sub-tab Surrounding Risk (tiket 38) - section Pega `RiskAround`. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { LAIN_SEKITAR as L, MEDAN_SISI, SISI_SEKITAR, TEKS_SEKITAR } from '../labels'
import SubTabSekitar, { adaJarakMinus, jarakMinus, pilihOccupation, sekitarKosong, sisiKosong } from './SubTabSekitar'
import { objekBaru } from './TabObject'

const O = objekBaru('1')
const HTML = renderToStaticMarkup(<SubTabSekitar o={O} ubah={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'SubTabSekitar.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('SubTabSekitar - tampilan', () => {
  it('empat blok berurutan Front, Left, Back, Right, lalu dua blok Other Description (tangkapan layar)', () => {
    const judul = [...HTML.matchAll(/<h5 class="nbf-objek__judul">([^<]*)<\/h5>/g)].map((m) => m[1])
    expect(judul).toEqual([...SISI_SEKITAR.map((s) => s.judul), L.judul.label, L.judul.label])
  })

  it('objek baru: Construction "Silahkan pilih"; Ownership / House keeping / Flood Area Status = Not Informed; paragraf faktor', () => {
    expect(HTML.split('<option value="" selected="">Silahkan pilih</option>')).toHaveLength(5)
    expect(HTML.split('<option value="2" selected="">Not Informed</option>')).toHaveLength(3)
    expect(HTML).toContain('<option value="0" selected="">Not Informed</option>')
    expect(HTML).toContain(`<p class="nbf-sekitar__faktor">${TEKS_SEKITAR.faktorRisiko}</p>`)
  })

  it('tiap sisi: Occupation, Construction, Distance (meter), Note', () => {
    const label = [...HTML.matchAll(/class="field__label"[^>]*>([^<]*)/g)].map((m) => m[1])
    const sisi = [MEDAN_SISI.occupation, MEDAN_SISI.construction, MEDAN_SISI.distance, MEDAN_SISI.note]
    expect(label.slice(0, 16)).toEqual([...sisi, ...sisi, ...sisi, ...sisi])
  })

  it('Other Description: Flood Area tersembunyi selama Flood Area Status bukan 0; tiga centang', () => {
    expect(HTML).toContain(`>${L.ownership.label}<`)
    expect(HTML).toContain(`>${L.floodAreaStatus.label}<`)
    expect(HTML).not.toContain(`>${L.floodArea.label}<`)
    // Flood Area Status "0" = Yes.
    const tampil = renderToStaticMarkup(
      <SubTabSekitar o={{ ...O, surroundingRisk: { ...sekitarKosong(), floodAreaStatus: '0' } }} ubah={() => {}} />,
    )
    expect(tampil).toContain(`>${L.floodArea.label}<`)
    for (const c of [L.productionProcess, L.hotWork, L.flammable]) expect(HTML).toContain(`/> ${c.label}`)
  })
})

describe('SubTabSekitar - aturan', () => {
  it('Distance minus = pesan NegativeIsNotAllowed; kosong / nol / desimal positif sah', () => {
    expect(jarakMinus('-0.5')).toBe(true)
    expect(jarakMinus('0')).toBe(false)
    expect(jarakMinus('12.25')).toBe(false)
    expect(jarakMinus('')).toBe(false)
    expect(adaJarakMinus({ ...sekitarKosong(), left: { ...sisiKosong(), distance: '-1' } })).toBe(true)
    const minus = renderToStaticMarkup(
      <SubTabSekitar o={{ ...O, surroundingRisk: { ...sekitarKosong(), back: { ...sisiKosong(), distance: '-3' } } }} ubah={() => {}} />,
    )
    expect(minus.split(TEKS_SEKITAR.jarakMinus)).toHaveLength(2)
  })

  it('memilih Occupation: nilai = OldID, Note = Name; Construction dan Distance tetap', () => {
    const s = pilihOccupation({ occupation: 'ab', construction: 'K1', distance: '4', note: '' }, { oldId: 'UJI-01', name: 'UJI NAMA' })
    expect(s).toEqual({ occupation: 'UJI-01', construction: 'K1', distance: '4', note: 'UJI NAMA' })
  })

  it('saran Occupation hanya saat mengetik, >= 2 karakter, sesudah jeda; jawaban lama dibuang', () => {
    expect(SUMBER).toContain('if (!ketik || sisi.occupation.trim().length < MIN_CARI)')
    expect(SUMBER).toContain('const MIN_CARI = 2')
    expect(SUMBER).toContain('if (n === nomor.current) setSaran(h.baris)')
  })
})
