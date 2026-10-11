// Tab Clauses tampilan polis NB (FIRE): baris T_CLAUSELIST (satu per argumen) dikelompokkan per SEQ_NO seperti nbfacin
// `kumpulkanKlausa`; halaman Pega (satu baris per klausa, tanpa SeqNo) tetap satu kartu per baris.

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { adaArgumenNB, klausaNB } from './polisNBKlausa'

const CL = 'OfferFacIn.ClauseList'
const h = (baris: Record<string, string>[]) => ({ nilai: {}, daftar: { [CL]: baris } }) as unknown as Halaman

describe('klausaNB', () => {
  it('tabel datar: SEQ_NO sama = satu klausa; kode sama ber-SEQ_NO beda tetap dua kartu', () => {
    const k = klausaNB(
      h([
        { SeqNo: '1', ClauseCode: 'UJI-A', ClauseTitle: 'UJI JUDUL A', ArgumentCount: '2' },
        { SeqNo: '1', ClauseCode: 'UJI-A', ClauseTitle: 'UJI JUDUL A', ArgumentCount: '2' },
        { SeqNo: '2', ClauseCode: 'UJI-B', ClauseLanguageID: 'Inggris' },
        { SeqNo: '3', ClauseCode: 'UJI-B' },
      ]),
    )
    expect(k.map((x) => x.clauseCode)).toEqual(['UJI-A', 'UJI-B', 'UJI-B'])
    expect(k[0]).toMatchObject({ clauseTitle: 'UJI JUDUL A', argumentCount: '2' })
    expect(k[1]!.clauseLanguageId).toBe('Inggris')
  })

  it('halaman Pega tanpa SeqNo: satu kartu per baris', () => {
    expect(klausaNB(h([{ ClauseCode: 'UJI-A' }, { ClauseCode: 'UJI-A' }]))).toHaveLength(2)
    expect(klausaNB(h([]))).toEqual([])
  })

  it('Total Argument tampil bila bilangan bulat > 0 (`adaArgumen` nbfacin)', () => {
    expect(adaArgumenNB('2')).toBe(true)
    expect(adaArgumenNB('010')).toBe(true)
    expect(adaArgumenNB('0')).toBe(false)
    expect(adaArgumenNB('')).toBe(false)
    expect(adaArgumenNB('1.5')).toBe(false)
  })
})
