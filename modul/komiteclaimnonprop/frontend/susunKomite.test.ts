// Tata letak layar komite Non Prop (pola Komite Claim Prop): bagian server dikelompokkan ke kartu berjudul, ringkasan
// di kepala, tangga sebagai langkah. Isi, label medan / grid, dan urutan di dalam kelompok tetap dari server.

import { describe, expect, it } from 'vitest'

import type { Anggota, Bagian, Layar } from './api'
import { kelompokkanBagian, langkahTangga, ringkasan } from './susunKomite'

const bagian = (kunci: string, extra: Partial<Bagian> = {}): Bagian => ({ kunci, ...extra })

describe('kelompokkanBagian', () => {
  it('bagian ShowTransfer masuk kartu berurutan; tangga dipisah; bagian tak dikenal tetap tampil di akhir', () => {
    const k = kelompokkanBagian([
      bagian('klaim', { judul: 'Claim Analysis' }),
      bagian('klaimKanan'),
      bagian('kejadian'),
      bagian('klaimAkseptasi'),
      bagian('totalKlaim'),
      bagian('alokasi'),
      bagian('dibayar'),
      bagian('spreading'),
      bagian('bayar'),
      bagian('bank'),
      bagian('bank2'),
      bagian('tangga'),
      bagian('teksKomite'),
      bagian('baruUJI', { judul: 'UJI Baru' }),
    ])
    expect(k.kelompok.map((x) => [x.judul, x.bagian.map((b) => b.kunci)])).toEqual([
      ['Claim Analysis', ['klaim', 'klaimKanan']],
      ['Loss Details', ['kejadian']],
      ['Claim Acceptation', ['klaimAkseptasi', 'totalKlaim']],
      ['Allocation', ['alokasi', 'dibayar']],
      ['Spreading', ['spreading']],
      ['Payment & Bank', ['bayar', 'bank', 'bank2']],
      ['Committee Notes', ['teksKomite']],
      ['UJI Baru', ['baruUJI']],
    ])
    expect(k.tangga?.kunci).toBe('tangga')
  })

  it('kelompok tanpa bagian tidak tampil (S36 tanpa IsPrevious: tanpa kartu Payment & Bank)', () => {
    expect(kelompokkanBagian([bagian('klaim'), bagian('spreading')]).kelompok.map((x) => x.judul)).toEqual([
      'Claim Analysis',
      'Spreading',
    ])
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

describe('langkahTangga', () => {
  it('disetujui / ditolak / sedang berjalan (menunggu PERTAMA) / menunggu', () => {
    expect(langkahTangga([anggota(1, '1'), anggota(2, '0'), anggota(3, '0')]).map((l) => l.keadaan)).toEqual([
      'setuju',
      'berjalan',
      'menunggu',
    ])
    expect(langkahTangga([anggota(1, '2'), anggota(2, '2')]).map((l) => l.keadaan)).toEqual(['tolak', 'tolak'])
  })
})

describe('ringkasan', () => {
  const layar = (tangga: Anggota[]): Layar => ({
    kasus: {
      id: 'KMTNP-UJI001',
      klaimId: 'CLMNP-UJI001',
      adjustmentId: 'UJI-ADJ-1',
      komiteLoop: tangga.length,
      komiteCount: 1,
      acceptStatus: '',
      statusWork: '',
      pembuatNama: 'UJI',
      tangga,
    },
    judul: [],
    bagian: [
      bagian('klaim', {
        medan: [
          { label: 'Policy No', nilai: 'UJI-POL-1', jenis: 'teks' },
          { label: 'Policy No Ceding', nilai: 'UJI-POLCED', jenis: 'teks' },
          { label: 'Insured Name', nilai: 'UJI TERTANGGUNG', jenis: 'teks' },
        ],
      }),
      bagian('klaimKanan', {
        medan: [
          { label: 'PIC Name', nilai: 'UJI', jenis: 'teks' },
          { label: '', nilai: 'UJI-K-1 / CLMNP-UJI001', jenis: 'teks' },
        ],
      }),
    ],
    isian: {
      nilai: {
        acceptStatus: '',
        comment: '',
        isSubjectivity: false,
        subjectivityNote: '',
        usulTutup: false,
        usulCadang: false,
      },
      terbuka: true,
      pilihanTerima: [],
      pilihanSubjectivityNote: [],
      label: { acceptStatus: '', comment: '', isSubjectivity: '', subjectivityNote: '', usulTutup: '', usulCadang: '' },
    },
    tombol: [],
    bolehKerja: true,
  })

  it('nomor klaim (sel S9 tanpa label), polis, tertanggung, tingkat berjalan', () => {
    const r = ringkasan(layar([anggota(1, '1'), anggota(2, '0'), anggota(3, '0')]))
    expect(r).toEqual({
      noKlaim: 'UJI-K-1 / CLMNP-UJI001',
      polis: 'UJI-POL-1',
      tertanggung: 'UJI TERTANGGUNG',
      tingkat: { ke: 2, dari: 3, jabatan: 'UJI-JABATAN-2' },
    })
  })

  it('tanpa baris menunggu: tanpa tingkat berjalan; tanpa NoClaim: ID klaim', () => {
    const l = layar([anggota(1, '1')])
    l.bagian = l.bagian.filter((b) => b.kunci !== 'klaimKanan')
    const r = ringkasan(l)
    expect(r.tingkat).toBeNull()
    expect(r.noKlaim).toBe('CLMNP-UJI001')
  })
})
