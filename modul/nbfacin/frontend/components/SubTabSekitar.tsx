// Sub-tab Surrounding Risk baris objek FIRE (tiket 38).
//
// Port `NB FacIn\Section\RiskAround.xml`: empat blok sisi Front / Left / Back / Right - Occupation (autocomplete
// data page `D_BrowseOccupationFacInFIRE` → RD `BrowseOccupationFacInFIRE_RD`, TYPE 'FIRE', nilai `.OldID`, memilih
// mengisi Note dengan `.Name`), Construction (dropdown), Distance (meter) (pxNumber min 0, 2 desimal), Note
// (baca-saja); lalu Other Description - Ownership, House keeping Status, Flood Area Status, Flood Area (tampil bila
// Flood Area Status = 0), Housekeeping Remark, dan tiga centang. Tidak ada medan wajib.
//
// Tata letak = tangkapan layar work owner 03-10-2026: dua blok "Other Description" berdampingan - kiri Ownership +
// House keeping Status sebaris, Flood Area Status, (Flood Area), Housekeeping Remark; kanan paragraf faktor risiko
// (teks dari tangkapan layar) + tiga centang.
//
// Keputusan agent (tiket 38): J-2 saran Occupation muncul setelah 2 karakter, jeda 400 ms; isian bebas tetap boleh
// (Pega mengizinkan). J-3 daftar dropdown = aturan properti `DDL\*.xml` (lihat labels.ts); Ownership / House keeping
// / Flood Area Status berawal "Not Informed" (pilihan pertama, tanpa baris kosong); Flood Area berpilihan kosong
// "Silahkan pilih" seperti Construction.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { cariOccupation, type BarisOccupation, type ObjekFire, type SisiRisiko, type SurroundingRisk } from '../api'
import {
  FLOOD_STATUS_TAMPIL,
  LAIN_SEKITAR as L,
  MEDAN_SISI,
  CONSTRUCTION_KOSONG,
  OPSI_CONSTRUCTION,
  OPSI_FLOOD_AREA,
  OPSI_FLOOD_STATUS,
  OPSI_HOUSEKEEPING,
  OPSI_OWNERSHIP,
  SISI_SEKITAR,
  TEKS_SEKITAR,
} from '../labels'

type KunciSisi = (typeof SISI_SEKITAR)[number]['kunci']

const MIN_CARI = 2
const JEDA_MS = 400

/** Sisi kosong. */
export const sisiKosong = (): SisiRisiko => ({ occupation: '', construction: '', distance: '', note: '' })

/** Surrounding Risk objek baru: status = pilihan pertama daftar ("Not Informed"). */
export function sekitarKosong(): SurroundingRisk {
  return {
    front: sisiKosong(), left: sisiKosong(), back: sisiKosong(), right: sisiKosong(),
    housekeepingStatus: OPSI_HOUSEKEEPING[0]!.value, floodAreaStatus: OPSI_FLOOD_STATUS[0]!.value, floodArea: '',
    housekeepingRemark: '',
  }
}

/** Ownership objek baru = pilihan pertama ("Not Informed"). */
export const OWNERSHIP_AWAL = OPSI_OWNERSHIP[0]!.value

/** Distance minus (`NegativeIsNotAllowed`). */
export function jarakMinus(distance: string): boolean {
  const n = Number(distance)
  return distance.trim() !== '' && Number.isFinite(n) && n < 0
}

/** Ada sisi dengan Distance minus. */
export function adaJarakMinus(s: SurroundingRisk): boolean {
  return SISI_SEKITAR.some((x) => jarakMinus(s[x.kunci].distance))
}

/** Memilih saran Occupation: nilai = OldID, Note = Name. */
export function pilihOccupation(sisi: SisiRisiko, b: BarisOccupation): SisiRisiko {
  return { ...sisi, occupation: b.oldId, note: b.name }
}

/** Isian Occupation dengan saran dari server. */
function IsianOccupation({ label, sisi, ubah }: { label: string; sisi: SisiRisiko; ubah: (s: SisiRisiko) => void }) {
  const [saran, setSaran] = useState<BarisOccupation[]>([])
  const [ketik, setKetik] = useState(false)
  const nomor = useRef(0)

  useEffect(() => {
    if (!ketik || sisi.occupation.trim().length < MIN_CARI) {
      setSaran([])
      return
    }
    const n = ++nomor.current
    const jadwal = window.setTimeout(() => {
      cariOccupation(sisi.occupation).then(
        (h) => {
          if (n === nomor.current) setSaran(h.baris)
        },
        () => {
          if (n === nomor.current) setSaran([])
        },
      )
    }, JEDA_MS)
    return () => window.clearTimeout(jadwal)
  }, [sisi.occupation, ketik])

  return (
    <div className="nbf-sekitar__occupation">
      <Field
        label={label}
        value={sisi.occupation}
        onChange={(v) => {
          setKetik(true)
          ubah({ ...sisi, occupation: v })
        }}
      />
      {saran.length > 0 && (
        <ul className="nbf-tambah-risk__saran" role="listbox" aria-label={label}>
          {saran.map((b) => (
            <li key={b.oldId}>
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                onClick={() => {
                  setKetik(false)
                  setSaran([])
                  ubah(pilihOccupation(sisi, b))
                }}
              >
                {b.oldId} · {b.name}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

/** Medan tampil-saja (Note). */
function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

function Centang({ label, checked, onChange }: { label: string; checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <label className="nbf-inward__pilihan">
      <input type="checkbox" checked={checked} onChange={(e) => onChange(e.target.checked)} /> {label}
    </label>
  )
}

export default function SubTabSekitar({ o, ubah }: { o: ObjekFire; ubah: (o: ObjekFire) => void }) {
  const s = o.surroundingRisk
  const setS = (baru: Partial<SurroundingRisk>) => ubah({ ...o, surroundingRisk: { ...s, ...baru } })
  const setSisi = (k: KunciSisi, sisi: SisiRisiko) => setS({ [k]: sisi } as Partial<SurroundingRisk>)

  return (
    <div className="nbf-objek__isi nbf-ringkas">
      {/* Empat sisi = kartu 2 x 2 (Front | Left, Back | Right); tata letak dirapikan work owner 05-10-2026. */}
      <div className="nbf-sekitar__grid">
        {SISI_SEKITAR.map((x) => {
          const sisi = s[x.kunci]
          return (
            <section key={x.kunci} className="nbf-sekitar__sisi">
              <h5 className="nbf-objek__judul">{x.judul}</h5>
              <IsianOccupation label={MEDAN_SISI.occupation} sisi={sisi} ubah={(v) => setSisi(x.kunci, v)} />
              <div className="nbf-sekitar__sebaris">
                <Pilih
                  label={MEDAN_SISI.construction}
                  value={sisi.construction}
                  onChange={(v) => setSisi(x.kunci, { ...sisi, construction: v })}
                  opsi={OPSI_CONSTRUCTION}
                  kosong={CONSTRUCTION_KOSONG}
                />
                <Field
                  label={MEDAN_SISI.distance}
                  type="number"
                  value={sisi.distance}
                  onChange={(v) => setSisi(x.kunci, { ...sisi, distance: v })}
                  error={jarakMinus(sisi.distance) ? TEKS_SEKITAR.jarakMinus : undefined}
                />
              </div>
              <Tampil label={MEDAN_SISI.note} nilai={sisi.note} />
            </section>
          )
        })}
      </div>

      <div className="nbf-sekitar__grid">
        <section className="nbf-sekitar__sisi">
          <h5 className="nbf-objek__judul">{L.judul.label}</h5>
          <div className="nbf-sekitar__sebaris">
            <Pilih label={L.ownership.label} value={o.ownership} onChange={(v) => ubah({ ...o, ownership: v })} opsi={OPSI_OWNERSHIP} />
            <Pilih label={L.housekeepingStatus.label} value={s.housekeepingStatus} onChange={(v) => setS({ housekeepingStatus: v })} opsi={OPSI_HOUSEKEEPING} />
          </div>
          <div className="nbf-sekitar__sebaris">
            <Pilih label={L.floodAreaStatus.label} value={s.floodAreaStatus} onChange={(v) => setS({ floodAreaStatus: v })} opsi={OPSI_FLOOD_STATUS} />
            {s.floodAreaStatus === FLOOD_STATUS_TAMPIL && (
              <Pilih label={L.floodArea.label} value={s.floodArea} onChange={(v) => setS({ floodArea: v })} opsi={OPSI_FLOOD_AREA} kosong={CONSTRUCTION_KOSONG} />
            )}
          </div>
          <Area label={L.housekeepingRemark.label} value={s.housekeepingRemark} onChange={(v) => setS({ housekeepingRemark: v })} baris={2} />
        </section>
        <section className="nbf-sekitar__sisi">
          <h5 className="nbf-objek__judul">{L.judul.label}</h5>
          <p className="nbf-sekitar__faktor">{TEKS_SEKITAR.faktorRisiko}</p>
          <div className="nbf-sekitar__centang">
            <Centang label={L.productionProcess.label} checked={o.isProductionProcess} onChange={(v) => ubah({ ...o, isProductionProcess: v })} />
            <Centang label={L.hotWork.label} checked={o.isHotWorkProcess} onChange={(v) => ubah({ ...o, isHotWorkProcess: v })} />
            <Centang label={L.flammable.label} checked={o.isFlammableItem} onChange={(v) => ubah({ ...o, isFlammableItem: v })} />
          </div>
        </section>
      </div>
    </div>
  )
}
