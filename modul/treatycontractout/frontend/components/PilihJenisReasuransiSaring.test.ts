// Uji pemilih jenis reasuransi yang dapat difilter — Treaty Limit (`pxAutoComplete`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilJenisReasuransiAnakTreatyLimit, PILIHAN_REINS_ANAK_TREATY_LIMIT } from '../api'
import { opsiSaringJenisReasuransi } from './PilihJenisReasuransiSaring'

const SUMBER = readFileSync(join(__dirname, 'PilihJenisReasuransiSaring.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

const DAFTAR = [
  { id: '10004', note: 'UJI QS (R/I)', tipe: '4' },
  { id: '10028', note: 'UJI QS (OR)', tipe: '4' },
  { id: '10248', note: 'UJI SPL (OR)', tipe: '4' },
]

describe('opsi yang disaring ketikan — pyUseForSearch hanya pada nama', () => {
  it('kata kosong = seluruh daftar, urutan server', () => {
    expect(opsiSaringJenisReasuransi(DAFTAR, '', false).map((o) => o.value)).toEqual(['10004', '10028', '10248'])
  })
  it('menyaring pada Note, tanpa peduli huruf besar-kecil', () => {
    expect(opsiSaringJenisReasuransi(DAFTAR, 'qs (', false).map((o) => o.value)).toEqual(['10004', '10028'])
    expect(opsiSaringJenisReasuransi(DAFTAR, ' (or) ', false).map((o) => o.value)).toEqual(['10028', '10248'])
  })
  it('ID TIDAK dicari (pyUseForSearch false pada .ID / .CARI1)', () => {
    expect(opsiSaringJenisReasuransi(DAFTAR, '10028', true)).toEqual([])
  })
  it('induk: hanya Note tampil; anak: ID ikut tampil (CARI1 pyShow true)', () => {
    expect(opsiSaringJenisReasuransi(DAFTAR, 'SPL', false)).toEqual([{ value: '10248', label: 'UJI SPL (OR)' }])
    expect(opsiSaringJenisReasuransi(DAFTAR, 'SPL', true)).toEqual([
      { value: '10248', label: 'UJI SPL (OR)', keterangan: '10248' },
    ])
  })
})

describe('sumber daftar', () => {
  it('baris anak: daftar dari NAMA induk (TreatyContractSetReinsTypeList); lainnya daftar induk tiket 02', () => {
    expect(SUMBER).toContain('ambilJenisReasuransiAnakTreatyLimit(namaIndukAnak)')
    expect(SUMBER).toContain('ambilJenisReasuransiTreaty()')
    // Ganti induk = daftar dimuat ulang.
    expect(SUMBER).toContain('}, [namaIndukAnak])')
  })
  it('memakai PilihSaring; nol literal ID porsi dan nol Number()', () => {
    expect(SUMBER).toContain('<PilihSaring')
    for (const id of ['10004', '10011', '10012', '10021', '10022', '10025', '10026', '10028', '10248', '10249', '10018', '10217']) {
      expect(SUMBER).not.toContain(id)
    }
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseInt|parseFloat/)
  })
  it('penanda sumber pilihan = konstanta Go models.PilihanReinsAnakTreatyLimit', () => {
    const go = readFileSync(join(__dirname, '..', '..', 'backend', 'models', 'tco_klausul.go'), 'utf8')
    expect(go).toContain(`const PilihanReinsAnakTreatyLimit = "${PILIHAN_REINS_ANAK_TREATY_LIMIT}"`)
  })
})

describe('GET jenis-reasuransi/anak-treaty-limit', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })
  it('mengirim NAMA induk ke jalur yang didaftarkan handler', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response('{"daftar":[],"total":0}', { status: 200 }))))
    await ambilJenisReasuransiAnakTreatyLimit('2019 QS 101M TRT')
    const url = String(vi.mocked(fetch).mock.calls[0]?.[0])
    expect(url).toContain('/api/treaty-contract-out/jenis-reasuransi/anak-treaty-limit')
    expect(url).toContain('namaInduk=2019')
    const rute = readFileSync(join(__dirname, '..', '..', 'backend', 'handlers', 'rute_treaty_contract_out.go'), 'utf8')
    expect(rute).toContain('"GET /api/treaty-contract-out/jenis-reasuransi/anak-treaty-limit"')
  })
})
