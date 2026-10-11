import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { ReactElement } from 'react'
import { describe, expect, it } from 'vitest'

import KonfirmasiProses from '../components/KonfirmasiProses'
import { PROSES_EDM } from '../labels'
import { adaPerubahan, jumlahDicentang, kalimatKonfirmasiProses, PILIHAN_KOSONG } from '../tampilan'

// Save → "Process Endorsement" di bawah layar + konfirmasi; Back dihapus - keputusan work owner 08-10-2026.

const SUMBER = readFileSync(join(__dirname, 'InputEDMLife.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('tombol Process Endorsement', () => {
  it('di bawah tabel peserta + paging, di atas rekap dan Keputusan', () => {
    const i = (s: string) => {
      const n = SUMBER.indexOf(s)
      expect(n, s).toBeGreaterThan(0)
      return n
    }
    expect(i('</header>')).toBeLessThan(i('<table className="inbox__tabel">'))
    expect(i('<table className="inbox__tabel">')).toBeLessThan(i('<Halaman halaman={halaman}'))
    expect(i('<Halaman halaman={halaman}')).toBeLessThan(i('<div className="edm-kasus__proses">'))
    expect(i('<div className="edm-kasus__proses">')).toBeLessThan(i('{kasus.rekap.length > 0 && ('))
    expect(i('{kasus.rekap.length > 0 && (')).toBeLessThan(i('<KeputusanEDM'))
  })

  it('kondisi tampil/nonaktif sama dengan Save lama; klik HANYA membuka konfirmasi', () => {
    const blok = SUMBER.slice(SUMBER.indexOf('{terbukaKasus && (\n        <div className="edm-kasus__proses">'))
    const tombol = blok.slice(0, blok.indexOf('</div>'))
    expect(tombol).toContain('disabled={kasus.sudahSimpan || menyimpan || tanpaPerubahan}')
    expect(tombol).toContain('onClick={() => setKonfirmasi(true)}')
    expect(tombol).toContain('PROSES_EDM.processEndorsement')
    expect(tombol).not.toContain('simpan()')
    expect(PROSES_EDM.processEndorsement).toBe('Process Endorsement')
  })

  it('Continue menjalankan simpan yang sama (simpanKasus(kasusId, pilihan)); Cancel hanya menutup', () => {
    const k = SUMBER.slice(SUMBER.indexOf('<KonfirmasiProses'), SUMBER.indexOf('/>', SUMBER.indexOf('<KonfirmasiProses')))
    expect(k).toContain('kalimat={kalimatKonfirmasiProses(kasus, pilihan)}')
    expect(k).toContain('onTutup={() => setKonfirmasi(false)}')
    expect(k).toContain('void simpan()')
    expect(SUMBER).toContain('await simpanKasus(kasusId, pilihan)')
    expect(SUMBER.split('simpan()').length - 1).toBe(1) // satu-satunya pemanggil: Continue
    expect(SUMBER).not.toContain('UMUM_EDM.kembali')
  })
})

type PropsModal = { judul: string; labelBatal: string; onTutup: () => void; aksi: ReactElement<{ onClick: () => void; children: string }> }

describe('KonfirmasiProses', () => {
  const buat = () => {
    const panggil: string[] = []
    const el = KonfirmasiProses({ kalimat: ['UJI-1'], onLanjut: () => panggil.push('lanjut'), onTutup: () => panggil.push('tutup') })
    return { panggil, p: (el as ReactElement<PropsModal>).props }
  }

  it('judul, Cancel, Continue', () => {
    const { p } = buat()
    expect(p.judul).toBe('Process Endorsement?')
    expect(p.labelBatal).toBe('Cancel')
    expect(p.aksi.props.children).toBe('Continue')
  })

  it('Cancel tidak memproses; Continue memproses', () => {
    const a = buat()
    a.p.onTutup()
    expect(a.panggil).toEqual(['tutup'])
    const b = buat()
    b.p.aksi.props.onClick()
    expect(b.panggil).toEqual(['lanjut'])
  })
})

describe('kalimat konfirmasi per EDM Type', () => {
  const REKAP = 'The summary will be recalculated. This can only be done once.'

  it('Perubahan Data: jumlah centang, nol, DELETE ALL, peserta baru', () => {
    const k = { edmType: '1', cacah: { Old: 10 } }
    expect(kalimatKonfirmasiProses(k, { ...PILIHAN_KOSONG, pilih: ['a', 'b', 'c'] })).toEqual([
      '3 participant(s) will be marked Delete and their amounts reversed.',
      REKAP,
    ])
    expect(kalimatKonfirmasiProses(k, PILIHAN_KOSONG)).toEqual(['No participants are marked Delete.', REKAP])
    expect(jumlahDicentang({ pilih: [], semua: true, kecuali: ['x'] }, k.cacah)).toBe(9)
    expect(kalimatKonfirmasiProses({ edmType: '1', cacah: { Old: 4, New: 2 } }, { pilih: [], semua: true, kecuali: [] })).toEqual([
      '4 participant(s) will be marked Delete and their amounts reversed.',
      '2 new participant(s) will be kept.',
      REKAP,
    ])
  })

  it('Batal: seluruh peserta Old; tanpa angka bila cacah tidak ada', () => {
    expect(kalimatKonfirmasiProses({ edmType: '3', cacah: { Old: 7 } }, PILIHAN_KOSONG)).toEqual([
      'All 7 existing participant(s) will be marked Batal and their amounts reversed.',
      REKAP,
    ])
    expect(kalimatKonfirmasiProses({ edmType: '3', cacah: {} }, PILIHAN_KOSONG)).toEqual([
      'All existing participants will be marked Batal and their amounts reversed.',
      REKAP,
    ])
  })
})

describe('proteksi tanpa perubahan (keputusan work owner 08-10-2026)', () => {
  const satu = (cacah: Record<string, number>) => ({ edmType: '1', cacah })

  it('Perubahan Data: nonaktif tanpa centang dan tanpa New; aktif bila salah satu ada', () => {
    expect(adaPerubahan(satu({ Old: 3 }), PILIHAN_KOSONG)).toBe(false)
    expect(adaPerubahan(satu({ Old: 3 }), { pilih: [], semua: true, kecuali: ['a', 'b', 'c'] })).toBe(false)
    expect(adaPerubahan(satu({ Old: 3 }), { ...PILIHAN_KOSONG, pilih: ['a'] })).toBe(true)
    expect(adaPerubahan(satu({ Old: 3 }), { pilih: [], semua: true, kecuali: ['a'] })).toBe(true)
    expect(adaPerubahan(satu({ Old: 3, New: 1 }), PILIHAN_KOSONG)).toBe(true)
  })

  it('Batal selalu aktif', () => {
    expect(adaPerubahan({ edmType: '3', cacah: {} }, PILIHAN_KOSONG)).toBe(true)
  })

  it('teks bantuan di samping tombol nonaktif; pesan sama dengan server', () => {
    const blok = SUMBER.slice(SUMBER.indexOf('<div className="edm-kasus__proses">'))
    const isi = blok.slice(0, blok.indexOf('</button>'))
    expect(isi).toContain('{!kasus.sudahSimpan && tanpaPerubahan && (')
    expect(isi).toContain('{PROSES_EDM.tanpaPerubahan}')
    expect(isi.indexOf('edm-kasus__bantuan')).toBeLessThan(isi.indexOf('<button'))
    expect(SUMBER).toContain('const tanpaPerubahan = !adaPerubahan(kasus, pilihan)')
    const GO = readFileSync(join(__dirname, '..', '..', 'backend', 'models', 'edm_perubahan.go'), 'utf8')
    expect(GO).toContain(`const PesanTanpaPerubahan = "${PROSES_EDM.tanpaPerubahan}"`)
  })
})

describe('Process Endorsement langkah final: CSV terkunci sesudahnya (keputusan work owner 08-10-2026)', () => {
  const awal = SUMBER.indexOf('<div className="edm-kasus__kepalapeserta">')
  const toolbar = SUMBER.slice(awal, SUMBER.indexOf('{periksa !== null', awal))

  it('Upload CSV, View Upload, Add CSV Data nonaktif bila sudahSimpan; kondisi lama tetap', () => {
    expect(toolbar).toContain('disabled={kasus.csvTerkunci || kasus.sudahSimpan}')
    expect(toolbar).toContain('disabled={kasus.sudahSimpan}')
    expect(toolbar).toContain('disabled={kasus.csvTerkunci || kasus.sudahSimpan || berkasCSV === null || menambah}')
    expect(toolbar.split("aria-describedby={kasus.sudahSimpan ? 'edm-csv-bantuan' : undefined}").length - 1).toBe(4) // + Download Template (10-10-2026)
  })

  it('teks bantuan sudah diproses tampil bila sudahSimpan; pesan sama dengan server', () => {
    expect(toolbar).toContain('{kasus.sudahSimpan && (')
    expect(toolbar).toContain('{PROSES_EDM.sudahDiproses}')
    const GO = readFileSync(join(__dirname, '..', '..', 'backend', 'models', 'edm_perubahan.go'), 'utf8')
    expect(GO).toContain(`const PesanSudahDiproses = "${PROSES_EDM.sudahDiproses}"`)
  })
})
