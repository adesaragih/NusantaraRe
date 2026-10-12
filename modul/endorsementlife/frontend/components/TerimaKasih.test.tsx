import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { ReactElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { HasilPutusanEDM } from '../api'
import { HASIL_EDM } from '../labels'
import TerimaKasih, { barisHasil } from './TerimaKasih'

// Popup hasil Confirm - permintaan work owner 09-10-2026.

const LAYAR = readFileSync(join(__dirname, '..', 'pages', 'InputEDMLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
const KEPUTUSAN = readFileSync(join(__dirname, 'KeputusanEDM.tsx'), 'utf8').replace(/\r\n/g, '\n')

const HASIL: HasilPutusanEDM = { status: 'Resolved-Completed', noEndorsement: 'RNML-QRL2.10.26.00042/01', peserta: 9, rekapWarisan: 1 }
const KASUS = { id: 'EDML-828', edmType: '1', edmTypePerubahanData: 'UJI-SUBJENIS', edmTypeBatal: '' }

type PropsModal = { judul: string; labelBatal: string; onTutup: () => void; children: unknown }

function props(h: HasilPutusanEDM, k: typeof KASUS, onTutup: () => void = () => {}) {
  return (TerimaKasih({ hasil: h, kasus: k, onTutup }) as ReactElement<PropsModal>).props
}

describe('popup hasil Confirm', () => {
  it('judul, nomor, EDM Type, Sub-type, Case ID; tombol OK', () => {
    const p = props(HASIL, KASUS)
    expect(p.judul).toBe('Endorsement confirmed')
    expect(p.labelBatal).toBe('OK')
    expect(barisHasil(HASIL, KASUS)).toEqual([
      [HASIL_EDM.noEndorsement, 'RNML-QRL2.10.26.00042/01'],
      [HASIL_EDM.edmType, 'Perubahan Data'],
      [HASIL_EDM.subType, 'UJI-SUBJENIS'],
      [HASIL_EDM.caseId, 'EDML-828'],
    ])
    const html = renderToStaticMarkup(<>{p.children}</>)
    expect(html).toContain('<dt>Endorsement No</dt><dd>RNML-QRL2.10.26.00042/01</dd>')
    expect(html).toContain('<dt>Case ID</dt><dd>EDML-828</dd>')
  })

  it('Sub-type hanya bila terisi; Batal memakai sub-jenis Batal', () => {
    expect(barisHasil(HASIL, { ...KASUS, edmTypePerubahanData: '' }).map(([l]) => l)).not.toContain(HASIL_EDM.subType)
    const batal = barisHasil(HASIL, { ...KASUS, edmType: '3', edmTypePerubahanData: 'X', edmTypeBatal: 'UJI-BATAL' })
    expect(batal).toContainEqual([HASIL_EDM.edmType, 'Batal'])
    expect(batal).toContainEqual([HASIL_EDM.subType, 'UJI-BATAL'])
  })

  it('nomor kosong → "—" (popup tetap tampil)', () => {
    expect(barisHasil({ ...HASIL, noEndorsement: '' }, KASUS)[0]).toEqual([HASIL_EDM.noEndorsement, '—'])
  })

  it('OK/X/Esc → onTutup sekali; kotak masuk hanya lewat penutupan popup', () => {
    let n = 0
    props(HASIL, KASUS, () => n++).onTutup()
    expect(n).toBe(1)
    expect(LAYAR).toContain('<TerimaKasih hasil={putusan} kasus={kasus} onTutup={onTutup} />')
  })

  it('Accept → popup; Decline → langsung kotak masuk tanpa nomor', () => {
    expect(LAYAR).toContain("onSelesai={(h, status) => (status === '1' ? setPutusan(h) : onTutup())}")
    expect(KEPUTUSAN).toContain('onSelesai(await putuskanKasus(kasusId, { status, comment: komentar }), status)')
  })
})
