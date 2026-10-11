// Copy Spreading PER LOKASI (tiket 73 P5, K73-40/50/51) - port `NB FacIn\Section\PropertyItemListCoverageSpreading_IsUW.xml`
// kontainer `(IsUWSpread)` L659: grid template `SpreadingListLoc.pxResults` (kolom "Type Treaty" L1571 · "% Share" L1735 · "Type
// Facultative" L1876), tombol Add L2025 / Delete L3290 / "Copy To All Spreading" L4444 (ketiganya `pyVisible=ALWAYS`). Dipakai tab
// Spreading FIRE dan Aneka di detail tiap lokasi, hanya tahap UW.
//
// - Type Facultative `.FlagSpreading` (pxDropdown L2934): opsi properti `DDL\FlagSpreading.xml`
//   (ASM-FW-GISFW-DATA-SPREADINGRISK!FLAGSPREADING) "SF-MR 2023" / "Open Fac Market"; nonaktif bila `.TreatyType!="10015"` (L3014),
//   wajib bila `=="10015"` (L3020). Tidak disimpan (K73-51).
// - Type Treaty berubah → `ResetFlagSpreading_ACT` (L2464 / L2724): Type Facultative SEMUA baris dikosongkan (L290).
// - Copy To All Spreading → `POST …/salin-lokasi` (`CopyToAllLocSpreading_ACT`): baris spreading lokasi ini diganti template; tidak
//   menyimpan. Efek IsFlagUW tidak diport (K73-50).

import { useState } from 'react'

import { Field, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import type { TemplateSpreadingLokasi, TreatySpreading } from '../api'
import { OPSI_FLAG_SPREADING, SPREADING_LOKASI as L, TEKS_SPREADING as T } from '../labels'
import { templateLebih, treatyKembar } from './TabSpreading'

/** ID treaty yang mewajibkan Type Facultative (`PropertyItemListCoverageSpreading_IsUW.xml` L3014 / L3020). */
export const TREATY_WAJIB_FAKULTATIF = '10015'

const angkaSah = (v: string) => v.trim() === '' || desimalSah(v.trim())
const barisKosong = (): TemplateSpreadingLokasi => ({ treatyType: '', treatyName: '', sharePercentage: '', flagSpreading: '' })

/** Baris yang wajib Type Facultative tetapi kosong. */
export const fakultatifKurang = (t: TemplateSpreadingLokasi[]) =>
  t.some((x) => x.treatyType.trim() === TREATY_WAJIB_FAKULTATIF && x.flagSpreading.trim() === '')

export default function CopySpreadingLokasi({
  treaty,
  menghitung,
  onSalin,
}: {
  treaty: TreatySpreading[]
  menghitung: boolean
  onSalin: (template: TemplateSpreadingLokasi[]) => void
}) {
  const [template, setTemplate] = useState<TemplateSpreadingLokasi[]>([])
  const opsiTreaty = treaty.map((t) => ({ value: t.id, label: t.name }))
  const opsiFlag = OPSI_FLAG_SPREADING.map((v) => ({ value: v, label: v }))
  const lebih = templateLebih(template)
  const kembar = treatyKembar(template)
  const kurang = fakultatifKurang(template)
  const angkaSalah = template.some((x) => !angkaSah(x.sharePercentage))

  return (
    <div className="nbf-spread__copy nbf-spread__copy--lokasi">
      <h5 className="nbf-akum__kartu-judul">{L.judul}</h5>
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {L.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setTemplate([...template, barisKosong()])}>
                  {L.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {template.length === 0 && (
              <tr>
                <td colSpan={4}>{T.kosong}</td>
              </tr>
            )}
            {template.map((t, i) => (
              <tr key={i}>
                <td>
                  <Pilih
                    label={L.kolom[0]}
                    value={t.treatyType}
                    kosong={T.pilihTreaty}
                    opsi={opsiTreaty}
                    onChange={(id) =>
                      // ResetFlagSpreading_ACT: Type Facultative semua baris dikosongkan.
                      setTemplate(
                        template.map((x, k) => ({
                          ...x,
                          flagSpreading: '',
                          ...(k === i ? { treatyType: id, treatyName: treaty.find((y) => y.id === id)?.name ?? '' } : {}),
                        })),
                      )
                    }
                  />
                </td>
                <td>
                  <Field
                    label={L.kolom[1]}
                    value={t.sharePercentage}
                    onChange={(s) => setTemplate(template.map((x, k) => (k === i ? { ...x, sharePercentage: s } : x)))}
                    error={angkaSah(t.sharePercentage) ? undefined : T.angka}
                  />
                </td>
                <td>
                  {t.treatyType.trim() === TREATY_WAJIB_FAKULTATIF ? (
                    <Pilih
                      label={L.kolom[2]}
                      value={t.flagSpreading}
                      kosong={T.pilihTreaty}
                      opsi={opsiFlag}
                      required
                      error={t.flagSpreading.trim() === '' ? L.fakultatifWajib : undefined}
                      onChange={(f) => setTemplate(template.map((x, k) => (k === i ? { ...x, flagSpreading: f } : x)))}
                    />
                  ) : (
                    <span className="muted">—</span>
                  )}
                </td>
                <td className="table__actions">
                  <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setTemplate(template.filter((_, k) => k !== i))}>
                    {L.hapus}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {lebih && <div className="field__error">{T.shareLebih}</div>}
      {kembar && <div className="field__error">{T.treatySama}</div>}
      <div>
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary btn--sm"
          disabled={template.length === 0 || lebih || kembar || kurang || angkaSalah || menghitung}
          onClick={() => onSalin(template)}
        >
          {L.salin}
        </button>
      </div>
    </div>
  )
}
