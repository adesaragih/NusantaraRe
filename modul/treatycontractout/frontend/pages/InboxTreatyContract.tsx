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
// salin tahun treaty dibuang (AC 72). `ReinsType` membuka `PanelKontrakTahun`
// (tiket 04), `List Description` membuka `PanelKlausulTahun` (tiket 08).
//
// [keputusan work owner 30-09-2026] Panel yang dibuka dari baris tahun hanya
// SATU sekaligus, dan selama ia terbuka tabel utama DISEMBUNYIKAN — fokus pada
// grid yang diklik; menutup panel mengembalikan tabel. Catatan pengembang di
// layar (label bersilang OQ-TCO-05, "simpan dulu") dibuang.
//
// ⚠️ OQ-TCO-05: label tahun bersilang antara grid dan form di korpus; keduanya
// dibawa apa adanya.

import { useCallback, useEffect, useRef, useState } from 'react'

import { TAHUN_TCO } from '../labels'
import PanelKlausulTahun from '../components/PanelKlausulTahun'
import PanelKontrakTahun from '../components/PanelKontrakTahun'
import { PROPORSI_TCO, labelProporsi } from '../proporsi'
import { Field, FieldTanggal, Gagal, Halaman, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { formatDate } from '../../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import {
  ambilAkhirBawaanTahun,
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

/** Tahun (YYYY) sebuah tanggal isian; kosong bila belum tanggal sah. */
export function tahunDari(tanggal: string): string {
  const iso = keInputTanggal(tanggal.trim())
  return iso === '' ? '' : iso.slice(0, 4)
}

/**
 * Start Date diisi. Form BARU [keputusan work owner 30-09-2026]: Underwriting
 * Year dan Transaction Year mengambil tahun Start Date — keduanya tetap dapat
 * diubah sesudahnya. Form yang sudah ber-ID tidak diisi ulang.
 */
export function isiDariMulai(f: FormTahun, mulai: string): FormTahun {
  const tahun = f.id === '' ? tahunDari(mulai) : ''
  return tahun === '' ? { ...f, startDate: mulai } : { ...f, startDate: mulai, treatyYear: tahun, underwritingYear: tahun }
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

/** Panel yang dibuka dari satu baris tahun — `ReinsType` atau `List Description`. */
export type RinciTahun = { jenis: 'kontrak' | 'klausul'; tahun: TahunTreaty }

export default function InboxTreatyContract() {
  const [hal, setHal] = useState<HalamanTahunTreaty | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [form, setForm] = useState<FormTahun | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [grup, setGrup] = useState<GrupTreaty[]>([])
  const [galatGrup, setGalatGrup] = useState<unknown>(null)
  // Tiket 04: tombol `ReinsType` b20778 membuka editor kontrak tahun itu
  // (Pega: `BrowseReinsTypeYear` + `showHarness` popup `InboxTreatyContractReinsType`).
  // Tiket 08: tombol `List Description` b22196 membuka layar klausul tahun itu
  // (Pega: `showHarness` `InboxTreatyContractDescription`). SATU keadaan untuk
  // keduanya: membuka yang satu menutup yang lain.
  const [rinci, setRinci] = useState<RinciTahun | null>(null)
  // Posisi gulir tabel saat panel dibuka - dikembalikan saat panel ditutup.
  const gulirTabel = useRef(0)
  // Urutan permintaan End Date bawaan - jawaban basi dibuang.
  const urutanAkhir = useRef(0)

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
      setInfo(TAHUN_TCO.tersimpan)
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

  // Form BARU: Start Date mengisi tahun (di sini) dan End Date (server, aturan
  // yang sama dengan kontrak); semuanya tetap dapat diubah. Hanya jawaban
  // TERAKHIR yang dipakai, dan hanya bila Start Date belum berubah lagi.
  function ubahMulai(v: string): void {
    setForm((f) => (f === null ? f : isiDariMulai(f, v)))
    const iso = keInputTanggal(v.trim())
    if (form === null || form.id !== '' || iso === '') return
    const ke = ++urutanAkhir.current
    // End Date yang diketik pemakai SELAMA permintaan berjalan tidak ditimpa.
    const akhirSaatMinta = form.endDate
    ambilAkhirBawaanTahun(iso)
      .then((akhir) => {
        if (ke !== urutanAkhir.current) return
        setForm((f) =>
          f === null || f.id !== '' || f.startDate !== v || f.endDate !== akhirSaatMinta ? f : { ...f, endDate: akhir },
        )
      })
      .catch((e: unknown) => {
        if (ke !== urutanAkhir.current) return
        setGalatSimpan(e)
      })
  }

  function bukaRinci(r: RinciTahun): void {
    if (typeof window !== 'undefined') gulirTabel.current = window.scrollY
    setRinci(r)
    // Tabel panjang tersembunyi: mulai dari atas panel, bukan dari tengah layar.
    if (typeof window !== 'undefined') window.scrollTo({ top: 0 })
  }

  // Fokus pada panel yang diklik: tabel utama, form, dan penomoran halaman
  // TIDAK dirender; `Tutup` panel mengembalikannya (state tabel tetap).
  if (rinci !== null) {
    const tutup = () => {
      setRinci(null)
      // Kembali ke baris yang tadi diklik, sesudah tabel dirender ulang.
      const y = gulirTabel.current
      if (typeof window !== 'undefined') window.requestAnimationFrame(() => window.scrollTo({ top: y }))
    }
    return (
      <section className="inbox tco">
        {rinci.jenis === 'kontrak' ? (
          <PanelKontrakTahun key={rinci.tahun.id} tahun={rinci.tahun} onTutup={tutup} />
        ) : (
          <PanelKlausulTahun key={rinci.tahun.id} tahun={rinci.tahun} onTutup={tutup} />
        )}
      </section>
    )
  }

  return (
    <section className="inbox tco">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TAHUN_TCO.judul}</h2>
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setGalatSimpan(null)
            setInfo(null)
            setForm(formKosong())
          }}
        >
          {TAHUN_TCO.add}
        </button>
      </header>
      {info !== null && <p role="status">{info}</p>}

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
            {/* OQ-TCO-04 (keputusan work owner 30-09-2026): dua nilai, bukan master jenis reasuransi. */}
            <Pilih
              label={TAHUN_TCO.formReinsuranceType}
              value={form.proportion}
              onChange={ubah('proportion')}
              opsi={PROPORSI_TCO.map((p) => ({ value: p.value, label: p.label }))}
            />
            <FieldTanggal label={TAHUN_TCO.formStartDate} value={form.startDate} onChange={ubahMulai} />
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
        </section>
      )}

      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length === 0 && <Kosong pesan={TAHUN_TCO.kosong} />}
      {hal !== null && hal.baris.length > 0 && (
        <div className="tco-tabel">
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
                  <td>{selTahun(labelProporsi(b.proportion))}</td>
                  <td className="table__actions">
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setGalatSimpan(null)
                        setInfo(null)
                        setForm(formDari(b))
                      }}
                    >
                      {TAHUN_TCO.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        bukaRinci({ jenis: 'kontrak', tahun: b })
                      }}
                    >
                      {TAHUN_TCO.reinsType}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        bukaRinci({ jenis: 'klausul', tahun: b })
                      }}
                    >
                      {TAHUN_TCO.listDescription}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {hal !== null && (
        <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />
      )}
    </section>
  )
}
