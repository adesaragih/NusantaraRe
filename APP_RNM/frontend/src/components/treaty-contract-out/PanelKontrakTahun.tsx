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
// ⛔ Tanggal akhir bawaan dihitung SERVER (`SetTanggalTreatyContract`,
// termasuk anomali OQ-TCO-10) — satu tempat, tidak ditulis ulang di sini.

import { useCallback, useEffect, useState } from 'react'

import { JENIS_REASURANSI_TCO, KONTRAK_TCO, TAHUN_TCO } from '../../assets/labels.treaty-contract-out'
import { formatDate } from '../../lib/format'
import { keInputTanggal } from '../../lib/tanggalInput'
import {
  ambilAkhirBawaanKontrak,
  ambilKontrakTahun,
  simpanKontrakTahun,
  type KontrakMasuk,
  type KontrakTreaty,
  type TahunTreaty,
} from '../../services/api'
import { Field, FieldTanggal, Gagal, Kosong, Memuat } from '../ui/dasar'
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

export default function PanelKontrakTahun({ tahun, onTutup }: { tahun: TahunTreaty; onTutup?: () => void }) {
  const [daftar, setDaftar] = useState<KontrakTreaty[] | null>(null)
  const [form, setForm] = useState<FormKontrak | null>(null)
  const [asal, setAsal] = useState<FormKontrak>(formKontrakKosong())
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)
  // Tiket 05: `Reinsurer List` b11308 membuka panel reinsurer kombinasi kontrak itu.
  const [kontrakReinsurer, setKontrakReinsurer] = useState<string | null>(null)

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

  function ubahMulai(v: string): void {
    setForm((f) => (f === null ? f : { ...f, treatyStartDate: v }))
    const iso = keInputTanggal(v)
    if (iso === '') return
    ambilAkhirBawaanKontrak(tahun.id, iso)
      .then((akhir) => {
        setForm((f) => (f === null ? f : { ...f, treatyEndDate: akhir }))
      })
      .catch((e: unknown) => {
        setGalat(e)
      })
  }

  function ubahAkhir(v: string): void {
    setForm((f) => (f === null ? f : { ...f, treatyEndDate: v }))
  }

  async function simpan(): Promise<void> {
    if (form === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      const k = await simpanKontrakTahun(tahun.id, keMasukKontrak(form))
      buka(formKontrakDari(k))
      setInfo(KONTRAK_TCO.tersimpan)
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
      <p className="polis__catatan" role="note">
        {KONTRAK_TCO.catatanLabelBersilang}
      </p>

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
            <FieldTanggal label={KONTRAK_TCO.formStartDate} value={form.treatyStartDate} onChange={ubahMulai} />
            <FieldTanggal label={KONTRAK_TCO.formEndDate} value={form.treatyEndDate} onChange={ubahAkhir} />
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
        <button type="button" className="btn btn--primary" onClick={() => buka(formKontrakKosong())}>
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
              <tr key={k.id} className="inbox__baris">
                <td>{sel(k.reinsTypeName)}</td>
                <td>{sel(formatDate(k.treatyStartDate))}</td>
                <td>{sel(formatDate(k.treatyEndDate))}</td>
                <td className="table__actions">
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(formKontrakDari(k))}>
                    {KONTRAK_TCO.edit}
                  </button>{' '}
                  <button type="button" className="btn btn--ghost btn--sm" disabled title={`${TAHUN_TCO.menungguTiket} 07`}>
                    {KONTRAK_TCO.businessList}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setKontrakReinsurer(k.id)
                    }}
                  >
                    {KONTRAK_TCO.reinsurerList}
                  </button>{' '}
                  <button type="button" className="btn btn--ghost btn--sm" disabled title={`${TAHUN_TCO.menungguTiket} 10`}>
                    {KONTRAK_TCO.delete}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {kontrakReinsurer !== null && (
        <PanelReinsurerKombinasi
          key={kontrakReinsurer}
          tahunID={tahun.id}
          kontrakID={kontrakReinsurer}
          onTutup={() => {
            setKontrakReinsurer(null)
          }}
        />
      )}
    </section>
  )
}
