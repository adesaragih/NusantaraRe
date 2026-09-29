// Panel reinsurer pada kombinasi kontrak — tiket 05 Treaty Contract Out.
//
// Meniru `Section/ViewDetailTreatyReinsurerGrid1.xml` (disertakan
// `InputTreatyContractReinsType.xml` b13311; dibuka tombol baris kontrak
// `Reinsurer List` b11308): grid `BrowseDetailTreatyReisurer_RD` (`Add` b1730;
// kolom b2418–b3138; tombol baris `Edit` b4491, `Delete` b4936, `Security
// Reinsurer` b5277), baris `Total Share -->>` b6186, form (`ID` b7842,
// `Reins.ID` b8042, `Reinsurer` b8226, `%Share` b8522, `%Comm` b8800, `Rating`
// b9076, `Operator Name` b11100, `Save` b11405, `Error` b12131, `Informasi` b12868).
//
// ⛔ Share dan komisi TEKS sepanjang jalan — tidak pernah angka JavaScript.
// Total share datang dari server (desimal persis, ADR-0003).

import { useCallback, useEffect, useState } from 'react'

import { REINSURER_TCO, TAHUN_TCO } from '../../assets/labels.treaty-contract-out'
import {
  ambilReinsurerKombinasi,
  cariReinsurerMaster,
  simpanReinsurerKombinasi,
  type DaftarReinsurer,
  type KombinasiTreaty,
  type ReinsurerMaster,
  type ReinsurerMasuk,
  type ReinsurerTreaty,
} from '../../services/api'
import { Field, Gagal, Kosong, Memuat, Pilih } from '../ui/dasar'

/** Isian form — hanya medan yang tampil di form Pega. */
export interface FormReinsurer {
  id: string
  reinsurerId: string
  /** Nama perusahaan dari master; hanya dibaca. */
  name: string
  pctShare: string
  ricomm: string
  stdRating: string
  operatorName: string
}

/** `NewTreatyReinsurerDetail_Act` langkah 3. */
export function formReinsurerKosong(): FormReinsurer {
  return { id: '', reinsurerId: '', name: '', pctShare: '', ricomm: '', stdRating: '', operatorName: '' }
}

/** `SetUbahTreatyReinsurerList_Act` langkah 3 (baris `.ID == Param.ID`). */
export function formReinsurerDari(r: ReinsurerTreaty): FormReinsurer {
  return {
    id: r.id, reinsurerId: r.reinsurerId, name: r.name, pctShare: r.pctShare, ricomm: r.ricomm,
    stdRating: r.stdRating, operatorName: r.operatorName,
  }
}

/** Badan simpan. Share/komisi dikirim apa adanya (di-trim) — server menormalkan koma. */
export function keMasukReinsurer(f: FormReinsurer): ReinsurerMasuk {
  return {
    id: f.id,
    reinsurerId: f.reinsurerId,
    pctShare: f.pctShare.trim(),
    ricomm: f.ricomm.trim(),
    stdRating: f.stdRating.trim(),
  }
}

/** Kepala kombinasi (`OutputData.HASIL1/3/2`). */
export function labelKombinasi(k: KombinasiTreaty): string {
  return [k.treatyYear, k.treatyGroupName || k.treatyGroupId, k.reinsTypeName || k.reinsTypeId]
    .filter((v) => v.trim() !== '')
    .join(' · ')
}

export default function PanelReinsurerKombinasi({
  tahunID,
  kontrakID,
  onTutup,
}: {
  tahunID: string
  kontrakID: string
  onTutup?: () => void
}) {
  const [daftar, setDaftar] = useState<DaftarReinsurer | null>(null)
  const [form, setForm] = useState<FormReinsurer | null>(null)
  const [cari, setCari] = useState('')
  const [pilihan, setPilihan] = useState<ReinsurerMaster[]>([])
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilReinsurerKombinasi(tahunID, kontrakID))
    } catch (e) {
      setGalat(e)
    }
  }, [tahunID, kontrakID])

  useEffect(() => {
    void muat()
  }, [muat])

  function buka(f: FormReinsurer): void {
    setGalat(null)
    setInfo(null)
    setCari('')
    setPilihan([])
    setForm(f)
  }

  function ubah<K extends keyof FormReinsurer>(k: K) {
    return (v: string) => {
      setForm((f) => (f === null ? f : { ...f, [k]: v }))
    }
  }

  function cariMaster(teks: string): void {
    setCari(teks)
    if (teks.trim().length < 2) {
      setPilihan([])
      return
    }
    cariReinsurerMaster(teks)
      .then(setPilihan)
      .catch((e: unknown) => {
        setGalat(e)
      })
  }

  function pilihMaster(id: string): void {
    const m = pilihan.find((p) => p.id === id)
    setForm((f) => (f === null ? f : { ...f, reinsurerId: id, name: m?.clientName ?? '' }))
  }

  async function simpan(): Promise<void> {
    if (form === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      const h = await simpanReinsurerKombinasi(tahunID, kontrakID, keMasukReinsurer(form))
      buka(formReinsurerDari(h.reinsurer))
      setInfo(REINSURER_TCO.tersimpan)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <section className="panel">
      <header className="inbox__kepala">
        <h3 className="panel__title">
          {REINSURER_TCO.kombinasi}: {daftar === null ? '…' : labelKombinasi(daftar.kombinasi)}
        </h3>
        {onTutup !== undefined && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {REINSURER_TCO.tutup}
          </button>
        )}
      </header>

      {galat !== null && (
        <div>
          <strong>{REINSURER_TCO.error}</strong>
          <Gagal galat={galat} />
        </div>
      )}
      {info !== null && (
        <p role="status">
          {REINSURER_TCO.informasi}: {info}
        </p>
      )}

      {form !== null && (
        <div className="panel">
          <div className="form-grid">
            <Field label={REINSURER_TCO.formId} value={form.id} onChange={() => undefined} readOnly />
            <Field label={REINSURER_TCO.formReinsId} value={form.reinsurerId} onChange={() => undefined} readOnly />
            <Field label={REINSURER_TCO.cariReinsurer} value={cari} onChange={cariMaster} />
            <Pilih
              label={REINSURER_TCO.formReinsurer}
              value={form.reinsurerId}
              onChange={pilihMaster}
              opsi={
                pilihan.length > 0
                  ? pilihan.map((p) => ({ value: p.id, label: p.clientName }))
                  : form.reinsurerId === ''
                    ? []
                    : [{ value: form.reinsurerId, label: form.name }]
              }
              required
            />
            <Field label={REINSURER_TCO.formShare} value={form.pctShare} onChange={ubah('pctShare')} required />
            <Field label={REINSURER_TCO.formComm} value={form.ricomm} onChange={ubah('ricomm')} required />
            <Field label={REINSURER_TCO.formRating} value={form.stdRating} onChange={ubah('stdRating')} />
            <Field label={REINSURER_TCO.formOperatorName} value={form.operatorName} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void simpan()}>
              {REINSURER_TCO.save}
            </button>
          </div>
        </div>
      )}

      <div className="aksi-baris">
        <button type="button" className="btn btn--primary" onClick={() => buka(formReinsurerKosong())}>
          {REINSURER_TCO.add}
        </button>
      </div>
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && daftar.daftar.length === 0 && <Kosong pesan={REINSURER_TCO.kosong} />}
      {daftar !== null && daftar.daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{REINSURER_TCO.kolomReinsId}</th>
              <th>{REINSURER_TCO.kolomReinsurer}</th>
              <th>{REINSURER_TCO.kolomShare}</th>
              <th>{REINSURER_TCO.kolomComm}</th>
              <th>{REINSURER_TCO.kolomRating}</th>
              <th>{REINSURER_TCO.kolomOperatorName}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.daftar.map((r) => (
              <tr key={r.id} className="inbox__baris">
                <td>{r.reinsurerId}</td>
                <td>{r.name}</td>
                <td>{r.pctShare}</td>
                <td>{r.ricomm}</td>
                <td>{r.stdRating}</td>
                <td>{r.operatorName}</td>
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formReinsurerDari(r))}>
                    {REINSURER_TCO.edit}
                  </button>{' '}
                  <button type="button" className="btn btn--ghost btn--sm" disabled title={`${TAHUN_TCO.menungguTiket} 10`}>
                    {REINSURER_TCO.delete}
                  </button>{' '}
                  <button type="button" className="btn btn--ghost btn--sm" disabled title={`${TAHUN_TCO.menungguTiket} 06`}>
                    {REINSURER_TCO.securityReinsurer}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr>
              <td colSpan={2}>{REINSURER_TCO.totalShare}</td>
              <td>{daftar.totalShare}</td>
              <td colSpan={4} />
            </tr>
          </tfoot>
        </table>
      )}
    </section>
  )
}
