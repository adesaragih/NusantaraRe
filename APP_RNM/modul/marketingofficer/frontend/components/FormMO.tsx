// Form tambah dan ubah Marketing Officer - padanan form Pega `InputMarketingOfficer` (caption VERBATIM: Code, Name
// Marketing, Set as a leader, Branch, Sub Branch, Leader, Active) ditambah Login Account dari `M_LOGIN_GO`.
// Code, Name Marketing, dan Marketing Code hanya ditampilkan: backend yang mengisinya.

import { useEffect, useState } from 'react'

import { Field, Gagal, Memuat, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { ambilPilihan, tambah, ubah, type BarisMO, type Isian, type MarketingOfficer, type Pilihan } from '../api'
import { isianDari, isianKosong, leaderUntuk, periksa, setelBranch, setelSubBranch, subBranchUntuk } from '../aturan'
import { MO } from '../labels'

const tetap = () => undefined

export default function FormMO({
  baris,
  onTutup,
  onTersimpan,
}: {
  /** `null` = tambah. */
  baris: BarisMO | null
  onTutup: () => void
  onTersimpan: (m: MarketingOfficer) => void
}) {
  const baru = baris === null
  const [pilihan, setPilihan] = useState<Pilihan | null>(null)
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [isi, setIsi] = useState<Isian>(() => (baris === null ? isianKosong() : isianDari(baris)))
  const [galatLokal, setGalatLokal] = useState<string | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let hidup = true
    ambilPilihan().then(
      (p) => {
        if (hidup) setPilihan(p)
      },
      (g: unknown) => {
        if (hidup) setGalatMuat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [])

  const ubahIsi = (sebagian: Partial<Isian>) => {
    setIsi((x) => ({ ...x, ...sebagian }))
  }

  const kirim = () => {
    if (sibuk || pilihan === null) return
    const awal = periksa(isi, baru)
    setGalatLokal(awal)
    setGalatSimpan(null)
    if (awal !== null) return
    setSibuk(true)
    const janji = baris === null ? tambah(isi) : ubah(baris.id, isi)
    janji.then(
      (m) => {
        setSibuk(false)
        onTersimpan(m)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatSimpan(g)
      },
    )
  }

  const akun = pilihan?.akun.find((a) => a.loginId === isi.aksesLogin)
  const sub = pilihan?.subBranch.find((c) => c.id === isi.branchDetailId)
  const opsiAkun: Opsi[] = (pilihan?.akun ?? []).map((a) => ({ value: a.loginId, label: `${a.nama} (${a.loginId})` }))
  const opsiLeader: Opsi[] = pilihan === null ? [] : leaderUntuk(pilihan, baris?.id ?? null).map((l) => ({
    value: l.id,
    label: l.subBranch === '' ? `${l.nama} (${l.id})` : `${l.nama} (${l.id}, ${l.subBranch})`,
  }))
  const opsiBranch: Opsi[] = (pilihan?.branch ?? []).map((c) => ({ value: c.id, label: `${c.nama} (${c.id})` }))
  const opsiSub: Opsi[] =
    pilihan === null ? [] : subBranchUntuk(pilihan, isi.branchParent).map((c) => ({ value: c.id, label: `${c.nama} (${c.id})` }))
  const teamGroup =
    sub !== undefined ? sub.kanwilGroup : baris !== null && isi.branchDetailId === baris.branchDetailId ? baris.teamGroup : ''

  return (
    <Modal
      judul={baru ? MO.judulBaru : `${MO.judulUbah} — ${baris.id}`}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={MO.batal}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || pilihan === null}>
          {sibuk ? MO.menyimpan : MO.simpan}
        </button>
      }
    >
      {galatMuat !== null && <Gagal galat={galatMuat} />}
      {pilihan === null && galatMuat === null && <Memuat pesan={MO.memuatPilihan} />}
      {pilihan !== null && (
        <>
          {galatLokal !== null && (
            <div className="alert alert--error" role="alert">
              {galatLokal}
            </div>
          )}
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}
          <div className="form-grid">
            <Field label={MO.code} value={baris?.id ?? ''} placeholder={MO.codeOtomatis} onChange={tetap} readOnly />
            <Pilih
              label={MO.akun}
              value={isi.aksesLogin}
              onChange={(v) => {
                ubahIsi({ aksesLogin: v })
              }}
              opsi={opsiAkun}
              kosong={MO.pilih}
              required={baru}
            />
            <div>
              <Field
                label={MO.namaMarketing}
                value={baris?.clientName ?? akun?.nama ?? ''}
                onChange={tetap}
                readOnly
              />
              <p className="muted marketingofficer__catatan">{MO.catatanNama}</p>
            </div>
            <Field label={MO.email} value={akun?.email ?? baris?.emailAkun ?? ''} onChange={tetap} readOnly />
            <div>
              <Field
                label={MO.marketingCode}
                value={baris?.clientId ?? ''}
                placeholder={MO.codeOtomatis}
                onChange={tetap}
                readOnly
              />
              <p className="muted marketingofficer__catatan">{MO.catatanMarketingCode}</p>
            </div>
            <label className="marketingofficer__centang">
              <input
                type="checkbox"
                checked={isi.leader}
                onChange={(e) => {
                  ubahIsi({ leader: e.target.checked })
                }}
              />
              {MO.setLeader}
            </label>
            {!isi.leader && (
              <Pilih
                label={MO.leader}
                value={isi.leaderId}
                onChange={(v) => {
                  ubahIsi({ leaderId: v })
                }}
                opsi={opsiLeader}
                kosong={MO.pilih}
                required
              />
            )}
            <Pilih
              label={MO.branch}
              value={isi.branchParent}
              onChange={(v) => {
                setIsi((x) => setelBranch(x, pilihan, v))
              }}
              opsi={opsiBranch}
              kosong={MO.pilih}
            />
            <Pilih
              label={MO.subBranch}
              value={isi.branchDetailId}
              onChange={(v) => {
                setIsi((x) => setelSubBranch(x, pilihan, v))
              }}
              opsi={opsiSub}
              kosong={MO.pilih}
            />
            <Field label={MO.teamGroup} value={teamGroup} onChange={tetap} readOnly />
            <label className="marketingofficer__centang">
              <input
                type="checkbox"
                checked={isi.aktif}
                onChange={(e) => {
                  ubahIsi({ aktif: e.target.checked })
                }}
              />
              {MO.active}
            </label>
          </div>
        </>
      )}
    </Modal>
  )
}
