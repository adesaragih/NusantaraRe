// Sub-tab FEA (Fire Extinguisher Availability) baris objek FIRE (tiket 41).
//
// Port grid `NB FacIn\Section\FEAList.xml` (`.FEAList` baris objek: APAR · Sprinkler · Smoke Detector & Alarm ·
// Hydrant · Private Truck Brigade · Others Info; Add di kepala = addRow, Delete per baris = deleteRow, tanpa activity).
// Form baris dibuka: section asli `OfferFEAList!InputFEA` TIDAK ada di korpus - medan mengikuti padanannya
// `InputFEA_IsUW` (kelas sama), dibuat dapat diisi.
//
// Keputusan agent (tiket 41): M-1 form = medan `InputFEA_IsUW` (APAR / Sprinkler / Smoke Detector & Alarm / Hydrant /
// Private Truck Brigade (Unit), Private Team Fire Brigade, Team & SOP Safety, Team & SOP Risk Management, Others
// Info). M-2 jumlah unit = bilangan bulat >= 0 (pxNumber tanpa batas di XML); galat menahan Save. M-3 tombol Save di
// atas grid FEA Pega (`SaveFacIn_Act`) tidak diulang - Save tab Object sudah menyimpan seluruh objek. M-4 tiga
// dropdown = aturan properti `DDL\PrivateFireBrigade.xml` dst. (Have / Not Have / No Info).

import { useState } from 'react'

import { Area, Field, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisFEA } from '../api'
import { FORM_FEA as F, GRID_FEA, GRID_OBJEK, OPSI_FIRE_BRIGADE, OPSI_SOP_RISIKO, OPSI_SOP_SAFETY, TEKS_FEA, TEKS_INWARD, TEKS_OBJEK } from '../labels'

/** Medan jumlah unit. */
const MEDAN_UNIT = ['apar', 'sprinkler', 'smokeDetector', 'hydrant', 'privateTruckBrigade'] as const

/** Baris FEA kosong (Add). */
export const feaBaru = (): BarisFEA => ({
  apar: '', sprinkler: '', smokeDetector: '', hydrant: '', privateTruckBrigade: '', privateFireBrigade: '', teamSopSafety: '',
  teamSopRiskManagement: '', info: '',
})

/** Jumlah unit sah: kosong atau bilangan bulat >= 0 (M-2). */
export const unitSah = (v: string) => v.trim() === '' || /^\d+$/.test(v.trim())

/** Ada baris FEA bergalat. */
export function adaGalatFEA(fea: BarisFEA[]): boolean {
  return fea.some((b) => MEDAN_UNIT.some((k) => !unitSah(b[k])))
}

function FormFEA({ b, ubah }: { b: BarisFEA; ubah: (b: BarisFEA) => void }) {
  const set = (k: keyof BarisFEA) => (v: string) => ubah({ ...b, [k]: v })
  const unit = (k: (typeof MEDAN_UNIT)[number]) => (
    <Field key={k} label={F[k].label} value={b[k]} onChange={set(k)} error={unitSah(b[k]) ? undefined : TEKS_FEA.unit} />
  )
  return (
    <div className="nbf-objek__isi">
      <div className="nbf-opp__kolom">
        <div className="nbf-opp__tumpuk">{MEDAN_UNIT.map(unit)}</div>
        <div className="nbf-opp__tumpuk">
          <Pilih label={F.privateFireBrigade.label} value={b.privateFireBrigade} onChange={set('privateFireBrigade')} opsi={OPSI_FIRE_BRIGADE} />
          <Pilih label={F.teamSopSafety.label} value={b.teamSopSafety} onChange={set('teamSopSafety')} opsi={OPSI_SOP_SAFETY} />
          <Pilih
            label={F.teamSopRiskManagement.label}
            value={b.teamSopRiskManagement}
            onChange={set('teamSopRiskManagement')}
            opsi={OPSI_SOP_RISIKO}
          />
          <Area label={F.info.label} value={b.info} onChange={set('info')} baris={3} />
        </div>
      </div>
    </div>
  )
}

export default function SubTabFEA({ fea, ubah }: { fea: BarisFEA[]; ubah: (fea: BarisFEA[]) => void }) {
  const [terbuka, setTerbuka] = useState<number[]>([])
  return (
    <div className="nbf-objek__isi">
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_FEA.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    ubah([...fea, feaBaru()])
                    setTerbuka((t) => [...t, fea.length])
                  }}
                >
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {fea.length === 0 && (
              <tr>
                <td colSpan={8}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {fea.map((b, n) => [
              <tr key={`b-${n}`}>
                <td>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    aria-label={TEKS_OBJEK.bukaBaris}
                    aria-expanded={terbuka.includes(n)}
                    onClick={() => setTerbuka((t) => (t.includes(n) ? t.filter((x) => x !== n) : [...t, n]))}
                  >
                    {terbuka.includes(n) ? '▾' : '▸'}
                  </button>
                </td>
                <td className="nbf-angka">{b.apar}</td>
                <td className="nbf-angka">{b.sprinkler}</td>
                <td className="nbf-angka">{b.smokeDetector}</td>
                <td className="nbf-angka">{b.hydrant}</td>
                <td className="nbf-angka">{b.privateTruckBrigade}</td>
                <td>{b.info}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah(fea.filter((_, x) => x !== n))
                      setTerbuka((t) => t.filter((x) => x !== n).map((x) => (x > n ? x - 1 : x)))
                    }}
                  >
                    {GRID_OBJEK.hapus}
                  </button>
                </td>
              </tr>,
              terbuka.includes(n) && (
                <tr key={`d-${n}`} className="nbf-objek__detail">
                  <td colSpan={8}>
                    <FormFEA b={b} ubah={(baru) => ubah(fea.map((x, k) => (k === n ? baru : x)))} />
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>
    </div>
  )
}
