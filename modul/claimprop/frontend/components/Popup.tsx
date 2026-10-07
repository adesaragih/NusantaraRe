// Pop-up harness Claim Prop (showHarness Target=popup): MasterTreatyIn, ListPolicyNoTreaty_Harness, CauseofLoss_Harness,
// CatastrofeList (local action), SummaryOutSClaim. Judul kolom VERBATIM section XML; baris dibaca dari server dan
// pilihan dikirim balik sebagai aksi (server membaca ulang barisnya).

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pilihanKasus } from '../api'
import { CP } from '../labels'
import { tampilAngka } from '../nilai'

export type JenisPopup = 'master' | 'polis' | 'sebab' | 'katastrofe' | 'ringkasanOS'

interface BarisMaster {
  treatyId: string
  classOfBusiness: string
  classOfBusinessId: string
  treatyContractName: string
  sob: string
  ceding: string
  treatyType: string
  proportionType: string
  treatyGroup: string
  treatyGroupId: string
  treatyYear: string
}

interface BarisPolis {
  policyNo: string
  quarter: string
  sourceOfBusinessName: string
  treatyGroup: string
}

interface BarisSebab {
  id: string
  description: string
}

interface BarisKatastrofe {
  id: string
  stsKatastrofe: string
  nonKatastrofeType: string
  note: string
}

interface Ringkasan {
  baris: {
    claimNo: string
    policyNo: string
    acceptedNo: string
    currency: string
    treatyGross: string
    rnmShare: string
    kursValue: string
    convertValue: string
  }[]
  totalLabel: string
  totalNilai: string
}

const JUDUL: Record<JenisPopup, string> = {
  master: CP.popMaster,
  polis: CP.popPolis,
  sebab: CP.popSebab,
  katastrofe: CP.popKatastrofe,
  ringkasanOS: CP.popRingkasan,
}

export default function Popup({
  id,
  jenis,
  pelaku,
  sts,
  onPilih,
  onTutup,
}: {
  id: string
  jenis: JenisPopup
  pelaku: string
  /** Catastrophe / NonKatastrofeType klaim (form katastrofe baru, read-only). */
  sts: { sts: string; non: string }
  onPilih: (aksi: string, param: string) => void
  onTutup: () => void
}) {
  const [cari, setCari] = useState('')
  const [data, setData] = useState<unknown>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [formBaru, setFormBaru] = useState(false)
  const [catatan, setCatatan] = useState('')

  useEffect(() => {
    const t = setTimeout(() => {
      pilihanKasus<unknown>(id, jenis, 0, cari).then(
        (d) => {
          setData(d)
          setGalat(null)
        },
        (g: unknown) => setGalat(g),
      )
    }, 250)
    return () => clearTimeout(t)
  }, [id, jenis, cari])

  const pakaiCari = jenis === 'master' || jenis === 'sebab' || jenis === 'katastrofe'
  let isi
  if (galat) isi = <Gagal galat={galat} />
  else if (data === null) isi = <Memuat pesan={CP.memuat} />
  else if (jenis === 'master') {
    isi = (
      <table className="claimprop__tabel">
        <thead>
          <tr>
            {[
              '',
              'Treaty ID',
              'Class of Business',
              'Contract Name',
              'Source of Business',
              'Insured Name',
              'Treaty Type',
              'Proportion Type',
              'Treaty Group',
              'Treaty Year',
            ].map((k, i) => (
              <th key={i}>{k}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {(data as BarisMaster[]).map((b) => (
            <tr key={b.treatyId + b.treatyGroupId + b.classOfBusinessId}>
              <td>
                <button
                  type="button"
                  className="btn btn--sm btn--primary"
                  onClick={() => onPilih('SetValueToClaim', `${b.treatyId}|${b.treatyGroupId}|${b.classOfBusinessId}`)}
                >
                  {CP.pilih}
                </button>
              </td>
              <td>{b.treatyId}</td>
              <td>{b.classOfBusiness}</td>
              <td>{b.treatyContractName}</td>
              <td>{b.sob}</td>
              <td>{b.ceding}</td>
              <td>{b.treatyType}</td>
              <td>{b.proportionType}</td>
              <td>{b.treatyGroup}</td>
              <td>{b.treatyYear}</td>
            </tr>
          ))}
        </tbody>
      </table>
    )
  } else if (jenis === 'polis') {
    isi = (
      <table className="claimprop__tabel">
        <thead>
          <tr>
            <th>No Polis</th>
            <th>Quarter</th>
            <th>SOB</th>
            <th>Treaty Group</th>
          </tr>
        </thead>
        <tbody>
          {(data as BarisPolis[]).map((b) => (
            <tr
              key={b.policyNo}
              className="claimprop__baris-rinci"
              onDoubleClick={() => onPilih('CheckNoPolicy', b.policyNo)}
            >
              <td>{b.policyNo}</td>
              <td>{b.quarter}</td>
              <td>{b.sourceOfBusinessName}</td>
              <td>{b.treatyGroup}</td>
            </tr>
          ))}
        </tbody>
      </table>
    )
  } else if (jenis === 'sebab') {
    isi = (
      <table className="claimprop__tabel">
        <thead>
          <tr>
            <th>Cause of Loss</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {(data as BarisSebab[]).map((b) => (
            <tr key={b.id}>
              <td>{b.description}</td>
              <td>
                <button
                  type="button"
                  className="btn btn--sm btn--primary"
                  onClick={() => onPilih('GetNameCauseofLoss', b.id)}
                >
                  {CP.pilih}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    )
  } else if (jenis === 'katastrofe') {
    isi = formBaru ? (
      <div className="claimprop__form-katastrofe">
        <label className="claimprop__medan">
          <span className="claimprop__nama-medan">Catastrophe</span>
          <span className="claimprop__tampil">{sts.sts}</span>
        </label>
        {sts.sts === 'Non-Catastrophe' && (
          <label className="claimprop__medan">
            <span className="claimprop__nama-medan" />
            <span className="claimprop__tampil">{sts.non}</span>
          </label>
        )}
        <label className="claimprop__medan">
          <span className="claimprop__nama-medan">{CP.note}</span>
          <textarea className="claimprop__area" value={catatan} onChange={(e) => setCatatan(e.target.value)} />
        </label>
        <label className="claimprop__medan">
          <span className="claimprop__nama-medan">{CP.userInput}</span>
          <span className="claimprop__tampil">{pelaku}</span>
        </label>
        <div className="claimprop__grid-alat">
          <button type="button" className="btn btn--sm" onClick={() => setFormBaru(false)}>
            {CP.cancel}
          </button>
          <button type="button" className="btn btn--sm btn--primary" onClick={() => onPilih('SaveCatasrtope', catatan)}>
            {CP.save}
          </button>
        </div>
      </div>
    ) : (
      <>
        <div className="claimprop__grid-alat">
          <button type="button" className="btn btn--sm" onClick={() => setFormBaru(true)}>
            {CP.addNew}
          </button>
        </div>
        <table className="claimprop__tabel">
          <thead>
            <tr>
              <th />
              <th />
              <th>{CP.note}</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {(data as BarisKatastrofe[]).map((b) => (
              <tr key={b.id}>
                <td>{b.stsKatastrofe}</td>
                <td>{b.nonKatastrofeType}</td>
                <td>{b.note}</td>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() => onPilih('SetCatastrope', b.id)}
                  >
                    {CP.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else {
    const r = data as Ringkasan
    isi = (
      <>
        <div className="claimprop__label">Summary Outstanding Claim</div>
        <table className="claimprop__tabel">
          <thead>
            <tr>
              {[
                'Claim No',
                'Policy No',
                'Accepted No',
                'Currency',
                'Treaty Gross  100%',
                'RNM Share Value',
                'Kurs Value',
                'Convert Value',
              ].map((k) => (
                <th key={k}>{k}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {r.baris.map((b, i) => (
              <tr key={i}>
                <td>{b.claimNo}</td>
                <td>{b.policyNo}</td>
                <td>{b.acceptedNo}</td>
                <td>{b.currency}</td>
                <td className="claimprop__angka">{tampilAngka(b.treatyGross)}</td>
                <td className="claimprop__angka">{tampilAngka(b.rnmShare)}</td>
                <td className="claimprop__angka">{tampilAngka(b.kursValue)}</td>
                <td className="claimprop__angka">{tampilAngka(b.convertValue)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <div className="claimprop__grid-kaki">
          <span className="claimprop__label">{CP.totalOutstanding}</span>
          {r.totalLabel !== '' && <span className="claimprop__tampil">{r.totalLabel}</span>}
          <span className="claimprop__angka">{tampilAngka(r.totalNilai)}</span>
        </div>
      </>
    )
  }

  return (
    <Modal
      judul={JUDUL[jenis]}
      onTutup={onTutup}
      labelBatal={CP.tutup}
      penuh={jenis === 'master' || jenis === 'ringkasanOS'}
    >
      {pakaiCari && !formBaru && (
        <input
          className="claimprop__input claimprop__cari"
          placeholder={CP.cariPopup}
          value={cari}
          onChange={(e) => setCari(e.target.value)}
        />
      )}
      {isi}
    </Modal>
  )
}
