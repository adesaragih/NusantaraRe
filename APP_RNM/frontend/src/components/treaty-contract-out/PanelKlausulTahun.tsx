// Layar klausul satu tahun treaty — tiket 08 Treaty Contract Out.
//
// Padanan harness `InboxTreatyContractDescription` (b359) yang dibuka tombol
// `List Description` baris tahun (b22196): kepala tahun hanya dibaca, dua grid
// jenis `BrowseTreatyDesc_RD` — `For Non XOL` (IsXOL 0) dan `For XOL`
// (IsXOL 1) — dan tombol `Show` per jenis yang membuka panelnya.
//
// ⛔ Setiap jenis yang dibuka memegang isiannya sendiri (AC 29): beberapa
// jenis boleh terbuka bersamaan, menutup satu tidak menyentuh yang lain.
// ⚠️ OQ-TCO-05: label `Underwriting Year` menunjuk `.TreatyYear` dan
// `Transaction Year` menunjuk `.UnderwritingYear` — dibawa apa adanya.

import { useEffect, useState } from 'react'

import { KLAUSUL_TCO, MENU_TCO } from '../../assets/labels.treaty-contract-out'
import { formatDate } from '../../lib/format'
import { ambilJenisKlausul, type JenisKlausul, type TahunTreaty } from '../../services/api'
import { Field, Gagal, Kosong, Memuat } from '../ui/dasar'
import PanelJenisKlausul from './PanelJenisKlausul'

/** Buka/tutup satu jenis tanpa menyentuh jenis lain yang terbuka. */
export function alihJenis(terbuka: readonly string[], id: string): string[] {
  return terbuka.includes(id) ? terbuka.filter((t) => t !== id) : [...terbuka, id]
}

function GridJenis({
  judul,
  daftar,
  galat,
  onShow,
}: {
  judul: string
  daftar: JenisKlausul[] | null
  galat: unknown
  onShow: (j: JenisKlausul) => void
}) {
  return (
    <section className="panel">
      <h3 className="panel__title">{judul}</h3>
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={KLAUSUL_TCO.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{KLAUSUL_TCO.kolomId}</th>
              <th>{KLAUSUL_TCO.kolomDescriptionName}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.map((j) => (
              <tr key={j.id} className="inbox__baris">
                <td>{j.id}</td>
                <td>{j.descName}</td>
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => onShow(j)}>
                    {KLAUSUL_TCO.show}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}

export default function PanelKlausulTahun({ tahun, onTutup }: { tahun: TahunTreaty; onTutup?: () => void }) {
  const [nonXol, setNonXol] = useState<JenisKlausul[] | null>(null)
  const [xol, setXol] = useState<JenisKlausul[] | null>(null)
  const [galatNonXol, setGalatNonXol] = useState<unknown>(null)
  const [galatXol, setGalatXol] = useState<unknown>(null)
  const [terbuka, setTerbuka] = useState<string[]>([])

  useEffect(() => {
    ambilJenisKlausul('0').then(setNonXol).catch(setGalatNonXol)
    ambilJenisKlausul('1').then(setXol).catch(setGalatXol)
  }, [])

  const semua = [...(nonXol ?? []), ...(xol ?? [])]
  return (
    <section className="panel">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU_TCO.inboxTreatyContractDescription}</h2>
        {onTutup !== undefined && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {KLAUSUL_TCO.tutup}
          </button>
        )}
      </header>
      <div className="form-grid">
        <Field label={KLAUSUL_TCO.headerTreatyGroupId} value={tahun.treatyGroupId} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerUnderwritingYear} value={tahun.treatyYear} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerTransactionYear} value={tahun.underwritingYear} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerStartDate} value={formatDate(tahun.startDate)} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerEndDate} value={formatDate(tahun.endDate)} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerTreatyDescription} value={tahun.treatyGroupName} onChange={() => undefined} readOnly />
        <Field label={KLAUSUL_TCO.headerProportionType} value={tahun.proportion} onChange={() => undefined} readOnly />
      </div>
      <GridJenis judul={KLAUSUL_TCO.gridNonXol} daftar={nonXol} galat={galatNonXol} onShow={(j) => setTerbuka((t) => alihJenis(t, j.id))} />
      <GridJenis judul={KLAUSUL_TCO.gridXol} daftar={xol} galat={galatXol} onShow={(j) => setTerbuka((t) => alihJenis(t, j.id))} />
      {terbuka.map((id) => {
        const j = semua.find((s) => s.id === id)
        return j === undefined ? null : (
          <PanelJenisKlausul key={`${tahun.id}/${id}`} tahunID={tahun.id} jenis={j} onTutup={() => setTerbuka((t) => alihJenis(t, id))} />
        )
      })}
    </section>
  )
}
