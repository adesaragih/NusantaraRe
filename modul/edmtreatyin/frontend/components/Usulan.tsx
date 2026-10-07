// `Section/ListSuggestEDM` (SUB_SECTION `DetailPolicyTreatyInAddendum`, wadah NOHEADER). Asal pola: blok usulan
// `modul/nbtreatyin/frontend/pages/LayarKasus.tsx` (06-10-2026: radio Approval berbingkai, kolom Date tanggal + jam,
// paginasi 5 baris).
//
//   Approval  `.IsApproved` pxRadioButtons, WAJIB - change/click -> runActivity SetDueTo_act -> refresh
//             CekLimitTreatyAcc_Act -> refresh Protection_Act (`onApproval`: aksi backend `SetDueTo` lalu `Protection`)
//   Suggest   `.Suggest` pxTextArea, WAJIB - change -> refresh thisSection (tanpa activity)
//   grid      `.SuggestList` Date / PIC / Approval / Suggest, hanya-baca, pyPageSizeOther 5
//
// ⛔ Beda dengan NB `ListSuggest`: TANPA medan Production Date.

import { useState } from 'react'

import { Area } from '../../../../inti/frontend/components/ui/dasar'
import { POLIS, daftar, nilai, type Halaman } from '../api'
import { KOLOM_USULAN, PILIHAN_APPROVAL } from '../labels'
import { BARIS_PER_HALAMAN_USULAN, irisan } from '../paginasi'
import { sajikanTanggalJam } from '../sajian'
import Paginasi from './Paginasi'
import Wadah from './Wadah'

export const JALUR_APPROVAL = POLIS + 'IsApproved'
export const JALUR_SUGGEST = POLIS + 'Suggest'
const USULAN = POLIS + 'SuggestList'

export default function Usulan({
  halaman: h,
  boleh,
  onUbah,
  onApproval,
}: {
  halaman: Halaman
  /** Pelaku boleh bekerja di kasus ini - isian Approval / Suggest tampil. */
  boleh: boolean
  onUbah: (jalur: string, v: string) => void
  onApproval: (v: string) => void
}) {
  const [hal, setHal] = useState(1)
  const usulan = daftar(h, USULAN)
  const pesan = (j: string) => h.pesan?.[j]?.join('; ')
  return (
    <Wadah>
      {boleh && (
        <div className="edmt__kolom edmt__usulan">
          <div className="field">
            <span className="field__label">
              {KOLOM_USULAN.putusan}
              <span className="field__req">*</span>
            </span>
            <div className="edmt__radio" role="radiogroup" aria-label={KOLOM_USULAN.putusan}>
              {PILIHAN_APPROVAL.map((o) => (
                <label key={o.value}>
                  <input
                    type="radio"
                    name="edmt-approval"
                    checked={nilai(h, JALUR_APPROVAL) === o.value}
                    onChange={() => onApproval(o.value)}
                  />
                  {o.label}
                </label>
              ))}
            </div>
            {pesan(JALUR_APPROVAL) && <div className="field__error">{pesan(JALUR_APPROVAL)}</div>}
          </div>
          <Area
            label={KOLOM_USULAN.catatan}
            value={nilai(h, JALUR_SUGGEST)}
            required
            error={pesan(JALUR_SUGGEST)}
            onChange={(x) => onUbah(JALUR_SUGGEST, x)}
          />
        </div>
      )}
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th scope="col">{KOLOM_USULAN.tanggal}</th>
              <th scope="col">{KOLOM_USULAN.operator}</th>
              <th scope="col">{KOLOM_USULAN.putusan}</th>
              <th scope="col">{KOLOM_USULAN.catatan}</th>
            </tr>
          </thead>
          <tbody>
            {irisan(usulan, hal, BARIS_PER_HALAMAN_USULAN).map((b, i) => (
              <tr key={i}>
                {/* tanggal + jam (keputusan work owner 06-10-2026 di NB untuk grid yang sama) */}
                <td>{sajikanTanggalJam(b.Date)}</td>
                <td>{b.OperatorName ?? ''}</td>
                <td>{PILIHAN_APPROVAL.find((o) => o.value === b.IsApproved)?.label ?? b.IsApproved ?? ''}</td>
                <td className="edmt__catatan">{b.Suggest ?? ''}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {/* pyGridPaginator - pyPageMode Numeric, pyPageSizeOther 5 */}
      <Paginasi jumlahBaris={usulan.length} ukuran={BARIS_PER_HALAMAN_USULAN} hal={hal} onHal={setHal} />
    </Wadah>
  )
}
