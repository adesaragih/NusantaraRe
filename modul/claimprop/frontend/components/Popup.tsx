// Pop-up harness Claim Prop (showHarness Target=popup): MasterTreatyIn, ListPolicyNoTreaty_Harness, CauseofLoss_Harness,
// CatastrofeList (local action), SummaryOutSClaim. Judul kolom VERBATIM section XML; baris dibaca dari server dan
// pilihan dikirim balik sebagai aksi (server membaca ulang barisnya).

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pilihanKasus } from '../api'
import { CP } from '../labels'
import { tampilAngka } from '../nilai'
import { kunciMaster } from './kunciBaris'
import { jumlahHalaman, potongHalaman } from './susun'

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

/** Kolom popup master: judul VERBATIM section MasterTreatyInList; `kunci` = filter per kolom (parameter kueri). */
const KOLOM_MASTER: { label: string; kunci?: string; nilai: (b: BarisMaster) => string }[] = [
  { label: 'Treaty ID', kunci: 'treatyId', nilai: (b) => b.treatyId },
  { label: 'Class of Business', kunci: 'classOfBusiness', nilai: (b) => b.classOfBusiness },
  { label: 'Contract Name', kunci: 'contractName', nilai: (b) => b.treatyContractName },
  { label: 'Source of Business', kunci: 'sob', nilai: (b) => b.sob },
  { label: 'Insured Name', kunci: 'insuredName', nilai: (b) => b.ceding },
  { label: 'Treaty Type', kunci: 'treatyType', nilai: (b) => b.treatyType },
  { label: 'Proportion Type', nilai: (b) => b.proportionType },
  { label: 'Treaty Group', kunci: 'treatyGroup', nilai: (b) => b.treatyGroup },
  { label: 'Treaty Year', kunci: 'treatyYear', nilai: (b) => b.treatyYear },
]
const PER_HALAMAN_MASTER = 50
const BATAS_MASTER = 500 // models.BatasMaster

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
  const [saring, setSaring] = useState<Record<string, string>>({})
  const [hal, setHal] = useState(1)

  useEffect(() => {
    // Jawaban basi diabaikan: permintaan saringan lama (lebih lebar, lebih lambat) tidak boleh menimpa hasil saringan
    // terbaru (temuan work owner 08-10-2026: filter Treaty ID menampilkan treaty lain).
    let aktif = true
    const t = setTimeout(() => {
      pilihanKasus<unknown>(id, jenis, 0, cari, saring).then(
        (d) => {
          if (!aktif) return
          setData(d)
          setGalat(null)
        },
        (g: unknown) => {
          if (aktif) setGalat(g)
        },
      )
    }, 300)
    return () => {
      aktif = false
      clearTimeout(t)
    }
  }, [id, jenis, cari, saring])

  const pakaiCari = jenis === 'sebab' || jenis === 'katastrofe'
  let isi
  if (galat) isi = <Gagal galat={galat} />
  else if (data === null) isi = <Memuat pesan={CP.memuat} />
  else if (jenis === 'master') {
    const semua = data as BarisMaster[]
    const nHal = jumlahHalaman(semua.length, PER_HALAMAN_MASTER)
    const halIni = Math.min(hal, nHal)
    const awal = (halIni - 1) * PER_HALAMAN_MASTER
    isi = (
      <>
        <div className="claimprop__pager">
          <span>
            {CP.menampilkan} {semua.length === 0 ? 0 : awal + 1}–{Math.min(awal + PER_HALAMAN_MASTER, semua.length)}{' '}
            {CP.dari} {semua.length}
            {semua.length >= BATAS_MASTER && ` (${CP.batasMaster})`}
          </span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni <= 1}
            onClick={() => setHal(halIni - 1)}
            aria-label="Previous"
          >
            ‹
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni >= nHal}
            onClick={() => setHal(halIni + 1)}
            aria-label="Next"
          >
            ›
          </button>
        </div>
        <table className="claimprop__tabel">
          <thead>
            <tr>
              <th />
              {KOLOM_MASTER.map((k) => (
                <th key={k.label}>{k.label}</th>
              ))}
            </tr>
            <tr>
              <th />
              {KOLOM_MASTER.map((k) => (
                <th key={k.label}>
                  {k.kunci && (
                    <input
                      className="field__input claimprop__input--sel"
                      placeholder={CP.saring}
                      aria-label={`${CP.saring} ${k.label}`}
                      value={saring[k.kunci] ?? ''}
                      onChange={(e) => {
                        const kunci = k.kunci ?? ''
                        setSaring((s) => ({ ...s, [kunci]: e.target.value }))
                        setHal(1)
                      }}
                    />
                  )}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {semua.length === 0 && (
              <tr>
                <td className="muted" colSpan={KOLOM_MASTER.length + 1}>
                  —
                </td>
              </tr>
            )}
            {potongHalaman(semua, PER_HALAMAN_MASTER, halIni).map((b) => (
              <tr key={kunciMaster(b)}>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() =>
                      onPilih('SetValueToClaim', `${b.treatyId}|${b.treatyGroupId}|${b.classOfBusinessId}`)
                    }
                  >
                    {CP.pilih}
                  </button>
                </td>
                {KOLOM_MASTER.map((k) => (
                  <td key={k.label}>{k.nilai(b)}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else if (jenis === 'polis') {
    isi = (
      <table className="claimprop__tabel claimprop__tabel--rapat">
        <thead>
          <tr>
            <th>No Polis</th>
            <th>Quarter</th>
            <th>SOB</th>
            <th>Treaty Group</th>
          </tr>
        </thead>
        <tbody>
          {(data as BarisPolis[]).map((b, i) => (
            <tr
              key={`${b.policyNo}|${b.quarter}|${i}`}
              className="inbox__baris"
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
      <div className="form-grid">
        <label className="field">
          <span className="field__label">Catastrophe</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.sts} />
        </label>
        {sts.sts === 'Non-Catastrophe' && (
          <label className="field">
            <span className="field__label">&nbsp;</span>
            <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.non} />
          </label>
        )}
        <label className="field field--lebar">
          <span className="field__label">{CP.note}</span>
          <textarea
            className="field__input claimprop__area"
            value={catatan}
            onChange={(e) => setCatatan(e.target.value)}
          />
        </label>
        <label className="field">
          <span className="field__label">{CP.userInput}</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={pelaku} />
        </label>
        <div className="field--lebar claimprop__tombol claimprop__tombol--akhir">
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setFormBaru(false)}>
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
      lebar={jenis !== 'master' && jenis !== 'ringkasanOS'}
    >
      {pakaiCari && !formBaru && (
        <input
          className="field__input claimprop__cari"
          placeholder={CP.cariPopup}
          value={cari}
          onChange={(e) => setCari(e.target.value)}
        />
      )}
      {isi}
    </Modal>
  )
}
