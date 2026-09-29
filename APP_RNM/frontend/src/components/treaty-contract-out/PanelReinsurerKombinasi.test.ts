// Uji panel reinsurer — tiket 05 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { formReinsurerDari, formReinsurerKosong, keMasukReinsurer, labelKombinasi } from './PanelReinsurerKombinasi'

const KODE = readFileSync(join(__dirname, 'PanelReinsurerKombinasi.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

describe('form reinsurer', () => {
  it('baris baru: NewTreatyReinsurerDetail_Act mengosongkan seluruh medan tampil', () => {
    expect(formReinsurerKosong()).toEqual({
      id: '', reinsurerId: '', name: '', pctShare: '', ricomm: '', stdRating: '', operatorName: '',
    })
  })
  it('Edit: SetUbahTreatyReinsurerList_Act menyalin ID, reinsurer, share, komisi, rating', () => {
    const f = formReinsurerDari({
      id: '1000005', treatyYear: '2026', treatyGroupId: '10001', treatyGroupName: 'UJI GRUP', reinsTypeId: '10003',
      reinsTypeName: 'UJI QS', reinsurerId: 'UJI-R1', clientId: 'UJI-C1', name: 'UJI REAS', ricomm: '12.5',
      pctShare: '33.33333333', iuDate: '', userId: 'UJI-ADMIN', startDate: '', endDate: '', statusOn: '',
      stdRating: 'A', operatorName: 'UJI-ADMIN', tglUpdate: '',
    })
    expect(f).toEqual({
      id: '1000005', reinsurerId: 'UJI-R1', name: 'UJI REAS', pctShare: '33.33333333', ricomm: '12.5',
      stdRating: 'A', operatorName: 'UJI-ADMIN',
    })
  })
  it('share dan komisi dikirim sebagai TEKS apa adanya — server yang menormalkan koma', () => {
    const m = keMasukReinsurer({ ...formReinsurerKosong(), reinsurerId: 'UJI-R1', pctShare: ' 33,33333333 ',
      ricomm: '12,5', stdRating: ' A ' })
    expect(m).toEqual({ id: '', reinsurerId: 'UJI-R1', pctShare: '33,33333333', ricomm: '12,5', stdRating: 'A' })
  })
  it('kepala kombinasi menyebut tahun, grup, jenis', () => {
    expect(labelKombinasi({ treatyYear: '2026', treatyGroupId: '10001', treatyGroupName: 'UJI GRUP',
      reinsTypeId: '10003', reinsTypeName: 'UJI QS' })).toBe('2026 · UJI GRUP · UJI QS')
  })
})

describe('paritas dan uang', () => {
  it('seluruh label dari REINSURER_TCO', () => {
    for (const k of ['add', 'kolomReinsId', 'kolomReinsurer', 'kolomShare', 'kolomComm', 'kolomRating',
      'kolomOperatorName', 'edit', 'delete', 'securityReinsurer', 'totalShare', 'formId', 'formReinsId',
      'formReinsurer', 'formShare', 'formComm', 'formRating', 'formOperatorName', 'save', 'informasi'] as const) {
      expect(KODE).toContain(`REINSURER_TCO.${k}`)
    }
  })
  it('uang TIDAK PERNAH menjadi angka JavaScript', () => {
    expect(KODE).not.toMatch(/Number\(|parseFloat|parseInt|toFixed|\+\s*l\.pctShare/)
    expect(KODE).toContain('daftar.totalShare')
  })
  it('reinsurer dipilih dari master, bukan diketik (Reins.ID hanya dibaca)', () => {
    expect(KODE).toContain('cariReinsurerMaster(')
    expect(KODE).toMatch(/label=\{REINSURER_TCO\.formReinsId\}[^/]*readOnly/)
  })
  it('Delete dan Security Reinsurer berdiri menunggu tiketnya', () => {
    expect(KODE).toContain('`${TAHUN_TCO.menungguTiket} 10`')
    expect(KODE).toContain('`${TAHUN_TCO.menungguTiket} 06`')
  })
})
