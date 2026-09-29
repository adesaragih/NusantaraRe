// Uji layar kontrak dari menu — tiket 04 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { labelTahun } from './InboxTreatyContractReinsType'

const KODE = readFileSync(join(__dirname, 'InboxTreatyContractReinsType.tsx'), 'utf8')

describe('layar kontrak dari menu', () => {
  it('label pilihan tahun menyebut tahun, grup, dan jenis; ID bila semuanya kosong', () => {
    const t = {
      id: '1000001', treatyYear: '2026', underwritingYear: '2026', treatyGroupId: '10001',
      treatyGroupName: 'UJI GRUP', proportion: '10003', startDate: '', endDate: '', userId: '', tglUpdate: '',
    }
    expect(labelTahun(t)).toBe('2026 · UJI GRUP · 10003')
    expect(labelTahun({ ...t, treatyYear: '', treatyGroupName: '', proportion: '' })).toBe('1000001')
  })
  it('memakai editor yang SAMA dengan popup dari baris tahun', () => {
    expect(KODE).toContain("import PanelKontrakTahun from '../components/PanelKontrakTahun'")
    expect(KODE).toContain('<PanelKontrakTahun key={terpilih.id} tahun={terpilih} />')
  })
})
