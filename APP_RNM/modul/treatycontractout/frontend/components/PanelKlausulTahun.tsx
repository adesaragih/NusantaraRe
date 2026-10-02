// Layar klausul satu tahun treaty — tiket 08 Treaty Contract Out.
//
// Padanan harness `InboxTreatyContractDescription` (b359) yang dibuka tombol
// `List Description` baris tahun (b22196): kepala tahun hanya dibaca, grid
// jenis `BrowseTreatyDesc_RD` `IsXOL = 0` — berjudul "Treaty Desc" — dan
// tombol `Show` per jenis.
//
// [keputusan work owner 30-09-2026] Grid `For XOL` (`IsXOL = 1`) dibuang —
// masternya kosong di DEV. `Show` membuka SATU jenis dalam popup: menekan
// Treaty Limit menampilkan Treaty Limit saja, menekan Cash Loss Limit
// menggantinya dengan Cash Loss Limit saja. Menggantikan AC 29 ("beberapa
// jenis boleh terbuka bersamaan").
// ⚠️ OQ-TCO-05: label `Underwriting Year` menunjuk `.TreatyYear` dan
// `Transaction Year` menunjuk `.UnderwritingYear` — dibawa apa adanya.

import { useEffect, useState } from 'react'

import { JUDUL_TAMPIL_TCO, KLAUSUL_TCO } from '../labels'
import { formatDate } from '../../../../inti/frontend/lib/format'
import { ambilJenisKlausul, type JenisKlausul, type TahunTreaty } from '../api'
import { Field, Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import PanelJenisKlausul from './PanelJenisKlausul'
import { labelProporsi } from '../proporsi'

/** Jenis yang tampil di popup: menekan jenis yang sama menutupnya, jenis lain menggantinya. */
export function alihJenisTunggal(terbuka: string | null, id: string): string | null {
  return terbuka === id ? null : id
}

export default function PanelKlausulTahun({ tahun, onTutup }: { tahun: TahunTreaty; onTutup?: () => void }) {
  const [daftar, setDaftar] = useState<JenisKlausul[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [terbuka, setTerbuka] = useState<string | null>(null)

  useEffect(() => {
    ambilJenisKlausul('0').then(setDaftar).catch(setGalat)
  }, [])

  const jenis = terbuka === null ? undefined : daftar?.find((j) => j.id === terbuka)
  return (
    <section className="panel">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{JUDUL_TAMPIL_TCO.deskripsi}</h2>
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
        <Field label={KLAUSUL_TCO.headerProportionType} value={labelProporsi(tahun.proportion)} onChange={() => undefined} readOnly />
      </div>
      <section className="panel">
        <h3 className="panel__title">{KLAUSUL_TCO.gridTreatyDesc}</h3>
        {galat !== null && <Gagal galat={galat} />}
        {daftar === null && galat === null && <Memuat />}
        {daftar !== null && daftar.length === 0 && <Kosong pesan={KLAUSUL_TCO.kosong} />}
        {daftar !== null && daftar.length > 0 && (
          <div className="tco-tabel">
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
                  <tr key={j.id} className={j.id === terbuka ? 'inbox__baris belah__baris--aktif' : 'inbox__baris'}>
                    <td>{j.id}</td>
                    <td>{j.descName}</td>
                    <td className="table__actions">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        aria-haspopup="dialog"
                        onClick={() => {
                          setTerbuka((t) => alihJenisTunggal(t, j.id))
                        }}
                      >
                        {KLAUSUL_TCO.show}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
      {jenis !== undefined && (
        <Modal
          judul={`${jenis.id} — ${jenis.descName}`}
          labelBatal={KLAUSUL_TCO.tutup}
          lebar
          onTutup={() => {
            setTerbuka(null)
          }}
        >
          <PanelJenisKlausul key={`${tahun.id}/${jenis.id}`} tahunID={tahun.id} jenis={jenis} tanpaJudul />
        </Modal>
      )}
    </section>
  )
}
