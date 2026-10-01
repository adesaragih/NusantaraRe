import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { isiDariDataPolis, type DataPolis } from '../api'
import { LABEL_DATA_POLIS } from '../labels'
import { kolomWajibDataPolis, typeRetro } from './FormDataPolis'

// Uji data polis layar Input Premium Detail — tiket 03 bagian 2.

const BERKAS = readFileSync(join(__dirname, 'FormDataPolis.tsx'), 'utf8')
const DETAIL = readFileSync(join(__dirname, 'PremiumListDetail.tsx'), 'utf8')

const contoh: DataPolis = {
  caseId: 'NBLF-1',
  tahap: 'Input Premium Detail',
  bolehDisimpan: true,
  type: 'QR',
  productNameId: 'UJI-P1',
  productName: 'UJI-PRODUK',
  sourceOfBusiness: 'UJI-S1',
  sobName: 'UJI-SOB',
  cedingCo: 'UJI-L01',
  cedingCoName: 'UJI-CEDING',
  policyHolder: 'UJI-C1',
  policyHolderName: 'UJI-PEMEGANG',
  riSlipRnm: '',
  proRateType: '1',
  moId: 'UJI-M1',
  marketingCode: 'UJI-MC',
  marketingName: 'UJI-MARKETING',
  annuityInterest: '0.05',
  premiumRefundFactor: '1',
  retroId: '',
  retroName: '',
  securityReinsurerId: '',
  securityReinsurer: '',
  wpc: null,
  peringatan: null,
  pilihan: { type: [], proRateType: [] },
}

describe('kolom wajib data polis', () => {
  it('QR lengkap — tidak ada yang kurang', () => {
    expect(kolomWajibDataPolis(isiDariDataPolis(contoh))).toEqual([])
  })

  it('TP/TR menuntut R/I SLIP dan Billing Name', () => {
    const isi = { ...isiDariDataPolis(contoh), type: 'TP' }
    expect(kolomWajibDataPolis(isi)).toEqual([LABEL_DATA_POLIS.riSlip, LABEL_DATA_POLIS.billing])
    expect(typeRetro('TR')).toBe(true)
    expect(typeRetro('QP')).toBe(false)
  })

  it('badan simpan tanpa medan milik server', () => {
    const isi = isiDariDataPolis(contoh)
    for (const k of ['caseId', 'tahap', 'bolehDisimpan', 'wpc', 'pilihan']) {
      expect(Object.keys(isi)).not.toContain(k)
    }
  })
})

describe('struktur layar', () => {
  it('Save Data terkunci selama ada yang kosong, dan dinyatakan tanpa perhitungan', () => {
    expect(BERKAS).toContain('disabled={sibuk || kurang.length > 0}')
    expect(BERKAS).toContain('data.peringatan')
  })

  it('Choose Product Name hanya sebelum bernomor', () => {
    expect(BERKAS).toContain("!bernomor && tombol(LABEL_DATA_POLIS.pilihProduk, 'produk')")
  })

  it('terpasang di halaman Premium List Detail', () => {
    expect(DETAIL).toContain('<FormDataPolis')
  })

  it('label VERBATIM ShowLifePremiumDetail.xml', () => {
    expect(LABEL_DATA_POLIS.proRateType).toBe('Premium Payment Method')
    expect(LABEL_DATA_POLIS.billing).toBe('Billing Name')
    expect(LABEL_DATA_POLIS.simpan).toBe('Save Data')
  })
})
