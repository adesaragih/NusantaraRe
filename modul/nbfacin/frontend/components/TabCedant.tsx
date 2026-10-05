// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Inw Fac
// Cedant Panels" (baris 60510-67588):
// - Source of Business (`.OfferFacIn.QuotationData.SobName`, tampil);
// - Share Cedant Type (`.OfferFacIn.ShareCedantType`, radio Gross "0" / Share RNM "1" = `DDL\ShareCedantType.xml`); wajib
//   bila proteksi berlaku; `SetShareOfCeding` (ShareOfCeding "100%" / PercentShare + "%") dijalankan backend saat Save;
// - % Share RNM, Total TSI RNM, Total Premi RNM (tampil - hitungan backend);
// - grid `.OfferFacIn.CedingCedantList`: Ceding · % Share · Add (`AddCedantList_act`) / Delete; tombol "Choose Ceding"
//   (`Section\CedingCedant.xml` -> harness `CedingCedant`, jendela "Choose Ceding" -> `SetDataSobCeding_Act`: baris =
//   `.ID` / `.ClientName`) memakai popup Ceding Company tab General (gambar layar work owner 05-10-2026);
// - Save -> `PUT …/cedant`; proteksi `ProtectShareCedant_Act` (pesan verbatim) diperiksa juga di layar.
//
// Keputusan agent (tiket 49): C-1 popup baris `CedingCedant` tidak dibuat - dua isiannya (Choose Ceding, % Share) ada di
// baris grid; grid CurrencyList (TSI / Premi per cedant) dan proteksi "Total premi share cedant must be equal total premi
// RNM!" menunggu tab Payment (`.Policy.Payment.Premium`; disetujui work owner 05-10-2026). C-2 % Share RNM baca-saja di
// sini - diubah di tab Spreading (satu jalur hitung). C-3 proteksi hanya bila backend menyatakan `wajib` (gerbang
// `ProtectShareCedant_Act`: IsGroup == "Group" atau FlagOnGoingPolicy == 2).

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import { ambilCedant, simpanCedant, type BarisCedant, type TampilanCedant } from '../api'
import { CEDANT as C, OPSI_SHARE_CEDANT_TYPE, TEKS_CEDANT as T } from '../labels'
import PopupPilihAgent from './PopupPilihAgent'
import { lebihDari100 } from './TabSpreading'

const DESIMAL = 4

const barisKosong = (): BarisCedant => ({ cedingCo: '', cedingCoName: '', shareCeding: '' })

/**
 * Add (`AddCedantList_act`): grid kosong -> salin Ceding Co tab General (`QuotationData.CedingCoList`, atau CedingCo
 * tunggal - dilipat backend ke `cedingUmum`); tanpa keduanya -> satu baris kosong; grid berisi -> tambah baris kosong.
 */
export function tambahCedant(baris: BarisCedant[], cedingUmum: BarisCedant[]): BarisCedant[] {
  if (baris.length > 0) return [...baris, barisKosong()]
  if (cedingUmum.length === 0) return [barisKosong()]
  return cedingUmum.map((c) => ({ cedingCo: c.cedingCo, cedingCoName: c.cedingCoName, shareCeding: '' }))
}

/** % Share sah: desimal > 0 dan <= 100 (`.ShareCeding<=0||.ShareCeding>100` ditolak); kosong = tidak sah. */
export function shareSah(v: string): boolean {
  const t = v.trim()
  if (!desimalSah(t) || t.startsWith('-')) return false
  return /[1-9]/.test(t) && !lebihDari100(t)
}

/** Ada % Share > 100 (`SetTSIPremiCedant_Act` - ditolak backend selalu, tanpa gerbang `wajib`). */
export function adaShareLebih(baris: BarisCedant[]): boolean {
  return baris.some((b) => desimalSah(b.shareCeding.trim()) && lebihDari100(b.shareCeding.trim()))
}

/** Pesan proteksi `ProtectShareCedant_Act` (urutan langkah Pega, tanpa pesan kembar). */
export function galatCedant(baris: BarisCedant[]): string[] {
  if (baris.length === 0) return [T.listKosong]
  const pesan: string[] = []
  if (baris.some((b) => !shareSah(b.shareCeding))) pesan.push(T.shareTidakSah)
  if (baris.some((b) => b.cedingCoName.trim() === '')) pesan.push(T.cedingKosong)
  return pesan
}

/** Medan tampil-saja (label kiri, teks kanan). */
function Tampil({ label, nilai, angka }: { label: string; nilai: string; angka?: boolean }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className={'nbf-inward__teks' + (angka ? ' nbf-angka' : '')}>{nilai}</div>
    </div>
  )
}

export default function TabCedant({ caseId }: { caseId: string }) {
  const [v, setV] = useState<TampilanCedant | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [cobaSimpan, setCobaSimpan] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  // Baris yang sedang memilih ceding (popup Choose Ceding).
  const [pilihUntuk, setPilihUntuk] = useState<number | null>(null)

  useEffect(() => {
    let batal = false
    ambilCedant(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const ubah = (baru: Partial<TampilanCedant>) => {
    setV({ ...v, ...baru })
    setTersimpan(false)
  }
  const ubahBaris = (n: number, b: Partial<BarisCedant>) => ubah({ cedant: v.cedant.map((x, k) => (k === n ? { ...x, ...b } : x)) })
  const typeKosong = v.wajib && v.shareCedantType === ''
  const pesan = v.wajib ? galatCedant(v.cedant) : []
  const shareLebih = adaShareLebih(v.cedant)

  async function simpan() {
    if (v === null) return
    setCobaSimpan(true)
    if (typeKosong || pesan.length > 0 || shareLebih) return
    setMenyimpan(true)
    setGalat(null)
    try {
      setV(await simpanCedant(caseId, { shareCedantType: v.shareCedantType, cedant: v.cedant }))
      setTersimpan(true)
      setCobaSimpan(false)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  return (
    <div className="nbf-objek nbf-cedant">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}

      <section className="nbf-lapis nbf-cedant__atas nbf-labelkiri">
        <Tampil label={C.sob} nilai={v.sobName} />
        <div className="field">
          <span className="field__label">
            {C.shareCedantType}
            {v.wajib && <span className="field__req">*</span>}
          </span>
          <div className="nbf-inward__radio" role="radiogroup" aria-label={C.shareCedantType}>
            {OPSI_SHARE_CEDANT_TYPE.map((o) => (
              <label key={o.value} className="nbf-inward__pilihan">
                <input
                  type="radio"
                  name="shareCedantType"
                  value={o.value}
                  checked={v.shareCedantType === o.value}
                  onChange={() => ubah({ shareCedantType: o.value })}
                />{' '}
                {o.label}
              </label>
            ))}
          </div>
          {cobaSimpan && typeKosong && <div className="field__error">{T.typeWajib}</div>}
        </div>
        <Tampil label={C.percentShare} nilai={formatNumber(v.percentShare, DESIMAL)} angka />
        <Tampil label={C.totalTsi} nilai={formatNumber(v.totalTsiRnm, DESIMAL)} angka />
        <Tampil label={C.totalPremi} nilai={formatNumber(v.totalPremiRnm, DESIMAL)} angka />
      </section>

      <section className="nbf-lapis nbf-cedant__grid">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                <th scope="col">{C.kolomCeding}</th>
                <th scope="col">{C.kolomShare}</th>
                <th scope="col" className="table__actions">
                  <button type="button" className="btn btn--sm" onClick={() => ubah({ cedant: tambahCedant(v.cedant, v.cedingUmum) })}>
                    {C.tambah}
                  </button>
                </th>
              </tr>
            </thead>
            <tbody>
              {v.cedant.length === 0 && (
                <tr>
                  <td colSpan={3}>{T.kosong}</td>
                </tr>
              )}
              {v.cedant.map((b, n) => (
                <tr key={n}>
                  <td>
                    <div className="nbf-cedant__ceding">
                      <span>{b.cedingCoName}</span>
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => setPilihUntuk(n)}>
                        {C.pilihCeding}
                      </button>
                    </div>
                  </td>
                  <td>
                    <input
                      className="nbf-angka"
                      aria-label={C.kolomShare}
                      inputMode="decimal"
                      value={b.shareCeding}
                      onChange={(e) => ubahBaris(n, { shareCeding: e.target.value })}
                    />
                  </td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--sm" onClick={() => ubah({ cedant: v.cedant.filter((_, k) => k !== n) })}>
                      {C.hapus}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {shareLebih && <div className="field__error">{T.nilaiLebih}</div>}
        {cobaSimpan &&
          pesan.map((p) => (
            <div key={p} className="field__error">
              {p}
            </div>
          ))}
      </section>

      {pilihUntuk !== null && (
        <PopupPilihAgent
          judul={C.judulPilih}
          onTutup={() => setPilihUntuk(null)}
          onPilih={(a) => {
            ubahBaris(pilihUntuk, { cedingCo: a.id, cedingCoName: a.name })
            setPilihUntuk(null)
          }}
        />
      )}

      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? T.menyimpan : T.simpan}
        </button>
      </div>
    </div>
  )
}
