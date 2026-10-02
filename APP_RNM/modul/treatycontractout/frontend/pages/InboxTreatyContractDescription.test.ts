// Uji layar klausul dari menu — tiket 08 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const KODE = readFileSync(join(__dirname, 'InboxTreatyContractDescription.tsx'), 'utf8')

describe('layar klausul dari menu', () => {
  it('memakai panel yang SAMA dengan tombol List Description baris tahun', () => {
    expect(KODE).toContain("import PanelKlausulTahun from '../components/PanelKlausulTahun'")
    expect(KODE).toContain('<PanelKlausulTahun key={terpilih.id} tahun={terpilih} />')
  })
  it('judul = judul tampilan (keputusan work owner 02-10-2026, bukan nama harness); label pilihan tahun sama dengan tiket 04', () => {
    expect(KODE).toContain('JUDUL_TAMPIL_TCO.deskripsi')
    expect(KODE).not.toContain('{MENU_TCO.inboxTreatyContractDescription}')
    expect(KODE).toContain("import { labelTahun } from './InboxTreatyContractReinsType'")
  })
})
