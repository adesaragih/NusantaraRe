// Editor kontrak treaty di dalam tahun treaty — tiket 04 Treaty Contract Out.
//
// Meniru harness `InboxTreatyContractReinsType` (judul `ReinsType` b1670) →
// `PanggilReinsType` → `Section/InputTreatyContractReinsType.xml`: kepala
// tahun (b1145/b1358), form kontrak (`ReinsType` b2652, `Start Date` b2905,
// `End Date` b3244, `Save` b3618, `Undo` b5343, `Modified Date` b4363,
// `Username` b4547, `Information` b6400), grid `BrowseTreatyContract_RD`
// (`Add` b8528; kolom b9164/b9304/b9444; tombol baris `Edit` b10519,
// `Business List` b10842, `Reinsurer List` b11308, `Delete` b11809).
//
// Di Pega layar ini popup dari tombol `ReinsType` baris tahun treaty
// (`InputTreatyContract.xml` b20778 → `showHarness` popup). Di sini ia panel
// yang dibuka dari baris itu, atau dari butir menu dengan pemilih tahun.
//
// ⛔ [keputusan work owner 30-09-2026] Start Date dan End Date kontrak HANYA
// DIBACA dan = tanggal TAHUN treaty induknya ("datanya diambil dari depan") —
// server yang menetapkannya (`KontrakTreatyTCO.Simpan`); isian tanggal dan
// End Date bawaan +1 tahun di form ini dibuang.
//
// [keputusan work owner 30-09-2026] `Business List` dan `Reinsurer List`
// membuka SATU panel rinci sekaligus, tepat di bawah baris kontraknya: membuka
// yang lain (kontrak lain, atau daftar lain) menutup yang terbuka, menekan
// tombol yang sama menutupnya.

import { Fragment, useCallback, useEffect, useState } from 'react'

import { JENIS_REASURANSI_TCO, KONTRAK_TCO } from '../labels'
import { formatDate } from '../../../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../../../inti/frontend/lib/tanggalInput'
import {
  ambilDampakHapusKontrak,
  ambilKontrakTahun,
  hapusKontrak,
  simpanKontrakTahun,
  type KontrakMasuk,
  type KontrakTreaty,
  type TahunTreaty,
  type DampakHapusTCO,
} from '../api'
import { Field, Gagal, Kosong, Memuat } from '../../../../../inti/frontend/components/ui/dasar'
import KonfirmasiHapusTCO from './KonfirmasiHapusTCO'
import PanelBusinessKombinasi from './PanelBusinessKombinasi'
import PanelReinsurerKombinasi from './PanelReinsurerKombinasi'
import PilihJenisReasuransi from './PilihJenisReasuransi'

/** Isian form — nama medan mengikuti `InputTreatyContract.*`. */
export interface FormKontrak {
  id: string
  reinsTypeId: string
  treatyStartDate: string
  treatyEndDate: string
  /** Hanya dibaca. */
  tglUpdate: string
  userId: string
}

/** Baris baru — `NewInputTreatyContract_Act` b502–b587. */
export function formKontrakKosong(): FormKontrak {
  return { id: '', reinsTypeId: '', treatyStartDate: '', treatyEndDate: '', tglUpdate: '', userId: '' }
}

/** Dari baris grid — `SetUbahTreatyContract` b314–b467. */
export function formKontrakDari(k: KontrakTreaty): FormKontrak {
  return {
    id: k.id, reinsTypeId: k.reinsTypeId, treatyStartDate: k.treatyStartDate, treatyEndDate: k.treatyEndDate,
    tglUpdate: k.tglUpdate, userId: k.userId,
  }
}

/** Tanggal kontrak = tanggal tahun treaty induknya (keputusan work owner 30-09-2026). */
export function denganTanggalTahun(f: FormKontrak, tahun: Pick<TahunTreaty, 'startDate' | 'endDate'>): FormKontrak {
  return { ...f, treatyStartDate: tahun.startDate, treatyEndDate: tahun.endDate }
}

function tanggalKirim(v: string): string {
  const t = v.trim()
  if (t === '') return ''
  const iso = keInputTanggal(t)
  return iso === '' ? t : iso
}

/** Badan simpan. Nama jenis TIDAK dikirim — server mengambilnya dari master. */
export function keMasukKontrak(f: FormKontrak): KontrakMasuk {
  return {
    id: f.id,
    reinsTypeId: f.reinsTypeId,
    treatyStartDate: tanggalKirim(f.treatyStartDate),
    treatyEndDate: tanggalKirim(f.treatyEndDate),
  }
}

function sel(v: string): string {
  return v.trim() === '' ? '—' : v
}

/** Panel rinci satu kontrak yang terbuka — Business List atau Reinsurer List. */
export type RinciKontrak = { daftar: 'business' | 'reinsurer'; kontrakID: string }

/** Tombol daftar satu baris: tombol yang sama menutup, yang lain menggantikan. */
export function alihRinci(
  terbuka: RinciKontrak | null,
  daftar: RinciKontrak['daftar'],
  kontrakID: string,
): RinciKontrak | null {
  return terbuka !== null && terbuka.daftar === daftar && terbuka.kontrakID === kontrakID
    ? null
    : { daftar, kontrakID }
}

export default function PanelKontrakTahun({ tahun, onTutup }: { tahun: TahunTreaty; onTutup?: () => void }) {
  const [daftar, setDaftar] = useState<KontrakTreaty[] | null>(null)
  const [form, setForm] = useState<FormKontrak | null>(null)
  const [asal, setAsal] = useState<FormKontrak>(formKontrakKosong())
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)
  // Tiket 05: `Reinsurer List` b11308 membuka panel reinsurer kombinasi kontrak itu.
  // Tiket 07: `Business List` b10842 membuka panel bisnis kombinasi kontrak itu.
  // SATU keadaan untuk keduanya (keputusan work owner 30-09-2026).
  const [rinci, setRinci] = useState<RinciKontrak | null>(null)
  // Tiket 10: popup Ya/Batal sebelum kaskade hapus kontrak.
  const [konfirmasi, setKonfirmasi] = useState<{ kontrak: KontrakTreaty; dampak: DampakHapusTCO | null; galat: unknown } | null>(null)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilKontrakTahun(tahun.id))
    } catch (e) {
      setGalat(e)
    }
  }, [tahun.id])

  useEffect(() => {
    void muat()
  }, [muat])

  function buka(f: FormKontrak): void {
    setGalat(null)
    setInfo(null)
    setAsal(f)
    setForm(f)
  }

  async function simpan(): Promise<void> {
    if (form === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      await simpanKontrakTahun(tahun.id, keMasukKontrak(form))
      // Simpan berhasil: form ditutup (keputusan work owner 30-09-2026).
      setForm(null)
      setInfo(KONTRAK_TCO.tersimpan)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  function mintaHapus(k: KontrakTreaty): void {
    setInfo(null)
    setKonfirmasi({ kontrak: k, dampak: null, galat: null })
    ambilDampakHapusKontrak(tahun.id, k.id)
      .then((d) => {
        setKonfirmasi((c) => (c === null || c.kontrak.id !== k.id ? c : { ...c, dampak: d }))
      })
      .catch((e: unknown) => {
        setKonfirmasi((c) => (c === null ? c : { ...c, galat: e }))
      })
  }

  async function yaHapus(): Promise<void> {
    if (konfirmasi === null || konfirmasi.dampak === null || sibuk) return
    setSibuk(true)
    try {
      const pesan = await hapusKontrak(tahun.id, konfirmasi.kontrak.id, konfirmasi.dampak)
      if (form?.id === konfirmasi.kontrak.id) setForm(null)
      if (rinci?.kontrakID === konfirmasi.kontrak.id) setRinci(null)
      setKonfirmasi(null)
      setInfo(pesan)
      await muat()
    } catch (e) {
      // Temuan /code-review: angka popup dimuat ulang sesudah galat (mis. 409
      // "jumlah berubah") - Ya tidak lagi mengirim angka basi berulang-ulang.
      const id = konfirmasi.kontrak.id
      setKonfirmasi((c) => (c === null ? c : { ...c, galat: e, dampak: null }))
      ambilDampakHapusKontrak(tahun.id, id)
        .then((d) => {
          setKonfirmasi((c) => (c === null || c.kontrak.id !== id ? c : { ...c, dampak: d }))
        })
        .catch(() => undefined)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <section className="panel">
      <header className="inbox__kepala">
        <h3 className="panel__title">{KONTRAK_TCO.judul}</h3>
        {onTutup !== undefined && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {KONTRAK_TCO.tutup}
          </button>
        )}
      </header>
      <div className="form-grid">
        <Field label={KONTRAK_TCO.headerUnderwritingYear} value={tahun.underwritingYear} onChange={() => undefined} readOnly />
        <Field label={KONTRAK_TCO.headerReinsType} value={tahun.treatyGroupName} onChange={() => undefined} readOnly />
      </div>

      {galat !== null && <Gagal galat={galat} />}
      {info !== null && (
        <p role="status">
          {KONTRAK_TCO.information}: {info}
        </p>
      )}

      {form !== null && (
        <div className="panel">
          <div className="form-grid">
            <Field label={KONTRAK_TCO.formId} value={form.id} onChange={() => undefined} readOnly />
            <PilihJenisReasuransi
              label={JENIS_REASURANSI_TCO.reinsType}
              value={form.reinsTypeId}
              onChange={(v) => {
                setForm((f) => (f === null ? f : { ...f, reinsTypeId: v }))
              }}
            />
            <Field label={KONTRAK_TCO.formStartDate} value={formatDate(form.treatyStartDate)} onChange={() => undefined} readOnly />
            <Field label={KONTRAK_TCO.formEndDate} value={formatDate(form.treatyEndDate)} onChange={() => undefined} readOnly />
            <Field label={KONTRAK_TCO.formModifiedDate} value={form.tglUpdate} onChange={() => undefined} readOnly />
            <Field label={KONTRAK_TCO.formUsername} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void simpan()}>
              {KONTRAK_TCO.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              disabled={sibuk}
              onClick={() => {
                setGalat(null)
                setForm(asal)
              }}
            >
              {KONTRAK_TCO.undo}
            </button>
          </div>
        </div>
      )}

      <div className="aksi-baris">
        <button type="button" className="btn btn--primary" onClick={() => buka(denganTanggalTahun(formKontrakKosong(), tahun))}>
          {KONTRAK_TCO.add}
        </button>
      </div>
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={KONTRAK_TCO.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{KONTRAK_TCO.kolomReinsType}</th>
              <th>{KONTRAK_TCO.kolomTreatyStart}</th>
              <th>{KONTRAK_TCO.kolomTreatyEnd}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.map((k) => (
              <Fragment key={k.id}>
                <tr className={rinci?.kontrakID === k.id ? 'inbox__baris belah__baris--aktif' : 'inbox__baris'}>
                  <td>{sel(k.reinsTypeName)}</td>
                  <td>{sel(formatDate(k.treatyStartDate))}</td>
                  <td>{sel(formatDate(k.treatyEndDate))}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(denganTanggalTahun(formKontrakDari(k), tahun))}>
                      {KONTRAK_TCO.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      aria-expanded={rinci?.daftar === 'business' && rinci.kontrakID === k.id}
                      onClick={() => {
                        setRinci((r) => alihRinci(r, 'business', k.id))
                      }}
                    >
                      {KONTRAK_TCO.businessList}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      aria-expanded={rinci?.daftar === 'reinsurer' && rinci.kontrakID === k.id}
                      onClick={() => {
                        setRinci((r) => alihRinci(r, 'reinsurer', k.id))
                      }}
                    >
                      {KONTRAK_TCO.reinsurerList}
                    </button>{' '}
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => mintaHapus(k)}>
                      {KONTRAK_TCO.delete}
                    </button>
                  </td>
                </tr>
                {rinci?.kontrakID === k.id && (
                  <tr className="inbox__rinci">
                    <td colSpan={4}>
                      {rinci.daftar === 'business' ? (
                        <PanelBusinessKombinasi
                          key={`business/${k.id}`}
                          tahunID={tahun.id}
                          kontrakID={k.id}
                          onTutup={() => {
                            setRinci(null)
                          }}
                        />
                      ) : (
                        <PanelReinsurerKombinasi
                          key={`reinsurer/${k.id}`}
                          tahunID={tahun.id}
                          kontrakID={k.id}
                          onTutup={() => {
                            setRinci(null)
                          }}
                        />
                      )}
                    </td>
                  </tr>
                )}
              </Fragment>
            ))}
          </tbody>
        </table>
      )}

      {konfirmasi !== null && (
        <KonfirmasiHapusTCO
          jenis="kontrak"
          nama={konfirmasi.kontrak.reinsTypeName || konfirmasi.kontrak.reinsTypeId}
          dampak={konfirmasi.dampak}
          galat={konfirmasi.galat}
          sibuk={sibuk}
          onYa={() => void yaHapus()}
          onBatal={() => {
            setKonfirmasi(null)
          }}
        />
      )}
    </section>
  )
}
