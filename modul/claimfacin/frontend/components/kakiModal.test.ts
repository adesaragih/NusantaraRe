// Disalin dari `modul/claimnonprop/frontend/components/kakiModal.test.ts` (pola, bukan impor; asal Claim Prop): keputusan
// work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Laporan work owner 09-10-2026 pada popup komite: "cancel kok ada 2?" - Cancel section ikut tergambar di isi, padahal
// `Modal` sudah punya tombol batal di kakinya. Tombol di akhir isi modal pindah ke kaki; tombol penutup (Claim Fac In:
// tombol TANPA aksi yang aktif) tidak digambar ulang, labelnya menjadi label tombol batal `Modal`.

import { describe, expect, it } from 'vitest'

import type { Tata } from '../api'
import { pisahKakiModal, tombolPenutup } from './susun'

const medan = (jalur: string, label: string): Tata => ({ jenis: 'medan', jalur, label, kendali: 'area' })
const label = (teks: string): Tata => ({ jenis: 'label', label: teks })
const tombol = (id: string, teks: string, aksi = ''): Tata =>
  aksi === '' ? { jenis: 'tombol', id, label: teks } : { jenis: 'tombol', id, label: teks, aksi }
const bagian = (...anak: Tata[]): Tata => ({ jenis: 'bagian', anak })
const sebaris = (...anak: Tata[]): Tata => ({ jenis: 'bagian', letak: 'sebaris', anak })

describe('tombolPenutup', () => {
  it('tombol tanpa aksi yang aktif = penutup; tombol OQ (tanpa aksi, nonaktif) dan tombol beraksi bukan', () => {
    expect(tombolPenutup(tombol('BatalPolis', 'Cancel'))).toBe(true)
    expect(tombolPenutup({ ...tombol('SendRejectClaimToKomite2', 'Yes'), nonaktif: true })).toBe(false)
    expect(tombolPenutup(tombol('SetInputParam', 'Submit', 'SetInputParam'))).toBe(false)
    expect(tombolPenutup(medan('Remarks', 'Remarks'))).toBe(false)
  })
})

describe('pisahKakiModal', () => {
  it('ViewPolis: baris sebaris Cancel / Submit di akhir pindah ke kaki; Cancel tidak digambar ulang', () => {
    const m = pisahKakiModal([
      sebaris(
        medan('SearchType', 'Search Type'),
        medan('SearchName', ''),
        tombol('SearchPolis', 'Search', 'SearchPolis'),
      ),
      bagian(medan('InputParam.CARI4', 'Detail')),
      sebaris(tombol('BatalPolis', 'Cancel'), tombol('SetInputParam', 'Submit', 'SetInputParam')),
    ])
    expect(m.isi).toHaveLength(2)
    expect(m.kaki.map((t) => t.id)).toEqual(['SetInputParam'])
    expect(m.batal).toBe('Cancel')
  })

  it('PreventRejectClaim: label konfirmasi tetap di isi, Yes ke kaki, batal bawaan Modal', () => {
    const m = pisahKakiModal([
      medan('TempCommiteClaim.Remarks', 'Remarks'),
      bagian(label('Are you sure want close this claim?'), sebaris(tombol('CloseClaim', 'Yes', 'CloseClaim'))),
    ])
    expect(m.isi.map((t) => t.label)).toEqual(['Remarks', 'Are you sure want close this claim?'])
    expect(m.kaki.map((t) => t.id)).toEqual(['CloseClaim'])
    expect(m.batal).toBeUndefined()
  })

  it('Comittee: tombol Send di akhir pindah ke kaki', () => {
    const m = pisahKakiModal([
      medan('Remarks', 'Remarks'),
      tombol('KirimKomite', 'Send Claim to Committee', 'KirimKomite'),
    ])
    expect(m.isi.map((t) => t.label)).toEqual(['Remarks'])
    expect(m.kaki.map((t) => t.id)).toEqual(['KirimKomite'])
  })

  it('Reject Claim: tombol Yes nonaktif (naJika IsReject = 1) tetap di kaki, bukan penutup', () => {
    const oq: Tata = { ...tombol('SendRejectClaimToKomite2', 'Yes'), nonaktif: true }
    const m = pisahKakiModal([
      medan('Remarks', 'Remarks'),
      label('Are you sure want to reject this claim ?'),
      sebaris(oq),
    ])
    expect(m.kaki.map((t) => t.id)).toEqual(['SendRejectClaimToKomite2'])
    expect(m.batal).toBeUndefined()
  })

  it('ProtectDOL: bagian berjudul berisi tombol tidak dipindah (Cancel di dalamnya menutup lewat penutup konteks)', () => {
    const lanjut: Tata = {
      jenis: 'bagian',
      label: 'Lanjutkan?',
      anak: [sebaris(tombol('BatalDOL', 'Cancel'), tombol('SubmitDOL2', 'Submit', 'Submit'))],
    }
    const m = pisahKakiModal([lanjut])
    expect(m.isi).toHaveLength(1)
    expect(m.kaki).toEqual([])
  })

  it('sebaris berisi medan tidak dipindah', () => {
    const m = pisahKakiModal([sebaris(medan('Pct', 'Pct'), tombol('X', 'X', 'Aksi'))])
    expect(m.isi).toHaveLength(1)
    expect(m.kaki).toEqual([])
  })
})
