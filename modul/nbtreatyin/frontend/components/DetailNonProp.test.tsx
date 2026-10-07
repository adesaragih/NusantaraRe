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

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

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

  it('layar memberi sunting HANYA bila admin pemegang berkas (keputusan work owner 07-10-2026)', () => {
    const sumber = readFileSync(join(__dirname, '../pages/LayarKasus.tsx'), 'utf8')
    const blok = sumber.slice(sumber.indexOf('<DetailNonProp'), sumber.indexOf('/>', sumber.indexOf('<DetailNonProp')))
    expect(blok).toContain('sunting={ubahAdmin}')
    expect(blok).not.toContain('sunting={boleh}')
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

// Screenshot layar Pega NonProp (work owner 06-10-2026: "tampilan NB Treaty Non Prop belum sesuai").
describe('DetailNonProp mengikuti screenshot layar Pega', () => {
  it('spreading: kolom % RNM Share (.SplitRNMSharePct) di depan, "% Share" menjadi "% Proportion Share"', () => {
    const h = hal(
      { 'TreatyIn.FacultativeShare': '0' },
      {
        'PolicyTreatyIn.SpreadingRiskList': [
          { TreatyName: 'UJI-QS', SplitRNMSharePct: '14', SharePercentage: '73.68', PremiumSpreaded: '100', ClaimPercentage: '73.68', ClaimSpreaded: '0' },
        ],
      },
    )
    const html = render(h, false)
    const kepala = [...html.matchAll(/<th scope="col"[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])
    expect(kepala).toEqual(
      expect.arrayContaining(['Type Treaty', '% RNM Share', '% Proportion Share', 'Premium', 'Claim']),
    )
    expect(kepala.filter((x) => x === '% Proportion Share')).toHaveLength(2)
    expect(html).not.toContain('>% Share<')
    expect(html).toContain('<td>UJI-QS</td><td class="nbti__angka">14,00</td><td class="nbti__angka">73,68</td>')
    // baris Total: sel % RNM Share kosong, total lain di bawah kolomnya
    expect(html).toContain('<td>Total</td><td></td>')
  })

  it('grid total sebaris dalam satu wadah; grid tanpa patah baris', () => {
    const h = hal(
      { 'TreatyIn.FacultativeShare': '0' },
      {
        'TreatyIn.TotalLimitIOONP': [{ Currency: 'IDR', Value: '1' }],
        'TreatyIn.TotalLimitDeductblNP': [{ Currency: 'IDR', Value: '2' }],
        'TreatyIn.TotalLimitMDPNP': [{ Currency: 'IDR', Value: '3' }],
      },
    )
    const html = render(h, false)
    const wadah = html.split('<div class="nbti__total-np">')[1]!.split('</table></div></div>')[0]!
    expect((wadah.match(/<div class="table-wrap nbti__grid-np">/g) ?? []).length).toBe(3)
    expect(html).toContain('Total Deductible')
  })

  it('Installment sebaris; angsuran per mata uang dapat dilipat, rincian bernomor, saldo dua desimal', () => {
    const h = hal(
      { 'TreatyIn.FacultativeShare': '0', 'PolicyTreatyIn.Installment': '1' },
      {
        'PolicyTreatyIn.ListInstallment': [{ Currency: 'UJI' }],
        'PolicyTreatyIn.ListInstallment(1).InstallmentList': [
          { DueDate: '2026-07-01', InstallmentPercentage: '25', Currency: 'UJI', Premium: '10.5', PremiumAfterPPN: '7', PremiumAfterTax: '8' },
        ],
      },
    )
    const html = render(h, false)
    expect(html).toContain('<div class="nbti__kolom nbti__kolom-np">')
    expect(html).toContain('<details class="nbti__rinci-np" open=""><summary>UJI</summary>')
    expect(html).toContain('<td class="nbti__urut-np">1</td><td>01-07-2026</td><td class="nbti__angka">25</td>')
    expect(html).toContain('<td class="nbti__angka">10,50</td><td class="nbti__angka">7,00</td><td class="nbti__angka">8,00</td>')
  })
})
