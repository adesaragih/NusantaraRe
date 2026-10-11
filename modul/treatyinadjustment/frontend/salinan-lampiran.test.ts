// Ringkasan salinan lampiran master di pesan Save draf penyesuaian
// (`TreatyInEDMSetValue` [8] `TreatyRevisionCopyAttachment`).

import { describe, expect, it } from 'vitest'

import { gabungPesanSalinan, pesanSalinanLampiran } from './komponen/salinanLampiran'

const PESAN = 'Data Sudah Disimpan Dengan ID : 1001001/R02'

describe('ringkasan salinan lampiran', () => {
  it('tanpa salinan, pesan prosedur tidak berubah', () => {
    expect(gabungPesanSalinan(PESAN, undefined)).toBe(PESAN)
    expect(gabungPesanSalinan(PESAN, null)).toBe(PESAN)
  })

  it('menyebut sumber dan cacah, lalu berkas yang gagal beserta alasannya', () => {
    const s = {
      sumber: '1001001/R01',
      tersalin: 1,
      berkas: [
        { nama: 'Slip.pdf', berhasil: true, pesan: 'Tersalin' },
        { nama: 'Bordero.xlsx', berhasil: false, pesan: 'File exceeds 25 MB — not copied.' },
      ],
    }
    expect(pesanSalinanLampiran(s)).toBe(
      'Attachments from 1001001/R01 were copied too: 1 of 2 files. Not copied: Bordero.xlsx (File exceeds 25 MB — not copied.)',
    )
    expect(gabungPesanSalinan(PESAN, s).startsWith(PESAN + ' Attachments from 1001001/R01')).toBe(true)
  })

  it('galat menyeluruh dari server ditampilkan apa adanya', () => {
    const s = { sumber: '2002002', tersalin: 0, berkas: [], pesan: 'Original ID 2002002 is not the source of 1001001/R02 — its attachments were not copied.' }
    expect(pesanSalinanLampiran(s)).toBe(s.pesan)
  })

  it('semua tersalin: tanpa daftar gagal', () => {
    const s = { sumber: '1001001', tersalin: 1, berkas: [{ nama: 'a.pdf', berhasil: true, pesan: 'Tersalin' }] }
    expect(pesanSalinanLampiran(s)).toBe('Attachments from 1001001 were copied too: 1 of 1 files.')
  })
})
