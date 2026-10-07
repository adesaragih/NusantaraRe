// Panel business - tiket 07/08. Harness `InboxBusinessLifeReinsurers` → section
// `InputBusinessLifeReinsurers.xml`, dibuka tombol baris `Business List` panel kontrak
// (`InputRetroLimitReinsurers.xml` b12999).
//
// Kepala b1066/b1282/b1497 (⚠️ `ID Reins Type` berisi ID KONTRAK). Wadah `Add` + grid b8260
// `HASILD10 = 2` ber-`ALWAYS` → selalu tampil (RALAT R6). Grid RD `BrowseTreatyBusiness_Life_RD` (urut
// `.TGLUPDATE ASC` - server). Form (wadah b2022): `BUSINESS CODE` b3675 (ro), `BUSINESS NAME`
// autocomplete b3862, `R/I RATE` autocomplete b4314, `View Rate` b4732 (tampil bila `RIRATEID != ''`),
// `Save` b5739, `Cancel` b5954. Baris: `Edit` b10824, `Delete` b11154, `View Rate` b11469,
// `Copy to all Reinstype` b12020. `Add` b8484.
//
// K1 keputusan work owner 01-10-2026 (OQ-MCRL-13): `R/I RATE` (tabel `M_RATE_LIFE_SUMMARY`) dan `Rate List`
// (view `RATE_LIFE`) dibaca saja; business BARU dapat disimpan (`RIRATEID` wajib, `SaveBusinessLife_Act`
// b589, diisi autocomplete). Server menolak RIRATEID pilihan baru yang tidak ada di view ringkasan.

import { useCallback, useEffect, useState } from 'react'

import { ambilBusiness, simpanBusiness, type Business, type BusinessMasuk, type JawabanBusiness, type Kontrak } from '../api'
import { BUSINESS_MCRL, UMUM_MCRL } from '../labels'
import { operatorKini, sel, selWaktu, waktuKini } from '../tampilan'
import { Field, Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { KepalaPanel } from './Bingkai'
import { useCariBusiness, useCariRate } from './cariMaster'
import KonfirmasiHapus from './KonfirmasiHapus'
import ModalRate from './ModalRate'
import PratinjauSalin from './PratinjauSalin'
import { useHapus } from './useHapus'

/** Isian form - `InputBusinessLife.*`. */
export interface FormBusiness {
  id: string
  bizCode: string
  bizName: string
  riRateId: string
  riRate: string
  /** `Inputor` / `Modified Date` (ro). */
  userId: string
  tglUpdate: string
}

/**
 * `NewInputBusinessLife_Act`: langkah 1 b268 `Page-New` ber-`//` (TIDAK jalan); langkah 2 b412 mengosongkan
 * `ID, BIZCODE, BIZNAME, RIRATEID, RIRATE`, `USERID ← OperatorID`, `TGLUPDATE ← @CurrentDateTime()`.
 */
export function formBusinessBaru(operator: string, kini: string): FormBusiness {
  return { id: '', bizCode: '', bizName: '', riRateId: '', riRate: '', userId: operator, tglUpdate: kini }
}

/** `SetBusinessListLife_Act` b501: `ID ← IDPEGA`, `BIZCODE`, `BIZNAME`, `USERID ← OperatorID`, `RIRATEID`, `RIRATE`. */
export function formBusinessDari(b: Business, operator: string): FormBusiness {
  return {
    id: b.id, bizCode: b.bizCode, bizName: b.bizName, riRateId: b.riRateId, riRate: b.riRate,
    userId: operator, tglUpdate: b.tglUpdate,
  }
}

export function keBusinessMasuk(f: FormBusiness): BusinessMasuk {
  return { id: f.id, bizCode: f.bizCode, bizName: f.bizName, riRateId: f.riRateId, riRate: f.riRate }
}

/** Tombol `View Rate` form b4732 tampil bila `InputBusinessLife.RIRATEID != ''`. */
/**
 * Pilihan Business Name yang BELUM dipakai business lain di kontrak ini - satu
 * Business Name sekali per kontrak (keputusan work owner 04-10-2026). Milik
 * business yang sedang diubah (`idSendiri`) tetap tersedia. Server menolak hal
 * yang sama (`PesanBusinessGanda`).
 */
export function bisnisTersedia(
  pilihan: readonly OpsiSaring[],
  daftar: readonly Business[],
  idSendiri: string,
): OpsiSaring[] {
  const terpakai = new Set(daftar.filter((b) => b.id !== idSendiri).map((b) => b.bizCode.trim()))
  return pilihan.filter((o) => !terpakai.has(o.value.trim()))
}

export function viewRateTampil(f: FormBusiness): boolean {
  return f.riRateId.trim() !== ''
}

export default function PanelBusiness({ kontrak, onTutup }: { kontrak: Kontrak; onTutup: () => void }) {
  const [jawab, setJawab] = useState<JawabanBusiness | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [form, setForm] = useState<FormBusiness | null>(null)
  const [galatForm, setGalatForm] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [pesan, setPesan] = useState<string | null>(null)
  const [rate, setRate] = useState<string | null>(null)
  const [salin, setSalin] = useState<Business | null>(null)
  const masterBisnis = useCariBusiness(setGalatForm)
  const masterRate = useCariRate(setGalatForm)

  const muat = useCallback(async () => {
    try {
      const j = await ambilBusiness(kontrak.id)
      setJawab(j)
      setGalat(null)
    } catch (e) {
      setGalat(e)
    }
  }, [kontrak.id])

  useEffect(() => {
    void muat()
  }, [muat])

  const hapusan = useHapus('business', async (id, p) => {
    if (form?.id === id) setForm(null)
    setPesan(p)
    await muat()
  })

  function buka(f: FormBusiness): void {
    setGalatForm(null)
    setPesan(null)
    masterBisnis.reset()
    masterRate.reset()
    setForm(f)
  }

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatForm(null)
    try {
      await simpanBusiness(kontrak.id, keBusinessMasuk(form))
      // `SaveBusinessLife_Act` langkah 6 b1187: `DATASHOW = ""` - form tertutup.
      setForm(null)
      await muat()
    } catch (e) {
      setGalatForm(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const induk = jawab?.kontrak ?? kontrak
  const daftar = jawab?.daftar ?? []

  return (
    <section className="panel">
      <KepalaPanel
        judul={BUSINESS_MCRL.judul}
        medan={[
          [BUSINESS_MCRL.idTreatyYear, induk.idTreatyYear],
          // `ID Reins Type` (ID kontrak) SENGAJA tidak ditampilkan - keputusan
          // work owner 04-10-2026.
          [BUSINESS_MCRL.reinsType, induk.reinsTypeName],
        ]}
        onTutup={onTutup}
      />

      {form !== null && (
        <div className="panel">
          {galatForm !== null && <Gagal galat={galatForm} />}
          <div className="form-grid">
            <Field label={BUSINESS_MCRL.formBusinessCode} value={form.bizCode} onChange={() => undefined} readOnly />
            <PilihSaring
              label={BUSINESS_MCRL.formBusinessName}
              value={form.bizCode}
              teksTerpilih={form.bizName || form.bizCode}
              opsi={bisnisTersedia(masterBisnis.pilihan, daftar, form.id)}
              memuat={masterBisnis.memuat}
              onCari={masterBisnis.cari}
              onPilih={(o) => {
                setForm((f) => (f === null ? f : { ...f, bizCode: o.value, bizName: o.label }))
              }}
              required
            />
            <PilihSaring
              label={BUSINESS_MCRL.formRiRate}
              value={form.riRateId}
              teksTerpilih={form.riRate || form.riRateId}
              opsi={masterRate.pilihan}
              memuat={masterRate.memuat}
              onCari={masterRate.cari}
              onPilih={(o) => {
                setForm((f) => (f === null ? f : { ...f, riRateId: o.value, riRate: o.label }))
              }}
              required
            />
            {/* Modified Date dan Inputor di PALING AKHIR - seragam di semua form (04-10-2026). */}
            <Field label={BUSINESS_MCRL.formModifiedDate} value={selWaktu(form.tglUpdate)} onChange={() => undefined} readOnly />
            <Field label={BUSINESS_MCRL.formInputor} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {BUSINESS_MCRL.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatForm(null)
              }}
            >
              {BUSINESS_MCRL.cancel}
            </button>
            {/* `View Rate` di POJOK KANAN baris tombol form (keputusan work owner 04-10-2026). */}
            {viewRateTampil(form) && (
              <button
                type="button"
                className="btn btn--ghost mcrl-aksi-kanan"
                onClick={() => {
                  setRate(form.riRateId)
                }}
              >
                {BUSINESS_MCRL.viewRateForm}
              </button>
            )}
          </div>
        </div>
      )}

      {pesan !== null && <p role="status">{pesan}</p>}

      <div className="aksi-baris mcrl-aksi-grid">
        <button type="button" className="btn btn--primary" onClick={() => buka(formBusinessBaru(operatorKini(), waktuKini(new Date())))}>
          {BUSINESS_MCRL.add}
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
                <th>{BUSINESS_MCRL.kolomBusinessName}</th>
                <th>{BUSINESS_MCRL.kolomRiRate}</th>
                <th>{BUSINESS_MCRL.kolomInputor}</th>
                <th>{BUSINESS_MCRL.kolomUpdateDate}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {daftar.map((b) => (
                <tr key={b.id} className="inbox__baris">
                  <td>{sel(b.bizName)}</td>
                  <td>{sel(b.riRate)}</td>
                  <td>{sel(b.userId)}</td>
                  <td>{selWaktu(b.tglUpdate)}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formBusinessDari(b, operatorKini()))}>
                      {BUSINESS_MCRL.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setPesan(null)
                        hapusan.minta(b.id, b.bizName || b.id)
                      }}
                    >
                      {BUSINESS_MCRL.delete}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setRate(b.riRateId)
                      }}
                    >
                      {BUSINESS_MCRL.viewRate}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setPesan(null)
                        setSalin(b)
                      }}
                    >
                      {BUSINESS_MCRL.copyToAll}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {rate !== null && (
        <ModalRate
          idusedby={rate}
          onTutup={() => {
            setRate(null)
          }}
        />
      )}
      {salin !== null && (
        <PratinjauSalin
          key={salin.id}
          business={salin}
          onSelesai={(p) => {
            setSalin(null)
            setPesan(p)
          }}
          onBatal={() => {
            setSalin(null)
          }}
        />
      )}
      {hapusan.konfirmasi !== null && (
        <KonfirmasiHapus
          judul={BUSINESS_MCRL.delete}
          labelBatal={BUSINESS_MCRL.cancel}
          jenis="business"
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
