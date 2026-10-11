import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { ReactElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisInbox } from '../api'
import { kasusEDMBisaDibuka } from '../tampilan'
import { BarisKotakMasuk, tombolBukaKasus } from './InboxEndorsementLife'

// Seluruh baris kotak masuk membuka kasus - keputusan work owner 08-10-2026.

const SUMBER = readFileSync(join(__dirname, 'InboxEndorsementLife.tsx'), 'utf8')

const B: BarisInbox = {
  caseId: 'EDML-828', endorsementNo: '', type: 'QR', edmType: '1', policyNo: 'UJI-PL-1', ceding: 'UJI-CEDING',
  createDate: '2026-10-08 09:00:00', createOperator: 'UJI-AKUN-1', status: '',
}

type PropsTr = {
  className: string
  tabIndex: number
  'aria-label': string
  onClick: () => void
  onKeyDown: (e: { key: string; preventDefault: () => void }) => void
}

function baris(onBuka: (id: string) => void): PropsTr {
  return (BarisKotakMasuk({ b: B, onBuka }) as ReactElement<PropsTr>).props
}

describe('baris kotak masuk sebagai pembuka kasus', () => {
  it('klik baris memanggil onBuka(caseId) tepat satu kali', () => {
    const dipanggil: string[] = []
    baris((id) => dipanggil.push(id)).onClick()
    expect(dipanggil).toEqual(['EDML-828'])
  })

  it('Enter / Space membuka; tombol lain tidak; default dicegah', () => {
    const dipanggil: string[] = []
    const p = baris((id) => dipanggil.push(id))
    let dicegah = 0
    for (const key of ['Enter', ' ', 'a', 'Tab', 'Escape']) {
      p.onKeyDown({ key, preventDefault: () => dicegah++ })
    }
    expect(dipanggil).toEqual(['EDML-828', 'EDML-828'])
    expect(dicegah).toBe(2)
    expect(tombolBukaKasus('Enter') && tombolBukaKasus(' ') && !tombolBukaKasus('x')).toBe(true)
  })

  it('baris dapat difokus, berlabel Case ID, kelas inbox__baris; Case ID bukan tombol (nol pembuka ganda)', () => {
    const p = baris(() => {})
    expect(p.className).toContain('inbox__baris')
    expect(p.tabIndex).toBe(0)
    expect(p['aria-label']).toBe('Case ID EDML-828')
    const html = renderToStaticMarkup(
      <table>
        <tbody>
          <BarisKotakMasuk b={B} onBuka={() => {}} />
        </tbody>
      </table>,
    )
    expect(html).not.toContain('<button')
    // Case ID teks biasa, tanpa garis bawah (08-10-2026).
    expect(html).toContain('<td>EDML-828</td>')
    expect(html).not.toContain('edm-tautan')
    expect(html).not.toContain('underline')
    expect(html.match(/<td>/g)?.length).toBe(9)
    expect(SUMBER.split('onBuka(b.caseId)').length - 1).toBe(2) // onClick + onKeyDown baris saja
  })
})

describe('kasus tertutup tetap terdaftar, tidak dapat dibuka (09-10-2026)', () => {
  const TUTUP: BarisInbox = { ...B, caseId: 'EDML-829', endorsementNo: 'UJI-PL-1/01', status: 'Resolved-Completed' }

  it('kasusEDMBisaDibuka hanya untuk status kosong', () => {
    expect(kasusEDMBisaDibuka('')).toBe(true)
    expect(kasusEDMBisaDibuka('  ')).toBe(true)
    expect(kasusEDMBisaDibuka('Resolved-Completed')).toBe(false)
    expect(kasusEDMBisaDibuka('Resolved-Rejected')).toBe(false)
  })

  it('baris tertutup: kelas abu, title, tanpa onClick/onKeyDown/tabIndex; NO_ENDORS dan status tampil', () => {
    const el = BarisKotakMasuk({ b: TUTUP, onBuka: () => {} }) as ReactElement<Record<string, unknown>>
    expect(el.props.className).toBe('edm-inbox__tutup')
    expect(el.props.title).toBe('Case closed — cannot be opened')
    expect(el.props.onClick).toBeUndefined()
    expect(el.props.onKeyDown).toBeUndefined()
    expect(el.props.tabIndex).toBeUndefined()
    const html = renderToStaticMarkup(
      <table>
        <tbody>
          <BarisKotakMasuk b={TUTUP} onBuka={() => {}} />
        </tbody>
      </table>,
    )
    expect(html).toContain('<td>UJI-PL-1/01</td>')
    expect(html).toContain('<td>Resolved-Completed</td>')
    expect(html.match(/<td>/g)?.length).toBe(9)
  })

  it('baris terbuka tetap dapat dibuka', () => {
    expect(baris(() => {}).tabIndex).toBe(0)
  })
})
