// `Section/DetailPolicyTreatyIn` sel `.StartDate` (Statement Period): aksi refresh
// `change` dengan `pyPreDataTransform` `SystemSetOneYear_DT` - mengubah tanggal
// mulai di layar admin mengisi ulang tanggal akhir (+1 tahun) lewat backend.

import { describe, expect, it } from 'vitest'

import { MEDAN_ADMIN_UMUM } from './medan'

describe('tanggal mulai layar admin', () => {
  it('Statement Period memicu aksi SystemSetOneYear (SystemSetOneYear_DT)', () => {
    const m = MEDAN_ADMIN_UMUM.find((x) => x.jalur === 'PolicyTreatyIn.StartDate')
    // aksi = action set sel, berurutan (medan.ts `Aksi[]`)
    expect(m?.aksi).toEqual([{ aksi: 'SystemSetOneYear' }])
  })
})
