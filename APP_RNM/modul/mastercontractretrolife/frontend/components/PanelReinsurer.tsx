// Panel reinsurer - tiket 05. Harness `InboxRetroLifeReinsurersList` → section
// `InputSecurityLifeReinsurers.xml`, dibuka tombol baris `Reinsurer List` panel kontrak
// (`InputRetroLimitReinsurers.xml` b14034 + `CountingPercentShare_Act`).
//
// Kepala b1029/b1245/b1446 (⚠️ `ID Reins Type` berisi ID KONTRAK - label VERBATIM). Grid RD
// `BrowseDetailTreatyReisurerLife_RD` (dua saringan induk, urut `.ID DESC` - server). Form (wadah
// b1972): `REINSURER NAME` autocomplete b3746 (`REINS ID` b4570 MATI `1 = 2` - tidak dirender; nilainya
// diisi autocomplete), `(%) SHARE` b4784, `(%) DISCOUNT` b5072, `(%) OVR COMM` b5364, `Save` b5656,
// `Cancel` b5942. Baris: `Edit` b12099, `Security Reinsurer` b12441, `Delete` b13422. `Add` b8689.
// `Total Share -->>` b14339 = jumlah `PCTSHARE` dari server (desimal persis).
//
// ⛔ Tiket 05 AC 18/19: total ≠ 100 DITANDAI mencolok, tidak memblokir penyimpanan.

import { useCallback, useEffect, useState } from 'react'

import { ambilReinsurer, simpanReinsurer, type JawabanReinsurer, type Kontrak, type Reinsurer, type ReinsurerMasuk } from '../api'
import { REINSURER_MCRL, UMUM_MCRL } from '../labels'
import { jepitHalaman, operatorKini, potongHalaman, sel, selAngka, selWaktu } from '../tampilan'
import { Field, Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { KepalaPanel, Penomoran } from './Bingkai'
import { useCariReinsurer } from './cariMaster'
import KonfirmasiHapus from './KonfirmasiHapus'
import PanelSecurity from './PanelSecurity'
import { useHapus } from './useHapus'

/** Isian form - `InputSecurityLife.*`. */
export interface FormReinsurer {
  id: string
  reinsurerId: string
  reinsurerName: string
  pctShare: string
  /** Kolom `COMMISION` - label `(%) DISCOUNT`. */
  komisi: string
  ovrComm: string
  /** `Inputor` (ro). */
  userId: string
}

/** `NewInputSecurityLife_Act`: 1 b233 `Page-New` (semua kosong); `Inputor` = operator (`DATASHOW2 = 0`, b2503). */
export function formReinsurerBaru(operator: string): FormReinsurer {
  return { id: '', reinsurerId: '', reinsurerName: '', pctShare: '', komisi: '', ovrComm: '', userId: operator }
}

/**
 * `SetSecurityLife_Act` b513: `USERID ← OperatorID`; salin `NAME ← .REINSURERNAME`,
 * `REINSURERID_LIFE ← .REINSURERID`, `PCTSHARE`, `COMMISION`, `OVR_COMM`, `ID`.
 */
export function formReinsurerDari(r: Reinsurer, operator: string): FormReinsurer {
  return {
    id: r.id, reinsurerId: r.reinsurerId, reinsurerName: r.reinsurerName, pctShare: r.pctShare,
    komisi: r.komisi, ovrComm: r.ovrComm, userId: operator,
  }
}

/** Badan simpan; angka dikirim apa adanya (di-trim) - server menormalkan koma dan memeriksa 0..100. */
export function keReinsurerMasuk(f: FormReinsurer): ReinsurerMasuk {
  return {
    id: f.id, reinsurerName: f.reinsurerName, reinsurerId: f.reinsurerId,
    pctShare: f.pctShare.trim(), komisi: f.komisi.trim(), ovrComm: f.ovrComm.trim(),
  }
}

/** Kelas baris total - mencolok bila total ≠ 100 (tiket 05 AC 19). */
export function kelasTotal(totalBukan100: boolean): string {
  return totalBukan100 ? 'mcrl-total mcrl-total--bukan100' : 'mcrl-total'
}

export default function PanelReinsurer({ kontrak, onTutup }: { kontrak: Kontrak; onTutup: () => void }) {
  const [jawab, setJawab] = useState<JawabanReinsurer | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  const [form, setForm] = useState<FormReinsurer | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [pesan, setPesan] = useState<string | null>(null)
  const [security, setSecurity] = useState<Reinsurer | null>(null)
  const master = useCariReinsurer(setGalatForm)

  const muat = useCallback(async () => {
    try {
      const j = await ambilReinsurer(kontrak.id)
      setJawab(j)
      setGalat(null)
      setHalaman((h) => jepitHalaman(h, j.daftar.length))
    } catch (e) {
      setGalat(e)
    }
  }, [kontrak.id])

  useEffect(() => {
    void muat()
  }, [muat])

  const hapusan = useHapus('reinsurer', async (id, p) => {
    if (form?.id === id) setForm(null)
    setPesan(p)
    await muat()
  })

  function buka(f: FormReinsurer): void {
    setGalatForm(null)
    setPesan(null)
    master.reset()
    setForm(f)
  }

  function ubah(k: 'pctShare' | 'komisi' | 'ovrComm') {
    return (v: string) => {
      setForm((f) => (f === null ? f : { ...f, [k]: v }))
    }
  }

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatForm(null)
    try {
      await simpanReinsurer(kontrak.id, keReinsurerMasuk(form))
      // `SaveSecurityLife_Act` langkah 6 b1335: `DATASHOW = ""` - form tertutup; total dihitung ulang.
      setForm(null)
      await muat()
    } catch (e) {
      setGalatForm(e)
    } finally {
      setMenyimpan(false)
    }
  }

  // Popup `Security Reinsurer` di atas panel ini: panel ini tidak dirender, keadaannya tetap.
  if (security !== null) {
    return (
      <PanelSecurity
        key={security.id}
        reinsurer={security}
        onTutup={() => {
          setSecurity(null)
        }}
      />
    )
  }

  const induk = jawab?.kontrak ?? kontrak
  const daftar = jawab?.daftar ?? []

  return (
    <section className="panel">
      <KepalaPanel
        judul={REINSURER_MCRL.judul}
        medan={[
          [REINSURER_MCRL.idTreatyYear, induk.idTreatyYear],
          [REINSURER_MCRL.idReinsType, induk.id],
          [REINSURER_MCRL.reinsType, induk.reinsTypeName],
        ]}
        onTutup={onTutup}
      />

      {form !== null && (
        <div className="panel">
          {galatForm !== null && <Gagal galat={galatForm} />}
          <div className="form-grid">
            <Field label={REINSURER_MCRL.formInputor} value={form.userId} onChange={() => undefined} readOnly />
            <PilihSaring
              label={REINSURER_MCRL.formReinsurerName}
              value={form.reinsurerId}
              teksTerpilih={form.reinsurerName || form.reinsurerId}
              opsi={master.pilihan}
              memuat={master.memuat}
              onCari={master.cari}
              onPilih={(o) => {
                setForm((f) => (f === null ? f : { ...f, reinsurerId: o.value, reinsurerName: o.label }))
              }}
              required
            />
            <Field label={REINSURER_MCRL.formShare} value={form.pctShare} onChange={ubah('pctShare')} required />
            <Field label={REINSURER_MCRL.formDiscount} value={form.komisi} onChange={ubah('komisi')} required />
            <Field label={REINSURER_MCRL.formOvrComm} value={form.ovrComm} onChange={ubah('ovrComm')} required />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {REINSURER_MCRL.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatForm(null)
              }}
            >
              {REINSURER_MCRL.cancel}
            </button>
          </div>
        </div>
      )}

      {pesan !== null && <p role="status">{pesan}</p>}

      <div className="aksi-baris mcrl-aksi-grid">
        <button type="button" className="btn btn--primary" onClick={() => buka(formReinsurerBaru(operatorKini()))}>
          {REINSURER_MCRL.add}
        </button>
        <Penomoran halaman={halaman} total={daftar.length} onPindah={setHalaman} />
      </div>
      {jawab === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {jawab !== null && daftar.length === 0 && <Kosong pesan={UMUM_MCRL.kosong} />}
      {daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{REINSURER_MCRL.kolomId}</th>
              <th>{REINSURER_MCRL.kolomReinsurerName}</th>
              <th>{REINSURER_MCRL.kolomShare}</th>
              <th>{REINSURER_MCRL.kolomDiscount}</th>
              <th>{REINSURER_MCRL.kolomOvrComm}</th>
              <th>{REINSURER_MCRL.kolomInputor}</th>
              <th>{REINSURER_MCRL.kolomUpdateDate}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {potongHalaman(daftar, halaman).map((r) => (
              <tr key={r.id} className="inbox__baris">
                <td>{sel(r.id)}</td>
                <td>{sel(r.reinsurerName)}</td>
                <td>{selAngka(r.pctShare)}</td>
                <td>{selAngka(r.komisi)}</td>
                <td>{selAngka(r.ovrComm)}</td>
                <td>{sel(r.userId)}</td>
                <td>{selWaktu(r.tglUpdate)}</td>
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formReinsurerDari(r, operatorKini()))}>
                    {REINSURER_MCRL.edit}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setSecurity(r)
                    }}
                  >
                    {REINSURER_MCRL.securityReinsurer}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setPesan(null)
                      hapusan.minta(r.id, r.reinsurerName || r.id)
                    }}
                  >
                    {REINSURER_MCRL.delete}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {jawab !== null && (
        <p className={kelasTotal(jawab.totalBukan100)} role={jawab.totalBukan100 ? 'alert' : undefined}>
          <span>{REINSURER_MCRL.totalShare}</span> <strong>{selAngka(jawab.totalShare)}</strong>
          {jawab.totalBukan100 && <span className="mcrl-total__tanda">{REINSURER_MCRL.totalBukan100}</span>}
        </p>
      )}

      {hapusan.konfirmasi !== null && (
        <KonfirmasiHapus
          judul={REINSURER_MCRL.delete}
          labelBatal={REINSURER_MCRL.cancel}
          jenis="reinsurer"
          nama={hapusan.konfirmasi.nama}
          dampak={hapusan.konfirmasi.dampak}
          galat={hapusan.konfirmasi.galat}
          sibuk={hapusan.sibuk}
          onYa={() => void hapusan.ya()}
          onBatal={hapusan.batal}
        />
      )}
    </section>
  )
}
