// Grid security di bawah seorang reinsurer — tiket 06 Treaty Contract Out.
//
// Meniru bagian security `Section/InputTreatyContractReinsType.xml` (b14885,
// tampil bila `HASILD21 == 1`) yang dibuka tombol baris reinsurer `Security
// Reinsurer` (`ViewDetailTreatyReinsurerGrid1.xml` b5277): grid (`Reas
// Security` b16088, `Security Name` b16228, `Percent Share` b16368; `Edit`
// b17252, `Delete` b17559), `Add` b15459, form (`Security ID` b19468 baca-saja,
// `Security Name` b19648 wajib dari pemilih `BrowseAgentReinsSOA_RD`, `%Share`
// b19888), `Save` b20246.
//
// ⛔ Baris dicocokkan menurut ID — mengganti nama security mengubah baris yang
// SAMA (User story 14). ⛔ `%Share` TEKS sepanjang jalan.
//
// [keputusan work owner 30-09-2026] Isian `Cari security` DIBUANG; `Security
// Name` sendiri adalah dropdown yang dapat difilter dengan mengetik
// (`PilihSaring` + `useCariReinsurerMaster`, dicari ke server — master > 100 nama).

import { useCallback, useEffect, useState } from 'react'

import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../../inti/frontend/lib/format'
import { PilihSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { SECURITY_TCO } from '../labels'
import {
  ambilSecurity,
  hapusSecurity,
  simpanSecurity,
  type DaftarSecurity,
  type SecurityMasuk,
  type SecurityReinsurer,
} from '../api'
import { Field, Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useCariReinsurerMaster } from './cariReinsurerMaster'

/** Isian form — `Security ID` (hanya dibaca), nama dari master, `%Share`. */
export interface FormSecurity {
  id: string
  reasSecurity: string
  clientName: string
  pctShare: string
}

/** `InputNewSecurityReinsurer`: form kosong, mode sisip. */
export function formSecurityKosong(): FormSecurity {
  return { id: '', reasSecurity: '', clientName: '', pctShare: '' }
}

/** `ShowEditSecurityReinsurer`: form dari baris, mode ubah. */
export function formSecurityDari(s: SecurityReinsurer): FormSecurity {
  return { id: s.id, reasSecurity: s.reasSecurity, clientName: s.clientName, pctShare: s.pctShare }
}

/** Badan simpan — nama tampil TIDAK dikirim (server membacanya dari master). */
export function keMasukSecurity(f: FormSecurity): SecurityMasuk {
  return { id: f.id, reasSecurity: f.reasSecurity.trim(), pctShare: f.pctShare.trim() }
}

export default function PanelSecurityReinsurer({
  tahunID,
  kontrakID,
  reinsurerID,
  onTutup,
}: {
  tahunID: string
  kontrakID: string
  reinsurerID: string
  onTutup?: () => void
}) {
  const [daftar, setDaftar] = useState<DaftarSecurity | null>(null)
  const [form, setForm] = useState<FormSecurity | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)
  const master = useCariReinsurerMaster(setGalat)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilSecurity(tahunID, kontrakID, reinsurerID))
    } catch (e) {
      setGalat(e)
    }
  }, [tahunID, kontrakID, reinsurerID])

  useEffect(() => {
    void muat()
  }, [muat])

  function buka(f: FormSecurity): void {
    setGalat(null)
    setInfo(null)
    // Jawaban cari form sebelumnya tidak boleh mengisi form ini.
    master.reset()
    setForm(f)
  }

  async function simpan(): Promise<void> {
    if (form === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      await simpanSecurity(tahunID, kontrakID, reinsurerID, keMasukSecurity(form))
      // Simpan berhasil: form ditutup (keputusan work owner 30-09-2026).
      setForm(null)
      setInfo(SECURITY_TCO.tersimpan)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  async function hapus(s: SecurityReinsurer): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      setInfo(await hapusSecurity(tahunID, kontrakID, reinsurerID, s.id))
      if (form?.id === s.id) setForm(null)
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
        <h4 className="panel__title">
          {SECURITY_TCO.judul}: {daftar === null ? '…' : daftar.reinsurer.name || daftar.reinsurer.reinsurerId}
        </h4>
        {onTutup !== undefined && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {SECURITY_TCO.tutup}
          </button>
        )}
      </header>

      {galat !== null && (
        <div>
          <strong>{SECURITY_TCO.error}</strong>
          <Gagal galat={galat} />
        </div>
      )}
      {info !== null && (
        <p role="status">
          {SECURITY_TCO.informasi}: {info}
        </p>
      )}

      <div className="aksi-baris">
        <button type="button" className="btn btn--primary" onClick={() => buka(formSecurityKosong())}>
          {SECURITY_TCO.add}
        </button>
      </div>
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && daftar.daftar.length === 0 && <Kosong pesan={SECURITY_TCO.kosong} />}
      {daftar !== null && daftar.daftar.length > 0 && (
        <div className="tco-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{SECURITY_TCO.kolomReasSecurity}</th>
                <th>{SECURITY_TCO.kolomSecurityName}</th>
                <th>{SECURITY_TCO.kolomPercentShare}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {daftar.daftar.map((s) => (
                <tr key={s.id} className="inbox__baris">
                  <td>{s.reasSecurity}</td>
                  <td>{s.clientName}</td>
                  <td>{formatNumber(s.pctShare, DESIMAL_TAK_DIBATASI)}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formSecurityDari(s))}>
                      {SECURITY_TCO.edit}
                    </button>{' '}
                    <button type="button" className="btn btn--ghost btn--sm" disabled={sibuk} onClick={() => void hapus(s)}>
                      {SECURITY_TCO.delete}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {form !== null && (
        <div className="panel">
          <div className="form-grid">
            <Field label={SECURITY_TCO.formSecurityId} value={form.reasSecurity} onChange={() => undefined} readOnly />
            <PilihSaring
              label={SECURITY_TCO.formSecurityName}
              value={form.reasSecurity}
              teksTerpilih={form.clientName || form.reasSecurity}
              opsi={master.pilihan}
              memuat={master.memuat}
              onCari={master.cari}
              onPilih={(o) => {
                setForm((f) => (f === null ? f : { ...f, reasSecurity: o.value, clientName: o.label }))
              }}
              required
            />
            <Field
              label={SECURITY_TCO.formShare}
              value={form.pctShare}
              onChange={(v) => {
                setForm((f) => (f === null ? f : { ...f, pctShare: v }))
              }}
              required
            />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void simpan()}>
              {SECURITY_TCO.save}
            </button>{' '}
            <button type="button" className="btn btn--ghost" onClick={() => setForm(null)}>
              {SECURITY_TCO.cancel}
            </button>
          </div>
        </div>
      )}
    </section>
  )
}
