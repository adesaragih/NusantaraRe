// Popup Historical Survey Report - `Harness/HistoricalSurveyReport` (admin, `Section/HistoricalSurveyReportDtl`)
// dan `Harness/HistoricalSurveyReportUW` (atasan, `Section/HistoricalSurveyReportDtlUW`, hanya-baca). Grid
// `.QuotationData.SurveyReportList` (keputusan work owner 06-10-2026, disimpan T_POLIS_SURVEY):
//
//   admin   Insured Name hanya-baca; Add (addRow), Delete (deleteRow) per baris, Submit (`SetSurveyReport_Act`:
//           nol Obj-Save) - hasil DIPEGANG layar dan tersimpan lewat Save, sama dengan Pega. Baris disunting
//           langsung di grid (pengganti form masterDetail `InputHistoricalSurveyReportDtl`, medannya sama).
//   atasan  hanya-baca, tanpa kolom Loss Prevention, tanpa Add / Delete / Submit.
//
// `.DateofSurvey` wajib (`InputHistoricalSurveyReportDtl`): Submit ditolak selama ada baris tanpa tanggal.

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { daftar, nilai, POLIS, type Baris, type Halaman } from '../api'
import { JUDUL, KOLOM_SURVEI, PESAN, PILIHAN_REMARKS_SURVEI, TOMBOL } from '../labels'
import { teksPilihan } from '../medan'
import { sajikan } from '../sajian'
import { barisTanpaTanggal, DAFTAR_SURVEI, SAJIAN_LOSS_PREVENTION } from '../survei'
import InputAngka from './InputAngka'

/** Pilihan Remarks; nilai tersimpan di luar daftar ditambahkan sebagai pilihan sendiri (tidak dibuang). */
function opsiRemarks(v: string) {
  return v !== '' && !PILIHAN_REMARKS_SURVEI.some((o) => o.value === v)
    ? [...PILIHAN_REMARKS_SURVEI, { value: v, label: v }]
    : PILIHAN_REMARKS_SURVEI
}

export default function SurveiHistoris({
  halaman,
  admin,
  sunting,
  onSetel,
  onTutup,
}: {
  halaman: Halaman
  /** Layar admin (kolom Loss Prevention ada); atasan memakai versi UW. */
  admin: boolean
  /** Admin yang boleh menyunting: Add / Delete / Submit dan isian grid. */
  sunting: boolean
  onSetel: (baris: Baris[]) => void
  onTutup: () => void
}) {
  const [baris, setBaris] = useState<Baris[]>(() => daftar(halaman, DAFTAR_SURVEI))
  const [kurang, setKurang] = useState<number[]>([])
  const ubah = (i: number, k: string, v: string) => setBaris((b) => b.map((r, j) => (j === i ? { ...r, [k]: v } : r)))
  const kirim = () => {
    const n = barisTanpaTanggal(baris)
    setKurang(n)
    if (n.length === 0) {
      onSetel(baris)
      onTutup()
    }
  }
  const kolom = 3 + (admin ? 1 : 0) + (sunting ? 1 : 0)

  return (
    <Modal
      judul={JUDUL.surveiHistoris}
      onTutup={onTutup}
      lebar
      aksi={
        sunting ? (
          <button type="button" className="btn btn--primary" onClick={kirim}>
            {TOMBOL.submit}
          </button>
        ) : undefined
      }
    >
      <div className="nbti__kolom nbti__survei-kepala">
        <div className="field nbti__medan">
          <span className="field__label">{KOLOM_SURVEI.tertanggung}</span>
          <span className="nbti__nilai">{nilai(halaman, POLIS + 'InsuredName')}</span>
        </div>
      </div>
      {kurang.length > 0 && (
        <div className="alert alert--warn">
          {PESAN.surveiTanpaTanggal} {kurang.join(', ')}.
        </div>
      )}
      {sunting && (
        <div className="nbti__aksi">
          <button type="button" className="btn btn--sm" onClick={() => setBaris((b) => [...b, {}])}>
            {TOMBOL.add}
          </button>
        </div>
      )}
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th scope="col">
                {KOLOM_SURVEI.tanggal}
                {sunting && <span className="field__req">*</span>}
              </th>
              <th scope="col">{KOLOM_SURVEI.oleh}</th>
              {admin && <th scope="col" className="nbti__angka">{KOLOM_SURVEI.lossPrevention}</th>}
              <th scope="col">{KOLOM_SURVEI.remarks}</th>
              {sunting && <th scope="col" />}
            </tr>
          </thead>
          <tbody>
            {baris.map((b, i) => (
              <tr key={i}>
                <td>
                  {sunting ? (
                    <input
                      className="field__input"
                      type="date"
                      aria-label={KOLOM_SURVEI.tanggal}
                      value={(b.DateofSurvey ?? '').slice(0, 10)}
                      onChange={(e) => ubah(i, 'DateofSurvey', e.target.value)}
                    />
                  ) : (
                    sajikan(b.DateofSurvey ?? '', 'tanggal')
                  )}
                </td>
                <td>
                  {sunting ? (
                    <input
                      className="field__input"
                      aria-label={KOLOM_SURVEI.oleh}
                      value={b.SurveyedBy ?? ''}
                      onChange={(e) => ubah(i, 'SurveyedBy', e.target.value)}
                    />
                  ) : (
                    (b.SurveyedBy ?? '')
                  )}
                </td>
                {admin && (
                  <td className="nbti__angka">
                    {sunting ? (
                      <InputAngka
                        label={KOLOM_SURVEI.lossPrevention}
                        value={b.LossPrevention ?? ''}
                        sajian={SAJIAN_LOSS_PREVENTION}
                        onChange={(v) => ubah(i, 'LossPrevention', v)}
                        onBlur={() => undefined}
                      />
                    ) : (
                      sajikan(b.LossPrevention ?? '', SAJIAN_LOSS_PREVENTION)
                    )}
                  </td>
                )}
                <td>
                  {sunting ? (
                    <select
                      className="field__input"
                      aria-label={KOLOM_SURVEI.remarks}
                      value={b.Remarks ?? ''}
                      onChange={(e) => ubah(i, 'Remarks', e.target.value)}
                    >
                      <option value="">{TOMBOL.pilihKosong}</option>
                      {opsiRemarks(b.Remarks ?? '').map((o) => (
                        <option key={o.value} value={o.value}>
                          {o.label}
                        </option>
                      ))}
                    </select>
                  ) : (
                    teksPilihan(PILIHAN_REMARKS_SURVEI, b.Remarks ?? '')
                  )}
                </td>
                {sunting && (
                  <td>
                    <button
                      type="button"
                      className="btn btn--sm btn--danger"
                      onClick={() => setBaris((x) => x.filter((_, j) => j !== i))}
                    >
                      {TOMBOL.delete}
                    </button>
                  </td>
                )}
              </tr>
            ))}
            {baris.length === 0 && (
              <tr>
                <td colSpan={kolom} className="muted">
                  —
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </Modal>
  )
}
