// Tautan berkas antar-tab (tombol View polis Claim Prop -> berkas NB / EDM Treaty In, 08-10-2026): form GET
// tersembunyi mengirim PARAM_MODUL + PARAM_KASUS; App membacanya dengan bacaTautanKasus.

import { describe, expect, it } from 'vitest'

import { bacaTautanKasus, PARAM_KASUS, PARAM_MODUL } from './tautanKasus'

const kirimForm = (modul: string, kasus: string) =>
  '?' + new URLSearchParams({ [PARAM_MODUL]: modul, [PARAM_KASUS]: kasus }).toString()

describe('tautan berkas antar-tab', () => {
  it('bolak-balik utuh, termasuk ID kunci Pega berspasi dan bertitik', () => {
    expect(bacaTautanKasus(kirimForm('nbtreatyin', 'UJI-KELAS NB-1.2'))).toEqual({
      modul: 'nbtreatyin',
      id: 'UJI-KELAS NB-1.2',
    })
  })

  it('tanpa modul atau tanpa kasus = bukan tautan berkas', () => {
    expect(bacaTautanKasus('')).toBeNull()
    expect(bacaTautanKasus('?modul=nbtreatyin')).toBeNull()
    expect(bacaTautanKasus('?kasus=UJI-1')).toBeNull()
    expect(bacaTautanKasus('?modul=%20&kasus=UJI-1')).toBeNull()
  })
})
