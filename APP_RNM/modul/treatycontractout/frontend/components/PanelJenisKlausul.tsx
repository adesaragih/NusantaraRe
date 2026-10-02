// Panel satu jenis klausul — tiket 08 Treaty Contract Out.
//
// Dirakit dari ATURAN yang dikirim server (`models.AturanKlausulTCO`): medan
// form, wajib-isi, turunan, alasan ditahan. Satu komponen untuk ke-25 jenis —
// padanan 25 section `GridTreatyArrangement*` / `GridTreatyArr*List` yang
// masing-masing punya `Add` / `Edit` / `Save` / `Show Child` / `Close Child`.
//
// ⛔ Setiap grid aturan memegang isiannya SENDIRI; `Cancel` membuang isian grid
// itu saja. Satu jenis tampil sekaligus, di popup (`PanelKlausulTahun`,
// keputusan work owner 30-09-2026 — menggantikan AC 29).
// ⛔ Nilai uang/persen TEKS sepanjang jalan; Rp/Usd anak dihitung server.
//
// Tiket 11: jenis berkurs menampilkan kurs berlaku tahun itu (`testingKurs`);
// tanpa kurs, pesan server tampil dan `Add` nonaktif (Pega: `DATASHOW = ""`).
// Pratinjau konversi Rp ↔ Usd dihitung SERVER (`konversiKurs`) — bukan JS.
//
// [keputusan work owner 30-09-2026] Nilai desimal TAMPIL dengan pemisah ribuan
// (`formatNumber`, gaya Indonesia, nol digit dibuang); isian form tetap teks
// mentah. Catatan pengembang ("dihitung server") tidak tampil di layar.

import { useCallback, useEffect, useRef, useState } from 'react'

import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../../inti/frontend/lib/format'
import { KLAUSUL_TCO, KURS_TCO, LABEL_MEDAN_KHUSUS, LABEL_MEDAN_KLAUSUL } from '../labels'
import {
  ambilKlausul,
  ambilKursTahun,
  konversiKurs,
  PILIHAN_REINS_ANAK_TREATY_LIMIT,
  simpanKlausul,
  type AturanKlausul,
  type DaftarKlausul,
  type JenisKlausul,
  type Klausul,
  type KlausulMasuk,
  type KursTahun,
} from '../api'
import { Field, Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import PilihJenisReasuransi from './PilihJenisReasuransi'
import PilihJenisReasuransiSaring from './PilihJenisReasuransiSaring'
import PilihMasterKlausul from './PilihMasterKlausul'

/** Label medan: penimpaan per jenis/subjenis, lalu bawaan, lalu nama medan. */
export function labelMedan(a: Pick<AturanKlausul, 'jenis' | 'subjenis'>, medan: string): string {
  const kunci = a.subjenis !== '' ? `${a.jenis}/${a.subjenis}` : a.jenis
  const khusus: Readonly<Record<string, Readonly<Record<string, string>>>> = LABEL_MEDAN_KHUSUS
  const bawaan: Readonly<Record<string, string>> = LABEL_MEDAN_KLAUSUL
  return khusus[kunci]?.[medan] ?? bawaan[medan] ?? medan
}

/** Isian form — satu entri per medan aturan, teks. */
export type FormKlausul = { id: string; medan: Record<string, string> }

/** Form kosong untuk aturan itu (`NewTreatyArr*`). */
export function formKlausulKosong(a: AturanKlausul): FormKlausul {
  return { id: '', medan: Object.fromEntries(a.medan.map((m) => [m, ''])) }
}

/** Form dari baris (`SetTreatyArr*_Act`) — hanya medan milik aturan. */
export function formKlausulDari(a: AturanKlausul, k: Klausul): FormKlausul {
  return { id: k.id, medan: Object.fromEntries(a.medan.map((m) => [m, k.medan[m] ?? ''])) }
}

/** Badan simpan; medan turunan (Rp/Usd anak) tidak pernah dikirim. */
export function keMasukKlausul(a: AturanKlausul, descId: string, f: FormKlausul, induk: string): KlausulMasuk {
  const turunan = new Set(a.turunan ?? [])
  const medan: Record<string, string> = {}
  for (const m of a.medan) {
    if (!turunan.has(m)) medan[m] = (f.medan[m] ?? '').trim()
  }
  return { id: f.id, descId, anak: a.anak, subjenis: a.subjenis, parentReinsTypeId: a.anak ? induk : '', medan }
}

/**
 * Medan klausul yang DESIMAL (`models.NilaiMedanKlausul`: `d(k.Rp)` …) — hanya
 * ini yang diberi pemisah ribuan; `Line`, `Layer`, kode, dan teks apa adanya.
 */
export const MEDAN_DESIMAL_KLAUSUL: ReadonlySet<string> = new Set([
  'Rp', 'Usd', 'Pct', 'PctMe', 'CoIns_Min', 'CoIns_Max', 'TreatyLimit',
])

/** Isi sel grid satu medan: desimal berpemisah ribuan, sisanya apa adanya. */
export function tampilMedanKlausul(medan: string, nilai: string): string {
  return MEDAN_DESIMAL_KLAUSUL.has(medan) ? formatNumber(nilai, DESIMAL_TAK_DIBATASI) : nilai
}

/** Baris exclusion milik subjenis itu (subjenis diturunkan server). */
export function barisSubjenis(daftar: Klausul[], subjenis: string): Klausul[] {
  return subjenis === '' ? daftar : daftar.filter((k) => k.subjenis === subjenis)
}

/** Aturan induk jenis (satu, atau satu per subjenis ExclutionTreaty). */
export function aturanInduk(j: JenisKlausul): AturanKlausul[] {
  return j.aturan.filter((a) => !a.anak)
}

/**
 * Tab subjenis: jenis berinduk lebih dari satu (10013 Exclusion Treaty — Occupation, Clause, Object,
 * Periode) tampil sebagai tab, bukan bertumpuk ke bawah [keputusan work owner 02-10-2026]. Kosong =
 * tanpa tab.
 */
export function tabSubjenis(j: JenisKlausul): string[] {
  const induk = aturanInduk(j)
  return induk.length > 1 ? induk.map((a) => a.subjenis) : []
}

/** Aturan anak jenis, bila ada. */
export function aturanAnak(j: JenisKlausul): AturanKlausul | undefined {
  return j.aturan.find((a) => a.anak)
}

/** Jenis menuntut kurs bila salah satu aturannya berkurs (tiket 11). */
export function jenisBerkurs(j: JenisKlausul): boolean {
  return j.aturan.some((a) => a.berkurs)
}

/**
 * Konversi yang dipicu satu medan: `RpKeUsd` — Rp → Usd skala 8
 * (`HitungRpUsd_depan`); `DuaArah` — Rp → Usd skala 4, Usd → Rp
 * (`CalculateTSIExcludeTreaty`). `null` bila medan itu tidak memicu apa pun.
 */
export function rencanaKonversi(
  konversi: string,
  medan: string,
): { dari: 'Rp' | 'Usd'; ke: 'Rp' | 'Usd'; skala: '4' | '8' } | null {
  if (konversi === 'RpKeUsd' && medan === 'Rp') return { dari: 'Rp', ke: 'Usd', skala: '8' }
  if (konversi === 'DuaArah' && medan === 'Rp') return { dari: 'Rp', ke: 'Usd', skala: '4' }
  if (konversi === 'DuaArah' && medan === 'Usd') return { dari: 'Usd', ke: 'Rp', skala: '8' }
  return null
}

/**
 * Pemilih ReinsTypeID satu aturan. Treaty Limit ikut XML: `pxAutoComplete` di
 * grid induk (`GridTreatyArrangementTreatyLimit.xml` b3025, daftar induk) dan
 * anak (`GridTreatyArrTreatyLimitList.xml` b2892, jenis porsi saja tanpa induk —
 * dari penanda aturan `pilihanReins`). ⛔ SETIAP baris anak membawa penanda itu
 * [keputusan work owner 02-10-2026: ReinsType anak semua jenis = anak Treaty
 * Limit]. Induk jenis lain: dropdown daftar induk tiket 02.
 */
export function pemilihReinsType(a: AturanKlausul): 'saring-induk' | 'saring-anak' | 'dropdown' {
  if (a.pilihanReins === PILIHAN_REINS_ANAK_TREATY_LIMIT) return 'saring-anak'
  return a.jenis === 'TreatyLimit' ? 'saring-induk' : 'dropdown'
}

/**
 * `Add` tampil? Jenis satu baris (`satuBaris`: `GridTreatyArrangementMinLOL.xml` b2232,
 * `…MaxCoinsPanel.xml` b2212, `…MInLOLMB.xml` b2262 — `Add` hanya bila `ID == ''`) menyembunyikannya
 * selama daftarnya belum dimuat atau sudah berisi; barisnya diubah lewat `Edit`.
 */
export function addTampil(a: Pick<AturanKlausul, 'satuBaris'>, dimuat: boolean, cacahBaris: number): boolean {
  return a.satuBaris !== true || (dimuat && cacahBaris === 0)
}

function FormMedan({
  tahunID,
  aturan,
  form,
  onUbah,
}: {
  tahunID: string
  aturan: AturanKlausul
  form: FormKlausul
  onUbah: (medan: string, nilai: string) => void
}) {
  const turunan = new Set(aturan.turunan ?? [])
  // Hanya jawaban konversi TERAKHIR yang dipakai (ketikan cepat).
  const urutan = useRef(0)

  function ubahDanKonversi(medan: string, nilai: string): void {
    onUbah(medan, nilai)
    const r = rencanaKonversi(aturan.konversi, medan)
    if (r === null || nilai.trim() === '') return
    const ke = ++urutan.current
    konversiKurs(tahunID, r.dari, nilai, r.skala)
      .then((h) => {
        if (ke === urutan.current) onUbah(r.ke, r.ke === 'Usd' ? h.usd : h.rp)
      })
      // Ketikan setengah jadi bukan desimal sah; simpan tetap diperiksa server.
      .catch(() => undefined)
  }
  return (
    <div className="form-grid">
      {aturan.medan.map((m) => {
        const label = labelMedan(aturan, m)
        if (m === 'ReinsTypeID') {
          const pemilih = pemilihReinsType(aturan)
          if (pemilih === 'dropdown') {
            return <PilihJenisReasuransi key={m} label={label} value={form.medan[m] ?? ''} onChange={(v) => onUbah(m, v)} />
          }
          return (
            <PilihJenisReasuransiSaring
              key={m}
              label={label}
              value={form.medan[m] ?? ''}
              onChange={(v) => onUbah(m, v)}
              anak={pemilih === 'saring-anak'}
            />
          )
        }
        if (m === 'ID_Occupation' || m === 'ID_Clause') {
          // Satu dropdown yang dapat difilter, tanpa kotak Search terpisah [keputusan work owner 02-10-2026].
          const namaMedan = m === 'ID_Occupation' ? 'Occupation' : 'Clause'
          return (
            <PilihMasterKlausul
              key={m}
              master={m === 'ID_Occupation' ? 'occupation' : 'clause'}
              label={label}
              value={form.medan[m] ?? ''}
              nama={form.medan[namaMedan] ?? ''}
              onPilih={(id, nama) => {
                onUbah(m, id)
                onUbah(namaMedan, nama)
              }}
              required
            />
          )
        }
        if (m === 'Occupation' || m === 'Clause') {
          // Nama dari master (server); hanya dibaca.
          return <Field key={m} label={label} value={form.medan[m] ?? ''} onChange={() => undefined} readOnly />
        }
        if (turunan.has(m)) {
          // Hanya dibaca dan TIDAK dikirim (`keMasukKlausul`) - aman diberi pemisah ribuan.
          return <Field key={m} label={label} value={tampilMedanKlausul(m, form.medan[m] ?? '')} onChange={() => undefined} readOnly />
        }
        return (
          <Field
            key={m}
            label={label}
            value={form.medan[m] ?? ''}
            onChange={(v) => ubahDanKonversi(m, v)}
            required={aturan.wajib.includes(m)}
          />
        )
      })}
    </div>
  )
}

/** Grid + form satu aturan (induk, atau anak satu induk). */
function GridAturan({
  tahunID,
  jenis,
  aturan,
  induk,
  kursAda,
  onShowChild,
}: {
  tahunID: string
  jenis: JenisKlausul
  aturan: AturanKlausul
  induk: string
  /** Tiket 11: false = jenis berkurs tanpa kurs berlaku — `Add` nonaktif. */
  kursAda: boolean
  onShowChild?: (k: Klausul) => void
}) {
  const [daftar, setDaftar] = useState<DaftarKlausul | null>(null)
  const [form, setForm] = useState<FormKlausul | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [info, setInfo] = useState<string | null>(null)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilKlausul(tahunID, jenis.id, induk))
    } catch (e) {
      setGalat(e)
    }
  }, [tahunID, jenis.id, induk])

  useEffect(() => {
    void muat()
  }, [muat])

  async function simpan(): Promise<void> {
    if (form === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setInfo(null)
    try {
      const h = await simpanKlausul(tahunID, keMasukKlausul(aturan, jenis.id, form, induk))
      // Simpan berhasil: form ditutup (keputusan work owner 30-09-2026).
      setForm(null)
      setInfo(h.peringatan !== '' ? h.peringatan : KLAUSUL_TCO.tersimpan)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  if (aturan.ditahan !== '') {
    return (
      <p className="polis__catatan" role="note">
        {aturan.jenis}: {aturan.ditahan}
      </p>
    )
  }

  const baris = barisSubjenis(daftar?.daftar ?? [], aturan.subjenis)
  return (
    <div className="panel">
      <h4 className="panel__title">
        {aturan.jenis}
        {aturan.subjenis !== '' ? ` — ${aturan.subjenis}` : ''}
      </h4>
      {galat !== null && <Gagal galat={galat} />}
      {info !== null && <p role="status">{info}</p>}
      {form !== null && (
        <>
          <FormMedan
            tahunID={tahunID}
            aturan={aturan}
            form={form}
            onUbah={(m, v) => {
              setForm((f) => (f === null ? f : { ...f, medan: { ...f.medan, [m]: v } }))
            }}
          />
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void simpan()}>
              {KLAUSUL_TCO.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalat(null)
              }}
            >
              {KLAUSUL_TCO.cancel}
            </button>
          </div>
        </>
      )}
      <div className="aksi-baris">
        {addTampil(aturan, daftar !== null, baris.length) && (
          <button
            type="button"
            className="btn btn--primary"
            disabled={aturan.berkurs && !kursAda}
            onClick={() => {
              setInfo(null)
              setForm(formKlausulKosong(aturan))
            }}
          >
            {KLAUSUL_TCO.add}
          </button>
        )}
        {aturan.anak && daftar !== null && (
          <span>
            {' '}
            {KLAUSUL_TCO.totalPct}: {formatNumber(daftar.totalPct, DESIMAL_TAK_DIBATASI)}
            {daftar.peringatan !== '' ? ` — ${daftar.peringatan}` : ''}
          </span>
        )}
      </div>
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && baris.length === 0 && <Kosong pesan={KLAUSUL_TCO.kosong} />}
      {baris.length > 0 && (
        <div className="tco-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {aturan.medan.map((m) => (
                  <th key={m}>{labelMedan(aturan, m)}</th>
                ))}
                <th>{KLAUSUL_TCO.formModifiedDate}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {baris.map((k) => (
                <tr key={k.id} className="inbox__baris">
                  {aturan.medan.map((m) => (
                    <td key={m}>{m === 'ReinsTypeID' ? k.reinsTypeName || k.reinsTypeId : tampilMedanKlausul(m, k.medan[m] ?? '')}</td>
                  ))}
                  <td>{k.tglUpdate}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => {
                        setInfo(null)
                        setForm(formKlausulDari(aturan, k))
                      }}>
                      {KLAUSUL_TCO.edit}
                    </button>
                    {onShowChild !== undefined && (
                      <>
                        {' '}
                        <button type="button" className="btn btn--ghost btn--sm" onClick={() => onShowChild(k)}>
                          {KLAUSUL_TCO.showChild}
                        </button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

export default function PanelJenisKlausul({
  tahunID,
  jenis,
  onTutup,
  tanpaJudul = false,
}: {
  tahunID: string
  jenis: JenisKlausul
  onTutup?: () => void
  /** true di dalam popup: judul jenis sudah di kepala popup. */
  tanpaJudul?: boolean
}) {
  const [indukTerpilih, setIndukTerpilih] = useState<Klausul | null>(null)
  const anak = aturanAnak(jenis)
  const tab = tabSubjenis(jenis)
  const [tabAktif, setTabAktif] = useState(tab[0] ?? '')
  const berkurs = jenisBerkurs(jenis)
  const [kurs, setKurs] = useState<KursTahun | null>(null)
  const [galatKurs, setGalatKurs] = useState<unknown>(null)

  useEffect(() => {
    if (!berkurs) return
    ambilKursTahun(tahunID).then(setKurs).catch(setGalatKurs)
  }, [berkurs, tahunID])
  return (
    <section className="panel">
      {!tanpaJudul && (
        <header className="inbox__kepala">
          <h3 className="panel__title">
            {jenis.id} — {jenis.descName}
          </h3>
          {onTutup !== undefined && (
            <button type="button" className="btn btn--ghost btn--sm" onClick={onTutup}>
              {KLAUSUL_TCO.tutup}
            </button>
          )}
        </header>
      )}
      {jenis.catatan !== '' && (
        <p className="polis__catatan" role="note">
          {jenis.catatan}
        </p>
      )}
      {berkurs && kurs !== null && (
        <p role="status">
          {KURS_TCO.kurs}: {formatNumber(kurs.kurs, DESIMAL_TAK_DIBATASI)} ({KURS_TCO.berlaku} {kurs.mulai} {KURS_TCO.sampai}{' '}
          {kurs.akhir})
        </p>
      )}
      {berkurs && galatKurs !== null && <Gagal galat={galatKurs} />}
      {tab.length > 0 && <StripTab tab={tab} aktif={tabAktif} onPilih={setTabAktif} />}
      {aturanInduk(jenis)
        .filter((a) => tab.length === 0 || a.subjenis === tabAktif)
        .map((a) => (
          <GridAturan
            key={`${a.jenis}/${a.subjenis}`}
            tahunID={tahunID}
            jenis={jenis}
            aturan={a}
            induk="00"
            kursAda={kurs !== null}
            onShowChild={anak !== undefined ? (k) => setIndukTerpilih(k) : undefined}
          />
        ))}
      {anak !== undefined && indukTerpilih !== null && (
        <div>
          <GridAturan
            key={indukTerpilih.reinsTypeId}
            tahunID={tahunID}
            jenis={jenis}
            aturan={anak}
            induk={indukTerpilih.reinsTypeId}
            kursAda={kurs !== null}
          />
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setIndukTerpilih(null)}>
            {KLAUSUL_TCO.closeChild}
          </button>
        </div>
      )}
    </section>
  )
}
