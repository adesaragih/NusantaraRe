// Tab Clauses kasus FIRE (tiket 47) - port `NB FacIn\Section\InputDtlClause_FacIn.xml` (layout S27 tombol Choose Clause,
// S29 grid `.ClauseList` -> `InputClauseFire_FacIn` per baris + tombol hapus) dengan gambar layar Pega DEV (work owner
// 05-10-2026):
// - Choose Clause (`SearchClauseFireSQL_PreAct` + harness `ChooseClauseFire` "Pilih Klausula Fire"): saring Language +
//   Key word, Search, grid Clause ID / Title / Description berkotak centang (multi-pilih), Submit
//   (`SearchClauseFireSQL_PostAct`): klausa terpilih yang belum ada ditambahkan ke ClauseList;
// - per baris: Clause Code, Total Argument (bila > 0), Title, Description (baca-saja), See Clause and Argument (local
//   action `InputClauseFire_ViewDtl` "Isi Klasula": isi klausa + grid argumen, Isi Argumen dapat diisi;
//   `ReplaceClauseArgumentFireAct`: isi = isi asli dengan setiap `_&<nomor>` diganti nilainya), hapus baris;
// - Save (sel 57, `SaveFacIn_Act`) menyimpan ClauseList kasus.
//
// Keputusan agent (tiket 47): C-1 isi klausa ditampilkan sebagai TEKS (bukan HTML rich text) - aman dari skrip. C-2
// penggantian argumen mengikuti urutan Pega (nomor menaik, ganti-semua harfiah; `_&1` juga mengenai awal `_&10`).
// C-3 hasil Choose Clause berhalaman 10 baris. C-4 (selesai 05-10-2026) Language = PromptList `.ClauseLanguageID`
// (kode 0 Indonesia / 1 Inggris / 2 Dual Bahasa) - yang dikirim ke pencarian adalah KODE-nya.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKlausaKasus,
  argumenKlausa,
  cariKlausa,
  simpanKlausaKasus,
  type ArgumenKlausa,
  type HasilKlausa,
  type KlausaKasus,
} from '../api'
import { BAHASA_KLAUSA_AWAL, KLAUSA as K, OPSI_BAHASA_KLAUSA, TEKS_FORM_OPPORTUNITY, TEKS_KLAUSA as T } from '../labels'
import PagerHalaman from './PagerHalaman'

/** Kode bahasa ClauseLanguage (`SearchClauseFireSQL_PostAct`: Indonesia "0", Inggris "1", lainnya "2"). */
export const kodeBahasa = (language: string) => (language === 'Indonesia' ? '0' : language === 'Inggris' ? '1' : '2')

/** Submit Choose Clause: tambahkan klausa terpilih yang kodenya belum ada (pencegah ganda PostAct langkah 3.1). */
export function tambahKlausa(daftar: KlausaKasus[], terpilih: HasilKlausa[]): KlausaKasus[] {
  const ada = new Set(daftar.map((k) => k.clauseCode))
  const baru: KlausaKasus[] = []
  for (const h of terpilih) {
    if (ada.has(h.id)) continue
    ada.add(h.id)
    baru.push({
      clauseCode: h.id,
      clauseTitle: h.title,
      clauseDescription: h.info,
      clauseLanguage: kodeBahasa(h.language),
      clauseLanguageId: h.language,
      clauseContent: h.text,
      clauseContentTemp: h.text,
      argumentCount: h.argumentCount,
      argumentList: [],
    })
  }
  return [...daftar, ...baru]
}

/** `ReplaceClauseArgumentFireAct` langkah 2-3: isi asli, lalu setiap `_&<nomor>` diganti Isi Argumen (urutan Pega). */
export function gantiArgumen(isiAsli: string, argumen: ArgumenKlausa[]): string {
  let isi = isiAsli
  for (const a of argumen) isi = isi.split(`_&${a.argumentNumber}`).join(a.argumentValue)
  return isi
}

/** Total Argument tampil bila > 0 (`.ArgumentCount > 0`, sel 5) - teks bilangan bulat, tanpa konversi angka. */
export const adaArgumen = (n: string) => /^0*[1-9][0-9]*$/.test(n.trim())

/** Baris hasil per halaman popup Choose Clause (C-3). */
export const UKURAN_HALAMAN_KLAUSA = 10

/**
 * Popup Choose Clause. `sudahAda` = kode di ClauseList kasus: tampil TERCENTANG (`SearchClauseFireSQL_PreAct` langkah 6
 * `.IsSelected = "true"` bila `.ID` sudah ada) dan dikunci - Submit Pega hanya menambah, tidak membuang.
 */
function PopupPilihKlausa({
  sudahAda,
  onTutup,
  onSubmit,
}: {
  sudahAda: ReadonlySet<string>
  onTutup: () => void
  onSubmit: (terpilih: HasilKlausa[]) => void
}) {
  const [bahasa, setBahasa] = useState(BAHASA_KLAUSA_AWAL)
  const [kata, setKata] = useState('')
  const [saring, setSaring] = useState({ bahasa: BAHASA_KLAUSA_AWAL, kata: '' })
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<{ baris: HasilKlausa[]; total: number; halaman: number; ukuran: number } | null>(
    null,
  )
  const [galat, setGalat] = useState<unknown>(null)
  // Pilihan bertahan lintas halaman (kunci = Clause ID).
  const [pilih, setPilih] = useState<Map<string, HasilKlausa>>(new Map())

  useEffect(() => {
    let batal = false
    setData(null)
    setGalat(null)
    cariKlausa(saring.bahasa, saring.kata, halaman).then(
      (h) => {
        if (!batal) setData(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [saring, halaman])

  function centang(h: HasilKlausa) {
    setPilih((m) => {
      const n = new Map(m)
      if (n.has(h.id)) n.delete(h.id)
      else n.set(h.id, h)
      return n
    })
  }

  return (
    <Modal
      judul={T.judulPilih}
      onTutup={onTutup}
      lebar
      aksi={
        <button
          type="button"
          className="btn btn--primary"
          disabled={pilih.size === 0}
          onClick={() => onSubmit([...pilih.values()])}
        >
          {K.submit.label}
        </button>
      }
    >
      <form
        className="nbf-klausa__saring"
        onSubmit={(e) => {
          e.preventDefault()
          setHalaman(1)
          setSaring({ bahasa, kata })
        }}
      >
        {/* Ganti Language = cari ulang + kosongkan centang: Title / isi klausa hasil cari bergantung bahasa (work owner
            05-10-2026: pilih Inggris tetapi isi masih TEXTINA, karena hasil lama belum dicari ulang). */}
        <Pilih
          label={K.language.label}
          value={bahasa}
          onChange={(v) => {
            setBahasa(v)
            setPilih(new Map())
            setHalaman(1)
            setSaring({ bahasa: v, kata })
          }}
          opsi={OPSI_BAHASA_KLAUSA}
        />
        <div className="field">
          <label className="field__label">{K.keyword.label}</label>
          <input className="field__input" type="search" value={kata} onChange={(e) => setKata(e.target.value)} />
        </div>
        <button type="submit" className="btn btn--primary">
          {K.cari.label}
        </button>
        {pilih.size > 0 && <span className="nbf-akum__jumlah">{T.dipilih(pilih.size)}</span>}
      </form>
      {data === null && galat === null && <Memuat />}
      <Gagal galat={galat} />
      {data !== null && data.baris.length === 0 && <Kosong pesan={T.tanpaHasil} />}
      {data !== null && data.baris.length > 0 && (
        <>
          <div className="nbf-cov-wrap">
            <table className="nbf-tabel">
              <thead>
                <tr>
                  <th scope="col" />
                  {K.kolomPilih.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {data.baris.map((h) => (
                  <tr key={h.id} className={pilih.has(h.id) || sudahAda.has(h.id) ? 'nbf-baris--terbuka' : undefined}>
                    <td>
                      <input
                        type="checkbox"
                        aria-label={h.id}
                        title={sudahAda.has(h.id) ? T.sudahAda : undefined}
                        checked={pilih.has(h.id) || sudahAda.has(h.id)}
                        disabled={sudahAda.has(h.id)}
                        onChange={() => centang(h)}
                      />
                    </td>
                    <td>
                      <code className="nbf-akum__kode">{h.id}</code>
                    </td>
                    <td>{h.title}</td>
                    <td>{h.info}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {data.total > data.ukuran && (
            <PagerHalaman total={data.total} halaman={data.halaman} ukuran={data.ukuran} onPindah={setHalaman} />
          )}
        </>
      )}
    </Modal>
  )
}

function PopupArgumen({
  klausa,
  onTutup,
  onTerapkan,
}: {
  klausa: KlausaKasus
  onTutup: () => void
  onTerapkan: (k: KlausaKasus) => void
}) {
  const [argumen, setArgumen] = useState<ArgumenKlausa[] | null>(
    klausa.argumentList.length > 0 ? klausa.argumentList : null,
  )
  const [galat, setGalat] = useState<unknown>(null)

  // ViewClauseArgFireSQL: argumen dimuat dari master hanya bila ArgumentList baris masih kosong.
  useEffect(() => {
    if (argumen !== null) return
    let batal = false
    argumenKlausa(klausa.clauseCode).then(
      (h) => {
        if (!batal) setArgumen(h.baris)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [argumen, klausa.clauseCode])

  const pratinjau = argumen ? gantiArgumen(klausa.clauseContentTemp, argumen) : klausa.clauseContent
  return (
    <Modal
      judul={`${K.judulArgumen} · ${klausa.clauseCode}`}
      onTutup={onTutup}
      lebar
      aksi={
        <button
          type="button"
          className="btn btn--primary"
          disabled={argumen === null}
          onClick={() =>
            argumen &&
            onTerapkan({
              ...klausa,
              argumentList: argumen,
              clauseContent: gantiArgumen(klausa.clauseContentTemp, argumen),
            })
          }
        >
          {K.submit.label}
        </button>
      }
    >
      <div className="nbf-klausa__isi">{pratinjau}</div>
      <Gagal galat={galat} />
      {argumen === null && galat === null && <Memuat />}
      {/* Tanpa argumen (M_ARGCLAUSEFIRE tidak ada di DEV, K47-6): popup hanya menampilkan isi klausa, seperti Pega. */}
      {argumen !== null && argumen.length > 0 && (
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {K.kolomArgumen.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {argumen.map((a, i) => (
                <tr key={a.argumentNumber}>
                  <td>{a.argumentNumber}</td>
                  <td>{a.argumentDescription}</td>
                  <td>
                    <input
                      className="field__input"
                      aria-label={`${K.kolomArgumen[2]} ${a.argumentNumber}`}
                      value={a.argumentValue}
                      onChange={(e) =>
                        setArgumen((d) =>
                          d ? d.map((x, k) => (k === i ? { ...x, argumentValue: e.target.value } : x)) : d,
                        )
                      }
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Modal>
  )
}

export default function TabClauses({ caseId }: { caseId: string }) {
  const [daftar, setDaftar] = useState<KlausaKasus[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pilih, setPilih] = useState(false)
  const [argumen, setArgumen] = useState<number | null>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)

  useEffect(() => {
    let batal = false
    ambilKlausaKasus(caseId).then(
      (h) => {
        if (!batal) setDaftar(h.baris)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [caseId])

  function ubah(baru: KlausaKasus[]) {
    setDaftar(baru)
    setTersimpan(false)
  }

  async function simpan() {
    if (daftar === null) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanKlausaKasus(caseId, daftar)
      setDaftar(h.baris)
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (daftar === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  return (
    <div className="nbf-objek nbf-klausa">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      <div className="nbf-klausa__atas">
        <button type="button" className="btn btn--primary" onClick={() => setPilih(true)}>
          {K.pilihKlausa.label}
        </button>
      </div>
      {daftar.length === 0 && <Kosong pesan={T.kosong} />}
      <div className="nbf-klausa__daftar">
        {daftar.map((k, i) => (
          <article key={k.clauseCode} className="nbf-klausa__kartu">
            <div className="nbf-klausa__kepala">
              <div className="nbf-klausa__identitas">
                <span className="nbf-klausa__label">{K.clauseCode.label}</span>
                <code className="nbf-cov-kode">{k.clauseCode}</code>
                {k.clauseLanguageId !== '' && <span className="nbf-akum__jumlah">{k.clauseLanguageId}</span>}
                {adaArgumen(k.argumentCount) && (
                  <span className="nbf-akum__jumlah">
                    {K.totalArgumen.label}: {k.argumentCount}
                  </span>
                )}
              </div>
              <div className="nbf-klausa__aksi">
                <button type="button" className="btn btn--ghost btn--sm" onClick={() => setArgumen(i)}>
                  {K.lihatArgumen.label}
                </button>
                <button
                  type="button"
                  className="btn btn--ghost btn--sm nbf-klausa__hapus"
                  aria-label={`${T.hapus} ${k.clauseCode}`}
                  onClick={() => ubah(daftar.filter((_, x) => x !== i))}
                >
                  {T.hapus}
                </button>
              </div>
            </div>
            <div className="nbf-klausa__judul">
              <span className="nbf-klausa__label">{K.title.label}</span>
              <strong>{k.clauseTitle}</strong>
            </div>
            <div>
              <span className="nbf-klausa__label">{K.description.label}</span>
              <p className="nbf-klausa__deskripsi">{k.clauseDescription}</p>
            </div>
          </article>
        ))}
      </div>
      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan}>
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : K.simpan}
        </button>
      </div>
      {pilih && (
        <PopupPilihKlausa
          sudahAda={new Set(daftar.map((k) => k.clauseCode))}
          onTutup={() => setPilih(false)}
          onSubmit={(terpilih) => {
            ubah(tambahKlausa(daftar, terpilih))
            setPilih(false)
          }}
        />
      )}
      {argumen !== null && daftar[argumen] && (
        <PopupArgumen
          klausa={daftar[argumen]}
          onTutup={() => setArgumen(null)}
          onTerapkan={(baru) => {
            ubah(daftar.map((x, k) => (k === argumen ? baru : x)))
            setArgumen(null)
          }}
        />
      )}
    </div>
  )
}
