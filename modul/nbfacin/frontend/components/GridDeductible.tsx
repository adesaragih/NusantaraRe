// Grid Deductible satu coverage FIRE (tiket 45, tahap C3).
//
// Port `.DeductibleList` di `NB FacIn\Section\CoverageItem.xml` (judul "Deductible", visible `!isClaim`; Add =
// addRow + `addDeductible_ACT`, Delete = deleteRow) dan form flow action `InputDtlDeductibleFire_FacIn` -> section
// `addDeductible.xml`:
// - Pct Deductible dan MinMax tampil bila Type Deductible != 7 (In Amount); MinMax baca-saja bila Type 0 / 7;
// - Currency baca-saja bila Type 0 (NIL);
// - Type Deductible 2 / Pct Deductible 2 tampil bila MinMax = 3 (Or);
// - Condition (teks) tampil bila Condition = 5 (Other).
// Daftar = aturan properti `DATA-DEDUCTIBLE!*` (`DDL\`); Currency = daftar mata uang (RD `BrowseCurrency_RD`).
//
// Keputusan agent (tiket 45): R-1 kepala kolom grid memakai label form (di XML kosong kecuali "Time Excess (Days)") -
// DIRALAT work owner 08-10-2026: kepala kosong seperti Pega, hanya "Time Excess (Days)".
// R-2 deductible baru bermata uang item (`addDeductible_ACT` membaca `Local.Curr` = mata uang item) `[dugaan]`.
// R-3 penanda beda mata uang (`SaveDeductible` / `addDeductible_ACT` -> IsBedaCurr, FlagCurrency), Descriptions (B2B),
// View Old Deductible (EDM) dan Copy Deductible To All / For BI = tahap berikut. R-4 deductible tersimpan lewat Save
// tab Coverage (tombol Save di form Pega = `ObjSave` halaman yang sama).

import { useState } from 'react'

import { Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import type { Deductible } from '../api'
import {
  DEDUCTIBLE_KOSONG,
  FORM_DEDUCTIBLE as F,
  GRID_OBJEK,
  JUDUL_DEDUCTIBLE,
  KOLOM_TIME_EXCESS,
  OPSI_KONDISI_DEDUCTIBLE,
  OPSI_MINMAX,
  OPSI_TYPE_DEDUCTIBLE,
  OPSI_TYPE_DEDUCTIBLE2,
  TEKS_COVERAGE,
  TEKS_INWARD,
  TEKS_OBJEK,
} from '../labels'
import { angkaSah } from './FormCoverage'
import IsianUang from './IsianUang'

/** Deductible baru bermata uang item (R-2). */
export const deductibleBaru = (currency: string): Deductible => ({
  typeDeductible: '', pctDeductible: '', minMax: '', currency, typeDeductible2: '', pctDeductible2: '', condition: '', amount: '',
  inputCondition: '', timeExcess: '',
})

/** Teks label pilihan (atau kode apa adanya bila di luar daftar). */
const teks = (opsi: readonly Opsi[], v: string) => opsi.find((o) => o.value === v)?.label ?? v

/** Medan angka deductible. */
const MEDAN_ANGKA = ['pctDeductible', 'pctDeductible2', 'amount', 'timeExcess'] as const

/** Ada deductible dengan angka tidak sah. */
export function adaGalatDeductible(d: Deductible[] | undefined): boolean {
  return (d ?? []).some((x) => MEDAN_ANGKA.some((k) => !angkaSah(x[k])))
}

function FormDeductible({ d, ubah, mataUang }: { d: Deductible; ubah: (d: Deductible) => void; mataUang: string[] }) {
  const set = (k: keyof Deductible) => (v: string) => ubah({ ...d, [k]: v })
  const err = (k: (typeof MEDAN_ANGKA)[number]) => (angkaSah(d[k]) ? undefined : TEKS_COVERAGE.angka)
  const opsiUang: Opsi[] = mataUang.map((m) => ({ value: m, label: m }))
  const t = d.typeDeductible
  // Tata letak (work owner 08-10-2026 "rapikan dong biar lebih user friendly"): kartu berbaris tiga kolom menurut
  // kelompok - jenis deductible / nilai / kondisi. Medan, syarat tampil, dan baca-saja tetap addDeductible.
  return (
    <div className="nbf-dedgrid__kartu">
      <div className="nbf-dedgrid__baris">
        <Pilih label={F.typeDeductible.label} value={t} onChange={set('typeDeductible')} opsi={OPSI_TYPE_DEDUCTIBLE} kosong={DEDUCTIBLE_KOSONG} />
        {t !== '7' && <Field label={F.pctDeductible.label} value={d.pctDeductible} onChange={set('pctDeductible')} error={err('pctDeductible')} />}
        {t !== '7' &&
          (t === '0' ? (
            <Field label={F.minMax.label} value={teks(OPSI_MINMAX, d.minMax)} onChange={() => {}} readOnly />
          ) : (
            <Pilih label={F.minMax.label} value={d.minMax} onChange={set('minMax')} opsi={OPSI_MINMAX} kosong={DEDUCTIBLE_KOSONG} />
          ))}
      </div>
      {d.minMax === '3' && (
        <div className="nbf-dedgrid__baris">
          <Pilih
            label={F.typeDeductible2.label}
            value={d.typeDeductible2}
            onChange={set('typeDeductible2')}
            opsi={OPSI_TYPE_DEDUCTIBLE2}
            kosong={DEDUCTIBLE_KOSONG}
          />
          <Field label={F.pctDeductible2.label} value={d.pctDeductible2} onChange={set('pctDeductible2')} error={err('pctDeductible2')} />
        </div>
      )}
      <div className="nbf-dedgrid__baris">
        {t === '0' ? (
          <Field label={F.currency.label} value={d.currency} onChange={() => {}} readOnly />
        ) : (
          <Pilih label={F.currency.label} value={d.currency} onChange={set('currency')} opsi={opsiUang} kosong={DEDUCTIBLE_KOSONG} />
        )}
        <IsianUang label={F.amount.label} value={d.amount} onChange={set('amount')} error={err('amount')} />
        <Field label={F.timeExcess.label} value={d.timeExcess} onChange={set('timeExcess')} error={err('timeExcess')} />
      </div>
      <div className="nbf-dedgrid__baris">
        <Pilih label={F.condition.label} value={d.condition} onChange={set('condition')} opsi={OPSI_KONDISI_DEDUCTIBLE} kosong={DEDUCTIBLE_KOSONG} />
        {d.condition === '5' && (
          <div className="nbf-dedgrid__lebar">
            <Field label={F.inputCondition.label} value={d.inputCondition} onChange={set('inputCondition')} />
          </div>
        )}
      </div>
    </div>
  )
}

/** Sel grid: kosong tampil "—" supaya baris yang belum diisi tetap terbaca. */
const isi = (v: string) => (v === '' ? '—' : v)

export default function GridDeductible({
  daftar,
  currency,
  mataUang,
  ubah,
}: {
  daftar: Deductible[]
  /** Mata uang item (nilai awal deductible baru). */
  currency: string
  mataUang: string[]
  ubah: (d: Deductible[]) => void
}) {
  const [terbuka, setTerbuka] = useState<number[]>([])
  return (
    <div className="nbf-objek__isi nbf-dedgrid">
      <div className="nbf-dedgrid__kepala">
        <h5 className="nbf-objek__judul">
          {JUDUL_DEDUCTIBLE}
          <span className="nbf-dedgrid__jumlah">{daftar.length}</span>
        </h5>
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--sm nbf-dedgrid__tambah"
          onClick={() => {
            ubah([...daftar, deductibleBaru(currency)])
            setTerbuka((x) => [...x, daftar.length])
          }}
        >
          + {GRID_OBJEK.tambah}
        </button>
      </div>
      <div className="table-wrap">
        <table className="nbf-tabel nbf-dedgrid__tabel">
          {/* Lebar kolom tetap (work owner 08-10-2026 "rapikan jarak antar kolom nya"): angka / kode sempit, teks lebar. */}
          <colgroup>
            <col className="nbf-dedgrid__k-buka" />
            <col className="nbf-dedgrid__k-angka" />
            <col />
            <col className="nbf-dedgrid__k-kode" />
            <col className="nbf-dedgrid__k-angka" />
            <col />
            <col className="nbf-dedgrid__k-kode" />
            <col className="nbf-dedgrid__k-uang" />
            <col />
            <col className="nbf-dedgrid__k-hari" />
            <col className="nbf-dedgrid__k-aksi" />
          </colgroup>
          {/* Kepala kolom = CoverageItem.xml: kosong kecuali "Time Excess (Days)" (work owner 08-10-2026 "judul kolom yg
              diatas itu tidak usah dibuat, hanya Time Excess (Days) aja"). */}
          <thead>
            <tr>
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" />
              <th scope="col" className="nbf-angka">{KOLOM_TIME_EXCESS}</th>
              <th scope="col" className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.length === 0 && (
              <tr>
                <td colSpan={11} className="nbf-dedgrid__kosong">{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {daftar.map((d, n) => [
              <tr key={`b-${n}`} className={terbuka.includes(n) ? 'nbf-dedgrid__aktif' : undefined}>
                <td>
                  <button
                    type="button"
                    className="btn nbf-tema__tombol btn--ghost btn--sm nbf-dedgrid__panah"
                    aria-label={TEKS_OBJEK.bukaBaris}
                    aria-expanded={terbuka.includes(n)}
                    onClick={() => setTerbuka((x) => (x.includes(n) ? x.filter((y) => y !== n) : [...x, n]))}
                  >
                    {terbuka.includes(n) ? '▾' : '▸'}
                  </button>
                </td>
                <td className="nbf-angka">{isi(d.pctDeductible)}</td>
                <td>{isi(teks(OPSI_TYPE_DEDUCTIBLE, d.typeDeductible))}</td>
                <td>{isi(teks(OPSI_MINMAX, d.minMax))}</td>
                {/* CoverageItem.xml: sel `.PctDeductible2` / `.TypeDeductible2` hanya bila `.MinMax = 3` (Or). */}
                <td className="nbf-angka">{isi(d.minMax === '3' ? d.pctDeductible2 : '')}</td>
                <td>{isi(d.minMax === '3' ? teks(OPSI_TYPE_DEDUCTIBLE2, d.typeDeductible2) : '')}</td>
                <td>{isi(d.currency)}</td>
                <td className="nbf-angka">{isi(angka4(d.amount))}</td>
                <td>{isi(d.condition === '5' ? d.inputCondition : teks(OPSI_KONDISI_DEDUCTIBLE, d.condition))}</td>
                <td className="nbf-angka">{isi(d.timeExcess)}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm nbf-dedgrid__hapus"
                    onClick={() => {
                      ubah(daftar.filter((_, x) => x !== n))
                      setTerbuka((x) => x.filter((y) => y !== n).map((y) => (y > n ? y - 1 : y)))
                    }}
                  >
                    {GRID_OBJEK.hapus}
                  </button>
                </td>
              </tr>,
              terbuka.includes(n) && (
                <tr key={`d-${n}`} className="nbf-objek__detail">
                  <td colSpan={11}>
                    <FormDeductible d={d} ubah={(baru) => ubah(daftar.map((x, k) => (k === n ? baru : x)))} mataUang={mataUang} />
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
