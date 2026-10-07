// Popup `Harness/BusinessAndSOBListEDM` -> `Section/BusinessAndSOBListEDM` ("Business And SOB List"). Asal pola:
// `modul/nbtreatyin/frontend/components/PilihBisnis.tsx` (06-10-2026: Modal penuh, judul kolom dapat diurut,
// saringan per kolom, 50 baris per halaman, tombol Choose per baris).
//
//   judul jendela   showHarness `pyWindowName` "Business And SOB List"
//   isi             POST /kasus/{id}/bisnis dengan isian layar saat tombol diklik (S1 `.EDMType = 3` RD
//                   `BrowseTREATY_IN_EDM`; S6 selain itu `TreatyLoadMasterJoinEdmChooseBusiness` - dipilih server)
//   grid            Choose · ID Revision · Previous ID · Treaty Contract Name · Proportion Type · SOB · Ceding ·
//                   Start Date · End Date (`KOLOM_BISNIS`); urut awal ID Revision turun; saringan daftar nilai
//   tombol Choose   `EDMChooseBusiness_Act` -> refresh -> opener reload -> window.close (`onPilih`). XML tidak membaca
//                   baris yang diklik (semua kueri memakai `OldData.NoOffer`) - permintaan tidak membawa ID baris.

import { useMemo, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { daftarBisnis, type Halaman } from '../api'
import {
  URUT_AWAL_BISNIS,
  nilaiSaringBisnis,
  saringUrutBisnis,
  urutBerikut,
  type SaringanKolom,
  type UrutBisnis,
} from '../gridbisnis'
import { JUDUL, KOLOM_BISNIS, TOMBOL } from '../labels'
import { BARIS_PER_HALAMAN_GRID, irisan } from '../paginasi'
import { sajikan } from '../sajian'
import Paginasi from './Paginasi'

/** Kolom tanggal (`.Commencement` / `.Termination`) - satu format tanggal seluruh sistem. */
const KOLOM_TANGGAL = new Set(['Commencement', 'Termination'])
const sel = (kolom: string, v: string) => (KOLOM_TANGGAL.has(kolom) ? sajikan(v, 'tanggal') || v : v)

export default function PilihBisnis({
  id,
  halaman,
  sibuk,
  onPilih,
  onTutup,
}: {
  id: string
  /** Isian layar saat tombol `Choose Business` diklik - dikirim sekali, saat popup dibuka. */
  halaman: Halaman
  sibuk: boolean
  onPilih: () => void
  onTutup: () => void
}) {
  const [awal] = useState(halaman)
  const { data: baris, galat } = useAmbil(() => daftarBisnis(id, awal), [id, awal])
  const [saringan, setSaringan] = useState<SaringanKolom>({})
  const [urut, setUrut] = useState<UrutBisnis>(URUT_AWAL_BISNIS)
  const [hal, setHal] = useState(1)
  const tampil = useMemo(() => (baris === null ? [] : saringUrutBisnis(baris, saringan, urut)), [baris, saringan, urut])

  return (
    <Modal judul={JUDUL.pilihBisnis} onTutup={onTutup} penuh>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col" />
                  {KOLOM_BISNIS.map((k) => {
                    const aktif = urut.kolom === k.kolom
                    return (
                      <th
                        scope="col"
                        key={k.kolom}
                        aria-sort={aktif ? (urut.arah === 'naik' ? 'ascending' : 'descending') : 'none'}
                      >
                        <button
                          type="button"
                          className="edmt__urut"
                          onClick={() => {
                            setUrut(urutBerikut(urut, k.kolom))
                            setHal(1)
                          }}
                        >
                          {k.judul}
                          {aktif && <span aria-hidden="true">{urut.arah === 'naik' ? ' ▲' : ' ▼'}</span>}
                        </button>
                      </th>
                    )
                  })}
                </tr>
                <tr>
                  <td />
                  {KOLOM_BISNIS.map((k) => (
                    <td key={k.kolom}>
                      {/* pyColumnFilteringDropDown true: daftar nilai kolom */}
                      <select
                        className="field__input edmt__saring-kolom"
                        aria-label={`${TOMBOL.filter} ${k.judul}`}
                        value={saringan[k.kolom] ?? ''}
                        onChange={(e) => {
                          setSaringan((s) => ({ ...s, [k.kolom]: e.target.value }))
                          setHal(1)
                        }}
                      >
                        <option value="" />
                        {nilaiSaringBisnis(baris, k.kolom).map((v) => (
                          <option key={v} value={v}>
                            {sel(k.kolom, v)}
                          </option>
                        ))}
                      </select>
                    </td>
                  ))}
                </tr>
              </thead>
              <tbody>
                {tampil.length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_BISNIS.length + 1}>—</td>
                  </tr>
                )}
                {irisan(tampil, hal, BARIS_PER_HALAMAN_GRID).map((b, i) => (
                  <tr key={`${b.ID ?? ''}-${i}`}>
                    <td>
                      <button type="button" className="btn btn--primary" disabled={sibuk} onClick={onPilih}>
                        {TOMBOL.choose}
                      </button>
                    </td>
                    {KOLOM_BISNIS.map((k) => (
                      <td key={k.kolom}>{sel(k.kolom, b[k.kolom] ?? '')}</td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {/* pyPageMode Numeric, pyPageSize 50 */}
          <Paginasi jumlahBaris={tampil.length} ukuran={BARIS_PER_HALAMAN_GRID} hal={hal} onHal={setHal} />
        </>
      )}
    </Modal>
  )
}
