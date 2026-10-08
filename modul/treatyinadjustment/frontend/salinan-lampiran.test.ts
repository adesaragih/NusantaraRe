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
        { nama: 'Bordero.xlsx', berhasil: false, pesan: 'Berkas melebihi 25 MB — tidak disalin.' },
      ],
    }
    expect(pesanSalinanLampiran(s)).toBe(
      'Lampiran 1001001/R01 ikut disalin: 1 dari 2 berkas. Tidak tersalin: Bordero.xlsx (Berkas melebihi 25 MB — tidak disalin.)',
    )
    expect(gabungPesanSalinan(PESAN, s).startsWith(PESAN + ' Lampiran 1001001/R01')).toBe(true)
  })

  it('galat menyeluruh dari server ditampilkan apa adanya', () => {
    const s = { sumber: '2002002', tersalin: 0, berkas: [], pesan: 'ID Original 2002002 bukan asal 1001001/R02 — lampirannya tidak disalin.' }
    expect(pesanSalinanLampiran(s)).toBe(s.pesan)
  })

  it('semua tersalin: tanpa daftar gagal', () => {
    const s = { sumber: '1001001', tersalin: 1, berkas: [{ nama: 'a.pdf', berhasil: true, pesan: 'Tersalin' }] }
    expect(pesanSalinanLampiran(s)).toBe('Lampiran 1001001 ikut disalin: 1 dari 1 berkas.')
  })
})
