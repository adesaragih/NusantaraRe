import { useState, type FormEvent } from 'react'

import { DETAIL } from '../assets/labels'
import { BelumTersedia } from '../components/ui/dasar'
import { PanelTotalPeserta } from '../components/PanelTotalPeserta'

import {
  ambilKlaimLife,
  kodeStatusGalat,
  pesanGalat,
  periksaBolehTutup,
  ubahTanggalKejadian,
  type HasilPeriksaTutup,
  tampilUang,
  tolakBarisAdjustment,
  serahkanKeKomite,
  tambahPutaran,
  simpanAdjustment,
  bolehSimpanAdjustment,
  bolehPutaranBaru,
  bolehSerahkanDiLayar,
  dampakHapusKlaim,
  hapusKlaim,
  STATUS_OUTSTANDING,
  type DampakHapus,
  type Klaim,
} from '../services/api'

// ============================================================================
// pages/KlaimLife.tsx — halaman "buka satu klaim Life" (tiket 01 AC-2).
//
// Cara membaca berkas ini bila baru mengenal React:
//   - Komponen = fungsi yang MENGEMBALIKAN tampilan. Bagian yang mirip HTML di
//     dalam `return (...)` disebut JSX; karena berkas ini TypeScript, namanya TSX.
//   - `useState` = kotak penyimpan nilai. Setiap kali nilainya diganti lewat
//     fungsi set-nya (misalnya `setKlaim`), React menggambar ulang tampilan.
//   - `<string>` di belakang useState memberi tahu TypeScript isi kotak itu apa.
//   - `{ ... }` di dalam JSX = "sisipkan nilai JavaScript di sini".
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika dan tidak pernah dibungkus
// Number(). Ia datang sebagai teks desimal dan ditampilkan apa adanya
// (ADR-U-0003, ADR-U-0016). Penjelasannya di services/api.ts bab 3.
// ============================================================================
export default function KlaimLife() {
  // Empat kotak keadaan halaman ini.
  const [id, setId] = useState<string>('') // yang sedang diketik pengguna
  const [klaim, setKlaim] = useState<Klaim | null>(null) // klaim yang berhasil dibaca
  const [galat, setGalat] = useState<string | null>(null) // pesan bila gagal
  const [memuat, setMemuat] = useState<boolean>(false) // sedang menunggu server?
  const [menolak, setMenolak] = useState<string | null>(null) // baris yang sedang ditolak
  const [menyerahkan, setMenyerahkan] = useState<string | null>(null) // baris ke Komite
  const [memutar, setMemutar] = useState<string | null>(null) // peserta yang dibuka putarannya
  const [mengaksep, setMengaksep] = useState<string | null>(null) // peserta yang diaksep
  // Galat tanggal kejadian PER PESERTA, bukan satu untuk seluruh layar:
  // `ValidasiDOL_Act` menempelkan pesannya pada `.DATE_OF_LOSS` peserta,
  // dan satu kotak galat bersama akan menuding peserta yang salah.
  const [galatTanggal, setGalatTanggal] = useState<Record<string, string>>({})

  // Peserta yang permintaannya sedang terbang. ⛔ Bukan satu penanda untuk
  // seluruh layar: dua peserta boleh diubah berbarengan, dan satu penanda
  // bersama akan mengunci kotak peserta lain tanpa sebab yang terlihat.
  const [tanggalSibuk, setTanggalSibuk] = useState<Record<string, boolean>>({})

  // Hasil gerbang tutup. null = belum pernah diperiksa.
  const [gerbangTutup, setGerbangTutup] = useState<HasilPeriksaTutup | null>(null)
  const [tutupSibuk, setTutupSibuk] = useState(false)

  /** Memeriksa kesiapan tutup; TIDAK menutup apa pun. */
  async function periksaTutup(): Promise<void> {
    if (klaim === null || tutupSibuk) return
    setTutupSibuk(true)
    try {
      setGerbangTutup(await periksaBolehTutup(klaim.id))
    } catch (e) {
      setGerbangTutup({
        boleh: false,
        penghalang: [
          {
            urutan: 0,
            nomorSertifikat: '',
            pesan: pesanGalat(e) ?? 'Kesiapan tutup tidak dapat diperiksa.',
          },
        ],
      })
    } finally {
      setTutupSibuk(false)
    }
  }

  /** Mengirim tanggal kejadian satu peserta; validasinya di server. */
  async function ubahTanggal(pesertaID: string, nilai: string): Promise<void> {
    // Galat lama dibersihkan LEBIH DULU, sebelum keluar karena kotak
    // dikosongkan. Kalau tidak, pesan penolakan tanggal sebelumnya menempel
    // di layar untuk kotak yang sudah tidak berisi apa-apa.
    setGalatTanggal((lama) => {
      const { [pesertaID]: _dibuang, ...sisa } = lama
      return sisa
    })
    // Kosong tidak dikirim: backend menjawab 400 "tanggal kejadian belum
    // diisi", dan menyuruh orang menunggu perjalanan bolak-balik hanya untuk
    // dimarahi soal kotak yang baru saja ia kosongkan bukan pertolongan.
    if (klaim === null || nilai === '') return
    // ⛔ Menolak permintaan kedua selagi yang pertama terbang. Tanpa ini,
    // dua pengubahan beruntun berlomba: yang kedua tiba lebih dulu, lalu
    // `ambilKlaimLife` milik yang PERTAMA mendarat dan menimpa layar dengan
    // potret yang sudah kedaluwarsa - tanggal yang baru saja disimpan tampak
    // hilang, dan pemakai mengetiknya lagi.
    if (tanggalSibuk[pesertaID] === true) return
    setTanggalSibuk((lama) => ({ ...lama, [pesertaID]: true }))
    try {
      await ubahTanggalKejadian(klaim.id, pesertaID, nilai)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (e) {
      // Kalimat server diteruskan apa adanya: 422 menyebut jendela valuasi
      // mana yang dilanggar, dan kalimat itu yang berguna.
      setGalatTanggal((lama) => ({
        ...lama,
        [pesertaID]: pesanGalat(e) ?? 'Tanggal kejadian ditolak.',
      }))
    } finally {
      setTanggalSibuk((lama) => ({ ...lama, [pesertaID]: false }))
    }
  }
  const [dampak, setDampak] = useState<DampakHapus | null>(null) // isi popup konfirmasi
  const [menghapus, setMenghapus] = useState<boolean>(false) // permintaan hapus berjalan

  // Membuka popup: MEMBACA dampak, tidak menghapus apa pun.
  async function bukaKonfirmasiHapus() {
    if (!klaim) return
    setGalat(null)
    try {
      setDampak(await dampakHapusKlaim(klaim.id))
    } catch {
      setGalat('Gagal menghitung dampak penghapusan.')
    }
  }

  // ⛔ Batal hanya menutup popup. Tidak ada yang perlu dibatalkan, sebab
  // membuka popup tidak menulis apa pun.
  function batalHapus() {
    setDampak(null)
    setGalat(null)
  }

  async function lanjutkanHapus() {
    if (!klaim) return
    setGalat(null)
    setMenghapus(true)
    try {
      await hapusKlaim(klaim.id)
      setKlaim(null)
      setDampak(null)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 501
          ? 'Penghapusan klaim adalah penanda, bukan hapus fisik (ADR-U-0031). Kolom penandanya belum diputuskan work owner.'
          : kode === 409
            ? 'Klaim sudah diserahkan ke Komite dan tidak dapat dihapus.'
            : kode === 403
              ? 'Hanya ReasLifeAdmin yang dapat menghapus klaim.'
              : 'Gagal menghapus klaim.',
      )
      setDampak(null)
    } finally {
      setMenghapus(false)
    }
  }

  // Menolak satu baris, lalu MEMBACA ULANG klaimnya dari server.
  //
  // Dibaca ulang, bukan diubah di layar: status klaim adalah TURUNAN dari
  // seluruh barisnya, dan menebaknya di sisi klien berarti dua sumber
  // kebenaran yang dapat berbeda.
  async function aksep(pesertaID: string) {
    if (!klaim) return
    setGalat(null)
    setMengaksep(pesertaID)
    try {
      const nomor = await simpanAdjustment(klaim.id, pesertaID)
      setKlaim(await ambilKlaimLife(klaim.id))
      setGalat(`Baris diaksep dengan nomor ${nomor}.`)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya pemegang tahap klaim ini yang dapat mengaksep barisnya.'
          : kode === 409
            ? (pesanGalat(err) ?? 'Baris tidak dalam keadaan yang dapat diaksep.')
            : kode === 501
              ? (pesanGalat(err) ?? 'Akseptasi belum dapat disimpan.')
              : 'Gagal mengaksep baris.',
      )
    } finally {
      setMengaksep(null)
    }
  }

  async function putaranBaru(pesertaID: string) {
    if (!klaim) return
    setGalat(null)
    setMemutar(pesertaID)
    try {
      await tambahPutaran(klaim.id, pesertaID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya ReasLifeSPV yang dapat membuka putaran berikutnya.'
          : kode === 409
            ? 'Putaran berikutnya hanya lahir sesudah baris terakhir ditolak.'
            : kode === 501
              ? 'Jejak audit belum dapat direkam; tempatnya belum diputuskan.'
              : 'Gagal membuka putaran berikutnya.',
      )
    } finally {
      setMemutar(null)
    }
  }

  async function serahkan(pesertaID: string, adjID: string) {
    if (!klaim) return
    setGalat(null)
    setMenyerahkan(adjID)
    try {
      await serahkanKeKomite(klaim.id, pesertaID, adjID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite.'
          : kode === 409
            ? 'Baris sudah diserahkan, atau bukan lagi Outstanding.'
            : kode === 422
              ? // ⛔ Pesan 422 dicetak APA ADANYA. Salah satunya berbunyi
                // "Name of bank cannot be empty" - kalimat sistem lama, yang
                // dikenali pengguna lama. Menggantinya dengan kalimat kita
                // sendiri memutus pengenalan itu.
                (pesanGalat(err) ?? 'Penyerahan ditolak.')
              : kode === 501
                ? 'Penyerahan belum dapat disimpan; tempatnya belum diputuskan.'
                : 'Gagal menyerahkan baris ke Komite.',
      )
    } finally {
      setMenyerahkan(null)
    }
  }

  async function tolak(adjID: string) {
    if (!klaim) return
    setGalat(null)
    setMenolak(adjID)
    try {
      await tolakBarisAdjustment(klaim.id, adjID)
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        kode === 403
          ? 'Hanya ReasLifeAdmin yang dapat menolak baris.'
          : kode === 409
            ? 'Baris sudah diputus dan tidak dapat ditolak lagi.'
            : kode === 422
              ? 'Klaim belum bernomor.'
              : kode === 501
                ? 'Jejak audit belum dapat direkam; tempatnya belum diputuskan.'
                : 'Gagal menolak baris.',
      )
    } finally {
      setMenolak(null)
    }
  }

  // Dipanggil saat tombol "Buka" ditekan (form dikirim).
  async function cari(e: FormEvent<HTMLFormElement>) {
    e.preventDefault() // jangan biarkan browser memuat ulang halaman
    setGalat(null)
    setKlaim(null)
    setMemuat(true)
    try {
      setKlaim(await ambilKlaimLife(id.trim()))
    } catch (err: unknown) {
      setGalat(kodeStatusGalat(err) === 404 ? 'Klaim tidak ada.' : 'Gagal membaca klaim.')
    } finally {
      setMemuat(false) // apa pun hasilnya, berhenti "memuat"
    }
  }

  return (
    <section>
      <h2>Klaim Life</h2>

      <form onSubmit={cari}>
        <label htmlFor="idKlaim">Pengenal klaim</label>{' '}
        <input
          id="idKlaim"
          value={id}
          onChange={(e) => setId(e.target.value)}
          placeholder="mis. UJI-KLAIM-1"
        />{' '}
        <button type="submit" disabled={memuat || id.trim() === ''}>
          {memuat ? 'Memuat…' : 'Buka'}
        </button>
      </form>

      {/* `galat && (...)` = tampilkan hanya bila galat terisi. */}
      {galat && <p role="alert">{galat}</p>}

      {klaim && (
        <article>
          <h3>{klaim.nomorKlaim || klaim.id}</h3>
          <dl>
            <dt>Nomor polis</dt>
            <dd>{klaim.nomorPolis || '—'}</dd>
            <dt>Nama bisnis</dt>
            <dd>{klaim.namaBisnis || '—'}</dd>
            <dt>Jumlah baris</dt>
            <dd>{klaim.cacahBaris}</dd>
            {/* Status klaim adalah TURUNAN dari baris-barisnya, bukan kolom
                tersimpan. Yang ditampilkan kata, bukan angka - dan bukan pula
                terjemahan `kodeStatus`, yang hanya cerminan baris terakhir. */}
            <dt>Status klaim</dt>
            <dd>{klaim.statusTurunan}</dd>
          </dl>

          {/* `.map` = ulangi blok di bawah untuk SETIAP peserta.
              `key` wajib ada supaya React tahu baris mana yang berubah. */}
          <p>
            <button type="button" onClick={() => void bukaKonfirmasiHapus()}>
              Hapus klaim…
            </button>
          </p>

          {/* Popup konfirmasi: angka dulu, keputusan kemudian. Tiap jenis
              disebut sendiri - satu angka total menyembunyikan tingkat mana
              yang ternyata lebih besar dari dugaan. */}
          {dampak && (
            <aside>
              <h4>Hapus klaim {klaim.nomorKlaim || klaim.id}?</h4>
              <p>Yang akan ikut terhapus:</p>
              <ul>
                <li>Header klaim: {dampak.header}</li>
                <li>Peserta: {dampak.peserta}</li>
                <li>Baris adjustment: {dampak.adjustment}</li>
                <li>Spreading: {dampak.spreading}</li>
                <li>Spreading retro: {dampak.spreadingRetro}</li>
                <li>Dokumen: {dampak.dokumen}</li>
                <li>Baris work: {dampak.barisWork}</li>
              </ul>
              <p>Total {dampak.total} baris.</p>
              {/* ⛔ Baris datar warisan BUKAN bagian daftar di atas dan tidak
                  dijumlahkan ke totalnya. Nasibnya saat klaim dihapus belum
                  diputuskan work owner; menaruhnya di bawah judul "yang akan
                  ikut terhapus" berarti menjawab pertanyaan itu diam-diam. */}
              <p>
                ⚠️ Selain itu terdapat {dampak.barisDatarWarisan} baris datar warisan
                ber-CASEID sama. <strong>Nasibnya belum diputuskan</strong>: apakah ikut
                terhapus atau ditinggal karena hilir sudah membacanya.
              </p>
              <button
                type="button"
                disabled={menghapus}
                onClick={() => void lanjutkanHapus()}
              >
                {menghapus ? 'Menghapus…' : 'Ya, hapus'}
              </button>{' '}
              <button type="button" onClick={batalHapus}>
                Batal
              </button>
            </aside>
          )}

          {klaim.peserta.map((p) => (
            <section key={p.id}>
              <h4>
                Peserta {p.nomorSertifikat || p.id}{' '}
                <small>({p.baris.length} baris)</small>
              </h4>

              {/* Tombol `Edit Date` b14115 -> `pyLocalAction
                  ShowEditClaimLife` b14144 -> `EditDateClaimLife_Section`.
                  Validasinya `ValidasiDOL_Act`, dipanggil section yang sama
                  di b11177 dan b11298 - dan ia berjalan di SERVER, bukan di
                  sini: jendela valuasi ada di baris polis, bukan di layar.

                  ⛔ Tanggal kejadian milik PESERTA, bukan klaim.
                  `ValidasiDOL_Act` berkelas Int-LIFE_PREMIUM_DETAIL dan
                  menempelkan galatnya pada `.DATE_OF_LOSS` peserta - jadi
                  kontrolnya pun berdiri per peserta. */}
              <p>
                <label>
                  {DETAIL.ubahTanggal}{' '}
                  <input
                    type="date"
                    // ⛔ `key` memaksa kotaknya LAHIR ULANG ketika tanggal
                    // yang tersimpan berubah. Tanpa itu `defaultValue` hanya
                    // dibaca sekali, sehingga sesudah penyimpanan berhasil
                    // kotaknya tetap memperlihatkan ketikan lama - dan
                    // pemakai tidak punya cara tahu mana yang tersimpan.
                    key={`${p.id}:${p.tanggalKejadian}`}
                    defaultValue={p.tanggalKejadian}
                    disabled={tanggalSibuk[p.id] === true}
                    onChange={(e) => {
                      void ubahTanggal(p.id, e.target.value)
                    }}
                  />
                </label>
                {galatTanggal[p.id] !== undefined && (
                  <span role="alert"> {galatTanggal[p.id]}</span>
                )}
              </p>

              {/* ⛔ Kelima total ini milik PESERTA, bukan klaim.
                  `ClaimLifeDetailGCNM.xml` berkelas
                  `Int-LIFE_PREMIUM_DETAIL` (b84) dan medannya terikat
                  properti berawalan TITIK pada halaman itu
                  (`.TotalShareRNM` b20921 dst) - titik berarti "halaman
                  yang sedang berjalan", dan halaman itu peserta.

                  ⛔ Ronde pertama menaruhnya di tingkat klaim berjudul
                  "Total klaim". Keliru, dan sebabnya sama dengan butir av:
                  labelnya dibaca, IKATANNYA tidak.

                  Nilainya tetap belum ada - rule penghitungnya
                  (`CheckTotalAdjustmentClaim`) nol berkasnya di ekspor,
                  dan menebak angka uang melanggar ADR-U-0003. OQ-H. */}
              <PanelTotalPeserta />

              {/* `Find Disease` b5061 `pxButton` -> `showHarness` b5071
                  `Diagnose_Harness`. Popup diagnosis itu milik kelompok
                  MEDIS: sumbernya `DISEASE_LIFE` yang 97.586 baris, dan
                  pencariannya harus berbatas. Tombolnya DINYATAKAN di sini
                  supaya layarnya tidak tampak lengkap padahal tidak. */}
              <BelumTersedia apa={DETAIL.cariPenyakit} />

              {/* ⛔ Penolakan bukan akhir: klaim TIDAK terminal (ADR-U-0011).
                  Yang terminal adalah baris, dan baris berikutnya memulai
                  putaran baru dengan angka yang diperbaiki. */}
              {/* ⭐ Jalur akseptasi Claim Life sendiri - SaveAdjustment_Act.
                  Prasyaratnya dari XML: peserta DIPILIH, baris terakhir belum
                  bernomor, dan masih Outstanding. Perannya TIDAK diperiksa di
                  layar: pohon XML membuktikan tombolnya tidak bergerbang
                  peran, dan yang menggerbanginya pemegang tahap. */}
              {bolehSimpanAdjustment(p, p.baris) && (
                <p>
                  <button
                    type="button"
                    disabled={mengaksep === p.id}
                    onClick={() => void aksep(p.id)}
                  >
                    {mengaksep === p.id ? 'Menyimpan…' : 'Save Adjustment'}
                  </button>
                </p>
              )}

              {bolehPutaranBaru(p.baris) && (
                <p>
                  <button
                    type="button"
                    disabled={memutar === p.id}
                    onClick={() => void putaranBaru(p.id)}
                  >
                    {memutar === p.id ? 'Membuka…' : 'Putaran berikutnya'}
                  </button>
                </p>
              )}

              {p.baris.length === 0 ? (
                <p>Belum ada baris adjustment.</p>
              ) : (
                <table>
                  <thead>
                    <tr>
                      <th>Baris</th>
                      <th>Status</th>
                      <th>Jumlah klaim</th>
                      <th>Nomor akseptasi</th>
                      <th>Tindakan</th>
                      <th>Komite</th>
                    </tr>
                  </thead>
                  <tbody>
                    {p.baris.map((b) => (
                      <tr key={b.id}>
                        <td>{b.id}</td>
                        <td>
                          {/* Tiket 01 AC-3: status sebagai KATA — bukan nama field,
                              bukan angka. Kode mentahnya (b.kodeStatus) tetap ada di
                              kontrak API bagi yang menelusuri, tetapi tidak dicetak. */}
                          {b.status}
                        </td>
                        <td>{tampilUang(b.jumlahKlaim) || '—'}</td>
                        <td>{b.nomorAkseptasi || '—'}</td>
                        <td>
                          {/* Gerbang XML: `pyPosition=='ReasLifeAdmin' &&
                              CLAIM_NO != '' && .STS_REJECT == 0`. Dua syarat
                              terakhir dapat diperiksa di sini; PERANNYA tidak —
                              layar tidak tahu peran siapa pun, dan yang
                              menegakkannya services (403). Menyembunyikan
                              tombol di sini adalah kenyamanan, bukan pagar. */}
                          {b.status === STATUS_OUTSTANDING && klaim.nomorKlaim ? (
                            <button
                              type="button"
                              disabled={menolak === b.id}
                              onClick={() => void tolak(b.id)}
                            >
                              {menolak === b.id ? 'Menolak…' : 'Reject Outstanding'}
                            </button>
                          ) : (
                            '—'
                          )}
                        </td>
                        <td>
                          {/* Wewenangnya bergantung Type (QP/QR hanya SPV), dan
                              layar TIDAK tahu peran siapa pun. Yang diperiksa di
                              sini hanya keadaan baris; perannya ditegakkan
                              services (403). Menyembunyikan tombol adalah
                              kenyamanan, bukan pagar. */}
                          {bolehSerahkanDiLayar(b) ? (
                            <button
                              type="button"
                              disabled={menyerahkan === b.id}
                              onClick={() => void serahkan(p.id, b.id)}
                            >
                              {menyerahkan === b.id ? 'Menyerahkan…' : 'Send ke Komite'}
                            </button>
                          ) : b.komiteId ? (
                            <span title={b.komiteId}>sudah diserahkan</span>
                          ) : (
                            '—'
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </section>
          ))}
        </article>
      )}

      {klaim !== null && (
        <section className="os__aksi">
          {/* `CloseClaim_Section.xml` b1081. ⛔ Tombolnya BELUM MENUTUP:
              di Pega satu klik menjalankan `refresh` -> `ProtectCloseClaim_act`
              b1101 (yang memeriksa LALU memanggil FinishAssignment) dan
              `closeContainer` b1129. Yang ada di sini baru pemeriksaannya,
              dan labelnya mengatakannya - tombol yang menjanjikan lebih
              daripada yang ia lakukan adalah cacat yang paling mahal
              ditemukan belakangan. */}
          <button type="button" disabled={tutupSibuk} onClick={() => void periksaTutup()}>
            Periksa kesiapan: {DETAIL.tutupKlaim}
          </button>

          {gerbangTutup !== null && gerbangTutup.boleh && (
            <p role="status">
              Seluruh peserta sudah diaksep. {DETAIL.konfirmasiTutup}{' '}
              <strong>Penutupannya sendiri belum terpasang</strong> — tahap
              tujuan sesudah tutup belum dibaca dari alurnya.
            </p>
          )}
          {gerbangTutup !== null && !gerbangTutup.boleh && (
            <div role="alert">
              <p>Klaim belum dapat ditutup:</p>
              <ul>
                {gerbangTutup.penghalang.map((p) => (
                  <li key={`${p.urutan}:${p.nomorSertifikat}`}>{p.pesan}</li>
                ))}
              </ul>
            </div>
          )}
        </section>
      )}
    </section>
  )
}
