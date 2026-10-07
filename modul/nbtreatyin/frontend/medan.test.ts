// Uji definisi medan layar - syarat VERBATIM section Pega (INVENTARIS bab 5;
// sel + wadah `pyContainerVisibleWhen` dibaca ulang 2026-10-03).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  KELOMPOK_UANG_ADMIN,
  KELOMPOK_UANG_ATASAN,
  MEDAN_ADMIN_UANG,
  MEDAN_ADMIN_UMUM,
  MEDAN_ATASAN_TOTAL,
  MEDAN_ATASAN_UANG,
  MEDAN_ATASAN_UMUM,
  medanTampil,
  saldoNegatif,
  teksPilihan,
  type Medan,
} from './medan'

const P = 'PolicyTreatyIn.'
const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })
const label = (ms: { label: string }[]) => ms.map((m) => m.label)
const tampil = (ms: Medan[], nilai: Record<string, string>) => label(medanTampil(ms, hal(nilai)))
const cari = (ms: Medan[], jalur: string, lbl?: string) => {
  const m = ms.find((x) => x.jalur === P + jalur && (lbl === undefined || x.label === lbl))
  if (!m) throw new Error(`medan ${jalur} tidak ada`)
  return m
}

describe('medan layar NB Treaty In - syarat tampil', () => {
  it('saldo: "Balance Due To You" bila BalanceDueTo < 0, "Due To Us" bila >= 0 (kosong = 0)', () => {
    expect(saldoNegatif(hal({ [P + 'BalanceDueTo']: '-0.01' }))).toBe(true)
    expect(saldoNegatif(hal({ [P + 'BalanceDueTo']: '-0' }))).toBe(false)
    expect(saldoNegatif(hal({}))).toBe(false)
    const neg = tampil(MEDAN_ADMIN_UANG, { [P + 'BalanceDueTo']: '-5' })
    expect(neg).toContain('Balance Due To You')
    expect(neg).not.toContain('Balance Due To Us')
    const pos = tampil(MEDAN_ADMIN_UANG, { [P + 'BalanceDueTo']: '5' })
    expect(pos).toContain('Balance Due To Us')
    expect(pos).toContain('Balance Before Tax')
  })

  it('admin FlagRetroTreaty tampil bila ClaimType != "XOL Retro" (DetailPolicyTreatyIn pyVisible)', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, {})).toContain('Overiding Commision')
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'ClaimType']: 'XOL Retro' })).not.toContain('Overiding Commision')
  })

  // Audit silang P3: sel tanpa pyLabelFieldValue memakai teks yang TAMPIL di XML -
  // pyCheckboxCaption kotak centang, LABEL "Q" / "/" / "U/Y" / "Of" di depan sel.
  it('teks tampil XML: caption kotak centang dan LABEL pendamping sel tanpa label', () => {
    expect(cari(MEDAN_ADMIN_UMUM, 'FlagRetroTreaty').label).toBe('Overiding Commision')
    expect(cari(MEDAN_ADMIN_UMUM, 'FlagPPH').label).toBe('With Tax')
    expect(cari(MEDAN_ATASAN_UMUM, 'FlagPPH').label).toBe('Include Tax')
    for (const ms of [MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UMUM]) {
      expect(cari(ms, 'Quartal').label).toBe('Q')
      expect(cari(ms, 'YearOfQuartal').label).toBe('/')
      expect(cari(ms, 'LayerPartType').label).toBe('Of')
    }
  })

  it('admin FlagPPH dan Type Tax di wadah ".ClaimType != XOL Retro"; Type Tax juga bila FlagPPH true', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, {})).toContain('With Tax')
    expect(tampil(MEDAN_ADMIN_UMUM, {})).not.toContain('Type Tax')
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'true' })).toContain('Type Tax')
    const retro = tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'true', [P + 'ClaimType']: 'XOL Retro' })
    expect(retro).not.toContain('With Tax')
    expect(retro).not.toContain('Type Tax')
  })

  it('Survey Report disembunyikan untuk NonProportional', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, { 'Quotation.ProportionalType': 'NonProportional' })).not.toContain('Survey Report')
  })

  it('Quartal / YearOfQuartal / U/Y hanya Proportional, Layer* hanya NonProportional (wadah, kedua layar)', () => {
    for (const ms of [MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UMUM]) {
      const prop = tampil(ms, { [P + 'QuotationData.ProportionalType']: 'Proportional' })
      expect(prop).toEqual(expect.arrayContaining(['Q', '/', 'U/Y']))
      expect(prop).not.toContain('LayerType')
      const np = tampil(ms, { [P + 'QuotationData.ProportionalType']: 'NonProportional' })
      expect(np).not.toContain('Q')
      expect(np).not.toContain('U/Y')
      expect(np).toEqual(expect.arrayContaining(['LayerType', 'Layer', 'Of', 'LayerPart']))
    }
  })

  it('bagian uang: admin bila IsNewPolicyNonProp != 1 && IsNewPolicyListFormat != 1; atasan bila IsNewPolicyNonProp != 1', () => {
    expect(tampil(MEDAN_ADMIN_UANG, {})).toContain('Premi Ogp')
    expect(tampil(MEDAN_ADMIN_UANG, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
    expect(tampil(MEDAN_ADMIN_UANG, { [P + 'IsNewPolicyListFormat']: '1' })).toEqual([])
    expect(tampil(MEDAN_ATASAN_UANG, { [P + 'IsNewPolicyListFormat']: '1' })).toContain('Premi Ogp')
    expect(tampil(MEDAN_ATASAN_UANG, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
    expect(tampil(MEDAN_ATASAN_TOTAL, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
  })

  it('atasan: Ceding Company, Marketing Officer, Type Tax NOTBLANK', () => {
    const kosong = tampil(MEDAN_ATASAN_UMUM, {})
    for (const l of ['Ceding Company', 'Marketing Officer', 'Type Tax']) expect(kosong).not.toContain(l)
    const isi = tampil(MEDAN_ATASAN_UMUM, {
      [P + 'CedingCoName']: 'UJI-CEDING',
      [P + 'MarketingOfficer']: 'UJI-MO',
      [P + 'TypeTax']: 'Inclusive',
    })
    for (const l of ['Ceding Company', 'Marketing Officer', 'Type Tax']) expect(isi).toContain(l)
  })

  it('atasan: No Polis dan Production Date di wadah "pyWorkPage.FlagViewPolicy = 1" (nol pengisi di korpus)', () => {
    expect(tampil(MEDAN_ATASAN_UMUM, {})).not.toContain('No Polis')
    expect(tampil(MEDAN_ATASAN_UMUM, {})).not.toContain('Production Date')
    expect(tampil(MEDAN_ATASAN_UMUM, { FlagViewPolicy: '1' })).toEqual(expect.arrayContaining(['No Polis', 'Production Date']))
  })

  it('elemen mati (1=2 / NEVER) tidak dibangun: Class Of Business dan Due To tidak ada di layar admin', () => {
    const semua = label(MEDAN_ADMIN_UMUM)
    expect(semua).not.toContain('Class Of Business')
    expect(semua).not.toContain('Due To Us / You')
  })
})

describe('medan layar NB Treaty In - dapat diisi', () => {
  it('AC 52: atasan tidak mengisi satu pun medan General / uang (DueTo, FlagPPH, No Offer Slip pyDisabled)', () => {
    const terbuka = [...MEDAN_ATASAN_UMUM, ...MEDAN_ATASAN_UANG, ...MEDAN_ATASAN_TOTAL].filter(
      (m) => m.jenis !== 'tampil' && !m.kunci,
    )
    expect(terbuka.map((m) => m.jalur)).toEqual([])
    expect(cari(MEDAN_ATASAN_UMUM, 'FlagPPH').kunci).toBe(true)
  })

  it('layar atasan: label VERBATIM (tanpa "(%)")', () => {
    expect(label(MEDAN_ATASAN_UANG)).toContain('Deduction In A (OGP)')
    expect(label(MEDAN_ATASAN_UANG)).not.toContain('(%) Deduction In A (OGP)')
  })
})

describe('medan layar NB Treaty In - aksi sel (action set XML, berurutan)', () => {
  it('persen/hasil OGP-ONP admin: Count*_Act(Data) LALU CountOGPONP_Act', () => {
    const harap: [string, string, string][] = [
      ['RiCommOgp', 'CountResult1', 'Pct'],
      ['ResultOgp1', 'CountResult1', 'Amount'],
      ['OveriddingCommOgp', 'CountResult2Ogp', 'Pct'],
      ['ResultOgp2', 'CountResult2Ogp', 'Amount'],
      ['RiCommOnp', 'CountResult1Onp', 'Pct'],
      ['ResultOnp1', 'CountResult1Onp', 'Amount'],
      ['OveriddingCommOnp', 'CountResult2Onp', 'Pct'],
      ['ResultOnp2', 'CountResult2Onp', 'Amount'],
    ]
    for (const [m, aksi, param] of harap) {
      expect(cari(MEDAN_ADMIN_UANG, m).aksi).toEqual([{ aksi, param }, { aksi: 'CountOGPONP' }])
    }
  })

  it('Gross* -> CalculatePremi_Act(Action); Premi/Claim/Salvage/Excess/Deduction -> CountOGPONP_Act; Outstanding tanpa aksi', () => {
    expect(cari(MEDAN_ADMIN_UANG, 'GrossPremium').aksi).toEqual([{ aksi: 'CalculatePremi', param: 'PREMIUM' }])
    expect(cari(MEDAN_ADMIN_UANG, 'GrossClaim').aksi).toEqual([{ aksi: 'CalculatePremi', param: 'CLAIM' }])
    for (const m of ['PremiOgp', 'PremiOnp', 'Claim', 'SalvageValue', 'ExcessLoss', 'Deduction1', 'Deduction2']) {
      expect(cari(MEDAN_ADMIN_UANG, m).aksi).toEqual([{ aksi: 'CountOGPONP' }])
    }
    expect(cari(MEDAN_ADMIN_UANG, 'OutstandingClaim').aksi).toBeUndefined()
  })

  // Audit silang putaran 3: action set sel `Section/DetailPolicyTreatyIn` -
  // `.EndDate` change -> refresh ProtectDate; `.FlagPPH` change -> runActivity
  // RemoveTypeTax_ACT; `.QuotationData.MOID` change -> refresh CheckDataMkt;
  // `.IDCurrency` change -> refresh SetCurrency_act(CURR=.IDCurrency).
  it('EndDate, FlagPPH, Marketing Officer, Currency admin: satu aksi XML masing-masing', () => {
    expect(cari(MEDAN_ADMIN_UMUM, 'EndDate').aksi).toEqual([{ aksi: 'ProtectDate' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'FlagPPH').aksi).toEqual([{ aksi: 'RemoveTypeTax' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'QuotationData.MOID').aksi).toEqual([{ aksi: 'CheckDataMkt' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'IDCurrency').aksi).toEqual([{ aksi: 'SetCurrency' }])
  })
})

describe('medan layar NB Treaty In - penyajian (AC 85, K3, K14)', () => {
  const PERSEN = ['RiCommOgp', 'OveriddingCommOgp', 'RiCommOnp', 'OveriddingCommOnp']

  it('AC 85 + K3: setiap angka uang membawa kode mata uang - termasuk Deduction1/2 (pxCurrency)', () => {
    for (const ms of [MEDAN_ADMIN_UANG, MEDAN_ATASAN_UANG]) {
      for (const m of ms) {
        const persen = PERSEN.some((x) => m.jalur === P + x)
        expect(m.mataUang === undefined, `${m.jalur} mataUang`).toBe(persen)
      }
    }
    expect(cari(MEDAN_ADMIN_UANG, 'Deduction1').mataUang).toBe(P + 'Currency')
    expect(cari(MEDAN_ATASAN_UANG, 'Deduction2').mataUang).toBe(P + 'Currency')
  })

  // RALAT 06-10-2026 (perintah work owner): bagian uang seragam - nol tampil "0", terisi 4 desimal.
  it('bagian uang: 4 desimal dan nol polos di admin maupun atasan; total spreading tetap', () => {
    const S = { desimal: 4, nolPolos: true, formatSaatSunting: true }
    for (const m of ['GrossPremium', 'RiCommOgp', 'PremiOgp', 'ResultOnp2', 'Claim', 'Deduction1']) {
      expect(cari(MEDAN_ADMIN_UANG, m).sajian).toEqual(S)
    }
    expect(cari(MEDAN_ADMIN_UANG, 'NetPremium').sajian).toEqual({ desimal: 4, nolPolos: true })
    expect(cari(MEDAN_ADMIN_UANG, 'PPNValue').sajian).toEqual({ desimal: 4, nolPolos: true })
    expect(cari(MEDAN_ATASAN_UANG, 'GrossPremium').sajian).toEqual({ desimal: 4, nolPolos: true })
    expect(cari(MEDAN_ATASAN_UANG, 'RiCommOnp').sajian).toEqual({ desimal: 4, nolPolos: true })
    expect(cari(MEDAN_ATASAN_TOTAL, 'TotalPremium').sajian).toEqual({ desimal: 2 })
  })

  it('atasan Quartal / YearOfQuartal: pyFormatType number, 0 desimal, tanpa pemisah ribuan', () => {
    expect(cari(MEDAN_ATASAN_UMUM, 'Quartal').sajian).toEqual({ desimal: 0, ribuan: false })
    expect(cari(MEDAN_ATASAN_UMUM, 'YearOfQuartal').sajian).toEqual({ desimal: 0, ribuan: false })
    expect(cari(MEDAN_ADMIN_UMUM, 'YearOfQuartal').sajian).toBeUndefined()
  })
})

// W6 audit silang P3: wadah bagian uang tanpa judul buatan. Satu-satunya judul di
// dalamnya LABEL sel ber-`pyIncludeLabel=true` (Heading 4) "OGP" / "ONP" di atas
// wadah S24/S25 (admin) dan S93/S94 (atasan); pengelompokan sel = wadah XML.
describe('bagian uang: kelompok wadah XML dan LABEL OGP / ONP', () => {
  const jalur = (ms: Medan[]) => ms.map((m) => m.jalur.replace(P, ''))
  it('admin: S20 Gross, S24 OGP, S25 ONP, S26-S28 klaim/saldo/potongan', () => {
    expect(KELOMPOK_UANG_ADMIN.map((k) => k.judul ?? '')).toEqual(['', 'OGP', 'ONP', ''])
    expect(jalur((KELOMPOK_UANG_ADMIN[1]?.medan ?? []))).toEqual(['PremiOgp', 'RiCommOgp', 'ResultOgp1', 'OveriddingCommOgp', 'ResultOgp2'])
    expect(jalur((KELOMPOK_UANG_ADMIN[2]?.medan ?? []))).toEqual(['PremiOnp', 'RiCommOnp', 'ResultOnp1', 'OveriddingCommOnp', 'ResultOnp2'])
    expect(KELOMPOK_UANG_ADMIN.flatMap((k) => k.medan)).toEqual(MEDAN_ADMIN_UANG)
  })
  it('atasan: S91 Gross, S93 OGP (sampai saldo), S94 ONP (sampai PPN)', () => {
    expect(KELOMPOK_UANG_ATASAN.map((k) => k.judul ?? '')).toEqual(['', 'OGP', 'ONP'])
    expect(jalur((KELOMPOK_UANG_ATASAN[1]?.medan ?? [])).slice(-5)).toEqual(['NetPremium', 'BalanceDueTo', 'BalanceBeforeTax', 'BalanceBeforePPH', 'BalanceDueTo'])
    expect(jalur((KELOMPOK_UANG_ATASAN[2]?.medan ?? []))[0]).toBe('PremiOnp')
    expect(jalur((KELOMPOK_UANG_ATASAN[2]?.medan ?? [])).slice(-2)).toEqual(['PPHValue', 'PPNValue'])
    expect(KELOMPOK_UANG_ATASAN.flatMap((k) => k.medan)).toEqual(MEDAN_ATASAN_UANG)
  })
})

// Sel `.ClaimType` / `.ClaimPaymentType` DetailPolicyTreatyIn: pxDropdown, pyListSource `associated`,
// pyHasNoSelection true; DetailDeptHeadTreatyIn_UW: dropdown Read-only. Isi daftar = prompt values property
// (screenshot work owner 05-10-2026): nilai standar disimpan, teks prompt ditampilkan.
describe('dropdown Claim Type dan Payment Type (prompt values property)', () => {
  const nilaiOpsi = (m: Medan) => (m.opsi ?? []).map((o) => [o.value, o.label])
  it('admin: dropdown dengan nilai standar dan teks prompt', () => {
    const tipe = cari(MEDAN_ADMIN_UMUM, 'ClaimType')
    expect(tipe.jenis).toBe('pilihan')
    expect(nilaiOpsi(tipe)).toEqual([
      ['SOA', 'SOA'],
      ['CashLoss', 'Cash Loss'],
      ['XOL', 'XOL'],
      ['XOL Retro', 'XOL Retro'],
    ])
    const bayar = cari(MEDAN_ADMIN_UMUM, 'ClaimPaymentType')
    expect(bayar.jenis).toBe('pilihan')
    expect(nilaiOpsi(bayar)).toEqual([
      ['Claim', 'Claim'],
      ['AdjusterFee', 'Adjuster Fee'],
      ['Salvage', 'Salvage'],
      ['Adjustment', 'Adjustment'],
      ['Retro', 'XOL Retro'],
    ])
  })

  it('atasan: tetap hanya-baca, menampilkan teks prompt', () => {
    for (const j of ['ClaimType', 'ClaimPaymentType']) {
      const m = cari(MEDAN_ATASAN_UMUM, j)
      expect(m.jenis).toBe('tampil')
      expect(m.opsi).toBe(cari(MEDAN_ADMIN_UMUM, j).opsi)
    }
    const bayar = cari(MEDAN_ATASAN_UMUM, 'ClaimPaymentType').opsi ?? []
    expect(teksPilihan(bayar, 'Retro')).toBe('XOL Retro')
    expect(teksPilihan(bayar, 'AdjusterFee')).toBe('Adjuster Fee')
    // nilai warisan di luar daftar tampil apa adanya, tidak dibuang
    expect(teksPilihan(bayar, 'Lama')).toBe('Lama')
    expect(teksPilihan(bayar, '')).toBe('')
  })
})

// Sel `.QuotationData.IsSurveyReport` DetailPolicyTreatyIn: pxRadioButtons, wajib bila ProportionalType !=
// NonProportional; DetailDeptHeadTreatyIn_UW: radio Read-only. Prompt values property (screenshot work owner
// 06-10-2026): Yes / No.
describe('radio Survey Report (prompt values property)', () => {
  const pasangan = (m: Medan) => (m.opsi ?? []).map((o) => [o.value, o.label])
  it('admin: radio Yes / No, tetap di wadah bukan NonProportional', () => {
    const m = cari(MEDAN_ADMIN_UMUM, 'QuotationData.IsSurveyReport')
    expect(m.jenis).toBe('radio')
    expect(pasangan(m)).toEqual([
      ['Yes', 'Yes'],
      ['No', 'No'],
    ])
    expect(tampil(MEDAN_ADMIN_UMUM, { 'Quotation.ProportionalType': 'NonProportional' })).not.toContain('Survey Report')
  })

  it('atasan: hanya-baca dengan daftar yang sama', () => {
    const m = cari(MEDAN_ATASAN_UMUM, 'QuotationData.IsSurveyReport')
    expect(m.jenis).toBe('tampil')
    expect(m.opsi).toBe(cari(MEDAN_ADMIN_UMUM, 'QuotationData.IsSurveyReport').opsi)
  })
})

// Sel `.StatementType`: DetailPolicyTreatyIn pxDropdown pyListSource `associated` (Editable); DetailDeptHeadTreatyIn_UW
// dropdown Read-only. Prompt values property (screenshot work owner 06-10-2026).
describe('dropdown Statement Type (prompt values property)', () => {
  it('admin: dropdown dengan nilai standar dan teks prompt', () => {
    const m = cari(MEDAN_ADMIN_UMUM, 'StatementType')
    expect(m.jenis).toBe('pilihan')
    expect((m.opsi ?? []).map((o) => [o.value, o.label])).toEqual([
      ['SOA', 'Statement of Account'],
      ['LPC', 'Loss Participation Clause'],
      ['PC', 'Profit Commission'],
      ['SC', 'Sliding Scale'],
    ])
  })

  it('atasan: hanya-baca, menampilkan teks prompt', () => {
    const m = cari(MEDAN_ATASAN_UMUM, 'StatementType')
    expect(m.jenis).toBe('tampil')
    expect(m.opsi).toBe(cari(MEDAN_ADMIN_UMUM, 'StatementType').opsi)
    expect(teksPilihan(m.opsi ?? [], 'LPC')).toBe('Loss Participation Clause')
  })
})

// Sel `.TypeTax` DetailPolicyTreatyIn: pxRadioButtons pyListSource `associated`, wajib bila `.FlagPPH = true`;
// DetailDeptHeadTreatyIn_UW: pxDisplayText (teks nilai apa adanya - tidak diubah). Prompt values property
// (screenshot work owner 06-10-2026): Inclusive / Exclusive.
describe('radio Type Tax (prompt values property)', () => {
  it('admin: radio Inclusive / Exclusive, tampil hanya bila With Tax dicentang', () => {
    const m = cari(MEDAN_ADMIN_UMUM, 'TypeTax')
    expect(m.jenis).toBe('radio')
    expect((m.opsi ?? []).map((o) => [o.value, o.label])).toEqual([
      ['Inclusive', 'Inclusive'],
      ['Exclusive', 'Exclusive'],
    ])
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'false' })).not.toContain('Type Tax')
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'true' })).toContain('Type Tax')
  })

  it('admin: memilih Type Tax menghitung ulang pajak (keputusan work owner 06-10-2026; XML hanya postValue)', () => {
    expect(cari(MEDAN_ADMIN_UMUM, 'TypeTax').aksi).toEqual([{ aksi: 'HitungPajak' }])
  })

  it('atasan: tetap teks tampil (pxDisplayText)', () => {
    const m = cari(MEDAN_ATASAN_UMUM, 'TypeTax')
    expect(m.jenis).toBe('tampil')
    expect(m.opsi).toBeUndefined()
  })
})
