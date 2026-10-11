// Tab Object FIRE PER HALAMAN untuk case besar (tiket 75 F3, keputusan work owner 10-10-2026 K75-3: "Ya, tampil per halaman dari
// server") - dipakai bila `jumlahLokasi` > `AMBANG_HALAMAN_OBJEK`; case kecil tetap `TabObject` (alur yang sudah berjalan tidak
// diubah). Kabel `GET …/objek/halaman` + `PATCH …/objek` (sesi c3).
//
// Grid, kolom, sub-tab detail, dan aturan isian SAMA dengan `TabObject` (port `ObjectList.xml` / `Property_FlowAction`); bedanya:
// - grid dibaca per halaman (10 baris, seperti Pega) + kotak cari di server;
// - perubahan dicatat per `id` (tetap terbawa saat pindah halaman), lokasi baru di daftar terpisah di atas grid, hapus ditandai
//   dulu (Batal hapus) - Save mengirim HANYA perubahan (`PATCH`), lalu halaman dibaca ulang;
// - batas per Save mengikuti backend (ubah + hapus ≤ 500, tambah ≤ 500);
// - 409 dari server tampil apa adanya (lokasi berubah sejak dimuat, lokasi ber-spreading K75-4, unggahan CSV berjalan).

import { useEffect, useRef, useState, type ReactNode } from 'react'

import { BelumTersedia, Gagal, Halaman, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilObjekHalaman, perbaruiObjek, type ObjekBerId, type ObjekFire } from '../api'
import { GRID_OBJEK, SIMPAN_OBJEK, SUBTAB_OBJEK, TEKS_FORM_OPPORTUNITY, TEKS_INWARD, TEKS_OBJEK, TEKS_OBJEK_HALAMAN as H } from '../labels'
import SubTabItem, { adaGalatItem } from './SubTabItem'
import SubTabFEA, { adaGalatFEA } from './SubTabFEA'
import { SubTabKerugian, SubTabKlaimInternal, adaGalatKerugian } from './SubTabKerugian'
import SubTabOkupasi from './SubTabOkupasi'
import SubTabSekitar, { adaJarakMinus } from './SubTabSekitar'
import { lantaiMinus, ObjectAddress, objekBaru, rapikanObjek } from './TabObject'

type SubTab = (typeof SUBTAB_OBJEK)[number]

const JEDA_CARI_MS = 400

/** Batas per Save (backend `PATCH …/objek`). */
export const BATAS_UBAH_PATCH = 500
export const BATAS_TAMBAH_PATCH = 500

/** Satu objek bermasalah untuk Save (aturan `TabObject.simpan`). */
export const objekSalah = (o: ObjekFire) =>
  o.objectType === '' ||
  lantaiMinus(o.numberOfFloor) ||
  adaJarakMinus(o.surroundingRisk) ||
  adaGalatItem(o.items) ||
  adaGalatFEA(o.fea) ||
  adaGalatKerugian(o.lossRecords)

/** Sub-tab yang memuat pesan galat objek (urutan `TabObject`). */
const subTabSalah = (o: ObjekFire): SubTab =>
  o.objectType === '' || lantaiMinus(o.numberOfFloor)
    ? SUBTAB_OBJEK[0]
    : adaJarakMinus(o.surroundingRisk)
      ? SUBTAB_OBJEK[1]
      : adaGalatItem(o.items)
        ? SUBTAB_OBJEK[2]
        : adaGalatFEA(o.fea)
          ? SUBTAB_OBJEK[4]
          : SUBTAB_OBJEK[5]

/** Badan PATCH dari catatan perubahan layar (lokasi yang ditandai hapus tidak ikut "ubah"). */
export function badanPatch(ubahan: Map<string, ObjekBerId>, tambahan: ObjekFire[], hapusan: Set<string>) {
  return { ubah: [...ubahan.values()].filter((o) => !hapusan.has(o.id)), tambah: tambahan, hapus: [...hapusan] }
}

export default function TabObjectHalaman({ caseId, insuredName = '' }: { caseId: string; insuredName?: string }) {
  const [kotak, setKotak] = useState('')
  const [cari, setCari] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [muat, setMuat] = useState(0)
  const [data, setData] = useState<{ total: number; lokasi: ObjekBerId[] } | null>(null)
  const [ubahan, setUbahan] = useState<Map<string, ObjekBerId>>(new Map())
  const [tambahan, setTambahan] = useState<{ kunci: number; data: ObjekFire }[]>([])
  const [hapusan, setHapusan] = useState<Set<string>>(new Set())
  const [terbuka, setTerbuka] = useState<Record<string, SubTab>>({})
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState('')
  const [menyimpan, setMenyimpan] = useState(false)
  const [cobaSimpan, setCobaSimpan] = useState(false)
  const kunciBaru = useRef(1)

  useEffect(() => {
    const t = window.setTimeout(() => {
      setCari(kotak.trim())
      setHalaman(1)
    }, JEDA_CARI_MS)
    return () => window.clearTimeout(t)
  }, [kotak])

  useEffect(() => {
    let batal = false
    setData(null)
    ambilObjekHalaman(caseId, halaman, GRID_OBJEK.ukuran, cari).then(
      (h) => {
        if (!batal) setData({ total: h.total, lokasi: h.lokasi.map((o) => ({ ...rapikanObjek(o), id: o.id })) })
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId, halaman, cari, muat])

  const jumlahPerubahan = ubahan.size + tambahan.length + hapusan.size
  const lewatBatas = ubahan.size + hapusan.size > BATAS_UBAH_PATCH || tambahan.length > BATAS_TAMBAH_PATCH
  const balik = (k: string) => setTerbuka((t) => (t[k] === undefined ? { ...t, [k]: SUBTAB_OBJEK[0] } : Object.fromEntries(Object.entries(t).filter(([x]) => x !== k))))

  function ubahLama(o: ObjekBerId) {
    setUbahan((m) => new Map(m).set(o.id, o))
    setPesan('')
  }

  function ubahBaru(kunci: number, d: ObjekFire) {
    setTambahan((t) => t.map((x) => (x.kunci === kunci ? { ...x, data: d } : x)))
    setPesan('')
  }

  function tambah() {
    const kunci = kunciBaru.current++
    setTambahan((t) => [...t, { kunci, data: objekBaru(String((data?.total ?? 0) + t.length + 1)) }])
    setTerbuka((t) => ({ ...t, [`baru:${kunci}`]: SUBTAB_OBJEK[0] }))
    setPesan('')
  }

  function alihHapus(id: string) {
    setHapusan((s) => {
      const b = new Set(s)
      if (b.has(id)) b.delete(id)
      else b.add(id)
      return b
    })
    setPesan('')
  }

  async function simpan() {
    setCobaSimpan(true)
    const salahLama = [...ubahan.values()].filter((o) => !hapusan.has(o.id) && objekSalah(o))
    const salahBaru = tambahan.filter((x) => objekSalah(x.data))
    if (salahLama.length > 0 || salahBaru.length > 0) {
      setTerbuka((t) => ({
        ...t,
        ...Object.fromEntries(salahLama.map((o) => [`id:${o.id}`, subTabSalah(o)])),
        ...Object.fromEntries(salahBaru.map((x) => [`baru:${x.kunci}`, subTabSalah(x.data)])),
      }))
      return
    }
    if (lewatBatas) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await perbaruiObjek(
        caseId,
        badanPatch(
          ubahan,
          tambahan.map((x) => x.data),
          hapusan,
        ),
      )
      setUbahan(new Map())
      setTambahan([])
      setHapusan(new Set())
      setTerbuka({})
      setCobaSimpan(false)
      setPesan(H.tersimpan(h.jumlahLokasi))
      setMuat((n) => n + 1)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  // Fungsi render biasa (bukan komponen di dalam komponen - isian kehilangan fokus tiap render).
  function detail(k: string, o: ObjekFire, ubah: (d: ObjekFire) => void) {
    const sub = terbuka[k] ?? SUBTAB_OBJEK[0]
    return (
      <div className="nbf-lapis nbf-lapis--objek">
        <StripTab tab={SUBTAB_OBJEK} aktif={sub} onPilih={(t) => setTerbuka((s) => ({ ...s, [k]: t }))} />
        {sub === SUBTAB_OBJEK[0] ? (
          <ObjectAddress o={o} ubah={ubah} galatType={cobaSimpan && o.objectType === ''} />
        ) : sub === SUBTAB_OBJEK[1] ? (
          <SubTabSekitar o={o} ubah={ubah} />
        ) : sub === SUBTAB_OBJEK[2] ? (
          <SubTabItem items={o.items} ubah={(items) => ubah({ ...o, items })} tandaiWajib={cobaSimpan} />
        ) : sub === SUBTAB_OBJEK[5] ? (
          <SubTabKerugian rows={o.lossRecords} lossRatio={o.lossRatio} insuredName={insuredName} ubah={(lossRecords) => ubah({ ...o, lossRecords })} tandaiWajib={cobaSimpan} />
        ) : sub === SUBTAB_OBJEK[6] ? (
          <SubTabKlaimInternal rows={o.internalLossRecords} />
        ) : sub === SUBTAB_OBJEK[4] ? (
          <SubTabFEA fea={o.fea} ubah={(fea) => ubah({ ...o, fea })} />
        ) : sub === SUBTAB_OBJEK[3] ? (
          <SubTabOkupasi caseId={caseId} occupations={o.occupations} ubah={(occupations) => ubah({ ...o, occupations })} />
        ) : (
          <BelumTersedia apa={`${TEKS_INWARD.isiTab} ${sub}`} />
        )}
      </div>
    )
  }

  function baris(k: string, o: ObjekFire, aksi: ReactNode, kelas?: string) {
    const buka = terbuka[k] !== undefined
    return (
      <tr key={k} className={[buka ? 'nbf-baris--terbuka' : '', kelas ?? ''].join(' ').trim() || undefined}>
        <td>
          <button
            type="button"
            className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
            aria-label={TEKS_OBJEK.bukaBaris}
            aria-expanded={buka}
            onClick={() => balik(k)}
          >
            ▸
          </button>
        </td>
        <td>
          <input type="checkbox" checked={o.isTopRisk} disabled aria-label={GRID_OBJEK.kolom[0].label} />
        </td>
        <td>{o.objectNo}</td>
        <td>{o.objectName}</td>
        <td>{o.riskLocation}</td>
        <td className="table__actions">{aksi}</td>
      </tr>
    )
  }

  const kepala = (
    <thead>
      <tr>
        <th scope="col" />
        {GRID_OBJEK.kolom.map((k) => (
          <th key={k.sel} scope="col">
            {k.label}
          </th>
        ))}
        <th scope="col" className="table__actions">
          <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={tambah}>
            {GRID_OBJEK.tambah}
          </button>
        </th>
      </tr>
    </thead>
  )

  return (
    <div className="nbf-objek nbf-obj-tab">
      <Gagal galat={galat} />
      {pesan !== '' && <div className="alert alert--ok">{pesan}</div>}
      <div className="toolbar">
        <input className="field__input kelola-user__cari" type="search" aria-label={H.cari} placeholder={H.cari} value={kotak} onChange={(e) => setKotak(e.target.value)} />
        {data !== null && <span className="muted">{H.jumlah(data.total)}</span>}
        {jumlahPerubahan > 0 && <span className="muted">{H.belumDisimpan(jumlahPerubahan)}</span>}
      </div>

      {tambahan.length > 0 && (
        <section className="nbf-lapis">
          <h5 className="nbf-objek__judul">{H.baru}</h5>
          <div className="nbf-cov-wrap">
            <table className="nbf-tabel">
              {kepala}
              <tbody>
                {tambahan.map((x) => {
                  const k = `baru:${x.kunci}`
                  return [
                    baris(
                      k,
                      x.data,
                      <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => setTambahan((t) => t.filter((y) => y.kunci !== x.kunci))}>
                        {GRID_OBJEK.hapus}
                      </button>,
                    ),
                    terbuka[k] !== undefined && (
                      <tr key={`${k}-detail`} className="nbf-objek__detail">
                        <td colSpan={6}>{detail(k, x.data, (d) => ubahBaru(x.kunci, d))}</td>
                      </tr>
                    ),
                  ]
                })}
              </tbody>
            </table>
          </div>
        </section>
      )}

      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          {kepala}
          <tbody>
            {data === null && (
              <tr>
                <td colSpan={6}>
                  <Memuat />
                </td>
              </tr>
            )}
            {data !== null && data.lokasi.length === 0 && (
              <tr>
                <td colSpan={6}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {data?.lokasi.map((asal) => {
              const o = ubahan.get(asal.id) ?? asal
              const k = `id:${asal.id}`
              const dihapus = hapusan.has(asal.id)
              return [
                baris(
                  k,
                  o,
                  <>
                    {dihapus && <span className="muted">{H.dihapus} </span>}
                    <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => alihHapus(asal.id)}>
                      {dihapus ? H.batalHapus : GRID_OBJEK.hapus}
                    </button>
                  </>,
                  dihapus ? 'nbf-objek__dihapus' : undefined,
                ),
                terbuka[k] !== undefined && !dihapus && (
                  <tr key={`${k}-detail`} className="nbf-objek__detail">
                    <td colSpan={6}>{detail(k, o, (d) => ubahLama({ ...d, id: asal.id }))}</td>
                  </tr>
                ),
              ]
            })}
          </tbody>
        </table>
      </div>
      {data !== null && data.total > GRID_OBJEK.ukuran && <Halaman halaman={halaman} ukuran={GRID_OBJEK.ukuran} total={data.total} onPindah={setHalaman} />}
      {lewatBatas && <div className="field__error">{H.terlaluBanyak}</div>}
      <div className="nbf-objek__kaki">
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary"
          onClick={() => void simpan()}
          disabled={menyimpan || jumlahPerubahan === 0 || lewatBatas}
        >
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_OBJEK.label}
        </button>
      </div>
    </div>
  )
}
