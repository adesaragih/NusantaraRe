// Tombol dan popup Historical Survey Report (keputusan work owner 06-10-2026, membatalkan K7):
//
//   tombol  `Section/DetailPolicyTreatyIn` sel 22 / `DetailDeptHeadTreatyIn_UW`, LABEL "Survey Report", showHarness
//           popup `HistoricalSurveyReport` / `HistoricalSurveyReportUW` (pyWindowName "Historical Survey Report")
//           pyVisible       `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`
//           pyDisabledWhen  `.QuotationData.IsSurveyReport=='No' || .QuotationData.IsSurveyReport==''`
//   grid    `.QuotationData.SurveyReportList` - disimpan T_POLIS_SURVEY; `.DateofSurvey` wajib
//           (`Section/InputHistoricalSurveyReportDtl`)

import { nilai, POLIS, type Baris, type Halaman } from './api'
import { bukanNonProp } from './medan'
import type { Sajian } from './sajian'

export const DAFTAR_SURVEI = POLIS + 'QuotationData.SurveyReportList'

export const tampilTombolSurvei = bukanNonProp

export function aktifTombolSurvei(h: Halaman): boolean {
  const v = nilai(h, POLIS + 'QuotationData.IsSurveyReport')
  return v !== 'No' && v !== ''
}

/** Nomor baris (1..n) yang belum berisi Date of Survey. */
export function barisTanpaTanggal(baris: Baris[]): number[] {
  return baris.flatMap((b, i) => ((b.DateofSurvey ?? '').trim() === '' ? [i + 1] : []))
}

/** Sel `.LossPrevention` pxNumber `pyDecimalPlaces` 4. */
export const SAJIAN_LOSS_PREVENTION: Sajian = { desimal: 4 }
