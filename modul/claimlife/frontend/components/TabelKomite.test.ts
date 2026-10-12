import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BarisKomite } from '../api'
import { selKomite, tingkatKomite } from './TabelKomite'
import { TABEL_KOMITE } from '../labels'

// Tabel komite di bawah inbox Claim Life (perintah work owner 09-10-2026: menu Komite Claim Life dihapus) - perilaku
// Inbox Komite tiket 01 yang dipindah ke sini.

const SUMBER = readFileSync(join(__dirname, 'TabelKomite.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS: BarisKomite = {
  kasusId: 'KMTLF-UJI1',
  tglUpdate: '',
  statusWork: '',
  klaimId: 'UJI-K',
  nomorKlaim: 'UJI-CLM',
  tingkatBerjalan: 2,
  komiteLoop: 3,
  nilaiKlaim: '1500000.25',
  mataUang: 'IDR',
  statusBaris: 'Outstanding',
}

describe('tabel komite inbox Claim Life', () => {
  it('tingkat berjalan dari seluruh tingkat', () => {
    expect(tingkatKomite(BARIS)).toBe('2 / 3')
  })
  it('sel kosong ditandai', () => {
    expect(selKomite('')).toBe('—')
    expect(selKomite('1500000.25')).toBe('1500000.25')
  })
  it('uang tidak lewat float; layar tidak menyaring ulang; kasus dibuka di tempat modul Komite Claim Life', () => {
    // \b: formatNumber( (format teks desimal, 12-10-2026) bukan konversi ke float.
    expect(SUMBER).not.toMatch(/\bNumber\(|parseFloat|parseInt/)
    expect(SUMBER).not.toContain('.filter(')
    expect(SUMBER).toContain('onBukaModul?.(MODUL_KOMITE, id)')
  })
})

describe('teks tabel komite berbahasa Inggris (11-10-2026)', () => {
  it('judul, kolom, dan pesan sejajar Inbox Claim Life; tanpa sisa kata Indonesia di layar', () => {
    expect(TABEL_KOMITE.judul).toBe('Committee Inbox')
    expect([TABEL_KOMITE.nomorKlaim, TABEL_KOMITE.tingkat, TABEL_KOMITE.nilaiKlaim, TABEL_KOMITE.statusBaris]).toEqual([
      'Parent Claim',
      'Level',
      'Claim Amount',
      'Row Status',
    ])
    for (const v of Object.values(TABEL_KOMITE)) expect(v).not.toMatch(/\b(Tidak|tidak|kasus|Kasus|dari|sejak|hari ini|Klaim)\b/)
    // 12-10-2026: pencacah "n of total" diganti paginasi Halaman inti (seperti Inbox Claim Life).
    expect(SUMBER).toContain('<Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />')
    expect(SUMBER).toContain('{TABEL_KOMITE.sejak}')
    expect(SUMBER).not.toMatch(/ dari \{| sejak \{/)
  })
})

describe('Committee Inbox dirapikan (12-10-2026)', () => {
  it('format: waktu dd/mm/yyyy HH:mm, Claim Amount ribuan, Work Status kosong = Waiting for decision', async () => {
    const { waktuKomite, nilaiKomite, statusWorkKomite } = await import('./TabelKomite')
    expect(waktuKomite('2026-10-12 17:02:58')).toBe('12/10/2026 17:02')
    expect(waktuKomite('2026-10-12T17:02:58Z')).toBe('12/10/2026 17:02')
    expect(waktuKomite('')).toBe('—')
    expect(nilaiKomite('2342342')).toBe('2.342.342')
    expect(nilaiKomite('1500000.25')).toBe('1.500.000,25')
    expect(nilaiKomite('')).toBe('—')
    expect(statusWorkKomite('')).toBe('Waiting for decision')
    expect(statusWorkKomite('Resolved-Completed')).toBe('Resolved-Completed')
  })
  it('kulit tabel inbox, Level tengah, angka kanan, Row Status chip; judul halaman di mode utama', () => {
    expect(SUMBER).toContain('<table className="inbox__tabel claimlife__tabel">')
    expect(SUMBER).toContain('<td className="claimlife__tengah">{tingkatKomite(b)}</td>')
    expect(SUMBER).toContain('<td className="claimlife__angka">{nilaiKomite(b.nilaiKlaim)}</td>')
    expect(SUMBER).toContain('<span className="claimlife__chip">{selKomite(b.statusBaris)}</span>')
    expect(SUMBER).toContain('<h2 id="claimlife-judul-komite" className="inbox__judul">')
  })
  it('akun tanpa antrian Claim Life: kotak ber-kulit + TabelKomite utama; "No queue for your role" dibuang', () => {
    const inbox = readFileSync(join(__dirname, '..', 'pages', 'InboxClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(inbox).not.toContain('No queue for your role')
    expect(inbox).toContain('<section className="inbox claimlife__akar">\n        <TabelKomite')
    expect(inbox).toContain('kosong={<Kosong pesan={TABEL_KOMITE.kosongAkun} />}')
    expect(TABEL_KOMITE.kosongAkun).toBe('No Claim Life or committee cases for your account.')
    // Kalimat kosong HANYA bila daftar komite pun kosong (sesudah dibaca).
    expect(SUMBER).toContain('return selesai && kosong !== undefined ? <>{kosong}</> : null')
  })
})
