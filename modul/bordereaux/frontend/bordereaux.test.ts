import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, ambilLama, buka, cariCedant, cariTreaty, hapus, salinLama, simpan, submit, unggahCsv, type BerkasLama } from './api'
import {
  businessUntuk,
  catatanLama,
  FILTER_KOSONG,
  formatAngka,
  keIso,
  keKabel,
  kunciKombinasi,
  potong,
  ringkasSalin,
  saringLama,
  susunKepala,
  tampilanChart,
  teksPersen,
} from './aturan'
import { BDX, MENU_BDX, teksStatus } from './labels'

describe('label Bordereaux', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 913, golongan MASTER TREATY', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/913_m_nav_menu_bordereaux.sql`, 'utf8')
    expect(sql).toContain(`'bordereaux', '${MENU_BDX.kelompok}', 'MASTER TREATY', 'bordereaux'`)
  })

  it('caption Pega PortalBordereaux_Sec dan InputBordereaux', () => {
    expect([BDX.inputData, BDX.searchFilter, BDX.applyFilter, BDX.filterBusiness, BDX.filterReffNoSoa, BDX.reffNoBdx]).toEqual([
      'Input Data', 'Search / Filter', 'Apply Filter', 'BUSINESS', 'Reff No Soa', 'Reff No BDX',
    ])
    expect([BDX.reportStart, BDX.reportEnd, BDX.cedingCo, BDX.treatyName, BDX.position, BDX.status]).toEqual([
      'Bordereaux Report Start', 'Bordereaux Report End', 'Ceding Co', 'Treaty Name', 'Position', 'Status',
    ])
    expect([BDX.typeBusiness, BDX.chooseMasterTreaty, BDX.masterId, BDX.sobName, BDX.reffNoOfSoa, BDX.reffNoOfBdx]).toEqual([
      'Type Business', 'Choose Master Treaty', 'Master ID', 'SOB Name', 'Reff No of SOA', 'Reff No of Bordereaux',
    ])
    expect([BDX.tabDetails, BDX.tabSummary, BDX.tabSubmit, BDX.unggahPetunjuk, BDX.template, BDX.uploadCsv]).toEqual([
      'Details', 'Summary', 'Submit', 'Please upload file with .csv format', 'Template', 'Upload CSV',
    ])
    expect([BDX.cedant, BDX.contractName, BDX.reinsuranceType, BDX.sourceOfBusiness, BDX.choose]).toEqual([
      'Cedant', 'Contract Name', 'Reinsurance Type', 'Source of Business', 'Choose',
    ])
    expect(BDX.ringkasan['PREMIUM']).toEqual(['Total Premium Reinsurer', 'Total Premium NusantaraRe'])
    expect([teksStatus(''), teksStatus('Accept')]).toEqual(['Draft', 'Accept'])
  })
})

describe('aturan Bordereaux', () => {
  it('format angka Indonesia dua desimal, tanggal kabel <-> input', () => {
    expect([formatAngka('1234567.5'), formatAngka('999.995'), formatAngka(''), formatAngka('x')]).toEqual(['1.234.567,50', '1.000,00', '', 'x'])
    expect([keKabel('2026-07-01'), keIso('30-09-2026'), keKabel(''), keIso('2026-09-30')]).toEqual(['01-07-2026', '2026-09-30', '', ''])
  })

  it('kombinasi dan SUBROGATION = BONDING', () => {
    expect(kunciKombinasi('PREMIUM', 'FIRE')).toBe('PREMIUM|FIRE')
    expect([businessUntuk('SUBROGATION', 'FIRE'), businessUntuk('CLAIM', 'FIRE')]).toEqual(['BONDING', 'FIRE'])
  })
})

describe('chart Bordereaux', () => {
  const irisan = [
    { business: 'FIRE', type: 'PREMIUM', cedingId: 'UJI-AG-1', cedingName: 'UJI CEDANT SATU', jumlah: 2 },
    { business: 'FIRE', type: 'CLAIM', cedingId: 'UJI-AG-1', cedingName: 'UJI CEDANT SATU', jumlah: 1 },
    { business: 'FIRE', type: 'PREMIUM', cedingId: 'UJI-AG-2', cedingName: 'UJI CEDANT DUA', jumlah: 3 },
    { business: 'BONDING', type: 'SUBROGATION', cedingId: 'UJI-AG-2', cedingName: 'UJI CEDANT DUA', jumlah: 2 },
  ]

  it('kelompok Business ditumpuk per Type, terbanyak dulu', () => {
    const v = tampilanChart(irisan, {})
    expect([v.tingkat, v.total, v.bisaBuka]).toEqual(['business', 8, true])
    expect(v.batang.map((b) => [b.kunci, b.jumlah, b.persen])).toEqual([
      ['FIRE', 6, 75],
      ['BONDING', 2, 25],
    ])
    expect(v.batang[0]!.bagian.map((s) => [s.kunci, s.jumlah, s.seri, s.lebar])).toEqual([
      ['PREMIUM', 5, 0, (5 / 6) * 100],
      ['CLAIM', 1, 1, (1 / 6) * 100],
    ])
    expect(v.batang[1]!.bagian[0]!.seri).toBe(2)
  })

  it('dibuka per Ceding satu Business; jumlah sama = urut nama', () => {
    const c = tampilanChart(irisan, { business: 'FIRE' })
    expect([c.tingkat, c.total, c.bisaBuka]).toEqual(['ceding', 6, false])
    expect(c.batang.map((b) => [b.label, b.jumlah])).toEqual([
      ['UJI CEDANT DUA', 3],
      ['UJI CEDANT SATU', 3],
    ])
    expect(teksPersen(33.333)).toBe('33,3%')
  })
})

describe('klien Bordereaux', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')}`)
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar({ ...FILTER_KOSONG, ceding: 'uji', start: '01-07-2026' }, 2)
    await buka('BDX-2026.10.04.07123')
    await unggahCsv('PREMIUM', 'FIRE', 'a;b')
    await simpan({ bdxId: '', type: 'PREMIUM', business: 'FIRE', masterId: 'UJI', reportStart: '', reportEnd: '', reffNoSoa: '', reffNoBdx: '', baris: [] })
    await submit('BDX-UJI', true, '')
    await hapus('BDX-UJI')
    await cariCedant('uji')
    await cariTreaty('UJI-AG-1')
    await ambilLama()
    await salinLama(['BDX-UJI'])
    expect(panggil).toEqual([
      'GET /api/bordereaux?halaman=2&ceding=uji&start=01-07-2026',
      'GET /api/bordereaux/berkas/BDX-2026.10.04.07123',
      'POST /api/bordereaux/unggah-csv',
      'POST /api/bordereaux/simpan',
      'POST /api/bordereaux/berkas/BDX-UJI/submit',
      'POST /api/bordereaux/berkas/BDX-UJI/hapus',
      'GET /api/bordereaux/cedant?q=uji',
      'GET /api/bordereaux/master-treaty?ceding=UJI-AG-1',
      'GET /api/bordereaux/lama',
      'POST /api/bordereaux/lama/salin',
    ])
  })
})

describe('kepala grid format Excel (05-10-2026)', () => {
  it('grup bersebelahan digabung, kolom tanpa grup dua baris, kolom sembunyi tidak tampil', () => {
    const k = (kolom: string, grup?: string, sembunyi?: boolean) => ({
      kolom, judul: kolom, jenis: 'teks' as const, kepala: `*${kolom}`, grup, sembunyi,
    })
    const s = susunKepala([k('COB'), k('TREATY_TYPE', undefined, true), k('POI_START', 'POI'), k('POI_END', 'POI'), k('UW')])
    expect(s.excel).toBe(true)
    expect(s.duaBaris).toBe(true)
    expect(s.tampil.map((x) => x.kolom)).toEqual(['COB', 'POI_START', 'POI_END', 'UW'])
    expect(s.baris1.map((x) => [x.teks, x.colSpan, x.rowSpan])).toEqual([['*COB', 1, 2], ['POI', 2, 1], ['*UW', 1, 2]])
    expect(s.baris2.map((x) => x.kolom)).toEqual(['POI_START', 'POI_END'])
  })

  it('tanpa grup dan tanpa kepala Excel: satu baris, judul CSV', () => {
    const s = susunKepala([{ kolom: 'A', judul: 'Judul A', jenis: 'angka' }])
    expect([s.excel, s.duaBaris, s.baris1[0]?.teks, s.baris1[0]?.rowSpan, s.baris1[0]?.angka]).toEqual([false, false, 'Judul A', 1, true])
  })
})

describe('Copy Old Data', () => {
  const lama = (isi: Partial<BerkasLama>): BerkasLama => ({
    id: 'BDX-UJI.1',
    type: 'PREMIUM',
    business: 'FIRE',
    ceding: 'UJI CEDANT SATU',
    treaty: 'UJI KONTRAK SATU',
    status: '',
    barisJson: 0,
    barisTabel: 0,
    komentar: 0,
    riwayat: 0,
    tanpaHeader: false,
    ...isi,
  })

  it('catatan = yang akan disalin: header, detail bila tabel nol, riwayat bila tabel kosong', () => {
    expect(catatanLama(lama({ tanpaHeader: true, barisJson: 2, komentar: 1 }))).toBe('header · 2 detail rows · 1 history row')
    expect(catatanLama(lama({ barisJson: 3, barisTabel: 3, komentar: 2 }))).toBe('2 history rows')
    expect(catatanLama(lama({ barisJson: 1, komentar: 2, riwayat: 2 }))).toBe('1 detail row')
    expect(catatanLama(lama({}))).toBe(BDX.lamaTakTerbaca)
  })

  it('saring menurut ID, ceding, treaty, Type, atau Business', () => {
    const d = [lama({ id: 'BDX-A' }), lama({ id: 'BDX-B', ceding: 'UJI CEDANT DUA', business: 'ENGINEERING' })]
    expect(saringLama(d, ' dua ').map((b) => b.id)).toEqual(['BDX-B'])
    expect(saringLama(d, 'engineering').map((b) => b.id)).toEqual(['BDX-B'])
    expect(saringLama(d, '')).toHaveLength(2)
  })

  it('Process Copy dikirim per kelompok dan diringkas per status', () => {
    expect(potong(['a', 'b', 'c'], 2)).toEqual([['a', 'b'], ['c']])
    expect(potong([], 10)).toEqual([])
    expect(
      ringkasSalin([
        { id: 'a', status: 'disalin', pesan: [] },
        { id: 'b', status: 'disalin', pesan: [] },
        { id: 'c', status: 'ditolak', pesan: ['x'] },
      ]),
    ).toEqual({ disalin: 2, sudahAda: 0, ditolak: 1, gagal: 0 })
  })
})
