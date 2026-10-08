// Nilai layar Komite Claim Prop - fungsi murni (diuji tanpa DOM). Tampilan angka dan tanggal DISALIN dari pola Claim
// Prop (`modul/claimprop/frontend/nilai.ts`, `ketikTanggal.ts`), tidak diimpor (batas modul): angka 4 desimal
// (`pyDecimalPlaces 4`) pemisah Indonesia, tanggal `dd-mm-yyyy` (+ ` hh:mm`).

import { formatNumber } from '../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../inti/frontend/lib/tanggalInput'
import type { JenisNilai, Keputusan } from './api'

/** Desimal tampilan angka. */
export const DESIMAL_TAMPIL = 4

/** Nilai halaman -> tampilan menurut jenisnya. Tidak terbaca = apa adanya. */
export function tampil(v: string, jenis: JenisNilai | ''): string {
  const s = (v ?? '').trim()
  if (s === '') return ''
  switch (jenis) {
    case 'angka':
      return formatNumber(s, DESIMAL_TAMPIL)
    case 'tanggal':
    case 'tanggalJam': {
      const iso = keInputTanggal(s)
      if (iso === '') return s
      const dmy = `${iso.slice(8, 10)}-${iso.slice(5, 7)}-${iso.slice(0, 4)}`
      if (jenis === 'tanggal') return dmy
      const jam = /[T ](\d{2}):(\d{2})/.exec(s)
      return jam ? `${dmy} ${jam[1]}:${jam[2]}` : dmy
    }
    default:
      return s
  }
}

/** Kode `.AcceptStatus` (`SetDataAcceptationTreaty_Act` S2-S3: 1 Approve, 2 Reject). */
export const KEPUTUSAN = { setuju: '1', tolak: '2' } as const

/** Isian Subjectivity tampil (`pyVisible` `.AcceptStatus = 1 && .TransferType = 2`; layar ini hanya TT 2). */
export function tampilSubjectivity(k: Keputusan): boolean {
  return k.acceptStatus === KEPUTUSAN.setuju
}

/** Subjectivity Note tampil dan wajib (`.IsSubjectivity = true`). */
export function tampilCatatanSubjectivity(k: Keputusan): boolean {
  return tampilSubjectivity(k) && k.isSubjectivity
}

/** Pesan wajib isi bawaan Pega. */
export const PESAN_KOSONG = 'Value cannot be blank'

/**
 * Validasi klien (pyRequired / pyRequiredWhen) - sama dengan `models.PeriksaIsian` di server; server tetap menegakkan.
 * Kunci = medan isian.
 */
export function periksa(k: Keputusan, terbuka: boolean): Partial<Record<keyof Keputusan, string>> {
  const out: Partial<Record<keyof Keputusan, string>> = {}
  if (k.acceptStatus !== KEPUTUSAN.setuju && k.acceptStatus !== KEPUTUSAN.tolak) out.acceptStatus = PESAN_KOSONG
  if (k.comment.trim() === '') out.comment = PESAN_KOSONG
  if (terbuka && tampilCatatanSubjectivity(k) && k.subjectivityNote.trim() === '') out.subjectivityNote = PESAN_KOSONG
  return out
}

/** Isian yang dikirim: isian bertingkat-1 yang nonaktif / tersembunyi tidak membawa nilai baru. */
export function isianKirim(k: Keputusan, terbuka: boolean): Keputusan {
  const subj = terbuka && tampilSubjectivity(k) && k.isSubjectivity
  return {
    acceptStatus: k.acceptStatus,
    comment: k.comment,
    isSubjectivity: subj,
    subjectivityNote: subj ? k.subjectivityNote : '',
    usulTutup: terbuka ? k.usulTutup : false,
    usulCadang: terbuka ? k.usulCadang : false,
  }
}
