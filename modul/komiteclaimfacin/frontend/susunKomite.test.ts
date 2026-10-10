// Tata letak layar komite Fac In (pola Komite Claim Prop / Non Prop, keputusan work owner 09-10-2026): satu kartu per
// bagian server berurutan, teks komite kartu rincian terakhir, tangga + keputusan DI BAWAH; grid berincian (expand pane),
// sel tautan View Retro, langkah tangga dari grid "List of Committee" (termasuk perluasan KCF-02 yang belum disimpan).

import { describe, expect, it } from 'vitest'

import type { Anggota, Bagian, Grid, Layar } from './api'
import {
  adaRincian,
  bagianBerisi,
  barisGrid,
  keadaanKasus,
  kunciModal,
  langkahTangga,
  rincianBaris,
  susunBagian,
} from './susunKomite'

const medan = [{ label: 'UJI', nilai: 'UJI', jenis: 'teks' as const }]
const grid = (extra: Partial<Grid> = {}): Grid => ({ judul: '', kolom: [], baris: [], ...extra })
const bagian = (kunci: string, extra: Partial<Bagian> = {}): Bagian => ({ kunci, medan, ...extra })

describe('susunBagian', () => {
  it('urutan server; teks komite kartu terakhir; tangga dipisah ke bawah', () => {
    const s = susunBagian([
      bagian('polis', { judul: 'Policy Detail' }),
      bagian('objek', { judul: 'Object Detail', medan: undefined, grid: [grid()] }),
      bagian('klaim', { judul: 'Claim Details' }),
      bagian('adjustment', { medan: undefined, grid: [grid()] }),
      bagian('riwayat', { judul: 'History Adjustment', medan: undefined, grid: [grid()] }),
      bagian('total', { judul: 'Total Adjustment', medan: undefined, grid: [grid()] }),
      bagian('tangga', { judul: 'List of Committee', medan: undefined, grid: [grid()] }),
      bagian('teksKomite'),
      bagian('baruUJI', { judul: 'UJI Baru' }),
    ])
    expect(s.kartu.map((b) => b.kunci)).toEqual([
      'polis',
      'objek',
      'klaim',
      'adjustment',
      'riwayat',
      'total',
      'baruUJI',
      'teksKomite',
    ])
    expect(s.tangga?.kunci).toBe('tangga')
  })

  it('bagian tanpa medan dan tanpa grid dilewati (Claim Details tanpa adjustment, teks komite kosong)', () => {
    const s = susunBagian([
      bagian('polis', { judul: 'Policy Detail' }),
      { kunci: 'klaim', judul: 'Claim Details' },
      { kunci: 'teksKomite' },
    ])
    expect(s.kartu.map((b) => b.kunci)).toEqual(['polis'])
    expect(bagianBerisi({ kunci: 'x', medan: [], grid: [] })).toBe(false)
  })

  it('TT3 / TT4: tanpa tangga (LS41 hanya TransferType 2), teks komite tetap tampil', () => {
    const s = susunBagian([bagian('teksKomite')])
    expect(s.kartu.map((b) => b.kunci)).toEqual(['teksKomite'])
    expect(s.tangga).toBeUndefined()
  })
})

describe('grid berincian dan tautan', () => {
  const g = grid({
    baris: [{ No: '1' }, { No: '2' }, { No: '3' }],
    rincian: [[bagian('spreadingPolis')], null, []],
  })

  it('baris ber-rincian dapat dibuka; null / kosong = tidak', () => {
    expect(rincianBaris(g, 0)?.map((b) => b.kunci)).toEqual(['spreadingPolis'])
    expect(rincianBaris(g, 1)).toBeNull()
    expect(rincianBaris(g, 2)).toBeNull()
    expect(rincianBaris(g, 9)).toBeNull()
    expect(adaRincian(g)).toBe(true)
    expect(adaRincian(grid({ baris: [{ No: '1' }] }))).toBe(false)
  })

  it('baris null dari server = daftar kosong', () => {
    expect(barisGrid(grid({ baris: null }))).toEqual([])
    expect(adaRincian(grid({ baris: null, rincian: [[bagian('x')]] }))).toBe(false)
  })

  it('tautan View Retro hanya bila kuncinya ada di Layar.modal', () => {
    const modal = { 'retro:1:1:2': [bagian('retro', { judul: 'ShowRetro' })] }
    expect(kunciModal('retro:1:1:2', modal)).toBe('retro:1:1:2')
    expect(kunciModal('', modal)).toBeNull()
    expect(kunciModal(undefined, modal)).toBeNull()
    expect(kunciModal('retro:9:9:9', modal)).toBeNull()
    expect(kunciModal('retro:1:1:2', undefined)).toBeNull()
  })
})

const anggota = (urut: number, keputusan: string): Anggota => ({
  id: `L${urut}`,
  urut,
  operatorId: `UJI-WB-${urut}`,
  jabatan: `UJI-JABATAN-${urut}`,
  keputusan,
  komentar: '',
  tanggal: '',
})

const LABEL: Record<string, string> = { '0': 'Waiting', '1': 'Approved', '2': 'Reject' }

/** Grid tangga seperti `models.barisTangga`: tersimpan lalu perluasan (menunggu). */
const bagianTangga = (t: Anggota[], perluasan: string[] = []): Bagian => ({
  kunci: 'tangga',
  judul: 'List of Committee',
  grid: [
    grid({
      baris: [
        ...t.map((a, i) => ({
          No: String(i + 1),
          jabatan: a.jabatan,
          keputusan: LABEL[a.keputusan] ?? a.keputusan,
          tanggal: a.tanggal,
          komentar: a.komentar,
        })),
        ...perluasan.map((j, i) => ({
          No: String(t.length + i + 1),
          jabatan: j,
          keputusan: 'Waiting',
          tanggal: '',
          komentar: '',
        })),
      ],
    }),
  ],
})

describe('langkahTangga', () => {
  it('disetujui / ditolak / berjalan (menunggu PERTAMA) / menunggu, label Status VERBATIM server', () => {
    const t = [anggota(1, '1'), anggota(2, '0'), anggota(3, '0')]
    const l = langkahTangga(bagianTangga(t), t)
    expect(l.map((x) => x.keadaan)).toEqual(['setuju', 'berjalan', 'menunggu'])
    expect(l.map((x) => x.status)).toEqual(['Approved', 'Waiting', 'Waiting'])
    expect(l.map((x) => x.no)).toEqual(['1', '2', '3'])
    const tolak = [anggota(1, '2'), anggota(2, '2')]
    expect(langkahTangga(bagianTangga(tolak), tolak).map((x) => x.keadaan)).toEqual(['tolak', 'tolak'])
  })

  it('perluasan tingkat 1 (KCF-02, belum disimpan) tampil sebagai langkah menunggu', () => {
    const t = [anggota(1, '0')]
    const l = langkahTangga(bagianTangga(t, ['UJI-JABATAN-2', 'UJI-JABATAN-3']), t)
    expect(l.map((x) => [x.jabatan, x.keadaan])).toEqual([
      ['UJI-JABATAN-1', 'berjalan'],
      ['UJI-JABATAN-2', 'menunggu'],
      ['UJI-JABATAN-3', 'menunggu'],
    ])
  })

  it('tanpa bagian tangga (TT3 / TT4) = tanpa langkah', () => {
    expect(langkahTangga(undefined, [anggota(1, '0')])).toEqual([])
  })
})

describe('keadaanKasus', () => {
  const layar = (tangga: Anggota[], extra: Partial<Layar> = {}, count = 1, statusWork = ''): Layar => ({
    kasus: {
      id: 'KMT-UJI001',
      klaimId: 'CLM-UJI001',
      adjustmentId: 'UJI-ADJ-1',
      transferType: '2',
      komiteLoop: tangga.length,
      komiteCount: count,
      acceptStatus: '',
      usulTutup: '',
      usulCadang: '',
      tahap: 'Komite_Flow',
      statusWork,
      pembuatId: 'UJI-OP',
      pembuatNama: 'UJI',
      tglCreate: '',
      tglUpdate: '',
      tangga,
    },
    judul: ['CLAIM COMMITTEE -', 'ADJUSTMENT'],
    ubin: [],
    bagian: [],
    isian: {
      nilai: { acceptStatus: '', comment: '', usulTutup: false, usulCadang: false },
      tampilUsul: true,
      terbuka: true,
      pilihanTerima: [],
      label: {},
    },
    tombol: [],
    bolehKerja: false,
    ...extra,
  })

  it('pemegang = giliran anda', () => {
    expect(keadaanKasus(layar([anggota(1, '0')], { bolehKerja: true }))).toEqual({ jenis: 'anda' })
  })

  it('bukan pemegang = menunggu jabatan tingkat berjalan (KomiteCount)', () => {
    const t = [anggota(1, '1'), anggota(2, '0'), anggota(3, '0')]
    expect(keadaanKasus(layar(t, {}, 2))).toEqual({ jenis: 'tunggu', jabatan: 'UJI-JABATAN-2' })
    // KomiteCount tanpa baris menunggu: cadangan baris menunggu pertama
    expect(keadaanKasus(layar(t, {}, 9))).toEqual({ jenis: 'tunggu', jabatan: 'UJI-JABATAN-2' })
  })

  it('kasus tertutup / tanpa baris menunggu = selesai', () => {
    expect(keadaanKasus(layar([anggota(1, '0')], {}, 1, 'Resolved-Completed'))).toEqual({ jenis: 'selesai' })
    expect(keadaanKasus(layar([anggota(1, '2'), anggota(2, '2')], {}, 2))).toEqual({ jenis: 'selesai' })
    const l = layar([])
    l.kasus.tangga = null
    expect(keadaanKasus(l)).toEqual({ jenis: 'selesai' })
  })
})
