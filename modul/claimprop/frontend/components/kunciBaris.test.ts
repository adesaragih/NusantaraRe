// Key baris popup master unik walau Treaty ID + Treaty Group + Class of Business sama (DEV: 500 baris, 386 key tiga
// kolom) - key kembar membuat baris lama tertinggal saat daftar berganti.

import { describe, expect, it } from 'vitest'

import { kunciMaster, type KolomMaster } from './kunciBaris'

const dasar: KolomMaster = {
  treatyId: 'UJI-M1',
  classOfBusiness: 'UJI-COB',
  classOfBusinessId: 'UJI-COBID',
  treatyContractName: 'UJI-KONTRAK',
  sob: 'UJI-SOB',
  ceding: 'UJI-CEDING',
  treatyType: 'UJI-QS',
  proportionType: 'Proportional',
  treatyGroup: 'UJI-GRUP',
  treatyGroupId: 'UJI-TG',
  treatyYear: '2026',
}

describe('kunci baris popup master', () => {
  it('baris yang hanya berbeda tahun / SOB tetap ber-key berbeda', () => {
    const k = new Set([
      kunciMaster(dasar),
      kunciMaster({ ...dasar, treatyYear: '2025' }),
      kunciMaster({ ...dasar, sob: 'UJI-SOB-2' }),
    ])
    expect(k.size).toBe(3)
  })

  it('pemisah tidak membuat dua baris berbeda bertabrakan', () => {
    expect(kunciMaster({ ...dasar, sob: 'A', ceding: 'BC' })).not.toBe(
      kunciMaster({ ...dasar, sob: 'AB', ceding: 'C' }),
    )
  })
})
