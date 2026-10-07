// Aturan popup "Choose Cause of decline" (tiket 57) - port `NB FacIn\Section\InputCauseOfDecline.xml` (kolom aktif),
// `Activity\DisableCauseOfDecline.xml` (satu pilihan per prioritas), `CekCauseOfDecline_ACT.xml` (validasi, pesan verbatim)
// dan `SetCauseOfDecline_ACT.xml` (teks alasan -> Comment). Backend menerapkan aturan yang sama saat Submit.
//
// Keputusan agent (tiket 57): K57-3 memilih Non - T/C mengosongkan pilihan T/C (Pega hanya menonaktifkan kontrolnya,
// centang lama ikut tersimpan lalu alasannya ditimpa); K57-4 semua pesan validasi ditampilkan (Pega menimpa `Local.Errmsg`
// sehingga hanya pesan terakhir yang tampil); K57-5 wajib minimal prioritas I atau satu Non - T/C (Pega: tombol Decline
// nonaktif selama `ReasonOfDecline` kosong - yang hanya terisi dari prioritas I / Non - T/C).

import type { PilihanSebabTc } from '../api'
import { SEBAB_DECLINE as S, TEKS_SEBAB_DECLINE as T } from '../labels'

const termAda = (daftar: readonly string[], term: string) => daftar.includes(term)

/** Kolom Option aktif (pyDisabledWhen `.TERMCONDITION != POI / Extension of coverage / Risk Exposure`). */
export const optionAktif = (term: string) => termAda(S.term.option, term)
/** Kolom Remarks aktif (Occupancy / TSI (in IDR) / Clauses / Lack of data for acceptance / Others T/C, atau Option = Others). */
export const remarksAktif = (term: string, description1: string) => termAda(S.term.remarks, term) || description1 === S.term.lainnya
/** Kolom % RNM aktif. */
export const pctRnmAktif = (term: string) => termAda(S.term.pctRnm, term)
/** Kolom % SOB aktif. */
export const pctSobAktif = (term: string) => termAda(S.term.pctSob, term)

/** Pilihan kosong untuk satu T/C. */
export const pilihanKosong = (termCondition: string): PilihanSebabTc => ({
  termCondition,
  prioritas: '',
  description1: '',
  description2: '',
  pctRnm: '',
  pctSob: '',
})

/**
 * Centang / lepas satu prioritas pada satu baris: satu baris per prioritas (DisableCauseOfDecline), satu baris tidak boleh
 * I dan II sekaligus. Melepas mengosongkan isian baris itu. `pctOtomatis` = % RNM bawaan (SetCODRate) untuk prioritas I.
 */
export function centangPrioritas(daftar: PilihanSebabTc[], idx: number, p: 'I' | 'II', pctOtomatis = ''): PilihanSebabTc[] {
  const kini = daftar[idx]
  if (!kini) return daftar
  if (kini.prioritas === p) return daftar.map((x, i) => (i === idx ? pilihanKosong(x.termCondition) : x))
  return daftar.map((x, i) => {
    if (i === idx) return { ...x, prioritas: p, pctRnm: p === 'I' && pctOtomatis !== '' ? pctOtomatis : x.pctRnm }
    return x.prioritas === p ? pilihanKosong(x.termCondition) : x
  })
}

/** Validasi `CekCauseOfDecline_ACT` atas baris T/C yang dicentang; pesan verbatim, unik, berurutan. */
export function validasiSebab(daftar: PilihanSebabTc[], nonTc: string): string[] {
  const pesan: string[] = []
  const tambah = (p: string) => {
    if (!pesan.includes(p)) pesan.push(p)
  }
  for (const b of daftar) {
    if (b.prioritas === '') continue
    const t = b.termCondition
    if ((t === S.term.rate || t === S.term.riComm || t === 'OJK Discount') && (b.pctRnm === '' || b.pctSob === '')) tambah(S.pesan.pctRnmSob)
    if (t === S.term.lossRatio && b.pctRnm === '') tambah(S.pesan.pctRnm)
    if (termAda(S.term.remarksWajib, t) && b.description2 === '') tambah(S.pesan.remaks)
    if (termAda(S.term.option, t) && b.description1 === '') tambah(S.pesan.option)
    if ((t === 'Extension of coverage' || t === 'Risk Exposure') && b.description1 === S.term.lainnya && b.description2 === '') tambah(S.pesan.remarks)
  }
  if (nonTc === '' && !daftar.some((b) => b.prioritas === 'I')) tambah(T.kosong)
  return pesan
}

/** Teks satu baris T/C (`SetCauseOfDecline_ACT` 5.4): TERM + " " + Option + " " + Remarks + " " + %RNM[" % RNM "] + %SOB[" % SOB "]. */
export function teksSebabTc(b: PilihanSebabTc): string {
  return `${b.termCondition} ${b.description1} ${b.description2} ${b.pctRnm}${b.pctRnm !== '' ? ' % RNM ' : ''}${b.pctSob}${b.pctSob !== '' ? ' % SOB ' : ''}`
}

/** Comment Decline = ReasonOfDecline + "\n" + ReasonOfDecline2 (langkah 8); Non - T/C menimpa ReasonOfDecline (langkah 7). */
export function komentarDecline(daftar: PilihanSebabTc[], nonTc: string): string {
  const satu = daftar.find((b) => b.prioritas === 'I')
  const dua = daftar.find((b) => b.prioritas === 'II')
  const alasan1 = nonTc !== '' ? nonTc : satu ? teksSebabTc(satu) : ''
  const alasan2 = dua ? teksSebabTc(dua) : ''
  return `${alasan1}\n${alasan2}`
}
