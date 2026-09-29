// Layar tahun treaty — tiket 03 Treaty Contract Out.
//
// Meniru `Harness/InboxTreatyContract.xml` → `Section/GridTreatyContract.xml`
// (judul b1057) → `Section/InputTreatyContract.xml`: grid `BrowseTreatyYear_RD`
// (b17355; sort `.ID DESC` b672) dengan tombol baris `Edit` b19939 /
// `ReinsType` b20778 / `List Description` b22196, tombol `Add` b16387, dan
// form lengkap "Input New Data" dari `Section/InputDtlTreatyContact.xml`
// (disertakan b1206; `Save` b10332 → `SaveTreatyYear_Act`, `Cancel` b10622).
//
// ⛔ Identitas tidak pernah diketik (AC 5): medan `ID` hanya dibaca; baris
// baru dikirim TANPA id (POST), pembaruan lewat PUT /{id}.
//
// ⛔ Gerbang ada di SERVER: periode terbalik (422), anti-dobel (409), wajib
// isi (422). Layar menampilkan kalimatnya apa adanya (`Gagal`).
//
// ⛔ Tombol `Copy` b20459 dan form From/To b2374/b3818 TIDAK dibawa — fitur
// salin tahun treaty dibuang (AC 72). `ReinsType` dan `List Description`
// berdiri, menyebut tiket yang ditunggunya (04, 08) — bukan disembunyikan.
// Tiket 04 (29-09-2026): `ReinsType` kini membuka `PanelKontrakTahun`.
// Tiket 08 (29-09-2026): `List Description` kini membuka `PanelKlausulTahun`.
//
// ⚠️ OQ-TCO-05: label tahun bersilang antara grid dan form di korpus; keduanya
// dibawa apa adanya dan catatannya tampil di layar.
//
// Tiket 12: panel lampiran `GridTreatyArrangementAttachment` (b13074) tampil di
// form tahun yang sudah ber-ID - `components/PanelLampiranTahun`.

import { useCallback, useEffect, useState } from 'react'

import { LAMPIRAN_TCO, TAHUN_TCO } from '../labels'
import PanelKlausulTahun from '../components/PanelKlausulTahun'
import PanelKontrakTahun from '../components/PanelKontrakTahun'
import PanelLampiranTahun from '../components/PanelLampiranTahun'
import PilihJenisReasuransi from '../components/PilihJenisReasuransi'
import { Field, FieldTanggal, Gagal, Halaman, Kosong, Memuat, Pilih } from '../../../inti/components/ui/dasar'
import { formatDate } from '../../../inti/lib/format'
import { keInputTanggal } from '../../../inti/lib/tanggalInput'
import {
  ambilGrupTreaty,
  ambilTahunTreaty,
  simpanTahunTreaty,
  type GrupTreaty,
  type HalamanTahunTreaty,
  type TahunTreaty,
  type TahunTreatyMasuk,
} from '../api'

/** Isian form — nama medan mengikuti `InputTreatyYear.*`. */
export interface FormTahun {
  id: string
  treatyYear: string
  underwritingYear: string
  treatyGroupId: string
  treatyGroupName: string
  proportion: string
  /** Bentuk apa pun yang `keInputTanggal` kenal; dikirim sebagai YYYY-MM-DD. */
  startDate: string
  endDate: string
  /** Hanya dibaca — `Modified Date` b9097 dan `Username` b9282. */
  tglUpdate: string
  userId: string
}

/** Form baris BARU — `NewInputTreatyYear_Act` mengosongkan ID/TreatyYear/TreatyGroupID (b469–b512). */
export function formKosong(): FormTahun {
  return {
    id: '', treatyYear: '', underwritingYear: '', treatyGroupId: '', treatyGroupName: '',
    proportion: '', startDate: '', endDate: '', tglUpdate: '', userId: '',
  }
}

/** Form dari baris terpilih — `SetTreatyYear_Act` menyalin seluruh medan (b1110–b1353). */
export function formDari(t: TahunTreaty): FormTahun {
  return {
    id: t.id, treatyYear: t.treatyYear, underwritingYear: t.underwritingYear,
    treatyGroupId: t.treatyGroupId, treatyGroupName: t.treatyGroupName, proportion: t.proportion,
    startDate: t.startDate, endDate: t.endDate, tglUpdate: t.tglUpdate, userId: t.userId,
  }
}

/**
 * Badan permintaan dari form. Tanggal diseragamkan ke YYYY-MM-DD; teks yang
 * bukan tanggal dikirim APA ADANYA supaya server yang menolaknya dengan
 * pesan yang tepat — layar tidak menebak.
 */
export function keMasuk(f: FormTahun): TahunTreatyMasuk {
  const tanggal = (v: string): string => {
    const t = v.trim()
    if (t === '') return ''
    const iso = keInputTanggal(t)
    return iso === '' ? t : iso
  }
  return {
    id: f.id,
    treatyYear: f.treatyYear.trim(),
    underwritingYear: f.underwritingYear.trim(),
    treatyGroupId: f.treatyGroupId,
    treatyGroupName: f.treatyGroupName,
    proportion: f.proportion,
    startDate: tanggal(f.startDate),
    endDate: tanggal(f.endDate),
  }
}

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function selTahun(v: string): string {
  return v.trim() === '' ? '—' : v
}

/** Nama grup dari daftar master untuk ID terpilih; kosong bila tidak dikenal. */
export function namaGrup(daftar: readonly GrupTreaty[], id: string): string {
  return daftar.find((g) => g.id === id)?.treatyGroupName ?? ''
}

const UKURAN = 20

export default function InboxTreatyContract() {
  const [hal, setHal] = useState<HalamanTahunTreaty | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [form, setForm] = useState<FormTahun | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [grup, setGrup] = useState<GrupTreaty[]>([])
  const [galatGrup, setGalatGrup] = useState<unknown>(null)
  // Tiket 04: tombol `ReinsType` b20778 membuka editor kontrak tahun itu
  // (Pega: `BrowseReinsTypeYear` + `showHarness` popup `InboxTreatyContractReinsType`).
  const [tahunKontrak, setTahunKontrak] = useState<TahunTreaty | null>(null)
  // Tiket 08: tombol `List Description` b22196 membuka layar klausul tahun itu
  // (Pega: `showHarness` `InboxTreatyContractDescription`).
  const [tahunKlausul, setTahunKlausul] = useState<TahunTreaty | null>(null)

  const muat = useCallback(async (h: number) => {
    setSibuk(true)
    setGalat(null)
    try {
      setHal(await ambilTahunTreaty(h, UKURAN))
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
      setGalat(e)
      setHal(null)
    } finally {
      setSibuk(false)
    }
  }, [])

  useEffect(() => {
    void muat(halaman)
  }, [halaman, muat])

  // Master grup treaty dimuat saat form dibuka — kosong dinyatakan (503).
  useEffect(() => {
    if (form === null) return
    let hidup = true
    void (async () => {
      try {
        const d = await ambilGrupTreaty()
        if (hidup) setGrup(d.daftar)
      } catch (e) {
        if (hidup) setGalatGrup(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [form === null])

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatSimpan(null)
    try {
      await simpanTahunTreaty(keMasuk(form))
      setForm(null)
      await muat(halaman)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const ubah = (k: keyof FormTahun) => (v: string) => {
    setForm((f) => (f === null ? f : { ...f, [k]: v }))
  }

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TAHUN_TCO.judul}</h2>
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setGalatSimpan(null)
            setForm(formKosong())
          }}
        >
          {TAHUN_TCO.add}
        </button>
      </header>
      <p className="polis__catatan" role="note">
        {TAHUN_TCO.catatanLabelBersilang}
      </p>

      {form !== null && (
        <section className="panel">
          <h3 className="panel__title">{TAHUN_TCO.inputNewData}</h3>
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}
          {galatGrup !== null && <Gagal galat={galatGrup} />}
          <div className="form-grid">
            <Field label={TAHUN_TCO.formId} value={form.id} onChange={() => undefined} readOnly />
            <Pilih
              label={TAHUN_TCO.formTreatyGroup}
              value={form.treatyGroupId}
              onChange={(v) => {
                setForm((f) =>
                  f === null ? f : { ...f, treatyGroupId: v, treatyGroupName: namaGrup(grup, v) },
                )
              }}
              opsi={grup.map((g) => ({ value: g.id, label: g.treatyGroupName }))}
              required
            />
            <PilihJenisReasuransi
              label={TAHUN_TCO.formReinsuranceType}
              value={form.proportion}
              onChange={ubah('proportion')}
            />
            <FieldTanggal label={TAHUN_TCO.formStartDate} value={form.startDate} onChange={ubah('startDate')} />
            <FieldTanggal label={TAHUN_TCO.formEndDate} value={form.endDate} onChange={ubah('endDate')} />
            <Field label={TAHUN_TCO.formUnderwritingYear} value={form.treatyYear} onChange={ubah('treatyYear')} required />
            <Field label={TAHUN_TCO.formTransactionYear} value={form.underwritingYear} onChange={ubah('underwritingYear')} />
            <Field label={TAHUN_TCO.formModifiedDate} value={form.tglUpdate} onChange={() => undefined} readOnly />
            <Field label={TAHUN_TCO.formUsername} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {TAHUN_TCO.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatSimpan(null)
              }}
            >
              {TAHUN_TCO.cancel}
            </button>
          </div>
          {/* Tiket 12: panel lampiran (`InputTreatyContract.xml` b13074) melekat
              pada tahun treaty yang SUDAH ber-ID; lampiran bersifat opsional dan
              tidak menjadi syarat tersimpannya tahun treaty (AC 55). */}
          {form.id !== '' ? (
            <PanelLampiranTahun tahunID={form.id} />
          ) : (
            <p className="polis__catatan" role="note">
              {LAMPIRAN_TCO.simpanDulu}
            </p>
          )}
        </section>
      )}

      {tahunKontrak !== null && (
        <PanelKontrakTahun
          key={tahunKontrak.id}
          tahun={tahunKontrak}
          onTutup={() => {
            setTahunKontrak(null)
          }}
        />
      )}

      {tahunKlausul !== null && (
        <PanelKlausulTahun
          key={tahunKlausul.id}
          tahun={tahunKlausul}
          onTutup={() => {
            setTahunKlausul(null)
          }}
        />
      )}

      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length === 0 && <Kosong pesan={TAHUN_TCO.kosong} />}
      {hal !== null && hal.baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{TAHUN_TCO.kolomUnderwritingYear}</th>
              <th>{TAHUN_TCO.kolomTransactionYear}</th>
              <th>{TAHUN_TCO.kolomStartDate}</th>
              <th>{TAHUN_TCO.kolomEndDate}</th>
              <th>{TAHUN_TCO.kolomTreatyGroup}</th>
              <th>{TAHUN_TCO.kolomReinsuranceType}</th>
              <th className="table__actions" />
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr key={b.id} className="inbox__baris">
                <td>{selTahun(b.underwritingYear)}</td>
                <td>{selTahun(b.treatyYear)}</td>
                <td>{selTahun(formatDate(b.startDate))}</td>
                <td>{selTahun(formatDate(b.endDate))}</td>
                <td>{selTahun(b.treatyGroupName)}</td>
                <td>{selTahun(b.proportion)}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setGalatSimpan(null)
                      setForm(formDari(b))
                    }}
                  >
                    {TAHUN_TCO.edit}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setTahunKontrak(b)
                    }}
                  >
                    {TAHUN_TCO.reinsType}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      setTahunKlausul(b)
                    }}
                  >
                    {TAHUN_TCO.listDescription}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {hal !== null && (
        <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />
      )}
    </section>
  )
}
