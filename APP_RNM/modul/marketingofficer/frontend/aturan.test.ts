import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, ambilPilihan, tambah, ubah, type BarisMO, type Pilihan } from './api'
import {
  isianDari,
  isianKosong,
  leaderUntuk,
  periksa,
  saring,
  setelBranch,
  setelSubBranch,
  subBranchUntuk,
  tandaAkun,
} from './aturan'
import { MO } from './labels'

const baris = (sebagian: Partial<BarisMO>): BarisMO => ({
  id: '10000101', clientId: 'CON-1001', clientName: 'UJI Leader Satu', clientId2: 'LEADER', moLeader: 'UJI Leader Satu',
  moStatus: '1', branchParent: 'UJI-P00', branchDetailId: 'UJI-B01', branchDetailName: 'UJI Cabang Satu', teamGroup: '1',
  branchStatus: '', aksesLogin: 'UJI-LEAD01', userUpdate: 'UJI-ADMIN', tanggal: '2026-10-03 09:00', statusAkun: 'aktif',
  emailAkun: '', ...sebagian,
})

const pilihan: Pilihan = {
  akun: [],
  leader: [{ id: '10000101', nama: 'UJI Leader Satu', subBranch: '' }, { id: '10000105', nama: 'UJI Leader Dua', subBranch: '' }],
  branch: [{ id: 'UJI-P00', nama: 'UJI Pusat', induk: 'UJI-P00', kanwilGroup: '' }],
  subBranch: [
    { id: 'UJI-B01', nama: 'UJI Cabang Satu', induk: 'UJI-P00', kanwilGroup: '1' },
    { id: 'UJI-X01', nama: 'UJI Cabang Lain', induk: 'UJI-PX', kanwilGroup: '3' },
  ],
}

describe('aturan layar Marketing Officer', () => {
  it('isian dari baris: leader dikenali dari CLIENTID2, anggota membawa ID leadernya', () => {
    expect(isianDari(baris({}))).toEqual({ aksesLogin: 'UJI-LEAD01', leader: true, leaderId: '', branchParent: 'UJI-P00', branchDetailId: 'UJI-B01', aktif: true })
    expect(isianDari(baris({ clientId2: '10000101', moStatus: '2' }))).toMatchObject({ leader: false, leaderId: '10000101', aktif: false })
  })

  it('periksa: akun wajib saat tambah; leader wajib kecuali Set as a leader', () => {
    expect(periksa(isianKosong(), true)).toBe(MO.galatAkun)
    expect(periksa({ ...isianKosong(), aksesLogin: 'UJI-MKT01' }, true)).toBe(MO.galatLeader)
    expect(periksa({ ...isianKosong(), aksesLogin: 'UJI-MKT01', leader: true }, true)).toBeNull()
    // Baris lama Pega tanpa akun boleh diubah tanpa akun.
    expect(periksa({ ...isianKosong(), leaderId: '10000101' }, false)).toBeNull()
  })

  it('Sub Branch disaring Branch; ganti Branch mengosongkan Sub Branch yang bukan miliknya; Sub Branch menyetel Branch', () => {
    expect(subBranchUntuk(pilihan, 'UJI-P00').map((c) => c.id)).toEqual(['UJI-B01'])
    expect(subBranchUntuk(pilihan, '')).toHaveLength(2)
    const isi = { ...isianKosong(), branchParent: 'UJI-P00', branchDetailId: 'UJI-B01' }
    expect(setelBranch(isi, pilihan, 'UJI-PX')).toMatchObject({ branchParent: 'UJI-PX', branchDetailId: '' })
    expect(setelBranch(isi, pilihan, 'UJI-P00').branchDetailId).toBe('UJI-B01')
    expect(setelSubBranch(isianKosong(), pilihan, 'UJI-X01')).toMatchObject({ branchParent: 'UJI-PX', branchDetailId: 'UJI-X01' })
  })

  it('pilihan Leader tanpa baris itu sendiri', () => {
    expect(leaderUntuk(pilihan, '10000101').map((l) => l.id)).toEqual(['10000105'])
    expect(leaderUntuk(pilihan, null)).toHaveLength(2)
  })

  it('saring: status lalu kata di kolom yang tampil', () => {
    const d = [baris({}), baris({ id: '10000103', clientName: 'UJI Pega Lama', clientId: 'UJI-KONTAK-3', moStatus: '2', aksesLogin: '' })]
    expect(saring(d, '', 'aktif').map((b) => b.id)).toEqual(['10000101'])
    expect(saring(d, '', 'nonaktif').map((b) => b.id)).toEqual(['10000103'])
    expect(saring(d, 'kontak-3', 'semua').map((b) => b.id)).toEqual(['10000103'])
    expect(saring(d, 'uji-lead01 satu', 'semua').map((b) => b.id)).toEqual(['10000101'])
  })

  it('tanda akun: tanpa akun, tidak ada di user list, nonaktif; akun aktif tanpa tanda', () => {
    expect(tandaAkun({ statusAkun: '' })).toBe(MO.tanpaAkun)
    expect(tandaAkun({ statusAkun: 'tidak-ada' })).toBe(MO.akunTidakAda)
    expect(tandaAkun({ statusAkun: 'nonaktif' })).toBe(MO.akunNonaktif)
    expect(tandaAkun({ statusAkun: 'aktif' })).toBe('')
  })
})

describe('klien API Marketing Officer', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('jalur, metode, dan badan setiap panggilan', async () => {
    const tertangkap: { url: string; init: RequestInit }[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        tertangkap.push({ url, init })
        return Promise.resolve(new Response('{}', { status: 200 }))
      }),
    )
    const isi = { ...isianKosong(), aksesLogin: 'UJI-MKT01', leaderId: '10000101' }
    await ambilDaftar()
    await ambilPilihan()
    await tambah(isi)
    await ubah('1/2', isi)
    expect(tertangkap.map((t) => `${t.init.method ?? 'GET'} ${t.url}`)).toEqual([
      'GET /api/marketing-officer',
      'GET /api/marketing-officer/pilihan',
      'POST /api/marketing-officer',
      'PUT /api/marketing-officer/1%2F2',
    ])
    expect(JSON.parse(String(tertangkap[2]?.init.body))).toEqual(isi)
  })
})
