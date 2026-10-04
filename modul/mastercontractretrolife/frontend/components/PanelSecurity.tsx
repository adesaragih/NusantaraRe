// Panel security reinsurer - tiket 06. Harness `InboxSecurityReinsurerLife` → section
// `InputSecurityReinsurerLife.xml`, dibuka tombol baris `Security Reinsurer` panel reinsurer
// (`InputSecurityLifeReinsurers.xml` b12441).
//
// Kepala b1051/b1267/b1483 = reinsurer induk. Grid RD `BrowseSecurityReinsurer_Life_RD` (tiga saringan
// induk, dikerjakan server). Form (wadah b2038 `DATASHOW = 1`): `SECURITY REINSURER NAME`
// autocomplete b3820, `(%) SHARE` b4861, `Save` b5151, `Cancel` b5432. Baris: `Edit` b10779,
// `Delete` b11097. `Add` b8179 → form KOSONG (OQ-MCRL-12: activity Pega mengosongkan halaman
// reinsurer, bukan halaman security).
//
// ⛔ Tiket 06 AC: `(%) SHARE` = persen DARI share induk; eksposur terhadap treaty (share × share induk,
// dihitung server) tampil BERDAMPINGAN dengan label yang membedakan (OQ-MCRL-08).

import { useCallback, useEffect, useState } from 'react'

import { ambilSecurity, simpanSecurity, type JawabanSecurity, type Reinsurer, type SecurityMasuk, type SecurityReinsurer } from '../api'
import { REINSURER_MCRL, SECURITY_MCRL, UMUM_MCRL } from '../labels'
import { operatorKini, sel, selAngka, selWaktu } from '../tampilan'
import { Field, Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { KepalaPanel } from './Bingkai'
import { useCariReinsurer } from './cariMaster'
import KonfirmasiHapus from './KonfirmasiHapus'
import { kelasTotal, reinsurerTersedia } from './aturanDaftar'
import { useHapus } from './useHapus'

/** Isian form - `InputSecurityReinsurerLife.*`. */
export interface FormSecurity {
  id: string
  reinsurerId: string
  reinsurerName: string
  pctShare: string
  /** `Inputor` (ro). */
  userId: string
}

/** `Add` b8179 - form kosong (OQ-MCRL-12); `Inputor` = operator (`DATASHOW2 = 0`, b2567). */
export function formSecurityBaru(operator: string): FormSecurity {
  return { id: '', reinsurerId: '', reinsurerName: '', pctShare: '', userId: operator }
}

/** `SetSecurityReinsurerLife_Act` b464: salin `REINSURERNAME`, `REINSURERID`, `PCTSHARE`, `ID`. */
export function formSecurityDari(s: SecurityReinsurer): FormSecurity {
  return { id: s.id, reinsurerId: s.reinsurerId, reinsurerName: s.reinsurerName, pctShare: s.pctShare, userId: s.userId }
}

/** Badan simpan; share dikirim apa adanya (di-trim) - server menormalkan koma dan memeriksa 0..100. */
export function keSecurityMasuk(f: FormSecurity): SecurityMasuk {
  return { id: f.id, reinsurerName: f.reinsurerName, reinsurerId: f.reinsurerId, pctShare: f.pctShare.trim() }
}

export default function PanelSecurity({ reinsurer, onTutup }: { reinsurer: Reinsurer; onTutup: () => void }) {
  const [jawab, setJawab] = useState<JawabanSecurity | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [form, setForm] = useState<FormSecurity | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [pesan, setPesan] = useState<string | null>(null)
  const master = useCariReinsurer(setGalatForm)

  const muat = useCallback(async () => {
    try {
      const j = await ambilSecurity(reinsurer.id)
      setJawab(j)
      setGalat(null)
    } catch (e) {
      setGalat(e)
    }
  }, [reinsurer.id])

  useEffect(() => {
    void muat()
  }, [muat])

  const hapusan = useHapus('security', async (id, p) => {
    if (form?.id === id) setForm(null)
    setPesan(p)
    await muat()
  })

  function buka(f: FormSecurity): void {
    setGalatForm(null)
    setPesan(null)
    master.reset()
    setForm(f)
  }

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatForm(null)
    try {
      await simpanSecurity(reinsurer.id, keSecurityMasuk(form))
      // `SaveSecurityReinsurerLife_Act` langkah 6 b1232: `DATASHOW = ""` - form tertutup.
      setForm(null)
      await muat()
    } catch (e) {
      setGalatForm(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const induk = jawab?.induk ?? reinsurer
  const daftar = jawab?.daftar ?? []

  return (
    <section className="panel">
      <KepalaPanel
        judul={SECURITY_MCRL.judul}
        // Kepala hanya Reinsurer Name - ID Reinsurer & PCT Share disembunyikan (work owner 04-10-2026).
        medan={[[SECURITY_MCRL.reinsurerName, induk.reinsurerName]]}
        onTutup={onTutup}
      />

      {form !== null && (
        <div className="panel">
          {galatForm !== null && <Gagal galat={galatForm} />}
          <div className="form-grid">
            <PilihSaring
              label={SECURITY_MCRL.formSecurityReinsurerName}
              value={form.reinsurerId}
              teksTerpilih={form.reinsurerName || form.reinsurerId}
              opsi={reinsurerTersedia(master.pilihan, daftar, form.id)}
              memuat={master.memuat}
              onCari={master.cari}
              onPilih={(o) => {
                setForm((f) => (f === null ? f : { ...f, reinsurerId: o.value, reinsurerName: o.label }))
              }}
              required
            />
            <Field
              label={SECURITY_MCRL.formShare}
              value={form.pctShare}
              onChange={(v) => {
                setForm((f) => (f === null ? f : { ...f, pctShare: v }))
              }}
              required
            />
            {/* Inputor di PALING AKHIR - seragam di semua form (04-10-2026). */}
            <Field label={SECURITY_MCRL.formInputor} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {SECURITY_MCRL.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatForm(null)
              }}
            >
              {SECURITY_MCRL.cancel}
            </button>
          </div>
        </div>
      )}

      {pesan !== null && <p role="status">{pesan}</p>}

      <div className="aksi-baris mcrl-aksi-grid">
        <button type="button" className="btn btn--primary" onClick={() => buka(formSecurityBaru(operatorKini()))}>
          {SECURITY_MCRL.add}
        </button>
      </div>
      {jawab === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {jawab !== null && daftar.length === 0 && <Kosong pesan={UMUM_MCRL.kosong} />}
      {daftar.length > 0 && (
        <div className="mcrl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{SECURITY_MCRL.kolomId}</th>
                <th>{SECURITY_MCRL.kolomReinsurerName}</th>
                <th>{SECURITY_MCRL.kolomShare}</th>
                <th>{SECURITY_MCRL.kolomEksposur}</th>
                <th>{SECURITY_MCRL.kolomInputor}</th>
                <th>{SECURITY_MCRL.kolomUpdateDate}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {daftar.map((s) => (
                <tr key={s.id} className="inbox__baris">
                  <td>{sel(s.id)}</td>
                  <td>{sel(s.reinsurerName)}</td>
                  <td>{selAngka(s.pctShare)}</td>
                  <td>{selAngka(jawab?.eksposur[s.id] ?? '')}</td>
                  <td>{sel(s.userId)}</td>
                  <td>{selWaktu(s.tglUpdate)}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formSecurityDari(s))}>
                      {SECURITY_MCRL.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setPesan(null)
                        hapusan.minta(s.id, s.reinsurerName || s.id)
                      }}
                    >
                      {SECURITY_MCRL.delete}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {/* Total share security - sama dengan Reinsurer List (keputusan work owner 04-10-2026). */}
      {jawab !== null && daftar.length > 0 && jawab.totalShare !== undefined && (
        <p className={kelasTotal(jawab.totalBukan100 === true)} role={jawab.totalBukan100 === true ? 'alert' : undefined}>
          <span>{REINSURER_MCRL.totalShare}</span> <strong>{selAngka(jawab.totalShare)}</strong>
          {jawab.totalBukan100 === true && <span className="mcrl-total__tanda">{REINSURER_MCRL.totalBukan100}</span>}
        </p>
      )}

      {hapusan.konfirmasi !== null && (
        <KonfirmasiHapus
          judul={SECURITY_MCRL.delete}
          labelBatal={SECURITY_MCRL.cancel}
          jenis="security"
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
