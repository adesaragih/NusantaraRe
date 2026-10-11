import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { UNGGAH_EDM } from '../labels'
import { ringkasanCSV, TabelPesanCSV } from './UnggahCSV'

// Download Template + modal Upload CSV (keputusan work owner 10-10-2026).

const LAYAR = readFileSync(join(__dirname, '..', 'pages', 'InputEDMLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
const MODAL = readFileSync(join(__dirname, 'UnggahCSV.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('Download Template', () => {
  it('tombol di kiri Upload CSV, ghost sm, mati pada kondisi yang sama', () => {
    const iUnduh = LAYAR.indexOf('{UNGGAH_EDM.downloadTemplate}')
    const iUnggah = LAYAR.indexOf('{UNGGAH_EDM.uploadCsv}')
    expect(iUnduh > 0 && iUnduh < iUnggah).toBe(true)
    const tombol = LAYAR.slice(LAYAR.lastIndexOf('<button', iUnduh), iUnduh)
    expect(tombol).toContain('className="btn btn--ghost btn--sm"')
    expect(tombol).toContain('disabled={kasus.csvTerkunci || kasus.sudahSimpan}')
    expect(tombol).toContain('unduhTemplatCSV(kasusId)')
    const unggah = LAYAR.slice(LAYAR.lastIndexOf('<button', iUnggah), iUnggah)
    expect(unggah).toContain('disabled={kasus.csvTerkunci || kasus.sudahSimpan}')
    expect(UNGGAH_EDM.downloadTemplate).toBe('Download Template')
  })
})

describe('modal Upload CSV disederhanakan (keputusan work owner 10-10-2026)', () => {
  it('daftar kolom wajib tidak lagi tampil di modal (kolom sudah ada di templat CSV)', () => {
    expect(MODAL).not.toContain('Your file must include these columns')
    expect(MODAL).not.toContain('DaftarKolomWajib')
    expect(MODAL).not.toContain('edm-csv__chip')
    expect(LAYAR).not.toContain('kolomWajib=')
    expect('kolomWajib' in UNGGAH_EDM).toBe(false)
    // Yang tetap: petunjuk templat + tombol kecil pada galat 422.
    expect(MODAL).toContain('{UNGGAH_EDM.petunjukTemplat}')
  })

  it('kalimat ringkasan: ada masalah (merah) dan semua lolos (hijau); kolom asing abu bila ada', () => {
    const masalah = { total: 20, ditolak: 3, pesan: [{ baris: 2, kolom: 'SEX', pesan: 'SEX is empty. Please fill it in.' }], terpotong: false, diabaikan: ['NO'] }
    expect(ringkasanCSV(masalah)).toEqual({ kalimat: '3 of 20 rows have a problem. Fix them in your file and upload again.', masalah: true })
    const html = renderToStaticMarkup(<TabelPesanCSV hasil={masalah} />)
    expect(html).toContain('<p class="edm-csv__ringkasan edm-csv__ringkasan--masalah" role="alert">3 of 20 rows have a problem.')
    expect(html).toContain('These columns are not used and will be ignored: NO')
    const lolos = { total: 5, ditolak: 0, pesan: [], terpotong: false, diabaikan: [] }
    const h2 = renderToStaticMarkup(<TabelPesanCSV hasil={lolos} />)
    expect(h2).toContain('<p class="edm-csv__ringkasan edm-csv__ringkasan--lolos" role="status">All 5 rows are valid. Click Add CSV Data to save them.</p>')
    expect(h2).not.toContain('will be ignored')
    expect(h2).not.toContain('<table')
    // Blok lama Rows read / Rows rejected / Ignored columns tidak ada lagi.
    expect(html).not.toContain('Rows read')
    expect(html).not.toContain('Rows rejected')
  })

  it('tabel masalah dua kolom Row / Problem dengan keterangan nomor baris spreadsheet', () => {
    const hasil = { total: 2, ditolak: 1, pesan: [{ baris: 3, kolom: 'DOB', pesan: 'DOB must be a date in DD/MM/YYYY format, e.g. 22/11/1990.' }], terpotong: false, diabaikan: [] }
    const html = renderToStaticMarkup(<TabelPesanCSV hasil={hasil} />)
    expect(html).toContain('Row numbers match the rows in your spreadsheet.')
    expect(html).toContain('<thead><tr><th>Row</th><th>Problem</th></tr></thead>')
    expect(html).toContain('<td class="edm-angka">3</td><td>DOB must be a date in DD/MM/YYYY format, e.g. 22/11/1990.</td>')
    expect(html).not.toContain('<th>Column</th>')
  })

  it('galat berkas/judul (422): petunjuk templat + tombol Download Template; pesan dari server tanpa diolah', () => {
    expect(MODAL).toContain('kodeStatusGalat(galat) === 422')
    expect(MODAL).toContain('{UNGGAH_EDM.petunjukTemplat}')
    expect(MODAL).toContain('unduhTemplatCSV(kasusId)')
    expect(MODAL).not.toContain('invalid endorsement input')
    expect(UNGGAH_EDM.petunjukTemplat).toBe('Download the template to see the correct format.')
  })
})
