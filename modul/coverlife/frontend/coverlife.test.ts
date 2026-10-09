import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, tambah, ubah } from './api'
import { BATAS_COVER, BATAS_NOTE, jumlahHalaman, periksaIsian, UKURAN_HALAMAN } from './aturan'
import { CVL, MENU_CVL } from './labels'

const HALAMAN = readFileSync(`${__dirname}/pages/CoverLife.tsx`, 'utf8')
/** Kode saja - komentar kepala berkas menyebut apa yang SENGAJA tidak dibuat. */
const KODE = HALAMAN.split('\n')
  .filter((b) => !b.trim().startsWith('//'))
  .join('\n')

/** Label VERBATIM dari `InboxCoverLife.xml` (nomor baris di labels.ts). */
const XML_LABEL = ['COVER', 'ID', 'Cover', 'Note', 'Edit', 'Save', 'Cancel']

describe('label Cover Life', () => {
  it('nama menu = M_NAV_MENU.LABEL slot menu modul 957 (C4), golongan MASTER TREATY URUTAN 16, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../backend/migrations/957_menu_coverlife.sql`, 'utf8')
    expect(MENU_CVL.kelompok).toBe('Cover Life')
    expect(sql).toContain(`'coverlife', '${MENU_CVL.kelompok}', 'MASTER TREATY', 'coverlife', 16, '1'`)
  })

  it('label layar VERBATIM dari XML', () => {
    const nilai = Object.values(CVL).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan Cover Life', () => {
  it('grid 50 baris (b4101), halaman minimal 1', () => {
    expect(UKURAN_HALAMAN).toBe(50)
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(50)).toBe(1)
    expect(jumlahHalaman(51)).toBe(2)
  })

  it('Cover wajib (b843), Note opsional (b1020), batas lebar kolom, huruf TIDAK diubah', () => {
    expect(periksaIsian('   ', '')).toBe(CVL.galatWajib)
    expect(periksaIsian('uji jiwa', '')).toBeNull()
    expect(periksaIsian('A'.repeat(BATAS_COVER), 'A'.repeat(BATAS_NOTE))).toBeNull()
    expect(periksaIsian('€'.repeat(67), '')).toBe(CVL.galatPanjang('Cover', 200))
    expect(periksaIsian('UJI', 'A'.repeat(BATAS_NOTE + 1))).toBe(CVL.galatPanjang('Note', 1000))
    expect(KODE).not.toMatch(/toUpperCase|hurufBesar/)
  })
})

describe('layar InboxCoverLife', () => {
  it('form TANPA medan ID; Cover pxTextInput wajib (b880), Note pxTextArea (b1058) tidak wajib', () => {
    expect(KODE).not.toMatch(/value=\{ubahID\}/)
    expect(HALAMAN).toMatch(/<input\s+className="field__input"\s+type="text"\s+value=\{cover\}[^>]*required/)
    expect(HALAMAN).toMatch(/<textarea\s+className="field__input coverlife__teks"\s+value=\{note\}/)
    const note = HALAMAN.slice(HALAMAN.indexOf('<textarea'), HALAMAN.indexOf('/>', HALAMAN.indexOf('<textarea')))
    expect(note).not.toContain('required')
  })

  it('Save selalu; Cancel HANYA saat Edit (DATASHOW = IsEdit b5831); Edit mengisi ID, Cover, Note (EditList_DT b3807)', () => {
    expect(HALAMAN).toContain('{sibuk ? CVL.menyimpan : CVL.save}')
    expect(HALAMAN).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{CVL\.batal\}/)
    expect(HALAMAN).toContain("ubahID === '' ? tambah(isi, note) : ubah(ubahID, isi, note)")
    expect(HALAMAN).toMatch(/setUbahID\(c\.id\)\s*setCover\(c\.cover\)\s*setNote\(c\.note\)/)
  })

  it('grid kolom ID, Cover, lalu kolom Edit - Note TIDAK tampil; nol saring / urut / Delete / Upload / detail', () => {
    const kepala = KODE.slice(KODE.indexOf('<thead>'), KODE.indexOf('</thead>'))
    expect(kepala.indexOf('{CVL.id}')).toBeLessThan(kepala.indexOf('{CVL.cover}'))
    expect(kepala).not.toContain('{CVL.note}')
    expect(kepala).not.toContain('onClick')
    const badan = KODE.slice(KODE.indexOf('<tbody>'), KODE.indexOf('</tbody>'))
    expect(badan).not.toContain('{c.note}')
    for (const kata of ['Delete', 'hapus', 'Upload', 'unggah', 'Detail', 'rincian', 'type="search"', 'saring', 'arah']) {
      expect(KODE, kata).not.toContain(kata)
    }
  })

  it('View only: form dan Edit hanya bila bolehUbah', () => {
    expect(HALAMAN).toContain('useBolehUbah(NAMA_CVL)')
    expect(HALAMAN.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(HALAMAN).toContain('{bolehUbah && <th')
  })
})

describe('klien Cover Life', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go (tanpa DELETE, tanpa saring)', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${String(init.body ?? '')}`.trim())
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(1)
    await ambilDaftar(2)
    await tambah('UJI JIWA', 'uji catatan')
    await ubah('100002', 'UJI JIWA 2', '')
    expect(panggil).toEqual([
      'GET /api/cover-life',
      'GET /api/cover-life?halaman=2',
      'POST /api/cover-life {"cover":"UJI JIWA","note":"uji catatan"}',
      'PUT /api/cover-life/100002 {"cover":"UJI JIWA 2","note":""}',
    ])
  })
})
