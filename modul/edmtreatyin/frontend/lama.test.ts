import { describe, expect, it } from 'vitest'

import type { DokumenLama } from './api'
import { idBolehDisalin, ringkasSalin, saringLama, tampilCopyOld } from './lama'

const dok = (id: string, boleh: boolean, noPolis = 'UJI-POL-1'): DokumenLama => ({
  id,
  noPolis,
  edmNo: `${noPolis}/E01`,
  prodKe: 1,
  edmType: '1',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  tglProd: '',
  bolehDisalin: boleh,
  alasan: boleh ? [] : ['UJI alasan'],
})

describe('Copy Old (perintah work owner 07-10-2026)', () => {
  it('tombol hanya bila server menyatakan superadmin', () => {
    expect(tampilCopyOld(null)).toBe(false)
    expect(tampilCopyOld({ copyOld: false })).toBe(false)
    expect(tampilCopyOld({ copyOld: true })).toBe(true)
  })

  it('Search menyaring EDM Number / Policy Number / EDM No / SOB / Ceding tanpa membedakan huruf', () => {
    const d = [dok('EDMT-1', true), dok('EDMT-2', true, 'UJI-POL-2')]
    expect(saringLama(d, 'pol-2').map((x) => x.id)).toEqual(['EDMT-2'])
    expect(saringLama(d, ' ').map((x) => x.id)).toEqual(['EDMT-1', 'EDMT-2'])
    expect(saringLama(d, 'ceding')).toHaveLength(2)
  })

  it('Process Copy hanya mengirim yang dicentang DAN boleh disalin, urutan daftar', () => {
    const d = [dok('EDMT-1', true), dok('EDMT-2', false), dok('EDMT-3', true)]
    expect(idBolehDisalin(d, new Set(['EDMT-3', 'EDMT-2', 'EDMT-1']))).toEqual(['EDMT-1', 'EDMT-3'])
  })

  it('ringkasan hasil per status', () => {
    expect(
      ringkasSalin([
        { id: 'EDMT-1', status: 'disalin', pesan: [] },
        { id: 'EDMT-2', status: 'ditolak', pesan: ['UJI'] },
        { id: 'EDMT-3', status: 'disalin', pesan: [] },
      ]),
    ).toEqual({ disalin: 2, sudahAda: 0, ditolak: 1, gagal: 0 })
  })
})
