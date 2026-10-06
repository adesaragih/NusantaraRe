// Audit silang putaran 3 (W2, W6): subsection `DetailPoliciesNonProportional` ->
// `DetailPolicyTreatyInNonProportional` DIRENDER (react-dom/server). Harapan dari XML
// (dibaca ulang 04-10-2026):
//
//   judul wadah  S1 "Limits", S19 "Share", S41 "Total Spreaded", S51 "Share Facultative"
//                (`pyIncludeHeader=true`); S73 (SpreadingRiskList) dan S74 (angsuran)
//                `NOHEADER` - TANPA judul
//   W2           grid SpreadingRiskList terbuka (Add/Delete, %Share) bila layar boleh
//                disunting DAN FacultativeShare 0/'' - di layar atasan juga (S88)
//   format       LimitShareSummaryList `.Limit` pxNumber TANPA pyDecimalPlaces (pola inti);
//                LimitFacShareSummaryList seluruh sel tanpa Format (apa adanya);
//                TotalLimitIOONP / TotalFacShareRnmNP kolom pertama FIELD kosong

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Baris, Halaman } from '../api'
import DetailNonProp from './DetailNonProp'

// Sel angka rata kanan (`nbti__angka`, perintah work owner 06-10-2026); format nilainya tidak berubah.

const hal = (nilai: Record<string, string>, daftar: Record<string, Baris[]> = {}): Halaman => ({ nilai, daftar })

const render = (h: Halaman, sunting: boolean) =>
  renderToStaticMarkup(
    <DetailNonProp
      halaman={h}
      sunting={sunting}
      tempat={{}}
      opsiSpreading={[]}
      onUbahBaris={() => {}}
      onRefresh={() => {}}
    />,
  )

const judulWadah = (html: string) => [...html.matchAll(/<h4[^>]*>([^<]*)<\/h4>/g)].map((m) => m[1])

describe('DetailNonProp = DetailPolicyTreatyInNonProportional', () => {
  const spreading: Record<string, Baris[]> = {
    'PolicyTreatyIn.SpreadingRiskList': [{ TreatyType: 'UJI-SPR', SharePercentage: '100', PremiumSpreaded: '2700' }],
  }
  const dasar = hal({ 'TreatyIn.FacultativeShare': '0' }, spreading)

  it('judul wadah hanya yang ber-pyIncludeHeader; spreading dan angsuran tanpa judul', () => {
    expect(judulWadah(render(dasar, false))).toEqual(['Limits', 'Share', 'Total Spreaded'])
    expect(judulWadah(render(hal({ 'TreatyIn.FacultativeShare': '5' }), false))).toEqual([
      'Limits',
      'Share',
      'Total Spreaded',
      'Share Facultative',
    ])
    const html = render(dasar, false)
    expect(html).not.toContain('Spreading Risk')
  })

  // Perintah work owner 06-10-2026: Add / Delete dibuang dan Treaty Type hanya-baca ("ga boleh di ubah lagi");
  // yang tetap tersunting (W2: layar boleh disunting dan FacultativeShare 0) hanya %Share / %Share Claim.
  it('W2: tanpa Add / Delete, Treaty Type hanya-baca; %Share tersunting bila boleh dan FacultativeShare 0', () => {
    const html = render(dasar, true)
    expect(html).not.toContain('>Add</button>')
    expect(html).not.toContain('>Delete</button>')
    expect(html).not.toContain('<select')
    expect(html).toContain('UJI-SPR')
    expect(html).toMatch(/<input[^>]*nbti__angka/)
    expect(render(dasar, false)).not.toMatch(/<input[^>]*nbti__angka/)
    const fak = hal({ 'TreatyIn.FacultativeShare': '5' }, spreading)
    expect(render(fak, true)).not.toMatch(/<input[^>]*nbti__angka/)
  })

  it('format sel grid master: Share .Limit tanpa desimal tetap, Share Facultative apa adanya', () => {
    const h = hal(
      { 'TreatyIn.FacultativeShare': '5' },
      {
        'TreatyIn.LimitShareSummaryList': [{ Note: 'UJI', Limit: '1000', Limit2: '7' }],
        'TreatyIn.LimitFacShareSummaryList': [{ Note: 'UJI-F', Limit: '2500.5', Limit2: '3' }],
      },
    )
    const html = render(h, false)
    expect(html).toContain('<td class="nbti__angka">1.000</td><td class="nbti__angka">7,00</td>')
    expect(html).toContain('<td>UJI-F</td><td class="nbti__angka">2500.5</td><td class="nbti__angka">3</td>')
  })

  it('TotalLimitIOONP dan TotalFacShareRnmNP: kolom pertama FIELD kosong, judul di atas .Value', () => {
    const h = hal(
      { 'TreatyIn.FacultativeShare': '5' },
      {
        'TreatyIn.TotalLimitIOONP': [{ Currency: 'IDR', Value: '5000' }],
        'TreatyIn.TotalFacShareRnmNP': [{ Currency: 'USD', Value: '7' }],
      },
    )
    const html = render(h, false)
    // kepala kolom angka rata kanan, sejajar dengan selnya
    expect(html).toContain('<th scope="col"></th><th scope="col"></th><th scope="col" class="nbti__angka">Total Limit</th>')
    expect(html).toContain('<td></td><td>IDR</td><td class="nbti__angka">5.000,00</td>')
    expect(html).toContain('<td></td><td>USD</td><td class="nbti__angka">7,00</td>')
  })
})
