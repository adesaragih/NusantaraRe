// Tempat berperan tiket 05 yang dibaca LAYAR - kode SAMA dengan
// `backend/models/peran_tempat.go` (ke-12 tempat; pemetaan peran KOSONG sampai
// IAM menjawab, K12/K16). Tampil/tidaknya tiap tempat bagi pelaku dihitung
// backend (`Layar.tempat`); layar tidak menebak peran (AC 81, 82, 91).

import { POLIS, nilai, type Halaman } from './api'

/** `Section/ListSuggest` `.ProductionDate` - syarat tampil (`pyVisible`), dua identitas. */
export const TEMPAT_PRODUKSI_TAMPIL = [
  'LISTSUGGEST_PRODUCTIONDATE_TAMPIL_OPERATOR_3',
  'LISTSUGGEST_PRODUCTIONDATE_TAMPIL_OPERATOR_4',
] as const

/** `Section/DetailPoliciesNonProportional` LABEL "NON EDM" (`OperatorID.pxInsName`). */
export const TEMPAT_LABEL_NON_EDM = 'DETAILPOLICIESNONPROPORTIONAL_LABEL_NON_EDM'

/** `Section/DetailPoliciesNonProportional` LABEL "EDM" (`OperatorID.pxInsName`). */
export const TEMPAT_LABEL_EDM = 'DETAILPOLICIESNONPROPORTIONAL_LABEL_EDM'

export type Tempat = Record<string, boolean> | undefined

/** Tempat terbuka bagi pelaku (tak terdaftar / tertunda = false). */
export const tempatTerbuka = (tempat: Tempat, kode: string) => tempat?.[kode] === true

/** `.ProductionDate` tampil: `.IsApproved == 1 && (<identitas-3> || <identitas-4>)`
 *  - SAMA dengan `models.TanggalProduksiTampil`. Wajibnya dari `Layar.medanWajib`. */
export const tampilTanggalProduksi = (h: Halaman, tempat: Tempat) =>
  nilai(h, POLIS + 'IsApproved') === '1' && TEMPAT_PRODUKSI_TAMPIL.some((k) => tempatTerbuka(tempat, k))
