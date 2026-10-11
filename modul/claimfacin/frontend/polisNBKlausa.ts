// Tampilan polis NB Fac In HANYA-BACA - tab Clauses FIRE (perintah work owner 11-10-2026 "tampilan view polis yang belum
// tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Rupa = kartu klausa `modul/nbfacin/frontend/components/
// TabClauses.tsx` (pola, bukan impor): Clause Code, bahasa, Total Argument (bila > 0), Title, Description. Tombol Choose
// Clause, See Clause and Argument, hapus, dan Save tidak dibawa. Label VERBATIM `modul/nbfacin/frontend/labels.ts` `KLAUSA`.
//
// T_CLAUSELIST menyimpan SATU BARIS PER ARGUMEN (kolom klausa berulang, nbfacin migrasi 197); baris ber-SEQ_NO sama = satu
// klausa (nbfacin repository/klausa.go `kumpulkanKlausa`). Halaman Pega (`.ClauseList` satu baris per klausa) tanpa SeqNo.

import type { Halaman } from './api'
import { AKAR_NB, daftarNB } from './polisNB'

/** `KLAUSA` (label kartu). */
export const KLAUSA_NB = {
  clauseCode: 'Clause Code',
  totalArgumen: 'Total Argument',
  title: 'Title',
  description: 'Description',
} as const
/** `TEKS_KLAUSA.kosong` tanpa ajakan Choose Clause (tampilan hanya-baca). */
export const TEKS_KLAUSA_NB = { kosong: 'Belum ada klausa.' } as const

export interface KlausaNB {
  clauseCode: string
  clauseTitle: string
  clauseDescription: string
  clauseLanguageId: string
  argumentCount: string
}

/** Kartu klausa polis: baris berurutan ber-SeqNo sama dilipat (baris pertama mewakili); tanpa SeqNo satu kartu per baris. */
export function klausaNB(h: Halaman): KlausaNB[] {
  const out: KlausaNB[] = []
  let seqAkhir = ''
  for (const b of daftarNB(h, AKAR_NB + '.ClauseList')) {
    const seq = (b['SeqNo'] ?? '').trim()
    if (seq !== '' && seq === seqAkhir) continue
    seqAkhir = seq
    out.push({
      clauseCode: b['ClauseCode'] ?? '',
      clauseTitle: b['ClauseTitle'] ?? '',
      clauseDescription: b['ClauseDescription'] ?? '',
      clauseLanguageId: b['ClauseLanguageID'] ?? '',
      argumentCount: b['ArgumentCount'] ?? '',
    })
  }
  return out
}

/** Total Argument tampil bila > 0 (`adaArgumen` nbfacin, `.ArgumentCount > 0` sel 5) - teks bilangan bulat. */
export const adaArgumenNB = (n: string) => /^0*[1-9][0-9]*$/.test(n.trim())
