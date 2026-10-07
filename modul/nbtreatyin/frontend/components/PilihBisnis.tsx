// Popup `Harness/BusinessAndSOBList` - grid AKTIF `Section/BusinessAndSOBList` (RD `BrowseTreatyJoinEDM`,
// view TREATYINDETAILJOINEDM; grid lama RD `BrowseTreatyInDetail` ber-`1=2`), audit silang P3 W1.
//
//   judul jendela   showHarness `pyWindowName` "Business And SOB List"
//   isi             POST /kasus/{id}/bisnis dengan isian layar saat tombol diklik (`pySubmitData=Yes`):
//                   server menyaring `.PROPORTIONTYPE = .QuotationData.ProportionalType`, <= 500 baris
//   grid            judul 25 kolom VERBATIM (`KOLOM_BISNIS`); urut di klien (`pyGridSorting`); saring per kolom
//                   (`pyGridFiltering`) DIKETIK dan dicari di SERVER sebelum batas 500 (keputusan work owner
//                   06-10-2026, pengganti daftar nilai `pyColumnFilteringDropDown` yang hanya memuat 500 baris
//                   pertama), 50 baris per halaman
//                   (`pyPageMode Numeric`, `pyPageSize 50`) - nol kotak cari dan tombol Filter
//   tombol Choose   `SetValue_Act(ID=.ID)` -> `InputPolicyTreatyInDetail_preACT` (`onPilih`)
//
// Kolom `*VALUE` = pxCurrency tanpa `pyDecimalPlaces` -> sajian pola inti (K14, `sajian.ts`).

import { useEffect, useMemo, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbilBatal } from '../ambil'
import { daftarBisnis, type BarisKontrak, type Halaman } from '../api'
import {
  BARIS_PER_HALAMAN_BISNIS,
  URUT_AWAL_BISNIS,
  kolomUangBisnis,
  saringUrutBisnis,
  urutBerikut,
  type SaringanKolom,
  type UrutBisnis,
} from '../gridbisnis'
import { JUDUL, KOLOM_BISNIS, TOMBOL } from '../labels'
import { irisan } from '../paginasi'
import { sajikan } from '../sajian'
import Paginasi from './Paginasi'

const sel = (kolom: string, v: string) => (kolomUangBisnis(kolom) ? sajikan(v, {}) : v)

export default function PilihBisnis({
  id,
  halaman,
  onPilih,
  onTutup,
}: {
  id: string
  /** Isian layar saat tombol `Choose Business` diklik - dikirim sekali, saat popup dibuka. */
  halaman: Halaman
  onPilih: (id: string) => void
  onTutup: () => void
}) {
  // isian layar SAAT DIBUKA (bukan sunting sesudahnya); saringan dikirim ke server 400 ms sesudah ketikan terakhir
  const [saringan, setSaringan] = useState<SaringanKolom>({})
  const [saringanTunda, setSaringanTunda] = useState<SaringanKolom>({})
  const [baris, setBaris] = useState<BarisKontrak[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(true)
  const [urut, setUrut] = useState<UrutBisnis>(URUT_AWAL_BISNIS)
  const [hal, setHal] = useState(1)
  useEffect(() => {
    const t = setTimeout(() => setSaringanTunda(saringan), 400)
    return () => clearTimeout(t)
  }, [saringan])
  // baris lama TETAP tampil selama memuat ulang - kotak saringan tidak hilang dan fokus ketik tidak lepas
  useAmbilBatal(
    () => {
      setMemuat(true)
      setGalat(null)
      return daftarBisnis(id, halaman, saringanTunda)
    },
    (b) => {
      setBaris(b)
      setMemuat(false)
    },
    (e) => {
      setGalat(e)
      setMemuat(false)
    },
    [id, JSON.stringify(saringanTunda)],
  )
  // saringan sudah dikerjakan server; di klien hanya urut
  const tampil = useMemo(() => (baris === null ? [] : saringUrutBisnis(baris, {}, urut)), [baris, urut])

  return (
    <Modal judul={JUDUL.pilihBisnis} onTutup={onTutup} penuh>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <>
          {memuat && <Memuat />}
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
                        className={kolomUangBisnis(k.kolom) ? 'nbti__angka' : undefined}
                        aria-sort={aktif ? (urut.arah === 'naik' ? 'ascending' : 'descending') : 'none'}
                      >
                        <button
                          type="button"
                          className="nbti__urut"
                          aria-label={k.judul || k.kolom}
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
                      <input
                        className="field__input nbti__saring-kolom"
                        aria-label={`${TOMBOL.filter} ${k.judul || k.kolom}`}
                        value={saringan[k.kolom] ?? ''}
                        onChange={(e) => {
                          setSaringan((s) => ({ ...s, [k.kolom]: e.target.value }))
                          setHal(1)
                        }}
                      />
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
                {irisan(tampil, hal, BARIS_PER_HALAMAN_BISNIS).map((b, i) => (
                  <tr key={`${b.ID ?? ''}-${i}`}>
                    <td>
                      <button type="button" className="btn btn--primary" onClick={() => onPilih(b.ID ?? '')}>
                        {TOMBOL.choose}
                      </button>
                    </td>
                    {KOLOM_BISNIS.map((k) => (
                      <td key={k.kolom} className={kolomUangBisnis(k.kolom) ? 'nbti__angka' : undefined}>
                        {sel(k.kolom, b[k.kolom] ?? '')}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {/* pyPageMode Numeric, pyPageSize 50 */}
          <Paginasi jumlahBaris={tampil.length} ukuran={BARIS_PER_HALAMAN_BISNIS} hal={hal} onHal={setHal} />
        </>
      )}
    </Modal>
  )
}
