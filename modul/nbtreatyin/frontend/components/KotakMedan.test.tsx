// Medan radio (`pxRadioButtons`, mis. Survey Report): DIRENDER (react-dom/server). Nilai standar disimpan, teks
// prompt ditampilkan; nilai tersimpan di luar daftar tidak dibuang (ditambahkan sebagai pilihan sendiri).

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Medan } from '../medan'
import KotakMedan from './KotakMedan'

const medan: Medan = {
  jalur: 'PolicyTreatyIn.QuotationData.IsSurveyReport',
  label: 'Survey Report',
  jenis: 'radio',
  opsi: [
    { value: 'Yes', label: 'Yes' },
    { value: 'No', label: 'No' },
  ],
}

const render = (v: string, hanyaBaca = false, wajib = false) =>
  renderToStaticMarkup(
    <KotakMedan
      medan={medan}
      halaman={{ nilai: { [medan.jalur]: v }, daftar: {} }}
      wajib={wajib}
      hanyaBaca={hanyaBaca}
      opsiMataUang={[]}
      opsiMO={[]}
      onUbah={() => undefined}
      onSelesai={() => undefined}
    />,
  )

const radio = (html: string) => [...html.matchAll(/<input[^>]*type="radio"[^>]*>/g)].map((m) => m[0])

describe('KotakMedan radio', () => {
  it('dua pilihan Yes / No dengan label; nilai tersimpan tercentang', () => {
    const html = render('No')
    const r = radio(html)
    expect(r).toHaveLength(2)
    expect(r.filter((x) => x.includes('checked'))).toHaveLength(1)
    expect(r[1]).toContain('checked')
    expect(html).toContain('>Yes</label>')
    expect(html).toContain('>No</label>')
    expect(html).toContain('Survey Report')
  })

  it('wajib: tanda bintang di label', () => {
    expect(render('', false, true)).toContain('field__req')
  })

  it('nilai di luar daftar tetap tampil sebagai pilihan tercentang', () => {
    const r = radio(render('Lama'))
    expect(r).toHaveLength(3)
    expect(r[2]).toContain('checked')
  })

  it('hanya-baca: teks prompt, tanpa radio', () => {
    const html = render('Yes', true)
    expect(radio(html)).toEqual([])
    expect(html).toContain('Yes')
  })
})

// Isian angka (permintaan work owner 06-10-2026): rata kanan, papan angka, kosong/nol = placeholder "0", tidak
// difokus = format sel (bagian uang 4 desimal); nilai hanya-baca angka rata kanan, nol tampil "0" redup.
describe('KotakMedan angka', () => {
  const uang: Medan = { jalur: 'PolicyTreatyIn.PremiOgp', label: 'Premi Ogp', jenis: 'angka', sajian: { desimal: 4, nolPolos: true, formatSaatSunting: true } }
  const tampil = (m: Medan, v: string, hanyaBaca = false) =>
    renderToStaticMarkup(
      <KotakMedan
        medan={m}
        halaman={{ nilai: { [m.jalur]: v }, daftar: {} }}
        wajib={false}
        hanyaBaca={hanyaBaca}
        opsiMataUang={[]}
        opsiMO={[]}
        onUbah={() => undefined}
        onSelesai={() => undefined}
      />,
    )
  const isian = (html: string) => /<input[^>]*>/.exec(html)?.[0] ?? ''

  it('isian angka: rata kanan, papan angka, label terhubung', () => {
    const i = isian(tampil(uang, '1234.5'))
    expect(i).toContain('nbti__angka')
    expect(i).toContain('inputMode="decimal"')
    expect(i).toContain('value="1.234,5000"')
    expect(tampil(uang, '1234.5')).toMatch(/<label[^>]*for="[^"]+"/)
  })

  it('nol dan kosong: isian kosong dengan placeholder 0', () => {
    for (const v of ['0', '']) {
      const i = isian(tampil(uang, v))
      expect(i).toContain('placeholder="0"')
      expect(i).toContain('value=""')
    }
  })

  it('hanya-baca: angka rata kanan; nol redup', () => {
    expect(tampil(uang, '12', true)).toContain('nbti__nilai--angka')
    const nol = tampil(uang, '0', true)
    expect(nol).toContain('nbti__nilai--nol')
    expect(nol).toMatch(/>0<\/span>/)
  })
})
