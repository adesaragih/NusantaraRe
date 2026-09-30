// Panel business pada kombinasi kontrak — tiket 07 Treaty Contract Out.
//
// Meniru `Section/ViewDetailTreatyBusinessGrid.xml` (disertakan
// `InputTreatyContractReinsType.xml` b14064; dibuka tombol baris kontrak
// `Business List` b10842): kepala `Business List` b1988, `Add` b2785, grid
// `BrowseTreatyBusiness_RD` (`Treaty Group` b3422, `Business ID` b3566,
// `Business Name` b3710; `Edit` b4547, `Delete` b4826), form (`Business Name`
// b6241, `Active` b6499, `Business Code` b6680, `Save` b6966), `Information`
// b9319/b10064, `Close List` b10889.
//
// ⛔ Baris NONAKTIF tetap tampil dengan statusnya (AC 22) — grid Pega
// menyaringnya keluar (`BrowseTreatyBusiness_RD` `.IsActive = 1`).

import { useCallback, useEffect, useState } from 'react'

import { BUSINESS_TCO } from '../labels'
import {
  ambilBusinessKombinasi,
  ambilBusinessMaster,
  hapusBusinessKombinasi,
  simpanBusinessKombinasi,
  type BusinessMaster,
  type BusinessMasuk,
  type BusinessTreaty,
  type DaftarBusiness,
} from '../api'
import { Field, Gagal, Kosong, Memuat, Pilih } from '../../../inti/components/ui/dasar'

/** Isian form. */
export interface FormBusiness {
  id: string
  bizCode: string
  isActive: string
}

/** `NewTreatyBusinessDetail_Act` langkah 2 — ID kosong; Active wajib dipilih. */
export function formBusinessKosong(): FormBusiness {
  return { id: '', bizCode: '', isActive: '' }
}

/** `SetUbahTreatyBusinessList_Act` langkah 3 — ID, IsActive, BizCode (BIZNAME dari master). */
export function formBusinessDari(b: BusinessTreaty): FormBusiness {
  return { id: b.id, bizCode: b.bizCode, isActive: b.isActive }
}

/** Badan simpan. */
export function keMasukBusiness(f: FormBusiness): BusinessMasuk {
  return { id: f.id, bizCode: f.bizCode.trim(), isActive: f.isActive.trim() }
}

/** Status tampil — '1' aktif, selain itu nonaktif (hilir membaca `isactive='1'`). */
export function labelAktif(b: Pick<BusinessTreaty, 'aktif'>): string {
  return b.aktif ? BUSINESS_TCO.aktif : BUSINESS_TCO.nonaktif
}

export default function PanelBusinessKombinasi({
  tahunID,
  kontrakID,
  onTutup,
}: {
  tahunID: string
  kontrakID: string
  onTutup?: () => void
}) {
  const [daftar, setDaftar] = useState<DaftarBusiness | null>(null)
  const [master, setMaster] = useState<BusinessMaster[]>([])
  const [form, setForm] = useState<FormBusiness | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilBusinessKombinasi(tahunID, kontrakID))
    } catch (e) {
      setGalat(e)
    }
  }, [tahunID, kontrakID])

  useEffect(() => {
    void muat()
    ambilBusinessMaster()
      .then(setMaster)
      .catch((e: unknown) => {
        setGalat(e)
      })
  }, [muat])

  function buka(f: FormBusiness): void {
    setGalat(null)
    setInfo(null)
    setForm(f)
  }

  async function jalankan(kerja: () => Promise<string>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      setInfo(await kerja())
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
          {BUSINESS_TCO.judul}: {daftar?.kombinasi.reinsTypeName ?? '…'}
        </h3>
        {onTutup !== undefined && (
          <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
            {BUSINESS_TCO.closeList}
          </button>
        )}
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {info !== null && (
        <p role="status">
          {BUSINESS_TCO.information}: {info}
        </p>
      )}

      {form !== null && (
        <div className="panel">
          <div className="form-grid">
            <Pilih
              label={BUSINESS_TCO.formBusinessName}
              value={form.bizCode}
              onChange={(v) => {
                setForm((f) => (f === null ? f : { ...f, bizCode: v }))
              }}
              opsi={master.map((m) => ({ value: m.id, label: m.note }))}
              required
            />
            <Pilih
              label={BUSINESS_TCO.formActive}
              value={form.isActive}
              onChange={(v) => {
                setForm((f) => (f === null ? f : { ...f, isActive: v }))
              }}
              opsi={[
                { value: '1', label: BUSINESS_TCO.aktif },
                { value: '0', label: BUSINESS_TCO.nonaktif },
              ]}
              required
            />
            <Field label={BUSINESS_TCO.formBusinessCode} value={form.bizCode} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button
              type="button"
              className="btn btn--primary"
              disabled={sibuk}
              onClick={() =>
                void jalankan(async () => {
                  await simpanBusinessKombinasi(tahunID, kontrakID, keMasukBusiness(form))
                  // Simpan berhasil: form ditutup (keputusan work owner 30-09-2026).
                  setForm(null)
                  return BUSINESS_TCO.tersimpan
                })
              }
            >
              {BUSINESS_TCO.save}
            </button>
          </div>
        </div>
      )}

      <div className="aksi-baris">
        <button type="button" className="btn btn--primary" onClick={() => buka(formBusinessKosong())}>
          {BUSINESS_TCO.add}
        </button>
      </div>
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && daftar.daftar.length === 0 && <Kosong pesan={BUSINESS_TCO.kosong} />}
      {daftar !== null && daftar.daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{BUSINESS_TCO.kolomTreatyGroup}</th>
              <th>{BUSINESS_TCO.kolomBusinessId}</th>
              <th>{BUSINESS_TCO.kolomBusinessName}</th>
              <th>{BUSINESS_TCO.formActive}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {daftar.daftar.map((b) => (
              <tr key={b.id} className={b.aktif ? 'inbox__baris' : 'inbox__baris inbox__baris--pasif'}>
                <td>{b.treatyGroupName}</td>
                <td>{b.id}</td>
                <td>{b.bizName}</td>
                <td>{labelAktif(b)}</td>
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formBusinessDari(b))}>
                    {BUSINESS_TCO.edit}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => void jalankan(() => hapusBusinessKombinasi(tahunID, kontrakID, b.id))}
                  >
                    {BUSINESS_TCO.delete}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
