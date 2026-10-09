import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, tambah, ubah } from './api'
import { BATAS_CAUSE_OF_LOSS, jumlahHalaman, periksaCauseOfLoss, UKURAN_HALAMAN } from './aturan'
import { COL, MENU_COL } from './labels'

const HALAMAN = readFileSync(`${__dirname}/pages/CauseOfLossLife.tsx`, 'utf8')

/** Label VERBATIM dari `InboxCauseofLossLife.xml` (nomor baris di labels.ts). */
const XML_LABEL = ['CAUSE OF LOSS', 'ID', 'Cause of Loss', 'Edit', 'Save', 'Cancel']

describe('label Cause Of Loss Life', () => {
  it('nama menu = M_NAV_MENU.LABEL slot menu modul 955 (K5), golongan MASTER TREATY URUTAN 14, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../backend/migrations/955_menu_causeoflosslife.sql`, 'utf8')
    expect(MENU_COL.kelompok).toBe('Cause Of Loss Life')
    expect(sql).toContain(`'causeoflosslife', '${MENU_COL.kelompok}', 'MASTER TREATY', 'causeoflosslife', 14, '1'`)
  })

  it('label layar VERBATIM dari XML', () => {
    const nilai = Object.values(COL).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan Cause Of Loss Life', () => {
  it('grid 10 baris (b4057), halaman minimal 1', () => {
    expect(UKURAN_HALAMAN).toBe(10)
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(10)).toBe(1)
    expect(jumlahHalaman(11)).toBe(2)
  })

  it('Cause of Loss wajib (b996), paling panjang 200 byte (batas kolom), huruf TIDAK diubah', () => {
    expect(periksaCauseOfLoss('   ')).toBe(COL.galatWajib)
    expect(periksaCauseOfLoss('uji sakit')).toBeNull()
    expect(periksaCauseOfLoss('A'.repeat(BATAS_CAUSE_OF_LOSS))).toBeNull()
    expect(periksaCauseOfLoss('A'.repeat(BATAS_CAUSE_OF_LOSS + 1))).toBe(COL.galatPanjang(200))
    expect(periksaCauseOfLoss('€'.repeat(67))).toBe(COL.galatPanjang(200))
    const kode = HALAMAN.split('\n')
      .filter((b) => !b.trim().startsWith('//'))
      .join('\n')
    expect(kode).not.toMatch(/toUpperCase|hurufBesar/)
  })
})

describe('layar InboxCauseofLossLife', () => {
  it('ID tidak dapat diisi (disabled b872); Cause of Loss pxTextInput wajib (b1029)', () => {
    expect(HALAMAN).toMatch(/\{COL\.id\}<\/span>\s*<input[^>]*value=\{ubahID\}[^>]*readOnly disabled/)
    expect(HALAMAN).toMatch(/<input\s+className="field__input"\s+type="text"\s+value=\{nama\}/)
    expect(HALAMAN).not.toContain('<textarea')
  })

  it('Save selalu; Cancel HANYA saat Edit (DATASHOW = IsEdit b5787); Edit mengisi form (EditList_DT b3770)', () => {
    expect(HALAMAN).toContain('{sibuk ? COL.menyimpan : COL.save}')
    expect(HALAMAN).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{COL\.batal\}/)
    expect(HALAMAN).toContain("ubahID === '' ? tambah(isi) : ubah(ubahID, isi)")
    expect(HALAMAN).toMatch(/setUbahID\(c\.id\)\s*setNama\(c\.causeOfLoss\)/)
  })

  it('grid kolom ID, Cause of Loss, lalu kolom Edit; nol saring / urut / Delete / Upload / detail (XML tidak memuatnya)', () => {
    const kepala = HALAMAN.slice(HALAMAN.indexOf('<thead>'), HALAMAN.indexOf('</thead>'))
    expect(kepala.indexOf('{COL.id}')).toBeLessThan(kepala.indexOf('{COL.causeOfLoss}'))
    expect(kepala).not.toContain('onClick')
    // Kode saja - komentar kepala berkas menyebut apa yang SENGAJA tidak dibuat.
    const kode = HALAMAN.split('\n')
      .filter((b) => !b.trim().startsWith('//'))
      .join('\n')
    for (const kata of ['Delete', 'hapus', 'Upload', 'unggah', 'Detail', 'rincian', 'type="search"', 'saring', 'arah']) {
      expect(kode, kata).not.toContain(kata)
    }
  })

  it('View only: form dan Edit hanya bila bolehUbah', () => {
    expect(HALAMAN).toContain('useBolehUbah(NAMA_COL)')
    expect(HALAMAN.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(HALAMAN).toContain('{bolehUbah && <th')
  })
})

describe('klien Cause Of Loss Life', () => {
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
    await tambah('UJI SAKIT')
    await ubah('100002', 'UJI SAKIT 2')
    expect(panggil).toEqual([
      'GET /api/cause-of-loss-life',
      'GET /api/cause-of-loss-life?halaman=2',
      'POST /api/cause-of-loss-life {"causeOfLoss":"UJI SAKIT"}',
      'PUT /api/cause-of-loss-life/100002 {"causeOfLoss":"UJI SAKIT 2"}',
    ])
  })
})
