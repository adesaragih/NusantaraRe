// Disalin dari `modul/claimprop/frontend/components/kakiModal.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Laporan work owner 09-10-2026 pada popup "Komite klaim Treaty": "cancel kok ada 2?" - Cancel section ikut tergambar
// di isi, padahal `Modal` sudah punya tombol batal di kakinya. Tombol di akhir isi modal pindah ke kaki; tombol penutup
// (`TutupModal`) tidak digambar ulang, labelnya menjadi label tombol batal `Modal`.

import { describe, expect, it } from 'vitest'

import type { Tata } from '../api'
import { pisahKakiModal } from './susun'

const medan = (jalur: string, label: string): Tata => ({ jenis: 'medan', jalur, label, kendali: 'area' })
const label = (teks: string): Tata => ({ jenis: 'label', label: teks })
const tombol = (id: string, teks: string, aksi: string): Tata => ({ jenis: 'tombol', id, label: teks, aksi })
const bagian = (...anak: Tata[]): Tata => ({ jenis: 'bagian', anak })

describe('pisahKakiModal', () => {
  it('popup komite: Send ke kaki, Cancel section tidak digambar ulang', () => {
    const m = pisahKakiModal([
      bagian(medan('Tgl', 'Date'), medan('PIC', 'PIC')),
      bagian(label('Adjustment')),
      bagian(medan('Remarks', 'Remarks')),
      bagian(
        bagian(tombol('SendClaimToCommittee', 'Send Claim to Committee', 'AddKomiteTreatyChild')),
        tombol('CancelKomite', 'Cancel', 'TutupModal'),
      ),
    ])
    expect(m.isi.map((t) => t.label)).toEqual(['Date', 'PIC', 'Adjustment', 'Remarks'])
    expect(m.kaki.map((t) => t.id)).toEqual(['SendClaimToCommittee'])
    expect(m.batal).toBe('Cancel')
  })

  it('tombol Send belum tampil: kaki kosong, batal tetap satu', () => {
    const m = pisahKakiModal([medan('Remarks', 'Remarks'), bagian(tombol('CancelKomite', 'Cancel', 'TutupModal'))])
    expect(m.kaki).toEqual([])
    expect(m.batal).toBe('Cancel')
  })

  it('label tombol penutup dipakai (No) dan tombol di tengah isi tidak dipindah', () => {
    const m = pisahKakiModal([
      tombol('Pilih', 'Choose', 'PilihSesuatu'),
      medan('Remarks', 'Remarks'),
      label('Are you sure want close this claim?'),
      tombol('CloseYes', 'Yes', 'CloseClaimProp'),
      tombol('CloseNo2', 'No', 'TutupModal'),
    ])
    expect(m.isi.map((t) => t.id ?? t.label)).toEqual(['Pilih', 'Remarks', 'Are you sure want close this claim?'])
    expect(m.kaki.map((t) => t.id)).toEqual(['CloseYes'])
    expect(m.batal).toBe('No')
  })

  it('tanpa tombol penutup: batal bawaan Modal', () => {
    expect(pisahKakiModal([medan('Remarks', 'Remarks')]).batal).toBeUndefined()
  })
})
