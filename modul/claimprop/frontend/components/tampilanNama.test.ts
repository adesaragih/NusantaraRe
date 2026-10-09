// Consultant / Adjuster dipilih dan ditampilkan lewat NAMA, ID tetap disimpan (perintah work owner 08-10-2026: "yang
// dropdown hanya dari namanya aja; untuk ID dihapus dari tampilan, tapi tetap simpan ID").

import { describe, expect, it } from 'vitest'

import type { Pilihan } from '../api'
import { butirSaring, teksTerpilih } from './tampilanNama'

const opsi: Pilihan[] = [
  { nilai: 'UJI-ADJ-1', label: 'UJI NAMA SATU' },
  { nilai: 'UJI-ADJ-2', label: 'UJI NAMA DUA' },
]
const halaman: Record<string, string> = { 'ClaimData.ConsultantName': 'UJI NAMA SATU' }
const ambil = (jalur: string) => halaman[jalur] ?? ''

describe('medan ber-tampilan (Consultant / Adjuster)', () => {
  it('butir dropdown hanya nama: tanpa ID sebagai label maupun keterangan; nilai tetap ID', () => {
    const b = butirSaring({ tampilan: 'ClaimData.ConsultantName' }, opsi)
    expect(b).toEqual([
      { value: 'UJI-ADJ-1', label: 'UJI NAMA SATU' },
      { value: 'UJI-ADJ-2', label: 'UJI NAMA DUA' },
    ])
  })

  it('teks terpilih = nama dari jalur tampilan; ID tidak pernah tampil', () => {
    expect(teksTerpilih({ tampilan: 'ClaimData.ConsultantName' }, 'UJI-ADJ-1', ambil, [])).toBe('UJI NAMA SATU')
    // nama belum terisi di halaman: diambil dari daftar pilihan; tidak ada di daftar pun tetap bukan ID
    expect(teksTerpilih({ tampilan: 'ClaimData.Lain' }, 'UJI-ADJ-2', ambil, opsi)).toBe('UJI NAMA DUA')
    expect(teksTerpilih({ tampilan: 'ClaimData.Lain' }, 'UJI-ADJ-9', ambil, opsi)).toBe('')
  })

  it('medan tanpa tampilan (Province) tidak berubah: nilai sebagai label, nama sebagai keterangan', () => {
    expect(butirSaring({}, [{ nilai: 'UJI-P1', label: 'UJI PROVINSI' }])).toEqual([
      { value: 'UJI-P1', label: 'UJI-P1', keterangan: 'UJI PROVINSI' },
    ])
    expect(teksTerpilih({}, 'UJI-P1', ambil, [])).toBe('UJI-P1')
  })
})
