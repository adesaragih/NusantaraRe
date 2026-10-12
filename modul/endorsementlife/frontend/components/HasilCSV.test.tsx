import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { barisUnggahan, KEPALA_UNDUH_CSV } from '../csv'
import { HASIL_CSV } from '../kolomKorpus'
import { kolomGrid } from '../tampilan'
import HasilCSV, { kolomViewUpload } from './HasilCSV'

// View Upload = kolom grid Participants (keputusan work owner 10-10-2026 "Sama dengan grid peserta").

const SUMBER = readFileSync(join(__dirname, 'HasilCSV.tsx'), 'utf8').replace(/\r\n/g, '\n')

function kepala(html: string): string[] {
  const t = html.slice(html.indexOf('<thead>'), html.indexOf('</thead>'))
  return [...t.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1] ?? '')
}

describe('View Upload memakai kolom grid Participants', () => {
  it('kolom tepat kolomGrid(): urutan dan label sama, tanpa kolom No, tanpa PLAN/POLICY_HOLDER', () => {
    expect(kolomViewUpload().map((k) => [k.label, k.kolom])).toEqual(kolomGrid().map(([l, k]) => [l, k]))
    const html = renderToStaticMarkup(<HasilCSV baris={[{ POLICY_NO: 'UJI-P1' }]} onTutup={() => {}} />)
    expect(kepala(html)).toEqual(kolomGrid().map(([l]) => l))
    expect(kepala(html)).toEqual(['POLICY NO', 'CERTIFICATE NO', 'NAME OF INSURED', 'DATE OF BIRTH'])
    expect(kolomViewUpload().map((k) => k.kolom)).not.toContain('PLAN')
    expect(kolomViewUpload().map((k) => k.kolom)).not.toContain('POLICY_HOLDER')
    expect(SUMBER).not.toContain('import { HASIL_CSV')
    expect(SUMBER).not.toContain('HASIL_CSV.map')
    expect(HASIL_CSV.length).toBe(38) // bukti korpus tetap
  })

  it('baris dari berkas ";" dan "," tampil apa adanya (tanggal DD/MM/YYYY tidak diurai)', () => {
    for (const teks of [
      'POLICY_NO,CERTIFICATE_NO,NAME_OF_INSURED,DOB,SUM_INSURED\nUJI-P1,UJI-C1,UJI-NAMA-1,22/11/1990,1500000.50\n',
      'POLICY_NO;CERTIFICATE_NO;NAME_OF_INSURED;DOB;SUM_INSURED\nUJI-P1;UJI-C1;UJI-NAMA-1;22/11/1990;1500000,50\n',
    ]) {
      const html = renderToStaticMarkup(<HasilCSV baris={barisUnggahan(teks)} onTutup={() => {}} />)
      for (const v of ['UJI-P1', 'UJI-C1', 'UJI-NAMA-1', '22/11/1990']) expect(html).toContain(`>${v}<`)
      expect(html).not.toContain('1500000')
    }
  })

  it('Generate Data Detail tetap KEPALA_UNDUH_CSV (23 kolom CSVPropHeaders b276, tidak berubah)', () => {
    expect(KEPALA_UNDUH_CSV.length).toBe(23)
    expect(SUMBER).toContain('rakitCSV(KEPALA_UNDUH_CSV, baris)')
  })
})
