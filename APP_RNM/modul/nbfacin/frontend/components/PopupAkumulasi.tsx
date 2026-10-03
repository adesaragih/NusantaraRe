// Popup Choose Accumulation (tiket 46, tahap 1) - port harness `NB FacIn\Harness\ChooseAccumulation_FacIn.xml` ->
// section `SearchRiskAccumCov.xml` jalur Report Definition (grid sel 90, `SearchRiskAccumulation_RD`):
// - awal: Zip Code = zip lokasi risiko (`SearchAccumAct` langkah 4 `SearchAccumulation.PostalCode = .RiskZipCode`);
// - Filter (sel 40, `GetDataAccumulation_act`) -> `GET …/akumulasi`; Clear Column (sel 41, `SearchAccumulationPre_Act`)
//   mengosongkan saringan dan hasil;
// - Choose (sel 119, `SetDataAccum_Act`): zip dari ID harus sama dengan zip lokasi risiko; beda -> pesan, popup tetap
//   terbuka; sama -> AccumulationCode = ID, AccumulationDescription = Note, popup tertutup.
//
// Keputusan agent (tiket 46): K-1 tahap 1 hanya medan yang dipakai RD dan dapat diisi teks (Accumulation Code, Road,
// Zip Code, CZone). Policy No / City / District / Area (jalur SQL grid sel 48), Country / Province / Accum. Type
// (autocomplete), Key Word (aturan `.Keyword` belum ada di korpus), Add New, Summary, Accumulation Risk, Risk
// Accumulation report, Adjustment = tahap berikut. K-2 tanpa paging 10 baris (hasil maks. 500, RD pyMaxRecords).

import { useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariAkumulasi, type BarisAkumulasi, type SaringAkumulasi } from '../api'
import { POPUP_AKUMULASI as P, TEKS_AKUMULASI } from '../labels'

/**
 * Zip di dalam Accumulation Code - `SetDataAccum_Act` langkah 2-3: `@substring(ID,4,9)` tanpa "-", kecuali karakter
 * ke-3 (0-based) bukan "-" -> `@substring(ID,3,8)`. Substring gaya Java (awal inklusif, akhir eksklusif); ID terlalu
 * pendek -> kosong (Pega melempar galat) `[dugaan]`.
 */
export function zipDariId(id: string): string {
  const sub = (a: number, b: number) => (id.length >= b ? id.slice(a, b) : null)
  const c3 = sub(3, 4)
  if (c3 !== null && c3 !== '-') return sub(3, 8) ?? ''
  return (sub(4, 9) ?? '').replaceAll('-', '')
}

const SARING_KOSONG: SaringAkumulasi = { id: '', note: '', postalCode: '', czone: '' }

export default function PopupAkumulasi({
  zipRisiko,
  onTutup,
  onPilih,
}: {
  /** `.Property.RiskLocation.ASMZipCode` objek pemilik coverage. */
  zipRisiko: string
  onTutup: () => void
  onPilih: (kode: string, alamat: string) => void
}) {
  const [saring, setSaring] = useState<SaringAkumulasi>({ ...SARING_KOSONG, postalCode: zipRisiko })
  const [hasil, setHasil] = useState<BarisAkumulasi[] | null>(null)
  const [memuat, setMemuat] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState('')
  const set = (k: keyof SaringAkumulasi) => (v: string) => setSaring((s) => ({ ...s, [k]: v }))

  async function filter() {
    setMemuat(true)
    setGalat(null)
    setPesan('')
    try {
      setHasil((await cariAkumulasi(saring)).baris)
    } catch (err) {
      setGalat(err)
    } finally {
      setMemuat(false)
    }
  }

  function pilih(b: BarisAkumulasi) {
    const dariId = zipDariId(b.id)
    if (dariId !== zipRisiko) {
      setPesan(TEKS_AKUMULASI.zipBeda(zipRisiko, dariId))
      return
    }
    onPilih(b.id, b.note)
  }

  return (
    <Modal judul={P.judul} onTutup={onTutup} lebar>
      <h4 className="nbf-objek__judul">{P.cari.label}</h4>
      <div className="nbf-akum__saring">
        <div className="nbf-akum__isian nbf-labelkiri">
          <Field label={P.accumulationCode.label} value={saring.id} onChange={set('id')} />
          <Field label={P.road.label} value={saring.note} onChange={set('note')} />
          <div className="nbf-akum__dua nbf-labelkiri">
            <Field label={P.zipCode.label} value={saring.postalCode} onChange={set('postalCode')} />
            <Field label={P.czone.label} value={saring.czone} onChange={set('czone')} />
          </div>
        </div>
        <div className="nbf-akum__tombol">
          <button type="button" className="btn btn--sm" onClick={() => void filter()} disabled={memuat}>
            {P.filter.label}
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              setSaring(SARING_KOSONG)
              setHasil(null)
              setPesan('')
              setGalat(null)
            }}
          >
            {P.bersih.label}
          </button>
        </div>
      </div>
      {pesan && <div className="alert alert--error">{pesan}</div>}
      {memuat && <Memuat />}
      <Gagal galat={galat} />
      {hasil !== null && !memuat && hasil.length > 0 && (
        <div className="table-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {P.kolom.map((k) => (
                  <th key={k.sel} scope="col">
                    {k.label}
                  </th>
                ))}
                <th scope="col" />
              </tr>
            </thead>
            <tbody>
              {hasil.map((b) => (
                <tr key={b.id}>
                  <td>{b.id}</td>
                  <td>{b.accumulationName}</td>
                  <td>{b.note}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => pilih(b)}>
                      {P.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {hasil !== null && hasil.length === 0 && !memuat && <Kosong pesan={TEKS_AKUMULASI.tanpaHasil} />}
    </Modal>
  )
}
