// Audit 8 Oktober 2026 — "jangan semisal isi data table kosong tp bisa
// ditambah". Penyesuaian yang isi tabnya TIDAK tersedia (`terdarat = false`)
// dibuka dalam mode LIHAT di panel New: grid kosong tidak dapat ditambah
// lalu disimpan sebagai data separuh. Di Pega penyesuaian selalu membawa
// salinan dokumen utuh, jadi grid tidak pernah kosong karena datanya hilang.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PENYESUAIAN } from './labelsPenyesuaian'
import { MEDAN_KANAN_BARU, MEDAN_KIRI_BARU } from './labelsPenyesuaian'
import { modeEfektif, simpanTampil, terapkanSetEdit } from './pages/PenyesuaianKontrak'

const HALAMAN = readFileSync(join(__dirname, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')

describe('penyesuaian tanpa isi tab — panel New terkunci', () => {
  it('terdarat=false memaksa mode LIHAT untuk panel New dan deret tombol; panel Old tidak berubah', () => {
    expect(HALAMAN).toContain("const modeBaru: ModeLayar = terkunciPendaratan ? '1' : mode")
    expect(HALAMAN).toMatch(/bacaSaja=\{false\}\s*mode=\{modeBaru\}/)
    expect(HALAMAN).toMatch(/<DeretTombol\s*mode=\{modeBaru\}/)
    expect(HALAMAN).toMatch(/bacaSaja\s*mode=\{mode\}/)
  })

  it('mode lihat = Save tidak tampil (syarat tombol Save)', () => {
    expect(simpanTampil('1', {})).toBe(false)
    expect(simpanTampil('0', {})).toBe(true)
  })

  it('pemakai diberi tahu sebabnya — tanpa istilah teknis', () => {
    expect(HALAMAN).toContain('{PENYESUAIAN.belumTerdarat}')
    expect(PENYESUAIAN.belumTerdarat).toContain('Treaty In')
  })
})


// ⛔ KELUHAN 8 OKTOBER 2026 — *"itu seharusnya yang bisa diedit hanya
// Commencement, Termination, Accounting Mode, RNM as Treaty Leader, Effective
// Date, Is Pro Rate, Bordereaux, Rate of Exchange … kenapa setelah saya coba
// tidak bisa edit juga"*.
//
// Penguncinya `RevisionState` TERSIMPAN, dan membacanya begitu adalah salah
// tafsir kami: di Pega properti itu disetel `SetTreatyIn_Act[7]` yang
// berprasyarat `param.revisionstate==1` — parameter TOMBOL, dan tombol Edit
// mengirimnya kosong (`Section/InputTreatyInOffer`).
describe('tombol Edit membuka suntingan — `RevisionState` tersimpan tidak menguncinya', () => {
  it('⛔ Edit tetap mode SUNTING walau dokumennya ber-RevisionState = 1', () => {
    // Nilai ini nyata: tombol Revision layar Treaty In menyimpannya
    // (`revisi_kontrak.go`), jadi tiap penyesuaian turunannya membawanya.
    expect(modeEfektif('0', false, { RevisionState: '1' })).toBe('0')
    expect(modeEfektif('0', true, { RevisionState: '1' })).toBe('0')
    expect(modeEfektif('0', false, {})).toBe('0')
  })

  it('tombol View tetap BACA — yang menentukan tombolnya, bukan dokumennya', () => {
    expect(modeEfektif('1', false, { RevisionState: '1' })).toBe('1')
    expect(modeEfektif('1', false, { RevisionState: '' })).toBe('1')
  })

  it('⭐ `terapkanSetEdit` menyetel ViewState 0 tanpa melihat RevisionState', () => {
    const p = (medan: Record<string, string>) =>
      ({ id: 'X', idAsal: 'X', terdarat: true, baru: { medan, larik: {} }, lama: { medan: {}, larik: {} } }) as never
    expect(terapkanSetEdit(p({ RevisionState: '1' }), 'X').baru.medan.ViewState).toBe('0')
    expect(terapkanSetEdit(p({}), 'X').baru.medan.ViewState).toBe('0')
  })

  // ⚠️ Kedelapan medan kepala yang pemilik proses sebut BOLEH disunting —
  // mereka `bacaSajaBila: modeLihat`, jadi mode '0' membuka SEMUANYA.
  // ⭐ 9 Oktober 2026 — empat medan kiri dikembalikan ke `pyReadOnlyCondition`
  // Pega (keputusan pemakai "Ikuti Pega"); diuji terpisah di bawah.
  const PEGA = ['ContractRefNo', 'TeritorialScope', 'BordereauxNote', 'TreatyContractName']
  it('⛔ kedelapan medan kepala terbuka di mode sunting', () => {
    const boleh = [...MEDAN_KIRI_BARU, ...MEDAN_KANAN_BARU].filter((m) => m.selaluBacaSaja !== true && !PEGA.includes(m.kunci))
    expect(boleh.map((m) => m.kunci).sort()).toEqual(
      ['AccountingMode', 'AccountingModeNonProp', 'Bordeaux', 'Commencement', 'EDMEffective', 'IsProRate', 'Termination', 'TreatyLeader'],
    )
    // Dan tiap satunya terkunci HANYA oleh mode lihat — diuji perilakunya,
    // bukan identitas fungsinya: yang penting hasilnya, bukan rujukannya.
    for (const m of boleh) {
      expect(m.bacaSajaBila?.({ ViewState: '1' })).toBe(true)
      expect(m.bacaSajaBila?.({ ViewState: '0' })).toBe(false)
    }
  })

  it('⭐ empat medan kepala kiri mengikuti pyReadOnlyCondition Pega', () => {
    const cari = (k: string) => MEDAN_KIRI_BARU.find((m) => m.kunci === k)!
    for (const k of ['TreatyContractName', 'ContractRefNo']) {
      expect(cari(k).bacaSajaBila?.({ ViewState: '0', IsEditData: '0', EDMMaterialType: '1' })).toBe(false)
      expect(cari(k).bacaSajaBila?.({ ViewState: '0', IsEditData: '1' })).toBe(true)
      expect(cari(k).bacaSajaBila?.({ ViewState: '1' })).toBe(true)
    }
    for (const k of ['TeritorialScope', 'BordereauxNote']) {
      // Non Material: dapat disunting; Material: terkunci.
      expect(cari(k).bacaSajaBila?.({ ViewState: '0', IsEditData: '0', EDMMaterialType: '2' })).toBe(false)
      expect(cari(k).bacaSajaBila?.({ ViewState: '0', IsEditData: '0', EDMMaterialType: '1' })).toBe(true)
      expect(cari(k).bacaSajaBila?.({ ViewState: '0', IsEditData: '1', EDMMaterialType: '2' })).toBe(true)
    }
  })
})
