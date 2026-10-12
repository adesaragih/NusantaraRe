import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { RINCI_PESERTA, RINCI_RETRO, RINCI_SPREADING } from '../kolomKorpus'
import { RINCI_EDM } from '../labels'
import { BAGIAN_RINCI, kelompokRinci, MEDAN_RINCI_LEBAR, selAngka } from '../tampilan'
import { MedanKorpus, TabelKorpus } from './TabelKorpus'

// Popup Details peserta dirapikan - keputusan work owner 08-10-2026 (tampilan saja).

describe('kelompokRinci: setiap medan rincian tampil tepat satu kali', () => {
  for (const t of ['QR', 'QP', 'TP', 'TR'] as const) {
    it(t, () => {
      const bagian = kelompokRinci(RINCI_PESERTA[t])
      const tampil = bagian.flatMap((b) => b.kolom.map((k) => k.kolom))
      expect([...tampil].sort()).toEqual(RINCI_PESERTA[t].map((k) => k.kolom).sort())
      expect(new Set(tampil).size).toBe(tampil.length)
      expect(bagian.map((b) => b.judul)).toEqual([
        RINCI_EDM.bagianPeserta,
        RINCI_EDM.bagianPeriode,
        RINCI_EDM.bagianPertanggungan,
        RINCI_EDM.bagianPremi,
      ])
    })
  }

  it('tidak ada medan di dua bagian', () => {
    const semua = BAGIAN_RINCI.flatMap((b) => b.kolom)
    expect(new Set(semua).size).toBe(semua.length)
  })
})

describe('angka spreading/retro diformat dan rata kanan', () => {
  it('kolom angka ber-jenis n, label tidak berubah', () => {
    expect(RINCI_SPREADING.map((k) => [k.kolom, k.jenis])).toEqual([
      ['treatyTypeName', 't'],
      ['retrocadedShare', 'n'],
    ])
    expect(RINCI_RETRO.filter((k) => k.jenis === 't').map((k) => k.kolom)).toEqual(['reinsurerName'])
  })

  it('format id: 5850038.75, 2047513.5625, 0.58', () => {
    expect(selAngka('5850038.75')).toBe('5.850.038,75')
    expect(selAngka('2047513.5625')).toBe('2.047.513,5625')
    expect(selAngka('0.58')).toBe('0,58')
    const html = renderToStaticMarkup(
      <TabelKorpus
        kolom={RINCI_RETRO}
        baris={[{ reinsurerName: 'UJI-REAS', percentShare: '0.58', amount: '5850038.75', premiumSpreadedGross: '2047513.5625' }]}
        kunci={(_, i) => String(i)}
      />,
    )
    expect(html).toContain('UJI-REAS')
    expect(html).toMatch(/class="[^"]*edm-angka[^"]*">5\.850\.038,75</)
    expect(html).toMatch(/class="[^"]*edm-angka[^"]*">2\.047\.513,5625</)
    expect(html).toMatch(/class="[^"]*edm-angka[^"]*">0,58</)
  })
})

describe('teks panjang tidak bertabrakan (perbaikan 08-10-2026)', () => {
  const CSS = readFileSync(join(__dirname, '..', 'endorsementlife.css'), 'utf8').replace(/\r\n/g, '\n')
  const aturan = (pemilih: string) => {
    const i = CSS.indexOf(pemilih + ' {')
    expect(i, pemilih).toBeGreaterThanOrEqual(0)
    return CSS.slice(i, CSS.indexOf('}', i))
  }

  it('kisi edm-rinci: item boleh menyusut, nilai dipecah; medan lebar selebar kisi', () => {
    expect(aturan('.edm-rinci__grid > .edm-kepala__medan')).toContain('min-width: 0')
    expect(aturan('.edm-rinci__grid dd')).toContain('overflow-wrap: anywhere')
    expect(aturan('.edm-rinci__grid dd')).toContain('min-width: 0')
    expect(aturan('.edm-rinci__lebar')).toContain('grid-column: 1 / -1')
  })

  it('tabel retro: nama dibungkus, angka tidak', () => {
    const r = aturan('.edm-rinci__lapisan td:not(.edm-angka)')
    expect(r).toContain('white-space: normal')
    expect(r).toContain('overflow-wrap: anywhere')
  })

  it('Certificate No dan Policy Holder ditandai lebar di Participant (keempat tipe)', () => {
    expect(MEDAN_RINCI_LEBAR).toEqual(['CERTIFICATE_NO', 'POLICY_HOLDER', 'PLAN'])
    for (const t of ['QR', 'QP', 'TP', 'TR'] as const) {
      const peserta = kelompokRinci(RINCI_PESERTA[t]).find((b) => b.judul === RINCI_EDM.bagianPeserta)
      expect(peserta, t).toBeDefined()
      const html = renderToStaticMarkup(
        <MedanKorpus
          kolom={peserta?.kolom ?? []}
          nilai={{ CERTIFICATE_NO: 'UJI-CERT-SANGATPANJANGTANPASPASI0001', POLICY_HOLDER: 'UJI-PEMEGANG' }}
          className="edm-kepala edm-rinci__grid"
          lebar={MEDAN_RINCI_LEBAR}
          kelasLebar="edm-rinci__lebar"
        />,
      )
      const lebar = [...html.matchAll(/<div class="edm-kepala__medan edm-rinci__lebar"><dt>([^<]*)<\/dt>/g)].map((m) => m[1])
      const label = (k: string) => RINCI_PESERTA[t].find((x) => x.kolom === k)?.label
      expect(lebar, t).toContain(label('CERTIFICATE_NO'))
      expect(lebar, t).toContain(label('POLICY_HOLDER'))
      // medan lain tetap setengah kisi
      expect(html.split('class="edm-kepala__medan"').length - 1, t).toBe(peserta!.kolom.length - lebar.length)
    }
  })

  it('MedanKorpus tanpa lebar: pemakai lain tidak berubah', () => {
    const html = renderToStaticMarkup(<MedanKorpus kolom={RINCI_PESERTA.QR.slice(0, 3)} nilai={{}} />)
    expect(html).not.toContain('edm-rinci__lebar')
    expect(html).toContain('<dl class="edm-kepala">')
  })
})
