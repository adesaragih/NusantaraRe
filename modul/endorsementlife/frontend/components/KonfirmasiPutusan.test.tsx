import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { ReactElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KONFIRMASI_PUTUSAN_EDM } from '../labels'
import KonfirmasiPutusan, { kutipanKomentar, PANJANG_KUTIPAN, teksKonfirmasiPutusan } from './KonfirmasiPutusan'

// Submit → "Confirm" + popup konfirmasi (keputusan work owner 08-10-2026).

const SUMBER = readFileSync(join(__dirname, 'KeputusanEDM.tsx'), 'utf8').replace(/\r\n/g, '\n')

type PropsModal = {
  judul: string
  labelBatal: string
  onTutup: () => void
  aksi: ReactElement<{ onClick: () => void; children: string; className: string }>
  children: unknown
}

function buat(status: string, komentar = '') {
  const panggil: string[] = []
  const el = KonfirmasiPutusan({ status, komentar, onLanjut: () => panggil.push('kirim'), onTutup: () => panggil.push('tutup') })
  return { panggil, el, p: (el as ReactElement<PropsModal>).props }
}

describe('tombol Confirm keputusan', () => {
  it('label Confirm, kelas utama, nonaktif tanpa Status atau saat mengirim; selalu tampil', () => {
    expect(KONFIRMASI_PUTUSAN_EDM.confirm).toBe('Confirm')
    const aksi = SUMBER.slice(SUMBER.indexOf('<div className="edm-aksi">'), SUMBER.indexOf('</div>', SUMBER.indexOf('<div className="edm-aksi">')))
    expect(aksi).toContain('className="btn btn--primary btn--sm"')
    expect(aksi).toContain("disabled={status === '' || mengirim}")
    expect(aksi).toContain('KONFIRMASI_PUTUSAN_EDM.confirm')
    expect(SUMBER).not.toContain("{status !== '' && (")
    expect(SUMBER).not.toContain('PUTUSAN_EDM.submit')
  })

  it('klik HANYA membuka popup; Yes satu-satunya pemanggil kirim(); kirim tetap putuskanKasus(status, comment)', () => {
    const aksi = SUMBER.slice(SUMBER.indexOf('<div className="edm-aksi">'), SUMBER.indexOf('</div>', SUMBER.indexOf('<div className="edm-aksi">')))
    expect(aksi).toContain('onClick={() => setTanya(true)}')
    expect(aksi).not.toContain('kirim()')
    expect(SUMBER.split('void kirim()').length - 1).toBe(1)
    const popup = SUMBER.slice(SUMBER.indexOf('<KonfirmasiPutusan'), SUMBER.indexOf('/>', SUMBER.indexOf('<KonfirmasiPutusan')))
    expect(popup).toContain('onTutup={() => setTanya(false)}')
    expect(popup).toContain('void kirim()')
    expect(SUMBER).toContain('await putuskanKasus(kasusId, { status, comment: komentar })')
  })
})

describe('popup konfirmasi', () => {
  it('Accept: judul/isi; Yes utama', () => {
    const { p } = buat('1')
    expect(p.judul).toBe('Confirm endorsement?')
    expect(renderToStaticMarkup(<>{p.children}</>)).toContain("This will create the endorsement number and finish the case. You can&#x27;t undo this.")
    expect(p.aksi.props.children).toBe('Yes')
    expect(p.aksi.props.className).toBe('btn btn--primary')
    expect(p.labelBatal).toBe('Cancel')
  })

  it('Decline: judul/isi', () => {
    const { p } = buat('2')
    expect(p.judul).toBe('Decline endorsement?')
    expect(teksKonfirmasiPutusan('2').isi).toBe("This will decline the endorsement and close the case. You can't undo this.")
    expect(p.aksi.props.className).toBe('btn btn--danger')
  })

  it('Cancel tidak mengirim; Yes mengirim', () => {
    const a = buat('1')
    a.p.onTutup()
    expect(a.panggil).toEqual(['tutup'])
    const b = buat('1')
    b.p.aksi.props.onClick()
    expect(b.panggil).toEqual(['kirim'])
  })

  it('kutipan Comment bila terisi (dipangkas); tanpa Comment tanpa kutipan', () => {
    expect(renderToStaticMarkup(<>{buat('1', 'UJI-KOMENTAR').p.children}</>)).toContain('UJI-KOMENTAR')
    expect(renderToStaticMarkup(<>{buat('1', '  ').p.children}</>)).not.toContain('edm-konfirmasi__kutipan')
    const panjang = 'x'.repeat(PANJANG_KUTIPAN + 10)
    expect(kutipanKomentar(panjang)).toBe(`${'x'.repeat(PANJANG_KUTIPAN)}…`)
  })
})
