// Tab Correspondence (tiket 56) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Korespondensi } from '../api'
import { KORESPONDENSI as K, TEKS_KORESPONDENSI as T } from '../labels'
import { BATAS_BITA_PESAN, bersihkanHtml, bitaPesan, pesanKosong } from './htmlKorespondensi'
import TabCorrespondence, { daftarNamaPenerima, galatPenerima, tanggalKorespondensi } from './TabCorrespondence'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const SUMBER = readFileSync(join(__dirname, 'TabCorrespondence.tsx'), 'utf8').replace(/\r\n/g, '\n')
const HALAMAN = readFileSync(join(__dirname, '../pages/InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe.skipIf(!existsSync(NB + 'Section\\Correspondence.xml'))('label = korpus', () => {
  it('grid (Correspondence) dan editor (CorrespondenceContent)', () => {
    const grid = readFileSync(NB + 'Section\\Correspondence.xml', 'utf-8')
    for (const t of K.kolom) expect(grid, t).toContain(`<pyValue>${t}</pyValue>`)
    for (const t of [K.tambah, K.hapus]) expect(grid, t).toContain(`<pyLabel>${t}</pyLabel>`)
    const isi = readFileSync(NB + 'Section\\CorrespondenceContent.xml', 'utf-8')
    for (const t of [K.sender, K.subject, K.message]) expect(isi, t).toContain(`>${t}<`)
    expect(isi).toContain(`<pyLabelFieldValue>${K.senderMail}</pyLabelFieldValue>`)
    expect(isi).toContain(`<pyValue>${K.penerima}</pyValue>`)
    for (const t of [K.simpan, K.kirim]) expect(isi, t).toContain(`<pyLabel>${t}</pyLabel>`)
    // Lampiran selalu tersembunyi di Pega -> tidak dibangun.
    expect(isi).toContain('<pyContainerVisibleWhen>1=2</pyContainerVisibleWhen>')
  })
  it('Pega mengirim To = pengirim, penerima di CC (dibetulkan: To = penerima)', () => {
    const xml = readFileSync(NB + 'Activity\\SendCorrespondence.xml', 'utf-8')
    expect(xml).toContain('<To>local.Email</To>')
    expect(xml).toContain('<CC>Local.Recipient</CC>')
  })
})

describe('pembersih HTML pesan', () => {
  it('tag yang diizinkan tetap, atribut dibuang', () => {
    expect(bersihkanHtml('<b style="color:red">a</b><i>b</i><u>c</u><ul><li>x</li></ul><br/><div class="k">d</div>')).toBe(
      '<b>a</b><i>b</i><u>c</u><ul><li>x</li></ul><br><div>d</div>',
    )
  })
  it('skrip, gaya, gambar, tautan, komentar dibuang; teks di-escape', () => {
    expect(bersihkanHtml('<script>alert(1)</script>ok<style>p{}</style>')).toBe('ok')
    expect(bersihkanHtml('<img src=x onerror=alert(1)>a<a href="javascript:x">b</a><!-- c -->')).toBe('ab')
    expect(bersihkanHtml('1 < 2 & 3 > 2 &nbsp; &amp; &#169; &bogus')).toBe('1 &lt; 2 &amp; 3 &gt; 2 &nbsp; &amp; &#169; &amp;bogus')
    expect(bersihkanHtml('<b title=">">x</b>')).toBe('<b>"&gt;x</b>')
    expect(bersihkanHtml('<SCRIPT>x</SCRIPT><B>y</B>')).toBe('<b>y</b>')
  })
  it('ukuran pesan = bita UTF-8 HTML bersih (batas 4000 = MESSAGE VARCHAR2(4000) backend)', () => {
    expect(BATAS_BITA_PESAN).toBe(4000)
    expect(bitaPesan('<b style="x">é</b>')).toBe('<b>é</b>'.length + 1)
    expect(bitaPesan('a'.repeat(4001)) > BATAS_BITA_PESAN).toBe(true)
    expect(SUMBER).toContain('disabled={sibuk !== null || kepanjangan}')
  })
  it('pesan kosong', () => {
    expect(pesanKosong('<div><br></div>&nbsp;')).toBe(true)
    expect(pesanKosong('<b>a</b>')).toBe(false)
  })
})

describe('TabCorrespondence - aturan dan tampilan', () => {
  it('tanggal "06-10-2026 17:28" -> "06/10/26 17:28" (gambar layar Pega)', () => {
    expect(tanggalKorespondensi('06-10-2026 17:28')).toBe('06/10/26 17:28')
    expect(tanggalKorespondensi('')).toBe('')
  })
  it('penerima: wajib minimal satu alamat, semua alamat sah', () => {
    expect(galatPenerima([])).toBe(T.penerimaKosong)
    expect(galatPenerima([{ recipient: 'UJI', recipientMail: '' }])).toBe(T.penerimaKosong)
    expect(galatPenerima([{ recipient: 'UJI', recipientMail: 'uji@contoh.test' }])).toBe('')
    expect(galatPenerima([{ recipient: '', recipientMail: 'uji@contoh' }])).toBe(`${T.emailSalah}uji@contoh`)
    expect(daftarNamaPenerima([{ recipient: 'UJI-A', recipientMail: 'a@contoh.test' }, { recipient: '', recipientMail: 'b@contoh.test' }])).toBe(
      'UJI-A, b@contoh.test',
    )
  })
  it('render grid: kolom, status Not Sent / Sent, Remove hanya untuk draf', () => {
    const k = (no: string, terkirim: boolean): Korespondensi => ({
      no, tanggal: '06-10-2026 17:28', sender: 'UJI-PENGIRIM', senderMail: 'pengirim@contoh.test', subject: 'UJI-SUBJEK',
      message: '<b>halo</b>', penerima: [{ recipient: 'UJI-PENERIMA', recipientMail: 'p@contoh.test' }], terkirim,
    })
    const html = renderToStaticMarkup(<TabCorrespondence caseId="NB-UJI" awal={{ daftar: [k('1', false), k('2', true)], pesan: [] }} />)
    for (const t of K.kolom) expect(html).toContain(`>${t}<`)
    expect(html).toContain('>06/10/26 17:28<')
    expect(html).toContain(`>${K.belumTerkirim}<`)
    expect(html).toContain(`>${K.terkirim}<`)
    expect(html.match(new RegExp(`${K.hapus}<`, 'g'))).toHaveLength(1)
    expect(html).toContain('>UJI-PENERIMA<')
  })
  it('tab hanya untuk POSITION IT (work owner 06-10-2026; Pega S95 selalu tampil); simpan / kirim lewat backend', () => {
    expect(HALAMAN).toContain("(t !== 'Correspondence' || dariServer?.aksesKorespondensi === true)")
    expect(HALAMAN).toContain("tab === 'Correspondence' && dariServer?.aksesKorespondensi === true ? (\n              <TabCorrespondence caseId={kasus.caseId} />")
    expect(SUMBER).toContain('ganti(await kirimKorespondensi(caseId, k.no, isiDari(k)))')
    expect(SUMBER).toContain('ganti(await simpanKorespondensi(caseId, k.no, isiDari(k)))')
  })
})
