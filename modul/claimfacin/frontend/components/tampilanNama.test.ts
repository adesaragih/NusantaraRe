// Disalin dari `modul/claimnonprop/frontend/components/tampilanNama.test.ts` (pola, bukan impor; asal Claim Prop):
// keputusan work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac
// In tahap 1). Consultant / Adjuster dipilih dan ditampilkan lewat NAMA, ID tetap disimpan (perintah work owner
// 08-10-2026: "yang dropdown hanya dari namanya aja; untuk ID dihapus dari tampilan, tapi tetap simpan ID").

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
})

describe('sel autocomplete tanpa tampilan (Claim Fac In)', () => {
  it('Occupation / Object Name: label = nama, kode pilihan sebagai keterangan; nilai tetap kode', () => {
    expect(butirSaring({ sumber: 'okupasi' }, [{ nilai: 'UJI-OKP-1', label: 'UJI OKUPASI' }])).toEqual([
      { value: 'UJI-OKP-1', label: 'UJI OKUPASI', keterangan: 'UJI-OKP-1' },
    ])
  })

  it('nilai = label (Coverage, wilayah): nama coverage tambahan sebagai keterangan, selainnya tanpa keterangan', () => {
    expect(
      butirSaring({ sumber: 'covFire' }, [
        { nilai: 'UJI-COV', label: 'UJI-COV', tambahan: { CoverageNote: 'UJI CATATAN COVERAGE' } },
      ]),
    ).toEqual([{ value: 'UJI-COV', label: 'UJI-COV', keterangan: 'UJI CATATAN COVERAGE' }])
    expect(butirSaring({ sumber: 'provinsi' }, [{ nilai: 'UJI PROVINSI', label: 'UJI PROVINSI' }])).toEqual([
      { value: 'UJI PROVINSI', label: 'UJI PROVINSI', keterangan: undefined },
    ])
    expect(teksTerpilih({}, 'UJI PROVINSI', ambil, [])).toBe('UJI PROVINSI')
  })

  it('Name of Bank: nama bank, cabang dan nomor rekening sebagai keterangan; nilai = kunci rekening', () => {
    expect(
      butirSaring({ sumber: 'rekening' }, [
        {
          nilai: 'UJI-REK|UJI BANK',
          label: 'UJI BANK',
          tambahan: { BranchOfBank: 'UJI CABANG', NoAccount: 'UJI-REK' },
        },
      ]),
    ).toEqual([{ value: 'UJI-REK|UJI BANK', label: 'UJI BANK', keterangan: 'UJI CABANG - UJI-REK' }])
  })
})
