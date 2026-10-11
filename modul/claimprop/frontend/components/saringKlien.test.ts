// Laporan work owner 10-10-2026 "aksi ini tidak tersedia di layar kasus saat ini: centang Close Without Payment":
// centang di pop-up PreventRejectClaimProp hanya postValue, jadi baris Yes / No harus bertukar di layar (Pega
// `pyIsClientWhen = true`). Pohon di bawah = bentuk kiriman server `LayarTutupKlaim`.

import { describe, expect, it } from 'vitest'

import type { SyaratKlien, Tata } from '../api'
import { pisahKakiModal, saringKlien } from './susun'

const CWP = 'TempCommiteClaim.AllocationShareSalvage'
const KOMITE = 'Komite.CARI1'

const jika = (t: Tata, jalur: string, nilai: string, beda = false): Tata => {
  const s: SyaratKlien = beda ? { jalur, nilai, beda } : { jalur, nilai }
  return { ...t, tampilJika: s }
}
const medan = (jalur: string, label: string, kendali = 'area'): Tata => ({ jenis: 'medan', jalur, label, kendali })
const label = (teks: string): Tata => ({ jenis: 'label', label: teks })
const tombol = (id: string, teks: string, aksi: string): Tata => ({ jenis: 'tombol', id, label: teks, aksi })
const bagian = (...anak: Tata[]): Tata => ({ jenis: 'bagian', anak })

const popUp: Tata[] = [
  jika(medan(CWP, 'Close Without Payment', 'centang'), KOMITE, ''),
  jika(
    bagian(
      medan('TempCommiteClaim.CircumtansesCouseOfLoss', 'Chronology'),
      medan('TempCommiteClaim.Remarks', 'Remarks'),
      jika(
        bagian(
          label('Are you sure want close this claim without payment?'),
          tombol('CloseWithoutPayment', 'Yes', 'SendCloseClaimToKomite'),
          tombol('CloseNo', 'No', 'TutupModal'),
        ),
        CWP,
        'true',
      ),
      jika(
        bagian(
          label('Are you sure want close this claim?'),
          tombol('CloseYes', 'Yes', 'CloseClaimProp'),
          tombol('CloseNo2', 'No', 'TutupModal'),
        ),
        CWP,
        'true',
        true,
      ),
    ),
    KOMITE,
    '',
  ),
  jika(
    bagian(
      label('Success Create Request to Committee'),
      { ...medan(KOMITE, 'Komite Number', 'teks'), hanyaBaca: true },
      tombol('CloseSukses', 'Close', 'TutupModal'),
    ),
    KOMITE,
    '',
    true,
  ),
]

const keadaan = (nilai: Record<string, string>) => pisahKakiModal(saringKlien(popUp, (j) => nilai[j] ?? ''))

describe('saringKlien - pop-up Close Claim', () => {
  it('belum dicentang: Yes biasa (CloseClaimProp)', () => {
    const m = keadaan({})
    expect(m.kaki.map((t) => t.aksi)).toEqual(['CloseClaimProp'])
    expect(m.isi.map((t) => t.label)).toContain('Are you sure want close this claim?')
    expect(m.batal).toBe('No')
  })

  it('dicentang: Yes tanpa bayar (SendCloseClaimToKomite), tanpa server', () => {
    const m = keadaan({ [CWP]: 'true' })
    expect(m.kaki.map((t) => t.aksi)).toEqual(['SendCloseClaimToKomite'])
    expect(m.isi.map((t) => t.label)).toEqual([
      'Close Without Payment',
      'Chronology',
      'Remarks',
      'Are you sure want close this claim without payment?',
    ])
  })

  it('centang dilepas lagi ("false"): kembali ke Yes biasa', () => {
    expect(keadaan({ [CWP]: 'false' }).kaki.map((t) => t.aksi)).toEqual(['CloseClaimProp'])
  })

  // perintah work owner 10-10-2026: centang disembunyikan, nomor kasus komite tampil
  it('sesudah Yes tanpa bayar (Komite.CARI1 terisi): tanpa centang, pesan sukses + Komite Number + Close', () => {
    const m = keadaan({ [CWP]: 'true', [KOMITE]: 'TKMT-UJI001' })
    expect(m.kaki).toEqual([])
    expect(m.batal).toBe('Close')
    expect(m.isi.map((t) => t.label)).toEqual(['Success Create Request to Committee', 'Komite Number'])
    expect(m.isi.find((t) => t.jalur === KOMITE)?.hanyaBaca).toBe(true)
  })
})
