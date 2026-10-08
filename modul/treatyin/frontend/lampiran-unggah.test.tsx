// Panel Attachment — tombol `Upload file` (`Section/WorkAttachments.xml`),
// HIDUP sejak 8 Oktober 2026: unggahan masuk `M_ATTACHMENTTREATY_2`.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import PanelLampiran from './components/PanelLampiran'
import type { BarisKategoriLampiran } from './api'

const KATEGORI: BarisKategoriLampiran[] = [
  { kode: '00002', nama: 'Approval Email', cacah: 1, dipastikan: true },
  { kode: '00008', nama: 'Letter of Acknowledgment / LOA', cacah: 0, dipastikan: true },
]
const tombolUpload = (html: string) => [...html.matchAll(/<button[^>]*>Upload<\/button>/g)].map((m) => m[0])

describe('⭐ Upload file — syarat tampil ekspor', () => {
  it('ViewState != 1 / RevisionState = 1 (bisaUnggah) + kontrak ber-ID: tombol HIDUP per kategori', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" bisaUnggah />)
    const t = tombolUpload(html)
    expect(t).toHaveLength(2)
    for (const b of t) expect(b).not.toContain('disabled')
  })

  it('mode lihat tanpa revisi: tombol TIDAK dirender', () => {
    expect(tombolUpload(renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" />))).toHaveLength(0)
  })

  it('kontrak belum tersimpan: tombol mati, dengan alasannya', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="" bisaUnggah />)
    for (const b of tombolUpload(html)) expect(b).toContain('disabled')
    expect(html).toContain('Simpan kontrak lebih dulu untuk mengunggah lampiran.')
  })

  it('Refresh tampil; Download All (`pyCondition never`) tidak', () => {
    const html = renderToStaticMarkup(<PanelLampiran kategori={KATEGORI} berkas={[]} idKontrak="1002305" />)
    expect(html).toContain('>Refresh</button>')
    expect(html).not.toContain('Download All')
  })
})
