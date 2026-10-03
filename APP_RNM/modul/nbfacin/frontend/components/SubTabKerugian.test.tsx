// Paritas sub-tab Loss Record / Loss Record Internal (tiket 42). Data uji sintetis; uang = teks desimal.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GRID_KERUGIAN, GRID_KLAIM_INTERNAL, LOSS_RATIO, TEKS_INWARD } from '../labels'
import { SubTabKerugian, SubTabKlaimInternal, adaGalatKerugian, kerugianBaru, lossRatioKosong, uangSah } from './SubTabKerugian'

const judulKolom = (html: string) =>
  [...html.matchAll(/<th[^>]*>([\s\S]*?)<\/th>|<th[^>]*\/>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))

describe('SubTabKerugian', () => {
  it('grid kosong + Tambah; empat loss ratio tampil', () => {
    const html = renderToStaticMarkup(<SubTabKerugian rows={[]} lossRatio={lossRatioKosong()} insuredName="UJI" ubah={() => {}} />)
    expect(judulKolom(html)).toEqual(['', ...GRID_KERUGIAN.map((k) => k.label), 'Tambah'])
    expect(html).toContain(TEKS_INWARD.kosong)
    for (const l of LOSS_RATIO) expect(html).toContain(`>${l.label}</span>`)
  })

  it('baris: tanggal dd/mm/yyyy, uang 4 desimal gaya Indonesia; loss ratio sesuai desimal Pega', () => {
    const html = renderToStaticMarkup(
      <SubTabKerugian
        rows={[{ ...kerugianBaru('UJI TERTANGGUNG'), dateOfLoss: '05-02-2026', currency: 'USD', amount: '1500000.12345', claim: '250.5' }]}
        lossRatio={{ oneYearAmount: '1234.5678', oneYearPercent: '12.345', threeFiveYearAmount: '0', threeFiveYearPercent: '0' }}
        insuredName="UJI TERTANGGUNG"
        ubah={() => {}}
      />,
    )
    expect(html).toContain('<td>05/02/2026</td><td>UJI TERTANGGUNG</td>')
    expect(html).toContain('<td class="nbf-angka">1.500.000,1235</td><td class="nbf-angka">250,5</td>')
    expect(html).toContain('>1.234,568</div>')
    expect(html).toContain('>12,35</div>')
  })

  it('catatan baru: Insured Name = tertanggung case, Total Claim "0"; uang tak sah menahan Save', () => {
    const r = kerugianBaru('UJI')
    expect([r.coinsName, r.claim, r.amount]).toEqual(['UJI', '0', ''])
    expect(['', '0', '12.5'].every(uangSah)).toBe(true)
    expect(uangSah('1,5')).toBe(false)
    expect(adaGalatKerugian([{ ...r, currency: 'USD', preventionOfLoss: 'abc' }])).toBe(true)
    // N-9: Currency wajib per catatan.
    expect(adaGalatKerugian([r])).toBe(true)
    expect(adaGalatKerugian([{ ...r, currency: 'USD' }])).toBe(false)
  })
})

describe('SubTabKlaimInternal', () => {
  it('grid baca-saja: sebelas kolom, tanpa Tambah / Hapus, kosong = No items', () => {
    const html = renderToStaticMarkup(<SubTabKlaimInternal rows={[]} />)
    expect(judulKolom(html)).toEqual(GRID_KLAIM_INTERNAL.map((k) => k.label))
    expect(html).not.toContain('Tambah')
    expect(html).toContain(TEKS_INWARD.kosong)
  })

  it('baris: Year = tahun Date of Loss', () => {
    const html = renderToStaticMarkup(
      <SubTabKlaimInternal
        rows={[{ dateOfLoss: '05-02-2026', locationNo: '1', location: 'UJI', currency: 'USD', premium: '100', osClaim: '0', acceptedClaim: '0', incurredClaim: '0', lossRatio: '0', remark: '' }]}
      />,
    )
    expect(html).toContain('<td>2026</td><td>05/02/2026</td><td>1</td><td>UJI</td><td>USD</td><td class="nbf-angka">100</td>')
  })
})

describe('SubTabKerugian - medan uang berformat ribuan', () => {
  it('Total Claim dan Prevention Of Loss = IsianUang', async () => {
    const { readFileSync } = await import('node:fs')
    const { join } = await import('node:path')
    const s = readFileSync(join(__dirname, 'SubTabKerugian.tsx'), 'utf8')
    expect(s).toContain('<IsianUang label={F.claim.label}')
    expect(s).toMatch(/<IsianUang\s*label=\{F\.preventionOfLoss\.label\}/)
  })
})
