// Disalin dari `modul/claimprop/frontend/components/kunciBaris.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
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
  proportionType: 'NonProportional',
  treatyGroup: 'UJI-GRUP',
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
