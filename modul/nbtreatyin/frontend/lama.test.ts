import { describe, expect, it } from 'vitest'

import type { DokumenLama } from './api'
import { idBolehDisalin, ringkasSalin, saringLama, tampilCopyOld } from './lama'

const dok = (id: string, boleh: boolean, insured = 'UJI TERTANGGUNG'): DokumenLama => ({
  id,
  noOffer: 'UJI-M1',
  noPolis: 'UJI-POL-1',
  insuredName: insured,
  businessName: 'UJI GRUP',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  tglProd: '',
  bolehDisalin: boleh,
  alasan: boleh ? [] : ['UJI alasan'],
})

describe('Copy Old NB (perintah work owner 07-10-2026)', () => {
  it('tombol hanya bila server menyatakan superadmin', () => {
    expect(tampilCopyOld(null)).toBe(false)
    expect(tampilCopyOld({ copyOld: false })).toBe(false)
    expect(tampilCopyOld({ copyOld: true })).toBe(true)
  })

  it('Search menyaring NB Number / Master ID / Policy / Insured / Group Business / SOB / Ceding', () => {
    const d = [dok('NB-1', true), dok('NB-2', true, 'UJI LAIN')]
    expect(saringLama(d, 'lain').map((x) => x.id)).toEqual(['NB-2'])
    expect(saringLama(d, 'grup')).toHaveLength(2)
    expect(saringLama(d, ' ').map((x) => x.id)).toEqual(['NB-1', 'NB-2'])
  })

  it('Process Copy hanya mengirim yang dicentang DAN boleh disalin, urutan daftar', () => {
    const d = [dok('NB-1', true), dok('NB-2', false), dok('NB-3', true)]
    expect(idBolehDisalin(d, new Set(['NB-3', 'NB-2', 'NB-1']))).toEqual(['NB-1', 'NB-3'])
  })

  it('ringkasan hasil per status', () => {
    expect(
      ringkasSalin([
        { id: 'NB-1', status: 'disalin', pesan: [] },
        { id: 'NB-2', status: 'gagal', pesan: ['UJI'] },
      ]),
    ).toEqual({ disalin: 1, sudahAda: 0, ditolak: 0, gagal: 1 })
  })
})
