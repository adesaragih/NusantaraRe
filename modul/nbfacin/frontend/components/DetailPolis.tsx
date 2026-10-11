// Detail polis bersama (tiket 71 / 70): dipakai menu View Policy dan View Details popup Following Policy Number. Versi dari
// kolom relasional JSON_POLIS; versi sistem baru = layar case mode hanya lihat (render-prop), polis lama Pega = ringkasan
// FACINPRODUCTION / CEDING_FACINPRODUCTION. Tanpa DATA_JSON (K70-1, K71-1).

import { useEffect, useState, type ReactNode } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { angka4 } from '../angka4'
import { ringkasanPolis, versiPolis, type BarisFollowing, type RingkasanPolis, type VersiPolis } from '../api'
import { FOLLOWING as F, TEKS_FOLLOWING, TEKS_VIEW_POLICY as T, VIEW_POLICY as V } from '../labels'
import { tanggalFollowing } from './PopupFollowing'

/** Policy Type kode -> label (`D:\migrasi\RNM\DDL\PolicyType.xml`: 0 Individual, 1 Master, 2 Declaration). */
export const labelPolicyType = (kode: string) =>
  ({ '0': 'Individual Policy', '1': 'Master Policy', '2': 'Declaration Policy' })[kode] ?? kode

/** Ringkasan polis lama Pega (FACINPRODUCTION / CEDING_FACINPRODUCTION). */
export function RingkasanPolisLama({ v }: { v: RingkasanPolis }) {
  const k = v.kepala
  const K = V.kepala
  const kepala: [string, string][] = [
    [K.policyNo, k.policyNo],
    [K.noEndors, k.noEndors],
    [K.caseId, k.caseId],
    [K.insuredName, k.insuredName],
    [K.qqName, k.qqName],
    [K.sourceOfBusiness, k.sourceOfBusiness],
    [K.cedingCo, k.cedingCo],
    [K.periode, [tanggalFollowing(k.startPeriod), tanggalFollowing(k.endPeriod)].filter((x) => x !== '').join(' – ')],
    [K.policyType, labelPolicyType(k.policyType)],
    [K.reffNumber, k.reffNumber],
    [K.groupName, k.groupName],
    [K.statusBusiness, k.statusBusiness],
    [K.noMasterPolis, k.noMasterPolis],
  ]
  return (
    <div className="nbf-vp__ringkasan">
      <p className="muted">{T.polisLama}</p>
      <dl className="nbf-vp__kepala">
        {kepala.map(([l, n]) => (
          <div key={l} className="nbf-vp__medan">
            <dt>{l}</dt>
            <dd>{n || '—'}</dd>
          </div>
        ))}
      </dl>
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {V.kolomRingkasan.map((c) => (
                <th key={c} scope="col">
                  {c}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {v.baris.length === 0 && (
              <tr>
                <td colSpan={V.kolomRingkasan.length}>{TEKS_FOLLOWING.kosong}</td>
              </tr>
            )}
            {v.baris.map((b, n) => (
              <tr key={n}>
                <td>{b.objectNo}</td>
                <td>{b.riskLocation}</td>
                <td>{b.objectItem}</td>
                <td>{b.coverage}</td>
                <td>{b.accumulationCode}</td>
                <td>{b.currency}</td>
                <td className="nbf-angka">{angka4(b.tsi100)}</td>
                <td className="nbf-angka">{angka4(b.tsiRnm)}</td>
                <td className="nbf-angka">{b.rate}</td>
                <td className="nbf-angka">{angka4(b.premiRnm)}</td>
                <td>{b.treaty}</td>
                <td className="nbf-angka">{angka4(b.pctShare)}</td>
                <td className="nbf-angka">{angka4(b.riComm)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <table className="nbf-tabel nbf-vp__ceding">
        <thead>
          <tr>
            {V.ceding.map((c) => (
              <th key={c} scope="col">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {v.ceding.map((c, n) => (
            <tr key={n}>
              <td>{c.cedingName}</td>
              <td className="nbf-angka">{angka4(c.pctShare)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/**
 * Mode detail polis (View Policy & View Details Following Policy): versi (View Data) + isi versi terpilih. Versi yang case-nya ADA
 * di sistem baru -> `tampilKasus` (layar case mode hanya lihat, dipasok pemanggil supaya berkas ini tidak mengimpor halaman case);
 * polis lama Pega -> ringkasan FACINPRODUCTION (K71-1; sampai migrasi tiket 72). Work owner 09-10-2026 "PERBAIKI, INI KENAPA":
 * View Details Following dulu selalu membuka layar case, sehingga polis lama Pega gagal dimuat ("Isi tab Object tidak dapat
 * dimuat saat ini").
 */
export function DetailPolis({
  baris,
  judul,
  onKembali,
  tampilKasus,
}: {
  baris: BarisFollowing
  /** Judul kepala di tengah; '' = tanpa judul (popup sudah berjudul). */
  judul: string
  onKembali: () => void
  tampilKasus: (caseId: string) => ReactNode
}) {
  const [versi, setVersi] = useState<VersiPolis[] | null>(null)
  const [pilih, setPilih] = useState('')
  const [ringkasan, setRingkasan] = useState<RingkasanPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [ketuk, setKetuk] = useState(0)

  useEffect(() => {
    let batal = false
    setGalat(null)
    versiPolis(baris.policyNo).then(
      (h) => {
        if (batal) return
        setVersi(h.versi)
        setPilih((p) => (h.versi.some((x) => x.prodke === p) ? p : (h.versi[0]?.prodke ?? '')))
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [baris.policyNo, ketuk])

  const kini = versi?.find((x) => x.prodke === pilih) ?? null

  useEffect(() => {
    setRingkasan(null)
    if (kini === null || kini.sistemBaru) return
    let batal = false
    ringkasanPolis(baris.policyNo, kini.prodke).then(
      (h) => {
        if (!batal) setRingkasan(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [baris.policyNo, kini?.prodke, kini?.sistemBaru, ketuk])

  return (
    <div className="nbf-vp__detail">
      {/* Kepala detail (work owner 09-10-2026): judul DI TENGAH, Back / Refresh (dan View Data bila > 1 versi) di KANAN ATAS. */}
      <header className="nbf-vp__kepala-detail">
        <span aria-hidden="true" />
        {judul !== '' ? <h2 className="inbox__judul nbf-tema__judul nbf-vp__judul-detail">{judul}</h2> : <span />}
        <div className="nbf-vp__aksi-detail">
        {versi !== null && versi.length > 1 && (
          <label className="field nbf-vp__versi">
            <span className="field__label">{V.viewData}</span>
            <select className="field__input nbf-tema__isian" value={pilih} onChange={(e) => setPilih(e.target.value)}>
              {versi.map((x) => (
                <option key={x.prodke} value={x.prodke}>
                  {x.caseId}
                </option>
              ))}
            </select>
          </label>
        )}
        <button type="button" className="btn btn--sm nbf-tema__tombol" onClick={onKembali}>
          {F.back}
        </button>
        <button type="button" className="btn btn--sm nbf-tema__tombol" onClick={() => setKetuk((n) => n + 1)}>
          {F.refresh}
        </button>
        </div>
      </header>
      <Gagal galat={galat} />
      {versi === null && !galat && <Memuat />}
      {kini !== null && kini.sistemBaru && <div key={kini.caseId}>{tampilKasus(kini.caseId)}</div>}
      {kini !== null && !kini.sistemBaru && (ringkasan ? <RingkasanPolisLama v={ringkasan} /> : !galat && <Memuat />)}
    </div>
  )
}
