// Uji layar tahun treaty — tiket 03 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { TAHUN_TCO } from '../../assets/labels.treaty-contract-out'
import type { TahunTreaty } from '../../services/api'
import { formDari, formKosong, keMasuk, namaGrup, selTahun } from './InboxTreatyContract'

const SUMBER = readFileSync(join(__dirname, 'InboxTreatyContract.tsx'), 'utf8')
const KODE = SUMBER.split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS: TahunTreaty = {
  id: '1000001', treatyYear: '2026', underwritingYear: '2026', treatyGroupId: '10001',
  treatyGroupName: 'UJI GRUP', proportion: '10003', startDate: '2026-01-01', endDate: '2026-12-31',
  userId: 'UJI-OP', tglUpdate: '2026-09-28 10:00:00',
}

describe('form tahun treaty', () => {
  it('baris baru: ID kosong — identitas dibuat server (AC 5)', () => {
    const f = formKosong()
    expect(f.id).toBe('')
    expect(keMasuk(f).id).toBe('')
  })
  it('Edit menyalin seluruh medan baris (SetTreatyYear_Act b1110–b1353)', () => {
    const f = formDari(BARIS)
    expect(f).toMatchObject({ id: '1000001', treatyGroupName: 'UJI GRUP', proportion: '10003', startDate: '2026-01-01' })
  })
  it('tanggal dikirim YYYY-MM-DD dari bentuk kabel DD-MM-YYYY maupun ISO', () => {
    const f = { ...formDari(BARIS), startDate: '01-01-2026', endDate: '2026-12-31' }
    const m = keMasuk(f)
    expect(m.startDate).toBe('2026-01-01')
    expect(m.endDate).toBe('2026-12-31')
  })
  it('teks yang bukan tanggal dikirim APA ADANYA — server yang menolak dengan pesan', () => {
    expect(keMasuk({ ...formKosong(), startDate: 'kapan-kapan' }).startDate).toBe('kapan-kapan')
    expect(keMasuk({ ...formKosong(), startDate: '  ' }).startDate).toBe('')
  })
  it('tahun dikirim sebagai TEKS, di-trim, tidak lewat Number()', () => {
    expect(keMasuk({ ...formKosong(), treatyYear: ' 2026 ' }).treatyYear).toBe('2026')
    expect(KODE).not.toMatch(/Number\(|parseInt|parseFloat/)
  })
  it('sel kosong ditandai', () => {
    expect(selTahun('')).toBe('—')
    expect(selTahun('2026')).toBe('2026')
  })
  it('nama grup dari master, kosong bila tidak dikenal', () => {
    const grup = [{ id: '10001', treatyGroupName: 'UJI GRUP' }]
    expect(namaGrup(grup, '10001')).toBe('UJI GRUP')
    expect(namaGrup(grup, '99')).toBe('')
  })
})

describe('paritas layar', () => {
  it('enam kepala kolom grid dan empat tombol dari label, tidak diketik ulang', () => {
    for (const k of ['kolomUnderwritingYear', 'kolomTransactionYear', 'kolomStartDate', 'kolomEndDate',
      'kolomTreatyGroup', 'kolomReinsuranceType', 'add', 'edit', 'reinsType', 'listDescription',
      'formId', 'formTreatyGroup', 'formReinsuranceType', 'formStartDate', 'formEndDate',
      'formUnderwritingYear', 'formTransactionYear', 'formModifiedDate', 'formUsername', 'save', 'cancel'] as const) {
      expect(KODE, k).toContain(`TAHUN_TCO.${k}`)
      expect(KODE).not.toContain(`'${TAHUN_TCO[k]}'`)
    }
  })
  it('kolom Reinsurance Type grid menampilkan proportion (b19724)', () => {
    expect(KODE).toContain('selTahun(b.proportion)')
  })
  it('nol tombol Copy dan nol form salin (AC 72)', () => {
    expect(KODE).not.toMatch(/Copy|salin|BrowseCopyData/)
  })
  it('ReinsType dan List Description berdiri menunggu tiketnya, bukan disembunyikan', () => {
    expect(KODE).toContain('TAHUN_TCO.reinsType')
    expect(KODE).toContain('TAHUN_TCO.listDescription')
    expect(KODE).toContain('menungguTiket')
  })
  it('ID, Modified Date, Username hanya dibaca', () => {
    expect((KODE.match(/readOnly/g) ?? []).length).toBe(3)
  })
  it('layar tidak menyaring daftar dan memakai pemilih jenis reasuransi tiket 02', () => {
    expect(KODE).not.toContain('.filter(')
    expect(KODE).toContain("from '../../components/treaty-contract-out/PilihJenisReasuransi'")
  })
  it('tiket 12: panel lampiran hanya untuk tahun yang sudah ber-ID, bukan syarat simpan', () => {
    expect(KODE).toContain('<PanelLampiranTahun tahunID={form.id} />')
    expect(KODE).toMatch(/form\.id !== '' \?/)
    expect(KODE).toContain('LAMPIRAN_TCO.simpanDulu')
    // Simpan tahun tidak menunggu lampiran: fungsi simpan tidak menyebut lampiran.
    const awal = KODE.indexOf('async function simpan')
    const simpan = KODE.slice(awal, KODE.indexOf('\n  }\n', awal))
    expect(simpan).toContain('simpanTahunTreaty(')
    expect(simpan).not.toMatch(/[Ll]ampiran/)
  })
})
