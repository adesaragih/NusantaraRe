// Tab Object layar Inward Facultative, Class of Business MARINE CARGO (tiket 74 P1). Kabel `GET`/`PUT …/objek-cargo` (sesi c3).
// Bukti `NB\` = `D:\migrasi\RNM\NB FacIn\`; spesifikasi `docs/issues/74-marine-cargo.md` S1-S2.
//
// 1. Grid `.OfferFacIn.CargoList` (`NB\Section\InputDtlObject_FacIn.xml` kontainer L1994 `IsMarineCargo`, grid L2849): kolom
//    No. / Top Risk / SailDate / Goods / Ship Name / Conveyance; "+" di kepala (L4086 addRow AFTER), Delete per baris (NB
//    deleteRow L5527; jalur EDM `SetFlagDeleteEDM_Act` bukan di sini). Baris dibuka = expand pane `InputDtlCargo_FacIn` (L6260).
// 2. Form cargo `InputDtlCargo_FacIn.xml`: Top Risk + Trading / Goods (kiri) + Conveyance / Packing (kanan), dropdown dari master
//    (nilai .ID, tampil .Description); Note diisi seperti `ProtectObjectMarine` (Note = Description; Conveyance juga Name).
// 3. `MarineCargoDtl.xml` dua kolom: kiri Survey Agent (L1184) / LC (L4517) / Ship Data (L7257); kanan Sail Date … Max Age
//    Detail (L12783) + Additional Ship (L15307).
//
// Penyimpangan sadar (dicatat di tiket 74):
// - Popup Select Ship / Select Agent memakai satu kotak cari (Pega: filter per kolom grid).
// - Tahap UW: Pega membuat Top Risk tetap dapat diubah (`InputDtlCargo_FacIn_IsUW.xml` L500); di sistem baru seluruh case
//   baca-saja di UW kecuali Spreading (tiket 59), jadi ikut baca-saja.

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Halaman, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilObjekCargo,
  cariAgentSurvei,
  cariKapal,
  cariMaxAge,
  daftarMasterCargo,
  simpanObjekCargo,
  type BarisAgentSurvei,
  type BarisKapal,
  type BarisMasterCargo,
  type JenisMasterCargo,
  type KapalTambahan,
  type ObjekCargo,
} from '../api'
import {
  FORM_CARGO as F,
  GRID_CARGO,
  GRID_OBJEK,
  POPUP_AGENT_CARGO,
  POPUP_KAPAL_CARGO,
  SIMPAN_OBJEK,
  TEKS_ANEKA,
  TEKS_CARGO,
  TEKS_FORM_OPPORTUNITY,
} from '../labels'
import { pesanProteksiDari } from './BlokEmailType'
import { cargoBaru, kapalTambahanBaru, rapikanCargo, tampilTanggal, terapkanAgent, terapkanKapal, terapkanKapalTambahan } from './cargo'
import IsianSaran from './IsianSaran'
import TanggalDMY from './TanggalDMY'

const JEDA_MS = 400

type Master = Record<JenisMasterCargo, BarisMasterCargo[]>
const MASTER_KOSONG: Master = { trading: [], goods: [], conveyance: [], packing: [] }

function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

/** Dropdown master (`pxDropdown`, teks tanpa pilihan "Choose", wajib). */
function PilihMaster({
  id,
  label,
  nilai,
  opsi,
  onPilih,
}: {
  id: string
  label: string
  nilai: string
  opsi: BarisMasterCargo[]
  onPilih: (b: BarisMasterCargo | null) => void
}) {
  return (
    <div className="field">
      <label className="field__label" htmlFor={id}>
        {label}
        <span className="field__req">*</span>
      </label>
      <select
        id={id}
        className="field__input"
        value={nilai}
        onChange={(e) => onPilih(opsi.find((o) => o.id === e.target.value) ?? null)}
      >
        <option value="">{F.pilih}</option>
        {/* Nilai tersimpan yang tidak ada di master tetap tampil (tidak diam-diam terganti). */}
        {nilai !== '' && !opsi.some((o) => o.id === nilai) && <option value={nilai}>{nilai}</option>}
        {opsi.map((o) => (
          <option key={o.id} value={o.id}>
            {o.description}
          </option>
        ))}
      </select>
    </div>
  )
}

/** Kotak cari + tabel pilih bersama kedua popup (Select Ship / Select Agent). */
function useCari<T>(cari: (q: string) => Promise<{ baris: T[] }>) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<T[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  useEffect(() => {
    let batal = false
    const t = setTimeout(() => {
      cari(kotak.trim()).then(
        (h) => {
          if (batal) return
          setHasil(h.baris)
          setHalaman(1)
        },
        (err: unknown) => {
          if (!batal) setGalat(err)
        },
      )
    }, JEDA_MS)
    return () => {
      batal = true
      clearTimeout(t)
    }
  }, [kotak, cari])
  return { kotak, setKotak, hasil, galat, halaman, setHalaman }
}

/** Popup Select Ship - `NB\Section\SelectShip.xml` (grid `BrowseShip_RD`, tombol Select → `SelectShip_act`). */
export function PopupKapalCargo({ onTutup, onPilih }: { onTutup: () => void; onPilih: (k: BarisKapal) => void }) {
  const c = useCari(cariKapal)
  const u = GRID_OBJEK.ukuran
  const baris = c.hasil ?? []
  const tampil = baris.slice((c.halaman - 1) * u, c.halaman * u)
  return (
    <Modal judul={POPUP_KAPAL_CARGO.judul} onTutup={onTutup} penuh>
      <input className="field__input" type="search" aria-label={TEKS_CARGO.cari} placeholder={TEKS_CARGO.cari} value={c.kotak} onChange={(e) => c.setKotak(e.target.value)} />
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid">
          <thead>
            <tr>
              {POPUP_KAPAL_CARGO.kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" />
            </tr>
          </thead>
          <tbody>
            {tampil.map((k) => (
              <tr key={k.id}>
                {[k.register ?? '', k.imo ?? '', k.id, k.nama, k.formerShipName ?? '', k.shipTypeName ?? '', k.flag, k.konstruksi, k.klasifikasi,
                  k.grt, k.nrt, k.dwt, k.tahun, k.yMake2 ?? '', k.rebuilt ?? '', k.remark ?? ''].map((v, i) => (
                  <td key={i}>{v}</td>
                ))}
                <td>
                  <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => onPilih(k)}>
                    {POPUP_KAPAL_CARGO.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {c.hasil === null && c.galat === null && <Memuat />}
      {c.hasil !== null && baris.length === 0 && <p className="muted">{TEKS_CARGO.kosongCari}</p>}
      {baris.length > u && <Halaman halaman={c.halaman} ukuran={u} total={baris.length} onPindah={c.setHalaman} />}
      <Gagal galat={c.galat} />
    </Modal>
  )
}

/** Popup Select Agent - `NB\Section\SelectAgent.xml` (`BrowseLLOYDAGENT_RD`: ID / Name / Address / City / Country). */
export function PopupAgentCargo({ onTutup, onPilih }: { onTutup: () => void; onPilih: (a: BarisAgentSurvei) => void }) {
  const c = useCari(cariAgentSurvei)
  const u = GRID_OBJEK.ukuran
  const baris = c.hasil ?? []
  const tampil = baris.slice((c.halaman - 1) * u, c.halaman * u)
  return (
    <Modal judul={POPUP_AGENT_CARGO.judul} onTutup={onTutup} penuh>
      <input className="field__input" type="search" aria-label={TEKS_CARGO.cari} placeholder={TEKS_CARGO.cari} value={c.kotak} onChange={(e) => c.setKotak(e.target.value)} />
      <table className="nbf-tabel nbf-aneka-grid">
        <thead>
          <tr>
            {POPUP_AGENT_CARGO.kolom.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
            <th scope="col" />
          </tr>
        </thead>
        <tbody>
          {tampil.map((a) => (
            <tr key={a.id}>
              <td>{a.id}</td>
              <td>{a.name}</td>
              <td>{a.address}</td>
              <td>{a.city}</td>
              <td>{a.country}</td>
              <td>
                <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => onPilih(a)}>
                  {POPUP_AGENT_CARGO.pilih}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {c.hasil === null && c.galat === null && <Memuat />}
      {c.hasil !== null && baris.length === 0 && <p className="muted">{TEKS_CARGO.kosongCari}</p>}
      {baris.length > u && <Halaman halaman={c.halaman} ukuran={u} total={baris.length} onPindah={c.setHalaman} />}
      <Gagal galat={c.galat} />
    </Modal>
  )
}

/** Grid Additional Ship (`MarineCargoDtl.xml` L15875): Register ID (autocomplete kapal) / Ship Name / Voyage No + Add / Delete. */
function GridKapalTambahan({ daftar, ubah }: { daftar: KapalTambahan[]; ubah: (d: KapalTambahan[]) => void }) {
  const ganti = (i: number, t: KapalTambahan) => ubah(daftar.map((x, j) => (j === i ? t : x)))
  // Hasil cari terakhir per id: `GET /kapal?cari=` mencari NAMA, jadi kapal terpilih diambil dari hasil yang sudah tampil.
  const terakhir = useRef(new Map<string, BarisKapal>())
  const cari = (q: string) =>
    cariKapal(q).then((h) => {
      for (const k of h.baris) terakhir.current.set(k.id, k)
      return { baris: h.baris.map((k) => ({ id: k.id, label: `${k.id} - ${k.nama}` })) }
    })
  return (
    <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi nbf-cargo-kapal">
      <thead>
        <tr>
          {F.kolomKapalTambahan.map((k) => (
            <th key={k} scope="col">
              {k}
            </th>
          ))}
          <th scope="col">
            <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => ubah([...daftar, kapalTambahanBaru()])}>
              {F.tambah}
            </button>
          </th>
        </tr>
      </thead>
      <tbody>
        {daftar.map((t, i) => (
          <tr key={i}>
            <td>
              <IsianSaran
                label={F.kolomKapalTambahan[0]}
                value={t.id}
                onKetik={(v) => ganti(i, { ...t, id: v })}
                cari={cari}
                onPilih={(s) => {
                  const k = terakhir.current.get(s.id)
                  if (k) ganti(i, terapkanKapalTambahan(t, k))
                }}
              />
            </td>
            <td>
              <input className="field__input" aria-label={F.kolomKapalTambahan[1]} value={t.nmShip} onChange={(e) => ganti(i, { ...t, nmShip: e.target.value })} />
            </td>
            <td>
              <input className="field__input" aria-label={F.kolomKapalTambahan[2]} value={t.voyageNo} onChange={(e) => ganti(i, { ...t, voyageNo: e.target.value })} />
            </td>
            <td>
              <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => ubah(daftar.filter((_, j) => j !== i))}>
                {F.hapus}
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

/** Form satu cargo - `InputDtlCargo_FacIn.xml` + `MarineCargoDtl.xml`. */
export function FormCargo({ c, ubah, master, idx }: { c: ObjekCargo; ubah: (c: ObjekCargo) => void; master: Master; idx: number }) {
  const [popup, setPopup] = useState<'kapal' | 'agent' | null>(null)
  const id = (m: string) => `nbfacin-cargo-${idx}-${m}`
  const k = c.ship
  return (
    <div className="nbf-ringkas nbf-aneka-form">
      <section className="nbf-alamat__blok nbf-aneka-form__kartu">
        <label className="nbf-inward__pilihan">
          <input type="checkbox" checked={c.isTopRisk} onChange={(e) => ubah({ ...c, isTopRisk: e.target.checked })} /> {F.topRisk}
        </label>
        <div className="nbf-aneka-form__grid">
          <PilihMaster id={id('trading')} label={F.trading} nilai={c.tradingId} opsi={master.trading}
            onPilih={(b) => ubah({ ...c, tradingId: b?.id ?? '', tradingNote: b?.description ?? '' })} />
          <PilihMaster id={id('conveyance')} label={F.conveyance} nilai={c.conveyanceId} opsi={master.conveyance}
            onPilih={(b) => ubah({ ...c, conveyanceId: b?.id ?? '', conveyanceName: b?.description ?? '', conveyanceNote: b?.description ?? '' })} />
          <PilihMaster id={id('goods')} label={F.goods} nilai={c.goodId} opsi={master.goods}
            onPilih={(b) => ubah({ ...c, goodId: b?.id ?? '', goodNote: b?.description ?? '' })} />
          <PilihMaster id={id('packing')} label={F.packing} nilai={c.packingId} opsi={master.packing}
            onPilih={(b) => ubah({ ...c, packingId: b?.id ?? '', packingNote: b?.description ?? '' })} />
        </div>
      </section>

      <div className="nbf-aneka-form__baris">
        <div className="nbf-aneka-form">
          <section className="nbf-alamat__blok nbf-aneka-form__kartu">
            <div className="nbf-alamat__kepala">
              <div className="nbf-opp__tombol">
                <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPopup('agent')}>
                  {F.selectAgent}
                </button>
              </div>
              <h5 className="nbf-objek__judul">{F.grupAgent}</h5>
            </div>
            <div className="nbf-aneka-form__grid">
              <Tampil label={F.fullName} nilai={c.surveyAgent.fullName} />
              <Tampil label={F.address} nilai={c.surveyAgent.address} />
              <Tampil label={F.city} nilai={c.surveyAgent.city} />
              <Tampil label={F.country} nilai={c.surveyAgent.country} />
            </div>
          </section>

          <section className="nbf-alamat__blok nbf-aneka-form__kartu">
            <h5 className="nbf-objek__judul">{F.grupLc}</h5>
            <label className="nbf-inward__pilihan">
              <input type="checkbox" checked={c.lc.lc} onChange={(e) => ubah({ ...c, lc: { ...c.lc, lc: e.target.checked } })} /> {F.lcStatus}
            </label>
            <div className="nbf-aneka-form__grid">
              <Field label={F.lcBank} value={c.lc.bank} onChange={(v) => ubah({ ...c, lc: { ...c.lc, bank: v } })} />
              <Field label={F.lcNumber} value={c.lc.number} onChange={(v) => ubah({ ...c, lc: { ...c.lc, number: v } })} />
              <Field label={F.lcCondition} value={c.lc.condition} onChange={(v) => ubah({ ...c, lc: { ...c.lc, condition: v } })} />
              <Field label={F.lcRemark} value={c.lc.remark} onChange={(v) => ubah({ ...c, lc: { ...c.lc, remark: v } })} />
              <TanggalDMY
                label={F.lcEndDate}
                value={c.lc.endDate}
                onChange={(v) => ubah({ ...c, lc: { ...c.lc, endDate: v } })}
                labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
                pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
              />
            </div>
          </section>

          <section className="nbf-alamat__blok nbf-aneka-form__kartu nbf-aneka-form__kartu--kapal">
            <div className="nbf-alamat__kepala">
              <div className="nbf-opp__tombol">
                <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setPopup('kapal')}>
                  {F.selectShip}
                </button>
              </div>
              <h5 className="nbf-objek__judul">{F.grupKapal}</h5>
            </div>
            <div className="nbf-aneka-form__grid">
              <Tampil label={F.shipName} nilai={k.nmShip} />
              <Tampil label={F.shipType} nilai={k.shipTypeName} />
              <Tampil label={F.grt} nilai={k.grt} />
              <Tampil label={F.nrt} nilai={k.nrt} />
              <Tampil label={F.dwt} nilai={k.dwt} />
              <Tampil label={F.age} nilai={k.age} />
              <Tampil label={F.flagName} nilai={k.flagName} />
              <Tampil label={F.construction} nilai={k.construction} />
              <Tampil label={F.yMake1} nilai={k.yMake1} />
              <Tampil label={F.remark} nilai={k.remark} />
            </div>
          </section>
        </div>

        <div className="nbf-aneka-form">
          <section className="nbf-alamat__blok nbf-aneka-form__kartu">
            <div className="nbf-aneka-form__grid">
              <TanggalDMY
                label={F.sailDate}
                value={c.sailDate}
                onChange={(v) => ubah({ ...c, sailDate: v })}
                labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
                pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
              />
              <Field label={F.fromRute} value={c.fromRute} onChange={(v) => ubah({ ...c, fromRute: v })} required />
              <Field label={F.toRute} value={c.toRute} onChange={(v) => ubah({ ...c, toRute: v })} required />
              <div className="field">
                <label className="field__label" htmlFor={id('bl')}>
                  {F.blNumber}
                </label>
                <textarea id={id('bl')} className="field__input" rows={2} value={c.blNumber} onChange={(e) => ubah({ ...c, blNumber: e.target.value })} />
              </div>
              <Field label={F.invoiceNumber} value={c.invoiceNumber} onChange={(v) => ubah({ ...c, invoiceNumber: v })} />
              <TanggalDMY
                label={F.invoiceDate}
                value={c.invoiceDate}
                onChange={(v) => ubah({ ...c, invoiceDate: v })}
                labelKalender={TEKS_FORM_OPPORTUNITY.kalender}
                pesanFormat={TEKS_FORM_OPPORTUNITY.formatTanggal}
              />
              <IsianSaran
                label={F.maxAge}
                value={c.maxAge}
                onKetik={(v) => ubah({ ...c, maxAge: v, maxAgeDetail: '' })}
                cari={(q) => cariMaxAge(q).then((h) => ({ baris: h.baris.map((b) => ({ id: b.id, label: b.id, ekstra: b.warranty })) }))}
                onPilih={(s) => ubah({ ...c, maxAge: s.id, maxAgeDetail: s.ekstra ?? '' })}
                muatSaatFokus
                wajibDariDaftar
              />
              <div className="field">
                <span className="field__label">{F.maxAgeDetail}</span>
                <div className="nbf-inward__teks">{c.maxAgeDetail}</div>
              </div>
            </div>
          </section>

        </div>
      </div>

      {/* Additional Ship selebar form (work owner 11-10-2026 "INI PERBAIKI DONG, MASA GINI": di kolom kanan Pega grid ini terjepit -
          Ship Name dan saran kapal terpotong). Urutan tetap sesudah blok pelayaran. */}
      <section className="nbf-alamat__blok nbf-aneka-form__kartu">
        <h5 className="nbf-objek__judul">{F.grupKapalTambahan}</h5>
        <GridKapalTambahan daftar={c.additionalShip} ubah={(additionalShip) => ubah({ ...c, additionalShip })} />
      </section>

      {popup === 'kapal' && (
        <PopupKapalCargo
          onTutup={() => setPopup(null)}
          onPilih={(b) => {
            ubah(terapkanKapal(c, b))
            setPopup(null)
          }}
        />
      )}
      {popup === 'agent' && (
        <PopupAgentCargo
          onTutup={() => setPopup(null)}
          onPilih={(a) => {
            ubah(terapkanAgent(c, a))
            setPopup(null)
          }}
        />
      )}
    </div>
  )
}

export default function TabObjectCargo({ caseId }: { caseId: string }) {
  const [cargo, setCargo] = useState<ObjekCargo[] | null>(null)
  const [master, setMaster] = useState<Master>(MASTER_KOSONG)
  const [buka, setBuka] = useState<number | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [proteksi, setProteksi] = useState<string[]>([])
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)

  useEffect(() => {
    let batal = false
    ambilObjekCargo(caseId).then(
      (h) => {
        if (!batal) setCargo(h.cargo.map(rapikanCargo))
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    for (const j of ['trading', 'goods', 'conveyance', 'packing'] as const) {
      daftarMasterCargo(j).then(
        (h) => {
          if (!batal) setMaster((m) => ({ ...m, [j]: h.baris }))
        },
        () => undefined,
      )
    }
    return () => {
      batal = true
    }
  }, [caseId])

  if (cargo === null) return galat !== null ? <Gagal galat={galat} /> : <Memuat />

  const daftar = cargo
  function ganti(baru: ObjekCargo[]) {
    setCargo(baru)
    setTersimpan(false)
  }

  async function simpan() {
    setMenyimpan(true)
    setGalat(null)
    setProteksi([])
    setTersimpan(false)
    try {
      const h = await simpanObjekCargo(caseId, daftar)
      setCargo(h.cargo.map(rapikanCargo))
      setTersimpan(true)
    } catch (err) {
      const p = pesanProteksiDari(err)
      if (p !== null) setProteksi(p)
      else setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const u = GRID_OBJEK.ukuran
  const jumlahHalaman = Math.max(1, Math.ceil(daftar.length / u))
  const hal = Math.min(halaman, jumlahHalaman)
  const awal = (hal - 1) * u
  const tampil = daftar.slice(awal, awal + u)
  const kolom = GRID_CARGO.length + 2

  return (
    <div className="nbf-objek nbf-obj-tab">
      <Gagal galat={galat} />
      {proteksi.length > 0 && (
        <div className="alert alert--error" role="alert">
          {TEKS_ANEKA.ditolak}
          <ul>
            {proteksi.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {tersimpan && <div className="alert alert--ok">{TEKS_ANEKA.tersimpan}</div>}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-aneka-grid nbf-aneka-grid--aksi nbf-aneka-grid--lokasi">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_CARGO.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button
                  type="button"
                  className="btn nbf-tema__tombol btn--ghost btn--sm"
                  onClick={() => {
                    ganti([...daftar, cargoBaru()])
                    setBuka(daftar.length)
                    setHalaman(Math.ceil((daftar.length + 1) / u))
                  }}
                >
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {daftar.length === 0 && (
              <tr>
                <td colSpan={kolom}>{TEKS_CARGO.kosong}</td>
              </tr>
            )}
            {tampil.map((c, n) => {
              const i = awal + n
              return [
                <tr key={i} className={buka === i ? 'nbf-baris--terbuka' : undefined}>
                  <td>
                    <button
                      type="button"
                      className={'btn btn--ghost btn--sm nbf-buka' + (buka === i ? ' nbf-buka--terbuka' : '')}
                      aria-label={TEKS_ANEKA.bukaBaris}
                      aria-expanded={buka === i}
                      onClick={() => setBuka(buka === i ? null : i)}
                    >
                      ▸
                    </button>
                  </td>
                  <td>{i + 1}</td>
                  <td>
                    <input type="checkbox" checked={c.isTopRisk} disabled aria-label={GRID_CARGO[1]} />
                  </td>
                  <td>{tampilTanggal(c.sailDate)}</td>
                  <td>{c.goodNote === '' ? <span className="nbf-aneka-grid__kosong">{TEKS_ANEKA.belumDiisi}</span> : c.goodNote}</td>
                  <td>{c.ship.nmShip}</td>
                  <td>{c.conveyanceNote}</td>
                  <td className="table__actions">
                    {c.flagDelete !== '1' && (
                      <button
                        type="button"
                        className="btn nbf-tema__tombol btn--ghost btn--sm"
                        onClick={() => {
                          ganti(daftar.filter((_, j) => j !== i))
                          setBuka(null)
                        }}
                      >
                        {GRID_OBJEK.hapus}
                      </button>
                    )}
                  </td>
                </tr>,
                buka === i && (
                  <tr key={`${i}-detail`} className="nbf-objek__detail">
                    <td colSpan={kolom}>
                      <div className="nbf-lapis nbf-lapis--objek">
                        <FormCargo c={c} idx={i} master={master} ubah={(baru) => ganti(daftar.map((x, j) => (j === i ? baru : x)))} />
                      </div>
                    </td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
      {daftar.length > u && <Halaman halaman={hal} ukuran={u} total={daftar.length} onPindah={setHalaman} />}
      <div className="nbf-objek__kaki">
        <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_OBJEK.label}
        </button>
      </div>
    </div>
  )
}
