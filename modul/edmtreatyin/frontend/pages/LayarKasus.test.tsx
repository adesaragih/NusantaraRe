// Layar kasus endorsemen DIRENDER (react-dom/server) dari keadaan awal `awal` (fixture UJI-, tanpa fetch: render
// statis tidak menjalankan efek; `fetch` distub dan diperiksa tidak terpanggil). Harapan dari XML
// `DetailPolicyTreatyInAddendum`, `DetailPolicyTreatyInAddGeneralEditable`, `DetailPolicyTreatyInAddPremi`,
// `ListSuggestEDM` (dibaca 06-10-2026).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import type { Acuan, Baris, Halaman, Layar, TombolKirim } from '../api'
import LayarKasus from './LayarKasus'

const acuan: Acuan = {
  mataUang: [{ nilai: '1', label: 'UJI-IDR' }],
  mo: [{ nilai: 'UJI-MO-1', label: 'UJI MO' }],
  jenisReas: [{ nilai: '7', label: 'UJI-QS' }],
  jenisEdm: [{ nilai: '3', label: '3' }],
}

function layar(
  posisi: string,
  nilai: Record<string, string>,
  daftar: Record<string, Baris[]> = {},
  o: Partial<Layar> = {},
): Layar {
  const halaman: Halaman = {
    nilai: { 'Quotation.ProportionalType': 'Proportional', 'PolicyTreatyIn.IsNewPolicyNonProp': '0', ...nilai },
    daftar,
  }
  return {
    kasus: {
      id: 'UJI-EDMT-1',
      position: posisi === 'ReasTreatyInAdmin' ? '4' : '5',
      statusWork: '',
      positionNote: posisi,
      noPolis: '',
      prodKe: 1,
      oldPolisId: 'UJI-POLIS-0',
      generasiTertutup: false,
      createOp: 'UJI-AKUN',
      tglCreate: '2026-10-06 09:00:00',
    },
    halaman,
    bolehKerja: true,
    tombol: '',
    medanWajib: null,
    ...o,
  }
}

const render = (ly: Layar) =>
  renderToStaticMarkup(<LayarKasus id="UJI-EDMT-1" onKembali={() => {}} awal={{ layar: ly, acuan }} />)

const ambil = vi.fn(() => Promise.reject(new Error('UJI: fetch tidak boleh dipanggil render statis')))
beforeEach(() => vi.stubGlobal('fetch', ambil))
afterEach(() => {
  expect(ambil).not.toHaveBeenCalled()
  vi.unstubAllGlobals()
})

describe('header DetailPolicyTreatyInAddendum', () => {
  it('admin: judul "Input Realitation", Choose Business, header terkunci kecuali sel Auto, Remark dapat diisi', () => {
    const html = render(
      layar('ReasTreatyInAdmin', {
        'PolicyTreatyIn.EDMNo': 'UJI-POLIS/E01',
        'PolicyTreatyIn.EDMType': '3',
        'PolicyTreatyIn.IDCurrency': '1',
        'PolicyTreatyIn.QuotationData.MOID': 'UJI-MO-1',
      }),
    )
    expect(html).toContain('Input Realitation — UJI-EDMT-1')
    expect(html).toContain('>Choose Business</button>')
    expect(html).toContain('data-jalur="PolicyTreatyIn.EDMNo">UJI-POLIS/E01<')
    // dropdown hanya-baca: teks pilihan acuan
    expect(html).toContain('data-jalur="PolicyTreatyIn.IDCurrency">UJI-IDR<')
    expect(html).toContain('data-jalur="PolicyTreatyIn.EDMType">3<')
    // Marketing Officer pxDropdown Auto
    expect(html).toMatch(/Marketing Officer[\s\S]*?<select/)
    // Remark pxTextArea: admin dapat mengisi
    expect(html).toMatch(/Remark[\s\S]*?<textarea/)
    expect(html).not.toContain('Survey Report')
    expect(html).not.toContain('Claim Payment Type')
  })

  it('atasan: tanpa Choose Business, Remark hanya-baca, sel header Auto TERKUNCI (keputusan work owner 07-10-2026)', () => {
    const html = render(layar('ReasTreatyInSecHead', { 'PolicyTreatyIn.Remark': 'UJI-CATATAN' }))
    expect(html).toContain('Acceptance by Sec Treaty — UJI-EDMT-1')
    expect(html).not.toContain('Choose Business')
    expect(html).toContain('data-jalur="PolicyTreatyIn.Remark">UJI-CATATAN<')
    // XML tidak mengunci With Tax / Type Tax / Overiding Commision / Marketing Officer per posisi; WO: atasan TIDAK
    // boleh mengubahnya
    expect(html).toMatch(/<input type="checkbox"[^>]*disabled[^>]*\/>With Tax/)
    expect(html).toMatch(/<input type="checkbox"[^>]*disabled[^>]*\/>Overiding Commision/)
  })

  it('tanpa bolehKerja: pemberitahuan hanya-baca, tanpa isian Approval dan tanpa tombol kaki', () => {
    const html = render(layar('ReasTreatyInAdmin', {}, {}, { bolehKerja: false, tombol: 'kirim' }))
    expect(html).toContain('hanya-baca')
    expect(html).not.toContain('Choose Business')
    expect(html).not.toContain('edmt-approval')
    expect(html).not.toContain('>Submit</button>')
    expect(html).not.toContain('>Save</button>')
  })
})

describe('tab Old Data / New Data / Value Difference (IsNewPolicyNonProp 0)', () => {
  it('strip tab VERBATIM, tab awal Old Data hanya-baca berawalan OldData', () => {
    const html = render(layar('ReasTreatyInAdmin', { 'PolicyTreatyIn.OldData.PremiOgp': '1500' }))
    const tab = [...html.matchAll(/role="tab" aria-selected="(true|false)"[^>]*>([^<]*)</g)].map(
      (m) => `${m[2]}:${m[1]}`,
    )
    expect(tab).toEqual(['Old Data:true', 'New Data:false', 'Value Difference:false'])
    expect(html).toContain('data-jalur="PolicyTreatyIn.OldData.PremiOgp">1.500,0000<')
    expect(html).toContain('<h4>OGP</h4>')
    expect(html).toContain('<h4>ONP</h4>')
    expect(html).toContain('% Deduction In A (OGP)')
    expect(html).toContain('Installment Data Information')
    // tab hanya-baca: tanpa isian angka, tanpa Add / Calculate
    expect(html).not.toMatch(/<input[^>]*edmt__angka/)
    expect(html).not.toContain('>Add</button>')
    expect(html).not.toContain('Calculate Value Difference')
  })

  it('Old Data = PropOldData2 (awalan OldData.TreatyDifference) bila OldData.EDMNo terisi', () => {
    const html = render(
      layar('ReasTreatyInSecHead', {
        'PolicyTreatyIn.OldData.EDMNo': 'UJI-POLIS/E01',
        'PolicyTreatyIn.OldData.TreatyDifference.PremiOgp': '7',
      }),
    )
    expect(html).toContain('data-jalur="PolicyTreatyIn.OldData.TreatyDifference.PremiOgp">7,0000<')
  })

  it('tanpa tombol Delete di grid spreading mana pun; Type Treaty hanya-baca menampilkan teks acuan', () => {
    const html = render(
      layar(
        'ReasTreatyInAdmin',
        {},
        { 'PolicyTreatyIn.OldData.SpreadingRiskList': [{ TreatyType: '7', SharePercentage: '100' }] },
      ),
    )
    expect(html).toContain('<td>UJI-QS</td>')
    expect(html).not.toContain('Delete')
  })
})

describe('cabang NonProp baru (IsNewPolicyNonProp 1) = AddPremi', () => {
  const xol = { Currency: 'UJI-USD', GrossPremi: '1000', NetPremi: '900' }
  const ly = layar(
    'ReasTreatyInAdmin',
    { 'PolicyTreatyIn.IsNewPolicyNonProp': '1', 'PolicyTreatyIn.Installment': '2' },
    {
      'PolicyTreatyIn.OldData.TreatyXOLList': [xol],
      'PolicyTreatyIn.TreatyXOLList': [xol],
      'PolicyTreatyIn.TreatyXOLDifferenceList': [{ Currency: 'UJI-USD', GrossPremi: '0' }],
      'PolicyTreatyIn.TreatyXOLList(1).ValueList': [
        { LayerType: 'UJI-layer', Layer: '1', Currency: 'UJI-USD', GrossPremi: '400' },
      ],
      'PolicyTreatyIn.ListInstallment': [{ Currency: 'UJI-USD', Premium: '900' }],
      'PolicyTreatyIn.ListInstallment(1).InstallmentList': [{ DueDate: '2026-11-01', Premium: '450' }],
    },
  )

  it('tiga wadah XOL berjudul + angsuran per mata uang; tanpa tab', () => {
    const html = render(ly)
    const judul = [...html.matchAll(/<h4 class="panel__title">([^<]*)<\/h4>/g)].map((m) => m[1])
    expect(judul).toEqual(['Previous Premium', 'Current Premium', 'Total Difference', 'Installment Data Information'])
    expect(html).not.toContain('role="tablist"')
    expect(html).toContain('Total For Currency')
  })

  it('rincian per layer (ValueList) terbuka di bawah barisnya, hanya-baca, 4 desimal', () => {
    const html = render(ly)
    expect(html).toContain('<td>UJI-layer</td><td class="edmt__angka">1</td><td>Part of</td>')
    expect(html).toContain('<td class="edmt__angka">400,0000</td>')
    expect(html).toContain('aria-expanded="true"')
    // rincian angsuran: Payment Date tanggal sistem, saldo 2 desimal
    expect(html).toContain('<td>01-11-2026</td>')
    expect(html).toContain('<td class="edmt__angka">450,00</td>')
  })

  // perintah work owner 07-10-2026 (screenshot Current Premium): "JANGAN NULL TAPI 0" + "BUAT 4 ANGKA BELAKANG KOMA"
  it('grid XOL dan rinciannya: angka 4 desimal, pajak kosong tampil 0 rata kanan', () => {
    const html = render(ly)
    expect(html).toContain('<td class="edmt__angka">1.000,0000</td>')
    expect(html).toContain('<td class="edmt__angka">900,0000</td>')
    // induk: PPN, PPh, Net Premium After PPN / PPh / Tax kosong -> 0 (Deduction juga kosong di fixture)
    const induk = html.match(/<td class="edmt__angka">1\.000,0000<\/td>(.*?)<\/tr>/)?.[1] ?? ''
    expect([...induk.matchAll(/<td class="edmt__angka">0<\/td>/g)]).toHaveLength(6)
    // rincian layer: Deduction, PPN, PPh, Net Premium, After PPN / PPh / Tax kosong -> 0
    const rinci = html.match(/<td class="edmt__angka">400,0000<\/td>(.*?)<\/tr>/)?.[1] ?? ''
    expect([...rinci.matchAll(/<td class="edmt__angka">0<\/td>/g)]).toHaveLength(7)
  })

  it('Installment admin: isian pyMaxLength 2', () => {
    expect(render(ly)).toMatch(/<input[^>]*maxLength="2"[^>]*value="2"/i)
  })
})

describe('ListSuggestEDM dan tombol Submit', () => {
  const t = (tombol: TombolKirim, posisi = 'ReasTreatyInAdmin') => render(layar(posisi, {}, {}, { tombol }))

  it('Approval radio Accept / Reject wajib, Suggest wajib, grid Date / PIC / Approval / Suggest, tanpa Production Date', () => {
    const html = render(
      layar(
        'ReasTreatyInAdmin',
        { 'PolicyTreatyIn.IsApproved': '1' },
        {
          'PolicyTreatyIn.SuggestList': [
            { Date: '2026-10-06 17:52:59', OperatorName: 'UJI PIC', IsApproved: '0', Suggest: 'UJI' },
          ],
        },
      ),
    )
    expect(html).toMatch(/Approval<span class="field__req">\*<\/span>/)
    expect(html).toMatch(/Suggest<span class="field__req">\*<\/span>/)
    expect(html).toContain('>Accept</label>')
    expect(html).toContain('>Reject</label>')
    expect(html).toContain('<td>06-10-2026 17:52:59</td><td>UJI PIC</td><td>Reject</td>')
    expect(html).not.toContain('Production Date')
  })

  it('Submit menurut Layar.tombol; admin selalu punya Save', () => {
    expect(t('')).not.toContain('>Submit</button>')
    expect(t('')).toContain('>Save</button>')
    expect(t('kirim')).toContain('>Submit</button>')
    expect(t('konfirmasi-tolak')).toContain('>Submit</button>')
    expect(t('nomor-polis', 'ReasTreatyInDeptHead')).toContain('>Submit</button>')
    expect(t('nomor-polis', 'ReasTreatyInDeptHead')).not.toContain('>Save</button>')
  })

  it('pesan halaman (Layar.pesan) tampil', () => {
    expect(render(layar('ReasTreatyInAdmin', {}, {}, { pesan: ['UJI pesan layar'] }))).toContain('UJI pesan layar')
  })
})

describe('popup (kode sumber)', () => {
  const src = readFileSync(join(__dirname, 'LayarKasus.tsx'), 'utf8')

  it('ShowPolicyNoTreaty_SC: Modal tanpaTutup, satu aksi OK yang mengirim, nilai PolicyTreatyIn.PolicyNo', () => {
    const mulai = src.indexOf('{nomor && (')
    const popup = src.slice(mulai, src.indexOf('</Modal>', mulai))
    expect(mulai).toBeGreaterThan(0)
    expect(popup).toContain('tanpaTutup')
    expect(popup).toContain('{TOMBOL.ok}')
    expect(popup).toContain('kirim()')
    expect(popup).toContain("nilai(h, POLIS + 'PolicyNo')")
  })

  it('PolicyTreatyInDeclineConfirm: Yes mengirim, No menutup', () => {
    const mulai = src.indexOf('{konfirmasi && (')
    const popup = src.slice(mulai, src.indexOf('</Modal>', mulai))
    expect(popup).toContain('labelBatal={TOMBOL.no}')
    expect(popup).toContain('{TOMBOL.yes}')
    expect(popup).toContain('{KONFIRMASI_TOLAK}')
  })

  it('Approval: aksi backend SetDueTo lalu Protection (satu urutan)', () => {
    expect(src).toContain("refresh([{ aksi: 'SetDueTo' }, { aksi: 'Protection' }])")
  })
})
